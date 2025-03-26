package controllers

import (
	"html/template"
	"net/http"

	"main/models"
	"main/utils"
)

type Category struct {
	Name   string
	Wishes []models.Wish
}

func Home(w http.ResponseWriter, r *http.Request) {
	wishes, err := models.GetWishes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories := groupByCategory(wishes)
	categories = sortCategories(categories)

	funcMap := template.FuncMap{
		"formatPrice": utils.FormatPrice,
	}

	t, err := template.New("layout.tmpl").Funcs(funcMap).ParseFiles("views/layout.tmpl", "views/home.tmpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, struct{ Categories []Category }{Categories: categories})
}
