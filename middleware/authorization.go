package middleware

import (
	"back-menchaca/config"
	"github.com/gofiber/fiber/v2"
)

func AutorizarPorPermiso() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extraer permisos desde el token
		permInterface := c.Locals("permisos")
		permisosSlice, ok := permInterface.([]interface{})
		if !ok {
			// Puede ser solo un string
			if singlePerm, ok := permInterface.(string); ok {
				permisosSlice = []interface{}{singlePerm}
			} else {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"statusCode": 401,
					"message":    "Permisos inválidos en el token",
				})
			}
		}

		metodo := c.Method()
		ruta := c.OriginalURL()

		// Buscar en la tabla si alguno de los permisos del token es válido para esta ruta y método
		var permitido bool
		for _, permiso := range permisosSlice {
			err := config.DB.QueryRow(`
				SELECT permitido FROM permisos 
				WHERE permiso = $1 AND metodo = $2 AND $3 ILIKE ruta
			`, permiso, metodo, ruta).Scan(&permitido)

			if err == nil && permitido {
				// Acceso concedido
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"statusCode": 403,
			"message":    "Acceso denegado. No tienes permiso suficiente.",
		})
	}
}

