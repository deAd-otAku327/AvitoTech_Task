package handlers

import (
	"encoding/json"
	"merch_shop/pkg/response"
	"net/http"
)

func (c *Controller) SendCoin() http.HandlerFunc {
	type sendCoinRequest struct {
		ToUser string `json:"to_user"`
		Amount int    `json:"amount"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		request := sendCoinRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			response.MakeErrorResponseJSON(w, http.StatusBadRequest, err)
			return
		}

		servErr := c.service.SendCoin(r.Context(), request.ToUser, request.Amount)
		if servErr != nil {
			response.MakeErrorResponseJSON(w, servErr.Code(), servErr)
			return
		}

		response.MakeResponseJSON(w, http.StatusOK, nil)
	}
}
