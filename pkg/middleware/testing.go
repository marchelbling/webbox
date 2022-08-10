package middleware

import "net/http"

// NoOpHandler returns a handler that does nothing. For testing purposes.
func NoOpHandler() http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		// no-op
	}
	return http.HandlerFunc(fn)
}
