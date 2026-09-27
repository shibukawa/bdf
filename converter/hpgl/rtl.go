package hpgl

import (
	"image"
	"image/color"
	"math"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/imgconv"
)

// HP RTL (and the raster graphics of PCL): images sent a row at a time,
// compressed, as palette indexes or direct colors. Coordinates of the
// raster context ("RTL space") are plotter units from the top left corner
// of the logical page, x across and y down; for a plotter, x runs across
// the width of the plot size and y along its length, which HP-GL/2 turns
// so that its x axis is the longer side; for a PCL printer they are the
// coordinates of the logical page.

type rtl struct {
	c *converter

	native float64 // native resolution (dpi): the unit of cursor moves
	res    float64 // graphics resolution (dpi): the size of a pixel
	srcW   int     // source raster width and height (pixels; 0 unset)
	srcH   int
	dstW   float64 // destination raster width and height (decipoints)
	dstH   float64
	comp   int
	cap    cad.Point // current active position (RTL space)
	// the plot size of a plotter (plotter units): length and width
	psL, psW float64

	// color set-up
	encoding     int // 0 indexed by plane, 1 indexed by pixel, 2 direct by plane, 3 direct by pixel, 4 plane by plane
	bitsIndex    int
	bits         [3]int
	white, black [3]float64
	palette      []bdf.Color
	rgb          [3]float64 // *v#A, B, C
	transparent  bool       // source transparency: white pixels leave the page as it is
	simple       int        // simple color mode (*r#U): 0 none, 1, 3, -3, -4
	assigned     bool       // *v#I changed the palette since the context was entered

	img *raster // the image being received
}

func newRTL(c *converter) *rtl {
	r := &rtl{c: c}
	r.reset()
	return r
}

// reset is ESC E: black and white, one plane, the resolution and the
// cursor at their defaults.
func (r *rtl) reset() {
	r.end()
	r.res = 0
	r.srcW, r.srcH = 0, 0
	r.dstW, r.dstH = 0, 0
	r.comp = 0
	r.cap = cad.Point{}
	r.encoding = 0
	r.bitsIndex = 1
	r.bits = [3]int{8, 8, 8}
	r.white = [3]float64{255, 255, 255}
	r.black = [3]float64{}
	r.palette = []bdf.Color{bdf.RGB(255, 255, 255), bdf.RGB(0, 0, 0)}
	r.transparent = true
	r.simple = 0
	r.psL, r.psW = 0, 0
}

// pageAdvanced moves the cursor to the top of the new page.
func (r *rtl) pageAdvanced() { r.cap = cad.Point{} }

// plotSize records the plot size of HP-GL/2 (plotter units).
func (r *rtl) plotSize(length, width float64) { r.psL, r.psW = length, width }

// resolution returns the graphics resolution.
func (r *rtl) resolution() float64 {
	if r.res > 0 {
		return r.res
	}
	if r.c.printer() {
		return 75 // PCL
	}
	return 300
}

// nativeUnit returns the size of a cursor unit in plotter units.
func (r *rtl) nativeUnit() float64 {
	if r.c.printer() {
		return 1016 / r.c.pcl.units
	}
	if r.native > 0 {
		return 1016 / r.native
	}
	return 1016.0 / 300
}

// toDrawing returns the transform of RTL space into the drawing.
func (r *rtl) toDrawing() canvas.Matrix {
	if r.c.printer() {
		return r.c.pcl.logical()
	}
	l, w := r.psL, r.psW
	if l == 0 {
		l, w = defaultHardW, defaultHardH
	}
	if l >= w {
		return canvas.Matrix{0, 1, 1, 0, 0, 0}
	}
	return canvas.Matrix{-1, 0, 0, 1, w, 0}
}

// view returns the transform of the drawing into view coordinates (y up)
// that shows raster images upright, for a page that has only them.
func (r *rtl) view() (canvas.Matrix, bool) {
	if r.c.printer() {
		return canvas.Matrix{}, false
	}
	// RTL space seen with y up
	m := canvas.Matrix{1, 0, 0, -1, 0, 0}.Mul(invert(r.toDrawing()))
	return m, true
}

