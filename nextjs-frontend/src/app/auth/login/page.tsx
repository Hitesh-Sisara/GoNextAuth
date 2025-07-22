// File: src/app/auth/login/page.tsx

"use client";

import { EmailStep } from "@/components/auth/email-step";
import { LoginChoice } from "@/components/auth/login-choice";
import { OTPStep } from "@/components/auth/otp-step";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { AuthService } from "@/lib/api/auth-service";
import { TokenManager } from "@/lib/auth/token-manager";
import { authConfig, brandConfig } from "@/lib/config/app-config";
import { useAuth, useGuestOnly } from "@/lib/hooks/use-auth";
import { Loader2 } from "lucide-react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { toast } from "sonner";

type LoginStep = "email" | "choice" | "otp";

export default function LoginPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { login } = useAuth();
  const { isReady } = useGuestOnly();

  const [currentStep, setCurrentStep] = useState<LoginStep>("email");
  const [email, setEmail] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [isResending, setIsResending] = useState(false);

  const redirectPath =
    searchParams.get("redirect") || authConfig.defaultRedirect;

  useEffect(() => {
    // Show success message if coming from signup
    const message = searchParams.get("message");
    if (message === "signup-success") {
      toast.success("Account created successfully! Please sign in.");
    } else if (message === "verification-success") {
      toast.success("Email verified successfully! Please sign in.");
    } else if (message === "password-reset-success") {
      toast.success("Password reset successfully! Please sign in.");
    }
  }, [searchParams]);

  const handleEmailSubmit = async (emailValue: string) => {
    setEmail(emailValue);
    setCurrentStep("choice");
  };

  const handlePasswordLogin = async (password: string) => {
    try {
      setIsLoading(true);
      const success = await login({ email, password });
      if (success) {
        router.push(redirectPath);
      }
    } catch (error: unknown) {
      console.error("Login error:", error);
      toast.error(
        error instanceof Error
          ? error.message
          : "Login failed. Please try again."
      );
    } finally {
      setIsLoading(false);
    }
  };

  const handleOTPLogin = async () => {
    try {
      setIsLoading(true);
      const response = await AuthService.initiateEmailLogin({ email });
      if (response.success) {
        setCurrentStep("otp");
        toast.success("Verification code sent to your email");
      } else {
        throw new Error(response.message || "Failed to send OTP");
      }
    } catch (error: unknown) {
      console.error("OTP login error:", error);
      toast.error(
        error instanceof Error
          ? error.message
          : "Failed to send verification code"
      );
    } finally {
      setIsLoading(false);
    }
  };

  // Handle forgot password with current email
  const handleForgotPassword = () => {
    if (email) {
      // Pass the current email to forgot password page
      router.push(`/auth/forgot-password?email=${encodeURIComponent(email)}`);
    } else {
      router.push("/auth/forgot-password");
    }
  };

  const handleOTPSubmit = async (otp: string) => {
    try {
      setIsLoading(true);
      const response = await AuthService.completeOTPLogin({ email, otp });

      if (response.success && response.data) {
        const { user, access_token, refresh_token, expires_in } = response.data;

        // Store tokens
        TokenManager.setTokens(access_token, refresh_token, expires_in);
        TokenManager.setUser(user);

        // Update auth store
        const { useAuthStore } = await import("@/lib/store/auth-store");
        const { setUser } = useAuthStore.getState();
        setUser(user);

        toast.success("Login successful!");
        router.push(redirectPath);
      } else {
        throw new Error(response.message || "OTP verification failed");
      }
    } catch (error: unknown) {
      console.error("OTP verification error:", error);
      toast.error(
        error instanceof Error ? error.message : "Invalid verification code"
      );
    } finally {
      setIsLoading(false);
    }
  };

  const handleResendOTP = async () => {
    try {
      setIsResending(true);
      const response = await AuthService.resendOTP({
        email,
        otp_type: "login",
      });

      if (response.success) {
        toast.success("New verification code sent to your email");
      } else {
        throw new Error(response.message || "Failed to resend code");
      }
    } catch (error: unknown) {
      console.error("Resend OTP error:", error);
      toast.error(
        error instanceof Error
          ? error.message
          : "Failed to resend verification code"
      );
    } finally {
      setIsResending(false);
    }
  };

  const handleGoogleSuccess = async (authResponse: {
    user: {
      id: string;
      email: string;
      name?: string;
      avatar_url?: string;
      created_at: string;
      updated_at: string;
    };
    access_token: string;
    refresh_token: string;
    expires_in: number;
  }) => {
    try {
      const { user, access_token, refresh_token, expires_in } = authResponse;

      // Store tokens
      TokenManager.setTokens(access_token, refresh_token, expires_in);
      TokenManager.setUser(user);

      // Update auth store
      const { useAuthStore } = await import("@/lib/store/auth-store");
      const { setUser } = useAuthStore.getState();
      setUser({
        ...user,
        first_name: user.name?.split(" ")[0] || "",
        last_name: user.name?.split(" ")[1] || "",
        is_email_verified: true, // Google OAuth users are verified
        is_active: true,
        role: "user",
        provider: "google",
      });

      toast.success("Google Sign-In successful!");
      router.push(redirectPath);
    } catch (error: unknown) {
      console.error("Google auth error:", error);
      toast.error("Google authentication failed");
    }
  };

  const goBack = () => {
    if (currentStep === "choice") {
      setCurrentStep("email");
    } else if (currentStep === "otp") {
      setCurrentStep("choice");
    }
  };

  if (!isReady) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  const getTitle = () => {
    switch (currentStep) {
      case "email":
        return `Sign in to ${brandConfig.name}`;
      case "choice":
        return `Sign in to ${brandConfig.name}`;
      case "otp":
        return "Enter sign-in code";
      default:
        return `Sign in to ${brandConfig.name}`;
    }
  };

  const getDescription = () => {
    switch (currentStep) {
      case "email":
        return "Welcome back! Please sign in to continue";
      case "choice":
        return "";
      case "otp":
        return "Check your email for the verification code";
      default:
        return "Welcome back! Please sign in to continue";
    }
  };

  const renderStep = () => {
    switch (currentStep) {
      case "email":
        return (
          <EmailStep
            onSubmit={handleEmailSubmit}
            onGoogleSuccess={handleGoogleSuccess}
            isLoading={isLoading}
            buttonText="Continue"
            description="Enter your email address to sign in"
            showGoogleOption={true}
          />
        );

      case "choice":
        return (
          <LoginChoice
            email={email}
            onPasswordLogin={handlePasswordLogin}
            onOTPLogin={handleOTPLogin}
            onForgotPassword={handleForgotPassword}
            isLoading={isLoading}
            title={`Sign in to ${brandConfig.name}`}
          />
        );

      case "otp":
        return (
          <OTPStep
            email={email}
            onSubmit={handleOTPSubmit}
            onResend={handleResendOTP}
            isLoading={isLoading}
            isResending={isResending}
            title="Enter sign-in code"
            description="We've sent a sign-in code to"
          />
        );

      default:
        return null;
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1 text-center">
          <CardTitle className="text-2xl font-bold">{getTitle()}</CardTitle>
          {getDescription() && (
            <CardDescription>{getDescription()}</CardDescription>
          )}
        </CardHeader>

        <CardContent>
          {renderStep()}

          {currentStep !== "email" && (
            <div className="mt-6 text-center">
              <button
                onClick={goBack}
                className="text-sm text-gray-600 hover:text-gray-500"
                disabled={isLoading}
              >
                ← Go back
              </button>
            </div>
          )}

          <div className="mt-6 text-center text-sm">
            Don&apos;t have an account?{" "}
            <Link
              href="/auth/signup"
              className="text-blue-600 hover:text-blue-500 font-medium"
            >
              Sign up
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
