package database

import (
	"database/sql"
	"fmt"

	"github.com/rustamyusupov/wishes/internal/models"
)

var predefinedCurrencies = map[string]string{
	"USD": "$",
	"EUR": "€",
	"RUB": "₽",
}

func InitializeCurrencies() error {
	db := Connect()
	defer db.Close()

	for code, symbol := range predefinedCurrencies {
		_, err := db.Exec("INSERT OR IGNORE INTO currencies (code, symbol) VALUES (?, ?)", code, symbol)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetCurrencies() ([]models.Currency, error) {
	db := Connect()
	defer db.Close()

	rows, err := db.Query(`SELECT id, code, symbol FROM currencies ORDER BY code`)
	if err != nil {
		return nil, fmt.Errorf("failed to query currencies: %w", err)
	}
	defer rows.Close()

	var currencies []models.Currency
	for rows.Next() {
		var currency models.Currency
		if err := rows.Scan(&currency.ID, &currency.Code, &currency.Symbol); err != nil {
			return nil, fmt.Errorf("failed to scan currency: %w", err)
		}
		currencies = append(currencies, currency)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating currency rows: %w", err)
	}

	return currencies, nil
}

func GetCurrencyByCode(code string) (models.Currency, error) {
	db := Connect()
	defer db.Close()

	var currency models.Currency
	err := db.QueryRow(`SELECT id, code, symbol FROM currencies WHERE code = ?`, code).
		Scan(&currency.ID, &currency.Code, &currency.Symbol)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Currency{}, fmt.Errorf("currency not found: %w", err)
		}
		return models.Currency{}, fmt.Errorf("failed to get currency: %w", err)
	}

	return currency, nil
}
