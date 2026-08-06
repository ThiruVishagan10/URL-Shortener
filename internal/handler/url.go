package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

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

	url, err := h.service.Create(
		r.Context(),
		request.URL,
	)

	if err != nil {
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
