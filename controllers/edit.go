package controllers

import (
	"html/template"
	"net/http"

	"main/models"
)

func Edit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	wish, err := models.GetWish(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currencies, err := models.GetCurrencies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	currencyOptions := getOptions(currencies)

	categories, err := models.GetCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	categoryOptions := getOptions(categories)

	t, err := template.ParseFiles("views/layout.tmpl", "views/edit.tmpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, struct {
		Wish       models.Wish
		Currencies []Option
		Categories []Option
	}{
		Wish:       wish,
		Currencies: currencyOptions,
		Categories: categoryOptions,
	})
}
