package services

import "auth/db/repositories"

// this is where we are calling database operations
// and manipulating with data

func CreateMessage(user_id int, msg string) error {
	return repositories.InsertMessage(user_id, msg)
}
