package utils

import (
	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	StatusCode int         
	IntCode    string      
	CodeModule string      
	Status     string     
	Message    string      
	Data       interface{} 
	From       string     
}

var GenericResponseCatalog = map[string]APIResponse{
	"01": {StatusCode: fiber.StatusOK, Status: "S01", Message: "Operación realizada exitosamente"},
	"02": {StatusCode: fiber.StatusBadRequest, Status: "A01", Message: "Datos de entrada inválidos"},
	"03": {StatusCode: fiber.StatusUnauthorized, Status: "A02", Message: "No autorizado"},
	"04": {StatusCode: fiber.StatusForbidden, Status: "F01", Message: "Acceso denegado por permisos"},
	"05": {StatusCode: fiber.StatusNotFound, Status: "W01", Message: "Recurso no encontrado"},
	"06": {StatusCode: fiber.StatusInternalServerError, Status: "F02", Message: "Error interno del servidor"},
	"07": {StatusCode: fiber.StatusConflict, Status: "A03", Message: "Conflicto con los datos existentes"},
}

// Función central para responder de forma estándar
func Responder(c *fiber.Ctx, intCode string, codeModule string, from string, data interface{}, overrideMessage ...string) error {
	base, ok := GenericResponseCatalog[intCode]
	if !ok {
		// Código no reconocido, usar fallback
		base = APIResponse{
			StatusCode: fiber.StatusInternalServerError,
			Status:     "F",
			Message:    "Código de respuesta no válido",
		}
		intCode = "99"
	}

	response := fiber.Map{
		"statusCode": base.StatusCode,
		"intCode":    codeModule + intCode, // Ej: ANT01
		"status":     base.Status,
		"message":    base.Message,
		"from":       from,
	}

	if len(overrideMessage) > 0 {
		response["message"] = overrideMessage[0]
	}
	if data != nil {
		response["data"] = data
	}

	return c.Status(base.StatusCode).JSON(response)
}
