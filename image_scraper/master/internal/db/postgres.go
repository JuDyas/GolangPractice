package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitPostgres() *sql.DB {
	var (
		db    *sql.DB
		dbURL = os.Getenv("POSTGRES_URL")
		err   error
	)

	if dbURL == "" {
		log.Println("POSTGRES_URL environment variable not set")
		dbURL = "postgres://user:password@postgres:5432/image_scraper?sslmode=disable"
	}

	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("connect to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("ping database:", err)
	}

	fmt.Println("Successfully connected to PostgreSQL database")
	return db
}
