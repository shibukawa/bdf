package pptx

import (
	"math"
	"sort"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/woff2"
)

// canvas is an Object under construction. Font references are placeholders
// until the whole deck has been laid out, because embedded font programs are
// subsets of every character used anywhere; child objects (vertical text
// columns) are encoded before their parent.
type canvas struct {
	c        *converter
	obj      *bdf.Object
	fonts    []fontUse
	fontIdx  map[fontUse]bdf.FontRef
	images   map[bdf.Hash]bdf.ImageRef
	children []*canvas
	drawn    bool

	encoded bool
	hash    bdf.Hash
	bbox    bdf.Rect
}

// fontUse is what a FONT instruction asks for: a resolved face (nil when no
// font file is available) and how the request looked, for the system-font
// fallback and synthetic styles.
type fontUse struct {
	face        *fontdb.Face
	requested   string
	generic     string
	bold        bool
	italic      bool
	synthBold   bool
	synthItalic bool
}

func (c *converter) newCanvas() *canvas {
	cv := &canvas{c: c, obj: bdf.NewObject(), fontIdx: map[fontUse]bdf.FontRef{}, images: map[bdf.Hash]bdf.ImageRef{}}
	c.canvases = append(c.canvases, cv)
	return cv
}

func (cv *canvas) font(u fontUse) bdf.FontRef {
	if r, ok := cv.fontIdx[u]; ok {
		return r
	}
	r := cv.obj.AddFont(bdf.SystemFont("\x00pending", 0, 0))
	cv.fonts = append(cv.fonts, u)
	cv.fontIdx[u] = r
	return r
}

func (cv *canvas) image(h bdf.Hash) bdf.ImageRef {
	if r, ok := cv.images[h]; ok {
		return r
	}
	r := cv.obj.AddImage(h)
	cv.images[h] = r
	return r
}

// child adds a child object and returns its canvas and reference.
func (cv *canvas) child(bbox bdf.Rect) (*canvas, bdf.ObjRef) {
	ch := cv.c.newCanvas()
	ref := cv.obj.AddObject(bdf.Hash{}, bbox)
	cv.children = append(cv.children, ch)
	return ch, ref
}

func f32(v float64) float32 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return float32(v)
}

// matrix is an affine transform [a b c d e f] (Canvas order).
type matrix [6]float64

var identity = matrix{1, 0, 0, 1, 0, 0}

func (m matrix) mul(n matrix) matrix {
	return matrix{
		m[0]*n[0] + m[2]*n[1], m[1]*n[0] + m[3]*n[1],
		m[0]*n[2] + m[2]*n[3], m[1]*n[2] + m[3]*n[3],
		m[0]*n[4] + m[2]*n[5] + m[4], m[1]*n[4] + m[3]*n[5] + m[5],
	}
}

func (m matrix) apply(x, y float64) (float64, float64) {
	return m[0]*x + m[2]*y + m[4], m[1]*x + m[3]*y + m[5]
}

func translate(x, y float64) matrix { return matrix{1, 0, 0, 1, x, y} }
func scale(x, y float64) matrix     { return matrix{x, 0, 0, y, 0, 0} }
func rotate(deg float64) matrix {
	s, c := math.Sincos(deg * math.Pi / 180)
	return matrix{c, s, -s, c, 0, 0}
}

// transform emits m (as TRANSLATE when it is one).
func (cv *canvas) transform(m matrix) {
	if m[0] == 1 && m[1] == 0 && m[2] == 0 && m[3] == 1 {
		if m[4] != 0 || m[5] != 0 {
			cv.obj.Translate(f32(m[4]), f32(m[5]))
		}
		return
	}
	cv.obj.Transform(f32(m[0]), f32(m[1]), f32(m[2]), f32(m[3]), f32(m[4]), f32(m[5]))
}

