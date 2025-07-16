package routes

import (
	"back-menchaca/handlers"
	"back-menchaca/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRecetasRoutes(app fiber.Router) {
	rec := app.Group("/recetas")

	rec.Post("/",middleware.JWTProtected("crear_recetas"), handlers.CrearReceta)
	rec.Get("/get",middleware.JWTProtected("ver_recetas"), handlers.ObtenerRecetas)

	
	rec.Post("/recetaget", middleware.JWTProtected("solicitar_cita"), handlers.ObtenerRecetaPorID)
	rec.Put("/update", handlers.ActualizarReceta)
	rec.Delete("/delete", handlers.EliminarReceta)
}
