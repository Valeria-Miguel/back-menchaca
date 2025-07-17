package routes

import (
	
	"back-menchaca/handlers"
	"back-menchaca/middleware"
	"github.com/gofiber/fiber/v2"
)

func ConsultasRoutes(app fiber.Router) {
	consultas := app.Group("/consultas")

	consultas.Post("/",middleware.JWTProtected("solicitar_cita"), handlers.AgendarConsulta)
	consultas.Get("/",middleware.JWTProtected("ver_citas"), handlers.ObtenerConsultas)
	consultas.Post("/getConsl",middleware.JWTProtected("solicitar_cita"), handlers.ObtenerConsultaPorID)
	consultas.Put("/update", handlers.ActualizarConsulta)
	consultas.Delete("/delete", handlers.EliminarConsulta)
	consultas.Post("/paciente/", middleware.JWTProtected("solicitar_cita"), handlers.ObtenerConsultasPaciente)
	consultas.Post("/doctor/", middleware.JWTProtected("solicitar_cita"), handlers.ObtenerConsultasPorEmpleado)

}
