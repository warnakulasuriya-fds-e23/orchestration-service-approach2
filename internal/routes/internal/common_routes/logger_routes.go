package common_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/controllers/common_controllers"
)

// SetupLoggerRoutes sets up health-related routes
func SetupLoggerRoutes(router *gin.Engine) {
	loggerController := common_controllers.LoggerController{}

	// Health check endpoint
	router.POST("/logger", loggerController.CheckHealth)
}
