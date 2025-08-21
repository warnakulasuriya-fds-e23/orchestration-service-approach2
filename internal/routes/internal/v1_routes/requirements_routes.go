package v1_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/controllers/v1_controllers"
)

func SetupRequirementsRoutes(v1Group *gin.RouterGroup) {
	controller := v1_controllers.RequirementsController{}
	v1Group.GET("/reload-requirements-file", controller.ReloadRequirementsFile)
	v1Group.GET("/configured-requirements-file-path", controller.GetConfiguredRequirementsFilePath)
}
