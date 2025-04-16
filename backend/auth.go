package main

import (
	"errors"
	log "log/slog"
	"net/http"

	"tailscale.com/client/local"
)

var CheckRequestAuth func(r *http.Request) (bool, error)

func SetAuthMethod(authMethodName string) {

	switch authMethodName {
	case "tailscale":
		lc := &local.Client{}
		CheckRequestAuth = func(r *http.Request) (bool, error) {
			_, err := lc.WhoIs(r.Context(), r.RemoteAddr)
			if err != nil {
				if errors.Is(err, local.ErrPeerNotFound) {
					return false, nil
				}
				return false, err
			}

			return true, nil
		}
		log.Info("Tailscale authentification enabled")
	}
}
