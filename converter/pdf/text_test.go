package pdf

import (
	"fmt"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

func TestStrokedTextLineWidth(t *testing.T) {
	// The line width and the dashes of the graphics state are in user
	// space for glyphs as for paths: a faux bold (Tr 2) set with the size
	// in Tm and a font size of 1, as Word and some Japanese producers do,
	// strokes as thin as one with the size in Tf. Each stroke is listed
	// with the LINE and DASH set before it, in the space it is drawn in.
	data := pagePDF(`/Font << /F1 5 0 R /F2 6 0 R >>`, [2]string{"", `
BT /F1 24 Tf 20 700 Td 2 Tr 1.5 w (Hm) Tj ET
BT /F1 1 Tf 24 0 0 24 20 650 Tm 2 Tr 1.5 w [6 3] 1.5 d (Hm) Tj ET
BT /F1 1 Tf 50 Tz 24 0 0 24 20 600 Tm 1 Tr 1.5 w [] 0 d (Hm) Tj ET
BT /F2 1 Tf 100 Tz 20 0 0 20 300 700 Tm 1 Tr 2 w <0465> Tj ET
2 w 20 500 m 200 500 l S`},
		`<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>`,
		`<< /Type /Font /Subtype /Type0 /BaseFont /Ryumin-Light-Identity-V /Encoding /Identity-V /DescendantFonts [7 0 R] >>`,
		`<< /Type /Font /Subtype /CIDFontType0 /BaseFont /Ryumin-Light /CIDSystemInfo << /Registry (Adobe) /Ordering (Japan1) /Supplement 2 >>
			/FontDescriptor << /Type /FontDescriptor /FontName /Ryumin-Light /Flags 6 /FontBBox [0 -141 1000 859] /ItalicAngle 0 /Ascent 859 /Descent -141 /CapHeight 709 /StemV 69 >> >>`)
	_, r := convertBytes(t, data, nil)
	o, err := r.Object(r.Manifest.Views[0].Pages[0].Layers[0].Obj)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	line, dash := "1", "[]/0"
	var stack [][2]string
	if err := o.Walk(func(in bdf.Instr) {
		switch in.Op {
		case bdf.OpSave:
			stack = append(stack, [2]string{line, dash})
		case bdf.OpRestore:
			line, dash = stack[len(stack)-1][0], stack[len(stack)-1][1]
			stack = stack[:len(stack)-1]
		case bdf.OpLine:
			line = fmt.Sprintf("%.4g", in.Args[0])
		case bdf.OpDash:
			dash = fmt.Sprintf("%.4g/%.4g", in.Args[0], in.Args[1])
		case bdf.OpStrokeText:
			got = append(got, "text "+line+" "+dash)
		case bdf.OpUseAt: // a vertical run, drawn by its child
			got = append(got, "vertical "+line+" "+dash)
		case bdf.OpStrokePath:
			got = append(got, "path "+line+" "+dash)
		}
	}); err != nil {
		t.Fatal(err)
	}
	want := strings.Join([]string{
		"text 1.5 []/0",
		"text 0.0625 [0.25 0.125]/0.0625",
		// Tz 50 scales the block by 0.5 across: the mean of both ways
		"text 0.08839 []/0",
		"vertical 0.1 []/0",
		"path 2 []/0",
	}, "\n")
	if s := strings.Join(got, "\n"); s != want {
		t.Errorf("strokes:\n%s\nwant:\n%s", s, want)
	}
}
