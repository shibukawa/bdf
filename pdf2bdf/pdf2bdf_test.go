package pdf2bdf

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/shibukawa/bdf"
)

func TestLexer(t *testing.T) {
	l := &lexer{b: []byte(`/Name#20x 1 -2.5 (a\(b\)\n\101) <48656c6c6f> [1 /B (c)] << /K 1 >> Tj BI /W 2 /H 1 /BPC 8 /CS /G ID ab EI Q`)}
	var ops []string
	var operands []types.Object
	for {
		k, o, op := l.next()
		if k == tokEOF {
			break
		}
		if k == tokOperator {
			ops = append(ops, op)
			if op == "BI" {
				// skip the inline image dict manually
				for {
					k2, _, op2 := l.next()
					if k2 == tokEOF || op2 == "ID" {
						break
					}
				}
				data := l.readInlineImageData(2)
				if string(data) != "ab" {
					t.Fatalf("inline data = %q", data)
				}
			}
			continue
		}
		operands = append(operands, o)
	}
	if got := strings.Join(ops, " "); got != "Tj BI Q" {
		t.Fatalf("ops = %q", got)
	}
	if operands[0].(types.Name).Value() != "Name x" {
		t.Fatalf("name = %v", operands[0])
	}
	if string(literalBytes(operands[3])) != "a(b)\nA" {
		t.Fatalf("literal = %q", literalBytes(operands[3]))
	}
	if string(literalBytes(operands[4])) != "Hello" {
		t.Fatalf("hex = %q", literalBytes(operands[4]))
	}
	if a := operands[5].(types.Array); len(a) != 3 {
		t.Fatalf("array = %v", a)
	}
	if d := operands[6].(types.Dict); d["K"] != types.Integer(1) {
		t.Fatalf("dict = %v", d)
	}
}

func TestCMap(t *testing.T) {
	cm := parseCMap([]byte(`
/CIDInit /ProcSet findresource begin begincmap
1 begincodespacerange <0000> <FFFF> endcodespacerange
2 beginbfchar <0003> <0041> <0004> <00660069> endbfchar
1 beginbfrange <0010> <0012> <0061> endbfrange
1 beginbfrange <0020> <0021> [<0078> <0079>] endbfrange
endcmap`))
	if cm.toUnicode(3) != "A" || cm.toUnicode(4) != "fi" || cm.toUnicode(0x11) != "b" || cm.toUnicode(0x21) != "y" {
		t.Fatalf("mappings: %q %q %q %q", cm.toUnicode(3), cm.toUnicode(4), cm.toUnicode(0x11), cm.toUnicode(0x21))
	}
	codes := cm.decode([]byte{0, 3, 0, 0x11})
	if len(codes) != 2 || codes[0].code != 3 || codes[1].code != 0x11 || codes[0].nbytes != 2 {
		t.Fatalf("decode = %+v", codes)
	}
	enc := parseCMap([]byte(`1 begincodespacerange <00> <80> <8140> <9ffc> endcodespacerange 1 begincidrange <8140> <8150> 633 endcidrange`))
	codes = enc.decode([]byte{0x41, 0x81, 0x42})
	if len(codes) != 2 || codes[0].nbytes != 1 || codes[1].nbytes != 2 || codes[1].cid != 635 {
		t.Fatalf("mixed decode = %+v", codes)
	}
}

func TestSFNTRebuild(t *testing.T) {
	data, err := os.ReadFile("../fixture/fonts/DejaVuSans-sub.ttf")
	if err != nil {
		t.Skip("fixture font missing")
	}
	sf, err := parseSFNT(data)
	if err != nil {
		t.Fatal(err)
	}
	gidA, ok := sf.cmap['A']
	if !ok || gidA == 0 {
		t.Fatal("no cmap entry for A")
	}
	out := sf.rebuild(map[uint32]uint16{'Z': gidA, 0xE000: gidA, 0x1F600: gidA}, "Test", 700, true)
	sf2, err := parseSFNT(out)
	if err != nil {
		t.Fatal(err)
	}
	if sf2.cmap['Z'] != gidA || sf2.cmap[0xE000] != gidA || sf2.cmap[0x1F600] != gidA {
		t.Fatalf("rebuilt cmap = %v", sf2.cmap)
	}
	if _, has := sf2.cmap['A']; has {
		t.Fatal("old mapping survived")
	}
	for _, tag := range []string{"glyf", "loca", "head", "hhea", "hmtx", "maxp", "cmap", "name", "OS/2", "post"} {
		if sf2.tables[tag] == nil {
			t.Fatalf("missing table %s", tag)
		}
	}
	if be16(sf2.tables["OS/2"], 4) != 700 {
		t.Fatal("weight not written")
	}
}

