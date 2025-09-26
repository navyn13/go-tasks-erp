package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func CompareHashPassword(storedPasswordHash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(password))
	return err == nil

}
