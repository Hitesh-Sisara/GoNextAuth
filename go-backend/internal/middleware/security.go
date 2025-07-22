// File: internal/middleware/security.go

package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy (adjust based on your needs)
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' https://accounts.google.com https://apis.google.com; " +
			"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
			"font-src 'self' https://fonts.gstatic.com; " +
			"connect-src 'self' https://accounts.google.com https://www.googleapis.com; " +
			"img-src 'self' data: https:; " +
			"frame-src https://accounts.google.com;"

		c.Header("Content-Security-Policy", csp)

		c.Next()
	}
}
