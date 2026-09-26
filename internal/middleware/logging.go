package middleware

// https://www.alexedwards.net/blog/making-and-using-middleware

import (
	"log/slog"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request",
			"addr", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL.String(),
			"proto", r.Proto,
		)
		next.ServeHTTP(w, r)
	})
}
