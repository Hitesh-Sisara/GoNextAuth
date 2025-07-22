// File: src/components/auth/login-choice.tsx

"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { brandConfig } from "@/lib/config/app-config";
import { Eye, EyeOff, Loader2, Mail } from "lucide-react";
import { useState } from "react";

interface LoginChoiceProps {
  email: string;
  onPasswordLogin: (password: string) => Promise<void>;
  onOTPLogin: () => Promise<void>;
  onForgotPassword?: () => void;
  isLoading: boolean;
  title?: string;
}

export function LoginChoice({
  email,
  onPasswordLogin,
  onOTPLogin,
  onForgotPassword,
  isLoading,
  title = `Sign in to ${brandConfig.name}`,
}: LoginChoiceProps) {
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);

  const handlePasswordSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (password.trim()) {
      await onPasswordLogin(password);
    }
  };

  const handleOTPSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onOTPLogin();
  };

  return (
    <div className="space-y-6">
      {/* Display email */}
      <div className="space-y-2">
        <Label>Email</Label>
        <div className="text-sm font-medium text-gray-900 bg-gray-50 px-3 py-2 rounded-md border">
          {email}
        </div>
      </div>

      {/* Password login form */}
      <form onSubmit={handlePasswordSubmit} className="space-y-4">
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <Label htmlFor="password">Password</Label>
            {onForgotPassword && (
              <button
                type="button"
                onClick={onForgotPassword}
                className="text-sm text-blue-600 hover:text-blue-500"
                disabled={isLoading}
              >
                Forgot your password?
              </button>
            )}
          </div>
          <div className="relative">
            <Input
              id="password"
              type={showPassword ? "text" : "password"}
              placeholder="Your password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
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
        </div>

        <Button
          type="submit"
          className="w-full"
          disabled={isLoading || !password.trim()}
        >
          {isLoading ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Signing in...
            </>
          ) : (
            "Sign in"
          )}
        </Button>
      </form>

      {/* Separator */}
      <div className="relative">
        <div className="absolute inset-0 flex items-center">
          <Separator className="w-full" />
        </div>
        <div className="relative flex justify-center text-xs uppercase">
          <span className="bg-white px-2 text-gray-500">OR</span>
        </div>
      </div>

      {/* OTP login option */}
      <form onSubmit={handleOTPSubmit}>
        <Button
          type="submit"
          variant="outline"
          className="w-full"
          disabled={isLoading}
        >
          {isLoading ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Sending code...
            </>
          ) : (
            <>
              <Mail className="mr-2 h-4 w-4" />
              Email sign-in code
            </>
          )}
        </Button>
      </form>
    </div>
  );
}
