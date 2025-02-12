package middleware

import (
	"context"
	"errors"
	"merch_shop/pkg/response"
	"merch_shop/pkg/tokenizer"
	"net/http"

	"github.com/gorilla/mux"
)

type key int

const userIDKey key = 0

var errBadTokenClaims error = errors.New("failed to extract token claims")

func Auth(t *tokenizer.Tokenizer) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenCookie, err := r.Cookie("token")
			if err != nil {
				response.MakeErrorResponseJSON(w, http.StatusUnauthorized, err)
				return
			}

			token, err := t.VerifyToken(tokenCookie.Value)
			if err != nil {
				response.MakeErrorResponseJSON(w, http.StatusForbidden, err)
				return
			}

			userID, err := token.Claims.GetSubject()
			if err != nil || userID == "" {
				response.MakeErrorResponseJSON(w, http.StatusUnauthorized, errBadTokenClaims)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userIDKey, userID)))
		})
	}
}
