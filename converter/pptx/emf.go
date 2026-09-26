package pptx

import (
	"math"
)

// stock objects (SelectObject with the high bit set)
func stockObject(i uint32) any {
	switch i & 0x7fffffff {
	case 0:
		return &gdiBrush{color: white, hatch: -1}
	case 1:
		return &gdiBrush{color: rgba{0.75, 0.75, 0.75, 1}, hatch: -1}
	case 2:
		return &gdiBrush{color: rgba{0.5, 0.5, 0.5, 1}, hatch: -1}
	case 3:
		return &gdiBrush{color: rgba{0.25, 0.25, 0.25, 1}, hatch: -1}
	case 4:
		return &gdiBrush{color: black, hatch: -1}
	case 5:
		return &gdiBrush{null: true}
	case 6:
		return &gdiPen{color: white}
	case 7:
		return &gdiPen{color: black}
	case 8:
		return &gdiPen{null: true}
	case 10, 11, 16:
		return &gdiFont{height: 12, weight: 400, face: "Courier New"}
	case 12, 13, 14, 17:
		return &gdiFont{height: 12, weight: 400, face: "Arial"}
	case 18:
		return &gdiBrush{color: white, hatch: -1}
	case 19:
		return &gdiPen{color: black}
	}
	return nil
}

func (g *gdi) selectObject(o any) {
	switch v := o.(type) {
	case *gdiPen:
		g.st.pen = v
	case *gdiBrush:
		g.st.brush = v
	case *gdiFont:
		g.st.font = v
	}
}

func brushFromLog(style uint32, c rgba, hatch uint32) *gdiBrush {
	switch style {
	case 1: // BS_NULL
		return &gdiBrush{null: true}
	case 2: // BS_HATCHED
		return &gdiBrush{color: c, hatch: int(hatch % 6)}
	}
	return &gdiBrush{color: c, hatch: -1}
}

// patternBrush approximates a bitmap brush by the average color of the bitmap.
func patternBrush(bmi, bits []byte) *gdiBrush {
	img, err := decodeDIB(bmi, bits)
	if err != nil || img == nil {
		return &gdiBrush{color: rgba{0.5, 0.5, 0.5, 1}, hatch: -1}
	}
	var r, gg, b, n float64
	for i := 0; i+3 < len(img.Pix); i += 4 {
		r += float64(img.Pix[i])
		gg += float64(img.Pix[i+1])
		b += float64(img.Pix[i+2])
		n++
	}
	if n == 0 {
		return &gdiBrush{color: rgba{0.5, 0.5, 0.5, 1}, hatch: -1}
	}
	return &gdiBrush{color: rgba{r / n / 255, gg / n / 255, b / n / 255, 1}, hatch: -1}
}

func slice(b []byte, off, n uint32) []byte {
	if n == 0 || uint64(off)+uint64(n) > uint64(len(b)) {
		return nil
	}
	return b[off : off+n]
}