// enter is ESC % # A: the raster context; move reports that the cursor
// goes to the HP-GL/2 pen.
func (r *rtl) enter(move bool) {
	if move {
		r.cap = apply(invert(r.toDrawing()), r.c.gl.pos)
	}
	// a plotter shares one palette between the two contexts
	if !r.c.printer() && r.simple == 0 {
		p := &r.c.gl.pens
		n := min(len(r.palette), p.n)
		copy(r.palette, p.colors[:n])
	}
	r.assigned = false
}

// leave hands the palette to HP-GL/2 when raster commands changed it.
func (r *rtl) leave() {
	if r.c.printer() || !r.assigned {
		return
	}
	p := &r.c.gl.pens
	for i := 0; i < len(r.palette) && i < p.n; i++ {
		p.colors[i] = r.palette[i]
	}
	p.valid = false
}

// capPoint returns the cursor in the drawing (to move the pen to on
// entering HP-GL/2).
func (r *rtl) capPoint() (cad.Point, bool) {
	return apply(r.toDrawing(), r.cap), true
}

// cursor carries out the commands that move the cursor, which end a
// raster image.
func (r *rtl) cursor(g, cmd byte, v float64, rel bool) {
	if !valid(v) {
		return
	}
	r.end()
	var d float64
	switch {
	case g == 'p' && (cmd == 'X' || cmd == 'Y'):
		d = v * r.nativeUnit()
	case g == 'a' && (cmd == 'H' || cmd == 'V'):
		d = v * 1016 / 720 // decipoints
	default:
		return
	}
	horizontal := cmd == 'X' || cmd == 'H'
	switch {
	case horizontal && rel:
		r.cap.X += d
	case horizontal:
		r.cap.X = d
	case rel:
		r.cap.Y += d
	default:
		r.cap.Y = d
	}
}

// command carries out the raster commands (ESC * r, b, t, v).
func (r *rtl) command(g, cmd byte, v float64, rel bool, data []byte) {
	c := r.c
	if !valid(v) {
		return
	}
	switch g {
	case 'r':
		switch cmd {
		case 'A':
			r.start(int(v))
		case 'C', 'B':
			r.end()
		case 'S':
			if r.img == nil {
				r.srcW = clampInt(v, 0, 65535)
			}
		case 'T':
			if r.img == nil {
				r.srcH = clampInt(v, 0, 65535)
			}
		case 'U':
			if r.img == nil {
				r.simpleColor(int(v))
			}
		}
	case 'b':
		switch cmd {
		case 'M':
			r.comp = int(v)
		case 'W', 'V':
			c.rasters++
			r.transfer(cmd == 'V', data)
		case 'Y':
			if r.img != nil {
				r.img.skip(clampInt(v, 0, 1<<20))
			} else {
				r.cap.Y += v * r.nativeUnit()
			}
		}
	case 't':
		switch cmd {
		case 'R':
			if r.img == nil && v > 0 && v <= 4800 {
				r.res = v
			}
		case 'H':
			r.dstW = math.Max(0, v)
		case 'V':
			r.dstH = math.Max(0, v)
		}
	case 'v':
		if r.img != nil && cmd != 'W' {
			return
		}
		switch cmd {
		case 'W':
			if r.img == nil {
				r.configure(data)
			}
		case 'A', 'B', 'C':
			r.rgb[cmd-'A'] = v
		case 'I':
			r.assign(int(v))
		case 'N':
			r.transparent = v == 0
		}
	}
}

func clampInt(v float64, lo, hi int) int {
	if v < float64(lo) {
		return lo
	}
	if v > float64(hi) {
		return hi
	}
	return int(v)
}

// configure is configure image data (ESC * v # W).
func (r *rtl) configure(d []byte) {
	if len(d) < 6 {
		return
	}
	if d[1] > 4 {
		return
	}
	r.encoding = int(d[1])
	r.bitsIndex = max(1, min(8, int(d[2])))
	for i := 0; i < 3; i++ {
		r.bits[i] = int(d[3+i])
		r.black[i] = 0
		r.white[i] = math.Pow(2, float64(r.bits[i])) - 1
	}
	if len(d) >= 18 {
		for i := 0; i < 3; i++ {
			r.white[i] = float64(int16(uint16(d[6+2*i])<<8 | uint16(d[7+2*i])))
			r.black[i] = float64(int16(uint16(d[12+2*i])<<8 | uint16(d[13+2*i])))
		}
	}
	r.simple = 0
	r.palette = defaultPalette(r.bitsIndex)
}

