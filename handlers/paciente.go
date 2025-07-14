package handlers

import (
	"database/sql"
	"strings"
	"github.com/gofiber/fiber/v2"
	"back-menchaca/config"
	"back-menchaca/models"
	"back-menchaca/utils"
    "log"
	"github.com/pquerna/otp/totp"
)

const modPac = "PAC"

func CrearPaciente(c *fiber.Ctx) error {
    var p models.Paciente

    if err := c.BodyParser(&p); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "statusCode": fiber.StatusBadRequest,
            "intCode": "A01",
            "message": "Datos inválidos: " ,
            "from": "paciente-service",
        })
    }

    // Validación mejorada
    if err := utils.ValidarPaciente(p); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "statusCode": fiber.StatusBadRequest,
            "intCode": "A01",
            "message": err.Error(),
            "from": "paciente-service",
        })
    }

    // Sanitización
    p.Nombre = utils.SanitizarInput(p.Nombre)
    p.Appaterno = utils.SanitizarInput(p.Appaterno)
    p.Apmaterno = utils.SanitizarInput(p.Apmaterno)
    p.Correo = strings.ToLower(utils.SanitizarInput(p.Correo))

    // Verificación de correo único
    var count int
    if err := config.DB.QueryRow(
        `SELECT COUNT(*) FROM (
            SELECT correo FROM Paciente 
            UNION 
            SELECT correo FROM Empleado
        ) AS usuarios WHERE correo = $1`, p.Correo).Scan(&count); err != nil {
        
        log.Printf("Error verificando correo: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "statusCode": fiber.StatusInternalServerError,
            "intCode": "A03",
            "message": "Error al verificar correo",
            "from": "paciente-service",
        })
    }

    if count > 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "statusCode": fiber.StatusBadRequest,
            "intCode": "A01",
            "message": "El correo ya está registrado",
            "from": "paciente-service",
        })
    }

    // Hash de contraseña
    hashed, err := utils.HashPassword(p.Contrasena)
    if err != nil {
        log.Printf("Error hashing password: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "statusCode": fiber.StatusInternalServerError,
            "intCode": "A03",
            "message": "Error al proteger la contraseña",
            "from": "paciente-service",
        })
    }

    // Generar secreto MFA
    mfaKey, err := totp.Generate(totp.GenerateOpts{
        Issuer:      "Menchaca System",
        AccountName: p.Correo,
        SecretSize:  20,
    })
    if err != nil {
        log.Printf("Error generando secreto MFA: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "statusCode": fiber.StatusInternalServerError,
            "intCode": "A03",
            "message": "Error configurando autenticación de dos factores",
            "from": "paciente-service",
        })
    }

    // Creación en BD con transacción
    tx, err := config.DB.Begin()
    if err != nil {
        log.Printf("Error iniciando transacción: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "statusCode": fiber.StatusInternalServerError,
            "intCode": "A03",
            "message": "Error al iniciar transacción",
            "from": "paciente-service",
        })
    }
    defer tx.Rollback()

    // Insertar paciente
    query := `INSERT INTO Paciente 
              (nombre, appaterno, apmaterno, correo, contraseña, mfa_secret, mfa_enabled) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id_paciente`
    
    err = tx.QueryRow(query, 
        p.Nombre, p.Appaterno, p.Apmaterno, p.Correo, hashed, mfaKey.Secret(), true,
    ).Scan(&p.ID)
    
    if err != nil {
        log.Printf("Error insertando paciente: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "statusCode": fiber.StatusInternalServerError,
            "intCode": "A03",
            "message": "Error al crear paciente en la base de datos",
            "from": "paciente-service",
        })
    }

    if err := tx.Commit(); err != nil {
        log.Printf("Error cometiendo transacción: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "statusCode": fiber.StatusInternalServerError,
            "intCode": "A03",
            "message": "Error al confirmar registro",
            "from": "paciente-service",
        })
    }

    // Limpiar datos sensibles antes de responder
    p.Contrasena = ""
    
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "statusCode": fiber.StatusCreated,
        "intCode": "S01",
        "message": "Paciente creado exitosamente",
        "from": "paciente-service",
        "data": fiber.Map{
            "id":        p.ID,
            "nombre":    p.Nombre,
            "correo":    p.Correo,
            "mfaSecret": mfaKey.Secret(),
            "mfaUrl":    mfaKey.URL(),
        },
    })
}


