package cmd

import (
	"context"
	"errors"
	"fmt"
	"jwt-mock/handlers"
	"jwt-mock/jwks"
	"jwt-mock/service"
	"jwt-mock/store"
	"jwt-mock/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// ConfigKeyType is a special type for setting config key
type ConfigKeyType string

const (
	// ConfigKey is the key for config set in the context for root command.
	ConfigKey ConfigKeyType = "config"

	// server is allowed this  much time to shutdown
	serverShutdownTimeout = time.Second * 5
)

// RootCmd is the root command for this CLI.
// It expects the command context to be set as follows:
// (1) Context must provide config value via "config" key
// (2) Context's Done() channel must be used to signal cancellation for this command to exit correctly.
var RootCmd = &cobra.Command{
	Use:   "jwt-mock",
	Short: "JWT Mock is a server used to mock an authorization server in JWT-based authentication services.",
	RunE: func(cmd *cobra.Command, args []string) error {
		configRaw := cmd.Context().Value(ConfigKey)

		config, ok := configRaw.(*types.Config)
		if !ok {
			return errors.New("invalid config type in context")
		}

		logger, err := zap.NewDevelopment()
		if err != nil {
			cmd.Println("Error creating logger", err)
			return err
		}

		logger.Info("Config", zap.Stringer("config", config))

		mainHandler := gin.Default()
		mainGroup := mainHandler.Group("")

		certGenerator := service.NewCertificateGenerator(config.GetCertificateDuration())
		keyGenerator := jwks.NewGenerator(certGenerator, service.NewRSAKeyGenerator(), config.KeyLength)
		keyStore, err := store.NewKeyStore(keyGenerator)
		if err != nil {
			logger.Error("Error while initializing key store", zap.Error(err))
			return err
		}

		// register handlers for server
		jwksHandler := handlers.NewJWKSHandler(keyStore, logger)
		jwksHandler.RegisterDefaultPaths(mainGroup)

		jwtHandler := handlers.NewJWTHandler(keyStore, logger)
		jwtHandler.RegisterDefaultPaths(mainGroup)

		s := &http.Server{
			Addr:    fmt.Sprintf(":%d", config.Port),
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
					logger.Error("Error while shutting down server", zap.Error(err))
				}
			}
		}()

		// handle cancellation and shutdown server with timeout before exiting
		<-cmd.Context().Done()

		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
		defer timeoutCancel()

		if err := s.Shutdown(timeoutCtx); err != nil {
			logger.Error("Error while shutting down server", zap.Error(err))
			return err
		}

		logger.Info("Server shutdown complete")

		return nil
	},
}
