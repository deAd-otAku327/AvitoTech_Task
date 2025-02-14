package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func ResponseTimeLimit(respTime time.Duration) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), respTime)
			defer cancel()

			next.ServeHTTP(w, r.WithContext(ctx))
		})

	}
}
