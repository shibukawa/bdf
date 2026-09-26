package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
	_ "github.com/shibukawa/bdf/converter/all" // every input format
	"github.com/shibukawa/bdf/imgconv"
)

// stringList is a repeatable string flag.
type stringList []string

func (l *stringList) String() string     { return strings.Join(*l, ",") }
func (l *stringList) Set(v string) error { *l = append(*l, v); return nil }

// generate converts a file in one of the registered input formats into a
// BDF document.
func generate(args []string) {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	var names []string
	for _, f := range converter.Formats() {
		names = append(names, f.Name)
	}
	format := fs.String("format", "auto", "input format: auto (detected from the content) or one of "+strings.Join(names, ", "))
	title := fs.String("title", "", "document title (default: from the input); the same as -dc title=...")
	var dcFlags stringList
	fs.Var(&dcFlags, "dc", "Dublin Core element as name=value, e.g. creator=Alice (repeatable; replaces the element read from the input, name= removes it)")
	pages := fs.String("pages", "", "pages, slides or sheets to convert, e.g. 1-3,5 (default: all)")
	quiet := fs.Bool("q", false, "do not print warnings")
	images := fs.String("images", "convert", "raster images: keep (store as is) or convert (try WebP, keep when smaller)")
	quality := fs.Int("quality", 80, "lossy WebP quality (1-100)")
	noSubset := fs.Bool("no-subset", false, "embed whole fonts instead of the glyphs in use")
	noWOFF2 := fs.Bool("no-woff2", false, "store embedded fonts as TrueType/OpenType instead of WOFF2")
	ignoreFSType := fs.Bool("ignore-fstype", false, "embed fonts whose OS/2 fsType forbids embedding or subsetting (only with the rights to do so)")
	kind := fs.String("kind", "fixed", "PDF: view kind, fixed or flow")
	noShare := fs.Bool("no-share", false, "PDF: do not move the instruction prefix pages have in common into a shared object")
	fonts := fs.String("fonts", "embed", "PowerPoint, Excel, metafiles: embed (subset and embed the fonts used for layout) or system (refer to fonts by name)")
	var fontDirs stringList
	fs.Var(&fontDirs, "font-dir", "PowerPoint, Excel, metafiles: directory searched for fonts before the system ones (repeatable)")
	noSystemFonts := fs.Bool("no-system-fonts", false, "PowerPoint, Excel, metafiles: use only the fonts under -font-dir")
	hidden := fs.Bool("hidden", false, "PowerPoint, Excel: include hidden slides or sheets (the same as -param hidden=true)")
	var paramFlags stringList
	fs.Var(&paramFlags, "param", "format-specific option as name=value (repeatable; see the formats below)")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: bdf generate [flags] <input> <out.bdf | outdir/>\n  an output path ending with / writes the split form")
		fs.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\ninput formats:")
		for _, f := range converter.Formats() {
			fmt.Fprintf(os.Stderr, "  %-6s %s (%s)\n", f.Name, f.Description, strings.Join(f.Extensions, " "))
			for _, p := range f.Params {
				fmt.Fprintf(os.Stderr, "         -param %s=…: %s\n", p.Name, p.Usage)
			}
		}
	}
	fs.Parse(args)
	if fs.NArg() != 2 {
		fs.Usage()
		os.Exit(2)
	}
	in, out := fs.Arg(0), fs.Arg(1)

	imgOpts := imgconv.Options{Quality: *quality}
	switch *images {
	case "keep":
		imgOpts.Mode = imgconv.Keep
	case "convert":
		imgOpts.Mode = imgconv.Convert
	default:
		usageError("-images must be keep or convert")
	}
	var sel []int
	if *pages != "" {
		var err error
		if sel, err = converter.PageRange(*pages, 1<<30); err != nil {
			usageError(err.Error())
		}
	}
	dc, err := parseDC(dcFlags)
	if err != nil {
		usageError(err.Error())
	}
	if *title != "" {
		dc.Title = bdf.DCValues{*title}
	}
	opts := &converter.Options{Title: dc.Title.First(), Pages: sel, Images: imgOpts,
		FontDirs: fontDirs, NoSystemFonts: *noSystemFonts, NoSubset: *noSubset, NoWOFF2: *noWOFF2, IgnoreFSType: *ignoreFSType,
		Params: map[string]string{"kind": *kind, "no-share": strconv.FormatBool(*noShare), "hidden": strconv.FormatBool(*hidden)}}
	switch *fonts {
	case "embed":
	case "system":
		opts.SystemFonts = true
	default:
		usageError("-fonts must be embed or system")
	}
	for _, kv := range paramFlags {
		name, value, ok := strings.Cut(kv, "=")
		if !ok || name == "" {
			usageError("-param " + kv + ": want name=value")
		}
		opts.Params[name] = value
	}
	var f *converter.Format
	if *format == "auto" {
		f, err = converter.DetectFile(in)
		check(err)
		if f == nil {
			switch ext := strings.ToLower(filepath.Ext(in)); ext {
			case ".ppt":
				usageError(in + ": legacy .ppt files are not supported; save as .pptx first")
			case ".xls":
				usageError(in + ": legacy .xls files are not supported; save as .xlsx first")
			case ".vsd", ".vss", ".vst":
				usageError(in + ": legacy binary Visio files (" + ext + ") are not supported; save as .vsdx first")
			}
			usageError(in + ": unknown input format (want one of " + strings.Join(names, ", ") + ")")
		}
	} else if f = converter.Lookup(*format); f == nil {
		usageError("-format " + *format + ": unknown format (want auto or one of " + strings.Join(names, ", ") + ")")
	}
	res, err := converter.ConvertFile(in, f.Name, opts)
	check(err)
	doc, warnings := res.Doc, res.Warnings
	for _, name := range bdf.DCTerms {
		if v := *dc.Field(name); v != nil {
			*doc.Meta.DC.Field(name) = v
		}
	}
	if !*quiet {
		for _, w := range warnings {
			fmt.Fprintln(os.Stderr, "warning:", w)
		}
	}
	if strings.HasSuffix(out, "/") || strings.HasSuffix(out, string(filepath.Separator)) {
		check(doc.WriteSplit(out))
	} else {
		check(writeSingle(doc, out))
	}
	fmt.Fprintf(os.Stderr, "%s: %s, %d warning(s)\n", out, res.Summary, len(warnings))
}

// parseDC reads -dc name=value flags into the elements they replace. An
// element named with an empty value is set to an empty, non-nil list, which
// removes it.
func parseDC(flags []string) (bdf.DublinCore, error) {
	var dc bdf.DublinCore
	for _, kv := range flags {
		name, value, _ := strings.Cut(kv, "=")
		f := dc.Field(name)
		if f == nil {
			return dc, fmt.Errorf("-dc %s: unknown element %q (want one of %s)", kv, name, strings.Join(bdf.DCTerms, ", "))
		}
		if *f == nil {
			*f = bdf.DCValues{}
		}
		if value != "" {
			*f = append(*f, value)
		}
	}
	return dc, nil
}

func usageError(msg string) {
	fmt.Fprintln(os.Stderr, "bdf generate:", msg)
	os.Exit(2)
}
