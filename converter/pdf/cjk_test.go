package pdf

import (
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// cjkPDF shows text in non-embedded CID fonts through predefined CMaps and
// without ToUnicode, as older Japanese PDFs do.
func cjkPDF(content string) []byte {
	cidFont := func(name, ordering string) string {
		return `<< /Type /Font /Subtype /CIDFontType0 /BaseFont /` + name + ` /CIDSystemInfo << /Registry (Adobe) /Ordering (` + ordering + `) /Supplement 2 >>
			/FontDescriptor << /Type /FontDescriptor /FontName /` + name + ` /Flags 6 /FontBBox [0 -141 1000 859] /ItalicAngle 0 /Ascent 859 /Descent -141 /CapHeight 709 /StemV 69 >>
			/DW 1000 /W [231 [500]] >>`
	}
	return buildPDF([]any{
		/* 1 */ `<< /Type /Catalog /Pages 2 0 R >>`,
		/* 2 */ `<< /Type /Pages /Kids [3 0 R] /Count 1 >>`,
		/* 3 */ `<< /Type /Page /Parent 2 0 R /MediaBox [0 0 400 400] /Contents 4 0 R
			/Resources << /Font << /F1 5 0 R /F2 6 0 R /F3 7 0 R /F4 8 0 R >> >> >>`,
		/* 4 */ [2]string{"", content},
		/* 5 */ `<< /Type /Font /Subtype /Type0 /BaseFont /Ryumin-Light-90ms-RKSJ-H /Encoding /90ms-RKSJ-H /DescendantFonts [9 0 R] >>`,
		/* 6 */ `<< /Type /Font /Subtype /Type0 /BaseFont /GothicBBB-Medium-UniJIS-UCS2-H /Encoding /UniJIS-UCS2-H /DescendantFonts [10 0 R] >>`,
		/* 7 */ `<< /Type /Font /Subtype /Type0 /BaseFont /STSong-Light-UniGB-UTF16-H /Encoding /UniGB-UTF16-H /DescendantFonts [11 0 R] >>`,
		/* 8 */ `<< /Type /Font /Subtype /Type0 /BaseFont /Ryumin-Light-Identity-H /Encoding /Identity-H /DescendantFonts [9 0 R] >>`,
		/* 9 */ cidFont("Ryumin-Light", "Japan1"),
		/* 10 */ cidFont("GothicBBB-Medium", "Japan1"),
		/* 11 */ cidFont("STSong-Light", "GB1"),
	}, "/Root 1 0 R")
}

func TestPredefinedCMaps(t *testing.T) {
	// あいう in Shift_JIS with a half-width "A"; 日本 in UCS-2; 中文 in
	// UTF-16; 亜 by CID with Identity-H.
	content := `BT /F1 20 Tf 20 350 Td <82a082a282a441> Tj ET
BT /F2 20 Tf 20 300 Td <65e5672c> Tj ET
BT /F3 20 Tf 20 250 Td <4e2d6587> Tj ET
BT /F4 20 Tf 20 200 Td <0465> Tj ET`
	res, r := convertBytes(t, cjkPDF(content), nil)
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	txt := pageText(t, r)
	for _, want := range []string{"あいうA", "日本", "中文", "亜"} {
		if !strings.Contains(txt, want) {
			t.Errorf("text %q lacks %q", txt, want)
		}
	}
	// The fonts are drawn by the viewer's CJK fonts: mincho → serif.
	o, err := r.Object(r.Manifest.Views[0].Pages[0].Layers[0].Obj)
	if err != nil {
		t.Fatal(err)
	}
	fams := map[string]bool{}
	for _, f := range o.Fonts {
		if f.Kind == bdf.FontSystem {
			fams[f.Family] = true
		}
	}
	if !fams["serif"] || !fams["sans-serif"] {
		t.Errorf("families %v", fams)
	}
	// Widths come from W by CID: the half-width A (CID 264) is not in W, so 1 em.
	var advances []float32
	if err := o.Walk(func(in bdf.Instr) {
		if in.Op == bdf.OpFillText {
			advances = append(advances, in.Args[3].(float32))
		}
	}); err != nil {
		t.Fatal(err)
	}
	if len(advances) == 0 || advances[0] != 80 {
		t.Errorf("advances %v", advances)
	}
}

func TestConvertCJKFixture(t *testing.T) {
	res, r := convert(t, "cjk-cmaps.pdf")
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	txt := pageText(t, r)
	for _, want := range []string{"明朝体の日本語 Shift_JIS", "ゴシック体：UCS-2 の文字列", "简体中文的文本", "繁體中文的文字", "한국어 텍스트",
		"縦書き「かぎ括弧」、ー。", "縦組み（丸括弧）ー", "たてがきのテキスト"} {
		if !strings.Contains(txt, want) {
			t.Errorf("text lacks %q:\n%s", want, txt)
		}
	}
	// A vertical line is one ALT_TEXT run on a USE_AT of its child, whose
	// glyphs a system font draws in their vertical forms.
	page, err := r.Object(r.Manifest.Views[0].Pages[0].Layers[0].Obj)
	if err != nil {
		t.Fatal(err)
	}
	var alts []string
	var draws []string
	if err := page.Walk(func(in bdf.Instr) {
		if in.Op == bdf.OpMark && in.Args[0].(uint64) == uint64(bdf.MarkAltText) {
			alts = append(alts, in.Args[1].(string))
		}
		if in.Op == bdf.OpUseAt {
			child, err := r.Object(page.Objects[in.Args[0].(uint64)])
			if err != nil {
				t.Fatal(err)
			}
			var sb strings.Builder
			_ = child.Walk(func(ci bdf.Instr) {
				if ci.Op == bdf.OpFillText {
					sb.WriteString(ci.Args[0].(string))
				}
			})
			draws = append(draws, sb.String())
		}
	}); err != nil {
		t.Fatal(err)
	}
	if len(alts) != 3 || alts[0] != "縦書き「かぎ括弧」、ー。" {
		t.Errorf("ALT_TEXT %q", alts)
	}
	if len(draws) != 3 || draws[0] != "縦書き﹁かぎ括弧﹂︑ー︒" || draws[1] != "縦組み︵丸括弧︶ー" {
		t.Errorf("drawn %q", draws)
	}
}

func TestEmbeddedCMapParent(t *testing.T) {
	// An embedded CMap over a predefined one: its own entries win, the rest
	// and the code space come from the parent.
	cm := parseCMap([]byte(`/CIDInit /ProcSet findresource begin 12 dict begin begincmap
/CMapName /Custom-RKSJ def /90ms-RKSJ-H usecmap
1 begincidchar <82a0> 9999 endcidchar endcmap end end`))
	if cm.cid(0x82a0) != 9999 || cm.cid(0x82a2) != 845 || cm.cid(0x41) != 264 {
		t.Errorf("cids %d %d %d", cm.cid(0x82a0), cm.cid(0x82a2), cm.cid(0x41))
	}
	if got := cm.decode([]byte{0x41, 0x82, 0xa2}); len(got) != 2 || got[0].nbytes != 1 || got[1].nbytes != 2 {
		t.Errorf("decode %+v", got)
	}
	if cm.vertical {
		t.Error("horizontal parent")
	}
	// A vertical parent makes it vertical unless it sets /WMode itself.
	v := parseCMap([]byte(`/UniJIS-UCS2-V usecmap 1 begincidrange <e000> <e0ff> 20000 endcidrange`))
	if !v.vertical || v.cid(0x3001) != 7887 || v.cid(0xe001) != 20001 {
		t.Errorf("vertical %v cid %d %d", v.vertical, v.cid(0x3001), v.cid(0xe001))
	}
	h := parseCMap([]byte(`/WMode 0 def /UniJIS-UCS2-V usecmap`))
	if h.vertical {
		t.Error("/WMode 0 overrides the parent")
	}
	// Identity as the parent: unmapped codes are CIDs.
	id := parseCMap([]byte(`/Identity-H usecmap 1 begincidchar <0005> 77 endcidchar`))
	if id.cid(5) != 77 || id.cid(0x1234) != 0x1234 {
		t.Errorf("identity %d %d", id.cid(5), id.cid(0x1234))
	}
}