// defaultPalette is the palette configure image data creates.
func defaultPalette(bits int) []bdf.Color {
	n := 1 << bits
	p := make([]bdf.Color, n)
	switch bits {
	case 1:
		p[0], p[1] = bdf.RGB(255, 255, 255), bdf.RGB(0, 0, 0)
	case 3:
		for i := range p {
			p[i] = bdf.RGB(uint8(i&1*255), uint8(i>>1&1*255), uint8(i>>2&1*255))
		}
	default:
		for i := range p {
			p[i] = defaultColor(i)
		}
	}
	return p
}

// simpleColor is ESC * r # U: a fixed palette of planes.
func (r *rtl) simpleColor(v int) {
	switch v {
	case 1:
		r.bitsIndex, r.palette = 1, defaultPalette(1)
	case 3:
		r.bitsIndex, r.palette = 3, defaultPalette(3)
	case -3:
		r.bitsIndex = 3
		r.palette = make([]bdf.Color, 8)
		for i := range r.palette {
			r.palette[i] = bdf.RGB(uint8((1-i&1)*255), uint8((1-i>>1&1)*255), uint8((1-i>>2&1)*255))
		}
	case -4:
		r.bitsIndex = 4
		r.palette = make([]bdf.Color, 16)
		for i := range r.palette {
			if i&1 != 0 {
				r.palette[i] = bdf.RGB(0, 0, 0)
				continue
			}
			r.palette[i] = bdf.RGB(uint8((1-i>>1&1)*255), uint8((1-i>>2&1)*255), uint8((1-i>>3&1)*255))
		}
	default:
		return
	}
	r.encoding, r.simple = 0, v
}

// assign is assign color index (ESC * v # I).
func (r *rtl) assign(i int) {
	defer func() { r.rgb = [3]float64{} }()
	if i < 0 || i >= len(r.palette) {
		return
	}
	var c [3]uint8
	for k := 0; k < 3; k++ {
		c[k] = r.primary(r.rgb[k], k)
	}
	r.palette[i] = bdf.RGB(c[0], c[1], c[2])
	r.assigned = true
}

// primary maps a primary value by the black and white references.
func (r *rtl) primary(v float64, k int) uint8 {
	b, w := r.black[k], r.white[k]
	if w == b {
		return 0
	}
	t := (v - b) / (w - b)
	return uint8(math.Round(math.Max(0, math.Min(1, t)) * 255))
}

// --- receiving an image ---

// raster is an image being received, scaled down by an integer factor as
// it arrives.
type raster struct {
	r      *rtl
	origin cad.Point // top left corner (RTL space)
	pw, ph float64   // pixel size (plotter units)
	planes int       // planes per row
	bpp    int       // bits per pixel of a row (pixel encodings)
	width  int       // pixels per row (0: taken from the rows)
	height int       // rows (0: unbounded)
	rows   int       // rows received
	plane  int       // the plane being received
	cur    [][]byte  // the planes of the row being received
	seed   [][]byte  // seed rows of delta row compression
	prev   []color.NRGBA
	acc    accumulator
	bad    bool // an encoding not drawn: its rows are left empty
}

