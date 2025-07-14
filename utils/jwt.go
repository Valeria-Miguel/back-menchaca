package utils

import (
	"os"
	"time"
	"back-menchaca/config"
	"log"
	"github.com/golang-jwt/jwt/v5"
)

func GetPermisosPorRol(rol string) ([]string, error) {
	rows, err := config.DB.Query(`
		SELECT permiso FROM permisos WHERE rol = $1
	`, rol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permisos []string
	for rows.Next() {
		var permiso string
		if err := rows.Scan(&permiso); err != nil {
			log.Println("Error scan permisos:", err)
			continue
		}
		permisos = append(permisos, permiso)
	}

	return permisos, nil
}
func GenerateJWT(email, rol string) (string, error) {
	permisos, err := GetPermisosPorRol(rol)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"email":    email,
		"rol":      rol,
		"permisos": permisos,
		"exp":      time.Now().Add(60 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}



func GenerateRefreshToken(email, rol string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"rol":   rol,
		"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(), 
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
}
// GenerateTempToken genera un token temporal para verificación MFA
func GenerateTempToken(email, rol string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"rol":   rol,
		"exp":   time.Now().Add(5 * time.Minute).Unix(), // 
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("TEMP_SECRET")))
}

