package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nayyara-cropsey/jwt-mock/handlers"
	"github.com/nayyara-cropsey/jwt-mock/jwks"
	"github.com/nayyara-cropsey/jwt-mock/log"
	"github.com/nayyara-cropsey/jwt-mock/service"
	"github.com/nayyara-cropsey/jwt-mock/types"
)

const (
	// server is allowed this  much time to shutdown
	serverShutdownTimeout = time.Second * 5
)

func Execute(ctx context.Context, cfg *types.Config, logger *log.Logger) error {
	logger.Infof("Config: %v", cfg)

	certGenerator := service.NewCertificateGenerator(cfg.GetCertificateDuration())
	keyGenerator := jwks.NewGenerator(certGenerator, service.NewRSAKeyGenerator(), cfg.KeyLength)
	keyStore, err := service.NewKeyStore(keyGenerator)
	if err != nil {
		logger.Errorf("Error while initializing key store: %v", err)
		return err
	}

	mainHandler := handlers.NewHandler(keyStore, logger)
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mainHandler,

		// add timeout to avoid long I/O waits
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	// start server
	go func() {
		if err := s.ListenAndServe(); err != nil {
			// http.ErrServerClosed is expected after a successful server shutdown
			if !errors.Is(err, http.ErrServerClosed) {
				logger.Errorf("Error while shutting down server: %v", err)
			}
		}
	}()

	// handle cancellation and shutdown server with timeout before exiting
	<-ctx.Done()

	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer timeoutCancel()

	if err := s.Shutdown(timeoutCtx); err != nil {
		logger.Errorf("Error while shutting down server: %v", err)
		return err
	}

	logger.Info("Server shutdown complete")

	return nil
}
