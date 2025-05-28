package repositories

import (
	"auth/db"
	"auth/db/models"
)

// sender_id, channel_id, msg
func InsertChannelMsg(user_id int, channel_id int, message string) error {
	msg := models.Message{
		SenderID:  uint(user_id),
		ChannelID: uint(channel_id),
		Content:   message,
	}

	return db.DB.Create(&msg).Error
}

func GetMessagesFromChannelID(channelID uint) ([]models.MessageWithUser, error) {
	var msgs []models.MessageWithUser
	err := db.DB.Raw(`
        SELECT m.id, m.content, u.username, m.created_at
        FROM messages m
        JOIN users u ON m.sender_id = u.id
        WHERE m.channel_id = ?
        ORDER BY m.created_at ASC
    `, channelID).Scan(&msgs).Error
	return msgs, err
}
