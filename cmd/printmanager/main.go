// Command printmanager is the BDF Print Manager: virtual printers that
// show slips on the tablets of a location (a warehouse) instead of
// printing them, on Popcorn Web. The manager is print/manager; this is
// its web application: configuration, pages, and the timer of a process.
//
// It runs as a process (a container, a server of a site) or as a
// function (AWS Lambda through the Lambda Web Adapter): with
// manager.serverless, the records must be in a store several instances
// share, and a scheduler calls POST /api/tick instead of the timer.
package main

import (
	"context"
	"log"

	"github.com/shibukawa/bdf/cmd/printmanager/handlers"
	_ "github.com/shibukawa/bdf/converter/all"
	"github.com/shibukawa/popcornweb/pw"
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
	ctx := context.Background()
	app, err := handlers.Open(ctx, pw.ConfigContext[handlers.ManagerConfig](nil))
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
