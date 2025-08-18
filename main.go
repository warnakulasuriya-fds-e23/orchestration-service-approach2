package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/routes"
)

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
	// Initialize Gin router
	router := gin.Default()

	// Add middleware for logging and recovery
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Setup all routes
	routes.SetupRoutes(router)

	// Start server on port 8080
	router.Run(":5000")
}
