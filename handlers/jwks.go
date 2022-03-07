package handlers

import (
	"net/http"

	"github.com/nayyara-cropsey/jwt-mock/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// JWKSDefaultPath is the default path for JWKS handlers.
const JWKSDefaultPath = "/.well-known/jwks.json"

// JWKSHandler provides handlers for JWKS operations and stores state of the current JWKS.
type JWKSHandler struct {
	keyStore *service.KeyStore
	logger   *zap.Logger
}

// NewJWKSHandler is the preferred way to create a JWKSHandler instance.
func NewJWKSHandler(keyStore *service.KeyStore, logger *zap.Logger) *JWKSHandler {
	return &JWKSHandler{
		keyStore: keyStore,
		logger:   logger,
	}
}

// RegisterDefaultPaths registers the default paths for JWKS operations.
func (j *JWKSHandler) RegisterDefaultPaths(api *gin.RouterGroup) {
	api.GET(JWKSDefaultPath, j.Get)
	api.POST(JWKSDefaultPath, j.Post)
}

// Get returns a JSON web key set for the authorization server.
func (j *JWKSHandler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, j.keyStore.GetJWKS())
}

// Post forces a new JSON web key set to be created / the key set to be refreshed.
func (j *JWKSHandler) Post(c *gin.Context) {
	if err := j.keyStore.GenerateNew(); err != nil {
		j.logger.Error("failed to generate new key set", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{
			Message: "Failed to refresh JWK set",
			Error:   err.Error(),
		})

		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
