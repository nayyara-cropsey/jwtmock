package service

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/nayyara-samuel/jwt-mock/types"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const (
	RS256 = "RS256"
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
		ID:            id.String(),
		Key:           key,
		SigningMethod: jwt.SigningMethodRS256,
		Algorithm:     RS256,
		PublicKey:     &key.PublicKey,
	}, nil
}
