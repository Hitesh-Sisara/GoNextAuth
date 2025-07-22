// src/components/auth/auth-guard.tsx

"use client";

import { useAuth } from "@/lib/hooks/use-auth";
import { Loader2 } from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

interface AuthGuardProps {
  children: React.ReactNode;
  requireAuth?: boolean;
  redirectTo?: string;
}

export function AuthGuard({
  children,
  requireAuth = true,
  redirectTo,
}: AuthGuardProps) {
  const { isAuthenticated, isLoading, user } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading) {
      if (requireAuth && !isAuthenticated) {
        // Redirect to login if auth is required but user is not authenticated
        router.push(redirectTo || "/auth/login");
      } else if (!requireAuth && isAuthenticated) {
        // Redirect to dashboard if user is authenticated but accessing guest-only pages
        router.push(redirectTo || "/dashboard");
      }
    }
  }, [isAuthenticated, isLoading, requireAuth, router, redirectTo]);

  // Show loading while checking authentication
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <Loader2 className="h-8 w-8 animate-spin mx-auto mb-4 text-blue-600" />
          <p className="text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  // For protected routes: only render if authenticated
  if (requireAuth && !isAuthenticated) {
    return null;
  }

  // For guest-only routes: only render if not authenticated
  if (!requireAuth && isAuthenticated) {
    return null;
  }

  return <>{children}</>;
}

// Convenience wrapper components
export function ProtectedRoute({
  children,
  redirectTo,
}: {
  children: React.ReactNode;
  redirectTo?: string;
}) {
  return (
    <AuthGuard requireAuth={true} redirectTo={redirectTo}>
      {children}
    </AuthGuard>
  );
}

export function GuestOnlyRoute({
  children,
  redirectTo,
}: {
  children: React.ReactNode;
  redirectTo?: string;
}) {
  return (
    <AuthGuard requireAuth={false} redirectTo={redirectTo}>
      {children}
    </AuthGuard>
  );
}
