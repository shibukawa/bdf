package visio

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// A Visio drawing is made of ShapeSheets: the shapes of pages and masters,
// the style sheets and the page sheets. A sheet has cells, each with a
// value V in internal units (inches, radians), the formula F that computed
// it and a display unit U, at its top level and in the rows of its
// sections. Both file formats (.vsdx and the XML drawings of Visio 2003 to
// 2010) are read into this model.
//
// A cell a sheet does not set is inherited: a shape takes what it lacks from
// the master shape it is an instance of, then from the style sheets its
// LineStyle, FillStyle and TextStyle name (each covers its own group of
// cells), which inherit from their parent styles in turn. The values are
// the ones Visio last computed, so formulas are not evaluated; the cells a
// dynamic theme sets hold "Themed" instead and are resolved from the theme
// (see theme.go).

// cell is the value of a ShapeSheet cell.
type cell struct {
	v, f, u string
}

// row is a row of a section, with its type for geometry rows (MoveTo …).
type row struct {
	t     string
	del   bool
	cells map[string]cell
}

// section is a section of a sheet: cells of its own (the NoFill … of a
// geometry section) and rows by index (IX) or name (N).
type section struct {
	del   bool
	cells map[string]cell
	rows  map[string]*row
	order []string // row keys in document order
}

// secKey identifies a section: geometry sections are indexed (IX), the
// others are single (ix -1).
type secKey struct {
	name string
	ix   int
}

// sheet is the content of a ShapeSheet.
type sheet struct {
	cells    map[string]cell
	sections map[secKey]*section
	text     *ooxml.Node // Text
	foreign  *ooxml.Node // ForeignData
}

func newSheet() sheet {
	return sheet{cells: map[string]cell{}, sections: map[secKey]*section{}}
}

func (s *sheet) section(name string, ix int) *section {
	k := secKey{name, ix}
	sec := s.sections[k]
	if sec == nil {
		sec = &section{cells: map[string]cell{}, rows: map[string]*row{}}
		s.sections[k] = sec
	}
	return sec
}

func (sec *section) row(key string) *row {
	r := sec.rows[key]
	if r == nil {
		r = &row{cells: map[string]cell{}}
		sec.rows[key] = r
		sec.order = append(sec.order, key)
	}
	return r
}

// style is a style sheet.
type style struct {
	sheet
	id               int
	name             string // universal name (NameU)
	line, fill, text *style // the styles it inherits from
}

// shape is a shape of a page or a master.
type shape struct {
	sheet
	id     int
	name   string // universal name
	typ    string // Shape, Group, Foreign, Guide
	del    bool
	inh    *shape // the master shape it is an instance of
	master *master
	// its style sheets; nil when the shape names none (the master shape's
	// then apply)
	lineStyle, fillStyle, textStyle *style
	kids                            []*shape
	part                            string // the part it is in (for relationships)
}

// master is a master: the shapes its instances inherit from.
type master struct {
	id     int
	name   string
	part   string // the part its images are related to
	shapes []*shape
	byID   map[int]*shape
}

// page is a drawing page or a background page.
type page struct {
	id         int
	name       string
	background bool
	backPage   int // ID of its background page (-1: none)
	sheet      *sheet
	part       string
	shapes     []*shape
	theme      *theme
}

// maxChain bounds inheritance chains (malformed documents may have cycles).
const maxChain = 32

// category is the style sheet group a cell belongs to.
type category int

const (
	catOther category = iota
	catLine
	catFill
	catText
)

// cellCategory tells which of a shape's style sheets a cell comes from.
func cellCategory(name string) category {
	switch name {
	case "LineWeight", "LineColor", "LinePattern", "Rounding", "EndArrowSize", "BeginArrow", "EndArrow", "LineCap",
		"BeginArrowSize", "LineColorTrans", "CompoundType", "LineGradientDir", "LineGradientAngle", "LineGradientEnabled",
		"QuickStyleLineColor", "QuickStyleLineMatrix":
		return catLine
	case "FillForegnd", "FillBkgnd", "FillPattern", "ShdwForegnd", "ShdwPattern", "FillForegndTrans", "FillBkgndTrans",
		"ShdwForegndTrans", "ShapeShdwType", "ShapeShdwOffsetX", "ShapeShdwOffsetY", "ShapeShdwObliqueAngle",
		"ShapeShdwScaleFactor", "ShapeShdwBlur", "ShapeShdwShow", "FillGradientDir", "FillGradientAngle",
		"FillGradientEnabled", "RotateGradientWithShape", "UseGroupGradient",
		"QuickStyleFillColor", "QuickStyleFillMatrix", "QuickStyleEffectsMatrix", "QuickStyleShadowColor":
		return catFill
	case "LeftMargin", "RightMargin", "TopMargin", "BottomMargin", "VerticalAlign", "TextBkgnd", "DefaultTabStop",
		"TextDirection", "TextBkgndTrans", "QuickStyleFontColor", "QuickStyleFontMatrix":
		return catText
	}
	return catOther
}

