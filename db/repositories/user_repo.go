package repositories

import (
	"auth/db"
	"auth/db/models"
	"auth/util"
	"errors"
	"fmt"
	"log"
)

var ErrUserExists = errors.New("username already exists")

func InsertUser(username, password string) error {
	_, err := FindUserByUsername(username)
	if err == nil {
		log.Printf("User with %s username already exists, please use different username", username)
		log.Printf("%s: ", err)
		return ErrUserExists
	}

	hashed_pw, err := util.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := models.User{
		Username: username,
		Password: hashed_pw,
	}
	return db.DB.Create(&user).Error
}

func FindUserByUsername(username string) (models.User, error) {
	var user models.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

func InsertPGPKey(username string, pgpKey string) error {
	user, err := FindUserByUsername(username)
	if err != nil {
		log.Printf("Failed to find user with username %s: %v", username, err)
		return fmt.Errorf("user not found: %w", err)
	}

	user.PGPKey = pgpKey
	if err := db.DB.Save(&user).Error; err != nil {
		log.Printf("Failed to save PGP key for %s: %v", username, err)
		return fmt.Errorf("failed to save PGP key: %w", err)
	}

	log.Printf("PGP key successfully inserted for user %s", username)
	return nil
}

func UpdateUser(user models.User) error {
	return db.DB.Save(&user).Error
}

func DeleteUser(user models.User) error {
	return db.DB.Delete(&user).Error
}
