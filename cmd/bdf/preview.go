package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/image/rowenc"
	"github.com/shibukawa/bdf/raster/imagebdf"
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
func (p *previewFlags) write(doc *bdf.Document, ropts imagebdf.Options, quiet bool) {
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
func fontFlags(fs *flag.FlagSet) func() imagebdf.Options {
	var dirs stringList
	fs.Var(&dirs, "font-dir", "directory searched for fonts before the system ones, for text in fonts the document refers to by name (repeatable)")
	noSystem := fs.Bool("no-system-fonts", false, "use only the fonts under -font-dir")
	return func() imagebdf.Options { return imagebdf.Options{FontDirs: dirs, NoSystemFonts: *noSystem} }
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
func convertForThumbnail(path, view string, allowPlaintext bool, ropts imagebdf.Options, quiet bool) *bdf.Document {
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
	pages := fs.String("pages", "", "batch pages: all or 1-3,5; output path must contain {page}")
	workers := fs.Int("workers", 0, "batch render workers: 0 chooses automatically, 1..16 sets a maximum")
	scale := fs.Float64("scale", 1, "pixels per unit (a pt page at 1 is drawn at 72 dpi)")
	antialias := fs.String("antialias", "default", "shape, text and clip edge smoothing: default, on or off (faster; independent of scale)")
	mono := fs.String("mono", "", "a monochrome image instead of colour: gray (8 bits a pixel) or bilevel (black and white; 1 bit a pixel in a PNG or TIFF)")
	band := fs.Int("band", imagebdf.DefaultBandRows, "draw a page in bands of this many rows, written to a PNG or TIFF file as each is drawn, with the memory of a band instead of the page (not with -workers above 1); 0 draws the page at once")
	width := fs.Float64("width", 1000, "sheets: width of the region from A1 to draw, in units")
	height := fs.Float64("height", 1000, "sheets: height of the region from A1 to draw, in units")
	allow := fs.Bool("allow-plaintext", false, "draw an encrypted document (read with $"+passwordEnv+"); the image is not encrypted")
	quiet := fs.Bool("q", false, "do not print warnings")
	ropts := fontFlags(fs)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: bdf render [flags] <file.bdf | dir> <out.png | out.jpg | out.webp | out.tif>")
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() != 2 {
		fs.Usage()
		os.Exit(2)
	}
	if *workers < 0 || *workers > 16 {
		badUsage("render", "-workers must be 0..16")
	}
	if *band < 0 {
		badUsage("render", "-band must be 0 or more")
	}
	var pageSet, workersSet bool
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "page" {
			pageSet = true
		}
		if f.Name == "workers" {
			workersSet = true
		}
	})
	if *pages != "" && pageSet {
		badUsage("render", "use either -page or -pages")
	}
	if *pages == "" && workersSet {
		badUsage("render", "-workers requires -pages")
	}
	if *pages != "" && !strings.Contains(fs.Arg(1), "{page}") {
		badUsage("render", "-pages output path must contain {page}")
	}
	var selection converter.Pages
	if *pages != "" && *pages != "all" {
		var err error
		selection, err = converter.ParsePages(*pages)
		if err != nil {
			badUsage("render", err.Error())
		}
	}
	format := thumbnail.FormatOf(fs.Arg(1))
	if format == "" {
		badUsage("render", fs.Arg(1)+": the image must be a .png, .jpg, .webp or .tif file")
	}
	if !(*scale > 0) || *scale > 16 {
		badUsage("render", "-scale must be above 0 and at most 16")
	}
	o := ropts()
	// CLI outputs are consumed once; retaining their images only adds a copy.
	o.ImageCacheBytes = -1
	switch *antialias {
	case "default":
		o.Antialias = imagebdf.AntialiasDefault
	case "on":
		o.Antialias = imagebdf.AntialiasOn
	case "off":
		o.Antialias = imagebdf.AntialiasOff
	default:
		badUsage("render", "-antialias must be default, on or off")
	}
	var m imagebdf.Mono
	switch *mono {
	case "":
	case "gray":
		m = imagebdf.Gray
	case "bilevel":
		m = imagebdf.Bilevel
	default:
		badUsage("render", "-mono must be gray or bilevel")
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
	o.BandRows = *band
	r := imagebdf.New(doc, &o)
	// a page is drawn in bands and written as it is drawn when the file
	// takes rows (PNG, TIFF) and the pages are drawn one at a time
	streamed := *band > 0 && v.Kind != bdf.ViewSheet && (format == thumbnail.PNG || format == thumbnail.TIFF) && (*pages == "" || *workers == 1)
	if *pages != "" {
		if v.Kind == bdf.ViewSheet {
			badUsage("render", "-pages does not apply to sheet views")
		}
		for _, span := range selection {
			if span.From > len(v.Pages) || span.To > len(v.Pages) {
				badUsage("render", "-pages selects a page outside this view")
			}
		}
		numbers := selection.Numbers(len(v.Pages))
		for i := range numbers {
			numbers[i]--
		}
		if streamed {
			var out *rowsFile
			err := r.PagesRows(context.Background(), v, numbers, *scale, m, func(page, y0 int, rows image.Image) error {
				if y0 == 0 {
					var err error
					if out, err = openRowsFile(strings.ReplaceAll(fs.Arg(1), "{page}", strconv.Itoa(page+1)), format, v, page, *scale, m); err != nil {
						return err
					}
				}
				return out.write(y0, rows)
			})
			if err != nil && out != nil {
				out.abandon()
			}
			check(err)
			if !*quiet {
				for _, warning := range r.Warnings() {
					fmt.Fprintln(os.Stderr, "warning:", warning)
				}
			}
			return
		}
		var encoded bytes.Buffer
		check(r.PagesMono(context.Background(), v, numbers, *scale, *workers, m, func(page int, img image.Image) error {
			encoded.Reset()
			if err := encodeImage(&encoded, img, format, *scale); err != nil {
				return err
			}
			return os.WriteFile(strings.ReplaceAll(fs.Arg(1), "{page}", strconv.Itoa(page+1)), encoded.Bytes(), 0o644)
		}))
		if !*quiet {
			for _, warning := range r.Warnings() {
				fmt.Fprintln(os.Stderr, "warning:", warning)
			}
		}
		return
	}
	var err error
	var b bytes.Buffer
	switch {
	case v.Kind == bdf.ViewSheet:
		w, h := int(*width**scale+0.5), int(*height**scale+0.5)
		img, e := r.RegionMono(v, 0, bdf.Rect{W: float32(*width), H: float32(*height)}, max(w, 1), max(h, 1), m)
		if err = e; err == nil {
			err = encodeImage(&b, img, format, *scale)
		}
	case streamed:
		out, e := openRowsFile(fs.Arg(1), format, v, *page-1, *scale, m)
		check(e)
		if err = r.PageRows(v, *page-1, *scale, m, out.write); err != nil {
			out.abandon()
		}
	default:
		img, e := r.PageMono(v, *page-1, *scale, m)
		if err = e; err == nil {
			err = encodeImage(&b, img, format, *scale)
		}
	}
	check(err)
	if !*quiet {
		for _, w := range r.Warnings() {
			fmt.Fprintln(os.Stderr, "warning:", w)
		}
	}
	if !streamed {
		check(os.WriteFile(fs.Arg(1), b.Bytes(), 0o644))
	}
}

// encodeImage writes an image of render: a PNG or TIFF as the rows of a
// page are written (see rowsFile), so that a page drawn whole and one
// drawn in bands are the same file; the rest as thumbnail.EncodeImage.
func encodeImage(w io.Writer, img image.Image, format string, scale float64) error {
	if format != thumbnail.PNG && format != thumbnail.TIFF {
		return thumbnail.EncodeImage(w, img, format)
	}
	kind := rowenc.KindOf(img)
	if kind == rowenc.RGBA {
		if _, ok := img.(*image.RGBA); !ok {
			return thumbnail.EncodeImage(w, img, format)
		}
	}
	return rowenc.Encode(w, format, img, kind, rowsOptions(scale))
}

// rowsOptions is how the PNG and TIFF files of render are written: the
// PNG compressed as thumbnail.EncodeImage compresses one, the resolution
// noted from the scale.
func rowsOptions(scale float64) *rowenc.Options {
	return &rowenc.Options{Level: 9, DPI: scale * 72}
}

// rowsFile is an image file written from the rows of a page as they are
// drawn (see rowenc), and closed with the last row.
type rowsFile struct {
	name string
	f    *os.File
	w    *bufio.Writer
	enc  rowenc.Writer
	h    int
	rows int
}

// openRowsFile starts an image file of a page of a view at scale with the
// pixels of the mono.
func openRowsFile(name, format string, v *bdf.View, page int, scale float64, m imagebdf.Mono) (*rowsFile, error) {
	w, h, err := imagebdf.PagePixels(v, page, scale)
	if err != nil {
		return nil, err
	}
	kind := rowenc.RGBA
	switch m {
	case imagebdf.Gray:
		kind = rowenc.Gray
	case imagebdf.Bilevel:
		kind = rowenc.Bilevel
	}
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	out := &rowsFile{name: name, f: f, w: bufio.NewWriterSize(f, 1<<20), h: h}
	if out.enc, err = rowenc.New(out.w, format, w, h, kind, rowsOptions(scale)); err != nil {
		out.abandon()
		return nil, err
	}
	return out, nil
}

// write writes the rows of a band, and ends the file with the last.
func (out *rowsFile) write(_ int, rows image.Image) error {
	if err := out.enc.WriteRows(rows); err != nil {
		return err
	}
	out.rows += rows.Bounds().Dy()
	if out.rows < out.h {
		return nil
	}
	if err := out.enc.Close(); err != nil {
		return err
	}
	if err := out.w.Flush(); err != nil {
		return err
	}
	return out.f.Close()
}

// abandon removes a file that could not be written whole.
func (out *rowsFile) abandon() {
	out.f.Close()
	os.Remove(out.name)
}

func badUsage(cmd, msg string) {
	fmt.Fprintf(os.Stderr, "bdf %s: %s\n", cmd, msg)
	os.Exit(2)
}
