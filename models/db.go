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

func InitializeDB() {
	db := Connect()
	defer db.Close()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create categories table: %s", err.Error())
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS currencies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT NOT NULL UNIQUE,
			symbol TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create currencies table: %s", err.Error())
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS wishes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			link TEXT NOT NULL,
			name TEXT NOT NULL,
			category_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (category_id) REFERENCES categories(id)
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create wishes table: %s", err.Error())
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS prices (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			wish_id INTEGER NOT NULL,
			price REAL NOT NULL,
			currency_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (wish_id) REFERENCES wishes(id) ON DELETE CASCADE,
			FOREIGN KEY (currency_id) REFERENCES currencies(id)
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create prices table: %s", err.Error())
	}

	err = InitializeCategories()
	if err != nil {
		log.Fatalf("Failed to initialize categories: %s", err.Error())
	}

	err = InitializeCurrencies()
	if err != nil {
		log.Fatalf("Failed to initialize currencies: %s", err.Error())
	}
}
