package handlers

import (
	"net/http"

	"github.com/nayyara-cropsey/jwtmock/log"
)

// JWKSDefaultPath is the default path for JWKS handlers.
const JWKSDefaultPath = "/.well-known/jwks.json"

// JWKSHandler provides handlers for JWKS operations and stores state of the current JWKS.
type JWKSHandler struct {
	keyStore keyStore
	logger   *log.Logger
}

// NewJWKSHandler is the preferred way to create a JWKSHandler instance.
func NewJWKSHandler(keyStore keyStore, logger *log.Logger) *JWKSHandler {
	return &JWKSHandler{
		keyStore: keyStore,
		logger:   logger,
	}
}

// RegisterDefaultPaths registers the default paths for JWKS operations.
func (j *JWKSHandler) RegisterDefaultPaths(api *http.ServeMux) {
	api.HandleFunc(JWKSDefaultPath, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			j.Get(w, r)
		case http.MethodPost:
			j.Post(w, r)
		default:
			notFoundResponse(w)
		}
	})
}

// Get returns a JSON web key set for the authorization server.
func (j *JWKSHandler) Get(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := jsonMarshal(w, j.keyStore.GetJWKS()); err != nil {
		j.logger.Errorf("Failed write JSON response: %v", err)
		return
	}
}

// Post forces a new JSON web key set to be created / the key set to be refreshed.
func (j *JWKSHandler) Post(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := j.keyStore.GenerateNew(); err != nil {
		j.logger.Errorf("Failed to generate new key set: %v", err)

		w.WriteHeader(http.StatusInternalServerError)

		if err = jsonMarshal(w, errorResponse{
			Message: "Failed to refresh JWK set",
			Error:   err.Error(),
		}); err != nil {
			j.logger.Errorf("Failed write JSON response: %v", err)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
