package sfnt

import (
	"encoding/binary"
	"math"
	"os"
	"testing"
)

// pointPen collects every point a glyph outline passes, control points
// included, and counts contours.
type pointPen struct {
	r        Rect
	contours int
	open     bool
	bad      bool
}

func newPointPen() *pointPen {
	return &pointPen{r: Rect{math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)}}
}

func (p *pointPen) pt(x, y float64) {
	if !p.open {
		p.bad = true
	}
	p.r.XMin, p.r.XMax = math.Min(p.r.XMin, x), math.Max(p.r.XMax, x)
	p.r.YMin, p.r.YMax = math.Min(p.r.YMin, y), math.Max(p.r.YMax, y)
}

func (p *pointPen) MoveTo(x, y float64) {
	if p.open {
		p.bad = true
	}
	p.open = true
	p.contours++
	p.pt(x, y)
}
func (p *pointPen) LineTo(x, y float64)         { p.pt(x, y) }
func (p *pointPen) QuadTo(cx, cy, x, y float64) { p.pt(cx, cy); p.pt(x, y) }
func (p *pointPen) CubeTo(ax, ay, bx, by, x, y float64) {
	p.pt(ax, ay)
	p.pt(bx, by)
	p.pt(x, y)
}
func (p *pointPen) Close() {
	if !p.open {
		p.bad = true
	}
	p.open = false
}

func TestTrueTypeOutlines(t *testing.T) {
	data, err := os.ReadFile("../../fixture/testdata/fonts/DejaVuSans-sub.ttf")
	if err != nil {
		t.Skip(err)
	}
	f, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	o := NewOutlines(f)
	drawn := 0
	for gid := 0; gid < f.NumGlyphs; gid++ {
		p := newPointPen()
		if !o.Outline(uint16(gid), p) {
			t.Fatalf("glyph %d: not read", gid)
		}
		if p.bad || p.open {
			t.Fatalf("glyph %d: contours not opened and closed in pairs", gid)
		}
		want, ok := o.Bounds(uint16(gid))
		if !ok || p.contours == 0 {
			continue
		}
		drawn++
		// the glyph header holds the bounds of the points, control points included
		if p.r != want {
			t.Errorf("glyph %d: points span %v, header says %v", gid, p.r, want)
		}
	}
	if drawn == 0 {
		t.Fatal("no glyph drew anything")
	}
}

func TestCFFOutlines(t *testing.T) {
	data, err := os.ReadFile("testdata/STIXTwoMath-cff-subset.otf")
	if err != nil {
		t.Skip(err)
	}
	f, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	o := NewOutlines(f)
	drawn := 0
	for gid := 0; gid < f.NumGlyphs; gid++ {
		p := newPointPen()
		if !o.Outline(uint16(gid), p) {
			t.Fatalf("glyph %d: not read", gid)
		}
		if p.bad || p.open {
			t.Fatalf("glyph %d: contours not opened and closed in pairs", gid)
		}
		want, ok := o.Bounds(uint16(gid))
		if !ok {
			continue
		}
		drawn++
		// the control points enclose the curves
		const eps = 1e-6
		if p.r.XMin > want.XMin+eps || p.r.YMin > want.YMin+eps || p.r.XMax < want.XMax-eps || p.r.YMax < want.YMax-eps {
			t.Errorf("glyph %d: points span %v, outline bounds %v", gid, p.r, want)
		}
	}
	if drawn == 0 {
		t.Fatal("no glyph drew anything")
	}
}

// TestCompositesBounded draws a glyph of four components that are glyphs
// of four components, nine levels deep: 262144 triangles from 600 bytes.
func TestCompositesBounded(t *testing.T) {
	const levels, fan = 9, 4
	var glyf, loca []byte
	add := func(g []byte) {
		loca = binary.BigEndian.AppendUint32(loca, uint32(len(glyf)))
		glyf = append(glyf, g...)
		for len(glyf)%4 != 0 {
			glyf = append(glyf, 0)
		}
	}
	for l := range levels {
		g := words(0xffff, 0, 0, 100, 100)
		for k := range fan {
			flags := argsAreXY
			if k < fan-1 {
				flags |= moreComponents
			}
			g = append(append(g, words(flags, l+1)...), 0, 0)
		}
		add(g)
	}
	triangle := append(words(1, 0, 0, 100, 100, 2, 0), onCurve, onCurve, onCurve)
	add(append(triangle, words(0, 100, 0, 0, 0, 100)...))
	loca = binary.BigEndian.AppendUint32(loca, uint32(len(glyf)))
	head := make([]byte, 54)
	binary.BigEndian.PutUint16(head[50:], 1)
	maxp := words(1, 0, levels+1)
	f, err := Parse(Build(map[string][]byte{"glyf": glyf, "loca": loca, "head": head, "maxp": maxp}, false))
	if err != nil {
		t.Fatal(err)
	}
	o := NewOutlines(f)
	p := newPointPen()
	if o.Outline(0, p) {
		t.Error("the glyph of 262144 components is drawn")
	}
	if p.contours > maxOutline {
		t.Errorf("%d contours drawn", p.contours)
	}
	// the glyphs of fewer components are drawn
	p = newPointPen()
	if !o.Outline(3, p) || p.contours != fan*fan*fan*fan*fan*fan {
		t.Errorf("six levels: %d contours", p.contours)
	}
}

// TestCFFWithoutFontDICT reads a CID-keyed font whose FDArray is empty:
// its glyphs have no Private DICT to run with.
func TestCFFWithoutFontDICT(t *testing.T) {
	number := func(v int) []byte { return []byte{29, byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)} }
	index := func(item []byte) []byte { return append([]byte{0, 1, 1, 1, byte(1 + len(item))}, item...) }
	header, name, none := []byte{1, 0, 4, 1}, index([]byte("A")), []byte{0, 0}
	charStrings := index([]byte{14})
	top := func(at int) []byte {
		d := append(number(at), 17)                                   // CharStrings
		d = append(append(d, number(at+len(charStrings))...), 12, 36) // FDArray, empty
		return append(append(d, number(at+len(charStrings)+2)...), 12, 37)
	}
	at := len(header) + len(name) + len(index(top(0))) + 2*len(none)
	var cff []byte
	for _, b := range [][]byte{header, name, index(top(at)), none, none, charStrings, none, none} {
		cff = append(cff, b...)
	}
	f, err := Parse(Build(map[string][]byte{"CFF ": cff, "head": make([]byte, 54), "maxp": {0, 0, 0x50, 0, 0, 1}}, true))
	if err != nil {
		t.Fatal(err)
	}
	o := NewOutlines(f)
	if _, ok := o.Bounds(0); ok {
		t.Error("the glyph has bounds")
	}
	if o.Outline(0, newPointPen()) {
		t.Error("the glyph is drawn")
	}
}
