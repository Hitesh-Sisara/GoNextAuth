// File: src/components/ui/loading.tsx

"use client";

import { cn } from "@/lib/utils";
import { AlertCircle, CheckCircle, Loader2 } from "lucide-react";
import React from "react";

// Basic spinner component
export interface SpinnerProps {
  size?: "sm" | "md" | "lg" | "xl";
  className?: string;
}

export const Spinner: React.FC<SpinnerProps> = ({ size = "md", className }) => {
  const sizeClasses = {
    sm: "h-4 w-4",
    md: "h-6 w-6",
    lg: "h-8 w-8",
    xl: "h-12 w-12",
  };

  return (
    <Loader2 className={cn("animate-spin", sizeClasses[size], className)} />
  );
};

// Loading overlay component
export interface LoadingOverlayProps {
  isVisible: boolean;
  message?: string;
  className?: string;
}

export const LoadingOverlay: React.FC<LoadingOverlayProps> = ({
  isVisible,
  message = "Loading...",
  className,
}) => {
  if (!isVisible) return null;

  return (
    <div
      className={cn(
        "fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50",
        className
      )}
    >
      <div className="bg-white rounded-lg p-6 flex flex-col items-center space-y-4 max-w-sm mx-4">
        <Spinner size="lg" />
        <p className="text-gray-700 text-center">{message}</p>
      </div>
    </div>
  );
};

// Inline loading component
export interface InlineLoadingProps {
  isLoading: boolean;
  message?: string;
  size?: "sm" | "md" | "lg";
  className?: string;
}

export const InlineLoading: React.FC<InlineLoadingProps> = ({
  isLoading,
  message,
  size = "md",
  className,
}) => {
  if (!isLoading) return null;

  return (
    <div className={cn("flex items-center space-x-2", className)}>
      <Spinner size={size} />
      {message && <span className="text-gray-600">{message}</span>}
    </div>
  );
};

// Loading button component
export interface LoadingButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  isLoading: boolean;
  loadingText?: string;
  children: React.ReactNode;
  variant?:
    | "default"
    | "destructive"
    | "outline"
    | "secondary"
    | "ghost"
    | "link";
  size?: "default" | "sm" | "lg" | "icon";
}

export const LoadingButton: React.FC<LoadingButtonProps> = ({
  isLoading,
  loadingText,
  children,
  disabled,
  className,
  ...props
}) => {
  return (
    <button
      {...props}
      disabled={disabled || isLoading}
      className={cn(
        "inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
        "disabled:opacity-50 disabled:pointer-events-none",
        "bg-primary text-primary-foreground hover:bg-primary/90",
        "h-10 py-2 px-4",
        className
      )}
    >
      {isLoading && <Spinner size="sm" className="mr-2" />}
      {isLoading ? loadingText || "Loading..." : children}
    </button>
  );
};

// Loading state component for forms
export interface LoadingStateProps {
  isLoading: boolean;
  isSuccess?: boolean;
  isError?: boolean;
  loadingMessage?: string;
  successMessage?: string;
  errorMessage?: string;
  className?: string;
}

export const LoadingState: React.FC<LoadingStateProps> = ({
  isLoading,
  isSuccess,
  isError,
  loadingMessage = "Processing...",
  successMessage = "Success!",
  errorMessage = "Something went wrong",
  className,
}) => {
  if (!isLoading && !isSuccess && !isError) return null;

  return (
    <div
      className={cn("flex items-center space-x-2 p-3 rounded-md", className)}
    >
      {isLoading && (
        <>
          <Spinner size="sm" />
          <span className="text-gray-600">{loadingMessage}</span>
        </>
      )}

      {isSuccess && (
        <>
          <CheckCircle className="h-5 w-5 text-green-600" />
          <span className="text-green-700">{successMessage}</span>
        </>
      )}

      {isError && (
        <>
          <AlertCircle className="h-5 w-5 text-red-600" />
          <span className="text-red-700">{errorMessage}</span>
        </>
      )}
    </div>
  );
};

// Page loading skeleton
export const PageLoadingSkeleton: React.FC = () => {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center space-y-4">
        <Spinner size="xl" className="mx-auto" />
        <div className="space-y-2">
          <div className="h-4 bg-gray-200 rounded animate-pulse w-48 mx-auto"></div>
          <div className="h-3 bg-gray-200 rounded animate-pulse w-32 mx-auto"></div>
        </div>
      </div>
    </div>
  );
};

// Form field loading skeleton
export const FieldLoadingSkeleton: React.FC<{ className?: string }> = ({
  className,
}) => {
  return (
    <div className={cn("space-y-2", className)}>
      <div className="h-4 bg-gray-200 rounded animate-pulse w-24"></div>
      <div className="h-10 bg-gray-200 rounded animate-pulse w-full"></div>
    </div>
  );
};

// Card loading skeleton
export const CardLoadingSkeleton: React.FC<{ className?: string }> = ({
  className,
}) => {
  return (
    <div className={cn("border rounded-lg p-6 space-y-4", className)}>
      <div className="space-y-2">
        <div className="h-6 bg-gray-200 rounded animate-pulse w-3/4"></div>
        <div className="h-4 bg-gray-200 rounded animate-pulse w-1/2"></div>
      </div>
      <div className="space-y-2">
        <div className="h-3 bg-gray-200 rounded animate-pulse w-full"></div>
        <div className="h-3 bg-gray-200 rounded animate-pulse w-5/6"></div>
        <div className="h-3 bg-gray-200 rounded animate-pulse w-4/6"></div>
      </div>
    </div>
  );
};

// Loading dots animation
export const LoadingDots: React.FC<{ className?: string }> = ({
  className,
}) => {
  return (
    <div className={cn("flex space-x-1", className)}>
      {[0, 1, 2].map((i) => (
        <div
          key={i}
          className="w-2 h-2 bg-current rounded-full animate-pulse"
          style={{
            animationDelay: `${i * 0.15}s`,
            animationDuration: "1s",
          }}
        />
      ))}
    </div>
  );
};

// Progress bar component
export interface ProgressBarProps {
  progress: number; // 0-100
  className?: string;
  showPercentage?: boolean;
}

export const ProgressBar: React.FC<ProgressBarProps> = ({
  progress,
  className,
  showPercentage = false,
}) => {
  const clampedProgress = Math.min(Math.max(progress, 0), 100);

  return (
    <div className={cn("w-full", className)}>
      <div className="flex justify-between items-center mb-1">
        {showPercentage && (
          <span className="text-sm text-gray-600">
            {Math.round(clampedProgress)}%
          </span>
        )}
      </div>
      <div className="w-full bg-gray-200 rounded-full h-2">
        <div
          className="bg-blue-600 h-2 rounded-full transition-all duration-300 ease-out"
          style={{ width: `${clampedProgress}%` }}
        />
      </div>
    </div>
  );
};
