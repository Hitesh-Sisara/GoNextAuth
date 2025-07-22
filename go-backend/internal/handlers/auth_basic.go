// File: internal/handlers/auth_basic.go

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/models"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/services"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
	otpService  *services.OTPService
}

func NewAuthHandler(authService *services.AuthService, otpService *services.OTPService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		otpService:  otpService,
	}
}

// Signup godoc (legacy)
// @Summary User registration (legacy)
// @Description Register a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.SignupRequest true "Signup request"
// @Success 201 {object} utils.APIResponse{data=models.UserResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 409 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var req models.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SingleErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.authService.Signup(c.Request.Context(), req, c)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			utils.SingleErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.SingleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Account created successfully. Please check your email for verification.", user.ToUserResponse())
}

// Login godoc
// @Summary User login with password
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} utils.APIResponse{data=models.AuthResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SingleErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	authResponse, err := h.authService.Login(c.Request.Context(), req, c)
	if err != nil {
		if err.Error() == "invalid email or password" || err.Error() == "account is deactivated" {
			utils.UnauthorizedResponse(c, err.Error())
			return
		}
		utils.SingleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", authResponse)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} utils.APIResponse{data=models.AuthResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SingleErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	authResponse, err := h.authService.RefreshToken(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "invalid refresh token" || err.Error() == "refresh token is invalid or expired" {
			utils.UnauthorizedResponse(c, err.Error())
			return
		}
		utils.SingleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token refreshed successfully", authResponse)
}

// Logout godoc
// @Summary User logout
// @Description Revoke refresh token and logout user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Logout request"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SingleErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate refresh token format
	if req.RefreshToken == "" {
		utils.SingleErrorResponse(c, http.StatusBadRequest, "Refresh token is required")
		return
	}

	fmt.Printf("Logout request received for token: %s...\n", req.RefreshToken[:20])

	// Create context with timeout for better performance
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Try to get user ID for activity logging (with safety checks)
	var userID int
	if h.authService != nil {
		if jwtService := h.authService.GetJWTService(); jwtService != nil {
			if claims, err := jwtService.ValidateToken(req.RefreshToken, "refresh"); err == nil {
				userID = claims.UserID
				fmt.Printf("Logout for user ID: %d\n", userID)
			} else {
				fmt.Printf("Could not extract user ID from token: %v\n", err)
			}
		}
	}

	// Perform logout (revoke refresh token)
	err := h.authService.Logout(ctx, req.RefreshToken, c)
	if err != nil {
		fmt.Printf("Logout service error: %v (continuing with success response)\n", err)
		// Don't return error to client - logout should always succeed from client perspective
	} else {
		fmt.Printf("Logout service completed successfully\n")
	}

	// Log activity in background (non-blocking) with safety checks
	if userID > 0 && h.authService != nil {
		if activityService := h.authService.GetActivityService(); activityService != nil {
			go func() {
				bgCtx, bgCancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer bgCancel()

				if err := activityService.LogActivity(bgCtx, userID, models.ActivityLogout, c, nil); err != nil {
					fmt.Printf("Failed to log logout activity: %v\n", err)
				} else {
					fmt.Printf("Logout activity logged successfully\n")
				}
			}()
		}
	}

	// Always return success for logout
	fmt.Printf("Logout completed, sending success response\n")
	utils.SuccessResponse(c, http.StatusOK, "Logged out successfully", nil)
}
