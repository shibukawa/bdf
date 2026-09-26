package pptx

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"math"
	"unicode/utf16"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/imgconv"
)

// Windows metafiles (EMF and WMF) are replayed into the object instead of
// being stored as images, which browsers cannot decode: a GDI state
// machine (map modes, world transform, pens, brushes, fonts, clipping,
// paths) turns the drawing records into BDF paths, text and images. EMF+
// records embedded in comments are skipped; Office writes the equivalent
// EMF records next to them.

type gdiPen struct {
	null     bool
	width    float64 // logical units; 0 = hairline
	cosmetic bool
	color    rgba
	style    uint32
}

type gdiBrush struct {
	null  bool
	color rgba
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
	textColor rgba
	bkColor   rgba
	bkOpaque  bool
	textAlign uint32
	winding   bool
	mapMode   int
	winOrg    [2]float64
	winExt    [2]float64
	vpOrg     [2]float64
	vpExt     [2]float64
	world     matrix
	clips     [][][2]float64 // output-space polygons
	cur       [2]float64
}

type gdi struct {
	s      *slideCtx
	cv     *canvas
	st     gdiState
	stack  []gdiState
	out    matrix  // device → output points
	pxmm   float64 // device pixels per millimetre (metric map modes)
	path   *pathB
	open   bool
	openM  matrix
	openC  int
	clipV  int
	wmf    bool
	objs   map[uint32]any
	slots  []any // WMF object table
	warned map[string]bool
	rc     *recolor // picture recoloring of the blip
}

func newGDI(s *slideCtx, cv *canvas) *gdi {
	g := &gdi{s: s, cv: cv, objs: map[uint32]any{}, warned: map[string]bool{}, pxmm: 96 / 25.4}
	g.st = gdiState{
		pen: &gdiPen{color: black}, brush: &gdiBrush{color: white, hatch: -1},
		font: &gdiFont{height: 12, weight: 400, face: "Arial"}, textColor: black, bkColor: white, bkOpaque: true,
		mapMode: 1, winExt: [2]float64{1, 1}, vpExt: [2]float64{1, 1}, world: identity,
	}
	return g
}

// matrix returns logical → output.
func (g *gdi) matrix() matrix {
	st := &g.st
	var page matrix
	switch st.mapMode {
	case 1: // MM_TEXT
		page = translate(st.vpOrg[0]-st.winOrg[0], st.vpOrg[1]-st.winOrg[1])
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
		page = translate(st.vpOrg[0], st.vpOrg[1]).mul(scale(sx, sy)).mul(translate(-st.winOrg[0], -st.winOrg[1]))
	default: // metric modes: y grows upwards
		mm := map[int]float64{2: 0.1, 3: 0.01, 4: 0.254, 5: 0.0254, 6: 25.4 / 1440}[st.mapMode]
		if mm == 0 {
			mm = 0.1
		}
		k := mm * g.pxmm
		page = translate(st.vpOrg[0], st.vpOrg[1]).mul(scale(k, -k)).mul(translate(-st.winOrg[0], -st.winOrg[1]))
	}
	return g.out.mul(page).mul(st.world)
}

func (g *gdi) warn(key, msg string) {
	if !g.warned[key] {
		g.warned[key] = true
		g.s.c.warnOnce("mf:"+key, "%s", msg)
	}
}

// begin opens a drawing group with the current clip and transform.
func (g *gdi) begin() matrix {
	m := g.matrix()
	if g.open && m == g.openM && g.openC == g.clipV {
		return m
	}
	g.end()
	g.cv.obj.Save()
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
		g.cv.obj.ClipPath(g.cv.obj.AddPath(p.p), 0)
	}
	g.cv.transform(m)
	g.open, g.openM, g.openC = true, m, g.clipV
	return m
}

func (g *gdi) end() {
	if g.open {
		g.cv.obj.Restore()
		g.open = false
	}
}

func matScale(m matrix) float64 { return math.Sqrt(math.Abs(m[0]*m[3] - m[1]*m[2])) }

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
	ref := g.cv.obj.AddPath(p)
	rule := byte(1)
	if g.st.winding {
		rule = 0
	}
	if fill {
		g.setBrush()
		g.cv.obj.FillPath(ref, rule)
	}
	if stroke {
		g.setPen(m)
		g.cv.obj.StrokePath(ref)
	}
	g.cv.drawn = true
}

