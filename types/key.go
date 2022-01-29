package types

import "github.com/dgrijalva/jwt-go"

// SigningKey represents a generic key used to sign JWTs.
type SigningKey struct {
	ID            string
	Key           interface{}
	SigningMethod jwt.SigningMethod
	Algorithm     string
	PublicKey     interface{}
}
