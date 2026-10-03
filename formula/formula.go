// Package formula lays out mathematical formulas written in LaTeX or
// MathML and returns them as vector graphics: a BDF object of filled and
// stroked paths, which any renderer of BDF objects draws — into an image
// (raster/imagebdf), on an Ebitengine screen (raster/ebitenginebdf), or on
// a browser's canvas as the page of a document (@bdfkit/render).
//
//	ts, _ := formula.New(nil)
//	f := formula.ParseTeX(`x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}`)
//	obj := ts.Layout(f, formula.Style{Size: 24, Display: true}).Object()
//	img, _ := imagebdf.Object(obj, 2, nil)
//
// Formulas are laid out the way TeX and MathML Core lay them out, with the
// parameters and the glyph variants of the OpenType MATH table of a formula
// font: script sizes and shifts, fraction and radical gaps, the spacing
// between atoms, large operators, and delimiters, radicals and accents that
// grow. It is the engine that the converters of this module lay the
// formulas of Word, PowerPoint, Excel, HTML, EPUB, Markdown and draw.io
// documents out with. The formula font is STIX Two Math, which the package
// embeds, unless Options name another.
//
// STIX Two Math is licensed under the SIL Open Font License 1.1
// (fonts/OFL.txt; Copyright 2001-2021 The STIX Fonts Project Authors, with
// Reserved Font Name "TM Math"). A program that imports this package holds
// the font; fonts/README.md says where the file comes from.
package formula

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/fontdb"
	"github.com/shibukawa/bdf/internal/mathlayout"
)

//go:embed fonts/STIXTwoMath-Regular.otf
var fontFS embed.FS

// Options are the fonts of a Typesetter.
type Options struct {
	// Math is the formula font, a TrueType or OpenType font with a MATH
	// table (STIX Two Math, Latin Modern Math, Cambria Math …). nil is
	// the embedded STIX Two Math.
	Math []byte
	// Text are fonts for the characters the formula font lacks, such as
	// East Asian text in \text{…}, tried in order.
	Text [][]byte
}

// Typesetter lays formulas out with one set of fonts. It is safe for
// concurrent use.
type Typesetter struct {
	mu     sync.Mutex
	engine *mathlayout.Engine
}

// New returns a typesetter. A nil opts lays formulas out with the embedded
// formula font alone.
func New(opts *Options) (*Typesetter, error) {
	if opts == nil {
		opts = &Options{}
	}
	var fsys fs.FS = fontFS
	mathPath := "fonts/STIXTwoMath-Regular.otf"
	if opts.Math != nil || len(opts.Text) > 0 {
		m := memFS{}
		if opts.Math != nil {
			mathPath = "math.otf"
			m[mathPath] = opts.Math
		} else {
			m[mathPath], _ = fontFS.ReadFile(mathPath)
		}
		for i, data := range opts.Text {
			m[fmt.Sprintf("text-%04d.otf", i)] = data
		}
		fsys = m
	}
	db := fontdb.New(fsys, nil, false)
	faces := map[*fontdb.Face]*mathlayout.Face{}
	load := func(f *fontdb.Face) *mathlayout.Face {
		if fa, ok := faces[f]; ok {
			return fa
		}
		var fa *mathlayout.Face
		if l, err := f.Load(); err == nil {
			fa = &mathlayout.Face{Loaded: l}
		}
		faces[f] = fa
		return fa
	}
	var fonts mathlayout.Fonts
	var text []*fontdb.Face
	for _, f := range db.Faces {
		switch {
		case f.Path != mathPath:
			text = append(text, f)
		case f.Math && fonts.Math == nil:
			fonts.Math = load(f)
		}
	}
	if fonts.Math == nil {
		return nil, errors.New("formula: the formula font cannot be read or has no MATH table")
	}
	if len(text) < len(opts.Text) {
		return nil, errors.New("formula: a text font cannot be read")
	}
	fonts.Text = func(family string, r rune, bold, italic bool) *mathlayout.Face {
		for _, f := range text {
			if f.HasRune(r) {
				if fa := load(f); fa != nil {
					return fa
				}
			}
		}
		return nil
	}
	return &Typesetter{engine: mathlayout.New(fonts)}, nil
}

// Formula is a formula read from LaTeX or MathML.
type Formula struct {
	node    mathlayout.Node
	display bool
}

// ParseTeX reads a formula written in LaTeX's math mode (without the $ or
// \[ \] around it), with the commands of amsmath and amssymb that MathJax
// and KaTeX pages use. It does not fail: what it does not know it keeps as
// text.
func ParseTeX(src string) *Formula {
	return &Formula{node: mathlayout.ParseTeX(src)}
}

