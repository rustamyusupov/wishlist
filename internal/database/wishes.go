package database

import (
	"fmt"

	"github.com/rustamyusupov/wishes/internal/models"
)

func GetWishes() ([]models.Wish, error) {
	db := GetDB()

	query := `
		SELECT w.id, w.link, w.name, p.price, cur.symbol, cat.name, w.created_at
		FROM wishes w
		JOIN categories cat ON w.category_id = cat.id
		JOIN (
			SELECT wish_id, price, currency_id, MAX(created_at) as max_date
			FROM prices
			GROUP BY wish_id
		) latest_p ON latest_p.wish_id = w.id
		JOIN prices p ON p.wish_id = w.id AND p.created_at = latest_p.max_date
		JOIN currencies cur ON p.currency_id = cur.id
		ORDER BY cat.name, w.name
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishes []models.Wish
	for rows.Next() {
		var wish models.Wish
		err := rows.Scan(&wish.ID, &wish.Link, &wish.Name, &wish.Price, &wish.Currency, &wish.Category, &wish.CreatedAt)
		if err != nil {
			return nil, err
		}
		wishes = append(wishes, wish)
	}

	return wishes, nil
}

func GetWishByID(id string) (models.Wish, error) {
	db := GetDB()

	query := `
		SELECT w.id, w.link, w.name, p.price, cur.code, cat.name, w.created_at
		FROM wishes w
		JOIN categories cat ON w.category_id = cat.id
		JOIN (
			SELECT wish_id, price, currency_id, MAX(created_at) as max_date
			FROM prices
			GROUP BY wish_id
			HAVING wish_id = ?
		) latest_p ON latest_p.wish_id = w.id
		JOIN prices p ON p.wish_id = w.id AND p.created_at = latest_p.max_date
		JOIN currencies cur ON p.currency_id = cur.id
		WHERE w.id = ?
	`

	var wish models.Wish
	err := db.QueryRow(query, id, id).Scan(&wish.ID, &wish.Link, &wish.Name, &wish.Price, &wish.Currency, &wish.Category, &wish.CreatedAt)
	if err != nil {
		return models.Wish{}, err
	}

	return wish, nil
}

func CreateWish(link, name string, categoryID int, price float64, currencyID int) (int, error) {
	db := GetDB()

	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		`INSERT INTO wishes (link, name, category_id) VALUES (?, ?, ?)`,
		link, name, categoryID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create wish: %w", err)
	}

	wishID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get wish ID: %w", err)
	}

	if price > 0 {
		_, err = tx.Exec(
			`INSERT INTO prices (wish_id, price, currency_id) VALUES (?, ?, ?)`,
			wishID, price, currencyID,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to add price: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return int(wishID), nil
}

func UpdateWish(id int, link, name string, categoryID int) error {
	db := GetDB()

	_, err := db.Exec(
		`UPDATE wishes SET link = ?, name = ?, category_id = ? WHERE id = ?`,
		link, name, categoryID, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update wish: %w", err)
	}

	return nil
}

func DeleteWish(id int) error {
	db := GetDB()

	_, err := db.Exec(`DELETE FROM wishes WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete wish: %w", err)
	}

	return nil
}
