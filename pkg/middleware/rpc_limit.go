package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func RpsLimit(rps int) mux.MiddlewareFunc {
	reqChannel := make(chan any, rps)
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			for len(reqChannel) > 0 {
				<-reqChannel
			}
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case reqChannel <- struct{}{}:
				next.ServeHTTP(w, r)
			default:
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
		})
	}

}
