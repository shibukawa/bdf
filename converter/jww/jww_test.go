package jww

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// The test drawings are made by test/jww/gen.py from Jw_cad's published
// data format; text is laid out with the test fonts of the PowerPoint
// converter.

func testOptions() *Options {
	return &Options{FontDirs: []string{"../pptx/testdata/fonts"}, NoSystemFonts: true}
}

func convert(t *testing.T, name string, opts *Options) (*Result, *bdf.Reader) {
	t.Helper()
	if opts == nil {
		opts = testOptions()
	}
	res, err := ConvertFile(filepath.Join("testdata", name), opts)
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

func plainText(t *testing.T, r *bdf.Reader) string {
	t.Helper()
	h, err := bdf.ParseHash(r.Manifest.Views[0].TextIndex)
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

func TestShapes(t *testing.T) {
	res, r := convert(t, "shapes.jww", nil)
	if len(res.Warnings) > 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	m := r.Manifest
	if len(m.Views) != 1 || len(m.Views[0].Pages) != 1 || m.Meta.Source != "jww" {
		t.Fatalf("manifest: %+v", m)
	}
	// A3
	if res.W != 420 || res.H != 297 {
		t.Errorf("page %v × %v mm", res.W, res.H)
	}
	if got := m.Meta.DC.Language.First(); got != "ja" {
		t.Errorf("language %q", got)
	}
	if got := m.Meta.DC.Description.First(); got != "テスト図" {
		t.Errorf("memo %q", got)
	}
	text := plainText(t, r)
	for _, want := range []string{"文字のサイズ 5mm", "Jw_cad 2026 ハンカク", "字間 あり", "ヨコナガ", "イタリック", "ボールド",
		"カイテン 30°", "タテガキ", "4000", "A"} {
		if !strings.Contains(text, want) {
			t.Errorf("text lacks %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "HIDDEN") {
		t.Errorf("text of a hidden layer group is drawn:\n%s", text)
	}
	// four inserts of the block with its text
	if n := strings.Count(text, "A"); n < 4 {
		t.Errorf("%d block texts", n)
	}
}

// TestOld reads a file of Jw_cad 3: no widths per figure, no SXF tables,
// 4-byte times in the block definitions.
func TestOld(t *testing.T) {
	res, r := convert(t, "old.jww", nil)
	if len(res.Warnings) > 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	if res.W != 297 || res.H != 210 {
		t.Errorf("page %v × %v mm", res.W, res.H)
	}
	if text := plainText(t, r); !strings.Contains(text, "キュウ バージョン") {
		t.Errorf("text %q", text)
	}
}

func TestColors(t *testing.T) {
	bg := func(c Colors) bdf.Color {
		opts := testOptions()
		opts.Colors = c
		_, r := convert(t, "shapes.jww", opts)
		o, err := r.Object(r.Manifest.Views[0].Pages[0].Layers[0].Obj)
		if err != nil {
			t.Fatal(err)
		}
		var first bdf.Color
		found := false
		o.Walk(func(in bdf.Instr) {
			if !found && in.Op == bdf.OpFillColor {
				first, found = bdf.Color(in.Args[0].(uint64)), true
			}
		})
		return first
	}
	if got := bg(Screen); got != bdf.RGB(255, 255, 255) {
		t.Errorf("screen background %08x", got)
	}
	if got := bg(Mono); got != bdf.RGB(255, 255, 255) {
		t.Errorf("mono background %08x", got)
	}
}

func TestLineTypes(t *testing.T) {
	c := &converter{h: &header{dpi: 300, lineTypes: map[int]lineType{
		2: {mask: 0x99999999, unit: 4, prtPitch: 10},
		3: {mask: 0xffffffff, unit: 32, prtPitch: 10},
	}}}
	bit := 10 * 25.4 / 300
	dash, off := c.dash(2)
	if len(dash) != 2 || abs(dash[0]-2*bit) > 1e-9 || abs(dash[1]-2*bit) > 1e-9 || abs(off-bit) > 1e-9 {
		t.Errorf("1001: %v %v", dash, off)
	}
	if dash, _ := c.dash(3); dash != nil {
		t.Errorf("a pattern of ones is continuous: %v", dash)
	}
	if dash, _ := c.dash(1); dash != nil {
		t.Errorf("line type 1: %v", dash)
	}
}

func TestFullWidth(t *testing.T) {
	for _, c := range []struct {
		r    rune
		want bool
	}{{'A', false}, {'ｱ', false}, {'あ', true}, {'字', true}, {'°', true}, {'±', true}, {'ü', false}} {
		if got := fullWidth(c.r); got != c.want {
			t.Errorf("fullWidth(%q) = %v", c.r, got)
		}
	}
}

func TestNotJWW(t *testing.T) {
	if _, err := Convert(bytes.NewReader([]byte("hello")), 5, nil); err == nil {
		t.Error("no error for a file that is not a JWW file")
	}
	// a truncated file converts what it has, with a warning
	full := mustRead(t, "testdata/shapes.jww")
	cut := full[:len(full)-600]
	res, err := Convert(bytes.NewReader(cut), int64(len(cut)), testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Warnings) == 0 {
		t.Error("no warning for a truncated file")
	}
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// TestBlockCycle draws a block that inserts itself three times once.
func TestBlockCycle(t *testing.T) {
	h := &header{version: 700, dpi: 300, lineTypes: map[int]lineType{}}
	for g := range h.groups {
		h.groups[g].state = 2
		for l := range h.groups[g].layers {
			h.groups[g].layers[l] = 2
		}
	}
	self := &blockDef{number: 1}
	self.list = []any{&sen{a: point{0, 0}, b: point{1, 1}}, &block{number: 1, sx: 1, sy: 1}, &block{number: 1, sx: 1, sy: 1}, &block{number: 1, sx: 1, sy: 1}}
	c := &converter{opts: &Options{}, h: h, d: &data{blocks: map[int]*blockDef{1: self}}, warned: map[string]bool{}}
	out := &cad.Drawing{}
	done := make(chan struct{})
	go func() {
		c.figures(out, []any{&block{number: 1, sx: 1, sy: 1}}, canvas.Identity, 0)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("a block that inserts itself did not finish")
	}
	if len(out.Items) != 1 || !c.warned["cycle"] {
		t.Errorf("%d items, warned %v", len(out.Items), c.warned)
	}
}
