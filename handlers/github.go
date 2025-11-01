package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/HiyawNT/GitGram/models"
    "github.com/HiyawNT/GitGram/services"
    "log"
    "net/http"
)

func GitHubWebhook(c *gin.Context) {
    var payload models.PushEvent

    if err := c.ShouldBindJSON(&payload); err != nil {
        log.Println("Error parsing payload:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Example: send each commit as message
    for _, commit := range payload.Commits {
        message := "📢 New commit in " + payload.Repository.FullName +
            "\n👤 By: " + payload.Pusher.Name +
            "\n📝 " + commit.Message +
            "\n🔗 " + commit.URL
        // Example chat ID (replace with dynamic storage later)
        services.SendMessage(123456789, message)
    }

    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
