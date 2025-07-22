// File: cmd/server/main.go

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/config"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/database"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/handlers"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/routes"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/services"

	_ "github.com/Hitesh-Sisara/GoNextAuth/docs" // Import generated docs
)

// @title GoNextAuth API
// @version 1.0
// @description Complete authentication API with JWT tokens, email verification, Google OAuth, and multi-step authentication flow.
// @termsOfService https://GoNextAuth.com/terms

// @contact.name API Support
// @contact.url https://GoNextAuth.com/support
// @contact.email support@GoNextAuth.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	database.InitDB()
	defer database.CloseDB()

	// Run migrations
	database.RunMigrations()

	// Initialize services
	jwtService := services.NewJWTService(cfg)
	emailService := services.NewEmailService(cfg)
	otpService := services.NewOTPService(cfg, emailService)
	googleService := services.NewGoogleService(cfg)
	activityService := services.NewActivityService()
	authService := services.NewAuthService(cfg, jwtService, otpService, emailService, googleService, activityService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, otpService)
	healthHandler := handlers.NewHealthHandler(cfg)

	// Setup routes
	router := routes.SetupRoutes(cfg, authHandler, healthHandler, jwtService)

	// Create HTTP server
	server := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 %s server starting on port %s", cfg.Server.AppName, cfg.Server.Port)
		log.Printf("📖 API Documentation: http://localhost:%s/docs/index.html", cfg.Server.Port)
		log.Printf("🔗 Health Check: http://localhost:%s/health", cfg.Server.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Start background cleanup routine for expired OTPs
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			if err := otpService.CleanupExpiredOTPs(ctx); err != nil {
				log.Printf("Failed to cleanup expired OTPs: %v", err)
			}
			cancel()
		}
	}()

	// Start background cleanup routine for old activity logs (keep for 90 days)
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			if err := activityService.CleanupOldActivity(ctx, 90); err != nil {
				log.Printf("Failed to cleanup old activity logs: %v", err)
			}
			cancel()
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}
