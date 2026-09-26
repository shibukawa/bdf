package xlsx

import (
	"encoding/xml"
	"sort"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// Charts are drawn from the values cached in the chart part. Excel writes
// them, but other producers often leave them out and let Excel compute
// them from the formulas of the series ('Sheet 1'!$B$2:$B$9). Missing
// caches are filled here from the cells of the workbook before a chart is
// drawn.

// fillChartCaches fills the missing value caches of the charts a drawing
// part relates to.
func (c *converter) fillChartCaches(drawingPart string) {
	rels := c.pkg.Rels(drawingPart)
	ids := make([]string, 0, len(rels))
	for id := range rels {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		r := rels[id]
		if r.External || !strings.HasSuffix(r.Type, "/chart") || c.chartsFilled[r.Target] {
			continue
		}
		c.chartsFilled[r.Target] = true
		root, err := c.pkg.XML(r.Target)
		if err != nil {
			continue
		}
		c.fillRefs(root)
	}
}

func (c *converter) fillRefs(n *ooxml.Node) {
	for i, k := range n.Kids {
		switch k.Name {
		case "numRef", "strRef":
			if k.Name == "numRef" && (n.Name == "cat" || n.Name == "xVal") && k.Child("numCache") == nil {
				// categories given as numbers that are text in the sheet are
				// text categories
				if vals, ok := c.refValues(k.Child("f").Content(), true); ok && !anyValue(vals) {
					k = &ooxml.Node{Space: k.Space, Name: "strRef", Attrs: k.Attrs, Kids: k.Kids}
					n.Kids[i] = k
				}
			}
			cacheName := "numCache"
			if k.Name == "strRef" {
				cacheName = "strCache"
			}
			if cache := k.Child(cacheName); cache != nil && len(cache.Children("pt")) > 0 {
				continue
			}
			vals, ok := c.refValues(k.Child("f").Content(), k.Name == "numRef")
			if !ok {
				continue
			}
			var kids []*ooxml.Node
			for _, e := range k.Kids {
				if e.Name != cacheName {
					kids = append(kids, e)
				}
			}
			cache := &ooxml.Node{Name: cacheName}
			if k.Name == "numRef" {
				cache.Kids = append(cache.Kids, &ooxml.Node{Name: "formatCode", Text: "General"})
			}
			cache.Kids = append(cache.Kids, &ooxml.Node{Name: "ptCount", Attrs: attrs("val", strconv.Itoa(len(vals)))})
			for i, v := range vals {
				if v == nil {
					continue
				}
				cache.Kids = append(cache.Kids, &ooxml.Node{Name: "pt", Attrs: []xml.Attr{{Name: xml.Name{Local: "idx"}, Value: strconv.Itoa(i)}},
					Kids: []*ooxml.Node{{Name: "v", Text: *v}}})
			}
			k.Kids = append(kids, cache)
		default:
			c.fillRefs(k)
		}
	}
}

// refValues evaluates a reference to a range of one sheet into cell values
// (nil for empty cells): numbers as numbers for numeric caches, the
// displayed text otherwise.
func (c *converter) refValues(f string, numeric bool) ([]*string, bool) {
	f = strings.TrimSpace(f)
	i := strings.LastIndexByte(f, '!')
	if i < 0 {
		return nil, false
	}
	name := f[:i]
	if strings.HasPrefix(name, "'") && strings.HasSuffix(name, "'") && len(name) >= 2 {
		name = strings.ReplaceAll(name[1:len(name)-1], "''", "'")
	}
	rg, ok := parseRange(f[i+1:])
	if !ok || (rg.r1-rg.r0+1)*(rg.c1-rg.c0+1) > 100000 {
		return nil, false
	}
	ws := c.sheetNamed(name)
	if ws == nil {
		return nil, false
	}
	var out []*string
	for r := rg.r0; r <= rg.r1; r++ {
		for col := rg.c0; col <= rg.c1; col++ {
			cl := ws.cellAt(r, col)
			if cl == nil || cl.kind == cellBlank {
				out = append(out, nil)
				continue
			}
			var s string
			switch {
			case numeric && (cl.kind == cellNum || cl.kind == cellBool):
				s = strconv.FormatFloat(cl.num, 'g', -1, 64)
			case numeric:
				out = append(out, nil)
				continue
			case cl.kind == cellNum:
				s = c.st.numFormat(c.st.xf(cl.style).numFmt).formatNumber(cl.num, c.date1904).text()
			case cl.kind == cellBool:
				s = "FALSE"
				if cl.num != 0 {
					s = "TRUE"
				}
			default:
				s = cl.text.plain
			}
			out = append(out, &s)
		}
	}
	return out, true
}

// sheetNamed returns the parsed worksheet with a name (parsed once).
func (c *converter) sheetNamed(name string) *worksheet {
	if ws, ok := c.parsed[name]; ok {
		return ws
	}
	var ws *worksheet
	for _, s := range c.sheets {
		if s.name == name && !s.chart {
			if w, err := c.readWorksheet(s.part); err == nil {
				w.name = s.name
				w.sortRows()
				ws = w
			}
			break
		}
	}
	c.parsed[name] = ws
	return ws
}

func anyValue(vals []*string) bool {
	for _, v := range vals {
		if v != nil {
			return true
		}
	}
	return false
}
