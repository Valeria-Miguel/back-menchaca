package main

import (
	"log"
	"time"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
	"back-menchaca/middleware"
	"back-menchaca/config"
	"back-menchaca/routes"
)


func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error cargando archivo .env")
	}

	config.ConnectDB()

	app := fiber.New()

	
	app.Use(middleware.Logger())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:4200",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        200,                
		Expiration: 1 * time.Minute,   
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"statusCode": 429,
				"message":    "Demasiadas solicitudes, intenta más tarde.",
			})
		},
	}))


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


log.Println(" Servidor iniciado en http://localhost:3000")
log.Fatal(app.Listen(":3000"))
}

