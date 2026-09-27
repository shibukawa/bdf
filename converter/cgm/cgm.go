// Package cgm converts Computer Graphics Metafiles (CGM, ISO/IEC 8632)
// into BDF documents: the binary encoding (Part 3), which CAD programs,
// technical illustration tools and WebCGM write, and the clear text
// encoding (Part 4), of versions 1 to 4. Metafiles compressed with gzip
// (.cgz) are read too.
//
// Each picture of the metafile becomes a page of a fixed view. Its size
// comes from the metric scaling of the picture (millimetres per VDC unit)
// or a device viewport in millimetres; an abstract picture is fitted to
// 297 mm, the long side of A4. Lines, markers, filled areas (hollow,
// solid, hatched, patterned and interpolated), closed figures, arcs,
// ellipses, conics, Bézier curves and B-splines, text (with its path,
// alignment, restricted boxes and character sets), cell arrays and tiles
// are drawn with converter/internal/cad. See docs/design.md §3.20.
package cgm

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/imgconv"
	"github.com/shibukawa/bdf/internal/fontdb"
)

// Options controls the conversion.
type Options struct {
	// Title overrides the document title.
	Title string
	// Pages selects 1-based pictures (see converter.Pages); nil converts
	// all of them.
	Pages conv.Pages
	// Images controls how cell arrays and tiles are stored (see imgconv).
	Images imgconv.Options
	// FontFS holds fonts that are not in the local file system; it is
	// searched before FontDirs (see converter.Options.FontFS).
	FontFS fs.FS
	// FontDirs are searched for fonts before the system font directories.
	FontDirs []string
	// NoSystemFonts restricts font lookup to FontFS and FontDirs.
	NoSystemFonts bool
	// SystemFonts refers to fonts by family name instead of embedding the
	// fonts used for layout.
	SystemFonts bool
	// NoSubset embeds whole fonts instead of the glyphs in use.
	NoSubset bool
	// NoWOFF2 stores embedded fonts as TrueType/OpenType instead of WOFF2.
	NoWOFF2 bool
	// IgnoreFSType embeds fonts whose OS/2 fsType forbids embedding or
	// subsetting. Set it only when you hold the rights to embed the fonts.
	IgnoreFSType bool
	// NoTextIndex skips building the text index part.
	NoTextIndex bool
	// Warn receives non-fatal problems; when nil they are collected in
	// Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of a conversion.
type Result struct {
	Doc      *bdf.Document
	Warnings []string
	// Encoding is "binary" or "clear text".
	Encoding string
	// Pictures is the number of pictures in the metafile, Pages the number
	// converted.
	Pictures, Pages int
	EmbeddedFonts   int
}

// ConvertFile converts a metafile.
func ConvertFile(path string, opts *Options) (*Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return Convert(f, st.Size(), opts)
}

const maxSize = 1 << 30

type converter struct {
	opts     *Options
	doc      *bdf.Document
	warnings []string
	warned   map[string]bool
	quiet    bool
}

func (c *converter) warnOnce(key, format string, args ...any) {
	if c.quiet || c.warned[key] {
		return
	}
	c.warned[key] = true
	msg := fmt.Sprintf(format, args...)
	if c.opts.Warn != nil {
		c.opts.Warn(msg)
		return
	}
	c.warnings = append(c.warnings, msg)
}

// readElements splits a metafile into its elements.
func readElements(data []byte) (elems []element, text, truncated bool, err error) {
	switch encodingOf(data) {
	case "binary":
		r := &binReader{data: data}
		for {
			e, ok := r.next()
			if !ok {
				break
			}
			elems = append(elems, e)
			if e.code == eEndMetafile {
				break
			}
		}
		return elems, false, r.truncated, nil
	case "clear text":
		r := &textReader{s: string(data)}
		for {
			e, ok := r.next()
			if !ok {
				break
			}
			elems = append(elems, e)
			if e.code == eEndMetafile {
				break
			}
		}
		return elems, true, false, nil
	}
	return nil, false, false, errors.New("cgm: not a CGM metafile")
}

// gunzip returns the metafile of a gzip file.
func gunzip(data []byte) ([]byte, error) {
	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(io.LimitReader(zr, maxSize+1))
	if err != nil && len(b) == 0 {
		return nil, err
	}
	if len(b) > maxSize {
		return nil, fmt.Errorf("the metafile is larger than %d bytes", maxSize)
	}
	return b, nil
}

// Convert converts a metafile read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	if size > maxSize {
		return nil, fmt.Errorf("cgm: the file is larger than %d bytes", maxSize)
	}
	data := make([]byte, size)
	if _, err := r.ReadAt(data, 0); err != nil && err != io.EOF {
		return nil, fmt.Errorf("cgm: %w", err)
	}
	if bytes.HasPrefix(data, []byte{0x1f, 0x8b}) {
		d, err := gunzip(data)
		if err != nil {
			return nil, fmt.Errorf("cgm: %w", err)
		}
		data = d
	}
	elems, text, truncated, err := readElements(data)
	if err != nil {
		return nil, err
	}
	doc := bdf.NewDocument()
	c := &converter{opts: opts, doc: doc, warned: map[string]bool{}}
	res := &Result{Doc: doc, Encoding: "binary"}
	if text {
		res.Encoding = "clear text"
	}
	if truncated {
		c.warnOnce("truncated", "the metafile is truncated; what it holds is drawn")
	}
	doc.Meta.Source = "cgm"

	// the strings of text first, to tell how undeclared multibyte text is
	// encoded
	scan := newInterp(&converter{opts: opts, doc: doc, warned: map[string]bool{}, quiet: true}, text)
	scan.scan = true
	scan.run(iterate(elems))
	fallback, name := guessEncoding(scan.strs)
	switch {
	case name == "shift_jis" || scan.declares("B", set94Multi) || scan.declares("@", set94Multi):
		doc.Meta.DC.Language = bdf.DCValues{"ja"}
	case scan.declares("A", set94Multi):
		doc.Meta.DC.Language = bdf.DCValues{"zh-Hans"}
	case scan.declares("C", set94Multi):
		doc.Meta.DC.Language = bdf.DCValues{"ko"}
	}

	db := fontdb.New(opts.FontFS, opts.FontDirs, !opts.NoSystemFonts)
	set := fontset.New(db, func(msg string) { c.warnOnce("font:"+msg, "%s", msg) })
	if len(db.Faces) == 0 && !opts.SystemFonts {
		c.warnOnce("no fonts", "no fonts found; text is laid out with estimated metrics and not embedded")
	}
	in := newInterp(c, text)
	in.fonts = &cad.Fonts{Set: set}
	in.fallback = fallback
	sel := map[int]bool{}
	pages := opts.Pages.Numbers(scan.count)
	for _, n := range pages {
		if n < 1 || n > scan.count {
			return nil, fmt.Errorf("cgm: page %d out of range (1-%d)", n, scan.count)
		}
		sel[n] = true
	}
	if pages != nil {
		in.selected = func(n int) bool { return sel[n] }
	}
	in.run(iterate(elems))
	res.Pictures = in.count
	if in.count == 0 {
		return nil, errors.New("cgm: the metafile has no pictures")
	}
	if t := strings.TrimSpace(in.title); t != "" {
		doc.Meta.DC.Title = bdf.DCValues{t}
	}
	if opts.Title != "" {
		doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}

	// the pages, in the order asked for
	order := in.pictures
	if pages != nil {
		byNumber := map[int]*picture{}
		for _, p := range in.pictures {
			byNumber[p.number] = p
		}
		order = nil
		done := map[int]bool{}
		for _, n := range pages {
			if p := byNumber[n]; p != nil && !done[n] {
				order = append(order, p)
				done[n] = true
			}
		}
	}
	view := doc.NewView("pages", bdf.ViewFixed, doc.Meta.DC.Title.First())
	cvs := canvas.NewBuilder(doc, set)
	pl := &cad.Plotter{Fonts: in.fonts, Thin: 0.1 * mm}
	type placed struct {
		page *bdf.Page
		cv   *canvas.Canvas
	}
	var laid []placed
	for _, pic := range order {
		pw, ph := float32(pic.w*mm), float32(pic.h*mm)
		page := view.AddPage(pw, ph)
		cv := cvs.New()
		cv.Obj.SetBBox(0, 0, pw, ph)
		cv.Obj.FillColor(pic.bg)
		cv.Obj.FillRect(0, 0, pw, ph)
		cv.Drawn = true
		pl.Plot(cv, pic.d, pic.m, cad.Rect{}.Add(cad.Point{}).Add(cad.Point{X: float64(pw), Y: float64(ph)}))
		laid = append(laid, placed{page, cv})
	}
	if !opts.SystemFonts {
		res.EmbeddedFonts = set.Embed(doc, fontset.EmbedOptions{NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType})
	}
	cvs.Encode()
	set.ReportMissing()
	for _, l := range laid {
		l.page.Layers = []bdf.Layer{{Role: bdf.RoleBody, Obj: l.cv.Hash()}}
	}
	if !opts.NoTextIndex {
		if _, err := doc.BuildTextIndex(view); err != nil {
			c.warnOnce("text index", "text index: %v", err)
		}
	}
	res.Pages = len(view.Pages)
	res.Warnings = c.warnings
	return res, nil
}

