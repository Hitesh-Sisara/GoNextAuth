// File: src/components/auth/email-step.tsx

"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { AuthService } from "@/lib/api/auth-service";
import { validateEmail } from "@/lib/utils";
import { Loader2 } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

interface EmailStepProps {
  onSubmit: (email: string) => Promise<void> | void;
  onGoogleSuccess?: (authResponse: any) => Promise<void> | void;
  onGoogleSignIn?: () => Promise<void> | void;
  isLoading: boolean;
  buttonText?: string;
  description?: string;
  showGoogleOption?: boolean;
}

export function EmailStep({
  onSubmit,
  onGoogleSuccess,
  onGoogleSignIn,
  isLoading,
  buttonText = "Continue",
  description = "Enter your email address to continue",
  showGoogleOption = true,
}: EmailStepProps) {
  const [email, setEmail] = useState("");
  const [emailError, setEmailError] = useState("");
  const [isGoogleLoading, setIsGoogleLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!email.trim()) {
      setEmailError("Email is required");
      return;
    }

    if (!validateEmail(email)) {
      setEmailError("Please enter a valid email address");
      return;
    }

    setEmailError("");
    await onSubmit(email);
  };

  const handleGoogleAuth = async () => {
    try {
      setIsGoogleLoading(true);

      if (onGoogleSignIn) {
        await onGoogleSignIn();
        return;
      }

      // Get Google auth URL from backend
      const response = await AuthService.getGoogleAuthURL();

      if (response.success && response.data?.auth_url) {
        // Store current page for redirect after auth
        localStorage.setItem("google_auth_redirect", window.location.pathname);
        // Redirect to Google OAuth
        window.location.href = response.data.auth_url;
      } else {
        throw new Error("Failed to get Google auth URL");
      }
    } catch (error: any) {
      console.error("Google auth error:", error);
      toast.error("Google sign-in failed. Please try again.");
    } finally {
      setIsGoogleLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      {description && (
        <p className="text-sm text-gray-600 text-center">{description}</p>
      )}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="email">Email address</Label>
          <Input
            id="email"
            type="email"
            placeholder="Enter your email"
            value={email}
            onChange={(e) => {
              setEmail(e.target.value);
              if (emailError) setEmailError("");
            }}
            disabled={isLoading || isGoogleLoading}
            required
          />
          {emailError && <p className="text-sm text-red-600">{emailError}</p>}
        </div>

        <Button
          type="submit"
          className="w-full"
          disabled={isLoading || isGoogleLoading || !email.trim()}
        >
          {isLoading ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Loading...
            </>
          ) : (
            buttonText
          )}
        </Button>
      </form>

      {showGoogleOption && (
        <>
          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <Separator className="w-full" />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-white px-2 text-gray-500">OR</span>
            </div>
          </div>

          <Button
            type="button"
            variant="outline"
            className="w-full"
            onClick={handleGoogleAuth}
            disabled={isLoading || isGoogleLoading}
          >
            {isGoogleLoading ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Connecting to Google...
              </>
            ) : (
              <>
                <svg className="mr-2 h-4 w-4" viewBox="0 0 24 24">
                  <path
                    fill="#4285F4"
                    d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                  />
                  <path
                    fill="#34A853"
                    d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                  />
                  <path
                    fill="#FBBC05"
                    d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                  />
                  <path
                    fill="#EA4335"
                    d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                  />
                </svg>
                Continue with Google
              </>
            )}
          </Button>
        </>
      )}
    </div>
  );
}
