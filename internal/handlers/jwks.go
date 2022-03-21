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
func (h *JWKSHandler) RegisterDefaultPaths(api *http.ServeMux) {
	api.HandleFunc(JWKSDefaultPath, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.Get(w, r)
		case http.MethodPost:
			h.Post(w, r)
		default:
			notFoundResponse(w)
		}
	})
}

// Get returns a JSON web key set for the authorization server.
func (h *JWKSHandler) Get(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := jsonMarshal(w, h.keyStore.GetJWKS()); err != nil {
		h.logger.Errorf("Failed write JSON response: %v", err)
		return
	}
}

// Post forces a new JSON web key set to be created / the key set to be refreshed.
func (h *JWKSHandler) Post(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := h.keyStore.GenerateNew(); err != nil {
		h.logger.Errorf("Failed to generate new key set: %v", err)

		w.WriteHeader(http.StatusInternalServerError)

		if err = jsonMarshal(w, errorResponse{
			Message: "Failed to refresh JWK set",
			Error:   err.Error(),
		}); err != nil {
			h.logger.Errorf("Failed write JSON response: %v", err)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
