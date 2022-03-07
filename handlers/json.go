package handlers

import (
	"net/http"

	"gopkg.in/square/go-jose.v2/json"
)

func jsonMarshal(w http.ResponseWriter, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func jsonUnmarshal(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func notFoundResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
}
