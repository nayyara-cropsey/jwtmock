package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"math/big"
	"time"
)

const (
	orgName = "JWT Mock"
)

// CertificateGenerator is used to generate certificates.
type CertificateGenerator struct {
	lifeTime time.Duration
}

// NewCertificateGenerator is the preferred way to instantiate a certificate generator.
func NewCertificateGenerator(lifeTime time.Duration) *CertificateGenerator {
	return &CertificateGenerator{lifeTime: lifeTime}
}

// CreateParent generates a X.509 parent/root certificate for use in authorization servers.
func (c *CertificateGenerator) CreateParent() (*x509.Certificate, error) {
	serialNumber, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		return nil, fmt.Errorf("serial number: %w", err)
	}

	return &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{orgName},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(c.lifeTime),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
	}, nil
}

// CreateChild - creates a child certificate with given key and parent certificate.
func (c *CertificateGenerator) CreateChild(parent *x509.Certificate, key interface{}) (*x509.Certificate, error) {
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("must be an RSA Key")
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, parent, parent, &rsaKey.PublicKey, rsaKey)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return x509.ParseCertificate(derBytes)
}