// playEMF replays an enhanced metafile into x,y,w,h: its frame (the
// picture's extent in 0.01 mm on the reference device) fills the rectangle.
func (g *gdi) playEMF(b []byte, x, y, w, h float64) {
	if len(b) < 88 {
		return
	}
	bl, bt, br, bb := lei32(b, 8), lei32(b, 12), lei32(b, 16), lei32(b, 20)
	fl, ft, fr, fb := lei32(b, 24), lei32(b, 28), lei32(b, 32), lei32(b, 36)
	devX, devY, mmX, mmY := lei32(b, 72), lei32(b, 76), lei32(b, 80), lei32(b, 84)
	pxmmX, pxmmY := 96/25.4, 96/25.4
	if devX > 0 && mmX > 0 && devY > 0 && mmY > 0 {
		pxmmX, pxmmY = devX/mmX, devY/mmY
	}
	g.pxmm = pxmmX
	// frame in device pixels
	x0, y0, x1, y1 := fl/100*pxmmX, ft/100*pxmmY, fr/100*pxmmX, fb/100*pxmmY
	if x1-x0 <= 0 || y1-y0 <= 0 {
		x0, y0, x1, y1 = bl, bt, br+1, bb+1
	}
	if x1-x0 <= 0 || y1-y0 <= 0 {
		return
	}
	g.out = translate(x, y).mul(scale(w/(x1-x0), h/(y1-y0))).mul(translate(-x0, -y0))
	pathFigure := false
	for p := 0; p+8 <= len(b); {
		typ, size := le32(b, p), int(le32(b, p+4))
		if size < 8 || p+size > len(b) {
			break
		}
		r := b[p : p+size]
		p += size
		pt := func(off int) [2]float64 { return [2]float64{lei32(r, off), lei32(r, off+4)} }
		points := func(off, n int, short bool) [][2]float64 {
			out := make([][2]float64, 0, n)
			for i := 0; i < n; i++ {
				if short {
					if off+i*4+4 > len(r) {
						break
					}
					out = append(out, [2]float64{lei16(r, off+i*4), lei16(r, off+i*4+2)})
				} else {
					if off+i*8+8 > len(r) {
						break
					}
					out = append(out, pt(off+i*8))
				}
			}
			return out
		}
		switch typ {
		case 14: // EOF
			return
		case 2, 3, 4, 85, 86, 87: // POLYBEZIER, POLYGON, POLYLINE (+16)
			short := typ >= 85
			pts := points(28, int(le32(r, 24)), short)
			if len(pts) == 0 {
				continue
			}
			kind := typ
			if short {
				kind -= 83
			}
			switch kind {
			case 2:
				g.draw(bezierPath(pts[0], pts[1:], true), false, true)
			case 3:
				g.draw(polyPath(pts, true), true, true)
			case 4:
				g.draw(polyPath(pts, false), false, true)
			}
			if g.path != nil {
				pathFigure = kind != 3
			}
		case 5, 6, 88, 89: // POLYBEZIERTO, POLYLINETO (+16)
			short := typ >= 88
			pts := points(28, int(le32(r, 24)), short)
			if len(pts) == 0 {
				continue
			}
			bezier := typ == 5 || typ == 88
			if g.path != nil {
				if !pathFigure {
					g.path.moveTo(g.st.cur[0], g.st.cur[1])
					pathFigure = true
				}
				if bezier {
					g.path.p.Verbs = append(g.path.p.Verbs, bezierPath(g.st.cur, pts, false).Verbs...)
					g.path.p.Args = append(g.path.p.Args, bezierPath(g.st.cur, pts, false).Args...)
				} else {
					for _, q := range pts {
						g.path.lineTo(q[0], q[1])
					}
				}
			} else if bezier {
				g.draw(bezierPath(g.st.cur, pts, true), false, true)
			} else {
				g.draw(polyPath(append([][2]float64{g.st.cur}, pts...), false), false, true)
			}
			g.st.cur = pts[len(pts)-1]
		case 7, 8, 90, 91: // POLYPOLYLINE, POLYPOLYGON (+16)
			short := typ >= 90
			n, total := int(le32(r, 24)), int(le32(r, 28))
			if n <= 0 || n > len(r) {
				continue
			}
			counts := make([]int, n)
			for i := range counts {
				counts[i] = int(le32(r, 32+i*4))
			}
			all := points(32+n*4, total, short)
			closed := typ == 8 || typ == 91
			p2 := newPath()
			k := 0
			for _, c := range counts {
				for i := 0; i < c && k < len(all); i++ {
					if i == 0 {
						p2.moveTo(all[k][0], all[k][1])
					} else {
						p2.lineTo(all[k][0], all[k][1])
					}
					k++
				}
				if closed {
					p2.close()
				}
			}
			g.draw(p2.p, closed, true)
		case 9:
			g.st.winExt = [2]float64{lei32(r, 8), lei32(r, 12)}
		case 10:
			g.st.winOrg = pt(8)
		case 11:
			g.st.vpExt = [2]float64{lei32(r, 8), lei32(r, 12)}
		case 12:
			g.st.vpOrg = pt(8)
		case 17:
			g.st.mapMode = int(le32(r, 8))
		case 18:
			g.st.bkOpaque = le32(r, 8) == 2
		case 19:
			g.st.winding = le32(r, 8) == 2
		case 22:
			g.st.textAlign = le32(r, 8)
		case 24:
			g.st.textColor = colorRef(le32(r, 8))
		case 25:
			g.st.bkColor = colorRef(le32(r, 8))
		case 27: // MOVETOEX
			g.st.cur = pt(8)
			if g.path != nil {
				g.path.moveTo(g.st.cur[0], g.st.cur[1])
				pathFigure = true
			}
		case 30: // INTERSECTCLIPRECT
			g.clipRect(lei32(r, 8), lei32(r, 12), lei32(r, 16), lei32(r, 20))
		case 33:
			g.save()
		case 34:
			g.restore(int(int32(le32(r, 8))))
		case 35: // SETWORLDTRANSFORM
			g.st.world = matrix{lef32(r, 8), lef32(r, 12), lef32(r, 16), lef32(r, 20), lef32(r, 24), lef32(r, 28)}
		case 36: // MODIFYWORLDTRANSFORM
			x := matrix{lef32(r, 8), lef32(r, 12), lef32(r, 16), lef32(r, 20), lef32(r, 24), lef32(r, 28)}
			switch le32(r, 32) {
			case 1:
				g.st.world = identity
			case 2:
				g.st.world = g.st.world.mul(x)
			case 3:
				g.st.world = x.mul(g.st.world)
			case 4:
				g.st.world = x
			}
		case 37: // SELECTOBJECT
			h := le32(r, 8)
			if h&0x80000000 != 0 {
				g.selectObject(stockObject(h))
			} else {
				g.selectObject(g.objs[h])
			}
		case 38: // CREATEPEN
			style := le32(r, 12)
			wd := lei32(r, 16)
			g.objs[le32(r, 8)] = &gdiPen{null: style&0xf == 5, width: wd, cosmetic: wd == 0, color: colorRef(le32(r, 24)), style: style}
		case 95: // EXTCREATEPEN
			style := le32(r, 28)
			wd := float64(le32(r, 32))
			geometric := style&0x10000 != 0
			pen := &gdiPen{null: style&0xf == 5, width: wd, cosmetic: !geometric, color: colorRef(le32(r, 40)), style: style}
			if le32(r, 36) == 1 { // BS_NULL
				pen.null = true
			}
			g.objs[le32(r, 8)] = pen
		case 39: // CREATEBRUSHINDIRECT
			g.objs[le32(r, 8)] = brushFromLog(le32(r, 12), colorRef(le32(r, 16)), le32(r, 20))
		case 93, 94: // CREATEMONOBRUSH, CREATEDIBPATTERNBRUSHPT
			g.objs[le32(r, 8)] = patternBrush(slice(r, le32(r, 16), le32(r, 20)), slice(r, le32(r, 24), le32(r, 28)))
		case 40: // DELETEOBJECT
			delete(g.objs, le32(r, 8))
		case 82: // EXTCREATEFONTINDIRECTW
			f := &gdiFont{height: lei32(r, 12), escapement: lei32(r, 20), weight: int(lei32(r, 28)),
				italic: len(r) > 32 && r[32] != 0, underline: len(r) > 33 && r[33] != 0, strike: len(r) > 34 && r[34] != 0}
			if len(r) >= 40 {
				f.face = string(utf16String(r[40:], 32))
			}
			g.objs[le32(r, 8)] = f
		case 42:
			g.ellipse(lei32(r, 8), lei32(r, 12), lei32(r, 16), lei32(r, 20))
		case 43:
			g.rect(lei32(r, 8), lei32(r, 12), lei32(r, 16), lei32(r, 20))
		case 44:
			g.roundRect(lei32(r, 8), lei32(r, 12), lei32(r, 16), lei32(r, 20), lei32(r, 24), lei32(r, 28))
		case 45, 46, 47: // ARC, CHORD, PIE
			g.arc(int(typ-45), lei32(r, 8), lei32(r, 12), lei32(r, 16), lei32(r, 20), lei32(r, 24), lei32(r, 28), lei32(r, 32), lei32(r, 36))
		case 54: // LINETO
			q := pt(8)
			if g.path != nil {
				if !pathFigure {
					g.path.moveTo(g.st.cur[0], g.st.cur[1])
					pathFigure = true
				}
				g.path.lineTo(q[0], q[1])
			} else {
				g.draw(polyPath([][2]float64{g.st.cur, q}, false), false, true)
			}
			g.st.cur = q
		case 59: // BEGINPATH
			g.path = newPath()
			pathFigure = false
		case 60: // ENDPATH
		case 61: // CLOSEFIGURE
			if g.path != nil {
				g.path.close()
				pathFigure = false
			}
		case 62, 63, 64: // FILLPATH, STROKEANDFILLPATH, STROKEPATH
			if g.path != nil {
				p2 := g.path.p
				g.path = nil
				g.draw(p2, typ != 64, typ != 62)
			}
		case 67: // SELECTCLIPPATH
			if g.path != nil {
				m := g.matrix()
				var poly [][2]float64
				a := 0
				for _, v := range g.path.p.Verbs {
					n := [...]int{2, 2, 4, 6, 0, 4, 8, 5, 5}[v]
					if v <= 3 && n >= 2 {
						x, y := m.apply(float64(g.path.p.Args[a+n-2]), float64(g.path.p.Args[a+n-1]))
						poly = append(poly, [2]float64{x, y})
					}
					a += n
				}
				g.path = nil
				if len(poly) >= 3 {
					if le32(r, 8) == 5 {
						g.st.clips = nil
					}
					g.st.clips = append(append([][][2]float64(nil), g.st.clips...), poly)
					g.clipV++
				}
			}
		case 68: // ABORTPATH
			g.path = nil
		case 75: // EXTSELECTCLIPRGN
			cb, mode := le32(r, 8), le32(r, 12)
			if cb == 0 {
				if mode == 5 {
					g.resetClip()
				}
				continue
			}
			rgn := slice(r, 16, cb)
			if len(rgn) < 32 {
				continue
			}
			n := int(le32(rgn, 8))
			var polys [][2]float64
			for i := 0; i < n && 32+i*16+16 <= len(rgn); i++ {
				o := 32 + i*16
				l, t, rr, bb := lei32(rgn, o), lei32(rgn, o+4), lei32(rgn, o+8), lei32(rgn, o+12)
				// region rectangles are in device units
				for _, c := range [][2]float64{{l, t}, {rr, t}, {rr, bb}, {l, bb}} {
					x, y := g.out.apply(c[0], c[1])
					polys = append(polys, [2]float64{x, y})
				}
			}
			if n == 1 && len(polys) == 4 && (mode == 5 || mode == 1 || mode == 2) {
				if mode != 1 {
					g.st.clips = nil
				}
				g.st.clips = append(append([][][2]float64(nil), g.st.clips...), polys)
				g.clipV++
			} else if mode == 5 {
				// several rectangles: keep their bounding box
				minX, minY, maxX, maxY := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
				for _, c := range polys {
					minX, minY = math.Min(minX, c[0]), math.Min(minY, c[1])
					maxX, maxY = math.Max(maxX, c[0]), math.Max(maxY, c[1])
				}
				g.st.clips = [][][2]float64{{{minX, minY}, {maxX, minY}, {maxX, maxY}, {minX, maxY}}}
				g.clipV++
			}
		case 76, 77: // BITBLT, STRETCHBLT
			dx, dy, dw, dh := lei32(r, 24), lei32(r, 28), lei32(r, 32), lei32(r, 36)
			rop := le32(r, 40)
			sx, sy := lei32(r, 44), lei32(r, 48)
			bmi := slice(r, le32(r, 84), le32(r, 88))
			bits := slice(r, le32(r, 92), le32(r, 96))
			sw, sh := dw, dh
			if typ == 77 {
				sw, sh = lei32(r, 100), lei32(r, 104)
			}
			if bmi == nil {
				g.patBlt(rop, dx, dy, dw, dh)
				continue
			}
			g.dib(bmi, bits, sx, sy, sw, sh, dx, dy, dw, dh, false)
		case 81: // STRETCHDIBITS
			dx, dy := lei32(r, 24), lei32(r, 28)
			sx, sy, sw, sh := lei32(r, 32), lei32(r, 36), lei32(r, 40), lei32(r, 44)
			bmi := slice(r, le32(r, 48), le32(r, 52))
			bits := slice(r, le32(r, 56), le32(r, 60))
			dw, dh := lei32(r, 72), lei32(r, 76)
			if bmi != nil {
				g.dib(bmi, bits, sx, sy, sw, sh, dx, dy, dw, dh, true)
			}
		case 80: // SETDIBITSTODEVICE
			dx, dy := lei32(r, 24), lei32(r, 28)
			sx, sy, sw, sh := lei32(r, 32), lei32(r, 36), lei32(r, 40), lei32(r, 44)
			bmi := slice(r, le32(r, 48), le32(r, 52))
			bits := slice(r, le32(r, 56), le32(r, 60))
			if bmi != nil {
				g.dib(bmi, bits, sx, sy, sw, sh, dx, dy, sw, sh, true)
			}
		case 83, 84: // EXTTEXTOUTA, EXTTEXTOUTW
			if len(r) < 76 {
				continue
			}
			ref := pt(36)
			n := int(le32(r, 44))
			offStr, opts, offDx := int(le32(r, 48)), le32(r, 52), int(le32(r, 72))
			if opts&0x10 != 0 { // ETO_GLYPH_INDEX
				g.warn("glyphidx", "metafile text drawn by glyph index is not supported")
				continue
			}
			if n <= 0 || offStr >= len(r) {
				continue
			}
			var s []rune
			if typ == 84 {
				s = utf16String(r[offStr:], n)
			} else {
				end := min(len(r), offStr+n)
				s = []rune(decodeANSI(r[offStr:end], 0))
			}
			var dxs []float64
			step := 4
			if opts&0x2000 != 0 { // ETO_PDY
				step = 8
			}
			if offDx > 0 {
				for i := 0; i < n && offDx+i*step+4 <= len(r); i++ {
					dxs = append(dxs, lei32(r, offDx+i*step))
				}
				if len(s) != n {
					dxs = nil
				}
			}
			var opaque *[4]float64
			if opts&0x2 != 0 { // ETO_OPAQUE
				opaque = &[4]float64{lei32(r, 56), lei32(r, 60), lei32(r, 64), lei32(r, 68)}
			}
			g.text(ref[0], ref[1], s, dxs, opaque)
		case 70: // GDICOMMENT: EMF+ records are skipped
			if len(r) >= 20 && string(r[12:16]) == "EMF+" && le16(r, 16) == 0x4001 && le16(r, 18)&1 == 0 {
				// an EMF+ header without the dual flag: nothing else is drawn
				g.warn("emfplus", "EMF+ only metafiles are not supported")
			}
		}
	}
}

