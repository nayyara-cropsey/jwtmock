package service

import (
	"sync"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/nayyara-cropsey/jwtmock"
	"github.com/nayyara-cropsey/jwtmock/internal/jwks"
)

// KeyStore is used to keep state about current JWKS and signing key.
type KeyStore struct {
	generator *jwks.Generator

	key    *jwtmock.SigningKey
	jwkSet *jwk.Set

	m sync.Mutex
}

// NewKeyStore is the preferred way to instantiate a key store.
func NewKeyStore(generator *jwks.Generator) (*KeyStore, error) {
	k := &KeyStore{
		generator: generator,
	}

	if err := k.GenerateNew(); err != nil {
		return nil, err
	}

	return k, nil
}

// GenerateNew generates a new pair of JWKS and signing key.
func (k *KeyStore) GenerateNew() error {
	k.m.Lock()
	defer k.m.Unlock()

	jwkSet, key, err := k.generator.GenerateJWKSet()
	if err != nil {
		return err
	}

	k.jwkSet = jwkSet
	k.key = key

	return nil
}

// GetJWKS returns the currently stored JWKS.
func (k *KeyStore) GetJWKS() *jwk.Set {
	return k.jwkSet
}

// GetSigningKey returns the currently stored signing key.
func (k *KeyStore) GetSigningKey() *jwtmock.SigningKey {
	return k.key
}
