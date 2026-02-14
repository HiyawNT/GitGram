package main

import (
	"log"
	"os"

	"github.com/HiyawNT/GitGram/handlers"
	"github.com/HiyawNT/GitGram/services"
	"github.com/HiyawNT/GitGram/storage"
	"github.com/HiyawNT/GitGram/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load .env
	utils.LoadEnv()
	// Init Db
	subscription.InitDB()

	// Init Telegram Bot
	services.InitBot()

	// Start listening for Telegram Updates in GORoutine
	go handlers.HandleTelegramUpdates()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	// GitHub Webhook endpoint
	r.POST("/webhook", handlers.GitHubWebhook)

	r.GET("/health", func(ctx *gin.Context) {

		ctx.JSON(200, gin.H{"status": "ﮩ٨ـﮩﮩ٨ـ♥️ﮩ٨ـﮩﮩ٨ـOK"})
	})
	log.Println("Server started at port " + port)
	r.Run(":" + port)
}
