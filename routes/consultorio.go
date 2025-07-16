package routes

import (
	"back-menchaca/handlers"
	"back-menchaca/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupConsultorioRoutes(app fiber.Router) {
	consultorio := app.Group("/consultorios")

	consultorio.Post("/", handlers.CrearConsultorio)
	consultorio.Get("/get",middleware.JWTProtected("solicitar_cita"), handlers.ObtenerConsultorios)
	consultorio.Post("/getconsultorio", handlers.ObtenerConsultorioPorID)
	consultorio.Put("/update", handlers.ActualizarConsultorio)
	consultorio.Delete("/delete", handlers.EliminarConsultorio)
}
