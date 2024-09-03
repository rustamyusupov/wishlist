package main

import (
	"log"

	"net/http"

	"main/controllers"
)

func main() {
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("GET /assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("GET /", controllers.GetIndex)
	http.HandleFunc("GET /new", controllers.GetNew)

	http.HandleFunc("POST /api/wishes", controllers.Post)

	// TODO: edit wish / update action

	log.Println("🚀 Starting up on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
