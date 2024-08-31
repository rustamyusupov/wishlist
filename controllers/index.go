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

func GetIndex(w http.ResponseWriter, r *http.Request) {
	wishes, err := models.GetWishes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories := groupByCategory(wishes)
	categories = sortCategories(categories)

	t, err := template.ParseFiles("views/layout.tmpl", "views/index.tmpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, struct{ Categories []Category }{Categories: categories})
}
