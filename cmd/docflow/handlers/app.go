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

	"github.com/shibukawa/bdf/cmd/docflow/help"
	"github.com/shibukawa/bdf/docflow"
	// the inputs and outputs: printers, CUPS servers and /ipp/print;
	// mail, S3 and its trigger, HTTP; the notify and render nodes
	_ "github.com/shibukawa/bdf/docflow/io/http"
	_ "github.com/shibukawa/bdf/docflow/io/ipp"
	_ "github.com/shibukawa/bdf/docflow/io/notify"
	_ "github.com/shibukawa/bdf/docflow/io/s3"
	_ "github.com/shibukawa/bdf/docflow/io/smtp"
	_ "github.com/shibukawa/bdf/docflow/render"
	"github.com/shibukawa/bdf/print/printconv"
	"github.com/shibukawa/popcornweb/plugin/auth"
	"github.com/shibukawa/popcornweb/pw"
)

// App is the manager of the configured store, and the pages around it.
type App struct {
	cfg    DocflowConfig
	policy policy
	st     *docflow.Store
	m      *docflow.Manager
	ports  *ports
	log    *slog.Logger
}

// Open opens the store and makes the manager. oidc is who of the people
// the IdP signs in gets in, and who is an admin (ConfigureAuth sets the
// sign-in itself).
func Open(ctx context.Context, cfg DocflowConfig, oidc OIDCConfig) (*App, error) {
	log := slog.Default()
	secret := []byte(cfg.Secret)
	if len(secret) == 0 {
		if !pw.Development() {
			return nil, errors.New("docflow: docflow.secret (DOCFLOW_SECRET) is required: 32 bytes or more, the same for every instance")
		}
		secret = make([]byte, 32)
		rand.Read(secret)
		log.Warn("docflow: no docflow.secret; a random one is used, and the devices sign in again after a restart")
	}
	tokens := map[string]string{}
	for _, t := range cfg.SubmitTokens {
		name, secret, ok := strings.Cut(t, ":")
		if !ok || name == "" || secret == "" {
			return nil, fmt.Errorf("docflow: submit token %q: want name:secret", name)
		}
		tokens[name] = secret
	}
	if cfg.AdminToken == "" {
		log.Info("docflow: no docflow.admin_token; the API (/api) answers the administrators signed in to the pages only")
	}
	if oidc.GroupsClaim == "" {
		oidc.GroupsClaim = "groups"
	}
	st, err := docflow.OpenStore(ctx, cfg.Records, cfg.Documents)
	if err != nil {
		return nil, err
	}
	if cfg.Serverless && !st.Shared {
		st.Close()
		return nil, fmt.Errorf("docflow: %s is one process's store; serverless needs one that instances share (dynamodb, firestore, s3, r2, gcs)", cfg.Records)
	}
	var previous [][]byte
	for _, p := range cfg.PreviousSecrets {
		if len(p) < 32 {
			return nil, errors.New("docflow: a previous secret has fewer than 32 bytes")
		}
		previous = append(previous, []byte(p))
	}
	a := &App{cfg: cfg, policy: policy{oidc}, st: st, log: log}
	a.ports = newPorts(a)
	conv := printconv.Converter{FontDirs: cfg.FontDirs, NoSystemFonts: cfg.NoSystemFonts}
	a.m, err = docflow.New(st, docflow.Options{
		Converter: conv, Secret: secret, PreviousSecrets: previous, FontDirs: cfg.FontDirs, NoSystemFonts: cfg.NoSystemFonts,
		AdminToken: cfg.AdminToken, SubmitTokens: tokens, TickToken: cfg.TickToken, Viewer: a.viewer,
		MaxDocument: int64(cfg.MaxDocumentMB) << 20, SecureCookies: cfg.SecureCookies, TabletPath: "/tablet", Log: log,
		JobKeep: cfg.JobKeep, OnError: cfg.OnError, CredentialDir: cfg.CredentialDir,
		// a process sends deliveries as jobs come, and listens on the
		// ports of the entry nodes; a function leaves the deliveries to
		// the scheduler's tick and has no ports
		Background: !cfg.Serverless, NoPorts: cfg.Serverless, OnChange: a.ports.kick,
	})
	if err != nil {
		st.Close()
		return nil, err
	}
	if len(tokens) == 0 {
		log.Warn("docflow: no docflow.submit_tokens; anyone who reaches the printers may send them jobs")
	}
	// plugin/auth's state of logins, and its accounts, are this app's
	authStore.Store(st)
	auth.SetAccountResolver(a.resolveAccount)
	return a, nil
}

// Close closes the ports and the store.
func (a *App) Close() error {
	a.ports.closeAll()
	authStore.CompareAndSwap(a.st, nil)
	return a.st.Close()
}

