// Command gen-webp translates tools/webp/webp.wasm into the Go package
// imgconv/internal/webpw with the wasm2go fork, using the options that keep
// regenerated output diff-friendly (symbol names, per-subject files, named
// data addresses). Run it through tools/gen-codecs.sh.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/wasm2go/transpile"
)

func main() {
	in := flag.String("i", "", "input wasm")
	outDir := flag.String("out-dir", "", "output directory")
	pkg := flag.String("pkg", "webpw", "package name")
	importPath := flag.String("import", "github.com/shibukawa/bdf/imgconv/internal/webpw", "import path")
	chunks := flag.Int("chunks", 1, "number of chunk packages")
	flag.Parse()
	if *in == "" || *outDir == "" {
		fmt.Fprintln(os.Stderr, "usage: gen-webp -i webp.wasm -out-dir DIR")
		os.Exit(2)
	}
	f, err := os.Open(*in)
	if err != nil {
		fail(err)
	}
	mod, err := transpile.Parse(f)
	f.Close()
	if err != nil {
		fail(err)
	}
	// Always emit the multi-package layout so -group-files applies even to a
	// module below the CLI's single-file threshold.
	restore := transpile.SetMultiPackageThreshold(0)
	defer restore()
	res, err := transpile.Translate(os.Stderr, mod, transpile.Options{
		Package:          *pkg,
		OutputImportPath: *importPath,
		PureOnly:         true,
		SymbolNames:      true,
		Chunks:           *chunks,
		GroupFiles:       true,
		AddrConsts:       true,
		EntryExports:     []string{"_initialize", "encode", "malloc", "free"},
	})
	if err != nil {
		fail(err)
	}
	if err := os.RemoveAll(*outDir); err != nil {
		fail(err)
	}
	write := func(rel string, data []byte) {
		p := filepath.Join(*outDir, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			fail(err)
		}
		if err := os.WriteFile(p, data, 0o644); err != nil {
			fail(err)
		}
	}
	for rel, data := range res.Files {
		write(rel, data)
	}
	for rel, data := range res.AuxFiles {
		write(rel, data)
	}
	for name, data := range res.Sidecars {
		write(name, data)
	}
	fmt.Fprintf(os.Stderr, "gen-webp: %d files\n", len(res.Files)+len(res.AuxFiles)+len(res.Sidecars))
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "gen-webp:", err)
	os.Exit(1)
}
