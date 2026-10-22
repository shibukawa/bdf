package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shibukawa/bdf/text/markdownbdf"
)

// markdownCmd writes a document as Markdown, with its images, in a directory.
func markdownCmd(args []string) {
	fs := flag.NewFlagSet("markdown", flag.ExitOnError)
	var views, skip stringList
	fs.Var(&views, "view", "id of a view to write (repeatable; default: the scroll view of a document that has one, else all views)")
	fs.Var(&skip, "skip", "role of the layers to leave out (repeatable; default: header, footer, master and background; none keeps them all)")
	internal := fs.Bool("internal-links", false, "keep links into the document itself (#page=N), which have no target in Markdown")
	running := fs.Bool("keep-running", false, "keep the running heads and feet of pages: the text repeated at the top or bottom of pages, and page numbers")
	tables := fs.String("tables", "inline", "the tables the document holds (bdf tables; the data of charts): inline (after the figure, up to 100 rows, the whole as a CSV file beside the text when there are more), files (CSV files only, the tables of sheets too) or none")
	allow := fs.Bool("allow-plaintext", false, "write an encrypted document (read with $"+passwordEnv+"); the Markdown is plain text")
	quiet := fs.Bool("q", false, "do not print warnings")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: bdf markdown [flags] <file.bdf | dir> <outdir>\n  writes the document as Markdown: <outdir>/<name>.md, named after the input, with its images beside it")
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() != 2 {
		fs.Usage()
		os.Exit(2)
	}
	in, out := fs.Arg(0), fs.Arg(1)
	doc := loadDocument(in, *allow)
	opts := &markdownbdf.Options{Views: views, InternalLinks: *internal, KeepRunning: *running, Tables: *tables}
	if len(skip) > 0 {
		opts.SkipRoles = []string{}
		for _, r := range skip {
			if r != "none" {
				opts.SkipRoles = append(opts.SkipRoles, r)
			}
		}
	}
	name := filepath.Base(filepath.Clean(in))
	name = strings.TrimSuffix(name, filepath.Ext(name))
	if name == "" || name == "." || name == string(filepath.Separator) {
		name = "document"
	}
	res, err := markdownbdf.WriteDir(out, name, doc, opts)
	check(err)
	if !*quiet {
		for _, w := range res.Warnings {
			fmt.Fprintln(os.Stderr, "warning:", w)
		}
	}
	fmt.Fprintf(os.Stderr, "%s: %d view(s), %d image(s), %d table(s), %d warning(s)\n", filepath.Join(out, name+".md"), res.Views, res.Images, res.Tables, len(res.Warnings))
}
