package util

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

var ErrPasswordMismatch = errors.New("password does not match")

func ComparePw(hashedPassword, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		// If the error is not nil, the passwords don't match
		return ErrPasswordMismatch
	}
	return nil
}
