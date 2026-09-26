package visio

import (
	"fmt"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// Visio 2013 and later save drawings as Open Packaging Conventions packages
// ([MS-VSDX]): the document part (styles, colors, fonts), the masters and
// pages parts that list the master and page parts, and a DrawingML theme.
// Sheets are Cell, Section and Row elements.

// readVSDX reads a .vsdx (.vsdm, .vstx …) package.
func readVSDX(p *ooxml.Package, warn func(key, msg string)) (*document, error) {
	d := &document{pkg: p, styles: map[int]*style{}, colors: map[int]bdf.Color{}, faces: map[int]string{},
		masters: map[int]*master{}, warn: warn}
	docPart := "visio/document.xml"
	if r, ok := p.RelOfType("", "/document"); ok {
		docPart = r.Target
	}
	root, err := p.XML(docPart)
	if err != nil {
		return nil, err
	}
	if root.Name != "VisioDocument" {
		return nil, fmt.Errorf("%s is not a Visio document part", docPart)
	}
	l := &loader{d: d, read: readVSDXSheet}
	l.readColors(root.Child("Colors"))
	for _, f := range root.Path("FaceNames").Children("FaceName") {
		if id := attrInt(f, "ID", -1); id >= 0 {
			d.faces[id] = f.AttrStr("NameU", f.AttrStr("Name", ""))
		}
	}
	l.readStyles(root.Path("StyleSheets").Children("StyleSheet"))
	d.dc = p.CoreProperties()
	if r, ok := p.RelOfType(docPart, "/theme"); ok {
		if n, err := p.XML(r.Target); err == nil {
			d.theme = parseTheme(n)
		} else {
			warn("theme", fmt.Sprintf("theme: %v", err))
		}
	}

	if r, ok := p.RelOfType(docPart, "/masters"); ok {
		if ms, err := p.XML(r.Target); err == nil {
			for _, n := range ms.Children("Master") {
				m := &master{id: attrInt(n, "ID", -1), name: n.AttrStr("NameU", n.AttrStr("Name", ""))}
				if t, ok := p.Target(r.Target, n.Child("Rel").RelID("id")); ok {
					m.part = t.Target
					l.part = t.Target
					if mc, err := p.XML(t.Target); err == nil {
						l.readMaster(m, mc.Child("Shapes"))
					} else {
						warn("master", fmt.Sprintf("master %s: %v", m.name, err))
					}
				}
				if m.byID == nil {
					m.byID = map[int]*shape{}
				}
				d.masters[m.id] = m
			}
		}
	}

	r, ok := p.RelOfType(docPart, "/pages")
	if !ok {
		return nil, fmt.Errorf("the document has no pages part")
	}
	ps, err := p.XML(r.Target)
	if err != nil {
		return nil, err
	}
	for _, n := range ps.Children("Page") {
		pg := &page{id: attrInt(n, "ID", -1), name: n.AttrStr("Name", n.AttrStr("NameU", "")),
			background: n.AttrBool("Background", false), backPage: attrInt(n, "BackPage", -1)}
		sh := newSheet()
		readVSDXSheet(&sh, n.Child("PageSheet"))
		pg.sheet = &sh
		if t, ok := p.Target(r.Target, n.Child("Rel").RelID("id")); ok {
			pg.part = t.Target
			l.part = t.Target
			if pc, err := p.XML(t.Target); err == nil {
				pg.shapes = l.readShapes(pc.Child("Shapes"), nil)
			} else {
				warn("page", fmt.Sprintf("page %s: %v", pg.name, err))
			}
			// a page may have a theme of its own
			if tr, ok := p.RelOfType(t.Target, "/theme"); ok {
				if tn, err := p.XML(tr.Target); err == nil {
					pg.theme = parseTheme(tn)
				}
			}
		}
		d.pages = append(d.pages, pg)
	}
	return d, nil
}

// readVSDXSheet reads the Cell, Section, Text and ForeignData children of
// a sheet element.
func readVSDXSheet(sh *sheet, n *ooxml.Node) {
	for _, k := range n.Elements() {
		switch k.Name {
		case "Cell":
			sh.cells[k.AttrStr("N", "")] = vsdxCell(k)
		case "Section":
			name := k.AttrStr("N", "")
			ix := -1
			if name == "Geometry" {
				ix = attrInt(k, "IX", 0)
			}
			sec := sh.section(name, ix)
			sec.del = k.AttrBool("Del", false)
			for _, e := range k.Elements() {
				switch e.Name {
				case "Cell":
					sec.cells[e.AttrStr("N", "")] = vsdxCell(e)
				case "Row":
					key, ok := e.Attr("IX")
					if !ok {
						key = e.AttrStr("N", "")
					}
					r := sec.row(key)
					r.t = e.AttrStr("T", r.t)
					r.del = e.AttrBool("Del", false)
					for _, c := range e.Children("Cell") {
						r.cells[c.AttrStr("N", "")] = vsdxCell(c)
					}
				}
			}
		case "Text":
			sh.text = k
		case "ForeignData":
			sh.foreign = k
		}
	}
}

func vsdxCell(n *ooxml.Node) cell {
	return cell{v: n.AttrStr("V", ""), f: n.AttrStr("F", ""), u: n.AttrStr("U", "")}
}
