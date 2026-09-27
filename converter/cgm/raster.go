package cgm

import (
	"bytes"
	"image"
	"image/color"
	"math"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/tiff"
	"github.com/shibukawa/bdf/imgconv"
)

// bitReader reads the packed colour lists of the binary encoding.
type bitReader struct {
	b   []byte
	pos int // bits
}

func (r *bitReader) read(n int) (uint64, bool) {
	if n <= 0 || n > 64 || r.pos+n > 8*len(r.b) {
		return 0, false
	}
	var v uint64
	for n > 0 {
		i, off := r.pos/8, r.pos%8
		take := min(8-off, n)
		bits := uint64(r.b[i]>>(8-off-take)) & (1<<take - 1)
		v = v<<take | bits
		r.pos += take
		n -= take
	}
	return v, true
}

// cells reads the nx × ny colours of a cell array or a pattern table,
// rows first, and passes each to emit. In the binary encoding they are
// packed at the local colour precision (0 is the colour precision of the
// metafile), each row starting on a word boundary, or (mode 0, cell
// arrays) run-length encoded as counts and colours; in clear text they are
// listed. Cells the element lacks are colour 0.
func (in *interp) cells(p params, nx, ny, prec, mode, limit int, emit func(x, y int, c colour)) bool {
	if nx <= 0 || ny <= 0 || nx > limit/ny {
		in.warnOnce("cells", "cell arrays and patterns of more than %d cells or none are not drawn", limit)
		return false
	}
	direct := in.st.colourMode == 1
	bp, ok := p.(*binParams)
	if !ok {
		for y := 0; y < ny; y++ {
			for x := 0; x < nx; x++ {
				var c colour
				if p.more() {
					c = in.readColour(p)
				}
				emit(x, y, c)
			}
		}
		return true
	}
	bits := prec
	if bits <= 0 {
		bits = in.pr.colourIndexBits
		if direct {
			bits = in.pr.colourBits
		}
	}
	if bits > 32 {
		in.warnOnce("cell precision", "cell arrays of %d-bit colours are not drawn", bits)
		return false
	}
	comps := 1
	if direct {
		comps = in.pr.components
	}
	cell := func(r *bitReader) (colour, bool) {
		if !direct {
			v, ok := r.read(bits)
			return colour{index: int(v)}, ok
		}
		var v [4]float64
		for i := 0; i < comps; i++ {
			c, ok := r.read(bits)
			if !ok {
				return colour{}, false
			}
			v[i] = float64(c)
		}
		max := math.Exp2(float64(bits)) - 1
		if bits != in.pr.colourBits {
			// the components rescaled to the colour precision the extent is in
			f := bp.colourMax() / max
			for i := range v {
				v[i] *= f
			}
			max = bp.colourMax()
		}
		return colour{direct: true, c: in.rgb(v, comps, max)}, true
	}
	data := bp.rest()
	short := false
	if mode == 0 {
		r := &bitReader{b: data}
		for y := 0; y < ny; y++ {
			x := 0
			for x < nx {
				n, ok := r.read(in.pr.intBits)
				if !ok {
					break
				}
				c, ok := cell(r)
				if !ok {
					break
				}
				for ; n > 0 && x < nx; n-- {
					emit(x, y, c)
					x++
				}
			}
			if x < nx {
				short = true
			}
			for ; x < nx; x++ {
				emit(x, y, colour{})
			}
			r.pos = (r.pos + 15) &^ 15
		}
	} else {
		rowBits := nx * bits * comps
		row := (rowBits + 7) / 8
		// rows start on word boundaries; some writers pack them tighter or
		// align them to 4 bytes
		stride := (row + 1) &^ 1
		for _, s := range []int{(row + 1) &^ 1, row, (row + 3) &^ 3} {
			if s*ny == len(data) || s*(ny-1)+row == len(data) {
				stride = s
				break
			}
		}
		short = stride*(ny-1)+row > len(data)
		for y := 0; y < ny; y++ {
			r := &bitReader{}
			if start := y * stride; start < len(data) {
				r.b = data[start:min(start+row, len(data))]
			}
			for x := 0; x < nx; x++ {
				c, _ := cell(r)
				emit(x, y, c)
			}
		}
	}
	if short {
		in.warnOnce("short cells", "cell arrays with fewer colours than cells are drawn with the missing cells in colour 0")
	}
	return true
}

