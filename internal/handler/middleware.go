package handler

import (
	"net/http"

	"golang.org/x/time/rate"
)

var (
	rateLimiter = make(map[string]*rate.Limiter)
)

func (h *Handler) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			http.Error(rw, "X-Forwarded-For header missing", http.StatusBadRequest)
			return
		}

		limiter, ok := rateLimiter[ip]
		if !ok {
			limiter = rate.NewLimiter(1, 5)
			rateLimiter[ip] = limiter
		}

		if !limiter.Allow() {
			http.Error(rw, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}
