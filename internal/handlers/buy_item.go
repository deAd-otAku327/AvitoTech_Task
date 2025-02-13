package handlers

import (
	"merch_shop/pkg/response"
	"net/http"

	"github.com/gorilla/mux"
)

func (c *Controller) BuyItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemID := mux.Vars(r)["item"]

		err := c.service.BuyItem(r.Context(), itemID)
		if err != nil { // few statuses.
			response.MakeErrorResponseJSON(w, http.StatusInternalServerError, err)
			return
		}
	}
}
