package handlers

import (
	"jwt-mock/jwt"
	"jwt-mock/store"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type jwtResponse struct {
	Token string `json:"token"`
}

// JWTDefaultPath is the default path for JWT handlers.
const JWTDefaultPath = "/generate-jwt"

// JWTHandler provides handlers for working with JWTs
type JWTHandler struct {
	keyStore *store.KeyStore
	logger   *zap.Logger
}

// NewJWTHandler is the preferred way to create a JWTHandler instance.
func NewJWTHandler(keyStore *store.KeyStore, logger *zap.Logger) *JWTHandler {
	return &JWTHandler{
		keyStore: keyStore,
		logger:   logger,
	}
}

// RegisterDefaultPaths registers the default paths for JWKS operations.
func (j *JWTHandler) RegisterDefaultPaths(api *gin.RouterGroup) {
	api.POST(JWTDefaultPath, j.Post)
}

// Post creates a signed JWT with the provided claims.
func (j *JWTHandler) Post(c *gin.Context) {
	var claims jwt.Claims
	if err := c.BindJSON(&claims); err != nil {
		j.logger.Error("failed to read claims", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{
			Message: "Failed to read claims",
			Error:   err.Error(),
		})

		return
	}

	signingKey := j.keyStore.GetSigningKey()
	token, err := jwt.CreateToken(claims, signingKey)
	if err != nil {
		j.logger.Error("failed to generate JWT", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{
			Message: "Failed to generate JWT",
			Error:   err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, jwtResponse{Token: token})
}
