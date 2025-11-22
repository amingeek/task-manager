// backend/utils/jwt.go

package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username string, expirationHours int) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		LogError("JWT_SECRET not configured", errors.New("missing environment variable"))
		return "", errors.New("JWT_SECRET is not set in environment")
	}

	if len(secret) < 32 {
		LogWarn("JWT_SECRET is too short", "minimum 32 characters recommended")
	}

	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "task-manager",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		LogError("Failed to generate token", err)
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// بررسی روش امضا
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		LogError("Token parsing failed", err)
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func RefreshToken(userID uint, username string) (string, error) {
	// مدت اعتبار توکن جدید: 7 روز
	return GenerateToken(userID, username, 24*7)
}
