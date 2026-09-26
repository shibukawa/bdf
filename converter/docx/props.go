package docx

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// Measures in WordprocessingML are mostly twentieths of a point (twips);
// font sizes are half-points, borders eighths of a point, drawings EMU.

// twips reads a twips attribute as points; strict documents may use
// universal measures ("12pt", "2cm").
func twips(n *ooxml.Node, name string, def float64) float64 {
	v, ok := n.Attr(name)
	if !ok {
		return def
	}
	return measure(v, 20, def)
}

// measure parses a number in units per point (twips: 20), or a universal
// measure.
func measure(v string, perPt, def float64) float64 {
	v = strings.TrimSpace(v)
	for _, u := range []struct {
		suffix string
		pt     float64
	}{{"pt", 1}, {"in", 72}, {"cm", 72 / 2.54}, {"mm", 72 / 25.4}, {"pc", 12}, {"pi", 12}} {
		if strings.HasSuffix(v, u.suffix) {
			if f, err := strconv.ParseFloat(strings.TrimSuffix(v, u.suffix), 64); err == nil {
				return f * u.pt
			}
			return def
		}
	}
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		return f / perPt
	}
	return def
}

// onOff reads an on/off property element (<w:b/>, <w:b w:val="0"/>); nil
// is off.
func onOff(n *ooxml.Node) bool {
	if n == nil {
		return false
	}
	return n.AttrBool("val", true)
}

func val(n *ooxml.Node) string { return n.AttrStr("val", "") }

// border is a line of a paragraph, table or cell border.
type border struct {
	style string    // single, double, dotted, dashed, thick …
	width float64   // points
	space float64   // distance from the text, points
	color bdf.Color //
}

// none reports whether the border draws nothing.
func (b *border) none() bool {
	return b == nil || b.style == "" || b.style == "none" || b.style == "nil" || b.width <= 0
}

func (c *converter) parseBorder(n *ooxml.Node) *border {
	if n == nil {
		return nil
	}
	b := &border{style: val(n), width: float64(n.AttrInt("sz", 4)) / 8, space: float64(n.AttrInt("space", 0)), color: bdf.RGB(0, 0, 0)}
	if b.style == "nil" || b.style == "none" {
		return b
	}
	// Word draws its hairlines at least a quarter point wide.
	b.width = math.Max(b.width, 0.25)
	if col, ok := c.color(n.AttrStr("color", "auto"), n.AttrStr("themeColor", ""), n.AttrStr("themeTint", ""), n.AttrStr("themeShade", "")); ok {
		b.color = col
	}
	return b
}

// shading returns the color a w:shd paints (false when it paints none).
func (c *converter) shading(n *ooxml.Node) (bdf.Color, bool) {
	if n == nil {
		return 0, false
	}
	fill, fok := c.color(n.AttrStr("fill", "auto"), n.AttrStr("themeFill", ""), n.AttrStr("themeFillTint", ""), n.AttrStr("themeFillShade", ""))
	fg, gok := c.color(n.AttrStr("color", "auto"), n.AttrStr("themeColor", ""), n.AttrStr("themeTint", ""), n.AttrStr("themeShade", ""))
	pattern := val(n)
	switch {
	case pattern == "solid":
		if !gok {
			fg = bdf.RGB(0, 0, 0)
		}
		return fg, true
	case strings.HasPrefix(pattern, "pct"):
		// a percentage of the pattern color over the fill
		p, err := strconv.Atoi(strings.TrimPrefix(pattern, "pct"))
		if err != nil {
			return fill, fok
		}
		if !gok {
			fg = bdf.RGB(0, 0, 0)
		}
		if !fok {
			fill = bdf.RGB(255, 255, 255)
		}
		return mix(fill, fg, float64(p)/100), true
	case pattern == "nil":
		return 0, false
	}
	return fill, fok
}

func mix(a, b bdf.Color, t float64) bdf.Color {
	ch := func(x, y bdf.Color, s uint) uint8 {
		return uint8(math.Round(float64(uint8(x>>s))*(1-t) + float64(uint8(y>>s))*t))
	}
	return bdf.RGB(ch(a, b, 24), ch(a, b, 16), ch(a, b, 8))
}

