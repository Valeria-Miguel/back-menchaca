package middleware

import (
	"back-menchaca/config"
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"runtime"
	"strings"
	"time"
	"log"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		var bodyCopy map[string]interface{}
		if c.Request().Body() != nil {
			json.NewDecoder(bytes.NewReader(c.Body())).Decode(&bodyCopy)
		}

		// Continuar con la ejecución
		err := c.Next()

		status := c.Response().StatusCode()
		role := c.Locals("role") // Si estás usando roles en JWT

		logEntry := map[string]interface{}{
			"timestamp":     time.Now(),
			"method":        c.Method(),
			"path":          c.OriginalURL(),
			"status":        status,
			"response_time": time.Since(start).Milliseconds(),
			"ip":            c.IP(),
			"user_agent":    c.Get("User-Agent"),
			"level":         getLevel(status),
			"role":          role,
			"system": map[string]interface{}{
				"goVersion": strings.TrimPrefix(runtime.Version(), "go"),
			},
			"query": c.Queries(),
			"body":  bodyCopy,
		}

		// Insertar log en Supabase
		go func(entry map[string]interface{}) {
			db := config.DB
			_, err := db.Exec(`INSERT INTO logs (
				timestamp, method, path, status, response_time, ip, user_agent, level, role, system, query, body
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)`,
				entry["timestamp"], entry["method"], entry["path"], entry["status"], entry["response_time"],
				entry["ip"], entry["user_agent"], entry["level"], entry["role"],
				toJSON(entry["system"]), toJSON(entry["query"]), toJSON(entry["body"]),
			)
			
			if err != nil {
				// Ahora sí, esto usará el paquete "log"
				log.Printf("Error insertando log: %v", err)
			}
		}(logEntry)


		return err
	}
}

func getLevel(status int) string {
	if status >= 500 {
		return "error"
	} else if status >= 400 {
		return "warn"
	}
	return "info"
}

func toJSON(value interface{}) []byte {
	data, _ := json.Marshal(value)
	return data
}
