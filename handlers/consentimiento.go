package handlers

import (
	"back-menchaca/config"
	"github.com/gofiber/fiber/v2"
	"time"
	"back-menchaca/utils"
)

// GET /api/aviso-privacidad
func ObtenerAvisoPrivacidad(c *fiber.Ctx) error {
	aviso := `
<h2>Aviso de Privacidad</h2>
<p>Este hospital garantiza la protección de sus datos personales conforme a lo establecido por la Ley Federal de Protección de Datos Personales.</p>
<p>Los datos que se recolectan serán utilizados únicamente para fines médicos y administrativos internos.</p>
<p>Para más información, puede contactar a nuestro departamento legal.</p>
`
	return c.Type("html").SendString(aviso)
}

// POST /api/consentimiento
func RegistrarConsentimiento(c *fiber.Ctx) error {
	var body struct {
		IDPaciente int `json:"id_paciente"`
	}
	if err := c.BodyParser(&body); err != nil || body.IDPaciente == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ID de paciente inválido"})
	}

	if !utils.ExisteID("Paciente", "id_paciente", body.IDPaciente) {
		return c.Status(404).JSON(fiber.Map{"error": "Paciente no encontrado"})
	}

	_, err := config.DB.Exec(`INSERT INTO Consentimientos (id_paciente, fecha_hora) VALUES ($1, $2)`,
		body.IDPaciente, time.Now())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al registrar consentimiento"})
	}

	return c.JSON(fiber.Map{"mensaje": "Consentimiento registrado correctamente"})
}
