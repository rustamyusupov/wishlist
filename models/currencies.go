package models

import (
	"database/sql"
)

type Currency struct {
	Id     int    `json:"id"`
	Code   string `json:"code"`
	Symbol string `json:"symbol"`
}

var PredefinedCurrencies = map[string]string{
	"USD": "$",
	"EUR": "€",
	"RUB": "₽",
}

func GetCurrencies() ([]string, error) {
	db := Connect()
	defer db.Close()

	query := `SELECT code FROM currencies ORDER BY code`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var currencies []string
	for rows.Next() {
		var currency string
		err := rows.Scan(&currency)
		if err != nil {
			return nil, err
		}
		currencies = append(currencies, currency)
	}

	return currencies, nil
}

func GetCurrencyID(tx *sql.Tx, currencyCode string) (int, error) {
	var currencyID int

	err := tx.QueryRow("SELECT id FROM currencies WHERE code = ?", currencyCode).Scan(&currencyID)
	if err != nil {
		return 0, err
	}

	return currencyID, nil
}

func InitializeCurrencies() error {
	db := Connect()
	defer db.Close()

	for code, symbol := range PredefinedCurrencies {
		_, err := db.Exec("INSERT OR IGNORE INTO currencies (code, symbol) VALUES (?, ?)", code, symbol)
		if err != nil {
			return err
		}
	}

	return nil
}