// sectionCategory tells which style sheet a section's rows come from.
func sectionCategory(name string) category {
	switch name {
	case "LineGradient":
		return catLine
	case "FillGradient":
		return catFill
	case "Character", "Paragraph", "Tabs":
		return catText
	}
	return catOther
}

// styleOf returns the style sheet of a category that a shape uses: its own
// or its master shape's.
func (s *shape) styleOf(cat category) *style {
	for k, n := s, 0; k != nil && n < maxChain; k, n = k.inh, n+1 {
		var st *style
		switch cat {
		case catLine:
			st = k.lineStyle
		case catFill:
			st = k.fillStyle
		case catText:
			st = k.textStyle
		}
		if st != nil {
			return st
		}
	}
	return nil
}

func (st *style) parent(cat category) *style {
	switch cat {
	case catLine:
		return st.line
	case catFill:
		return st.fill
	case catText:
		return st.text
	}
	return nil
}

// styleChain lists the style sheets a shape inherits a category of cells
// from, most specific first. Cells of no category come from any of them.
func (s *shape) styleChain(cat category) []*style {
	if cat == catOther {
		var out []*style
		for _, c := range []category{catLine, catFill, catText} {
			out = append(out, s.styleChain(c)...)
		}
		return out
	}
	var out []*style
	for st, n := s.styleOf(cat), 0; st != nil && n < maxChain; st, n = st.parent(cat), n+1 {
		out = append(out, st)
	}
	return out
}

// get looks up a top-level cell: the shape's own, its master shape's, then
// its style sheets'.
func (s *shape) get(name string) (cell, bool) {
	for k, n := s, 0; k != nil && n < maxChain; k, n = k.inh, n+1 {
		if c, ok := k.cells[name]; ok {
			return c, true
		}
	}
	for _, st := range s.styleChain(cellCategory(name)) {
		if c, ok := st.cells[name]; ok {
			return c, true
		}
	}
	return cell{}, false
}

// formula returns the formula of a top-level cell where it is set: "Inh"
// only says that the formula is inherited.
func (s *shape) formula(name string) string {
	for k, n := s, 0; k != nil && n < maxChain; k, n = k.inh, n+1 {
		if c, ok := k.cells[name]; ok && c.f != "Inh" {
			return c.f
		}
	}
	for _, st := range s.styleChain(cellCategory(name)) {
		if c, ok := st.cells[name]; ok && c.f != "Inh" {
			return c.f
		}
	}
	return ""
}

// rowCell looks up a cell of a row of a single section (Character,
// Paragraph …) along the same chain.
func (s *shape) rowCell(sec, key, name string) (cell, bool) {
	look := func(sh *sheet) (cell, bool, bool) {
		x := sh.sections[secKey{sec, -1}]
		if x == nil {
			return cell{}, false, false
		}
		if x.del {
			return cell{}, false, true
		}
		r := x.rows[key]
		if r == nil {
			return cell{}, false, false
		}
		if r.del {
			return cell{}, false, true
		}
		c, ok := r.cells[name]
		return c, ok, false
	}
	for k, n := s, 0; k != nil && n < maxChain; k, n = k.inh, n+1 {
		if c, ok, stop := look(&k.sheet); ok || stop {
			return c, ok
		}
	}
	for _, st := range s.styleChain(sectionCategory(sec)) {
		if c, ok, stop := look(&st.sheet); ok || stop {
			return c, ok
		}
	}
	return cell{}, false
}

// rowKeys returns the keys of the rows of a single section in order, with
// the rows inherited from the master shape and the style sheets.
func (s *shape) rowKeys(sec string) []string {
	var out []string
	seen := map[string]bool{}
	add := func(sh *sheet) bool {
		x := sh.sections[secKey{sec, -1}]
		if x == nil {
			return false
		}
		if x.del {
			return true
		}
		for _, k := range x.order {
			if !seen[k] && !x.rows[k].del {
				seen[k] = true
				out = append(out, k)
			}
		}
		return false
	}
	for k, n := s, 0; k != nil && n < maxChain; k, n = k.inh, n+1 {
		if add(&k.sheet) {
			return out
		}
	}
	for _, st := range s.styleChain(sectionCategory(sec)) {
		if add(&st.sheet) {
			return out
		}
	}
	return out
}

