// File: internal/middleware/cors.go

package middleware

import (
	"net/http"
	"strings"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/config"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware configures permissive CORS settings for development
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Always set CORS headers for all requests

		// 1. Handle Origin - Allow configured origins or any localhost for development
		if origin != "" {
			allowedOrigins := cfg.CORS.AllowedOrigins
			isAllowed := false

			// Check if origin is in allowed list
			for _, allowedOrigin := range allowedOrigins {
				if allowedOrigin == "*" || origin == allowedOrigin {
					isAllowed = true
					break
				}
			}

			// For development: also allow any localhost/127.0.0.1 with any port
			if !isAllowed && (strings.HasPrefix(origin, "http://localhost:") ||
				strings.HasPrefix(origin, "http://127.0.0.1:") ||
				strings.HasPrefix(origin, "https://localhost:") ||
				strings.HasPrefix(origin, "https://127.0.0.1:")) {
				isAllowed = true
			}

			if isAllowed {
				c.Header("Access-Control-Allow-Origin", origin)
			}
		} else {
			// If no origin header, allow all (for tools like Postman)
			c.Header("Access-Control-Allow-Origin", "*")
		}

		// 2. Allow all methods
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS,HEAD")

		// 3. Allow all headers - reflect what client requested or use wildcard
		requestedHeaders := c.GetHeader("Access-Control-Request-Headers")
		if requestedHeaders != "" {
			c.Header("Access-Control-Allow-Headers", requestedHeaders)
		} else {
			// Allow common headers if no specific headers requested
			c.Header("Access-Control-Allow-Headers", "*")
		}

		// 4. Allow credentials
		c.Header("Access-Control-Allow-Credentials", "true")

		// 5. Expose headers that client can access
		c.Header("Access-Control-Expose-Headers", "Content-Length,Content-Type,Authorization,X-Requested-With")

		// 6. Cache preflight for 24 hours
		c.Header("Access-Control-Max-Age", "86400")

		// 7. Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			// Log preflight request for debugging
			if cfg.Server.GinMode == "debug" {
				c.Header("X-Debug-Preflight", "handled")
			}

			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// SimpleCORSMiddleware - Ultra permissive CORS for development (alternative)
func SimpleCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set permissive CORS headers for all requests
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
