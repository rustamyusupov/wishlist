package main

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

func InitDB() error {
	db := Connect()

	stmt := `CREATE TABLE IF NOT EXISTS wishes (
		id INTEGER,
		link TEXT,
		name TEXT,
		price REAL,
		currency TEXT,
		category TEXT,
		created_at TIMESTAMP DEFAULT DATETIME
	);`

	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

func Close() {
	db.Close()
}
