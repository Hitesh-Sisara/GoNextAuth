// File: internal/routes/routes.go (Updated with enhanced security)

package routes

import (
	"net/http"
	"time"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/config"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/handlers"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/middleware"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(
	cfg *config.Config,
	authHandler *handlers.AuthHandler,
	healthHandler *handlers.HealthHandler,
	jwtService *services.JWTService,
) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Create router
	router := gin.New()

	// Global middleware - ORDER MATTERS!
	router.Use(middleware.LoggerMiddleware())
	router.Use(gin.Recovery())

	// Security middleware
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.SanitizeInput())
	router.Use(middleware.RequestSizeMiddleware(10 << 20)) // 10MB limit

	// Choose CORS middleware based on environment
	if cfg.Server.GinMode == "debug" {
		router.Use(middleware.SimpleCORSMiddleware())
	} else {
		router.Use(middleware.CORSMiddleware(cfg))
	}

	// Global rate limiting (more permissive)
	router.Use(middleware.RateLimitMiddleware(100, time.Minute))

	// Universal OPTIONS handler
	router.OPTIONS("/*path", handlePreflight)

	// Health check routes (no auth required)
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)

	// Swagger documentation
	if cfg.Server.GinMode != "release" {
		router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// API routes
	api := router.Group("/api/" + cfg.Server.APIVersion)
	{
		// Authentication routes (with stricter rate limiting)
		auth := api.Group("/auth")
		auth.Use(middleware.AuthRateLimitMiddleware())
		auth.Use(middleware.ContentTypeMiddleware("application/json"))
		{
			// Multi-step authentication flow
			signup := auth.Group("/signup")
			{
				signup.POST("/initiate",
					middleware.OTPRateLimitMiddleware(),
					authHandler.InitiateSignup)
				signup.POST("/complete", authHandler.CompleteSignup)
			}

			login := auth.Group("/login")
			{
				login.POST("/email",
					middleware.OTPRateLimitMiddleware(),
					authHandler.InitiateEmailLogin)
				login.POST("/otp", authHandler.CompleteOTPLogin)
			}

			// Google OAuth routes with enhanced security
			google := auth.Group("/google")
			{
				google.GET("/url", authHandler.GoogleAuthURL)
				google.GET("/callback", authHandler.GoogleCallback)
				google.POST("/token", authHandler.GoogleAuthToken)
				if cfg.Server.GinMode == "debug" {
					google.GET("/debug", authHandler.DebugGoogleConfig)
				}
			}

			// Password management with stricter rate limiting
			auth.POST("/forgot-password",
				middleware.OTPRateLimitMiddleware(),
				authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)

			// OTP management
			auth.POST("/resend-otp",
				middleware.OTPRateLimitMiddleware(),
				authHandler.ResendOTP)

			// Legacy and additional auth endpoints
			auth.POST("/signup", authHandler.Signup)
			auth.POST("/login", authHandler.Login)
			auth.POST("/verify-email", authHandler.VerifyEmail)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)

			// Protected routes (auth required)
			protected := auth.Group("")
			protected.Use(middleware.AuthMiddleware(jwtService))
			{
				protected.GET("/profile", authHandler.GetProfile)
			}
		}

		// Admin routes (future expansion)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(jwtService))
		// admin.Use(middleware.AdminRoleMiddleware()) // TODO: Implement role-based access
		{
			// Admin endpoints can be added here
		}
	}

	return router
}

// handlePreflight handles all OPTIONS/preflight requests
func handlePreflight(c *gin.Context) {
	// Additional debug info for preflight requests
	if gin.Mode() == gin.DebugMode {
		c.Header("X-Preflight-Path", c.Request.URL.Path)
		c.Header("X-Preflight-Method", c.GetHeader("Access-Control-Request-Method"))
		c.Header("X-Preflight-Headers", c.GetHeader("Access-Control-Request-Headers"))
	}

	c.Status(http.StatusNoContent)
}
