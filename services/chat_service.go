package services

import "auth/db/repositories"

func CreateMessage(user_id int, channel_id int, msg string) error {
	return repositories.InsertChannelMsg(user_id, channel_id, msg)
}
