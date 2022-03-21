package jwtmock

const (
	// Bearer is a constant for the bearer token
	Bearer = "Bearer"

	// ClientCredentials is a constant for the client credentials grant type
	ClientCredentials = "client_credentials"
)

// ClientRegistration is used to register a new client for machine-to-machine auth.
type ClientRegistration struct {
	ID     string `json:"client_id"`
	Secret string `json:"client_secret"`
	Scope  string `json:"scope"`
}

// ClientTokenRequest is a request to obtain a JWT token for the given client - passed as form-URL encoded
type ClientTokenRequest struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	GrantType    string `mapstructure:"grant_type"`
	Audience     string `mapstructure:"audience"`
}

// ClientTokenClaims are claims in the JWT for the client
type ClientTokenClaims struct {
	Issuer          string `mapstructure:"iss"`
	Subject         string `mapstructure:"sub"`
	Audience        string `mapstructure:"aud"`
	IssuedAt        int64  `mapstructure:"iat"`
	Expires         int64  `mapstructure:"exp"`
	AuthorizedParty string `mapstructure:"azp"`
	Scope           string `mapstructure:"scope"`
	GrantType       string `mapstructure:"gty"`
}

// ClientTokenResponse is the response for token endpoint
type ClientTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}
