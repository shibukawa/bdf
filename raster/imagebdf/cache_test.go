package imagebdf

import (
	"bytes"
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/shibukawa/bdf"
)

func TestCompletedImagesOwnedAndInvalidated(t *testing.T) {
	o := bdf.NewObject()
	o.FillColor(0xff0000ff).FillRect(.3, .3, 4.5, 4.5)
	doc, v := onePage(o, 8, 8)
	bg := color.RGBA{255, 255, 255, 255}
	r := New(doc, &Options{Background: &bg, NoSystemFonts: true})
	fresh := func() *image.RGBA {
		t.Helper()
		x, e := New(doc, &Options{Background: &bg, NoSystemFonts: true, ImageCacheBytes: -1, TextCacheBytes: -1}).Page(v, 0, 1)
		if e != nil {
			t.Fatal(e)
		}
		return x
	}
	draw := func() *image.RGBA {
		t.Helper()
		x, e := r.Page(v, 0, 1)
		if e != nil {
			t.Fatal(e)
		}
		return x
	}
	first := draw()
	want := fresh()
	clear(first.Pix)
	second := draw()
	assertImagesEqual(t, second, want)
	clear(second.Pix)
	assertImagesEqual(t, draw(), want)
	if len(r.images.items) == 0 {
		t.Fatal("completed image was not retained")
	}
	v.Pages[0].W = 12
	assertImagesEqual(t, draw(), fresh())
	bg.G = 12
	assertImagesEqual(t, draw(), fresh())
	replacement := bdf.NewObject()
	replacement.FillColor(0x0000ffff).FillRect(0, 0, 5, 5)
	rd := bdf.NewDocument()
	replacementHash, _ := rd.AddObject(replacement)
	doc.Part(v.Pages[0].Layers[0].Obj).Data = rd.Part(replacementHash).Data
	r.Invalidate()
	if r.fonts.Doc != doc || len(r.images.items) != 0 || len(r.textMasks.items) != 0 {
		t.Fatal("Invalidate did not reset the fonts and caches")
	}
	assertImagesEqual(t, draw(), fresh())
}

func TestImageCacheEntryBudgetAndKeys(t *testing.T) {
	o := bdf.NewObject()
	o.FillColor(0xff0000ff).FillRect(.3, .3, 5, 5)
	doc, v := onePage(o, 8, 8)
	r := New(doc, &Options{ImageCacheBytes: 1100, NoSystemFonts: true})
	for _, mode := range []AntialiasMode{AntialiasOn, AntialiasOff} {
		r.opts.Antialias = mode
		for _, scale := range []float64{1, 2, 3} {
			img, e := r.Page(v, 0, scale)
			if e != nil {
				t.Fatal(e)
			}
			want, e := New(doc, &Options{Antialias: mode, ImageCacheBytes: -1, NoSystemFonts: true}).Page(v, 0, scale)
			if e != nil {
				t.Fatal(e)
			}
			assertImagesEqual(t, img, want)
			if r.images.bytes > 1100 {
				t.Fatal("budget exceeded")
			}
		}
	}
	if len(r.images.items) > 1 {
		t.Fatalf("oversize entries retained: %d", len(r.images.items))
	}
	r.opts.Antialias = AntialiasOn
	if _, ok := r.imageKey(1, v.Pages[0], bdf.Rect{X: float32(math.NaN()), W: 8, H: 8}, 8, 8, 0, 0); ok {
		t.Fatal("NaN cache key accepted")
	}
	k := textMaskKey{m: identity, size: math.NaN()}
	r.keepTextMask(k, &mask{r: image.Rect(0, 0, 8, 8)})
	if len(r.textMasks.items) != 0 {
		t.Fatal("NaN text key retained")
	}
	disabled := New(doc, &Options{ImageCacheBytes: -1, TextCacheBytes: -1, NoSystemFonts: true})
	if _, e := disabled.Page(v, 0, 1); e != nil {
		t.Fatal(e)
	}
	if len(disabled.images.items) != 0 || len(disabled.textMasks.items) != 0 {
		t.Fatal("disabled caches retained entries")
	}
}

