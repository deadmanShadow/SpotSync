import * as React from "react";
import { ShieldOff, ArrowLeft } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { $user, isAdmin } from "@/store/authStore";

interface AdminGuardProps {
  children: React.ReactNode;
}

/**
 * React admin guard.
 *
 * Renders a "checking" state until the persistent auth store is hydrated,
 * then either renders the protected children (admin) or the forbidden
 * state (non-admin / signed out). The forbidden state replaces the old
 * "super-admin" wording with a friendlier "admin role required" message.
 */
export function AdminGuard({ children }: AdminGuardProps) {
  const [status, setStatus] = React.useState<"checking" | "ready" | "forbidden">(
    "checking",
  );

  React.useEffect(() => {
    const user = $user.get();
    if (!user) {
      // Not signed in — bounce to /login.
      window.location.replace("/login");
      return;
    }
    if (!isAdmin()) {
      setStatus("forbidden");
      return;
    }
    setStatus("ready");
  }, []);

  if (status === "checking") {
    return (
      <div className="flex min-h-[60vh] items-center justify-center">
        <div className="flex flex-col items-center gap-3 text-slate-400">
          <span className="inline-block h-9 w-9 animate-spin rounded-full border-2 border-slate-700 border-t-emerald-500" />
          <p className="text-sm">Verifying admin access…</p>
        </div>
      </div>
    );
  }

  if (status === "forbidden") {
    return (
      <div className="flex min-h-[60vh] items-center justify-center px-4 py-10">
        <Card className="max-w-md w-full">
          <CardContent className="flex flex-col items-center text-center p-8 gap-2">
            <span className="inline-flex h-14 w-14 items-center justify-center rounded-2xl bg-red-500/15 text-red-300 border border-red-500/30 mb-2">
              <ShieldOff className="h-7 w-7" />
            </span>
            <h2 className="font-display text-2xl font-semibold text-slate-50">
              Administrator access required
            </h2>
            <p className="text-sm text-slate-400 max-w-sm">
              Your account doesn't have the <strong>admin</strong> role. Make
              sure you signed in with the administrator account — there is no
              separate super-admin tier in SpotSync.
            </p>
            <Button
              type="button"
              variant="outline"
              className="mt-4"
              onClick={() => window.location.replace("/")}
            >
              <ArrowLeft className="h-4 w-4" />
              Return Home
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return <>{children}</>;
}
