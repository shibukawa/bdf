package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shibukawa/bdf/conformance"
)

// conformanceCmd writes the conformance suite: the inputs other
// implementations check themselves with, and what this one does with them.
func conformanceCmd(args []string) {
	fs := flag.NewFlagSet("conformance", flag.ExitOnError)
	testdata := fs.String("testdata", "testdata", "the directory of test documents to include as accept cases (\"\" for none)")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: bdf conformance [flags] <dir>")
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() != 1 {
		fs.Usage()
		os.Exit(2)
	}
	s, err := conformance.Build(*testdata)
	check(err)
	check(s.Write(fs.Arg(0)))
	fmt.Printf("%d cases written to %s\n", len(s.Cases), fs.Arg(0))
}
