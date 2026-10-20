package ipp

import (
	"crypto/x509"
	"errors"
	"testing"
)

type memTrust struct{ m map[string]string }

func (b *memTrust) Load() (map[string]string, error) { return b.m, nil }
func (b *memTrust) Save(m map[string]string) error   { b.m = m; return nil }

func TestTrustBackendAndPinned(t *testing.T) {
	a := &x509.Certificate{Raw: []byte("certificate a")}
	b := &x509.Certificate{Raw: []byte("certificate b")}
	be := &memTrust{}
	s := &TrustStore{Backend: be}
	if err := s.check("mfp1", a); err != nil {
		t.Fatal(err)
	}
	if be.m["mfp1"] != Fingerprint(a) {
		t.Fatalf("backend has %v", be.m)
	}
	if err := (&TrustStore{Backend: be}).check("mfp1", b); !errors.Is(err, ErrCertificateChanged) {
		t.Fatalf("changed certificate: %v", err)
	}
	pinned := &TrustStore{Backend: be, Pinned: map[string]string{"mfp2": Fingerprint(b)}}
	if err := pinned.check("mfp2", b); err != nil {
		t.Fatal(err)
	}
	if err := pinned.check("mfp2", a); !errors.Is(err, ErrCertificateChanged) {
		t.Fatalf("pinned: %v", err)
	}
	if _, ok := be.m["mfp2"]; ok {
		t.Fatal("a pinned host was recorded")
	}
}
