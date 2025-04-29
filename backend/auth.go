package main

import (
	"errors"
	log "log/slog"
	"net/http"

	"github.com/gorilla/handlers"
	"tailscale.com/client/local"
)

var authHandlerFunc func(func(w http.ResponseWriter, r *http.Request)) http.Handler

func tailscaleAuthHandler(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	lc := &local.Client{}

	// If we're behind a reverse-proxy, we need to make sure that the remote address is correct
	return handlers.ProxyHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("auth: handling request from remote address", "ip", r.RemoteAddr)
		_, err := lc.WhoIs(r.Context(), r.RemoteAddr)
		if err != nil {
			if errors.Is(err, local.ErrPeerNotFound) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			log.Error("Error looking up request identity through Tailscale", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.HandlerFunc(handler).ServeHTTP(w, r)
	}))
}

func identityAuthHandler(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(handler)
}

func SetAuthMethod(authMethodName string) {

	switch authMethodName {
	case "tailscale":
		authHandlerFunc = tailscaleAuthHandler
		log.Info("Tailscale authentification enabled")
	case "none":
		authHandlerFunc = identityAuthHandler
		log.Warn("Authentification disabled! Please enable this before exposing this service to the internet")
	}
}
