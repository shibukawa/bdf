package fakeprinter

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Server is a printer (or a CUPS) served over IPP on a port.
type Server struct {
	// URL is the printer's URL (ipp://host:port/ipp/print, or ipps://),
	// or the CUPS server's (ipp://host:port).
	URL  string
	Addr string
	srv  *http.Server
	ln   net.Listener
}

// Close stops the server.
func (s *Server) Close() error { return s.srv.Close() }

// ServeIPP serves the printer over IPP at addr ("127.0.0.1:0" for any
// free port of the loopback), at /ipp/print.
func (p *Printer) ServeIPP(addr string) (*Server, error) {
	return serve(p, addr, nil, "/ipp/print")
}

// ServeIPPS serves the printer over ipps at addr with a self-signed
// certificate for localhost (a fault can change it: NewCertificate,
// ExpiredCertificate).
func (p *Printer) ServeIPPS(addr string) (*Server, error) {
	p.lock()
	c := p.cert()
	p.mu.Unlock()
	if c.cur.Load() == nil {
		if err := c.renew(p.certHosts(), false); err != nil {
			return nil, err
		}
	}
	return serve(p, addr, &tls.Config{GetCertificate: c.get}, "/ipp/print")
}

// Certificate is the certificate the printer presents over ipps now, or
// nil.
func (p *Printer) Certificate() *x509.Certificate {
	p.lock()
	c := p.cert()
	p.mu.Unlock()
	cur := c.cur.Load()
	if cur == nil {
		return nil
	}
	cert, _ := x509.ParseCertificate(cur.Certificate[0])
	return cert
}

// Serve serves the CUPS server at addr.
func (c *CUPS) Serve(addr string) (*Server, error) { return serve(c, addr, nil, "") }

func serve(h http.Handler, addr string, tc *tls.Config, path string) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	scheme := "ipp"
	if tc != nil {
		scheme = "ipps"
		ln = tls.NewListener(ln, tc)
	}
	s := &Server{URL: scheme + "://" + ln.Addr().String() + path, Addr: ln.Addr().String(), ln: ln,
		srv: &http.Server{Handler: h, ReadHeaderTimeout: 10 * time.Second, ErrorLog: quietLog()}}
	go s.srv.Serve(ln)
	return s, nil
}

// Start serves the printer over IPP on the loopback until the test ends,
// and returns its URL.
func Start(t testing.TB, p *Printer) string {
	t.Helper()
	s, err := p.ServeIPP("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s.URL
}

// StartTLS serves the printer over ipps on the loopback until the test
// ends, and returns its URL.
func StartTLS(t testing.TB, p *Printer) string {
	t.Helper()
	s, err := p.ServeIPPS("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s.URL
}

// StartSocket serves the printer on a raw port (9100) of the loopback
// until the test ends.
func StartSocket(t testing.TB, p *Printer) *Listener {
	t.Helper()
	l, err := p.ListenSocket("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(l.Close)
	return l
}

// StartLPD serves the printer as an LPD queue on the loopback until the
// test ends.
func StartLPD(t testing.TB, p *Printer, queue string) *Listener {
	t.Helper()
	l, err := p.ListenLPD("127.0.0.1:0", queue)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(l.Close)
	return l
}

// Raw is a printer served on a raw transport: what it received, how it
// fails (the Printer's), and where it is (the Listener's URL).
type Raw struct {
	*Printer
	*Listener
}

// NewSocket starts a printer that takes anything on a raw port (9100) of
// the loopback, until the test ends.
func NewSocket(t testing.TB) *Raw {
	t.Helper()
	p := New()
	p.Name = "socket"
	return &Raw{p, StartSocket(t, p)}
}

// NewLPD starts a printer that takes anything as an LPD queue on the
// loopback, until the test ends.
func NewLPD(t testing.TB, queue string) *Raw {
	t.Helper()
	p := New()
	p.Name = queue
	return &Raw{p, StartLPD(t, p, queue)}
}

// StartSerial makes the printer read a pseudo-terminal until the test
// ends, and returns the terminal's name (what serial:// names); the test
// is skipped where there are none.
func StartSerial(t testing.TB, p *Printer) string {
	t.Helper()
	far, name, err := OpenPTY()
	if err != nil {
		t.Skip("no pseudo-terminals:", err)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		p.ServeSerial(far, 300*time.Millisecond)
	}()
	t.Cleanup(func() {
		far.Close()
		<-done
	})
	return name
}

// StartCUPS serves the CUPS server on the loopback until the test ends,
// and returns its URL (ipp://127.0.0.1:port).
func StartCUPS(t testing.TB, c *CUPS) string {
	t.Helper()
	s, err := c.Serve("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s.URL
}

// writeFile writes a document a printer kept into its directory.
func writeFile(dir, name string, b []byte) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, name), b, 0o644)
}

// quietLog drops the server's own complaints (a client that refuses the
// certificate, a connection a fault cut): the tests look at what the
// clients see.
func quietLog() *log.Logger { return log.New(io.Discard, "", 0) }
