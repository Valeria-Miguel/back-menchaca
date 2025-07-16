package middleware

import (
	"back-menchaca/config"
	//"bytes"
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

        // Continuar con la ejecución
        err := c.Next()

        // Preparar datos del log
        status := c.Response().StatusCode()
        role := c.Locals("role")
        
        // Leer el cuerpo de la solicitud de forma segura
        var bodyCopy interface{}
        if len(c.Body()) > 0 {
            if err := json.Unmarshal(c.Body(), &bodyCopy); err != nil {
                bodyCopy = string(c.Body())
            }
        }

        logEntry := map[string]interface{}{
            "timestamp":     time.Now(),
            "method":        c.Method(),
            "path":          c.Path(),
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

        // Insertar log en Supabase de forma sincrónica (sin goroutine)
        if _, err := config.DB.Exec(`
            INSERT INTO logs (
                timestamp, method, path, status, response_time, ip, 
                user_agent, level, role, system, query, body
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
            logEntry["timestamp"],
            logEntry["method"],
            logEntry["path"],
            logEntry["status"],
            logEntry["response_time"],
            logEntry["ip"],
            logEntry["user_agent"],
            logEntry["level"],
            logEntry["role"],
            toJSON(logEntry["system"]),
            toJSON(logEntry["query"]),
            toJSON(logEntry["body"]),
        ); err != nil {
            log.Printf("Error insertando log: %v", err)
        }

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
