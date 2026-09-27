package font

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/internal/sfnt"
	"github.com/shibukawa/bdf/woff2"
)

// The test fonts are made by test/font/gen.py from STIX Two (SIL Open Font
// License 1.1); features.ttf is a font of rectangles with a lookup of each
// kind. The UI text is laid out with the fonts of the PowerPoint tests.

func testData(t testing.TB, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func testOptions(params map[string]string) *conv.Options {
	return &conv.Options{FontDirs: []string{"../pptx/testdata/fonts"}, NoSystemFonts: true, Params: params}
}

// convertTest converts a test font, failing on warnings unless allowed.
func convertTest(t *testing.T, name string, o *conv.Options, warnings ...string) *bdf.Document {
	t.Helper()
	if o == nil {
		o = testOptions(nil)
	}
	var got []string
	o.Warn = func(msg string) { got = append(got, msg) }
	res, err := convert(testData(t, name), o, o.Warn)
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range got {
		ok := false
		for _, want := range warnings {
			ok = ok || strings.Contains(w, want)
		}
		if !ok {
			t.Errorf("%s: warning: %s", name, w)
		}
	}
	// the document writes and reads back
	var buf bytes.Buffer
	if err := res.Doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	if _, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len())); err != nil {
		t.Fatal(err)
	}
	return res.Doc
}

// viewText returns the text of a view.
func viewText(t *testing.T, doc *bdf.Document, id string) string {
	t.Helper()
	st, err := doc.SearchText()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range st.Views {
		if v.ID == id {
			var b strings.Builder
			for _, p := range v.Pages {
				b.WriteString(p.Text)
				b.WriteByte('\n')
			}
			return b.String()
		}
	}
	t.Fatalf("no view %s", id)
	return ""
}

func viewIDs(doc *bdf.Document) []string {
	var ids []string
	for _, v := range doc.Views {
		ids = append(ids, v.ID)
	}
	return ids
}

func TestDetect(t *testing.T) {
	for _, name := range []string{"stix.ttf", "stix.otf", "bold.ttc", "stix.woff", "stix.woff2", "features.ttf"} {
		b := testData(t, name)
		if !detect(b[:min(len(b), 1024)]) {
			t.Errorf("%s is not detected", name)
		}
		if f := conv.Detect(bytes.NewReader(b), int64(len(b))); f == nil || f.Name != "font" {
			t.Errorf("%s is detected as %v", name, f)
		}
	}
	for _, b := range [][]byte{
		[]byte("\x89PNG\r\n\x1a\n\x00\x00\x00\x0dIHDR"),
		append([]byte{0, 1, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0}, bytes.Repeat([]byte{1, 2, 3, 4}, 12)...),
		append([]byte("OTTO\x01\x00"), make([]byte, 20)...),
		[]byte("ttcf\x00\x05\x00\x00\x00\x00\x00\x01"),
	} {
		if detect(b) {
			t.Errorf("%q is detected as a font", b[:8])
		}
	}
}

func TestViews(t *testing.T) {
	doc := convertTest(t, "stix.ttf", nil)
	if ids := viewIDs(doc); !slices.Equal(ids, []string{"overview", "characters", "glyphs", "features"}) {
		t.Fatalf("views %v", ids)
	}
	if doc.Meta.Source != "font" || doc.Meta.DC.Title.First() != "STIX Two Text Regular" || !strings.Contains(doc.Meta.DC.Rights.First(), "STIX") {
		t.Errorf("meta %+v", doc.Meta)
	}
	for _, v := range doc.Views {
		if v.Kind != bdf.ViewScroll || len(v.Pages) == 0 || v.TextIndex == "" {
			t.Errorf("view %s: %s, %d pages, text index %q", v.ID, v.Kind, len(v.Pages), v.TextIndex)
		}
		for _, p := range v.Pages {
			if p.W != pageW || len(p.Layers) != 1 {
				t.Errorf("view %s: page %+v", v.ID, p)
			}
		}
	}
	ov := viewText(t, doc, "overview")
	for _, want := range []string{"STIX Two Text Regular", "The quick brown fox jumps over the lazy dog.", "Τάχιστη αλώπηξ",
		"Съешь же", "SIL Open Font License", "Installable", "TrueType outlines", "Latin (latn)",
		"kern mark mkmk", "GPOS"} {
		if !strings.Contains(ov, want) {
			t.Errorf("the overview lacks %q", want)
		}
	}
	// the characters, from their alternative text
	ch := viewText(t, doc, "characters")
	for _, want := range []string{"Basic Latin", "Greek and Coptic", "95 of 95", "004x @ A B C D E F G", "α β γ", "Ж"} {
		if !strings.Contains(ch, want) {
			t.Errorf("the characters lack %q", want)
		}
	}
	gl := viewText(t, doc, "glyphs")
	for _, want := range []string{"906 glyphs", ".notdef", "a.smcp", "U+0041"} {
		if !strings.Contains(gl, want) {
			t.Errorf("the glyphs lack %q", want)
		}
	}
	fe := viewText(t, doc, "features")
	for _, want := range []string{"Small Capitals", "Standard Ligatures", "on by default", "-100"} {
		if !strings.Contains(fe, want) {
			t.Errorf("the features lack %q", want)
		}
	}
}

