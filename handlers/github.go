package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/HiyawNT/GitGram/models"
	services "github.com/HiyawNT/GitGram/services"
	"github.com/HiyawNT/GitGram/storage"
	"github.com/gin-gonic/gin"
)

// GitHubWebhook handles incoming GitHub push event webhooks
func GitHubWebhook(c *gin.Context) {
	var pushEvent models.PushEvent

	// Bind JSON payload to PushEvent struct
	if err := c.ShouldBindJSON(&pushEvent); err != nil {
		log.Printf("Error binding JSON payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid payload format",
		})
		return
	}

	// Log the received push event
	log.Printf("Received push event for repository: %s", pushEvent.Repository.FullName)

	// Extract repository full name
	repository := pushEvent.Repository.FullName
	if repository == "" {
		log.Println("Error: Repository name is empty")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Repository name is missing",
		})
		return
	}

	// Get list of subscribers for this repository
	subscribers := storage.GetSubscribers(repository)
	if len(subscribers) == 0 {
		log.Printf("No subscribers found for repository: %s", repository)
		c.JSON(http.StatusOK, gin.H{
			"message": "No subscribers",
		})
		return
	}

	log.Printf("Found %d subscriber(s) for repository: %s", len(subscribers), repository)

	// Format the push event message
	message := formatPushEventMessage(pushEvent)

	// Send message to all subscribers
	successCount := 0
	for _, chatID := range subscribers {
		 err := services.SendMessage(chatID, message)
		if err != nil {
			log.Printf("Failed to send message to chat %d: %v", chatID, err)
		} else {
			successCount++
		}
	}

	log.Printf("Successfully sent notifications to %d/%d subscribers", successCount, len(subscribers))

	c.JSON(http.StatusOK, gin.H{
		"message":            "Webhook processed",
		"subscribers":        len(subscribers),
		"notifications_sent": successCount,
	})
}

// formatPushEventMessage creates a formatted message from a push event
func formatPushEventMessage(event models.PushEvent) string {
	var msg strings.Builder

	// Header with repository and branch
	branch := extractBranchName(event.Ref)
	msg.WriteString(fmt.Sprintf("🔔 *New Push to %s*\n", event.Repository.FullName))
	msg.WriteString(fmt.Sprintf("Branch: `%s`\n", branch))
	msg.WriteString(fmt.Sprintf("Pushed by: %s\n\n", event.Pusher.Name))

	// Commit information
	commitCount := len(event.Commits)
	if commitCount == 0 {
		msg.WriteString("No commits in this push.\n")
	} else {
		msg.WriteString(fmt.Sprintf("📝 *%d Commit", commitCount))
		if commitCount != 1 {
			msg.WriteString("s")
		}
		msg.WriteString(":*\n\n")

		// List commits (limit to first 5 to avoid message being too long)
		displayCount := commitCount
		if displayCount > 5 {
			displayCount = 5
		}

		for i := 0; i < displayCount; i++ {
			commit := event.Commits[i]
			shortID := commit.ID
			if len(shortID) > 7 {
				shortID = shortID[:7]
			}

			// Get first line of commit message
			commitMsg := strings.Split(commit.Message, "\n")[0]
			if len(commitMsg) > 100 {
				commitMsg = commitMsg[:97] + "..."
			}

			msg.WriteString(fmt.Sprintf("`%s` %s\n", shortID, commitMsg))

			// Add commit author if different from pusher
			if commit.Author.Name != "" && commit.Author.Name != event.Pusher.Name {
				msg.WriteString(fmt.Sprintf("   by %s\n", commit.Author.Name))
			}
		}

		if commitCount > 5 {
			msg.WriteString(fmt.Sprintf("\n...and %d more commit", commitCount-5))
			if commitCount-5 != 1 {
				msg.WriteString("s")
			}
			msg.WriteString("\n")
		}
	}

	// Add repository link
	msg.WriteString(fmt.Sprintf("\n[View Repository](%s)", event.Repository.HTMLURL))

	// Add compare link if we have commits
	if len(event.Commits) > 0 {
		// GitHub compare URL format: https://github.com/owner/repo/compare/before...after
		firstCommit := event.Commits[0].ID
		lastCommit := event.Commits[len(event.Commits)-1].ID
		if len(event.Commits) == 1 {
			msg.WriteString(fmt.Sprintf(" | [View Commit](%s)", event.Commits[0].URL))
		} else {
			compareURL := fmt.Sprintf("%s/compare/%s...%s",
				event.Repository.HTMLURL,
				firstCommit[:12],
				lastCommit[:12])
			msg.WriteString(fmt.Sprintf(" | [View Changes](%s)", compareURL))
		}
	}

	return msg.String()
}

// extractBranchName extracts the branch name from the ref
// Example: refs/heads/main -> main
func extractBranchName(ref string) string {
	parts := strings.Split(ref, "/")
	if len(parts) >= 3 && parts[0] == "refs" && parts[1] == "heads" {
		return strings.Join(parts[2:], "/")
	}
	return ref
}