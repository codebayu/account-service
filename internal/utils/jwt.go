package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UUID      string `json:"uuid"`
	Role      int    `json:"role"`
	IsRefresh bool   `json:"isRefresh"`
	jwt.RegisteredClaims
}

func GenerateToken(uuid string, duration time.Duration, isRefresh bool) (string, time.Time, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret" // Fallback for development
	}

	expiresAt := time.Now().Add(duration)
	claims := &JWTClaims{
		UUID:      uuid,
		Role:      0, // Default role
		IsRefresh: isRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}
