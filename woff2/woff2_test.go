package woff2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/andybalholm/brotli"
)

type point struct {
	x, y int
	on   bool
}

// glyph is the decoded form both glyf encodings are compared in.
type glyph struct {
	contours  [][]point
	composite []byte
	instr     []byte
	bbox      [4]int16
	overlap   bool
}

// parseGlyf decodes a plain glyf/loca pair.
func parseGlyf(t *testing.T, glyf, loca []byte, n int, long bool) []glyph {
	t.Helper()
	off := func(i int) int {
		if long {
			return int(binary.BigEndian.Uint32(loca[i*4:]))
		}
		return int(binary.BigEndian.Uint16(loca[i*2:])) * 2
	}
	out := make([]glyph, n)
	for g := 0; g < n; g++ {
		d := glyf[off(g):off(g+1)]
		if len(d) == 0 {
			continue
		}
		nc := int16(binary.BigEndian.Uint16(d))
		if nc == 0 {
			continue
		}
		gl := &out[g]
		for i := range gl.bbox {
			gl.bbox[i] = int16(binary.BigEndian.Uint16(d[2+i*2:]))
		}
		if nc < 0 {
			p := 10
			haveInstr := false
			for {
				f := binary.BigEndian.Uint16(d[p:])
				haveInstr = haveInstr || f&compHaveInstr != 0
				p += 4
				if f&compArgsAreWords != 0 {
					p += 4
				} else {
					p += 2
				}
				switch {
				case f&compHaveScale != 0:
					p += 2
				case f&compHaveXYScale != 0:
					p += 4
				case f&compHave2x2 != 0:
					p += 8
				}
				if f&compMore == 0 {
					break
				}
			}
			gl.composite = d[10:p]
			if haveInstr {
				n := int(binary.BigEndian.Uint16(d[p:]))
				gl.instr = d[p+2 : p+2+n]
			}
			continue
		}
		ends := make([]int, nc)
		for i := range ends {
			ends[i] = int(binary.BigEndian.Uint16(d[10+i*2:]))
		}
		p := 10 + 2*int(nc)
		il := int(binary.BigEndian.Uint16(d[p:]))
		gl.instr = d[p+2 : p+2+il]
		p += 2 + il
		np := ends[nc-1] + 1
		var flags []byte
		for len(flags) < np {
			f := d[p]
			p++
			flags = append(flags, f)
			if f&flagRepeat != 0 {
				for r := int(d[p]); r > 0; r-- {
					flags = append(flags, f)
				}
				p++
			}
		}
		gl.overlap = flags[0]&flagOverlap != 0
		read := func(short, same byte) []int {
			v, out := 0, make([]int, np)
			for i, f := range flags {
				switch {
				case f&short != 0:
					dv := int(d[p])
					p++
					if f&same == 0 {
						dv = -dv
					}
					v += dv
				case f&same == 0:
					v += int(int16(binary.BigEndian.Uint16(d[p:])))
					p += 2
				}
				out[i] = v
			}
			return out
		}
		xs := read(flagXShort, flagXSame)
		ys := read(flagYShort, flagYSame)
		start := 0
		for _, e := range ends {
			var c []point
			for i := start; i <= e; i++ {
				c = append(c, point{xs[i], ys[i], flags[i]&flagOnCurve != 0})
			}
			gl.contours = append(gl.contours, c)
			start = e + 1
		}
	}
	return out
}

// reader is a byte stream for the decoder below.
type reader struct {
	b []byte
	p int
}

func (r *reader) u8() int  { v := r.b[r.p]; r.p++; return int(v) }
func (r *reader) u16() int { v := binary.BigEndian.Uint16(r.b[r.p:]); r.p += 2; return int(v) }
func (r *reader) take(n int) []byte {
	v := r.b[r.p : r.p+n]
	r.p += n
	return v
}
func (r *reader) u255() int {
	switch c := r.u8(); c {
	case 253:
		return r.u16()
	case 254:
		return r.u8() + 506
	case 255:
		return r.u8() + 253
	default:
		return c
	}
}

