// File: src/app/auth/reset-password/page.tsx

"use client";

import { Loader2 } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect } from "react";

export default function ResetPasswordPage() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    // Redirect to the new forgot password flow
    const email = searchParams.get("email");
    if (email) {
      router.replace(
        `/auth/forgot-password?email=${encodeURIComponent(email)}`
      );
    } else {
      router.replace("/auth/forgot-password");
    }
  }, [router, searchParams]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center">
        <Loader2 className="h-8 w-8 animate-spin mx-auto mb-4" />
        <p className="text-gray-600">Redirecting...</p>
      </div>
    </div>
  );
}