// start is start raster graphics (ESC * r # A): 0 and 2 at the left edge,
// 1 and 3 at the cursor; 2 and 3 scaled to the destination size.
func (r *rtl) start(mode int) {
	if r.img != nil {
		return
	}
	res := r.resolution()
	img := &raster{r: r, pw: 1016 / res, ph: 1016 / res, width: r.srcW, height: r.srcH}
	img.origin = cad.Point{Y: r.cap.Y}
	if mode == 1 || mode == 3 {
		img.origin.X = r.cap.X
	}
	if (mode == 2 || mode == 3) && r.srcW > 0 && r.srcH > 0 && (r.dstW > 0 || r.dstH > 0) {
		dw, dh := r.dstW*1016/720, r.dstH*1016/720
		switch {
		case dw == 0:
			dw = dh * float64(r.srcW) / float64(r.srcH)
		case dh == 0:
			dh = dw * float64(r.srcH) / float64(r.srcW)
		}
		img.pw, img.ph = dw/float64(r.srcW), dh/float64(r.srcH)
	}
	switch r.encoding {
	case 1:
		img.planes, img.bpp = 1, r.bitsIndex
	case 3:
		img.planes, img.bpp = 1, r.bits[0]+r.bits[1]+r.bits[2]
	case 2:
		img.planes = r.bits[0] + r.bits[1] + r.bits[2]
	default:
		img.planes = r.bitsIndex
	}
	if img.planes < 1 || img.planes > 24 || img.bpp > 48 {
		r.c.warnOnce("rtlmode", "raster images in this pixel encoding are not drawn")
		img.planes, img.bpp, img.bad = 1, 0, true
	}
	img.cur = make([][]byte, img.planes)
	img.seed = make([][]byte, img.planes)
	img.acc.init(r, img)
	r.img = img
}

// end is end raster graphics: the image is stored and drawn, the cursor
// moves below it, and the compression method returns to 0.
func (r *rtl) end() {
	img := r.img
	if img == nil {
		return
	}
	r.img = nil
	if img.plane > 0 {
		img.emit()
	}
	if img.height > 0 && img.rows < img.height {
		img.skip(img.height - img.rows)
	}
	r.cap.Y = img.origin.Y + float64(img.rows)*img.ph
	r.cap.X = img.origin.X
	r.comp = 0
	pix := img.acc.image()
	if pix == nil {
		return
	}
	c := r.c
	b := pix.Bounds()
	w := float64(img.acc.inW) * img.pw
	h := float64(img.rows) * img.ph
	if w <= 0 || h <= 0 {
		return
	}
	opts := c.opts.Images
	tw, th := opts.FitSize(b.Dx(), b.Dy(), w*plu, h*plu)
	if tw != b.Dx() || th != b.Dy() {
		pix = imgconv.Resize(pix, tw, th)
	}
	res, _ := imgconv.EncodePixels(pix, true, opts)
	hash := c.doc.AddImage(res.Data)
	c.images++
	if !c.budget() {
		return
	}
	// the stored pixels onto the image in RTL space, then the drawing
	pw, ph := pix.Bounds().Dx(), pix.Bounds().Dy()
	m := r.toDrawing().Mul(canvas.Matrix{w / float64(pw), 0, 0, h / float64(ph), img.origin.X, img.origin.Y})
	c.gl.flush()
	c.cur.d.Image(&cad.Image{Image: hash, W: pw, H: ph, M: m, Smooth: true})
}

// transfer receives a plane (V) or a row (W) of raster data.
func (r *rtl) transfer(plane bool, data []byte) {
	if r.img == nil {
		r.start(0)
	}
	img := r.img
	switch r.comp {
	case 0, 1, 2, 3, 9:
		img.decodePlane(r.comp, data)
		if plane && img.plane+1 < img.planes {
			img.plane++
		} else {
			img.emit()
		}
	case 5:
		img.adaptive(data)
	default:
		r.c.warnOnce("rtlcomp", "raster data compressed with method %d is not drawn", r.comp)
		if !plane {
			img.skip(1)
		}
	}
}

// maxRasterWidth is the widest raster row, in pixels (the limit of HP RTL).
const maxRasterWidth = 65535

// rowBytes returns the bytes of a plane's row: those of the source width,
// or at most those of the widest row when it is not set.
func (img *raster) rowBytes() (n int, known bool) {
	bpp := max(1, img.bpp)
	if img.width <= 0 {
		return (maxRasterWidth*bpp + 7) / 8, false
	}
	return (img.width*bpp + 7) / 8, true
}

