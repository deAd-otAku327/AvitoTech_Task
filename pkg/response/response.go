package response

import (
	"encoding/json"
	"net/http"
)

const (
	contentTypeHeader = "Content-Type"
	contentType       = "application/json"
)

type errorResponse struct {
	Error string `json:"error"`
}

func MakeResponseJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set(contentTypeHeader, contentType)
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data) //nolint:errcheck
	}
}

func MakeErrorResponseJSON(w http.ResponseWriter, code int, err error) {
	MakeResponseJSON(w, code, errorResponse{Error: err.Error()})
}