// fontParts returns the font parts of a document parsed, with their hashes.
func fontParts(t *testing.T, doc *bdf.Document) map[bdf.Hash]*sfnt.Font {
	t.Helper()
	out := map[bdf.Hash]*sfnt.Font{}
	for _, p := range doc.Parts() {
		if p.Type != bdf.PartFont {
			continue
		}
		data := p.Data
		if woff2.IsWOFF(data) {
			var err error
			if data, err = woff2.Decode(data); err != nil {
				t.Fatal(err)
			}
		}
		f, err := sfnt.Parse(data)
		if err != nil {
			t.Fatal(err)
		}
		out[p.Hash] = f
	}
	return out
}

func TestGlyphFont(t *testing.T) {
	fl, err := load(testData(t, "features.ttf"))
	if err != nil {
		t.Fatal(err)
	}
	fc := fl.faces[0]
	f, err := sfnt.Parse(glyphProgram(fc))
	if err != nil {
		t.Fatal(err)
	}
	for _, tag := range []string{"GSUB", "GPOS", "GDEF"} {
		if f.Tables[tag] != nil {
			t.Errorf("the glyph font has %s", tag)
		}
	}
	if f.NumGlyphs != fc.f.NumGlyphs || len(f.Cmap) != f.NumGlyphs-1 {
		t.Fatalf("%d glyphs, %d characters", f.NumGlyphs, len(f.Cmap))
	}
	for g := 1; g < f.NumGlyphs; g++ {
		if f.Cmap[uint32(puaRune(uint16(g)))] != uint16(g) {
			t.Errorf("glyph %d is not mapped", g)
		}
	}
	if puaRune(0xFFFD) != 0xFFFFD || puaRune(0xFFFE) != 0x100000 {
		t.Error("private use characters of the last glyphs")
	}
}

func TestSampleFont(t *testing.T) {
	for _, name := range []string{"stix.ttf", "stix.otf"} {
		fl, err := load(testData(t, name))
		if err != nil {
			t.Fatal(err)
		}
		fc := fl.faces[0]
		data, ok := sampleProgram(fc, "af", true)
		if !ok {
			t.Fatalf("%s: no sample font", name)
		}
		f, err := sfnt.Parse(data)
		if err != nil {
			t.Fatal(err)
		}
		if f.NumGlyphs != fc.f.NumGlyphs || f.Tables["GSUB"] == nil || f.Tables["GPOS"] == nil || len(f.Cmap) != len(fc.f.Cmap) {
			t.Fatalf("%s: %d glyphs, GSUB %v, %d characters", name, f.NumGlyphs, f.Tables["GSUB"] != nil, len(f.Cmap))
		}
		out := sfnt.NewOutlines(f)
		has := func(r rune) bool {
			_, ok := out.Bounds(fc.glyph(r))
			return ok
		}
		// a and f are kept, with what substitutions make of a (a.smcp);
		// Z is emptied
		if !has('a') || !has('f') || has('Z') {
			t.Errorf("%s: a %v, f %v, Z %v", name, has('a'), has('f'), has('Z'))
		}
		if g := slices.Index(fc.glyphs, "a.smcp"); g < 0 {
			t.Errorf("%s: no a.smcp", name)
		} else if _, ok := out.Bounds(uint16(g)); !ok {
			t.Errorf("%s: a.smcp is dropped", name)
		}
		if len(data) >= len(fc.data) {
			t.Errorf("%s: the sample font (%d bytes) is not smaller than the font (%d)", name, len(data), len(fc.data))
		}
	}
}

