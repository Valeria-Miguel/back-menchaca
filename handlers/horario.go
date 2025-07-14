package handlers

import (
	"back-menchaca/config"
	"back-menchaca/models"
	"back-menchaca/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

const modHor = "HOR"

func CrearHorario(c *fiber.Ctx) error {
	var h models.Horario
	if err := c.BodyParser(&h); err != nil {
		return utils.Responder(c, "02", modHor, "horario-service", nil, "Datos inválidos")
	}

	if err := utils.ValidarHorario(h.Turno, h.IDEmpleado, h.IDConsultorio); err != nil {
		return utils.Responder(c, "02", modHor, "horario-service", nil, err.Error())
	}
	h.Turno = utils.SanitizarInput(strings.ToLower(h.Turno))

	var empExists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM Empleado WHERE id_empleado=$1)", h.IDEmpleado).Scan(&empExists)
	if err != nil || !empExists {
		return utils.Responder(c, "02", modHor, "horario-service", nil, "Empleado no encontrado")
	}

	var consExists bool
	err = config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM Consultorios WHERE id_consultorio=$1)", h.IDConsultorio).Scan(&consExists)
	if err != nil || !consExists {
		return utils.Responder(c, "02", modHor, "horario-service", nil, "Consultorio no encontrado")
	}

	query := `INSERT INTO Horarios (id_consultorio, turno, id_empleado) 
	          VALUES ($1, $2, $3) RETURNING id_horario`
	err = config.DB.QueryRow(query, h.IDConsultorio, h.Turno, h.IDEmpleado).Scan(&h.ID)
	if err != nil {
		return utils.Responder(c, "06", modHor, "horario-service", nil, "Error al crear horario: "+err.Error())
	}

	return utils.Responder(c, "01", modHor, "horario-service", h)
}

func ObtenerHorarios(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT id_horario, id_consultorio, turno, id_empleado FROM Horarios")
	if err != nil {
		return utils.Responder(c, "06", modHor, "horario-service", nil, "Error al obtener horarios")
	}
	defer rows.Close()

	var lista []models.Horario
	for rows.Next() {
		var h models.Horario
		if err := rows.Scan(&h.ID, &h.IDConsultorio, &h.Turno, &h.IDEmpleado); err == nil {
			lista = append(lista, h)
		}
	}
	return utils.Responder(c, "01", modHor, "horario-service", lista)
}

func ObtenerHorarioPorID(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_horario"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", modHor, "horario-service", nil, "ID inválido")
	}

	var h models.Horario
	err := config.DB.QueryRow(
		"SELECT id_horario, id_consultorio, turno, id_empleado FROM Horarios WHERE id_horario = $1",
		body.ID).Scan(&h.ID, &h.IDConsultorio, &h.Turno, &h.IDEmpleado)

	if err != nil {
		return utils.Responder(c, "05", modHor, "horario-service", nil, "Horario no encontrado")
	}
	return utils.Responder(c, "01", modHor, "horario-service", h)
}

func ActualizarHorario(c *fiber.Ctx) error {
	var h models.Horario
	if err := c.BodyParser(&h); err != nil || h.ID == 0 {
		return utils.Responder(c, "02", modHor, "horario-service", nil, "Datos inválidos")
	}

	var actual models.Horario
	err := config.DB.QueryRow(
		"SELECT id_consultorio, turno, id_empleado FROM Horarios WHERE id_horario=$1", h.ID,
	).Scan(&actual.IDConsultorio, &actual.Turno, &actual.IDEmpleado)
	if err != nil {
		return utils.Responder(c, "05", modHor, "horario-service", nil, "Horario no encontrado")
	}

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

	var consExists bool
	err = config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM Consultorios WHERE id_consultorio=$1)", h.IDConsultorio).Scan(&consExists)
	if err != nil || !consExists {
		return utils.Responder(c, "02", modHor, "horario-service", nil, "Consultorio no existe")
	}

	var empExists bool
	err = config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM Empleado WHERE id_empleado=$1)", h.IDEmpleado).Scan(&empExists)
	if err != nil || !empExists {
		return utils.Responder(c, "02", modHor, "horario-service", nil, "Empleado no existe")
	}

	_, err = config.DB.Exec(
		"UPDATE Horarios SET id_consultorio=$1, turno=$2, id_empleado=$3 WHERE id_horario=$4",
		h.IDConsultorio, h.Turno, h.IDEmpleado, h.ID,
	)
	if err != nil {
		return utils.Responder(c, "06", modHor, "horario-service", nil, "Error al actualizar horario: "+err.Error())
	}

	return utils.Responder(c, "01", modHor, "horario-service", fiber.Map{"mensaje": "Horario actualizado correctamente"})
}


func EliminarHorario(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_horario"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", modHor, "horario-service", nil, "ID inválido")
	}

	_, err := config.DB.Exec("DELETE FROM Horarios WHERE id_horario=$1", body.ID)
	if err != nil {
		return utils.Responder(c, "06", modHor, "horario-service", nil, "Error al eliminar horario")
	}
	return utils.Responder(c, "01", modHor, "horario-service", fiber.Map{"mensaje": "Horario eliminado"})
}
