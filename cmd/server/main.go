package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/thiruvishagan10/URL-Shortener/internal/config"
	"github.com/thiruvishagan10/URL-Shortener/internal/database"
	"github.com/thiruvishagan10/URL-Shortener/internal/handler"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbPool, err := database.NewPostgresPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("POST /api/urls", handler.CreateURL)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	log.Printf("Server running on http://localhost:%s", cfg.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
