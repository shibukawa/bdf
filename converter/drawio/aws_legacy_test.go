package drawio

import (
	"fmt"
	"html"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestRewriteAWSLegacyStyle(t *testing.T) {
	icon := &awsLegacyTarget{style: "shape=mxgraph.aws4.resourceIcon;resIcon=mxgraph.aws4.ec2;fillColor=#ED7100;strokeColor=#ffffff;fontColor=#232F3E;dashed=0;aspect=fixed", w: 1, h: 1}
	group := &awsLegacyTarget{style: "shape=mxgraph.aws4.group;grIcon=mxgraph.aws4.group_region;fillColor=none;strokeColor=#00A4A6;fontColor=#147EBA;dashed=1;container=1;verticalAlign=top;align=left;spacingLeft=30"}
	frame := &awsLegacyTarget{style: "rounded=1;arcSize=10", keepColors: true}
	const iconKeys = "shape=mxgraph.aws4.resourceIcon;resIcon=mxgraph.aws4.ec2;fillColor=#ED7100;strokeColor=#ffffff"
	for _, tc := range []struct {
		style string
		t     *awsLegacyTarget
		want  string
	}{
		// label and layout keys stay; the old icon's keys go
		{"outlineConnect=0;dashed=0;verticalLabelPosition=bottom;verticalAlign=top;align=center;html=1;shape=mxgraph.aws3.ec2;fillColor=#F58534;gradientColor=none;fontSize=14;rotation=30;flipH=1;spacingTop=4;whiteSpace=wrap;",
			icon, "outlineConnect=0;verticalLabelPosition=bottom;verticalAlign=top;align=center;html=1;fontSize=14;rotation=30;flipH=1;spacingTop=4;whiteSpace=wrap;" + iconKeys + ";fontColor=#232F3E;dashed=0;aspect=fixed"},
		// a font color of the cell's own stays
		{"shape=mxgraph.aws2.compute_and_networking.ec2;fontColor=#FF0000;strokeColor=none;", icon,
			"fontColor=#FF0000;" + iconKeys + ";dashed=0;aspect=fixed"},
		// the keys of the 3D shapes go; a leading ";" (no default style) and named styles stay
		{";text;shape=mxgraph.aws3d.ec2;strokeColor2=#292929;strokeWidth=1;aspect=fixed;fillColor=#ECECEC", icon,
			";text;" + iconKeys + ";fontColor=#232F3E;dashed=0;aspect=fixed"},
		// groups place their labels next to their icons
		{"dashed=0;html=1;shape=mxgraph.aws.groups.region;fillColor=none;gradientColor=none;verticalAlign=middle;align=center;spacingLeft=8;fontSize=16", group,
			"html=1;fontSize=16;shape=mxgraph.aws4.group;grIcon=mxgraph.aws4.group_region;fillColor=none;strokeColor=#00A4A6;fontColor=#147EBA;dashed=1;container=1;verticalAlign=top;align=left;spacingLeft=30"},
		// a frame keeps its colors
		{"dashed=0;html=1;shape=mxgraph.aws.groups.rrect;fillColor=#F2F2F2;gradientColor=none;strokeColor=#FF0000;", frame,
			"dashed=0;html=1;fillColor=#F2F2F2;gradientColor=none;strokeColor=#FF0000;rounded=1;arcSize=10"},
	} {
		if got := rewriteAWSLegacyStyle(tc.style, tc.t); got != tc.want {
			t.Errorf("%s\n got %s\nwant %s", tc.style, got, tc.want)
		}
	}
}

// awsLegacyModel is a model with an old icon (EC2 of the AWS17 icons,
// taller than wide) that has children: a text beside it, a relative child
// and an edge, all in the icon's coordinates; a second icon below it; an
// edge between the icons; an old group; an old icon without a counterpart.
const awsLegacyModel = `<mxGraphModel><root><mxCell id="0"/><mxCell id="1" parent="0"/>
<mxCell id="ec2" value="EC2" style="verticalLabelPosition=bottom;verticalAlign=top;html=1;shape=mxgraph.aws3.ec2;fillColor=#F58534;gradientColor=none;" vertex="1" parent="1"><mxGeometry x="10" y="20" width="76.5" height="93" as="geometry"/></mxCell>
<mxCell id="note" value="note" style="text;html=1;" vertex="1" parent="ec2"><mxGeometry x="84" y="36" width="80" height="20" as="geometry"/></mxCell>
<mxCell id="badge" value="" style="ellipse;" vertex="1" parent="ec2"><mxGeometry x="1" y="0.5" width="10" height="10" relative="1" as="geometry"><mxPoint x="-5" y="-5" as="offset"/></mxGeometry></mxCell>
<mxCell id="tick" style="endArrow=none;" edge="1" parent="ec2"><mxGeometry relative="1" as="geometry"><mxPoint x="0" y="100" as="sourcePoint"/><mxPoint x="40" y="100" as="targetPoint"/><Array as="points"><mxPoint x="20" y="110"/></Array></mxGeometry></mxCell>
<mxCell id="s3" value="S3" style="verticalLabelPosition=bottom;verticalAlign=top;html=1;shape=mxgraph.aws3.s3;fillColor=#E05243;gradientColor=none;" vertex="1" parent="1"><mxGeometry x="10" y="220" width="76.5" height="93" as="geometry"/></mxCell>
<mxCell id="e" style="endArrow=block;" edge="1" parent="1" source="ec2" target="s3"><mxGeometry relative="1" as="geometry"/></mxCell>
<mxCell id="region" value="Region" style="shape=mxgraph.aws.groups.region;fillColor=none;" vertex="1" parent="1"><mxGeometry x="300" y="0" width="200" height="150" as="geometry"/></mxCell>
<mxCell id="swf" value="SWF" style="shape=mxgraph.aws3.swf;fillColor=#D9A741;" vertex="1" parent="1"><mxGeometry x="300" y="200" width="76.5" height="93" as="geometry"/></mxCell>
</root></mxGraphModel>`

func awsLegacyTestModel(t *testing.T) *model {
	t.Helper()
	root, err := parseXML([]byte(awsLegacyModel))
	if err != nil {
		t.Fatal(err)
	}
	return parseModel(root)
}

func TestSubstituteAWSLegacy(t *testing.T) {
	var warnings []string
	c := &converter{opts: &Options{Warn: func(s string) { warnings = append(warnings, s) }}, warned: map[string]bool{}}
	before := newView(awsLegacyTestModel(t), nil, nil)
	m := awsLegacyTestModel(t)
	if n := c.substituteAWSLegacy(m); n != 3 {
		t.Errorf("%d substituted, want 3", n)
	}
	v := newView(m, nil, nil)
	st := func(v *graphView, id string) *cellState { return v.state(v.m.cells[id]) }
	near := func(a, b rect) bool {
		return math.Abs(a.x-b.x) < 1e-9 && math.Abs(a.y-b.y) < 1e-9 && math.Abs(a.w-b.w) < 1e-9 && math.Abs(a.h-b.h) < 1e-9
	}

	// the icon is the square centered in the old bounds, with its label below it
	ec2 := st(v, "ec2")
	if want := (rect{10, 28.25, 76.5, 76.5}); !near(ec2.bounds(), want) {
		t.Errorf("ec2 bounds %v, want %v", ec2.bounds(), want)
	}
	if s := ec2.style; s["shape"] != "mxgraph.aws4.resourceIcon" || s["resIcon"] != "mxgraph.aws4.ec2" ||
		s["fillColor"] != "#ED7100" || s["gradientColor"] != "" || s["verticalLabelPosition"] != "bottom" {
		t.Errorf("ec2 style %v", s)
	}
	if ec2.absoluteOffset.y != 76.5 {
		t.Errorf("label offset %v, want the fitted height", ec2.absoluteOffset)
	}
	// its children stay where they were
	for _, id := range []string{"note", "badge"} {
		if a, b := st(before, id).bounds(), st(v, id).bounds(); !near(a, b) {
			t.Errorf("%s moved from %v to %v", id, a, b)
		}
	}
	if a, b := st(before, "tick").edgePoints(), st(v, "tick").edgePoints(); fmt.Sprint(a) != fmt.Sprint(b) {
		t.Errorf("tick moved from %v to %v", a, b)
	}
	// the edge between the icons runs between the fitted bounds
	pts := st(v, "e").edgePoints()
	if len(pts) != 2 || math.Abs(pts[0].y-104.75) > 1e-9 || math.Abs(pts[1].y-228.25) > 1e-9 {
		t.Errorf("edge %v, want from y 104.75 to 228.25", pts)
	}
	// groups keep their bounds
	if r := st(v, "region"); r.bounds() != (rect{300, 0, 200, 150}) || r.style["shape"] != "mxgraph.aws4.group" {
		t.Errorf("region %v %v", r.bounds(), r.style)
	}
	// a name without a counterpart stays, with a warning newShape does not repeat
	if s := st(v, "swf"); s.style["shape"] != "mxgraph.aws3.swf" || s.bounds() != (rect{300, 200, 76.5, 93}) {
		t.Errorf("swf %v %v", s.bounds(), s.style)
	}
	c.newShape(st(v, "swf"))
	if len(warnings) != 1 || !strings.Contains(warnings[0], `"mxgraph.aws3.swf" of an older AWS icon set has no current AWS counterpart`) {
		t.Errorf("warnings %q", warnings)
	}
}

// TestSubstituteAWSLegacyResolved checks that names the converter draws
// itself are left alone.
func TestSubstituteAWSLegacyResolved(t *testing.T) {
	shapeRegistry["mxgraph.aws3.ec2"] = shapeRegistry["rectangle"]
	defer delete(shapeRegistry, "mxgraph.aws3.ec2")
	c := &converter{opts: &Options{Warn: func(string) {}}, warned: map[string]bool{}}
	m := awsLegacyTestModel(t)
	if n := c.substituteAWSLegacy(m); n != 2 {
		t.Errorf("%d substituted, want 2", n)
	}
	if s := m.cells["ec2"].styleStr; !strings.Contains(s, "shape=mxgraph.aws3.ec2") || m.cells["ec2"].geo.h != 93 {
		t.Errorf("ec2 rewritten: %s", s)
	}
}

// awsLegacyJSShapes are the aws4 shapes draw.io implements in JavaScript
// (mxAWS4.js) that the targets use.
var awsLegacyJSShapes = map[string]bool{"mxgraph.aws4.resourceIcon": true, "mxgraph.aws4.productIcon": true,
	"mxgraph.aws4.group": true, "mxgraph.aws4.groupCenter": true, "mxgraph.aws4.group2": true}

// TestAWSLegacyTargets checks the generated table: every target draws with
// an embedded stencil or an aws4 JavaScript shape whose icon is an embedded
// stencil, and every name belongs to an older AWS icon set.
func TestAWSLegacyTargets(t *testing.T) {
	used := make([]bool, len(awsLegacyTargets))
	gens := map[string][2]int{}
	for name, i := range awsLegacyShapes {
		if !isAWSLegacyShape(name) {
			t.Errorf("%s: not an older AWS shape name", name)
		}
		g := strings.Split(name, ".")[1]
		c := gens[g]
		if i < 0 {
			c[1]++
		} else {
			c[0]++
			if int(i) >= len(awsLegacyTargets) {
				t.Fatalf("%s: target %d out of range", name, i)
			}
			used[i] = true
		}
		gens[g] = c
		if tg, ok := lookupAWSLegacy(name); !ok || (tg == nil) != (i < 0) {
			t.Errorf("%s: lookup %v %v", name, tg, ok)
		}
	}
	t.Logf("names mapped/unmapped per generation: %v", gens)
	for i, tg := range awsLegacyTargets {
		if !used[i] {
			t.Errorf("target %d is not used", i)
		}
		kv := map[string]string{}
		for _, e := range strings.Split(tg.style, ";") {
			k, v, ok := strings.Cut(e, "=")
			if !ok || v == "" {
				t.Errorf("target %d: entry %q", i, e)
			}
			kv[k] = v
		}
		if (tg.w > 0) != (tg.h > 0) {
			t.Errorf("target %d: size %gx%g", i, tg.w, tg.h)
		}
		shape, ok := kv["shape"]
		switch {
		case !ok:
			// a group drawn as a rectangle, or a draw.io style
		case awsLegacyJSShapes[shape]:
			icon := kv["resIcon"] + kv["grIcon"]
			if icon == "" || lookupStencil(icon) == nil {
				t.Errorf("target %d: icon %q does not resolve", i, icon)
			}
		case lookupStencil(shape) == nil:
			t.Errorf("target %d: stencil %q does not resolve", i, shape)
		}
		// groups (and frames) keep their bounds, icons are fitted
		group := !ok || shape == "mxgraph.aws4.group" || shape == "mxgraph.aws4.groupCenter" || shape == "mxgraph.aws4.group2"
		if group != (tg.w == 0) {
			t.Errorf("target %d: size %gx%g for %s", i, tg.w, tg.h, tg.style)
		}
	}
}

// TestAWSLegacyDiagrams converts the test diagrams of the older AWS icons:
// every old shape is known, and all but those without a counterpart are
// drawn with a current one.
func TestAWSLegacyDiagrams(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("testdata", "aws_legacy", "*.drawio"))
	if err != nil || len(files) == 0 {
		t.Fatal("no test diagrams", err)
	}
	summary := regexp.MustCompile(`^(\d+) shapes? of (an )?older AWS icon sets? (is|are) drawn`)
	for _, f := range files {
		name := filepath.Base(f)
		t.Run(name, func(t *testing.T) {
			var warnings []string
			opts := testOptions()
			opts.Warn = func(s string) { warnings = append(warnings, s) }
			if _, err := ConvertFile(f, opts); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			old, unmapped := 0, 0
			for _, m := range regexp.MustCompile(`shape=(mxgraph\.aws[23d]*\.[^;"]+)`).FindAllStringSubmatch(string(data), -1) {
				old++
				if tg, ok := lookupAWSLegacy(m[1]); !ok {
					t.Errorf("unknown old shape %s", m[1])
				} else if tg == nil {
					unmapped++
				}
			}
			got := -1
			unsupported := regexp.MustCompile(`shape "([^"]+)" is not supported`)
			for _, w := range warnings {
				// (the aws4 shapes of mxAWS4.js may not be there yet)
				if m := unsupported.FindStringSubmatch(w); strings.Contains(w, "internal error") || m != nil && isAWSLegacyShape(m[1]) {
					t.Error(w)
				}
				if m := summary.FindStringSubmatch(w); m != nil {
					fmt.Sscan(m[1], &got)
				}
			}
			if got != old-unmapped {
				t.Errorf("%d shapes substituted, want %d; warnings %q", got, old-unmapped, warnings)
			}
		})
	}
}

