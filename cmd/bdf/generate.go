package main

import (
	"errors"
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
	fonts := fs.String("fonts", "embed", "PowerPoint, Excel, Word, CSV, metafiles: embed (subset and embed the fonts used for layout) or system (refer to fonts by name)")
	var fontDirs stringList
	fs.Var(&fontDirs, "font-dir", "PowerPoint, Excel, Word, CSV, metafiles: directory searched for fonts before the system ones (repeatable)")
	noSystemFonts := fs.Bool("no-system-fonts", false, "PowerPoint, Excel, Word, CSV, metafiles: use only the fonts under -font-dir")
	hidden := fs.Bool("hidden", false, "PowerPoint, Excel: include hidden slides or sheets (the same as -param hidden=true)")
	var paramFlags stringList
	fs.Var(&paramFlags, "param", "format-specific option as name=value (repeatable; see the formats below)")
	passwordFile := fs.String("password-file", "", "read the password of an encrypted input from this file (- for the standard input; default: $"+passwordEnv+")")
	encrypt := fs.String("encrypt", "auto", "encrypt the output with the password: auto (when the input needs it), always or never")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: bdf generate [flags] <input> <out.bdf | outdir/>\n  an output path ending with / writes the split form\n  a password-protected input is converted with its password and the output encrypted with it")
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
	if *encrypt != "auto" && *encrypt != "always" && *encrypt != "never" {
		usageError("-encrypt must be auto, always or never")
	}
	password, err := readPassword(*passwordFile)
	check(err)
	if *encrypt == "always" && password == "" {
		usageError("-encrypt always needs a password (-password-file or $" + passwordEnv + ")")
	}
	opts.Password = password
	name := ""
	if *format != "auto" {
		if converter.Lookup(*format) == nil {
			usageError("-format " + *format + ": unknown format (want auto or one of " + strings.Join(names, ", ") + ")")
		}
		name = *format
	}
	res, err := converter.ConvertFile(in, name, opts)
	switch {
	case errors.Is(err, converter.ErrUnknownFormat):
		switch ext := strings.ToLower(filepath.Ext(in)); ext {
		case ".ppt":
			usageError(in + ": legacy .ppt files are not supported; save as .pptx first")
		case ".xls":
			usageError(in + ": legacy .xls files are not supported; save as .xlsx first")
		case ".doc":
			usageError(in + ": legacy .doc files are not supported; save as .docx first")
		case ".vsd", ".vss", ".vst":
			usageError(in + ": legacy binary Visio files (" + ext + ") are not supported; save as .vsdx first")
		case ".dwg":
			usageError(in + ": AutoCAD DWG files are not supported; save as .dxf first")
		}
		usageError(in + ": unknown input format (want one of " + strings.Join(names, ", ") + ")")
	case errors.Is(err, converter.ErrPasswordRequired):
		check(fmt.Errorf("%s is encrypted: give its password with -password-file or $%s", in, passwordEnv))
	case errors.Is(err, converter.ErrWrongPassword):
		check(fmt.Errorf("the password does not open %s", in))
	}
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
	summary := res.Summary
	if *encrypt == "always" || *encrypt == "auto" && res.Protected {
		doc.Lock, err = bdf.NewPasswordLock(password, 0)
		check(err)
		summary += ", encrypted"
	}
	check(write(doc, out))
	fmt.Fprintf(os.Stderr, "%s: %s, %d warning(s)\n", out, summary, len(warnings))
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