// iterate returns a function that yields the elements in turn.
func iterate(elems []element) func() (element, bool) {
	i := 0
	return func() (element, bool) {
		if i >= len(elems) {
			return element{}, false
		}
		i++
		return elems[i-1], true
	}
}

func newInterp(c *converter, text bool) *interp {
	in := &interp{c: c, pr: defaultPrecisions(), text: text, colourModel: 1, st: newState()}
	if text {
		in.pr.colourBits = 255 // the largest component value
	}
	return in
}

// declares reports whether the character set list names a set.
func (in *interp) declares(final string, typ int) bool {
	for _, cs := range in.charsets {
		if cs.typ == typ && cs.final() == final {
			return true
		}
	}
	return false
}

// encodingOf tells the encoding of a metafile from its first element,
// BEGIN METAFILE: in the binary encoding a word of class 0, id 1 and the
// length of its name, followed by METAFILE VERSION; in clear text the
// name BEGMF.
func encodingOf(data []byte) string {
	if len(data) >= 4 {
		pos := 0
		// skip no-ops
		for pos+2 <= len(data) && data[pos] == 0 && data[pos+1]&0xe0 == 0 {
			n := int(data[pos+1] & 0x1f)
			pos += 2 + n + n%2
		}
		if pos+2 <= len(data) && data[pos] == 0 && data[pos+1]&0xe0 == 0x20 {
			r := &binReader{data: data[pos:]}
			if e, ok := r.next(); ok && e.code == eBeginMetafile {
				if next, ok := r.next(); ok && next.class == 1 {
					return "binary"
				}
				if r.pos >= len(r.data) {
					return "binary" // only the name
				}
			}
		}
	}
	// clear text: BEGMF after white space and comments
	r := &textReader{s: string(data[:min(len(data), 4096)])}
	r.skip()
	if rest := r.s[r.pos:]; len(rest) >= 5 {
		w := rest[:min(len(rest), 16)]
		if i := strings.IndexAny(w, " \t\r\n;'\"/("); i > 0 {
			w = w[:i]
		}
		if keyword(w) == "BEGMF" {
			return "clear text"
		}
	}
	return ""
}
