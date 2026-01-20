package middleware

import (
	"net/http"
)

// Chain chains multiple middleware functions together
func Chain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// Apply applies middleware to a handler
func Apply(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	return Chain(middlewares...)(handler)
}

// ApplyFunc applies middleware to a handler function
func ApplyFunc(handlerFunc http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) http.Handler {
	return Apply(handlerFunc, middlewares...)
}
