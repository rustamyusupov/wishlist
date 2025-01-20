package models

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
		SELECT id, link, name, price, currency, category, created_at
		FROM wishes
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

	query := `
		INSERT INTO wishes (name, link, price, currency, category)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := db.Exec(query, name, link, price, currency, category)
	if err != nil {
		return err
	}

	return nil
}

func GetWish(id string) (Wish, error) {
	db := Connect()
	defer db.Close()

	query := `
		SELECT id, link, name, price, currency, category, created_at
		FROM wishes
		WHERE id = $1
	`

	var wish Wish
	err := db.QueryRow(query, id).Scan(&wish.Id, &wish.Link, &wish.Name, &wish.Price, &wish.Currency, &wish.Category, &wish.CreatedAt)
	if err != nil {
		return Wish{}, err
	}

	return wish, nil
}

func UpdateWish(id, name, link, price, currency, category string) error {
	db := Connect()
	defer db.Close()

	query := `
		UPDATE wishes
		SET name = $1, link = $2, price = $3, currency = $4, category = $5
		WHERE id = $6
	`

	_, err := db.Exec(query, name, link, price, currency, category, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteWish(id string) error {
	db := Connect()
	defer db.Close()

	query := `
		DELETE FROM wishes
		WHERE id = $1
	`

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
