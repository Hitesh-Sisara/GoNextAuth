// File: internal/handlers/auth_multistep.go

package handlers

import (
	"net/http"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/models"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/utils"

	"github.com/gin-gonic/gin"
)

// InitiateSignup godoc
// @Summary Initiate signup process
// @Description Start signup process by sending OTP to email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.EmailLoginRequest true "Email request"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 409 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /auth/signup/initiate [post]
func (h *AuthHandler) InitiateSignup(c *gin.Context) {
	var req models.EmailLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SingleErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := h.authService.InitiateEmailSignup(c.Request.Context(), req.Email)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			utils.SingleErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.SingleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "OTP sent to your email address", nil)
}

// CompleteSignup godoc
// @Summary Complete signup process
// @Description Complete signup after email verification
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.CompleteSignupRequest true "Complete signup request"
// @Success 201 {object} utils.APIResponse{data=models.UserResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 409 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /auth/signup/complete [post]
func (h *AuthHandler) CompleteSignup(c *gin.Context) {
	var req models.CompleteSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SingleErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.authService.CompleteSignup(c.Request.Context(), req, c)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			utils.SingleErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.SingleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Account created successfully", user.ToUserResponse())
}

// InitiateEmailLogin godoc
// @Summary Initiate OTP-based login
// @Description Start OTP-based login process
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.EmailLoginRequest true "Email login request"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /auth/login/email [post]
func (h *AuthHandler) InitiateEmailLogin(c *gin.Context) {
	var req models.EmailLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SingleErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Fixed: Pass the email string instead of the entire request object
	err := h.authService.InitiateEmailLogin(c.Request.Context(), models.EmailLoginRequest{Email: req.Email})
	if err != nil {
		utils.SingleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "OTP sent to your email address", nil)
}

// CompleteOTPLogin godoc
// @Summary Complete OTP-based login
// @Description Complete login using OTP
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.OTPLoginRequest true "OTP login request"
// @Success 200 {object} utils.APIResponse{data=models.AuthResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /auth/login/otp [post]
func (h *AuthHandler) CompleteOTPLogin(c *gin.Context) {
	var req models.OTPLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SingleErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	authResponse, err := h.authService.CompleteOTPLogin(c.Request.Context(), req, c)
	if err != nil {
		utils.SingleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", authResponse)
}
