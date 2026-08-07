/**
 * Admin user-management API endpoints.
 *
 * Thin wrappers around `apiFetch` for the admin-only user-listing routes:
 *   - GET    /auth/users           -> list every registered user
 *   - GET    /auth/users/count     -> count of users with a given role
 *
 * All HTTP plumbing (URL composition, JSON encoding, Bearer token injection,
 * error translation) lives in `api.ts`. This file is purely contract.
 */

import { apiFetch } from "./api";
import type { User, UserRole } from "../types/auth";

/**
 * Fetch every registered user in the system, newest-first.
 *
 * Admin-only endpoint — `AdminOnly` middleware guards it on the backend.
 * Returns the full `User[]` already unwrapped from the API envelope.
 */
export async function getAllUsers(): Promise<User[]> {
  return apiFetch<User[]>("/auth/users", {
    method: "GET",
  });
}

/**
 * Count the number of users with the given role.
 *
 * Used by the admin dashboard to summarize drivers vs admins without
 * downloading the entire user roster.
 */
export async function countUsersByRole(role: UserRole): Promise<number> {
  const response = await apiFetch<{ count: number }>(
    `/auth/users/count?role=${encodeURIComponent(role)}`,
    {
      method: "GET",
    },
  );
  return typeof response?.count === "number" ? response.count : 0;
}

/**
 * Admin-only: delete a user by ID. Used by the admin dashboard to
 * permanently remove driver or admin accounts. The backend returns 404
 * if the user does not exist (translated to ApiError by api.ts).
 */
export async function deleteUser(id: number): Promise<void> {
  await apiFetch<void>(`/auth/users/${id}`, {
    method: "DELETE",
  });
}