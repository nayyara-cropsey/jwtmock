package main

import (
	"fmt"
	"jwt-mock/handlers"
	"jwt-mock/jwks"
	"jwt-mock/service"
	"jwt-mock/store"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("Error logger", err)
		return
	}

	g := gin.Default()
	mainGroup := g.Group("")

	certGenerator := service.NewCertificateGenerator(time.Hour * 24)
	keyGenerator := jwks.NewGenerator(certGenerator, service.NewRSAKeyGenerator(), 1024)
	keyStore, err := store.NewKeyStore(keyGenerator)
	if err != nil {
		logger.Error("Error while initializing key store", zap.Error(err))
		return
	}

	jwksHandler := handlers.NewJWKSHandler(keyStore, logger)
	jwksHandler.RegisterDefaultPaths(mainGroup)

	jwtHandler := handlers.NewJWTHandler(keyStore, logger)
	jwtHandler.RegisterDefaultPaths(mainGroup)

	s := &http.Server{
		Addr:    ":8081",
		Handler: g,
		// add timeout to avoid long I/O waits
		ReadTimeout:    2 * time.Minute,
		WriteTimeout:   2 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		logger.Error("Error while shutting down server", zap.Error(err))
		return
	}
}
