package handlers

import (
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/nayyara-cropsey/jwt-mock/types"
)

type keyStore interface {
	GetJWKS() *jwk.Set
	GenerateNew() error
	GetSigningKey() *types.SigningKey
}
