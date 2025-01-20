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
			id INTEGER PRIMARY KEY AUTOINCREMENT,
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

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM wishes`).Scan(&count)
	if err != nil {
		log.Fatalf("🔥 failed to count wishes: %s", err.Error())
	}

	if count == 0 {
		insertDefaultWishes(db)
	}
}

func insertDefaultWishes(db *sql.DB) {
	defaultWishes := []Wish{
		{Name: "Wish 1", Link: "http://exmpl.com", Price: 1.0, Currency: "$", Category: "Apparel"},
		{Name: "Wish 2", Link: "http://exmpl.com", Price: 2.0, Currency: "€", Category: "Cycling Gear"},
		{Name: "Wish 3", Link: "http://exmpl.com", Price: 3.0, Currency: "₽", Category: "Devices"},
		{Name: "Wish 4", Link: "http://exmpl.com", Price: 4.0, Currency: "$", Category: "Other"},
	}

	for _, wish := range defaultWishes {
		_, err := db.Exec(`INSERT INTO wishes (name, link, price, currency, category) VALUES (?, ?, ?, ?, ?)`, wish.Name, wish.Link, wish.Price, wish.Currency, wish.Category)
		if err != nil {
			log.Fatalf("🔥 failed to insert default wish: %s", err.Error())
		}
	}
}
