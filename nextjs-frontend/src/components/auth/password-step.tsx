// File: src/components/auth/password-step.tsx

"use client";

import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { validatePassword } from "@/lib/utils";
import { CheckCircle, Eye, EyeOff, Loader2, X } from "lucide-react";
import { useState } from "react";

interface PasswordStepProps {
  onSubmit: (password: string) => Promise<void> | void;
  isLoading: boolean;
  title?: string;
  buttonText?: string;
  showConfirmPassword?: boolean;
}

export function PasswordStep({
  onSubmit,
  isLoading,
  title = "Create your password",
  buttonText = "Continue",
  showConfirmPassword = true,
}: PasswordStepProps) {
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPasswordField, setShowConfirmPasswordField] =
    useState(false);
  const [passwordError, setPasswordError] = useState("");
  const [confirmPasswordError, setConfirmPasswordError] = useState("");

  const passwordValidation = validatePassword(password);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Reset errors
    setPasswordError("");
    setConfirmPasswordError("");

    // Validate password
    if (!password) {
      setPasswordError("Password is required");
      return;
    }

    if (!passwordValidation.isValid) {
      setPasswordError("Password does not meet requirements");
      return;
    }

    // Validate confirm password if shown
    if (showConfirmPassword) {
      if (!confirmPassword) {
        setConfirmPasswordError("Please confirm your password");
        return;
      }

      if (password !== confirmPassword) {
        setConfirmPasswordError("Passwords don't match");
        return;
      }
    }

    await onSubmit(password);
  };

  const handlePasswordChange = (value: string) => {
    setPassword(value);
    if (passwordError) setPasswordError("");
  };

  const handleConfirmPasswordChange = (value: string) => {
    setConfirmPassword(value);
    if (confirmPasswordError) setConfirmPasswordError("");
  };

  return (
    <div className="space-y-6">
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="password">Password</Label>
          <div className="relative">
            <Input
              id="password"
              type={showPassword ? "text" : "password"}
              placeholder="Enter your password"
              value={password}
              onChange={(e) => handlePasswordChange(e.target.value)}
              disabled={isLoading}
              required
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
          {passwordError && (
            <p className="text-sm text-red-600">{passwordError}</p>
          )}
        </div>

        {showConfirmPassword && (
          <div className="space-y-2">
            <Label htmlFor="confirmPassword">Confirm Password</Label>
            <div className="relative">
              <Input
                id="confirmPassword"
                type={showConfirmPasswordField ? "text" : "password"}
                placeholder="Confirm your password"
                value={confirmPassword}
                onChange={(e) => handleConfirmPasswordChange(e.target.value)}
                disabled={isLoading}
                required
              />
              <Button
                type="button"
                variant="ghost"
                size="sm"
                className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                onClick={() =>
                  setShowConfirmPasswordField(!showConfirmPasswordField)
                }
                disabled={isLoading}
              >
                {showConfirmPasswordField ? (
                  <EyeOff className="h-4 w-4" />
                ) : (
                  <Eye className="h-4 w-4" />
                )}
              </Button>
            </div>
            {confirmPasswordError && (
              <p className="text-sm text-red-600">{confirmPasswordError}</p>
            )}
          </div>
        )}

        <Button
          type="submit"
          className="w-full"
          disabled={isLoading || !password || !passwordValidation.isValid}
        >
          {isLoading ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Creating...
            </>
          ) : (
            buttonText
          )}
        </Button>
      </form>

      {/* Password requirements */}
      {password && (
        <Alert
          className={
            passwordValidation.isValid ? "border-green-200 bg-green-50" : ""
          }
        >
          <CheckCircle
            className={`h-4 w-4 ${
              passwordValidation.isValid ? "text-green-600" : "text-gray-400"
            }`}
          />
          <AlertDescription>
            <div className="space-y-1">
              <p className="font-medium">Password requirements:</p>
              <div className="grid grid-cols-1 gap-1 text-sm">
                <div className="flex items-center gap-2">
                  {passwordValidation.errors.minLength ? (
                    <X className="h-3 w-3 text-red-500" />
                  ) : (
                    <CheckCircle className="h-3 w-3 text-green-500" />
                  )}
                  <span>At least 8 characters</span>
                </div>
                <div className="flex items-center gap-2">
                  {passwordValidation.errors.hasUpperCase ? (
                    <X className="h-3 w-3 text-red-500" />
                  ) : (
                    <CheckCircle className="h-3 w-3 text-green-500" />
                  )}
                  <span>One uppercase letter</span>
                </div>
                <div className="flex items-center gap-2">
                  {passwordValidation.errors.hasLowerCase ? (
                    <X className="h-3 w-3 text-red-500" />
                  ) : (
                    <CheckCircle className="h-3 w-3 text-green-500" />
                  )}
                  <span>One lowercase letter</span>
                </div>
                <div className="flex items-center gap-2">
                  {passwordValidation.errors.hasNumbers ? (
                    <X className="h-3 w-3 text-red-500" />
                  ) : (
                    <CheckCircle className="h-3 w-3 text-green-500" />
                  )}
                  <span>One number</span>
                </div>
                <div className="flex items-center gap-2">
                  {passwordValidation.errors.hasSpecialChar ? (
                    <X className="h-3 w-3 text-red-500" />
                  ) : (
                    <CheckCircle className="h-3 w-3 text-green-500" />
                  )}
                  <span>One special character</span>
                </div>
              </div>
            </div>
          </AlertDescription>
        </Alert>
      )}
    </div>
  );
}
