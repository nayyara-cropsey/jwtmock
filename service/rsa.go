package service

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/nayyara-cropsey/jwt-mock/types"

	"github.com/google/uuid"
)

// RSAKeyGenerator generates key IDs and keys.
type RSAKeyGenerator struct{}

// NewRSAKeyGenerator is the preferred way to create a RSA key generator.
func NewRSAKeyGenerator() *RSAKeyGenerator {
	return &RSAKeyGenerator{}
}

// GenerateKey generates a RSA signing key.
func (k *RSAKeyGenerator) GenerateKey(length int) (*types.SigningKey, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("generate ID: %w", err)
	}

	key, err := rsa.GenerateKey(rand.Reader, length)
	if err != nil {
		return nil, err
	}

	return &types.SigningKey{
		ID:        id.String(),
		Key:       key,
		Algorithm: jwa.RS256,
		PublicKey: &key.PublicKey,
	}, nil
}
