package handlers

import (
	"merch_shop/pkg/response"
	"net/http"
)

func (c *Controller) GetInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info, servErr := c.service.GetInfo(r.Context())
		if servErr != nil {
			response.MakeErrorResponseJSON(w, servErr.Code(), servErr)
			return
		}

		response.MakeResponseJSON(w, http.StatusOK, info)
	}
}
