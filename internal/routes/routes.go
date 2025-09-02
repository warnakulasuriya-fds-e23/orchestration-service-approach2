package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/routes/internal/common_routes"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/routes/internal/v1_routes"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine) {
	// Setup health routes
	common_routes.SetupHealthRoutes(router)
	common_routes.SetupLoggerRoutes(router)

	v1 := router.Group("/api/v1") // Create a base group for API v1
	{
		// Setup authorization routes
		v1_routes.SetupAuthorizationRoutes(v1.Group("/authorization"))
		v1_routes.SetupRequirementsRoutes(v1.Group("/requirements"))
	}
}
