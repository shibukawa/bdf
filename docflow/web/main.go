package web

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"log/slog"
	"os"

	"github.com/shibukawa/popcornweb/plugin/auth"
	"github.com/shibukawa/popcornweb/pw"
	"github.com/shibukawa/popcornweb/pwconfig"
	"github.com/shibukawa/popcornweb/sessionconfig"
)

// Options compose the web application around Docflow's own pages and
// API, for a program of its own (cmd/docflow sets none): what it adds
// to the one server.
type Options struct {
	// Public is the program's own public tree, which Popcorn Web serves
	// at /public (its [assets] and manifest); nil for none. Docflow's
	// own scripts and style are not there: they are served at AssetsPath
	// by this package.
	Public fs.FS
	// Routes adds the program's own pages and API to the mux, after
	// Docflow's are there: a page of its own, an API under a prefix of
	// its own (not /api/<kind>/…, which are the nodes').
	Routes func(mux *pw.ServeMux, app *App)
	// AssetsDir serves Docflow's scripts and style from a directory (this
	// package's public, while they are edited) instead of the copy
	// embedded in the binary.
	AssetsDir string
}

// Main runs the web application as a program's main: RegisterConfig
// (and the program's own registrations) first. It parses the
// configuration, answers a framework action (--generate-config, the
// health probe …), sets the sign-in up, opens the store, starts the
// timer and the ports, and serves until the process is told to stop; an
// error ends the program.
func Main(o Options) {
	if err := Run(context.Background(), o); err != nil {
		log.Fatal(err)
	}
}

// Run is Main that returns its error (and nothing when a framework
// action was answered).
func Run(ctx context.Context, o Options) error {
	if err := pw.ParseConfig(); err != nil {
		return err
	}
	// a subcommand (RegisterCommands): gen-project writes a project
	if c, ok := pw.Command[GenProject](); ok {
		return GenerateProject(c, os.Stdout)
	}
	// --generate-config, the health probe …: answered before the store
	// is opened
	if handled, err := pw.RunFrameworkAction(); handled {
		return err
	}
	// the pages' sign-in (popcornweb's [auth] and [session]), before the
	// server reads them
	authCfg, _ := pwconfig.Bound[auth.Config]()
	sessionCfg, _ := pwconfig.Bound[sessionconfig.SessionConfig]()
	mc := pw.ConfigContext[Config](nil)
	oc := pw.ConfigContext[OIDCConfig](nil)
	if authCfg == nil || sessionCfg == nil {
		return errors.New("docflow: popcornweb's auth and session settings are not registered")
	}
	if err := ConfigureAuth(authCfg, sessionCfg, mc.Secret, mc.PreviousSecrets, mc.SecureCookies, oc); err != nil {
		return err
	}
	// the startup summary shows the settings as they were read
	if authCfg.Enabled {
		slog.Info("docflow: people sign in with OIDC", "issuer", authCfg.OIDC.Issuer, "session", sessionCfg.Backend,
			"admins", len(oc.Admins), "admin_groups", len(oc.AdminGroups))
	} else {
		slog.Warn("docflow: no OIDC issuer (AUTH_OIDC_ISSUER); nobody signs in to the pages (dev only)")
	}
	app, err := Open(ctx, mc, oc, o)
	if err != nil {
		return err
	}
	defer app.Close()
	stop := app.StartTimer(ctx)
	defer stop()
	// docflow.https_port serves the whole chain, for the devices on the
	// network
	pw.RegisterMiddleware(pw.SlotMetrics-4, "https-port", app.Outermost)
	stopHTTPS, err := app.StartHTTPS()
	if err != nil {
		return err
	}
	defer stopHTTPS()
	pw.RegisterMiddleware(pw.SlotAccessLog-5, "quiet-polls", QuietPolls)
	public := o.Public
	if public == nil {
		public = emptyFS{}
	}
	return pw.Run(ctx, app.Handler(), pw.WithPublicFS(public))
}

// emptyFS is the public tree of a program that has none.
type emptyFS struct{}

func (emptyFS) Open(name string) (fs.File, error) {
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}