// decodePlane decompresses a plane of the current row into cur, updating
// the plane's seed row.
func (img *raster) decodePlane(method int, data []byte) {
	p := img.plane
	seed := img.seed[p]
	n, known := img.rowBytes()
	var row []byte
	switch method {
	case 0:
		row = append([]byte(nil), data[:min(len(data), n)]...)
	case 1:
		for i := 0; i+1 < len(data) && len(row) < n; i += 2 {
			for k := 0; k <= int(data[i]); k++ {
				row = append(row, data[i+1])
			}
		}
	case 2:
		row = packbits(data, n)
	case 3:
		row = deltaRow(seed, data, n)
	case 9:
		row = replacementDelta(seed, data, n)
	}
	if len(row) > n {
		row = row[:n]
	} else if known && len(row) < n {
		row = append(row, make([]byte, n-len(row))...)
	}
	img.cur[p] = row
	img.seed[p] = row
}

// packbits decodes TIFF PackBits, up to max bytes.
func packbits(d []byte, max int) []byte {
	var out []byte
	for i := 0; i < len(d) && len(out) < max; {
		n := int(int8(d[i]))
		i++
		switch {
		case n >= 0:
			end := min(len(d), i+n+1)
			out = append(out, d[i:end]...)
			i = end
		case n != -128 && i < len(d):
			for k := 0; k < 1-n; k++ {
				out = append(out, d[i])
			}
			i++
		}
	}
	return out
}

// deltaRow decodes seed-row (delta row) compression, into a row of at
// most max bytes.
func deltaRow(seed, d []byte, max int) []byte {
	row := append([]byte(nil), seed...)
	pos := 0
	for i := 0; i < len(d); {
		cmd := d[i]
		i++
		count := int(cmd>>5) + 1
		off := int(cmd & 31)
		if off == 31 {
			for i < len(d) {
				b := d[i]
				i++
				off += int(b)
				if b != 255 {
					break
				}
			}
		}
		pos += off
		for k := 0; k < count && i < len(d); k++ {
			if pos >= max {
				return row
			}
			for pos >= len(row) {
				row = append(row, 0)
			}
			row[pos] = d[i]
			pos++
			i++
		}
	}
	return row
}

// replacementDelta decodes the replacement delta row compression of PCL
// (method 9): literal or repeated replacement bytes at offsets from the
// last byte replaced, into a row of at most max bytes.
func replacementDelta(seed, d []byte, max int) []byte {
	row := append([]byte(nil), seed...)
	pos := 0
	set := func(b byte) {
		if pos >= max {
			pos++
			return
		}
		for pos >= len(row) {
			row = append(row, 0)
		}
		row[pos] = b
		pos++
	}
	extend := func(i int, v int) (int, int) {
		for i < len(d) {
			b := d[i]
			i++
			v += int(b)
			if b != 255 {
				break
			}
		}
		return i, v
	}
	for i := 0; i < len(d); {
		cmd := d[i]
		i++
		var off, count int
		if cmd&0x80 != 0 {
			// repeated: 2 bits of offset, 5 of count (+2)
			off, count = int(cmd>>5&3), int(cmd&31)
			if off == 3 {
				i, off = extend(i, off)
			}
			if count == 31 {
				i, count = extend(i, count)
			}
			count += 2
			pos += off
			if i >= len(d) || pos >= max {
				break
			}
			b := d[i]
			i++
			for k := 0; k < count && pos < max; k++ {
				set(b)
			}
		} else {
			// literal: 4 bits of offset, 3 of count (+1)
			off, count = int(cmd>>3&15), int(cmd&7)
			if off == 15 {
				i, off = extend(i, off)
			}
			if count == 7 {
				i, count = extend(i, count)
			}
			count++
			pos += off
			for k := 0; k < count && i < len(d); k++ {
				set(d[i])
				i++
			}
			if pos >= max {
				break
			}
		}
	}
	return row
}

// adaptive decodes a block of adaptive compression (method 5).
func (img *raster) adaptive(d []byte) {
	for i := 0; i+2 < len(d); {
		cmd := d[i]
		n := int(d[i+1])<<8 | int(d[i+2])
		i += 3
		switch cmd {
		case 0, 1, 2, 3:
			end := min(len(d), i+n)
			img.decodePlane(int(cmd), d[i:end])
			i = end
			if img.plane+1 < img.planes {
				img.plane++
			} else {
				img.emit()
			}
		case 4:
			img.skip(n)
		case 5:
			for k := 0; k < n; k++ {
				img.repeat()
			}
		default:
			return
		}
	}
}

