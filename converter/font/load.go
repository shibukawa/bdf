package font

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"slices"

	"github.com/shibukawa/bdf/internal/otlayout"
	"github.com/shibukawa/bdf/internal/sfnt"
	"github.com/shibukawa/bdf/woff2"
)

// file is a font file: the fonts it holds and how it stores them.
type file struct {
	// container is "WOFF" or "WOFF2" for web fonts, "" for sfnt files.
	container  string
	collection bool
	faces      []*face
	size       int // bytes of the file

	// layout are the layout tables read, by the bytes they are read from:
	// the fonts of a collection share tables, and what a table may list
	// (otlayout.Table.Truncated) is then counted once.
	layout map[tableKey]any
}

// tableKey identifies a table of a file: its tag and its bytes.
type tableKey struct {
	tag   string
	start *byte
	size  int
}

// face is one font of a file.
type face struct {
	fl    *file
	index int // in the collection
	// size is that of the font on its own: the file itself, or for a font
	// of a collection or a web font the sfnt made of its tables.
	size int
	f    *sfnt.Font
	out  *sfnt.Outlines
	// ready is set once prepare has read what the views show.
	ready bool

	gsub, gpos *otlayout.Table
	gdef       *otlayout.GDEF

	names    []nameRec  // the name table
	namesCut bool       // the name table holds more than was read (see parseNames)
	glyphs   []string   // glyph names by glyph ID ("" when the font has none)
	chars    [][]rune   // characters by glyph ID
	runes    []rune     // characters the font maps, ascending
	advances []float64  // em, by glyph ID
	bounds   []glyphBox // by glyph ID
	upem     float64
}

// glyphBox is the bounding box of a glyph's outline in em (y up); ok is
// false for glyphs without an outline.
type glyphBox struct {
	x0, y0, x1, y1 float64
	ok             bool
}

// maxFaces bounds the fonts of a collection that are read.
const maxFaces = 64

// load reads a font file: TrueType, OpenType, a collection of either, WOFF
// or WOFF2. Its fonts are parsed; prepare reads what the views show of
// those that are shown.
func load(data []byte) (*file, error) {
	fl := &file{size: len(data), layout: map[tableKey]any{}}
	switch {
	case len(data) >= 4 && string(data[:4]) == "wOFF":
		d, err := unWOFF(data)
		if err != nil {
			return nil, err
		}
		fl.container, data = "WOFF", d
	case len(data) >= 4 && string(data[:4]) == "wOF2":
		d, err := woff2.Decode(data)
		if errors.Is(err, woff2.ErrDecodeNotAvailable) {
			return nil, errors.New("font: WOFF2 fonts cannot be read in this build (it has no Brotli decoder)")
		}
		if err != nil {
			return nil, fmt.Errorf("font: %w", err)
		}
		fl.container, data = "WOFF2", d
	}
	n := sfnt.NumFonts(data)
	if n == 0 {
		return nil, errors.New("font: not a font file")
	}
	fl.collection = len(data) >= 4 && string(data[:4]) == "ttcf"
	if n > maxFaces {
		n = maxFaces
	}
	for i := range n {
		f, err := sfnt.ParseIndex(data, i)
		if err == nil && f.NumGlyphs == 0 {
			err = errors.New("the font has no glyphs")
		}
		if err != nil {
			if n == 1 {
				return nil, fmt.Errorf("font: %w", err)
			}
			continue
		}
		fc := &face{fl: fl, index: i, f: f, size: len(data)}
		if fl.collection || fl.container != "" {
			// the size of the sfnt made of its tables (sfnt.Build)
			fc.size = 12 + 16*len(f.Tables)
			for _, t := range f.Tables {
				fc.size += (len(t) + 3) &^ 3
			}
		}
		fl.faces = append(fl.faces, fc)
	}
	if len(fl.faces) == 0 {
		return nil, errors.New("font: no font of the collection can be read")
	}
	return fl, nil
}

// table returns a layout table of the file read with read, which is called
// once for the fonts that share the table.
func table[T any](fl *file, tag string, b []byte, read func([]byte) T) T {
	if len(b) == 0 {
		return read(b)
	}
	k := tableKey{tag, &b[0], len(b)}
	if t, ok := fl.layout[k]; ok {
		return t.(T)
	}
	t := read(b)
	fl.layout[k] = t
	return t
}

