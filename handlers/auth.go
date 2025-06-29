package handlers

import (
	"database/sql"
	"strings"

	"github.com/gofiber/fiber/v2"
	"back-menchaca/config"
	"back-menchaca/utils"
)

func Login(c *fiber.Ctx) error {
	var input struct {
		Correo   string `json:"correo"`
		Password string `json:"contrasena"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Petición inválida"})
	}

	input.Correo = strings.ToLower(input.Correo)

	var (
		hashedPassword string
		id             int
		rol            string
		err            error
	)

	err = config.DB.QueryRow(`SELECT id_empleado, contraseña FROM Empleado WHERE correo=$1`, input.Correo).
		Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		err = config.DB.QueryRow(`SELECT id_paciente, contraseña FROM Paciente WHERE correo=$1`, input.Correo).
			Scan(&id, &hashedPassword)
		if err == sql.ErrNoRows {
			return c.Status(401).JSON(fiber.Map{"error": "Credenciales inválidas"})
		} else if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error en base de datos"})
		}
		rol = "paciente"
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error en base de datos"})
	} else {
		rol = "empleado"
	}

	if !utils.CheckPasswordHash(input.Password, hashedPassword) {
		return c.Status(401).JSON(fiber.Map{"error": "Contraseña incorrecta"})
	}

	token, err := utils.GenerateJWT(input.Correo, rol)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error generando token"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"rol":   rol,
	})
}
