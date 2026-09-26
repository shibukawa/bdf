package drawio

import (
	"strings"
)

// cell is an mxCell of a graph model.
type cell struct {
	id        string
	value     string            // the label (HTML when the style says html=1)
	attrs     map[string]string // attributes of a UserObject/object wrapper, nil for plain cells
	styleStr  string
	vertex    bool
	edge      bool
	visible   bool
	collapsed bool

	parentID, sourceID, targetID string
	parent, source, target       *cell
	children                     []*cell
	geo                          *geometry

	order int // position in document order
}

// geometry is an mxGeometry.
type geometry struct {
	x, y, w, h               float64
	relative                 bool
	sourcePoint, targetPoint *point
	offset                   *point
	points                   []point
	alternateBounds          *rect
}

// model is a parsed mxGraphModel.
type model struct {
	root    *cell
	cells   map[string]*cell
	attrs   map[string]string // mxGraphModel attributes (background, shadow, …)
	ordered []*cell
}

// parseModel builds the cell tree of an mxGraphModel element.
func parseModel(n *xmlNode) *model {
	m := &model{cells: map[string]*cell{}, attrs: n.attrs}
	if m.attrs == nil {
		m.attrs = map[string]string{}
	}
	root := n.child("root")
	for _, k := range root.kids {
		var c *cell
		switch k.name {
		case "mxCell":
			c = cellFrom(k, k)
		case "UserObject", "object":
			inner := k.child("mxCell")
			if inner == nil {
				inner = &xmlNode{name: "mxCell", attrs: map[string]string{}}
			}
			c = cellFrom(k, inner)
			c.attrs = map[string]string{}
			for a, v := range k.attrs {
				if a != "id" {
					c.attrs[a] = v
				}
			}
			c.value = k.attrs["label"]
		default:
			continue
		}
		if c.id == "" {
			continue
		}
		if _, dup := m.cells[c.id]; dup {
			continue
		}
		c.order = len(m.ordered)
		m.cells[c.id] = c
		m.ordered = append(m.ordered, c)
	}
	for _, c := range m.ordered {
		if p, ok := m.cells[c.parentID]; ok && p != c {
			c.parent = p
			p.children = append(p.children, c)
		} else if m.root == nil && c.parentID == "" {
			m.root = c
		}
		c.source = m.cells[c.sourceID]
		c.target = m.cells[c.targetID]
	}
	if m.root == nil {
		// no cell without a parent: take the first one
		for _, c := range m.ordered {
			if c.parent == nil {
				m.root = c
				break
			}
		}
	}
	// break parent cycles: cells not reachable from the root are dropped
	return m
}

func cellFrom(outer, inner *xmlNode) *cell {
	c := &cell{
		id:        outer.attrs["id"],
		value:     inner.attrs["value"],
		styleStr:  inner.attrs["style"],
		vertex:    inner.attrs["vertex"] == "1",
		edge:      inner.attrs["edge"] == "1",
		visible:   inner.attrs["visible"] != "0",
		collapsed: inner.attrs["collapsed"] == "1",
		parentID:  inner.attrs["parent"],
		sourceID:  inner.attrs["source"],
		targetID:  inner.attrs["target"],
	}
	if g := inner.child("mxGeometry"); g != nil {
		c.geo = parseGeometry(g)
	}
	return c
}

func parseGeometry(n *xmlNode) *geometry {
	g := &geometry{
		x: attrNum(n, "x", 0), y: attrNum(n, "y", 0),
		w: attrNum(n, "width", 0), h: attrNum(n, "height", 0),
		relative: n.attrs["relative"] == "1",
	}
	for _, k := range n.kids {
		switch k.name {
		case "mxPoint":
			p := point{attrNum(k, "x", 0), attrNum(k, "y", 0)}
			switch k.attrs["as"] {
			case "sourcePoint":
				g.sourcePoint = &p
			case "targetPoint":
				g.targetPoint = &p
			case "offset":
				g.offset = &p
			}
		case "Array":
			if k.attrs["as"] == "points" {
				for _, pt := range k.childrenNamed("mxPoint") {
					g.points = append(g.points, point{attrNum(pt, "x", 0), attrNum(pt, "y", 0)})
				}
			}
		case "mxRectangle":
			if k.attrs["as"] == "alternateBounds" {
				g.alternateBounds = &rect{attrNum(k, "x", 0), attrNum(k, "y", 0), attrNum(k, "width", 0), attrNum(k, "height", 0)}
			}
		}
	}
	return g
}

func attrNum(n *xmlNode, name string, def float64) float64 {
	v, ok := n.attrs[name]
	if !ok {
		return def
	}
	if f, ok := parseFloat(v); ok {
		return f
	}
	return def
}

// label returns the text of a cell's label with placeholders replaced
// (Graph.getLabel / replacePlaceholders): %name% stands for an attribute
// of the cell or, failing that, of its ancestors; %page%, %pagenumber% and
// %pagecount% for the page.
func (m *model) label(c *cell, st style, pg *pageInfo) string {
	v := c.value
	if c.attrs != nil && (c.attrs["placeholders"] == "1" || st.is("placeholders")) && strings.Contains(v, "%") {
		v = m.replacePlaceholders(c, v, pg)
	}
	return v
}

// pageInfo identifies the page being converted, for placeholders and links.
type pageInfo struct {
	name   string
	number int // 1-based
	count  int
}

func (m *model) replacePlaceholders(c *cell, s string, pg *pageInfo) string {
	var b strings.Builder
	for {
		i := strings.IndexByte(s, '%')
		if i < 0 {
			b.WriteString(s)
			break
		}
		j := strings.IndexByte(s[i+1:], '%')
		if j < 0 {
			b.WriteString(s)
			break
		}
		b.WriteString(s[:i])
		name := s[i+1 : i+1+j]
		if val, ok := m.placeholder(c, name, pg); ok {
			b.WriteString(val)
			s = s[i+2+j:]
		} else {
			b.WriteByte('%')
			s = s[i+1:]
		}
	}
	return b.String()
}

func (m *model) placeholder(c *cell, name string, pg *pageInfo) (string, bool) {
	if name == "" {
		return "%", true
	}
	for p := c; p != nil; p = p.parent {
		if p.attrs != nil {
			if v, ok := p.attrs[name]; ok {
				return v, true
			}
		}
	}
	if pg != nil {
		switch name {
		case "page":
			return pg.name, true
		case "pagenumber":
			return itoa(pg.number), true
		case "pagecount":
			return itoa(pg.count), true
		}
	}
	return "", false
}
