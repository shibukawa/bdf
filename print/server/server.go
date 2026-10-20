package server

import (
	"context"
	_ "embed"
	"errors"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/shibukawa/bdf/print/spool"
)

// Options are what a server takes besides its configuration.
type Options struct {
	// Converter converts the documents that are not .bdf; nil takes .bdf
	// only.
	Converter spool.Converter
	// FontFS adds fonts to the configuration's.
	FontFS fs.FS
	// Log receives what happens; nothing by default.
	Log *slog.Logger
	// Grace is how long Serve lets the documents being sent finish when
	// it stops: 30 seconds by default.
	Grace time.Duration
	// Now is the clock; time.Now by default.
	Now func() time.Time
	// IPPListener is the listener of the IPP server (Config.IPP) for
	// Serve; nil listens on the configuration's address.
	IPPListener net.Listener
}

// readTimeout is how long a request may take to arrive (tests shorten
// it).
var readTimeout = time.Minute

// Server is a print server: a spool and its HTTP API.
type Server struct {
	cfg  *Config
	o    Options
	sp   *spool.Spool
	log  *slog.Logger
	auth *authenticator
	// uploads and streams bound the submissions and the event streams at
	// once.
	uploads chan struct{}
	streams chan struct{}
	// done is closed when the server stops: the event streams end.
	done     chan struct{}
	stopOnce sync.Once
	handler  http.Handler
	// ipp shows the queues as IPP printers, when configured.
	ipp *ippServer
}

// New opens the spool of a configuration (spool.db, locked: a second
// server of the same file is refused) and makes the server. Serve starts
// it; Close closes a server that was not served.
func New(ctx context.Context, cfg *Config, o Options) (*Server, error) {
	if err := cfg.check(); err != nil {
		return nil, err
	}
	// a configuration made in Go, or by ParseConfig, has its secrets read
	// here (relative paths are the working directory's)
	if err := cfg.resolve(); err != nil {
		return nil, err
	}
	if o.Log == nil {
		o.Log = slog.New(slog.DiscardHandler)
	}
	if o.Grace <= 0 {
		o.Grace = 30 * time.Second
	}
	so, err := cfg.spoolOptions()
	if err != nil {
		return nil, err
	}
	so.Converter, so.FontFS, so.Log, so.Now = o.Converter, o.FontFS, o.Log, o.Now
	sp, err := spool.Open(ctx, so)
	if err != nil {
		return nil, err
	}
	// the fingerprints of the printers' certificates are kept in the
	// spool, the configuration's pinned
	ts := sp.TrustStore()
	ts.Pinned = cfg.pinned()
	for _, p := range so.Printers {
		if p.Client != nil && !p.Client.TLS.Insecure {
			p.Client.TLS.Store = ts
		}
	}
	s := &Server{cfg: cfg, o: o, sp: sp, log: o.Log, auth: newAuthenticator(cfg),
		uploads: make(chan struct{}, orDefault(cfg.Limits.Uploads, 8)), streams: make(chan struct{}, orDefault(cfg.Limits.Events, 64)),
		done: make(chan struct{})}
	s.handler = s.routes()
	if cfg.IPP != nil {
		if s.ipp, err = newIPPServer(s); err != nil {
			sp.Close()
			return nil, err
		}
	}
	return s, nil
}

func orDefault(n, def int) int {
	if n <= 0 {
		return def
	}
	return n
}

// Spool is the server's spool.
func (s *Server) Spool() *spool.Spool { return s.sp }

// Handler is the server's HTTP handler: the API, /events, /metrics,
// /healthz and the page of the state.
func (s *Server) Handler() http.Handler { return s.handler }

// IPPHandler is the handler of the IPP server (Config.IPP), nil without
// one: the printers at /ipp/print/<queue>, their icons and pages.
func (s *Server) IPPHandler() http.Handler {
	if s.ipp == nil {
		return nil
	}
	return s.ipp.handler()
}

// Run listens on the configuration's address and serves until ctx ends.
func (s *Server) Run(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.cfg.Listen)
	if err != nil {
		s.Close()
		return err
	}
	return s.Serve(ctx, ln)
}

// Serve starts the spool and serves a listener until ctx ends. Then it
// stops taking requests, ends the event streams, lets the documents being
// sent finish for Options.Grace, and closes the spool.
func (s *Server) Serve(ctx context.Context, ln net.Listener) error {
	s.sp.Start()
	hs := &http.Server{
		Handler: s.handler,
		// a request has a minute to arrive; a submission extends it for
		// its document
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       readTimeout,
		IdleTimeout:       2 * time.Minute,
		ErrorLog:          slog.NewLogLogger(s.log.Handler(), slog.LevelWarn),
	}
	errc := make(chan error, 2)
	go func() {
		if t := s.cfg.TLS; t != nil {
			errc <- hs.ServeTLS(ln, t.Cert, t.Key)
		} else {
			errc <- hs.Serve(ln)
		}
	}()
	s.log.Info("serving", "listen", ln.Addr().String(), "tls", s.cfg.TLS != nil, "spool", s.cfg.Spool)
	var ippHS *http.Server
	if s.ipp != nil {
		var err error
		if ippHS, err = s.ipp.serve(ctx, errc); err != nil {
			hs.Close()
			s.stop()
			s.sp.Shutdown(s.o.Grace)
			return err
		}
	}
	var err error
	select {
	case <-ctx.Done():
	case err = <-errc:
	}
	s.log.Info("stopping")
	s.stop()
	sctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	hs.Shutdown(sctx)
	if ippHS != nil {
		ippHS.Shutdown(sctx)
		s.ipp.stop()
	}
	cancel()
	if cerr := s.sp.Shutdown(s.o.Grace); err == nil || errors.Is(err, http.ErrServerClosed) {
		err = cerr
	}
	s.log.Info("stopped")
	return err
}

// Close closes the spool of a server that is not served.
func (s *Server) Close() error {
	s.stop()
	return s.sp.Close()
}

func (s *Server) stop() {
	s.stopOnce.Do(func() { close(s.done) })
}

// redact drops the user and password a URL may hold.
func redact(u string) string {
	p, err := url.Parse(u)
	if err != nil || p.User == nil {
		return u
	}
	p.User = nil
	return p.String()
}
