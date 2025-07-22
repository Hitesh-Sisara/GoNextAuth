// src/app/(auth)/layout.tsx

import { Toaster } from "@/components/ui/sonner";

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-gray-50">
      {/* Background Pattern */}
      <div className="absolute inset-0 bg-gradient-to-br from-blue-50 to-indigo-100 opacity-50" />

      {/* Content */}
      <div className="relative z-10">{children}</div>

      {/* Toast Notifications */}
      <Toaster position="top-right" />
    </div>
  );
}