func (g *gdi) setBrush() {
	br := g.st.brush
	if br.hatch < 0 {
		g.cv.obj.FillColor(g.rc.color(br.color).bdf())
		return
	}
	names := []string{"horz", "vert", "dnDiag", "upDiag", "cross", "diagCross"}
	f := fill{kind: fillPatt, patt: names[br.hatch%len(names)], fg: g.rc.color(br.color), bg: rgba{}}
	if g.st.bkOpaque {
		f.bg = g.rc.color(g.st.bkColor)
	}
	g.cv.setFill(f, 8, 8)
}

func (g *gdi) setPen(m matrix) {
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
	g.cv.obj.Line(f32(w), cap, join, 10)
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
		g.cv.obj.Dash(seg, 0)
	} else {
		g.cv.obj.Dash(nil, 0)
	}
	g.cv.obj.StrokeColor(g.rc.color(pen.color).bdf())
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
		x, y := m.apply(pt[0], pt[1])
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

// text draws a string at the logical reference point; dx are the logical
// advances of the characters (nil: the font's).
func (g *gdi) text(x, y float64, s []rune, dx []float64, opaque *[4]float64) {
	if len(s) == 0 {
		return
	}
	m := g.matrix()
	if opaque != nil {
		save := g.st.brush
		g.st.brush = &gdiBrush{color: g.st.bkColor, hatch: -1}
		pen := g.st.pen
		g.st.pen = &gdiPen{null: true}
		g.rect(opaque[0], opaque[1], opaque[2], opaque[3])
		g.st.brush, g.st.pen = save, pen
	}
	if g.st.textAlign&1 != 0 { // TA_UPDATECP
		x, y = g.st.cur[0], g.st.cur[1]
	}
	f := g.st.font
	// scale of the logical space along the text and across it
	sx := math.Hypot(m[0], m[1])
	sy := math.Hypot(m[2], m[3])
	em := math.Abs(f.height)
	st := &runStyle{size: 0, underline: "none", strike: "noStrike", caps: "none", latin: f.face, ea: f.face,
		bold: f.weight >= 600, italic: f.italic, fill: fill{kind: fillSolid, color: g.rc.color(g.st.textColor)}}
	if f.underline {
		st.underline = "sng"
	}
	if f.strike {
		st.strike = "sngStrike"
	}
	text := string(s)
	if isSymbolFont(f.face) {
		text = mapSymbolFont(f.face, text)
		st.latin, st.ea = "", ""
	}
	st.setKey()
	fc0 := g.s.c.faceFor(st, []rune(text)[0])
	if f.height > 0 {
		// cell height: the em is the cell minus the internal leading
		em = f.height / math.Max(fc0.asc+fc0.desc, 0.5)
	}
	if em == 0 {
		em = 12 / math.Max(sy, 1e-9)
	}
	st.size = em * sy
	st.setKey()
	pa := &para{}
	g.s.addChars(pa, st, text)
	adv := 0.0
	for i := range pa.items {
		pa.items[i].x = adv
		if i < len(dx) {
			pa.items[i].w = dx[i] * sx
		}
		adv += pa.items[i].w
	}
	// reference point to baseline
	asc, desc := fc0.asc*st.size, fc0.desc*st.size
	var oy float64
	switch g.st.textAlign & 24 {
	case 0: // TA_TOP
		oy = asc
	case 8: // TA_BOTTOM
		oy = -desc
	}
	var ox float64
	switch g.st.textAlign & 6 {
	case 6: // TA_CENTER
		ox = -adv / 2
	case 2: // TA_RIGHT
		ox = -adv
	}
	if g.st.textAlign&1 != 0 {
		// advance the current position
		g.st.cur[0] += adv / math.Max(sx, 1e-9)
	}
	px, py := m.apply(x, y)
	// the text stands upright in output space even when the logical y
	// axis points up; escapement turns it counter-clockwise
	rot := math.Atan2(m[1], m[0])*180/math.Pi - f.escapement/10
	g.end()
	g.cv.obj.Save()
	for _, poly := range g.st.clips {
		g.cv.obj.ClipPath(g.cv.obj.AddPath(polyPath(poly, true)), 0)
	}
	g.cv.transform(translate(px, py).mul(rotate(rot)))
	em2 := &textEmitter{cv: g.cv, m: identity}
	g.cv.obj.Mark(bdf.MarkBox, "")
	em2.emitItems(&textLine{items: pa.items}, ox, oy)
	g.cv.obj.Restore()
}

