package sxf

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
)

// The predefined colors, line types and line widths of SXF (the common
// predefined elements, 共通既定義要素編).
var (
	colorNames = []string{"black", "red", "green", "blue", "yellow", "magenta", "cyan", "white",
		"deeppink", "brown", "orange", "lightgreen", "lightblue", "lavender", "lightgray", "darkgray"}
	colorRGB = [][3]uint8{{0, 0, 0}, {255, 0, 0}, {0, 255, 0}, {0, 0, 255}, {255, 255, 0}, {255, 0, 255}, {0, 255, 255},
		{255, 255, 255}, {192, 0, 128}, {192, 128, 64}, {255, 128, 0}, {128, 192, 128}, {0, 128, 255}, {128, 64, 255},
		{192, 192, 192}, {128, 128, 128}}
	lineTypeNames = []string{"continuous", "dashed", "dashed spaced", "long dashed dotted", "long dashed double-dotted",
		"long dashed triplicate-dotted", "dotted", "chain", "chain double dash", "dashed dotted", "double-dashed dotted",
		"dashed double-dotted", "double-dashed double-dotted", "dashed triplicate-dotted", "double-dashed triplicate-dotted"}
	// the patterns (dash, gap, ...) in mm for lines 0.5 mm wide; they
	// scale with the width (JIS Z 8312)
	lineTypePitch = [][]float64{nil, {6, 1.5}, {6, 6}, {12, 1.5, 0.25, 1.5}, {12, 1.5, 0.25, 1.5, 0.25, 1.5},
		{12, 1.5, 0.25, 1.5, 0.25, 1.5, 0.25, 1.5}, {0.25, 1.5}, {12, 1.5, 3.5, 1.5}, {12, 1.5, 3.5, 1.5, 3.5, 1.5},
		{6, 1.5, 0.25, 1.5}, {6, 1.5, 6, 1.5, 0.25, 1.5}, {6, 1.5, 0.25, 1.5, 0.25, 1.5}, {6, 1.5, 6, 1.5, 0.25, 1.5, 0.25, 1.5},
		{6, 1.5, 0.25, 1.5, 0.25, 1.5, 0.25, 1.5}, {6, 1.5, 6, 1.5, 0.25, 1.5, 0.25, 1.5, 0.25, 1.5}}
	widths = []float64{0.13, 0.18, 0.25, 0.35, 0.5, 0.7, 1, 1.4, 2}
)

type layer struct {
	name    string
	visible bool
}

type sfig struct {
	name  string
	flag  int // 1, 2: parts in mathematical and geodetic coordinates; 3: group; 4: part
	items []*feature
}

type ccurve struct {
	color, typ, width int
	visible           bool
	members           []*feature
}

type sheet struct {
	name   string
	typ    int
	orient int
	w, h   float64
	items  []*feature
}

// model is what an SFC file defines, in the order it defines it.
type model struct {
	colors   map[int]bdf.Color // user colors, from 17
	types    map[int][]float64 // user line types, from 17 (mm)
	widths   map[int]float64   // user widths, from 11
	fonts    []string          // font codes from 1
	layers   []layer           // layer codes from 1
	sfigs    map[string]*sfig
	curves   []*ccurve // composite curves from 1
	sheet    *sheet
	title    map[string]string
	bg       *bdf.Color // the background of the drawing, from an attribute
	unknown  map[string]int
	versions string
}

