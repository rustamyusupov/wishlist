package models

import "time"

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Currency struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Symbol string `json:"symbol"`
}

type Price struct {
	ID         int       `json:"id"`
	WishID     int       `json:"wish_id"`
	Amount     float64   `json:"price"`
	CurrencyID int       `json:"currency_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type Wish struct {
	ID         int       `json:"id"`
	Link       string    `json:"link"`
	Name       string    `json:"name"`
	CategoryID int       `json:"-"`
	Category   string    `json:"category"`
	Price      float64   `json:"price"`
	Currency   string    `json:"currency"`
	Sort       int       `json:"sort"`
	CreatedAt  time.Time `json:"created_at"`
}

type Option struct {
	Label string
	Value string
}