// useCells reports whether a cell array or tile of nx × ny cells is drawn:
// not when it has no cells or too many, or the metafile has used up its
// cells.
func (in *interp) useCells(nx, ny int) bool {
	if nx <= 0 || ny <= 0 || nx > maxCells/ny {
		in.warnOnce("cells", "cell arrays and tiles of more than %d cells or none are not drawn", maxCells)
		return false
	}
	if in.cellsUsed += nx * ny; in.cellsUsed > maxTotalCells {
		in.warnOnce("total cells", "cell arrays and tiles beyond %d cells in all are not drawn", maxTotalCells)
		return false
	}
	return true
}

// cellImage stores an image of nx × ny cells, which fill sets, with the
// transparent cell colour transparent.
func (in *interp) cellImage(nx, ny int, fill func(set func(x, y int, c colour)) bool) (bdf.Hash, bool) {
	img := image.NewNRGBA(image.Rect(0, 0, nx, ny))
	tc := in.st.transparentCell
	ok := fill(func(x, y int, c colour) {
		if tc != nil && c == *tc {
			return
		}
		img.SetNRGBA(x, y, toColor(in.colourOf(c)|0xff))
	})
	if !ok {
		return bdf.Hash{}, false
	}
	res, _ := imgconv.EncodeImage(img, true, in.c.opts.Images)
	return in.c.doc.AddImage(res.Data), true
}

// cellArray draws a CELL ARRAY: rows of cells in the parallelogram whose
// first row runs from P to R, with Q the corner diagonal to P.
func (in *interp) cellArray(p params) {
	pp, q, r := p.point(), p.point(), p.point()
	nx, ny := p.int(), p.int()
	prec := p.int()
	mode := 1
	if _, ok := p.(*binParams); ok {
		mode = p.enum("")
	}
	if !in.useCells(nx, ny) {
		return
	}
	h, ok := in.cellImage(nx, ny, func(set func(x, y int, c colour)) bool {
		return in.cells(p, nx, ny, prec, mode, maxCells, set)
	})
	if !ok {
		return
	}
	fx, fy := float64(nx), float64(ny)
	m := canvas.Matrix{(r.X - pp.X) / fx, (r.Y - pp.Y) / fx, (q.X - r.X) / fy, (q.Y - r.Y) / fy, pp.X, pp.Y}
	in.out.Image(&cad.Image{Image: h, W: nx, H: ny, M: m})
}

// patternFill fills a path with a pattern of PATTERN TABLE, repeated
// every PATTERN SIZE from the fill reference point.
func (in *interp) patternFill(path *cad.Path, idx int, c bdf.Color) {
	def, ok := in.st.patterns[idx]
	if !ok {
		in.warnOnce("pattern:"+itoa(idx), "pattern %d is not defined; it is filled with the fill colour", idx)
		in.out.Fill(path, cad.Fill{Color: c}, true)
		return
	}
	key := [2]int{idx, in.generation}
	h, ok := in.patternImages[key]
	if !ok {
		h, _ = in.cellImage(def.nx, def.ny, func(set func(x, y int, c colour)) bool {
			for i, c := range def.cells {
				set(i%def.nx, i/def.nx, c)
			}
			return true
		})
		in.patternImages[key] = h
	}
	a := &in.st.a
	hv := cad.Point{Y: defaultPattern * in.vdcPerMM()}
	wv := cad.Point{X: defaultPattern * in.vdcPerMM()}
	if a.patternSizeSet {
		hv = cad.Point{X: a.patternSize[0], Y: a.patternSize[1]}
		wv = cad.Point{X: a.patternSize[2], Y: a.patternSize[3]}
	}
	if math.Abs(wv.X*hv.Y-wv.Y*hv.X) < 1e-12 {
		in.out.Fill(path, cad.Fill{Color: c}, true)
		return
	}
	ref := in.fillRef()
	fx, fy := float64(def.nx), float64(def.ny)
	m := canvas.Matrix{wv.X / fx, wv.Y / fx, -hv.X / fy, -hv.Y / fy, ref.X + hv.X, ref.Y + hv.Y}
	in.out.Fill(path, cad.Fill{Color: c, Pattern: &cad.Pattern{Image: h, M: m}}, true)
}

