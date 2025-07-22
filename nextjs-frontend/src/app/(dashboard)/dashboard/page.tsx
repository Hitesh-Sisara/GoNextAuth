// File: src/app/(dashboard)/dashboard/page.tsx

"use client";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { brandConfig } from "@/lib/config/app-config";
import { useAuth } from "@/lib/hooks/use-auth";
import { formatDateTime, getInitials } from "@/lib/utils";
import {
  Activity,
  Calendar,
  Loader2,
  LogOut,
  Mail,
  MapPin,
  Monitor,
  Phone,
  Shield,
  User,
} from "lucide-react";
import { useRouter } from "next/navigation";
import { useRef } from "react";
import { toast } from "sonner";

export default function DashboardPage() {
  const { user, logout, isLoading } = useAuth();
  const router = useRouter();
  const isLoggingOutRef = useRef(false);

  const handleLogout = async () => {
    // Prevent multiple logout attempts
    if (isLoggingOutRef.current || isLoading) {
      console.log("Logout already in progress, ignoring click");
      return;
    }

    try {
      console.log("Logout button clicked");
      isLoggingOutRef.current = true;

      await logout();

      console.log("Logout completed, redirecting to login");
      router.push("/auth/login");
    } catch (error) {
      console.error("Logout failed:", error);
      toast.error("Logout failed. Please try again.");
    } finally {
      isLoggingOutRef.current = false;
    }
  };

  if (!user) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  const getProviderBadge = (provider: string) => {
    switch (provider) {
      case "google":
        return <Badge variant="secondary">Google</Badge>;
      case "email":
        return <Badge variant="outline">Email</Badge>;
      default:
        return <Badge variant="outline">{provider}</Badge>;
    }
  };

  const isLoggingOut = isLoggingOutRef.current || isLoading;

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <h1 className="text-2xl font-bold text-gray-900">
                {brandConfig.name}
              </h1>
            </div>
            <div className="flex items-center space-x-4">
              <Avatar>
                <AvatarImage src={user.avatar_url} alt={user.first_name} />
                <AvatarFallback>
                  {getInitials(user.first_name, user.last_name)}
                </AvatarFallback>
              </Avatar>
              <div className="hidden md:block">
                <p className="text-sm font-medium text-gray-900">
                  {user.first_name} {user.last_name}
                </p>
                <p className="text-xs text-gray-500">{user.email}</p>
              </div>
              <Button
                variant="outline"
                size="sm"
                onClick={handleLogout}
                disabled={isLoggingOut}
              >
                {isLoggingOut ? (
                  <>
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                    Signing out...
                  </>
                ) : (
                  <>
                    <LogOut className="h-4 w-4 mr-2" />
                    Logout
                  </>
                )}
              </Button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h2 className="text-3xl font-bold text-gray-900">
            Welcome back, {user.first_name}!
          </h2>
          <p className="text-gray-600 mt-2">
            Here&apos;s an overview of your account and recent activity.
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Profile Card */}
          <Card className="lg:col-span-2">
            <CardHeader>
              <CardTitle className="flex items-center">
                <User className="h-5 w-5 mr-2" />
                Profile Information
              </CardTitle>
              <CardDescription>
                Your account details and preferences
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="flex items-center space-x-4">
                <Avatar className="h-16 w-16">
                  <AvatarImage src={user.avatar_url} alt={user.first_name} />
                  <AvatarFallback className="text-lg">
                    {getInitials(user.first_name, user.last_name)}
                  </AvatarFallback>
                </Avatar>
                <div>
                  <h3 className="text-lg font-semibold">
                    {user.first_name} {user.last_name}
                  </h3>
                  <p className="text-gray-600">{user.email}</p>
                  <div className="flex items-center space-x-2 mt-1">
                    {getProviderBadge(user.auth_provider)}
                    {user.is_email_verified ? (
                      <Badge
                        variant="default"
                        className="bg-green-100 text-green-800"
                      >
                        <Shield className="h-3 w-3 mr-1" />
                        Verified
                      </Badge>
                    ) : (
                      <Badge variant="destructive">Unverified</Badge>
                    )}
                  </div>
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="flex items-center space-x-3 p-3 bg-gray-50 rounded-lg">
                  <Mail className="h-5 w-5 text-gray-400" />
                  <div>
                    <p className="text-sm font-medium text-gray-900">Email</p>
                    <p className="text-sm text-gray-600">{user.email}</p>
                  </div>
                </div>

                {user.phone && (
                  <div className="flex items-center space-x-3 p-3 bg-gray-50 rounded-lg">
                    <Phone className="h-5 w-5 text-gray-400" />
                    <div>
                      <p className="text-sm font-medium text-gray-900">Phone</p>
                      <p className="text-sm text-gray-600">{user.phone}</p>
                    </div>
                  </div>
                )}

                <div className="flex items-center space-x-3 p-3 bg-gray-50 rounded-lg">
                  <Calendar className="h-5 w-5 text-gray-400" />
                  <div>
                    <p className="text-sm font-medium text-gray-900">
                      Member Since
                    </p>
                    <p className="text-sm text-gray-600">
                      {formatDateTime(user.created_at)}
                    </p>
                  </div>
                </div>

                <div className="flex items-center space-x-3 p-3 bg-gray-50 rounded-lg">
                  <Activity className="h-5 w-5 text-gray-400" />
                  <div>
                    <p className="text-sm font-medium text-gray-900">
                      Last Active
                    </p>
                    <p className="text-sm text-gray-600">
                      {formatDateTime(user.last_activity_at)}
                    </p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Quick Actions */}
          <Card>
            <CardHeader>
              <CardTitle>Quick Actions</CardTitle>
              <CardDescription>Manage your account settings</CardDescription>
            </CardHeader>
            <CardContent className="space-y-3">
              <Button className="w-full justify-start" variant="outline">
                <User className="h-4 w-4 mr-2" />
                Edit Profile
              </Button>
              <Button className="w-full justify-start" variant="outline">
                <Shield className="h-4 w-4 mr-2" />
                Security Settings
              </Button>
              <Button className="w-full justify-start" variant="outline">
                <Activity className="h-4 w-4 mr-2" />
                Activity Log
              </Button>
              <Button className="w-full justify-start" variant="outline">
                <Mail className="h-4 w-4 mr-2" />
                Email Preferences
              </Button>
            </CardContent>
          </Card>
        </div>

        {/* Recent Activity */}
        <Card className="mt-6">
          <CardHeader>
            <CardTitle className="flex items-center">
              <Activity className="h-5 w-5 mr-2" />
              Recent Activity
            </CardTitle>
            <CardDescription>
              Your recent account activity and sign-ins
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {/* Mock activity items - you would fetch these from your API */}
              <div className="flex items-center space-x-3 p-3 border rounded-lg">
                <div className="h-2 w-2 bg-green-500 rounded-full"></div>
                <Monitor className="h-4 w-4 text-gray-400" />
                <div className="flex-1">
                  <p className="text-sm font-medium">Signed in</p>
                  <p className="text-xs text-gray-500">
                    {formatDateTime(user.last_activity_at)} • Web Browser
                  </p>
                </div>
                <MapPin className="h-4 w-4 text-gray-400" />
              </div>

              <div className="flex items-center space-x-3 p-3 border rounded-lg">
                <div className="h-2 w-2 bg-blue-500 rounded-full"></div>
                <Shield className="h-4 w-4 text-gray-400" />
                <div className="flex-1">
                  <p className="text-sm font-medium">Email verified</p>
                  <p className="text-xs text-gray-500">
                    {formatDateTime(user.created_at)} • Account Creation
                  </p>
                </div>
              </div>

              <div className="flex items-center space-x-3 p-3 border rounded-lg">
                <div className="h-2 w-2 bg-purple-500 rounded-full"></div>
                <User className="h-4 w-4 text-gray-400" />
                <div className="flex-1">
                  <p className="text-sm font-medium">Account created</p>
                  <p className="text-xs text-gray-500">
                    {formatDateTime(user.created_at)} • Welcome to{" "}
                    {brandConfig.name}!
                  </p>
                </div>
              </div>
            </div>

            <div className="mt-4 pt-4 border-t">
              <Button variant="outline" className="w-full">
                View All Activity
              </Button>
            </div>
          </CardContent>
        </Card>
      </main>
    </div>
  );
}