// finalize builds the embedded fonts and encodes every canvas.
func (c *converter) finalize() {
	type faceKey = *fontdb.Face
	hashes := map[faceKey]bdf.Hash{}
	embedded := map[faceKey]bool{}
	if !c.opts.SystemFonts {
		faces := make([]faceKey, 0, len(c.faceRunes))
		for f := range c.faceRunes {
			faces = append(faces, f)
		}
		sort.Slice(faces, func(i, j int) bool {
			if faces[i].Path != faces[j].Path {
				return faces[i].Path < faces[j].Path
			}
			return faces[i].Index < faces[j].Index
		})
		for _, f := range faces {
			l, err := f.Load()
			if err != nil {
				continue
			}
			if !l.CanSubset(c.opts.IgnoreFSType) && !c.opts.NoSubset && l.Size() > maxWholeFont {
				c.warnf("font %s is not embedded: its outlines cannot be subset and the file is %d KB", f.Family, l.Size()/1024)
				continue
			}
			runes := make([]rune, 0, len(c.faceRunes[f]))
			for r := range c.faceRunes[f] {
				runes = append(runes, r)
			}
			data, ok := l.Program(runes, c.opts.NoSubset, c.opts.IgnoreFSType)
			if !ok {
				c.warnf("font %s is not embedded: its license does not allow embedding", f.Family)
				continue
			}
			if !c.opts.NoWOFF2 {
				if w, err := woff2.Encode(data); err == nil {
					data = w
				} else if err != woff2.ErrNotAvailable {
					c.warnf("font %s: not stored as WOFF2: %v", f.Family, err)
				}
			}
			hashes[f] = c.doc.AddFont(data)
			embedded[f] = true
		}
		c.embeddedFonts = len(hashes)
	}
	var encode func(cv *canvas) bdf.Hash
	encode = func(cv *canvas) bdf.Hash {
		if cv.encoded {
			return cv.hash
		}
		cv.encoded = true
		for i, u := range cv.fonts {
			cv.obj.UpdateFont(bdf.FontRef(i), fontRecord(u, hashes, embedded))
		}
		for i, ch := range cv.children {
			cv.obj.UpdateObject(bdf.ObjRef(i), encode(ch))
		}
		cv.hash, cv.bbox = c.doc.AddObject(cv.obj)
		return cv.hash
	}
	for _, cv := range c.canvases {
		encode(cv)
	}
}

// maxWholeFont is the largest font file embedded whole when it cannot be
// subset (CFF outlines).
const maxWholeFont = 2 << 20

func fontRecord(u fontUse, hashes map[*fontdb.Face]bdf.Hash, embedded map[*fontdb.Face]bool) bdf.Font {
	if u.face != nil && embedded[u.face] {
		// The face registered in the browser is the one that was measured;
		// only styles it lacks are synthesized.
		f := bdf.EmbeddedFont(hashes[u.face], 400, bdf.StyleNormal)
		if u.synthBold {
			f.Weight = 700
		}
		if u.synthItalic {
			f.Style = bdf.StyleItalic
		}
		f.Family = u.generic
		return f
	}
	family := cssQuote(u.requested)
	if u.face != nil && u.face.Family != u.requested {
		family += ", " + cssQuote(u.face.Family)
	}
	if family != "" {
		family += ", "
	}
	family += u.generic
	var w uint16 = 400
	if u.bold {
		w = 700
	}
	st := bdf.StyleNormal
	if u.italic {
		st = bdf.StyleItalic
	}
	return bdf.SystemFont(family, w, st)
}

func cssQuote(s string) string {
	if s == "" {
		return ""
	}
	out := []rune{'"'}
	for _, r := range s {
		if r == '"' || r == '\\' {
			out = append(out, '\\')
		}
		out = append(out, r)
	}
	return string(append(out, '"'))
}

// pathB builds a bdf.Path from float64 coordinates.
type pathB struct{ p *bdf.Path }

func newPath() *pathB { return &pathB{&bdf.Path{}} }

func (b *pathB) moveTo(x, y float64) *pathB { b.p.MoveTo(f32(x), f32(y)); return b }
func (b *pathB) lineTo(x, y float64) *pathB { b.p.LineTo(f32(x), f32(y)); return b }
func (b *pathB) close() *pathB              { b.p.Close(); return b }
