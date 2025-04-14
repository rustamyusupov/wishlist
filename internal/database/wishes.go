package database

import (
	"fmt"

	"github.com/rustamyusupov/wishlist/internal/models"
)

func GetWishlist() ([]models.Wish, error) {
	db := GetDB()

	query := `
		SELECT w.id, w.link, w.name, p.price, cur.symbol, cat.name, w.sort, w.created_at
		FROM wishlist w
		JOIN categories cat ON w.category_id = cat.id
		JOIN (
			SELECT wish_id, price, currency_id, MAX(created_at) as max_date
			FROM prices
			GROUP BY wish_id
		) latest_p ON latest_p.wish_id = w.id
		JOIN prices p ON p.wish_id = w.id AND p.created_at = latest_p.max_date
		JOIN currencies cur ON p.currency_id = cur.id
		ORDER BY cat.name, w.sort, w.name
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlist []models.Wish
	for rows.Next() {
		var wish models.Wish
		err := rows.Scan(&wish.ID, &wish.Link, &wish.Name, &wish.Price, &wish.Currency, &wish.Category, &wish.Sort, &wish.CreatedAt)
		if err != nil {
			return nil, err
		}
		wishlist = append(wishlist, wish)
	}

	return wishlist, nil
}

func GetWishByID(id string) (models.Wish, error) {
	db := GetDB()

	query := `
		SELECT w.id, w.link, w.name, p.price, cur.code, cat.name, w.sort, w.created_at
		FROM wishlist w
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
	err := db.QueryRow(query, id, id).Scan(&wish.ID, &wish.Link, &wish.Name, &wish.Price, &wish.Currency, &wish.Category, &wish.Sort, &wish.CreatedAt)
	if err != nil {
		return models.Wish{}, err
	}

	return wish, nil
}

func CreateWish(link, name string, categoryID int, price float64, currencyID int, sort int) (int, error) {
	db := GetDB()

	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		`INSERT INTO wishlist (link, name, category_id, sort) VALUES (?, ?, ?, ?)`,
		link, name, categoryID, sort,
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

func UpdateWish(id int, link, name string, categoryID int, sort int) error {
	db := GetDB()

	_, err := db.Exec(
		`UPDATE wishlist SET link = ?, name = ?, category_id = ?, sort = ? WHERE id = ?`,
		link, name, categoryID, sort, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update wish: %w", err)
	}

	return nil
}

func DeleteWish(id int) error {
	db := GetDB()

	_, err := db.Exec(`DELETE FROM wishlist WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete wish: %w", err)
	}

	return nil
}
