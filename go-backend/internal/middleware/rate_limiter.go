// File: internal/middleware/rate_limiter.go

package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/utils"

	"github.com/gin-gonic/gin"
)

// RateLimiter represents a rate limiter
type RateLimiter struct {
	requests map[string]*ClientData
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

// ClientData represents client request data
type ClientData struct {
	count     int
	lastReset time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*ClientData),
		limit:    limit,
		window:   window,
	}

	// Cleanup routine
	go rl.cleanup()

	return rl
}

// cleanup removes old entries
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		for ip, data := range rl.requests {
			if now.Sub(data.lastReset) > rl.window {
				delete(rl.requests, ip)
			}
		}
		rl.mutex.Unlock()
	}
}

// isAllowed checks if the request is allowed
func (rl *RateLimiter) isAllowed(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	if data, exists := rl.requests[clientIP]; exists {
		if now.Sub(data.lastReset) > rl.window {
			// Reset the window
			data.count = 1
			data.lastReset = now
			return true
		}

		if data.count >= rl.limit {
			return false
		}

		data.count++
		return true
	}

	// First request from this IP
	rl.requests[clientIP] = &ClientData{
		count:     1,
		lastReset: now,
	}
	return true
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(limit, window)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !limiter.isAllowed(clientIP) {
			utils.SingleErrorResponse(c, http.StatusTooManyRequests,
				fmt.Sprintf("Rate limit exceeded. Maximum %d requests per %v", limit, window))
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthRateLimitMiddleware creates a stricter rate limit for auth endpoints
func AuthRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(10, time.Minute) // 10 requests per minute for auth
}

// OTPRateLimitMiddleware creates a rate limit for OTP endpoints
func OTPRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(3, time.Minute) // 3 OTP requests per minute
}
