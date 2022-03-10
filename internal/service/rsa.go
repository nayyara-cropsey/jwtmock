package service

import (
	"crypto/rand"
	"crypto/rsa"
	mrand "math/rand"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/nayyara-cropsey/jwtmock"
)

// idLen is the ID length
const idLen = 16

// ID runes contains characters for generating an ID
var idRunes = []rune("123456abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RSAKeyGenerator generates key IDs and keys.
type RSAKeyGenerator struct{}

// NewRSAKeyGenerator is the preferred way to create a RSA key generator.
func NewRSAKeyGenerator() *RSAKeyGenerator {
	mrand.Seed(time.Now().UnixNano())

	return &RSAKeyGenerator{}
}

// GenerateKey generates a RSA signing key.
func (k *RSAKeyGenerator) GenerateKey(length int) (*jwtmock.SigningKey, error) {
	id := generateID(idLen)
	key, err := rsa.GenerateKey(rand.Reader, length)
	if err != nil {
		return nil, err
	}

	return &jwtmock.SigningKey{
		ID:        id,
		Key:       key,
		Algorithm: jwa.RS256,
		PublicKey: &key.PublicKey,
	}, nil
}

// generateID generates an random string of the given length
func generateID(n int) string {
	b := make([]rune, n)
	for i := range b {
		// nolint:gosec // ignore weak rand warning
		b[i] = idRunes[mrand.Intn(len(idRunes))]
	}
	return string(b)
}
