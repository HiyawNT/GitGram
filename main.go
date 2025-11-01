package main

import (
    "github.com/gin-gonic/gin"
    "github.com/HiyawNT/GitGram/handlers"
    "github.com/HiyawNT/GitGram/utils"
    "log"
    "os"
)

func main() {
    // Load .env
    utils.LoadEnv()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    r := gin.Default()
    
    // GitHub Webhook endpoint
    r.POST("/webhook", handlers.GitHubWebhook)

    log.Println("Server started at port " + port)
    r.Run(":" + port)
}
