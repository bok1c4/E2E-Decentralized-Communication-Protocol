package services

import "auth/db/repositories"

func CreateMessage(user_id int, msg string) error {
	return repositories.InsertMessage(user_id, msg)
}
