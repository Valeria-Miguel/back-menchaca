package routes

import (
	"back-menchaca/handlers"
	"back-menchaca/middleware"
	"github.com/gofiber/fiber/v2"
)

func HistorialRoutes(app fiber.Router) {
	historial := app.Group("/historial", middleware.JWTProtected("empleado"))

	historial.Post("/create", handlers.CrearHistorialClinico)
	historial.Get("/get", handlers.ObtenerHistorialesClinicos)
	historial.Post("/historialget", handlers.ObtenerHistorialClinicoPorID)
	historial.Put("/update", handlers.ActualizarHistorialClinico)
	historial.Delete("/delete", handlers.EliminarHistorialClinico)
}




