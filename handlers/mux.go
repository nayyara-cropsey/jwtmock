package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nayyara-cropsey/jwt-mock/service"

	"go.uber.org/zap"
)

// requestLog represents information about a request
type requestLog struct {
	method string
	path   string
	host   string
	took   time.Duration
}

// String returns a string representation of the request log
func (r *requestLog) String() string {
	return fmt.Sprintf("%v %v [%vms] %v", r.method, r.path, r.took.Milliseconds(), r.host)
}

// NewHandler the fully-wired HTTP handler with all routes registered.
func NewHandler(keyStore *service.KeyStore, logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()

	jwksHandler := NewJWKSHandler(keyStore, logger)
	jwksHandler.RegisterDefaultPaths(mux)

	jwtHandler := NewJWTHandler(keyStore, logger)
	jwtHandler.RegisterDefaultPaths(mux)

	// wrap mux with a handler that logs requests
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := &requestLog{
			method: r.Method,
			path:   r.URL.Path,
			host:   r.Host,
		}

		defer func() {
			logger.Debug(l.String())
		}()

		now := time.Now()
		mux.ServeHTTP(w, r)
		l.took = time.Since(now)
	})
}
