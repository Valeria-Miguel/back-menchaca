package routes

import (
	"github.com/gofiber/fiber/v2"
	"back-menchaca/handlers"
	"back-menchaca/middleware"
)
/*
func SetupPacienteRoutes(app fiber.Router) {
	app.Post("/pacientes", handlers.CrearPaciente)

	//paciente := app.Group("/pacientes", middleware.JWTProtected())
	
	paciente := app.Group("/pacientes", middleware.AutorizarPorPermiso())

	paciente.Get("/get", handlers.ObtenerPacientes)
	paciente.Post("/getpaciente", handlers.ObtenerPacientePorID)
	paciente.Put("/update", handlers.ActualizarPaciente)
	paciente.Delete("/delete", handlers.EliminarPaciente)
}
*/

func SetupPacienteRoutes(app fiber.Router) {
	app.Post("/pacientes", handlers.CrearPaciente) // Registro libre (sin protección)

	// Agrupamos todas las rutas protegidas
	paciente := app.Group("/pacientes",
		middleware.JWTProtected(),           // 1️⃣ Verifica JWT y extrae el rol
		middleware.AutorizarPorPermiso(),    // 2️⃣ Verifica en BD si ese rol tiene permiso para acceder
	)

	paciente.Get("/get", handlers.ObtenerPacientes)             // paciente puede
	paciente.Post("/getpaciente", handlers.ObtenerPacientePorID) // paciente puede
	paciente.Put("/update", handlers.ActualizarPaciente)         // paciente puede
	paciente.Delete("/delete", handlers.EliminarPaciente)        // solo empleado puede
}