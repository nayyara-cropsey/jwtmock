package types

import (
	//nolint:gosec // ignore weak cryptographic algorithm warning
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
)

// SigningKey represents a generic key used to sign JWTs.
type SigningKey struct {
	ID        string
	Key       interface{}
	Algorithm jwa.SignatureAlgorithm
	PublicKey interface{}
}

// GenerateJWK generates a JWK from the signing key
func (s *SigningKey) GenerateJWK(usage string, cert *x509.Certificate) (jwk.Key, error) {
	// generate JWK public key
	key, err := jwk.New(s.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("jwk: %w", err)
	}

	// encode cert as base-64 string
	certStr := base64.StdEncoding.EncodeToString(cert.Raw)

	// nolint:gosec // ignore weak cryptographic algorithm warning
	hash := sha1.New()
	hash.Write(cert.Raw)
	x5tSHA1 := hex.EncodeToString(hash.Sum(nil))

	vals := map[string]interface{}{
		jwk.KeyIDKey:              s.ID,
		jwk.X509CertThumbprintKey: x5tSHA1,
		jwk.X509CertChainKey:      []string{certStr},
		jwk.KeyUsageKey:           usage,
		jwk.AlgorithmKey:          s.Algorithm,
	}

	for k, v := range vals {
		if err = key.Set(k, v); err != nil {
			return nil, fmt.Errorf("jwk field %v: %w", k, err)
		}
	}

	return key, nil
}
