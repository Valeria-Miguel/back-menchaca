package handlers

import (
	"back-menchaca/config"
	"back-menchaca/models"
	"back-menchaca/utils"
	"github.com/gofiber/fiber/v2"
	"strings"

)

// POST /api/consultorios
func CrearConsultorio(c *fiber.Ctx) error {
	var cons models.Consultorio
	if err := c.BodyParser(&cons); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// Validar
	if err := utils.ValidarConsultorio(cons.Nombre, cons.Tipo); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Sanitizar
	cons.Nombre = utils.SanitizarInput(cons.Nombre)
	cons.Tipo = utils.SanitizarInput(cons.Tipo)

	// Insertar
	query := `INSERT INTO Consultorios (nombre, tipo) VALUES ($1, $2) RETURNING id_consultorio`
	err := config.DB.QueryRow(query, cons.Nombre, cons.Tipo).Scan(&cons.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al crear consultorio"})
	}
	return c.Status(201).JSON(cons)
}

// GET /api/consultorios
func ObtenerConsultorios(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT id_consultorio, nombre, tipo FROM Consultorios")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener consultorios"})
	}
	defer rows.Close()

	var lista []models.Consultorio
	for rows.Next() {
		var cons models.Consultorio
		if err := rows.Scan(&cons.ID, &cons.Nombre, &cons.Tipo); err == nil {
			lista = append(lista, cons)
		}
	}
	return c.JSON(lista)
}

// POST /api/consultorios/get
func ObtenerConsultorioPorID(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_consultorio"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	var cons models.Consultorio
	err := config.DB.QueryRow(
		"SELECT id_consultorio, nombre, tipo FROM Consultorios WHERE id_consultorio = $1",
		body.ID).Scan(&cons.ID, &cons.Nombre, &cons.Tipo)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Consultorio no encontrado"})
	}
	return c.JSON(cons)
}

// PUT /api/consultorios/update
func ActualizarConsultorio(c *fiber.Ctx) error {
	var cons models.Consultorio
	if err := c.BodyParser(&cons); err != nil || cons.ID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// Obtener datos actuales
	var actual models.Consultorio
	err := config.DB.QueryRow("SELECT nombre, tipo FROM Consultorios WHERE id_consultorio = $1", cons.ID).
		Scan(&actual.Nombre, &actual.Tipo)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Consultorio no encontrado"})
	}

	// Usar los datos nuevos si vienen, si no, mantener los existentes
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
		return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar consultorio"})
	}
	return c.JSON(fiber.Map{"mensaje": "Consultorio actualizado"})
}


// DELETE /api/consultorios/delete
func EliminarConsultorio(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_consultorio"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}
	_, err := config.DB.Exec("DELETE FROM Consultorios WHERE id_consultorio=$1", body.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar consultorio"})
	}
	return c.JSON(fiber.Map{"mensaje": "Consultorio eliminado"})
}
