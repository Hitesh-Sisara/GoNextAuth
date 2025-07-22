// File: src/components/auth/google-signin-button.tsx

"use client";

import { Button } from "@/components/ui/button";
import { AuthService } from "@/lib/api/auth-service";
import { useGoogleAuth } from "@/lib/auth/google-auth";
import { Loader2 } from "lucide-react";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { toast } from "sonner";

interface GoogleSignInButtonProps {
  onSuccess?: (authResponse: any) => void;
  onError?: (error: Error) => void;
  disabled?: boolean;
  className?: string;
  variant?: "default" | "outline" | "ghost";
  size?: "default" | "sm" | "lg";
  text?: string;
  redirectPath?: string;
  useServerFlow?: boolean;
}

export function GoogleSignInButton({
  onSuccess,
  onError,
  disabled = false,
  className = "",
  variant = "outline",
  size = "default",
  text = "Continue with Google",
  redirectPath,
  useServerFlow = true, // Default to server flow for better security
}: GoogleSignInButtonProps) {
  const { isLoaded, signIn } = useGoogleAuth();
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  const handleServerFlow = async () => {
    try {
      setIsLoading(true);

      // Store redirect path for callback
      if (redirectPath) {
        localStorage.setItem("google_auth_redirect", redirectPath);
      }

      // Get Google OAuth URL from backend
      const response = await AuthService.getGoogleAuthURL();

      if (response.success && response.data?.auth_url) {
        // Redirect to Google OAuth
        window.location.href = response.data.auth_url;
      } else {
        throw new Error(response.message || "Failed to get Google OAuth URL");
      }
    } catch (error: any) {
      console.error("Google OAuth URL error:", error);
      toast.error(error.message || "Failed to initiate Google Sign-In");
      setIsLoading(false);

      if (onError) {
        onError(error);
      }
    }
  };

  const handleClientFlow = async () => {
    if (!isLoaded) {
      toast.error("Google Sign-In is not ready yet");
      return;
    }

    try {
      setIsLoading(true);

      // Get access token from Google
      const accessToken = await signIn();

      // Authenticate with our backend
      const response = await AuthService.googleAuth({
        access_token: accessToken,
      });

      if (response.success && response.data) {
        if (onSuccess) {
          onSuccess(response.data);
        }
        toast.success("Google Sign-In successful!");

        if (redirectPath) {
          router.push(redirectPath);
        }
      } else {
        throw new Error(response.message || "Google authentication failed");
      }
    } catch (error: any) {
      console.error("Google Sign-In error:", error);
      const errorMessage = error.message || "Google Sign-In failed";
      toast.error(errorMessage);

      if (onError) {
        onError(error);
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleGoogleSignIn = () => {
    if (useServerFlow) {
      handleServerFlow();
    } else {
      handleClientFlow();
    }
  };

  return (
    <Button
      type="button"
      variant={variant}
      size={size}
      className={`w-full ${className}`}
      onClick={handleGoogleSignIn}
      disabled={disabled || isLoading || (!useServerFlow && !isLoaded)}
    >
      {isLoading ? (
        <>
          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
          Signing in...
        </>
      ) : (
        <>
          <svg className="mr-2 h-4 w-4" viewBox="0 0 24 24">
            <path
              fill="currentColor"
              d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
            />
            <path
              fill="currentColor"
              d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
            />
            <path
              fill="currentColor"
              d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
            />
            <path
              fill="currentColor"
              d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
            />
          </svg>
          {text}
        </>
      )}
    </Button>
  );
}
