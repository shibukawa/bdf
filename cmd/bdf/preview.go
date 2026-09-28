package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/raster"
	"github.com/shibukawa/bdf/thumbnail"
)

// plaintextNote explains why the thumbnail and text of an encrypted
// document are not written by default.
const plaintextNote = "the document is encrypted and they are not (-allow-plaintext writes them anyway)"

// previewFlags are the flags of the thumbnail and search text that
// generate writes next to the document.
type previewFlags struct {
	thumbnail      *string
	size           *int
	mode           *string
	sheetDPI       *float64
	text           *string
	allowPlaintext *bool
}

func addPreviewFlags(fs *flag.FlagSet) *previewFlags {
	return &previewFlags{
		thumbnail:      fs.String("thumbnail", "", "also write a thumbnail image of the document (.png, .jpg or .webp; see bdf thumbnail)"),
		size:           fs.Int("thumbnail-size", thumbnail.DefaultSize, "thumbnail size in pixels: the side of a cropped one, the longer side of a fitted one"),
		mode:           fs.String("thumbnail-mode", "auto", "thumbnail layout: auto, crop (a square from the top-left of the first page) or fit (the whole first page)"),
		sheetDPI:       fs.Float64("thumbnail-sheet-dpi", thumbnail.DefaultSheetDPI, sheetDPIUsage),
		text:           fs.String("text", "", "also write the text of the document for a search index, as JSON (- for the standard output; see bdf text)"),
		allowPlaintext: fs.Bool("allow-plaintext", false, "write the thumbnail and text of an encrypted document too (they are not encrypted)"),
	}
}

// write writes the thumbnail and text generate was asked for.
func (p *previewFlags) write(doc *bdf.Document, ropts raster.Options, quiet bool) {
	if *p.thumbnail == "" && *p.text == "" {
		return
	}
	if doc.Lock != nil && !*p.allowPlaintext {
		fmt.Fprintln(os.Stderr, "bdf: no thumbnail or text written:", plaintextNote)
		return
	}
	if *p.thumbnail != "" {
		mode, err := thumbnail.ParseMode(*p.mode)
		if err != nil {
			usageError(err.Error())
		}
		check(writeThumbnail(doc, *p.thumbnail, &thumbnail.Options{Size: *p.size, Mode: mode, SheetDPI: *p.sheetDPI, Raster: ropts}, quiet))
	}
	if *p.text != "" {
		check(writeText(doc, *p.text))
	}
}

// sheetDPIUsage explains -sheet-dpi and -thumbnail-sheet-dpi.
const sheetDPIUsage = "resolution of a sheet in the thumbnail: the square from A1 is the size at this dpi (a smaller thumbnail shows fewer cells), 96 to 480 units"

func writeThumbnail(doc *bdf.Document, out string, opts *thumbnail.Options, quiet bool) error {
	format := thumbnail.FormatOf(out)
	if format == "" {
		return fmt.Errorf("%s: the thumbnail must be a .png, .jpg or .webp file", out)
	}
	res, err := thumbnail.Make(doc, opts)
	if err != nil {
		return err
	}
	if !quiet {
		for _, w := range res.Warnings {
			fmt.Fprintln(os.Stderr, "warning: thumbnail:", w)
		}
	}
	var b bytes.Buffer
	if err := thumbnail.Encode(&b, res.Image, format); err != nil {
		return err
	}
	if err := os.WriteFile(out, b.Bytes(), 0o644); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "%s: %d×%d %s thumbnail (%s)\n", out, res.Image.Rect.Dx(), res.Image.Rect.Dy(), format, res.Mode)
	return nil
}

