package main

import (
	"log"

	"net/http"

	"main/controllers"
)

func main() {
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", controllers.GetIndex)
	http.HandleFunc("/add", controllers.GetAdd)

	log.Println("🚀 Starting up on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
