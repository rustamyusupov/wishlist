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
