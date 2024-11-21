package main

import (
	"mini-telegram/config"
	"mini-telegram/models"
	"mini-telegram/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin mode and initialize logger
	gin.SetMode(gin.ReleaseMode)
	config.LoadEnv()
	config.InitializeLogger()

	// Connect to the database
	config.ConnectDB()

	// Auto-migrate database schema asynchronously
	go func() {
		config.Logger.Info("Starting database migration...")
		if err := config.DB.AutoMigrate(&models.User{}, &models.LogEntry{}); err != nil {
			config.Logger.Fatalf("Migration failed: %v", err)
		}
		config.Logger.Info("Database migration completed.")
	}()

	// Start the Gin server
	port := config.GetEnv("APP_PORT", "8080")
	config.Logger.Infof("Starting server on port %s", port)
	router := routes.SetupRoutes()

	// Run server concurrently
	go func() {
		if err := router.Run(":" + port); err != nil {
			config.Logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Keep the main routine running
	select {}
}
