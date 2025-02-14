package middleware

import (
	"context"
	"merch_shop/pkg/tokenizer"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type key int

const UserIDKey key = 0

func Auth(t *tokenizer.Tokenizer) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenCookie, err := r.Cookie("token")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			token, err := t.VerifyToken(tokenCookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userID, err := token.Claims.GetSubject()
			if err != nil || userID == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			uid, err := strconv.Atoi(userID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserIDKey, uid)))
		})
	}
}
