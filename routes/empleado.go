package routes

import (
	"github.com/gofiber/fiber/v2"
	"back-menchaca/handlers"
	"back-menchaca/middleware"
)

func SetupEmpleadoRoutes(app fiber.Router) {
	empleado := app.Group("/empleados", middleware.JWTProtected("paciente"))

	empleado.Post("/", handlers.CrearEmpleado)
	empleado.Get("/get", handlers.ObtenerEmpleados)
	empleado.Post("/getempleado", handlers.ObtenerEmpleadoPorID)
	empleado.Put("/update", handlers.ActualizarEmpleado)
	empleado.Delete("/delete", handlers.EliminarEmpleado)
}
