package main

import (
	"fmt"
	"log"
	"net/http"

	"main/controllers"
)

const port = 3000

func MethodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			switch method := r.PostFormValue("_method"); method {
			case http.MethodPut:
				fallthrough
			case http.MethodPatch:
				fallthrough
			case http.MethodDelete:
				r.Method = method
			default:
			}
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("GET /assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("GET /", controllers.Index)
	http.HandleFunc("GET /new", controllers.New)
	http.HandleFunc("GET /edit/{id}", controllers.Edit)

	http.HandleFunc("POST /api/wishes", controllers.Post)
	http.HandleFunc("PUT /api/wishes/{id}", controllers.Put)
	http.HandleFunc("DELETE /api/wishes/{id}", controllers.Delete)

	log.Printf("🚀 Starting up on http://localhost:%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), MethodOverride(http.DefaultServeMux))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
