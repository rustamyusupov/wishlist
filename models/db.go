package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Connect() *sql.DB {
	if db != nil {
		return db
	}

	db, err := sql.Open("sqlite3", "./wishes.db")
	if err != nil {
		log.Fatalf("🔥 failed to connect to the database: %s", err.Error())
	}

	return db
}

func Migrate() {
	db := Connect()
	defer db.Close()

	query := `
		CREATE TABLE IF NOT EXISTS wishes (
			id INTEGER,
			link TEXT NOT NULL,
			name TEXT NOT NULL,
			price REAL NOT NULL,
			currency TEXT NOT NULL,
			category TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("🔥 failed to migrate the database: %s", err.Error())
	}
}
