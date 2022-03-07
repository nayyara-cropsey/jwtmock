package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/nayyara-cropsey/jwt-mock/types"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

const keyID = "kid"

var (
	// Ensure Claims complies with jwt.Claims interface
	_ jwt.Claims = Claims{}

	// ErrExpiredToken means the JWT token has expired.
	ErrExpiredToken = errors.New("token has expired (seconds)")

	// ErrBadTokenFuture means the JWT token was issued in the future.
	ErrBadTokenFuture = errors.New("token was issued in the future (seconds)")

	// ErrSubMissing means the JWT token has missing subject
	ErrSubMissing = errors.New("token subject is missing")
)

// Claims represents the type for JWT claims
type Claims map[string]interface{}

// internal type for validating require fields.
type requiredClaims struct {
	Subject   string `mapstructure:"sub"`
	IssuedAt  int64  `mapstructure:"iat"`
	ExpiredAt int64  `mapstructure:"exp"`
}

// Valid returns an error if this token is invalid
func (c Claims) Valid() error {
	r := requiredClaims{}
	err := mapstructure.Decode(c, &r)
	if err != nil {
		return fmt.Errorf("check required claims: %w", err)
	}

	expiresAt := time.Unix(r.ExpiredAt, 0)
	if time.Now().After(expiresAt) {
		return ErrExpiredToken
	}

	issuedAt := time.Unix(r.IssuedAt, 0)
	if time.Now().Before(issuedAt) {
		return ErrBadTokenFuture
	}

	if r.Subject == "" {
		return ErrSubMissing
	}

	return nil
}

// CreateToken generates a JWT token using the provided claims and signing key.
func CreateToken(claims Claims, signingKey *types.SigningKey) (string, error) {
	if err := claims.Valid(); err != nil {
		return "", fmt.Errorf("validation: %w", err)
	}

	jwtToken := jwt.NewWithClaims(signingKey.SigningMethod, claims)
	jwtToken.Header[keyID] = signingKey.ID

	return jwtToken.SignedString(signingKey.Key)
}
