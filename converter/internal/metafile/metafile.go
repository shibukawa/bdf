// Package metafile replays Windows metafiles (EMF and WMF), the vector
// pictures Office documents embed, into BDF objects. Browsers cannot decode
// metafiles as images, so a GDI state machine (map modes, world transform,
// pens, brushes, fonts, clipping, paths) turns the drawing records into BDF
// paths, text and images. EMF+ records embedded in comments are skipped;
// Office writes the equivalent EMF records next to them.
package metafile

import (
	"image"
	"image/color"
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/imgconv"
)

// Recolor changes the colors of a picture (DrawingML picture effects, for
// example): Color those of its records, Image the pixels of its bitmaps.
type Recolor interface {
	Color(c color.NRGBA) color.NRGBA
	Image(img *image.NRGBA)
}

// Options are what drawing a metafile needs from the document.
type Options struct {
	// Doc stores the bitmaps of the metafile.
	Doc *bdf.Document
	// Fonts picks and measures the fonts of its text.
	Fonts *fontset.Set
	// Images controls how its bitmaps are encoded.
	Images imgconv.Options
	// Recolor, when set, changes the colors of the picture.
	Recolor Recolor
	// Warn receives non-fatal problems, each at most once per metafile.
	Warn func(msg string)
}

// Kind reports "emf", "wmf" or "" for data that is no metafile.
func Kind(data []byte) string {
	switch {
	case len(data) >= 44 && le32(data, 0) == 1 && string(data[40:44]) == " EMF":
		return "emf"
	case len(data) >= 22 && le32(data, 0) == 0x9AC6CDD7,
		len(data) >= 18 && (le16(data, 0) == 1 || le16(data, 0) == 2) && le16(data, 2) == 9:
		return "wmf"
	}
	return ""
}

// Draw replays a metafile into the rectangle x,y,w,h of cv, clipped to it.
// Malformed records end the replay with a warning.
func Draw(cv *canvas.Canvas, data []byte, x, y, w, h float64, opts *Options) {
	g := newGDI(cv, opts)
	g.obj.Save()
	g.obj.ClipRect(f32(x), f32(y), f32(w), f32(h))
	defer func() {
		if r := recover(); r != nil {
			g.warn("a metafile could not be replayed completely")
		}
		g.end()
		g.obj.Restore()
	}()
	switch Kind(data) {
	case "emf":
		g.playEMF(data, x, y, w, h)
	case "wmf":
		g.playWMF(data, x, y, w, h)
	}
}

var (
	black = color.NRGBA{0, 0, 0, 255}
	white = color.NRGBA{255, 255, 255, 255}
)

func gray(v uint8) color.NRGBA { return color.NRGBA{v, v, v, 255} }

func toBDF(c color.NRGBA) bdf.Color { return bdf.RGBA(c.R, c.G, c.B, c.A) }

type gdiPen struct {
	null     bool
	width    float64 // logical units; 0 = hairline
	cosmetic bool
	color    color.NRGBA
	style    uint32
}

type gdiBrush struct {
	null  bool
	color color.NRGBA
	hatch int // -1 = solid
}

type gdiFont struct {
	height     float64
	escapement float64 // tenths of a degree
	weight     int
	italic     bool
	underline  bool
	strike     bool
	face       string
}

type gdiState struct {
	pen       *gdiPen
	brush     *gdiBrush
	font      *gdiFont
	textColor color.NRGBA
	bkColor   color.NRGBA
	bkOpaque  bool
	textAlign uint32
	winding   bool
	mapMode   int
	winOrg    [2]float64
	winExt    [2]float64
	vpOrg     [2]float64
	vpExt     [2]float64
	world     canvas.Matrix
	clips     [][][2]float64 // output-space polygons
	cur       [2]float64
}

type gdi struct {
	cv     *canvas.Canvas
	obj    *bdf.Object
	opts   *Options
	st     gdiState
	stack  []gdiState
	out    canvas.Matrix // device → output points
	pxmm   float64       // device pixels per millimetre (metric map modes)
	path   *pathB
	open   bool
	openM  canvas.Matrix
	openC  int
	clipV  int
	objs   map[uint32]any
	slots  []any // WMF object table
	warned map[string]bool
	hatch  map[hatchKey]bdf.Hash // tiles of hatched brushes
}

func newGDI(cv *canvas.Canvas, opts *Options) *gdi {
	g := &gdi{cv: cv, obj: cv.Obj, opts: opts, objs: map[uint32]any{}, warned: map[string]bool{},
		hatch: map[hatchKey]bdf.Hash{}, pxmm: 96 / 25.4}
	g.st = gdiState{
		pen: &gdiPen{color: black}, brush: &gdiBrush{color: white, hatch: -1},
		font: &gdiFont{height: 12, weight: 400, face: "Arial"}, textColor: black, bkColor: white, bkOpaque: true,
		mapMode: 1, winExt: [2]float64{1, 1}, vpExt: [2]float64{1, 1}, world: canvas.Identity,
	}
	return g
}

