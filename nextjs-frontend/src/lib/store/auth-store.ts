// File: src/lib/store/auth-store.ts

import { AuthService } from "@/lib/api/auth-service";
import { TokenManager } from "@/lib/auth/token-manager";
import { brandConfig } from "@/lib/config/app-config";
import { AuthState, LoginRequest, SignupRequest, User } from "@/types/auth";
import { toast } from "sonner";
import { create } from "zustand";
import { persist } from "zustand/middleware";

interface AuthStore extends AuthState {
  // Actions
  login: (credentials: LoginRequest) => Promise<boolean>;
  signup: (data: SignupRequest) => Promise<boolean>;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
  clearError: () => void;
  setUser: (user: User | null) => void;
  setLoading: (loading: boolean) => void;
  refreshUserProfile: () => Promise<void>;

  // New state to prevent duplicate operations
  isLoggingOut: boolean;
}

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      // Initial state
      user: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,
      isLoggingOut: false, // New state

      // Actions
      login: async (credentials: LoginRequest) => {
        try {
          set({ isLoading: true, error: null });

          const response = await AuthService.login(credentials);

          if (response.success && response.data) {
            const { user, access_token, refresh_token, expires_in } =
              response.data;

            // Store tokens
            TokenManager.setTokens(access_token, refresh_token, expires_in);
            TokenManager.setUser(user);

            // Update store
            set({
              user,
              isAuthenticated: true,
              isLoading: false,
              error: null,
            });

            toast.success(`Welcome back to ${brandConfig.name}!`);
            return true;
          } else {
            throw new Error(response.message || "Login failed");
          }
        } catch (error: any) {
          const errorMessage = error.message || "Login failed";
          set({
            isLoading: false,
            error: errorMessage,
            isAuthenticated: false,
            user: null,
          });
          toast.error(errorMessage);
          return false;
        }
      },

      signup: async (data: SignupRequest) => {
        try {
          set({ isLoading: true, error: null });

          const response = await AuthService.signup(data);

          if (response.success) {
            set({ isLoading: false, error: null });
            toast.success(
              `Welcome to ${brandConfig.name}! Please check your email for verification.`
            );
            return true;
          } else {
            throw new Error(response.message || "Signup failed");
          }
        } catch (error: any) {
          const errorMessage = error.message || "Signup failed";
          set({ isLoading: false, error: errorMessage });
          toast.error(errorMessage);
          return false;
        }
      },

      logout: async () => {
        const state = get();

        // Prevent multiple simultaneous logout requests
        if (state.isLoggingOut) {
          console.log("Logout already in progress, skipping duplicate request");
          return;
        }

        try {
          console.log("Starting logout process");
          set({ isLoggingOut: true, isLoading: true });

          const refreshToken = TokenManager.getRefreshToken();

          if (refreshToken) {
            try {
              console.log("Calling logout API");
              await AuthService.logout(refreshToken);
              console.log("Logout API call successful");
            } catch (error) {
              console.error(
                "Logout API error (continuing with local logout):",
                error
              );
              // Don't throw - we still want to clear local state
            }
          } else {
            console.log("No refresh token found, skipping API call");
          }
        } catch (error) {
          console.error("Logout error:", error);
        } finally {
          // Always clear tokens and state, even if API call fails
          console.log("Clearing tokens and local state");
          TokenManager.clearTokens();
          set({
            user: null,
            isAuthenticated: false,
            isLoading: false,
            isLoggingOut: false,
            error: null,
          });
          toast.success(`Successfully signed out from ${brandConfig.name}`);
          console.log("Logout process completed");
        }
      },

      checkAuth: async () => {
        // Prevent multiple simultaneous auth checks
        const state = get();
        if (state.isLoading) return;

        try {
          set({ isLoading: true });

          // Check if we have stored user data
          const storedUser = TokenManager.getUser();
          const accessToken = TokenManager.getAccessToken();

          if (!accessToken || TokenManager.isTokenExpired(accessToken)) {
            // Try to refresh token
            const refreshToken = TokenManager.getRefreshToken();
            if (refreshToken) {
              try {
                const response = await AuthService.refreshToken(refreshToken);
                if (response.success && response.data) {
                  const { user, access_token, refresh_token, expires_in } =
                    response.data;
                  TokenManager.setTokens(
                    access_token,
                    refresh_token,
                    expires_in
                  );
                  TokenManager.setUser(user);

                  set({
                    user,
                    isAuthenticated: true,
                    isLoading: false,
                    error: null,
                  });
                  return;
                }
              } catch (error) {
                console.error("Token refresh failed:", error);
              }
            }

            // No valid tokens, clear everything
            TokenManager.clearTokens();
            set({
              user: null,
              isAuthenticated: false,
              isLoading: false,
              error: null,
            });
            return;
          }

          // If we have a stored user and valid token, use stored user first
          if (storedUser) {
            set({
              user: storedUser,
              isAuthenticated: true,
              isLoading: false,
              error: null,
            });

            // Optional: Refresh profile in background (don't wait for it)
            setTimeout(async () => {
              try {
                const { isAuthenticated, user } =
                  await AuthService.checkAuthStatus();
                if (isAuthenticated && user) {
                  TokenManager.setUser(user);
                  set({ user });
                }
              } catch (error) {
                console.log("Background profile refresh failed:", error);
              }
            }, 1000);
            return;
          }

          // Fallback: verify with server if no stored user but valid token
          try {
            const { isAuthenticated, user } =
              await AuthService.checkAuthStatus();

            if (isAuthenticated && user) {
              TokenManager.setUser(user);
              set({
                user,
                isAuthenticated: true,
                isLoading: false,
                error: null,
              });
            } else {
              throw new Error("Authentication verification failed");
            }
          } catch (error) {
            TokenManager.clearTokens();
            set({
              user: null,
              isAuthenticated: false,
              isLoading: false,
              error: null,
            });
          }
        } catch (error) {
          console.error("Auth check error:", error);
          set({
            user: null,
            isAuthenticated: false,
            isLoading: false,
            error: null,
          });
        }
      },

      refreshUserProfile: async () => {
        try {
          const response = await AuthService.getProfile();
          if (response.success && response.data) {
            const user = response.data;
            TokenManager.setUser(user);
            set({ user });
          }
        } catch (error) {
          console.error("Failed to refresh user profile:", error);
        }
      },

      clearError: () => set({ error: null }),

      setUser: (user: User | null) => {
        if (user) {
          TokenManager.setUser(user);
        }
        set({ user, isAuthenticated: !!user });
      },

      setLoading: (isLoading: boolean) => set({ isLoading }),
    }),
    {
      name: `${brandConfig.name.toLowerCase()}-auth`,
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
        // Don't persist loading states
      }),
    }
  )
);
