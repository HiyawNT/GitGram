package models

import "time"

type Subscription struct {
	Id           int64     `json:"id" gorm:"primaryKey"`
	ChatID       int64     `json:"chat_id" gorm:"index"`
	RepoFullName string    `json:"repo_full_name" gorm:"index"`
	CreatedAt    time.Time `json:"created_at"`
}
