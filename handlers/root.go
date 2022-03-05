package handlers

import (
	"github.com/nayyara-samuel/jwt-mock/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewHandler the fully-wired HTTP handler with all routes registered.
func NewHandler(keyStore *service.KeyStore, logger *zap.Logger) *gin.Engine {
	mainHandler := gin.Default()
	mainGroup := mainHandler.Group("")

	// register handlers for server
	jwksHandler := NewJWKSHandler(keyStore, logger)
	jwksHandler.RegisterDefaultPaths(mainGroup)

	jwtHandler := NewJWTHandler(keyStore, logger)
	jwtHandler.RegisterDefaultPaths(mainGroup)

	return mainHandler
}
