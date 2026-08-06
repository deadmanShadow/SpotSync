/**
 * Native Fetch API client wrapper.
 *
 * Responsibilities:
 *   1. Compose the full URL from PUBLIC_API_BASE_URL (loaded via Vite's
 *      import.meta.env inside Astro).
 *   2. Inject `Authorization: Bearer <token>` automatically when a token
 *      is present in localStorage (key: "spotsync_token").
 *   3. Parse the unified `ApiResponse<T>` envelope and unwrap its `data`.
 *   4. Translate non-2xx responses into typed `ApiError` instances so the
 *      UI can switch on HTTP status (401 -> redirect, 409 -> zone full,
 *      500 -> generic error toast).
 *
 * This file deliberately has zero third-party dependencies — we use the
 * native `fetch` API so the bundle stays minimal and the cogs of the
 * request lifecycle are explicit and easy to debug.
 */

import type { ApiResponse } from "../types/api";

const API_BASE_URL: string =
  import.meta.env.PUBLIC_API_BASE_URL ?? "http://localhost:8080/api/v1";

/** LocalStorage key holding the JWT Bearer token written by authStore. */
export const TOKEN_STORAGE_KEY = "spotsync_token";

/**
 * Custom error class carrying the HTTP status code and the backend's
 * human-readable message so callers can show useful toasts.
 */
export class ApiError extends Error {
  public readonly status: number;
  public readonly errors?: string;

  constructor(status: number, message: string, errors?: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.errors = errors;
  }
}

/**
 * Read the current JWT token from localStorage. Returns `null` when the
 * user is not authenticated. Kept as a function (not a constant) so each
 * request reads the latest value at call time, which matters when the
 * user logs in or out in another tab / via the auth store.
 */
function getAuthToken(): string | null {
  if (typeof window === "undefined") {
    // Defensive: this module is only called from client-side scripts,
    // but guard against SSR contexts just in case.
    return null;
  }
  const token = window.localStorage.getItem(TOKEN_STORAGE_KEY);
  return token && token.length > 0 ? token : null;
}

/**
 * Build the final headers for a request. JSON is the default content type
 * for any non-GET request that has a body. The Authorization header is
 * added opportunistically when a token is present.
 */
function buildHeaders(
  extraHeaders: Record<string, string> | undefined,
  hasBody: boolean,
): Record<string, string> {
  const headers: Record<string, string> = {
    Accept: "application/json",
  };

  if (hasBody) {
    headers["Content-Type"] = "application/json";
  }

  const token = getAuthToken();
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  if (extraHeaders) {
    for (const [key, value] of Object.entries(extraHeaders)) {
      headers[key] = value;
    }
  }

  return headers;
}

/**
 * Core wrapper around the native `fetch` API. Always returns the parsed
 * `data` payload on success; throws an `ApiError` on any non-2xx response.
 *
 * @typeParam T - The shape of the `data` field inside `ApiResponse`.
 * @param endpoint - Path appended to the API base URL (e.g. `/auth/login`).
 * @param options - Standard `RequestInit` overrides (method, body, headers…).
 */
export async function apiFetch<T>(
  endpoint: string,
  options: RequestInit = {},
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;

  const method = (options.method ?? "GET").toUpperCase();
  const hasBody = options.body !== undefined && options.body !== null;

  const headers = buildHeaders(
    options.headers as Record<string, string> | undefined,
    hasBody,
  );

  let response: Response;
  try {
    response = await fetch(url, {
      ...options,
      method,
      headers,
      // We never want to short-circuit on HTTP errors — we want to read
      // the body ourselves so we can surface the backend's message.
    });
  } catch (networkError) {
    // Truly unreachable server, DNS failure, CORS preflight rejected, etc.
    const message =
      networkError instanceof Error
        ? networkError.message
        : "Network error: unable to reach the SpotSync API.";
    throw new ApiError(0, message);
  }

  // Some endpoints (e.g. DELETE) may return an empty body. Guard accordingly.
  const contentType = response.headers.get("Content-Type") ?? "";
  const isJson = contentType.includes("application/json");

  let payload: ApiResponse<T> | null = null;
  if (isJson) {
    try {
      payload = (await response.json()) as ApiResponse<T>;
    } catch {
      payload = null;
    }
  }

  if (!response.ok) {
    const status = response.status;
    const message =
      payload?.message ??
      defaultMessageForStatus(status) ??
      `Request failed with status ${status}.`;
    const errors = payload?.errors;
    throw new ApiError(status, message, errors);
  }

  // Success path: unwrap the envelope. If the backend forgot to send a
  // `data` field we still return an empty cast so callers don't crash.
  if (payload && "data" in payload) {
    return payload.data;
  }

  // Fallback for endpoints that respond with raw JSON (no envelope).
  return (await response.json()) as T;
}

/**
 * Convenience: default human-readable message for common HTTP statuses
 * when the backend hasn't provided one. Keeps the UX friendly even if
 * the server returns an empty body.
 */
function defaultMessageForStatus(status: number): string | null {
  switch (status) {
    case 400:
      return "The request was invalid. Please check your input and try again.";
    case 401:
      return "Your session has expired. Please sign in again.";
    case 403:
      return "You do not have permission to perform this action.";
    case 404:
      return "The requested resource was not found.";
    case 409:
      return "This action conflicts with the current state of the resource.";
    case 422:
      return "Validation failed for the submitted data.";
    case 500:
      return "The server encountered an error. Please try again later.";
    case 503:
      return "The service is temporarily unavailable. Please try again later.";
    default:
      return null;
  }
}
