// File: src/components/auth/otp-step.tsx

"use client";

import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
} from "@/components/ui/input-otp";
import { Label } from "@/components/ui/label";
import { CheckCircle, Loader2, Mail } from "lucide-react";
import { useEffect, useState } from "react";

interface OTPStepProps {
  email: string;
  onSubmit: (otp: string) => Promise<void> | void;
  onResend?: () => Promise<void> | void;
  isLoading: boolean;
  isResending?: boolean;
  title?: string;
  description?: string;
  otpType?: "verification" | "login" | "reset";
}

export function OTPStep({
  email,
  onSubmit,
  onResend,
  isLoading,
  isResending = false,
  title = "Enter verification code",
  description = "We've sent a verification code to",
  otpType = "verification",
}: OTPStepProps) {
  const [otp, setOtp] = useState("");
  const [canResend, setCanResend] = useState(false);
  const [countdown, setCountdown] = useState(60);

  useEffect(() => {
    // Auto-submit when OTP is complete
    if (otp.length === 6 && !isLoading) {
      handleSubmit();
    }
  }, [otp, isLoading]);

  useEffect(() => {
    // Countdown timer for resend button
    if (countdown > 0 && !canResend) {
      const timer = setTimeout(() => setCountdown(countdown - 1), 1000);
      return () => clearTimeout(timer);
    } else if (countdown === 0) {
      setCanResend(true);
    }
  }, [countdown, canResend]);

  const handleSubmit = async () => {
    if (otp.length === 6) {
      await onSubmit(otp);
    }
  };

  const handleResend = async () => {
    if (onResend) {
      await onResend();
      setCanResend(false);
      setCountdown(60);
    }
  };

  const getAlertMessage = () => {
    switch (otpType) {
      case "login":
        return "Enter the 6-digit sign-in code sent to your email.";
      case "reset":
        return "Enter the 6-digit password reset code sent to your email.";
      default:
        return "Enter the 6-digit verification code sent to your email. The code will expire in 10 minutes.";
    }
  };

  return (
    <div className="space-y-6">
      <div className="text-center space-y-2">
        <p className="text-sm text-gray-600">
          {description}
          <br />
          <span className="font-medium text-gray-900">{email}</span>
        </p>
      </div>

      <div className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="otp" className="text-center block">
            {title}
          </Label>
          <div className="flex justify-center">
            <InputOTP
              maxLength={6}
              value={otp}
              onChange={setOtp}
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
        </div>

        <Alert>
          <CheckCircle className="h-4 w-4" />
          <AlertDescription>{getAlertMessage()}</AlertDescription>
        </Alert>
      </div>

      <div className="space-y-4">
        <Button
          onClick={handleSubmit}
          className="w-full"
          disabled={isLoading || otp.length !== 6}
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

        {onResend && (
          <div className="text-center space-y-2">
            <p className="text-sm text-gray-600">Didn't receive the code?</p>
            <Button
              type="button"
              variant="ghost"
              onClick={handleResend}
              disabled={!canResend || isResending}
              className="text-blue-600 hover:text-blue-500"
            >
              {isResending ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Resending...
                </>
              ) : canResend ? (
                <>
                  <Mail className="mr-2 h-4 w-4" />
                  Resend code
                </>
              ) : (
                `Resend in ${countdown}s`
              )}
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
