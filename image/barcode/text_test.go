package barcode

import (
	"reflect"
	"testing"

	"github.com/shibukawa/bdf"
)

func TestHumanReadableDigits(t *testing.T) {
	for _, tc := range []struct{ kind, input, full string }{
		{EAN13, "490123456789", "4901234567894"},
		{EAN8, "9638507", "96385074"},
		{UPCA, "03600029145", "036000291452"},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			short := mustEncode(t, tc.kind, tc.input, nil)
			full := mustEncode(t, tc.kind, tc.full, nil)
			bars := mustEncode(t, tc.kind, tc.input, &Options{NoText: true})
			if short.Text != tc.full || !reflect.DeepEqual(short.Path(0, 0, 113, 30), full.Path(0, 0, 113, 30)) {
				t.Fatal("the displayed digits must include the automatically added check digit")
			}
			if row(short, 0) != row(bars, 0) || !reflect.DeepEqual(short.Rects(), bars.Rects()) {
				t.Fatal("human-readable text changed the encoded modules")
			}
			p := short.Path(0, 0, 113, 30)
			n := len(short.Rects())
			if len(p.Verbs) <= n || p.Verbs[n] != bdf.VerbMove || p.Args[3] >= 30 {
				t.Fatal("the digits need their own outlines below shortened bars")
			}
			for i, v := range bars.Path(0, 0, 113, 30).Verbs {
				if v != bdf.VerbRect || p.Verbs[i] != v {
					t.Fatal("the bar-only option must preserve the bars")
				}
			}
			// A standalone object includes both bars and digits without
			// depending on an embedded font or any other resource.
			o := short.Object(113, 30)
			if o.BBox.W != 113 || o.BBox.H != 30 || len(o.Deps()) != 0 {
				t.Fatalf("standalone object: bbox %v, dependencies %v", o.BBox, o.Deps())
			}
		})
	}
}

func TestDigitBoundsAndQuietZone(t *testing.T) {
	for _, quiet := range []int{0, NoQuiet, 2, 64} {
		for _, box := range [][4]float32{{10, 20, 113, 30}, {0, 0, 37, 10}, {-5, -8, 8, 100}} {
			s := mustEncode(t, EAN13, "012345678901", &Options{Quiet: quiet})
			p := s.Path(box[0], box[1], box[2], box[3])
			mx := box[2] / float32(s.W)
			first := 0
			for !s.At(first, 0) {
				first++
			}
			if p.Args[0] != box[0]+float32(first)*mx || p.Args[2] != mx {
				t.Fatal("the label changed the quiet zone or bar width")
			}
			at := 0
			check := func(x, y float32) {
				t.Helper()
				const eps = 0.0001
				if x < box[0]-eps || y < box[1]-eps || x > box[0]+box[2]+eps || y > box[1]+box[3]+eps {
					t.Fatalf("point (%g, %g) is outside %v", x, y, box)
				}
			}
			for _, v := range p.Verbs {
				switch v {
				case bdf.VerbRect:
					check(p.Args[at], p.Args[at+1])
					check(p.Args[at]+p.Args[at+2], p.Args[at+1]+p.Args[at+3])
					at += 4
				case bdf.VerbMove, bdf.VerbLine:
					check(p.Args[at], p.Args[at+1])
					at += 2
				case bdf.VerbQuad:
					check(p.Args[at], p.Args[at+1])
					check(p.Args[at+2], p.Args[at+3])
					at += 4
				case bdf.VerbClose:
				default:
					t.Fatalf("unexpected path verb %d", v)
				}
			}
		}
	}
	for _, kind := range []string{QR, Code128} {
		a := mustEncode(t, kind, "0123456789", nil)
		b := mustEncode(t, kind, "0123456789", &Options{NoText: true})
		if !reflect.DeepEqual(a.Path(0, 0, 100, 30), b.Path(0, 0, 100, 30)) {
			t.Fatalf("NoText changed a %s symbol", kind)
		}
	}
}
