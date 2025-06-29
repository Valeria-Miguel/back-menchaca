package routes

import (
    "back-menchaca/handlers"
    "github.com/gofiber/fiber/v2"
	"back-menchaca/middleware"
)


func ExpedienteRoutes(app fiber.Router) {
	expediente := app.Group("/expediente", middleware.JWTProtected("empleado"))

    expediente.Post("/", handlers.CrearExpediente)
    expediente.Get("/get", handlers.ObtenerExpedientes)
    expediente.Post("/getExp", handlers.ObtenerExpedientePorID)
    expediente.Put("/update", handlers.ActualizarExpediente)
    expediente.Delete("/delete", handlers.EliminarExpediente)
}
