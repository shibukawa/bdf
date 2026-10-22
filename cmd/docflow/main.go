// Command docflow is BDF Docflow, a document workflow engine on Popcorn
// Web: documents and data enter at endpoints and virtual printers, flow
// through a topology of nodes, and end on printers, on the tablets of a
// location (a warehouse), in archives and in exports. The engine is
// package docflow; this is its web application: configuration, pages, and the timer of a process.
//
// People sign in to its pages at an OIDC IdP (popcornweb's plugin/auth,
// set by handlers.ConfigureAuth), and see the jobs of their scope; the
// devices of /tablet keep their own sign-in.
//
// It runs as a process (a container, a server of a site) or as a
// function (AWS Lambda through the Lambda Web Adapter): with
// docflow.serverless, the records must be in a store several instances
// share, and a scheduler calls POST /api/tick instead of the timer.
package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/shibukawa/bdf/cmd/docflow/handlers"
	_ "github.com/shibukawa/bdf/converter/all"
	"github.com/shibukawa/popcornweb/plugin/auth"
	"github.com/shibukawa/popcornweb/pw"
	"github.com/shibukawa/popcornweb/pwconfig"
	"github.com/shibukawa/popcornweb/sessionconfig"
)

func main() {
	handlers.RegisterConfig()
	if err := pw.ParseConfig(); err != nil {
		log.Fatal(err)
	}
	// --generate-config, the health probe …: answered before the store
	// is opened
	if handled, err := pw.RunFrameworkAction(); handled {
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	// the pages' sign-in (popcornweb's [auth] and [session]), before the
	// server reads them
	authCfg, _ := pwconfig.Bound[auth.Config]()
	sessionCfg, _ := pwconfig.Bound[sessionconfig.SessionConfig]()
	mc := pw.ConfigContext[handlers.DocflowConfig](nil)
	oc := pw.ConfigContext[handlers.OIDCConfig](nil)
	if authCfg == nil || sessionCfg == nil {
		log.Fatal("docflow: popcornweb's auth and session settings are not registered")
	}
	if err := handlers.ConfigureAuth(authCfg, sessionCfg, mc.Secret, mc.PreviousSecrets, mc.SecureCookies, oc); err != nil {
		log.Fatal(err)
	}
	// the startup summary shows the settings as they were read
	if authCfg.Enabled {
		slog.Info("docflow: people sign in with OIDC", "issuer", authCfg.OIDC.Issuer, "session", sessionCfg.Backend,
			"admins", len(oc.Admins), "admin_groups", len(oc.AdminGroups))
	} else {
		slog.Warn("docflow: no OIDC issuer (AUTH_OIDC_ISSUER); nobody signs in to the pages (dev only)")
	}
	ctx := context.Background()
	app, err := handlers.Open(ctx, mc, oc)
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()
	stop := app.StartTimer(ctx)
	defer stop()
	if err := pw.Run(ctx, app.Handler(), pw.WithPublicFS(publicTree())); err != nil {
		log.Fatal(err)
	}
}
