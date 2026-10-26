// Command docflow is BDF Docflow, a document workflow engine on Popcorn
// Web: documents and data enter at uploaders and virtual printers, flow
// through a topology of nodes, and end on printers, on the tablets of a
// location (a warehouse), in archives and in downloaders. The engine is
// package docflow, the web application docflow/web: this is the one
// that serves every input and output of docflow/io. A program of its own
// composes the same application with kinds of node, settings and pages
// of its own (docflow/web.Options; examples/docflow-custom).
//
// People sign in to its pages at an OIDC IdP (popcornweb's plugin/auth,
// set by web.ConfigureAuth), and see the jobs of their scope; the
// devices of /tablet keep their own sign-in.
//
// It runs as a process (a container, a server of a site) or as a
// function (AWS Lambda through the Lambda Web Adapter): with
// docflow.serverless, the records must be in a store several instances
// share, and a scheduler calls POST /api/tick instead of the timer.
package main

import (
	_ "github.com/shibukawa/bdf/docflow/io/all"
	"github.com/shibukawa/bdf/docflow/web"
)

func main() {
	web.RegisterConfig()
	web.Main(web.Options{Public: publicTree()})
}
