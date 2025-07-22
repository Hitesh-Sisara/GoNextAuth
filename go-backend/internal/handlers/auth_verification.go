// File: internal/handlers/auth_verification.go

package handlers

import (
	"net/http"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/models"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/utils"

	"github.com/gin-gonic/gin"
)

// VerifyEmail godoc
// @Summary Verify email address
// @Description Verify user's email address using OTP
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.VerifyEmailRequest true "Email verification request"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req models.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SingleErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := h.authService.VerifyEmail(c.Request.Context(), req, c)
	if err != nil {
		utils.SingleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Email verified successfully", nil)
}

// ResendOTP godoc
// @Summary Resend OTP
// @Description Resend OTP for email verification or password reset
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.ResendOTPRequest true "Resend OTP request"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /auth/resend-otp [post]
func (h *AuthHandler) ResendOTP(c *gin.Context) {
	var req models.ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SingleErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := h.otpService.ResendOTP(c.Request.Context(), req.Email, req.OTPType)
	if err != nil {
		utils.SingleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "OTP sent successfully", nil)
}
