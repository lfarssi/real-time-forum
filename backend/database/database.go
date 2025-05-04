package database

import (
	"database/sql"
	"io"
	"log"
	"os"
)

var DB *sql.DB

// OpenDB connects to the SQLite database, runs migrations, and returns the database connection or an error.
func OpenDB() error {
	var err error
	DB, err = sql.Open("sqlite", "backend/database/forum.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return err
	}

	err = DB.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		return err
	}

	err = Migrate()
	if err != nil {
		log.Printf("Error running migration: %v", err)
		return err
	}

	return nil
}

// Migrate reads and executes SQL migration scripts from "sqlite.sql" to set up the database schema.
func Migrate() error {
	file, err := os.Open("backend/database/migration.sql")
	if err != nil {
		return err
	}
	defer file.Close()

	dataBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	dataString := string(dataBytes)

	_, err = DB.Exec(dataString)
	if err != nil {
		return err
	}

	return nil
}
