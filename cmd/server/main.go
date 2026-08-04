package main

import (
	"log"
	"net/http"

	"github.com/thiruvishagan10/URL-Shortener/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("POST /api/urls", handler.CreateURL)

	log.Println("Server running on PORT: 8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