// StartTimer ticks every TickInterval while the process runs, and keeps
// the ports of the entry nodes open; serverless, or with no interval, a
// scheduler calls POST /api/tick instead (and there are no ports).
func (a *App) StartTimer(ctx context.Context) (stop func()) {
	ctx, cancel := context.WithCancel(ctx)
	if !a.cfg.Serverless {
		go a.ports.run(ctx, a.cfg.TickInterval)
	}
	if a.cfg.Serverless || a.cfg.TickInterval <= 0 {
		return cancel
	}
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
// (ConfigureAuth); /admin/… and /audit for an administrator; /help is the
// manual.
func (a *App) Handler() http.Handler {
	mux := pw.NewServeMux()
	mh := docflow.Handler(a.m)
	for _, p := range docflow.Prefixes() {
		mux.Handle(p+"/", mh)
		mux.Handle(p, mh)
	}
	mux.Handle(docflow.InPath+"/", mh)
	mux.Handle(docflow.OutPath+"/", mh)
	mux.Handle("/oauth/", mh)
	mux.Handle("/api/", mh)
	mux.Handle("/device/", mh)
	// the pages: a sign-in at the top; the jobs of the person signed in;
	// the administrator's pages and audit log
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		if a.viewer(r) != nil {
			http.Redirect(w, r, "/jobs", http.StatusSeeOther)
			return
		}
		pw.WriteHTML(w, r, SignIn(SignInParams{Enabled: auth.Protection().Protected != nil}))
	})
	mux.HandleFunc("GET /jobs", a.page("", func(w http.ResponseWriter, r *http.Request, ac Account) {
		pw.WriteHTML(w, r, Jobs(JobsParams{Account: ac}))
	}))
	// the topology is everyone's page: a person with the role topology
	// edits it there, the others look (the page renders the editor's
	// controls by the role; the API checks the role itself)
	mux.HandleFunc("GET /topology", a.page("", func(w http.ResponseWriter, r *http.Request, ac Account) {
		pw.WriteHTML(w, r, Topology(TopologyParams{Account: ac}))
	}))
	// the pages of the role topology, one item each; /admin and
	// /admin/topology (as it was) are the topology
	mux.HandleFunc("GET /admin", a.page(docflow.RoleTopology, func(w http.ResponseWriter, r *http.Request, ac Account) {
		http.Redirect(w, r, "/topology", http.StatusSeeOther)
	}))
	mux.HandleFunc("GET /admin/topology", a.page(docflow.RoleTopology, func(w http.ResponseWriter, r *http.Request, ac Account) {
		http.Redirect(w, r, "/topology", http.StatusSeeOther)
	}))
	mux.HandleFunc("GET /admin/apps", a.page(docflow.RoleTopology, func(w http.ResponseWriter, r *http.Request, ac Account) {
		pw.WriteHTML(w, r, Apps(AppsParams{Account: ac}))
	}))
	mux.HandleFunc("GET /admin/inboxes", a.page(docflow.RoleTopology, func(w http.ResponseWriter, r *http.Request, ac Account) {
		pw.WriteHTML(w, r, Inboxes(InboxesParams{Account: ac}))
	}))
	mux.HandleFunc("GET /admin/jobs", a.page(docflow.RoleTopology, func(w http.ResponseWriter, r *http.Request, ac Account) {
		pw.WriteHTML(w, r, AdminJobs(AdminJobsParams{Account: ac}))
	}))
	mux.HandleFunc("GET /admin/deliveries", a.page(docflow.RoleTopology, func(w http.ResponseWriter, r *http.Request, ac Account) {
		pw.WriteHTML(w, r, Deliveries(DeliveriesParams{Account: ac}))
	}))
	// the page of the role audit
	mux.HandleFunc("GET /audit", a.page(docflow.RoleAudit, func(w http.ResponseWriter, r *http.Request, ac Account) {
		pw.WriteHTML(w, r, Audit(AuditParams{Account: ac}))
	}))
	mux.HandleFunc("GET /tablet", func(w http.ResponseWriter, r *http.Request) {
		pw.WriteHTML(w, r, Tablet(TabletParams{}))
	})
	// the manual, for everyone signed in (help.go)
	mux.HandleFunc("GET /help", a.page("", a.help))
	mux.HandleFunc("GET /help/{first}", a.page("", a.help))
	mux.HandleFunc("GET /help/{first}/{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, help.Path+"/"+r.PathValue("first"), http.StatusSeeOther)
	})
	mux.HandleFunc("GET /help/{first}/{second}", a.page("", a.help))
	return mux
}

// page serves a page of a person signed in; with a role, of a person
// with it (others get 403).
func (a *App) page(role string, fn func(http.ResponseWriter, *http.Request, Account)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		v := a.viewer(r)
		switch {
		case v == nil:
			http.Redirect(w, r, "/", http.StatusSeeOther)
		case role != "" && !v.Can(role):
			pw.WriteHTML(&statusWriter{ResponseWriter: w, status: http.StatusForbidden}, r, Forbidden(ForbiddenParams{Account: account(v), Role: role}))
		default:
			fn(w, r, account(v))
		}
	}
}

func account(v *docflow.Viewer) Account {
	name := firstOr(v.Name, v.Email)
	initial := "?"
	for _, r := range name {
		initial = strings.ToUpper(string(r))
		break
	}
	return Account{Name: name, Initial: initial, Email: v.Email, Admin: v.Admin, Topology: v.Can(docflow.RoleTopology), Audit: v.Can(docflow.RoleAudit)}
}

func firstOr(s, or string) string {
	if s == "" {
		return or
	}
	return s
}
