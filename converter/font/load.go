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
}

// face is one font of a file.
type face struct {
	index int // in the collection
	// data is the font on its own: the file itself, or for a font of a
	// collection or a web font the sfnt made of its tables.
	data []byte
	f    *sfnt.Font
	out  *sfnt.Outlines

	gsub, gpos *otlayout.Table
	gdef       *otlayout.GDEF

	names    []nameRec  // the name table
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
// or WOFF2.
func load(data []byte) (*file, error) {
	fl := &file{size: len(data)}
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
		if err != nil {
			if n == 1 {
				return nil, fmt.Errorf("font: %w", err)
			}
			continue
		}
		fc := &face{index: i, f: f, data: data}
		if fl.collection || fl.container != "" {
			fc.data = sfnt.Build(f.Tables, f.IsCFF)
		}
		fc.prepare()
		fl.faces = append(fl.faces, fc)
	}
	if len(fl.faces) == 0 {
		return nil, errors.New("font: no font of the collection can be read")
	}
	return fl, nil
}

// prepare reads what the views show from the tables.
func (fc *face) prepare() {
	f := fc.f
	fc.upem = float64(f.UnitsPerEm)
	fc.out = sfnt.NewOutlines(f)
	if t, err := otlayout.ParseGSUB(f.Tables["GSUB"]); err == nil {
		fc.gsub = t
	}
	if t, err := otlayout.ParseGPOS(f.Tables["GPOS"]); err == nil {
		fc.gpos = t
	}
	fc.gdef = otlayout.ParseGDEF(f.Tables["GDEF"])
	fc.names = parseNames(f.Tables["name"])
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
