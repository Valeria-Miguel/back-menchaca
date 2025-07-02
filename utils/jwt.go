package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(email, rol string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"rol":   rol,
		"exp":   time.Now().Add(5 * time.Minute).Unix(), // 5 minutos
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GenerateRefreshToken(email, rol string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"rol":   rol,
		"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 días
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
}