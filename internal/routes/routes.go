package routes

import (
	"net/http"

	"github.com/rustamyusupov/wishes/internal/handlers"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("POST /api/wishes", handlers.Post)
	mux.HandleFunc("PATCH /api/wishes/{id}", handlers.Patch)
	mux.HandleFunc("DELETE /api/wishes/{id}", handlers.Delete)

	mux.HandleFunc("GET /", handlers.Home)
	mux.HandleFunc("GET /new", handlers.New)
	mux.HandleFunc("GET /edit/{id}", handlers.Edit)

	return mux
}