// dark reports whether text on a background of c needs a light color.
func dark(c bdf.Color) bool {
	r, g, b := float64(uint8(c>>24)), float64(uint8(c>>16)), float64(uint8(c>>8))
	return 0.299*r+0.587*g+0.114*b < 96
}

// tabStop is a custom tab stop of a paragraph.
type tabStop struct {
	pos    float64
	align  string // left, center, right, decimal, bar
	leader string // none, dot, hyphen, underscore, heavy, middleDot
}

// pprops are resolved paragraph properties.
type pprops struct {
	jc                           string // left, center, right, both, distribute
	indL, indR, indFirst         float64
	indLCh, indRCh, indFirstCh   float64 // hundredths of a character (0: not set)
	before, after                float64
	beforeLn, afterLn            float64 // hundredths of a line (0: not set)
	beforeAuto, afterAuto        bool
	line                         float64 // 240ths of a line (auto) or points
	lineRule                     string  // auto, exact, atLeast
	contextual                   bool
	keepNext, keepLines          bool
	pageBreakBefore, widowCtl    bool
	snapToGrid                   bool
	tabs                         []tabStop
	numID                        string
	ilvl                         int
	outline                      int // outline level 0-8, -1 for body text
	bdr                          [5]*border
	shd                          bdf.Color
	hasShd                       bool
	kinsoku, overflowPunct       bool
	autoSpaceDE, autoSpaceDN     bool
	bidi                         bool
	suppressLineNumbers, framePr bool
}

// border sides of a paragraph
const (
	bTop = iota
	bLeft
	bBottom
	bRight
	bBetween
)

func defaultPProps() *pprops {
	return &pprops{jc: "left", line: 240, lineRule: "auto", snapToGrid: true, outline: -1,
		kinsoku: true, overflowPunct: true, autoSpaceDE: true, autoSpaceDN: true}
}

