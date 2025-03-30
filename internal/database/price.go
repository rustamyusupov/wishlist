package database

import (
	"database/sql"
	"fmt"

	"github.com/rustamyusupov/wishes/internal/models"
)

func GetLatestPrice(wishID int) (models.Price, error) {
	db := GetDB()

	query := `
		SELECT p.id, p.wish_id, p.price, p.currency_id, c.code, p.created_at
		FROM prices p
		JOIN currencies c ON p.currency_id = c.id
		WHERE p.wish_id = ?
		ORDER BY p.created_at DESC
		LIMIT 1
	`

	var price models.Price
	var currencyCode string
	var createdAt string

	err := db.QueryRow(query, wishID).Scan(
		&price.ID, &price.WishID, &price.Amount, &price.CurrencyID, &currencyCode, &createdAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Price{}, fmt.Errorf("no price found for wish: %w", err)
		}
		return models.Price{}, fmt.Errorf("failed to get latest price: %w", err)
	}

	return price, nil
}

func CreatePrice(wishID int, amount float64, currencyID int) (models.Price, error) {
	db := GetDB()

	result, err := db.Exec(
		`INSERT INTO prices (wish_id, price, currency_id) VALUES (?, ?, ?)`,
		wishID, amount, currencyID,
	)
	if err != nil {
		return models.Price{}, fmt.Errorf("failed to add price: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Price{}, fmt.Errorf("failed to get price ID: %w", err)
	}

	var price models.Price
	err = db.QueryRow(`
		SELECT id, wish_id, price, currency_id, created_at
		FROM prices WHERE id = ?
	`, id).Scan(&price.ID, &price.WishID, &price.Amount, &price.CurrencyID, &price.CreatedAt)

	if err != nil {
		return models.Price{}, fmt.Errorf("failed to retrieve new price: %w", err)
	}

	return price, nil
}

func UpdatePrice(id int, amount float64, currencyID int) error {
	db := GetDB()

	_, err := db.Exec(
		`UPDATE prices SET price = ?, currency_id = ? WHERE id = ?`,
		amount, currencyID, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update price: %w", err)
	}

	return nil
}

func DeletePricesByWishID(wishID int) error {
	db := GetDB()

	_, err := db.Exec(`DELETE FROM prices WHERE wish_id = ?`, wishID)
	if err != nil {
		return fmt.Errorf("failed to delete prices for wish: %w", err)
	}

	return nil
}
