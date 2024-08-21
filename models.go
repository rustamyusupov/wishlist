package main

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
	// TODO: check select last wish by created_at
	query := `
		WITH ranked_wishes AS (
				SELECT *, ROW_NUMBER() OVER (PARTITION BY id ORDER BY created_at DESC) AS rn
				FROM wishes
		)
		SELECT id, link, name, price, currency, category, created_at
		FROM ranked_wishes
		WHERE rn = 1;
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
