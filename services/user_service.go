package services

import (
	"auth/db/repositories"
	"auth/util"
	"errors"
	"fmt"
	"log"
)

func RegisterUser(username string, password string) error {
	return repositories.InsertUser(username, password)
}

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("incorrect password")
)

func LoginUser(username, password string) error {
	user, err := repositories.FindUserByUsername(username)
	if err != nil {
		log.Printf("Username %s not found", username)
		log.Printf("%s: ", err)
		return ErrUserNotFound
	}

	err = util.ComparePw(user.Password, password)
	if err != nil {
		if errors.Is(err, util.ErrPasswordMismatch) {
			log.Printf("Password mismatch for user: %s", username)
			return ErrWrongPassword
		}
		return fmt.Errorf("failed to compare password: %w", err)
	}

	log.Printf("User %s authenticated successfully", username)
	return nil
}
