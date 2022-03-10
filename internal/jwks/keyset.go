package jwks

import (
	// nolint:gosec //ignore warning about weak cryptographic primitive
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/nayyara-cropsey/jwtmock"
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
func (t *Generator) GenerateJWKSet() (*jwk.Set, *jwtmock.SigningKey, error) {
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
	// generate JWK public key
	key, err := jwk.New(signingKey.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("jwk: %w", err)
	}

	// encode cert as base-64 string
	certStr := base64.StdEncoding.EncodeToString(cert.Raw)

	// nolint:gosec // ignore weak cryptographic algorithm warning
	hash := sha1.New()
	hash.Write(cert.Raw)
	x5tSHA1 := hex.EncodeToString(hash.Sum(nil))

	vals := map[string]interface{}{
		jwk.KeyIDKey:              signingKey.ID,
		jwk.X509CertThumbprintKey: x5tSHA1,
		jwk.X509CertChainKey:      []string{certStr},
		jwk.KeyUsageKey:           signingUsage,
		jwk.AlgorithmKey:          signingKey.Algorithm,
	}

	for k, v := range vals {
		if err = key.Set(k, v); err != nil {
			return nil, nil, fmt.Errorf("jwk field %v: %w", k, err)
		}
	}

	keySet := &jwk.Set{Keys: []jwk.Key{key}}

	return keySet, signingKey, nil
}
