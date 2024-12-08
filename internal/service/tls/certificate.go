package tls

import (
	"crypto/tls"
	"errors"

	"github.com/monkeydioude/heyo/pkg/tiger/assert"
)

func NewCertFromCertKey(certFile, keyFile string) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	assert.NoError(errors.New("NewCertFromCertKey"), err)
	return cert
}
