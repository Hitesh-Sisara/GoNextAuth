// File: src/lib/config/app-config.ts

/**
 * Application configuration
 * Centralizes all app-wide settings and constants
 */

interface AppConfig {
  // Brand and Identity
  brandName: string;
  tagline: string;
  description: string;

  // URLs and Domains
  baseUrl: string;
  domain: string;

  // Social Media
  twitterHandle: string;

  // SEO and Metadata
  keywords: string[];

  // Auth Configuration
  tokenKeys: {
    accessToken: string;
    refreshToken: string;
    user: string;
  };

  // Other settings
  defaultRedirect: string;
}

// Default configuration values
const defaultConfig: AppConfig = {
  brandName: "GoNextAuth",
  tagline: "Your Time Management Solution",
  description:
    "Streamline your productivity with our ultimate time management and task organization platform.",
  baseUrl: "https://GoNextAuth.com",
  domain: "GoNextAuth.com",
  twitterHandle: "@GoNextAuth",
  keywords: [
    "time management",
    "productivity",
    "task organization",
    "scheduling",
  ],
  tokenKeys: {
    accessToken: "GoNextAuth_access_token",
    refreshToken: "GoNextAuth_refresh_token",
    user: "GoNextAuth_user",
  },
  defaultRedirect: "/dashboard",
};

// Create the configuration object with environment variable overrides
export const appConfig: AppConfig = {
  brandName: process.env.NEXT_PUBLIC_BRAND_NAME || defaultConfig.brandName,
  tagline: process.env.NEXT_PUBLIC_BRAND_TAGLINE || defaultConfig.tagline,
  description:
    process.env.NEXT_PUBLIC_BRAND_DESCRIPTION || defaultConfig.description,
  baseUrl: process.env.NEXT_PUBLIC_BASE_URL || defaultConfig.baseUrl,
  domain: process.env.NEXT_PUBLIC_DOMAIN || defaultConfig.domain,
  twitterHandle:
    process.env.NEXT_PUBLIC_TWITTER_HANDLE || defaultConfig.twitterHandle,
  keywords: process.env.NEXT_PUBLIC_SEO_KEYWORDS
    ? process.env.NEXT_PUBLIC_SEO_KEYWORDS.split(",").map((k) => k.trim())
    : defaultConfig.keywords,
  tokenKeys: {
    accessToken: process.env.NEXT_PUBLIC_TOKEN_PREFIX
      ? `${process.env.NEXT_PUBLIC_TOKEN_PREFIX}_access_token`
      : defaultConfig.tokenKeys.accessToken,
    refreshToken: process.env.NEXT_PUBLIC_TOKEN_PREFIX
      ? `${process.env.NEXT_PUBLIC_TOKEN_PREFIX}_refresh_token`
      : defaultConfig.tokenKeys.refreshToken,
    user: process.env.NEXT_PUBLIC_TOKEN_PREFIX
      ? `${process.env.NEXT_PUBLIC_TOKEN_PREFIX}_user`
      : defaultConfig.tokenKeys.user,
  },
  defaultRedirect:
    process.env.NEXT_PUBLIC_DEFAULT_REDIRECT || defaultConfig.defaultRedirect,
};

// Helper functions for commonly used values
export const getBrandName = () => appConfig.brandName;
export const getFullBrandName = () =>
  `${appConfig.brandName} - ${appConfig.tagline}`;
export const getBrandDescription = () => appConfig.description;

// Export individual config sections for convenience
export const brandConfig = {
  name: appConfig.brandName,
  tagline: appConfig.tagline,
  description: appConfig.description,
  fullName: getFullBrandName(),
};

export const urlConfig = {
  base: appConfig.baseUrl,
  domain: appConfig.domain,
};

export const socialConfig = {
  twitter: appConfig.twitterHandle,
};

export const seoConfig = {
  keywords: appConfig.keywords,
  title: getFullBrandName(),
  description: appConfig.description,
};

export const authConfig = {
  tokenKeys: appConfig.tokenKeys,
  defaultRedirect: appConfig.defaultRedirect,
};
