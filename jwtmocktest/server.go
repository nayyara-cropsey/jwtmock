package jwtmocktest

import (
	"fmt"

	"github.com/nayyara-samuel/jwt-mock/handlers"
	"github.com/nayyara-samuel/jwt-mock/jwks"
	"github.com/nayyara-samuel/jwt-mock/jwt"
	"github.com/nayyara-samuel/jwt-mock/service"

	"net/http/httptest"

	"time"

	"go.uber.org/zap"
)

var (
	// reasonable test defaults
	defaultCertLen = 24 * time.Hour
	defaultKeyLen  = 1024
)

// A Server is an HTTP server listening on a system-chosen port on the
// local loopback interface, for use in end-to-end HTTP tests.
type Server struct {
	URL string

	server   *httptest.Server
	keystore *service.KeyStore
}

// NewServer starts and returns a new Server.
// The caller should call Close when finished, to shut it down.
func NewServer() (*Server, error) {
	certGenerator := service.NewCertificateGenerator(defaultCertLen)
	keyGenerator := jwks.NewGenerator(certGenerator, service.NewRSAKeyGenerator(), defaultKeyLen)
	keyStore, err := service.NewKeyStore(keyGenerator)
	if err != nil {
		return nil, fmt.Errorf("init key store: %w", err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("logger: %w", err)
	}

	handler := handlers.NewHandler(keyStore, logger)
	server := httptest.NewServer(handler)

	return &Server{
		URL:      server.URL,
		server:   server,
		keystore: keyStore,
	}, nil
}

// Close shuts down the server and blocks until all outstanding
// requests on this server have completed.
func (s *Server) Close() {
	s.server.Close()
}

// GenerateJWT generates a JWT token for use in authorization header.
func (s *Server) GenerateJWT(claims jwt.Claims) (string, error) {
	signingKey := s.keystore.GetSigningKey()
	return jwt.CreateToken(claims, signingKey)
}
