package controllers

import (
	"html/template"
	"net/http"
)

func GetAdd(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/layout.tmpl", "views/add.tmpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}
