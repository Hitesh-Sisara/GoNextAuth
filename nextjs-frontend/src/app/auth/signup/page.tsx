// File: src/app/auth/signup/page.tsx

"use client";

import { EmailStep } from "@/components/auth/email-step";
import { OTPStep } from "@/components/auth/otp-step";
import { PasswordStep } from "@/components/auth/password-step";
import { ProfileStep } from "@/components/auth/profile-step";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { AuthService } from "@/lib/api/auth-service";
import { useGuestOnly } from "@/lib/hooks/use-auth";
import { Loader2 } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { toast } from "sonner";

type SignupStep = "email" | "otp" | "password" | "profile";

interface SignupData {
  email: string;
  password?: string;
  firstName?: string;
  lastName?: string;
  phone?: string;
}

export default function SignupPage() {
  const router = useRouter();
  const { isReady } = useGuestOnly();

  const [currentStep, setCurrentStep] = useState<SignupStep>("email");
  const [signupData, setSignupData] = useState<SignupData>({ email: "" });
  const [isLoading, setIsLoading] = useState(false);
  const [isResending, setIsResending] = useState(false);

  const handleEmailSubmit = async (email: string) => {
    try {
      setIsLoading(true);
      const response = await AuthService.initiateSignup({ email });

      if (response.success) {
        setSignupData({ ...signupData, email });
        setCurrentStep("otp");
        toast.success("Verification code sent to your email");
      } else {
        throw new Error(response.message || "Failed to send verification code");
      }
    } catch (error: any) {
      console.error("Email submit error:", error);
      if (error.message?.includes("already exists")) {
        toast.error("An account with this email already exists");
      } else {
        toast.error(error.message || "Failed to send verification code");
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleOTPSubmit = async (otp: string) => {
    try {
      setIsLoading(true);
      const response = await AuthService.verifyEmail({
        email: signupData.email,
        otp,
      });

      if (response.success) {
        setCurrentStep("password");
        toast.success("Email verified successfully!");
      } else {
        throw new Error(response.message || "Email verification failed");
      }
    } catch (error: any) {
      console.error("OTP verification error:", error);
      toast.error(error.message || "Invalid verification code");
    } finally {
      setIsLoading(false);
    }
  };

  const handlePasswordSubmit = async (password: string) => {
    setSignupData({ ...signupData, password });
    setCurrentStep("profile");
  };

  const handleProfileSubmit = async (profileData: {
    firstName: string;
    lastName: string;
    phone?: string;
  }) => {
    try {
      setIsLoading(true);

      const completeSignupData = {
        email: signupData.email,
        password: signupData.password!,
        first_name: profileData.firstName,
        last_name: profileData.lastName,
        phone: profileData.phone || "",
      };

      const response = await AuthService.completeSignup(completeSignupData);

      if (response.success) {
        toast.success("Account created successfully!");
        router.push("/auth/login?message=signup-success");
      } else {
        throw new Error(response.message || "Account creation failed");
      }
    } catch (error: any) {
      console.error("Profile submit error:", error);
      toast.error(error.message || "Failed to create account");
    } finally {
      setIsLoading(false);
    }
  };

  const handleResendOTP = async () => {
    try {
      setIsResending(true);
      const response = await AuthService.resendOTP({
        email: signupData.email,
        otp_type: "email_verification",
      });

      if (response.success) {
        toast.success("New verification code sent to your email");
      } else {
        throw new Error(response.message || "Failed to resend code");
      }
    } catch (error: any) {
      console.error("Resend OTP error:", error);
      toast.error(error.message || "Failed to resend verification code");
    } finally {
      setIsResending(false);
    }
  };

  const handleGoogleSignIn = async () => {
    try {
      setIsLoading(true);
      // Implement Google OAuth flow
      toast.info("Google Sign-In will be implemented soon");
    } catch (error: any) {
      console.error("Google sign-in error:", error);
      toast.error("Google sign-in failed");
    } finally {
      setIsLoading(false);
    }
  };

  const goBack = () => {
    if (currentStep === "otp") {
      setCurrentStep("email");
    } else if (currentStep === "password") {
      setCurrentStep("otp");
    } else if (currentStep === "profile") {
      setCurrentStep("password");
    }
  };

  if (!isReady) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  const getStepTitle = () => {
    switch (currentStep) {
      case "email":
        return "Create your account";
      case "otp":
        return "Verify your email";
      case "password":
        return "Secure your account";
      case "profile":
        return "Complete your profile";
      default:
        return "Create your account";
    }
  };

  const getStepDescription = () => {
    switch (currentStep) {
      case "email":
        return "Enter your email to get started";
      case "otp":
        return "Check your email for the verification code";
      case "password":
        return "Choose a strong password";
      case "profile":
        return "Just a few more details to get you started";
      default:
        return "Enter your email to get started";
    }
  };

  const renderStep = () => {
    switch (currentStep) {
      case "email":
        return (
          <EmailStep
            onSubmit={handleEmailSubmit}
            isLoading={isLoading}
            buttonText="Get started"
            description="Enter your email address to create your account"
            onGoogleSignIn={handleGoogleSignIn}
          />
        );

      case "otp":
        return (
          <OTPStep
            email={signupData.email}
            onSubmit={handleOTPSubmit}
            onResend={handleResendOTP}
            isLoading={isLoading}
            isResending={isResending}
            title="Verify your email"
            description="We've sent a verification code to"
          />
        );

      case "password":
        return (
          <PasswordStep
            onSubmit={handlePasswordSubmit}
            isLoading={isLoading}
            title="Create your password"
            buttonText="Continue"
          />
        );

      case "profile":
        return (
          <ProfileStep
            onSubmit={handleProfileSubmit}
            isLoading={isLoading}
            title="Tell us about yourself"
            buttonText="Create Account"
          />
        );

      default:
        return null;
    }
  };

  const getProgressPercentage = () => {
    switch (currentStep) {
      case "email":
        return 25;
      case "otp":
        return 50;
      case "password":
        return 75;
      case "profile":
        return 100;
      default:
        return 0;
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1 text-center">
          <CardTitle className="text-2xl font-bold">{getStepTitle()}</CardTitle>
          <CardDescription>{getStepDescription()}</CardDescription>

          {/* Progress bar */}
          <div className="w-full bg-gray-200 rounded-full h-2 mt-4">
            <div
              className="bg-blue-600 h-2 rounded-full transition-all duration-300"
              style={{ width: `${getProgressPercentage()}%` }}
            />
          </div>
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
            Already have an account?{" "}
            <Link
              href="/auth/login"
              className="text-blue-600 hover:text-blue-500 font-medium"
            >
              Sign in
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
