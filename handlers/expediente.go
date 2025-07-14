package handlers

import (
	"back-menchaca/config"
	"back-menchaca/models"
	"back-menchaca/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

const modExp = "EXP"

func CrearExpediente(c *fiber.Ctx) error {
	var e models.Expediente
	if err := c.BodyParser(&e); err != nil {
		return utils.Responder(c, "02", modExp, "expediente-service", nil, "Datos inválidos")
	}

	if !utils.ExisteID("Paciente", "id_paciente", e.IDPaciente) {
		return utils.Responder(c, "02", modExp, "expediente-service", nil, "ID de paciente no válido")
	}

	if err := utils.ValidarSeguro(e.Seguro); err != nil {
		return utils.Responder(c, "02", modExp, "expediente-service", nil, err.Error())
	}

	if e.FechaCreacion.IsZero() {
		e.FechaCreacion = time.Now()
	}

	query := `INSERT INTO Expediente (id_paciente, seguro, fecha_creacion) VALUES ($1, $2, $3) RETURNING id_expediente`
	err := config.DB.QueryRow(query, e.IDPaciente, e.Seguro, e.FechaCreacion).Scan(&e.ID)
	if err != nil {
		return utils.Responder(c, "06", modExp, "expediente-service", nil, "Error al crear expediente")
	}

	return utils.Responder(c, "01", modExp, "expediente-service", e)
}

func ObtenerExpedientes(c *fiber.Ctx) error {
	rows, err := config.DB.Query(`SELECT id_expediente, id_paciente, seguro, fecha_creacion FROM Expediente`)
	if err != nil {
		return utils.Responder(c, "06", modExp, "expediente-service", nil, "Error al obtener expedientes")
	}
	defer rows.Close()

	var expedientes []models.Expediente
	for rows.Next() {
		var e models.Expediente
		if err := rows.Scan(&e.ID, &e.IDPaciente, &e.Seguro, &e.FechaCreacion); err == nil {
			expedientes = append(expedientes, e)
		}
	}
	return utils.Responder(c, "01", modExp, "expediente-service", expedientes)
}

func ObtenerExpedientePorID(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_expediente"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", modExp, "expediente-service", nil, "ID inválido")
	}

	if !utils.ExisteIDExped(body.ID) {
		return utils.Responder(c, "05", modExp, "expediente-service", nil, "Expediente no encontrado")
	}

	var e models.Expediente
	err := config.DB.QueryRow(`SELECT id_expediente, id_paciente, seguro, fecha_creacion FROM Expediente WHERE id_expediente=$1`, body.ID).
		Scan(&e.ID, &e.IDPaciente, &e.Seguro, &e.FechaCreacion)
	if err != nil {
		return utils.Responder(c, "06", modExp, "expediente-service", nil, "Error al buscar expediente")
	}
	return utils.Responder(c, "01", modExp, "expediente-service", e)
}

func ActualizarExpediente(c *fiber.Ctx) error {
	var e models.Expediente
	if err := c.BodyParser(&e); err != nil || e.ID == 0 {
		return utils.Responder(c, "02", modExp, "expediente-service", nil, "Datos inválidos")
	}

	if !utils.ExisteIDExped(e.ID) {
		return utils.Responder(c, "05", modExp, "expediente-service", nil, "Expediente no encontrado")
	}

	var actual models.Expediente
	err := config.DB.QueryRow(`SELECT id_paciente, seguro, fecha_creacion FROM Expediente WHERE id_expediente=$1`, e.ID).
		Scan(&actual.IDPaciente, &actual.Seguro, &actual.FechaCreacion)
	if err != nil {
		return utils.Responder(c, "06", modExp, "expediente-service", nil, "Error al obtener expediente actual")
	}

	if e.IDPaciente == 0 {
		e.IDPaciente = actual.IDPaciente
	} else if !utils.ExisteID("Paciente", "id_paciente", e.IDPaciente) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de paciente no válido"})
	}

	if strings.TrimSpace(e.Seguro) == "" {
		e.Seguro = actual.Seguro
	} else {
		if err := utils.ValidarSeguro(e.Seguro); err != nil {
			return utils.Responder(c, "02", modExp, "expediente-service", nil, err.Error())
		}
	}

	if e.FechaCreacion.IsZero() {
		e.FechaCreacion = actual.FechaCreacion
	}

	_, err = config.DB.Exec(`UPDATE Expediente SET id_paciente=$1, seguro=$2, fecha_creacion=$3 WHERE id_expediente=$4`,
		e.IDPaciente, e.Seguro, e.FechaCreacion, e.ID)
	if err != nil {
		return utils.Responder(c, "06", modExp, "expediente-service", nil, "Error al actualizar expediente")
	}
	return utils.Responder(c, "01", modExp, "expediente-service", fiber.Map{"mensaje": "Expediente actualizado"})
}

func EliminarExpediente(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_expediente"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", modExp, "expediente-service", nil, "ID inválido")
	}

	if !utils.ExisteIDExped(body.ID) {
		return utils.Responder(c, "05", modExp, "expediente-service", nil, "Expediente no encontrado")
	}

	_, err := config.DB.Exec("DELETE FROM Expediente WHERE id_expediente=$1", body.ID)
	if err != nil {
		return utils.Responder(c, "06", modExp, "expediente-service", nil, "Error al eliminar expediente")
	}
	return utils.Responder(c, "01", modExp, "expediente-service", fiber.Map{"mensaje": "Expediente eliminado"})
}
