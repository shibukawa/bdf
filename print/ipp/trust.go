package ipp

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// TLSPolicy says how the certificate of an ipps printer is verified.
// The zero value verifies with the roots of the system, as any HTTPS
// client does. Printers mostly have certificates of their own making,
// so Store lets the first connection record the certificate's
// fingerprint and later ones require the same (trust on first use, as
// CUPS does); Insecure takes any certificate.
type TLSPolicy struct {
	// Store records the fingerprints of certificates the roots of the
	// system do not verify, by the host of the URL, and refuses a host
	// whose certificate changed. nil keeps none: such a certificate is
	// refused.
	Store *TrustStore
	// Insecure takes any certificate. For testing, or a network that is
	// trusted as a whole.
	Insecure bool
}

// ErrCertificateChanged is returned when a host's certificate is not
// the one the store recorded.
var ErrCertificateChanged = errors.New("ipp: the certificate of the printer changed since it was recorded")

// config returns the TLS configuration of connections to host (the host of
// the printer's URL, without the port), which keys the store: a TLS
// connection to an address names no server.
func (p TLSPolicy) config(host string) *tls.Config {
	if p.Insecure {
		return &tls.Config{InsecureSkipVerify: true}
	}
	if p.Store == nil {
		return nil
	}
	store := p.Store
	cfg := &tls.Config{InsecureSkipVerify: true}
	cfg.VerifyConnection = func(cs tls.ConnectionState) error {
		if len(cs.PeerCertificates) == 0 {
			return errors.New("ipp: the printer sent no certificate")
		}
		// First the roots of the system, as the default verification.
		opts := x509.VerifyOptions{DNSName: host, Intermediates: x509.NewCertPool()}
		for _, c := range cs.PeerCertificates[1:] {
			opts.Intermediates.AddCert(c)
		}
		if _, err := cs.PeerCertificates[0].Verify(opts); err == nil {
			return nil
		}
		return store.check(host, cs.PeerCertificates[0])
	}
	return cfg
}

// TrustStore keeps the fingerprints of the certificates of printers
// that the roots of the system do not verify, by host, in a JSON file.
type TrustStore struct {
	// Path is the file: $XDG_CONFIG_HOME/bdf/printers.json by default
	// (DefaultTrustPath).
	Path string
	// ReadOnly refuses hosts not yet recorded instead of recording them.
	ReadOnly bool

	mu    sync.Mutex
	known map[string]string // host → "sha256:…"
}

// DefaultTrustPath returns the file of the user's store.
func DefaultTrustPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "bdf", "printers.json"), nil
}

// Fingerprint returns "sha256:" and the hex of the SHA-256 of the
// certificate.
func Fingerprint(c *x509.Certificate) string {
	sum := sha256.Sum256(c.Raw)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func (s *TrustStore) load() error {
	if s.known != nil {
		return nil
	}
	s.known = map[string]string{}
	b, err := os.ReadFile(s.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	var f struct {
		Printers map[string]string `json:"printers"`
	}
	if err := json.Unmarshal(b, &f); err != nil {
		return fmt.Errorf("ipp: %s: %w", s.Path, err)
	}
	for k, v := range f.Printers {
		s.known[k] = v
	}
	return nil
}

func (s *TrustStore) save() error {
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o700); err != nil {
		return err
	}
	keys := make([]string, 0, len(s.known))
	for k := range s.known {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ordered := make(map[string]string, len(keys))
	for _, k := range keys {
		ordered[k] = s.known[k]
	}
	b, err := json.MarshalIndent(struct {
		Printers map[string]string `json:"printers"`
	}{ordered}, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.Path + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.Path)
}

// check compares the certificate of a host with the recorded
// fingerprint, recording it on first use.
func (s *TrustStore) check(host string, c *x509.Certificate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Path == "" {
		p, err := DefaultTrustPath()
		if err != nil {
			return err
		}
		s.Path = p
	}
	if err := s.load(); err != nil {
		return err
	}
	fp := Fingerprint(c)
	if want, ok := s.known[host]; ok {
		if want != fp {
			return fmt.Errorf("%w: %s has %s, the store has %s", ErrCertificateChanged, host, fp, want)
		}
		return nil
	}
	if s.ReadOnly {
		return fmt.Errorf("ipp: the certificate of %s (%s) is not verified by the system nor recorded", host, fp)
	}
	s.known[host] = fp
	return s.save()
}

// Forget removes the fingerprint recorded for a host.
func (s *TrustStore) Forget(host string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Path == "" {
		p, err := DefaultTrustPath()
		if err != nil {
			return err
		}
		s.Path = p
	}
	if err := s.load(); err != nil {
		return err
	}
	delete(s.known, host)
	return s.save()
}
