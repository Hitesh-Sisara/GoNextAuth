package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/config"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/database"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/utils"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	config *config.Config
}

func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{config: cfg}
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
}

// Health godoc
// @Summary Health check
// @Description Get application health status
// @Tags Health
// @Produce json
// @Success 200 {object} utils.APIResponse{data=HealthResponse}
// @Failure 503 {object} utils.APIResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	services := make(map[string]string)

	// Check database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := database.GetDB()
	if err := db.Ping(ctx); err != nil {
		services["database"] = "unhealthy"
		utils.ErrorResponseJSON(c, http.StatusServiceUnavailable, "Service unhealthy", HealthResponse{
			Status:    "unhealthy",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Version:   h.config.Server.APIVersion,
			Services:  services,
		})
		return
	}
	services["database"] = "healthy"

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   h.config.Server.APIVersion,
		Services:  services,
	}

	utils.SuccessResponse(c, http.StatusOK, "Service is healthy", response)
}

// Ready godoc
// @Summary Readiness check
// @Description Check if application is ready to serve requests
// @Tags Health
// @Produce json
// @Success 200 {object} utils.APIResponse
// @Failure 503 {object} utils.APIResponse
// @Router /ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	// Check database connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	db := database.GetDB()
	if err := db.Ping(ctx); err != nil {
		utils.SingleErrorResponse(c, http.StatusServiceUnavailable, "Service not ready")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Service is ready", nil)
}
