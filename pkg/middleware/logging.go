package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Logging(log *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			client := r.RemoteAddr
			startReq := time.Now()

			next.ServeHTTP(w, r)

			responseTime := time.Since(startReq)

			log.Info(
				fmt.Sprintf("%s %s", r.Method, r.URL.Path),
				slog.String("client", client),
				slog.String("resp_time", responseTime.String()),
			)
		})
	}
}
