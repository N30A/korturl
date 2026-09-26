package middleware

import (
	"net/http"
	"slices"
)

type Middleware func(http.Handler) http.Handler

func Chain(middlewares ...Middleware) Middleware {
	return func(final http.Handler) http.Handler {
		for i := range slices.Backward(middlewares) {
			final = middlewares[i](final)
		}
		return final
	}
}

func ChainFunc(middlewares ...Middleware) func(http.HandlerFunc) http.Handler {
	return func(final http.HandlerFunc) http.Handler {
		var handler http.Handler = final
		for i := range slices.Backward(middlewares) {
			handler = middlewares[i](handler)
		}
		return handler
	}
}
