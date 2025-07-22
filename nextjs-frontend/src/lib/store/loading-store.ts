// File: src/lib/store/loading-store.ts

import { create } from "zustand";
import { devtools } from "zustand/middleware";

interface LoadingState {
  // Global loading states
  isLoading: boolean;
  loadingMessage: string;

  // Specific loading states
  loadingStates: Record<string, boolean>;

  // Actions
  setLoading: (loading: boolean, message?: string) => void;
  setSpecificLoading: (key: string, loading: boolean) => void;
  isSpecificLoading: (key: string) => boolean;
  clearAllLoading: () => void;
}

export const useLoadingStore = create<LoadingState>()(
  devtools(
    (set, get) => ({
      isLoading: false,
      loadingMessage: "",
      loadingStates: {},

      setLoading: (loading: boolean, message: string = "") =>
        set(
          {
            isLoading: loading,
            loadingMessage: message,
          },
          false,
          "setLoading"
        ),

      setSpecificLoading: (key: string, loading: boolean) =>
        set(
          (state) => ({
            loadingStates: {
              ...state.loadingStates,
              [key]: loading,
            },
          }),
          false,
          `setSpecificLoading/${key}`
        ),

      isSpecificLoading: (key: string) => {
        const state = get();
        return state.loadingStates[key] || false;
      },

      clearAllLoading: () =>
        set(
          {
            isLoading: false,
            loadingMessage: "",
            loadingStates: {},
          },
          false,
          "clearAllLoading"
        ),
    }),
    {
      name: "loading-store",
    }
  )
);

// Loading keys for specific operations
export const LoadingKeys = {
  // Auth operations
  AUTH_LOGIN: "auth.login",
  AUTH_SIGNUP: "auth.signup",
  AUTH_LOGOUT: "auth.logout",
  AUTH_REFRESH: "auth.refresh",
  AUTH_GOOGLE: "auth.google",

  // OTP operations
  OTP_SEND: "otp.send",
  OTP_VERIFY: "otp.verify",
  OTP_RESEND: "otp.resend",

  // Password operations
  PASSWORD_FORGOT: "password.forgot",
  PASSWORD_RESET: "password.reset",

  // Profile operations
  PROFILE_FETCH: "profile.fetch",
  PROFILE_UPDATE: "profile.update",

  // Email operations
  EMAIL_VERIFY: "email.verify",
  EMAIL_RESEND: "email.resend",
} as const;

export type LoadingKey = (typeof LoadingKeys)[keyof typeof LoadingKeys];

// Hook for easier usage
export const useLoading = (key?: LoadingKey) => {
  const {
    isLoading,
    loadingMessage,
    setLoading,
    setSpecificLoading,
    isSpecificLoading,
  } = useLoadingStore();

  if (key) {
    return {
      isLoading: isSpecificLoading(key),
      setLoading: (loading: boolean) => setSpecificLoading(key, loading),
    };
  }

  return {
    isLoading,
    loadingMessage,
    setLoading,
  };
};
