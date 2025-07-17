package handlers

import (
	"back-menchaca/config"
	"back-menchaca/models"
	"back-menchaca/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
	"database/sql"
	"log"
)


const modExp = "EXP"

func CrearExpediente(c *fiber.Ctx) error {
	var e models.Expediente
	if err := c.BodyParser(&e); err != nil {
		return utils.Responder(c, "02", modExp, "expediente-service", nil, "Datos inválidos")
	}

	if !utils.ExisteID("Paciente", "id_paciente", e.IDPaciente) {
		return utils.Responder(c, "02", modExp, "expediente-service", nil, "ID de paciente no válido")
	}

	if err := utils.ValidarSeguro(e.Seguro); err != nil {
		return utils.Responder(c, "02", modExp, "expediente-service", nil, err.Error())
	}

	if e.FechaCreacion.IsZero() {
		e.FechaCreacion = time.Now()
	}

	query := `INSERT INTO Expediente (id_paciente, seguro, fecha_creacion) VALUES ($1, $2, $3) RETURNING id_expediente`
	err := config.DB.QueryRow(query, e.IDPaciente, e.Seguro, e.FechaCreacion).Scan(&e.ID)
	if err != nil {
		return utils.Responder(c, "06", modExp, "expediente-service", nil, "Error al crear expediente")
	}

	return utils.Responder(c, "01", modExp, "expediente-service", e)
}



func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func formatDate(t sql.NullTime) string {
	if t.Valid {
		return t.Time.Format("2006-01-02")
	}
	return ""
}

func ObtenerExpedientes(c *fiber.Ctx) error {
	type Antecedente struct {
		Tiene       string `json:"tiene"` // Puedes omitir o calcular esto si quieres
		Diagnostico string `json:"diagnostico"`
		Descripcion string `json:"descripcion"`
	}

	type ExpedienteDetallado struct {
		IDExpediente  int    `json:"id_expediente"`
		Seguro        string `json:"seguro"`
		FechaCreacion string `json:"fecha_creacion"`
		Paciente      struct {
			ID        int    `json:"id_paciente"`
			Nombre    string `json:"nombre"`
			Appaterno string `json:"appaterno"`
			Apmaterno string `json:"apmaterno"`
		} `json:"paciente"`
		Antecedentes []Antecedente `json:"antecedentes"`
	}

	var expedientes []ExpedienteDetallado

	rows, err := config.DB.Query(`
		SELECT e.id_expediente, e.id_paciente, e.seguro, e.fecha_creacion,
		       p.nombre, p.appaterno, p.apmaterno
		FROM Expediente e
		LEFT JOIN Paciente p ON e.id_paciente = p.id_paciente
	`)
	if err != nil {
		return utils.Responder(c, "06", modExp, "expediente-service", nil, "Error al obtener expedientes")
	}
	defer rows.Close()

	for rows.Next() {
		var e ExpedienteDetallado
		var seguro sql.NullString
		var fechaCreacion sql.NullTime
		var nombre sql.NullString
		var appaterno sql.NullString
		var apmaterno sql.NullString

		if err := rows.Scan(
			&e.IDExpediente,
			&e.Paciente.ID,
			&seguro,
			&fechaCreacion,
			&nombre,
			&appaterno,
			&apmaterno,
		); err != nil {
			log.Println("❌ Error al escanear expediente:", err)
			continue
		}

		// Convertir nulos
		e.Seguro = nullStringToString(seguro)
		e.FechaCreacion = formatDate(fechaCreacion)
		e.Paciente.Nombre = nullStringToString(nombre)
		e.Paciente.Appaterno = nullStringToString(appaterno)
		e.Paciente.Apmaterno = nullStringToString(apmaterno)

		// Obtener antecedentes
		antRows, err := config.DB.Query(`
			SELECT diagnostico, descripcion
			FROM Antecedentes
			WHERE id_expediente = $1
		`, e.IDExpediente)
		if err == nil {
			defer antRows.Close()
			for antRows.Next() {
				var ant Antecedente
				var diag sql.NullString
				var desc sql.NullString

				if err := antRows.Scan(&diag, &desc); err != nil {
					log.Println("❌ Error al escanear antecedente:", err)
					continue
				}

				ant.Diagnostico = nullStringToString(diag)
				ant.Descripcion = nullStringToString(desc)
				e.Antecedentes = append(e.Antecedentes, ant)
			}
		}

		expedientes = append(expedientes, e)
	}

	return utils.Responder(c, "01", modExp, "expediente-service", expedientes)
}


