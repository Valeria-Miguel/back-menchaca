package handlers

import (
	"fmt"
	"log"
	"os"
	"time"
    "database/sql"
	"back-menchaca/config"
	"back-menchaca/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func mapRolToIntCode(rol string) string {
	switch rol {
	case "paciente":
		return "P01"
	case "doctor":
		return "D01"
	case "enfermera":
		return "E01"
	case "administrador":
		return "A01"
	default:
		return "NOSE" // No especificado
	}
}

// Login maneja el inicio de sesión y soporte MFA
func Login(c *fiber.Ctx) error {
	var input struct {
		Correo     string `json:"correo" validate:"required,email"`
		Contrasena string `json:"contrasena" validate:"required,min=6"`
		TOTP       string `json:"totp,omitempty"` // Opcional en primer paso
	}

	
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"statusCode": fiber.StatusBadRequest,
			"intCode":    "A01",
			"message":    "Datos de entrada inválidos",
			"from":       "auth-service",
		})
	}

	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"statusCode": fiber.StatusBadRequest,
			"intCode":    "A01",
			"message":    "Validación fallida: " + err.Error(),
			"from":       "auth-service",
		})
	}

	var (
		id string 
		rol        string
		hash       string
		mfaEnabled bool
		mfaSecret  sql.NullString
	)

	log.Printf("Intentando login con correo: %s", input.Correo)

	err := config.DB.QueryRow(`
    SELECT id_empleado, tipo_empleado as rol, contraseña, mfa_enabled, mfa_secret 
    FROM empleado WHERE correo=$1`, input.Correo).Scan(&id, &rol, &hash, &mfaEnabled, &mfaSecret)

	// Si falla, intentar en pacientes
	if err != nil {
		err = config.DB.QueryRow(`
			SELECT id_paciente, 'paciente' as rol, contraseña, mfa_enabled, mfa_secret 
			FROM paciente WHERE correo=$1`, input.Correo).Scan(&id, &rol, &hash, &mfaEnabled, &mfaSecret)


		if err != nil {
			log.Println("Tampoco en paciente:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"statusCode": fiber.StatusUnauthorized,
				"intCode":    "A01",
				"message":    "Credenciales inválidas",
				"from":       "auth-service",
			})
		}
	}

	// Mapear rol a código interno
	intCodeRol := mapRolToIntCode(rol)

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input.Contrasena)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"statusCode": fiber.StatusUnauthorized,
			"intCode":    "A01",
			"message":    "Credenciales inválidas",
			"from":       "auth-service",
		})
	}

	// Convertir mfaSecret a string si es válido
	secret := ""
	if mfaSecret.Valid {
		secret = mfaSecret.String
	}

	// Si MFA está activado pero no se ha enviado TOTP, pedir MFA
	if mfaEnabled && input.TOTP == "" {
		tempToken, err := utils.GenerateTempToken(input.Correo, rol)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"statusCode": fiber.StatusInternalServerError,
				"intCode":    "A03",
				"message":    "Error generando token temporal",
				"from":       "auth-service",
			})
		}

		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"statusCode": fiber.StatusAccepted,
			"intCode":    "M01",
			"message":    "Se requiere autenticación de dos factores",
			"from":       "auth-service",
			"data": fiber.Map{
				"tempToken":  tempToken,
				"mfaRequired": true,
			},
		})
	}

    // Si MFA aún no ha sido activado, generarlo automáticamente
