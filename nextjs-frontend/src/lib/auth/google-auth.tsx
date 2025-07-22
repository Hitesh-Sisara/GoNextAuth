// File: src/lib/auth/google-auth.tsx (Updated with enhanced features)

"use client";

import { createContext, useContext, useEffect, useState } from "react";
import { toast } from "sonner";

declare global {
  interface Window {
    google: any;
  }
}

interface GoogleAuthContextType {
  isLoaded: boolean;
  signIn: () => Promise<string>;
  signOut: () => Promise<void>;
  isSignedIn: boolean;
  error: string | null;
  initiateServerFlow: () => Promise<string | null>;
}

const GoogleAuthContext = createContext<GoogleAuthContextType | null>(null);

interface GoogleAuthProviderProps {
  children: React.ReactNode;
  clientId: string;
}

export function GoogleAuthProvider({
  children,
  clientId,
}: GoogleAuthProviderProps) {
  const [isLoaded, setIsLoaded] = useState(false);
  const [isSignedIn, setIsSignedIn] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const initializeGoogleAuth = async () => {
      try {
        setError(null);

        // Load Google Identity Services script
        if (!window.google) {
          await loadGoogleScript();
        }

        // Initialize Google Identity Services
        window.google.accounts.id.initialize({
          client_id: clientId,
          callback: () => {}, // We'll handle this in the sign-in methods
          auto_select: false,
          cancel_on_tap_outside: true,
        });

        setIsLoaded(true);
      } catch (error) {
        console.error("Failed to initialize Google Auth:", error);
        setError("Failed to initialize Google authentication");
        setIsLoaded(true); // Set loaded even on error to prevent infinite loading
      }
    };

    if (clientId) {
      initializeGoogleAuth();
    } else {
      setError("Google Client ID not configured");
      setIsLoaded(true);
    }
  }, [clientId]);

  const loadGoogleScript = (): Promise<void> => {
    return new Promise((resolve, reject) => {
      if (window.google) {
        resolve();
        return;
      }

      const script = document.createElement("script");
      script.src = "https://accounts.google.com/gsi/client";
      script.async = true;
      script.defer = true;

      script.onload = () => {
        // Wait a bit for the API to be ready
        setTimeout(() => resolve(), 100);
      };

      script.onerror = () =>
        reject(new Error("Failed to load Google Identity Services script"));

      document.head.appendChild(script);
    });
  };

  // Client-side flow for getting access token
  const signIn = async (): Promise<string> => {
    return new Promise((resolve, reject) => {
      if (!window.google) {
        const error = new Error("Google Auth not initialized");
        setError(error.message);
        reject(error);
        return;
      }

      if (!clientId) {
        const error = new Error("Google Client ID not configured");
        setError(error.message);
        reject(error);
        return;
      }

      try {
        setError(null);

        // Use Google OAuth2 popup flow
        const client = window.google.accounts.oauth2.initTokenClient({
          client_id: clientId,
          scope: "openid email profile",
          callback: (response: any) => {
            if (response.error) {
              const errorMsg = `Google sign-in failed: ${response.error}`;
              setError(errorMsg);
              reject(new Error(errorMsg));
              return;
            }

            if (response.access_token) {
              setIsSignedIn(true);
              setError(null);
              resolve(response.access_token);
            } else {
              const errorMsg = "Failed to get access token from Google";
              setError(errorMsg);
              reject(new Error(errorMsg));
            }
          },
          error_callback: (error: any) => {
            let errorMsg = "Google sign-in failed";

            if (error.type === "popup_closed") {
              errorMsg = "Google sign-in was cancelled";
            } else if (error.type === "popup_blocked") {
              errorMsg = "Popup was blocked. Please allow popups and try again";
            } else if (error.message) {
              errorMsg = `Google sign-in failed: ${error.message}`;
            }

            setError(errorMsg);
            reject(new Error(errorMsg));
          },
        });

        // Request access token
        client.requestAccessToken({ prompt: "consent" });
      } catch (error: any) {
        console.error("Google sign-in failed:", error);
        const errorMsg = error.message || "Google sign-in failed";
        setError(errorMsg);
        reject(new Error(errorMsg));
      }
    });
  };

  // Server-side flow for better security
  const initiateServerFlow = async (): Promise<string | null> => {
    try {
      setError(null);

      // This would typically call your backend to get the Google OAuth URL
      // For now, we'll construct it here, but in production you should get it from your backend
      const scopes = encodeURIComponent("openid email profile");
      const redirectUri = encodeURIComponent(
        `${window.location.origin}/auth/google/callback`
      );
      const state = encodeURIComponent(
        Math.random().toString(36).substring(2, 15)
      );

      // Store state for validation (in a real app, this should be more secure)
      sessionStorage.setItem("google_oauth_state", state);

      const authUrl =
        `https://accounts.google.com/o/oauth2/v2/auth?` +
        `client_id=${clientId}&` +
        `redirect_uri=${redirectUri}&` +
        `scope=${scopes}&` +
        `response_type=code&` +
        `access_type=offline&` +
        `prompt=consent&` +
        `state=${state}`;

      return authUrl;
    } catch (error: any) {
      console.error("Failed to initiate server flow:", error);
      const errorMsg = error.message || "Failed to initiate Google sign-in";
      setError(errorMsg);
      return null;
    }
  };

  const signOut = async (): Promise<void> => {
    if (!window.google) {
      throw new Error("Google Auth not initialized");
    }

    try {
      setError(null);
      window.google.accounts.id.disableAutoSelect();
      setIsSignedIn(false);

      // Clear any stored state
      sessionStorage.removeItem("google_oauth_state");

      toast.success("Signed out from Google");
    } catch (error: any) {
      console.error("Google sign-out failed:", error);
      const errorMsg = error.message || "Google sign-out failed";
      setError(errorMsg);
      throw new Error(errorMsg);
    }
  };

  const value: GoogleAuthContextType = {
    isLoaded,
    signIn,
    signOut,
    isSignedIn,
    error,
    initiateServerFlow,
  };

  return (
    <GoogleAuthContext.Provider value={value}>
      {children}
    </GoogleAuthContext.Provider>
  );
}

export function useGoogleAuth() {
  const context = useContext(GoogleAuthContext);
  if (!context) {
    throw new Error("useGoogleAuth must be used within GoogleAuthProvider");
  }
  return context;
}

// Utility function to validate Google callback state
export function validateGoogleCallbackState(receivedState: string): boolean {
  const storedState = sessionStorage.getItem("google_oauth_state");
  if (!storedState || storedState !== receivedState) {
    console.error("Invalid or missing OAuth state parameter");
    return false;
  }

  // Clear the state after validation
  sessionStorage.removeItem("google_oauth_state");
  return true;
}

// Error messages for better UX
export const GOOGLE_AUTH_ERRORS = {
  NOT_INITIALIZED:
    "Google authentication is not initialized. Please try again.",
  NO_CLIENT_ID: "Google authentication is not properly configured.",
  POPUP_BLOCKED:
    "Popup was blocked. Please allow popups for this site and try again.",
  POPUP_CLOSED:
    "Google sign-in was cancelled. Please try again if you want to continue.",
  NETWORK_ERROR:
    "Network error occurred. Please check your connection and try again.",
  UNKNOWN_ERROR:
    "An unexpected error occurred during Google sign-in. Please try again.",
} as const;