// ParseMathML reads the first math element of src (presentation MathML,
// as HTML holds it).
func ParseMathML(src string) (*Formula, error) {
	doc, err := html.Parse(strings.NewReader(src))
	if err != nil {
		return nil, err
	}
	var find func(n *html.Node) *html.Node
	find = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && n.Data == "math" {
			return n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if m := find(c); m != nil {
				return m
			}
		}
		return nil
	}
	m := find(doc)
	if m == nil {
		return nil, errors.New("formula: no math element")
	}
	f := &Formula{}
	f.node, f.display = mathlayout.ParseMathML(m)
	return f, nil
}

// Display reports whether the formula asks for display style (a MathML
// math element with display="block").
func (f *Formula) Display() bool { return f.display }

// Text returns the formula in a linear notation, such as
// "x=(−b±√(b^2−4ac))/(2a)": what a screen reader, a search or a copy takes
// for the formula.
func (f *Formula) Text() string { return mathlayout.Linear(f.node) }

// Style is how a formula is set.
type Style struct {
	// Size is the font size in the units of the result (12 when it is 0).
	Size float64
	// Display sets the formula in display style, as on a line of its own:
	// larger operators with their limits above and below, taller
	// fractions. A formula that asks for it (Formula.Display) gets it
	// too.
	Display bool
	// Color is the color of the formula where it sets none itself (\color);
	// black when it is 0.
	Color bdf.Color
	Bold  bool
}

// Layout lays a formula out.
func (t *Typesetter) Layout(f *Formula, st Style) *Layout {
	t.mu.Lock()
	defer t.mu.Unlock()
	box := t.engine.Layout(f.node, mathlayout.Style{Size: st.Size, Display: st.Display || f.display, Color: st.Color, Bold: st.Bold})
	return &Layout{t: t, box: box, text: mathlayout.Linear(f.node)}
}

// Layout is a formula laid out. Its origin is the left end of its baseline
// and y grows downwards, in the units of Style.Size.
type Layout struct {
	t    *Typesetter
	box  *mathlayout.Box
	text string
}

// Width is the advance of the formula: where what follows it on the line
// starts.
func (l *Layout) Width() float64 { return l.box.W }

// Height is how far the formula reaches above its baseline.
func (l *Layout) Height() float64 { return l.box.H }

// Depth is how far the formula reaches below its baseline.
func (l *Layout) Depth() float64 { return l.box.D }

// Ink returns the extent of what the formula draws, which may reach past
// its width, height and depth (italic overhangs, accents).
func (l *Layout) Ink() bdf.Rect { return l.box.Ink() }

// Text returns the formula in a linear notation (see Formula.Text).
func (l *Layout) Text() string { return l.text }

// memFS is the fonts of Options as a file system for fontdb.
type memFS map[string][]byte

func (m memFS) Open(name string) (fs.File, error) {
	if name == "." {
		return &memDir{m: m}, nil
	}
	data, ok := m[name]
	if !ok {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	return &memFile{Reader: bytes.NewReader(data), name: name, size: int64(len(data))}, nil
}

type memFile struct {
	*bytes.Reader
	name string
	size int64
}

func (f *memFile) Stat() (fs.FileInfo, error) { return memInfo{name: f.name, size: f.size}, nil }
func (f *memFile) Close() error               { return nil }

type memDir struct {
	m    memFS
	read bool
}

func (d *memDir) Stat() (fs.FileInfo, error) { return memInfo{name: ".", dir: true}, nil }
func (d *memDir) Read([]byte) (int, error)   { return 0, io.EOF }
func (d *memDir) Close() error               { return nil }
func (d *memDir) ReadDir(n int) ([]fs.DirEntry, error) {
	if d.read {
		if n > 0 {
			return nil, io.EOF
		}
		return nil, nil
	}
	d.read = true
	var out []fs.DirEntry
	for name, data := range d.m {
		out = append(out, fs.FileInfoToDirEntry(memInfo{name: name, size: int64(len(data))}))
	}
	return out, nil
}

type memInfo struct {
	name string
	size int64
	dir  bool
}

func (i memInfo) Name() string { return i.name }
func (i memInfo) Size() int64  { return i.size }
func (i memInfo) Mode() fs.FileMode {
	if i.dir {
		return fs.ModeDir | 0o555
	}
	return 0o444
}
func (i memInfo) ModTime() time.Time { return time.Time{} }
func (i memInfo) IsDir() bool        { return i.dir }
func (i memInfo) Sys() any           { return nil }
