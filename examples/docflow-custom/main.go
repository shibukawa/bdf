// Command docflow-custom is BDF Docflow composed by a program of its own:
// the web application of docflow/web, with every input and output of
// docflow/io, and a kind of node of this program's own, the folder exit
// (kinds/folder), which the editor lays out from the fields it
// registers. It is what a program that imports bdf and adds kinds,
// settings and pages of its own looks like: a Popcorn Web project
// (popcornweb.toml beside go.mod), built with pw build.
package main

import (
	_ "github.com/shibukawa/bdf/docflow/io/all"
	"github.com/shibukawa/bdf/docflow/web"

	// the kinds of this program's own register themselves
	_ "github.com/shibukawa/bdf/examples/docflow-custom/kinds/folder"
)

func main() {
	// Docflow's settings ([docflow], [oidc]) and the program's own, before
	// the configuration is parsed (web.Main parses it)
	web.RegisterConfig()
	web.Main(web.Options{Public: publicTree()})
}
