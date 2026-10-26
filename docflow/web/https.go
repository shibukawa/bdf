package web

// The devices of displays. A device makes its key with WebCrypto, which
// a browser gives only to a secure context: https, or localhost. A
// tablet on the network reaches this computer at its address, where
// plain http is not one, so it cannot enroll there; https_port serves
// the same pages and endpoints over TLS, and an enrollment's URL goes
// there, at an address the device reaches.

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/popcornweb/pw"
)

// devCertFile is the development certificate's file in the user's cache
// (bdf/), with its key: one for every Docflow of the user, so that a
// device that trusted it once trusts it after a restart, and at another
// worktree's port.
const devCertFile = "docflow-dev-tls.pem"

// checkHTTPS refuses settings of the devices' origin and of the HTTPS
// listener that cannot work.
func checkHTTPS(cfg Config) error {
	if cfg.DeviceURL != "" {
		u, err := url.Parse(cfg.DeviceURL)
		if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" || strings.Trim(u.Path, "/") != "" || u.RawQuery != "" {
			return fmt.Errorf("docflow: docflow.device_url %q: want an origin, such as https://docflow.example.com", cfg.DeviceURL)
		}
	}
	switch {
	case cfg.HTTPSPort == 0:
		return nil
	case cfg.HTTPSPort < 0 || cfg.HTTPSPort > 65535:
		return fmt.Errorf("docflow: docflow.https_port %d is not a port", cfg.HTTPSPort)
	case cfg.Serverless:
		return errors.New("docflow: docflow.https_port: a serverless deployment has no listeners of its own")
	case (cfg.TLSCertFile == "") != (cfg.TLSKeyFile == ""):
		return errors.New("docflow: docflow.tls_cert_file and docflow.tls_key_file go together")
	case cfg.TLSCertFile == "" && !pw.Development():
		return errors.New("docflow: docflow.https_port needs docflow.tls_cert_file and docflow.tls_key_file outside development")
	}
	return nil
}

// Outermost is a middleware for the outermost slot of the framework's
// chain (main registers it): it keeps the chain inside it, which the
// HTTPS listener serves, so that a request over TLS gets what one over
// http does (the public tree, the sessions, the access log).
func (a *App) Outermost(next http.Handler) http.Handler {
	a.chain.Store(&next)
	return next
}

// StartHTTPS listens on docflow.https_port, when it is set, until stop:
// with tls_cert_file and tls_key_file, or in development with a
// certificate of this computer signed by itself. A development run takes
// the next free port when that one is taken, as the framework does with
// the HTTP one.
func (a *App) StartHTTPS() (stop func(), err error) {
	if a.cfg.HTTPSPort == 0 {
		return func() {}, nil
	}
	var cert tls.Certificate
	if a.cfg.TLSCertFile != "" {
		cert, err = tls.LoadX509KeyPair(a.cfg.TLSCertFile, a.cfg.TLSKeyFile)
	} else {
		dir, derr := os.UserCacheDir()
		if derr == nil {
			dir = filepath.Join(dir, "bdf")
		}
		cert, err = devCertificate(dir, a.lanAddr(), time.Now())
	}
	if err != nil {
		return nil, fmt.Errorf("docflow: https_port: %w", err)
	}
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(a.cfg.HTTPSPort))
	for port := a.cfg.HTTPSPort + 1; err != nil && pw.Development() && port <= a.cfg.HTTPSPort+10 && port <= 65535; port++ {
		if shifted, serr := net.Listen("tcp", ":"+strconv.Itoa(port)); serr == nil {
			a.log.Warn("docflow: https_port is taken; this development run listens on the next free one", "configured", a.cfg.HTTPSPort, "port", port)
			ln, err = shifted, nil
		}
	}
	if err != nil {
		return nil, fmt.Errorf("docflow: https_port: %w", err)
	}
	a.httpsPort = ln.Addr().(*net.TCPAddr).Port
	srv := &http.Server{
		Handler: http.HandlerFunc(a.serveChain), ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 2 * time.Minute,
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS12},
		// a device that does not trust the certificate yet refuses each
		// handshake
		ErrorLog: slog.NewLogLogger(a.log.Handler(), slog.LevelDebug),
	}
	go func() {
		if err := srv.ServeTLS(ln, "", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Warn("docflow: https_port", "err", err)
		}
	}()
	host := "localhost"
	if ip := a.lanAddr(); ip != "" {
		host = ip
	}
	a.log.Info("docflow: HTTPS for the devices", "url", "https://"+net.JoinHostPort(host, strconv.Itoa(a.httpsPort)), "self_signed", a.cfg.TLSCertFile == "")
	return func() { srv.Close() }, nil
}

