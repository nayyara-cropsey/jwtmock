package jwks

import (
	"crypto/x509"

	"github.com/nayyara-cropsey/jwtmock"
)

type certGenerator interface {
	CreateParent() (*x509.Certificate, error)
	CreateChild(parent *x509.Certificate, key interface{}) (*x509.Certificate, error)
}

type keyGenerator interface {
	GenerateKey(length int) (*jwtmock.SigningKey, error)
}
