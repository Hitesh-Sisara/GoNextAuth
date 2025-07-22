// File: internal/handlers/auth_oauth.go

package handlers

import (
	"net/http"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/models"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/utils"

	"github.com/gin-gonic/gin"
)

// GoogleAuth godoc
// @Summary Google OAuth authentication
// @Description Authenticate user with Google OAuth
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.GoogleAuthRequest true "Google auth request"
// @Success 200 {object} utils.APIResponse{data=models.AuthResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /auth/google [post]
func (h *AuthHandler) GoogleAuth(c *gin.Context) {
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
