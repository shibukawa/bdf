// Command gen-newstroke writes the stroke font the KiCad converter draws
// text with (converter/kicad/newstroke.txt.gz): NewStroke by Vladimir
// Uryvaev, the plotter font KiCad was designed with, as its author released
// it under CC0. It reads the font's sources (glyph libraries in KiCad's
// legacy symbol library format and the list that maps Unicode code points to
// glyphs, accents and composites) and compiles them the way the font's own
// script (fontconv.awk) does, into one line per code point from U+0020.
//
// A line holds the glyph's left and right bearing followed by its strokes,
// each point two characters (the value plus 'R'), strokes separated by
// " R": the encoding of the Hershey fonts that KiCad reads. Only the
// author's release is used: the newstroke_font.cpp that KiCad ships carries
// KiCad's GPL notice and CJK glyphs made later from Source Han Sans, and the
// converter draws CJK text with TrueType fonts instead.
//
// Usage (from the repository root):
//
//	go run ./tools/gen-newstroke [-src URL-or-dir] [-out dir]
//
// -src is the release archive (newstroke-font.tgz) or a directory holding
// its font.lib, symbol.lib and charlist.txt; by default the archive on the
// font's home page. The release's README, which states the license, is
// written next to the table as NEWSTROKE.txt.
package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// release is the archive the table was made from.
const release = "https://vovanium.ru/_media/sledy/newstroke/newstroke-font.tgz"

func main() {
	src := flag.String("src", release, "newstroke-font.tgz (URL or file) or a directory with font.lib, symbol.lib and charlist.txt")
	out := flag.String("out", "converter/kicad", "directory to write newstroke.txt.gz to")
	flag.Parse()
	files, err := load(*src, "font.lib", "symbol.lib", "charlist.txt", "README.txt")
	if err != nil {
		log.Fatal(err)
	}
	c := newCompiler()
	// the order of fontconv.awk's command line: later definitions win
	for _, lib := range []string{"symbol.lib", "font.lib"} {
		c.readLib(files[lib])
	}
	lines, err := c.compile(files["charlist.txt"])
	if err != nil {
		log.Fatal(err)
	}
	for _, w := range c.warnings {
		log.Print(w)
	}
	var buf bytes.Buffer
	zw, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	for _, l := range lines {
		zw.Write([]byte(l))
		zw.Write([]byte{'\n'})
	}
	zw.Close()
	if err := os.WriteFile(filepath.Join(*out, "newstroke.txt.gz"), buf.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
	readme := "The stroke font in newstroke.txt.gz is compiled by tools/gen-newstroke from\n" +
		"the NewStroke release " + release + "\n(font.lib, symbol.lib and charlist.txt). Its README follows.\n\n" +
		strings.ReplaceAll(string(files["README.txt"]), "\r\n", "\n")
	if err := os.WriteFile(filepath.Join(*out, "NEWSTROKE.txt"), []byte(readme), 0o644); err != nil {
		log.Fatal(err)
	}
	log.Printf("%d glyphs (U+0020–U+%04X), %d bytes", len(lines), 0x20+len(lines)-1, buf.Len())
}

// load reads the named files of the release: from a directory, or from the
// archive at a URL or path.
func load(src string, names ...string) (map[string][]byte, error) {
	out := map[string][]byte{}
	if st, err := os.Stat(src); err == nil && st.IsDir() {
		for _, n := range names {
			b, err := os.ReadFile(filepath.Join(src, n))
			if err != nil {
				return nil, err
			}
			out[n] = b
		}
		return out, nil
	}
	var r io.Reader
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		resp, err := http.Get(src)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s: %s", src, resp.Status)
		}
		r = resp.Body
	} else {
		f, err := os.Open(src)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
	}
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", src, err)
	}
	tr := tar.NewReader(zr)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		n := filepath.Base(h.Name)
		for _, want := range names {
			if n == want {
				b, err := io.ReadAll(tr)
				if err != nil {
					return nil, err
				}
				out[n] = b
			}
		}
	}
	for _, n := range names {
		if _, ok := out[n]; !ok {
			return nil, fmt.Errorf("%s: no %s", src, n)
		}
	}
	return out, nil
}

// The constants of fontconv.awk.
const (
	code0     = 'R'
	libScale  = 50 // library units per font unit
	base      = 9  // added to y
	symdef    = "DEL"
	capHeight = -21
	xHeight   = -14
	symHeight = -16
	supOffset = -13
	subOffset = 6
)

// transform is the mirroring and offset a glyph name's first character
// asks for.
type transform struct{ sx, sy, oy int }

var transforms = map[byte]transform{
	'!': {-1, +1, 0},         // reversed
	'-': {+1, -1, xHeight},   // small, upside down
	'=': {+1, -1, capHeight}, // capital, upside down
	'~': {+1, -1, symHeight}, // symbol, upside down
	'+': {-1, -1, xHeight},   // small, turned
	'%': {-1, -1, capHeight}, // capital, turned
	'*': {-1, -1, symHeight}, // symbol, turned
	'^': {+1, +1, supOffset}, // superscript
	'`': {-1, +1, supOffset}, // superscript, reversed
	'.': {+1, +1, subOffset}, // subscript
	',': {-1, +1, subOffset}, // subscript, reversed
}