func TestRestricted(t *testing.T) {
	doc := convertTest(t, "restricted.ttf", nil, "does not allow embedding")
	if n := len(fontParts(t, doc)); n > 3 {
		t.Errorf("%d font parts: the font is embedded", n)
	}
	// the glyphs are paths, with their characters as alternative text
	if ch := viewText(t, doc, "characters"); !strings.Contains(ch, "004x @ A B C D E F G") {
		t.Errorf("characters %q", ch[:min(200, len(ch))])
	}
	if ov := viewText(t, doc, "overview"); !strings.Contains(ov, "Restricted") || !strings.Contains(ov, "The quick brown fox") {
		t.Error("the overview of a restricted font")
	}
	o := testOptions(nil)
	o.IgnoreFSType = true
	doc = convertTest(t, "restricted.ttf", o)
	found := false
	for _, f := range fontParts(t, doc) {
		found = found || f.NumGlyphs == 292 && f.Tables["GSUB"] == nil
	}
	if !found {
		t.Error("the glyph font is not embedded with IgnoreFSType")
	}
}

func TestCollection(t *testing.T) {
	doc := convertTest(t, "bold.ttc", nil)
	want := []string{"overview-1", "characters-1", "glyphs-1", "features-1", "overview-2", "characters-2", "glyphs-2", "features-2"}
	if ids := viewIDs(doc); !slices.Equal(ids, want) {
		t.Fatalf("views %v", ids)
	}
	if doc.Views[4].Title != "STIX Two Text Bold: Overview" {
		t.Errorf("title %q", doc.Views[4].Title)
	}
	// one font of the collection is shown as a font of its own
	doc = convertTest(t, "bold.ttc", testOptions(map[string]string{"font": "2"}))
	if ids := viewIDs(doc); !slices.Equal(ids, []string{"overview", "characters", "glyphs", "features"}) || doc.Meta.DC.Title.First() != "STIX Two Text Bold" {
		t.Fatalf("views %v of %q", ids, doc.Meta.DC.Title.First())
	}
	if _, err := convert(testData(t, "bold.ttc"), testOptions(map[string]string{"font": "3"}), func(string) {}); err == nil {
		t.Error("font 3 of 2")
	}
}

// TestSharedGlyphs makes a collection of one font twice: the fonts share
// their tables, and one glyph font draws both.
func TestSharedGlyphs(t *testing.T) {
	data := collection(testData(t, "features.ttf"), 2)
	res, err := convert(data, testOptions(nil), func(msg string) { t.Error(msg) })
	if err != nil {
		t.Fatal(err)
	}
	glyphFonts := 0
	for _, f := range fontParts(t, res.Doc) {
		if f.Tables["GSUB"] == nil && f.NumGlyphs == 15 {
			glyphFonts++
		}
	}
	if glyphFonts != 1 || len(res.Doc.Views) != 8 {
		t.Errorf("%d glyph fonts, %d views", glyphFonts, len(res.Doc.Views))
	}
}

// collection makes a TTC of n fonts that are font, sharing its tables.
func collection(font []byte, n int) []byte {
	numTables := int(binary.BigEndian.Uint16(font[4:]))
	hdr := 12 + 4*n
	dir := 12 + 16*numTables
	var b bytes.Buffer
	b.WriteString("ttcf")
	binary.Write(&b, binary.BigEndian, []uint16{1, 0})
	binary.Write(&b, binary.BigEndian, uint32(n))
	for i := range n {
		binary.Write(&b, binary.BigEndian, uint32(hdr+i*dir))
	}
	shift := uint32(hdr + n*dir)
	for range n {
		d := append([]byte(nil), font[:dir]...)
		for k := range numTables {
			off := binary.BigEndian.Uint32(d[12+k*16+8:])
			binary.BigEndian.PutUint32(d[12+k*16+8:], off-uint32(dir)+shift)
		}
		b.Write(d)
	}
	b.Write(font[dir:])
	return b.Bytes()
}

