package main

import (
	"log"
	"net/http"
)

func main() {
	InitDB()
	defer Close()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", IndexPage)

	log.Println("🚀 Starting up on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
