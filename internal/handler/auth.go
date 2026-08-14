package handler

import (
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