// apply layers a w:pPr over the properties.
func (c *converter) applyPPr(p *pprops, n *ooxml.Node) {
	for _, k := range n.Elements() {
		switch k.Name {
		case "keepNext":
			p.keepNext = onOff(k)
		case "keepLines":
			p.keepLines = onOff(k)
		case "pageBreakBefore":
			p.pageBreakBefore = onOff(k)
		case "widowControl":
			p.widowCtl = onOff(k)
		case "contextualSpacing":
			p.contextual = onOff(k)
		case "snapToGrid":
			p.snapToGrid = onOff(k)
		case "kinsoku":
			p.kinsoku = onOff(k)
		case "overflowPunct":
			p.overflowPunct = onOff(k)
		case "autoSpaceDE":
			p.autoSpaceDE = onOff(k)
		case "autoSpaceDN":
			p.autoSpaceDN = onOff(k)
		case "bidi":
			p.bidi = onOff(k)
		case "suppressLineNumbers":
			p.suppressLineNumbers = onOff(k)
		case "framePr":
			p.framePr = true
		case "numPr":
			if id := k.Child("numId"); id != nil {
				p.numID = val(id)
			}
			if l := k.Child("ilvl"); l != nil {
				p.ilvl = int(l.AttrInt("val", 0))
			}
		case "outlineLvl":
			p.outline = int(k.AttrInt("val", 9))
			if p.outline > 8 {
				p.outline = -1
			}
		case "jc":
			p.jc = jcValue(val(k))
		case "pBdr":
			for i, side := range [][]string{{"top"}, {"left", "start"}, {"bottom"}, {"right", "end"}, {"between"}} {
				for _, name := range side {
					if b := k.Child(name); b != nil {
						p.bdr[i] = c.parseBorder(b)
					}
				}
			}
		case "shd":
			p.shd, p.hasShd = c.shading(k)
		case "tabs":
			for _, t := range k.Children("tab") {
				pos := twips(t, "pos", 0)
				kept := p.tabs[:0]
				for _, o := range p.tabs {
					if math.Abs(o.pos-pos) > 0.05 {
						kept = append(kept, o)
					}
				}
				p.tabs = kept
				align := val(t)
				switch align {
				case "clear":
					continue
				case "start":
					align = "left"
				case "end":
					align = "right"
				case "num":
					align = "left"
				}
				p.tabs = append(p.tabs, tabStop{pos: pos, align: align, leader: t.AttrStr("leader", "none")})
			}
			sortTabs(p.tabs)
		case "spacing":
			if _, ok := k.Attr("before"); ok {
				p.before = twips(k, "before", 0)
			}
			if _, ok := k.Attr("after"); ok {
				p.after = twips(k, "after", 0)
			}
			if _, ok := k.Attr("beforeLines"); ok {
				p.beforeLn = k.AttrFloat("beforeLines", 0)
			}
			if _, ok := k.Attr("afterLines"); ok {
				p.afterLn = k.AttrFloat("afterLines", 0)
			}
			if _, ok := k.Attr("beforeAutospacing"); ok {
				p.beforeAuto = k.AttrBool("beforeAutospacing", false)
			}
			if _, ok := k.Attr("afterAutospacing"); ok {
				p.afterAuto = k.AttrBool("afterAutospacing", false)
			}
			if _, ok := k.Attr("line"); ok {
				p.line = k.AttrFloat("line", 240)
				p.lineRule = k.AttrStr("lineRule", "auto")
			} else if r, ok := k.Attr("lineRule"); ok {
				p.lineRule = r
			}
		case "ind":
			for _, a := range []struct {
				dst   *float64
				names []string
			}{{&p.indL, []string{"left", "start"}}, {&p.indR, []string{"right", "end"}}} {
				for _, name := range a.names {
					if _, ok := k.Attr(name); ok {
						*a.dst = twips(k, name, 0)
					}
				}
			}
			if _, ok := k.Attr("firstLine"); ok {
				p.indFirst = twips(k, "firstLine", 0)
			}
			if _, ok := k.Attr("hanging"); ok {
				p.indFirst = -twips(k, "hanging", 0)
			}
			for _, a := range []struct {
				dst   *float64
				names []string
			}{{&p.indLCh, []string{"leftChars", "startChars"}}, {&p.indRCh, []string{"rightChars", "endChars"}}} {
				for _, name := range a.names {
					if _, ok := k.Attr(name); ok {
						*a.dst = k.AttrFloat(name, 0)
					}
				}
			}
			if _, ok := k.Attr("firstLineChars"); ok {
				p.indFirstCh = k.AttrFloat("firstLineChars", 0)
			}
			if _, ok := k.Attr("hangingChars"); ok {
				p.indFirstCh = -k.AttrFloat("hangingChars", 0)
			}
		}
	}
}

func jcValue(v string) string {
	switch v {
	case "start", "left":
		return "left"
	case "end", "right":
		return "right"
	case "center":
		return "center"
	case "both", "mediumKashida", "highKashida", "lowKashida":
		return "both"
	case "distribute", "thaiDistribute":
		return "distribute"
	}
	return "left"
}

func sortTabs(t []tabStop) {
	for i := 1; i < len(t); i++ {
		for j := i; j > 0 && t[j].pos < t[j-1].pos; j-- {
			t[j], t[j-1] = t[j-1], t[j]
		}
	}
}

// Toggle properties: character styles flip what paragraph styles set.
const (
	tB = 1 << iota
	tI
	tCaps
	tSmallCaps
	tStrike
	tDStrike
	tVanish
	tOutline
	tShadow
	tEmboss
	tImprint
)

// rprops are resolved run properties.
type rprops struct {
	fonts     [4]string // ascii, hAnsi, eastAsia, cs
	themes    [4]string // theme fonts overriding them
	hint      string
	toggles   int // tB …
	set       int // toggles specified (while resolving a style chain)
	sz        float64
	color     bdf.Color
	autoColor bool
	highlight string
	shd       bdf.Color
	hasShd    bool
	u         string
	uColor    bdf.Color
	hasUColor bool
	vertAlign string
	position  float64 // points, raises
	spacing   float64 // points
	scale     float64 // horizontal scale (w:w), 1 = 100%
	em        string  // emphasis mark
	lang      [3]string
	webHidden bool
	rtl       bool
	tcy       bool // horizontal in vertical text (tate-chu-yoko)
	tcyFit    bool // compressed to the width of the line
}

