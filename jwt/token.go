package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/nayyara-cropsey/jwt-mock/types"

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

	headers := jws.NewHeaders()
	if err := headers.Set(jws.KeyIDKey, signingKey.ID); err != nil {
		return "", fmt.Errorf("JWS headers key: %w", err)
	}

	options := []jwt.Option{jwt.WithHeaders(headers)}
	for k, v := range claims {
		options = append(options, jwt.WithClaimValue(k, v))
	}

	token := jwt.New()
	signedToken, err := jwt.Sign(token, signingKey.Algorithm, signingKey.Key, options...)
	if err != nil {
		return "", fmt.Errorf("sign: %w", err)
	}

	return string(signedToken), nil
}
