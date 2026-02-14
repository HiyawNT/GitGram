package handlers

import (
	"fmt"
	"github.com/HiyawNT/GitGram/services"
	"github.com/HiyawNT/GitGram/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

func HandleTelegramUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := services.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID

		// Handle commands
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				handleStart(chatID)
			case "subscribe":
				handleSubscribe(chatID, update.Message.CommandArguments())
			case "unsubscribe":
				handleUnsubscribe(chatID, update.Message.CommandArguments())
			case "list_subscriptions":
				handleList(chatID)
			case "help":
				handleHelp(chatID)
			default:
				services.SendMessage(chatID, "Unknown command. Use /help to see available commands.")
			}
		}
	}
}
func handleStart(chatID int64) {
	message := "👋 Welcome to GitGram!\n\n" +
		"I am your Trusty Bot That will notify you about GitHub repository events.\n\n" +
		"Use /help to see available commands."
	services.SendMessage(chatID, message)
}

func handleSubscribe(chatID int64, args string) {
	args = strings.TrimSpace(args)
	if args == "" {
		services.SendMessage(chatID, " Usage: /subscribe <owner/repo>\n\nExample: /subscribe octocat/Hello-World")
		return
	}

	// Validate repo format (basic validation)
	if !strings.Contains(args, "/") {
		services.SendMessage(chatID, " Invalid repository format. Use: owner/repo\n\nExample: /subscribe octocat/Hello-World")
		return
	}

	repoFullName := args

	// Check if already subscribed
	exists, err := subscription.CheckSubscriptionExists(chatID, repoFullName)
	if err != nil {
		log.Println("Error checking subscription:", err)
		services.SendMessage(chatID, " An error occurred. Please try again.")
		return
	}

	if exists {
		services.SendMessage(chatID, fmt.Sprintf("ℹ️ You're already subscribed to %s", repoFullName))
		return
	}

	// Add subscription
	err = subscription.AddSubscription(chatID, repoFullName)
	if err != nil {
		log.Println("Error adding subscription:", err)
		services.SendMessage(chatID, " Failed to add subscription. Please try again.")
		return
	}

	message := fmt.Sprintf(" Successfully subscribed to %s\n\n"+
		"You'll receive notifications when there are pushes to this repository.\n\n"+
		"🔗 Make sure to add a webhook in your GitHub repository:\n"+
		"Repository → Settings → Webhooks → Add webhook\n"+
		"Payload URL: <your-server-url>/webhook", repoFullName)
	services.SendMessage(chatID, message)
}

func handleUnsubscribe(chatID int64, args string) {
	args = strings.TrimSpace(args)
	if args == "" {
		services.SendMessage(chatID, " Usage: /unsubscribe <owner/repo>\n\nExample: /unsubscribe octocat/Hello-World")
		return
	}

	repoFullName := args

	// Check if subscription exists
	exists, err := subscription.CheckSubscriptionExists(chatID, repoFullName)
	if err != nil {
		log.Println("Error checking subscription:", err)
		services.SendMessage(chatID, " An error occurred. Please try again.")
		return
	}

	if !exists {
		services.SendMessage(chatID, fmt.Sprintf("ℹ️ You're not subscribed to %s", repoFullName))
		return
	}

	// Remove subscription
	err = subscription.RemoveSubscription(chatID, repoFullName)
	if err != nil {
		log.Println("Error removing subscription:", err)
		services.SendMessage(chatID, " Failed to remove subscription. Please try again.")
		return
	}

	services.SendMessage(chatID, fmt.Sprintf(" Successfully unsubscribed from %s", repoFullName))
}

func handleList(chatID int64) {
	subscriptions, err := subscription.GetSubscriptionsByChatID(chatID)
	if err != nil {
		log.Println("Error fetching subscriptions:", err)
		services.SendMessage(chatID, " An error occurred. Please try again.")
		return
	}

	if len(subscriptions) == 0 {
		services.SendMessage(chatID, " You have no active subscriptions.\n\nUse /subscribe <owner/repo> to add one.")
		return
	}

	message := " Your active subscriptions:\n\n"
	for i, sub := range subscriptions {
		message += fmt.Sprintf("%d. %s\n", i+1, sub.RepoFullName)
	}
	message += "\nUse /unsubscribe <owner/repo> to remove a subscription."
	services.SendMessage(chatID, message)
}

func handleHelp(chatID int64) {
	message := " Available Commands:\n\n" +
		"/start - Start the bot\n" +
		"/subscribe <owner/repo> - Subscribe to a repository\n" +
		"/unsubscribe <owner/repo> - Unsubscribe from a repository\n" +
		"/list - List your subscriptions\n" +
		"/help - Show this help message\n\n" +
		" Example:\n" +
		"/subscribe octocat/Hello-World"
	services.SendMessage(chatID, message)
}
