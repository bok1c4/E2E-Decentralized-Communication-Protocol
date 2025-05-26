package repositories

import (
	"auth/db"
	"auth/db/models"
)

func InsertMessage(user_id int, message string) error {
	msg := models.Message{
		UserID:  user_id,
		Content: message,
	}

	return db.DB.Create(&msg).Error
}

func GetRecentMessages() ([]models.MessageWithUser, error) {
	var msgs []models.MessageWithUser
	err := db.DB.Raw(`
        SELECT m.id, m.content, u.username, m.created_at
        FROM messages m
        JOIN users u ON m.user_id = u.id
        ORDER BY m.created_at
    `).Scan(&msgs).Error
	return msgs, err
}