func ObtenerExpedientePorID(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_expediente"`
	}

	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		log.Printf("❌ Error al parsear body o ID inválido: %v", err)
		return utils.Responder(c, "02", modExp, "expediente-service", nil, "ID inválido")
	}

	if !utils.ExisteIDExped(body.ID) {
		return utils.Responder(c, "05", modExp, "expediente-service", nil, "Expediente no encontrado")
	}

	// Struct para antecedentes
	type Antecedente struct {
		Tiene       string `json:"tiene"`
		Diagnostico string `json:"diagnostico"`
		Descripcion string `json:"descripcion"`
	}

	// Struct para respuesta detallada
	type ExpedienteDetallado struct {
		IDExpediente int    `json:"id_expediente"`
		Seguro       string `json:"seguro"`
		FechaCreacion string `json:"fecha_creacion"`
		Paciente      struct {
			ID        int    `json:"id_paciente"`
			Nombre    string `json:"nombre"`
			Appaterno string `json:"appaterno"`
			Apmaterno string `json:"apmaterno"`
		} `json:"paciente"`
		Antecedentes []Antecedente `json:"antecedentes"`
	}

	var exp ExpedienteDetallado
	var fecha sql.NullString // <- Cambio aquí

	// Consulta
	err := config.DB.QueryRow(`
		SELECT e.id_expediente, e.id_paciente, e.seguro, e.fecha_creacion,
		       COALESCE(p.nombre, '') AS nombre,
		       COALESCE(p.appaterno, '') AS appaterno,
		       COALESCE(p.apmaterno, '') AS apmaterno
		FROM Expediente e
		LEFT JOIN Paciente p ON e.id_paciente = p.id_paciente
		WHERE e.id_expediente = $1
	`, body.ID).Scan(
		&exp.IDExpediente,
		&exp.Paciente.ID,
		&exp.Seguro,
		&fecha, // <- aquí escaneamos en NullString
		&exp.Paciente.Nombre,
		&exp.Paciente.Appaterno,
		&exp.Paciente.Apmaterno,
	)
	if err != nil {
		log.Printf("❌ Error al obtener expediente y paciente: %v", err)
		return utils.Responder(c, "06", modExp, "expediente-service", nil, "Error al buscar expediente")
	}

	// Asignar la fecha (con manejo de NULL)
	exp.FechaCreacion = fecha.String

	// Obtener antecedentes
	rows, err := config.DB.Query(`
		SELECT diagnostico, descripcion
		FROM Antecedentes
		WHERE id_expediente = $1
	`, body.ID)
	if err != nil {
		log.Printf("⚠️ Error al consultar antecedentes: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var diagnostico string
			var descripcion sql.NullString
			if err := rows.Scan(&diagnostico, &descripcion); err == nil {
				exp.Antecedentes = append(exp.Antecedentes, Antecedente{
					Tiene:       "Sí",
					Diagnostico: diagnostico,
					Descripcion: descripcion.String,
				})
			}
		}
	}

	return utils.Responder(c, "01", modExp, "expediente-service", exp)
}


func ActualizarExpediente(c *fiber.Ctx) error {
	var e models.Expediente
	if err := c.BodyParser(&e); err != nil || e.ID == 0 {
		return utils.Responder(c, "02", modExp, "expediente-service", nil, "Datos inválidos")
	}

	if !utils.ExisteIDExped(e.ID) {
		return utils.Responder(c, "05", modExp, "expediente-service", nil, "Expediente no encontrado")
	}

	var actual models.Expediente
	err := config.DB.QueryRow(`SELECT id_paciente, seguro, fecha_creacion FROM Expediente WHERE id_expediente=$1`, e.ID).
		Scan(&actual.IDPaciente, &actual.Seguro, &actual.FechaCreacion)
	if err != nil {
		return utils.Responder(c, "06", modExp, "expediente-service", nil, "Error al obtener expediente actual")
	}

	if e.IDPaciente == 0 {
		e.IDPaciente = actual.IDPaciente
	} else if !utils.ExisteID("Paciente", "id_paciente", e.IDPaciente) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de paciente no válido"})
	}

	if strings.TrimSpace(e.Seguro) == "" {
		e.Seguro = actual.Seguro
	} else {
		if err := utils.ValidarSeguro(e.Seguro); err != nil {
			return utils.Responder(c, "02", modExp, "expediente-service", nil, err.Error())
		}
	}

	if e.FechaCreacion.IsZero() {
		e.FechaCreacion = actual.FechaCreacion
	}

	_, err = config.DB.Exec(`UPDATE Expediente SET id_paciente=$1, seguro=$2, fecha_creacion=$3 WHERE id_expediente=$4`,
		e.IDPaciente, e.Seguro, e.FechaCreacion, e.ID)
	if err != nil {
		return utils.Responder(c, "06", modExp, "expediente-service", nil, "Error al actualizar expediente")
	}
	return utils.Responder(c, "01", modExp, "expediente-service", fiber.Map{"mensaje": "Expediente actualizado"})
}

func EliminarExpediente(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_expediente"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", modExp, "expediente-service", nil, "ID inválido")
	}

	if !utils.ExisteIDExped(body.ID) {
		return utils.Responder(c, "05", modExp, "expediente-service", nil, "Expediente no encontrado")
	}

	_, err := config.DB.Exec("DELETE FROM Expediente WHERE id_expediente=$1", body.ID)
	if err != nil {
		return utils.Responder(c, "06", modExp, "expediente-service", nil, "Error al eliminar expediente")
	}
	return utils.Responder(c, "01", modExp, "expediente-service", fiber.Map{"mensaje": "Expediente eliminado"})
}