type point struct{ x, y float64 }

// glyph is a glyph of the libraries.
type glyph struct {
	strokes [][]point
	// left and right bearing (the "~" pins, or the "P" and "S" pins)
	ml, mr float64
	// anchors by pin name, where accents attach
	anchors map[string]point
}

type compiler struct {
	glyphs   map[string]*glyph
	warnings []string
	// ofx is the horizontal offset of the accent placed last (the awk
	// script's global of that name, which "+w" and "+p" lines read)
	ofx float64
}

func newCompiler() *compiler { return &compiler{glyphs: map[string]*glyph{}} }

// readLib reads the glyphs of a library: each DEF is a glyph, its P
// records its strokes and its X records (pins) its bearings and anchors.
func (c *compiler) readLib(b []byte) {
	var g *glyph
	sc := bufio.NewScanner(bytes.NewReader(b))
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for sc.Scan() {
		f := strings.Fields(sc.Text())
		if len(f) == 0 {
			continue
		}
		switch f[0] {
		case "DEF":
			if len(f) < 2 {
				continue
			}
			g = &glyph{anchors: map[string]point{"-": {}}}
			c.glyphs[f[1]] = g
		case "P":
			if g == nil || len(f) < 2 {
				continue
			}
			n, _ := strconv.Atoi(f[1])
			var pts []point
			o := point{-100, -100}
			for i := 5; i+1 < len(f) && i < 5+2*n; i += 2 {
				x := num(f[i]) / libScale
				y := -num(f[i+1]) / libScale
				if (point{x, y}) != o {
					pts = append(pts, point{x, y})
					o = point{x, y}
				}
			}
			g.strokes = append(g.strokes, pts)
		case "X":
			if g == nil || len(f) < 5 {
				continue
			}
			name, x, y := f[1], num(f[3])/libScale, -num(f[4])/libScale
			g.anchors[name] = point{x, y}
			switch {
			case x > 0 && name == "~" || name == "S":
				g.mr = x
			case x <= 0 && name == "~" || name == "P":
				g.ml = x
			}
		case "ENDDEF":
			g = nil
		}
	}
}

