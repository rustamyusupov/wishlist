package models

func GetCategories() ([]string, error) {
	db := Connect()
	defer db.Close()
	query := `
		SELECT DISTINCT category
		FROM wishes;
	`

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