// geoRow is a merged row of a geometry section.
type geoRow struct {
	t     string
	cells []map[string]cell // most specific first
	own   bool              // cells[0] is the shape's own
	w, h  float64           // the shape's size, for inherited formulas
}

// num returns the value of a cell of the row. An inherited value whose
// formula only scales the shape's size is evaluated with the shape's own
// size (see formula.go).
func (r *geoRow) num(name string) float64 {
	for i, m := range r.cells {
		if c, ok := m[name]; ok {
			if i > 0 || !r.own {
				if v, ok := evalSize(c.f, r.w, r.h); ok {
					return v
				}
			}
			return num(c.v)
		}
	}
	return 0
}

func (r *geoRow) str(name string) string {
	for _, m := range r.cells {
		if c, ok := m[name]; ok {
			return c.v
		}
	}
	return ""
}

// geoSection is a merged geometry section.
type geoSection struct {
	ix    int
	cells []map[string]cell
	rows  []*geoRow
}

func (g *geoSection) flag(name string) bool {
	for _, m := range g.cells {
		if c, ok := m[name]; ok {
			return num(c.v) != 0
		}
	}
	return false
}

// geometry merges the geometry sections of a shape and its master shape by
// index, and their rows by index.
func (s *shape) geometry() []*geoSection {
	var chain []*sheet
	for k, n := s, 0; k != nil && n < maxChain; k, n = k.inh, n+1 {
		chain = append(chain, &k.sheet)
	}
	var ixs []int
	seen := map[int]bool{}
	for i := len(chain) - 1; i >= 0; i-- {
		for k := range chain[i].sections {
			if k.name == "Geometry" && !seen[k.ix] {
				seen[k.ix] = true
				ixs = append(ixs, k.ix)
			}
		}
	}
	sortInts(ixs)
	var out []*geoSection
	for _, ix := range ixs {
		g := &geoSection{ix: ix}
		var secs []*section
		deleted, ownSec := false, false
		for i, sh := range chain {
			if sec := sh.sections[secKey{"Geometry", ix}]; sec != nil {
				if i == 0 {
					ownSec = true
				}
				if sec.del {
					deleted = true
					break
				}
				secs = append(secs, sec)
				g.cells = append(g.cells, sec.cells)
			}
		}
		if deleted || len(secs) == 0 {
			continue
		}
		// rows in index order; a row deleted in a more specific sheet is gone
		var keys []string
		kseen := map[string]bool{}
		for i := len(secs) - 1; i >= 0; i-- {
			for _, k := range secs[i].order {
				if !kseen[k] {
					kseen[k] = true
					keys = append(keys, k)
				}
			}
		}
		sortRowKeys(keys)
		for _, k := range keys {
			gr := &geoRow{}
			for i, sec := range secs {
				r := sec.rows[k]
				if r == nil {
					continue
				}
				if r.del {
					gr = nil
					break
				}
				if gr.t == "" {
					gr.t = r.t
				}
				if len(gr.cells) == 0 && i == 0 && ownSec {
					gr.own = true
				}
				gr.cells = append(gr.cells, r.cells)
			}
			if gr != nil && gr.t != "" {
				g.rows = append(g.rows, gr)
			}
		}
		out = append(out, g)
	}
	return out
}

// textNode returns the text of a shape (its own or its master shape's).
func (s *shape) textNode() *ooxml.Node {
	for k, n := s, 0; k != nil && n < maxChain; k, n = k.inh, n+1 {
		if k.text != nil {
			return k.text
		}
	}
	return nil
}

// foreignNode returns the embedded picture or object of a shape.
func (s *shape) foreignNode() (*ooxml.Node, *shape) {
	for k, n := s, 0; k != nil && n < maxChain; k, n = k.inh, n+1 {
		if k.foreign != nil {
			return k.foreign, k
		}
	}
	return nil, nil
}

// num parses a cell value; malformed and missing values are 0.
func num(v string) float64 {
	f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return 0
	}
	return f
}

func itoa(i int) string { return strconv.Itoa(i) }

func atoi(v string, def int) int {
	i, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		if f, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
			return int(f)
		}
		return def
	}
	return i
}

func sortInts(a []int) {
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j] < a[j-1]; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}

// sortRowKeys orders row keys numerically (named rows keep their order
// after the indexed ones).
func sortRowKeys(keys []string) {
	less := func(a, b string) bool {
		ia, ea := strconv.Atoi(a)
		ib, eb := strconv.Atoi(b)
		switch {
		case ea == nil && eb == nil:
			return ia < ib
		case ea == nil:
			return true
		}
		return false
	}
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && less(keys[j], keys[j-1]); j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
}
