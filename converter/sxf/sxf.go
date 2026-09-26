// Package sxf converts SXF drawings, the CAD data exchange format of the
// electronic deliveries of Japanese public works (電子納品), into BDF
// documents: the STEP files (.p21, and .p2z, a P21 file in a zip) and the
// feature comment files (.sfc) of SXF Ver.2 to Ver.3.1.
//
// A drawing becomes one page of its sheet, grown to hold what lies outside
// it, drawn on the background the drawing names (black when it does not
// say, as SXF viewers show drawings; white paper with the light option).
// Compound figures are placed at their scales, and the colors, line types
// and widths are SXF's predefined ones or the drawing's own. See
// docs/design.md §3.14.
package sxf

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"golang.org/x/text/encoding/japanese"
)

// Options controls the conversion.
type Options struct {
	// Title overrides the document title.
	Title string
	// Light draws the drawing on white paper (white lines in black)
	// instead of the background it names.
	Light bool
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
	Doc           *bdf.Document
	Warnings      []string
	Format        string  // "sfc" or "p21"
	W, H          float64 // the page in mm
	EmbeddedFonts int
}

// ConvertFile converts an SXF file.
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
	warnings []string
	warned   map[string]bool
}

func (c *converter) warnf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if c.opts.Warn != nil {
		c.opts.Warn(msg)
		return
	}
	c.warnings = append(c.warnings, msg)
}

func (c *converter) warnOnce(key, format string, args ...any) {
	if c.warned[key] {
		return
	}
	c.warned[key] = true
	c.warnf(format, args...)
}

// decode converts the text of a file into UTF-8: Shift_JIS, as SXF
// specifies, unless it is valid UTF-8 and not valid Shift_JIS.
func decode(b []byte) string {
	ascii := true
	for _, c := range b {
		if c >= 0x80 {
			ascii = false
			break
		}
	}
	if ascii {
		return string(b)
	}
	s, err := japanese.ShiftJIS.NewDecoder().Bytes(b)
	if err == nil && !bytes.ContainsRune(s, utf8.RuneError) {
		return string(s)
	}
	if utf8.Valid(b) {
		return string(b)
	}
	if err == nil {
		return string(s)
	}
	return strings.ToValidUTF8(string(b), "�")
}

// unzip returns the P21 (or SFC) file of a .p2z archive, or its only file.
func unzip(data []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	for _, f := range zr.File {
		name := strings.ToLower(f.Name)
		if strings.HasSuffix(name, ".p21") || strings.HasSuffix(name, ".sfc") || len(zr.File) == 1 {
			if f.UncompressedSize64 > maxSize {
				return nil, fmt.Errorf("%s is larger than %d bytes", f.Name, maxSize)
			}
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(io.LimitReader(rc, maxSize))
		}
	}
	return nil, fmt.Errorf("the archive holds no P21 file")
}

// draw draws the drawing of an SFC or P21 file with rend and returns the
// format, the size of the sheet (0 when it has none) and the title.
func (c *converter) draw(text string, rend *renderer) (format string, sw, sh float64, title string, err error) {
	switch {
	case strings.Contains(text, "/*SXF"):
		feats, err := readSFC(text)
		if err != nil {
			if len(feats) == 0 {
				return "", 0, 0, "", err
			}
			c.warnf("%v; the rest of the file is ignored", err)
		}
		m := buildModel(feats)
		if m.sheet == nil {
			return "", 0, 0, "", errors.New("the file has no sheet")
		}
		if m.bg != nil && !rend.light {
			rend.bg = *m.bg
		}
		sr := &sfcRenderer{renderer: rend, m: m}
		sr.items(m.sheet.items, canvas.Identity)
		for _, name := range slices.Sorted(maps.Keys(m.unknown)) {
			c.warnOnce("feature:"+name, "%d feature(s) %s are not drawn", m.unknown[name], name)
		}
		if w, h, ok := m.sheet.size(); ok {
			sw, sh = w, h
		}
		return "sfc", sw, sh, m.title["title"], nil
	case strings.Contains(text, "ISO-10303-21"):
		p, err := parseP21(text)
		if err != nil {
			return "", 0, 0, "", err
		}
		sw, sh, title = p.draw(rend)
		return "p21", sw, sh, title, nil
	}
	return "", 0, 0, "", errors.New("not an SXF file")
}

