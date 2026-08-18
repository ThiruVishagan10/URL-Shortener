package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thiruvishagan10/URL-Shortener/internal/auth"
)

const (
	oauthStateCookie    = "oauth_state"
	oauthVerifierCookie = "oauth_code_verifier"

	oauthCookieMaxAge = 5 * 60
)

type AuthHandler struct {
	googleAuth *auth.GoogleOAuthProvider
}

func NewAuthHandler(
	googleAuth *auth.GoogleOAuthProvider,
) *AuthHandler {
	return &AuthHandler{
		googleAuth: googleAuth,
	}
}

func (h *AuthHandler) GoogleLogin(
	w http.ResponseWriter,
	r *http.Request,
) {
	authRequest, err := h.googleAuth.Begin()
	if err != nil {
		http.Error(
			w,
			"Failed to initialize Google authentication",
			http.StatusInternalServerError,
		)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookie,
		Value:    authRequest.State,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   oauthCookieMaxAge,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     oauthVerifierCookie,
		Value:    authRequest.CodeVerifier,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   oauthCookieMaxAge,
	})

	http.Redirect(
		w,
		r,
		authRequest.URL,
		http.StatusFound,
	)
}

func (h *AuthHandler) GoogleCallback(
	w http.ResponseWriter,
	r *http.Request,
) {
	query := r.URL.Query()

	code := query.Get("code")
	state := query.Get("state")

	if code == "" {
		http.Error(w, "Authorization code missing", http.StatusBadRequest)
		return
	}

	if state == "" {
		http.Error(w, "OAuth state missing", http.StatusBadRequest)
		return
	}

	stateCookie, err := r.Cookie(oauthStateCookie)
	if err != nil {
		http.Error(w, "OAuth Cookie missing", http.StatusBadRequest)
		return
	}

	if state != stateCookie.Value {
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	verifierCookie, err := r.Cookie(oauthVerifierCookie)
	if err != nil {
		http.Error(w, "OAuth verifier missing", http.StatusBadRequest)
		return
	}

	identity, err := h.googleAuth.Authenticate(
		r.Context(),
		code,
		verifierCookie.Value,
	)
	if err != nil {
		http.Error(w, "Google authentication failed", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Subject   string  `json:"subject"`
		Email     string  `json:"email"`
		Name      string  `json:"name"`
		AvatarURL *string `json:"avatar_url,omitempty"`
	}{
		Subject:   identity.Subject,
		Email:     identity.Email,
		Name:      identity.Name,
		AvatarURL: identity.AvatarURL,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
