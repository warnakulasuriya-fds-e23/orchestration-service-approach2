package v1_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/controllers/v1_controllers"
)

func SetupAuthorizationRoutes(v1Group *gin.RouterGroup) {
	controller := v1_controllers.AuthorizationController{}
	v1Group.POST("/authorize-for-door-access", controller.AuthorizeUserForDoorAccess)
}
