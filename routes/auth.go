package routes

import (
	"github.com/gofiber/fiber/v2"
	"back-menchaca/handlers"
)

func SetupAuthRoutes(app fiber.Router) {
	auth := app.Group("/auth")
	auth.Post("/login", handlers.Login)
	auth.Post("/refresh", handlers.RefreshToken)
}