// Si MFA aún no ha sido activado, generarlo automáticamente
if !mfaEnabled {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Menchaca System",
		AccountName: input.Correo,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"statusCode": fiber.StatusInternalServerError,
			"intCode":    "A03",
			"message":    "Error generando secreto MFA",
			"from":       "auth-service",
		})
	}

	// Guardar nuevo secreto y activar MFA
	table := "Empleado"
	if rol == "paciente" {
		table = "Paciente"
	}
	_, err = config.DB.Exec(
		fmt.Sprintf("UPDATE %s SET mfa_secret=$1, mfa_enabled=true WHERE correo=$2", table),
		key.Secret(),
		input.Correo,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"statusCode": fiber.StatusInternalServerError,
			"intCode":    "A03",
			"message":    "Error guardando secreto MFA",
			"from":       "auth-service",
		})
	}

	// Generar tempToken para MFA
	tempToken, err := utils.GenerateTempToken(input.Correo, rol)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"statusCode": fiber.StatusInternalServerError,
			"intCode":    "A03",
			"message":    "Error generando token temporal MFA",
			"from":       "auth-service",
		})
	}

	// Responder con QR y tempToken para que el frontend pueda continuar
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"statusCode": fiber.StatusAccepted,
		"intCode":    "MFA01",
		"message":    "Autenticación de dos factores configurada. Escanea el QR para activarla.",
		"from":       "auth-service",
		"data": fiber.Map{
			"qrUrl":       key.URL(),
			"secret":      key.Secret(),
			"mfaConfigured": true,
			"tempToken":   tempToken,
		},
	})
}


	// Validar código TOTP si MFA está activado
	if mfaEnabled {
		if !utils.ValidateTOTP(input.TOTP, secret) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"statusCode": fiber.StatusUnauthorized,
				"intCode":    "A02",
				"message":    "Código de autenticación inválido",
				"from":       "auth-service",
			})
		}
	}
	log.Printf("Generando JWT para id: %s, email: %s, rol: %s", id, input.Correo, rol)
	if id == "" {
		log.Println("ERROR: id vacío, no se puede generar token")
		// Maneja error
	}

	// Generar tokens
	accessToken, err := utils.GenerateJWT(id, input.Correo, rol)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"statusCode": fiber.StatusInternalServerError,
			"intCode":    "A03",
			"message":    "Error generando token de acceso",
			"from":       "auth-service",
		})
	}

	refreshToken, err := utils.GenerateRefreshToken(id, input.Correo, rol)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"statusCode": fiber.StatusInternalServerError,
			"intCode":    "A03",
			"message":    "Error generando token de refresco",
			"from":       "auth-service",
		})
	}

	// Configurar cookie segura
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: "Lax",
		Path:     "/",
	})

	return c.JSON(fiber.Map{
		"statusCode": fiber.StatusOK,
		"intCode":    intCodeRol,
		"message":    "Autenticación exitosa",
		"from":       "auth-service",
		"data": fiber.Map{
			"token":        accessToken,
			"tokenType":    "Bearer",
			"expiresIn":    1800,
			"refreshToken": refreshToken,
		},
	})
}


