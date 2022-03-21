package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nayyara-cropsey/jwtmock/log"
)

// statusCapturingResponseWriter captures response status
type statusCapturingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusCapturingResponseWriter) WriteHeader(statusCode int) {
	s.status = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

// requestLog represents information about a request
type requestLog struct {
	method string
	path   string
	host   string
	status int
	took   time.Duration
}

// String returns a string representation of the request log
func (r *requestLog) String() string {
	return fmt.Sprintf("%v %v [%vms] %v %v", r.method, r.path,
		r.took.Milliseconds(), r.status, http.StatusText(r.status))
}

// NewHandler the fully-wired HTTP handler with all routes registered.
func NewHandler(keyStore keyStore, clientRepo clientRepo, logger *log.Logger) http.Handler {
	mux := http.NewServeMux()

	jwksHandler := NewJWKSHandler(keyStore, logger)
	jwksHandler.RegisterDefaultPaths(mux)

	jwtHandler := NewJWTHandler(keyStore, logger)
	jwtHandler.RegisterDefaultPaths(mux)

	clientsHandler := NewClientsHandler(keyStore, clientRepo, logger)
	clientsHandler.RegisterDefaultPaths(mux)

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
		ws := &statusCapturingResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		mux.ServeHTTP(ws, r)

		l.status = ws.status
		l.took = time.Since(now)
	})
}
