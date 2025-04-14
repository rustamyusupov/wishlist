package routes

import (
	"net/http"

	"github.com/rustamyusupov/wishlist/internal/auth"
	"github.com/rustamyusupov/wishlist/internal/handlers"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("POST /api/wishlist", handlers.Post)
	mux.HandleFunc("PATCH /api/wishlist/{id}", handlers.Patch)
	mux.HandleFunc("DELETE /api/wishlist/{id}", handlers.Delete)

	mux.HandleFunc("GET /", handlers.Home)
	mux.HandleFunc("GET /new", handlers.New)
	mux.HandleFunc("GET /edit/{id}", handlers.Edit)

	mux.HandleFunc("GET /login", handlers.Login)
	mux.HandleFunc("POST /login", handlers.Login)
	mux.HandleFunc("GET /logout", handlers.Logout)

	return auth.AuthMiddleware(mux)
}