// tileArray is the geometry of a tile array (BEGIN TILE ARRAY): the tiles
// that follow fill it along the path direction, then line after line.
type tileArray struct {
	at                 cad.Point
	path, line         cad.Point // a cell along the path and the line progression
	tilesPath          int
	cellsPath, cellsLn int // per tile
	offPath, offLine   int
	n                  int // tiles drawn
	clip               bool
}

// beginTileArray reads BEGIN TILE ARRAY.
func (in *interp) beginTileArray(p params) {
	at := p.point()
	dir := p.enum("0 90 180 270")
	prog := p.enum("90 270")
	tp, tl := p.int(), p.int()
	cp, cl := p.int(), p.int()
	sp, sl := p.real(), p.real()
	op, ol := p.int(), p.int()
	np, nl := p.int(), p.int()
	_ = tl
	ang := float64(dir) * math.Pi / 2
	u := cad.Point{X: math.Cos(ang), Y: math.Sin(ang)}
	v := cad.Point{X: -u.Y, Y: u.X}
	if prog == 1 {
		v = v.Mul(-1)
	}
	t := &tileArray{at: at, path: u.Mul(sp), line: v.Mul(sl), tilesPath: max(tp, 1), cellsPath: cp, cellsLn: cl, offPath: op, offLine: ol}
	in.tiles = t
	if !in.drawing() || np <= 0 || nl <= 0 {
		return
	}
	// the displayed cells
	a, b := t.path.Mul(float64(np)), t.line.Mul(float64(nl))
	clip := (&cad.Path{}).Polyline([]cad.Point{at, at.Add(a), at.Add(a).Add(b), at.Add(b)}, true)
	in.out.Begin(clip)
	t.clip = true
}

func (in *interp) endTileArray() {
	if t := in.tiles; t != nil && t.clip {
		in.out.End()
	}
	in.tiles = nil
}

// tile draws a TILE or BITONAL TILE of the tile array: JPEG and PNG tiles
// as they are, CCITT fax (T.4, T.6) and uncompressed tiles decoded.
func (in *interp) tile(c code, p params) {
	t := in.tiles
	if t == nil {
		in.warnOnce("tile", "tiles outside a tile array are not drawn")
		return
	}
	i, j := t.n%t.tilesPath, t.n/t.tilesPath
	t.n++
	comp := p.index()
	pad := p.int()
	var fg, bg colour
	bits := 1
	if c == eBitonalTile {
		bg, fg = in.readColour(p), in.readColour(p)
	} else {
		bits = p.int()
		if bits <= 0 {
			bits = in.pr.colourBits
			if in.st.colourMode == 0 {
				bits = in.pr.colourIndexBits
			}
		}
	}
	bp, ok := p.(*binParams)
	if !in.drawing() || !ok {
		return
	}
	bp.string() // the method-specific parameters (an SDR)
	data := bp.rest()
	w, h := t.cellsPath, t.cellsLn
	var img image.Image
	var stored []byte
	switch comp {
	case 0, 1: // the whole tile in the background or foreground colour
		col := bg
		if comp == 1 {
			col = fg
		}
		o := t.at.Add(t.path.Mul(float64(i*w - t.offPath))).Add(t.line.Mul(float64(j*h - t.offLine)))
		quad := (&cad.Path{}).Polyline([]cad.Point{o, o.Add(t.path.Mul(float64(w))), o.Add(t.path.Mul(float64(w))).Add(t.line.Mul(float64(h))), o.Add(t.line.Mul(float64(h)))}, true)
		in.out.Fill(quad, cad.Fill{Color: in.colourOf(col)}, false)
		return
	case 2, 3, 4: // T.6, T.4 1-D, T.4 2-D
		if !in.useCells(w, h) {
			return
		}
		rows, damaged := tiff.DecodeFax(data, comp == 2, comp == 4, w, h)
		if damaged {
			in.warnOnce("fax", "damaged fax tiles are drawn with the damaged rows blank")
		}
		img = bilevel(rows, w, h, in.colourOf(bg), in.colourOf(fg))
	case 5: // uncompressed
		if !in.useCells(w, h) {
			return
		}
		img = in.bitmapTile(data, w, h, bits, pad, c == eBitonalTile, bg, fg)
	case 7, 9: // baseline JPEG, PNG
		stored = data
		if cfg, _, err := image.DecodeConfig(bytes.NewReader(data)); err == nil {
			w, h = cfg.Width, cfg.Height
		}
	default:
		in.warnOnce("tile:"+itoa(comp), "tiles of compression type %d are not drawn", comp)
		return
	}
	var hash bdf.Hash
	if stored != nil {
		res, err := imgconv.Optimize(stored, in.c.opts.Images)
		if err != nil {
			in.warnOnce("tile image", "tiles whose images could not be read are not drawn")
			return
		}
		hash = in.c.doc.AddImage(res.Data)
	} else {
		if img == nil {
			return
		}
		res, _ := imgconv.EncodePixels(img, true, in.c.opts.Images)
		hash = in.c.doc.AddImage(res.Data)
	}
	if w <= 0 || h <= 0 {
		return
	}
	o := t.at.Add(t.path.Mul(float64(i*t.cellsPath - t.offPath))).Add(t.line.Mul(float64(j*t.cellsLn - t.offLine)))
	// the image fills the tile's cells
	sx, sy := float64(t.cellsPath)/float64(w), float64(t.cellsLn)/float64(h)
	m := canvas.Matrix{t.path.X * sx, t.path.Y * sx, t.line.X * sy, t.line.Y * sy, o.X, o.Y}
	in.out.Image(&cad.Image{Image: hash, W: w, H: h, M: m, Smooth: stored != nil})
}

