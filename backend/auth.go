package main

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	log "log/slog"
	"math/big"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

// `identityAuthHandler` handles the auth requests by passing them right through, with no operation done
func identityAuthHandler(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(handler)
}

type JWK struct {
	KeyID     string `json:"kid"`
	Algorithm string `json:"alg"`
	KeyType   string `json:"kty"`
	N         string `json:"n"`
	E         string `json:"e"`
}

func getKeycloakPublicKeys(hostname, realmName string) (map[string]*rsa.PublicKey, error) {

	log.Debug("Downloading Keycloak public keys",
		"keycloakServer", hostname,
		"realm", realmName)

	// https://www.keycloak.org/securing-apps/oidc-layers#_certificate_endpoint
	resp, err := http.Get(fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", hostname, realmName))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var certs struct {
		Keys []JWK `json:"keys"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&certs); err != nil {
		return nil, err
	}

	pubkeys := make(map[string]*rsa.PublicKey)

	for _, key := range certs.Keys {
		nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
		if err != nil {
			return nil, err
		}

		eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
		if err != nil {
			return nil, err
		}

		// Build the RSA public key
		pubKey := &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: int(new(big.Int).SetBytes(eBytes).Int64()),
		}
		pubkeys[key.KeyID] = pubKey
	}

	return pubkeys, nil
}

func validateToken(tokenString string) (*jwt.Token, error) {
	// Parse the token to get the header and validate signature
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Ensure the signing method is RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		fmt.Println("Header:")
		fmt.Println(token.Header)

		// Return the public key for validation
		// return publicKey, nil
		return publicKeys[token.Header["kid"].(string)], nil
	})

	return token, err
}

func keycloakAuthHandler(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Headers", "Authorization")
			return
		}

		bearerToken, contains := r.Header["Authorization"]
		// The user did not provide an auth token
		if !contains {
			http.Error(w, "error: User did not provide the token in the \"Authorization\" header", http.StatusBadRequest)
			return
		}

		bearer := strings.TrimPrefix(bearerToken[0], "Bearer: ")

		log.Debug("Incoming auth request", "token", bearer)

		tok, err := validateToken(bearer)
		if err != nil {
			log.Error("Error validating token", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// TODO: Don't trust anyone as admin
		if tok.Valid {
			http.HandlerFunc(handler).ServeHTTP(w, r)
		}
	})
}

var (
	publicKeys       map[string]*rsa.PublicKey
	keycloakHostname string
	keycloakRealmId  string
)

func SetAuthMethod(authMethodName string) error {

	switch authMethodName {
	case "tailscale":
		authHandlerFunc = tailscaleAuthHandler
		log.Info("Tailscale authentification enabled")
	case "none":
		authHandlerFunc = identityAuthHandler
		log.Warn("Authentification disabled! Please enable this before exposing this service to the internet")
	case "keycloak":
		authHandlerFunc = keycloakAuthHandler

		var err error
		publicKeys, err = getKeycloakPublicKeys(keycloakHostname, keycloakRealmId)
		if err != nil {
			return err
		}

		log.Info("Keycloak authentification enabled")
	default:
		log.Error("invalid auth handler specified: " + authMethodName)
	}

	return nil
}