func TestWebFonts(t *testing.T) {
	for _, name := range []string{"stix.woff", "stix.woff2"} {
		fl, err := load(testData(t, name))
		if errors.Is(err, woff2.ErrDecodeNotAvailable) || err != nil && strings.Contains(err.Error(), "Brotli") {
			t.Logf("%s: %v", name, err)
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		fc := fl.faces[0]
		if fl.container == "" || fc.f.NumGlyphs != 292 || len(fc.runes) != 95 || fc.gsub == nil {
			t.Errorf("%s: %s, %d glyphs, %d characters", name, fl.container, fc.f.NumGlyphs, len(fc.runes))
		}
	}
}

func TestVariableAndColor(t *testing.T) {
	ov := viewText(t, convertTest(t, "variable.ttf", nil), "overview")
	for _, want := range []string{"Axis wght Weight", "400 to 700, default 400", "Instance Bold", "the default instance"} {
		if !strings.Contains(ov, want) {
			t.Errorf("the overview of the variable font lacks %q", want)
		}
	}
	doc := convertTest(t, "color.ttf", nil)
	ov = viewText(t, doc, "overview")
	for _, want := range []string{"COLRv0 (3 glyphs)", "CPAL (2 palettes of 2 colors)", "Palette 1"} {
		if !strings.Contains(ov, want) {
			t.Errorf("the overview of the color font lacks %q", want)
		}
	}
	// the glyph font keeps the color tables
	found := false
	for _, f := range fontParts(t, doc) {
		found = found || f.Tables["COLR"] != nil && f.Tables["GSUB"] == nil
	}
	if !found {
		t.Error("the glyph font has no COLR table")
	}
}

func TestFeatureExamples(t *testing.T) {
	fl, err := load(testData(t, "features.ttf"))
	if err != nil {
		t.Fatal(err)
	}
	c := &converter{fc: fl.faces[0], examples: defaultExamples}
	c.plan = c.fc.plan("")
	groups := map[string]*featureGroup{}
	for _, g := range c.fc.featureGroups() {
		groups[g.table+" "+g.tag] = g
	}
	glyphs := func(ps []placed) []uint16 {
		var out []uint16
		for _, p := range ps {
			out = append(out, p.g)
		}
		return out
	}
	check := func(key string, want [][2][]uint16) []item {
		t.Helper()
		items, total := c.featureExamples(groups[key])
		if total != len(want) || len(items) != len(want) {
			t.Fatalf("%s: %d examples of %d, want %d", key, len(items), total, len(want))
		}
		for i, w := range want {
			if b, a := glyphs(items[i].before), glyphs(items[i].after); !slices.Equal(b, w[0]) || !slices.Equal(a, w[1]) {
				t.Errorf("%s: example %d is %v → %v, want %v → %v", key, i, b, a, w[0], w[1])
			}
		}
		return items
	}
	check("GSUB liga", [][2][]uint16{{{7, 7, 8}, {13}}, {{7, 8}, {12}}})
	if it := check("GSUB salt", [][2][]uint16{{{2}, {9, 10, 11}}}); !it[0].alts {
		t.Error("salt: alternates are not choices")
	}
	check("GSUB calt", [][2][]uint16{{{2}, {9}}})
	check("GSUB rclt", [][2][]uint16{{{4}, {6}}})
	// the pairs of the sample text first (d e), then the others, largest
	// adjustments first
	items := check("GPOS kern", [][2][]uint16{{{2, 3}, {2, 3}}, {{5, 6}, {5, 6}}, {{4, 6}, {4, 6}}, {{4, 7}, {4, 7}}, {{5, 7}, {5, 7}}})
	if items[0].note != "-50" || items[0].after[1].x != 0.45 {
		t.Errorf("kern a b: %+v", items[0])
	}
	items = check("GPOS mark", [][2][]uint16{{{2, 14}, {2, 14}}, {{3, 14}, {3, 14}}})
	if m := items[0].after[1]; m.x != 0.15 || m.y != -0.05 {
		t.Errorf("mark on a at %v, %v", m.x, m.y)
	}
	if items := check("GPOS palt", [][2][]uint16{{{2}, {2}}}); !items[0].boxes || items[0].advA != 0.48 || items[0].offset != -0.01 {
		t.Errorf("palt %+v", items[0])
	}
	g := groups["GSUB cv01"]
	if g.uiName != "Open a" || string(g.chars) != "ab" {
		t.Errorf("cv01: %q %q", g.uiName, string(g.chars))
	}
	if g := groups["GSUB locl"]; scriptList(g.scripts) != "latn (TRK)" {
		t.Errorf("locl scripts %q", scriptList(g.scripts))
	}
}

func TestOneStrike(t *testing.T) {
	const glyphs = 2
	// strikes of 20, 64 and 160 ppem, each with the data of its glyphs
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, []uint16{1, 1})
	binary.Write(&b, binary.BigEndian, uint32(3))
	strike := func(ppem uint16) []byte {
		var s bytes.Buffer
		binary.Write(&s, binary.BigEndian, []uint16{ppem, 72})
		start := 4 + 4*(glyphs+1)
		binary.Write(&s, binary.BigEndian, []uint32{uint32(start), uint32(start + 3), uint32(start + 6)})
		s.Write([]byte{byte(ppem), 1, 2, byte(ppem), 3, 4})
		return s.Bytes()
	}
	var strikes [][]byte
	for _, ppem := range []uint16{20, 64, 160} {
		strikes = append(strikes, strike(ppem))
	}
	off := 8 + 4*3
	for _, s := range strikes {
		binary.Write(&b, binary.BigEndian, uint32(off))
		off += len(s)
	}
	for _, s := range strikes {
		b.Write(s)
	}
	got := oneStrike(b.Bytes(), glyphs)
	want := append([]byte{0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 12}, strikes[1]...)
	if !bytes.Equal(got, want) {
		t.Errorf("one strike %v, want %v", got, want)
	}
	if oneStrike(strikes[0], glyphs) != nil {
		t.Error("a table of one strike is rewritten")
	}
}

