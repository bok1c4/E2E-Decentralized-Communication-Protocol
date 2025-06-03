package repositories

import (
	"auth/db"
	"auth/db/models"
	"fmt"
)

func FetchOpenChannels() ([]models.Channel, error) {
	var channels []models.Channel
	err := db.DB.Where("is_direct = ?", false).Find(&channels).Error
	return channels, err
}

func FindChanBetweenTwoUsers(userID1, userID2 uint) (*models.Channel, error) {
	var channel models.Channel

	err := db.DB.
		Joins("JOIN channel_users cu ON cu.channel_id = channels.id").
		Where("channels.is_direct = ?", true).
		Where("cu.user_id IN ?", []uint{userID1, userID2}).
		Group("channels.id").
		Having("COUNT(DISTINCT cu.user_id) = 2").
		First(&channel).Error
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func CreateChanBetweenTwoUsers(userID1, userID2 uint) (*models.Channel, error) {
	var users []models.User
	err := db.DB.Where("id IN ?", []uint{userID1, userID2}).Find(&users).Error
	if err != nil {
		return nil, err
	}

	if len(users) != 2 {
		return nil, fmt.Errorf("one or both users not found")
	}

	channel := models.Channel{
		IsDirect: true,
		Users:    []*models.User{&users[0], &users[1]},
	}

	err = db.DB.Create(&channel).Error
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func GetUsersInChannel(channelID uint) ([]models.User, error) {
	var users []models.User

	err := db.DB.
		Joins("JOIN channel_users cu ON cu.user_id = users.id").
		Where("cu.channel_id = ?", channelID).
		Find(&users).Error

	return users, err
}

func IsDirectChannel(chanID uint) (bool, error) {
	var isDirect bool

	err := db.DB.
		Table("channels").
		Select("is_direct").
		Where("id = ?", chanID).
		Scan(&isDirect).Error
	if err != nil {
		return false, err
	}

	return isDirect, nil
}