// TestAWSLegacyExport writes the test diagrams with the substituted
// styles and bounds into $DRAWIO_AWS_LEGACY_OUT, for comparing draw.io's
// rendering of the old icons with that of their substitutes:
//
//	DRAWIO_AWS_LEGACY_OUT=/tmp/aws go test -run TestAWSLegacyExport ./converter/drawio/
//	draw.io -x -f png -s 2 -o /tmp/aws/architecture.png /tmp/aws/architecture.drawio
func TestAWSLegacyExport(t *testing.T) {
	dir := os.Getenv("DRAWIO_AWS_LEGACY_OUT")
	if dir == "" {
		t.Skip("DRAWIO_AWS_LEGACY_OUT is not set")
	}
	files, _ := filepath.Glob(filepath.Join("testdata", "aws_legacy", "*.drawio"))
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		df, err := readFile(data)
		if err != nil {
			t.Fatal(err)
		}
		var b strings.Builder
		b.WriteString("<mxfile>")
		for _, p := range df.pages {
			m := parseModel(p.model)
			c := &converter{opts: &Options{Warn: func(string) {}}, warned: map[string]bool{}}
			c.substituteAWSLegacy(m)
			fmt.Fprintf(&b, `<diagram id="%s" name="%s">%s</diagram>`, html.EscapeString(p.id), html.EscapeString(p.name), modelXML(m))
		}
		b.WriteString("</mxfile>\n")
		if err := os.WriteFile(filepath.Join(dir, filepath.Base(f)), []byte(b.String()), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// modelXML writes a model back as an mxGraphModel element (cells, styles
// and geometry; the attributes of UserObject wrappers are dropped).
func modelXML(m *model) string {
	var b strings.Builder
	b.WriteString("<mxGraphModel><root>")
	num := func(v float64) string { return fmtNum(v) }
	pt := func(p *point, as string) {
		if p != nil {
			fmt.Fprintf(&b, `<mxPoint x="%s" y="%s" as="%s"/>`, num(p.x), num(p.y), as)
		}
	}
	for _, c := range m.ordered {
		fmt.Fprintf(&b, `<mxCell id="%s"`, html.EscapeString(c.id))
		for _, a := range [][2]string{{"value", c.value}, {"style", c.styleStr}, {"parent", c.parentID}, {"source", c.sourceID}, {"target", c.targetID}} {
			if a[1] != "" {
				fmt.Fprintf(&b, ` %s="%s"`, a[0], html.EscapeString(a[1]))
			}
		}
		for _, a := range []struct {
			name string
			on   bool
		}{{"vertex", c.vertex}, {"edge", c.edge}, {"collapsed", c.collapsed}} {
			if a.on {
				fmt.Fprintf(&b, ` %s="1"`, a.name)
			}
		}
		if !c.visible {
			b.WriteString(` visible="0"`)
		}
		b.WriteString(">")
		if g := c.geo; g != nil {
			fmt.Fprintf(&b, `<mxGeometry x="%s" y="%s" width="%s" height="%s"`, num(g.x), num(g.y), num(g.w), num(g.h))
			if g.relative {
				b.WriteString(` relative="1"`)
			}
			b.WriteString(` as="geometry">`)
			pt(g.sourcePoint, "sourcePoint")
			pt(g.targetPoint, "targetPoint")
			pt(g.offset, "offset")
			if len(g.points) > 0 {
				b.WriteString(`<Array as="points">`)
				for i := range g.points {
					pt(&g.points[i], "")
				}
				b.WriteString("</Array>")
			}
			b.WriteString("</mxGeometry>")
		}
		b.WriteString("</mxCell>")
	}
	b.WriteString("</root></mxGraphModel>")
	return strings.ReplaceAll(b.String(), ` as=""`, "")
}
