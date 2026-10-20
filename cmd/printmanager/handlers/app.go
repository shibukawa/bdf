package handlers

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/shibukawa/bdf/print/manager"
	"github.com/shibukawa/bdf/print/printconv"
	"github.com/shibukawa/popcornweb/pw"
)

// App is the manager of the configured store, and the pages around it.
type App struct {
	cfg ManagerConfig
	st  *manager.Store
	m   *manager.Manager
	log *slog.Logger
}

// Open opens the store and makes the manager.
func Open(ctx context.Context, cfg ManagerConfig) (*App, error) {
	log := slog.Default()
	secret := []byte(cfg.Secret)
	if len(secret) == 0 {
		if !pw.Development() {
			return nil, errors.New("printmanager: manager.secret (MANAGER_SECRET) is required: 32 bytes or more, the same for every instance")
		}
		secret = make([]byte, 32)
		rand.Read(secret)
		log.Warn("printmanager: no manager.secret; a random one is used, and the devices sign in again after a restart")
	}
	tokens := map[string]string{}
	for _, t := range cfg.SubmitTokens {
		name, secret, ok := strings.Cut(t, ":")
		if !ok || name == "" || secret == "" {
			return nil, fmt.Errorf("printmanager: submit token %q: want name:secret", name)
		}
		tokens[name] = secret
	}
	if cfg.AdminToken == "" {
		log.Warn("printmanager: no manager.admin_token; the API (/api) answers nobody")
	}
	st, err := manager.OpenStore(ctx, cfg.Records, cfg.Documents)
	if err != nil {
		return nil, err
	}
	if cfg.Serverless && !st.Shared {
		st.Close()
		return nil, fmt.Errorf("printmanager: %s is one process's store; serverless needs one that instances share (dynamodb, firestore, s3, r2, gcs)", cfg.Records)
	}
	conv := printconv.Converter{FontDirs: cfg.FontDirs, NoSystemFonts: cfg.NoSystemFonts}
	m, err := manager.New(st, manager.Options{
		Converter: conv, Secret: secret, FontDirs: cfg.FontDirs, NoSystemFonts: cfg.NoSystemFonts,
		AdminToken: cfg.AdminToken, SubmitTokens: tokens, TickToken: cfg.TickToken,
		MaxDocument: int64(cfg.MaxDocumentMB) << 20, SecureCookies: cfg.SecureCookies, TabletPath: "/tablet", Log: log,
	})
	if err != nil {
		st.Close()
		return nil, err
	}
	if len(tokens) == 0 {
		log.Warn("printmanager: no manager.submit_tokens; anyone who reaches the printers may send them jobs")
	}
	return &App{cfg: cfg, st: st, m: m, log: log}, nil
}

// Close closes the store.
func (a *App) Close() error { return a.st.Close() }

// StartTimer ticks every TickInterval while the process runs; serverless,
// or with no interval, a scheduler calls POST /api/tick instead.
func (a *App) StartTimer(ctx context.Context) (stop func()) {
	if a.cfg.Serverless || a.cfg.TickInterval <= 0 {
		return func() {}
	}
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		t := time.NewTicker(a.cfg.TickInterval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if _, err := a.m.Tick(ctx); err != nil && ctx.Err() == nil {
					a.log.Warn("tick", "err", err)
				}
			}
		}
	}()
	return cancel
}

// Handler is the application: the pages, their scripts and style (static,
// served at StaticPath), and the manager's printers, API and devices'
// endpoints.
func (a *App) Handler(static fs.FS) (http.Handler, error) {
	sh, err := staticHandler(static)
	if err != nil {
		return nil, err
	}
	mux := pw.NewServeMux()
	mux.Handle("GET "+StaticPath, sh)
	mh := manager.Handler(a.m)
	mux.Handle(manager.IPPPath+"/", mh)
	mux.Handle("/api/", mh)
	mux.Handle("/device/", mh)
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		pw.WriteHTML(w, r, Home(HomeParams{}))
	})
	mux.HandleFunc("GET /tablet", func(w http.ResponseWriter, r *http.Request) {
		pw.WriteHTML(w, r, Tablet(TabletParams{}))
	})
	mux.HandleFunc("GET /admin", func(w http.ResponseWriter, r *http.Request) {
		pw.WriteHTML(w, r, Admin(AdminParams{}))
	})
	return mux, nil
}
