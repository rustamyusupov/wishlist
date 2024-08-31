package controllers

import (
	"html/template"
	"net/http"

	"main/models"
)

type Category struct {
	Name   string
	Wishes []models.Wish
}

type Wishlist struct {
	Title      string
	Categories []Category
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	wishes, err := models.GetWishes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories := groupByCategory(wishes)
	categories = sortCategories(categories)

	wishlist := Wishlist{
		Title:      "Wishlist",
		Categories: categories,
	}

	t, err := template.ParseFiles("views/layout.html", "views/header.html", "views/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, wishlist)
}
