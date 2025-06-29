package handlers

import (
	"back-menchaca/config"
	"back-menchaca/models"
	"back-menchaca/utils"
	"database/sql"
	"github.com/gofiber/fiber/v2"
)

func CrearHistorialClinico(c *fiber.Ctx) error {
	var h models.HistorialClinico
	if err := c.BodyParser(&h); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	if h.IDExpediente == 0 || h.IDConsulta == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Faltan campos obligatorios"})
	}

	if !utils.ExisteIDHisotial("Expediente", "id_expediente", h.IDExpediente) {
		return c.Status(400).JSON(fiber.Map{"error": "ID de expediente no válido"})
	}
	if !utils.ExisteIDHisotial("Consultas", "id_consulta", h.IDConsulta) {
		return c.Status(400).JSON(fiber.Map{"error": "ID de consulta no válido"})
	}

	err := config.DB.QueryRow(`
		INSERT INTO Historial_Clinico (id_expediente, id_consultas)
		VALUES ($1, $2) RETURNING id_historial`,
		h.IDExpediente, h.IDConsulta,
	).Scan(&h.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al crear historial clínico"})
	}
	return c.Status(201).JSON(h)
}

func ObtenerHistorialesClinicos(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT id_historial, id_expediente, id_consultas FROM Historial_Clinico")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener historiales"})
	}
	defer rows.Close()

	var historiales []models.HistorialClinico
	for rows.Next() {
		var h models.HistorialClinico
		if err := rows.Scan(&h.ID, &h.IDExpediente, &h.IDConsulta); err == nil {
			historiales = append(historiales, h)
		}
	}
	return c.JSON(historiales)
}

func ObtenerHistorialClinicoPorID(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_historial"`
	}

	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	var h models.HistorialClinico
	err := config.DB.QueryRow(
		"SELECT id_historial, id_expediente, id_consultas FROM Historial_Clinico WHERE id_historial = $1",
		body.ID).Scan(&h.ID, &h.IDExpediente, &h.IDConsulta)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Historial no encontrado"})
	}

	return c.JSON(h)
}

func ActualizarHistorialClinico(c *fiber.Ctx) error {
	var h models.HistorialClinico
	if err := c.BodyParser(&h); err != nil || h.ID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	var actual models.HistorialClinico
	err := config.DB.QueryRow("SELECT id_expediente, id_consultas FROM Historial_Clinico WHERE id_historial = $1", h.ID).
		Scan(&actual.IDExpediente, &actual.IDConsulta)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{"error": "Historial no encontrado"})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al buscar historial"})
	}

	if h.IDExpediente == 0 {
		h.IDExpediente = actual.IDExpediente
	}
	if h.IDConsulta == 0 {
		h.IDConsulta = actual.IDConsulta
	}

	_, err = config.DB.Exec(
		`UPDATE Historial_Clinico SET id_expediente=$1, id_consultas=$2 WHERE id_historial=$3`,
		h.IDExpediente, h.IDConsulta, h.ID,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar historial clínico"})
	}
	return c.JSON(fiber.Map{"mensaje": "Historial actualizado"})
}

func EliminarHistorialClinico(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_historial"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	_, err := config.DB.Exec("DELETE FROM Historial_Clinico WHERE id_historial = $1", body.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar historial clínico"})
	}
	return c.JSON(fiber.Map{"mensaje": "Historial eliminado"})
}