func ObtenerPacientes(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT id_paciente, nombre, appaterno, apmaterno, correo FROM Paciente")
	if err != nil {
		return utils.Responder(c, "06", modPac, "paciente-service", nil, "Error al obtener pacientes")
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

	return utils.Responder(c, "01", modPac, "paciente-service", pacientes)
}

func ObtenerPacientePorID(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_paciente"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", modPac, "paciente-service", nil, "ID inválido")
	}

	var p models.Paciente
	err := config.DB.QueryRow("SELECT id_paciente, nombre, appaterno, apmaterno, correo FROM Paciente WHERE id_paciente = $1", body.ID).
		Scan(&p.ID, &p.Nombre, &p.Appaterno, &p.Apmaterno, &p.Correo)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.Responder(c, "05", modPac, "paciente-service", nil, "Paciente no encontrado")
		}
		return utils.Responder(c, "06", modPac, "paciente-service", nil, "Error al buscar paciente")
	}

	return utils.Responder(c, "01", modPac, "paciente-service", p)
}

func ActualizarPaciente(c *fiber.Ctx) error {
	var p models.Paciente
	if err := c.BodyParser(&p); err != nil || p.ID == 0 {
		return utils.Responder(c, "02", modPac, "paciente-service", nil, "Datos inválidos")
	}

	var current models.Paciente
	err := config.DB.QueryRow(
		"SELECT nombre, appaterno, apmaterno, correo FROM Paciente WHERE id_paciente = $1", p.ID,
	).Scan(&current.Nombre, &current.Appaterno, &current.Apmaterno, &current.Correo)
	if err != nil {
		return utils.Responder(c, "05", modPac, "paciente-service", nil, "Paciente no encontrado")
	}

	//mantener los campos no enviados
	if p.Nombre != "" && !utils.ValidarTextoLetras(p.Nombre) {
		return utils.Responder(c, "02", modPac, "paciente-service", nil, "Nombre inválido")
	}
	if p.Appaterno != "" && !utils.ValidarTextoLetras(p.Appaterno) {
		return utils.Responder(c, "02", modPac, "paciente-service", nil, "Apellido paterno inválido")
	}
	if p.Apmaterno != "" && !utils.ValidarTextoLetras(p.Apmaterno) {
		return utils.Responder(c, "02", modPac, "paciente-service", nil, "Apellido materno inválido")
	}
	if p.Correo != "" && !utils.ValidarCorreo(p.Correo) {
		return utils.Responder(c, "02", modPac, "paciente-service", nil, "Correo inválido")
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
		return utils.Responder(c, "06", modPac, "paciente-service", nil, "Error al actualizar paciente")
	}

	return utils.Responder(c, "01", modPac, "paciente-service", fiber.Map{"mensaje": "Paciente actualizado"})
}

func EliminarPaciente(c *fiber.Ctx) error {
	var body struct {
		ID int `json:"id_paciente"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == 0 {
		return utils.Responder(c, "02", modPac, "paciente-service", nil, "ID inválido")
	}

	_, err := config.DB.Exec("DELETE FROM Paciente WHERE id_paciente = $1", body.ID)
	if err != nil {
		return utils.Responder(c, "06", modPac, "paciente-service", nil, "Error al eliminar paciente")
	}

	return utils.Responder(c, "01", modPac, "paciente-service", fiber.Map{"mensaje": "Paciente eliminado"})
}

