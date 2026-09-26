package metafile

import (
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

// decodeANSI decodes text of a GDI font charset.
func decodeANSI(b []byte, charset byte) string {
	var enc encoding.Encoding
	switch charset {
	case 128: // SHIFTJIS_CHARSET
		enc = japanese.ShiftJIS
	case 129, 130: // HANGUL, JOHAB
		enc = korean.EUCKR
	case 134: // GB2312
		enc = simplifiedchinese.GBK
	case 136: // CHINESEBIG5
		enc = traditionalchinese.Big5
	case 161:
		enc = charmap.Windows1253
	case 162:
		enc = charmap.Windows1254
	case 177:
		enc = charmap.Windows1255
	case 178:
		enc = charmap.Windows1256
	case 186:
		enc = charmap.Windows1257
	case 204:
		enc = charmap.Windows1251
	case 238:
		enc = charmap.Windows1250
	case 2: // SYMBOL_CHARSET: the codes are the symbol font's
		r := make([]rune, len(b))
		for i, c := range b {
			r[i] = rune(c)
		}
		return string(r)
	default:
		enc = charmap.Windows1252
	}
	out, err := enc.NewDecoder().Bytes(b)
	if err != nil {
		return string(b)
	}
	return string(out)
}

// wmfFont carries the charset its text is encoded in.
type wmfFont struct {
	gdiFont
	charset byte
}

// wmfBox reads the picture box of a Windows metafile: the bounding box of
// its placeable header, or else the first window origin and extent. start
// is the offset of the first record and inch the logical units per inch
// (0 when the metafile does not say).
func wmfBox(b []byte) (bl, bt, br, bb float64, start int, inch float64, ok bool) {
	p := 0
	placeable := le32(b, 0) == 0x9AC6CDD7
	if placeable {
		bl, bt, br, bb = lei16(b, 6), lei16(b, 8), lei16(b, 10), lei16(b, 12)
		inch = float64(le16(b, 14))
		p = 22
	}
	if p+18 > len(b) {
		return
	}
	p += int(le16(b, p+2)) * 2
	start = p
	if !placeable || br == bl || bb == bt {
		// the first window origin and extent define the picture
		var org, ext [2]float64
		haveExt := false
		for q := start; q+6 <= len(b); {
			size := int(le32(b, q)) * 2
			fn := le16(b, q+4)
			if size < 6 || q+size > len(b) {
				break
			}
			switch fn {
			case 0x020B:
				org = [2]float64{lei16(b, q+8), lei16(b, q+6)}
			case 0x020C:
				if !haveExt {
					ext = [2]float64{lei16(b, q+8), lei16(b, q+6)}
					haveExt = true
				}
			}
			q += size
		}
		if !haveExt {
			ext = [2]float64{1000, 1000}
		}
		bl, bt, br, bb = org[0], org[1], org[0]+ext[0], org[1]+ext[1]
	}
	ok = br != bl && bb != bt
	return
}

// playWMF replays a Windows metafile into x,y,w,h: the bounding box of
// the placeable header (or the first window extent) fills the rectangle.
func (g *gdi) playWMF(b []byte, x, y, w, h float64) {
	bl, bt, br, bb, start, _, ok := wmfBox(b)
	if !ok {
		return
	}
	// device space = the bounding box coordinates
	g.st.mapMode = 8
	g.st.winOrg, g.st.vpOrg = [2]float64{bl, bt}, [2]float64{bl, bt}
	g.st.winExt, g.st.vpExt = [2]float64{br - bl, bb - bt}, [2]float64{br - bl, bb - bt}
	g.out = canvas.Translate(x, y).Mul(canvas.Scale(w/(br-bl), h/(bb-bt))).Mul(canvas.Translate(-bl, -bt))
	charset := byte(0)
	for p := start; p+6 <= len(b); {
		size := int(le32(b, p)) * 2
		fn := le16(b, p+4)
		if size < 6 || p+size > len(b) {
			break
		}
		r := b[p+6 : p+size] // parameters
		p += size
		arg := func(i int) float64 { return lei16(r, i*2) }
		addObj := func(o any) {
			for i, s := range g.slots {
				if s == nil {
					g.slots[i] = o
					return
				}
			}
			g.slots = append(g.slots, o)
		}
		switch fn {
		case 0x0000:
			return
		case 0x020B: // SETWINDOWORG
			g.st.winOrg = [2]float64{arg(1), arg(0)}
		case 0x020C: // SETWINDOWEXT
			g.st.winExt = [2]float64{arg(1), arg(0)}
		case 0x020D:
			g.st.vpOrg = [2]float64{arg(1), arg(0)}
		case 0x020E:
			g.st.vpExt = [2]float64{arg(1), arg(0)}
		case 0x0103: // SETMAPMODE: keep the anisotropic mapping of the box
		case 0x0102:
			g.st.bkOpaque = le16(r, 0) == 2
		case 0x0106:
			g.st.winding = le16(r, 0) == 2
		case 0x012E:
			g.st.textAlign = uint32(le16(r, 0))
		case 0x0209:
			g.st.textColor = colorRef(le32(r, 0))
		case 0x0201:
			g.st.bkColor = colorRef(le32(r, 0))
		case 0x001E:
			g.save()
		case 0x0127:
			g.restore(int(int16(le16(r, 0))))
		case 0x02FA: // CREATEPENINDIRECT
			style := uint32(le16(r, 0))
			wd := arg(1)
			addObj(&gdiPen{null: style&0xf == 5, width: wd, cosmetic: wd == 0, color: colorRef(le32(r, 6)), style: style})
		case 0x02FC: // CREATEBRUSHINDIRECT
			addObj(brushFromLog(uint32(le16(r, 0)), colorRef(le32(r, 2)), uint32(le16(r, 6))))
		case 0x02FB: // CREATEFONTINDIRECT
			f := &wmfFont{gdiFont: gdiFont{height: arg(0), escapement: arg(2), weight: int(arg(4))}}
			if len(r) >= 18 {
				f.italic, f.underline, f.strike, f.charset = r[10] != 0, r[11] != 0, r[12] != 0, r[13]
				name := r[18:]
				for i, c := range name {
					if c == 0 {
						name = name[:i]
						break
					}
				}
				f.face = decodeANSI(name, f.charset)
			}
			addObj(f)
		case 0x0142: // DIBCREATEPATTERNBRUSH
			if len(r) > 4 {
				addObj(patternBrush(r[4:], dibBits(r[4:])))
			} else {
				addObj(&gdiBrush{color: gray(128), hatch: -1})
			}
		case 0x01F9, 0x00F7, 0x06FF: // CREATEPATTERNBRUSH, CREATEPALETTE, CREATEREGION
			addObj(struct{}{})
		case 0x012D: // SELECTOBJECT
			if i := int(le16(r, 0)); i < len(g.slots) {
				switch o := g.slots[i].(type) {
				case *wmfFont:
					g.st.font = &o.gdiFont
					charset = o.charset
				default:
					g.selectObject(o)
				}
			}
		case 0x01F0: // DELETEOBJECT
			if i := int(le16(r, 0)); i < len(g.slots) {
				g.slots[i] = nil
			}
		case 0x0324, 0x0325: // POLYGON, POLYLINE
			n := int(le16(r, 0))
			var pts [][2]float64
			for i := 0; i < n && 2+i*4+4 <= len(r); i++ {
				pts = append(pts, [2]float64{lei16(r, 2+i*4), lei16(r, 4+i*4)})
			}
			if len(pts) > 0 {
				g.draw(polyPath(pts, fn == 0x0324), fn == 0x0324, true)
			}
		case 0x0538: // POLYPOLYGON
			n := int(le16(r, 0))
			o := 2 + n*2
			pp := newPath()
			for k := 0; k < n; k++ {
				c := int(le16(r, 2+k*2))
				for i := 0; i < c && o+4 <= len(r); i++ {
					if i == 0 {
						pp.moveTo(lei16(r, o), lei16(r, o+2))
					} else {
						pp.lineTo(lei16(r, o), lei16(r, o+2))
					}
					o += 4
				}
				pp.close()
			}
			g.draw(pp.p, true, true)
		case 0x041B: // RECTANGLE: bottom right top left
			g.rect(arg(3), arg(2), arg(1), arg(0))
		case 0x0418: // ELLIPSE
			g.ellipse(arg(3), arg(2), arg(1), arg(0))
		case 0x061C: // ROUNDRECT: height width bottom right top left
			g.roundRect(arg(5), arg(4), arg(3), arg(2), arg(1), arg(0))
		case 0x0817, 0x0830, 0x081A: // ARC, CHORD, PIE
			kind := map[uint16]int{0x0817: 0, 0x0830: 1, 0x081A: 2}[fn]
			g.arc(kind, arg(7), arg(6), arg(5), arg(4), arg(3), arg(2), arg(1), arg(0))
		case 0x0214: // MOVETO
			g.st.cur = [2]float64{arg(1), arg(0)}
		case 0x0213: // LINETO
			q := [2]float64{arg(1), arg(0)}
			g.draw(polyPath([][2]float64{g.st.cur, q}, false), false, true)
			g.st.cur = q
		case 0x0416: // INTERSECTCLIPRECT: bottom right top left
			g.clipRect(arg(3), arg(2), arg(1), arg(0))
		case 0x0922: // PATBLT: rop(2 words) height width y x
			g.patBlt(le32(r, 0), arg(5), arg(4), arg(3), arg(2))
		case 0x0521: // TEXTOUT
			n := int(le16(r, 0))
			o := 2 + (n+1)/2*2
			if 2+n > len(r) {
				continue
			}
			g.text(lei16(r, o+2), lei16(r, o), []rune(decodeANSI(r[2:2+n], charset)), nil, nil)
		case 0x0A32: // EXTTEXTOUT: y x count options [rect] string [dx]
			yy, xx, n, opts := arg(0), arg(1), int(le16(r, 4)), le16(r, 6)
			o := 8
			var opaque *[4]float64
			if opts&0x6 != 0 {
				if opts&0x2 != 0 {
					opaque = &[4]float64{lei16(r, 8), lei16(r, 10), lei16(r, 12), lei16(r, 14)}
				}
				o += 8
			}
			if n <= 0 || o+n > len(r) {
				continue
			}
			s := []rune(decodeANSI(r[o:o+n], charset))
			o += (n + 1) / 2 * 2
			var dxs []float64
			if len(s) == n && o+n*2 <= len(r) {
				for i := 0; i < n; i++ {
					dxs = append(dxs, lei16(r, o+i*2))
				}
			}
			g.text(xx, yy, s, dxs, opaque)
		case 0x0F43: // STRETCHDIB: rop(2) usage srcH srcW ySrc xSrc dstH dstW yDst xDst DIB
			if len(r) > 22 {
				bmi := r[22:]
				g.dib(bmi, dibBits(bmi), arg(6), arg(5), arg(4), arg(3), arg(10), arg(9), arg(8), arg(7), true)
			}
		case 0x0B41: // DIBSTRETCHBLT: rop(2) srcH srcW ySrc xSrc dstH dstW yDst xDst DIB
			if len(r) > 20 {
				bmi := r[20:]
				g.dib(bmi, dibBits(bmi), arg(5), arg(4), arg(3), arg(2), arg(9), arg(8), arg(7), arg(6), true)
			} else if len(r) >= 18 {
				g.patBlt(le32(r, 0), arg(8), arg(7), arg(6), arg(5))
			}
		case 0x0940: // DIBBITBLT: rop(2) ySrc xSrc h w yDst xDst DIB
			if len(r) > 16 {
				bmi := r[16:]
				g.dib(bmi, dibBits(bmi), arg(3), arg(2), arg(5), arg(4), arg(7), arg(6), arg(5), arg(4), true)
			}
		}
	}
}
