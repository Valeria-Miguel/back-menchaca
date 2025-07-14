package routes

import (
	"github.com/gofiber/fiber/v2"
	"back-menchaca/handlers"
	"back-menchaca/middleware"
)

func SetupAuthRoutes(app fiber.Router) {
	auth := app.Group("/auth")
	auth.Post("/login", handlers.Login)
	auth.Post("/refresh", handlers.RefreshToken)
	auth.Post("/verify-mfa", handlers.VerifyMFA)
	
	// Nueva ruta para activar MFA
	auth.Post("/mfa/activate", middleware.JWTProtected(), handlers.ActivateMFA)
	
}