package main

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/routes"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/utils"
)

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := io.ReadAll(tee)
		c.Request.Body = io.NopCloser(&buf)
		log.Println("Request middleware log start")
		log.Println(string(body))
		log.Println(c.Request.Header)
		log.Println("Request middleware log end")
		c.Next()
	}
}

func main() {
	_, err := os.Stat(".env")
	if err == nil {
		log.Println("discovered .env file")
		err := godotenv.Load()
		if err != nil {
			log.Println("however failed to load .env file")
		} else {
			log.Println(".env successfully loaded")
		}
	}
	// Initialize requirements manager
	requirementsManager := utils.GetRequirementsManager()
	if !requirementsManager.IsInitialized {
		log.Println("failed to initialize requirements manager")
		return
	}

	// Initialize Gin router
	router := gin.Default()

	// Add middleware for logging and recovery
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(RequestLoggerMiddleware())

	// Setup all routes
	routes.SetupRoutes(router)

	// Start server on port 8080
	router.Run(":9090")
}
