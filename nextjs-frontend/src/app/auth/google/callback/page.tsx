"use client";

import { AuthService } from "@/lib/api/auth-service";
import { TokenManager } from "@/lib/auth/token-manager";
import { Loader2 } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";

export default function GoogleCallbackPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [isProcessing, setIsProcessing] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Use ref to prevent duplicate processing
  const hasProcessed = useRef(false);
  const processingRef = useRef(false);

  useEffect(() => {
    const handleGoogleCallback = async () => {
      // Prevent duplicate processing
      if (hasProcessed.current || processingRef.current) {
        console.log("Callback already processed or in progress");
        return;
      }

      processingRef.current = true;

      try {
        // Get parameters from URL
        const code = searchParams.get("code");
        const error = searchParams.get("error");
        const state = searchParams.get("state");

        console.log("Callback params:", {
          code: !!code,
          error,
          state: !!state,
        });

        // Handle OAuth errors
        if (error) {
          throw new Error(
            error === "access_denied"
              ? "Google sign-in was cancelled"
              : `Google OAuth error: ${error}`
          );
        }

        // Check if authorization code is present
        if (!code) {
          throw new Error("Authorization code not found in callback URL");
        }

        // Check if state is present
        if (!state) {
          throw new Error("State parameter not found in callback URL");
        }

        // Mark as processed before making the API call
        hasProcessed.current = true;

        // Exchange code for tokens via backend - include state parameter
        console.log("Sending code and state to backend...");
        const response = await AuthService.googleCallback(code, state);

        if (response.success && response.data) {
          const { user, access_token, refresh_token, expires_in } =
            response.data;

          // Store tokens
          TokenManager.setTokens(access_token, refresh_token, expires_in);
          TokenManager.setUser(user);

          // Update auth store
          const { useAuthStore } = await import("@/lib/store/auth-store");
          const { setUser } = useAuthStore.getState();
          setUser(user);

          // Show success message
          toast.success("Google Sign-In successful!");

          // Redirect to dashboard or intended page
          const redirectPath =
            localStorage.getItem("google_auth_redirect") || "/dashboard";
          localStorage.removeItem("google_auth_redirect"); // Clean up

          console.log("Redirecting to:", redirectPath);

          // Use replace instead of push to prevent back button issues
          router.replace(redirectPath);
        } else {
          throw new Error(response.message || "Google authentication failed");
        }
      } catch (error: any) {
        console.error("Google callback error:", error);
        setError(error.message || "Google authentication failed");

        // Show user-friendly error message
        const userMessage = error.message?.includes("invalid_grant")
          ? "The Google sign-in session has expired. Please try signing in again."
          : error.message || "Google authentication failed";

        toast.error(userMessage);

        // Redirect to login page after a delay
        setTimeout(() => {
          router.replace("/auth/login?error=google_auth_failed");
        }, 3000);
      } finally {
        setIsProcessing(false);
        processingRef.current = false;
      }
    };

    // Only run if we have search params and haven't processed yet
    if (searchParams.get("code") && !hasProcessed.current) {
      handleGoogleCallback();
    } else if (!searchParams.get("code") && !searchParams.get("error")) {
      // No valid params, redirect to login
      setError("Invalid callback parameters");
      setTimeout(() => {
        router.replace("/auth/login");
      }, 2000);
    }
  }, [searchParams, router]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full bg-white p-6 rounded-lg shadow-md text-center">
        {isProcessing ? (
          <>
            <Loader2 className="h-8 w-8 animate-spin mx-auto mb-4 text-blue-600" />
            <h2 className="text-xl font-semibold text-gray-900 mb-2">
              Completing Google Sign-In
            </h2>
            <p className="text-gray-600">
              Please wait while we authenticate your account...
            </p>
          </>
        ) : error ? (
          <>
            <div className="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg
                className="w-6 h-6 text-red-600"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </div>
            <h2 className="text-xl font-semibold text-gray-900 mb-2">
              Authentication Failed
            </h2>
            <p className="text-gray-600 mb-4">{error}</p>
            <p className="text-sm text-gray-500">
              Redirecting to login page...
            </p>
          </>
        ) : (
          <>
            <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg
                className="w-6 h-6 text-green-600"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M5 13l4 4L19 7"
                />
              </svg>
            </div>
            <h2 className="text-xl font-semibold text-gray-900 mb-2">
              Sign-In Successful!
            </h2>
            <p className="text-gray-600">Redirecting to your dashboard...</p>
          </>
        )}
      </div>
    </div>
  );
}