// decodeGlyf decodes a transformed glyf table following the reference decoder.
func decodeGlyf(t *testing.T, b []byte) []glyph {
	t.Helper()
	hdr := &reader{b: b}
	hdr.u16()
	option := hdr.u16()
	n := hdr.u16()
	hdr.u16()
	var sizes [7]int
	for i := range sizes {
		sizes[i] = int(binary.BigEndian.Uint32(b[8+i*4:]))
	}
	streams := make([]*reader, 7)
	p := 36
	for i, s := range sizes {
		streams[i] = &reader{b: b[p : p+s]}
		p += s
	}
	nContours, nPoints, flagS, glyphS, compS, bboxS, instrS := streams[0], streams[1], streams[2], streams[3], streams[4], streams[5], streams[6]
	bitmap := bboxS.take(((n + 31) >> 5) << 2)
	var overlap []byte
	if option&1 != 0 {
		overlap = b[p : p+(n+7)>>3]
		p += (n + 7) >> 3
	}
	if p != len(b) {
		t.Fatalf("transformed glyf has %d trailing bytes", len(b)-p)
	}
	out := make([]glyph, n)
	for g := 0; g < n; g++ {
		nc := int16(nContours.u16())
		hasBBox := bitmap[g>>3]&(0x80>>(g&7)) != 0
		gl := &out[g]
		switch {
		case nc == 0:
			if hasBBox {
				t.Fatalf("glyph %d: empty glyph with bbox", g)
			}
			continue
		case nc < 0:
			if nc != -1 {
				t.Fatalf("glyph %d: nContours %d", g, nc)
			}
			start := compS.p
			haveInstr := false
			for {
				f := compS.u16()
				compS.u16()
				haveInstr = haveInstr || f&compHaveInstr != 0
				if f&compArgsAreWords != 0 {
					compS.take(4)
				} else {
					compS.take(2)
				}
				switch {
				case f&compHaveScale != 0:
					compS.take(2)
				case f&compHaveXYScale != 0:
					compS.take(4)
				case f&compHave2x2 != 0:
					compS.take(8)
				}
				if f&compMore == 0 {
					break
				}
			}
			gl.composite = compS.b[start:compS.p]
			if haveInstr {
				gl.instr = instrS.take(glyphS.u255())
			}
			if !hasBBox {
				t.Fatalf("glyph %d: composite without bbox", g)
			}
		default:
			x, y := 0, 0
			for c := 0; c < int(nc); c++ {
				np := nPoints.u255()
				var pts []point
				for i := 0; i < np; i++ {
					flag := flagS.u8()
					on := flag>>7 == 0
					flag &= 0x7f
					var dx, dy int
					sign := func(f, v int) int {
						if f&1 != 0 {
							return v
						}
						return -v
					}
					switch {
					case flag < 10:
						dy = sign(flag, (flag&14)<<7+glyphS.u8())
					case flag < 20:
						dx = sign(flag, ((flag-10)&14)<<7+glyphS.u8())
					case flag < 84:
						b0, b1 := flag-20, glyphS.u8()
						dx = sign(flag, 1+(b0&0x30)+b1>>4)
						dy = sign(flag>>1, 1+(b0&0x0c)<<2+b1&0x0f)
					case flag < 120:
						b0 := flag - 84
						dx = sign(flag, 1+(b0/12)<<8+glyphS.u8())
						dy = sign(flag>>1, 1+((b0%12)>>2)<<8+glyphS.u8())
					case flag < 124:
						b0, b1, b2 := glyphS.u8(), glyphS.u8(), glyphS.u8()
						dx = sign(flag, b0<<4+b1>>4)
						dy = sign(flag>>1, (b1&0x0f)<<8+b2)
					default:
						dx = sign(flag, glyphS.u16())
						dy = sign(flag>>1, glyphS.u16())
					}
					x += dx
					y += dy
					pts = append(pts, point{x, y, on})
				}
				gl.contours = append(gl.contours, pts)
			}
			gl.instr = instrS.take(glyphS.u255())
			gl.overlap = overlap != nil && overlap[g>>3]&(0x80>>(g&7)) != 0
		}
		if hasBBox {
			bb := bboxS.take(8)
			for i := range gl.bbox {
				gl.bbox[i] = int16(binary.BigEndian.Uint16(bb[i*2:]))
			}
		} else {
			for i, c := range gl.contours {
				for j, pt := range c {
					if i == 0 && j == 0 || pt.x < int(gl.bbox[0]) {
						gl.bbox[0] = int16(pt.x)
					}
					if i == 0 && j == 0 || pt.y < int(gl.bbox[1]) {
						gl.bbox[1] = int16(pt.y)
					}
					if i == 0 && j == 0 || pt.x > int(gl.bbox[2]) {
						gl.bbox[2] = int16(pt.x)
					}
					if i == 0 && j == 0 || pt.y > int(gl.bbox[3]) {
						gl.bbox[3] = int16(pt.y)
					}
				}
			}
		}
	}
	for i, s := range streams {
		if s.p != len(s.b) {
			t.Fatalf("stream %d: %d of %d bytes used", i, s.p, len(s.b))
		}
	}
	return out
}

