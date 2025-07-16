package config

import (
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	
	// Añadir parámetros importantes a la cadena de conexión
	if !strings.Contains(dsn, "?") {
		dsn += "?"
	} else {
		dsn += "&"
	}
	dsn += "binary_parameters=yes&connect_timeout=5"

	log.Println("DSN utilizado:", dsn)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error conectando a Supabase: ", err)
	}

	// Configuración óptima del pool de conexiones
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(30 * time.Minute)
	DB.SetConnMaxIdleTime(5 * time.Minute)

	// Verificar SSL
	var ssl string
	err = DB.QueryRow("SHOW ssl").Scan(&ssl)
	if err != nil {
		log.Println("⚠️ Advertencia al verificar SSL:", err)
	} else {
		log.Println("🔐 SSL en uso:", ssl)
	}

	// Verificar conexión
	err = DB.Ping()
	if err != nil {
		log.Fatal("Error haciendo ping a Supabase: ", err)
	}

	log.Println("✅ Conexión a Supabase exitosa")
}