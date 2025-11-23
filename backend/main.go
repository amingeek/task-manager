// backend/main.go
package main

import (
	"log"
	"os"
	"task-manager/config"
	"task-manager/routes"
	"task-manager/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	config.ConnectDatabase()

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.LogInfo("Server starting", "port", port)

	if err := router.Run(":" + port); err != nil {
		utils.LogError("Failed to start server", err)
		log.Fatal("Failed to start server:", err)
	}
}
