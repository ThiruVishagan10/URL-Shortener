package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"

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

	// Rate limiter configuration:
	// 10 requests per second with a burst capacity of 20.
	rateLimiter := middleware.NewIPRateLimiter(
		rate.Limit(10),
		20,
	)

	// URL dependencies
	repo := repository.NewURLRepository(dbPool)
	svc := service.NewURLService(repo)
	urlHandler := handler.NewURLHandler(svc, cfg.AppBaseURL)

	mux := http.NewServeMux()

	// Google authentication dependencies
	googleAuth := auth.NewGoogleOAuthProvider(cfg)

	userRepo := repository.NewUserRepository(dbPool)
	userSvc := service.NewUserService(userRepo)

	frontSvc := cfg.FRONTEND_URL

	// Session dependencies
	sessionRepo := repository.NewSessionRepository(dbPool)
	sessionSvc := service.NewSessionService(sessionRepo)

	authHandler := handler.NewAuthHandler(
		googleAuth,
		userSvc,
		sessionSvc,
		frontSvc,
	)

	authMiddleware := middleware.NewAuthMiddleware(sessionSvc)

	// Root
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>URL Shortener API</title>
</head>
<body>
  <main>
    <h1>The frontend is still under development</h1>
    <p>We are sorry, but the URL Shortener frontend is not available just yet.</p>
    <p>In the meantime, feel free to test the application with the APIs below.</p>
    <ul>
      <li><a href="/health">GET /health</a> &mdash; service health check</li>
      <li><a href="/auth/google">GET /auth/google</a> &mdash; sign in with Google</li>
      <li><code>POST /api/urls</code> &mdash; create a short URL</li>
      <li><code>GET /api/urls</code> &mdash; list your URLs</li>
      <li><code>GET /api/urls/{shortID}</code> &mdash; look up a short URL</li>
    </ul>
    <p>See the repository README for the complete API reference and request examples.</p>
  </main>
</body>
</html>`))
	})

	// Current authenticated user
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
					http.Error(
						w,
						"User ID is missing from context",
						http.StatusInternalServerError,
					)
					return
				}

				w.Header().Set("Content-Type", "text/plain")

				_, _ = w.Write(
					[]byte("Authenticated user: " + userID),
				)
			}),
		),
	)

	// Application routes

	// Health check
	mux.HandleFunc(
		"GET /health",
		handler.Health,
	)

	// Create short URL
	mux.Handle(
		"POST /api/urls",
		rateLimiter.Middleware(
			authMiddleware.RequireAuth(
				http.HandlerFunc(urlHandler.Create),
			),
		),
	)

	// Get URL by short ID
	mux.Handle(
		"GET /api/urls/{shortID}",
		authMiddleware.OptionalAuth(
			http.HandlerFunc(urlHandler.GetByShortID),
		),
	)

	// User URLs list
	mux.Handle(
		"GET /api/urls",
		authMiddleware.RequireAuth(
			http.HandlerFunc(urlHandler.GetByUserID),
		),
	)

	// Actual redirection of shortened URLs
	mux.Handle(
		"GET /{shortID}",
		authMiddleware.OptionalAuth(
			http.HandlerFunc(urlHandler.Redirect),
		),
	)

	// Update URL visibility
	mux.Handle(
		"PATCH /api/urls/{shortID}",
		authMiddleware.RequireAuth(
			http.HandlerFunc(urlHandler.UpdateVisibility),
		),
	)

	// Delete URL
	mux.Handle(
		"DELETE /api/urls/{shortID}",
		authMiddleware.RequireAuth(
			http.HandlerFunc(urlHandler.Delete),
		),
	)

	// Google OAuth login
	mux.Handle(
		"GET /auth/google",
		rateLimiter.Middleware(
			http.HandlerFunc(authHandler.GoogleLogin),
		),
	)

	// Google OAuth callback
	mux.HandleFunc(
		"GET /auth/google/callback",
		authHandler.GoogleCallback,
	)

	// Logout
	mux.HandleFunc(
		"POST /api/auth/logout",
		authHandler.Logout,
	)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.CORS(cfg.FRONTEND_URL)(mux),
	}

	log.Printf("Server listening on :%s (public base URL: %s)", cfg.Port, cfg.AppBaseURL)

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
