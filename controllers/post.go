package controllers

import (
	"main/models"
	"net/http"
)

func Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = models.CreateWish(r.FormValue("name"), r.FormValue("link"), r.FormValue("price"), r.FormValue("currency"), r.FormValue("category"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
