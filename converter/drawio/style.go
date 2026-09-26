package drawio

import (
	"strconv"
	"strings"
)

// style is a resolved cell style: key → value, as mxStylesheet.getCellStyle
// leaves it (numbers stay strings here and are parsed on use).
type style map[string]string

// get returns the value of key, or def when it is not set.
func (s style) get(key, def string) string {
	if v, ok := s[key]; ok {
		return v
	}
	return def
}

// has reports whether key is set.
func (s style) has(key string) bool { _, ok := s[key]; return ok }

// num returns the numeric value of key (parseFloat semantics: a leading
// number is enough), or def when it is not set or not a number.
func (s style) num(key string, def float64) float64 {
	v, ok := s[key]
	if !ok {
		return def
	}
	if f, ok := parseFloat(v); ok {
		return f
	}
	return def
}

// flag reports whether key is "1" (or "true"), or def when it is not set.
func (s style) flag(key string, def bool) bool {
	v, ok := s[key]
	if !ok {
		return def
	}
	v = strings.TrimSpace(v)
	if v == "true" {
		return true
	}
	if f, ok := parseFloat(v); ok {
		return f != 0
	}
	return false
}

// is reports whether key is set to "1" (mxUtils.getValue(style, key, def) == 1).
func (s style) is(key string) bool { return s.flag(key, false) }

func (s style) clone() style {
	c := make(style, len(s)+4)
	for k, v := range s {
		c[k] = v
	}
	return c
}

// parseFloat parses the leading number of s the way JavaScript's parseFloat
// does ("12px" → 12); ok is false when s does not start with a number.
func parseFloat(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f, finite(f)
	}
	end := 0
	seenDigit, seenDot, seenExp := false, false, false
	for end < len(s) {
		c := s[end]
		switch {
		case c >= '0' && c <= '9':
			seenDigit = true
		case (c == '+' || c == '-') && (end == 0 || s[end-1] == 'e' || s[end-1] == 'E'):
		case c == '.' && !seenDot && !seenExp:
			seenDot = true
		case (c == 'e' || c == 'E') && seenDigit && !seenExp:
			seenExp = true
		default:
			goto done
		}
		end++
	}
done:
	for end > 0 {
		if f, err := strconv.ParseFloat(s[:end], 64); err == nil {
			return f, finite(f)
		}
		end--
	}
	return 0, false
}

// Stylesheet: draw.io's styles/default.xml (the named styles and the
// defaults of vertices and edges).

var defaultVertexStyle = style{
	"shape": "label", "perimeter": "rectanglePerimeter", "fontSize": "12", "fontFamily": "Helvetica",
	"align": "center", "verticalAlign": "middle", "fillColor": "default", "strokeColor": "default", "fontColor": "default",
}

var defaultEdgeStyle = style{
	"shape": "connector", "labelBackgroundColor": "default", "endArrow": "classic", "fontSize": "11",
	"fontFamily": "Helvetica", "align": "center", "verticalAlign": "middle", "rounded": "1",
	"strokeColor": "default", "fontColor": "default",
}

