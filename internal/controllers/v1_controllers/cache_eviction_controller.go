package v1_controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/utils/authorizationscache"
)

type CacheEvictionController struct{}

func (cec *CacheEvictionController) EvictUserFromCache(c *gin.Context) {
	userName := c.Param("userName")
	if userName == "" {
		c.JSON(400, gin.H{"error": "User name is required"})
		return
	}

	cachedAuth := authorizationscache.GetAuthorizationsCacheInstance()

	message := cachedAuth.EvictFromCache(userName)
	if message == "No authorized doors found for user" {
		c.JSON(404, gin.H{"error": message})
		return
	}

	c.JSON(200, gin.H{"message": message})
}
