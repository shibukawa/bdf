package visio

import (
	"fmt"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// Visio 2003 to 2010 save drawings as one XML file (DatadiagramML, the
// .vdx format): the styles, masters and pages are all in it, and cells are
// elements named after them, grouped by the section they belong to
// (<XForm><PinX F="…">1.5</PinX>…</XForm>). The groups of single cells are
// flattened into the sheet's cells and the others become sections, so that
// both formats share the sheet model.

// vdxSections maps the elements of indexed or named rows to the sections
// they make up.
var vdxSections = map[string]string{
	"Char": "Character", "Para": "Paragraph", "Tabs": "Tabs", "Geom": "Geometry", "Prop": "Property",
	"User": "User", "Hyperlink": "Hyperlink", "Connection": "Connection", "Field": "Field", "Scratch": "Scratch",
	"Control": "Control", "Act": "Actions", "Layer": "Layer", "SmartTagDef": "SmartTag",
	"FillGradient": "FillGradient", "LineGradient": "LineGradient",
}

// readVDX reads an XML drawing.
func readVDX(data []byte, warn func(key, msg string)) (*document, error) {
	root, err := ooxml.Parse(data)
	if err != nil {
		return nil, err
	}
	if root.Name != "VisioDocument" {
		return nil, fmt.Errorf("not a Visio XML drawing")
	}
	d := &document{styles: map[int]*style{}, colors: map[int]bdf.Color{}, faces: map[int]string{},
		masters: map[int]*master{}, warn: warn}
	l := &loader{d: d, read: readVDXSheet}
	l.readColors(root.Child("Colors"))
	for _, f := range root.Child("FaceNames").Children("FaceName") {
		d.faces[attrInt(f, "ID", -1)] = f.AttrStr("Name", f.AttrStr("NameU", ""))
	}
	l.readStyles(root.Child("StyleSheets").Children("StyleSheet"))
	d.dc = vdxProperties(root.Child("DocumentProperties"))
	for _, n := range root.Child("Masters").Children("Master") {
		m := &master{id: attrInt(n, "ID", -1), name: n.AttrStr("NameU", n.AttrStr("Name", ""))}
		l.readMaster(m, n.Child("Shapes"))
		d.masters[m.id] = m
	}
	for _, n := range root.Child("Pages").Children("Page") {
		pg := &page{id: attrInt(n, "ID", -1), name: n.AttrStr("Name", n.AttrStr("NameU", "")),
			background: n.AttrBool("Background", false), backPage: attrInt(n, "BackPage", -1)}
		sh := newSheet()
		readVDXSheet(&sh, n.Child("PageSheet"))
		pg.sheet = &sh
		pg.shapes = l.readShapes(n.Child("Shapes"), nil)
		d.pages = append(d.pages, pg)
	}
	if len(d.pages) == 0 {
		return nil, fmt.Errorf("the drawing has no pages")
	}
	return d, nil
}

// vdxProperties reads the document properties as Dublin Core.
func vdxProperties(n *ooxml.Node) bdf.DublinCore {
	var dc bdf.DublinCore
	for _, e := range []struct {
		f    *bdf.DCValues
		name string
	}{
		{&dc.Title, "Title"}, {&dc.Creator, "Creator"}, {&dc.Subject, "Subject"}, {&dc.Description, "Desc"},
		{&dc.Created, "TimeCreated"}, {&dc.Modified, "TimeSaved"},
	} {
		if s := strings.TrimSpace(n.Child(e.name).Content()); s != "" {
			*e.f = append(*e.f, s)
		}
	}
	dc.Subject = append(dc.Subject, bdf.SplitKeywords(n.Child("Keywords").Content())...)
	return dc
}

// isCellElement reports whether an element is a cell: it holds a value and
// no elements.
func isCellElement(n *ooxml.Node) bool { return len(n.Kids) == 0 }

func vdxCell(n *ooxml.Node) cell {
	return cell{v: strings.TrimSpace(n.Content()), f: n.AttrStr("F", ""), u: n.AttrStr("Unit", "")}
}

// rowKey returns the key of an indexed (IX) or named (NameU) row.
func rowKey(n *ooxml.Node) (string, bool) {
	if v, ok := n.Attr("IX"); ok {
		return v, true
	}
	if v, ok := n.Attr("NameU"); ok {
		return v, true
	}
	if v, ok := n.Attr("Name"); ok {
		return v, true
	}
	return "", false
}

// readVDXSheet reads the children of a sheet element of an XML drawing.
func readVDXSheet(sh *sheet, n *ooxml.Node) {
	for _, k := range n.Elements() {
		switch k.Name {
		case "Shapes", "Connects":
			continue
		case "Text":
			sh.text = k
			continue
		case "ForeignData":
			sh.foreign = k
			continue
		}
		if isCellElement(k) {
			sh.cells[k.Name] = vdxCell(k)
			continue
		}
		name, listed := vdxSections[k.Name]
		key, hasKey := rowKey(k)
		switch {
		case k.Name == "Geom":
			sec := sh.section("Geometry", attrInt(k, "IX", 0))
			sec.del = k.AttrBool("Del", false)
			for _, e := range k.Elements() {
				if isCellElement(e) {
					sec.cells[e.Name] = vdxCell(e)
					continue
				}
				rk, _ := rowKey(e)
				r := sec.row(rk)
				r.t = e.Name
				r.del = e.AttrBool("Del", false)
				for _, c := range e.Elements() {
					r.cells[c.Name] = vdxCell(c)
				}
			}
		case k.Name == "Tabs":
			// the stops of a Tabs row are Tab elements
			sec := sh.section("Tabs", -1)
			r := sec.row(key)
			for _, t := range k.Children("Tab") {
				i := attrInt(t, "IX", 0) + 1
				for _, c := range t.Elements() {
					r.cells[c.Name+itoa(i)] = vdxCell(c)
				}
			}
		case listed || hasKey:
			if !listed {
				name = k.Name
			}
			sec := sh.section(name, -1)
			r := sec.row(key)
			r.del = k.AttrBool("Del", false)
			for _, c := range k.Elements() {
				r.cells[c.Name] = vdxCell(c)
			}
		default:
			// a group of single cells (XForm, Line, Fill, TextBlock …)
			for _, c := range k.Elements() {
				if isCellElement(c) {
					sh.cells[c.Name] = vdxCell(c)
				}
			}
		}
	}
}
