// File: src/app/layout.tsx

import { AuthProvider } from "@/components/auth/auth-provider";
import { ErrorBoundary } from "@/components/common/error-boundary";
import { Toaster } from "@/components/ui/sonner";
import { GoogleAuthProvider } from "@/lib/auth/google-auth";
import {
  brandConfig,
  seoConfig,
  socialConfig,
  urlConfig,
} from "@/lib/config/app-config";
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: seoConfig.title,
  description: seoConfig.description,
  keywords: seoConfig.keywords,
  authors: [{ name: `${brandConfig.name} Team` }],
  creator: brandConfig.name,
  openGraph: {
    type: "website",
    locale: "en_US",
    url: urlConfig.base,
    title: seoConfig.title,
    description: brandConfig.description,
    siteName: brandConfig.name,
  },
  twitter: {
    card: "summary_large_image",
    title: seoConfig.title,
    description: brandConfig.description,
    creator: socialConfig.twitter,
  },
  robots: {
    index: true,
    follow: true,
  },
};

// Get Google Client ID from environment variables
const GOOGLE_CLIENT_ID = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID || "";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <ErrorBoundary>
          <GoogleAuthProvider clientId={GOOGLE_CLIENT_ID}>
            <AuthProvider>{children}</AuthProvider>
          </GoogleAuthProvider>
          <Toaster position="top-right" />
        </ErrorBoundary>
      </body>
    </html>
  );
}
