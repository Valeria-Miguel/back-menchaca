package routes

import (
	"github.com/gofiber/fiber/v2"
	"back-menchaca/handlers"
)

func SetupLogRoutes(app fiber.Router) {
	logs := app.Group("/logs")
	logs.Get("/", handlers.GetLogs)
}
