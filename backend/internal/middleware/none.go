package middleware

import "net/http"

type NoAuthHandler struct{}

func (*NoAuthHandler) RequireAuth(handler http.Handler) http.HandlerFunc {
	return handler.ServeHTTP
}
