package drawio

import (
	"fmt"
	"html"
	"math"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// The shapes of mxAWS4.js.
var aws4ShapeNames = []string{"mxgraph.aws4.productIcon", "mxgraph.aws4.resourceIcon", "mxgraph.aws4.group",
	"mxgraph.aws4.groupCenter", "mxgraph.aws4.group2"}

// aws4Styles are style variants with the keys the AWS shapes read,
// including names that are no stencils and values that are not numbers.
var aws4Styles = []string{
	"resIcon=mxgraph.aws4.lambda;prIcon=mxgraph.aws4.athena;grIcon=mxgraph.aws4.group_aws_cloud_alt;strokeColor=#ffffff;fillColor=#ED7100;",
	"resIcon=mxgraph.aws4.piop;prIcon=;grIcon=mxgraph.aws4.resourceIcon;",
	"resIcon=mxgraph.flowchart.decision;prIcon=mxgraph.aws4;grIcon=stencil(abc);",
	"grIcon=mxgraph.aws4.group_region;grIconSize=abc;grStroke=0;",
	"grIcon=mxgraph.aws4.group_region;grIconSize=;grStroke=2;",
	"grIcon=mxgraph.aws4.group_region;grIconSize=1e308;strokeOpacity=abc;",
	"grIcon=mxgraph.aws4.group_region;grIconSize=-25;strokeOpacity=-50;",
	"grIcon=mxgraph.aws4.group_region;grIconSize=Infinity;",
	"grIcon=mxgraph.aws4.group_region;grIconSize=0x19;pointerEvents=1;",
	"prIcon=mxgraph.aws4.athena;opacity=abc;gradientColor=#F78E04;gradientDirection=radial;",
	"prIcon=mxgraph.aws4.athena;opacity=0.5;fillColor=none;strokeColor=none;",
	"resIcon=mxgraph.aws4.users;grIcon=mxgraph.aws4.group_region;direction=south;flipH=1;rotation=45;shadow=1;",
	"resIcon=mxgraph.aws4.users;grIcon=mxgraph.aws4.group_region;direction=north;flipV=1;gradientColor=#ff0000;",
}

// TestAws4ShapesPaint paints the AWS shapes with their style keys and
// degenerate sizes (TestShapesPaint paints them with the common style
// variants); an internal error (a recovered panic) fails the test.
func TestAws4ShapesPaint(t *testing.T) {
	var b strings.Builder
	b.WriteString("<mxfile>")
	id := 0
	for _, name := range aws4ShapeNames {
		fmt.Fprintf(&b, `<diagram id="%s" name="%s"><mxGraphModel><root><mxCell id="0"/><mxCell id="1" parent="0"/>`, name, name)
		for _, st := range aws4Styles {
			for _, sz := range shapeSizes {
				id++
				style := html.EscapeString("shape=" + name + ";" + st)
				fmt.Fprintf(&b, `<mxCell id="v%d" value="label" style="%s" vertex="1" parent="1"><mxGeometry x="10" y="10" width="%g" height="%g" as="geometry"/></mxCell>`, id, style, sz[0], sz[1])
			}
		}
		b.WriteString(`</root></mxGraphModel></diagram>`)
	}
	b.WriteString("</mxfile>")
	res, err := Convert([]byte(b.String()), &Options{NoSystemFonts: true, NoTextIndex: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.Pages != len(aws4ShapeNames) {
		t.Errorf("%d pages, want %d", res.Pages, len(aws4ShapeNames))
	}
	for _, w := range res.Warnings {
		if strings.Contains(w, "internal error") || strings.Contains(w, "not supported") {
			t.Error(w)
		}
	}
}

// TestAws4ShapesRegistered checks that the shapes of mxAWS4.js are
// registered and are no stencils, and that the icons the palette names
// (Sidebar-AWS4.js) resolve in the aws4 library.
func TestAws4ShapesRegistered(t *testing.T) {
	for _, name := range aws4ShapeNames {
		if shapeRegistry[name] == nil {
			t.Errorf("shape %q is not registered", name)
		}
		if lookupStencil(name) != nil {
			t.Errorf("%q resolves to a stencil", name)
		}
	}
	var got []string
	for name := range shapeRegistry {
		if strings.HasPrefix(name, "mxgraph.aws4.") {
			got = append(got, name)
		}
	}
	if len(got) != len(aws4ShapeNames) {
		slices.Sort(got)
		t.Errorf("registered aws4 shapes %v", got)
	}
	for _, name := range []string{"ec2", "lambda", "s3", "athena", "marketplace", "general", "all_products",
		"group_aws_cloud_alt", "group_aws_cloud", "group_region", "group_vpc2", "group_security_group",
		"group_auto_scaling_group", "group_on_premise", "group_corporate_data_center", "group_elastic_beanstalk",
		"group_ec2_instance_contents", "group_spot_fleet", "group_aws_step_functions_workflow", "group_account",
		"group_iot_greengrass_deployment", "group_iot_greengrass", "group_availability_zone", "group_subnet",
		"group_vpc", "group_elastic_load_balancing"} {
		if lookupStencil("mxgraph.aws4."+name) == nil {
			t.Errorf("stencil mxgraph.aws4.%s is not found", name)
		}
	}
}

// paintAws4 paints a shape of mxAWS4.js in b and returns the decoded object.
func paintAws4(t *testing.T, name, st string, b rect) *bdf.ObjectPart {
	t.Helper()
	s := testShape(t, name, st)
	s.conv.opts = &Options{}
	s.bounds = b
	cv := s.conv.newCanvas()
	s.paint(newC2D(cv))
	part, err := bdf.DecodeObject(cv.obj.Encode())
	if err != nil {
		t.Fatal(err)
	}
	return part
}

// fillColors returns the fill color in effect for each filled path, in
// the order of painting (0 for a paint).
func fillColors(t *testing.T, part *bdf.ObjectPart) []bdf.Color {
	t.Helper()
	ins, err := part.Instructions()
	if err != nil {
		t.Fatal(err)
	}
	var cur bdf.Color
	var out []bdf.Color
	for _, in := range ins {
		switch in.Op {
		case bdf.OpFillColor:
			cur = bdf.Color(in.Args[0].(uint64))
		case bdf.OpFillPaint:
			cur = 0
		case bdf.OpFillPath:
			out = append(out, cur)
		}
	}
	return out
}

// TestAws4Geometry checks the boxes the AWS shapes paint against values
// worked out from the JavaScript. The a1 instance icon (44 by 44) reaches
// the four sides of its stencil, so its box is the stencil's.
func TestAws4Geometry(t *testing.T) {
	const icon = "mxgraph.aws4.a1_instance"
	for _, tc := range []struct {
		name, style string
		b           rect
		want        []rect // the boxes of the paths, in order
	}{
		// the middle 80%
		{"mxgraph.aws4.resourceIcon", "resIcon=" + icon, rect{10, 20, 78, 78},
			[]rect{{10, 20, 78, 78}, {17.8, 27.8, 62.4, 62.4}}},
		// a fixed aspect: centered in the 80% (80 by 48)
		{"mxgraph.aws4.resourceIcon", "resIcon=" + icon, rect{0, 0, 100, 60},
			[]rect{{0, 0, 100, 60}, {26, 6, 48, 48}}},
		{"mxgraph.aws4.resourceIcon", "resIcon=mxgraph.aws4.piop", rect{0, 0, 60, 60}, []rect{{0, 0, 60, 60}}},
		// the frame, the square inset by 1 as high as wide, the icon at
		// 1 + 15% with 70% - 2
		{"mxgraph.aws4.productIcon", "prIcon=" + icon + ";strokeColor=#ffffff;", rect{0, 0, 80, 110},
			[]rect{{0, 0, 80, 110}, {1, 1, 78, 78}, {13, 13, 54, 54}}},
		{"mxgraph.aws4.productIcon", "prIcon=" + icon + ";strokeColor=#ffffff;", rect{0, 0, 100, 80},
			[]rect{{0, 0, 100, 80}, {1, 1, 98, 98}, {16, 16, 68, 68}}},
		// grIconSize at the top left, or centered at the top
		{"mxgraph.aws4.group", "grIcon=" + icon, rect{5, 5, 130, 130},
			[]rect{{5, 5, 130, 130}, {5, 5, 25, 25}}},
		{"mxgraph.aws4.group", "grIcon=" + icon + ";grIconSize=40;", rect{0, 0, 130, 130},
			[]rect{{0, 0, 130, 130}, {0, 0, 40, 40}}},
		{"mxgraph.aws4.group", "grIcon=" + icon + ";grIconSize=abc;", rect{0, 0, 130, 130},
			[]rect{{0, 0, 130, 130}}},
		{"mxgraph.aws4.groupCenter", "grIcon=" + icon, rect{0, 0, 130, 130},
			[]rect{{0, 0, 130, 130}, {52.5, 0, 25, 25}}},
		// the 25 by 25 square with the icon in its middle 80%
		{"mxgraph.aws4.group2", "grIcon=" + icon, rect{0, 0, 130, 130},
			[]rect{{0, 0, 130, 130}, {0, 0, 25, 25}, {2.5, 2.5, 20, 20}}},
		{"mxgraph.aws4.group2", "grIcon=" + icon + ";grIconSize=40;", rect{0, 0, 130, 130},
			[]rect{{0, 0, 130, 130}, {0, 0, 25, 25}, {4, 4, 32, 32}}},
	} {
		part := paintAws4(t, tc.name, tc.style, tc.b)
		if len(part.Paths) != len(tc.want) {
			t.Errorf("%s %q: %d paths, want %d", tc.name, tc.style, len(part.Paths), len(tc.want))
			continue
		}
		for i, want := range tc.want {
			box, _ := pathBox(paintedPath(t, part, i))
			if math.Abs(box.x-want.x) > 1e-3 || math.Abs(box.y-want.y) > 1e-3 || math.Abs(box.w-want.w) > 1e-3 || math.Abs(box.h-want.h) > 1e-3 {
				t.Errorf("%s %q: path %d box %v, want %v", tc.name, tc.style, i, box, want)
			}
		}
	}
}

// TestAws4Colors checks the colors the shapes fill with and whether the
// groups stroke their frames.
func TestAws4Colors(t *testing.T) {
	white, black := bdf.RGB(255, 255, 255), bdf.RGB(0, 0, 0)
	orange, navy := bdf.RGB(0xED, 0x71, 0x00), bdf.RGB(0x23, 0x2F, 0x3E)
	b := rect{0, 0, 80, 80}
	for _, tc := range []struct {
		name, style string
		fills       []bdf.Color
		strokes     int
	}{
		// the category color and the icon in the stroke color
		{"mxgraph.aws4.resourceIcon", "resIcon=mxgraph.aws4.lambda;fillColor=#ED7100;strokeColor=#ffffff;", []bdf.Color{orange, white}, 0},
		// strokeColor=none removes the key: the icon is black
		{"mxgraph.aws4.resourceIcon", "resIcon=mxgraph.aws4.lambda;fillColor=#ED7100;strokeColor=none;", []bdf.Color{orange, black}, 0},
		// a gradient is a paint
		{"mxgraph.aws4.resourceIcon", "resIcon=mxgraph.aws4.lambda;fillColor=#ED7100;gradientColor=#ffffff;", []bdf.Color{0, black}, 0},
		// the frame in the stroke color, the square in the fill color
		// (white for none), the icon in the stroke color
		{"mxgraph.aws4.productIcon", "prIcon=mxgraph.aws4.athena;fillColor=#232F3E;strokeColor=#ffffff;", []bdf.Color{white, navy, white}, 0},
		{"mxgraph.aws4.productIcon", "prIcon=mxgraph.aws4.athena;fillColor=none;strokeColor=#232F3E;", []bdf.Color{navy, white, navy}, 0},
		{"mxgraph.aws4.productIcon", "prIcon=mxgraph.aws4.athena;fillColor=#D05C17;gradientColor=#F78E04;strokeColor=#ffffff;", []bdf.Color{white, 0, white}, 0},
		// strokeColor=none: nothing but the square
		{"mxgraph.aws4.productIcon", "prIcon=mxgraph.aws4.athena;fillColor=#232F3E;strokeColor=none;", []bdf.Color{navy}, 0},
		// grStroke=1 (the default) strokes the frame, other values do not
		{"mxgraph.aws4.group", "grIcon=mxgraph.aws4.group_region;fillColor=none;strokeColor=#ED7100;", []bdf.Color{orange}, 1},
		{"mxgraph.aws4.group", "grIcon=mxgraph.aws4.group_region;fillColor=#ffffff;strokeColor=#ED7100;grStroke=0;", []bdf.Color{white, orange}, 0},
		{"mxgraph.aws4.groupCenter", "grIcon=mxgraph.aws4.group_region;fillColor=none;strokeColor=#ED7100;grStroke=2;", []bdf.Color{orange}, 0},
		// group2 always strokes; the icon is white on the stroke color
		{"mxgraph.aws4.group2", "grIcon=mxgraph.aws4.group_region;fillColor=none;strokeColor=#ED7100;grStroke=0;", []bdf.Color{orange, white}, 1},
	} {
		part := paintAws4(t, tc.name, tc.style, b)
		fills := fillColors(t, part)
		ins, _ := part.Instructions()
		strokes := 0
		for _, in := range ins {
			if in.Op == bdf.OpStrokePath {
				strokes++
			}
		}
		// the icons may have more than one path: compare their first
		if len(fills) < len(tc.fills) || !slices.Equal(fills[:len(tc.fills)], tc.fills) || strokes != tc.strokes {
			t.Errorf("%s %q: fills %08x and %d strokes, want %08x and %d", tc.name, tc.style, fills, strokes, tc.fills, tc.strokes)
		}
	}
	// the group icons are filled at the stroke opacity
	part := paintAws4(t, "mxgraph.aws4.group", "grIcon=mxgraph.aws4.group_region;fillColor=none;strokeColor=#000000;strokeOpacity=40;", b)
	if fills := fillColors(t, part); len(fills) != 1 || fills[0] != bdf.RGBA(0, 0, 0, 102) {
		t.Errorf("group icon at strokeOpacity=40: fills %08x", fills)
	}
}

// TestAws4Testdata converts the test diagrams of testdata/shapes_aws4
// (compared with draw.io's exports when the shapes were ported) without
// warnings.
func TestAws4Testdata(t *testing.T) {
	files, _ := filepath.Glob(filepath.Join("testdata", "shapes_aws4", "*.drawio"))
	if len(files) == 0 {
		t.Fatal("no test diagrams")
	}
	for _, f := range files {
		res, err := ConvertFile(f, testOptions())
		if err != nil {
			t.Errorf("%s: %v", f, err)
			continue
		}
		for _, w := range res.Warnings {
			t.Errorf("%s: %s", f, w)
		}
	}
}
