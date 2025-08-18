package common_controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthController handles health check endpoints
type HealthController struct{}

// CheckHealth handles GET /health
func (h *HealthController) CheckHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "Service is running",
	})
}