// image draws a DIB into a logical rectangle (source rectangle in pixels).
func (g *gdi) image(bmi, bits []byte, sx, sy, sw, sh, dx, dy, dw, dh float64) {
	img, err := decodeDIB(bmi, bits)
	if err != nil || img == nil {
		g.warn("dib", "a metafile bitmap could not be decoded")
		return
	}
	g.rc.image(img, 0)
	res, _ := imgconv.EncodeImage(img, true, g.s.c.opts.Images)
	h := g.s.c.doc.AddImage(res.Data)
	ref := g.cv.image(h)
	g.begin()
	b := img.Bounds()
	if sw == 0 || sh == 0 {
		sw, sh = float64(b.Dx()), float64(b.Dy())
	}
	x0, y0, x1, y1 := dx, dy, dx+dw, dy+dh
	if dw < 0 || dh < 0 {
		// mirrored destination: flip around the rectangle
		g.cv.obj.Save()
		fx, fy := 1.0, 1.0
		if dw < 0 {
			fx = -1
		}
		if dh < 0 {
			fy = -1
		}
		g.cv.transform(translate(dx, dy).mul(scale(fx, fy)))
		g.cv.obj.ImageSub(ref, f32(sx), f32(sy), f32(sw), f32(sh), 0, 0, f32(math.Abs(dw)), f32(math.Abs(dh)))
		g.cv.obj.Restore()
	} else {
		g.cv.obj.ImageSub(ref, f32(sx), f32(sy), f32(sw), f32(sh), f32(x0), f32(y0), f32(x1-x0), f32(y1-y0))
	}
	g.cv.drawn = true
}

// decodeDIB decodes a device-independent bitmap (BITMAPINFO + bits).
func decodeDIB(bmi, bits []byte) (*image.NRGBA, error) {
	if len(bmi) < 12 {
		return nil, errBadDIB
	}
	hsize := int(le32(bmi, 0))
	var w, h, bpp, compression int
	if hsize == 12 { // BITMAPCOREHEADER
		w, h, bpp = int(le16(bmi, 4)), int(int16(le16(bmi, 6))), int(le16(bmi, 10))
	} else {
		if len(bmi) < 40 {
			return nil, errBadDIB
		}
		w, h = int(int32(le32(bmi, 4))), int(int32(le32(bmi, 8)))
		bpp, compression = int(le16(bmi, 14)), int(le32(bmi, 16))
	}
	if compression == 4 || compression == 5 { // BI_JPEG, BI_PNG
		img, _, err := image.Decode(bytes.NewReader(bits))
		if err != nil {
			return nil, err
		}
		return toNRGBA(img), nil
	}
	topDown := h < 0
	if h < 0 {
		h = -h
	}
	if w <= 0 || h <= 0 || w*h > 1<<26 {
		return nil, errBadDIB
	}
	// color table
	var pal []color.NRGBA
	if bpp <= 8 {
		n := 1 << bpp
		if hsize >= 40 {
			if used := int(le32(bmi, 32)); used > 0 && used < n {
				n = used
			}
		}
		entry := 4
		if hsize == 12 {
			entry = 3
		}
		off := hsize
		for i := 0; i < n && off+entry <= len(bmi); i++ {
			pal = append(pal, color.NRGBA{bmi[off+2], bmi[off+1], bmi[off], 255})
			off += entry
		}
	}
	masks := [3]uint32{0x7c00, 0x03e0, 0x001f}
	if bpp == 32 {
		masks = [3]uint32{0xff0000, 0xff00, 0xff}
	}
	if compression == 3 && len(bmi) >= hsize+12 {
		masks = [3]uint32{le32(bmi, hsize), le32(bmi, hsize+4), le32(bmi, hsize+8)}
		if hsize >= 52 {
			masks = [3]uint32{le32(bmi, 40), le32(bmi, 44), le32(bmi, 48)}
		}
	} else if compression != 0 {
		return nil, errBadDIB // RLE
	}
	stride := ((w*bpp + 31) / 32) * 4
	if len(bits) < stride*h {
		return nil, errBadDIB
	}
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	field := func(v, mask uint32) uint8 {
		if mask == 0 {
			return 0
		}
		shift := 0
		for mask&1 == 0 {
			mask >>= 1
			shift++
		}
		return uint8((v >> shift & mask) * 255 / mask)
	}
	for row := 0; row < h; row++ {
		src := bits[row*stride:]
		y := h - 1 - row
		if topDown {
			y = row
		}
		for x := 0; x < w; x++ {
			var c color.NRGBA
			switch bpp {
			case 1, 2, 4, 8:
				idx := int(src[x*bpp/8]>>(8-bpp-(x*bpp%8))) & (1<<bpp - 1)
				if idx < len(pal) {
					c = pal[idx]
				} else {
					c = color.NRGBA{0, 0, 0, 255}
				}
			case 16:
				v := uint32(le16(src, x*2))
				c = color.NRGBA{field(v, masks[0]), field(v, masks[1]), field(v, masks[2]), 255}
			case 24:
				c = color.NRGBA{src[x*3+2], src[x*3+1], src[x*3], 255}
			case 32:
				v := le32(src, x*4)
				c = color.NRGBA{field(v, masks[0]), field(v, masks[1]), field(v, masks[2]), 255}
			default:
				return nil, errBadDIB
			}
			img.SetNRGBA(x, y, c)
		}
	}
	return img, nil
}