func TestParams(t *testing.T) {
	doc := convertTest(t, "stix.ttf", testOptions(map[string]string{"text": "Hamburgefonstiv", "examples": "3"}))
	if ov := viewText(t, doc, "overview"); !strings.Contains(ov, "Hamburgefonstiv") {
		t.Error("the sample text is not shown")
	}
	if fe := viewText(t, doc, "features"); !strings.Contains(fe, "… and ") {
		t.Error("the examples are not limited")
	}
	if _, err := convert(testData(t, "stix.ttf"), testOptions(map[string]string{"examples": "many"}), func(string) {}); err == nil {
		t.Error("examples=many is accepted")
	}
}

// FuzzLoad reads damaged fonts and what the views show of them (a
// document is made only of the fonts that are small: drawing is slow and
// works with what this reads): it may fail, but must not panic.
func FuzzLoad(f *testing.F) {
	for _, name := range []string{"features.ttf", "stix.woff"} {
		f.Add(testData(f, name))
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		fl, err := load(b)
		if err != nil {
			return
		}
		small := true
		for i, fc := range fl.faces {
			if i == 4 || fc.f.NumGlyphs > 5000 {
				return // slow to read, and made of the same damage
			}
			c := &converter{fc: fc, examples: 50}
			c.plan = fc.plan("")
			fc.sections(fl, false)
			fc.tableList()
			fc.coverage()
			fc.scriptCounts()
			fc.palettes()
			for _, g := range fc.featureGroups() {
				c.featureExamples(g)
				for _, l := range g.lookups {
					c.lookupNote(g.t, l)
				}
			}
			glyphProgram(fc)
			sampleProgram(fc, c.plan.all(), true)
			small = small && fc.f.NumGlyphs <= 300 && len(fc.runes) <= 300 && len(fl.faces) <= 2
		}
		if small {
			convert(b, &conv.Options{NoSystemFonts: true, NoTextIndex: true, NoWOFF2: true}, func(string) {})
		}
	})
}
