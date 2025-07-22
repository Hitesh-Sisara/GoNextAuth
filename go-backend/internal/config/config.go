// File: internal/config/config.go

package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	AWS      AWSConfig
	Google   GoogleConfig
	OTP      OTPConfig
	CORS     CORSConfig
	Branding BrandingConfig
}

type ServerConfig struct {
	Port       string
	GinMode    string
	AppName    string
	APIVersion string
}

type DatabaseConfig struct {
	URL string
}

type JWTConfig struct {
	Secret        string
	AccessExpiry  int // minutes
	RefreshExpiry int // days
}

type AWSConfig struct {
	SES SESConfig
}

type SESConfig struct {
	Region       string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	SMTPHost     string
	SMTPPort     int
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type OTPConfig struct {
	Expiry int // minutes
	Length int
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

type BrandingConfig struct {
	Name           string
	PrimaryColor   string
	SecondaryColor string
	LogoURL        string
	SupportEmail   string
	Website        string
}

var AppConfig *Config

func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	smtpPort, _ := strconv.Atoi(getEnv("AWS_SES_SMTP_PORT", "587"))
	accessExpiry, _ := strconv.Atoi(getEnv("JWT_ACCESS_EXPIRY", "1440"))
	refreshExpiry, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY", "30"))
	otpExpiry, _ := strconv.Atoi(getEnv("OTP_EXPIRY", "10"))
	otpLength, _ := strconv.Atoi(getEnv("OTP_LENGTH", "6"))

	config := &Config{
		Server: ServerConfig{
			Port:       getEnv("PORT", "8080"),
			GinMode:    getEnv("GIN_MODE", "debug"),
			AppName:    getEnv("APP_NAME", "GoNextAuth"),
			APIVersion: getEnv("API_VERSION", "v1"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", ""),
			AccessExpiry:  accessExpiry,
			RefreshExpiry: refreshExpiry,
		},
		AWS: AWSConfig{
			SES: SESConfig{
				Region:       getEnv("AWS_SES_REGION", "ap-south-1"),
				SMTPUsername: getEnv("AWS_SES_SMTP_USERNAME", ""),
				SMTPPassword: getEnv("AWS_SES_SMTP_PASSWORD", ""),
				FromEmail:    getEnv("AWS_SES_FROM_EMAIL", ""),
				SMTPHost:     getEnv("AWS_SES_SMTP_HOST", ""),
				SMTPPort:     smtpPort,
			},
		},
		Google: GoogleConfig{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
		},
		OTP: OTPConfig{
			Expiry: otpExpiry,
			Length: otpLength,
		},
		CORS: CORSConfig{
			AllowedOrigins: getCORSOrigins(),
			AllowedMethods: getCORSMethods(),
			AllowedHeaders: getCORSHeaders(),
		},
		Branding: BrandingConfig{
			Name:           getEnv("BRAND_NAME", "GoNextAuth"),
			PrimaryColor:   getEnv("BRAND_PRIMARY_COLOR", "#6366f1"),
			SecondaryColor: getEnv("BRAND_SECONDARY_COLOR", "#8b5cf6"),
			LogoURL:        getEnv("BRAND_LOGO_URL", ""),
			SupportEmail:   getEnv("BRAND_SUPPORT_EMAIL", "support@GoNextAuth.com"),
			Website:        getEnv("BRAND_WEBSITE", "https://GoNextAuth.com"),
		},
	}

	// Validate required fields
	if config.Database.URL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if config.JWT.Secret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	if config.AWS.SES.SMTPUsername == "" || config.AWS.SES.SMTPPassword == "" {
		log.Fatal("AWS SES SMTP credentials are required")
	}
	if config.Google.ClientID == "" || config.Google.ClientSecret == "" {
		log.Println("Warning: Google OAuth credentials not configured")
	}

	AppConfig = config
	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getCORSOrigins returns CORS origins with smart defaults
func getCORSOrigins() []string {
	originsEnv := getEnv("CORS_ALLOWED_ORIGINS", "")

	if originsEnv != "" {
		return strings.Split(originsEnv, ",")
	}

	// Default permissive origins for development
	return []string{
		"*", // Allow all origins - you can restrict this for production
		"http://localhost:3000",
		"http://localhost:3001",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:3001",
		"https://localhost:3000",
		"https://accounts.google.com",
	}
}

// getCORSMethods returns all HTTP methods
func getCORSMethods() []string {
	methodsEnv := getEnv("CORS_ALLOWED_METHODS", "")

	if methodsEnv != "" {
		return strings.Split(methodsEnv, ",")
	}

	// Allow all common HTTP methods
	return []string{
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
		"OPTIONS",
		"HEAD",
	}
}

// getCORSHeaders returns permissive headers
func getCORSHeaders() []string {
	headersEnv := getEnv("CORS_ALLOWED_HEADERS", "")

	if headersEnv != "" {
		return strings.Split(headersEnv, ",")
	}

	// Allow all common headers
	return []string{
		"*", // Allow all headers
		"Accept",
		"Authorization",
		"Cache-Control",
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"X-Requested-With",
		"DNT",
		"If-Modified-Since",
		"Keep-Alive",
		"Origin",
		"User-Agent",
		"Sec-Fetch-Mode",
		"Sec-Fetch-Site",
		"Sec-Fetch-Dest",
	}
}
