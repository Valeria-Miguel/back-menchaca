package handlers

import (
	"back-menchaca/config"
	"back-menchaca/models"
	"back-menchaca/utils"
	"database/sql"
	"github.com/gofiber/fiber/v2"
)

const mod = "ANT"

func CrearAntecedente(c *fiber.Ctx) error {
	var a models.Antecedente
	if err := c.BodyParser(&a); err != nil {
		return utils.Responder(c, "02", mod, "antecedente-service", nil, "Datos inválidos")
	}

	if !utils.ExisteID("Expediente", "id_expediente", a.IDExpediente) {
		return utils.Responder(c, "02", mod, "antecedente-service", nil, "ID de expediente no válido")
	}

	if err := utils.ValidarAntecedente(a.Diagnostico, a.Descripcion, a.Fecha, a.IDExpediente); err != nil {
		return utils.Responder(c, "02", mod, "antecedente-service", nil, err.Error())
	}

	query := `INSERT INTO Antecedentes (id_expediente, diagnostico, descripcion, fecha)
	          VALUES ($1, $2, $3, $4) RETURNING id_antecedente`
	err := config.DB.QueryRow(query, a.IDExpediente, a.Diagnostico, a.Descripcion, a.Fecha).Scan(&a.ID)
	if err != nil {
		return utils.Responder(c, "06", mod, "antecedente-service", nil, "Error al crear antecedente")
	}

	return utils.Responder(c, "01", mod, "antecedente-service", a)
}

func ObtenerAntecedentes(c *fiber.Ctx) error {
	rows, err := config.DB.Query(`SELECT id_antecedente, id_expediente, diagnostico, descripcion, fecha FROM Antecedentes`)
	if err != nil {
		return utils.Responder(c, "06", mod, "antecedente-service", nil, "Error al obtener antecedentes")
	}
	defer rows.Close()

	var antecedentes []models.Antecedente
	for rows.Next() {
		var a models.Antecedente
		if err := rows.Scan(&a.ID, &a.IDExpediente, &a.Diagnostico, &a.Descripcion, &a.Fecha); err == nil {
			antecedentes = append(antecedentes, a)
		}
	}
	return utils.Responder(c, "01", mod, "antecedente-service", antecedentes)
}

func ObtenerAntecedentePorID(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_antecedente"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", mod, "antecedente-service", nil, "ID inválido")
	}

	var a models.Antecedente
	err := config.DB.QueryRow(`SELECT id_antecedente, id_expediente, diagnostico, descripcion, fecha FROM Antecedentes WHERE id_antecedente=$1`, body.ID).
		Scan(&a.ID, &a.IDExpediente, &a.Diagnostico, &a.Descripcion, &a.Fecha)
	if err == sql.ErrNoRows {
		return utils.Responder(c, "05", mod, "antecedente-service", nil, "Antecedente no encontrado")
	} else if err != nil {
		return utils.Responder(c, "06", mod, "antecedente-service", nil, "Error al buscar antecedente")
	}
	return utils.Responder(c, "01", mod, "antecedente-service", a)
}

func ActualizarAntecedente(c *fiber.Ctx) error {
	var a models.Antecedente
	if err := c.BodyParser(&a); err != nil || a.ID == 0 {
		return utils.Responder(c, "02", mod, "antecedente-service", nil, "Datos inválidos")
	}

	var actual models.Antecedente
	err := config.DB.QueryRow(`SELECT id_expediente, diagnostico, descripcion, fecha FROM Antecedentes WHERE id_antecedente=$1`, a.ID).
		Scan(&actual.IDExpediente, &actual.Diagnostico, &actual.Descripcion, &actual.Fecha)
	if err == sql.ErrNoRows {
		return utils.Responder(c, "05", mod, "antecedente-service", nil, "Antecedente no encontrado")
	} else if err != nil {
		return utils.Responder(c, "06", mod, "antecedente-service", nil, "Error al buscar antecedente")
	}

	if a.IDExpediente == 0 {
		a.IDExpediente = actual.IDExpediente
	} else if !utils.ExisteID("Expediente", "id_expediente", a.IDExpediente) {
		return utils.Responder(c, "02", mod, "antecedente-service", nil, "ID de expediente no válido")
	}

	if a.Diagnostico == "" {
		a.Diagnostico = actual.Diagnostico
	}
	if a.Descripcion == "" {
		a.Descripcion = actual.Descripcion
	}
	if a.Fecha.IsZero() {
		a.Fecha = actual.Fecha
	}

	if err := utils.ValidarAntecedente(a.Diagnostico, a.Descripcion, a.Fecha, a.IDExpediente); err != nil {
		return utils.Responder(c, "02", mod, "antecedente-service", nil, err.Error())
	}

	_, err = config.DB.Exec(`UPDATE Antecedentes SET id_expediente=$1, diagnostico=$2, descripcion=$3, fecha=$4 WHERE id_antecedente=$5`,
		a.IDExpediente, a.Diagnostico, a.Descripcion, a.Fecha, a.ID)
	if err != nil {
		return utils.Responder(c, "06", mod, "antecedente-service", nil, "Error al actualizar antecedente")
	}
	return utils.Responder(c, "01", mod, "antecedente-service", fiber.Map{"mensaje": "Antecedente actualizado"})
}

func EliminarAntecedente(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_antecedente"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", mod, "antecedente-service", nil, "ID inválido")
	}

	_, err := config.DB.Exec("DELETE FROM Antecedentes WHERE id_antecedente=$1", body.ID)
	if err != nil {
		return utils.Responder(c, "06", mod, "antecedente-service", nil, "Error al eliminar antecedente")
	}
	return utils.Responder(c, "01", mod, "antecedente-service", fiber.Map{"mensaje": "Antecedente eliminado"})
}
