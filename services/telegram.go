package services

import (
	"log"
	"os"
	"strings"

	"github.com/HiyawNT/GitGram/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func InitBot() error {
	token := os.Getenv("TELEGRAM_TOKEN")
	var err error
	Bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", Bot.Self.UserName)

	// Start listening for updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := Bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message != nil { // If we got a message
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				// Handle commands
				if update.Message.IsCommand() {
					handleCommand(update.Message)
				}
			}
		}
	}()
	
	return err
}

func SendMessage(chatID int64, message string) ( error){
	if Bot == nil {
		InitBot()
	}
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := Bot.Send(msg)
	if err != nil {
		log.Println("Failed to send message:", err)
		return  err
	}
	return  nil
}

func handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		SendMessage(message.Chat.ID, "Wellcome to  GitGram Bot your One stop To Github Notification")
	case "subscribe":
		args := strings.Split(message.Text, " ")
		if len(args) != 2 {
			SendMessage(message.Chat.ID, "Usage: /subscribe <repo_full_name>")
			return
		}
		repoFullName := args[1]
		storage.Subscribe(message.Chat.ID, repoFullName)
		SendMessage(message.Chat.ID, "Subscribed to "+repoFullName)
	case "unsubscribe":
		args := strings.Split(message.Text, " ")
		if len(args) != 2 {
			SendMessage(message.Chat.ID, "Usage: /unsubscribe <repo_full_name>")
			return
		}
		repoFullName := args[1]
		storage.Unsubscribe(message.Chat.ID, repoFullName)
		SendMessage(message.Chat.ID, "Unsubscribed from "+repoFullName)

	case "help":
		helpStr :=  `
		- /start - Welcome message and list of commands 
		- /help - Display help information 
		- /subscribe <owner/repo> 
		- Subscribe to a repository      
			 Example: /subscribe octocat/Hello-World   
		-  /unsubscribe <owner/repo>  - Unsubscribe from a repository 
		-  /list_subscriptions  - View your active subscriptions`
		SendMessage(message.Chat.ID, helpStr)
	default:
		SendMessage(message.Chat.ID, "Unknown command")
	}
}