// canvas.Matrix returns logical → output.
func (g *gdi) matrix() canvas.Matrix {
	st := &g.st
	var page canvas.Matrix
	switch st.mapMode {
	case 1: // MM_TEXT
		page = canvas.Translate(st.vpOrg[0]-st.winOrg[0], st.vpOrg[1]-st.winOrg[1])
	case 7, 8: // MM_ISOTROPIC, MM_ANISOTROPIC
		sx, sy := 1.0, 1.0
		if st.winExt[0] != 0 {
			sx = st.vpExt[0] / st.winExt[0]
		}
		if st.winExt[1] != 0 {
			sy = st.vpExt[1] / st.winExt[1]
		}
		if st.mapMode == 7 {
			m := math.Min(math.Abs(sx), math.Abs(sy))
			sx, sy = math.Copysign(m, sx), math.Copysign(m, sy)
		}
		page = canvas.Translate(st.vpOrg[0], st.vpOrg[1]).Mul(canvas.Scale(sx, sy)).Mul(canvas.Translate(-st.winOrg[0], -st.winOrg[1]))
	default: // metric modes: y grows upwards
		mm := map[int]float64{2: 0.1, 3: 0.01, 4: 0.254, 5: 0.0254, 6: 25.4 / 1440}[st.mapMode]
		if mm == 0 {
			mm = 0.1
		}
		k := mm * g.pxmm
		page = canvas.Translate(st.vpOrg[0], st.vpOrg[1]).Mul(canvas.Scale(k, -k)).Mul(canvas.Translate(-st.winOrg[0], -st.winOrg[1]))
	}
	return g.out.Mul(page).Mul(st.world)
}

func (g *gdi) warn(msg string) {
	if !g.warned[msg] {
		g.warned[msg] = true
		if g.opts.Warn != nil {
			g.opts.Warn(msg)
		}
	}
}

// color applies the picture's recoloring to a color of its records.
func (g *gdi) color(c color.NRGBA) color.NRGBA {
	if g.opts.Recolor == nil {
		return c
	}
	return g.opts.Recolor.Color(c)
}

// begin opens a drawing group with the current clip and transform.
func (g *gdi) begin() canvas.Matrix {
	m := g.matrix()
	if g.open && m == g.openM && g.openC == g.clipV {
		return m
	}
	g.end()
	g.obj.Save()
	for _, poly := range g.st.clips {
		p := newPath()
		for i, pt := range poly {
			if i == 0 {
				p.moveTo(pt[0], pt[1])
			} else {
				p.lineTo(pt[0], pt[1])
			}
		}
		p.close()
		g.obj.ClipPath(g.obj.AddPath(p.p), 0)
	}
	g.cv.Transform(m)
	g.open, g.openM, g.openC = true, m, g.clipV
	return m
}

func (g *gdi) end() {
	if g.open {
		g.obj.Restore()
		g.open = false
	}
}

func matScale(m canvas.Matrix) float64 { return math.Sqrt(math.Abs(m[0]*m[3] - m[1]*m[2])) }

// draw fills and/or strokes a path in logical coordinates.
func (g *gdi) draw(p *bdf.Path, fill, stroke bool) {
	if g.path != nil {
		// inside BeginPath: collect
		g.path.p.Verbs = append(g.path.p.Verbs, p.Verbs...)
		g.path.p.Args = append(g.path.p.Args, p.Args...)
		return
	}
	br, pen := g.st.brush, g.st.pen
	fill = fill && br != nil && !br.null
	stroke = stroke && pen != nil && !pen.null
	if !fill && !stroke {
		return
	}
	m := g.begin()
	ref := g.obj.AddPath(p)
	rule := byte(1)
	if g.st.winding {
		rule = 0
	}
	if fill {
		g.setBrush()
		g.obj.FillPath(ref, rule)
	}
	if stroke {
		g.setPen(m)
		g.obj.StrokePath(ref)
	}
	g.cv.Drawn = true
}

func (g *gdi) setBrush() {
	br := g.st.brush
	if br.hatch < 0 {
		g.obj.FillColor(toBDF(g.color(br.color)))
		return
	}
	var bg color.NRGBA
	if g.st.bkOpaque {
		bg = g.color(g.st.bkColor)
	}
	k := hatchKey{br.hatch % len(hatchBits), g.color(br.color), bg}
	h, ok := g.hatch[k]
	if !ok {
		h = g.opts.Doc.AddImage(hatchTile(k))
		g.hatch[k] = h
	}
	p := bdf.Pattern(g.cv.Image(h), bdf.RepeatBoth)
	p.Matrix = [6]float32{0.75, 0, 0, 0.75, 0, 0} // 96 dpi pixels
	g.obj.FillPaint(g.obj.AddPaint(p))
}

