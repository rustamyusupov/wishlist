package main

import (
	"fmt"
	"log"
	"net/http"

	"main/controllers"
	"main/models"
)

const port = 3000

func main() {
	models.InitializeDB()

	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("GET /assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("GET /", controllers.Home)
	http.HandleFunc("GET /new", controllers.New)
	http.HandleFunc("GET /edit/{id}", controllers.Edit)

	http.HandleFunc("POST /api/wishes", controllers.Post)
	http.HandleFunc("PATCH /api/wishes/{id}", controllers.Patch)
	http.HandleFunc("DELETE /api/wishes/{id}", controllers.Delete)

	log.Printf("🚀 Starting up on http://localhost:%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
