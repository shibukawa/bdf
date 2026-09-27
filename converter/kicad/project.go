package kicad

import (
	"encoding/json"
	"math"
	"path"
	"regexp"
	"strings"
	"time"
)

// readProject reads the project's name and text variables.
func (c *conv) readProject(in *input) {
	name := in.pro
	if name == "" {
		name = in.sch
	}
	if name == "" {
		name = in.pcb
	}
	c.project = strings.TrimSuffix(path.Base(name), path.Ext(name))
	c.textVars = map[string]string{}
	if in.proData == nil {
		return
	}
	var pro struct {
		TextVariables map[string]string `json:"text_variables"`
		Schematic     struct {
			PageLayout string `json:"page_layout_descr_file"`
			Drawing    struct {
				LineThickness *float64 `json:"default_line_thickness"` // mils
				TextOffset    *float64 `json:"text_offset_ratio"`
				LabelSize     *float64 `json:"label_size_ratio"`
				PinSymbolSize *float64 `json:"pin_symbol_size"` // mils
				JunctionSize  *int     `json:"junction_size_choice"`
				DashRatio     *float64 `json:"dashed_lines_dash_length_ratio"`
				GapRatio      *float64 `json:"dashed_lines_gap_length_ratio"`
			} `json:"drawing"`
			Meta struct {
				Version int `json:"version"`
			} `json:"meta"`
		} `json:"schematic"`
		NetSettings struct {
			Classes []struct {
				Name      string   `json:"name"`
				WireWidth *float64 `json:"wire_width"` // mils
				BusWidth  *float64 `json:"bus_width"`
			} `json:"classes"`
		} `json:"net_settings"`
	}
	if err := json.Unmarshal(in.proData, &pro); err != nil {
		c.warnf("%s: %v", in.pro, err)
		return
	}
	for k, v := range pro.TextVariables {
		c.textVars[k] = v
	}
	c.schSheetFile = pro.Schematic.PageLayout

	// SCHEMATIC_SETTINGS, with its limits
	s := &c.sch
	d := pro.Schematic.Drawing
	mils := func(v *float64, lo, hi float64, dst *float64) {
		if v != nil && !math.IsNaN(*v) {
			*dst = math.Min(math.Max(*v, lo), hi) * 0.0254
		}
	}
	ratio := func(v *float64, dst *float64) {
		if v != nil && !math.IsNaN(*v) {
			*dst = math.Min(math.Max(*v, 0), 2)
		}
	}
	mils(d.LineThickness, 5, 1000, &s.line)
	mils(d.PinSymbolSize, 0, 1000, &s.pinSymbol)
	ratio(d.TextOffset, &s.textOffset)
	ratio(d.LabelSize, &s.labelBox)
	if d.LabelSize == nil && d.TextOffset != nil && pro.Schematic.Meta.Version < 1 {
		// settings of before the label size had its own ratio
		s.labelBox = s.textOffset
	}
	if d.DashRatio != nil && *d.DashRatio > 0 && *d.DashRatio < 1000 {
		s.dash = *d.DashRatio
	}
	if d.GapRatio != nil && *d.GapRatio > 0 && *d.GapRatio < 1000 {
		s.gap = *d.GapRatio
	}
	for _, nc := range pro.NetSettings.Classes {
		if nc.Name == "Default" {
			mils(nc.WireWidth, 1, 1000, &s.wire)
			mils(nc.BusWidth, 1, 1000, &s.bus)
		}
	}
	choice := 3
	if d.JunctionSize != nil && *d.JunctionSize >= 0 && *d.JunctionSize < len(junctionSizes) {
		choice = *d.JunctionSize
	}
	s.junction = s.wire * junctionSizes[choice]
}

// schSettings are the drawing settings of a schematic (SCHEMATIC_SETTINGS
// and the Default net class of its project), in mm.
type schSettings struct {
	line       float64 // the pen of lines and texts without a width
	wire, bus  float64 // wires and buses without a width
	junction   float64 // the diameter of junctions without one
	textOffset float64 // labels' and pins' distance from their wire, of the text height
	labelBox   float64 // global labels' margin, of the text height
	pinSymbol  float64 // the size of pin decorations (0: from the pin's texts)
	dash, gap  float64 // dashes and gaps, of the line width
}

