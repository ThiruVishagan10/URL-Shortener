package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"

	apperrors "github.com/thiruvishagan10/URL-Shortener/internal/apperrors"
	"github.com/thiruvishagan10/URL-Shortener/internal/service"
)

type contextKey string

const userIDContextKey contextKey = "user_id"

type AuthMiddleware struct {
	sessionService service.SessionService
}

func NewAuthMiddleware(
	sessionService service.SessionService,
) *AuthMiddleware {
	return &AuthMiddleware{
		sessionService: sessionService,
	}
}

func (m *AuthMiddleware) RequireAuth(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				http.Error(
					w,
					"Authentication required",
					http.StatusUnauthorized,
				)
				return
			}

			http.Error(
				w,
				"Failed to read session",
				http.StatusInternalServerError,
			)
			return
		}

		session, err := m.sessionService.Find(
			r.Context(),
			cookie.Value,
		)
		if err != nil {
			if errors.Is(err, apperrors.ErrSessionNotFound) {
				http.Error(
					w,
					"Authentication required",
					http.StatusUnauthorized,
				)
				return
			}

			log.Printf("session validation error: %v", err)

			http.Error(
				w,
				"Failed to validate session",
				http.StatusInternalServerError,
			)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			userIDContextKey,
			session.UserID,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)

	return userID, ok
}
