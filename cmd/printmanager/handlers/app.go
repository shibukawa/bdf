package handlers

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/shibukawa/bdf/print/manager"
	"github.com/shibukawa/bdf/print/printconv"
	"github.com/shibukawa/popcornweb/plugin/auth"
	"github.com/shibukawa/popcornweb/pw"
)

// App is the manager of the configured store, and the pages around it.
type App struct {
	cfg    ManagerConfig
	policy policy
	st     *manager.Store
	m      *manager.Manager
	log    *slog.Logger
}

// Open opens the store and makes the manager. oidc is who of the people
// the IdP signs in gets in, and who is an admin (ConfigureAuth sets the
// sign-in itself).
func Open(ctx context.Context, cfg ManagerConfig, oidc OIDCConfig) (*App, error) {
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
		log.Info("printmanager: no manager.admin_token; the API (/api) answers the admins signed in to the pages only")
	}
	if oidc.GroupsClaim == "" {
		oidc.GroupsClaim = "groups"
	}
	st, err := manager.OpenStore(ctx, cfg.Records, cfg.Documents)
	if err != nil {
		return nil, err
	}
	if cfg.Serverless && !st.Shared {
		st.Close()
		return nil, fmt.Errorf("printmanager: %s is one process's store; serverless needs one that instances share (dynamodb, firestore, s3, r2, gcs)", cfg.Records)
	}
	a := &App{cfg: cfg, policy: policy{oidc}, st: st, log: log}
	conv := printconv.Converter{FontDirs: cfg.FontDirs, NoSystemFonts: cfg.NoSystemFonts}
	a.m, err = manager.New(st, manager.Options{
		Converter: conv, Secret: secret, FontDirs: cfg.FontDirs, NoSystemFonts: cfg.NoSystemFonts,
		AdminToken: cfg.AdminToken, SubmitTokens: tokens, TickToken: cfg.TickToken, Viewer: a.viewer,
		MaxDocument: int64(cfg.MaxDocumentMB) << 20, SecureCookies: cfg.SecureCookies, TabletPath: "/tablet", Log: log,
	})
	if err != nil {
		st.Close()
		return nil, err
	}
	if len(tokens) == 0 {
		log.Warn("printmanager: no manager.submit_tokens; anyone who reaches the printers may send them jobs")
	}
	// plugin/auth's state of logins, and its accounts, are this app's
	authStore.Store(st)
	auth.SetAccountResolver(a.resolveAccount)
	return a, nil
}

// Close closes the store.
func (a *App) Close() error {
	authStore.CompareAndSwap(a.st, nil)
	return a.st.Close()
}

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

// Handler is the application: the pages, and the manager's printers, API
// and devices' endpoints. The pages' scripts and style are the framework's
// /public mount's. Every page but /tablet is for a person signed in
// (ConfigureAuth); /admin for an admin.
func (a *App) Handler() http.Handler {
	mux := pw.NewServeMux()
	mh := manager.Handler(a.m)
	mux.Handle(manager.IPPPath+"/", mh)
	mux.Handle("/api/", mh)
	mux.Handle("/device/", mh)
	mux.HandleFunc("GET /{$}", a.home)
	mux.HandleFunc("GET /tablet", func(w http.ResponseWriter, r *http.Request) {
		pw.WriteHTML(w, r, Tablet(TabletParams{}))
	})
	mux.HandleFunc("GET /admin", func(w http.ResponseWriter, r *http.Request) {
		v := a.viewer(r)
		if v == nil || !v.Admin {
			w.Header().Set("Cache-Control", "no-store")
			pw.WriteHTML(&statusWriter{ResponseWriter: w, status: http.StatusForbidden}, r,
				Forbidden(ForbiddenParams{SignedIn: v != nil, Account: account(v)}))
			return
		}
		pw.WriteHTML(w, r, Admin(AdminParams{Account: account(v)}))
	})
	return mux
}

func account(v *manager.Viewer) Account {
	if v == nil {
		return Account{}
	}
	name := firstOr(v.Name, v.Email)
	initial := "?"
	for _, r := range name {
		initial = strings.ToUpper(string(r))
		break
	}
	return Account{Name: name, Initial: initial, Email: v.Email, Admin: v.Admin}
}

func firstOr(s, or string) string {
	if s == "" {
		return or
	}
	return s
}

var stateLabels = map[manager.SlipState]string{manager.SlipNew: "未開封", manager.SlipOpened: "開いた", manager.SlipDone: "完了", manager.SlipCanceled: "取り消し"}

var viaLabels = map[string]string{manager.ViaOIDC: "サインイン", manager.ViaSubmitToken: "送信のトークン", manager.ViaAdminToken: "管理のトークン", manager.ViaNone: "確かめていない"}

// home is everyone's page: the jobs they may see, and where they may
// print.
func (a *App) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	v := a.viewer(r)
	if v == nil {
		pw.WriteHTML(w, r, Home(HomeParams{NoSlips: true}))
		return
	}
	ctx := r.Context()
	slips, err := a.m.SlipsOf(ctx, v)
	if err != nil {
		a.log.Error("slips", "err", err)
		http.Error(w, "the jobs cannot be read", http.StatusInternalServerError)
		return
	}
	places, err := a.m.PlacesOf(ctx, v)
	if err != nil {
		a.log.Error("places", "err", err)
		http.Error(w, "the places cannot be read", http.StatusInternalServerError)
		return
	}
	p := HomeParams{SignedIn: true, Account: account(v), Groups: v.Groups, NoSlips: len(slips) == 0, CanPrint: len(places.Locations) > 0}
	for _, s := range slips {
		row := SlipRow{Job: int(s.Job), Name: firstOr(s.Name, "slip-"+s.ID), Location: s.Location, State: string(s.State),
			StateLabel: stateLabels[s.State], Pages: s.Pages, Received: s.Received.Local().Format("2006-01-02 15:04"),
			Claimed: s.ClaimedUser, Group: s.Group}
		if by := s.SubmittedBy; by != nil {
			row.By, row.Via, row.Mine = by.Label(), viaLabels[by.Via], by.Kind == manager.PrincipalUser && by.ID == v.ID
		}
		p.Slips = append(p.Slips, row)
	}
	for _, l := range places.Locations {
		p.Locations = append(p.Locations, Place{Name: l.Name, Title: l.Title, Kind: firstOr(l.Kind, manager.KindShared), Owner: l.Owner})
	}
	for _, pr := range places.Printers {
		p.Printers = append(p.Printers, Place{Name: pr.Name, Kind: firstOr(pr.Kind, manager.KindShared), Owner: pr.Owner})
	}
	pw.WriteHTML(w, r, Home(p))
}
