// File: src/lib/api/auth-service.ts

import {
  APIResponse,
  AuthResponse,
  ForgotPasswordRequest,
  LoginRequest,
  ResendOTPRequest,
  ResetPasswordRequest,
  SignupRequest,
  User,
  VerifyEmailRequest,
} from "@/types/auth";
import { apiClient } from "./client";

// New request types for multi-step auth
export interface EmailOnlyRequest {
  email: string;
}

export interface OTPLoginRequest {
  email: string;
  otp: string;
}

export interface CompleteSignupRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  phone?: string;
}

export interface GoogleAuthRequest {
  access_token: string;
}

export interface GoogleAuthURLResponse {
  auth_url: string;
}

export class AuthService {
  // Google OAuth endpoints
  static async getGoogleAuthURL(): Promise<APIResponse<GoogleAuthURLResponse>> {
    return apiClient.publicRequest("GET", "/auth/google/url", null);
  }

  static async googleAuth(
    data: GoogleAuthRequest
  ): Promise<APIResponse<AuthResponse>> {
    return apiClient.publicRequest("POST", "/auth/google/token", data);
  }

  // Updated to include state parameter
  static async googleCallback(
    code: string,
    state: string
  ): Promise<APIResponse<AuthResponse>> {
    return apiClient.publicRequest(
      "GET",
      `/auth/google/callback?code=${encodeURIComponent(
        code
      )}&state=${encodeURIComponent(state)}`,
      null
    );
  }

  // New multi-step authentication endpoints
  static async initiateSignup(data: EmailOnlyRequest): Promise<APIResponse> {
    return apiClient.publicRequest("POST", "/auth/signup/initiate", data);
  }

  static async completeSignup(
    data: CompleteSignupRequest
  ): Promise<APIResponse<User>> {
    return apiClient.publicRequest("POST", "/auth/signup/complete", data);
  }

  static async initiateEmailLogin(
    data: EmailOnlyRequest
  ): Promise<APIResponse> {
    return apiClient.publicRequest("POST", "/auth/login/email", data);
  }

  static async completeOTPLogin(
    data: OTPLoginRequest
  ): Promise<APIResponse<AuthResponse>> {
    return apiClient.publicRequest("POST", "/auth/login/otp", data);
  }

  // Legacy authentication endpoints (for backward compatibility)
  static async signup(data: SignupRequest): Promise<APIResponse<User>> {
    return apiClient.publicRequest("POST", "/auth/signup", data);
  }

  static async login(data: LoginRequest): Promise<APIResponse<AuthResponse>> {
    return apiClient.publicRequest("POST", "/auth/login", data);
  }

  static async verifyEmail(data: VerifyEmailRequest): Promise<APIResponse> {
    return apiClient.publicRequest("POST", "/auth/verify-email", data);
  }

  static async forgotPassword(
    data: ForgotPasswordRequest
  ): Promise<APIResponse> {
    return apiClient.publicRequest("POST", "/auth/forgot-password", data);
  }

  static async resetPassword(data: ResetPasswordRequest): Promise<APIResponse> {
    return apiClient.publicRequest("POST", "/auth/reset-password", data);
  }

  static async resendOTP(data: ResendOTPRequest): Promise<APIResponse> {
    return apiClient.publicRequest("POST", "/auth/resend-otp", data);
  }

  static async refreshToken(
    refreshToken: string
  ): Promise<APIResponse<AuthResponse>> {
    return apiClient.publicRequest("POST", "/auth/refresh", {
      refresh_token: refreshToken,
    });
  }

  static async logout(refreshToken: string): Promise<APIResponse> {
    return apiClient.publicRequest("POST", "/auth/logout", {
      refresh_token: refreshToken,
    });
  }

  // Protected endpoints
  static async getProfile(): Promise<APIResponse<User>> {
    return apiClient.get("/auth/profile");
  }

  static async updateProfile(data: Partial<User>): Promise<APIResponse<User>> {
    return apiClient.put("/auth/profile", data);
  }

  // Utility methods
  static async checkAuthStatus(): Promise<{
    isAuthenticated: boolean;
    user?: User;
  }> {
    try {
      const response = await this.getProfile();
      if (response.success && response.data) {
        return {
          isAuthenticated: true,
          user: response.data,
        };
      }
      return { isAuthenticated: false };
    } catch (error) {
      return { isAuthenticated: false };
    }
  }
}
