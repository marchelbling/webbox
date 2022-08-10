package middleware

import (
	"net/http"
)

type recovery struct {
	h http.Handler
}

func (r recovery) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			// FIXME print stack (see  https://github.com/gorilla/handlers/blob/master/recovery.go#L89-L95)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError))) // nolint: errcheck
		}
	}()
	r.h.ServeHTTP(w, req)
}

// Recover returns a middleware that will gracefully handle panics.
func Recover() Middleware {
	return func(h http.Handler) http.Handler {
		return recovery{h}
	}
}
