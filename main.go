package main

import (
	"log"

	"github.com/HiyawNT/GitGram/handlers"
	"github.com/HiyawNT/GitGram/utils"
	"github.com/HiyawNT/GitGram/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables from .env file
	utils.LoadEnv()

	// Initialize Telegram bot
	err := services.InitBot()
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}

	// Set Gin to release mode in production
	// gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "GitGram bot is running",
			"status":  "ok",
		})
	})

	// GitHub webhook endpoint
	router.POST("/webhook", handlers.GitHubWebhook)

	// Get port from environment or use default
	port := utils.GetEnv("PORT", "8080")

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", port)
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("GitGram bot is running. Press Ctrl+C to stop.")

	// Keep the main function running to receive Telegram updates
	select {}
}