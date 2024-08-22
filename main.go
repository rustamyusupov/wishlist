package main

import (
	"log"
	"main/controllers"
	"main/models"
	"net/http"
)

func main() {
	models.InitDB()
	defer models.Close()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// TODO: CRUD
	http.HandleFunc("/", controllers.GetIndex)

	log.Println("🚀 Starting up on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
