package handlers

import (
	"net/http"

	"github.com/rustamyusupov/wishlist/internal/auth"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if auth.IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		RenderTemplate(w, "login", nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		HandleError(w, err, http.StatusBadRequest, "Error parsing form")
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if auth.Login(w, r, email, password) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	RenderTemplate(w, "login", map[string]interface{}{
		"Error": "Invalid email or password",
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	auth.Logout(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
