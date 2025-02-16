package handlers

import (
	"encoding/json"
	"fmt"
	"merch_shop/pkg/response"
	"net/http"
)

func (c *Controller) Auth() http.HandlerFunc {
	type authRequest struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		request := authRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			response.MakeErrorResponseJSON(w, http.StatusBadRequest, err)
			return
		}

		token, servErr := c.service.AuthentificateUser(r.Context(), request.Username, request.Password)
		if servErr != nil {
			response.MakeErrorResponseJSON(w, servErr.Code(), servErr)
			return
		}

		w.Header().Set("Set-Cookie", fmt.Sprintf("token=%s", token))

		response.MakeResponseJSON(w, http.StatusOK, nil)
	}
}
