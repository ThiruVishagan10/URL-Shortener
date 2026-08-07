package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/thiruvishagan10/URL-Shortener/internal/config"
	"github.com/thiruvishagan10/URL-Shortener/internal/database"
	"github.com/thiruvishagan10/URL-Shortener/internal/handler"
	"github.com/thiruvishagan10/URL-Shortener/internal/repository"
	"github.com/thiruvishagan10/URL-Shortener/internal/service"
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

	repo := repository.NewURLRepository(dbPool)

	svc := service.NewURLService(repo)

	urlHandler := handler.NewURLHandler(svc)

	mux := http.NewServeMux()

	//Application routes
	mux.HandleFunc("GET /health", handler.Health) //Health Check
	mux.HandleFunc("POST /api/urls", urlHandler.Create) //Create short ID 
	mux.HandleFunc("GET /api/urls/{shortID}", urlHandler.GetByShortID) //Fetch particular ID
	mux.HandleFunc("GET /{shortID}", urlHandler.Redirect) //Redirect to original url through shortid

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	log.Printf("Server running on http://localhost:%s", cfg.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
