package main

import (
	"fmt"
	"html"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// Sheets are diagrams for reviewing the mapping in draw.io: every old name
// in a grid, drawn with the style of its first palette entry (or just its
// shape name) and labeled with the name, alone (<generation>.drawio) and
// next to its target (compare-<generation>.drawio, compare-overrides.drawio).

// paletteStyle returns the style of the first palette entry of an old
// name and its size, or a bare style and the stencil's size.
func (s *sources) paletteStyle(n oldName) (string, float64, float64) {
	for _, e := range s.old {
		if parseStyle(e.style)["shape"] == n.name {
			return e.style, e.w, e.h
		}
	}
	w, h := n.w, n.h
	if w <= 0 || h <= 0 {
		w, h = 60, 60
	}
	return "shape=" + n.name + ";html=1;", w, h
}

// writeSheets writes a sheet per generation into dir.
func writeSheets(s *sources, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	byGen := map[string][]oldName{}
	var gens []string
	for _, n := range s.names {
		if byGen[n.gen] == nil {
			gens = append(gens, n.gen)
		}
		byGen[n.gen] = append(byGen[n.gen], n)
	}
	for _, g := range gens {
		var b strings.Builder
		b.WriteString(`<mxfile><diagram id="` + g + `" name="` + g + `"><mxGraphModel><root><mxCell id="0"/><mxCell id="1" parent="0"/>` + "\n")
		const cols, cw, ch, box = 10, 130.0, 150.0, 80.0
		for i, n := range byGen[g] {
			style, w, h := s.paletteStyle(n)
			sc := math.Min(1, math.Min(box/w, box/h))
			w, h = w*sc, h*sc
			x := float64(i%cols)*cw + (cw-w)/2
			y := float64(i/cols)*ch + (box-h)/2
			label := n.name[strings.LastIndex(n.name, ".")+1:]
			fmt.Fprintf(&b, `<mxCell id="c%d" value="" style="%s" vertex="1" parent="1"><mxGeometry x="%g" y="%g" width="%g" height="%g" as="geometry"/></mxCell>`+"\n",
				i, html.EscapeString(style), round2(x), round2(y), round2(w), round2(h))
			fmt.Fprintf(&b, `<mxCell id="t%d" value="%s" style="text;html=1;align=center;verticalAlign=top;fontSize=9;whiteSpace=wrap;" vertex="1" parent="1"><mxGeometry x="%g" y="%g" width="%g" height="40" as="geometry"/></mxCell>`+"\n",
				i, html.EscapeString(label), float64(i%cols)*cw, float64(i/cols)*ch+box+4, cw)
		}
		b.WriteString("</root></mxGraphModel></diagram></mxfile>\n")
		if err := os.WriteFile(filepath.Join(dir, g+".drawio"), []byte(b.String()), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func round2(v float64) float64 { return math.Round(v*100) / 100 }

// writeSamples writes the converter's test diagrams of the old icons
// (converter/drawio/testdata/aws_legacy/<generation>.drawio): a sample of
// every generation's icons with the style and size of their palette
// entries, labeled with their titles, the first row joined by edges.
// Names without a counterpart are among them; groups are not (see
// groups.drawio).
func writeSamples(s *sources, ds []decision, dir string) error {
	const perGen = 36
	byGen := map[string][]decision{}
	var gens []string
	for _, d := range ds {
		if t := d.target; t != nil && (t.kind == "group" || t.kind == "rect" || t.kind == "builtin") {
			continue
		}
		if byGen[d.name.gen] == nil {
			gens = append(gens, d.name.gen)
		}
		byGen[d.name.gen] = append(byGen[d.name.gen], d)
	}
	for _, g := range gens {
		all := byGen[g]
		var pick []decision
		unmapped := 0
		stride := max(1, len(all)/perGen)
		for i, d := range all {
			if d.target == nil && unmapped < 2 {
				unmapped++
				pick = append(pick, d)
			} else if i%stride == 0 && d.target != nil && len(pick) < perGen {
				pick = append(pick, d)
			}
		}
		var b strings.Builder
		b.WriteString("<mxfile>\n  <diagram id=\"" + g + "\" name=\"" + g + "\">\n    <mxGraphModel>\n      <root>\n" +
			"        <mxCell id=\"0\"/>\n        <mxCell id=\"1\" parent=\"0\"/>\n")
		const cols, cw, ch, box = 6, 160.0, 160.0, 80.0
		for i, d := range pick {
			style, w, h := s.paletteStyle(d.name)
			if g == "aws" {
				style += "verticalLabelPosition=bottom;verticalAlign=top;"
			}
			sc := math.Min(1, math.Min(box/w, box/h))
			w, h = round2(w*sc), round2(h*sc)
			title := d.name.display
			if es := s.oldEntries(d.name); len(es) > 0 && es[0].title != "" {
				title = es[0].title
			}
			fmt.Fprintf(&b, "        <mxCell id=\"i%d\" value=\"%s\" style=\"%s\" vertex=\"1\" parent=\"1\">\n"+
				"          <mxGeometry x=\"%g\" y=\"%g\" width=\"%g\" height=\"%g\" as=\"geometry\"/>\n        </mxCell>\n",
				i, html.EscapeString(title), html.EscapeString(style),
				round2(float64(i%cols)*cw+(cw-w)/2), round2(float64(i/cols)*ch+(box-h)/2), w, h)
		}
		for i := 0; i+1 < min(cols, len(pick)); i++ {
			fmt.Fprintf(&b, "        <mxCell id=\"e%d\" style=\"edgeStyle=orthogonalEdgeStyle;html=1;endArrow=block;endFill=1;\" edge=\"1\" parent=\"1\" source=\"i%d\" target=\"i%d\">\n"+
				"          <mxGeometry relative=\"1\" as=\"geometry\"/>\n        </mxCell>\n", i, i, i+1)
		}
		b.WriteString("      </root>\n    </mxGraphModel>\n  </diagram>\n</mxfile>\n")
		if err := os.WriteFile(filepath.Join(dir, g+".drawio"), []byte(b.String()), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// writeCompareSheets writes a diagram per generation with every old name
// next to its target (drawn with the target's palette style, fitted to
// the old bounds), for looking at the mapping in draw.io.
func writeCompareSheets(s *sources, ds []decision, dir string) error {
	byGen := map[string][]decision{}
	var gens []string
	for _, d := range ds {
		if byGen[d.name.gen] == nil {
			gens = append(gens, d.name.gen)
		}
		byGen[d.name.gen] = append(byGen[d.name.gen], d)
		if d.how == "override" {
			byGen["overrides"] = append(byGen["overrides"], d)
		}
	}
	gens = append(gens, "overrides")
	for _, g := range gens {
		var b strings.Builder
		b.WriteString(`<mxfile><diagram id="` + g + `" name="` + g + `"><mxGraphModel><root><mxCell id="0"/><mxCell id="1" parent="0"/>` + "\n")
		const cols, cw, ch, box = 6, 220.0, 130.0, 70.0
		for i, d := range byGen[g] {
			style, w, h := s.paletteStyle(d.name)
			sc := math.Min(1, math.Min(box/w, box/h))
			w, h = w*sc, h*sc
			x0 := float64(i%cols) * cw
			y0 := float64(i/cols) * ch
			fmt.Fprintf(&b, `<mxCell id="o%d" value="" style="%s" vertex="1" parent="1"><mxGeometry x="%g" y="%g" width="%g" height="%g" as="geometry"/></mxCell>`+"\n",
				i, html.EscapeString(style), round2(x0+10+(box-w)/2), round2(y0+(box-h)/2), round2(w), round2(h))
			label := d.name.name[strings.LastIndex(d.name.name, ".")+1:] + " → -"
			if t := d.target; t != nil {
				label = d.name.name[strings.LastIndex(d.name.name, ".")+1:] + " → " + t.label()
				ns := t.compareStyle(style)
				nw, nh := w, h
				if aw, ah := t.aspect(); aw > 0 {
					f := math.Min(w/aw, h/ah)
					nw, nh = aw*f, ah*f
				}
				fmt.Fprintf(&b, `<mxCell id="n%d" value="" style="%s" vertex="1" parent="1"><mxGeometry x="%g" y="%g" width="%g" height="%g" as="geometry"/></mxCell>`+"\n",
					i, html.EscapeString(ns), round2(x0+110+(box-nw)/2), round2(y0+(box-nh)/2), round2(nw), round2(nh))
			}
			fmt.Fprintf(&b, `<mxCell id="t%d" value="%s" style="text;html=1;align=center;verticalAlign=top;fontSize=9;whiteSpace=wrap;" vertex="1" parent="1"><mxGeometry x="%g" y="%g" width="%g" height="40" as="geometry"/></mxCell>`+"\n",
				i, html.EscapeString(label), x0, y0+box+4, cw-10)
		}
		b.WriteString("</root></mxGraphModel></diagram></mxfile>\n")
		if err := os.WriteFile(filepath.Join(dir, "compare-"+g+".drawio"), []byte(b.String()), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// writeTargetSheet writes a diagram of the palette entries whose labels
// (as in overrides.txt) are listed, or of all of them, for choosing
// overrides.
func writeTargetSheet(targets []*target, labels []string, path string) error {
	var b strings.Builder
	b.WriteString(`<mxfile><diagram id="aws4" name="aws4"><mxGraphModel><root><mxCell id="0"/><mxCell id="1" parent="0"/>` + "\n")
	const cols, cw, ch, box = 8, 150.0, 130.0, 70.0
	i := 0
	for _, t := range targets {
		lab := t.label() + "@" + strings.ReplaceAll(t.category, " ", "_")
		if len(labels) > 0 && !contains(labels, t.label()) && !contains(labels, lab) {
			continue
		}
		w, h := t.entry.w, t.entry.h
		sc := math.Min(box/w, box/h)
		w, h = w*sc, h*sc
		x0, y0 := float64(i%cols)*cw, float64(i/cols)*ch
		fmt.Fprintf(&b, `<mxCell id="n%d" value="" style="%s" vertex="1" parent="1"><mxGeometry x="%g" y="%g" width="%g" height="%g" as="geometry"/></mxCell>`+"\n",
			i, html.EscapeString(t.entry.style), round2(x0+(cw-w)/2), round2(y0+(box-h)/2), round2(w), round2(h))
		fmt.Fprintf(&b, `<mxCell id="t%d" value="%s" style="text;html=1;align=center;verticalAlign=top;fontSize=9;whiteSpace=wrap;" vertex="1" parent="1"><mxGeometry x="%g" y="%g" width="%g" height="40" as="geometry"/></mxCell>`+"\n",
			i, html.EscapeString(lab+" ("+t.entry.title+")"), x0, y0+box+4, cw-6)
		i++
	}
	b.WriteString("</root></mxGraphModel></diagram></mxfile>\n")
	return os.WriteFile(path, []byte(b.String()), 0o644)
}
