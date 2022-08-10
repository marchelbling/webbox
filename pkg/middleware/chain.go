package middleware

import "net/http"

// Middleware represents... a middleware.
type Middleware func(http.Handler) http.Handler

// Chain represents a list of middlewares
type Chain struct {
	middlewares []Middleware
}

// NewChain creates a new Chain storing the given middlewares.
// Call Then() to actually call the middlewares.
func NewChain(middlewares ...Middleware) Chain {
	return Chain{append(([]Middleware)(nil), middlewares...)}
}

// Then takes a handler and applies all the previously chained middlewares
// in the order they were given.
// NewChain(m1, m2).Then(h) is equivalent to m1(m2(h))
func (c Chain) Then(h http.Handler) http.Handler {
	for i := range c.middlewares {
		h = c.middlewares[len(c.middlewares)-1-i](h)
	}
	return h
}
