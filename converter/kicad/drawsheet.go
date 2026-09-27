package kicad

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"image"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// defaultSheet is the drawing sheet KiCad puts on a page when the project
// names none: a double frame 10 mm from the paper's edges with grid
// references every 50 mm, and a title block in the bottom right corner.
const defaultSheet = `(kicad_wks
 (setup (textsize 1.5 1.5) (linewidth 0.15) (textlinewidth 0.15)
  (left_margin 10) (right_margin 10) (top_margin 10) (bottom_margin 10))
 (rect (start 110 34) (end 2 2))
 (rect (start 0 0 ltcorner) (end 0 0) (repeat 2) (incrx 2) (incry 2))
 (line (start 50 2 ltcorner) (end 50 0 ltcorner) (repeat 30) (incrx 50))
 (tbtext "1" (pos 25 1 ltcorner) (font (size 1.3 1.3)) (repeat 100) (incrx 50))
 (line (start 50 2 lbcorner) (end 50 0 lbcorner) (repeat 30) (incrx 50))
 (tbtext "1" (pos 25 1 lbcorner) (font (size 1.3 1.3)) (repeat 100) (incrx 50))
 (line (start 0 50 ltcorner) (end 2 50 ltcorner) (repeat 30) (incry 50))
 (tbtext "A" (pos 1 25 ltcorner) (font (size 1.3 1.3)) (justify center) (repeat 100) (incry 50))
 (line (start 0 50 rtcorner) (end 2 50 rtcorner) (repeat 30) (incry 50))
 (tbtext "A" (pos 1 25 rtcorner) (font (size 1.3 1.3)) (justify center) (repeat 100) (incry 50))
 (tbtext "Date: ${ISSUE_DATE}" (pos 87 6.9))
 (line (start 110 5.5) (end 2 5.5))
 (tbtext "${KICAD_VERSION}" (pos 109 4.1))
 (line (start 110 8.5) (end 2 8.5))
 (tbtext "Rev: ${REVISION}" (pos 24 6.9) (font bold))
 (tbtext "Size: ${PAPER}" (pos 109 6.9))
 (tbtext "Id: ${#}/${##}" (pos 24 4.1))
 (line (start 110 12.5) (end 2 12.5))
 (tbtext "Title: ${TITLE}" (pos 109 10.7) (font (size 2 2) bold italic))
 (tbtext "File: ${FILENAME}" (pos 109 14.3))
 (line (start 110 18.5) (end 2 18.5))
 (tbtext "Sheet: ${SHEETPATH}" (pos 109 17))
 (tbtext "${COMPANY}" (pos 109 20) (font bold))
 (tbtext "${COMMENT1}" (pos 109 23))
 (tbtext "${COMMENT2}" (pos 109 26))
 (tbtext "${COMMENT3}" (pos 109 29))
 (tbtext "${COMMENT4}" (pos 109 32))
 (line (start 90 8.5) (end 90 5.5))
 (line (start 26 8.5) (end 26 2)))`

// worksheet is a drawing sheet (.kicad_wks, or the page_layout files of
// KiCad 5).
type worksheet struct {
	n                        *node
	size                     int // bytes
	textW, textH             float64
	lineWidth, textLineWidth float64
	left, right, top, bottom float64
}

func parseWorksheet(b []byte) (*worksheet, error) {
	n, err := parse(b)
	if err != nil {
		return nil, err
	}
	if n.name != "kicad_wks" && n.name != "page_layout" && n.name != "drawing_sheet" {
		return nil, errSyntax
	}
	s := n.child("setup")
	ws := &worksheet{n: n, size: len(b), textW: 1.5, textH: 1.5, lineWidth: 0.15, textLineWidth: 0.15, left: 10, right: 10, top: 10, bottom: 10}
	if ts := s.child("textsize"); ts != nil {
		ws.textW, ws.textH = ts.num(0), ts.num(1)
	}
	ws.lineWidth = s.numOf("linewidth", ws.lineWidth)
	ws.textLineWidth = s.numOf("textlinewidth", ws.textLineWidth)
	ws.left = s.numOf("left_margin", ws.left)
	ws.right = s.numOf("right_margin", ws.right)
	ws.top = s.numOf("top_margin", ws.top)
	ws.bottom = s.numOf("bottom_margin", ws.bottom)
	return ws, nil
}

var defaultWorksheet, _ = parseWorksheet([]byte(defaultSheet))

// wsPage is a page a drawing sheet is drawn on.
type wsPage struct {
	w, h   float64
	first  bool      // the first page (page1only, notonpage1)
	color  bdf.Color // lines and text
	minPen float64   // the thinnest line drawn
	vars   func(string) (string, bool)
	bg, fg *cad.Drawing // items without and with variables
}

