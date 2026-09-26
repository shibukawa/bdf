package sxf

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
)

// The test drawings are made by test/sxf/gen.py: one drawing as SFC, as
// P21 written the way the SCADEC library writes it, and as that P21 file
// zipped. Text is laid out with the test fonts of the PowerPoint converter.

func testOptions() *Options {
	return &Options{FontDirs: []string{"../pptx/testdata/fonts"}, NoSystemFonts: true}
}

func convert(t *testing.T, name string, opts *Options) (*Result, *bdf.Reader, []byte) {
	t.Helper()
	if opts == nil {
		opts = testOptions()
	}
	res, err := ConvertFile(filepath.Join("testdata", name), opts)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := res.Doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	return res, r, buf.Bytes()
}

func plainText(t *testing.T, r *bdf.Reader) string {
	t.Helper()
	h, err := bdf.ParseHash(r.Manifest.Views[0].TextIndex)
	if err != nil {
		t.Fatal(err)
	}
	b, err := r.Part(h)
	if err != nil {
		t.Fatal(err)
	}
	idx, err := bdf.DecodeTextIndex(b)
	if err != nil {
		t.Fatal(err)
	}
	return bdf.PlainText(idx)
}

// firstFill returns the first fill color of the page: its background.
func firstFill(t *testing.T, r *bdf.Reader) bdf.Color {
	t.Helper()
	o, err := r.Object(r.Manifest.Views[0].Pages[0].Layers[0].Obj)
	if err != nil {
		t.Fatal(err)
	}
	var first bdf.Color
	found := false
	o.Walk(func(in bdf.Instr) {
		if !found && in.Op == bdf.OpFillColor {
			first, found = bdf.Color(in.Args[0].(uint64)), true
		}
	})
	return first
}

// draw draws a file into a drawing and returns its items and the warnings.
func draw(t *testing.T, text string, light bool) ([]string, []string) {
	t.Helper()
	c := &converter{opts: &Options{}, warned: map[string]bool{}}
	set := fontset.New(fontdb.New(nil, []string{"../pptx/testdata/fonts"}, false), func(string) {})
	dr := &cad.Drawing{}
	rend := &renderer{fonts: &cad.Fonts{Set: set}, out: dr, light: light, bg: bdf.RGB(0, 0, 0), warn: c.warnOnce}
	if _, _, _, _, err := c.draw(text, rend); err != nil {
		t.Fatal(err)
	}
	return dr.Describe(0.01), c.warnings
}

func readText(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return decode(b)
}

