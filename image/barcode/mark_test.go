package barcode

import (
	"slices"
	"testing"
)

// TestMarkRoundTrip checks that the symbol a mark makes again is the one
// it was made of, for every symbology and the options a mark keeps.
func TestMarkRoundTrip(t *testing.T) {
	for _, c := range []struct {
		kind, text string
		opts       Options
	}{
		{QR, "https://example.com/?q=日本語", Options{Level: "H"}},
		{QR, "12345", Options{Quiet: NoQuiet}},
		{EAN13, "490123456789", Options{}},
		{EAN13, "4901234567894", Options{NoText: true}},
		{EAN8, "4901234", Options{}},
		{UPCA, "01234567890", Options{Quiet: 3}},
		{Code39, "BDF-39", Options{}},
		{Code128, "bdf 128 0123456789", Options{}},
		{GS1128, "(01)04912345123459(10)ABC", Options{}},
		{ITF, "123456", Options{}},
		{ITF14, "1456789012345", Options{}},
		{NW7, "A1234B", Options{Quiet: NoQuiet}},
	} {
		s, err := Encode(c.kind, c.text, &c.opts)
		if err != nil {
			t.Fatal(c.kind, err)
		}
		m, err := ParseMark(s.Mark(1.5, 2, 30, 40))
		if err != nil {
			t.Fatal(c.kind, err)
		}
		if m.Box != [4]float32{1.5, 2, 30, 40} || m.Kind != s.Kind || m.NoText != c.opts.NoText {
			t.Fatalf("%s: the mark is %+v", c.kind, m)
		}
		again, err := m.Symbol()
		if err != nil {
			t.Fatal(c.kind, err)
		}
		if again.W != s.W || again.H != s.H || !slices.Equal(again.mods, s.mods) || again.Text != s.Text || again.HasText() != s.HasText() {
			t.Fatalf("%s %q: the symbol made again differs", c.kind, c.text)
		}
	}
	for _, bad := range []string{"", "{}", `{"kind":"qr","text":"x"}`, `{"kind":"pdf417","text":"x","box":[0,0,1,1]}`, `{"kind":"qr","text":"x","box":[0,0,0,1]}`} {
		if _, err := ParseMark(bad); err == nil {
			t.Errorf("%s: no error", bad)
		}
	}
}
