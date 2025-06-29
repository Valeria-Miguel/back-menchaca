package routes

import (
	"github.com/gofiber/fiber/v2"
	"back-menchaca/handlers"
	"back-menchaca/middleware"
)

func SetupPacienteRoutes(app fiber.Router) {
	app.Post("/pacientes", handlers.CrearPaciente)

	//paciente := app.Group("/pacientes", middleware.JWTProtected())
	paciente := app.Group("/pacientes", middleware.JWTProtected("paciente"))
	
	paciente.Get("/get", handlers.ObtenerPacientes)
	paciente.Post("/getpaciente", handlers.ObtenerPacientePorID)
	paciente.Put("/update", handlers.ActualizarPaciente)
	paciente.Delete("/delete", handlers.EliminarPaciente)
}
