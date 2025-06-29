package routes

import (
	"back-menchaca/handlers"
	"back-menchaca/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupConsultorioRoutes(app fiber.Router) {
	consultorio := app.Group("/consultorios", middleware.JWTProtected("empleado", "paciente"))

	consultorio.Post("/", handlers.CrearConsultorio)
	consultorio.Get("/get", handlers.ObtenerConsultorios)
	consultorio.Post("/getconsultorio", handlers.ObtenerConsultorioPorID)
	consultorio.Put("/update", handlers.ActualizarConsultorio)
	consultorio.Delete("/delete", handlers.EliminarConsultorio)
}
