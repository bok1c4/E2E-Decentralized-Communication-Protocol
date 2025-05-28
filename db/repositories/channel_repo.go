package repositories

import (
	"auth/db"
	"auth/db/models"
)

func FetchOpenChannels() ([]models.Channel, error) {
	var channels []models.Channel
	err := db.DB.Where("is_direct = ?", false).Find(&channels).Error
	return channels, err
}