// prepare reads what the views show from the tables (once).
func (fc *face) prepare() {
	if fc.ready {
		return
	}
	fc.ready = true
	f := fc.f
	fc.upem = float64(f.UnitsPerEm)
	fc.out = sfnt.NewOutlines(f)
	// a table that cannot be read is nil
	fc.gsub = table(fc.fl, "GSUB", f.Tables["GSUB"], func(b []byte) *otlayout.Table { t, _ := otlayout.ParseGSUB(b); return t })
	fc.gpos = table(fc.fl, "GPOS", f.Tables["GPOS"], func(b []byte) *otlayout.Table { t, _ := otlayout.ParseGPOS(b); return t })
	fc.gdef = table(fc.fl, "GDEF", f.Tables["GDEF"], otlayout.ParseGDEF)
	fc.names, fc.namesCut = parseNames(f.Tables["name"])
	fc.glyphs = glyphNames(f)
	n := f.NumGlyphs
	fc.chars = make([][]rune, n)
	for c, g := range f.Cmap {
		if int(g) < n && g != 0 && c <= 0x10FFFF {
			fc.chars[g] = append(fc.chars[g], rune(c))
			fc.runes = append(fc.runes, rune(c))
		}
	}
	for _, cs := range fc.chars {
		slices.Sort(cs)
	}
	slices.Sort(fc.runes)
	fc.advances = make([]float64, n)
	fc.bounds = make([]glyphBox, n)
	for g := range n {
		fc.advances[g] = float64(f.Advance(uint16(g))) / fc.upem
		if r, ok := fc.out.Bounds(uint16(g)); ok {
			fc.bounds[g] = glyphBox{r.XMin / fc.upem, r.YMin / fc.upem, r.XMax / fc.upem, r.YMax / fc.upem, true}
		}
	}
}

// glyph returns the glyph of a character, 0 when the font lacks it.
func (fc *face) glyph(r rune) uint16 {
	if g := fc.f.Cmap[uint32(r)]; int(g) < fc.f.NumGlyphs {
		return g
	}
	return 0
}

// valid reports whether glyphs a layout table names are in the font.
func (fc *face) valid(gs ...uint16) bool {
	for _, g := range gs {
		if int(g) >= fc.f.NumGlyphs {
			return false
		}
	}
	return true
}

// has reports whether the font maps every character of s that is not a
// space or a control character.
func (fc *face) has(s string) bool {
	for _, r := range s {
		if r > ' ' && fc.glyph(r) == 0 {
			return false
		}
	}
	return true
}

// unWOFF unpacks a WOFF 1.0 file into an sfnt.
func unWOFF(data []byte) ([]byte, error) {
	if len(data) < 44 {
		return nil, errors.New("font: WOFF file too short")
	}
	flavor := binary.BigEndian.Uint32(data[4:])
	if flavor == 0x74746366 {
		return nil, errors.New("font: WOFF font collections are not supported")
	}
	n := int(binary.BigEndian.Uint16(data[12:]))
	// tables may share their compressed bytes: what they decompress to is
	// bounded together, before any is decompressed
	total := 0
	for i := range n {
		if rec := 44 + i*20; rec+20 <= len(data) {
			total += int(binary.BigEndian.Uint32(data[rec+12:]))
		}
		if total > maxSize {
			return nil, fmt.Errorf("font: the tables of the WOFF file are larger than %d bytes", maxSize)
		}
	}
	tables := map[string][]byte{}
	for i := range n {
		rec := 44 + i*20
		if rec+20 > len(data) {
			return nil, errors.New("font: WOFF table directory truncated")
		}
		tag := string(data[rec : rec+4])
		off := int(binary.BigEndian.Uint32(data[rec+4:]))
		comp := int(binary.BigEndian.Uint32(data[rec+8:]))
		orig := int(binary.BigEndian.Uint32(data[rec+12:]))
		if off < 0 || comp < 0 || off+comp > len(data) || orig > 256<<20 {
			return nil, fmt.Errorf("font: WOFF table %q out of range", tag)
		}
		b := data[off : off+comp]
		if comp < orig {
			zr, err := zlib.NewReader(bytes.NewReader(b))
			if err != nil {
				return nil, fmt.Errorf("font: WOFF table %q: %w", tag, err)
			}
			u, err := io.ReadAll(io.LimitReader(zr, int64(orig)+1))
			if err != nil || len(u) != orig {
				return nil, fmt.Errorf("font: WOFF table %q does not decompress", tag)
			}
			b = u
		}
		tables[tag] = b
	}
	return sfnt.Build(tables, flavor == 0x4f54544f), nil
}
