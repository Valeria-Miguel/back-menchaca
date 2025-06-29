package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	
	"back-menchaca/config"
	"back-menchaca/routes"
)


func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error cargando archivo .env")
	}

	config.ConnectDB()
	app := fiber.New()
	api := app.Group("/api")
	routes.SetupAuthRoutes(api)
	routes.SetupPacienteRoutes(api)
	routes.SetupEmpleadoRoutes(api)
	routes.SetupConsultorioRoutes(api)
	routes.ConsultasRoutes(api)
	routes.SetupHorarioRoutes(api)
	routes.HistorialRoutes(api)
	routes.SetupRecetasRoutes(api)
	routes.ExpedienteRoutes(api)
	routes.AntecedentesRoutes(api)
	routes.ReportesRoutes(api)
	routes.AvisoRoutes(api)


log.Println("🚀 Servidor iniciado en http://localhost:3000")
log.Fatal(app.Listen(":3000"))
}

