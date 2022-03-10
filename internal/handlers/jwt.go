package handlers

import (
	"net/http"

	"github.com/nayyara-cropsey/jwtmock"
	"github.com/nayyara-cropsey/jwtmock/log"
)

type jwtResponse struct {
	Token string `json:"token"`
}

// JWTDefaultPath is the default path for JWT handlers.
const JWTDefaultPath = "/generate-jwt"

// JWTHandler provides handlers for working with JWTs
type JWTHandler struct {
	keyStore keyStore
	logger   *log.Logger
}

// NewJWTHandler is the preferred way to create a JWTHandler instance.
func NewJWTHandler(keyStore keyStore, logger *log.Logger) *JWTHandler {
	return &JWTHandler{
		keyStore: keyStore,
		logger:   logger,
	}
}

// RegisterDefaultPaths registers the default paths for JWKS operations.
func (j *JWTHandler) RegisterDefaultPaths(api *http.ServeMux) {
	api.HandleFunc(JWTDefaultPath, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			j.Post(w, r)
		default:
			notFoundResponse(w)
		}
	})
}

// Post creates a signed JWT with the provided claims.
func (j *JWTHandler) Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var claims jwtmock.Claims
	if err := jsonUnmarshal(r, &claims); err != nil {
		j.logger.Errorf("Failed to read claims: %v", err)

		w.WriteHeader(http.StatusBadRequest)

		if err = jsonMarshal(w, errorResponse{
			Message: "Failed to read claims",
			Error:   err.Error(),
		}); err != nil {
			j.logger.Errorf("Failed write JSON response: %v", err)
		}

		return
	}

	signingKey := j.keyStore.GetSigningKey()
	token, err := claims.CreateJWT(signingKey)
	if err != nil {
		j.logger.Errorf("Failed to generate JWT: %v", err)

		w.WriteHeader(http.StatusBadRequest)

		if err = jsonMarshal(w, errorResponse{
			Message: "Failed to generate JWT",
			Error:   err.Error(),
		}); err != nil {
			j.logger.Errorf("Failed write JSON response: %v", err)
		}

		return
	}

	if err := jsonMarshal(w, jwtResponse{Token: token}); err != nil {
		j.logger.Errorf("Failed write JSON response: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
