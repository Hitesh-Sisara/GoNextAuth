// File: internal/middleware/deduplication.go

package middleware

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/utils"

	"github.com/gin-gonic/gin"
)

// RequestTracker tracks ongoing requests to prevent duplicates
type RequestTracker struct {
	mu       sync.RWMutex
	requests map[string]*RequestInfo
}

type RequestInfo struct {
	InProgress bool
	Timestamp  time.Time
}

var tracker = &RequestTracker{
	requests: make(map[string]*RequestInfo),
}

// DeduplicationMiddleware prevents duplicate requests for specific endpoints
func DeduplicationMiddleware() gin.HandlerFunc {
	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			tracker.cleanup()
		}
	}()

	return func(c *gin.Context) {
		// Only apply to logout and callback endpoints
		path := c.Request.URL.Path
		if path != "/api/v1/auth/logout" && path != "/api/v1/auth/google/callback" {
			c.Next()
			return
		}

		// Create request signature based on client IP and path
		signature := createRequestSignature(c)

		tracker.mu.Lock()
		if info, exists := tracker.requests[signature]; exists && info.InProgress {
			tracker.mu.Unlock()
			fmt.Printf("Duplicate request detected for %s from %s\n", path, c.ClientIP())
			utils.SingleErrorResponse(c, http.StatusTooManyRequests, "Request already in progress, please wait")
			c.Abort()
			return
		}

		// Mark request as in progress
		tracker.requests[signature] = &RequestInfo{
			InProgress: true,
			Timestamp:  time.Now(),
		}
		tracker.mu.Unlock()

		fmt.Printf("Processing %s request from %s\n", path, c.ClientIP())

		// Process the request
		c.Next()

		// Mark request as completed
		tracker.mu.Lock()
		if info, exists := tracker.requests[signature]; exists {
			info.InProgress = false
			info.Timestamp = time.Now()
		}
		tracker.mu.Unlock()

		fmt.Printf("Completed %s request from %s\n", path, c.ClientIP())
	}
}

func createRequestSignature(c *gin.Context) string {
	// Create signature based on client IP and path
	// For logout, we want to prevent multiple simultaneous logouts from same client
	data := fmt.Sprintf("%s:%s", c.ClientIP(), c.Request.URL.Path)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

func (rt *RequestTracker) cleanup() {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	cutoff := time.Now().Add(-10 * time.Minute)
	cleaned := 0

	for key, info := range rt.requests {
		if info.Timestamp.Before(cutoff) {
			delete(rt.requests, key)
			cleaned++
		}
	}

	if cleaned > 0 {
		fmt.Printf("Cleaned up %d old request entries\n", cleaned)
	}
}
