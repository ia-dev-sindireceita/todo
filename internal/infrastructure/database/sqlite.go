package database

import (
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schema string

//go:embed seed.sql
var seed string

// NewSQLiteDB creates a new SQLite database connection
func NewSQLiteDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, err
	}

	// Create tables
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}

	// Seed demo users
	if _, err := db.Exec(seed); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
