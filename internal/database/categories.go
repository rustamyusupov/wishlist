package database

import (
	"database/sql"
	"fmt"

	"github.com/rustamyusupov/wishes/internal/models"
)

var predefinedCategories = []string{"Apparel", "Devices", "Equipment", "Other"}

func InitializeCategories() error {
	db := GetDB()

	for _, category := range predefinedCategories {
		_, err := db.Exec("INSERT OR IGNORE INTO categories (name) VALUES (?)", category)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetCategories() ([]models.Category, error) {
	db := GetDB()

	rows, err := db.Query(`SELECT id, name FROM categories ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating category rows: %w", err)
	}

	return categories, nil
}

func GetCategoryByName(name string) (models.Category, error) {
	db := GetDB()

	var category models.Category
	err := db.QueryRow(`SELECT id, name FROM categories WHERE name = ?`, name).Scan(&category.ID, &category.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Category{}, fmt.Errorf("category not found: %w", err)
		}
		return models.Category{}, fmt.Errorf("failed to get category: %w", err)
	}

	return category, nil
}
