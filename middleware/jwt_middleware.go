package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTProtected(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{
				"statusCode": 401,
				"intCode": "A01",
				"message": "Token requerido",
				"from": "auth-service",
			})
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"statusCode": 401,
				"intCode": "A01",
				"message": "Token inválido o expirado",
				"from": "auth-service",
			})
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("email", claims["email"])
		c.Locals("rol", claims["rol"])

		// Verificar si el rol está permitido
		if len(allowedRoles) > 0 {
			userRol := claims["rol"].(string)
			valid := false
			for _, r := range allowedRoles {
				if r == userRol {
					valid = true
					break
				}
			}
			if !valid {
				return c.Status(403).JSON(fiber.Map{
					"statusCode": 403,
					"intCode": "A02",
					"message": "Acceso denegado",
					"from": "auth-service",
				})
			}
		}

		return c.Next()
	}
}