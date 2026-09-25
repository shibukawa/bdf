// Command pdf2bdf converts a PDF file into a BDF document.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/shibukawa/bdf/pdf2bdf"
)

func main() {
	title := flag.String("title", "", "document title (default: PDF Info Title)")
	pages := flag.String("pages", "", "pages to convert, e.g. 1-3,5 (default: all)")
	kind := flag.String("kind", "fixed", "view kind: fixed or flow")
	quiet := flag.Bool("q", false, "do not print warnings")
	noSubset := flag.Bool("no-subset", false, "keep unused glyphs of embedded TrueType fonts")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: pdf2bdf [flags] in.pdf out.bdf|outdir/\n  an output path ending with / writes the split form")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(2)
	}
	in, out := flag.Arg(0), flag.Arg(1)
	opts := &pdf2bdf.Options{Title: *title, Kind: *kind, NoSubset: *noSubset}
	if *pages != "" {
		sel, err := pdf2bdf.PageRange(*pages, 1<<30)
		if err != nil {
			fmt.Fprintln(os.Stderr, "pdf2bdf:", err)
			os.Exit(2)
		}
		opts.Pages = sel
	}
	res, err := pdf2bdf.ConvertFile(in, opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "pdf2bdf:", err)
		os.Exit(1)
	}
	if !*quiet {
		for _, w := range res.Warnings {
			fmt.Fprintln(os.Stderr, "warning:", w)
		}
	}
	if strings.HasSuffix(out, "/") {
		err = res.Doc.WriteSplit(out)
	} else {
		var f *os.File
		if f, err = os.Create(out); err == nil {
			if err = res.Doc.WriteSingle(f); err == nil {
				err = f.Close()
			}
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "pdf2bdf:", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "%s: %d page(s), %d warning(s)\n", out, res.Pages, len(res.Warnings))
}
