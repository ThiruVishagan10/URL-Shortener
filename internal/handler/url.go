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
	"github.com/thiruvishagan10/URL-Shortener/internal/model"
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
	URL        string `json:"url"`
	Visibility string `json:"visibility"`
}

type CreateURLResponse struct {
	ID       string `json:"id"`
	ShortURL string `json:"short_url"`
}

type GetURLResponse struct {
	ID          string `json:"id"`
	ShortID     string `json:"short_id"`
	OriginalURL string `json:"original_url"`
	Visibility  string `json:"visibility"`
	CreatedAt   string `json:"created_at"`
}

type GetURLsResponse struct {
	URLs []GetURLResponse `json:"urls"`
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

	visibility := model.URLVisibilityPrivate

	if request.Visibility == "" {
		visibility = request.Visibility
	}

	if visibility != model.URLVisibilityPrivate &&
		visibility != model.URLVisibilityPublic {
		http.Error(w, "Invalid visibility", http.StatusBadRequest)
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
		request.Visibility,
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

	userID, _ := middleware.UserIDFromContext(r.Context())

	url, err := h.service.GetByShortID(r.Context(), shortID, userID)
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
		Visibility:  url.Visibility,
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

	userID, _ := middleware.UserIDFromContext(r.Context())

	url, err := h.service.GetByShortID(
		r.Context(),
		shortID,
		userID,
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

func (h *URLHandler) GetByUserID(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := middleware.UserIDFromContext(r.Context())

	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	urls, err := h.service.GetByUserID(r.Context(), userID)

	if err != nil {
		log.Printf("Failed to fetch user URLs: %v", err)

		http.Error(w, "Failed to fetch User URLs", http.StatusInternalServerError)
		return
	}

	response := GetURLsResponse{
		URLs: make([]GetURLResponse, 0, len(urls)),
	}

	for _, url := range urls {
		response.URLs = append(
			response.URLs,
			GetURLResponse{
				ID:          url.ID,
				ShortID:     url.ShortID,
				OriginalURL: url.OriginalURL,
				Visibility:  url.Visibility,
				CreatedAt:   url.CreatedAt.Format(time.RFC3339),
			},
		)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode user URLs: %v", err)
	}
}
