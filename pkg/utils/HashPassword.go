package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// ** Hash password
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return password
	} else {
		return string(hash)
	}
}
