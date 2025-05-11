package jwtutil

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"server/config"
	"server/models"
)

func init() {
	config.LoadEnv()
}

func GenerateJWT(user models.User) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