// serveChain serves a request of the HTTPS listener as the HTTP one
// does; before the framework built its chain, there is nothing to serve.
func (a *App) serveChain(w http.ResponseWriter, r *http.Request) {
	if h := a.chain.Load(); h != nil {
		(*h).ServeHTTP(w, r)
		return
	}
	http.Error(w, "starting", http.StatusServiceUnavailable)
}

// deviceOrigin is the origin an enrollment's URL opens on a device:
// docflow.device_url, else this server as the request reached it, with
// two changes for a device on the network: a loopback host (the admin's
// browser on this computer) becomes this computer's address, and plain
// http becomes https_port when there is one, for the device's key.
func (a *App) deviceOrigin(r *http.Request) string {
	if a.cfg.DeviceURL != "" {
		return strings.TrimSuffix(a.cfg.DeviceURL, "/")
	}
	scheme := "http"
	if r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
		scheme = "https"
	}
	host, port, err := net.SplitHostPort(r.Host)
	if err != nil {
		host, port = strings.Trim(r.Host, "[]"), ""
	}
	if loopback(host) {
		if ip := a.lanAddr(); ip != "" {
			host = ip
		}
	}
	if scheme == "http" && a.httpsPort != 0 {
		scheme, port = "https", strconv.Itoa(a.httpsPort)
	}
	switch {
	case port != "":
		host = net.JoinHostPort(host, port)
	case strings.Contains(host, ":"):
		host = "[" + host + "]"
	}
	return scheme + "://" + host
}

// lanAddress is this computer's IPv4 address on its network, the one a
// device there reaches it at: the address its default route leaves from
// (of a UDP socket "connected" to a documentation address, which sends
// nothing), else the first private one of an interface; "" when there is
// none.
func lanAddress() string {
	if c, err := net.Dial("udp4", "192.0.2.1:9"); err == nil {
		ip := c.LocalAddr().(*net.UDPAddr).IP
		c.Close()
		if ip.IsGlobalUnicast() {
			return ip.String()
		}
	}
	addrs, _ := net.InterfaceAddrs()
	for _, ad := range addrs {
		if n, ok := ad.(*net.IPNet); ok && n.IP.To4() != nil && n.IP.IsPrivate() {
			return n.IP.String()
		}
	}
	return ""
}

// devCertificate is a certificate of this computer for development,
// signed by itself, for the names and addresses a device may reach it
// at. It is kept in dir (none: made for this run), and made again when
// it ends within 30 days or does not name lan (lanAddress) any more.
func devCertificate(dir, lan string, now time.Time) (tls.Certificate, error) {
	path := ""
	if dir != "" {
		path = filepath.Join(dir, devCertFile)
		if b, err := os.ReadFile(path); err == nil {
			if c, err := tls.X509KeyPair(b, b); err == nil && c.Leaf.NotAfter.After(now.Add(30*24*time.Hour)) && (lan == "" || c.Leaf.VerifyHostname(lan) == nil) {
				return c, nil
			}
		}
	}
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return tls.Certificate{}, err
	}
	dns, ips := devNames()
	tpl := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: "BDF Docflow (development)"},
		// 397 days: what Apple's devices take of a server's certificate
		NotBefore: now.Add(-time.Hour), NotAfter: now.Add(397 * 24 * time.Hour),
		KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true, DNSNames: dns, IPAddresses: ips,
	}
	der, err := x509.CreateCertificate(rand.Reader, tpl, tpl, &key.PublicKey, key)
	if err != nil {
		return tls.Certificate{}, err
	}
	pk, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return tls.Certificate{}, err
	}
	b := append(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pk})...)
	if path != "" {
		if err := writePrivate(path, b); err != nil {
			slog.Warn("docflow: the development certificate is not kept; a device trusts it again after a restart", "err", err)
		}
	}
	return tls.X509KeyPair(b, b)
}

// devNames are the names and addresses of this computer: localhost, its
// host name, the host name of its DNS-SD services (ports.advertise), and
// the addresses of its interfaces.
func devNames() (dns []string, ips []net.IP) {
	dns = []string{"localhost"}
	if h, err := os.Hostname(); err == nil && h != "" {
		dns = append(dns, h)
		if !strings.Contains(h, ".") {
			dns = append(dns, h+".local")
		}
	}
	dns = append(dns, machineName()+"-docflow.local")
	ips = []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback}
	addrs, _ := net.InterfaceAddrs()
	for _, ad := range addrs {
		if n, ok := ad.(*net.IPNet); ok && n.IP.IsGlobalUnicast() {
			ips = append(ips, n.IP)
		}
	}
	return dns, ips
}

// writePrivate writes a file only its owner reads, whole or not at all
// (another process may read it at the same time).
func writePrivate(path string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".docflow-dev-tls-*")
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Rename(f.Name(), path)
	}
	if err != nil {
		os.Remove(f.Name())
	}
	return err
}
