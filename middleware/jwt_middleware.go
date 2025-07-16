package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strings"
)

func JWTProtected(requiredPerms ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{
				"statusCode": 401,
				"message":    "Token requerido",
				"from":       "auth-service",
			})
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"statusCode": 401,
				"message":    "Token inválido o expirado",
				"from":       "auth-service",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"statusCode": 401,
				"message":    "Claims inválidos",
			})
		}

		// Obtener permisos desde el token
		// Obtener permisos desde el token
		// Obtener permisos desde el token (manejo robusto)
permisosToken := map[string]bool{}
if permsRaw, ok := claims["permisos"]; ok {
	permsIface, ok := permsRaw.([]interface{})
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"statusCode": 401,
			"message":    "Formato inválido en permisos",
		})
	}

	for _, p := range permsIface {
		permStr, ok := p.(string)
		if !ok {
			continue
		}
		permisosToken[permStr] = true
	}
}


		
		// Verifica si el usuario tiene al menos uno de los permisos requeridos
		autorizado := false
		for _, reqPerm := range requiredPerms {
			if permisosToken[reqPerm] {
				autorizado = true
				break
			}
		}

		if !autorizado {
			return c.Status(403).JSON(fiber.Map{
				"statusCode": 403,
				"message":    "Permiso insuficiente",
			})
		}

		// Guardar info útil en el contexto
		c.Locals("email", claims["email"])
		c.Locals("permisos", permisosToken)

		return c.Next()
	}
}