func num(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

// split separates a glyph reference into its name and transform.
func split(ref string) (string, transform) {
	if ref != "" {
		if t, ok := transforms[ref[0]]; ok {
			return ref[1:], t
		}
	}
	return ref, transform{1, 1, 0}
}

// lookup resolves a glyph reference, falling back to subst.
func (c *compiler) lookup(ref, subst string) (string, *glyph, transform) {
	name, t := split(ref)
	if g, ok := c.glyphs[name]; ok {
		return name, g, t
	}
	if subst == "" {
		subst = symdef
	}
	name, t = split(subst)
	return name, c.glyphs[name], t
}

// ch encodes a coordinate as the awk script's sprintf("%c", x+82) does:
// the integer part.
func ch(v float64) byte { return byte(int(v + code0)) }

// graph encodes the strokes of a glyph moved by (ofx, ofy).
func (c *compiler) graph(ref, subst string, ofx, ofy float64) string {
	name, _ := split(ref)
	if _, ok := c.glyphs[name]; !ok {
		c.warnings = append(c.warnings, "glyph "+ref+" not found")
	}
	_, g, t := c.lookup(ref, subst)
	if g == nil {
		return ""
	}
	ofy += float64(t.oy) + base
	var sb strings.Builder
	var last [2]byte
	for i, s := range g.strokes {
		if i > 0 {
			sb.WriteString(" R")
		}
		for _, p := range s {
			cur := [2]byte{ch(p.x*float64(t.sx) + ofx), ch(p.y*float64(t.sy) + ofy)}
			if cur != last {
				sb.Write(cur[:])
				last = cur
			}
		}
	}
	return sb.String()
}

// metric encodes the bearings of a glyph.
func (c *compiler) metric(ref string) string {
	_, g, t := c.lookup(ref, "")
	ml, mr := g.ml, g.mr
	if t.sx < 0 {
		ml, mr = mr, ml
	}
	return string([]byte{ch(float64(t.sx) * ml), ch(float64(t.sx) * mr)})
}

// metric2 encodes the left bearing of one glyph and the right bearing of
// another, each moved.
func (c *compiler) metric2(ref1, ref2 string, ofx1, ofx2 float64) string {
	_, g1, t1 := c.lookup(ref1, "")
	_, g2, t2 := c.lookup(ref2, "")
	ml := g1.ml
	if t1.sx < 0 {
		ml = g1.mr
	}
	mr := g2.mr
	if t2.sx < 0 {
		mr = g2.ml
	}
	return string([]byte{ch(float64(t1.sx)*ml + ofx1), ch(float64(t2.sx)*mr + ofx2)})
}

// dist is the advance from one glyph of a composite to the next.
func (c *compiler) dist(ref1, ref2 string) float64 {
	_, g1, t1 := c.lookup(ref1, "")
	_, g2, t2 := c.lookup(ref2, "")
	return float64(t1.sx)*g1.mr - float64(t2.sx)*g2.ml
}

// anchorPair returns the anchors an accent's anchor spec names on the base
// glyph and on the accent: "NAME" names the same pin on both, "A=B" pin A
// of the base and pin B of the accent.
func anchorPair(base, accent, spec string) (string, string) {
	if a, b, ok := strings.Cut(spec, "="); ok {
		return base + " " + a, accent + " " + b
	}
	return base + " " + spec, accent + " " + spec
}

func (c *compiler) anchor(key string) (point, bool) {
	name, pin, _ := strings.Cut(key, " ")
	g, ok := c.glyphs[name]
	if !ok {
		return point{}, false
	}
	p, ok := g.anchors[pin]
	return p, ok
}

// offset returns where an accent goes on a base glyph.
func (c *compiler) offset(baseRef, accentRef, spec string) (float64, float64) {
	if spec == "" || spec == "#" {
		return 0, 0
	}
	bn, bt := split(baseRef)
	an, at := split(accentRef)
	bk, ak := anchorPair(bn, an, spec)
	bp, okb := c.anchor(bk)
	ap, oka := c.anchor(ak)
	if !okb || !oka {
		c.warnings = append(c.warnings, "anchor "+bk+" / "+ak+" not found")
		return 0, 0
	}
	ox := float64(bt.sx)*bp.x - float64(at.sx)*ap.x
	oy := float64(bt.sy)*bp.y + float64(bt.oy) - float64(at.sy)*ap.y - float64(at.oy)
	return ox, oy
}

// accent encodes an accent placed on a base glyph.
func (c *compiler) accent(baseRef, accentRef, spec string, cx float64) string {
	if spec == "#" {
		spec = ""
	}
	ox, oy := c.offset(baseRef, accentRef, spec)
	c.ofx = ox
	return " R" + c.graph(accentRef, "", cx+ox, oy)
}

// compose encodes a line's glyph with up to three accents: fields are the
// glyph, then accent and anchor pairs.
func (c *compiler) compose(f []string, cx float64) (string, float64) {
	at := func(i int) string {
		if i < len(f) {
			return f[i]
		}
		return ""
	}
	s := c.graph(at(1), symdef, cx, 0)
	ofx2 := c.ofx
	if at(2) != "" && at(2) != "#" {
		s += c.accent(at(1), at(2), at(3), cx)
		ofx2 = c.ofx
		if at(3) != "#" && at(4) != "" && at(4) != "#" {
			s += c.accent(at(1), at(4), at(5), cx)
			if at(5) != "#" && at(6) != "" && at(6) != "#" {
				s += c.accent(at(1), at(6), at(7), cx)
			}
		}
	}
	return s, ofx2
}

// compile reads the character list and returns the glyph of each code
// point from the start.
func (c *compiler) compile(list []byte) ([]string, error) {
	var out []string
	code := -1
	var comx float64
	var str, first, prev string
	emit := func(s string) {
		if code < 0 {
			return
		}
		out = append(out, s)
		code++
	}
	sc := bufio.NewScanner(bytes.NewReader(list))
	for sc.Scan() {
		f := strings.Fields(sc.Text())
		if len(f) == 0 {
			continue
		}
		switch f[0] {
		case "startchar":
			n, err := strconv.Atoi(f[1])
			if err != nil || n != 0x20 {
				return nil, fmt.Errorf("charlist: the table starts at U+0020, not %q", f[1])
			}
			code = n
		case "skipcodes":
			n, _ := strconv.Atoi(f[1])
			for range n {
				emit(c.metric(symdef) + c.graph(symdef, "", 0, 0))
			}
		case "+", "+w", "+p":
			s, ofx2 := c.compose(f, 0)
			var met string
			switch f[0] {
			case "+w":
				met = c.metric2(f[1], field(f, 2), 0, ofx2)
			case "+p":
				met = c.metric2(field(f, 2), f[1], ofx2, 0)
			default:
				met = c.metric(f[1])
			}
			emit(met + s)
		case "+(":
			comx = 0
			str, _ = c.compose(f, comx)
			first, prev = f[1], f[1]
		case "+|":
			comx += c.dist(prev, f[1])
			s, _ := c.compose(f, comx)
			str += " R" + s
			prev = f[1]
		case "+)":
			comx += c.dist(prev, f[1])
			s, _ := c.compose(f, comx)
			str += " R" + s
			emit(c.metric2(first, f[1], 0, comx) + str)
		}
	}
	return out, sc.Err()
}

func field(f []string, i int) string {
	if i < len(f) {
		return f[i]
	}
	return ""
}
