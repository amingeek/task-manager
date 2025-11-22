// backend/utils/password.go

package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// حداقل طول پسورد
	MinPasswordLength = 8
	// حداکثر طول پسورد
	MaxPasswordLength = 128
)

func HashPassword(password string) (string, error) {
	// بررسی طول پسورد
	if len(password) < MinPasswordLength {
		return "", &ValidationError{Message: "Password must be at least 8 characters long"}
	}
	if len(password) > MaxPasswordLength {
		return "", &ValidationError{Message: "Password must be at most 128 characters"}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		LogError("Failed to hash password", err)
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return &ValidationError{Message: "Password must be at least 8 characters long"}
	}
	if len(password) > MaxPasswordLength {
		return &ValidationError{Message: "Password must be at most 128 characters"}
	}
	return nil
}