type dibError struct{}

func (dibError) Error() string { return "bad DIB" }

var errBadDIB error = dibError{}

func le16(b []byte, i int) uint16 {
	if i < 0 || i+2 > len(b) {
		return 0
	}
	return binary.LittleEndian.Uint16(b[i:])
}

func le32(b []byte, i int) uint32 {
	if i < 0 || i+4 > len(b) {
		return 0
	}
	return binary.LittleEndian.Uint32(b[i:])
}

func lei32(b []byte, i int) float64 { return float64(int32(le32(b, i))) }
func lei16(b []byte, i int) float64 { return float64(int16(le16(b, i))) }
func lef32(b []byte, i int) float64 { return float64(math.Float32frombits(le32(b, i))) }

func colorRef(v uint32) rgba {
	return rgba{float64(v&0xff) / 255, float64(v>>8&0xff) / 255, float64(v>>16&0xff) / 255, 1}
}

func utf16String(b []byte, n int) []rune {
	if n*2 > len(b) {
		n = len(b) / 2
	}
	u := make([]uint16, n)
	for i := range u {
		u[i] = le16(b, i*2)
	}
	for len(u) > 0 && u[len(u)-1] == 0 {
		u = u[:len(u)-1]
	}
	return utf16.Decode(u)
}

// metafileKind reports "emf", "wmf" or "".
func metafileKind(data []byte) string {
	switch {
	case len(data) >= 44 && le32(data, 0) == 1 && string(data[40:44]) == " EMF":
		return "emf"
	case len(data) >= 22 && le32(data, 0) == 0x9AC6CDD7,
		len(data) >= 18 && (le16(data, 0) == 1 || le16(data, 0) == 2) && le16(data, 2) == 9:
		return "wmf"
	}
	return ""
}

// drawMetafile replays a metafile into the rectangle x,y,w,h.
func (s *slideCtx) drawMetafile(cv *canvas, data []byte, kind string, rc *recolor, x, y, w, h float64) {
	g := newGDI(s, cv)
	g.rc = rc
	cv.obj.Save()
	cv.obj.ClipRect(f32(x), f32(y), f32(w), f32(h))
	defer func() {
		if r := recover(); r != nil {
			s.c.warnOnce("mfpanic", "a metafile could not be replayed completely")
		}
		g.end()
		cv.obj.Restore()
	}()
	if kind == "emf" {
		g.playEMF(data, x, y, w, h)
	} else {
		g.wmf = true
		g.playWMF(data, x, y, w, h)
	}
}
