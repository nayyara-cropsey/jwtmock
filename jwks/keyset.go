package jwks

import (
	//nolint:gosec // ignore weak cryptographic algorithm warning
	"crypto/sha1"
	"crypto/x509"
	"fmt"

	"github.com/nayyara-cropsey/jwt-mock/types"

	"gopkg.in/square/go-jose.v2"
)

const signingUsage = "sig" // key is used for signing only

// Generator is atype for generating JWKS with a single singing key.
type Generator struct {
	certGen certGenerator
	keyGen  keyGenerator
	keyLen  int
}

// NewGenerator is the preferred way to instantiate a key generator.
func NewGenerator(certGen certGenerator, keyGen keyGenerator, keyLen int) *Generator {
	return &Generator{
		certGen: certGen,
		keyGen:  keyGen,
		keyLen:  keyLen,
	}
}

// GenerateJWKSet generates a JSON web key set for use in signing tokens.
func (t *Generator) GenerateJWKSet() (*jose.JSONWebKeySet, *types.SigningKey, error) {
	signingKey, err := t.keyGen.GenerateKey(t.keyLen)
	if err != nil {
		return nil, nil, fmt.Errorf("generate key: %w", err)
	}

	parentCert, err := t.certGen.CreateParent()
	if err != nil {
		return nil, nil, fmt.Errorf("parent cert: %w", err)
	}

	cert, err := t.certGen.CreateChild(parentCert, signingKey.Key)
	if err != nil {
		return nil, nil, fmt.Errorf("cert: %w", err)
	}

	// nolint:gosec // ignore weak cryptographic algorithm warning
	x5tSHA1 := sha1.Sum(cert.Raw)

	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				KeyID:                     signingKey.ID,
				Key:                       signingKey.PublicKey,
				CertificateThumbprintSHA1: x5tSHA1[:],
				Certificates:              []*x509.Certificate{cert},
				Use:                       signingUsage,
				Algorithm:                 signingKey.Algorithm,
			},
		},
	}, signingKey, nil
}
