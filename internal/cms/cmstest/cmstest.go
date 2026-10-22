// Package cmstest makes what tests of signatures need: a certificate
// authority of its own, certificates of keys that it issues, and a
// time-stamp authority (RFC 3161) over HTTP.
package cmstest

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"time"

	"github.com/shibukawa/bdf/internal/cms"
)

// Authority is a certificate authority.
type Authority struct {
	Key  crypto.Signer
	Cert *x509.Certificate
	next atomic.Int64
}

// Identity is a key and its certificate, with the chain to the root.
type Identity struct {
	Key   crypto.Signer
	Chain []*x509.Certificate
}

// NewAuthority makes a root certificate authority.
func NewAuthority(name string) *Authority {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: name, Organization: []string{"BDF tests"}},
		NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(24 * time.Hour),
		IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, key.Public(), key)
	if err != nil {
		panic(err)
	}
	cert, _ := x509.ParseCertificate(der)
	a := &Authority{Key: key, Cert: cert}
	a.next.Store(1)
	return a
}

// Pool is a pool of the root certificate.
func (a *Authority) Pool() *x509.CertPool {
	p := x509.NewCertPool()
	p.AddCert(a.Cert)
	return p
}

// Issue makes a key (ECDSA P-256, or RSA of 2048 bits) and a certificate
// of it for a name, valid from an hour ago for a day, for the usages.
func (a *Authority) Issue(name string, rsaKey bool, usage ...x509.ExtKeyUsage) *Identity {
	var key crypto.Signer
	var err error
	if rsaKey {
		key, err = rsa.GenerateKey(rand.Reader, 2048)
	} else {
		key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}
	if err != nil {
		panic(err)
	}
	return a.IssueKey(name, key, time.Now().Add(-time.Hour), time.Now().Add(24*time.Hour), usage...)
}

// IssueKey makes a certificate of a key, valid between two times.
func (a *Authority) IssueKey(name string, key crypto.Signer, from, to time.Time, usage ...x509.ExtKeyUsage) *Identity {
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(a.next.Add(1)), Subject: pkix.Name{CommonName: name, Organization: []string{"BDF tests"}},
		NotBefore: from, NotAfter: to, KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageContentCommitment, ExtKeyUsage: usage,
	}
	if len(usage) == 1 && usage[0] == x509.ExtKeyUsageTimeStamping {
		// (RFC 3161 §2.3: the extended key usage of a TSA is critical)
		v, _ := asn1.Marshal([]asn1.ObjectIdentifier{{1, 3, 6, 1, 5, 5, 7, 3, 8}})
		tmpl.ExtraExtensions = []pkix.Extension{{Id: asn1.ObjectIdentifier{2, 5, 29, 37}, Critical: true, Value: v}}
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, a.Cert, key.Public(), a.Key)
	if err != nil {
		panic(err)
	}
	cert, _ := x509.ParseCertificate(der)
	return &Identity{Key: key, Chain: []*x509.Certificate{cert, a.Cert}}
}

// PEM writes the certificates of an identity, and its key as PKCS #8.
func (id *Identity) PEM() (certs, key []byte) {
	for _, c := range id.Chain {
		certs = append(certs, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: c.Raw})...)
	}
	k, err := x509.MarshalPKCS8PrivateKey(id.Key)
	if err != nil {
		panic(err)
	}
	return certs, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: k})
}

// TSA is a time-stamp authority over HTTP.
type TSA struct {
	*httptest.Server
	Identity *Identity
	// Now is the time it stamps (the clock when nil).
	Now func() time.Time
	// Requests counts the requests it answered.
	Requests atomic.Int64
}

// NewTSA starts a time-stamp authority whose certificate the authority
// issues.
func NewTSA(a *Authority) *TSA {
	t := &TSA{Identity: a.Issue("BDF test TSA", false, x509.ExtKeyUsageTimeStamping)}
	t.Server = httptest.NewServer(http.HandlerFunc(t.serve))
	return t
}

func (t *TSA) serve(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<16))
	if err != nil || r.Header.Get("Content-Type") != "application/timestamp-query" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	digest, h, nonce, err := cms.ParseTimeStampRequest(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	now := time.Now()
	if t.Now != nil {
		now = t.Now()
	}
	serial := t.Requests.Add(1)
	token, err := cms.TimeStampToken(digest, h, now, big.NewInt(serial), nonce, asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 0, 1}, t.Identity.Key, t.Identity.Chain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := cms.TimeStampResponse(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/timestamp-reply")
	w.Write(resp)
}
