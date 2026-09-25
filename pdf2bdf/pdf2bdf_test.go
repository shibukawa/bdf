package pdf2bdf

import (
	"bytes"
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/woff2"
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
	out := sf.rebuild(map[uint32]uint16{'Z': gidA, 0xE000: gidA, 0x1F600: gidA}, fontInfo{family: "Test", weight: 700, italic: true})
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
			sf := fontPart(t, b)
			if sf.cmap == nil || len(sf.cmap) == 0 {
				t.Fatal("rebuilt font has no unicode cmap")
			}
		}
	}
	if embedded != 2 {
		t.Fatalf("embedded fonts = %d", embedded)
	}
}

var woff2KnownTags = strings.Fields("cmap head hhea hmtx maxp name OS/2 post cvt_ fpgm glyf loca prep CFF_ VORG EBDT EBLC gasp hdmx kern LTSH PCLT VDMX vhea vmtx BASE GDEF GPOS GSUB EBSC JSTF MATH CBDT CBLC COLR CPAL SVG_ sbix acnt avar bdat bloc bsln cvar fdsc feat fmtx fvar gvar hsty just lcar mort morx opbd prop trak Zapf Silf Glat Gloc Feat Sill")

// fontPart parses a font part. A WOFF2 part is unpacked into an sfnt without
// its (transformed) glyf and loca tables, which the tests do not look at.
func fontPart(t *testing.T, b []byte) *sfnt {
	t.Helper()
	if woff2.IsWOFF(b) {
		if string(b[:4]) != "wOF2" {
			t.Fatal("WOFF 1 font part")
		}
		n := int(be16(b, 12))
		p := 48
		base128 := func() int {
			v := 0
			for {
				c := b[p]
				p++
				v = v<<7 | int(c&0x7f)
				if c&0x80 == 0 {
					return v
				}
			}
		}
		type entry struct {
			tag   string
			len   int
			xform bool
		}
		var entries []entry
		for i := 0; i < n; i++ {
			flags := b[p]
			p++
			tag := ""
			if flags&0x3f == 0x3f {
				tag = string(b[p : p+4])
				p += 4
			} else {
				tag = strings.ReplaceAll(woff2KnownTags[flags&0x3f], "_", " ")
			}
			e := entry{tag: tag, len: base128()}
			if (tag == "glyf" || tag == "loca") && flags>>6 == 0 {
				e.len, e.xform = base128(), true
			}
			entries = append(entries, e)
		}
		data, err := io.ReadAll(brotli.NewReader(bytes.NewReader(b[p : p+int(be32(b, 20))])))
		if err != nil {
			t.Fatal(err)
		}
		tables := map[string][]byte{}
		off := 0
		for _, e := range entries {
			if !e.xform {
				tables[e.tag] = data[off : off+e.len]
			}
			off += e.len
		}
		b = buildSFNT(tables, string(b[4:8]) == "OTTO")
	}
	sf, err := parseSFNT(b)
	if err != nil {
		t.Fatal(err)
	}
	return sf
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
	out := sf.rebuild(map[uint32]uint16{'A': sf.cmap['A'], 0xE9: sf.cmap[0xE9]}, fontInfo{family: "Pruned", weight: 400})
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

// nameString returns an English Windows record of a name table.
func nameString(sf *sfnt, id uint16) string {
	nt := sf.tables["name"]
	count, strOff := int(be16(nt, 2)), int(be16(nt, 4))
	for i := 0; i < count; i++ {
		rec := 6 + i*12
		if be16(nt, rec) == 3 && be16(nt, rec+6) == id {
			off, l := strOff+int(be16(nt, rec+10)), int(be16(nt, rec+8))
			u := make([]rune, l/2)
			for k := range u {
				u[k] = rune(be16(nt, off+k*2))
			}
			return string(u)
		}
	}
	return ""
}

// fontParts returns the font parts of a converted document keyed by family name.
func fontParts(t *testing.T, r *bdf.Reader) map[string]*sfnt {
	t.Helper()
	out := map[string]*sfnt{}
	for _, e := range r.Manifest.Parts {
		if e.T != bdf.PartFont {
			continue
		}
		b, err := r.Part(e.H)
		if err != nil {
			t.Fatal(err)
		}
		// WOFF2 is stored as is; without the encoder (bdf_noconv) the plain
		// font is deflated by the container instead.
		if woff2.IsWOFF(b) != woff2.Available() {
			t.Errorf("font part %s: WOFF2 %v, encoder available %v", e.H, woff2.IsWOFF(b), woff2.Available())
		}
		if woff2.IsWOFF(b) == (e.Enc == bdf.EncDeflateRaw) {
			t.Errorf("font part %s stored as %s", e.H, e.Enc)
		}
		sf := fontPart(t, b)
		out[nameString(sf, 1)] = sf
	}
	return out
}

func TestConvertWebFonts(t *testing.T) {
	res, r := convert(t, "weasyprint-webfonts.pdf")
	var restricted, noSubset bool
	for _, w := range res.Warnings {
		switch {
		case strings.Contains(w, "LicRestricted") && strings.Contains(w, "Restricted License"):
			restricted = true
		case strings.Contains(w, "LicNoSubset") && strings.Contains(w, "forbids subsetting"):
			noSubset = true
		default:
			t.Errorf("unexpected warning: %s", w)
		}
	}
	if !restricted || !noSubset {
		t.Fatalf("license warnings missing: %v", res.Warnings)
	}
	text := pageText(t, r)
	for _, want := range []string{"OpenType CFF (.otf): The quick brown fox", "Web font (.woff2): jumps over the lazy dog", "café naïve résumé", "Restricted License: system font", "Preview & Print: embedded"} {
		if !strings.Contains(text, want) {
			t.Errorf("text index lacks %q", want)
		}
	}
	fonts := fontParts(t, r)
	if len(fonts) != 4 || fonts["LicRestricted"] != nil {
		t.Fatalf("embedded fonts: %v", slices.Sorted(maps.Keys(fonts)))
	}
	// The CFF fonts (embedded as FontFile3/OpenType) come back as OpenType CFF
	// holding only the glyphs the page uses, with Loma's copyright notice.
	for _, name := range []string{"LomaOTF", "LomaWOFF2"} {
		sf := fonts[name]
		if sf == nil || !sf.isCFF {
			t.Fatalf("%s: not an OpenType CFF font", name)
		}
		cf, err := parseCFF(sf.tables["CFF "])
		if err != nil {
			t.Fatal(err)
		}
		if cf.numGlyphs < 20 || cf.numGlyphs > 50 || sf.numGlyphs != cf.numGlyphs {
			t.Errorf("%s: %d glyphs (maxp %d)", name, cf.numGlyphs, sf.numGlyphs)
		}
		if !strings.Contains(nameString(sf, 0), "NECTEC") || sf.fsType != 0 || !sf.hasFSType {
			t.Errorf("%s: notice %q, fsType %#x", name, nameString(sf, 0), sf.fsType)
		}
		for _, r := range map[string]string{"LomaOTF": "The quick brown fox café", "LomaWOFF2": "jumps over the lazy dog"}[name] {
			if _, ok := sf.cmap[uint32(r)]; !ok {
				t.Errorf("%s: no glyph for %q", name, r)
			}
		}
	}
	// fsType carries over; the No subsetting font keeps every glyph outline.
	if sf := fonts["LicPreview"]; sf == nil || sf.fsType != fsPreviewPrint {
		t.Errorf("LicPreview fsType not carried over")
	}
	if sf := fonts["LicNoSubset"]; sf == nil || sf.fsType != fsNoSubsetting {
		t.Errorf("LicNoSubset fsType not carried over")
	}
	if !strings.Contains(nameString(fonts["LicPreview"], 0), "Bitstream") {
		t.Errorf("DejaVu copyright not carried over: %q", nameString(fonts["LicPreview"], 0))
	}

	// With -ignore-fstype the Restricted License font is embedded too.
	res2, err := ConvertFile(filepath.Join("testdata", "weasyprint-webfonts.pdf"), &Options{IgnoreFSType: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res2.Warnings) != 0 {
		t.Errorf("warnings with IgnoreFSType: %v", res2.Warnings)
	}
	n := 0
	for _, p := range res2.Doc.Parts() {
		if p.Type == bdf.PartFont {
			n++
		}
	}
	if n != 5 {
		t.Errorf("IgnoreFSType: %d fonts embedded, want 5", n)
	}
}

func TestConvertCairoCFF(t *testing.T) {
	res, r := convert(t, "cairo-cff.pdf")
	if len(res.Warnings) != 0 {
		t.Fatalf("warnings: %v", res.Warnings)
	}
	text := pageText(t, r)
	for _, want := range []string{"Type1C via cairo: Hello, world", "CIDFontType0C:", "ภาษาไทย", "café — “quotes” ½"} {
		if !strings.Contains(text, want) {
			t.Errorf("text index lacks %q", want)
		}
	}
	cid := 0
	for name, sf := range fontParts(t, r) {
		if !sf.isCFF || sf.hasFSType && sf.fsType != fsPreviewPrint {
			t.Errorf("%s: CFF %v, fsType %#x", name, sf.isCFF, sf.fsType)
		}
		cf, err := parseCFF(sf.tables["CFF "])
		if err != nil {
			t.Fatal(err)
		}
		if cf.isCID {
			cid++
		}
		if !strings.Contains(nameString(sf, 0), "NECTEC") {
			t.Errorf("%s: CFF Notice not carried over", name)
		}
	}
	if cid != 1 {
		t.Errorf("want one CID-keyed font, got %d", cid)
	}
}

// cffPrograms returns the FontFile3 programs of a testdata PDF (the CFF table
// for OpenType ones) keyed by FontName.
func cffPrograms(t *testing.T, name string) map[string][]byte {
	t.Helper()
	f, err := os.Open(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed
	ctx, err := api.ReadContext(f, conf)
	if err != nil {
		t.Fatal(err)
	}
	p := &pdf{ctx: ctx}
	out := map[string][]byte{}
	for nr := range ctx.XRefTable.Table {
		d := p.dict(*types.NewIndirectRef(nr, 0))
		if d == nil || p.name(d["Type"]) != "FontDescriptor" {
			continue
		}
		if sd := p.stream(d["FontFile3"]); sd != nil {
			data, _, err := p.decodeStream(sd)
			if err != nil {
				t.Fatal(err)
			}
			if string(data[:4]) == "OTTO" {
				sf, err := parseSFNT(data)
				if err != nil {
					t.Fatal(err)
				}
				data = sf.tables["CFF "]
			}
			out[p.name(d["FontName"])] = data
		}
	}
	return out
}

// cffParts reads the CharStrings INDEX and the local Subrs INDEXes of a CFF program.
func cffParts(t *testing.T, data []byte) (charStrings [][]byte, subrs [][][]byte) {
	t.Helper()
	_, pos, _ := cffReadIndex(data, int(data[2]))
	tops, _, _ := cffReadIndex(data, pos)
	top := cffParseDict(tops[0])
	charStrings, _, err := cffReadIndex(data, int(top[17][0]))
	if err != nil {
		t.Fatal(err)
	}
	private := func(v []float64) {
		pd := cffParseDict(data[int(v[1]) : int(v[1])+int(v[0])])
		if s := pd[19]; len(s) == 1 {
			items, _, err := cffReadIndex(data, int(v[1])+int(s[0]))
			if err != nil {
				t.Fatal(err)
			}
			subrs = append(subrs, items)
		}
	}
	if v := top[1236]; len(v) == 1 {
		fds, _, _ := cffReadIndex(data, int(v[0]))
		for _, fd := range fds {
			private(cffParseDict(fd)[18])
		}
	} else {
		private(top[18])
	}
	return charStrings, subrs
}

func TestSubsetCFF(t *testing.T) {
	progs := cffPrograms(t, "cairo-cff.pdf")
	for k, v := range cffPrograms(t, "weasyprint-webfonts.pdf") {
		progs[k] = v
	}
	if len(progs) != 4 {
		t.Fatalf("CFF programs: %v", slices.Sorted(maps.Keys(progs)))
	}
	prunedSubrs := 0
	for name, data := range progs {
		cf, err := parseCFF(data)
		if err != nil {
			t.Fatal(err)
		}
		keep := map[int]bool{1: true, cf.numGlyphs / 2: true, cf.numGlyphs - 1: true}
		out, order, err := subsetCFF(data, keep)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if order[0] != 0 || !slices.IsSorted(order) || len(order) < len(keep)+1 {
			t.Fatalf("%s: order %v", name, order)
		}
		for g := range keep {
			if !slices.Contains(order, g) {
				t.Fatalf("%s: glyph %d dropped", name, g)
			}
		}
		nf, err := parseCFF(out)
		if err != nil || nf.numGlyphs != len(order) || nf.isCID != cf.isCID {
			t.Fatalf("%s: subset parses as %+v, %v", name, nf, err)
		}
		oldCS, oldSubrs := cffParts(t, data)
		newCS, newSubrs := cffParts(t, out)
		for i, g := range order {
			if !bytes.Equal(newCS[i], oldCS[g]) {
				t.Errorf("%s: charstring of glyph %d changed", name, g)
			}
			if i > 0 && nf.charset[i] != cf.charset[g] {
				t.Errorf("%s: glyph %d is now SID/CID %d, was %d", name, g, nf.charset[i], cf.charset[g])
			}
		}
		// Subroutine numbers are kept; the ones no kept glyph calls are emptied.
		if len(newSubrs) != len(oldSubrs) {
			t.Fatalf("%s: %d Subrs INDEXes, was %d", name, len(newSubrs), len(oldSubrs))
		}
		for k := range oldSubrs {
			if len(newSubrs[k]) != len(oldSubrs[k]) {
				t.Fatalf("%s: Subrs count changed", name)
			}
			for i, s := range newSubrs[k] {
				if !bytes.Equal(s, oldSubrs[k][i]) {
					if !bytes.Equal(s, []byte{11}) {
						t.Fatalf("%s: subr %d rewritten", name, i)
					}
					prunedSubrs++
				}
			}
		}
		if len(out) >= len(data) {
			t.Errorf("%s: %d -> %d bytes", name, len(data), len(out))
		}
		t.Logf("%s: %d glyphs, %d bytes -> %d glyphs, %d bytes", name, cf.numGlyphs, len(data), nf.numGlyphs, len(out))
	}
	if prunedSubrs == 0 {
		t.Error("no subroutine was pruned")
	}
}

func TestFontLicense(t *testing.T) {
	for _, c := range []struct {
		fsType                               uint16
		restricted, bitmapOnly, noSubsetting bool
	}{
		{0x0000, false, false, false},
		{0x0002, true, false, false},
		{0x0004, false, false, false},
		{0x0006, false, false, false}, // the least restrictive usage bit wins
		{0x000a, false, false, false},
		{0x0104, false, false, true},
		{0x0200, false, true, false},
		{0x0302, true, true, true},
	} {
		l := fontLicense{known: true, fsType: c.fsType}
		if l.restricted() != c.restricted || l.bitmapOnly() != c.bitmapOnly || l.noSubsetting() != c.noSubsetting || l.outFSType() != c.fsType {
			t.Errorf("fsType %#04x: restricted %v bitmapOnly %v noSubsetting %v", c.fsType, l.restricted(), l.bitmapOnly(), l.noSubsetting())
		}
	}
	if l := (fontLicense{}); l.restricted() || l.outFSType() != fsPreviewPrint {
		t.Error("a font without OS/2 must be embeddable and marked Preview & Print")
	}
}

func TestT2Exec(t *testing.T) {
	num := func(v int) byte { return byte(v + 139) } // -107..107
	local := newSubrSet([][]byte{{11}, {num(5), 11}, {11}})
	global := newSubrSet([][]byte{{11}})
	// Two stems, then a hintmask whose mask byte 0xff must not be read as an
	// operand; subr 1 (biased -106) is called, then seac 'A' + grave (code 193).
	cs := []byte{num(10), num(20), num(30), num(40), 1, 19, 0xff, num(-106), 10, num(0), num(0), num(65), 247, 193 - 108, 14}
	st := &t2State{seac: [2]int{-1, -1}}
	done, err := st.exec(cs, local, global, 0)
	if err != nil || !done {
		t.Fatalf("exec = %v, %v", done, err)
	}
	if !slices.Equal(local.used, []bool{false, true, false}) || global.used[0] {
		t.Fatalf("used = %v %v", local.used, global.used)
	}
	if st.nStems != 2 || st.seac != [2]int{65, 193} {
		t.Fatalf("stems %d, seac %v", st.nStems, st.seac)
	}
	// Arithmetic can compute subroutine numbers: give up rather than guess.
	st = &t2State{seac: [2]int{-1, -1}}
	if _, err := st.exec([]byte{num(1), num(2), 12, 10, 10, 14}, local, global, 0); err == nil {
		t.Fatal("arithmetic operator was followed")
	}
	// A subroutine number out of range is an error too.
	if _, err := (&t2State{}).exec([]byte{num(50), 10, 14}, local, global, 0); err == nil {
		t.Fatal("out-of-range subroutine was followed")
	}
}
