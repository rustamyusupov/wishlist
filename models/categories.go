package models

import (
	"database/sql"
)

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var PredefinedCategories = []string{"Apparel", "Devices", "Equipment", "Other"}

func GetCategories() ([]string, error) {
	db := Connect()
	defer db.Close()

	query := `SELECT name FROM categories ORDER BY name`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func GetCategoryID(tx *sql.Tx, categoryName string) (int, error) {
	var categoryID int

	err := tx.QueryRow("SELECT id FROM categories WHERE name = ?", categoryName).Scan(&categoryID)
	if err != nil {
		return 0, err
	}

	return categoryID, nil
}

func InitializeCategories() error {
	db := Connect()
	defer db.Close()

	for _, category := range PredefinedCategories {
		_, err := db.Exec("INSERT OR IGNORE INTO categories (name) VALUES (?)", category)
		if err != nil {
			return err
		}
	}

	return nil
}
