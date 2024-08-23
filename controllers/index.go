package controllers

import (
	"html/template"
	"net/http"
	"path/filepath"

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

var tmpl *template.Template

func init() {
	if tmpl == nil {
		if tmpl == nil {
			tmpl = template.Must(tmpl.ParseGlob(filepath.Join("views", "*.html")))
		}
	}
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

	tmpl.ExecuteTemplate(w, "layout.html", wishlist)
}
