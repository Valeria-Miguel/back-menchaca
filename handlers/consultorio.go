package handlers

import (
	"back-menchaca/config"
	"back-menchaca/models"
	"back-menchaca/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

const modConsultorio = "Consulorio"
func CrearConsultorio(c *fiber.Ctx) error {
	var cons models.Consultorio
	if err := c.BodyParser(&cons); err != nil {
		return utils.Responder(c, "02", modConsultorio, "consultorio-service", nil, "Datos inválidos")
	}

	if err := utils.ValidarConsultorio(cons.Nombre, cons.Tipo); err != nil {
		return utils.Responder(c, "02", modConsultorio, "consultorio-service", nil, err.Error())
	}

	// Sanitizar
	cons.Nombre = utils.SanitizarInput(cons.Nombre)
	cons.Tipo = utils.SanitizarInput(cons.Tipo)

	query := `INSERT INTO Consultorios (nombre, tipo) VALUES ($1, $2) RETURNING id_consultorio`
	err := config.DB.QueryRow(query, cons.Nombre, cons.Tipo).Scan(&cons.ID)
	if err != nil {
		return utils.Responder(c, "06", modConsultorio, "consultorio-service", nil, "Error al crear consultorio")
	}
	return utils.Responder(c, "01", modConsultorio, "consultorio-service", cons)
}

func ObtenerConsultorios(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT id_consultorio, nombre, tipo FROM Consultorios")
	if err != nil {
		return utils.Responder(c, "06", modConsultorio, "consultorio-service", nil, "Error al obtener consultorios")
	}
	defer rows.Close()

	var lista []models.Consultorio
	for rows.Next() {
		var cons models.Consultorio
		if err := rows.Scan(&cons.ID, &cons.Nombre, &cons.Tipo); err == nil {
			lista = append(lista, cons)
		}
	}
	return utils.Responder(c, "01", modConsultorio, "consultorio-service", lista)
}

func ObtenerConsultorioPorID(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_consultorio"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", modConsultorio, "consultorio-service", nil, "ID inválido")
	}

	var cons models.Consultorio
	err := config.DB.QueryRow(
		"SELECT id_consultorio, nombre, tipo FROM Consultorios WHERE id_consultorio = $1",
		body.ID).Scan(&cons.ID, &cons.Nombre, &cons.Tipo)

	if err != nil {
		return utils.Responder(c, "05", modConsultorio, "consultorio-service", nil, "Consultorio no encontrado")
	}
	return utils.Responder(c, "01", modConsultorio, "consultorio-service", cons)
}

func ActualizarConsultorio(c *fiber.Ctx) error {
	var cons models.Consultorio
	if err := c.BodyParser(&cons); err != nil || cons.ID == 0 {
		return utils.Responder(c, "02", modConsultorio, "consultorio-service", nil, "Datos inválidos")
	}

	var actual models.Consultorio
	err := config.DB.QueryRow("SELECT nombre, tipo FROM Consultorios WHERE id_consultorio = $1", cons.ID).
		Scan(&actual.Nombre, &actual.Tipo)
	if err != nil {
		return utils.Responder(c, "05", modConsultorio, "consultorio-service", nil, "Consultorio no encontrado")
	}

	if strings.TrimSpace(cons.Nombre) != "" {
		actual.Nombre = utils.SanitizarInput(cons.Nombre)
	}
	if strings.TrimSpace(cons.Tipo) != "" {
		actual.Tipo = utils.SanitizarInput(cons.Tipo)
	}

	_, err = config.DB.Exec(
		"UPDATE Consultorios SET nombre=$1, tipo=$2 WHERE id_consultorio=$3",
		actual.Nombre, actual.Tipo, cons.ID,
	)
	if err != nil {
		return utils.Responder(c, "06", modConsultorio, "consultorio-service", nil, "Error al actualizar consultorio")
	}
	return utils.Responder(c, "01", modConsultorio, "consultorio-service", fiber.Map{"mensaje": "Consultorio actualizado"})
}

func EliminarConsultorio(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_consultorio"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", modConsultorio, "consultorio-service", nil, "ID inválido")
	}
	_, err := config.DB.Exec("DELETE FROM Consultorios WHERE id_consultorio=$1", body.ID)
	if err != nil {
		return utils.Responder(c, "06", modConsultorio, "consultorio-service", nil, "Error al eliminar consultorio")
	}
	return utils.Responder(c, "01", modConsultorio, "consultorio-service", fiber.Map{"mensaje": "Consultorio eliminado"})
}