func writeText(doc *bdf.Document, out string) error {
	st, err := doc.SearchText()
	if err != nil {
		return err
	}
	var w io.Writer = os.Stdout
	var f *os.File
	if out != "-" {
		if f, err = os.Create(out); err != nil {
			return err
		}
		w = f
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err = enc.Encode(st)
	if f != nil {
		err = errors.Join(err, f.Close())
	}
	return err
}

// loadDocument reads a document in either form for drawing or its text: an
// encrypted one with $BDF_PASSWORD, and only when allowed.
func loadDocument(path string, allowPlaintext bool) *bdf.Document {
	r, err := open(path)
	check(err)
	defer r.Close()
	if r.Locked() {
		check(fmt.Errorf("%s is encrypted: set $%s", path, passwordEnv))
	}
	if r.Encrypted() && !allowPlaintext {
		check(fmt.Errorf("%s: %s", path, plaintextNote))
	}
	d, err := r.ToDocument()
	check(err)
	return d
}

// fontFlags are the font lookup flags of the drawing commands.
func fontFlags(fs *flag.FlagSet) func() raster.Options {
	var dirs stringList
	fs.Var(&dirs, "font-dir", "directory searched for fonts before the system ones, for text in fonts the document refers to by name (repeatable)")
	noSystem := fs.Bool("no-system-fonts", false, "use only the fonts under -font-dir")
	return func() raster.Options { return raster.Options{FontDirs: dirs, NoSystemFonts: *noSystem} }
}

// thumbnailCmd draws the thumbnail of a document, or of an input in a format
// bdf generate converts.
func thumbnailCmd(args []string) {
	fs := flag.NewFlagSet("thumbnail", flag.ExitOnError)
	size := fs.Int("size", thumbnail.DefaultSize, "size in pixels: the side of a cropped thumbnail, the longer side of a fitted one")
	mode := fs.String("mode", "auto", "layout: auto (from the kind of document), crop (a square from the top-left of the first page) or fit (the whole first page)")
	view := fs.String("view", "", "id of the view to draw (default: the first)")
	sheetDPI := fs.Float64("sheet-dpi", thumbnail.DefaultSheetDPI, sheetDPIUsage)
	allow := fs.Bool("allow-plaintext", false, "draw an encrypted document or a password-protected input (read with $"+passwordEnv+"); the thumbnail is not encrypted")
	quiet := fs.Bool("q", false, "do not print warnings")
	ropts := fontFlags(fs)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: bdf thumbnail [flags] <file.bdf | dir | input> <out.png | out.jpg | out.webp>\n  an input in a format bdf generate converts (a PDF, a .docx …) is converted for the thumbnail: only its first page, unless -view")
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() != 2 {
		fs.Usage()
		os.Exit(2)
	}
	m, err := thumbnail.ParseMode(*mode)
	if err != nil {
		badUsage("thumbnail", err.Error())
	}
	in, ro := fs.Arg(0), ropts()
	var doc *bdf.Document
	if isStored(in) {
		doc = loadDocument(in, *allow)
	} else {
		doc = convertForThumbnail(in, *view, *allow, ro, *quiet)
	}
	check(writeThumbnail(doc, fs.Arg(1), &thumbnail.Options{Size: *size, Mode: m, View: *view, SheetDPI: *sheetDPI, Raster: ro}, *quiet))
}

// isStored reports that path is a BDF document: a directory (the split
// form) or a file that starts with the magic number. A path that cannot be
// read is left to open to explain.
func isStored(path string) bool {
	st, err := os.Stat(path)
	if err != nil || st.IsDir() {
		return true
	}
	f, err := os.Open(path)
	if err != nil {
		return true
	}
	defer f.Close()
	var head [4]byte
	_, err = io.ReadFull(f, head[:])
	return err == nil && head == bdf.Magic
}

// protectedNote explains why the thumbnail of a password-protected input is
// not drawn by default.
const protectedNote = "the input is password-protected and the thumbnail would not be (-allow-plaintext draws it anyway)"

// convertForThumbnail converts an input for its thumbnail, which shows the
// first page of the first view: the pages after it are not converted (the
// thumbnail is the same, and a long PDF takes seconds instead of
// milliseconds). Another view is converted whole. A password-protected
// input is opened with $BDF_PASSWORD, and drawn only when allowed.
func convertForThumbnail(path, view string, allowPlaintext bool, ropts raster.Options, quiet bool) *bdf.Document {
	opts := &converter.Options{FontDirs: ropts.FontDirs, NoSystemFonts: ropts.NoSystemFonts, Password: os.Getenv(passwordEnv)}
	if view == "" {
		opts.Pages = converter.PageList(1)
	}
	res, err := converter.ConvertFile(path, "", opts)
	switch {
	case errors.Is(err, converter.ErrUnknownFormat):
		badUsage("thumbnail", path+": neither a bdf document nor an input format bdf generate converts")
	case errors.Is(err, converter.ErrPasswordRequired):
		check(fmt.Errorf("%s is encrypted: set $%s", path, passwordEnv))
	case errors.Is(err, converter.ErrWrongPassword):
		check(fmt.Errorf("the password does not open %s", path))
	}
	check(err)
	if res.Protected && !allowPlaintext {
		check(fmt.Errorf("%s: %s", path, protectedNote))
	}
	if !quiet {
		for _, w := range res.Warnings {
			fmt.Fprintln(os.Stderr, "warning:", w)
		}
	}
	return res.Doc
}

// textCmd writes the text of a document for a search index.
func textCmd(args []string) {
	fs := flag.NewFlagSet("text", flag.ExitOnError)
	allow := fs.Bool("allow-plaintext", false, "read an encrypted document (with $"+passwordEnv+"); the text is written in plain text")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: bdf text [flags] <file.bdf | dir> [out.json]\n  writes the metadata and the text of each page as JSON (default: to the standard output)")
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() < 1 || fs.NArg() > 2 {
		fs.Usage()
		os.Exit(2)
	}
	out := "-"
	if fs.NArg() == 2 {
		out = fs.Arg(1)
	}
	check(writeText(loadDocument(fs.Arg(0), *allow), out))
}

// renderCmd draws a page of a document.
func renderCmd(args []string) {
	fs := flag.NewFlagSet("render", flag.ExitOnError)
	view := fs.String("view", "", "id of the view (default: the first)")
	page := fs.Int("page", 1, "page to draw, from 1 (fixed, flow and scroll views)")
	scale := fs.Float64("scale", 1, "pixels per unit (a pt page at 1 is drawn at 72 dpi)")
	width := fs.Float64("width", 1000, "sheets: width of the region from A1 to draw, in units")
	height := fs.Float64("height", 1000, "sheets: height of the region from A1 to draw, in units")
	allow := fs.Bool("allow-plaintext", false, "draw an encrypted document (read with $"+passwordEnv+"); the image is not encrypted")
	quiet := fs.Bool("q", false, "do not print warnings")
	ropts := fontFlags(fs)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: bdf render [flags] <file.bdf | dir> <out.png | out.jpg | out.webp>")
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() != 2 {
		fs.Usage()
		os.Exit(2)
	}
	format := thumbnail.FormatOf(fs.Arg(1))
	if format == "" {
		badUsage("render", fs.Arg(1)+": the image must be a .png, .jpg or .webp file")
	}
	if !(*scale > 0) || *scale > 16 {
		badUsage("render", "-scale must be above 0 and at most 16")
	}
	doc := loadDocument(fs.Arg(0), *allow)
	var v *bdf.View
	for _, vv := range doc.Views {
		if *view == "" || vv.ID == *view {
			v = vv
			break
		}
	}
	if v == nil {
		check(fmt.Errorf("%s: no view %q", fs.Arg(0), *view))
	}
	o := ropts()
	r := raster.New(doc, &o)
	var err error
	var b bytes.Buffer
	if v.Kind == bdf.ViewSheet {
		w, h := int(*width**scale+0.5), int(*height**scale+0.5)
		img, e := r.Region(v, 0, bdf.Rect{W: float32(*width), H: float32(*height)}, max(w, 1), max(h, 1))
		if err = e; err == nil {
			err = thumbnail.Encode(&b, img, format)
		}
	} else {
		img, e := r.Page(v, *page-1, *scale)
		if err = e; err == nil {
			err = thumbnail.Encode(&b, img, format)
		}
	}
	check(err)
	if !*quiet {
		for _, w := range r.Warnings() {
			fmt.Fprintln(os.Stderr, "warning:", w)
		}
	}
	check(os.WriteFile(fs.Arg(1), b.Bytes(), 0o644))
}

func badUsage(cmd, msg string) {
	fmt.Fprintf(os.Stderr, "bdf %s: %s\n", cmd, msg)
	os.Exit(2)
}
