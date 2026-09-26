package visio

import (
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// document is a drawing read into the sheet model.
type document struct {
	pkg     *ooxml.Package // nil for an XML drawing (.vdx)
	styles  map[int]*style
	colors  map[int]bdf.Color // the document's color table (indices from 24)
	faces   map[int]string    // font IDs of older documents
	masters map[int]*master
	pages   []*page
	theme   *theme
	dc      bdf.DublinCore
	warn    func(key, msg string)
}

// sheetReader reads the cells, sections, text and picture of a sheet
// element in one of the two formats.
type sheetReader func(sh *sheet, n *ooxml.Node)

// loader builds the shapes of a document.
type loader struct {
	d    *document
	read sheetReader
	part string // the part being read
}

func attrInt(n *ooxml.Node, name string, def int) int {
	v, ok := n.Attr(name)
	if !ok {
		return def
	}
	return atoi(v, def)
}

// readStyles reads the StyleSheet elements and links them to their parents.
func (l *loader) readStyles(list []*ooxml.Node) {
	type links struct{ line, fill, text int }
	parents := map[*style]links{}
	for _, n := range list {
		st := &style{sheet: newSheet(), id: attrInt(n, "ID", -1), name: n.AttrStr("NameU", n.AttrStr("Name", ""))}
		l.read(&st.sheet, n)
		l.d.styles[st.id] = st
		parents[st] = links{attrInt(n, "LineStyle", -1), attrInt(n, "FillStyle", -1), attrInt(n, "TextStyle", -1)}
	}
	for st, p := range parents {
		get := func(id int) *style {
			if id < 0 || id == st.id {
				return nil
			}
			return l.d.styles[id]
		}
		st.line, st.fill, st.text = get(p.line), get(p.fill), get(p.text)
	}
}

// readShapes reads the Shape elements of a Shapes element. m is the master
// whose shapes the MasterShape attributes refer to (nil outside instances).
func (l *loader) readShapes(parent *ooxml.Node, m *master) []*shape {
	var out []*shape
	for _, n := range parent.Children("Shape") {
		if s := l.readShape(n, m); s != nil {
			out = append(out, s)
		}
	}
	return out
}

func (l *loader) readShape(n *ooxml.Node, m *master) *shape {
	s := &shape{sheet: newSheet(), id: attrInt(n, "ID", -1), name: n.AttrStr("NameU", n.AttrStr("Name", "")),
		typ: n.AttrStr("Type", "Shape"), del: n.AttrBool("Del", false), part: l.part}
	l.read(&s.sheet, n)
	style := func(name string) *style {
		if id := attrInt(n, name, -1); id >= 0 {
			return l.d.styles[id]
		}
		return nil
	}
	s.lineStyle, s.fillStyle, s.textStyle = style("LineStyle"), style("FillStyle"), style("TextStyle")
	if id := attrInt(n, "Master", -1); id >= 0 {
		if mm := l.d.masters[id]; mm != nil {
			m = mm
			s.master = mm
			if len(mm.shapes) == 1 {
				s.inh = mm.shapes[0]
			} else {
				// A master of several shapes gives them to its instances
				// as subshapes.
				s.inh = &shape{sheet: newSheet(), typ: "Group", kids: mm.shapes}
			}
		}
	} else if id := attrInt(n, "MasterShape", -1); id >= 0 && m != nil {
		s.inh = m.byID[id]
	}
	if k := n.Child("Shapes"); k != nil {
		s.kids = l.readShapes(k, m)
	} else if s.inh != nil && len(s.inh.kids) > 0 {
		// subshapes it does not list are its master shape's
		for _, mk := range s.inh.kids {
			s.kids = append(s.kids, instanceOf(mk))
		}
	}
	return s
}

// instanceOf makes a shape that inherits everything from a master shape.
func instanceOf(m *shape) *shape {
	s := &shape{sheet: newSheet(), id: m.id, name: m.name, typ: m.typ, inh: m, part: m.part}
	for _, k := range m.kids {
		s.kids = append(s.kids, instanceOf(k))
	}
	return s
}

// readMaster reads the shapes of a master.
func (l *loader) readMaster(m *master, shapes *ooxml.Node) {
	m.byID = map[int]*shape{}
	m.shapes = l.readShapes(shapes, nil)
	var index func(list []*shape)
	index = func(list []*shape) {
		for _, s := range list {
			m.byID[s.id] = s
			index(s.kids)
		}
	}
	index(m.shapes)
}

// defaultColors is the color table indices 0-23 stand for.
var defaultColors = [24]bdf.Color{
	0x000000ff, 0xffffffff, 0xff0000ff, 0x00ff00ff, 0x0000ffff, 0xffff00ff, 0xff00ffff, 0x00ffffff,
	0x800000ff, 0x008000ff, 0x000080ff, 0x808000ff, 0x800080ff, 0x008080ff, 0xc0c0c0ff, 0xe6e6e6ff,
	0xcdcdcdff, 0xb3b3b3ff, 0x9a9a9aff, 0x808080ff, 0x666666ff, 0x4d4d4dff, 0x333333ff, 0x1a1a1aff,
}

// readColors reads the document's color table (ColorEntry IX RGB).
func (l *loader) readColors(colors *ooxml.Node) {
	for _, e := range colors.Children("ColorEntry") {
		if c, ok := hexColor(e.AttrStr("RGB", "")); ok {
			l.d.colors[attrInt(e, "IX", -1)] = c
		}
	}
}

// hexColor parses "#RRGGBB".
func hexColor(s string) (bdf.Color, bool) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(s) != 6 {
		return 0, false
	}
	v, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return 0, false
	}
	return bdf.Color(v<<8 | 0xff), true
}

// color reads a color cell value: "#RRGGBB" or an index into the color
// table.
func (d *document) color(v string) (bdf.Color, bool) {
	v = strings.TrimSpace(v)
	if strings.HasPrefix(v, "#") {
		return hexColor(v)
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			i = int(f)
		} else {
			return 0, false
		}
	}
	if i >= 0 && i < len(defaultColors) {
		return defaultColors[i], true
	}
	c, ok := d.colors[i]
	return c, ok
}

// fontName reads a font cell value: a family name, or the ID of a font in
// the font table of older documents.
func (d *document) fontName(v string) string {
	v = strings.TrimSpace(v)
	if i, err := strconv.Atoi(v); err == nil {
		if name, ok := d.faces[i]; ok {
			return name
		}
		return ""
	}
	return v
}

// pageByID returns a page by its ID.
func (d *document) pageByID(id int) *page {
	for _, p := range d.pages {
		if p.id == id {
			return p
		}
	}
	return nil
}
