package main

import (
	"log"
	"net/http"

	"github.com/rustamyusupov/wishlist/internal/database"
	"github.com/rustamyusupov/wishlist/internal/handlers"
	"github.com/rustamyusupov/wishlist/internal/routes"
)

func main() {
	if err := database.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	router := routes.SetupRoutes()

	if err := handlers.Initialize(); err != nil {
		log.Fatalf("Failed to initialize handlers: %v", err)
	}

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
