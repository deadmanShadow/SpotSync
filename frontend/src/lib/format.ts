/**
 * Formatting helpers shared across admin pages.
 *
 * Keeping these in one place means the reservations table, drivers page,
 * users page, and overview dashboard all render dates and numbers the
 * same way — no copy-pasted formatters, no inconsistent timestamps.
 */

const MONTHS = [
  "Jan",
  "Feb",
  "Mar",
  "Apr",
  "May",
  "Jun",
  "Jul",
  "Aug",
  "Sep",
  "Oct",
  "Nov",
  "Dec",
] as const;

/**
 * Format an ISO timestamp as `8-Aug-2026 15:38` (locale-independent).
 * Returns the raw input as a graceful fallback when the value is not a
 * valid Date so the UI never shows "Invalid Date".
 */
export function formatDateTime(iso: string | null | undefined): string {
  if (!iso) return "—";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  const day = d.getDate();
  const month = MONTHS[d.getMonth()];
  const year = d.getFullYear();
  const hh = String(d.getHours()).padStart(2, "0");
  const mm = String(d.getMinutes()).padStart(2, "0");
  return `${day}-${month}-${year} ${hh}:${mm}`;
}

/**
 * Format an ISO timestamp as `8-Aug-2026` (date only, no time).
 * Useful for registration dates and other date-only contexts.
 */
export function formatDate(iso: string | null | undefined): string {
  if (!iso) return "—";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  const day = d.getDate();
  const month = MONTHS[d.getMonth()];
  const year = d.getFullYear();
  return `${day}-${month}-${year}`;
}

/**
 * Escape an arbitrary string for safe injection into HTML.
 * Used by the small set of append-via-innerHTML helpers across the admin
 * pages, where we want to avoid both XSS and layout corruption.
 */
export function escapeHtml(input: string | null | undefined): string {
  return String(input ?? "")
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;");
}

/**
 * Two-letter initials for an avatar fallback. Falls back to `?` when the
 * source string is empty.
 */
export function initialsOf(name: string, email: string): string {
  const source = (name || email || "").trim();
  if (!source) return "?";
  return source
    .split(/\s+|@|\./)
    .filter((part) => part.length > 0)
    .map((part) => part.charAt(0).toUpperCase())
    .slice(0, 2)
    .join("");
}

/**
 * Compact, human-readable "5 minutes ago" relative timestamp.
 * Returns the absolute `formatDateTime` string when the input is older
 * than 7 days so the relative terms don't become misleading.
 */
export function formatRelativeTime(iso: string | null | undefined): string {
  if (!iso) return "—";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  const diff = Date.now() - d.getTime();
  const sec = Math.floor(diff / 1000);
  if (sec < 60) return "just now";
  const min = Math.floor(sec / 60);
  if (min < 60) return `${min} min ago`;
  const hr = Math.floor(min / 60);
  if (hr < 24) return `${hr} hr ago`;
  const day = Math.floor(hr / 24);
  if (day < 7) return `${day} day${day === 1 ? "" : "s"} ago`;
  return formatDateTime(iso);
}
