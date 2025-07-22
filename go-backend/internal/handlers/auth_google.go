// File: internal/handlers/auth_google.go

package handlers

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/models"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/utils"

	"github.com/gin-gonic/gin"
)

// GoogleAuthURL godoc
// @Summary Get Google OAuth URL
// @Description Get the Google OAuth authorization URL for frontend redirection
// @Tags Authentication
// @Produce json
// @Success 200 {object} utils.APIResponse{data=map[string]string}
// @Router /auth/google/url [get]
func (h *AuthHandler) GoogleAuthURL(c *gin.Context) {
	// Generate a secure state parameter with HMAC
	state, err := h.generateSecureState()
	if err != nil {
		fmt.Printf("Failed to generate state: %v\n", err)
		utils.SingleErrorResponse(c, http.StatusInternalServerError, "Failed to generate state parameter")
		return
	}

	fmt.Printf("Generated state: %s\n", state)

	baseURL := "https://accounts.google.com/o/oauth2/auth"
	params := url.Values{}
	params.Add("client_id", h.authService.GetGoogleClientID())
	params.Add("redirect_uri", h.authService.GetGoogleRedirectURL())
	params.Add("scope", "openid email profile")
	params.Add("response_type", "code")
	params.Add("access_type", "offline")
	params.Add("prompt", "consent")
	params.Add("state", state)

	authURL := baseURL + "?" + params.Encode()

	utils.SuccessResponse(c, http.StatusOK, "Google OAuth URL generated", map[string]string{
		"auth_url": authURL,
		"state":    state,
	})
}

// GoogleCallback godoc
// @Summary Handle Google OAuth callback
// @Description Handle the callback from Google OAuth and authenticate user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param code query string true "Authorization code from Google"
// @Param state query string true "State parameter for CSRF protection"
// @Success 200 {object} utils.APIResponse{data=models.AuthResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	error_param := c.Query("error")

	// Log received parameters for debugging
	fmt.Printf("Google callback received:\n")
	fmt.Printf("  - Code: %t\n", code != "")
	fmt.Printf("  - State: %s\n", state)
	fmt.Printf("  - Error: %s\n", error_param)

	// Handle OAuth errors
	if error_param != "" {
		errorMsg := "Google OAuth error"
		if error_param == "access_denied" {
			errorMsg = "Google sign-in was cancelled"
		}
		utils.SingleErrorResponse(c, http.StatusBadRequest, errorMsg)
		return
	}

	if code == "" {
		utils.SingleErrorResponse(c, http.StatusBadRequest, "Authorization code is required")
		return
	}

	if state == "" {
		utils.SingleErrorResponse(c, http.StatusBadRequest, "State parameter is required")
		return
	}

	// Validate state parameter using HMAC (stateless)
	if !h.validateState(state) {
		fmt.Printf("State validation failed for state: %s\n", state)
		utils.SingleErrorResponse(c, http.StatusBadRequest, "Invalid state parameter - possible CSRF attack")
		return
	}

	fmt.Printf("State validation successful, proceeding with Google auth\n")

	// Process the Google authentication
	authResponse, err := h.authService.GoogleCallbackAuth(c.Request.Context(), code, c)
	if err != nil {
		fmt.Printf("Google callback auth failed: %v\n", err)

		// Check if the error is related to invalid_grant (used authorization code)
		if strings.Contains(err.Error(), "invalid_grant") {
			utils.SingleErrorResponse(c, http.StatusBadRequest, "The authorization code has already been used or has expired. Please try signing in again.")
			return
		}

		// Handle other Google API errors
		if strings.Contains(err.Error(), "google token exchange failed") {
			utils.SingleErrorResponse(c, http.StatusBadRequest, "Google authentication failed. Please try again.")
			return
		}

		utils.SingleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Printf("Google authentication successful for user: %s\n", authResponse.User.Email)
	utils.SuccessResponse(c, http.StatusOK, "Google authentication successful", authResponse)
}

// GoogleAuthToken godoc
// @Summary Authenticate with Google access token
// @Description Authenticate user using Google access token (for frontend flows)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.GoogleAuthRequest true "Google auth request"
// @Success 200 {object} utils.APIResponse{data=models.AuthResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /auth/google/token [post]
func (h *AuthHandler) GoogleAuthToken(c *gin.Context) {
	var req models.GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SingleErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	authResponse, err := h.authService.GoogleAuth(c.Request.Context(), req, c)
	if err != nil {
		utils.SingleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Google authentication successful", authResponse)
}

// generateSecureState generates a HMAC-signed state parameter with timestamp
func (h *AuthHandler) generateSecureState() (string, error) {
	// Get JWT secret for HMAC signing
	secret := h.authService.GetJWTSecret()

	// Generate random bytes
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// Create timestamp (valid for 10 minutes)
	timestamp := time.Now().Unix()

	// Create payload: timestamp + random
	payload := fmt.Sprintf("%d.%s", timestamp, base64.URLEncoding.EncodeToString(randomBytes))

	// Create HMAC signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	signature := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	// Combine payload and signature
	state := fmt.Sprintf("%s.%s", payload, signature)

	fmt.Printf("Generated state components:\n")
	fmt.Printf("  - Timestamp: %d\n", timestamp)
	fmt.Printf("  - Payload: %s\n", payload)
	fmt.Printf("  - Final state: %s\n", state)

	return state, nil
}

// validateState validates the HMAC-signed state parameter
func (h *AuthHandler) validateState(state string) bool {
	fmt.Printf("Validating state: %s\n", state)

	// Get JWT secret for HMAC validation
	secret := h.authService.GetJWTSecret()

	// Split state into payload and signature
	parts := strings.Split(state, ".")
	if len(parts) != 3 {
		fmt.Printf("Invalid state format: expected 3 parts, got %d\n", len(parts))
		return false
	}

	payload := fmt.Sprintf("%s.%s", parts[0], parts[1])
	receivedSignature := parts[2]

	fmt.Printf("State parts:\n")
	fmt.Printf("  - Timestamp: %s\n", parts[0])
	fmt.Printf("  - Random: %s\n", parts[1])
	fmt.Printf("  - Signature: %s\n", receivedSignature)

	// Verify HMAC signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expectedSignature := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(receivedSignature), []byte(expectedSignature)) {
		fmt.Printf("State signature validation failed\n")
		fmt.Printf("  - Expected: %s\n", expectedSignature)
		fmt.Printf("  - Received: %s\n", receivedSignature)
		return false
	}

	// Check timestamp (10 minutes validity)
	timestamp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		fmt.Printf("Invalid timestamp in state: %v\n", err)
		return false
	}

	currentTime := time.Now().Unix()
	if currentTime-timestamp > 600 { // 10 minutes
		fmt.Printf("State expired: timestamp %d, now %d, diff %d seconds\n",
			timestamp, currentTime, currentTime-timestamp)
		return false
	}

	fmt.Printf("State validation successful\n")
	return true
}

// DebugGoogleConfig godoc
// @Summary Debug Google OAuth configuration
// @Description Show Google OAuth configuration for debugging
// @Tags Authentication
// @Produce json
// @Success 200 {object} utils.APIResponse{data=map[string]string}
// @Router /auth/google/debug [get]
func (h *AuthHandler) DebugGoogleConfig(c *gin.Context) {
	config := map[string]string{
		"client_id":    h.authService.GetGoogleClientID(),
		"redirect_uri": h.authService.GetGoogleRedirectURL(),
	}

	utils.SuccessResponse(c, http.StatusOK, "Google OAuth configuration", config)
}
