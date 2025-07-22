import { getServerSideTokens } from "@/lib/auth/token-manager";
import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";

const publicPaths = [
  "/auth/login",
  "/auth/signup",
  "/auth/forgot-password",
  "/auth/reset-password",
  "/auth/verify-email",
  "/auth/google/callback", // Make sure this is included
  "/",
  "/health",
];

const protectedPaths = ["/dashboard", "/profile", "/settings"];

// Rate limiting store (simple in-memory for middleware)
const rateLimitStore = new Map<string, { count: number; resetTime: number }>();

// Security headers
const securityHeaders = {
  "X-Frame-Options": "DENY",
  "X-Content-Type-Options": "nosniff",
  "X-XSS-Protection": "1; mode=block",
  "Referrer-Policy": "strict-origin-when-cross-origin",
  "Permissions-Policy": "camera=(), microphone=(), geolocation=()",
};

// Rate limiting function
function isRateLimited(ip: string, endpoint: string): boolean {
  const key = `${ip}:${endpoint}`;
  const now = Date.now();
  const limit = rateLimitStore.get(key);

  if (!limit || now > limit.resetTime) {
    rateLimitStore.set(key, { count: 1, resetTime: now + 60000 }); // 1 minute window
    return false;
  }

  if (limit.count >= 10) {
    // 10 requests per minute
    return true;
  }

  limit.count++;
  return false;
}

// Clean up old rate limit entries
function cleanupRateLimit() {
  const now = Date.now();
  for (const [key, value] of rateLimitStore.entries()) {
    if (now > value.resetTime) {
      rateLimitStore.delete(key);
    }
  }
}

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const ip = request.ip || request.headers.get("x-forwarded-for") || "unknown";

  // Clean up rate limit store periodically
  if (Math.random() < 0.01) {
    // 1% chance to clean up
    cleanupRateLimit();
  }

  // Rate limiting for auth endpoints (but not Google callback)
  if (pathname.startsWith("/auth/") && pathname !== "/auth/google/callback") {
    if (isRateLimited(ip, "auth")) {
      return new NextResponse("Rate limit exceeded", {
        status: 429,
        headers: securityHeaders,
      });
    }
  }

  // Security headers for all responses
  const response = NextResponse.next();
  Object.entries(securityHeaders).forEach(([key, value]) => {
    response.headers.set(key, value);
  });

  // IMPORTANT: Allow Google callback without any authentication checks
  if (pathname === "/auth/google/callback") {
    console.log("Google callback detected, allowing without auth check");
    return response;
  }

  // Health check endpoint
  if (pathname === "/health") {
    return response;
  }

  try {
    const cookies = request.headers.get("cookie") || "";
    const { accessToken, user } = getServerSideTokens(cookies);

    const isAuthenticated = !!(accessToken && user);
    const isPublicPath = publicPaths.some((path) => pathname.startsWith(path));
    const isProtectedPath = protectedPaths.some((path) =>
      pathname.startsWith(path)
    );

    // Log suspicious activity
    if (!isAuthenticated && isProtectedPath) {
      console.warn(`Unauthorized access attempt to ${pathname} from ${ip}`);
    }

    // If user is authenticated and trying to access auth pages, redirect to dashboard
    if (
      isAuthenticated &&
      pathname.startsWith("/auth/") &&
      pathname !== "/auth/google/callback" // Don't redirect callback
    ) {
      return NextResponse.redirect(new URL("/dashboard", request.url));
    }

    // If user is not authenticated and trying to access protected pages, redirect to login
    if (!isAuthenticated && isProtectedPath) {
      const loginUrl = new URL("/auth/login", request.url);
      loginUrl.searchParams.set("redirect", pathname);
      return NextResponse.redirect(loginUrl);
    }

    // For root path, redirect based on authentication status
    if (pathname === "/") {
      if (isAuthenticated) {
        return NextResponse.redirect(new URL("/dashboard", request.url));
      } else {
        return NextResponse.redirect(new URL("/auth/login", request.url));
      }
    }

    return response;
  } catch (error) {
    console.error("Middleware error:", error);

    // Log security incidents
    console.warn(`Middleware error for ${pathname} from ${ip}:`, error);

    // On error, allow public paths and redirect protected paths to login
    const isPublicPath = publicPaths.some((path) => pathname.startsWith(path));
    const isProtectedPath = protectedPaths.some((path) =>
      pathname.startsWith(path)
    );

    if (isProtectedPath) {
      const loginUrl = new URL("/auth/login", request.url);
      loginUrl.searchParams.set("redirect", pathname);
      return NextResponse.redirect(loginUrl);
    }

    if (pathname === "/") {
      return NextResponse.redirect(new URL("/auth/login", request.url));
    }

    return response;
  }
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - public assets
     */
    "/((?!api|_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)",
  ],
};
