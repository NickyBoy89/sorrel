package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"git.nicholasnovak.io/recipe_planning/backend/internal/auth"
	"git.nicholasnovak.io/recipe_planning/backend/internal/handlers"
	"github.com/coreos/go-oidc/v3/oidc"
)

type AuthMiddleware interface {
	RequireAuth(handler http.Handler) http.HandlerFunc
}

type OIDCAuthMiddleware struct {
	authClient *auth.Client
}

// NewAuthMiddleware creates a new authentication middleware with OIDC verification
func NewAuthMiddleware(c context.Context,
	authClient *auth.Client,
) *OIDCAuthMiddleware {
	return &OIDCAuthMiddleware{
		authClient: authClient,
	}
}
func (m *OIDCAuthMiddleware) RequireAuth(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session from cookie
		sessionID, err := r.Cookie("session_id")
		if err != nil {
			slog.Debug("session_id cookie was not found", "error", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		// Get session data from Redis
		sessionData, ok := handlers.SessionStore[sessionID.Value]
		if !ok {
			// Clear invalid session cookie
			slog.Debug("session was not found in store, resetting connection")
			r.AddCookie(&http.Cookie{Name: "session_id", MaxAge: -1, Path: "/", Secure: false, HttpOnly: true})
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		// Verify the access token using the OIDC provider
		token, err := m.authClient.Provider.Verifier(&oidc.Config{
			SkipClientIDCheck: true, // Access tokens don't require client ID check
		}).Verify(r.Context(), sessionData.AccessToken)

		if err != nil {
			slog.Debug("error verifying access token", "error", err)
			// The token is invalid - let's clean up and redirect
			delete(handlers.SessionStore, sessionID.Value)
			r.AddCookie(&http.Cookie{Name: "session_id", MaxAge: -1, Path: "/", Secure: false, HttpOnly: true})
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		// Extract claims from the token
		var claims map[string]any
		if err := token.Claims(&claims); err != nil {
			slog.Error("failed to validate claims", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Store the validated claims and session in the context

		// c.Set("user_session", sessionData)
		// c.Set("user_claims", claims)

		handler.ServeHTTP(w, r)
	}
}
