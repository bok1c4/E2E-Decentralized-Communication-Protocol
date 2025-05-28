package models

import (
	"time"
)

type Channel struct {
	ID        uint      `gorm:"primaryKey"`
	Name      *string   `gorm:"uniqueIndex"` // optional for DMs
	IsDirect  bool      `gorm:"not null;default:false"`
	Users     []*User   `gorm:"many2many:channel_users"`
	Messages  []Message `gorm:"foreignKey:ChannelID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChannelUser struct {
	ChannelID uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"primaryKey"`
	JoinedAt  time.Time `gorm:"autoCreateTime"`
}
