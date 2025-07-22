// src/app/(auth)/verify-email/page.tsx

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
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
} from "@/components/ui/input-otp";
import { Label } from "@/components/ui/label";
import { AuthService } from "@/lib/api/auth-service";
import { zodResolver } from "@hookform/resolvers/zod";
import { CheckCircle, Loader2, Mail } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

const verifyEmailSchema = z.object({
  otp: z.string().length(6, "OTP must be 6 digits"),
});

type VerifyEmailFormData = z.infer<typeof verifyEmailSchema>;

export default function VerifyEmailPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [isLoading, setIsLoading] = useState(false);
  const [isResending, setIsResending] = useState(false);
  const [email, setEmail] = useState("");
  const [canResend, setCanResend] = useState(false);
  const [countdown, setCountdown] = useState(60);

  const {
    setValue,
    handleSubmit,
    formState: { errors },
    watch,
  } = useForm<VerifyEmailFormData>({
    resolver: zodResolver(verifyEmailSchema),
  });

  const otpValue = watch("otp") || "";

  useEffect(() => {
    const emailParam = searchParams.get("email");
    if (emailParam) {
      setEmail(decodeURIComponent(emailParam));
    } else {
      // Redirect to signup if no email provided
      router.push("/auth/signup");
    }
  }, [searchParams, router]);

  useEffect(() => {
    // Countdown timer for resend button
    if (countdown > 0 && !canResend) {
      const timer = setTimeout(() => setCountdown(countdown - 1), 1000);
      return () => clearTimeout(timer);
    } else if (countdown === 0) {
      setCanResend(true);
    }
  }, [countdown, canResend]);

  const onSubmit = async (data: VerifyEmailFormData) => {
    if (!email) {
      toast.error("Email is required");
      return;
    }

    try {
      setIsLoading(true);

      const response = await AuthService.verifyEmail({
        email,
        otp: data.otp,
      });

      if (response.success) {
        toast.success("Email verified successfully!");
        // Redirect to login with success message
        router.push("/auth/login?message=verification-success");
      } else {
        throw new Error(response.message || "Verification failed");
      }
    } catch (error: any) {
      console.error("Verification error:", error);
      toast.error(error.message || "Invalid or expired OTP. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleResendOTP = async () => {
    try {
      setIsResending(true);

      const response = await AuthService.resendOTP({
        email,
        otp_type: "email_verification",
      });

      if (response.success) {
        toast.success("Verification code sent to your email");
        setCanResend(false);
        setCountdown(60);
      } else {
        throw new Error(response.message || "Failed to resend OTP");
      }
    } catch (error: any) {
      console.error("Resend OTP error:", error);
      toast.error(error.message || "Failed to resend verification code");
    } finally {
      setIsResending(false);
    }
  };

  const handleOTPChange = (value: string) => {
    setValue("otp", value);
    // Auto-submit when OTP is complete
    if (value.length === 6) {
      handleSubmit(onSubmit)();
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1 text-center">
          <div className="mx-auto w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center mb-4">
            <Mail className="h-6 w-6 text-blue-600" />
          </div>
          <CardTitle className="text-2xl font-bold">
            Verify your email
          </CardTitle>
          <CardDescription>
            We've sent a 6-digit verification code to
            <br />
            <span className="font-medium text-gray-900">{email}</span>
          </CardDescription>
        </CardHeader>

        <form onSubmit={handleSubmit(onSubmit)}>
          <CardContent className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="otp" className="text-center block">
                Enter verification code
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
              {errors.otp && (
                <p className="text-sm text-red-600 text-center">
                  {errors.otp.message}
                </p>
              )}
            </div>

            <Alert>
              <CheckCircle className="h-4 w-4" />
              <AlertDescription>
                Enter the 6-digit code sent to your email. The code will expire
                in 10 minutes.
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
                "Verify Email"
              )}
            </Button>

            <div className="text-center space-y-2">
              <p className="text-sm text-gray-600">Didn't receive the code?</p>
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

            <div className="text-center text-sm">
              Wrong email?{" "}
              <Button
                type="button"
                variant="ghost"
                onClick={() => router.push("/auth/signup")}
                className="text-blue-600 hover:text-blue-500 p-0 h-auto font-medium"
              >
                Sign up again
              </Button>
            </div>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}
