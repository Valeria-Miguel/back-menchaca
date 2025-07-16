package handlers

import (
	"back-menchaca/config"
	"back-menchaca/models"
	"back-menchaca/utils"
	"database/sql"
	 "fmt"
	"github.com/gofiber/fiber/v2"
)

const modConsul = "Consul"

func AgendarConsulta(c *fiber.Ctx) error {
	var cons models.Consulta
	if err := c.BodyParser(&cons); err != nil {
		return utils.Responder(c, "02", modConsul, "consulta-service", nil, "Datos inválidos")
	}

	if err := utils.ValidarConsulta(cons); err != nil {
		return utils.Responder(c, "02", modConsul, "consulta-service", nil, err.Error())
	}

	if !utils.ExisteID("Paciente", "id_paciente", cons.IDPaciente) {
		return utils.Responder(c, "02", modConsul, "consulta-service", nil, "ID de paciente no válido")
	}
	if !utils.ExisteID("Horarios", "id_horario", cons.IDHorario) {
		return utils.Responder(c, "02", modConsul, "consulta-service", nil, "ID de horario no válido")
	}
	if !utils.ExisteID("Consultorios", "id_consultorio", cons.IDConsultorio) {
		return utils.Responder(c, "02", modConsul, "consulta-service", nil, "ID de consultorio no válido")
	}

	err := config.DB.QueryRow(`
		INSERT INTO Consultas (id_paciente, tipo, id_receta, id_horario, id_consultorio, diagnostico, costo, fecha_hora)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id_consulta`,
		cons.IDPaciente, cons.Tipo, cons.IDReceta, cons.IDHorario, cons.IDConsultorio, cons.Diagnostico, cons.Costo, cons.FechaHora,
	).Scan(&cons.ID)

	if err != nil {
		return utils.Responder(c, "06", modConsul, "consulta-service", nil, "Error al agendar consulta: "+err.Error())
	}

	return utils.Responder(c, "01", modConsul, "consulta-service", cons)
}

func ObtenerConsultas(c *fiber.Ctx) error {
	rows, err := config.DB.Query(`SELECT id_consulta, id_paciente, tipo, id_receta, id_horario, id_consultorio, diagnostico, costo, fecha_hora FROM Consultas`) 
	if err != nil {
		return utils.Responder(c, "06", modConsul, "consulta-service", nil, "Error al obtener consultas")
	}
	defer rows.Close()

	var consultas []models.Consulta
	for rows.Next() {
		var cons models.Consulta
		if err := rows.Scan(&cons.ID, &cons.IDPaciente, &cons.Tipo, &cons.IDReceta, &cons.IDHorario, &cons.IDConsultorio, &cons.Diagnostico, &cons.Costo, &cons.FechaHora); err == nil {
			consultas = append(consultas, cons)
		}
	}

	return utils.Responder(c, "01", modConsul, "consulta-service", consultas)
}

func ObtenerConsultaPorID(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_consulta"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", modConsul, "consulta-service", nil, "ID inválido")
	}

	var cons models.Consulta
	err := config.DB.QueryRow(
		`SELECT id_consulta, id_paciente, tipo, id_receta, id_horario, id_consultorio, diagnostico, costo, fecha_hora
		 FROM Consultas WHERE id_consulta = $1`, body.ID,
	).Scan(&cons.ID, &cons.IDPaciente, &cons.Tipo, &cons.IDReceta, &cons.IDHorario, &cons.IDConsultorio, &cons.Diagnostico, &cons.Costo, &cons.FechaHora)

	if err == sql.ErrNoRows {
		return utils.Responder(c, "05", modConsul, "consulta-service", nil, "Consulta no encontrada")
	} else if err != nil {
		return utils.Responder(c, "06", modConsul, "consulta-service", nil, "Error al buscar consulta")
	}

	return utils.Responder(c, "01", modConsul, "consulta-service", cons)
}

