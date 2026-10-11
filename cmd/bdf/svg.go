package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shibukawa/bdf/vector/svgbdf"
)

// svgCmd writes the pages of a document as SVG files.
func svgCmd(args []string) {
	fs := flag.NewFlagSet("svg", flag.ExitOnError)
	inline := fs.Bool("inline", false, "put the fonts and images into each SVG, for files that stand alone (an SVG shown by an img element loads no other file)")
	allow := fs.Bool("allow-plaintext", false, "write an encrypted document (read with $"+passwordEnv+"); the SVGs are not encrypted")
	quiet := fs.Bool("q", false, "do not print warnings")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: bdf svg [flags] <file.bdf | dir> <outdir>")
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() != 2 {
		fs.Usage()
		os.Exit(2)
	}
	doc := loadDocument(fs.Arg(0), *allow)
	res, err := svgbdf.Write(fs.Arg(1), doc, &svgbdf.Options{Inline: *inline})
	check(err)
	if !*quiet {
		for _, w := range res.Warnings {
			fmt.Fprintln(os.Stderr, "warning:", w)
		}
	}
	if len(res.Pages) == 0 {
		check(fmt.Errorf("%s has no page to write", fs.Arg(0)))
	}
}