func TestFunctions(t *testing.T) {
	f := &pdfFunction{kind: 2, domain: []float64{0, 1}, c0: []float64{0, 0, 1}, c1: []float64{1, 0, 0}, n: 1}
	if out := f.eval(0.25); out[0] != 0.25 || out[2] != 0.75 {
		t.Fatalf("exp = %v", out)
	}
	st := &pdfFunction{kind: 3, domain: []float64{0, 1}, funcs: []*pdfFunction{f, f}, bounds: []float64{0.5}, encode: []float64{0, 1, 1, 0}}
	if out := st.eval(0.75); out[0] != 0.5 {
		t.Fatalf("stitch = %v", out)
	}
	if got, err := PageRange("1-2,5,7-", 8); err != nil || len(got) != 5 || got[4] != 8 {
		t.Fatalf("PageRange = %v %v", got, err)
	}
}

// convert converts a testdata PDF and returns the document with a reader over it.
func convert(t *testing.T, name string) (*Result, *bdf.Reader) {
	t.Helper()
	res, err := ConvertFile(filepath.Join("testdata", name), nil)
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
	return res, r
}

func pageText(t *testing.T, r *bdf.Reader) string {
	t.Helper()
	v := r.Manifest.Views[0]
	h, err := bdf.ParseHash(v.TextIndex)
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

func opCounts(t *testing.T, r *bdf.Reader, h bdf.Hash) map[byte]int {
	t.Helper()
	o, err := r.Object(h)
	if err != nil {
		t.Fatal(err)
	}
	counts := map[byte]int{}
	if err := o.Walk(func(in bdf.Instr) { counts[in.Op]++ }); err != nil {
		t.Fatal(err)
	}
	return counts
}

func TestConvertChromeSlides(t *testing.T) {
	res, r := convert(t, "chrome-slides.pdf")
	if res.Pages != 2 {
		t.Fatalf("pages = %d", res.Pages)
	}
	for _, w := range res.Warnings {
		if !strings.Contains(w, "soft masks") {
			t.Errorf("unexpected warning: %s", w)
		}
	}
	v := r.Manifest.Views[0]
	if v.Kind != bdf.ViewFixed || len(v.Pages) != 2 || v.Pages[0].W != 720 {
		t.Fatalf("view = %+v", v)
	}
	fonts, images := 0, 0
	for _, e := range r.Manifest.Parts {
		switch e.T {
		case bdf.PartFont:
			fonts++
		case bdf.PartImage:
			images++
		}
	}
	if fonts < 4 || images < 2 {
		t.Fatalf("fonts=%d images=%d", fonts, images)
	}
	text := pageText(t, r)
	for _, want := range []string{"Shapes, gradients and text", "quick brown fox", "ligatures", "a link to example.com", "R2C2"} {
		if !strings.Contains(text, want) {
			t.Errorf("text index lacks %q", want)
		}
	}
	// The page content is a form XObject reused as a shared object; links are emitted on page 2.
	counts := opCounts(t, r, v.Pages[1].Layers[0].Obj)
	if counts[bdf.OpLink] != 1 || counts[bdf.OpUse] < 1 {
		t.Fatalf("page 2 ops: %v", counts)
	}
}

func TestConvertChromeDoc(t *testing.T) {
	_, r := convert(t, "chrome-doc.pdf")
	text := pageText(t, r)
	for _, want := range []string{"BDF conversion", "Each page is a list of objects", "Justified text", "hyperlink"} {
		if !strings.Contains(text, want) {
			t.Errorf("text index lacks %q", want)
		}
	}
	// Block-level marked content (P, TH …) becomes MARK PARAGRAPH. Chrome tags most
	// text as NonStruct/Span, so only a few marks are expected here.
	counts := opCounts(t, r, r.Manifest.Views[0].Pages[0].Layers[0].Obj)
	obj, _ := r.Object(r.Manifest.Views[0].Pages[0].Layers[0].Obj)
	marks := 0
	for _, child := range obj.Objects {
		marks += opCounts(t, r, child)[bdf.OpMark]
	}
	if counts[bdf.OpMark]+marks < 3 {
		t.Fatalf("expected paragraph marks, got %d", counts[bdf.OpMark]+marks)
	}
	// Kerned glyphs are merged into runs rather than one FILL_TEXT per glyph.
	if n := counts[bdf.OpFillText] + func() int {
		s := 0
		for _, child := range obj.Objects {
			s += opCounts(t, r, child)[bdf.OpFillText]
		}
		return s
	}(); n > 120 {
		t.Fatalf("too many text runs: %d", n)
	}
}

func TestConvertReportlab(t *testing.T) {
	res, r := convert(t, "reportlab-mixed.pdf")
	if len(res.Warnings) != 0 {
		t.Fatalf("warnings: %v", res.Warnings)
	}
	text := pageText(t, r)
	for _, want := range []string{"Helvetica: The quick brown fox", "café naïve — “quotes” • bullet ½ ©", "right aligned text", "rotated 30 degrees", "word spacing six between words", "invisible text (searchable)"} {
		if !strings.Contains(text, want) {
			t.Errorf("text index lacks %q", want)
		}
	}
	counts := opCounts(t, r, r.Manifest.Views[0].Pages[0].Layers[0].Obj)
	if counts[bdf.OpLink] != 1 || counts[bdf.OpImage] != 2 || counts[bdf.OpFillText] < 8 {
		t.Fatalf("page 1 ops: %v", counts)
	}
	p2 := opCounts(t, r, r.Manifest.Views[0].Pages[1].Layers[0].Obj)
	if p2[bdf.OpStrokeText] != 1 || p2[bdf.OpTextStyle] < 1 {
		t.Fatalf("page 2 ops: %v", p2)
	}
	// Simple TrueType fonts with a symbolic cmap get rebuilt with a Unicode cmap.
	embedded := 0
	for _, e := range r.Manifest.Parts {
		if e.T == bdf.PartFont {
			embedded++
			b, _ := r.Part(e.H)
			sf, err := parseSFNT(b)
			if err != nil {
				t.Fatal(err)
			}
			if sf.cmap == nil || len(sf.cmap) == 0 {
				t.Fatal("rebuilt font has no unicode cmap")
			}
		}
	}
	if embedded != 2 {
		t.Fatalf("embedded fonts = %d", embedded)
	}
}

func TestDeterministic(t *testing.T) {
	var outs [2][]byte
	for i := range outs {
		res, err := ConvertFile(filepath.Join("testdata", "chrome-slides.pdf"), nil)
		if err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		if err := res.Doc.WriteSingle(&buf); err != nil {
			t.Fatal(err)
		}
		outs[i] = buf.Bytes()
	}
	if !bytes.Equal(outs[0], outs[1]) {
		t.Fatal("conversion is not deterministic")
	}
}

func TestPruneGlyphs(t *testing.T) {
	data, err := os.ReadFile("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf")
	if err != nil {
		data, err = os.ReadFile("../fixture/fonts/DejaVuSans-sub.ttf")
		if err != nil {
			t.Skip("no font available")
		}
	}
	sf, err := parseSFNT(data)
	if err != nil {
		t.Fatal(err)
	}
	before := len(sf.tables["glyf"])
	// é is a composite of e and acute in DejaVu; both components must survive.
	keep := map[uint16]bool{sf.cmap['A']: true, sf.cmap[0xE9]: true}
	sf.pruneGlyphs(keep)
	after := len(sf.tables["glyf"])
	if after >= before/4 {
		t.Fatalf("glyf not pruned: %d -> %d", before, after)
	}
	out := sf.rebuild(map[uint32]uint16{'A': sf.cmap['A'], 0xE9: sf.cmap[0xE9]}, "Pruned", 400, false)
	sf2, err := parseSFNT(out)
	if err != nil {
		t.Fatal(err)
	}
	if sf2.numGlyphs != sf.numGlyphs {
		t.Fatal("glyph count changed")
	}
	loca, head := sf2.tables["loca"], sf2.tables["head"]
	long := be16(head, 50) != 0
	off := func(g uint16) (uint32, uint32) {
		if long {
			return be32(loca, int(g)*4), be32(loca, int(g)*4+4)
		}
		return uint32(be16(loca, int(g)*2)) * 2, uint32(be16(loca, int(g)*2+2)) * 2
	}
	if s, e := off(sf.cmap['A']); e <= s {
		t.Fatal("kept glyph is empty")
	}
	if s, e := off(sf.cmap['Z']); e != s {
		t.Fatal("unused glyph still has data")
	}
	// Components of the composite: the accent glyph must still have outlines.
	if acute, ok := sf.cmap[0xB4]; ok {
		if s, e := off(acute); e <= s {
			t.Fatal("composite component was pruned")
		}
	}
}
