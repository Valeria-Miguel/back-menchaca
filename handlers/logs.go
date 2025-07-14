package handlers

import (
	"back-menchaca/config"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

func GetLogs(c *fiber.Ctx) error {
	rows, err := config.DB.Query(`SELECT * FROM logs ORDER BY timestamp DESC LIMIT 100`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error consultando logs"})
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error leyendo columnas"})
	}

	var results []map[string]interface{}

	for rows.Next() {
		// Creamos un slice de interfaces para los valores
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error escaneando fila"})
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]

			// Convertir []byte a string si aplica
			if b, ok := val.([]byte); ok {
				var decoded interface{}
				// Intentar decodificar JSON (para los campos JSONB como query, body, etc.)
				if json.Unmarshal(b, &decoded) == nil {
					row[col] = decoded
				} else {
					row[col] = string(b)
				}
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	return c.JSON(results)
}
