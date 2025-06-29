package routes

import (
	"back-menchaca/handlers"
	"github.com/gofiber/fiber/v2"
	"back-menchaca/middleware"

)

func AvisoRoutes(app fiber.Router) {
	aviso := app.Group("/consentimiento", middleware.JWTProtected("paciente"))

	aviso.Get("/aviso-privacidad", handlers.ObtenerAvisoPrivacidad)
	aviso.Post("/consentimiento", handlers.RegistrarConsentimiento)
}
