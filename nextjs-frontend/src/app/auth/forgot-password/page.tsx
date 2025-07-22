// File: src/app/auth/forgot-password/page.tsx

"use client";

import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
} from "@/components/ui/input-otp";
import { Label } from "@/components/ui/label";
import { AuthService } from "@/lib/api/auth-service";
import { useGuestOnly } from "@/lib/hooks/use-auth";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  ArrowLeft,
  CheckCircle,
  Eye,
  EyeOff,
  KeyRound,
  Loader2,
  Mail,
} from "lucide-react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

type ForgotPasswordStep = "email" | "otp" | "password" | "success";

const emailSchema = z.object({
  email: z.string().email("Please enter a valid email address"),
});

const otpSchema = z.object({
  otp: z.string().length(6, "OTP must be 6 digits"),
});

const passwordSchema = z
  .object({
    new_password: z
      .string()
      .min(8, "Password must be at least 8 characters")
      .regex(/[A-Z]/, "Password must contain at least one uppercase letter")
      .regex(/[a-z]/, "Password must contain at least one lowercase letter")
      .regex(/[0-9]/, "Password must contain at least one number")
      .regex(
        /[^A-Za-z0-9]/,
        "Password must contain at least one special character"
      ),
    confirm_password: z.string(),
  })
  .refine((data) => data.new_password === data.confirm_password, {
    message: "Passwords don't match",
    path: ["confirm_password"],
  });

type EmailFormData = z.infer<typeof emailSchema>;
type OTPFormData = z.infer<typeof otpSchema>;
type PasswordFormData = z.infer<typeof passwordSchema>;

