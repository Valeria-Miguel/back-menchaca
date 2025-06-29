package routes

import (
	"back-menchaca/handlers"
	"back-menchaca/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRecetasRoutes(app fiber.Router) {
	rec := app.Group("/recetas", middleware.JWTProtected("paciente", "empleados"))

	rec.Post("/", handlers.CrearReceta)
	rec.Get("/get", handlers.ObtenerRecetas)
	rec.Post("/recetaget", handlers.ObtenerRecetaPorID)
	rec.Put("/update", handlers.ActualizarReceta)
	rec.Delete("/delete", handlers.EliminarReceta)
}
