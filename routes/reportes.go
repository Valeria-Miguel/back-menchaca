package routes

import (
	"back-menchaca/handlers"
	"back-menchaca/middleware"
	"github.com/gofiber/fiber/v2"
)

func ReportesRoutes(app fiber.Router) {
	rep := app.Group("/reportes", middleware.JWTProtected("paciente", "empleados"))

	// POST
	rep.Post("/consultas-por-paciente-detalle", handlers.ReporteDetalleConsultasPorPaciente)
	rep.Post("/detalles-consulta-expediente", handlers.ReporteDetallesConsultaExpediente)

	// GET
	rep.Get("/consultas-por-area", handlers.ReporteConsultasPorArea)
	rep.Get("/consultas-por-turno", handlers.ReporteConsultasPorTurno)
	rep.Get("/ingresos-por-consultorio", handlers.ReporteIngresosPorConsultorio)
	rep.Get("/consultas-detalle-simple", handlers.ObtenerDetalleSimpleConsultas)
}