var namedStyles = func() map[string]style {
	text := style{"fillColor": "none", "gradientColor": "none", "strokeColor": "none", "align": "left", "verticalAlign": "top"}
	label := style{"fontStyle": "1", "align": "left", "verticalAlign": "middle", "spacing": "2", "spacingLeft": "52",
		"imageWidth": "42", "imageHeight": "42", "rounded": "1"}
	image := style{"shape": "image", "labelBackgroundColor": "default", "verticalAlign": "top", "verticalLabelPosition": "bottom"}
	fancy := style{"shadow": "1", "glass": "1"}
	extend := func(base style, kv ...string) style {
		s := base.clone()
		for i := 0; i+1 < len(kv); i += 2 {
			s[kv[i]] = kv[i+1]
		}
		return s
	}
	m := map[string]style{
		"text":      text,
		"edgeLabel": extend(text, "labelBackgroundColor", "default", "fontSize", "11"),
		"label":     label,
		"icon": extend(label, "align", "center", "imageAlign", "center", "verticalLabelPosition", "bottom",
			"verticalAlign", "top", "labelBackgroundColor", "default", "spacing", "0", "spacingLeft", "0",
			"spacingTop", "6", "fontStyle", "0", "imageWidth", "48", "imageHeight", "48"),
		"swimlane":     {"shape": "swimlane", "fontSize": "12", "fontStyle": "1", "startSize": "23"},
		"group":        {"verticalAlign": "top", "fillColor": "none", "strokeColor": "none", "gradientColor": "none", "pointerEvents": "0"},
		"ellipse":      {"shape": "ellipse", "perimeter": "ellipsePerimeter"},
		"rhombus":      {"shape": "rhombus", "perimeter": "rhombusPerimeter"},
		"triangle":     {"shape": "triangle", "perimeter": "trianglePerimeter"},
		"line":         {"shape": "line", "strokeWidth": "4", "labelBackgroundColor": "default", "verticalAlign": "top", "spacingTop": "8"},
		"image":        image,
		"roundImage":   extend(image, "perimeter", "ellipsePerimeter"),
		"rhombusImage": extend(image, "perimeter", "rhombusPerimeter"),
		"arrow":        {"shape": "arrow", "edgeStyle": "none", "fillColor": "default"},
		"fancy":        fancy,
	}
	for _, c := range []struct{ name, grad, fill, stroke string }{
		{"gray", "#B3B3B3", "#F5F5F5", "#666666"},
		{"blue", "#7EA6E0", "#DAE8FC", "#6C8EBF"},
		{"green", "#97D077", "#D5E8D4", "#82B366"},
		{"turquoise", "#67AB9F", "#D5E8D4", "#6A9153"},
		{"yellow", "#FFD966", "#FFF2CC", "#D6B656"},
		{"orange", "#FFA500", "#FFCD28", "#D79B00"},
		{"red", "#EA6B66", "#F8CECC", "#B85450"},
		{"pink", "#B5739D", "#E6D0DE", "#996185"},
		{"purple", "#8C6C9C", "#E1D5E7", "#9673A6"},
	} {
		m[c.name] = extend(fancy, "gradientColor", c.grad, "fillColor", c.fill, "strokeColor", c.stroke)
		m["plain-"+c.name] = style{"gradientColor": c.grad, "fillColor": c.fill, "strokeColor": c.stroke}
	}
	return m
}()

// parseStyle resolves a cell style string against the default style of
// the cell's kind (mxStylesheet.getCellStyle and Graph.postProcessCellStyle):
// "key=value" entries override, a bare name merges a named style, and
// "none" removes a key. A leading ";" drops the default style.
func parseStyle(str string, edge bool) style {
	base := defaultVertexStyle
	if edge {
		base = defaultEdgeStyle
	}
	var s style
	if str == "" || str[0] != ';' {
		s = base.clone()
	} else {
		s = style{}
	}
	for _, tok := range strings.Split(str, ";") {
		if tok == "" {
			continue
		}
		if k, v, ok := strings.Cut(tok, "="); ok {
			if v == "none" {
				delete(s, k)
			} else {
				s[k] = v
			}
			continue
		}
		if named, ok := namedStyles[tok]; ok {
			for k, v := range named {
				s[k] = v
			}
		}
	}
	replaceDefaultColors(s)
	return s
}

// replaceDefaultColors resolves "default" colors to the light theme's
// shape colors (Graph.replaceDefaultColors): white backgrounds, black
// foregrounds; "defaultFillColor=invert" and the like swap them.
func replaceDefaultColors(s style) {
	for _, k := range []struct {
		key string
		bg  bool
	}{
		{"fontColor", false}, {"fillColor", true}, {"gradientColor", false}, {"strokeColor", false},
		{"imageBorder", false}, {"imageBackground", true}, {"labelBorderColor", false},
		{"swimlaneFillColor", true}, {"labelBackgroundColor", true},
	} {
		if s[k.key] != "default" {
			continue
		}
		bg := k.bg
		if s["default"+strings.ToUpper(k.key[:1])+k.key[1:]] == "invert" {
			bg = !bg
		}
		if bg {
			s[k.key] = "#ffffff"
		} else {
			s[k.key] = "#000000"
		}
	}
}
