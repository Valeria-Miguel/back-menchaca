package handlers

import (
	"back-menchaca/config"
	"back-menchaca/models"
	"back-menchaca/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// POST /api/horarios
func CrearHorario(c *fiber.Ctx) error {
	var h models.Horario
	if err := c.BodyParser(&h); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// Validar formato
	if err := utils.ValidarHorario(h.Turno, h.IDEmpleado, h.IDConsultorio); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	h.Turno = utils.SanitizarInput(strings.ToLower(h.Turno))

	// Verificar existencia de empleado
	var empExists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM Empleado WHERE id_empleado=$1)", h.IDEmpleado).Scan(&empExists)
	if err != nil || !empExists {
		return c.Status(400).JSON(fiber.Map{"error": "El empleado no existe"})
	}

	// Verificar existencia de consultorio
	var consExists bool
	err = config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM Consultorios WHERE id_consultorio=$1)", h.IDConsultorio).Scan(&consExists)
	if err != nil || !consExists {
		return c.Status(400).JSON(fiber.Map{"error": "El consultorio no existe"})
	}

	// Insertar
	query := `INSERT INTO Horarios (id_consultorio, turno, id_empleado) 
	          VALUES ($1, $2, $3) RETURNING id_horario`
	err = config.DB.QueryRow(query, h.IDConsultorio, h.Turno, h.IDEmpleado).Scan(&h.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al crear horario: " + err.Error()})
	}
	return c.Status(201).JSON(h)
}

// GET /api/horarios
func ObtenerHorarios(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT id_horario, id_consultorio, turno, id_empleado FROM Horarios")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener horarios"})
	}
	defer rows.Close()

	var lista []models.Horario
	for rows.Next() {
		var h models.Horario
		if err := rows.Scan(&h.ID, &h.IDConsultorio, &h.Turno, &h.IDEmpleado); err == nil {
			lista = append(lista, h)
		}
	}
	return c.JSON(lista)
}

// POST /api/horarios/get
func ObtenerHorarioPorID(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_horario"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	var h models.Horario
	err := config.DB.QueryRow(
		"SELECT id_horario, id_consultorio, turno, id_empleado FROM Horarios WHERE id_horario = $1",
		body.ID).Scan(&h.ID, &h.IDConsultorio, &h.Turno, &h.IDEmpleado)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Horario no encontrado"})
	}
	return c.JSON(h)
}

// PUT /api/horarios/update
func ActualizarHorario(c *fiber.Ctx) error {
	var h models.Horario
	if err := c.BodyParser(&h); err != nil || h.ID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// Obtener datos actuales
	var actual models.Horario
	err := config.DB.QueryRow(
		"SELECT id_consultorio, turno, id_empleado FROM Horarios WHERE id_horario=$1", h.ID,
	).Scan(&actual.IDConsultorio, &actual.Turno, &actual.IDEmpleado)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Horario no encontrado"})
	}

	// Conservar los campos no enviados
	if h.IDConsultorio == 0 {
		h.IDConsultorio = actual.IDConsultorio
	}
	if h.Turno == "" {
		h.Turno = actual.Turno
	} else {
		h.Turno = utils.SanitizarInput(strings.ToLower(h.Turno))
	}
	if h.IDEmpleado == 0 {
		h.IDEmpleado = actual.IDEmpleado
	}

	// Verificar existencia del consultorio
	var consExists bool
	err = config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM Consultorios WHERE id_consultorio=$1)", h.IDConsultorio).Scan(&consExists)
	if err != nil || !consExists {
		return c.Status(400).JSON(fiber.Map{"error": "El consultorio no existe"})
	}

	// Verificar existencia del empleado
	var empExists bool
	err = config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM Empleado WHERE id_empleado=$1)", h.IDEmpleado).Scan(&empExists)
	if err != nil || !empExists {
		return c.Status(400).JSON(fiber.Map{"error": "El empleado no existe"})
	}

	// Ejecutar actualización
	_, err = config.DB.Exec(
		"UPDATE Horarios SET id_consultorio=$1, turno=$2, id_empleado=$3 WHERE id_horario=$4",
		h.IDConsultorio, h.Turno, h.IDEmpleado, h.ID,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar horario: " + err.Error()})
	}

	return c.JSON(fiber.Map{"mensaje": "Horario actualizado correctamente"})
}


// DELETE /api/horarios/delete
func EliminarHorario(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_horario"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	_, err := config.DB.Exec("DELETE FROM Horarios WHERE id_horario=$1", body.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar horario"})
	}
	return c.JSON(fiber.Map{"mensaje": "Horario eliminado"})
}
