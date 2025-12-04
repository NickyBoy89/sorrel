package auth

import (
	"context"
	"fmt"

	"git.nicholasnovak.io/recipe_planning/backend/internal/config"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// See https://medium.com/@rizkysr90/getting-started-with-oauth-2-0-in-golang-using-keycloak-8e61fdf3620b for a lot more inspiration

// Client struct holds all components needed for authentication
type Client struct {
	Provider *oidc.Provider        // Handles OIDC protocol operations with Keycloak
	OIDC     *oidc.IDTokenVerifier // Verifies JWT tokens from Keycloak
	Oauth    oauth2.Config         // Manages OAuth2 flow (authorization codes, tokens)
}

func New(ctx context.Context, config *config.Config) (*Client, error) {
	if config.KeycloakHostname == "" ||
		config.KeycloakRealm == "" ||
		config.KeycloakClientId == "" ||
		config.KeycloakClientSecret == "" {
		return nil, fmt.Errorf("config was missing at least one key to start Keycloak")
	}

	if config.RedirectURL == "" {
		return nil, fmt.Errorf("redirect URL was left empty")
	}

	// Construct the provider URL using Keycloak realm
	providerURL := fmt.Sprintf("%s/realms/%s", config.KeycloakHostname, config.KeycloakRealm)

	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}

	// Create ID token verifier
	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.KeycloakClientId,
	})

	// Configure an OpenID Connect aware OAuth2 client with specific scopes:
	// - oidc.ScopeOpenID: Required for OpenID Connect authentication, provides subject ID (sub)
	// - "roles": Keycloak-specific scope to get user roles in the token
	oauth2Config := oauth2.Config{
		ClientID:     config.KeycloakClientId,
		ClientSecret: config.KeycloakClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes: []string{
			oidc.ScopeOpenID, // Required for OIDC authentication
			"roles",          // Request user roles from Keycloak
		},
	}

	// Return initialized client with all required components
	return &Client{
		// oauth2Config: Used for OAuth2 operations like:
		// - Generating login URL (AuthCodeURL)
		// - Exchanging auth code for tokens (Exchange)
		// - Managing token refresh
		Oauth: oauth2Config,

		// verifier: Used to validate tokens:
		// - Verifies JWT signature
		// - Validates token claims (exp, iss, aud)
		// - Extracts user information
		OIDC: verifier,

		// provider: Keycloak OIDC provider that:
		// - Provides endpoint URLs (auth, token)
		// - Handles OIDC protocol details
		// - Manages provider metadata
		Provider: provider,
	}, nil
}

// AuthCodeURL generates the login URL for OAuth2 authorization code flow.
// It returns a URL that the user should be redirected to for authentication.
// The state parameter is a random string that will be validated in the callback
// to prevent CSRF attacks.
func (c *Client) AuthCodeURL(state string) string {
	return c.Oauth.AuthCodeURL(state)
}
