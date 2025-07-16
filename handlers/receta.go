package handlers

import (
	"back-menchaca/config"
	"back-menchaca/models"
	"back-menchaca/utils"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

const modRec = "REC"

func CrearReceta(c *fiber.Ctx) error {
	var r models.Receta
	if err := c.BodyParser(&r); err != nil {
		return utils.Responder(c, "02", modRec, "receta-service", nil, "Datos inválidos")
	}

	if err := utils.ValidarReceta(r.Medicamento, r.Dosis, r.IDConsultorio); err != nil {
		return utils.Responder(c, "02", modRec, "receta-service", nil, err.Error())
	}

	if r.Fecha.IsZero() {
		r.Fecha = time.Now()
	}

	query := `INSERT INTO Recetas (fecha, medicamento, dosis, id_consultorio)
			  VALUES ($1, $2, $3, $4) RETURNING id_receta`

	err := config.DB.QueryRow(query, r.Fecha, r.Medicamento, r.Dosis, r.IDConsultorio).Scan(&r.ID)
	if err != nil {
		return utils.Responder(c, "06", modRec, "receta-service", nil, "Error al crear receta")
	}

	return utils.Responder(c, "01", modRec, "receta-service", r)
}

func ObtenerRecetas(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT id_receta, fecha, medicamento, dosis, id_consultorio FROM Recetas")
	if err != nil {
		return utils.Responder(c, "06", modRec, "receta-service", nil, "Error al obtener recetas")
	}
	defer rows.Close()

	var recetas []models.Receta
	for rows.Next() {
		var r models.Receta
		if err := rows.Scan(&r.ID, &r.Fecha, &r.Medicamento, &r.Dosis, &r.IDConsultorio); err == nil {
			recetas = append(recetas, r)
		}
	}
	return utils.Responder(c, "01", modRec, "receta-service", recetas)
}

func ObtenerRecetaPorID(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_receta"`
	}

	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", modRec, "receta-service", nil, "ID inválido")
	}

	var r struct {
		ID           int
		Fecha        string
		Medicamento  string
		Dosis        string
		IDConsultorio int
		NombreConsultorio string
	}

	err := config.DB.QueryRow(
		`SELECT r.id_receta, r.fecha, r.medicamento, r.dosis, r.id_consultorio, c.nombre
		FROM Recetas r
		INNER JOIN Consultorios c ON r.id_consultorio = c.id_consultorio
		WHERE r.id_receta = $1`, 
		body.ID).Scan(&r.ID, &r.Fecha, &r.Medicamento, &r.Dosis, &r.IDConsultorio, &r.NombreConsultorio)

	if err != nil {
		if err == sql.ErrNoRows {
			return utils.Responder(c, "05", modRec, "receta-service", nil, "Receta no encontrada")
		}
		return utils.Responder(c, "06", modRec, "receta-service", nil, "Error al buscar receta")
	}

	return utils.Responder(c, "01", modRec, "receta-service", r)
}



func ActualizarReceta(c *fiber.Ctx) error {
	var r models.Receta
	if err := c.BodyParser(&r); err != nil || r.ID == 0 {
		return utils.Responder(c, "02", modRec, "receta-service", nil, "Datos inválidos")
	}

	var actual models.Receta
	err := config.DB.QueryRow("SELECT fecha, medicamento, dosis, id_consultorio FROM Recetas WHERE id_receta=$1", r.ID).
		Scan(&actual.Fecha, &actual.Medicamento, &actual.Dosis, &actual.IDConsultorio)
	if err == sql.ErrNoRows {
		return utils.Responder(c, "05", modRec, "receta-service", nil, "Receta no encontrada")
	} else if err != nil {
		return utils.Responder(c, "06", modRec, "receta-service", nil, "Error al consultar receta")
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
		return utils.Responder(c, "06", modRec, "receta-service", nil, "Error al actualizar receta")
	}

	return utils.Responder(c, "01", modRec, "receta-service", fiber.Map{"mensaje": "Receta actualizada"})
}

func EliminarReceta(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_receta"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", modRec, "receta-service", nil, "ID inválido")
	}

	_, err := config.DB.Exec("DELETE FROM Recetas WHERE id_receta=$1", body.ID)
	if err != nil {
		return utils.Responder(c, "06", modRec, "receta-service", nil, "Error al eliminar receta")
	}

	return utils.Responder(c, "01", modRec, "receta-service", fiber.Map{"mensaje": "Receta eliminada"})
}
