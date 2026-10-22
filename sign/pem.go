package sign

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// ParsePEM reads a private key (PKCS #8, EC or PKCS #1, in PEM) and the
// certificates of its chain (PEM, the key's first), as the bdf command and
// the manager of documents take them.
func ParsePEM(keyPEM, certPEM []byte) (crypto.Signer, []*x509.Certificate, error) {
	var key crypto.Signer
	for rest := keyPEM; key == nil; {
		var blk *pem.Block
		if blk, rest = pem.Decode(rest); blk == nil {
			return nil, nil, errors.New("sign: no private key in PEM (a PKCS #12 file is written as PEM by openssl pkcs12 -in key.p12 -nodes)")
		}
		var k any
		var err error
		switch blk.Type {
		case "PRIVATE KEY":
			k, err = x509.ParsePKCS8PrivateKey(blk.Bytes)
		case "EC PRIVATE KEY":
			k, err = x509.ParseECPrivateKey(blk.Bytes)
		case "RSA PRIVATE KEY":
			k, err = x509.ParsePKCS1PrivateKey(blk.Bytes)
		case "ENCRYPTED PRIVATE KEY":
			return nil, nil, errors.New("sign: the key is encrypted; write it without its password (openssl pkey -in key.pem -out plain.pem)")
		default:
			continue
		}
		if err != nil {
			return nil, nil, fmt.Errorf("sign: the key: %w", err)
		}
		s, ok := k.(crypto.Signer)
		if !ok {
			return nil, nil, errors.New("sign: the key does not sign")
		}
		key = s
	}
	certs, err := ParseCertificates(certPEM)
	if err != nil {
		return nil, nil, err
	}
	if len(certs) == 0 {
		return nil, nil, errors.New("sign: no certificate in PEM")
	}
	return key, certs, nil
}

// ParseCertificates reads the certificates of PEM.
func ParseCertificates(b []byte) ([]*x509.Certificate, error) {
	var out []*x509.Certificate
	for rest := b; ; {
		var blk *pem.Block
		if blk, rest = pem.Decode(rest); blk == nil {
			break
		}
		if blk.Type != "CERTIFICATE" {
			continue
		}
		c, err := x509.ParseCertificate(blk.Bytes)
		if err != nil {
			return nil, fmt.Errorf("sign: a certificate: %w", err)
		}
		out = append(out, c)
	}
	return out, nil
}
