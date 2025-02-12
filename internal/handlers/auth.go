package handlers

import (
	"merch_shop/internal/db"
	"merch_shop/pkg/tokenizer"
	"net/http"
)

func Auth(storage db.DB, t *tokenizer.Tokenizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
