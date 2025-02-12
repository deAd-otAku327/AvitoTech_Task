package handlers

import (
	"merch_shop/internal/service"
	"merch_shop/pkg/tokenizer"
	"net/http"
)

func Auth(service service.MerchShopService, t *tokenizer.Tokenizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
