package jwks

import (
	"crypto/x509"
	"jwt-mock/types"
)

type certGenerator interface {
	CreateParent() (*x509.Certificate, error)
	CreateChild(parent *x509.Certificate, key interface{}) (*x509.Certificate, error)
}

type keyGenerator interface {
	GenerateKey(length int) (*types.SigningKey, error)
}