func ActualizarConsulta(c *fiber.Ctx) error {
	var cons models.Consulta
	if err := c.BodyParser(&cons); err != nil || cons.ID == 0 {
		return utils.Responder(c, "02", modConsul, "consulta-service", nil, "Datos inválidos")
	}

	var actual models.Consulta
	err := config.DB.QueryRow(`SELECT id_paciente, tipo, id_receta, id_horario, id_consultorio, diagnostico, costo, fecha_hora FROM Consultas WHERE id_consulta=$1`, cons.ID).
		Scan(&actual.IDPaciente, &actual.Tipo, &actual.IDReceta, &actual.IDHorario, &actual.IDConsultorio, &actual.Diagnostico, &actual.Costo, &actual.FechaHora)
	
	if err == sql.ErrNoRows {
		return utils.Responder(c, "05", modConsul, "consulta-service", nil, "Consulta no encontrada")
	} else if err != nil {
		return utils.Responder(c, "06", modConsul, "consulta-service", nil, "Error al buscar consulta")
	}

	if cons.IDPaciente == 0 {
		cons.IDPaciente = actual.IDPaciente
	}
	if cons.Tipo == "" {
		cons.Tipo = actual.Tipo
	}
	if cons.IDReceta == nil {
		cons.IDReceta = actual.IDReceta
	}
	if cons.IDHorario == 0 {
		cons.IDHorario = actual.IDHorario
	}
	if cons.IDConsultorio == 0 {
		cons.IDConsultorio = actual.IDConsultorio
	}
	if cons.Diagnostico == "" {
		cons.Diagnostico = actual.Diagnostico
	}
	if cons.Costo == 0 {
		cons.Costo = actual.Costo
	}
	if cons.FechaHora.IsZero() {
		cons.FechaHora = actual.FechaHora
	}

	_, err = config.DB.Exec(`UPDATE Consultas SET id_paciente=$1, tipo=$2, id_receta=$3, id_horario=$4, id_consultorio=$5, diagnostico=$6, costo=$7, fecha_hora=$8 WHERE id_consulta=$9`,
		cons.IDPaciente, cons.Tipo, cons.IDReceta, cons.IDHorario, cons.IDConsultorio, cons.Diagnostico, cons.Costo, cons.FechaHora, cons.ID,
	)
	if err != nil {
		return utils.Responder(c, "06", modConsul, "consulta-service", nil, "Error al actualizar consulta")
	}

	return utils.Responder(c, "01", modConsul, "consulta-service", fiber.Map{"mensaje": "Consulta actualizada"})
}

func EliminarConsulta(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_consulta"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", modConsul, "consulta-service", nil, "ID inválido")
	}

	_, err := config.DB.Exec("DELETE FROM Consultas WHERE id_consulta=$1", body.ID)
	if err != nil {
		return utils.Responder(c, "06", modConsul, "consulta-service", nil, "Error al eliminar consulta")
	}
	return utils.Responder(c, "01", modConsul, "consulta-service", fiber.Map{"mensaje": "Consulta eliminada"})
}


func ObtenerConsultasPaciente(c *fiber.Ctx) error {
    var reqBody struct {
        IdPaciente int `json:"id_paciente"`
    }

    // Parsear el cuerpo de la solicitud
    if err := c.BodyParser(&reqBody); err != nil {
        fmt.Println("Error al parsear body:", err)
        return c.Status(400).JSON(fiber.Map{
            "statusCode": 400,
            "message":    "Body inválido o formato incorrecto",
            "error":      err.Error(),
        })
    }

    // Validar ID del paciente
    if reqBody.IdPaciente <= 0 {
        return c.Status(400).JSON(fiber.Map{
            "statusCode": 400,
            "message":    "ID de paciente inválido",
        })
    }

    // Ejecutar consulta SQL
    rows, err := config.DB.Query(`
        SELECT 
            id_consulta, 
            id_paciente, 
            tipo, 
            id_receta, 
            id_horario, 
            id_consultorio, 
            diagnostico, 
            costo, 
            fecha_hora 
        FROM Consultas 
        WHERE id_paciente = $1`, 
        reqBody.IdPaciente)
    if err != nil {
        fmt.Println("Error en consulta SQL:", err)
        return c.Status(500).JSON(fiber.Map{
            "statusCode": 500,
            "message":    "Error al obtener consultas",
            "error":      err.Error(),
        })
    }
    defer rows.Close()

    var consultas []models.Consulta
    for rows.Next() {
        var cons models.Consulta
        var (
            diagnostico sql.NullString
            idReceta    sql.NullInt64
            costo       sql.NullFloat64
        )
        
        // Escanear los valores, incluyendo los que pueden ser NULL
        if err := rows.Scan(
            &cons.ID,
            &cons.IDPaciente,
            &cons.Tipo,
            &idReceta,
            &cons.IDHorario,
            &cons.IDConsultorio,
            &diagnostico,
            &costo,
            &cons.FechaHora,
        ); err != nil {
            fmt.Println("Error al escanear fila:", err)
            continue
        }
        
        // Asignar valores NULLables con sus valores por defecto
        if diagnostico.Valid {
            cons.Diagnostico = diagnostico.String
        } else {
            cons.Diagnostico = "" // Valor por defecto para string
        }
        
        if idReceta.Valid {
            val := int(idReceta.Int64)
            cons.IDReceta = &val
        } else {
            cons.IDReceta = nil // Valor por defecto para puntero
        }
        
        if costo.Valid {
            cons.Costo = costo.Float64
        } else {
            cons.Costo = 0.0 // Valor por defecto para float
        }
        
        consultas = append(consultas, cons)
    }

    // Verificar si hubo errores después de iterar
    if err = rows.Err(); err != nil {
        fmt.Println("Error después de iterar filas:", err)
        return c.Status(500).JSON(fiber.Map{
            "statusCode": 500,
            "message":    "Error al procesar resultados",
            "error":      err.Error(),
        })
    }

    return c.Status(200).JSON(fiber.Map{
        "data":       consultas,
        "from":       "consulta-service",
        "intCode":    "Consul01",
        "message":    "Operación realizada exitosamente",
        "status":     "S01",
        "statusCode": 200,
    })
}


