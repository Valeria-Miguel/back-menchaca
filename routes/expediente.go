package routes

import (
    "back-menchaca/handlers"
    "github.com/gofiber/fiber/v2"
	"back-menchaca/middleware"
)


func ExpedienteRoutes(app fiber.Router) {
	expediente := app.Group("/expediente")

    expediente.Post("/", handlers.CrearExpediente)
    expediente.Get("/get",middleware.JWTProtected("solicitar_cita"), handlers.ObtenerExpedientes)
    expediente.Post("/getExp",middleware.JWTProtected("solicitar_cita"), handlers.ObtenerExpedientePorID)
    expediente.Put("/update", handlers.ActualizarExpediente)
    expediente.Delete("/delete", handlers.EliminarExpediente)
}