// decode unpacks a WOFF2 file into its (possibly transformed) tables.
func decode(t *testing.T, w []byte) (flavor uint32, tables map[string][]byte, origLen map[string]int, transformed map[string]bool) {
	t.Helper()
	if string(w[:4]) != "wOF2" || int(binary.BigEndian.Uint32(w[8:])) != len(w) || len(w)%4 != 0 {
		t.Fatalf("bad header or length")
	}
	if int(binary.BigEndian.Uint32(w[16:])) < len(w) {
		t.Fatal("totalSfntSize is smaller than the file")
	}
	n := int(binary.BigEndian.Uint16(w[12:]))
	r := &reader{b: w, p: 48}
	base128 := func() int {
		v := 0
		for i := 0; i < 5; i++ {
			c := r.u8()
			if i == 0 && c == 0x80 {
				t.Fatal("base128 with leading zero")
			}
			v = v<<7 | c&0x7f
			if c&0x80 == 0 {
				return v
			}
		}
		t.Fatal("base128 too long")
		return 0
	}
	type entry struct {
		tag       string
		orig, len int
		xform     bool
	}
	var entries []entry
	for i := 0; i < n; i++ {
		f := r.u8()
		tag := ""
		if f&0x3f == 0x3f {
			tag = string(r.take(4))
		} else {
			tag = knownTags[f&0x3f]
		}
		e := entry{tag: tag, orig: base128()}
		v := f >> 6
		e.xform = v != 0
		if tag == "glyf" || tag == "loca" {
			e.xform = v != 3
		}
		e.len = e.orig
		if e.xform {
			e.len = base128()
		}
		if i > 0 && entries[i-1].tag >= tag {
			t.Fatalf("tables not sorted: %s after %s", tag, entries[i-1].tag)
		}
		entries = append(entries, e)
	}
	compressed := int(binary.BigEndian.Uint32(w[20:]))
	data, err := io.ReadAll(brotli.NewReader(bytes.NewReader(w[r.p : r.p+compressed])))
	if err != nil {
		t.Fatal(err)
	}
	tables, origLen, transformed = map[string][]byte{}, map[string]int{}, map[string]bool{}
	p := 0
	for _, e := range entries {
		tables[e.tag] = data[p : p+e.len]
		origLen[e.tag] = e.orig
		transformed[e.tag] = e.xform
		p += e.len
	}
	if p != len(data) {
		t.Fatalf("decompressed stream has %d extra bytes", len(data)-p)
	}
	return binary.BigEndian.Uint32(w[4:]), tables, origLen, transformed
}

func sfntTables(font []byte) map[string][]byte {
	out := map[string][]byte{}
	n := int(binary.BigEndian.Uint16(font[4:]))
	for i := 0; i < n; i++ {
		rec := 12 + i*16
		off, l := binary.BigEndian.Uint32(font[rec+8:]), binary.BigEndian.Uint32(font[rec+12:])
		out[string(font[rec:rec+4])] = font[off : off+l]
	}
	return out
}

func testFonts(t *testing.T) map[string][]byte {
	fonts := map[string][]byte{}
	for _, p := range []string{"../fixture/testdata/fonts/DejaVuSans-sub.ttf", "../fixture/testdata/fonts/DejaVuSans-Bold-sub.ttf", "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", "/usr/share/fonts/truetype/dejavu/DejaVuSerif-Bold.ttf"} {
		if b, err := os.ReadFile(p); err == nil {
			fonts[p] = b
		}
	}
	if len(fonts) == 0 {
		t.Skip("no test fonts")
	}
	return fonts
}