func (g *gdi) setPen(m canvas.Matrix) {
	pen := g.st.pen
	w := pen.width
	if w <= 0 || pen.cosmetic {
		w = 0.75 / math.Max(matScale(m), 1e-9) // a hairline of 0.75 pt
	}
	var cap, join byte = 1, 1
	if !pen.cosmetic {
		switch pen.style & 0xf00 {
		case 0x100:
			cap = 2
		case 0x200:
			cap = 0
		}
		switch pen.style & 0xf000 {
		case 0x1000:
			join = 2
		case 0x2000:
			join = 0
		}
	}
	g.obj.Line(f32(w), cap, join, 10)
	var dash []float64
	switch pen.style & 0xf {
	case 1:
		dash = []float64{3, 1}
	case 2:
		dash = []float64{1, 1}
	case 3:
		dash = []float64{3, 1, 1, 1}
	case 4:
		dash = []float64{3, 1, 1, 1, 1, 1}
	}
	if len(dash) > 0 {
		seg := make([]float32, len(dash))
		for i, d := range dash {
			seg[i] = f32(d * math.Max(w, 0.75/math.Max(matScale(m), 1e-9)) * 2)
		}
		g.obj.Dash(seg, 0)
	} else {
		g.obj.Dash(nil, 0)
	}
	g.obj.StrokeColor(toBDF(g.color(pen.color)))
}

func polyPath(pts [][2]float64, closed bool) *bdf.Path {
	p := newPath()
	for i, pt := range pts {
		if i == 0 {
			p.moveTo(pt[0], pt[1])
		} else {
			p.lineTo(pt[0], pt[1])
		}
	}
	if closed {
		p.close()
	}
	return p.p
}

func bezierPath(start [2]float64, pts [][2]float64, move bool) *bdf.Path {
	p := newPath()
	if move {
		p.moveTo(start[0], start[1])
	}
	for i := 0; i+2 < len(pts); i += 3 {
		p.p.CubicTo(f32(pts[i][0]), f32(pts[i][1]), f32(pts[i+1][0]), f32(pts[i+1][1]), f32(pts[i+2][0]), f32(pts[i+2][1]))
	}
	return p.p
}

func (g *gdi) rect(l, t, r, b float64) {
	p := &bdf.Path{}
	p.Rect(f32(math.Min(l, r)), f32(math.Min(t, b)), f32(math.Abs(r-l)), f32(math.Abs(b-t)))
	g.draw(p, true, true)
}

func (g *gdi) ellipse(l, t, r, b float64) {
	p := &bdf.Path{}
	p.Ellipse(f32((l+r)/2), f32((t+b)/2), f32(math.Abs(r-l)/2), f32(math.Abs(b-t)/2), 0, 0, 2*math.Pi, false)
	g.draw(p, true, true)
}

func (g *gdi) roundRect(l, t, r, b, w, h float64) {
	p := newPath()
	x0, y0, x1, y1 := math.Min(l, r), math.Min(t, b), math.Max(l, r), math.Max(t, b)
	rx, ry := math.Min(math.Abs(w)/2, (x1-x0)/2), math.Min(math.Abs(h)/2, (y1-y0)/2)
	p.p.Ellipse(f32(x0+rx), f32(y0+ry), f32(rx), f32(ry), 0, math.Pi, 1.5*math.Pi, false)
	p.p.Ellipse(f32(x1-rx), f32(y0+ry), f32(rx), f32(ry), 0, 1.5*math.Pi, 2*math.Pi, false)
	p.p.Ellipse(f32(x1-rx), f32(y1-ry), f32(rx), f32(ry), 0, 0, 0.5*math.Pi, false)
	p.p.Ellipse(f32(x0+rx), f32(y1-ry), f32(rx), f32(ry), 0, 0.5*math.Pi, math.Pi, false)
	p.close()
	g.draw(p.p, true, true)
}

// arc draws Arc (0), Chord (1) or Pie (2) from the bounding box and the
// radial end points (counter-clockwise in device space).
func (g *gdi) arc(kind int, l, t, r, b, xs, ys, xe, ye float64) {
	cx, cy := (l+r)/2, (t+b)/2
	rx, ry := math.Abs(r-l)/2, math.Abs(b-t)/2
	if rx == 0 || ry == 0 {
		return
	}
	a0 := math.Atan2((ys-cy)/ry, (xs-cx)/rx)
	a1 := math.Atan2((ye-cy)/ry, (xe-cx)/rx)
	p := newPath()
	if kind == 2 {
		p.moveTo(cx, cy)
	}
	// counter-clockwise on screen is decreasing angle in y-down space
	if a1 >= a0 {
		a1 -= 2 * math.Pi
	}
	p.p.Ellipse(f32(cx), f32(cy), f32(rx), f32(ry), 0, f32(a0), f32(a1), true)
	if kind > 0 {
		p.close()
	}
	g.draw(p.p, kind > 0, true)
}

