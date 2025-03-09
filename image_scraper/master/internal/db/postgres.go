package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const migratePath = "file:///migrations"

func InitPostgres() (*sql.DB, error) {
	var (
		db *sql.DB
		//TODO ПЕРАДАВАТЬ
		dbURL = os.Getenv("POSTGRES_URL")
		err   error
	)

	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if err := runMigrations(dbURL, migratePath); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	fmt.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

func runMigrations(url string, migratePath string) error {
	m, err := migrate.New(migratePath, url)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
