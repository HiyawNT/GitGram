package handlers

import (
	"github.com/HiyawNT/GitGram/models"
	"github.com/HiyawNT/GitGram/services"
	"github.com/HiyawNT/GitGram/storage"
	"github.com/gin-gonic/gin"
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

	repoFullName := payload.Repository.FullName

	// Get all chat IDs subscribed to this repo
	chatIDs, err := subscription.GetSubscriptionByRepo(repoFullName)
	if err != nil {
		log.Println("Error fetching subscriptions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscriptions"})
		return
	}

	if len(chatIDs) == 0 {
		log.Printf("No subscriptions found for repo: %s", repoFullName)
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "No subscribers"})
		return
	}

	// Send notification to all subscribed chats
	for _, commit := range payload.Commits {
		message := "📢 New commit in " + payload.Repository.FullName +
			"\n👤 By: " + payload.Pusher.Name +
			"\n📝 " + commit.Message +
			"\n🔗 " + commit.URL

		for _, chatID := range chatIDs {
			services.SendMessage(chatID, message)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "subscribers": len(chatIDs)})
}