func TestShapes(t *testing.T) {
	res, r, _ := convert(t, "shapes.p21", nil)
	if len(res.Warnings) > 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	m := r.Manifest
	if len(m.Views) != 1 || len(m.Views[0].Pages) != 1 || m.Meta.Source != "p21" || res.Format != "p21" {
		t.Fatalf("manifest: %+v", m)
	}
	// A3, the empty group that carries the background left out
	if res.W != 420 || res.H != 297 {
		t.Errorf("page %v × %v mm", res.W, res.H)
	}
	if got := m.Meta.DC.Title.First(); got != "SXF テスト図" {
		t.Errorf("title %q", got)
	}
	if got := m.Meta.DC.Language.First(); got != "ja" {
		t.Errorf("language %q", got)
	}
	if got := firstFill(t, r); got != bdf.RGB(16, 24, 48) {
		t.Errorf("background %08x", got)
	}
	text := plainText(t, r)
	for _, want := range []string{"SXF テスト図", "文字1", "文字5", "文字9", "カイテン 30°", "スラント", "字間 アリ", "縦書き",
		"60", "60°", "R12", "20", "ラベル", "1", "A"} {
		if !strings.Contains(text, want) {
			t.Errorf("text lacks %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "ヒヒョウジ") {
		t.Errorf("text of a hidden layer is drawn:\n%s", text)
	}
	// the part is placed twice and once more in a group
	if n := strings.Count(text, "A"); n != 3 {
		t.Errorf("%d texts of the part", n)
	}
}

// TestSFCMatchesP21 draws the SFC and P21 files of the same drawing: SFC
// numbers its styles and places text at its anchors, P21 refers to styles
// and SCADEC moves text off its anchors, turns leader arrows round and
// flags angular dimension arcs clockwise.
func TestSFCMatchesP21(t *testing.T) {
	sfc, w1 := draw(t, readText(t, "shapes.sfc"), false)
	p21, w2 := draw(t, readText(t, "shapes.p21"), false)
	if len(w1) > 0 || len(w2) > 0 {
		t.Errorf("warnings: %q, %q", w1, w2)
	}
	if len(sfc) < 100 {
		t.Fatalf("%d items", len(sfc))
	}
	if len(sfc) != len(p21) {
		t.Errorf("%d items from the SFC file, %d from the P21 file", len(sfc), len(p21))
	}
	shown := 0
	for i := range min(len(sfc), len(p21)) {
		if sfc[i] != p21[i] && shown < 10 {
			t.Errorf("item %d:\nsfc %s\np21 %s", i, sfc[i], p21[i])
			shown++
		}
	}
}

func TestP2Z(t *testing.T) {
	_, _, a := convert(t, "shapes.p21", nil)
	_, _, b := convert(t, "shapes.p2z", nil)
	if !bytes.Equal(a, b) {
		t.Error("the zipped P21 file converts differently")
	}
}

func TestLight(t *testing.T) {
	opts := testOptions()
	opts.Light = true
	_, r, _ := convert(t, "shapes.sfc", opts)
	if got := firstFill(t, r); got != bdf.RGB(255, 255, 255) {
		t.Errorf("background %08x", got)
	}
	items, _ := draw(t, readText(t, "shapes.sfc"), true)
	white := fmt.Sprintf("color=%08x", uint32(bdf.RGB(255, 255, 255)))
	black := fmt.Sprintf("color=%08x", uint32(bdf.RGB(0, 0, 0)))
	blacks := 0
	for _, it := range items {
		if strings.Contains(it, white) {
			t.Errorf("white on white paper: %s", it)
			break
		}
		if strings.Contains(it, black) {
			blacks++
		}
	}
	if blacks < 20 {
		t.Errorf("%d black items", blacks)
	}
}

const stepHeader = `ISO-10303-21;
HEADER;
FILE_DESCRIPTION((''),'2;1');
FILE_NAME('t.p21','',(''),(''),'%s','','');
FILE_SCHEMA(('ASSOCIATIVE_DRAUGHTING'));
ENDSEC;
DATA;
`

const stepText = `#1=CARTESIAN_POINT(' ',(10.0,20.0));
#2=DIRECTION(' ',(1.0,0.0));
#3=AXIS2_PLACEMENT_2D(' ',#1,#2);
#4=PLANAR_EXTENT(' ',10.0,5.0);
#5=TEXT_LITERAL_WITH_EXTENT('$$SXF_topline left','AB',#3,'baseline left',.RIGHT.,#6,#4);
#6=EXTERNALLY_DEFINED_TEXT_FONT(IDENTIFIER('x'),#7);
#7=EXTERNAL_SOURCE(IDENTIFIER('scadec'));
#8=DRAUGHTING_PRE_DEFINED_COLOUR('white');
#9=TEXT_STYLE_FOR_DEFINED_FONT(#8);
#10=TEXT_STYLE(' ',#9);
#11=PRESENTATION_STYLE_ASSIGNMENT((#10));
#12=(ANNOTATION_OCCURRENCE()ANNOTATION_TEXT_OCCURRENCE()DRAUGHTING_ANNOTATION_OCCURRENCE()
GEOMETRIC_REPRESENTATION_ITEM()REPRESENTATION_ITEM(' ')STYLED_ITEM((#11),#5));
#13=DRAWING_SHEET_REVISION('s',(#12),$,'01');
ENDSEC;
END-ISO-10303-21;
`

// TestTextPoint places a text at its anchor, as the SXF specification
// says, unless SCADEC wrote the file: it moves the point to the bottom of
// the box for the top row.
func TestTextPoint(t *testing.T) {
	origin := func(system string) float64 {
		items, _ := draw(t, fmt.Sprintf(stepHeader, system)+stepText, false)
		if len(items) != 1 {
			t.Fatalf("items %q", items)
		}
		// the text's transform: the box's bottom plus 0.12 em to the baseline
		var m [6]float64
		if _, err := fmt.Sscanf(items[0][strings.Index(items[0], "m=["):], "m=[%g %g %g %g %g %g]", &m[0], &m[1], &m[2], &m[3], &m[4], &m[5]); err != nil {
			t.Fatal(items[0], err)
		}
		if m[4] != 10 {
			t.Errorf("%s: x %v", system, m[4])
		}
		return m[5] - 0.12*5
	}
	if got := origin("CAD 1.0"); math.Abs(got-15) > 0.01 {
		t.Errorf("the box's bottom at %v, want 15", got)
	}
	if got := origin("SCADEC_API_Ver3.30$$3.1"); math.Abs(got-20) > 0.01 {
		t.Errorf("SCADEC: the box's bottom at %v, want 20", got)
	}
	// a text not named by SXF sits on its baseline
	items, _ := draw(t, fmt.Sprintf(stepHeader, "CAD")+strings.Replace(stepText, "'$$SXF_topline left'", "' '", 1), false)
	if len(items) != 1 || !strings.Contains(items[0], "m=[5 0 0 5 10 20]") {
		t.Errorf("items %q", items)
	}
}

// TestHostile converts files made to exhaust the converter: each must
// finish in time.
func TestHostile(t *testing.T) {
	head := fmt.Sprintf(stepHeader, "CAD")
	var callouts, composites, invisible strings.Builder
	// a chain of callouts, each holding the next
	const chain = 20000
	for i := 1; i < chain; i++ {
		fmt.Fprintf(&callouts, "#%d=(DRAUGHTING_CALLOUT((#%d))REPRESENTATION_ITEM(' '));\n", i, i+1)
	}
	fmt.Fprintf(&callouts, "#%d=DRAWING_SHEET_REVISION('s',(#1),$,'01');\n", chain)
	// composite curves, each made of the one before twice
	composites.WriteString("#1=CARTESIAN_POINT(' ',(0.0,0.0));\n#2=CARTESIAN_POINT(' ',(1.0,1.0));\n#3=POLYLINE(' ',(#1,#2));\n")
	prev := 3
	for i := range 30 {
		id := 100 + 10*i
		fmt.Fprintf(&composites, "#%d=COMPOSITE_CURVE_SEGMENT(.CONTINUOUS.,.T.,#%d);\n#%d=COMPOSITE_CURVE(' ',(#%d,#%d),.T.);\n", id, prev, id+1, id, id)
		prev = id + 1
	}
	fmt.Fprintf(&composites, "#900=PRESENTATION_STYLE_ASSIGNMENT((#901));\n#901=CURVE_STYLE(' ',$,#902,$);\n#902=LENGTH_MEASURE_WITH_UNIT(#902,$);\n")
	fmt.Fprintf(&composites, "#903=(ANNOTATION_CURVE_OCCURRENCE()STYLED_ITEM((#900),#%d));\n#904=DRAWING_SHEET_REVISION('s',(#903,#903),$,'01');\n", prev)
	// a spline of a huge degree
	var ctrl []string
	for i := range 2000 {
		fmt.Fprintf(&composites, "#%d=CARTESIAN_POINT(' ',(%d.0,0.0));\n", 10000+i, i)
		ctrl = append(ctrl, fmt.Sprintf("#%d", 10000+i))
	}
	// an invisibility that lists the same layer many times
	invisible.WriteString("#1=CARTESIAN_POINT(' ',(0.0,0.0));\n#2=POLYLINE(' ',(#1,#1));\n#3=(ANNOTATION_CURVE_OCCURRENCE()STYLED_ITEM((),#2));\n")
	invisible.WriteString("#4=PRESENTATION_LAYER_ASSIGNMENT('l',' ',(" + strings.TrimSuffix(strings.Repeat("#3,", 5000), ",") + "));\n")
	invisible.WriteString("#5=INVISIBILITY((" + strings.TrimSuffix(strings.Repeat("#4,", 5000), ",") + "));\n#6=DRAWING_SHEET_REVISION('s',(#3),$,'01');\n")
	spline := head + composites.String()[:0] + strings.Join([]string{
		"#1=B_SPLINE_CURVE_WITH_KNOTS(' ',100000000,(" + strings.Join(ctrl, ",") + "),.UNSPECIFIED.,.F.,.F.,(),());",
		"#2=(ANNOTATION_CURVE_OCCURRENCE()STYLED_ITEM((),#1));",
		"#3=DRAWING_SHEET_REVISION('s',(#2),$,'01');",
	}, "\n")
	for i := range 2000 {
		spline += fmt.Sprintf("\n#%d=CARTESIAN_POINT(' ',(%d.0,0.0));", 10000+i, i)
	}
	// compound figures placing the one below ten times, around a long
	// polyline
	var sfc strings.Builder
	sfc.WriteString(head)
	xs := strings.TrimSuffix(strings.Repeat("0.0,1.0,", 1000), ",")
	fmt.Fprintf(&sfc, "/*SXF\n#1 = polyline_feature('1','8','1','1','2000','(%s)','(%s)')\nSXF*/\n", xs, xs)
	fmt.Fprintf(&sfc, "/*SXF\n#2 = sfig_org_feature(\\'f0\\','1')\nSXF*/\n")
	for level := 1; level <= 6; level++ {
		for range 10 {
			fmt.Fprintf(&sfc, "/*SXF\n#3 = sfig_locate_feature('0',\\'f%d\\','1','1','0','1','1')\nSXF*/\n", level-1)
		}
		fmt.Fprintf(&sfc, "/*SXF\n#4 = sfig_org_feature(\\'f%d\\','1')\nSXF*/\n", level)
	}
	fmt.Fprintf(&sfc, "/*SXF\n#5 = sfig_locate_feature('0',\\'f6\\','1','1','0','1','1')\nSXF*/\n")
	fmt.Fprintf(&sfc, "/*SXF\n#6 = drawing_sheet_feature(\\'s\\','3','1','0','0')\nSXF*/\n")
	want := map[string]string{"callouts": "nested deeper", "composites": "are cut", "figures": "too large"}
	for name, text := range map[string]string{
		"callouts":   head + callouts.String(),
		"composites": head + composites.String(),
		"spline":     spline,
		"invisible":  head + invisible.String(),
		"figures":    sfc.String(),
	} {
		t.Run(name, func(t *testing.T) {
			done := make(chan struct{})
			go func() {
				defer close(done)
				res, err := Convert(strings.NewReader(text), int64(len(text)), &Options{NoSystemFonts: true, NoTextIndex: true})
				if err == nil && want[name] != "" && !strings.Contains(strings.Join(res.Warnings, "\n"), want[name]) {
					t.Errorf("warnings %q", res.Warnings)
				}
			}()
			select {
			case <-done:
			case <-time.After(20 * time.Second):
				t.Fatal("did not finish")
			}
		})
	}
}

// TestCycle draws compound figures that place themselves.
func TestCycle(t *testing.T) {
	p21 := fmt.Sprintf(stepHeader, "CAD") + `#1=CARTESIAN_POINT(' ',(0.0,0.0));
#2=AXIS2_PLACEMENT_2D(' ',#1,$);
#3=SYMBOL_TARGET(' ',#2,1.0,1.0);
#4=DRAUGHTING_SUBFIGURE_REPRESENTATION('$$SXF_FM_loop',(#20,#21,#22,#30,#2),$);
#5=SYMBOL_REPRESENTATION_MAP(#2,#4);
#6=(ANNOTATION_SYMBOL()GEOMETRIC_REPRESENTATION_ITEM()MAPPED_ITEM(#5,#3)REPRESENTATION_ITEM(' '));
#7=PRESENTATION_STYLE_ASSIGNMENT((NULL_STYLE(.NULL.)));
#20=(ANNOTATION_OCCURRENCE()ANNOTATION_SUBFIGURE_OCCURRENCE()STYLED_ITEM((#7),#6));
#21=(ANNOTATION_OCCURRENCE()ANNOTATION_SUBFIGURE_OCCURRENCE()STYLED_ITEM((#7),#6));
#22=(ANNOTATION_OCCURRENCE()ANNOTATION_SUBFIGURE_OCCURRENCE()STYLED_ITEM((#7),#6));
#23=CARTESIAN_POINT(' ',(1.0,1.0));
#24=POLYLINE(' ',(#1,#23));
#30=(ANNOTATION_CURVE_OCCURRENCE()ANNOTATION_OCCURRENCE()STYLED_ITEM((#7),#24));
#40=DRAWING_SHEET_REVISION('s',(#20),$,'01');
ENDSEC;
END-ISO-10303-21;
`
	var sfc strings.Builder
	sfc.WriteString(fmt.Sprintf(stepHeader, "CAD"))
	for _, f := range []string{
		`line_feature('1','8','1','1','0','0','1','1')`,
		`sfig_locate_feature('0',\'loop\','0','0','0','1','1')`,
		`sfig_locate_feature('0',\'loop\','0','0','0','1','1')`,
		`sfig_locate_feature('0',\'loop\','0','0','0','1','1')`,
		`sfig_org_feature(\'loop\','1')`,
		`sfig_locate_feature('0',\'loop\','0','0','0','1','1')`,
		`drawing_sheet_feature(\'s\','4','0','0','0')`,
	} {
		fmt.Fprintf(&sfc, "/*SXF\n#10 = %s\nSXF*/\n", f)
	}
	for name, text := range map[string]string{"p21": p21, "sfc": sfc.String()} {
		done := make(chan struct{})
		var items, warnings []string
		go func() {
			items, warnings = draw(t, text, false)
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatalf("%s: a compound figure that places itself did not finish", name)
		}
		if len(items) != 1 || len(warnings) != 1 || !strings.Contains(warnings[0], "place themselves") {
			t.Errorf("%s: items %q, warnings %q", name, items, warnings)
		}
	}
}

func TestMalformed(t *testing.T) {
	for name, data := range map[string]string{
		"empty":           "",
		"junk":            "hello",
		"no data":         "ISO-10303-21;\nHEADER;\nENDSEC;\n",
		"no sheet":        fmt.Sprintf(stepHeader, "CAD") + "#1=CARTESIAN_POINT(' ',(0.0,0.0));\nENDSEC;\n",
		"unclosed string": fmt.Sprintf(stepHeader, "CAD") + "#1=CARTESIAN_POINT('x,(0.0,0.0));\n",
		"unclosed list":   fmt.Sprintf(stepHeader, "CAD") + "#1=CARTESIAN_POINT(' ',(0.0,0.0",
		"deep list":       fmt.Sprintf(stepHeader, "CAD") + "#1=X(" + strings.Repeat("(", 100000),
		"bad references":  fmt.Sprintf(stepHeader, "CAD") + "#1=DRAWING_SHEET_REVISION('s',(#2,#3,#99),$,'01');\n#2=(STYLED_ITEM((#1),#1));\n#3=STYLED_ITEM(' ',(#2),#3);\n",
		"sfc garbage":     "ISO-10303-21;\nDATA;\n/*SXF\n#10 = line_feature('1','x'\nSXF*/\n/*SXF\n#20 = drawing_sheet_feature(\\'s\\','9')\nSXF*/\n",
		"zip":             "PK\x03\x04garbage",
	} {
		t.Run(name, func(t *testing.T) {
			// errors are fine, panics are not
			Convert(strings.NewReader(data), int64(len(data)), testOptions())
		})
	}
	// cut files convert what they have
	for _, name := range []string{"shapes.sfc", "shapes.p21"} {
		b, err := os.ReadFile(filepath.Join("testdata", name))
		if err != nil {
			t.Fatal(err)
		}
		for _, n := range []int{len(b) / 3, len(b) / 2, len(b) - 100} {
			Convert(bytes.NewReader(b[:n]), int64(n), testOptions())
		}
	}
}

func TestDetect(t *testing.T) {
	for _, c := range []struct {
		name string
		want bool
	}{{"shapes.sfc", true}, {"shapes.p21", true}, {"shapes.p2z", true}} {
		b, err := os.ReadFile(filepath.Join("testdata", c.name))
		if err != nil {
			t.Fatal(err)
		}
		if got := Detect(b[:min(len(b), 1024)], bytes.NewReader(b), int64(len(b))); got != c.want {
			t.Errorf("%s: %v", c.name, got)
		}
	}
	ap214 := []byte("ISO-10303-21;\nHEADER;\nFILE_SCHEMA(('AUTOMOTIVE_DESIGN'));\n")
	if Detect(ap214, bytes.NewReader(ap214), int64(len(ap214))) {
		t.Error("a 3D STEP file is detected")
	}
}

func TestStepString(t *testing.T) {
	for in, want := range map[string]string{
		`plain`:                        "plain",
		`\X2\30C630B930C8\X0\ A`:       "テスト A",
		`\X4\0001F600\X0\`:             "😀",
		`caf\X\E9`:                     "café",
		`\S\a`:                         "á",
		`a\\b`:                         `a\b`,
		`\PA\x`:                        "x",
		`\X2\D83DDE00\X0\`:             "😀",
		`\X2\30C6`:                     "テ",
		`\X2\ZZ\X0\ok`:                 "ok",
		`\X2\30C630B9\X0\\X2\30C8\X0\`: "テスト",
	} {
		if got := decodeStepString(in); got != want {
			t.Errorf("%s: %q, want %q", in, got, want)
		}
	}
}

func TestArcEnds(t *testing.T) {
	const turn = 2 * math.Pi
	for _, c := range []struct {
		t0, t1 float64
		cw     bool
		sweep  float64
	}{
		{0, math.Pi / 2, false, math.Pi / 2},
		{0, math.Pi / 2, true, -3 * math.Pi / 2},
		{3 * math.Pi / 2, math.Pi / 2, false, math.Pi},
		{0, 0, false, turn},
		{0, 6.28318530717959, false, turn}, // 2π written a hair over
		{1, 1 + turn, true, -turn},
		{-math.Pi / 2, math.Pi / 2, false, math.Pi},
	} {
		a, b := arcEnds(c.t0, c.t1, c.cw)
		if math.Abs((b-a)-c.sweep) > 1e-9 {
			t.Errorf("%v → %v (cw %v): sweep %v, want %v", c.t0, c.t1, c.cw, b-a, c.sweep)
		}
	}
}
