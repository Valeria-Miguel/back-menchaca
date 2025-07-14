package handlers

import (
	"database/sql"

	"strings"
	"github.com/gofiber/fiber/v2"
	"back-menchaca/config"
	"back-menchaca/models"
	"back-menchaca/utils"
)

const modEmpl = "EMPL"
func CrearEmpleado(c *fiber.Ctx) error {
	var e models.Empleado

	if err := c.BodyParser(&e); err != nil {
		return utils.Responder(c, "02", modEmpl, "empleado-service", nil, "Datos inválidos")
	}

	if err := utils.ValidarEmpleado(e.Nombre, e.Appaterno, e.Tipo, e.Area, e.Correo, e.Contrasena); err != nil {
		return utils.Responder(c, "02", modEmpl, "empleado-service", nil, err.Error())
	}

	e.Nombre = utils.SanitizarInput(e.Nombre)
	e.Appaterno = utils.SanitizarInput(e.Appaterno)
	e.Apmaterno = utils.SanitizarInput(e.Apmaterno)
	e.Tipo = utils.SanitizarInput(e.Tipo)
	e.Area = utils.SanitizarInput(e.Area)
	e.Correo = utils.SanitizarInput(strings.ToLower(e.Correo))

	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM Empleado WHERE correo = $1", e.Correo).Scan(&count)
	if err != nil {
		return utils.Responder(c, "06", modEmpl, "empleado-service", nil, "Error al verificar correo")
	}
	if count > 0 {
		return utils.Responder(c, "07", modEmpl, "empleado-service", nil, "El correo ya está registrado")
	}

	hashed, err := utils.HashPassword(e.Contrasena)
	if err != nil {
		return utils.Responder(c, "06", modEmpl, "empleado-service", nil, "Error al encriptar contraseña")
	}

	err = config.DB.QueryRow(
		`INSERT INTO Empleado (nombre, appaterno, apmaterno, tipo_empleado, area, correo, contraseña)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id_empleado`,
		e.Nombre, e.Appaterno, e.Apmaterno, e.Tipo, e.Area, e.Correo, hashed,
	).Scan(&e.ID)
	if err != nil {
		return utils.Responder(c, "06", modEmpl, "empleado-service", nil, "Error al registrar empleado: "+err.Error())
	}

	e.Contrasena = ""
	return utils.Responder(c, "01", modEmpl, "empleado-service", e)
}

func ObtenerEmpleados(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT id_empleado, nombre, appaterno, apmaterno, tipo_empleado, area, correo FROM Empleado")
	if err != nil {
		return utils.Responder(c, "06", modEmpl, "empleado-service", nil, "Error al obtener empleados")
	}
	defer rows.Close()

	empleados := []models.Empleado{}
	for rows.Next() {
		var e models.Empleado
		if err := rows.Scan(&e.ID, &e.Nombre, &e.Appaterno, &e.Apmaterno, &e.Tipo, &e.Area, &e.Correo); err == nil {
			empleados = append(empleados, e)
		}
	}
	return utils.Responder(c, "01", modEmpl, "empleado-service", empleados)
}

func ObtenerEmpleadoPorID(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_empleado"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", modEmpl, "empleado-service", nil, "ID inválido")
	}

	var e models.Empleado
	err := config.DB.QueryRow(
		"SELECT id_empleado, nombre, appaterno, apmaterno, tipo_empleado, area, correo FROM Empleado WHERE id_empleado = $1",
		body.ID,
	).Scan(&e.ID, &e.Nombre, &e.Appaterno, &e.Apmaterno, &e.Tipo, &e.Area, &e.Correo)

	if err != nil {
		if err == sql.ErrNoRows {
			return utils.Responder(c, "05", modEmpl, "empleado-service", nil, "Empleado no encontrado")
		}
		return utils.Responder(c, "06", modEmpl, "empleado-service", nil, "Error al buscar empleado")
	}
	return utils.Responder(c, "01", modEmpl, "empleado-service", e)
}

func ActualizarEmpleado(c *fiber.Ctx) error {
	var e models.Empleado
	if err := c.BodyParser(&e); err != nil || e.ID == 0 {
		return utils.Responder(c, "02", modEmpl, "empleado-service", nil, "Datos inválidos")
	}

	var current models.Empleado
	err := config.DB.QueryRow(
		"SELECT nombre, appaterno, apmaterno, tipo_empleado, area, correo FROM Empleado WHERE id_empleado = $1", e.ID,
	).Scan(&current.Nombre, &current.Appaterno, &current.Apmaterno, &current.Tipo, &current.Area, &current.Correo)
	if err != nil {
		return utils.Responder(c, "05", modEmpl, "empleado-service", nil, "Empleado no encontrado")
	}

	if e.Nombre == "" {
		e.Nombre = current.Nombre
	}
	if e.Appaterno == "" {
		e.Appaterno = current.Appaterno
	}
	if e.Apmaterno == "" {
		e.Apmaterno = current.Apmaterno
	}
	if e.Tipo == "" {
		e.Tipo = current.Tipo
	}
	if e.Area == "" {
		e.Area = current.Area
	}
	if e.Correo == "" {
		e.Correo = current.Correo
	}

	e.Nombre = utils.SanitizarInput(e.Nombre)
	e.Appaterno = utils.SanitizarInput(e.Appaterno)
	e.Apmaterno = utils.SanitizarInput(e.Apmaterno)
	e.Tipo = utils.SanitizarInput(e.Tipo)
	e.Area = utils.SanitizarInput(e.Area)
	e.Correo = utils.SanitizarInput(strings.ToLower(e.Correo))

	_, err = config.DB.Exec(
		`UPDATE Empleado SET nombre=$1, appaterno=$2, apmaterno=$3, tipo_empleado=$4, area=$5, correo=$6
		 WHERE id_empleado=$7`,
		e.Nombre, e.Appaterno, e.Apmaterno, e.Tipo, e.Area, e.Correo, e.ID,
	)
	if err != nil {
		return utils.Responder(c, "06", modEmpl, "empleado-service", nil, "Error al actualizar empleado")
	}

	return utils.Responder(c, "01", modEmpl, "empleado-service", fiber.Map{"mensaje": "Empleado actualizado"})
}


func EliminarEmpleado(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_empleado"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", modEmpl, "empleado-service", nil, "ID inválido")
	}

	_, err := config.DB.Exec("DELETE FROM Empleado WHERE id_empleado = $1", body.ID)
	if err != nil {
		return utils.Responder(c, "06", modEmpl, "empleado-service", nil, "Error al eliminar empleado: "+err.Error())
	}

	return utils.Responder(c, "01", modEmpl, "empleado-service", fiber.Map{"mensaje": "Empleado eliminado"})
}
