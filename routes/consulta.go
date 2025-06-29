package routes

import (
	
	"back-menchaca/handlers"
	"back-menchaca/middleware"
	"github.com/gofiber/fiber/v2"
)

func ConsultasRoutes(app fiber.Router) {
	consultas := app.Group("/consultas", middleware.JWTProtected("empleado", "paciente"))

	consultas.Post("/", handlers.AgendarConsulta)
	consultas.Get("/", handlers.ObtenerConsultas)
	consultas.Post("/getConsl", handlers.ObtenerConsultaPorID)
	consultas.Put("/update", handlers.ActualizarConsulta)
	consultas.Delete("/delete", handlers.EliminarConsulta)
}
