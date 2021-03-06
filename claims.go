package jwtmock

import (
	"errors"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/mitchellh/mapstructure"
)

var (
	// ErrExpiredToken means the JWT token has expired.
	ErrExpiredToken = errors.New("token has expired (seconds)")

	// ErrBadTokenFuture means the JWT token was issued in the future.
	ErrBadTokenFuture = errors.New("token was issued in the future (seconds)")

	// ErrSubMissing means the JWT token has missing subject
	ErrSubMissing = errors.New("token subject is missing")
)

// SigningKey represents a generic key used to sign JWTs.
type SigningKey struct {
	ID        string
	Key       interface{}
	Algorithm jwa.SignatureAlgorithm
	PublicKey interface{}
}

// Claims represents the type for JWT claims
type Claims map[string]interface{}

// ClaimsFrom generates a claims object form the given struct
func ClaimsFrom(v interface{}) (Claims, error) {
	var claims Claims

	err := mapstructure.Decode(v, &claims)
	return claims, err
}

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

// CreateJWT generates a JWT token using the provided claims and signing key.
func (c Claims) CreateJWT(signingKey *SigningKey) (string, error) {
	if err := c.Valid(); err != nil {
		return "", fmt.Errorf("validation: %w", err)
	}

	headers := jws.NewHeaders()
	if err := headers.Set(jws.KeyIDKey, signingKey.ID); err != nil {
		return "", fmt.Errorf("JWS headers key: %w", err)
	}

	token := jwt.New()
	for k, v := range c {
		if err := token.Set(k, v); err != nil {
			return "", fmt.Errorf("set claim %v: %w", k, err)
		}
	}

	signedToken, err := jwt.Sign(token, signingKey.Algorithm, signingKey.Key, jwt.WithHeaders(headers))
	if err != nil {
		return "", fmt.Errorf("sign: %w", err)
	}

	return string(signedToken), nil
}
