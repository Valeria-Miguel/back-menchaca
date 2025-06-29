package handlers

import (
	"back-menchaca/config"
	"back-menchaca/models"
	"back-menchaca/utils"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CrearReceta(c *fiber.Ctx) error {
	var r models.Receta
	if err := c.BodyParser(&r); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	if err := utils.ValidarReceta(r.Medicamento, r.Dosis, r.IDConsultorio); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if r.Fecha.IsZero() {
		r.Fecha = time.Now()
	}

	query := `INSERT INTO Recetas (fecha, medicamento, dosis, id_consultorio)
			  VALUES ($1, $2, $3, $4) RETURNING id_receta`

	err := config.DB.QueryRow(query, r.Fecha, r.Medicamento, r.Dosis, r.IDConsultorio).Scan(&r.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al crear receta"})
	}

	return c.Status(201).JSON(r)
}

func ObtenerRecetas(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT id_receta, fecha, medicamento, dosis, id_consultorio FROM Recetas")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener recetas"})
	}
	defer rows.Close()

	var recetas []models.Receta
	for rows.Next() {
		var r models.Receta
		if err := rows.Scan(&r.ID, &r.Fecha, &r.Medicamento, &r.Dosis, &r.IDConsultorio); err == nil {
			recetas = append(recetas, r)
		}
	}
	return c.JSON(recetas)
}

func ObtenerRecetaPorID(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_receta"`
	}

	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	var r models.Receta
	err := config.DB.QueryRow(
		"SELECT id_receta, fecha, medicamento, dosis, id_consultorio FROM Recetas WHERE id_receta = $1",
		body.ID).Scan(&r.ID, &r.Fecha, &r.Medicamento, &r.Dosis, &r.IDConsultorio)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Receta no encontrada"})
	}

	return c.JSON(r)
}


func ActualizarReceta(c *fiber.Ctx) error {
	var r models.Receta
	if err := c.BodyParser(&r); err != nil || r.ID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	var actual models.Receta
	err := config.DB.QueryRow("SELECT fecha, medicamento, dosis, id_consultorio FROM Recetas WHERE id_receta=$1", r.ID).
		Scan(&actual.Fecha, &actual.Medicamento, &actual.Dosis, &actual.IDConsultorio)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{"error": "Receta no encontrada"})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al consultar receta"})
	}

	if r.Fecha.IsZero() {
		r.Fecha = actual.Fecha
	}
	if r.Medicamento == "" {
		r.Medicamento = actual.Medicamento
	}
	if r.Dosis == "" {
		r.Dosis = actual.Dosis
	}
	if r.IDConsultorio == 0 {
		r.IDConsultorio = actual.IDConsultorio
	}

	_, err = config.DB.Exec(`UPDATE Recetas SET fecha=$1, medicamento=$2, dosis=$3, id_consultorio=$4 WHERE id_receta=$5`,
		r.Fecha, r.Medicamento, r.Dosis, r.IDConsultorio, r.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar receta"})
	}

	return c.JSON(fiber.Map{"mensaje": "Receta actualizada"})
}

func EliminarReceta(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_receta"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	_, err := config.DB.Exec("DELETE FROM Recetas WHERE id_receta=$1", body.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar receta"})
	}
	return c.JSON(fiber.Map{"mensaje": "Receta eliminada"})
}
