package storage

import (
	"fmt"
	"sync"
)

// Subscription represents a user's subscription to a GitHub repository
type Subscription struct {
	ChatID     int64
	Repository string
}

// subscriptionStore holds the in-memory subscription data
// TODO: Replace with persistent storage (PostgreSQL, MongoDB, BoltDB, or SQLite)
// in production environment. Current implementation uses an in-memory map which
// will lose all data on application restart.
var (
	// Map structure: repository -> map[chatID]bool
	subscriptions = make(map[string]map[int64]bool)
	
	// Map structure: chatID -> map[repository]bool (for quick user lookup)
	userSubscriptions = make(map[int64]map[string]bool)
	mu                sync.RWMutex
)

// Subscribe adds a new subscription for a chat ID to a repository
func Subscribe(chatID int64, repository string) error {
	mu.Lock()
	defer mu.Unlock()

	// Initialize repository map if it doesn't exist
	if subscriptions[repository] == nil {
		subscriptions[repository] = make(map[int64]bool)
	}

	// Initialize user subscription map if it doesn't exist
	if userSubscriptions[chatID] == nil {
		userSubscriptions[chatID] = make(map[string]bool)
	}

	// Check if already subscribed
	if subscriptions[repository][chatID] {
		return fmt.Errorf("already subscribed to %s", repository)
	}

	// Add subscription
	subscriptions[repository][chatID] = true
	userSubscriptions[chatID][repository] = true

	return nil
}

// Unsubscribe removes a subscription for a chat ID from a repository
func Unsubscribe(chatID int64, repository string) error {
	mu.Lock()
	defer mu.Unlock()

	// Check if subscription exists
	if subscriptions[repository] == nil || !subscriptions[repository][chatID] {
		return fmt.Errorf("not subscribed to %s", repository)
	}

	// Remove subscription
	delete(subscriptions[repository], chatID)
	delete(userSubscriptions[chatID], repository)

	// Clean up empty maps
	if len(subscriptions[repository]) == 0 {
		delete(subscriptions, repository)
	}
	if len(userSubscriptions[chatID]) == 0 {
		delete(userSubscriptions, chatID)
	}

	return nil
}

// GetSubscriptions returns a list of repositories the given chat ID is subscribed to
func GetSubscriptions(chatID int64) []string {
	mu.RLock()
	defer mu.RUnlock()

	repos := make([]string, 0)
	if userSubscriptions[chatID] != nil {
		for repo := range userSubscriptions[chatID] {
			repos = append(repos, repo)
		}
	}

	return repos
}

// IsSubscribed checks if a chat ID is subscribed to a repository
func IsSubscribed(chatID int64, repository string) bool {
	mu.RLock()
	defer mu.RUnlock()

	if subscriptions[repository] == nil {
		return false
	}

	return subscriptions[repository][chatID]
}

// GetSubscribers returns a list of chat IDs subscribed to a repository
func GetSubscribers(repository string) []int64 {
	mu.RLock()
	defer mu.RUnlock()

	chatIDs := make([]int64, 0)
	if subscriptions[repository] != nil {
		for chatID := range subscriptions[repository] {
			chatIDs = append(chatIDs, chatID)
		}
	}

	return chatIDs
}

// GetAllSubscriptions returns all subscriptions (useful for debugging)
func GetAllSubscriptions() map[string][]int64 {
	mu.RLock()
	defer mu.RUnlock()

	result := make(map[string][]int64)
	for repo, chatMap := range subscriptions {
		chatIDs := make([]int64, 0, len(chatMap))
		for chatID := range chatMap {
			chatIDs = append(chatIDs, chatID)
		}
		result[repo] = chatIDs
	}

	return result
}