package common_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/controllers/common_controllers"
)

// SetupHealthRoutes sets up health-related routes
func SetupHealthRoutes(router *gin.Engine) {
	healthController := common_controllers.HealthController{}

	// Health check endpoint
	router.GET("/health", healthController.CheckHealth)
}