// buildModel reads the features of an SFC file: the tables they define,
// and the elements of each composite curve, compound figure and the sheet,
// which are the elements written since the previous one (SFC 1-3).
func buildModel(fs []*feature) *model {
	m := &model{colors: map[int]bdf.Color{}, types: map[int][]float64{}, widths: map[int]float64{},
		sfigs: map[string]*sfig{}, title: map[string]string{}, unknown: map[string]int{}}
	var pending []*feature
	for _, f := range fs {
		switch f.name {
		case "pre_defined_colour_feature", "pre_defined_font_feature":
			// fixed codes
		case "user_defined_colour_feature":
			m.colors[17+len(m.colors)] = bdf.RGB(uint8(f.int(0)), uint8(f.int(1)), uint8(f.int(2)))
		case "user_defined_font_feature":
			m.types[17+len(m.types)] = f.list(2)
		case "width_feature":
			w := f.num(0)
			predefined := false
			for _, p := range widths {
				if math.Abs(p-w) < 1e-6 {
					predefined = true
				}
			}
			if !predefined {
				m.widths[11+len(m.widths)] = w
			}
		case "text_font_feature":
			m.fonts = append(m.fonts, strings.TrimSpace(f.str(0)))
		case "layer_feature":
			m.layers = append(m.layers, layer{name: f.str(0), visible: f.int(1) != 0})
		case "drawing_attribute_feature":
			for i, k := range []string{"project", "construction", "contract", "title", "number", "type", "scale", "year", "month", "day", "contractor", "owner"} {
				m.title[k] = strings.TrimSpace(f.str(i))
			}
		case "composite_curve_org_feature":
			m.curves = append(m.curves, &ccurve{color: f.int(0), typ: f.int(1), width: f.int(2), visible: f.int(3) != 0, members: pending})
			pending = nil
		case "sfig_org_feature":
			name := f.str(0)
			m.sfigs[name] = &sfig{name: name, flag: f.int(1), items: pending}
			if c, ok := backgroundAttribute(name); ok {
				m.bg = &c
			}
			pending = nil
		case "drawing_sheet_feature":
			m.sheet = &sheet{name: f.str(0), typ: f.int(1), orient: f.int(2), w: f.num(3), h: f.num(4), items: pending}
			pending = nil
		default:
			pending = append(pending, f)
		}
	}
	if m.sheet == nil && len(pending) > 0 {
		// no sheet: what is left is drawn on a sheet of its extent
		m.sheet = &sheet{typ: -1, items: pending}
	}
	return m
}

// backgroundAttribute reads the background color attribute that CAD
// programs attach to a drawing as the name of a group:
// "$$ATRU$$n$$背景色$$色$$R_G_B".
func backgroundAttribute(name string) (bdf.Color, bool) {
	if !strings.HasPrefix(name, "$$ATRU$$") {
		return 0, false
	}
	parts := strings.Split(name, "$$")
	for i := 0; i+2 < len(parts); i++ {
		if parts[i] == "背景色" {
			rgb := strings.Split(parts[i+2], "_")
			if len(rgb) != 3 {
				return 0, false
			}
			var v [3]uint8
			for k, s := range rgb {
				n, err := strconv.Atoi(s)
				if err != nil {
					return 0, false
				}
				v[k] = uint8(n)
			}
			return bdf.RGB(v[0], v[1], v[2]), true
		}
	}
	return 0, false
}

// sheetSize returns the size of a sheet in mm (A0 to A4, portrait or
// landscape, or free).
func (s *sheet) size() (w, h float64, ok bool) {
	a := [][2]float64{{841, 1189}, {594, 841}, {420, 594}, {297, 420}, {210, 297}}
	if s.typ >= 0 && s.typ <= 4 {
		w, h = a[s.typ][0], a[s.typ][1]
		if s.orient == 1 {
			w, h = h, w
		}
		return w, h, true
	}
	if s.w > 0 && s.h > 0 {
		return s.w, s.h, true
	}
	return 0, 0, false
}

// color returns the color of a color code.
func (m *model) color(code int) (bdf.Color, bool) {
	switch {
	case code >= 1 && code <= 16:
		c := colorRGB[code-1]
		return bdf.RGB(c[0], c[1], c[2]), true
	case code >= 17:
		c, ok := m.colors[code]
		return c, ok
	}
	return 0, false
}

// width returns the width in mm of a width code.
func (m *model) width(code int) float64 {
	if code >= 1 && code <= 9 {
		return widths[code-1]
	}
	if w, ok := m.widths[code]; ok {
		return w
	}
	return 0.25
}

// pattern returns the dash pattern of a line type in mm (CAD convention:
// dashes positive, gaps negative), for a line of the width w in mm.
func (m *model) pattern(code int, w float64) []float64 {
	var pitch []float64
	switch {
	case code >= 2 && code <= 15:
		scale := math.Max(w, 0.13) / 0.5
		for _, v := range lineTypePitch[code-1] {
			pitch = append(pitch, v*scale)
		}
	case code >= 17:
		pitch = m.types[code]
	}
	if len(pitch) == 0 {
		return nil
	}
	elems := make([]float64, len(pitch))
	for i, v := range pitch {
		if i%2 == 1 {
			v = -math.Abs(v)
		} else {
			v = math.Abs(v)
		}
		elems[i] = v
	}
	return elems
}

func (m *model) font(code int) string {
	if code >= 1 && code <= len(m.fonts) {
		return m.fonts[code-1]
	}
	return ""
}

func (m *model) layerVisible(code int) bool {
	if code >= 1 && code <= len(m.layers) {
		return m.layers[code-1].visible
	}
	return true
}