// Convert converts a drawing read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	if size > maxSize {
		return nil, fmt.Errorf("sxf: the file is larger than %d bytes", maxSize)
	}
	data := make([]byte, size)
	if _, err := r.ReadAt(data, 0); err != nil && err != io.EOF {
		return nil, fmt.Errorf("sxf: %w", err)
	}
	if bytes.HasPrefix(data, []byte("PK\x03\x04")) {
		d, err := unzip(data)
		if err != nil {
			return nil, fmt.Errorf("sxf: %w", err)
		}
		data = d
	}
	text := decode(data)
	c := &converter{opts: opts, warned: map[string]bool{}}
	doc := bdf.NewDocument()
	doc.Meta.DC.Language = bdf.DCValues{"ja"}
	db := fontdb.New(opts.FontFS, opts.FontDirs, !opts.NoSystemFonts)
	set := fontset.New(db, func(msg string) { c.warnf("%s", msg) })
	if len(db.Faces) == 0 && !opts.SystemFonts {
		c.warnf("no fonts found; text is laid out with estimated metrics and not embedded")
	}
	dr := &cad.Drawing{}
	rend := &renderer{fonts: &cad.Fonts{Set: set}, out: dr, light: opts.Light, bg: bdf.RGB(0, 0, 0), warn: c.warnOnce}
	if opts.Light {
		rend.bg = bdf.RGB(255, 255, 255)
	}
	res := &Result{Doc: doc}
	format, sw, sh, title, err := c.draw(text, rend)
	if err != nil {
		return nil, fmt.Errorf("sxf: %w", err)
	}
	res.Format = format
	doc.Meta.Source = res.Format
	if t := strings.TrimSpace(title); t != "" {
		doc.Meta.DC.Title = bdf.DCValues{t}
	}
	if opts.Title != "" {
		doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}

	// the sheet, grown to hold what lies outside it
	var r0 cad.Rect
	if sw > 0 && sh > 0 {
		r0 = r0.Add(cad.Point{}).Add(cad.Point{X: sw, Y: sh})
	}
	if b := dr.Bounds(); b.Valid() {
		r0 = r0.Union(b)
	}
	if !r0.Valid() {
		r0 = r0.Add(cad.Point{}).Add(cad.Point{X: 420, Y: 297})
	}
	pw, ph := r0.W()*mm, r0.H()*mm
	view := doc.NewView("pages", bdf.ViewFixed, doc.Meta.DC.Title.First())
	page := view.AddPage(float32(pw), float32(ph))
	cvs := canvas.NewBuilder(doc, set)
	cv := cvs.New()
	cv.Obj.SetBBox(0, 0, page.W, page.H)
	cv.Obj.FillColor(rend.bg)
	cv.Obj.FillRect(0, 0, page.W, page.H)
	cv.Drawn = true
	pm := canvas.Matrix{mm, 0, 0, -mm, -r0.Min.X * mm, r0.Max.Y * mm}
	pl := &cad.Plotter{Fonts: rend.fonts, Thin: 0.1 * mm}
	pl.Plot(cv, dr, pm, cad.Rect{}.Add(cad.Point{}).Add(cad.Point{X: pw, Y: ph}))
	if !opts.SystemFonts {
		res.EmbeddedFonts = set.Embed(doc, fontset.EmbedOptions{NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType})
	}
	cvs.Encode()
	set.ReportMissing()
	page.Layers = []bdf.Layer{{Role: bdf.RoleBody, Obj: cv.Hash()}}
	if !opts.NoTextIndex {
		if _, err := doc.BuildTextIndex(view); err != nil {
			c.warnf("text index: %v", err)
		}
	}
	res.W, res.H = r0.W(), r0.H()
	res.Warnings = c.warnings
	return res, nil
}
