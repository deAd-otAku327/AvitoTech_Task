package handlers

import (
	"encoding/json"
	"errors"
	"merch_shop/pkg/response"
	"net/http"
)

var errInvalidCoinAmount = errors.New("coin amount is invalid")

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

		if request.Amount <= 0 {
			response.MakeErrorResponseJSON(w, http.StatusBadRequest, errInvalidCoinAmount)
			return
		}

		err = c.service.SendCoin(r.Context(), request.ToUser, request.Amount)
		if err != nil {
			response.MakeErrorResponseJSON(w, http.StatusInternalServerError, err)
			return
		}

		response.MakeResponseJSON(w, http.StatusOK, nil)
	}
}
