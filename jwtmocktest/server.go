package jwtmocktest

import (
	"fmt"

	"net/http/httptest"

	"github.com/nayyara-cropsey/jwtmock"
	"github.com/nayyara-cropsey/jwtmock/internal/handlers"
	"github.com/nayyara-cropsey/jwtmock/internal/jwks"
	"github.com/nayyara-cropsey/jwtmock/internal/service"
	"github.com/nayyara-cropsey/jwtmock/log"

	"time"
)

var (
	// reasonable test defaults
	defaultCertLen = 24 * time.Hour
	defaultKeyLen  = 1024
)

// A Server is an HTTP server listening on a system-chosen port on the
// local loopback interface, for use in end-to-end HTTP tests.
type Server struct {
	*httptest.Server

	keystore    *service.KeyStore
	clientsRepo *service.ClientRepo
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

	logger := log.NewLogger(log.WithLevel(log.Debug))
	clientRepo := service.NewClientRepo()
	handler := handlers.NewHandler(keyStore, clientRepo, logger)
	server := httptest.NewServer(handler)

	return &Server{
		Server:      server,
		keystore:    keyStore,
		clientsRepo: clientRepo,
	}, nil
}

// GenerateJWT generates a JWT token for use in authorization header.
func (s *Server) GenerateJWT(claims jwtmock.Claims) (string, error) {
	signingKey := s.keystore.GetSigningKey()
	return claims.CreateJWT(signingKey)
}

// RegisterClient registers a new client for subsequent token request.
func (s *Server) RegisterClient(registration jwtmock.ClientRegistration) error {
	return s.clientsRepo.Register(registration)
}
