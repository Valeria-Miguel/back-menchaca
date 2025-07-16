package routes

import (
	"back-menchaca/handlers"
	"back-menchaca/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupHorarioRoutes(app fiber.Router) {
	horario := app.Group("/horarios")
	
	horario.Post("/create", handlers.CrearHorario)
	horario.Get("/get",middleware.JWTProtected("solicitar_cita"), handlers.ObtenerHorarios)
	horario.Post("/gethorario", handlers.ObtenerHorarioPorID)
	horario.Put("/update", handlers.ActualizarHorario)
	horario.Delete("/delete", handlers.EliminarHorario)
}
