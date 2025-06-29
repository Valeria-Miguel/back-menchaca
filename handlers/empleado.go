package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"back-menchaca/config"
	"back-menchaca/models"
	"back-menchaca/utils"
)

func CrearEmpleado(c *fiber.Ctx) error {
	var e models.Empleado

	if err := c.BodyParser(&e); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	if err := utils.ValidarEmpleado(e.Nombre, e.Appaterno, e.Tipo, e.Area, e.Correo, e.Contrasena); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
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
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error al verificar correo"})
	}
	if count > 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "El correo ya está registrado"})
	}

	hashed, err := utils.HashPassword(e.Contrasena)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error al encriptar contraseña"})
	}

	err = config.DB.QueryRow(
		`INSERT INTO Empleado (nombre, appaterno, apmaterno, tipo_empleado, area, correo, contraseña)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id_empleado`,
		e.Nombre, e.Appaterno, e.Apmaterno, e.Tipo, e.Area, e.Correo, hashed,
	).Scan(&e.ID)
	
	if err != nil {
		
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
        "error": "Error al registrar empleado: " + err.Error(),
    })
	}
	e.Contrasena = ""
	return c.Status(http.StatusCreated).JSON(e)
}

func ObtenerEmpleados(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT id_empleado, nombre, appaterno, apmaterno, tipo_empleado, area, correo FROM Empleado")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener empleados"})
	}
	defer rows.Close()

	empleados := []models.Empleado{}
	for rows.Next() {
		var e models.Empleado
		if err := rows.Scan(&e.ID, &e.Nombre, &e.Appaterno, &e.Apmaterno, &e.Tipo, &e.Area, &e.Correo); err == nil {
			empleados = append(empleados, e)
		}
	}
	return c.JSON(empleados)
}

func ObtenerEmpleadoPorID(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_empleado"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var e models.Empleado
	err := config.DB.QueryRow(
		"SELECT id_empleado, nombre, appaterno, apmaterno, tipo_empleado, area, correo FROM Empleado WHERE id_empleado = $1",
		body.ID,
	).Scan(&e.ID, &e.Nombre, &e.Appaterno, &e.Apmaterno, &e.Tipo, &e.Area, &e.Correo)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Empleado no encontrado"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error al buscar empleado"})
	}

	return c.JSON(e)
}

func ActualizarEmpleado(c *fiber.Ctx) error {
	var e models.Empleado
	if err := c.BodyParser(&e); err != nil || e.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// Leer datos actuales
	var current models.Empleado
	err := config.DB.QueryRow(
		"SELECT nombre, appaterno, apmaterno, tipo_empleado, area, correo FROM Empleado WHERE id_empleado = $1", e.ID,
	).Scan(&current.Nombre, &current.Appaterno, &current.Apmaterno, &current.Tipo, &current.Area, &current.Correo)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Empleado no encontrado"})
	}

	// Para cada campo: si el enviado está vacío, usar el actual
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

	// Sanitizar antes de guardar
	e.Nombre = utils.SanitizarInput(e.Nombre)
	e.Appaterno = utils.SanitizarInput(e.Appaterno)
	e.Apmaterno = utils.SanitizarInput(e.Apmaterno)
	e.Tipo = utils.SanitizarInput(e.Tipo)
	e.Area = utils.SanitizarInput(e.Area)
	e.Correo = utils.SanitizarInput(strings.ToLower(e.Correo))

	// Ejecutar actualización
	_, err = config.DB.Exec(
		`UPDATE Empleado SET nombre=$1, appaterno=$2, apmaterno=$3, tipo_empleado=$4, area=$5, correo=$6
		 WHERE id_empleado=$7`,
		e.Nombre, e.Appaterno, e.Apmaterno, e.Tipo, e.Area, e.Correo, e.ID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar empleado"})
	}

	return c.JSON(fiber.Map{"mensaje": "Empleado actualizado"})
}


func EliminarEmpleado(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_empleado"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	_, err := config.DB.Exec("DELETE FROM Empleado WHERE id_empleado = $1", body.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error al eliminar empleado"})
	}
	return c.JSON(fiber.Map{"mensaje": "Empleado eliminado"})
}
