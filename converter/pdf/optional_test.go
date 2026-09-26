package pdf

import (
	"fmt"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// paintSeq lists what an object and the objects it USEs paint, in walk
// order: fills with their colour, clips, text runs with their x, links.
func paintSeq(t *testing.T, r *bdf.Reader, h bdf.Hash) string {
	t.Helper()
	var out []string
	var walk func(h bdf.Hash)
	walk = func(h bdf.Hash) {
		o, err := r.Object(h)
		if err != nil {
			t.Fatal(err)
		}
		fill := bdf.RGB(0, 0, 0)
		err = o.Walk(func(in bdf.Instr) {
			switch in.Op {
			case bdf.OpFillColor:
				fill = bdf.Color(in.Args[0].(uint64))
			case bdf.OpFillPath, bdf.OpFillRect:
				out = append(out, "fill "+fill.CSS())
			case bdf.OpClipPath, bdf.OpClipRect:
				out = append(out, "clip")
			case bdf.OpFillText:
				out = append(out, fmt.Sprintf("text %s@%.2f", in.Args[0], in.Args[1]))
			case bdf.OpLink:
				out = append(out, "link "+in.Args[4].(string))
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

// ocPDF is a page with optional content: group A is on, B off, C not listed
// (so in the base state, on), and membership dictionaries over them.
func ocPDF(content string) []byte {
	return buildPDF([]any{
		/* 1 */ `<< /Type /Catalog /Pages 2 0 R /OCProperties << /OCGs [6 0 R 7 0 R 8 0 R] /D << /ON [6 0 R] /OFF [7 0 R] >> >> >>`,
		/* 2 */ `<< /Type /Pages /Kids [3 0 R] /Count 1 >>`,
		/* 3 */ `<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R /Annots [14 0 R 15 0 R]
			/Resources << /Font << /F1 5 0 R >> /XObject << /Fm 12 0 R /FmHidden 13 0 R >>
			/Properties << /A 6 0 R /B 7 0 R /C 8 0 R /AllOn 9 0 R /VE 10 0 R /AnyOff 11 0 R >> >> >>`,
		/* 4 */ [2]string{"", content},
		/* 5 */ `<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>`,
		/* 6 */ `<< /Type /OCG /Name (A) >>`,
		/* 7 */ `<< /Type /OCG /Name (B) >>`,
		/* 8 */ `<< /Type /OCG /Name (C) >>`,
		/* 9 */ `<< /Type /OCMD /OCGs [6 0 R 7 0 R] /P /AllOn >>`,
		// The expression wins over /OCGs, which alone would hide.
		/* 10 */ `<< /Type /OCMD /OCGs 7 0 R /VE [/And 6 0 R [/Not 7 0 R]] >>`,
		/* 11 */ `<< /Type /OCMD /OCGs [6 0 R 7 0 R] /P /AnyOff >>`,
		/* 12 */ [2]string{"/Type /XObject /Subtype /Form /BBox [0 0 10 10]", "1 1 1 rg 0 0 10 10 re f"},
		/* 13 */ [2]string{"/Type /XObject /Subtype /Form /BBox [0 0 10 10] /OC 7 0 R", "0 0 0 rg 0 0 10 10 re f"},
		/* 14 */ `<< /Type /Annot /Subtype /Link /Rect [0 0 10 10] /A << /S /URI /URI (https://hidden.example/) >> /OC 7 0 R >>`,
		/* 15 */ `<< /Type /Annot /Subtype /Link /Rect [0 0 10 10] /A << /S /URI /URI (https://shown.example/) >> >>`,
	}, "/Root 1 0 R")
}

func TestOptionalContent(t *testing.T) {
	content := `1 0 0 rg /OC /A BDC 0 0 10 10 re f EMC
0 1 0 rg /OC /B BDC 10 0 10 10 re f EMC
/OC /C BDC 0 0 1 rg 20 0 10 10 re f EMC
/OC /AllOn BDC 1 1 0 rg 30 0 10 10 re f EMC
/OC /VE BDC 0 1 1 rg 40 0 10 10 re f EMC
/OC /AnyOff BDC 1 0 1 rg 50 0 10 10 re f EMC
/OC /B BDC /OC /A BDC 0.5 g 60 0 10 10 re f EMC EMC
q /OC /B BDC 100 100 50 50 re W n EMC 0 g 100 100 10 10 re f Q
/Fm Do /FmHidden Do
BT /F1 12 Tf 72 700 Td (one ) Tj /OC /B BDC (two ) Tj EMC (three) Tj ET`
	res, r := convertBytes(t, ocPDF(content), nil)
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	got := paintSeq(t, r, r.Manifest.Views[0].Pages[0].Layers[0].Obj)
	// The hidden text is not drawn but moves the pen: "three" starts where
	// it does when "two " is drawn.
	_, shown := convertBytes(t, ocPDF(`BT /F1 12 Tf 72 700 Td (one ) Tj 3 Tr (two ) Tj 0 Tr (three) Tj ET`), nil)
	var x string
	for _, s := range strings.Split(paintSeq(t, shown, shown.Manifest.Views[0].Pages[0].Layers[0].Obj), "|") {
		if after, ok := strings.CutPrefix(s, "text three@"); ok {
			x = after
		}
	}
	want := strings.Join([]string{
		"fill #ff0000", "fill #0000ff", "fill #00ffff", "fill #ff00ff",
		"clip", "fill #000000", // the hidden clip still applies
		"clip", "fill #ffffff", // the form without /OC, clipped to its bounding box
		"text one @0.00", "text three@" + x,
		"link https://shown.example/",
	}, "|")
	if got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
	if text := pageText(t, r); text != "one three" && text != "one  three" {
		t.Errorf("text = %q", text)
	}

	// Without /OCProperties, optional content is ignored: everything shows.
	data := strings.Replace(string(ocPDF(`/OC /B BDC 0 0 10 10 re f EMC`)), "/OCProperties", "/Unused", 1)
	_, r = convertBytes(t, []byte(data), nil)
	if got := paintSeq(t, r, r.Manifest.Views[0].Pages[0].Layers[0].Obj); !strings.HasPrefix(got, "fill #000000|") {
		t.Errorf("without /OCProperties: %s", got)
	}
}

func TestPageBox(t *testing.T) {
	data := buildPDF([]any{
		`<< /Type /Catalog /Pages 2 0 R >>`,
		`<< /Type /Pages /Kids [3 0 R] /Count 1 /MediaBox [0 0 200 100] >>`,
		`<< /Type /Page /Parent 2 0 R /CropBox [0 0 190 100] /BleedBox [5 5 195 95] /TrimBox [10 10 190 90]
			/ArtBox [20 20 30 30] /Contents 4 0 R >>`,
		[2]string{"", "0 0 200 100 re f"},
	}, "/Root 1 0 R")
	for _, c := range []struct {
		box  string
		w, h float32
	}{
		{"", 190, 100}, {"crop", 190, 100}, {"media", 200, 100},
		{"bleed", 185, 90}, // the crop box bounds it
		{"trim", 180, 80}, {"art", 10, 10},
	} {
		_, r := convertBytes(t, data, &Options{Box: c.box})
		if p := r.Manifest.Views[0].Pages[0]; p.W != c.w || p.H != c.h {
			t.Errorf("box %q: page %v × %v, want %v × %v", c.box, p.W, p.H, c.w, c.h)
		}
	}
	// A page without the box falls back to its crop box.
	noTrim := strings.Replace(string(data), "/TrimBox [10 10 190 90]", "", 1)
	_, r := convertBytes(t, []byte(noTrim), &Options{Box: "trim"})
	if p := r.Manifest.Views[0].Pages[0]; p.W != 190 || p.H != 100 {
		t.Errorf("no trim box: page %v × %v", p.W, p.H)
	}
	if _, err := Convert(strings.NewReader(string(data)), &Options{Box: "paper"}); err == nil {
		t.Error("an unknown box converted")
	}
}

func TestOptionalContentVertical(t *testing.T) {
	// Hidden vertical text moves the pen down the line, as drawn text does:
	// the third 亜 starts where it does when the second is drawn invisibly.
	build := func(content string) []byte {
		return buildPDF([]any{
			/* 1 */ `<< /Type /Catalog /Pages 2 0 R /OCProperties << /OCGs [6 0 R] /D << /OFF [6 0 R] >> >> >>`,
			/* 2 */ `<< /Type /Pages /Kids [3 0 R] /Count 1 >>`,
			/* 3 */ `<< /Type /Page /Parent 2 0 R /MediaBox [0 0 400 400] /Contents 4 0 R
				/Resources << /Font << /F1 5 0 R >> /Properties << /B 6 0 R >> >> >>`,
			/* 4 */ [2]string{"", content},
			/* 5 */ `<< /Type /Font /Subtype /Type0 /BaseFont /Ryumin-Light-Identity-V /Encoding /Identity-V /DescendantFonts [7 0 R] >>`,
			/* 6 */ `<< /Type /OCG /Name (B) >>`,
			/* 7 */ `<< /Type /Font /Subtype /CIDFontType0 /BaseFont /Ryumin-Light /CIDSystemInfo << /Registry (Adobe) /Ordering (Japan1) /Supplement 2 >>
				/FontDescriptor << /Type /FontDescriptor /FontName /Ryumin-Light /Flags 6 /FontBBox [0 -141 1000 859] /ItalicAngle 0 /Ascent 859 /Descent -141 /CapHeight 709 /StemV 69 >> >>`,
		}, "/Root 1 0 R")
	}
	// the text and the line position (the TRANSFORM before it) of each vertical run
	runs := func(r *bdf.Reader) []string {
		o, err := r.Object(r.Manifest.Views[0].Pages[0].Layers[0].Obj)
		if err != nil {
			t.Fatal(err)
		}
		var out []string
		var at string
		if err := o.Walk(func(in bdf.Instr) {
			switch in.Op {
			case bdf.OpTransform:
				at = fmt.Sprintf("%.2f,%.2f", in.Args[4], in.Args[5])
			case bdf.OpMark:
				if byte(in.Args[0].(uint64)) == bdf.MarkAltText {
					out = append(out, in.Args[1].(string)+"@"+at)
				}
			}
		}); err != nil {
			t.Fatal(err)
		}
		return out
	}
	_, hidden := convertBytes(t, build(`BT /F1 20 Tf 100 350 Td <0465> Tj /OC /B BDC <0465> Tj EMC <0465> Tj ET`), nil)
	_, shown := convertBytes(t, build(`BT /F1 20 Tf 100 350 Td <0465> Tj 3 Tr <0465> Tj 0 Tr <0465> Tj ET`), nil)
	got, ref := runs(hidden), runs(shown)
	if len(got) != 2 || len(ref) != 3 || got[0] != ref[0] || got[1] != ref[2] || ref[1] == ref[2] {
		t.Errorf("hidden: %q, drawn invisibly: %q", got, ref)
	}
}