func defaultRProps() *rprops {
	return &rprops{fonts: [4]string{"Times New Roman", "Times New Roman", "", "Times New Roman"}, sz: 10,
		autoColor: true, color: bdf.RGB(0, 0, 0), u: "none", scale: 1}
}

func (r *rprops) clone() *rprops { n := *r; return &n }

func (r *rprops) has(t int) bool { return r.toggles&t != 0 }

func (r *rprops) setToggle(t int, on bool) {
	if on {
		r.toggles |= t
	} else {
		r.toggles &^= t
	}
	r.set |= t
}

var toggleNames = map[string]int{"b": tB, "i": tI, "caps": tCaps, "smallCaps": tSmallCaps, "strike": tStrike,
	"dstrike": tDStrike, "vanish": tVanish, "outline": tOutline, "shadow": tShadow, "emboss": tEmboss, "imprint": tImprint}

// applyRPr layers a w:rPr over the properties.
func (c *converter) applyRPr(r *rprops, n *ooxml.Node) {
	for _, k := range n.Elements() {
		if t, ok := toggleNames[k.Name]; ok {
			r.setToggle(t, onOff(k))
			continue
		}
		switch k.Name {
		case "rFonts":
			for i, name := range []string{"ascii", "hAnsi", "eastAsia", "cs"} {
				if v, ok := k.Attr(name); ok {
					r.fonts[i] = v
					r.themes[i] = ""
				}
			}
			for i, name := range []string{"asciiTheme", "hAnsiTheme", "eastAsiaTheme", "cstheme"} {
				if v, ok := k.Attr(name); ok {
					r.themes[i] = v
				}
			}
			if v, ok := k.Attr("hint"); ok {
				r.hint = v
			}
		case "sz":
			if v := k.AttrFloat("val", 0); v > 0 {
				r.sz = v / 2
			}
		case "color":
			if col, ok := c.color(k.AttrStr("val", "auto"), k.AttrStr("themeColor", ""), k.AttrStr("themeTint", ""), k.AttrStr("themeShade", "")); ok {
				r.color, r.autoColor = col, false
			} else {
				r.autoColor = true
			}
		case "highlight":
			r.highlight = val(k)
		case "shd":
			r.shd, r.hasShd = c.shading(k)
		case "u":
			r.u = k.AttrStr("val", "single")
			if col, ok := c.color(k.AttrStr("color", "auto"), k.AttrStr("themeColor", ""), k.AttrStr("themeTint", ""), k.AttrStr("themeShade", "")); ok {
				r.uColor, r.hasUColor = col, true
			} else {
				r.hasUColor = false
			}
		case "vertAlign":
			r.vertAlign = val(k)
		case "position":
			r.position = k.AttrFloat("val", 0) / 2
		case "spacing":
			r.spacing = twips(k, "val", 0)
		case "w":
			if v := k.AttrFloat("val", 100); v > 0 {
				r.scale = v / 100
			}
		case "em":
			r.em = val(k)
		case "lang":
			for i, name := range []string{"val", "eastAsia", "bidi"} {
				if v, ok := k.Attr(name); ok {
					r.lang[i] = v
				}
			}
		case "webHidden":
			r.webHidden = onOff(k)
		case "eastAsianLayout":
			r.tcy, r.tcyFit = k.AttrBool("vert", false), k.AttrBool("vertCompress", false)
			if k.AttrBool("combine", false) {
				c.warnOnce("combine", "two lines in one (eastAsianLayout combine) are drawn as ordinary text")
			}
		case "rtl":
			r.rtl = onOff(k)
		}
	}
}

