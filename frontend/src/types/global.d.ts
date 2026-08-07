/**
 * Global ambient type declarations for SpotSync.
 *
 * - `window.showToast` is exposed by the layout's inline script so any
 *   React component can call into it without going through imports.
 *   This file declares the global signature so TypeScript knows the
 *   function exists on `window`.
 */

export {};

declare global {
  interface Window {
    showToast?: (
      message: string,
      type?: "success" | "error" | "info",
      duration?: number,
    ) => void;
  }
}