func TestEncodeTrueType(t *testing.T) {
	if !Available() {
		t.Skip("built without the Brotli encoder")
	}
	for name, font := range testFonts(t) {
		t.Run(name, func(t *testing.T) {
			w, err := Encode(font)
			if err != nil {
				t.Fatal(err)
			}
			if !IsWOFF(w) {
				t.Fatal("not recognized as WOFF")
			}
			flavor, tables, origLen, xform := decode(t, w)
			orig := sfntTables(font)
			if flavor != binary.BigEndian.Uint32(font) || len(tables) != len(orig) {
				t.Fatalf("flavor %x, %d tables (want %d)", flavor, len(tables), len(orig))
			}
			if !xform["glyf"] || !xform["loca"] || len(tables["loca"]) != 0 {
				t.Fatal("glyf/loca not transformed")
			}
			n := int(binary.BigEndian.Uint16(orig["maxp"][4:]))
			if origLen["loca"] != (n+1)*4 || binary.BigEndian.Uint16(tables["head"][50:]) != 1 {
				t.Fatal("loca must be rebuilt in the long format")
			}
			for tag, data := range orig {
				if tag != "glyf" && tag != "loca" && tag != "head" && !bytes.Equal(tables[tag], data) {
					t.Errorf("table %s changed", tag)
				}
			}
			want := parseGlyf(t, orig["glyf"], orig["loca"], n, binary.BigEndian.Uint16(orig["head"][50:]) != 0)
			got := decodeGlyf(t, tables["glyf"])
			composites := 0
			for g := range want {
				if want[g].composite != nil {
					composites++
				}
				if !reflect.DeepEqual(fmt.Sprint(want[g]), fmt.Sprint(got[g])) {
					t.Fatalf("glyph %d differs:\n want %v\n got  %v", g, want[g], got[g])
				}
			}
			t.Logf("%d glyphs (%d composite): %d bytes -> %d", n, composites, len(font), len(w))
			if len(w) >= len(font)*6/10 {
				t.Errorf("WOFF2 is %d bytes for a %d byte font", len(w), len(font))
			}
		})
	}
}

func TestEncodeFallsBackForUnusualGlyf(t *testing.T) {
	if !Available() {
		t.Skip("built without the Brotli encoder")
	}
	font := testFonts(t)["../fixture/testdata/fonts/DejaVuSans-sub.ttf"]
	if font == nil {
		t.Skip("fixture font missing")
	}
	// Truncate loca: the transform must give up and store both tables as they are.
	tables := sfntTables(font)
	loca := tables["loca"][:len(tables["loca"])-4]
	var b bytes.Buffer
	b.Write(font[:12])
	var body []byte
	n := int(binary.BigEndian.Uint16(font[4:]))
	off := 12 + 16*n
	for i := 0; i < n; i++ {
		rec := font[12+i*16 : 28+i*16]
		tag := string(rec[:4])
		data := tables[tag]
		if tag == "loca" {
			data = loca
		}
		var e [16]byte
		copy(e[:], rec[:8])
		binary.BigEndian.PutUint32(e[8:], uint32(off+len(body)))
		binary.BigEndian.PutUint32(e[12:], uint32(len(data)))
		b.Write(e[:])
		body = append(body, data...)
		for len(body)%4 != 0 {
			body = append(body, 0)
		}
	}
	b.Write(body)
	w, err := Encode(b.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	_, got, _, xform := decode(t, w)
	if xform["glyf"] || xform["loca"] || !bytes.Equal(got["glyf"], tables["glyf"]) || !bytes.Equal(got["loca"], loca) {
		t.Fatal("expected untransformed glyf and loca")
	}
}

func TestBase128And255(t *testing.T) {
	for v, want := range map[uint32]string{0: "00", 127: "7f", 128: "8100", 16384: "818000", 0xffffffff: "8fffffff7f"} {
		if got := fmt.Sprintf("%x", appendBase128(nil, v)); got != want {
			t.Errorf("base128(%d) = %s, want %s", v, got, want)
		}
	}
	for _, v := range []int{0, 252, 253, 505, 506, 761, 762, 65535} {
		r := &reader{b: append255(nil, v)}
		if got := r.u255(); got != v || r.p != len(r.b) {
			t.Errorf("255UInt16 %d round-trips to %d", v, got)
		}
	}
}