// color resolves a WordprocessingML color: an RGB value or "auto", or a
// theme color with a tint or shade; false for "auto".
func (c *converter) color(v, theme, tint, shade string) (bdf.Color, bool) {
	var col bdf.Color
	ok := false
	if theme != "" && theme != "none" {
		col, ok = c.themeColor(theme)
	}
	if !ok {
		v = strings.TrimPrefix(v, "#")
		if len(v) != 6 || strings.EqualFold(v, "auto") {
			if name, found := namedColors[strings.ToLower(v)]; found {
				return name, true
			}
			return 0, false
		}
		x, err := strconv.ParseUint(v, 16, 32)
		if err != nil {
			return 0, false
		}
		return bdf.RGB(uint8(x>>16), uint8(x>>8), uint8(x)), true
	}
	if t, err := strconv.ParseUint(tint, 16, 8); err == nil && tint != "" {
		col = lumMod(col, float64(t)/255, 1-float64(t)/255)
	}
	if s, err := strconv.ParseUint(shade, 16, 8); err == nil && shade != "" {
		col = lumMod(col, float64(s)/255, 0)
	}
	return col, true
}

// themeColor maps a WordprocessingML theme color name to the theme's
// scheme, through the document's color scheme mapping.
func (c *converter) themeColor(name string) (bdf.Color, bool) {
	scheme := map[string]string{"dark1": "dk1", "light1": "lt1", "dark2": "dk2", "light2": "lt2",
		"accent1": "accent1", "accent2": "accent2", "accent3": "accent3", "accent4": "accent4",
		"accent5": "accent5", "accent6": "accent6", "hyperlink": "hlink", "followedHyperlink": "folHlink"}
	if m, ok := c.clrMap[name]; ok {
		name = m
	}
	if s, ok := scheme[name]; ok {
		return c.r.ThemeColor(c.main, s)
	}
	return 0, false
}

// lumMod scales the HSL luminance of a color and adds an offset, as
// DrawingML's lumMod and lumOff (Word's theme tints and shades).
func lumMod(col bdf.Color, mod, off float64) bdf.Color {
	r, g, b := float64(uint8(col>>24))/255, float64(uint8(col>>16))/255, float64(uint8(col>>8))/255
	mx, mn := math.Max(r, math.Max(g, b)), math.Min(r, math.Min(g, b))
	l := (mx + mn) / 2
	var h, s float64
	if mx != mn {
		d := mx - mn
		if l > 0.5 {
			s = d / (2 - mx - mn)
		} else {
			s = d / (mx + mn)
		}
		switch mx {
		case r:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/d + 2
		default:
			h = (r-g)/d + 4
		}
		h /= 6
	}
	l = math.Min(1, math.Max(0, l*mod+off))
	hue := func(p, q, t float64) float64 {
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
	if s == 0 {
		r, g, b = l, l, l
	} else {
		q := l * (1 + s)
		if l >= 0.5 {
			q = l + s - l*s
		}
		p := 2*l - q
		r, g, b = hue(p, q, h+1.0/3), hue(p, q, h), hue(p, q, h-1.0/3)
	}
	q8 := func(v float64) uint8 { return uint8(math.Round(math.Min(1, math.Max(0, v)) * 255)) }
	return bdf.RGB(q8(r), q8(g), q8(b))
}

// namedColors are the highlight colors (ST_HighlightColor).
var namedColors = map[string]bdf.Color{
	"black": bdf.RGB(0, 0, 0), "blue": bdf.RGB(0, 0, 255), "cyan": bdf.RGB(0, 255, 255), "green": bdf.RGB(0, 255, 0),
	"magenta": bdf.RGB(255, 0, 255), "red": bdf.RGB(255, 0, 0), "yellow": bdf.RGB(255, 255, 0), "white": bdf.RGB(255, 255, 255),
	"darkblue": bdf.RGB(0, 0, 128), "darkcyan": bdf.RGB(0, 128, 128), "darkgreen": bdf.RGB(0, 128, 0),
	"darkmagenta": bdf.RGB(128, 0, 128), "darkred": bdf.RGB(128, 0, 0), "darkyellow": bdf.RGB(128, 128, 0),
	"darkgray": bdf.RGB(128, 128, 128), "lightgray": bdf.RGB(192, 192, 192),
}
