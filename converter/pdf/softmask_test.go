package pdf

import (
	"os"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// softMaskPDF masks a red square with an alpha mask (a half-transparent
// disc), a blue square and text with a luminosity mask (a gradient on a
// white backdrop), and sets a mask without q to end it with /SMask /None.
func softMaskPDF() []byte {
	content := `q 1 0 0 1 50 50 cm /GSa gs 1 0 0 rg 0 0 200 200 re f Q
q /GSl gs 0 0 1 rg 300 50 200 200 re f BT /F1 24 Tf 310 100 Td (masked) Tj ET Q
/GSa gs 0 1 0 rg 50 300 100 100 re f /GSn gs 0 g 200 300 100 100 re f`
	return buildPDF([]any{
		/* 1 */ `<< /Type /Catalog /Pages 2 0 R >>`,
		/* 2 */ `<< /Type /Pages /Kids [3 0 R] /Count 1 >>`,
		/* 3 */ `<< /Type /Page /Parent 2 0 R /MediaBox [0 0 600 500] /Contents 4 0 R
			/Resources << /Font << /F1 5 0 R >> /ExtGState << /GSa 6 0 R /GSl 7 0 R /GSn << /SMask /None >> >> >> >>`,
		/* 4 */ [2]string{"", content},
		/* 5 */ `<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>`,
		/* 6 */ `<< /Type /ExtGState /SMask << /S /Alpha /G 8 0 R >> >>`,
		/* 7 */ `<< /Type /ExtGState /SMask << /S /Luminosity /G 9 0 R /BC [1] >> >>`,
		/* 8 */ [2]string{"/Type /XObject /Subtype /Form /BBox [0 0 200 200] /Group << /S /Transparency >> /Resources << /ExtGState << /H << /ca 0.5 >> >> >>",
			"/H gs 0 g 100 0 m 200 100 l 100 200 l 0 100 l h f"},
		/* 9 */ [2]string{"/Type /XObject /Subtype /Form /BBox [300 50 500 250] /Group << /S /Transparency /CS /DeviceGray >> /Resources << /Shading << /Sh 10 0 R >> /Font << /F1 5 0 R >> >>",
			"/Sh sh BT /F1 10 Tf 310 60 Td (masktext) Tj ET"},
		/* 10 */ `<< /ShadingType 2 /ColorSpace /DeviceGray /Coords [300 0 500 0] /Function << /FunctionType 2 /Domain [0 1] /C0 [0] /C1 [1] /N 1 >> /Extend [true true] >>`,
	}, "/Root 1 0 R")
}

func TestSoftMask(t *testing.T) {
	res, r := convertBytes(t, softMaskPDF(), nil)
	for _, w := range res.Warnings {
		if strings.Contains(w, "mask") {
			t.Errorf("warning: %s", w)
		}
	}
	o, err := r.Object(r.Manifest.Views[0].Pages[0].Layers[0].Obj)
	if err != nil {
		t.Fatal(err)
	}
	var seq []string
	var kinds []uint64
	var backdrops []bdf.Color
	depth := 0
	if err := o.Walk(func(in bdf.Instr) {
		switch in.Op {
		case bdf.OpSave, bdf.OpGroupBegin, bdf.OpMaskBegin:
			depth++
		case bdf.OpRestore, bdf.OpGroupEnd, bdf.OpMaskEnd:
			depth--
		}
		switch in.Op {
		case bdf.OpGroupBegin:
			seq = append(seq, "G")
		case bdf.OpGroupEnd:
			seq = append(seq, "g")
		case bdf.OpMaskBegin:
			seq = append(seq, "M")
			kinds = append(kinds, in.Args[0].(uint64))
			backdrops = append(backdrops, bdf.Color(in.Args[1].(uint64)))
		case bdf.OpMaskEnd:
			seq = append(seq, "m")
		case bdf.OpUse:
			seq = append(seq, "U")
		case bdf.OpFillPath:
			seq = append(seq, "f")
		case bdf.OpFillText:
			seq = append(seq, "t")
		}
	}); err != nil {
		t.Fatal(err)
	}
	if depth != 0 {
		t.Errorf("unbalanced: depth %d", depth)
	}
	// The last square is painted after /SMask /None, outside any group.
	if got, want := strings.Join(seq, ""), "GfMUmgGftMUmgGfMUmgf"; got != want {
		t.Errorf("ops %s, want %s", got, want)
	}
	if len(kinds) != 3 || kinds[0] != uint64(bdf.MaskAlpha) || kinds[1] != uint64(bdf.MaskLuminosity) || kinds[2] != uint64(bdf.MaskAlpha) {
		t.Errorf("mask kinds %v", kinds)
	}
	if len(backdrops) > 1 && backdrops[1] != bdf.RGB(255, 255, 255) {
		t.Errorf("luminosity backdrop %v, want white", backdrops[1])
	}
	// Text inside a mask is not content.
	if txt := pageText(t, r); !strings.Contains(txt, "masked") || strings.Contains(txt, "masktext") {
		t.Errorf("text = %q", txt)
	}
}

func TestWriteSoftMaskPDF(t *testing.T) {
	if p := os.Getenv("BDF_WRITE_SOFTMASK"); p != "" {
		if err := os.WriteFile(p, softMaskPDF(), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
