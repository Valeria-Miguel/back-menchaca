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

func CrearPaciente(c *fiber.Ctx) error {
	var p models.Paciente

	if err := c.BodyParser(&p); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Datos inválidos",
		})
	}

	if err := utils.ValidarPaciente(p.Nombre, p.Appaterno, p.Correo, p.Contrasena); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	p.Nombre     = utils.SanitizarInput(p.Nombre)
	p.Appaterno  = utils.SanitizarInput(p.Appaterno)
	p.Apmaterno  = utils.SanitizarInput(p.Apmaterno)
	p.Correo     = utils.SanitizarInput(strings.ToLower(p.Correo)) // estandarizar y sanitizar

	p.Correo = strings.ToLower(p.Correo)

	var count int
	err := config.DB.QueryRow(
		"SELECT COUNT(*) FROM Paciente WHERE correo = $1", p.Correo,
	).Scan(&count)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al verificar correo en pacientes",
		})
	}
	if count > 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "El correo ya está registrado como paciente",
		})
	}

	err = config.DB.QueryRow(
		"SELECT COUNT(*) FROM Empleado WHERE correo = $1", p.Correo,
	).Scan(&count)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al verificar correo en empleados",
		})
	}
	if count > 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "El correo ya está registrado como empleado",
		})
	}

	hashed, err := utils.HashPassword(p.Contrasena)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al proteger la contraseña",
		})
	}

	query := `INSERT INTO Paciente (nombre, appaterno, apmaterno, correo, contraseña) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id_paciente`
	err = config.DB.QueryRow(query, p.Nombre, p.Appaterno, p.Apmaterno, p.Correo, hashed).Scan(&p.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al crear paciente",
		})
	}

	p.Contrasena = ""
	return c.Status(http.StatusCreated).JSON(p)
}



func ObtenerPacientes(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT id_paciente, nombre, appaterno, apmaterno, correo FROM Paciente")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener pacientes"})
	}
	defer rows.Close()

	var pacientes []models.Paciente
	for rows.Next() {
		var p models.Paciente
		if err := rows.Scan(&p.ID, &p.Nombre, &p.Appaterno, &p.Apmaterno, &p.Correo); err != nil {
			continue
		}
		pacientes = append(pacientes, p)
	}

	return c.JSON(pacientes)
}

func ObtenerPacientePorID(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_paciente"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var p models.Paciente
	err := config.DB.QueryRow("SELECT id_paciente, nombre, appaterno, apmaterno, correo FROM Paciente WHERE id_paciente = $1", body.ID).
		Scan(&p.ID, &p.Nombre, &p.Appaterno, &p.Apmaterno, &p.Correo)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Paciente no encontrado"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error al buscar paciente"})
	}

	return c.JSON(p)
}

func ActualizarPaciente(c *fiber.Ctx) error {
	var p models.Paciente
	if err := c.BodyParser(&p); err != nil || p.ID == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	var current models.Paciente
	err := config.DB.QueryRow(
		"SELECT nombre, appaterno, apmaterno, correo FROM Paciente WHERE id_paciente = $1", p.ID,
	).Scan(&current.Nombre, &current.Appaterno, &current.Apmaterno, &current.Correo)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Paciente no encontrado"})
	}

	//mantener los campos no enviados
	if p.Nombre != "" && !utils.ValidarTextoLetras(p.Nombre) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Nombre inválido"})
	}
	if p.Appaterno != "" && !utils.ValidarTextoLetras(p.Appaterno) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Apellido paterno inválido"})
	}
	if p.Apmaterno != "" && !utils.ValidarTextoLetras(p.Apmaterno) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Apellido materno inválido"})
	}
	if p.Correo != "" && !utils.ValidarCorreo(p.Correo) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Correo inválido"})
	}



	p.Nombre = utils.SanitizarInput(p.Nombre)
	p.Appaterno = utils.SanitizarInput(p.Appaterno)
	p.Apmaterno = utils.SanitizarInput(p.Apmaterno)
	p.Correo = utils.SanitizarInput(strings.ToLower(p.Correo))

	query := `UPDATE Paciente 
	          SET nombre=$1, appaterno=$2, apmaterno=$3, correo=$4 
	          WHERE id_paciente=$5`
	_, err = config.DB.Exec(query, p.Nombre, p.Appaterno, p.Apmaterno, p.Correo, p.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar paciente"})
	}

	return c.JSON(fiber.Map{"mensaje": "Paciente actualizado"})
}

func EliminarPaciente(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_paciente"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	_, err := config.DB.Exec("DELETE FROM Paciente WHERE id_paciente = $1", body.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error al eliminar paciente"})
	}

	return c.JSON(fiber.Map{"mensaje": "Paciente eliminado"})
}
