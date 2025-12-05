package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/NickyBoy89/sorrel/backend/internal/auth"
	"golang.org/x/oauth2"
)

type AuthHandler interface {
	LoginHandler(w http.ResponseWriter, r *http.Request)
	CallbackHandler(w http.ResponseWriter, r *http.Request)
}

type OIDCAuthHandler struct {
	authClient *auth.Client
}

func NewAuthHandler(authClient *auth.Client) *OIDCAuthHandler {
	return &OIDCAuthHandler{
		authClient: authClient,
	}
}

// generateRandomSecureString creates a random secure string
func generateRandomSecureString() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

var authStore = make(map[string]bool)

var SessionStore = make(map[string]auth.SessionData)

// LoginHandler initiates the OAuth2 authorization code flow with Keycloak.
// It generates a secure state parameter to prevent CSRF attacks and stores it
// in Redis for later verification during the callback phase.
//
// Returns:
// - 302: Redirects to Keycloak login page
// - 500: Internal Server Error if state generation or storage fails
func (a *OIDCAuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandomSecureString()
	if err != nil {
		slog.Error("failed to generate state", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Store state in session for later verification
	authStore[state] = true

	// Build authentication URL
	authURL := a.authClient.Oauth.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("response_type", "code"),
		oauth2.SetAuthURLParam("scope", "openid profile email"),
	)

	// Redirect to Keycloak login page
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (a *OIDCAuthHandler) CallbackHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form values for callback", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := a.validateStateSession(r); err != nil {
		slog.Error("failed to validate state session", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	oauthToken, err := a.tokenExchange(r)
	if err != nil {
		slog.Error("failed to exchange token", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	userInfo, err := a.validateAndGetClaimsIDToken(r, oauthToken)
	if err != nil {
		slog.Error("failed to validate and get claims id token", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sessionID, err := generateRandomSecureString()
	if err != nil {
		slog.Error("failed to generate session id", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// Create session data
	sessionData := auth.SessionData{
		AccessToken: oauthToken.AccessToken, // From Keycloak
		UserInfo: auth.UserInfo{
			Username: userInfo.Username,
			Email:    userInfo.Email,
		},
		CreatedAt: time.Now(),
	}
	// Store session
	SessionStore[sessionID] = sessionData

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		MaxAge:   60, // A minute
		Path:     "/",
		Secure:   true, // Disabled for development
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

type oidcClaims struct {
	Email    string `json:"email"`
	Username string `json:"preferred_username"`
}

// ValidateIDToken verifies the id token from the oauth2token
func (a *OIDCAuthHandler) validateAndGetClaimsIDToken(
	r *http.Request, oauth2Token *oauth2.Token) (*oidcClaims, error) {
	// Get and validate the ID token - this proves the user's identity
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no ID token found")
	}
	// Verify the ID token
	idToken, err := a.authClient.OIDC.Verify(r.Context(), rawIDToken)
	if err != nil {
		return nil, errors.New("failed to verify ID token")
	}
	claims := oidcClaims{}
	if err := idToken.Claims(&claims); err != nil {
		return nil, errors.New("failed to get user info")
	}
	return &claims, nil
}

func (a *OIDCAuthHandler) tokenExchange(r *http.Request) (*oauth2.Token, error) {
	authorizationCode := r.Form.Get("code")
	if authorizationCode == "" {
		return nil, errors.New("authorizationCode is required")
	}
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("grant_type", "authorization_code"),
	}
	oauth2Token, err := a.authClient.Oauth.Exchange(context.Background(), authorizationCode, opts...)
	if err != nil {
		return nil, err
	}
	return oauth2Token, nil
}

func (a *OIDCAuthHandler) validateStateSession(r *http.Request) error {
	// Get state from callback parameters
	stateParam := r.Form.Get("state")
	if stateParam == "" {
		return errors.New("missing state parameter in callback")
	}

	// Validate state match
	if _, ok := authStore[stateParam]; !ok {
		return errors.New("state parameter mismatch")
	}

	// Clean up used state from store
	delete(authStore, stateParam)

	return nil
}
