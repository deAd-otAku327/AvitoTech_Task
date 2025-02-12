package handlers

import (
	"merch_shop/internal/db"
	"net/http"
)

func GetInfo(storage db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
