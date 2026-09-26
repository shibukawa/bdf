package xlsx

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// The styles part (ECMA-376 Part 1 §18.8): cell formats (xf) index number
// formats, fonts, fills and borders; differential formats (dxf) are the
// partial formats of conditional formatting and table styles.

// rgb is an opaque color with components in 0..1.
type rgb struct{ R, G, B float64 }

func (c rgb) bdf() bdf.Color {
	q := func(v float64) uint8 { return uint8(math.Round(clamp01(v) * 255)) }
	return bdf.RGB(q(c.R), q(c.G), q(c.B))
}

func clamp01(v float64) float64 {
	if v < 0 || math.IsNaN(v) {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func hexRGB(v uint32) rgb {
	return rgb{float64(v>>16&0xff) / 255, float64(v>>8&0xff) / 255, float64(v&0xff) / 255}
}

var (
	black = rgb{0, 0, 0}
	white = rgb{1, 1, 1}
)

// defaultPalette is the indexed color table (§18.8.27); 64 is the system
// foreground and 65 the system background.
var defaultPalette = func() []rgb {
	v := []uint32{
		0x000000, 0xFFFFFF, 0xFF0000, 0x00FF00, 0x0000FF, 0xFFFF00, 0xFF00FF, 0x00FFFF,
		0x000000, 0xFFFFFF, 0xFF0000, 0x00FF00, 0x0000FF, 0xFFFF00, 0xFF00FF, 0x00FFFF,
		0x800000, 0x008000, 0x000080, 0x808000, 0x800080, 0x008080, 0xC0C0C0, 0x808080,
		0x9999FF, 0x993366, 0xFFFFCC, 0xCCFFFF, 0x660066, 0xFF8080, 0x0066CC, 0xCCCCFF,
		0x000080, 0xFF00FF, 0xFFFF00, 0x00FFFF, 0x800080, 0x800000, 0x008080, 0x0000FF,
		0x00CCFF, 0xCCFFFF, 0xCCFFCC, 0xFFFF99, 0x99CCFF, 0xFF99CC, 0xCC99FF, 0xFFCC99,
		0x3366FF, 0x33CCCC, 0x99CC00, 0xFFCC00, 0xFF9900, 0xFF6600, 0x666699, 0x969696,
		0x003366, 0x339966, 0x003300, 0x333300, 0x993300, 0x993366, 0x333399, 0x333333,
		0x000000, 0xFFFFFF,
	}
	out := make([]rgb, len(v))
	for i, c := range v {
		out[i] = hexRGB(c)
	}
	return out
}()

// themeNames maps the theme color indexes of SpreadsheetML to the color
// scheme: the first four are the light/dark pairs in the opposite order of
// the scheme (§18.8.3 note).
var themeNames = [...]string{"lt1", "dk1", "lt2", "dk2", "accent1", "accent2", "accent3", "accent4", "accent5", "accent6", "hlink", "folHlink"}

// colorRef is an unresolved CT_Color.
type colorRef struct {
	kind byte // colorNone, colorRGB, colorTheme, colorIndexed, colorAuto
	rgb  rgb
	idx  int
	tint float64
}

const (
	colorNone = iota
	colorRGB
	colorTheme
	colorIndexed
	colorAuto
)

func parseColor(n *ooxml.Node) colorRef {
	if n == nil {
		return colorRef{}
	}
	c := colorRef{tint: n.AttrFloat("tint", 0)}
	switch {
	case n.AttrBool("auto", false):
		c.kind = colorAuto
	case hasAttr(n, "rgb"):
		s := strings.TrimSpace(n.AttrStr("rgb", ""))
		if len(s) == 8 {
			s = s[2:] // ARGB: Excel ignores the alpha
		}
		if v, err := strconv.ParseUint(s, 16, 32); err == nil && len(s) == 6 {
			c.kind, c.rgb = colorRGB, hexRGB(uint32(v))
		}
	case hasAttr(n, "theme"):
		c.kind, c.idx = colorTheme, int(n.AttrInt("theme", 0))
	case hasAttr(n, "indexed"):
		c.kind, c.idx = colorIndexed, int(n.AttrInt("indexed", 0))
	}
	return c
}

func hasAttr(n *ooxml.Node, name string) bool {
	_, ok := n.Attr(name)
	return ok
}

// color resolves a color; def is used for none and auto.
func (st *styles) color(c colorRef, def rgb) rgb {
	var out rgb
	switch c.kind {
	case colorRGB:
		out = c.rgb
	case colorTheme:
		if c.idx < 0 || c.idx >= len(themeNames) {
			return def
		}
		var ok bool
		if out, ok = st.theme(themeNames[c.idx]); !ok {
			return def
		}
	case colorIndexed:
		switch {
		case c.idx == 64 || c.idx == 81:
			out = black // system foreground, tooltip text
		case c.idx == 65 || c.idx == 80:
			out = white // system background, tooltip background
		case c.idx >= 0 && c.idx < len(st.palette):
			out = st.palette[c.idx]
		default:
			return def
		}
	default:
		return def
	}
	if c.tint != 0 {
		out = applyTint(out, c.tint)
	}
	return out
}

// applyTint lightens (tint > 0) or darkens a color in HLS (§18.8.19).
func applyTint(c rgb, tint float64) rgb {
	h, l, s := rgbToHLS(c)
	if tint < 0 {
		l *= 1 + tint
	} else {
		l = l*(1-tint) + tint
	}
	return hlsToRGB(h, clamp01(l), s)
}

func rgbToHLS(c rgb) (h, l, s float64) {
	mx := math.Max(c.R, math.Max(c.G, c.B))
	mn := math.Min(c.R, math.Min(c.G, c.B))
	l = (mx + mn) / 2
	if mx == mn {
		return 0, l, 0
	}
	d := mx - mn
	if l > 0.5 {
		s = d / (2 - mx - mn)
	} else {
		s = d / (mx + mn)
	}
	switch mx {
	case c.R:
		h = (c.G - c.B) / d
		if c.G < c.B {
			h += 6
		}
	case c.G:
		h = (c.B-c.R)/d + 2
	default:
		h = (c.R-c.G)/d + 4
	}
	return h / 6, l, s
}

func hlsToRGB(h, l, s float64) rgb {
	if s == 0 {
		return rgb{l, l, l}
	}
	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q
	f := func(t float64) float64 {
		if t < 0 {
			t++
		}
		if t > 1 {
			t--
		}
		switch {
		case t < 1.0/6:
			return p + (q-p)*6*t
		case t < 0.5:
			return q
		case t < 2.0/3:
			return p + (q-p)*(2.0/3-t)*6
		}
		return p
	}
	return rgb{f(h + 1.0/3), f(h), f(h - 1.0/3)}
}

func mix(a, b rgb, t float64) rgb {
	return rgb{a.R + (b.R-a.R)*t, a.G + (b.G-a.G)*t, a.B + (b.B-a.B)*t}
}

// xfont is a font of the font table or a rich text run.
type xfont struct {
	name      string
	size      float64
	bold      bool
	italic    bool
	underline string // "", single, double, singleAccounting, doubleAccounting
	strike    bool
	color     colorRef
	vert      string // "", superscript, subscript
	scheme    string // "", minor, major
	family    int
}

func parseFont(n *ooxml.Node, def *xfont) xfont {
	f := xfont{name: "Calibri", size: 11}
	if def != nil {
		f = *def
	}
	if n == nil {
		return f
	}
	for _, k := range n.Kids {
		val := k.AttrStr("val", "")
		switch k.Name {
		case "name", "rFont":
			if val != "" {
				f.name = val
			}
		case "sz":
			if v := k.AttrFloat("val", 0); v > 0 && v < 1000 {
				f.size = v
			}
		case "b":
			f.bold = k.AttrBool("val", true)
		case "i":
			f.italic = k.AttrBool("val", true)
		case "u":
			f.underline = k.AttrStr("val", "single")
			if f.underline == "none" {
				f.underline = ""
			}
		case "strike":
			f.strike = k.AttrBool("val", true)
		case "color":
			f.color = parseColor(k)
		case "vertAlign":
			if val == "baseline" {
				val = ""
			}
			f.vert = val
		case "scheme":
			if val == "none" {
				val = ""
			}
			f.scheme = val
		case "family":
			f.family = int(k.AttrInt("val", 0))
		}
	}
	return f
}

// xfill is a pattern or gradient fill.
type xfill struct {
	pattern string // "" (none), solid, gray125 …
	fg, bg  colorRef
	grad    *gradFill
}

type gradFill struct {
	path                     bool
	degree                   float64
	left, right, top, bottom float64
	stops                    []gradStop
}

type gradStop struct {
	pos   float64
	color colorRef
}

func parseFill(n *ooxml.Node, dxf bool) xfill {
	var f xfill
	if pf := n.Child("patternFill"); pf != nil {
		f.pattern = pf.AttrStr("patternType", "")
		if f.pattern == "none" {
			f.pattern = ""
		}
		f.fg = parseColor(pf.Child("fgColor"))
		f.bg = parseColor(pf.Child("bgColor"))
		if dxf {
			// In a differential format a solid fill's color is its
			// background color, and a missing pattern type means solid.
			if f.pattern == "" && (f.bg.kind != colorNone || f.fg.kind != colorNone) {
				f.pattern = "solid"
			}
			if f.pattern == "solid" && f.bg.kind != colorNone {
				f.fg = f.bg
			}
		}
	}
	if gf := n.Child("gradientFill"); gf != nil {
		g := &gradFill{path: gf.AttrStr("type", "linear") == "path", degree: gf.AttrFloat("degree", 0),
			left: gf.AttrFloat("left", 0), right: gf.AttrFloat("right", 0), top: gf.AttrFloat("top", 0), bottom: gf.AttrFloat("bottom", 0)}
		for _, s := range gf.Children("stop") {
			g.stops = append(g.stops, gradStop{pos: s.AttrFloat("position", 0), color: parseColor(s.Child("color"))})
		}
		if len(g.stops) > 0 {
			f.grad = g
		}
	}
	return f
}

func (f xfill) empty() bool { return f.pattern == "" && f.grad == nil }

// borderSide is one line of a border.
type borderSide struct {
	style string // "" (none), thin, medium, dashed, dotted, thick, double, hair …
	color colorRef
}

type xborder struct {
	left, right, top, bottom borderSide
	diag                     borderSide
	diagUp, diagDown         bool
	vertical, horizontal     borderSide // inner lines (differential formats)
}

func parseBorder(n *ooxml.Node) xborder {
	side := func(name string) borderSide {
		e := n.Child(name)
		if e == nil {
			switch name {
			case "left":
				e = n.Child("start")
			case "right":
				e = n.Child("end")
			}
		}
		s := borderSide{style: e.AttrStr("style", ""), color: parseColor(e.Child("color"))}
		if s.style == "none" {
			s.style = ""
		}
		return s
	}
	return xborder{left: side("left"), right: side("right"), top: side("top"), bottom: side("bottom"), diag: side("diagonal"),
		diagUp: n.AttrBool("diagonalUp", false), diagDown: n.AttrBool("diagonalDown", false),
		vertical: side("vertical"), horizontal: side("horizontal")}
}

// xalign is the alignment of a cell format.
type xalign struct {
	h, v        string
	wrap        bool
	shrink      bool
	indent      int
	rotation    int // 0–180, 255 for stacked text
	justifyLast bool
}

func parseAlign(n *ooxml.Node) xalign {
	a := xalign{h: n.AttrStr("horizontal", "general"), v: n.AttrStr("vertical", "bottom"),
		wrap: n.AttrBool("wrapText", false), shrink: n.AttrBool("shrinkToFit", false),
		indent: int(n.AttrInt("indent", 0)), rotation: int(n.AttrInt("textRotation", 0)),
		justifyLast: n.AttrBool("justifyLastLine", false)}
	if a.rotation < 0 || a.rotation > 180 && a.rotation != 255 {
		a.rotation = 0
	}
	return a
}

// xf is a cell format.
type xf struct {
	numFmt int
	font   int
	fill   int
	border int
	align  xalign
}

// dxf is a differential format: only what it sets.
type dxf struct {
	font   *dxfFont
	fill   *xfill
	border *xborder
	numFmt *numFormat
}

type dxfFont struct {
	bold, italic, strike *bool
	underline            *string
	color                colorRef
}

type styles struct {
	c          *converter
	numFmts    map[int]string
	parsed     map[int]*numFormat
	fonts      []xfont
	fills      []xfill
	borders    []xborder
	cellXfs    []xf
	dxfs       []dxf
	palette    []rgb
	tableStyle map[string]*ooxml.Node // custom table styles by name
	defTable   string
	themeColor func(name string) (rgb, bool)
}

func (st *styles) theme(name string) (rgb, bool) { return st.themeColor(name) }

func loadStyles(c *converter, n *ooxml.Node, themeColor func(string) (rgb, bool)) *styles {
	st := &styles{c: c, numFmts: map[int]string{}, parsed: map[int]*numFormat{}, palette: defaultPalette,
		tableStyle: map[string]*ooxml.Node{}, themeColor: themeColor}
	if ic := n.Path("colors", "indexedColors"); ic != nil {
		pal := append([]rgb(nil), defaultPalette...)
		for i, k := range ic.Children("rgbColor") {
			if i >= len(pal) {
				break
			}
			if c := parseColor(k); c.kind == colorRGB {
				pal[i] = c.rgb
			}
		}
		st.palette = pal
	}
	for _, k := range n.Path("numFmts").Children("numFmt") {
		st.numFmts[int(k.AttrInt("numFmtId", -1))] = k.AttrStr("formatCode", "General")
	}
	for _, k := range n.Path("fonts").Children("font") {
		st.fonts = append(st.fonts, parseFont(k, nil))
	}
	if len(st.fonts) == 0 {
		st.fonts = []xfont{{name: "Calibri", size: 11, scheme: "minor"}}
	}
	for _, k := range n.Path("fills").Children("fill") {
		st.fills = append(st.fills, parseFill(k, false))
	}
	for _, k := range n.Path("borders").Children("border") {
		st.borders = append(st.borders, parseBorder(k))
	}
	var styleXfs []xf
	readXf := func(k *ooxml.Node) xf {
		return xf{numFmt: int(k.AttrInt("numFmtId", 0)), font: int(k.AttrInt("fontId", 0)),
			fill: int(k.AttrInt("fillId", 0)), border: int(k.AttrInt("borderId", 0)), align: parseAlign(k.Child("alignment"))}
	}
	for _, k := range n.Path("cellStyleXfs").Children("xf") {
		styleXfs = append(styleXfs, readXf(k))
	}
	for _, k := range n.Path("cellXfs").Children("xf") {
		x := readXf(k)
		// a format that does not apply its own alignment takes its style's
		if k.Child("alignment") == nil && !k.AttrBool("applyAlignment", false) {
			if id := int(k.AttrInt("xfId", 0)); id >= 0 && id < len(styleXfs) {
				x.align = styleXfs[id].align
			}
		}
		st.cellXfs = append(st.cellXfs, x)
	}
	if len(st.cellXfs) == 0 {
		st.cellXfs = []xf{{align: xalign{h: "general", v: "bottom"}}}
	}
	for _, k := range n.Path("dxfs").Children("dxf") {
		st.dxfs = append(st.dxfs, st.parseDxf(k))
	}
	ts := n.Child("tableStyles")
	st.defTable = ts.AttrStr("defaultTableStyle", "TableStyleMedium2")
	for _, k := range ts.Children("tableStyle") {
		st.tableStyle[decodeXString(k.AttrStr("name", ""))] = k
	}
	return st
}

func (st *styles) parseDxf(k *ooxml.Node) dxf {
	var d dxf
	if f := k.Child("font"); f != nil {
		df := &dxfFont{}
		for _, e := range f.Kids {
			switch e.Name {
			case "b":
				v := e.AttrBool("val", true)
				df.bold = &v
			case "i":
				v := e.AttrBool("val", true)
				df.italic = &v
			case "strike":
				v := e.AttrBool("val", true)
				df.strike = &v
			case "u":
				v := e.AttrStr("val", "single")
				if v == "none" {
					v = ""
				}
				df.underline = &v
			case "color":
				df.color = parseColor(e)
			}
		}
		d.font = df
	}
	if f := k.Child("fill"); f != nil {
		fl := parseFill(f, true)
		d.fill = &fl
	}
	if b := k.Child("border"); b != nil {
		br := parseBorder(b)
		d.border = &br
	}
	if nf := k.Child("numFmt"); nf != nil {
		d.numFmt = parseNumFormat(nf.AttrStr("formatCode", "General"), st.palette)
	}
	return d
}

// xf returns a cell format (the default one for bad indexes).
func (st *styles) xf(i int) *xf {
	if i < 0 || i >= len(st.cellXfs) {
		return &st.cellXfs[0]
	}
	return &st.cellXfs[i]
}

func (st *styles) font(i int) *xfont {
	if i < 0 || i >= len(st.fonts) {
		return &st.fonts[0]
	}
	return &st.fonts[i]
}

func (st *styles) fill(i int) xfill {
	if i < 0 || i >= len(st.fills) {
		return xfill{}
	}
	return st.fills[i]
}

func (st *styles) border(i int) *xborder {
	if i < 0 || i >= len(st.borders) {
		return &xborder{}
	}
	return &st.borders[i]
}

// numFormat returns a parsed number format by id.
func (st *styles) numFormat(id int) *numFormat {
	if nf, ok := st.parsed[id]; ok {
		return nf
	}
	code, ok := st.numFmts[id]
	if !ok {
		code = builtinFormat(id, st.c.locale)
	}
	nf := parseNumFormat(code, st.palette)
	st.parsed[id] = nf
	return nf
}