// skip adds empty rows (index 0) and zeroes the seed rows.
func (img *raster) skip(n int) {
	for p := range img.cur {
		img.cur[p] = nil
		img.seed[p] = nil
	}
	img.plane = 0
	for k := 0; k < n && !img.full(); k++ {
		img.addRow(nil)
	}
}

func (img *raster) full() bool {
	if img.acc.full() {
		img.r.c.warnOnce("rastersize", "a raster image is cut off: it is too large")
		return true
	}
	return img.height > 0 && img.rows >= img.height
}

// repeat adds the last row again.
func (img *raster) repeat() {
	if !img.full() {
		img.addRow(img.prev)
	}
}

// emit turns the planes of the row received into pixels and adds them.
func (img *raster) emit() {
	defer func() { img.plane = 0 }()
	if img.full() {
		return
	}
	if img.bad {
		img.addRow(nil)
		return
	}
	r := img.r
	width := img.width
	if width <= 0 {
		for _, p := range img.cur {
			bpp := max(1, img.bpp)
			width = max(width, len(p)*8/bpp)
		}
		width = min(width, maxRasterWidth)
	}
	row := make([]color.NRGBA, width)
	white := color.NRGBA{255, 255, 255, 255}
	bit := func(p []byte, i int) int {
		if i/8 >= len(p) {
			return 0
		}
		return int(p[i/8] >> (7 - i%8) & 1)
	}
	bits := func(p []byte, at, n int) int {
		v := 0
		for k := 0; k < n; k++ {
			v = v<<1 | bit(p, at+k)
		}
		return v
	}
	for x := 0; x < width; x++ {
		var c color.NRGBA
		switch r.encoding {
		case 1: // indexed by pixel
			c = r.indexColor(bits(img.cur[0], x*img.bpp, img.bpp))
		case 3: // direct by pixel
			at := x * img.bpp
			var v [3]uint8
			for k := 0; k < 3; k++ {
				v[k] = r.primary(float64(bits(img.cur[0], at, r.bits[k])), k)
				at += r.bits[k]
			}
			c = color.NRGBA{v[0], v[1], v[2], 255}
		case 2: // direct by plane
			at := 0
			var v [3]uint8
			for k := 0; k < 3; k++ {
				n := 0
				for b := 0; b < r.bits[k]; b++ {
					n = n<<1 | bit(img.cur[at], x)
					at++
				}
				v[k] = r.primary(float64(n), k)
			}
			c = color.NRGBA{v[0], v[1], v[2], 255}
		default: // indexed by plane: plane 1 is the least significant bit
			idx := 0
			for p := range img.cur {
				idx |= bit(img.cur[p], x) << p
			}
			c = r.indexColor(idx)
		}
		if r.transparent && c == white {
			c = color.NRGBA{}
		}
		row[x] = c
	}
	img.addRow(row)
}

func (r *rtl) indexColor(i int) color.NRGBA {
	if i < 0 || i >= len(r.palette) {
		return color.NRGBA{0, 0, 0, 255}
	}
	c := r.palette[i]
	return color.NRGBA{uint8(c >> 24), uint8(c >> 16), uint8(c >> 8), 255}
}

func (img *raster) addRow(row []color.NRGBA) {
	img.prev = row
	img.acc.add(row)
	img.rows++
}

// --- scaling down as the rows arrive ---

// accumulator averages blocks of factor × factor pixels into the stored
// image: images finer than the resolution cap are scaled down as they
// arrive, so that a plot of a metre of paper at 600 dpi is never held
// whole. Two-color images stay two-colored.
type accumulator struct {
	factor  int
	bilevel bool // a two-color palette image
	ink     color.NRGBA
	paper   color.NRGBA
	inW     int // widest row (source pixels)
	band    int // source rows in the current output row
	sum     []uint32
	out     [][]color.NRGBA
	outBi   [][]byte
}

