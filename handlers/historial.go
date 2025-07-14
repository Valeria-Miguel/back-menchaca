package handlers

import (
	"back-menchaca/config"
	"back-menchaca/models"
	"back-menchaca/utils"
	"database/sql"
	"github.com/gofiber/fiber/v2"
)

const modHis = "HIST"

func CrearHistorialClinico(c *fiber.Ctx) error {
	var h models.HistorialClinico
	if err := c.BodyParser(&h); err != nil {
		return utils.Responder(c, "02", modHis, "historial-service", nil, "Datos inválidos")
	}

	if h.IDExpediente == 0 || h.IDConsulta == 0 {
		return utils.Responder(c, "02", modHis, "historial-service", nil, "Faltan campos obligatorios")
	}

	if !utils.ExisteIDHisotial("Expediente", "id_expediente", h.IDExpediente) {
		return utils.Responder(c, "02", modHis, "historial-service", nil, "ID de expediente no válido")
	}
	if !utils.ExisteIDHisotial("Consultas", "id_consulta", h.IDConsulta) {
		return utils.Responder(c, "02", modHis, "historial-service", nil, "ID de consulta no válido")
	}

	err := config.DB.QueryRow(`
		INSERT INTO Historial_Clinico (id_expediente, id_consultas)
		VALUES ($1, $2) RETURNING id_historial`,
		h.IDExpediente, h.IDConsulta,
	).Scan(&h.ID)
	if err != nil {
		return utils.Responder(c, "06", modHis, "historial-service", nil, "Error al crear historial clínico")
	}
	return utils.Responder(c, "01", modHis, "historial-service", h)
}

func ObtenerHistorialesClinicos(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT id_historial, id_expediente, id_consultas FROM Historial_Clinico")
	if err != nil {
		return utils.Responder(c, "06", modHis, "historial-service", nil, "Error al obtener historiales")
	}
	defer rows.Close()

	var historiales []models.HistorialClinico
	for rows.Next() {
		var h models.HistorialClinico
		if err := rows.Scan(&h.ID, &h.IDExpediente, &h.IDConsulta); err == nil {
			historiales = append(historiales, h)
		}
	}
	return utils.Responder(c, "01", modHis, "historial-service", historiales)
}

func ObtenerHistorialClinicoPorID(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_historial"`
	}

	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", modHis, "historial-service", nil, "ID inválido")
	}

	var h models.HistorialClinico
	err := config.DB.QueryRow(
		"SELECT id_historial, id_expediente, id_consultas FROM Historial_Clinico WHERE id_historial = $1",
		body.ID).Scan(&h.ID, &h.IDExpediente, &h.IDConsulta)

	if err == sql.ErrNoRows {
		return utils.Responder(c, "05", modHis, "historial-service", nil, "Historial no encontrado")
	} else if err != nil {
		return utils.Responder(c, "06", modHis, "historial-service", nil, "Error al buscar historial")
	}

	return utils.Responder(c, "01", modHis, "historial-service", h)
}

func ActualizarHistorialClinico(c *fiber.Ctx) error {
	var h models.HistorialClinico
	if err := c.BodyParser(&h); err != nil || h.ID == 0 {
		return utils.Responder(c, "02", modHis, "historial-service", nil, "Datos inválidos")
	}

	var actual models.HistorialClinico
	err := config.DB.QueryRow("SELECT id_expediente, id_consultas FROM Historial_Clinico WHERE id_historial = $1", h.ID).
		Scan(&actual.IDExpediente, &actual.IDConsulta)
	if err == sql.ErrNoRows {
		return utils.Responder(c, "05", modHis, "historial-service", nil, "Historial no encontrado")
	} else if err != nil {
		return utils.Responder(c, "06", modHis, "historial-service", nil, "Error al buscar historial")
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
		return utils.Responder(c, "06", modHis, "historial-service", nil, "Error al actualizar historial clínico")
	}
	return utils.Responder(c, "01", modHis, "historial-service", fiber.Map{"mensaje": "Historial actualizado"})
}

func EliminarHistorialClinico(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_historial"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", modHis, "historial-service", nil, "ID inválido")
	}

	_, err := config.DB.Exec("DELETE FROM Historial_Clinico WHERE id_historial = $1", body.ID)
	if err != nil {
		return utils.Responder(c, "06", modHis, "historial-service", nil, "Error al eliminar historial clínico")
	}
	return utils.Responder(c, "01", modHis, "historial-service", fiber.Map{"mensaje": "Historial eliminado"})
}
