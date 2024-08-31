package models

func GetCurrencies() ([]string, error) {
	db := Connect()
	defer db.Close()
	query := `
		SELECT DISTINCT currency
		FROM wishes;
	`

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
