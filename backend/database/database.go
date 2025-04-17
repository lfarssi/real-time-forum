package database

import (
	"database/sql"
	"io"
	"log"
	"os"
)

// OpenDB connects to the SQLite database, runs migrations, and returns the database connection or an error.
func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "backend/database/forum.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	err = Migrate(db)
	if err != nil {
		log.Printf("Error running migration: %v", err)
		return nil, err
	}

	return db, nil
}

// Migrate reads and executes SQL migration scripts from "sqlite.sql" to set up the database schema.
func Migrate(db *sql.DB) error {
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

	_, err = db.Exec(dataString)
	if err != nil {
		return err
	}

	return nil
}
