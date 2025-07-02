package handlers

import (
	"os"
	"time"
	
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	
	"back-menchaca/utils"
	"back-menchaca/config"
)

var validate = validator.New()

func Login(c *fiber.Ctx) error {
	var input struct {
		Correo      string `json:"correo" validate:"required,email"`
		Contrasena  string `json:"contrasena" validate:"required,min=6"`
	}

	// Parsear y validar entrada
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"statusCode": 400,
			"intCode": "A01",
			"message": "Datos de entrada inválidos",
			"from": "auth-service",
		})
	}

	// Validar esquema
	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"statusCode": 400,
			"intCode": "A01",
			"message": "Validación fallida: " + err.Error(),
			"from": "auth-service",
		})
	}

	// Buscar usuario en la base de datos
	var (
		rol  string
		hash string
	)

	// Primero buscar como empleado
	err := config.DB.QueryRow(`
		SELECT 'empleado' as rol, contraseña FROM Empleado WHERE correo=$1
		UNION
		SELECT 'paciente' as rol, contraseña FROM Paciente WHERE correo=$1`, input.Correo).Scan(&rol, &hash)
	
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"statusCode": 401,
			"intCode": "A01",
			"message": "Credenciales inválidas",
			"from": "auth-service",
		})
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input.Contrasena)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"statusCode": 401,
			"intCode": "A01",
			"message": "Credenciales inválidas",
			"from": "auth-service",
		})
	}

	// Generar tokens
	accessToken, err := utils.GenerateJWT(input.Correo, rol)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"statusCode": 500,
			"intCode": "A03",
			"message": "Error generando token de acceso",
			"from": "auth-service",
		})
	}

	refreshToken, err := utils.GenerateRefreshToken(input.Correo, rol)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"statusCode": 500,
			"intCode": "A03",
			"message": "Error generando token de refresco",
			"from": "auth-service",
		})
	}

	// Configurar cookie para refresh token
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour), // 7 días
		HTTPOnly: true,
		Secure:   true, // Solo en HTTPS
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"statusCode": 200,
		"intCode": "S01",
		"message": "Autenticación exitosa",
		"from": "auth-service",
		"data": []fiber.Map{
			{
				"token":         accessToken,
				"tokenType":    "Bearer",
				"expiresIn":    30 * 60, // 30 minutos en segundos
				"refreshToken": refreshToken,
			},
		},
	})
}

func RefreshToken(c *fiber.Ctx) error {
	// Opción 1: Obtener de cookies
	refreshToken := c.Cookies("refresh_token")
	
	// Opción 2: Obtener del body (alternativa)
	if refreshToken == "" {
		var input struct {
			RefreshToken string `json:"refreshToken" validate:"required"`
		}
		
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"statusCode": 400,
				"intCode": "A01",
				"message": "Datos de entrada inválidos",
				"from": "auth-service",
			})
		}
		
		// Validar esquema
		if err := validate.Struct(input); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"statusCode": 400,
				"intCode": "A01",
				"message": "Validación fallida: " + err.Error(),
				"from": "auth-service",
			})
		}
		
		refreshToken = input.RefreshToken
	}

	if refreshToken == "" {
		return c.Status(401).JSON(fiber.Map{
			"statusCode": 401,
			"intCode": "A01",
			"message": "Refresh token requerido",
			"from": "auth-service",
		})
	}

	// Validar token
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"statusCode": 401,
			"intCode": "A01",
			"message": "Refresh token inválido o expirado",
			"from": "auth-service",
		})
	}

	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	rol := claims["rol"].(string)

	// Verificar si el usuario aún existe
	var userExists bool
	err = config.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM Empleado WHERE correo=$1 
			UNION 
			SELECT 1 FROM Paciente WHERE correo=$1
		)`, email).Scan(&userExists)
		
	if err != nil || !userExists {
		return c.Status(401).JSON(fiber.Map{
			"statusCode": 401,
			"intCode": "A01",
			"message": "Usuario no encontrado",
			"from": "auth-service",
		})
	}

	// Generar nuevo access token
	newToken, err := utils.GenerateJWT(email, rol)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"statusCode": 500,
			"intCode": "A03",
			"message": "Error generando nuevo token",
			"from": "auth-service",
		})
	}

	return c.JSON(fiber.Map{
		"statusCode": 200,
		"intCode": "S01",
		"message": "Token refrescado exitosamente",
		"from": "auth-service",
		"data": []fiber.Map{
			{
				"token":     newToken,
				"tokenType": "Bearer",
				"expiresIn": 30 * 60, // 30 minutos en segundos
			},
		},
	})
}