// patBlt fills a rectangle with the brush (or black / white).
func (g *gdi) patBlt(rop uint32, x, y, w, h float64) {
	save, pen := g.st.brush, g.st.pen
	switch rop {
	case 0x00F00021, 0x005A0049: // PATCOPY, PATINVERT
	case 0x00000042: // BLACKNESS
		g.st.brush = &gdiBrush{color: black, hatch: -1}
	case 0x00FF0062: // WHITENESS
		g.st.brush = &gdiBrush{color: white, hatch: -1}
	default:
		return
	}
	g.st.pen = &gdiPen{null: true}
	g.rect(x, y, x+w, y+h)
	g.st.brush, g.st.pen = save, pen
}

// dib draws a bitmap; bottomUpSrc converts a source rectangle measured
// from the bottom of a bottom-up DIB (StretchDIBits) to image rows.
func (g *gdi) dib(bmi, bits []byte, sx, sy, sw, sh, dx, dy, dw, dh float64, bottomUpSrc bool) {
	if bottomUpSrc && len(bmi) >= 12 {
		var h float64
		if le32(bmi, 0) == 12 {
			h = lei16(bmi, 6)
		} else {
			h = lei32(bmi, 8)
		}
		if h > 0 {
			sy = h - sy - sh
		}
	}
	g.image(bmi, bits, sx, sy, sw, sh, dx, dy, dw, dh)
}
