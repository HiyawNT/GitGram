package services

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

var Bot *tgbotapi.BotAPI

func InitBot() {
	token := os.Getenv("TELEGRAM_TOKEN")
	var err error
	Bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", Bot.Self.UserName)
}

func SendMessage(chatID int64, message string) {
	if Bot == nil {
		InitBot()
	}
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := Bot.Send(msg)
	if err != nil {
		log.Println("Failed to send message:", err)
	}
}
