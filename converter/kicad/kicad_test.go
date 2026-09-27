package kicad

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/internal/fontdb"
)

const testFonts = "../pptx/testdata/fonts"

// testConv returns a conversion that fails the test on warnings.
func testConv(t *testing.T) *conv {
	t.Helper()
	c := &conv{opts: &Options{Warn: func(m string) { t.Error("warning:", m) }}, doc: bdf.NewDocument(), warned: map[string]bool{}, sch: defaultSchSettings}
	set := fontset.New(fontdb.New(nil, []string{testFonts}, false), nil)
	c.fonts = &cad.Fonts{Set: set}
	c.plotter = &cad.Plotter{Fonts: c.fonts}
	c.cvs = canvas.NewBuilder(c.doc, set)
	return c
}

// schematic wraps items in a KiCad 8 schematic whose root sheet has the
// uuid rootUUID.
func schematic(items ...string) string {
	return `(kicad_sch (version 20231120) (generator "eeschema") (uuid "` + rootUUID + `") (paper "A4")
` + strings.Join(items, "\n") + "\n)"
}

const rootUUID = "00000000-0000-0000-0000-000000000001"

// drawSheet draws the first sheet of a schematic and describes its items.
func drawSheet(t *testing.T, src string, settings ...func(*schSettings)) []string {
	t.Helper()
	c := testConv(t)
	for _, f := range settings {
		f(&c.sch)
	}
	c.iu, c.defaultPen = 10000, schLineWidth
	sheets, err := c.loadSchematic("test.kicad_sch", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	p := &schPage{c: c, s: sheets[0], all: sheets, pageOf: map[*sheetInst]int{}, d: &cad.Drawing{}, count: len(sheets), v6refs: map[string]*node{}}
	p.w, p.h = paperSize(sheets[0].file.n.child("paper"))
	p.draw()
	return p.d.Describe(0.0001)
}

var strokeTextRE = regexp.MustCompile(`^\s*stroke-text "((?:[^"\\]|\\.)*)" (\S+) (\S+) (\S+) (\S+) color=\S+ width=(\S+)`)

// strokeExtents returns the extents of the centre lines of the strokes of
// each stroke text (the bounds less half the pen), by text.
func strokeExtents(t *testing.T, desc []string) map[string][4]float64 {
	t.Helper()
	out := map[string][4]float64{}
	for _, line := range desc {
		m := strokeTextRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		s, err := strconv.Unquote(`"` + m[1] + `"`)
		if err != nil {
			t.Fatal(err)
		}
		var v [5]float64
		for i := range v {
			v[i], _ = strconv.ParseFloat(m[i+2], 64)
		}
		h := v[4] / 2
		out[s] = [4]float64{v[0] + h, v[1] + h, v[2] - h, v[3] - h}
	}
	return out
}

// TestTextPlacement checks the strokes of texts against where KiCad 10
// draws them (its SVG export of a demo project, to 0.1 µm).
func TestTextPlacement(t *testing.T) {
	src := schematic(`(lib_symbols
  (symbol "t:ICL7660" (pin_names (offset 1.016)) (in_bom yes) (on_board yes)
    (property "Reference" "U" (at 5.08 10.16 0) (effects (font (size 1.778 1.778))))
    (property "Value" "ICL7660" (at 1.27 -11.43 0) (effects (font (size 1.778 1.778))))
    (symbol "ICL7660_0_1" (rectangle (start -13.97 13.97) (end 13.97 -13.97) (stroke (width 0) (type default)) (fill (type background))))
    (symbol "ICL7660_1_1"
      (pin input line (at -21.59 6.35 0) (length 7.62) (name "CAP+" (effects (font (size 1.524 1.524)))) (number "2" (effects (font (size 1.524 1.524)))))
      (pin power_out line (at 21.59 3.81 180) (length 7.62) (name "VOUT" (effects (font (size 1.524 1.524)))) (number "5" (effects (font (size 1.524 1.524))))))))`,
		`(label "12Vext" (at 54.61 63.5 0) (effects (font (size 1.524 1.524)) (justify left bottom)) (uuid "00000000-0000-0000-0000-000000000002"))`,
		`(symbol (lib_id "t:ICL7660") (at 201.93 73.66 0) (unit 1) (in_bom yes) (on_board yes) (dnp no) (uuid "00000000-0000-0000-0000-000000000003")
  (property "Reference" "U102" (at 187.96 63.5 0) (effects (font (size 1.778 1.778)) (justify left)))
  (property "Value" "ICL7660" (at 215.9 85.09 0) (effects (font (size 1.778 1.778)) (justify right)))
  (instances (project "t" (path "/`+rootUUID+`" (reference "U102") (unit 1)))))`)
	// the demo project's text offset ratio
	got := strokeExtents(t, drawSheet(t, src, func(s *schSettings) { s.textOffset = 0.3 }))
	for s, want := range map[string][4]float64{
		"12Vext":  {55.0004, 61.1338, 62.185, 62.6578},
		"U102":    {188.539, 62.5365, 194.7196, 64.3145},
		"ICL7660": {205.3304, 84.1265, 215.4056, 85.9045},
		"CAP+":    {189.3664, 66.483, 194.9545, 68.007},
		"VOUT":    {209.2683, 69.023, 214.6387, 70.547},
	} {
		g, ok := got[s]
		if !ok {
			t.Errorf("%s: not drawn (%v)", s, got)
			continue
		}
		for i := range want {
			if math.Abs(g[i]-want[i]) > 0.002 {
				t.Errorf("%s: strokes span %v, KiCad draws %v", s, g, want)
				break
			}
		}
	}
}

// TestLabelPlacement checks labels and notes at the default settings
// against KiCad 10's SVG export of the same sheet.
func TestLabelPlacement(t *testing.T) {
	src := schematic(
		`(label "L1a" (at 50 50 0) (effects (font (size 1.524 1.524)) (justify left bottom)) (uuid "00000000-0000-0000-0000-000000000002"))`,
		`(label "L2b" (at 50 70 90) (effects (font (size 2.54 2.54)) (justify left bottom)) (uuid "00000000-0000-0000-0000-000000000003"))`,
		`(label "L3c" (at 50 90 180) (effects (font (size 1.27 1.27)) (justify right bottom)) (uuid "00000000-0000-0000-0000-000000000004"))`,
		`(label "L4d" (at 50 110 0) (effects (font (size 2.54 2.54) (thickness 0.5)) (justify left bottom)) (uuid "00000000-0000-0000-0000-000000000005"))`,
		`(global_label "G1" (shape input) (at 100 50 0) (effects (font (size 1.27 1.27)) (justify left)) (uuid "00000000-0000-0000-0000-000000000008"))`,
		`(global_label "G2" (shape bidirectional) (at 100 70 180) (effects (font (size 1.27 1.27)) (justify right)) (uuid "00000000-0000-0000-0000-000000000009"))`,
		`(hierarchical_label "H1" (shape output) (at 150 50 0) (effects (font (size 1.27 1.27)) (justify left)) (uuid "00000000-0000-0000-0000-00000000000a"))`,
		`(text "T1\nline2" (at 150 70 0) (effects (font (size 1.524 1.524)) (justify left bottom)) (uuid "00000000-0000-0000-0000-000000000006"))`,
		`(text "T2" (at 150 90 90) (effects (font (size 1.27 1.27) italic) (justify left)) (uuid "00000000-0000-0000-0000-00000000000b"))`)
	desc := drawSheet(t, src)
	got := strokeExtents(t, desc)
	for s, want := range map[string][4]float64{
		"L1a": {50.463, 47.8624, 53.8013, 49.3864},
		"L2b": {46.4426, 63.6103, 48.9826, 69.295},
		"L3c": {46.8759, 88.2172, 49.6578, 89.4872},
		"L4d": {50.9336, 106.1212, 56.4974, 108.7821},
		"G1":  {101.7709, 49.4003, 103.7666, 50.6703},
		"G2":  {96.2334, 69.4003, 98.2291, 70.6703},
		"H1":  {151.863, 49.3095, 153.8588, 50.5795},
		"T1":  {150.2453, 65.5779, 152.4224, 67.1019},
		"T2":  {149.3095, 87.5936, 150.5795, 89.3776},
	} {
		g, ok := got[s]
		if !ok {
			t.Errorf("%s: not drawn", s)
			continue
		}
		for i := range want {
			if math.Abs(g[i]-want[i]) > 0.002 {
				t.Errorf("%s: strokes span %v, KiCad draws %v", s, g, want)
				break
			}
		}
	}
	// the outline of G1
	if !hasBox(desc, "stroke", [4]float64{100, 48.7296, 104.9249, 51.2704}) {
		t.Errorf("no global label outline in\n%s", strings.Join(desc, "\n"))
	}
}

// hasBox reports whether an item of the kind spans the box (to 1 µm).
func hasBox(desc []string, kind string, box [4]float64) bool {
	for _, line := range desc {
		f := strings.Fields(line)
		if len(f) < 5 || f[0] != kind {
			continue
		}
		ok := true
		for i := range box {
			v, err := strconv.ParseFloat(f[i+1], 64)
			if err != nil || math.Abs(v-box[i]) > 0.001 {
				ok = false
			}
		}
		if ok {
			return true
		}
	}
	return false
}

func TestSexpr(t *testing.T) {
	n, err := parse([]byte(" \ufeff(a b \"c \\\"d\\\"\\n\" (e 1 2.5 -3) (f) (hide yes) visible (g no) (h) (\"x y\" z))"))
	if err == nil {
		t.Fatal("a BOM before the list was accepted")
	}
	n, err = parse([]byte("\n\t(a b \"c \\\"d\\\"\\n\" (e 1 2.5 -3) (f) (hide yes) visible (g no) (h) (\"x y\" z) \"hide\")"))
	if err != nil {
		t.Fatal(err)
	}
	if n.name != "a" || n.arg(0) != "b" || n.arg(1) != "c \"d\"\n" || n.arg(2) != "visible" || n.arg(3) != "hide" || n.arg(4) != "" {
		t.Errorf("arguments: %q %q %q %q", n.name, n.arg(0), n.arg(1), n.arg(2))
	}
	if e := n.child("e"); e.num(0) != 1 || e.num(1) != 2.5 || e.int(2) != -3 || e.num(3) != 0 {
		t.Errorf("numbers: %v", e.items)
	}
	if n.child("x y").arg(0) != "z" {
		t.Error("a quoted name was not kept")
	}
	for name, want := range map[string]bool{"hide": true, "visible": true, "g": false, "h": true, "f": true, "missing": false} {
		if got := n.flag(name); got != want {
			t.Errorf("flag %s = %v", name, got)
		}
	}
	if !n.flagOr("missing", true) || n.flagOr("g", true) {
		t.Error("flagOr")
	}
	if n.hasAtom("hide") {
		// the quoted "hide" is a string, not the keyword
		t.Error("a quoted string was taken for a keyword")
	}
	var nilNode *node
	if nilNode.child("a") != nil || nilNode.arg(0) != "" || nilNode.num(0) != 0 || nilNode.flag("a") || len(nilNode.children("a")) != 0 {
		t.Error("a missing node is not empty")
	}
	for _, bad := range []string{"", "a", "(a", "(a \"b)", "(a \"b\\", "(a (b)", strings.Repeat("(", 300) + strings.Repeat(")", 300)} {
		if _, err := parse([]byte(bad)); err == nil {
			t.Errorf("%q: no error", bad)
		}
	}
	if n, _ := parse([]byte("(a 1e400 nan inf 3e9)")); n.num(0) != 0 || n.num(1) != 0 || n.num(2) != 0 || n.int(3) != 0 {
		t.Error("numbers out of range were not zero")
	}
}

func TestMarkup(t *testing.T) {
	for s, want := range map[string]string{
		"V_{CC}":         "VCC",
		"x^{2}+~{RESET}": "x2+RESET",
		"~{A}_{b^{c}}":   "Abc",
		"unclosed_{a":    "unclosed_{a",
		"a_{b^{c}":       "a_{bc",
		"a~b":            "a~b",
	} {
		if got := parseMarkup(s).plain(); got != want {
			t.Errorf("%q: %q, want %q", s, got, want)
		}
	}
}

func TestExpand(t *testing.T) {
	c := testConv(t)
	c.project = "demo"
	c.textVars = map[string]string{"A": "${B}", "B": "b", "LOOP": "${LOOP}x"}
	lookup := func(name string) (string, bool) {
		if name == "SHEETNAME" {
			return "Root", true
		}
		return "", false
	}
	for s, want := range map[string]string{
		"${A}/${SHEETNAME}/${PROJECTNAME}": "b/Root/demo",
		"${UNKNOWN}":                       "${UNKNOWN}",
		"$x ${":                            "$x ${",
	} {
		if got := c.expand(s, lookup); got != want {
			t.Errorf("%q: %q, want %q", s, got, want)
		}
	}
	if got := c.expand("${LOOP}", lookup); len(got) > 100 {
		t.Errorf("a variable that refers to itself expanded to %d bytes", len(got))
	}
}

func TestSymbolTransform(t *testing.T) {
	// a point right of and above the anchor: (1, 2) in the file's y-up
	// library, (1, -2) in KiCad's y-down one
	for _, tc := range []struct {
		angle  float64
		mirror string
		x, y   float64
	}{
		{0, "", 1, -2},
		{90, "", -2, -1},
		{180, "", -1, 2},
		{270, "", 2, 1},
		{0, "x", 1, 2},
		{0, "y", -1, -2},
		{90, "y", 2, -1},
	} {
		tr := symbolTransform(tc.angle, tc.mirror)
		x, y := tr.apply(1, -2)
		if x != tc.x || y != tc.y {
			t.Errorf("angle %g mirror %q: (1, -2) → (%g, %g), want (%g, %g)", tc.angle, tc.mirror, x, y, tc.x, tc.y)
		}
		if ix, iy := tr.invApply(x, y); ix != 1 || iy != -2 {
			t.Errorf("angle %g mirror %q: the inverse gives (%g, %g)", tc.angle, tc.mirror, ix, iy)
		}
	}
}

func TestLabelSpin(t *testing.T) {
	for _, tc := range []struct {
		angle float64
		h     hAlign
		want  spin
	}{
		{0, alignLeft, spinRight},
		{0, alignRight, spinLeft},
		{180, alignRight, spinLeft},
		{90, alignLeft, spinUp},
		{270, alignRight, spinBottom},
	} {
		if got := labelSpin(keepUpright(tc.angle), tc.h); got != tc.want {
			t.Errorf("angle %g %v: spin %v, want %v", tc.angle, tc.h, got, tc.want)
		}
	}
}
