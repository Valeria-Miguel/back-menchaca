package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	dsn := os.Getenv("DATABASE_URL")
	log.Println("DSN utilizado:", dsn)

	var ssl string
		err = DB.QueryRow("SHOW ssl").Scan(&ssl)
		if err != nil {
			log.Fatal("Error al verificar si está usando SSL:", err)
		}

		log.Println("🔐 SSL en uso:", ssl)


	if err != nil {
		log.Fatal(" Error conectando a Supabase: ", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal(" Error haciendo ping a Supabase: ", err)
	}
	log.Println("Conexión a Supabase exitosa.")
}
