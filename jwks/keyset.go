package jwks

import (
	"fmt"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/nayyara-cropsey/jwt-mock/types"
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
func (t *Generator) GenerateJWKSet() (*jwk.Set, *types.SigningKey, error) {
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

	// generate JWK from signing key
	key, err := signingKey.GenerateJWK(signingUsage, cert)
	if err != nil {
		return nil, nil, fmt.Errorf("jwk: %w", err)
	}

	keySet := &jwk.Set{Keys: []jwk.Key{key}}

	return keySet, signingKey, nil
}
