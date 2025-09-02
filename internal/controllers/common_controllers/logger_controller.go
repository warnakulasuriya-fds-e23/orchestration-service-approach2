package common_controllers

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

// LoggerController handles health check endpoints
type LoggerController struct{}

// CheckHealth handles GET /health
func (h *LoggerController) CheckHealth(c *gin.Context) {
	reqBodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error while reading request body": err.Error()})
		return
	}
	reqBody := gjson.ParseBytes(reqBodyBytes)
	log.Println("requsetBody :")
	log.Println(reqBody)
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "Service is running",
	})
}
