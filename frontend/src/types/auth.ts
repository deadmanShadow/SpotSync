/**
 * Domain type contracts for authentication.
 *
 * These interfaces mirror the Go backend DTOs so that the fetch layer
 * (src/services/api.ts) can pass strongly-typed payloads from end to end.
 */

/** Role discriminator for the two first-class user personas in SpotSync. */
export type UserRole = "driver" | "admin";

/** Authenticated user record returned by /auth/login and /auth/register. */
export interface User {
  id: number;
  full_name: string;
  email: string;
  role: UserRole;
  created_at: string;
}

/** Credentials captured on the login page. */
export interface LoginCredentials {
  email: string;
  password: string;
}

/** Credentials captured on the registration page (role is selected client-side). */
export interface RegisterPayload {
  full_name: string;
  email: string;
  password: string;
  role: UserRole;
}

/**
 * Auth response envelope returned by both /auth/login and /auth/register.
 * The token is a JWT Bearer token that the apiFetch wrapper attaches to
 * every subsequent request.
 */
export interface AuthResponse {
  token: string;
  user: User;
}
