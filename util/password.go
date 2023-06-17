package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword: returns the bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to generate the hash of the password: %v", err)
	}
	return string(hash), nil
}

// CheckPassword: Check if the provided password is correct
func CheckPassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
