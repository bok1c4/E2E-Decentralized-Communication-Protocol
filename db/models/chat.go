package models

import (
	"time"
)

type Message struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	UserID    int       `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"` // Auto-preload user info if needed
	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

type MessageWithUser struct {
	ID        int
	Username  string
	Content   string
	CreatedAt time.Time
}