var defaultSchSettings = schSettings{
	line:       schLineWidth,
	wire:       schLineWidth,
	bus:        schBusWidth,
	junction:   schJunctionDiam,
	textOffset: 0.15,
	labelBox:   0.375,
	pinSymbol:  0.635, // 25 mils
	dash:       12,
	gap:        3,
}

// junctionSizes are the junction sizes a project chooses from, times the
// wire width (none, smallest, small, default, large, largest).
var junctionSizes = []float64{0, 1.7, 4, 6, 9, 12}

// pinTextOffset is the distance of pins' texts from the pin: 24 mils times
// the text offset ratio, in whole mils.
func (s *schSettings) pinTextOffset() float64 {
	return math.Round(24*s.textOffset) * 0.0254
}

// readWorksheet reads the drawing sheet a project names (relative to the
// project, ${KIPRJMOD} standing for its directory), or returns the default
// one.
func (c *conv) readWorksheet(from, name string) *worksheet {
	name = strings.TrimSpace(name)
	if name == "" {
		return defaultWorksheet
	}
	name = strings.TrimPrefix(strings.TrimPrefix(name, "${KIPRJMOD}"), "$(KIPRJMOD)")
	name = strings.TrimLeft(strings.ReplaceAll(name, `\`, "/"), "/")
	if c.refs == nil {
		c.warnf("the drawing sheet %s was not given with the schematic: the default one is drawn", name)
		return defaultWorksheet
	}
	b, _, err := c.readRef(from, name)
	if err != nil {
		c.warnf("the drawing sheet %s is not among the files given: the default one is drawn", name)
		return defaultWorksheet
	}
	ws, err := parseWorksheet(b)
	if err != nil {
		c.warnf("drawing sheet %s: %v; the default one is drawn", name, err)
		return defaultWorksheet
	}
	return ws
}

var textVarRe = regexp.MustCompile(`\$\{([^{}]*)\}`)

// expand replaces the ${VARIABLE} references of a text: lookup resolves
// the names of the item's context (fields of a symbol, the page's
// variables); then the project's variables. Unknown ones stay as they are.
func (c *conv) expand(s string, lookup func(name string) (string, bool)) string {
	s = unescape(s)
	if !strings.Contains(s, "${") {
		return s
	}
	for range 4 { // variables may refer to variables
		out := textVarRe.ReplaceAllStringFunc(s, func(m string) string {
			name := m[2 : len(m)-1]
			if lookup != nil {
				if v, ok := lookup(name); ok {
					return v
				}
			}
			if v, ok := c.textVars[name]; ok {
				return v
			}
			switch name {
			case "PROJECTNAME":
				return c.project
			case "CURRENT_DATE":
				return time.Now().Format("2006-01-02")
			}
			return m
		})
		if out == s {
			break
		}
		s = out
	}
	return s
}

// kicadVersion is the text of ${KICAD_VERSION}: KiCad's name and the
// version that saved the file.
func kicadVersion(file *node) string {
	if v := file.str("generator_version"); v != "" {
		return "KiCad E.D.A. " + v
	}
	return "KiCad E.D.A."
}

// escapes are the {tokens} KiCad writes for characters that net names
// cannot hold, and shows as the characters (UnescapeString).
var escapes = map[string]string{
	"dblquote": `"`, "quote": "'", "lt": "<", "gt": ">", "backslash": `\`, "slash": "/", "bar": "|",
	"comma": ",", "colon": ":", "space": " ", "dollar": "$", "tab": "\t", "return": "\n", "brace": "{",
}

var escapeRe = regexp.MustCompile(`\{(dblquote|quote|lt|gt|backslash|slash|bar|comma|colon|space|dollar|tab|return|brace)\}`)

// unescape replaces the escapes of a text with their characters.
func unescape(s string) string {
	if !strings.Contains(s, "{") {
		return s
	}
	return escapeRe.ReplaceAllStringFunc(s, func(m string) string { return escapes[m[1:len(m)-1]] })
}
