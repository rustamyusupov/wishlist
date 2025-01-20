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

func AddWish(name, link, price, currency, category string) error {
	db := Connect()
	defer db.Close()

	var lastID int
	err := db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM wishes").Scan(&lastID)
	if err != nil {
		return err
	}

	newID := lastID + 1
	query := `
		INSERT INTO wishes (id, name, link, price, currency, category)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = db.Exec(query, newID, name, link, price, currency, category)
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
		ORDER BY created_at DESC
		LIMIT 1
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

	// TODO: check update last wish by created_at can be bug
	// create new wish
	// update price for old one in same category
	// you have two same wishes wo new one
	query := `
		UPDATE wishes
		SET name = $1, link = $2, price = $3, currency = $4, category = $5
		WHERE id = (
			SELECT id
			FROM wishes
			ORDER BY created_at DESC
			LIMIT 1
		)
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
