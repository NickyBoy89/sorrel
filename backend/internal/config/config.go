package config

type Config struct {
	PublicKey            string `json:"public_key"`
	PrivateKey           string `json:"private_key"`
	KeycloakHostname     string `json:"keycloak_hostname,omitempty"`
	KeycloakRealm        string `json:"keycloak_realm,omitempty"`
	KeycloakClientId     string `json:"keycloak_client_id,omitempty"`
	KeycloakClientSecret string `json:"keycloak_client_secret,omitempty"`
	RedirectURL          string `json:"redirect_url,omitempty"`
}
