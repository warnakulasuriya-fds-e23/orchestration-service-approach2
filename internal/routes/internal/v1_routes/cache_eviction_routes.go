package v1_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/controllers/v1_controllers"
)

func SetupCacheEvictionRoutes(v1Group *gin.RouterGroup) {
	controller := v1_controllers.CacheEvictionController{}
	v1Group.DELETE("/evict-user-from-cache/:userName", controller.EvictUserFromCache)
}
