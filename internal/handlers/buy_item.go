package handlers

import (
	"merch_shop/pkg/response"
	"net/http"

	"github.com/gorilla/mux"
)

func (c *Controller) BuyItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemID := mux.Vars(r)["item"]

		servErr := c.service.BuyItem(r.Context(), itemID)
		if servErr != nil {
			response.MakeErrorResponseJSON(w, servErr.Code(), servErr)
			return
		}
	}
}
