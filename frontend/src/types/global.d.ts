/**
 * Global ambient type declarations for SpotSync.
 *
 * - `window.showToast` is exposed by the layout's inline script so any
 *   React component or Astro page can call into it without going through
 *   imports. The recommended way to fire toasts in app code is to
 *   import the `toast` helper from `@/lib/toast` instead — this
 *   declaration is kept for the inline script that wires the DOM.
 *   This file declares the global signature so TypeScript knows the
 *   function exists on `window`.
 */

import type { ToastType } from "../lib/toast";

export {};

declare global {
  interface Window {
    showToast?: (
      message: string,
      type?: ToastType,
      duration?: number,
    ) => void;
  }
}
