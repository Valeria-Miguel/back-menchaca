package handlers

import (
	"back-menchaca/config"
	"back-menchaca/models"
	"back-menchaca/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

// POST /api/expedientes
func CrearExpediente(c *fiber.Ctx) error {
	var e models.Expediente
	if err := c.BodyParser(&e); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// Validar existencia del paciente
	if !utils.ExisteID("Paciente", "id_paciente", e.IDPaciente) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de paciente no válido"})
	}

	// Validar campo seguro
	if err := utils.ValidarSeguro(e.Seguro); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Fecha creación por defecto
	if e.FechaCreacion.IsZero() {
		e.FechaCreacion = time.Now()
	}

	query := `INSERT INTO Expediente (id_paciente, seguro, fecha_creacion) VALUES ($1, $2, $3) RETURNING id_expediente`
	err := config.DB.QueryRow(query, e.IDPaciente, e.Seguro, e.FechaCreacion).Scan(&e.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear expediente"})
	}
	return c.Status(fiber.StatusCreated).JSON(e)
}

// GET /api/expedientes
func ObtenerExpedientes(c *fiber.Ctx) error {
	rows, err := config.DB.Query(`SELECT id_expediente, id_paciente, seguro, fecha_creacion FROM Expediente`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener expedientes"})
	}
	defer rows.Close()

	var expedientes []models.Expediente
	for rows.Next() {
		var e models.Expediente
		if err := rows.Scan(&e.ID, &e.IDPaciente, &e.Seguro, &e.FechaCreacion); err == nil {
			expedientes = append(expedientes, e)
		}
	}
	return c.JSON(expedientes)
}

// POST /api/expedientes/get
func ObtenerExpedientePorID(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_expediente"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	// Validar existencia del expediente
	if !utils.ExisteIDExped(body.ID) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Expediente no encontrado"})
	}

	var e models.Expediente
	err := config.DB.QueryRow(`SELECT id_expediente, id_paciente, seguro, fecha_creacion FROM Expediente WHERE id_expediente=$1`, body.ID).
		Scan(&e.ID, &e.IDPaciente, &e.Seguro, &e.FechaCreacion)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al buscar expediente"})
	}
	return c.JSON(e)
}

// PUT /api/expedientes/update
func ActualizarExpediente(c *fiber.Ctx) error {
	var e models.Expediente
	if err := c.BodyParser(&e); err != nil || e.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// Validar existencia del expediente
	if !utils.ExisteIDExped(e.ID) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Expediente no encontrado"})
	}

	// Obtener datos actuales
	var actual models.Expediente
	err := config.DB.QueryRow(`SELECT id_paciente, seguro, fecha_creacion FROM Expediente WHERE id_expediente=$1`, e.ID).
		Scan(&actual.IDPaciente, &actual.Seguro, &actual.FechaCreacion)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener expediente actual"})
	}

	// Conservar datos no enviados
	if e.IDPaciente == 0 {
		e.IDPaciente = actual.IDPaciente
	} else if !utils.ExisteID("Paciente", "id_paciente", e.IDPaciente) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de paciente no válido"})
	}

	if strings.TrimSpace(e.Seguro) == "" {
		e.Seguro = actual.Seguro
	} else {
		if err := utils.ValidarSeguro(e.Seguro); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	}

	if e.FechaCreacion.IsZero() {
		e.FechaCreacion = actual.FechaCreacion
	}

	// Actualizar en DB
	_, err = config.DB.Exec(`UPDATE Expediente SET id_paciente=$1, seguro=$2, fecha_creacion=$3 WHERE id_expediente=$4`,
		e.IDPaciente, e.Seguro, e.FechaCreacion, e.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar expediente"})
	}
	return c.JSON(fiber.Map{"mensaje": "Expediente actualizado"})
}

// DELETE /api/expedientes/delete
func EliminarExpediente(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_expediente"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	// Validar existencia antes de eliminar
	if !utils.ExisteIDExped(body.ID) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Expediente no encontrado"})
	}

	_, err := config.DB.Exec("DELETE FROM Expediente WHERE id_expediente=$1", body.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al eliminar expediente"})
	}
	return c.JSON(fiber.Map{"mensaje": "Expediente eliminado"})
}
