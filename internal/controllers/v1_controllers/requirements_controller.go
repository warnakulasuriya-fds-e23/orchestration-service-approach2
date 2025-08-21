package v1_controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/utils"
)

type RequirementsController struct{}

func (rc *RequirementsController) ReloadRequirementsFile(c *gin.Context) {
	if err := utils.GetRequirementsManager().ReloadRequirements(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully reloaded access requirements"})
}

func (rc *RequirementsController) GetConfiguredRequirementsFilePath(c *gin.Context) {
	filePath := utils.GetRequirementsManager().GetConfiguredRequirementsFilePath()
	c.JSON(http.StatusOK, gin.H{"configured_requirements_file_path": filePath})
}
