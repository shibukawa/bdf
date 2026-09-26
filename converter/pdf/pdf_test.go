package pdf

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	"github.com/shibukawa/bdf/converter/internal/sfnt"
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

func TestFunctions(t *testing.T) {
	f := &pdfFunction{kind: 2, domain: []float64{0, 1}, c0: []float64{0, 0, 1}, c1: []float64{1, 0, 0}, n: 1}
	if out := f.eval(0.25); out[0] != 0.25 || out[2] != 0.75 {
		t.Fatalf("exp = %v", out)
	}
	st := &pdfFunction{kind: 3, domain: []float64{0, 1}, funcs: []*pdfFunction{f, f}, bounds: []float64{0.5}, encode: []float64{0, 1, 1, 0}}
	if out := st.eval(0.75); out[0] != 0.5 {
		t.Fatalf("stitch = %v", out)
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
	if m := r.Manifest.Meta; m.DC.Title.First() != "slides.html" || m.DC.Language.First() != "en-US" {
		t.Errorf("meta = %+v", m)
	}
	// The page content is a form XObject reused as a shared object; links are emitted on page 2.
	// The form's marked content belongs to the page's structure.
	seq := markSeq(t, r, v.Pages[1].Layers[0].Obj)
	for _, want := range []string{"HEADING 1|text Images, clips and blending", "TABLE|CELL A1 col|text R1C1|CELL B1 col|text R1C2|CELL A2|text R2C1|CELL B2|text R2C2|fill|END|PARAGRAPH Div|text a link to example.com"} {
		if !strings.Contains(seq, want) {
			t.Errorf("page 2 lacks %q:\n%s", want, seq)
		}
	}
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
	if m := r.Manifest.Meta; m.DC.Title.First() != "doc.html" || m.DC.Language.First() != "en-US" {
		t.Errorf("meta = %+v", m)
	}
	// Chrome tags all text as NonStruct; headings, paragraphs, the list and the
	// table come from the structure tree. Paragraph breaks are the only change
	// to the text.
	want := []string{
		"BDF conversion fixture",
		"BDF is a display-list format for browsers. Each page is a list of objects, and each object is a stream of instructions that map directly onto the Canvas 2D API. Fonts and images are stored as the browser already understands them, so the decoder stays small.",
		"Repeated content such as headers and footers is stored once and referenced by hash from every page. The body of each page is its own object, laid out inside the body rectangle so that a viewer can stack bodies vertically for a continuous reading mode.",
		"This paragraph is a note with a colored rule. It contains emphasis, strong text, monospace code and a hyperlink.",
		"First bullet item",
		"Second bullet item with more words to wrap onto the next line when the column is narrow enough to force it",
		"Third", "Name", "Value", "Note", "alpha", "1.00", "first", "beta", "2.50", "second", "gamma", "10.25", "third",
		"Second page",
		"Text on the second page, followed by a justified paragraph.",
		"Justified text uses word-spacing adjustments in the PDF content stream, which the converter must reproduce by positioning each word so that the right edge lines up. This paragraph is long enough to fill several lines and thereby exercise that code path with a variety of word lengths and spacing amounts.",
	}
	if got := strings.Split(text, "\n"); !slices.Equal(got, want) {
		t.Errorf("text lines:\n%s", strings.Join(got, "\n"))
	}
	seq := markSeq(t, r, r.Manifest.Views[0].Pages[0].Layers[0].Obj)
	for _, want := range []string{
		"HEADING 1|text BDF conversion |ALT_TEXT fixture",
		"PARAGRAPH P|text BDF is a display-list",
		"LIST|LIST_ITEM|text First bullet item|fill|LIST_ITEM|text Second bullet",
		"fill|LIST_ITEM|text Third|END|TABLE|CELL A1 col|text Name|CELL B1 col|text Value|CELL C1 col|text Note|CELL A2|text alpha",
		"CELL C4|text third|END|link 242.25 542.67 53.25 12.75 https://example.com/doc",
	} {
		if !strings.Contains(seq, want) {
			t.Errorf("page 1 lacks %q:\n%s", want, seq)
		}
	}
	if seq := markSeq(t, r, r.Manifest.Views[0].Pages[1].Layers[0].Obj); !strings.HasPrefix(seq, "HEADING 2|text Second page|PARAGRAPH P|") {
		t.Errorf("page 2: %s", seq)
	}
	counts := opCounts(t, r, r.Manifest.Views[0].Pages[0].Layers[0].Obj)
	obj, _ := r.Object(r.Manifest.Views[0].Pages[0].Layers[0].Obj)
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
	// The information dictionary becomes Dublin Core.
	dc := r.Manifest.Meta.DC
	if dc.Title.First() != "reportlab fixture" || dc.Creator.First() != "anonymous" || dc.Description.First() != "unspecified" ||
		len(dc.Subject) != 0 || dc.Created.First() != "2026-09-25T07:21:16+00:00" || dc.Modified.First() != "2026-09-25T07:21:16+00:00" {
		t.Errorf("dc = %+v", dc)
	}
	if v := r.Manifest.Views[0]; v.Title != "reportlab fixture" {
		t.Errorf("view title = %q", v.Title)
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
			if len(sf.Cmap) == 0 {
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
func fontPart(t *testing.T, b []byte) *sfnt.Font {
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
		b = sfnt.Build(tables, string(b[4:8]) == "OTTO")
	}
	sf, err := sfnt.Parse(b)
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

// nameString returns an English Windows record of a name table.
func nameString(sf *sfnt.Font, id uint16) string {
	nt := sf.Tables["name"]
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
func fontParts(t *testing.T, r *bdf.Reader) map[string]*sfnt.Font {
	t.Helper()
	out := map[string]*sfnt.Font{}
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
		if sf == nil || !sf.IsCFF {
			t.Fatalf("%s: not an OpenType CFF font", name)
		}
		cf, err := parseCFF(sf.Tables["CFF "])
		if err != nil {
			t.Fatal(err)
		}
		if cf.numGlyphs < 20 || cf.numGlyphs > 50 || sf.NumGlyphs != cf.numGlyphs {
			t.Errorf("%s: %d glyphs (maxp %d)", name, cf.numGlyphs, sf.NumGlyphs)
		}
		if !strings.Contains(nameString(sf, 0), "NECTEC") || sf.FSType != 0 || !sf.HasFSType {
			t.Errorf("%s: notice %q, fsType %#x", name, nameString(sf, 0), sf.FSType)
		}
		for _, r := range map[string]string{"LomaOTF": "The quick brown fox café", "LomaWOFF2": "jumps over the lazy dog"}[name] {
			if _, ok := sf.Cmap[uint32(r)]; !ok {
				t.Errorf("%s: no glyph for %q", name, r)
			}
		}
	}
	// fsType carries over; the No subsetting font keeps every glyph outline.
	if sf := fonts["LicPreview"]; sf == nil || sf.FSType != fsPreviewPrint {
		t.Errorf("LicPreview fsType not carried over")
	}
	if sf := fonts["LicNoSubset"]; sf == nil || sf.FSType != fsNoSubsetting {
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
		if !sf.IsCFF || sf.HasFSType && sf.FSType != fsPreviewPrint {
			t.Errorf("%s: CFF %v, fsType %#x", name, sf.IsCFF, sf.FSType)
		}
		cf, err := parseCFF(sf.Tables["CFF "])
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
				sf, err := sfnt.Parse(data)
				if err != nil {
					t.Fatal(err)
				}
				data = sf.Tables["CFF "]
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

func TestConvertType1(t *testing.T) {
	for _, c := range []struct {
		file  string
		fonts int
		text  []string
	}{
		{"cairo-type1.pdf", 3, []string{"Type1 via cairo: Hello, world", "café naïve Ångström — fi", "Courier 10 Pitch: mono"}},
		{"reportlab-type1.pdf", 1, []string{"Whole Type1 font: café naïve Ærø", "0123456789 (reportlab embeds every glyph)"}},
	} {
		res, r := convert(t, c.file)
		if len(res.Warnings) != 0 {
			t.Fatalf("%s: warnings: %v", c.file, res.Warnings)
		}
		text := pageText(t, r)
		for _, want := range c.text {
			if !strings.Contains(text, want) {
				t.Errorf("%s: text index lacks %q in %q", c.file, want, text)
			}
		}
		n := 0
		for _, e := range r.Manifest.Parts {
			if e.T == bdf.PartFont {
				n++
			}
		}
		if n != c.fonts {
			t.Fatalf("%s: %d fonts embedded, want %d", c.file, n, c.fonts)
		}
		fonts := fontParts(t, r)
		for name, sf := range fonts {
			cf, err := parseCFF(sf.Tables["CFF "])
			if err != nil || cf.isCID {
				t.Fatalf("%s %s: %v", c.file, name, err)
			}
			// Only the glyphs the page uses (plus .notdef and seac parts) are kept.
			if cf.numGlyphs > 45 {
				t.Errorf("%s %s: %d glyphs", c.file, name, cf.numGlyphs)
			}
			if !strings.Contains(nameString(sf, 0)+nameString(sf, 7), "Bitstream") || sf.FSType != fsPreviewPrint {
				t.Errorf("%s %s: notice %q, fsType %#x", c.file, name, nameString(sf, 0), sf.FSType)
			}
		}
	}
}

// type1Programs returns the FontFile programs of a testdata PDF.
func type1Programs(t *testing.T, name string) [][]byte {
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
	var out [][]byte
	for nr := range ctx.XRefTable.Table {
		d := p.dict(*types.NewIndirectRef(nr, 0))
		if d == nil || p.name(d["Type"]) != "FontDescriptor" {
			continue
		}
		if sd := p.stream(d["FontFile"]); sd != nil {
			data, _, err := p.decodeStream(sd)
			if err != nil {
				t.Fatal(err)
			}
			out = append(out, data)
		}
	}
	return out
}

func TestType1Forms(t *testing.T) {
	progs := type1Programs(t, "reportlab-type1.pdf")
	if len(progs) != 1 {
		t.Fatalf("%d Type1 programs", len(progs))
	}
	data := progs[0]
	want, cf, failed, err := type1ToCFF(data)
	if err != nil || failed != 0 {
		t.Fatalf("convert: %v, %d failed", err, failed)
	}
	if cf.numGlyphs < 200 || cf.notice == "" {
		t.Fatalf("%d glyphs, notice %q", cf.numGlyphs, cf.notice)
	}
	// PFB segments and hex (PFA) eexec data must give the same font.
	i := bytes.Index(data, []byte("eexec")) + len("eexec")
	for data[i] == '\r' || data[i] == '\n' {
		i++
	}
	clear, enc := data[:i], data[i:]
	seg := func(typ byte, b []byte) []byte {
		return append(binary.LittleEndian.AppendUint32([]byte{0x80, typ}, uint32(len(b))), b...)
	}
	pfb := append(append(seg(1, clear), seg(2, enc)...), 0x80, 3)
	var hex []byte
	for k, b := range enc {
		hex = fmt.Appendf(hex, "%02x", b)
		if k%32 == 31 {
			hex = append(hex, '\n')
		}
	}
	pfa := append(append([]byte(nil), clear...), hex...)
	for form, b := range map[string][]byte{"pfb": pfb, "pfa": pfa} {
		got, _, _, err := type1ToCFF(b)
		if err != nil || !bytes.Equal(got, want) {
			t.Errorf("%s: %v, same font %v", form, err, bytes.Equal(got, want))
		}
	}
	if _, _, _, err := type1ToCFF(data[:i-20]); err == nil {
		t.Error("a program cut before eexec was accepted")
	}
}

// t1cs encodes a Type 1 charstring: ints are operands, strings are operators.
func t1cs(items ...any) []byte {
	ops := map[string][]byte{"hstem": {1}, "rmoveto": {21}, "rlineto": {5}, "hsbw": {13}, "endchar": {14}, "callsubr": {10}, "return": {11},
		"closepath": {9}, "seac": {12, 6}, "div": {12, 12}, "callothersubr": {12, 16}, "pop": {12, 17}, "setcurrentpoint": {12, 33}}
	var b []byte
	for _, it := range items {
		switch v := it.(type) {
		case int:
			switch {
			case v >= -107 && v <= 107:
				b = append(b, byte(v+139))
			case v >= 108 && v <= 1131:
				b = append(b, byte((v-108)>>8+247), byte(v-108))
			case v <= -108 && v >= -1131:
				b = append(b, byte((-v-108)>>8+251), byte(-v-108))
			default:
				b = binary.BigEndian.AppendUint32(append(b, 255), uint32(int32(v)))
			}
		case string:
			b = append(b, ops[v]...)
		}
	}
	return b
}

// t2Outline decodes the Type 2 charstrings the converter writes (rmoveto,
// rlineto, rrcurveto, endchar) into a width and absolute points per contour.
func t2Outline(t *testing.T, cs []byte) (float64, [][][2]float64) {
	t.Helper()
	var stack []float64
	var contours [][][2]float64
	x, y, width := 0.0, 0.0, 0.0
	first := true
	for i := 0; i < len(cs); {
		v := cs[i]
		switch {
		case v >= 32 && v <= 246:
			stack = append(stack, float64(int(v)-139))
			i++
			continue
		case v >= 247 && v <= 250:
			stack = append(stack, float64((int(v)-247)*256+int(cs[i+1])+108))
			i += 2
			continue
		case v >= 251 && v <= 254:
			stack = append(stack, float64(-(int(v)-251)*256-int(cs[i+1])-108))
			i += 2
			continue
		case v == 28:
			stack = append(stack, float64(int16(binary.BigEndian.Uint16(cs[i+1:]))))
			i += 3
			continue
		case v == 255:
			stack = append(stack, float64(int32(binary.BigEndian.Uint32(cs[i+1:])))/65536)
			i += 5
			continue
		}
		i++
		n := map[byte]int{21: 2, 5: 2, 8: 6, 14: 0}[v]
		if first && (v == 21 || v == 14) && len(stack)%max(n, 1) == 1 || first && v == 14 && len(stack) == 1 {
			width, stack = stack[0], stack[1:]
		}
		first = false
		switch v {
		case 21:
			x, y = x+stack[0], y+stack[1]
			contours = append(contours, [][2]float64{{x, y}})
		case 5, 8:
			for k := 0; k+1 < len(stack); k += 2 {
				x, y = x+stack[k], y+stack[k+1]
				contours[len(contours)-1] = append(contours[len(contours)-1], [2]float64{x, y})
			}
		case 14:
			return width, contours
		default:
			t.Fatalf("unexpected operator %d", v)
		}
		stack = stack[:0]
	}
	t.Fatal("no endchar")
	return 0, nil
}

func TestType1Charstrings(t *testing.T) {
	font := &type1Font{
		subrs: [][]byte{
			t1cs(3, 0, "callothersubr", "pop", "pop", "setcurrentpoint", "return"), // 0: end flex
			t1cs(0, 1, "callothersubr", "return"),                                  // 1: start flex
			t1cs(0, 2, "callothersubr", "return"),                                  // 2: flex point
			t1cs("return"),                                                         // 3
			t1cs(0, 50, "hstem", "return"),                                         // 4: replacement hints
		},
		charString: map[string][]byte{
			// Flex from (100,0) to (400,0), a hint replacement, then a line up.
			"flex": t1cs(0, 500, "hsbw", 100, 0, "rmoveto", 1, "callsubr",
				150, 0, "rmoveto", 2, "callsubr", -100, 20, "rmoveto", 2, "callsubr", 50, 0, "rmoveto", 2, "callsubr",
				50, 0, "rmoveto", 2, "callsubr", 50, 0, "rmoveto", 2, "callsubr", 50, 0, "rmoveto", 2, "callsubr",
				50, -20, "rmoveto", 2, "callsubr", 50, 400, 0, 0, "callsubr",
				4, 1, 3, "callothersubr", "pop", "callsubr", 0, 100, "rlineto", "closepath", "endchar"),
			"A":     t1cs(20, 600, "hsbw", 0, 0, "rmoveto", 100, 0, "rlineto", "closepath", "endchar"),
			"grave": t1cs(5, 300, "hsbw", 0, 500, "rmoveto", 50, 0, "rlineto", "closepath", "endchar"),
			// seac: the accent's sidebearing point goes to sbx + adx - asb + its own sbx.
			"Agrave": t1cs(20, 600, "hsbw", 5, 150, 200, 65, 193, "seac"),
			// div and a line without a preceding moveto (starts at the sidebearing point).
			"div": t1cs(10, 1200, "hsbw", 7, 2, "div", 0, "rlineto", "endchar"),
		},
	}
	for name, want := range map[string][][][2]float64{
		"flex":   {{{100, 0}, {150, 20}, {200, 20}, {250, 20}, {300, 20}, {350, 20}, {400, 0}, {400, 100}}},
		"Agrave": {{{20, 0}, {120, 0}}, {{170, 700}, {220, 700}}},
		"div":    {{{10, 0}, {13.5, 0}}},
	} {
		cs, err := font.convert(name)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		width, got := t2Outline(t, cs)
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("%s: outline %v, want %v", name, got, want)
		}
		if w := map[string]float64{"flex": 500, "Agrave": 600, "div": 1200}[name]; width != w {
			t.Errorf("%s: width %v, want %v", name, width, w)
		}
	}
	// Multiple master blends and a seac whose base is a seac are not converted.
	font.charString["mm"] = t1cs(0, 500, "hsbw", 1, 2, 1, 14, "callothersubr", "endchar")
	font.charString["A"] = t1cs(0, 500, "hsbw", 0, 0, 0, 193, 193, "seac")
	font.charString["nested"] = t1cs(0, 500, "hsbw", 0, 0, 0, 65, 193, "seac")
	for _, name := range []string{"mm", "nested"} {
		if _, err := font.convert(name); err == nil {
			t.Errorf("%s was converted", name)
		}
	}
}

var markNames = map[byte]string{bdf.MarkParagraph: "PARAGRAPH", bdf.MarkLine: "LINE", bdf.MarkCell: "CELL", bdf.MarkBox: "BOX",
	bdf.MarkAltText: "ALT_TEXT", bdf.MarkWrap: "WRAP", bdf.MarkHeading: "HEADING", bdf.MarkList: "LIST", bdf.MarkListItem: "LIST_ITEM",
	bdf.MarkTable: "TABLE", bdf.MarkFigure: "FIGURE", bdf.MarkEnd: "END", bdf.MarkLang: "LANG"}

// markSeq lists the MARKs, text runs, path drawings, images and links of an
// object and the objects it USEs in walk order, joined by "|".
func markSeq(t *testing.T, r *bdf.Reader, h bdf.Hash) string {
	t.Helper()
	var out []string
	var walk func(h bdf.Hash)
	walk = func(h bdf.Hash) {
		o, err := r.Object(h)
		if err != nil {
			t.Fatal(err)
		}
		err = o.Walk(func(in bdf.Instr) {
			switch in.Op {
			case bdf.OpMark:
				out = append(out, strings.TrimSpace(markNames[byte(in.Args[0].(uint64))]+" "+in.Args[1].(string)))
			case bdf.OpFillText, bdf.OpStrokeText:
				out = append(out, "text "+in.Args[0].(string))
			case bdf.OpFillPath, bdf.OpFillRect:
				out = append(out, "fill")
			case bdf.OpStrokePath, bdf.OpStrokeRect:
				out = append(out, "stroke")
			case bdf.OpImage:
				out = append(out, "image")
			case bdf.OpLink:
				out = append(out, fmt.Sprintf("link %v %v %v %v %s", in.Args[0], in.Args[1], in.Args[2], in.Args[3], in.Args[4]))
			case bdf.OpUse, bdf.OpUseAt:
				walk(o.Objects[int(in.Args[0].(uint64))])
			}
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	walk(h)
	return strings.Join(out, "|")
}

// buildPDF assembles a PDF from objects numbered from 1; a [2]string is a
// stream (dictionary entries, data).
func buildPDF(objs []any, trailer string) []byte {
	var b bytes.Buffer
	b.WriteString("%PDF-1.7\n")
	offsets := make([]int, len(objs))
	for i, o := range objs {
		offsets[i] = b.Len()
		fmt.Fprintf(&b, "%d 0 obj\n", i+1)
		switch v := o.(type) {
		case string:
			b.WriteString(v)
		case [2]string:
			fmt.Fprintf(&b, "<< %s /Length %d >>\nstream\n%s\nendstream", v[0], len(v[1]), v[1])
		}
		b.WriteString("\nendobj\n")
	}
	xref := b.Len()
	fmt.Fprintf(&b, "xref\n0 %d\n0000000000 65535 f \n", len(objs)+1)
	for _, off := range offsets {
		fmt.Fprintf(&b, "%010d 00000 n \n", off)
	}
	fmt.Fprintf(&b, "trailer\n<< /Size %d %s >>\nstartxref\n%d\n%%%%EOF\n", len(objs)+1, trailer, xref)
	return b.Bytes()
}

// taggedPDF is a two-page tagged PDF with the structure types the converter
// maps, a role-mapped heading, languages, a figure and links. Without a
// parent tree, MCIDs are found by walking the structure tree.
func taggedPDF(parentTree bool) []byte {
	root := `<< /Type /StructTreeRoot /K 12 0 R /ParentTree 13 0 R /RoleMap << /MyHead /H2 >> >>`
	if !parentTree {
		root = `<< /Type /StructTreeRoot /K 12 0 R /RoleMap << /MyHead /H2 >> >>`
	}
	page1 := `/Artifact BMC BT /F1 10 Tf 72 770 Td (header) Tj ET EMC
/H1 <</MCID 0>> BDC BT /F1 24 Tf 72 720 Td /Span <</ActualText (T\\itle)>> BDC (Title) Tj EMC ET EMC
/P <</MCID 1>> BDC BT /F1 12 Tf 72 690 Td (Para one) Tj ET EMC
/LBody <</MCID 2>> BDC BT /F1 12 Tf 72 670 Td (Item one) Tj ET EMC
/LBody <</MCID 3>> BDC BT /F1 12 Tf 90 655 Td (Nested) Tj ET EMC
/Lbl <</MCID 12>> BDC BT /F1 12 Tf 72 640 Td (2.) Tj ET EMC
/P <</MCID 4>> BDC BT /F1 12 Tf 90 640 Td (Item two) Tj ET EMC
/TH <</MCID 5>> BDC 0.9 g 72 605 150 15 re f BT 0 g /F1 12 Tf 72 610 Td (Head) Tj ET EMC
/TH <</MCID 6>> BDC BT /F1 12 Tf 72 595 Td (Row) Tj ET EMC
/TD <</MCID 7>> BDC BT /F1 12 Tf 150 595 Td (Cell) Tj ET EMC
/Figure <</MCID 8>> BDC 0 0 1 rg 72 500 50 50 re f EMC
/Artifact BMC 0 G 72 480 m 300 480 l S EMC
/P <</MCID 9>> BDC BT 0 g /F1 12 Tf 72 450 Td (Mixed) Tj ET EMC
/Span <</MCID 10>> BDC BT /F1 12 Tf 110 450 Td (english) Tj ET EMC
/P <</MCID 11>> BDC BT /F1 12 Tf 160 450 Td (tail ) Tj /Span /Pr1 BDC (deutsch) Tj EMC ET EMC`
	page2 := `/MyHead <</MCID 0>> BDC BT /F1 18 Tf 72 720 Td (Chapter 2) Tj ET EMC
/P <</MCID 1>> BDC q 1 0 0 1 72 690 cm /Fm1 Do Q EMC
/P <</MCID 2>> BDC q 1 0 0 1 72 670 cm /Fm1 Do Q EMC`
	elem := func(s, parent, rest string) string {
		return "<< /Type /StructElem /S /" + s + " /P " + parent + " 0 R " + rest + " >>"
	}
	return buildPDF([]any{
		/* 1 */ `<< /Type /Catalog /Pages 2 0 R /Lang (ja) /MarkInfo << /Marked true >> /StructTreeRoot 11 0 R
			/Names << /Dests 9 0 R >> /Dests << /Legacy [3 0 R /Fit] >> >>`,
		/* 2 */ `<< /Type /Pages /Kids [3 0 R 4 0 R] /Count 2 >>`,
		/* 3 */ `<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 5 0 R /StructParents 0
			/Resources << /Font << /F1 7 0 R >> /Properties << /Pr1 << /Lang (de) >> >> >> /Annots [37 0 R 38 0 R 39 0 R] >>`,
		/* 4 */ `<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 6 0 R /StructParents 1
			/Resources << /Font << /F1 7 0 R >> /XObject << /Fm1 40 0 R >> >> >>`,
		/* 5 */ [2]string{"", page1},
		/* 6 */ [2]string{"", page2},
		/* 7 */ `<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>`,
		/* 8 */ `<< /Title (Tagged \(test\)) >>`,
		/* 9 */ `<< /Kids [10 0 R] >>`,
		/* 10 */ `<< /Limits [(chap2) (chap2)] /Names [(chap2) << /D [4 0 R /Fit] >>] >>`,
		/* 11 */ root,
		/* 12 */ `<< /Type /StructElem /S /Document /P 11 0 R /K [15 0 R 16 0 R 17 0 R 25 0 R 31 0 R 32 0 R 34 0 R 35 0 R 36 0 R] >>`,
		/* 13 */ `<< /Kids [14 0 R] >>`,
		/* 14 */ `<< /Limits [0 1] /Nums [0 [15 0 R 16 0 R 19 0 R 22 0 R 42 0 R 27 0 R 29 0 R 30 0 R 31 0 R 32 0 R 33 0 R 32 0 R 41 0 R] 1 [34 0 R 35 0 R 36 0 R]] >>`,
		/* 15 */ elem("H1", "12", "/K 0 /Pg 3 0 R"),
		/* 16 */ elem("P", "12", "/K 1 /Pg 3 0 R"),
		/* 17 */ elem("L", "12", "/K [18 0 R 23 0 R]"),
		/* 18 */ elem("LI", "17", "/K 19 0 R"),
		/* 19 */ elem("LBody", "18", "/K [2 20 0 R] /Pg 3 0 R"),
		/* 20 */ elem("L", "19", "/K 21 0 R"),
		/* 21 */ elem("LI", "20", "/K 22 0 R"),
		/* 22 */ elem("LBody", "21", "/K 3 /Pg 3 0 R"),
		/* 23 */ elem("LI", "17", "/K [41 0 R 24 0 R]"),
		/* 24 */ elem("LBody", "23", "/K 42 0 R"),
		/* 25 */ elem("Table", "12", "/K [26 0 R 28 0 R]"),
		/* 26 */ elem("TR", "25", "/K 27 0 R"),
		/* 27 */ elem("TH", "26", "/K 5 /Pg 3 0 R /A << /O /Table /ColSpan 2 >>"),
		/* 28 */ elem("TR", "25", "/K [29 0 R 30 0 R]"),
		/* 29 */ elem("TH", "28", "/K 6 /Pg 3 0 R /A [<< /O /Table /Scope /Row >> 0]"),
		/* 30 */ elem("TD", "28", "/K 7 /Pg 3 0 R"),
		/* 31 */ elem("Figure", "12", "/K 8 /Pg 3 0 R /Alt (A blue square)"),
		/* 32 */ elem("P", "12", "/K [9 33 0 R 11] /Pg 3 0 R"),
		/* 33 */ elem("Span", "32", "/K 10 /Pg 3 0 R /Lang (en)"),
		/* 34 */ elem("MyHead", "12", "/K 0 /Pg 4 0 R"),
		/* 35 */ elem("P", "12", "/K 1 /Pg 4 0 R"),
		/* 36 */ elem("P", "12", "/K 2 /Pg 4 0 R"),
		/* 37 */ `<< /Type /Annot /Subtype /Link /Rect [72 100 172 120] /A << /S /GoTo /D [4 0 R /XYZ 0 792 0] >> >>`,
		/* 38 */ `<< /Type /Annot /Subtype /Link /Rect [200 100 300 120] /Dest (chap2) >>`,
		/* 39 */ `<< /Type /Annot /Subtype /Link /Rect [320 100 420 120] /Dest /Legacy >>`,
		/* 40 */ [2]string{"/Type /XObject /Subtype /Form /BBox [0 0 200 20] /Resources << /Font << /F1 7 0 R >> >>", "BT /F1 12 Tf 0 5 Td (inform) Tj ET"},
		/* 41 */ elem("Lbl", "23", "/K 12 /Pg 3 0 R"),
		/* 42 */ elem("P", "24", "/K 4 /Pg 3 0 R"),
	}, "/Root 1 0 R /Info 8 0 R")
}

func convertBytes(t *testing.T, data []byte, opts *Options) (*Result, *bdf.Reader) {
	t.Helper()
	res, err := Convert(bytes.NewReader(data), opts)
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

func TestConvertTagged(t *testing.T) {
	res, r := convertBytes(t, taggedPDF(true), nil)
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	if m := r.Manifest.Meta; m.DC.Title.First() != "Tagged (test)" || m.DC.Language.First() != "ja" {
		t.Errorf("meta = %+v", m)
	}
	pages := r.Manifest.Views[0].Pages
	want1 := strings.Join([]string{
		"text header",
		"HEADING 1", `ALT_TEXT T\itle`, "text Title",
		"PARAGRAPH P", "text Para one",
		"LIST", "LIST_ITEM", "text Item one",
		"LIST", "LIST_ITEM", "text Nested",
		// The item's label and first paragraph stay together.
		"END", "LIST_ITEM", "text 2.", "text Item two",
		// The TH background draws no figure: its MARKs wait for the text.
		"fill", "END", "TABLE", "CELL A1:B1 col", "text Head",
		"CELL A2 row", "text Row",
		"CELL B2", "text Cell",
		"END", "FIGURE A blue square", "fill",
		"END", "stroke", // the artifact line is not part of the figure
		"PARAGRAPH P", "text Mixed",
		"LANG en", "text english",
		"LANG", "text tail ",
		"LANG de", "text deutsch",
		"link 72 100 100 20 #page=2", "link 200 100 100 20 #page=2", "link 320 100 100 20 #page=1",
	}, "|")
	if got := markSeq(t, r, pages[0].Layers[0].Obj); got != want1 {
		t.Errorf("page 1:\n got %s\nwant %s", got, want1)
	}
	// The form is drawn in two paragraphs: each use gets its own MARKs.
	want2 := "HEADING 2|text Chapter 2|PARAGRAPH P|text inform|PARAGRAPH P|text inform"
	if got := markSeq(t, r, pages[1].Layers[0].Obj); got != want2 {
		t.Errorf("page 2:\n got %s\nwant %s", got, want2)
	}
	_, r2 := convertBytes(t, taggedPDF(false), nil)
	for i, pg := range r2.Manifest.Views[0].Pages {
		if got, want := markSeq(t, r2, pg.Layers[0].Obj), markSeq(t, r, pages[i].Layers[0].Obj); got != want {
			t.Errorf("page %d without parent tree:\n got %s\nwant %s", i+1, got, want)
		}
	}
	lines := strings.Split(pageText(t, r), "\n")
	if want := []string{"header", `T\itle`, "Para one", "Item one", "Nested", "2. Item two", "Head", "Row", "Cell", "Mixed english tail deutsch", "Chapter 2", "inform", "inform"}; !slices.Equal(lines, want) {
		t.Errorf("text = %q", lines)
	}

	// Links point at page indexes of the output view; pages left out drop theirs.
	_, r = convertBytes(t, taggedPDF(true), &Options{Pages: []int{2, 1}})
	if got := markSeq(t, r, r.Manifest.Views[0].Pages[1].Layers[0].Obj); !strings.HasSuffix(got, "link 72 100 100 20 #page=1|link 200 100 100 20 #page=1|link 320 100 100 20 #page=2") {
		t.Errorf("reordered pages: %s", got)
	}
	_, r = convertBytes(t, taggedPDF(true), &Options{Pages: []int{1}})
	if got := markSeq(t, r, r.Manifest.Views[0].Pages[0].Layers[0].Obj); !strings.HasSuffix(got, "LANG de|text deutsch|link 320 100 100 20 #page=1") {
		t.Errorf("single page: %s", got)
	}
}

func TestConvertUntaggedMarks(t *testing.T) {
	// Without a structure tree, block tags still give paragraphs, after the
	// run before them; marked content /Lang gives the language of its runs.
	content := `BT /F1 12 Tf 72 700 Td (one ) Tj /P BMC (two) Tj EMC ET
BT /F1 12 Tf 72 680 Td /Span <</Lang (fr) /ActualText (tr\\ois)>> BDC (trois) Tj EMC ( four) Tj ET`
	data := buildPDF([]any{
		`<< /Type /Catalog /Pages 2 0 R >>`,
		`<< /Type /Pages /Kids [3 0 R] /Count 1 >>`,
		`<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >>`,
		[2]string{"", content},
		`<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>`,
	}, "/Root 1 0 R")
	_, r := convertBytes(t, data, nil)
	if m := r.Manifest.Meta; m.DC.Title.First() != "" || m.DC.Language.First() != "" {
		t.Errorf("meta = %+v", m)
	}
	want := `text one |PARAGRAPH P|text two|LANG fr|ALT_TEXT tr\ois|text trois|LANG|text  four`
	if got := markSeq(t, r, r.Manifest.Views[0].Pages[0].Layers[0].Obj); got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
}