func TestRegionAndDocumentObjectCache(t *testing.T) {
	o := bdf.NewObject()
	o.FillColor(0xff0000ff).FillRect(-3.25, -1.5, 8, 9)
	doc, v := onePage(o, 16, 16)
	hash := v.Pages[0].Layers[0].Obj
	r := New(doc, &Options{NoSystemFonts: true})
	for _, draw := range []func() (*image.RGBA, error){func() (*image.RGBA, error) { return r.Region(v, 0, bdf.Rect{X: 1, Y: 2, W: 8, H: 8}, 16, 16) }, func() (*image.RGBA, error) { return r.Object(hash, 1.5) }} {
		a, e := draw()
		if e != nil {
			t.Fatal(e)
		}
		want := cloneRGBA(a)
		clear(a.Pix)
		b, e := draw()
		if e != nil {
			t.Fatal(e)
		}
		assertImagesEqual(t, b, want)
	}
	standalone, e := Object(o, 1.5, &Options{NoSystemFonts: true})
	if e != nil {
		t.Fatal(e)
	}
	fromDoc, e := r.Object(hash, 1.5)
	if e != nil {
		t.Fatal(e)
	}
	assertImagesEqual(t, fromDoc, standalone)
	// A full-page region and PageSize use the same transform and result.
	a, e := r.PageSize(v, 0, 24, 24)
	if e != nil {
		t.Fatal(e)
	}
	b, e := r.Region(v, 0, bdf.Rect{W: 16, H: 16}, 24, 24)
	if e != nil {
		t.Fatal(e)
	}
	assertImagesEqual(t, a, b)
}

func TestTextMaskCacheCompatibility(t *testing.T) {
	for _, name := range []string{"pdf/chrome-doc", "pptx/features", "drawio/showcase"} {
		doc := openDoc(t, "../../testdata/"+name+".bdf")
		for _, mode := range []AntialiasMode{AntialiasOn, AntialiasOff} {
			for _, dpi := range []float64{72, 300} {
				opts := Options{FontDirs: testFonts, NoSystemFonts: true, Antialias: mode, ImageCacheBytes: -1}
				r := New(doc, &opts)
				opts.TextCacheBytes = -1
				reference := New(doc, &opts)
				want, e := reference.PageDPI(doc.Views[0], 0, dpi)
				if e != nil {
					t.Fatal(e)
				}
				for i := 0; i < 2; i++ {
					got, e := r.PageDPI(doc.Views[0], 0, dpi)
					if e != nil {
						t.Fatal(e)
					}
					assertImagesEqual(t, got, want)
				}
				if len(r.textMasks.items) == 0 {
					t.Fatalf("%s %.0f dpi did not cache text", name, dpi)
				}
				if r.textMasks.bytes > r.textMasks.limit {
					t.Fatal("text cache budget exceeded")
				}
			}
		}
	}
}

func TestInvalidateFontFS(t *testing.T) {
	files, e := filepath.Glob("../../converter/pptx/testdata/fonts/*.ttf")
	if e != nil || len(files) == 0 {
		t.Fatal("font fixture missing")
	}
	data, e := os.ReadFile(files[0])
	if e != nil {
		t.Fatal(e)
	}
	fsys := fstest.MapFS{}
	o := bdf.NewObject()
	o.Font(o.AddFont(bdf.SystemFont("sans-serif", 400, bdf.StyleNormal)), 18).FillColor(0x000000ff).FillText("ABC", 2, 22, 0)
	doc, v := onePage(o, 80, 32)
	opts := Options{FontFS: fsys, NoSystemFonts: true}
	r := New(doc, &opts)
	before, e := r.Page(v, 0, 1)
	if e != nil {
		t.Fatal(e)
	}
	fsys["font.ttf"] = &fstest.MapFile{Data: data}
	r.Invalidate()
	got, e := r.Page(v, 0, 1)
	if e != nil {
		t.Fatal(e)
	}
	want, e := New(doc, &opts).Page(v, 0, 1)
	if e != nil {
		t.Fatal(e)
	}
	assertImagesEqual(t, got, want)
	if bytes.Equal(before.Pix, got.Pix) {
		t.Fatal("FontFS update did not change the rendered text")
	}
}

func TestCacheEvictsLeastRecentlyUsed(t *testing.T) {
	c := lru[int, string]{limit: 2}
	c.put(1, "one", 1)
	c.put(2, "two", 1)
	if _, ok := c.get(1); !ok {
		t.Fatal("missing first entry")
	}
	c.put(3, "three", 1)
	if _, ok := c.get(2); ok {
		t.Fatal("recently unused entry was retained")
	}
	if v, ok := c.get(1); !ok || v != "one" {
		t.Fatal("recently used entry evicted")
	}
	c.put(4, "oversized", 3)
	if _, ok := c.get(4); ok || c.bytes > c.limit {
		t.Fatal("oversized entry exceeded budget")
	}
}
