package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/thiruvishagan10/URL-Shortener/internal/auth"
	"github.com/thiruvishagan10/URL-Shortener/internal/service"
)

const (
	oauthStateCookie    = "oauth_state"
	oauthVerifierCookie = "oauth_code_verifier"

	oauthCookieMaxAge = 5 * 60
)

type AuthHandler struct {
	googleAuth     *auth.GoogleOAuthProvider
	userService    service.UserService
	sessionService service.SessionService
	frontendURL    string
}

func NewAuthHandler(
	googleAuth *auth.GoogleOAuthProvider,
	userService service.UserService,
	sessionService service.SessionService,
	frontendURL string,
) *AuthHandler {
	return &AuthHandler{
		googleAuth:     googleAuth,
		userService:    userService,
		sessionService: sessionService,
		frontendURL:    frontendURL,
	}
}

func (h *AuthHandler) GoogleLogin(
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Println("GoogleLogin handler reached")

	authRequest, err := h.googleAuth.Begin()
	if err != nil {
		log.Printf("Google OAuth initialization failed: %v", err)

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

	log.Println("Redirecting to Google OAuth")

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
		http.Error(
			w,
			"Authorization code missing",
			http.StatusBadRequest,
		)
		return
	}

	if state == "" {
		http.Error(
			w,
			"OAuth state missing",
			http.StatusBadRequest,
		)
		return
	}

	stateCookie, err := r.Cookie(oauthStateCookie)
	// log.Printf("reason for missing OAuth state cookie : %v", err)

	if err != nil {
		http.Error(
			w,
			"OAuth state cookie missing",
			http.StatusBadRequest,
		)
		return
	}

	if state != stateCookie.Value {
		http.Error(
			w,
			"Invalid OAuth state",
			http.StatusBadRequest,
		)
		return
	}

	verifierCookie, err := r.Cookie(oauthVerifierCookie)
	if err != nil {
		http.Error(
			w,
			"OAuth verifier missing",
			http.StatusBadRequest,
		)
		return
	}

	identity, err := h.googleAuth.Authenticate(
		r.Context(),
		code,
		verifierCookie.Value,
	)
	if err != nil {
		http.Error(
			w,
			"Google authentication failed",
			http.StatusUnauthorized,
		)
		return
	}

	user, err := h.userService.FindOrCreate(
		r.Context(),
		identity.Subject,
		identity.Email,
		identity.Name,
		identity.AvatarURL,
	)

	if err != nil {
		http.Error(
			w,
			"Failed to create or retrieve user",
			http.StatusInternalServerError,
		)
		return
	}

	session, err := h.sessionService.Create(
		r.Context(),
		user.ID,
	)

	if err != nil {
		log.Printf("session creation error: %v", err)

		http.Error(
			w,
			"Failed to create session",
			http.StatusInternalServerError,
		)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  session.ExpiresAt,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     oauthVerifierCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.Redirect(
		w,
		r,
		"/",
		http.StatusFound,
	)
}

func (h *AuthHandler) Logout(
	w http.ResponseWriter,
	r *http.Request,
) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		http.Error(
			w,
			"Failed to read session",
			http.StatusInternalServerError,
		)
		return
	}

	if err := h.sessionService.Delete(
		r.Context(),
		cookie.Value,
	); err != nil {
		log.Printf("session deletion error: %v", err)

		http.Error(
			w,
			"Failed to sign out",
			http.StatusInternalServerError,
		)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusNoContent)
}