// VerifyMFA valida el token temporal y el código TOTP para completar la MFA
func VerifyMFA(c *fiber.Ctx) error {
	var input struct {
		TempToken string `json:"tempToken" validate:"required"`
		TOTP      string `json:"totp" validate:"required,len=6,numeric"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"statusCode": fiber.StatusBadRequest,
			"intCode":    "A01",
			"message":    "Datos de entrada inválidos: " + err.Error(),
			"from":       "auth-service",
		})
	}

	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"statusCode": fiber.StatusBadRequest,
			"intCode":    "A01",
			"message":    "Validación fallida: " + err.Error(),
			"from":       "auth-service",
		})
	}

	token, err := jwt.Parse(input.TempToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("TEMP_SECRET")), nil
	})

	if err != nil || !token.Valid {
		log.Printf("Error token temporal: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"statusCode": fiber.StatusUnauthorized,
			"intCode":    "A01",
			"message":    "Token temporal inválido o expirado",
			"from":       "auth-service",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"statusCode": fiber.StatusInternalServerError,
			"intCode":    "A03",
			"message":    "Error al procesar los claims del token",
			"from":       "auth-service",
		})
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"statusCode": fiber.StatusBadRequest,
			"intCode":    "A01",
			"message":    "Token no contiene email válido",
			"from":       "auth-service",
		})
	}

	rol, ok := claims["rol"].(string)
	if !ok || rol == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"statusCode": fiber.StatusBadRequest,
			"intCode":    "A01",
			"message":    "Token no contiene rol válido",
			"from":       "auth-service",
		})
	}

	// Obtener secreto MFA actual
	var mfaSecret string

	var id string
	err = config.DB.QueryRow(`
		SELECT id, mfa_secret FROM (
			SELECT id_paciente as id, correo, mfa_secret FROM Paciente WHERE correo = $1
			UNION
			SELECT id_empleado as id, correo, mfa_secret FROM Empleado WHERE correo = $1
		) AS usuarios LIMIT 1`, email).Scan(&id, &mfaSecret)


	isNewMFA := false
	if err != nil || mfaSecret == "" {
		// Generar nuevo secreto MFA
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "SistemaHospitalario",
			AccountName: email,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"statusCode": fiber.StatusInternalServerError,
				"intCode":    "A03",
				"message":    "Error generando secreto MFA",
				"from":       "auth-service",
			})
		}

		mfaSecret = key.Secret()
		isNewMFA = true
	}

	intCodeRol := mapRolToIntCode(rol)

	// Validar código TOTP
	valid, err := totp.ValidateCustom(input.TOTP, mfaSecret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})

	if err != nil || !valid {
		log.Printf("Validación MFA fallida - Código: %s, Secreto: %s, err: %v", input.TOTP, mfaSecret, err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"statusCode": fiber.StatusUnauthorized,
			"intCode":    "A02",
			"message":    "Código de autenticación inválido o expirado",
			"from":       "auth-service",
		})
	}

	// Guardar secreto MFA y activar MFA si es nuevo
	if isNewMFA {
		var updateQuery string
		if rol == "empleado" {
			updateQuery = `UPDATE Empleado SET mfa_secret=$1, mfa_enabled=true WHERE correo=$2`
		} else {
			updateQuery = `UPDATE Paciente SET mfa_secret=$1, mfa_enabled=true WHERE correo=$2`
		}

		_, err := config.DB.Exec(updateQuery, mfaSecret, email)
		if err != nil {
			log.Printf("Error guardando MFA en BD: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"statusCode": fiber.StatusInternalServerError,
				"intCode":    "A03",
				"message":    "Error guardando configuración MFA",
				"from":       "auth-service",
			})
		}
	}

	// Generar tokens finales
	accessToken, err := utils.GenerateJWT(id, email, rol)
	if err != nil {
		log.Printf("Error generando access token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"statusCode": fiber.StatusInternalServerError,
			"intCode":    "A03",
			"message":    "Error generando token de acceso",
			"from":       "auth-service",
		})
	}

	refreshToken, err := utils.GenerateRefreshToken(id, email, rol)
	if err != nil {
		log.Printf("Error generando refresh token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"statusCode": fiber.StatusInternalServerError,
			"intCode":    "A03",
			"message":    "Error generando token de refresco",
			"from":       "auth-service",
		})
	}

	// Configurar cookie para refresh token
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: "Lax",
		Path:     "/",
	})

	return c.JSON(fiber.Map{
		"statusCode": fiber.StatusOK,
		"intCode":    intCodeRol,
		"message":    "Autenticación MFA exitosa",
		"from":       "auth-service",
		"data": fiber.Map{
			"token":        accessToken,
			"tokenType":    "Bearer",
			"expiresIn":    1800,
			"refreshToken": refreshToken,
			"mfaActivated": isNewMFA,
		},
	})
}

// RefreshToken genera un nuevo access token dado un refresh token válido
func RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	log.Println("[DEBUG] Cookie refresh_token:", refreshToken)

	if refreshToken == "" {
		var input struct {
			RefreshToken string `json:"refreshToken" validate:"required"`
		}
		if err := c.BodyParser(&input); err != nil {
			log.Println("[ERROR] Error al parsear cuerpo:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"intCode":    "A01",
				"message":    "Datos de entrada inválidos",
				"from":       "auth-service",
			})
		}

		if err := validate.Struct(input); err != nil {
			log.Println("[ERROR] Validación fallida:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"intCode":    "A01",
				"message":    "Validación fallida: " + err.Error(),
				"from":       "auth-service",
			})
		}

		refreshToken = input.RefreshToken
		log.Println("[DEBUG] Refresh token recibido en body:", refreshToken)
	}

	if refreshToken == "" {
		log.Println("[ERROR] Refresh token vacío")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"statusCode": fiber.StatusUnauthorized,
			"intCode":    "A01",
			"message":    "Refresh token requerido",
			"from":       "auth-service",
		})
	}

	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	if err != nil || !token.Valid {
		log.Printf("[ERROR] Token inválido o expirado: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"statusCode": fiber.StatusUnauthorized,
			"intCode":    "A01",
			"message":    "Refresh token inválido o expirado",
			"from":       "auth-service",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("[ERROR] Claims no válidos")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"statusCode": fiber.StatusInternalServerError,
			"intCode":    "A03",
			"message":    "Error al procesar los claims del token",
			"from":       "auth-service",
		})
	}

	email, okEmail := claims["email"].(string)
	rol, okRol := claims["rol"].(string)
	log.Println("[DEBUG] Claims extraídos - Email:", email, "Rol:", rol)

	if !okEmail || !okRol || email == "" || rol == "" {
		log.Println("[ERROR] Email o rol inválido en el token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"statusCode": fiber.StatusUnauthorized,
			"intCode":    "A01",
			"message":    "Refresh token inválido",
			"from":       "auth-service",
		})
	}

	var id string
	err = config.DB.QueryRow(`
		SELECT id FROM (
			SELECT id_paciente AS id, correo FROM paciente WHERE correo = $1
			UNION
			SELECT id_empleado AS id, correo FROM empleado WHERE correo = $1
		) AS usuarios LIMIT 1`, email).Scan(&id)

	if err != nil || id == "" {
		log.Println("[ERROR] Usuario no encontrado en DB para el email:", email)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"statusCode": fiber.StatusUnauthorized,
			"intCode":    "A01",
			"message":    "Usuario no encontrado",
			"from":       "auth-service",
		})
	}

	log.Println("[DEBUG] Usuario encontrado. ID:", id)

	newToken, err := utils.GenerateJWT(id, email, rol)
	if err != nil {
		log.Println("[ERROR] Error generando nuevo token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"statusCode": fiber.StatusInternalServerError,
			"intCode":    "A03",
			"message":    "Error generando nuevo token",
			"from":       "auth-service",
		})
	}

	log.Println("[DEBUG] Nuevo token generado para:", email)

	return c.JSON(fiber.Map{
		"statusCode": fiber.StatusOK,
		"intCode":    "S01",
		"message":    "Token refrescado exitosamente",
		"from":       "auth-service",
		"data": fiber.Map{
			"token":     newToken,
			"tokenType": "Bearer",
			"expiresIn": 1800,
		},
	})
}

// ActivateMFA genera y activa MFA para un usuario autenticado
/*
func ActivateMFA(c *fiber.Ctx) error {
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"statusCode": fiber.StatusInternalServerError,
            "intCode":    "A03",
            "message":    "Error al obtener información del usuario",
            "from":       "auth-service",
        })
    }

    // Obtener claims del token
    claims, ok := userToken.Claims.(jwt.MapClaims)
    if !ok {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "statusCode": fiber.StatusInternalServerError,
            "intCode":    "A03",
            "message":    "Error al procesar los claims del token",
            "from":       "auth-service",
        })
    }

    // Extraer email y rol de los claims
    email, ok := claims["email"].(string)
    if !ok || email == "" {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "statusCode": fiber.StatusInternalServerError,
            "intCode":    "A03",
            "message":    "Error al obtener el email del usuario",
            "from":       "auth-service",
        })
    }

    rol, ok := claims["rol"].(string)
    if !ok || rol == "" {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "statusCode": fiber.StatusInternalServerError,
            "intCode":    "A03",
            "message":    "Error al obtener el rol del usuario",
            "from":       "auth-service",
        })
    }

    // Generar nuevo secreto MFA
    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      "Menchaca System",
        AccountName: email,
        SecretSize:  20,
    })
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "statusCode": fiber.StatusInternalServerError,
            "intCode":    "A03",
            "message":    "Error generando secreto MFA",
            "from":       "auth-service",
        })
    }

    // Determinar la tabla según el rol
    tableName := "Paciente"
    if rol == "empleado" {
        tableName = "Empleado"
    }

    // Actualizar el usuario en la base de datos (activando MFA)
    _, err = config.DB.Exec(
        fmt.Sprintf("UPDATE %s SET mfa_secret = $1, mfa_enabled = true WHERE correo = $2", tableName),
        key.Secret(),
        email,
    )
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "statusCode": fiber.StatusInternalServerError,
            "intCode":    "A03",
            "message":    "Error actualizando usuario en la base de datos",
            "from":       "auth-service",
        })
    }

    return c.JSON(fiber.Map{
        "statusCode": fiber.StatusOK,
        "intCode":    "S01",
        "message":    "MFA configurado exitosamente",
        "from":       "auth-service",
        "data": fiber.Map{
            "secret": key.Secret(),
            "qrUrl":  key.URL(),
        },
    })
}*/