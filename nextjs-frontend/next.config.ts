/** @type {import('next').NextConfig} */
const nextConfig = {
  // Disable React StrictMode in development to prevent duplicate effect calls
  // You can re-enable it once the Google auth issue is resolved
  reactStrictMode: process.env.NODE_ENV === "production",

  // Additional configurations
  experimental: {
    // Ensure middleware runs properly
    middlewareHost: true,
  },

  // Logging for debugging
  logging: {
    fetches: {
      fullUrl: process.env.NODE_ENV === "development",
    },
  },
};

module.exports = nextConfig;
