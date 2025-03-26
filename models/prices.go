package models

import "database/sql"

type Price struct {
	Id         int     `json:"id"`
	WishId     int     `json:"wish_id"`
	Price      float64 `json:"price"`
	CurrencyId int     `json:"currency_id"`
	CreatedAt  string  `json:"created_at"`
}

func GetLastPrice(wishId int) (float64, string, error) {
	db := Connect()
	defer db.Close()

	var price float64
	var currencyCode string

	query := `
		SELECT p.price, c.code
		FROM prices p
		JOIN currencies c ON p.currency_id = c.id
		WHERE p.wish_id = ?
		ORDER BY p.created_at DESC
		LIMIT 1
	`

	err := db.QueryRow(query, wishId).Scan(&price, &currencyCode)
	if err != nil {
		return 0, "", err
	}

	return price, currencyCode, nil
}

func AddPrice(tx *sql.Tx, wishId int, price float64, currencyId int) error {
	_, err := tx.Exec("INSERT INTO prices (wish_id, price, currency_id) VALUES (?, ?, ?)",
		wishId, price, currencyId)

	return err
}

func DeletePrices(tx *sql.Tx, wishId int) error {
	_, err := tx.Exec("DELETE FROM prices WHERE wish_id = ?", wishId)

	return err
}
