package v1_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/controllers/v1_controllers"
)

func SetupEventReceiveRoutes(v1Group *gin.RouterGroup) {
	controller := v1_controllers.EventReceiveController{}
	v1Group.POST("/receive-face-match-event", controller.ReceiveFaceMatchEvent)
}
