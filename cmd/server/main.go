package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/thiruvishagan10/URL-Shortener/internal/auth"
	"github.com/thiruvishagan10/URL-Shortener/internal/config"
	"github.com/thiruvishagan10/URL-Shortener/internal/database"
	"github.com/thiruvishagan10/URL-Shortener/internal/handler"
	"github.com/thiruvishagan10/URL-Shortener/internal/middleware"
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

	googleAuth := auth.NewGoogleOAuthProvider(cfg)
	userRepo := repository.NewUserRepository(dbPool)
	userSvc := service.NewUserService(userRepo)

	sessionRepo := repository.NewSessionRepository(dbPool)
	sessionSvc := service.NewSessionService(sessionRepo)
	authHandler := handler.NewAuthHandler(googleAuth, userSvc, sessionSvc)

	authMiddleware := middleware.NewAuthMiddleware(sessionSvc)

	mux.Handle(
		"GET /api/me",
		authMiddleware.RequireAuth(
			http.HandlerFunc(func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				userID, ok := middleware.UserIDFromContext(
					r.Context(),
				)

				if !ok {
					http.Error(w, "User ID is missing from context", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte("Authenticated user: " + userID))
			}),
		),
	)

	//Application routes
	mux.HandleFunc("GET /health", handler.Health)                      //Health Check
	mux.HandleFunc("POST /api/urls", urlHandler.Create)                //Create short ID
	mux.HandleFunc("GET /api/urls/{shortID}", urlHandler.GetByShortID) //Fetch particular ID
	mux.HandleFunc("GET /{shortID}", urlHandler.Redirect)

	//Google Auth
	mux.HandleFunc(
		"GET /auth/google",
		authHandler.GoogleLogin,
	) //Redirect to original url through shortid

	//Redirect
	mux.HandleFunc(
		"GET /auth/google/callback",
		authHandler.GoogleCallback,
	)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	log.Printf("Server running on http://localhost:%s", cfg.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
