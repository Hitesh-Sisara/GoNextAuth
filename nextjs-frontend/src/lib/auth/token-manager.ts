// File: src/lib/auth/token-manager.ts

import { authConfig } from "@/lib/config/app-config";
import Cookies from "js-cookie";

const { tokenKeys } = authConfig;

export class TokenManager {
  // Client-side token management
  static setTokens(
    accessToken: string,
    refreshToken: string,
    expiresIn: number
  ) {
    // Store in localStorage for client-side access
    if (typeof window !== "undefined") {
      localStorage.setItem(tokenKeys.accessToken, accessToken);
      localStorage.setItem(tokenKeys.refreshToken, refreshToken);
    }

    // Store in cookies for server-side access
    const expirationDate = new Date(Date.now() + expiresIn * 1000);

    Cookies.set(tokenKeys.accessToken, accessToken, {
      expires: expirationDate,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      path: "/",
    });

    Cookies.set(tokenKeys.refreshToken, refreshToken, {
      expires: 30, // 30 days
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      path: "/",
    });
  }

  static getAccessToken(): string | null {
    // Try localStorage first (client-side)
    if (typeof window !== "undefined") {
      const token = localStorage.getItem(tokenKeys.accessToken);
      if (token) return token;
    }

    // Fallback to cookies
    return Cookies.get(tokenKeys.accessToken) || null;
  }

  static getRefreshToken(): string | null {
    // Try localStorage first (client-side)
    if (typeof window !== "undefined") {
      const token = localStorage.getItem(tokenKeys.refreshToken);
      if (token) return token;
    }

    // Fallback to cookies
    return Cookies.get(tokenKeys.refreshToken) || null;
  }

  static clearTokens() {
    // Clear localStorage
    if (typeof window !== "undefined") {
      localStorage.removeItem(tokenKeys.accessToken);
      localStorage.removeItem(tokenKeys.refreshToken);
      localStorage.removeItem(tokenKeys.user);
    }

    // Clear cookies
    Cookies.remove(tokenKeys.accessToken, { path: "/" });
    Cookies.remove(tokenKeys.refreshToken, { path: "/" });
    Cookies.remove(tokenKeys.user, { path: "/" });
  }

  static setUser(user: any) {
    try {
      // Create a minimal user object to avoid cookie size issues
      const minimalUser = {
        id: user.id,
        email: user.email,
        first_name: user.first_name,
        last_name: user.last_name,
        is_email_verified: user.is_email_verified,
        auth_provider: user.auth_provider,
      };

      const userStr = JSON.stringify(minimalUser);

      // Check if the JSON string is too large for cookies (4KB limit)
      if (userStr.length > 3000) {
        console.warn(
          "User data too large for cookie storage, using localStorage only"
        );

        // Store only in localStorage if too large
        if (typeof window !== "undefined") {
          localStorage.setItem(tokenKeys.user, JSON.stringify(user));
        }
        return;
      }

      // Store in localStorage
      if (typeof window !== "undefined") {
        localStorage.setItem(tokenKeys.user, JSON.stringify(user));
      }

      // Store minimal user in cookies for SSR
      Cookies.set(tokenKeys.user, userStr, {
        expires: 30, // 30 days
        secure: process.env.NODE_ENV === "production",
        sameSite: "lax",
        path: "/",
      });
    } catch (error) {
      console.error("Failed to store user data:", error);

      // Fallback: clear any corrupted data and store only in localStorage
      Cookies.remove(tokenKeys.user, { path: "/" });
      if (typeof window !== "undefined") {
        localStorage.setItem(tokenKeys.user, JSON.stringify(user));
      }
    }
  }

  static getUser() {
    // Try localStorage first
    if (typeof window !== "undefined") {
      const userStr = localStorage.getItem(tokenKeys.user);
      if (userStr) {
        try {
          return JSON.parse(userStr);
        } catch (error) {
          console.error("Failed to parse user from localStorage:", error);
          localStorage.removeItem(tokenKeys.user);
        }
      }
    }

    // Fallback to cookies
    const userStr = Cookies.get(tokenKeys.user);
    if (userStr) {
      try {
        return JSON.parse(userStr);
      } catch (error) {
        console.error("Failed to parse user from cookie:", error);
        Cookies.remove(tokenKeys.user, { path: "/" });
        return null;
      }
    }

    return null;
  }

  static isTokenExpired(token: string): boolean {
    if (!token) return true;

    try {
      const payload = JSON.parse(atob(token.split(".")[1]));
      const currentTime = Date.now() / 1000;
      return payload.exp < currentTime;
    } catch {
      return true;
    }
  }

  static getTokenFromHeaders(headers: Headers): string | null {
    const authHeader = headers.get("authorization");
    if (authHeader && authHeader.startsWith("Bearer ")) {
      return authHeader.substring(7);
    }
    return null;
  }
}

// Server-side token utilities with better error handling
export const getServerSideTokens = (cookies: string) => {
  try {
    const cookieArray = cookies.split("; ");
    const cookieObj: Record<string, string> = {};

    cookieArray.forEach((cookie) => {
      const [key, value] = cookie.split("=");
      if (key && value) {
        try {
          cookieObj[key] = decodeURIComponent(value);
        } catch (error) {
          console.error(`Failed to decode cookie ${key}:`, error);
        }
      }
    });

    let user = null;
    if (cookieObj[tokenKeys.user]) {
      try {
        user = JSON.parse(cookieObj[tokenKeys.user]);
      } catch (error) {
        console.error("Failed to parse user from server-side cookie:", error);
        // Don't throw, just return null user
      }
    }

    return {
      accessToken: cookieObj[tokenKeys.accessToken] || null,
      refreshToken: cookieObj[tokenKeys.refreshToken] || null,
      user,
    };
  } catch (error) {
    console.error("Failed to parse server-side tokens:", error);
    return {
      accessToken: null,
      refreshToken: null,
      user: null,
    };
  }
};
