package http

import "net/http"

// StatusResponseWriter wraps a http.ResponseWriter to allow access to the status.
// If you need to access both the status and the response, see ResponseWriter instead.
type StatusResponseWriter struct {
	http.ResponseWriter
	Status int
}

// WriteHeader intercepts the status code then forwards it to the wrapped ResponseWriter.
func (s *StatusResponseWriter) WriteHeader(status int) {
	s.Status = status
	s.ResponseWriter.WriteHeader(status)
}