export default function ForgotPasswordPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { isReady } = useGuestOnly();

  const [currentStep, setCurrentStep] = useState<ForgotPasswordStep>("email");
  const [email, setEmail] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [isResending, setIsResending] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [canResend, setCanResend] = useState(false);
  const [countdown, setCountdown] = useState(60);

  // Forms
  const emailForm = useForm<EmailFormData>({
    resolver: zodResolver(emailSchema),
    defaultValues: { email: "" },
  });

  const otpForm = useForm<OTPFormData>({
    resolver: zodResolver(otpSchema),
  });

  const passwordForm = useForm<PasswordFormData>({
    resolver: zodResolver(passwordSchema),
  });

  const otpValue = otpForm.watch("otp") || "";

  useEffect(() => {
    // Check if email is passed from login page
    const emailParam = searchParams.get("email");
    if (emailParam) {
      const decodedEmail = decodeURIComponent(emailParam);
      setEmail(decodedEmail);
      emailForm.setValue("email", decodedEmail);
      // If email is provided, skip to OTP step after sending OTP
      handleEmailSubmit({ email: decodedEmail });
    }
  }, [searchParams]);

  useEffect(() => {
    // Countdown timer for resend button
    if (countdown > 0 && !canResend && currentStep === "otp") {
      const timer = setTimeout(() => setCountdown(countdown - 1), 1000);
      return () => clearTimeout(timer);
    } else if (countdown === 0) {
      setCanResend(true);
    }
  }, [countdown, canResend, currentStep]);

  const handleEmailSubmit = async (data: EmailFormData) => {
    try {
      setIsLoading(true);

      const response = await AuthService.forgotPassword(data);

      if (response.success) {
        setEmail(data.email);
        setCurrentStep("otp");
        setCanResend(false);
        setCountdown(60);
        toast.success("Password reset code sent to your email");
      } else {
        // For security, always show success message
        setEmail(data.email);
        setCurrentStep("otp");
        setCanResend(false);
        setCountdown(60);
        toast.success(
          "If an account exists, you will receive a password reset code"
        );
      }
    } catch (error: unknown) {
      console.error("Forgot password error:", error);
      // Always show success for security
      setEmail(data.email);
      setCurrentStep("otp");
      setCanResend(false);
      setCountdown(60);
      toast.success(
        "If an account exists, you will receive a password reset code"
      );
    } finally {
      setIsLoading(false);
    }
  };

  const handleOTPSubmit = async (data: OTPFormData) => {
    try {
      setIsLoading(true);

      // Just proceed to password step - we'll verify OTP when resetting password
      setCurrentStep("password");
      otpForm.setValue("otp", data.otp);
    } catch (error: unknown) {
      console.error("OTP verification error:", error);
      toast.error("Invalid verification code");
    } finally {
      setIsLoading(false);
    }
  };

  const handlePasswordSubmit = async (data: PasswordFormData) => {
    try {
      setIsLoading(true);

      const response = await AuthService.resetPassword({
        email,
        otp: otpForm.getValues("otp"),
        new_password: data.new_password,
      });

      if (response.success) {
        setCurrentStep("success");
        toast.success("Password reset successfully!");
      } else {
        throw new Error(response.message || "Password reset failed");
      }
    } catch (error: unknown) {
      console.error("Reset password error:", error);
      toast.error(
        (error as Error)?.message ||
          "Invalid or expired code. Please try again."
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
        otp_type: "password_reset",
      });

      if (response.success) {
        toast.success("New reset code sent to your email");
        setCanResend(false);
        setCountdown(60);
      } else {
        throw new Error(response.message || "Failed to resend code");
      }
    } catch (error: unknown) {
      console.error("Resend OTP error:", error);
      toast.error((error as Error)?.message || "Failed to resend reset code");
    } finally {
      setIsResending(false);
    }
  };

  const handleOTPChange = (value: string) => {
    otpForm.setValue("otp", value);
  };

  const goBack = () => {
    if (currentStep === "otp") {
      setCurrentStep("email");
    } else if (currentStep === "password") {
      setCurrentStep("otp");
    } else {
      router.push("/auth/login");
    }
  };

  const goToLogin = () => {
    router.push("/auth/login?message=password-reset-success");
  };

  if (!isReady) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  const renderEmailStep = () => (
    <form onSubmit={emailForm.handleSubmit(handleEmailSubmit)}>
      <CardContent className="space-y-4 pb-2">
        <div className="space-y-2">
          <Label htmlFor="email">Email address</Label>
          <Input
            id="email"
            type="email"
            placeholder="Enter your email"
            disabled={isLoading}
            {...emailForm.register("email")}
          />
          {emailForm.formState.errors.email && (
            <p className="text-sm text-red-600">
              {emailForm.formState.errors.email.message}
            </p>
          )}
        </div>
      </CardContent>

      <CardFooter className="flex flex-col space-y-4 pt-4">
        <Button type="submit" className="w-full" disabled={isLoading}>
          {isLoading ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Sending...
            </>
          ) : (
            "Send reset code"
          )}
        </Button>

        <div className="text-center">
          <Link
            href="/auth/login"
            className="inline-flex items-center text-sm text-blue-600 hover:text-blue-500"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to sign in
          </Link>
        </div>
      </CardFooter>
    </form>
  );

  const renderOTPStep = () => (
    <form onSubmit={otpForm.handleSubmit(handleOTPSubmit)}>
      <CardContent className="space-y-8">
        <div className="space-y-2">
          <Label htmlFor="otp" className="text-center block">
            Enter reset code
          </Label>
          <div className="flex justify-center">
            <InputOTP
              maxLength={6}
              value={otpValue}
              onChange={handleOTPChange}
              disabled={isLoading}
            >
              <InputOTPGroup>
                <InputOTPSlot index={0} />
                <InputOTPSlot index={1} />
                <InputOTPSlot index={2} />
                <InputOTPSlot index={3} />
                <InputOTPSlot index={4} />
                <InputOTPSlot index={5} />
              </InputOTPGroup>
            </InputOTP>
          </div>
          {otpForm.formState.errors.otp && (
            <p className="text-sm text-red-600 text-center">
              {otpForm.formState.errors.otp.message}
            </p>
          )}
        </div>

        <Alert className="mb-6">
          <Mail className="h-4 w-4" />
          <AlertDescription>
            We&apos;ve sent a 6-digit code to {email}. The code will expire in
            10 minutes.
          </AlertDescription>
        </Alert>
      </CardContent>

      <CardFooter className="flex flex-col space-y-4">
        <Button
          type="submit"
          className="w-full"
          disabled={isLoading || otpValue.length !== 6}
        >
          {isLoading ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Verifying...
            </>
          ) : (
            "Verify Code"
          )}
        </Button>

        <div className="text-center space-y-2">
          <p className="text-sm text-gray-600">Didn&apos;t receive the code?</p>
          <Button
            type="button"
            variant="ghost"
            onClick={handleResendOTP}
            disabled={!canResend || isResending}
            className="text-blue-600 hover:text-blue-500"
          >
            {isResending ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Resending...
              </>
            ) : canResend ? (
              "Resend code"
            ) : (
              `Resend in ${countdown}s`
            )}
          </Button>
        </div>
      </CardFooter>
    </form>
  );

  const renderPasswordStep = () => (
    <form onSubmit={passwordForm.handleSubmit(handlePasswordSubmit)}>
      <CardContent className="space-y-6">
        <div className="space-y-2">
          <Label htmlFor="new_password">New Password</Label>
          <div className="relative">
            <Input
              id="new_password"
              type={showPassword ? "text" : "password"}
              {...passwordForm.register("new_password")}
              placeholder="Enter your new password"
              disabled={isLoading}
            />
            <Button
              type="button"
              variant="ghost"
              size="sm"
              className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
              onClick={() => setShowPassword(!showPassword)}
              disabled={isLoading}
            >
              {showPassword ? (
                <EyeOff className="h-4 w-4" />
              ) : (
                <Eye className="h-4 w-4" />
              )}
            </Button>
          </div>
          {passwordForm.formState.errors.new_password && (
            <p className="text-sm text-red-600">
              {passwordForm.formState.errors.new_password.message}
            </p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="confirm_password">Confirm New Password</Label>
          <div className="relative">
            <Input
              id="confirm_password"
              type={showConfirmPassword ? "text" : "password"}
              {...passwordForm.register("confirm_password")}
              placeholder="Confirm your new password"
              disabled={isLoading}
            />
            <Button
              type="button"
              variant="ghost"
              size="sm"
              className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
              onClick={() => setShowConfirmPassword(!showConfirmPassword)}
              disabled={isLoading}
            >
              {showConfirmPassword ? (
                <EyeOff className="h-4 w-4" />
              ) : (
                <Eye className="h-4 w-4" />
              )}
            </Button>
          </div>
          {passwordForm.formState.errors.confirm_password && (
            <p className="text-sm text-red-600">
              {passwordForm.formState.errors.confirm_password.message}
            </p>
          )}
        </div>

        <Alert>
          <CheckCircle className="h-4 w-4" />
          <AlertDescription className="text-sm">
            Password must contain at least 8 characters including uppercase,
            lowercase, numbers, and special characters.
          </AlertDescription>
        </Alert>
      </CardContent>

      <CardFooter className="flex flex-col space-y-4">
        <Button type="submit" className="w-full" disabled={isLoading}>
          {isLoading ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Resetting password...
            </>
          ) : (
            "Reset Password"
          )}
        </Button>
      </CardFooter>
    </form>
  );

  const renderSuccessStep = () => (
    <>
      <CardContent className="space-y-4 text-center">
        <div className="mx-auto w-12 h-12 bg-green-100 rounded-full flex items-center justify-center mb-4">
          <CheckCircle className="h-6 w-6 text-green-600" />
        </div>
        <p className="text-gray-600">
          Your password has been successfully reset. You can now sign in with
          your new password.
        </p>
      </CardContent>

      <CardFooter>
        <Button onClick={goToLogin} className="w-full">
          Continue to Sign In
        </Button>
      </CardFooter>
    </>
  );

  const getStepInfo = () => {
    switch (currentStep) {
      case "email":
        return {
          title: "Forgot password?",
          description:
            "Enter your email address and we'll send you a code to reset your password",
          icon: <KeyRound className="h-6 w-6 text-blue-600" />,
        };
      case "otp":
        return {
          title: "Enter reset code",
          description: `We've sent a reset code to ${email}`,
          icon: <Mail className="h-6 w-6 text-blue-600" />,
        };
      case "password":
        return {
          title: "Create new password",
          description: "Enter your new password below",
          icon: <KeyRound className="h-6 w-6 text-blue-600" />,
        };
      case "success":
        return {
          title: "Password reset successful",
          description: "Your password has been updated",
          icon: <CheckCircle className="h-6 w-6 text-green-600" />,
        };
      default:
        return {
          title: "Forgot password?",
          description:
            "Enter your email address and we'll send you a code to reset your password",
          icon: <KeyRound className="h-6 w-6 text-blue-600" />,
        };
    }
  };

  const stepInfo = getStepInfo();

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1 text-center">
          <div className="mx-auto w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center mb-4">
            {stepInfo.icon}
          </div>
          <CardTitle className="text-2xl font-bold">{stepInfo.title}</CardTitle>
          <CardDescription>{stepInfo.description}</CardDescription>
        </CardHeader>

        {currentStep === "email" && renderEmailStep()}
        {currentStep === "otp" && renderOTPStep()}
        {currentStep === "password" && renderPasswordStep()}
        {currentStep === "success" && renderSuccessStep()}

        {currentStep !== "email" && currentStep !== "success" && (
          <div className="px-6 pb-6">
            <div className="text-center">
              <button
                onClick={goBack}
                className="text-sm text-gray-600 hover:text-gray-500"
                disabled={isLoading}
              >
                ← Go back
              </button>
            </div>
          </div>
        )}
      </Card>
    </div>
  );
}
