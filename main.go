package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/routes"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/utils"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/utils/tokenstorage"
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
	utils.CheckEnvs()
	// Initialize token storage
	_, err = tokenstorage.GetTokenStorage().GetAccessToken()
	if err != nil {
		log.Fatalf("failed to initialize token storage: %v", err)
		return
	}
	// Initialize requirements manager
	requirementsManager := utils.GetRequirementsManager()
	if !requirementsManager.IsInitialized {
		log.Println("failed to initialize requirements manager")
		return
	}

	// streamlisteners.StartAlertStreamListener(os.Getenv("ALERT_STREAM_ENDPOINT"))
	// Initialize Gin router
	router := gin.Default()

	// Add middleware for logging and recovery
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Setup all routes
	routes.SetupRoutes(router)

	// Start server on port specified in environment variables
	router.Run(":" + os.Getenv("PORT"))
}