// clipRect intersects the clip with a logical rectangle.
func (g *gdi) clipRect(l, t, r, b float64) {
	m := g.matrix()
	var poly [][2]float64
	for _, pt := range [][2]float64{{l, t}, {r, t}, {r, b}, {l, b}} {
		x, y := m.Apply(pt[0], pt[1])
		poly = append(poly, [2]float64{x, y})
	}
	g.st.clips = append(append([][][2]float64(nil), g.st.clips...), poly)
	g.clipV++
}

func (g *gdi) resetClip() {
	g.st.clips = nil
	g.clipV++
}

func (g *gdi) save() {
	st := g.st
	st.clips = append([][][2]float64(nil), g.st.clips...)
	g.stack = append(g.stack, st)
}

func (g *gdi) restore(n int) {
	if n < 0 {
		n = len(g.stack) + n
	} else {
		n = n - 1
	}
	if n < 0 || n >= len(g.stack) {
		return
	}
	g.st = g.stack[n]
	g.stack = g.stack[:n]
	g.clipV++
}

// image draws a DIB into a logical rectangle (source rectangle in pixels).
func (g *gdi) image(bmi, bits []byte, sx, sy, sw, sh, dx, dy, dw, dh float64) {
	img, err := decodeDIB(bmi, bits)
	if err != nil || img == nil {
		g.warn("a metafile bitmap could not be decoded")
		return
	}
	if g.opts.Recolor != nil {
		g.opts.Recolor.Image(img)
	}
	res, _ := imgconv.EncodeImage(img, true, g.opts.Images)
	h := g.opts.Doc.AddImage(res.Data)
	ref := g.cv.Image(h)
	g.begin()
	b := img.Bounds()
	if sw == 0 || sh == 0 {
		sw, sh = float64(b.Dx()), float64(b.Dy())
	}
	x0, y0, x1, y1 := dx, dy, dx+dw, dy+dh
	if dw < 0 || dh < 0 {
		// mirrored destination: flip around the rectangle
		g.obj.Save()
		fx, fy := 1.0, 1.0
		if dw < 0 {
			fx = -1
		}
		if dh < 0 {
			fy = -1
		}
		g.cv.Transform(canvas.Translate(dx, dy).Mul(canvas.Scale(fx, fy)))
		g.obj.ImageSub(ref, f32(sx), f32(sy), f32(sw), f32(sh), 0, 0, f32(math.Abs(dw)), f32(math.Abs(dh)))
		g.obj.Restore()
	} else {
		g.obj.ImageSub(ref, f32(sx), f32(sy), f32(sw), f32(sh), f32(x0), f32(y0), f32(x1-x0), f32(y1-y0))
	}
	g.cv.Drawn = true
}

// Size returns the size a metafile is meant to be shown at, in points: the
// frame of an EMF, or the bounding box of a placeable WMF in its units per
// inch. Other WMFs do not say; their window is taken to be in 96 dpi
// pixels. It returns false when the metafile has no picture box.
func Size(data []byte) (w, h float64, ok bool) {
	switch Kind(data) {
	case "emf":
		f, ok := emfFrame(data)
		if !ok {
			return 0, 0, false
		}
		return (f.box[2] - f.box[0]) / f.pxmm[0] * 72 / 25.4, (f.box[3] - f.box[1]) / f.pxmm[1] * 72 / 25.4, true
	case "wmf":
		bl, bt, br, bb, _, inch, ok := wmfBox(data)
		if !ok {
			return 0, 0, false
		}
		if inch <= 0 {
			inch = 96
		}
		return math.Abs(br-bl) / inch * 72, math.Abs(bb-bt) / inch * 72, true
	}
	return 0, 0, false
}

// Description returns the application and picture names an EMF header
// carries ("" when it has none, and for WMF).
func Description(data []byte) (app, picture string) {
	if Kind(data) != "emf" {
		return "", ""
	}
	n, off := le32(data, 60), le32(data, 64) // nDescription, offDescription
	s := slice(data, off, n*2)
	parts := strings.Split(string(utf16String(s, int(n))), "\x00")
	if len(parts) > 0 {
		app = strings.TrimSpace(parts[0])
	}
	if len(parts) > 1 {
		picture = strings.TrimSpace(parts[1])
	}
	return app, picture
}
