package handler

import (
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
}

func NewAuthHandler(
	googleAuth *auth.GoogleOAuthProvider,
	userService service.UserService,
	sessionService service.SessionService,
) *AuthHandler {
	return &AuthHandler{
		googleAuth:     googleAuth,
		userService:    userService,
		sessionService: sessionService,
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
