package models

import "time"

type Message struct {
	ID       uint `gorm:"primaryKey"`
	SenderID uint `gorm:"not null"`
	Sender   User `gorm:"foreignKey:SenderID"`

	ChannelID uint    `gorm:"not null"` // Always present
	Channel   Channel `gorm:"foreignKey:ChannelID"`

	Content   string `gorm:"not null"`
	CreatedAt time.Time
}

type MessageWithUser struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
