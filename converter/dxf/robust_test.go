package dxf

import (
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"golang.org/x/text/encoding/charmap"
)

// dxfText joins group codes and values into a text DXF file.
func dxfText(pairs ...string) []byte {
	return []byte(strings.Join(pairs, "\n") + "\n")
}

// entities wraps entities in a file with an R2000 header.
func entities(pairs ...string) []byte {
	head := []string{"0", "SECTION", "2", "HEADER", "9", "$ACADVER", "1", "AC1015", "0", "ENDSEC", "0", "SECTION", "2", "ENTITIES"}
	return dxfText(append(append(head, pairs...), "0", "ENDSEC", "0", "EOF")...)
}

// drawing draws the model space of a file into a drawing.
func drawing(t *testing.T, data []byte) (*converter, *cad.Drawing) {
	t.Helper()
	tags, err := readTags(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	c := &converter{opts: &Options{}, warned: map[string]bool{}, ltscale: 1}
	c.opts.Warn = func(m string) { t.Log("warning:", m) }
	c.d = parse(tags)
	c.fonts = &cad.Fonts{Set: fontset.New(fontdb.New([]string{"../pptx/testdata/fonts"}, false), nil)}
	out := &cad.Drawing{}
	var deferred []func(cad.Rect)
	c.sheets()
	c.draw(c.model, c.topCtx(out, false, false, &deferred))
	return c, out
}

// within fails the test when f does not return in time.
func within(t *testing.T, d time.Duration, what string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() { f(); close(done) }()
	select {
	case <-done:
	case <-time.After(d):
		t.Fatalf("%s did not finish in %v", what, d)
	}
}

func TestInfiniteAngles(t *testing.T) {
	for _, angles := range [][2]string{{"1e20", "0"}, {"0", "-inf"}, {"nan", "90"}} {
		data := entities("0", "ARC", "8", "0", "10", "0", "20", "0", "40", "1", "50", angles[0], "51", angles[1],
			"0", "ELLIPSE", "8", "0", "10", "0", "20", "0", "11", "1", "21", "0", "40", "0.5", "41", angles[0], "42", angles[1])
		within(t, 2*time.Second, "arcs with angles "+angles[0]+", "+angles[1], func() { drawing(t, data) })
	}
	if _, _, ok := arcAngles(1e300, 0, false); !ok {
		// huge but finite angles are reduced, not refused
		t.Error("a finite angle was refused")
	}
}

func TestBlockCycle(t *testing.T) {
	data := dxfText("0", "SECTION", "2", "BLOCKS",
		"0", "BLOCK", "8", "0", "2", "A", "70", "0", "10", "0", "20", "0",
		"0", "LINE", "8", "0", "10", "0", "20", "0", "11", "1", "21", "1",
		"0", "INSERT", "8", "0", "2", "A", "10", "1", "20", "0",
		"0", "INSERT", "8", "0", "2", "A", "10", "0", "20", "1",
		"0", "INSERT", "8", "0", "2", "A", "10", "1", "20", "1",
		"0", "ENDBLK", "0", "ENDSEC",
		"0", "SECTION", "2", "ENTITIES",
		"0", "INSERT", "8", "0", "2", "A", "10", "0", "20", "0",
		"0", "ENDSEC", "0", "EOF")
	var c *converter
	var out *cad.Drawing
	within(t, 5*time.Second, "a block that inserts itself", func() { c, out = drawing(t, data) })
	if len(out.Items) != 1 {
		t.Errorf("%d items, want the block's line once", len(out.Items))
	}
	if !c.warned["cycle:A"] {
		t.Error("no warning for the cycle")
	}
}

func TestHatchCounts(t *testing.T) {
	// counts far beyond the data
	for _, code := range []string{"95", "78", "97"} {
		pairs := []string{"0", "HATCH", "8", "0", "100", "AcDbHatch", "70", "0", "91", "1", "92", "0", "93", "1",
			"72", "4", "94", "3", "73", "0", "74", "0", "95", "4", "96", "0"}
		for i := 0; i+1 < len(pairs); i += 2 {
			if pairs[i] == code {
				pairs[i+1] = "50000000"
			}
		}
		if code != "95" {
			pairs = append(pairs, code, "50000000")
		}
		within(t, 2*time.Second, "a hatch with a count of 50000000 in "+code, func() { drawing(t, entities(pairs...)) })
	}
	// an array of inserts whose size overflows
	data := dxfText("0", "SECTION", "2", "BLOCKS", "0", "BLOCK", "8", "0", "2", "B", "10", "0", "20", "0",
		"0", "LINE", "8", "0", "10", "0", "20", "0", "11", "1", "21", "1", "0", "ENDBLK", "0", "ENDSEC",
		"0", "SECTION", "2", "ENTITIES", "0", "INSERT", "8", "0", "2", "B", "70", "2000000000", "71", "2000000000",
		"0", "ENDSEC", "0", "EOF")
	within(t, 2*time.Second, "a huge MINSERT", func() { drawing(t, data) })
}

// TestSplineEdgeBefore2010 reads a hatch of AutoCAD 2000 whose last edge is
// a spline: the 97 after it counts the source objects, not fit points.
func TestSplineEdgeBefore2010(t *testing.T) {
	data := entities("0", "HATCH", "8", "0", "100", "AcDbEntity", "100", "AcDbHatch",
		"10", "0", "20", "0", "30", "0", "2", "ANSI31", "70", "0", "71", "1", "91", "1",
		"92", "1", "93", "2",
		"72", "1", "10", "0", "20", "0", "11", "10", "21", "0",
		"72", "4", "94", "3", "73", "0", "74", "0", "95", "8", "96", "4",
		"40", "0", "40", "0", "40", "0", "40", "0", "40", "1", "40", "1", "40", "1", "40", "1",
		"10", "10", "20", "0", "10", "10", "20", "10", "10", "0", "20", "10", "10", "0", "20", "0",
		"97", "1", "330", "B1",
		"75", "1", "76", "1", "52", "0", "41", "1", "77", "0", "78", "1",
		"53", "45", "43", "0", "44", "0", "45", "-2.2450640303", "46", "2.2450640303", "79", "0", "98", "0")
	c, _ := drawing(t, data)
	h := c.readHatch(c.d.entities[0], false)
	if len(h.paths) != 1 || len(h.lines) != 1 || h.solid {
		t.Errorf("paths %d, pattern lines %d, solid %v", len(h.paths), len(h.lines), h.solid)
	}
}

func TestPointModes(t *testing.T) {
	for _, mode := range []string{"33", "65", "97", "34"} {
		data := dxfText("0", "SECTION", "2", "HEADER", "9", "$PDMODE", "70", mode, "9", "$PDSIZE", "40", "1", "0", "ENDSEC",
			"0", "SECTION", "2", "ENTITIES", "0", "POINT", "8", "0", "10", "5", "20", "5", "0", "ENDSEC", "0", "EOF")
		if _, out := drawing(t, data); len(out.Items) == 0 {
			t.Errorf("PDMODE %s draws nothing", mode)
		}
	}
}

func TestTextDefaults(t *testing.T) {
	// centred text without its second alignment point is centred on the
	// insertion point
	_, out := drawing(t, entities("0", "TEXT", "8", "0", "10", "100", "20", "200", "40", "2.5", "1", "Hello", "72", "1"))
	b := out.Bounds()
	if c := (b.Min.X + b.Max.X) / 2; c < 99 || c > 101 {
		t.Errorf("text centred on x %v, want 100", c)
	}
	// a text without a width factor has the width 1 whatever its style
	data := dxfText("0", "SECTION", "2", "TABLES", "0", "TABLE", "2", "STYLE",
		"0", "STYLE", "2", "NARROW", "70", "0", "40", "0", "41", "0.5", "3", "txt", "0", "ENDTAB", "0", "ENDSEC",
		"0", "SECTION", "2", "ENTITIES", "0", "TEXT", "8", "0", "10", "0", "20", "0", "40", "10", "1", "WWWW", "7", "NARROW",
		"0", "ENDSEC", "0", "EOF")
	_, narrow := drawing(t, data)
	_, plain := drawing(t, entities("0", "TEXT", "8", "0", "10", "0", "20", "0", "40", "10", "1", "WWWW"))
	if nw, pw := narrow.Bounds().W(), plain.Bounds().W(); nw < pw*0.99 {
		t.Errorf("width %v in a narrow style, %v without", nw, pw)
	}
}

func TestAttributeDefinition(t *testing.T) {
	_, out := drawing(t, entities("0", "ATTDEF", "8", "0", "10", "0", "20", "0", "40", "3", "1", "value", "2", "TAG", "3", "prompt", "70", "0"))
	if len(out.Items) != 1 || out.Items[0].Text() != "TAG" {
		t.Errorf("an attribute definition outside a block shows %v", out.Items)
	}
}

func TestWesternNotShiftJIS(t *testing.T) {
	for _, words := range [][]string{{"Gebäude", "Längsschnitt", "Maß"}, {"Grundriß Erdgeschoß", "Wände"}, {"Höhe", "Größe"}} {
		pairs := []string{"0", "SECTION", "2", "HEADER", "9", "$ACADVER", "1", "AC1015", "9", "$DWGCODEPAGE", "3", "ANSI_1252",
			"0", "ENDSEC", "0", "SECTION", "2", "ENTITIES"}
		for _, w := range words {
			enc, _ := charmap.Windows1252.NewEncoder().String(w)
			pairs = append(pairs, "0", "TEXT", "8", "0", "1", enc)
		}
		tags, _ := readTags(dxfText(append(pairs, "0", "ENDSEC", "0", "EOF")...), nil)
		if name := decodeStrings(tags); name != "ansi_1252" {
			t.Errorf("%v read as %s", words, name)
		}
	}
}

// TestR12Viewport reads the view of an R12 viewport from its MVIEW
// extended data.
func TestR12Viewport(t *testing.T) {
	data := dxfText("0", "SECTION", "2", "HEADER", "9", "$ACADVER", "1", "AC1009",
		"9", "$PLIMMIN", "10", "0", "20", "0", "9", "$PLIMMAX", "10", "420", "20", "297", "0", "ENDSEC",
		"0", "SECTION", "2", "ENTITIES",
		"0", "TEXT", "8", "0", "10", "40", "20", "50", "40", "5", "1", "MODEL TEXT",
		"0", "TEXT", "8", "OFF", "10", "40", "20", "40", "40", "5", "1", "FROZEN HERE",
		"0", "VIEWPORT", "8", "0", "67", "1", "10", "210", "20", "148.5", "40", "420", "41", "297", "68", "1", "69", "1",
		"0", "VIEWPORT", "8", "0", "67", "1", "10", "210", "20", "148.5", "40", "200", "41", "100", "68", "2", "69", "2",
		"1001", "ACAD", "1000", "MVIEW", "1002", "{", "1070", "16",
		"1010", "0", "1020", "0", "1030", "0", "1010", "0", "1020", "0", "1030", "1",
		"1040", "0", "1040", "100", "1040", "50", "1040", "50",
		"1040", "50", "1040", "0", "1040", "0", "1070", "0", "1070", "100", "1070", "1", "1070", "3", "1070", "0", "1070", "0",
		"1070", "0", "1070", "0", "1040", "0", "1040", "0", "1040", "0", "1040", "1", "1040", "1", "1040", "1", "1040", "1",
		"1070", "0", "1002", "{", "1003", "OFF", "1002", "}", "1002", "}",
		"0", "ENDSEC", "0", "EOF")
	_, r := convertData(t, data, nil)
	if n := len(r.Manifest.Views); n != 2 {
		t.Fatalf("%d views", n)
	}
	text := plainText(t, r, 1)
	if !strings.Contains(text, "MODEL TEXT") {
		t.Errorf("the viewport does not show model space: %q", text)
	}
	if strings.Contains(text, "FROZEN HERE") {
		t.Errorf("a layer frozen in the viewport is shown: %q", text)
	}
}