// point returns a position of the sheet (x y [corner]) on the page, moved
// by i repetitions of the item's increments.
func (ws *worksheet) point(n *node, it *node, i int, pg *wsPage) cad.Point {
	x, y := n.num(0), n.num(1)
	x += float64(i) * it.numOf("incrx", 0)
	y += float64(i) * it.numOf("incry", 0)
	switch n.arg(2) {
	case "ltcorner":
		return cad.Point{X: ws.left + x, Y: ws.top + y}
	case "lbcorner":
		return cad.Point{X: ws.left + x, Y: pg.h - ws.bottom - y}
	case "rtcorner":
		return cad.Point{X: pg.w - ws.right - x, Y: ws.top + y}
	}
	return cad.Point{X: pg.w - ws.right - x, Y: pg.h - ws.bottom - y}
}

func (ws *worksheet) inside(p cad.Point, pg *wsPage) bool {
	const eps = 1e-6
	return p.X >= ws.left-eps && p.X <= pg.w-ws.right+eps && p.Y >= ws.top-eps && p.Y <= pg.h-ws.bottom+eps
}

// legacyCodes are the % codes of KiCad 5 layouts and the variables they
// became.
var legacyCodes = map[byte]string{
	'D': "${ISSUE_DATE}", 'R': "${REVISION}", 'K': "${KICAD_VERSION}", 'Z': "${PAPER}", 'S': "${#}",
	'N': "${##}", 'F': "${FILENAME}", 'P': "${SHEETPATH}", 'Y': "${COMPANY}", 'T': "${TITLE}", 'L': "${LAYER}",
}

func convertLegacy(s string) string {
	if !strings.Contains(s, "%") {
		return s
	}
	var sb strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '%' && i+1 < len(s) {
			if v, ok := legacyCodes[s[i+1]]; ok {
				sb.WriteString(v)
				i++
				continue
			}
			if s[i+1] == 'C' && i+2 < len(s) && s[i+2] >= '0' && s[i+2] <= '9' {
				sb.WriteString("${COMMENT" + strconv.Itoa(int(s[i+2]-'0')+1) + "}")
				i += 2
				continue
			}
		}
		sb.WriteByte(s[i])
	}
	return sb.String()
}

// incrLabel increments a repeated label: a number by n, a letter by n
// letters.
func incrLabel(s string, n int) string {
	if v, err := strconv.Atoi(s); err == nil {
		return strconv.Itoa(v + n)
	}
	if r := []rune(s); len(r) == 1 && unicode.IsLetter(r[0]) {
		return string(r[0] + rune(n))
	}
	return s
}

// drawWorksheet draws a drawing sheet on a page.
func (c *conv) drawWorksheet(ws *worksheet, pg *wsPage) {
	if ws == nil {
		return
	}
	for _, it := range ws.n.lists() {
		switch it.str("option") {
		case "page1only":
			if !pg.first {
				continue
			}
		case "notonpage1":
			if pg.first {
				continue
			}
		}
		repeat := max(it.child("repeat").int(0), 1)
		repeat = min(repeat, 200)
		for i := range repeat {
			switch it.name {
			case "line", "rect":
				a := ws.point(it.child("start"), it, i, pg)
				b := ws.point(it.child("end"), it, i, pg)
				if i > 0 && (!ws.inside(a, pg) || !ws.inside(b, pg)) {
					continue
				}
				w := math.Max(it.numOf("linewidth", ws.lineWidth), pg.minPen)
				var path *cad.Path
				if it.name == "line" {
					path = (&cad.Path{}).Polyline([]cad.Point{a, b}, false)
				} else {
					path = (&cad.Path{}).Polyline([]cad.Point{a, {X: b.X, Y: a.Y}, b, {X: a.X, Y: b.Y}}, true)
				}
				pg.bg.Stroke(path, pen(w, "solid", pg.color))
			case "tbtext":
				p := ws.point(it.child("pos"), it, i, pg)
				if i > 0 && !ws.inside(p, pg) {
					continue
				}
				s := convertLegacy(it.arg(0))
				if i > 0 {
					s = incrLabel(s, i*max(it.child("incrlabel").int(0), 1))
				}
				c.wsText(ws, it, s, p, pg)
			case "polygon":
				p := ws.point(it.child("pos"), it, i, pg)
				if i > 0 && !ws.inside(p, pg) {
					continue
				}
				angle := it.numOf("rotate", 0)
				path := &cad.Path{}
				for _, pts := range it.children("pts") {
					var poly []cad.Point
					for _, q := range pts.children("xy") {
						x, y := rotatePt(q.num(0), q.num(1), angle)
						poly = append(poly, cad.Point{X: p.X + x, Y: p.Y + y})
					}
					if len(poly) >= 2 {
						path.Polyline(poly, true)
					}
				}
				if !path.Empty() {
					pg.bg.Fill(path, cad.Fill{Color: pg.color}, false)
				}
			case "bitmap":
				p := ws.point(it.child("pos"), it, i, pg)
				if i > 0 && !ws.inside(p, pg) {
					continue
				}
				c.wsBitmap(it, p, pg)
			}
		}
	}
}

