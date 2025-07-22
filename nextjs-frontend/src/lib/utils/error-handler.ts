// File: src/lib/utils/error-handler.ts

import { toast } from "sonner";

export interface APIError {
  success: false;
  message: string;
  error?:
    | {
        message: string;
        field?: string;
        code?: string;
      }
    | Array<{
        field: string;
        message: string;
        value?: string;
      }>;
}

export interface APIResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
  error?: any;
}

export class AuthError extends Error {
  constructor(message: string, public code?: string, public field?: string) {
    super(message);
    this.name = "AuthError";
  }
}

export class NetworkError extends Error {
  constructor(
    message: string,
    public status?: number,
    public statusText?: string
  ) {
    super(message);
    this.name = "NetworkError";
  }
}

export class ValidationError extends Error {
  constructor(message: string, public field: string, public value?: string) {
    super(message);
    this.name = "ValidationError";
  }
}

// Error message mapping for better UX
const ERROR_MESSAGES: Record<string, string> = {
  // Network errors
  NETWORK_ERROR: "Please check your internet connection and try again.",
  TIMEOUT_ERROR: "Request timed out. Please try again.",
  SERVER_ERROR: "Server is temporarily unavailable. Please try again later.",

  // Auth errors
  INVALID_CREDENTIALS: "Invalid email or password. Please try again.",
  ACCOUNT_LOCKED:
    "Your account has been temporarily locked. Please contact support.",
  EMAIL_NOT_VERIFIED: "Please verify your email address before signing in.",
  ACCOUNT_DISABLED: "Your account has been disabled. Please contact support.",
  TOKEN_EXPIRED: "Your session has expired. Please sign in again.",

  // OTP errors
  INVALID_OTP: "Invalid verification code. Please check and try again.",
  OTP_EXPIRED: "Verification code has expired. Please request a new one.",
  OTP_LIMIT_EXCEEDED: "Too many attempts. Please wait before trying again.",

  // Rate limiting
  RATE_LIMIT_EXCEEDED:
    "Too many requests. Please wait a moment before trying again.",

  // Validation errors
  INVALID_EMAIL: "Please enter a valid email address.",
  WEAK_PASSWORD:
    "Password must be at least 8 characters with uppercase, lowercase, numbers, and special characters.",
  PASSWORDS_DONT_MATCH: "Passwords don't match. Please try again.",

  // Google OAuth errors
  GOOGLE_AUTH_CANCELLED: "Google sign-in was cancelled.",
  GOOGLE_AUTH_FAILED: "Google sign-in failed. Please try again.",
  GOOGLE_AUTH_POPUP_BLOCKED:
    "Popup was blocked. Please allow popups and try again.",
};

// Get user-friendly error message
export const getErrorMessage = (error: any): string => {
  if (typeof error === "string") {
    return ERROR_MESSAGES[error] || error;
  }

  if (
    error instanceof AuthError ||
    error instanceof NetworkError ||
    error instanceof ValidationError
  ) {
    return ERROR_MESSAGES[error.message] || error.message;
  }

  if (error?.message) {
    return ERROR_MESSAGES[error.message] || error.message;
  }

  if (error?.error?.message) {
    return ERROR_MESSAGES[error.error.message] || error.error.message;
  }

  return "An unexpected error occurred. Please try again.";
};

// Handle API errors with toast notifications
export const handleAPIError = (error: any, context?: string): void => {
  console.error(`API Error${context ? ` (${context})` : ""}:`, error);

  let errorMessage = "An unexpected error occurred.";
  let errorTitle = "Error";

  if (error?.response?.status) {
    const status = error.response.status;

    switch (status) {
      case 400:
        errorTitle = "Invalid Request";
        errorMessage = getErrorMessage(
          error?.response?.data?.message ||
            "Please check your input and try again."
        );
        break;
      case 401:
        errorTitle = "Authentication Failed";
        errorMessage = getErrorMessage(
          error?.response?.data?.message || "Please sign in again."
        );
        break;
      case 403:
        errorTitle = "Access Denied";
        errorMessage = "You don't have permission to perform this action.";
        break;
      case 404:
        errorTitle = "Not Found";
        errorMessage = "The requested resource was not found.";
        break;
      case 409:
        errorTitle = "Conflict";
        errorMessage = getErrorMessage(
          error?.response?.data?.message ||
            "A conflict occurred with your request."
        );
        break;
      case 422:
        errorTitle = "Validation Error";
        const validationErrors = error?.response?.data?.error;
        if (Array.isArray(validationErrors) && validationErrors.length > 0) {
          errorMessage = validationErrors.map((e) => e.message).join(", ");
        } else {
          errorMessage = getErrorMessage(
            error?.response?.data?.message || "Please check your input."
          );
        }
        break;
      case 429:
        errorTitle = "Too Many Requests";
        errorMessage = "Please wait a moment before trying again.";
        break;
      case 500:
      case 502:
      case 503:
      case 504:
        errorTitle = "Server Error";
        errorMessage =
          "Server is temporarily unavailable. Please try again later.";
        break;
      default:
        errorMessage = getErrorMessage(
          error?.response?.data?.message || "An unexpected error occurred."
        );
    }
  } else if (
    error?.code === "NETWORK_ERROR" ||
    error?.name === "NetworkError"
  ) {
    errorTitle = "Connection Error";
    errorMessage = "Please check your internet connection and try again.";
  } else if (error?.code === "TIMEOUT") {
    errorTitle = "Request Timeout";
    errorMessage = "Request timed out. Please try again.";
  } else {
    errorMessage = getErrorMessage(error?.message || error);
  }

  // Show toast with error
  toast.error(errorMessage, {
    description: context,
    duration: 5000,
  });
};

// Handle success responses
export const handleAPISuccess = (
  message: string,
  description?: string
): void => {
  toast.success(message, {
    description,
    duration: 3000,
  });
};

// Handle validation errors specifically
export const handleValidationErrors = (
  errors: Array<{ field: string; message: string }>
): void => {
  const errorMessages = errors
    .map((error) => `${error.field}: ${error.message}`)
    .join("\n");

  toast.error("Please fix the following errors:", {
    description: errorMessages,
    duration: 5000,
  });
};

// Async error wrapper for better error handling
export const withErrorHandling = async <T>(
  asyncFunction: () => Promise<T>,
  context?: string
): Promise<T | null> => {
  try {
    return await asyncFunction();
  } catch (error) {
    handleAPIError(error, context);
    return null;
  }
};

// Retry mechanism with exponential backoff
export const withRetry = async <T>(
  asyncFunction: () => Promise<T>,
  maxRetries: number = 3,
  baseDelay: number = 1000
): Promise<T> => {
  let lastError;

  for (let attempt = 1; attempt <= maxRetries; attempt++) {
    try {
      return await asyncFunction();
    } catch (error) {
      lastError = error;

      if (attempt === maxRetries) {
        throw error;
      }

      // Exponential backoff
      const delay = baseDelay * Math.pow(2, attempt - 1);
      await new Promise((resolve) => setTimeout(resolve, delay));

      console.log(`Retry attempt ${attempt} failed, retrying in ${delay}ms...`);
    }
  }

  throw lastError;
};
