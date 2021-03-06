package handlers

import (
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/nayyara-cropsey/jwtmock"
)

type keyStore interface {
	GenerateNew() error
	GetJWKS() *jwk.Set
	GetSigningKey() *jwtmock.SigningKey
}

type clientRepo interface {
	Register(jwtmock.ClientRegistration) error
	GenerateToken(jwtmock.ClientTokenRequest, *jwtmock.SigningKey) (*jwtmock.ClientTokenResponse, error)
}
