package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func connect() *sql.DB {
	if db != nil {
		return db
	}

	dbPath := os.Getenv("DB_URL")
	if dbPath == "" {
		dbPath = "./wishes.db"
	}

	log.Printf("Connecting to database at: %s", dbPath)

	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("🔥 failed to connect to the database: %s", err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	return db
}

func Initialize() error {
	db = connect()

	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	if err := InitializeCategories(); err != nil {
		return fmt.Errorf("failed to initialize categories: %w", err)
	}

	if err := InitializeCurrencies(); err != nil {
		return fmt.Errorf("failed to initialize currencies: %w", err)
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}

func createTables() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create categories table: %w", err)
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
		return fmt.Errorf("failed to create currencies table: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS wishes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			link TEXT NOT NULL,
			name TEXT NOT NULL,
			category_id INTEGER NOT NULL,
			sort INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (category_id) REFERENCES categories(id)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create wishes table: %w", err)
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
		return fmt.Errorf("failed to create prices table: %w", err)
	}

	return nil
}
