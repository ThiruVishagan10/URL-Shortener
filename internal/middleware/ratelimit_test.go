package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/time/rate"
)

func TestIPRateLimiter(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(1), 1)

	handler := limiter.Middleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	t.Run("allows request within limit", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodGet,
			"/",
			nil,
		)

		req.RemoteAddr = "192.168.1.10:12345"

		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusOK,
				rec.Code,
			)
		}
	})

	t.Run("rejects request exceeding limit", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodGet,
			"/",
			nil,
		)

		req.RemoteAddr = "192.168.1.20:12345"

		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf(
				"expected first request status %d, got %d",
				http.StatusOK,
				rec.Code,
			)
		}

		req = httptest.NewRequest(
			http.MethodGet,
			"/",
			nil,
		)

		req.RemoteAddr = "192.168.1.20:12346"

		rec = httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusTooManyRequests {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusTooManyRequests,
				rec.Code,
			)
		}
	})
}
