package handlers

import (
	"net/http"

	"github.com/nayyara-cropsey/jwtmock"
	"github.com/nayyara-cropsey/jwtmock/log"
)

const (
	// ClientsDefaultPath is the default path for ClientsHandler handlers.
	ClientsDefaultPath = "/clients"

	// ClientDefaultTokenPath is the path for authenticating clients
	// nolint:gosec // ignore irrelevant warning
	ClientDefaultTokenPath = "/oauth/token"
)

// ClientsHandler provides handlers for working with API clients for machine-to-machine workflows.
type ClientsHandler struct {
	keyStore   keyStore
	clientRepo clientRepo

	logger *log.Logger
}

// NewClientsHandler is the preferred way to create a ClientsHandler instance.
func NewClientsHandler(keyStore keyStore, clientRepo clientRepo, logger *log.Logger) *ClientsHandler {
	return &ClientsHandler{
		keyStore:   keyStore,
		clientRepo: clientRepo,
		logger:     logger,
	}
}

// RegisterDefaultPaths registers the default paths for JWKS operations.
func (h *ClientsHandler) RegisterDefaultPaths(api *http.ServeMux) {
	api.HandleFunc(ClientsDefaultPath, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.Register(w, r)
		default:
			notFoundResponse(w)
		}
	})

	api.HandleFunc(ClientDefaultTokenPath, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.Token(w, r)
		default:
			notFoundResponse(w)
		}
	})
}

// Register register a client
func (h *ClientsHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var registration jwtmock.ClientRegistration
	if err := jsonUnmarshal(r, &registration); err != nil {
		h.logger.Errorf("Failed to read client registration: %v", err)

		w.WriteHeader(http.StatusBadRequest)

		if err = jsonMarshal(w, errorResponse{
			Message: "Failed to read client registration",
			Error:   err.Error(),
		}); err != nil {
			h.logger.Errorf("Failed write JSON response: %v", err)
		}

		return
	}

	if err := h.clientRepo.Register(registration); err != nil {
		h.logger.Errorf("Failed to generate JWT: %v", err)

		w.WriteHeader(http.StatusBadRequest)

		if err = jsonMarshal(w, errorResponse{
			Message: "Failed to register client",
			Error:   err.Error(),
		}); err != nil {
			h.logger.Errorf("Failed write JSON response: %v", err)
		}

		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// Token authenticates a client and generates a token
func (h *ClientsHandler) Token(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req jwtmock.ClientTokenRequest
	if err := formUnmarshal(r, &req); err != nil {
		h.logger.Errorf("Failed to read client token req: %v", err)

		w.WriteHeader(http.StatusBadRequest)

		if err = jsonMarshal(w, errorResponse{
			Message: "Failed to read client token req",
			Error:   err.Error(),
		}); err != nil {
			h.logger.Errorf("Failed write JSON response: %v", err)
		}

		return
	}

	signingKey := h.keyStore.GetSigningKey()
	resp, err := h.clientRepo.GenerateToken(req, signingKey)
	if err != nil {
		h.logger.Errorf("Failed to generate token: %v", err)

		w.WriteHeader(http.StatusBadRequest)

		if err = jsonMarshal(w, errorResponse{
			Message: "Failed to generate token",
			Error:   err.Error(),
		}); err != nil {
			h.logger.Errorf("Failed write JSON response: %v", err)
		}

		return
	}

	if err := jsonMarshal(w, resp); err != nil {
		h.logger.Errorf("Failed write JSON response: %v", err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
