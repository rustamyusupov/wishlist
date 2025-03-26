package models

import (
	"strconv"
)

type Wish struct {
	Id        int     `json:"Id"`
	Link      string  `json:"link"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
	Category  string  `json:"category"`
	CreatedAt string  `json:"created_at"`
}

func GetWishes() ([]Wish, error) {
	db := Connect()
	defer db.Close()

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

	var wishes []Wish
	for rows.Next() {
		var wish Wish
		err := rows.Scan(&wish.Id, &wish.Link, &wish.Name, &wish.Price, &wish.Currency, &wish.Category, &wish.CreatedAt)
		if err != nil {
			return nil, err
		}
		wishes = append(wishes, wish)
	}

	return wishes, nil
}

func AddWish(name, link, price, currency, category string) error {
	db := Connect()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	categoryID, err := GetCategoryID(tx, category)
	if err != nil {
		return err
	}

	result, err := tx.Exec("INSERT INTO wishes (name, link, category_id) VALUES (?, ?, ?)",
		name, link, categoryID)
	if err != nil {
		return err
	}

	wishID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	currencyID, err := GetCurrencyID(tx, currency)
	if err != nil {
		return err
	}

	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return err
	}

	err = AddPrice(tx, int(wishID), priceFloat, currencyID)
	if err != nil {
		return err
	}

	return nil
}

func GetWish(id string) (Wish, error) {
	db := Connect()
	defer db.Close()

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

	var wish Wish
	err := db.QueryRow(query, id, id).Scan(&wish.Id, &wish.Link, &wish.Name, &wish.Price, &wish.Currency, &wish.Category, &wish.CreatedAt)
	if err != nil {
		return Wish{}, err
	}

	return wish, nil
}

func UpdateWish(id, name, link, priceStr, currency, category string) error {
	wishId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	newPrice, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return err
	}

	lastPrice, lastCurrency, priceErr := GetLastPrice(wishId)

	db := Connect()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	categoryID, err := GetCategoryID(tx, category)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE wishes SET name = ?, link = ?, category_id = ? WHERE id = ?",
		name, link, categoryID, wishId)
	if err != nil {
		return err
	}

	currencyID, err := GetCurrencyID(tx, currency)
	if err != nil {
		return err
	}

	if priceErr != nil || newPrice != lastPrice || currency != lastCurrency {
		err = AddPrice(tx, wishId, newPrice, currencyID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func DeleteWish(id string) error {
	db := Connect()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	wishId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM wishes WHERE id = ?", id)
	if err != nil {
		return err
	}

	err = DeletePrices(tx, wishId)
	if err != nil {
		return err
	}

	return nil
}
