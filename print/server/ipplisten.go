package server

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// sniffListener serves ipp and ipps on one port, as CUPS does: a
// connection whose first byte begins a TLS handshake (0x16) is TLS.
type sniffListener struct {
	net.Listener
	tls *tls.Config
}

// connKey is the context key of the connection of a request.
type connKey struct{}

// sniffConn is a connection whose first byte tells whether it is TLS. It
// is told on the first read, by the server's goroutine of the
// connection, not by Accept.
type sniffConn struct {
	net.Conn
	cfg    *tls.Config
	once   sync.Once
	inner  net.Conn
	br     *bufio.Reader
	secure bool
	err    error
}

func (l *sniffListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return &sniffConn{Conn: c, cfg: l.tls}, nil
}

func (c *sniffConn) sniff() {
	c.once.Do(func() {
		c.Conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		c.br = bufio.NewReader(c.Conn)
		b, err := c.br.Peek(1)
		c.Conn.SetReadDeadline(time.Time{})
		if err != nil {
			c.err = err
			return
		}
		raw := &peekedConn{Conn: c.Conn, r: c.br}
		if b[0] == 0x16 && c.cfg != nil {
			tc := tls.Server(raw, c.cfg)
			c.inner, c.secure = tc, true
			return
		}
		c.inner = raw
	})
}

func (c *sniffConn) Read(p []byte) (int, error) {
	c.sniff()
	if c.err != nil {
		return 0, c.err
	}
	return c.inner.Read(p)
}

func (c *sniffConn) Write(p []byte) (int, error) {
	c.sniff()
	if c.err != nil {
		return 0, c.err
	}
	return c.inner.Write(p)
}

func (c *sniffConn) Close() error {
	if c.inner != nil {
		return c.inner.Close()
	}
	return c.Conn.Close()
}

// peekedConn reads what was peeked first.
type peekedConn struct {
	net.Conn
	r *bufio.Reader
}

func (c *peekedConn) Read(p []byte) (int, error) { return c.r.Read(p) }

// connContext puts the connection in the context of its requests.
func connContext(ctx context.Context, c net.Conn) context.Context {
	if sc, ok := c.(*sniffConn); ok {
		return context.WithValue(ctx, connKey{}, sc)
	}
	return ctx
}

// ippCertificate is the certificate of ipps: the configuration's, or one
// the server made, signed by itself, kept next to the spool and made once
// (clients that trust it on first use refuse another).
func (x *ippServer) certificate() (tls.Certificate, error) {
	if t := x.cfg.TLS; t != nil {
		return tls.LoadX509KeyPair(t.Cert, t.Key)
	}
	path := x.certPath()
	if b, err := os.ReadFile(path); err == nil {
		return tls.X509KeyPair(b, b)
	} else if !errors.Is(err, os.ErrNotExist) {
		return tls.Certificate{}, err
	}
	b, err := selfSigned(x.host, time.Now())
	if err != nil {
		return tls.Certificate{}, err
	}
	if err := os.WriteFile(path, b, 0o600); err != nil {
		return tls.Certificate{}, fmt.Errorf("ipp: keep the certificate: %w", err)
	}
	x.s.log.Info("ipp: made a certificate signed by itself", "file", path, "host", x.host)
	return tls.X509KeyPair(b, b)
}

// certPath is the file of the certificate the server makes.
func (x *ippServer) certPath() string {
	return filepath.Join(filepath.Dir(x.s.cfg.Spool), "ipp-tls.pem")
}

// selfSigned makes a key and a certificate for a host, signed by the key,
// for ten years, with the names and addresses of the machine; the PEM
// holds the key and the certificate.
func selfSigned(host string, now time.Time) ([]byte, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 127))
	if err != nil {
		return nil, err
	}
	tmpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: host, Organization: []string{"bdf printd"}},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	names := map[string]bool{host: true, "localhost": true}
	if h, err := os.Hostname(); err == nil && h != "" {
		names[h] = true
		short, _, _ := strings.Cut(h, ".")
		names[short+".local"] = true
	}
	for n := range names {
		if ip := net.ParseIP(n); ip != nil {
			tmpl.IPAddresses = append(tmpl.IPAddresses, ip)
		} else {
			tmpl.DNSNames = append(tmpl.DNSNames, n)
		}
	}
	for _, ip := range localAddrs() {
		tmpl.IPAddresses = append(tmpl.IPAddresses, ip)
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	kb, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	out := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: kb})
	return append(out, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})...), nil
}

// localAddrs are the machine's addresses, loopback ones included.
func localAddrs() []net.IP {
	var out []net.IP
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, a := range addrs {
		if ipn, ok := a.(*net.IPNet); ok && !ipn.IP.IsLinkLocalUnicast() {
			out = append(out, ipn.IP)
		}
	}
	return out
}

// fingerprint is the SHA-256 of a certificate, as clients that pin it
// show it.
func fingerprint(c tls.Certificate) string {
	if len(c.Certificate) == 0 {
		return ""
	}
	sum := sha256.Sum256(c.Certificate[0])
	return "sha256:" + hex.EncodeToString(sum[:])
}

// serve listens for IPP (Options.IPPListener or the configuration's
// address), with ipp and ipps on the port, and serves it; the error of
// the server goes to errc.
func (x *ippServer) serve(ctx context.Context, errc chan<- error) (*http.Server, error) {
	cert, err := x.certificate()
	if err != nil {
		return nil, fmt.Errorf("ipp: %w", err)
	}
	x.fingerprint = fingerprint(cert)
	ln := x.s.o.IPPListener
	if ln == nil {
		if ln, err = net.Listen("tcp", x.cfg.Listen); err != nil {
			return nil, fmt.Errorf("ipp: %w", err)
		}
	}
	x.port = ln.Addr().(*net.TCPAddr).Port
	sl := &sniffListener{Listener: ln, tls: &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS12, NextProtos: []string{"http/1.1"}}}
	hs := &http.Server{
		Handler:           x.handler(),
		ConnContext:       connContext,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       readTimeout,
		IdleTimeout:       2 * time.Minute,
		ErrorLog:          slog.NewLogLogger(x.s.log.Handler(), slog.LevelWarn),
	}
	wctx, cancel := context.WithCancel(ctx)
	x.cancel = cancel
	go x.watchQueues(wctx)
	go func() { errc <- hs.Serve(sl) }()
	x.s.log.Info("serving IPP", "listen", ln.Addr().String(), "host", x.host, "auth", x.cfg.Auth, "ipps-only", x.secureOnly,
		"queues", len(x.order), "fingerprint", x.fingerprint)
	if err := x.advertise(wctx); err != nil {
		x.s.log.Warn("ipp: DNS-SD", "err", err)
	}
	return hs, nil
}

// stop ends what serve began.
func (x *ippServer) stop() {
	if x.cancel != nil {
		x.cancel()
	}
	x.mu.Lock()
	for _, in := range x.incoming {
		if in.upload != nil {
			in.upload.Discard()
			in.upload = nil
		}
	}
	x.mu.Unlock()
}
