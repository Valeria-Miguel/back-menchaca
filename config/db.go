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
	if err != nil {
		log.Fatal(" Error conectando a Supabase: ", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal(" Error haciendo ping a Supabase: ", err)
	}
	log.Println("Conexión a Supabase exitosa.")
}