// maxRasterBytes bounds the memory of a stored image before it is scaled
// down to the resolution cap.
const maxRasterBytes = 256 << 20

func (a *accumulator) init(r *rtl, img *raster) {
	res := 1016 / math.Min(img.pw, img.ph)
	a.factor = 1
	if dpi := r.c.opts.Images.MaxDPI; dpi >= 0 {
		if dpi == 0 {
			dpi = imgconv.DefaultMaxDPI
		}
		a.factor = max(1, int(math.Floor(res/dpi)))
	}
	if img.width > 0 && img.height > 0 {
		for int64(img.width/a.factor+1)*int64(img.height/a.factor+1)*4 > maxRasterBytes {
			a.factor++
		}
	}
	a.bilevel = r.encoding == 0 && img.planes == 1
	if a.bilevel {
		a.paper = r.indexColor(0)
		a.ink = r.indexColor(1)
		if r.transparent && a.paper == (color.NRGBA{255, 255, 255, 255}) {
			a.paper = color.NRGBA{}
		}
	}
}

func (a *accumulator) add(row []color.NRGBA) {
	f := a.factor
	if len(row) > a.inW {
		a.inW = len(row)
	}
	ow := (a.inW + f - 1) / f
	if a.bilevel {
		for len(a.sum) < ow {
			a.sum = append(a.sum, 0)
		}
		for x, c := range row {
			if c.A != 0 && c == a.ink {
				a.sum[x/f]++
			}
		}
	} else {
		for len(a.sum) < 4*ow {
			a.sum = append(a.sum, 0)
		}
		for x, c := range row {
			o := 4 * (x / f)
			al := uint32(c.A)
			a.sum[o] += uint32(c.R) * al
			a.sum[o+1] += uint32(c.G) * al
			a.sum[o+2] += uint32(c.B) * al
			a.sum[o+3] += al
		}
	}
	a.band++
	if a.band == f {
		a.flushBand()
	}
}

// full reports whether the image holds as many pixels as it may.
func (a *accumulator) full() bool {
	ow := (a.inW + a.factor - 1) / a.factor
	return int64(len(a.out)+len(a.outBi))*int64(ow)*4 > maxRasterBytes
}

// flushBand turns the sums of a band of rows into an output row.
func (a *accumulator) flushBand() {
	if a.band == 0 {
		return
	}
	f := a.factor
	area := uint32(f * a.band)
	if a.bilevel {
		out := make([]byte, len(a.sum))
		for i, s := range a.sum {
			// ink covering two fifths of the block is ink
			if s > 0 && s*5 >= 2*area {
				out[i] = 1
			}
		}
		a.outBi = append(a.outBi, out)
	} else {
		out := make([]color.NRGBA, len(a.sum)/4)
		for i := range out {
			al := a.sum[4*i+3]
			if al == 0 {
				continue
			}
			out[i] = color.NRGBA{uint8(a.sum[4*i] / al), uint8(a.sum[4*i+1] / al), uint8(a.sum[4*i+2] / al), uint8(al / area)}
		}
		a.out = append(a.out, out)
	}
	clear(a.sum)
	a.band = 0
}

// image returns the stored image (nil when it is empty).
func (a *accumulator) image() image.Image {
	a.flushBand()
	f := a.factor
	w := (a.inW + f - 1) / f
	if a.bilevel {
		h := len(a.outBi)
		if w == 0 || h == 0 {
			return nil
		}
		img := image.NewPaletted(image.Rect(0, 0, w, h), color.Palette{a.paper, a.ink})
		ink := false
		for y, row := range a.outBi {
			for x, v := range row {
				if v != 0 {
					img.Pix[y*img.Stride+x] = 1
					ink = true
				}
			}
		}
		if !ink && a.paper.A == 0 {
			return nil
		}
		return img
	}
	h := len(a.out)
	if w == 0 || h == 0 {
		return nil
	}
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	visible := false
	for y, row := range a.out {
		for x, c := range row {
			img.SetNRGBA(x, y, c)
			if c.A != 0 {
				visible = true
			}
		}
	}
	if !visible {
		return nil
	}
	return img
}
