package middlewares

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
	"greenlight.damian.net/internal/errorsManager"
)

type Middlewares struct {
	ErrorManager *errorsManager.ErrorsManager
}

func New(em *errorsManager.ErrorsManager) *Middlewares {
	return &Middlewares{
		ErrorManager: em,
	}
}

func (m *Middlewares) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				m.ErrorManager.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (m *Middlewares) RateLimit(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			m.ErrorManager.RateLimitExceededResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