// bilevel makes an image of packed rows of 1 bits (foreground) and 0 bits
// (background).
func bilevel(rows []byte, w, h int, bg, fg bdf.Color) image.Image {
	pal := color.Palette{toColor(bg), toColor(fg)}
	img := image.NewPaletted(image.Rect(0, 0, w, h), pal)
	rb := (w + 7) / 8
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if i := y*rb + x/8; i < len(rows) && rows[i]&(0x80>>(x%8)) != 0 {
				img.SetColorIndex(x, y, 1)
			}
		}
	}
	return img
}

func toColor(c bdf.Color) color.NRGBA {
	return color.NRGBA{uint8(c >> 24), uint8(c >> 16), uint8(c >> 8), uint8(c)}
}

// bitmapTile decodes an uncompressed tile: cells of bits each (colour
// indexes or direct colours, or one bit of a bitonal tile), rows padded
// to multiples of pad bits.
func (in *interp) bitmapTile(data []byte, w, h, bits, pad int, bitonal bool, bg, fg colour) image.Image {
	if w <= 0 || h <= 0 || w > maxCells/h || bits <= 0 || bits > 32 {
		return nil
	}
	direct := in.st.colourMode == 1 && !bitonal
	comps := 1
	if direct {
		comps = in.pr.components
	}
	rowBits := w * bits * comps
	if pad > 0 {
		rowBits = (rowBits + pad - 1) / pad * pad
	}
	if bitonal {
		rows := make([]byte, (w+7)/8*h)
		r := &bitReader{b: data}
		rb := (w + 7) / 8
		for y := 0; y < h; y++ {
			r.pos = y * rowBits
			for x := 0; x < w; x++ {
				if v, ok := r.read(1); ok && v == 1 {
					rows[y*rb+x/8] |= 0x80 >> (x % 8)
				}
			}
		}
		return bilevel(rows, w, h, in.colourOf(bg), in.colourOf(fg))
	}
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	r := &bitReader{b: data}
	max := math.Exp2(float64(bits)) - 1
	for y := 0; y < h; y++ {
		r.pos = y * rowBits
		for x := 0; x < w; x++ {
			var c bdf.Color
			if direct {
				var v [4]float64
				for k := 0; k < comps; k++ {
					u, _ := r.read(bits)
					v[k] = float64(u) / max
				}
				c = rgbUnit(v, comps)
			} else {
				u, _ := r.read(bits)
				c = in.colourOf(colour{index: int(u)})
			}
			img.SetNRGBA(x, y, toColor(c))
		}
	}
	return img
}
