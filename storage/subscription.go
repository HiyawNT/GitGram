package subscription

import (
	"github.com/HiyawNT/GitGram/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// Init the DB

func InitDB() {

	var err error
	DB, err = gorm.Open(sqlite.Open("gitgram.db"), &gorm.Config{})
	if err != nil {
		log.Panic("Failed to Connect to database ", err)
	}

	//Time to Auto-migrate
	err = DB.AutoMigrate(&models.Subscription{})
	if err != nil {

		log.Panic("Failed to Migrate the Database ", err)
	}

	log.Println("Database Initialized successfully")

}

// Add New Subscription

func AddSubscription(ChatId int64, repoFullName string) error {

	subscription := models.Subscription{
		ChatID:       ChatId,
		RepoFullName: repoFullName,
	}
	result := DB.Create(&subscription)
	return result.Error

}

func RemoveSubscription(ChatId int64, repoFullName string) error {

	result := DB.Where("chat_id = ? AND repo_full_name = ?", ChatId, repoFullName).Delete(&models.Subscription{})

	return result.Error

}

func GetSubscriptionByRepo(repoFullName string) ([]int64, error) {
	var subscriptions []models.Subscription
	err := DB.Where("repo_full_name = ?", repoFullName).Find(&subscriptions).Error
	if err != nil {
		return nil, err
	}
	chatIDs := make([]int64, len(subscriptions))
	for i, sub := range subscriptions {
		chatIDs[i] = sub.ChatID
	}
	return chatIDs, nil

}

// GetSubscriptionsByChatID gets all repos a chat is subscribed to

func GetSubscriptionsByChatID(chatID int64) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := DB.Where("chat_id = ?", chatID).Find(&subscriptions).Error
	return subscriptions, err
}

func CheckSubscriptionExists(chatID int64, repoFullName string) (bool, error) {
	var count int64
	err := DB.Model(&models.Subscription{}).
		Where("chat_id = ? AND repo_full_name = ?", chatID, repoFullName).Count(&count).Error

	return count > 0, err
}
