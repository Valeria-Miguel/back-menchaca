package middleware

import (
	"back-menchaca/config"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
	"runtime"
	"strings"
    "fmt"
	"time"
    "math"
    "context"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		status := c.Response().StatusCode()

		logEntry := map[string]interface{}{
			"timestamp":     time.Now(),
			"method":        c.Method(),
			"path":          c.Path(),
			"status":        status,
			"response_time": time.Since(start).Milliseconds(),
			"ip":            c.IP(),
			"user_agent":    c.Get("User-Agent"),
			"level":         getLevel(status),
			"system": map[string]interface{}{
				"goVersion": strings.TrimPrefix(runtime.Version(), "go"),
			},
			"body": safeGetBody(c),
		}

		// Usar consulta directa con parámetros posicionales
		query := `
			INSERT INTO logs (
				timestamp, method, path, status, response_time, ip, 
				user_agent, level, system, body
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`

		_, execErr := config.DB.Exec(query,
			logEntry["timestamp"],
			logEntry["method"],
			logEntry["path"],
			logEntry["status"],
			logEntry["response_time"],
			logEntry["ip"],
			logEntry["user_agent"],
			logEntry["level"],
			toJSON(logEntry["system"]),
			toJSON(logEntry["body"]),
		)

		if execErr != nil {
			log.Printf("⚠️ Error insertando log (query directa): %v", execErr)
		}

		return err
	}
}

func saveLogWithRetry(logEntry map[string]interface{}) error {
	// Consulta SQL directa (sin prepared statement)
	query := `
		INSERT INTO logs (
			timestamp, method, path, status, response_time, ip, 
			user_agent, level, system, body
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	maxRetries := 2
	var lastError error

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Usar QueryRowContext en lugar de Exec para mejor manejo de errores
		_, err := config.DB.ExecContext(ctx, query,
			logEntry["timestamp"],
			logEntry["method"],
			logEntry["path"],
			logEntry["status"],
			logEntry["response_time"],
			logEntry["ip"],
			logEntry["user_agent"],
			logEntry["level"],
			toJSON(logEntry["system"]),
			toJSON(logEntry["body"]),
		)

		if err == nil {
			return nil
		}

		lastError = err

		// Manejar específicamente el error de prepared statement
		if strings.Contains(err.Error(), "unnamed prepared statement does not exist") {
			// Reintentar inmediatamente con nueva conexión
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Para otros errores, hacer backoff exponencial
		time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
	}

	return fmt.Errorf("after %d attempts: %v", maxRetries, lastError)
}

// safeGetBody obtiene el cuerpo de forma segura
func safeGetBody(c *fiber.Ctx) interface{} {
	defer func() {
		if r := recover(); r != nil {
			log.Println("⚠️ Recovered in body parsing:", r)
		}
	}()

	if len(c.Body()) == 0 {
		return nil
	}

	var body interface{}
	if err := json.Unmarshal(c.Body(), &body); err != nil {
		// Si falla el unmarshal, devolver como string (limitado a 1KB)
		bodyStr := string(c.Body())
		if len(bodyStr) > 1024 {
			return bodyStr[:1024] + "...[TRUNCATED]"
		}
		return bodyStr
	}
	return body
}


// getLevel mantiene tu lógica original
func getLevel(status int) string {
	if status >= 500 {
		return "error"
	} else if status >= 400 {
		return "warn"
	}
	return "info"
}

// toJSON convierte a JSON de forma segura
func toJSON(value interface{}) string {
	if value == nil {
		return "null"
	}
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf(`{"json_error":"%v"}`, err.Error())
	}
	return string(data)
}


// Función separada para guardar logs con reintentos
func saveLogEntry(entry map[string]interface{}) error {
	stmt := `
		INSERT INTO logs (
			timestamp, method, path, status, response_time, ip, 
			user_agent, level, request_id, system, body
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	maxRetries := 3
	var lastError error

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err := config.DB.ExecContext(ctx, stmt,
			entry["timestamp"],
			entry["method"],
			entry["path"],
			entry["status"],
			entry["response_time"],
			entry["ip"],
			entry["user_agent"],
			entry["level"],
			entry["request_id"],
			toJSON(entry["system"]),
			toJSON(entry["body"]),
		)

		if err == nil {
			return nil
		}

		lastError = err
		log.Printf("Intento %d fallido: %v", i+1, err)

		// Esperar antes de reintentar (backoff exponencial)
		time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
	}

	return fmt.Errorf("fallo después de %d intentos: %v", maxRetries, lastError)
}

// Sanitizar el body para evitar problemas
func sanitizeBody(body interface{}) interface{} {
	if body == nil {
		return nil
	}

	switch v := body.(type) {
	case string:
		if len(v) > 1000 {
			return v[:1000] + "... [TRUNCATED]"
		}
		return v
	case map[string]interface{}:
		// Eliminar campos sensibles
		delete(v, "password")
		delete(v, "token")
		delete(v, "access_token")
		delete(v, "refresh_token")
		return v
	default:
		return body
	}
}