// wsText draws a text of a drawing sheet.
func (c *conv) wsText(ws *worksheet, it *node, s string, p cad.Point, pg *wsPage) {
	f := it.child("font")
	e := effects{sizeX: ws.textW, sizeY: ws.textH, lineSpacing: 1, h: alignLeft, v: alignMiddle}
	if sz := f.child("size"); sz != nil {
		e.sizeX, e.sizeY = sz.num(0), sz.num(1)
	}
	e.bold = f.flag("bold")
	e.italic = f.flag("italic")
	e.face = f.str("face")
	if cl, ok := parseColor(f.child("color")); ok {
		e.color, e.hasColor = cl, true
	}
	if j := it.child("justify"); j != nil {
		for i := 0; j.arg(i) != ""; i++ {
			switch j.arg(i) {
			case "left":
				e.h = alignLeft
			case "center":
				e.h = alignCenter
			case "right":
				e.h = alignRight
			case "top":
				e.v = alignTop
			case "bottom":
				e.v = alignBottom
			}
		}
	}
	thick := f.numOf("linewidth", 0)
	if thick <= 0 {
		if e.bold {
			thick = boldThickness(math.Min(e.sizeX, e.sizeY))
		} else {
			thick = math.Max(ws.textLineWidth, pg.minPen)
		}
	}
	variable := strings.Contains(s, "${")
	s = c.expand(s, pg.vars)
	if strings.TrimSpace(s) == "" {
		return
	}
	// maxlen and maxheight shrink a text to fit
	if ml := it.numOf("maxlen", 0); ml > 0 {
		x0, _, w, _ := c.textBox(s, e)
		_ = x0
		if w > ml {
			e.sizeX *= ml / w
		}
	}
	color := pg.color
	if e.hasColor {
		color = e.color
	}
	d := pg.bg
	if variable {
		d = pg.fg
	}
	c.text(d, placedText{s: s, x: p.X, y: p.Y, angle: it.numOf("rotate", 0), e: e, color: color, thickness: thick, brk: cad.BreakBox})
}

// wsBitmap draws an image of a drawing sheet: PNG data, at its resolution
// (300 pixels per inch unless it says otherwise) times its scale, centered
// on its position.
func (c *conv) wsBitmap(it *node, p cad.Point, pg *wsPage) {
	bm := c.bitmaps[it]
	if bm == nil {
		bm = c.decodeBitmap(it)
		if c.bitmaps == nil {
			c.bitmaps = map[*node]*sheetBitmap{}
		}
		c.bitmaps[it] = bm
	}
	if bm.hash == (bdf.Hash{}) {
		return
	}
	w, h := float64(bm.w)*bm.mm, float64(bm.h)*bm.mm
	m := canvas.Translate(p.X-w/2, p.Y-h/2).Mul(canvas.Scale(bm.mm, bm.mm))
	pg.bg.Image(&cad.Image{Image: bm.hash, W: bm.w, H: bm.h, M: m, Smooth: true})
}

// sheetBitmap is an image of the drawing sheet, added to the document
// once: its size in pixels and mm a pixel (a zero hash when it cannot be
// read).
type sheetBitmap struct {
	hash bdf.Hash
	w, h int
	mm   float64
}

func (c *conv) decodeBitmap(it *node) *sheetBitmap {
	var data []byte
	if pd := it.child("pngdata"); pd != nil {
		var sb strings.Builder
		for _, d := range pd.children("data") {
			sb.WriteString(d.arg(0))
		}
		data, _ = hex.DecodeString(strings.NewReplacer(" ", "", "\n", "", "\t", "").Replace(sb.String()))
	} else {
		var sb strings.Builder
		for _, d := range it.children("data") {
			for i := 0; d.arg(i) != ""; i++ {
				sb.WriteString(d.arg(i))
			}
		}
		data, _ = base64.StdEncoding.DecodeString(sb.String())
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || cfg.Width <= 0 || cfg.Height <= 0 {
		c.warnOnce("ws-bitmap", "an image of the drawing sheet cannot be read and is left out")
		return &sheetBitmap{}
	}
	ppi := pngPPI(data)
	if ppi <= 0 {
		ppi = 300
	}
	scale := it.numOf("scale", 1)
	if scale <= 0 {
		scale = 1
	}
	return &sheetBitmap{hash: c.doc.AddImage(data), w: cfg.Width, h: cfg.Height, mm: 25.4 / ppi * scale}
}
