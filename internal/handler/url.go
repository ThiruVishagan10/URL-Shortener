package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	apperrors "github.com/thiruvishagan10/URL-Shortener/internal/apperrors"
	"github.com/thiruvishagan10/URL-Shortener/internal/middleware"
	"github.com/thiruvishagan10/URL-Shortener/internal/service"
)

type URLHandler struct {
	service service.URLService
}

func NewURLHandler(service service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

type CreateURLRequest struct {
	URL string `json:"url"`
}

type CreateURLResponse struct {
	ID       string `json:"id"`
	ShortURL string `json:"short_url"`
}

type GetURLResponse struct {
	ID          string `json:"id"`
	ShortID     string `json:"short_id"`
	OriginalURL string `json:"original_url"`
	CreatedAt   string `json:"created_at"`
}

func validateURL(rawURL string) error {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return err
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return errors.New("invalid url")
	}

	return nil
}

func (h *URLHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request CreateURLRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	if err := validateURL(request.URL); err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	userID, ok := middleware.UserIDFromContext(r.Context())

	if !ok {
		http.Error(
			w,
			"Authentication required",
			http.StatusUnauthorized,
		)
		return
	}

	url, err := h.service.Create(
		r.Context(),
		userID,
		request.URL,
	)

	if err != nil {
		log.Printf("URL Creation error: %v", err)

		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	response := CreateURLResponse{
		ID:       url.ShortID,
		ShortURL: "http://localhost:8080/" + url.ShortID, // TODO: Replace localhost with configurable base URL.
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Search By Short ID implementation
func (h *URLHandler) GetByShortID(
	w http.ResponseWriter,
	r *http.Request,
) {
	shortID := r.PathValue("shortID")

	url, err := h.service.GetByShortID(r.Context(), shortID)
	if err != nil {
		if errors.Is(err, apperrors.ErrURLNotFound) {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := GetURLResponse{
		ID:          url.ID,
		ShortID:     url.ShortID,
		OriginalURL: url.OriginalURL,
		CreatedAt:   url.CreatedAt.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// URLS redirection
func (h *URLHandler) Redirect(
	w http.ResponseWriter,
	r *http.Request,
) {
	shortID := r.PathValue("shortID")

	url, err := h.service.GetByShortID(
		r.Context(),
		shortID,
	)

	if err != nil {
		if errors.Is(err, apperrors.ErrURLNotFound) {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		url.OriginalURL,
		http.StatusFound,
	)
}
