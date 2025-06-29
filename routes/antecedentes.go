package routes

import (
    "back-menchaca/handlers"
    "github.com/gofiber/fiber/v2"
	"back-menchaca/middleware"
)

func AntecedentesRoutes(app fiber.Router) {
	ants := app.Group("/antecedentes", middleware.JWTProtected("paciente", "empleado"))
    ants.Post("/", handlers.CrearAntecedente)
    ants.Get("/get", handlers.ObtenerAntecedentes)
    ants.Post("/getant", handlers.ObtenerAntecedentePorID)
    ants.Put("/update", handlers.ActualizarAntecedente)
    ants.Delete("/delete", handlers.EliminarAntecedente)
}




