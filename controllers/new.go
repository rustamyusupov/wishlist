package controllers

import (
	"html/template"
	"main/models"
	"net/http"
)

type Option struct {
	Label string
	Value string
}

func GetNew(w http.ResponseWriter, r *http.Request) {
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

	t, err := template.ParseFiles("views/layout.tmpl", "views/new.tmpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, struct {
		Currencies []Option
		Categories []Option
	}{
		Currencies: currencyOptions,
		Categories: categoryOptions,
	})
}
