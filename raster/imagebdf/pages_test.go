package imagebdf

import (
	stdcontext "context"
	"errors"
	"github.com/shibukawa/bdf"
	"image"
	"slices"
	"testing"
)

func TestPagesOrderOwnershipAndCancellation(t *testing.T) {
	doc := bdf.NewDocument()
	v := doc.NewView("v", bdf.ViewFixed, "")
	for _, c := range []bdf.Color{0xff0000ff, 0x00ff00ff, 0x0000ffff} {
		o := bdf.NewObject()
		o.FillColor(c).FillRect(.3, .3, 8, 8)
		h, _ := doc.AddObject(o)
		v.AddPage(12, 12, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	}
	for _, mode := range []AntialiasMode{AntialiasOn, AntialiasOff} {
		for _, workers := range []int{0, 1, 4} {
			r := New(doc, &Options{NoSystemFonts: true, Antialias: mode})
			order := []int{}
			pages := []int{2, 0, 1, 2, 0, 2, 1}
			err := r.PagesDPI(stdcontext.Background(), v, pages, 300, workers, func(page int, img *image.RGBA) error {
				order = append(order, page)
				want, e := New(doc, &Options{NoSystemFonts: true, Antialias: mode, ImageCacheBytes: -1, TextCacheBytes: -1}).PageDPI(v, page, 300)
				if e != nil {
					return e
				}
				assertImagesEqual(t, img, want)
				clear(img.Pix)
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(order, pages) {
				t.Fatalf("order %v", order)
			}
			got, e := r.PageDPI(v, 2, 300)
			if e != nil {
				t.Fatal(e)
			}
			want, e := New(doc, &Options{NoSystemFonts: true, Antialias: mode}).PageDPI(v, 2, 300)
			if e != nil {
				t.Fatal(e)
			}
			assertImagesEqual(t, got, want)
		}
	}
	r := New(doc, nil)
	stop := errors.New("stop")
	calls := 0
	if e := r.Pages(stdcontext.Background(), v, nil, 1, 4, func(int, *image.RGBA) error { calls++; return stop }); !errors.Is(e, stop) || calls != 1 {
		t.Fatalf("callback error: %v, %d", e, calls)
	}
	ctx, cancel := stdcontext.WithCancel(stdcontext.Background())
	cancel()
	if e := r.Pages(ctx, v, nil, 1, 4, func(int, *image.RGBA) error { t.Fatal("called after cancellation"); return nil }); !errors.Is(e, stdcontext.Canceled) {
		t.Fatal(e)
	}
	ctx, cancel = stdcontext.WithCancel(stdcontext.Background())
	defer cancel()
	calls = 0
	if e := r.Pages(ctx, v, nil, 1, 2, func(int, *image.RGBA) error { calls++; cancel(); return nil }); !errors.Is(e, stdcontext.Canceled) || calls != 1 {
		t.Fatalf("cancellation: %v calls=%d", e, calls)
	}
	calls = 0
	if e := r.Pages(stdcontext.Background(), v, []int{0, 99}, 1, 4, func(int, *image.RGBA) error { calls++; return nil }); !errors.Is(e, ErrNoPage) || calls != 0 {
		t.Fatalf("invalid selection: %v, %d", e, calls)
	}
	if e := r.Pages(stdcontext.Background(), v, []int{}, 1, 4, func(int, *image.RGBA) error { t.Fatal("empty selection"); return nil }); e != nil {
		t.Fatal(e)
	}
	if e := r.Pages(stdcontext.Background(), v, nil, 1, 17, func(int, *image.RGBA) error { return nil }); e == nil {
		t.Fatal("invalid workers accepted")
	}
}

func TestPagesSharedDocumentFontsAndWarnings(t *testing.T) {
	for _, name := range []string{"pdf/chrome-doc", "pptx/features", "drawio/showcase"} {
		doc := openDoc(t, "../../testdata/"+name+".bdf")
		v := doc.Views[0]
		opts := Options{NoSystemFonts: true, FontDirs: testFonts, ImageCacheBytes: -1}
		r := New(doc, &opts)
		want, e := New(doc, &opts).Page(v, 0, 1)
		if e != nil {
			t.Fatal(e)
		}
		e = r.Pages(stdcontext.Background(), v, []int{0, 0, 0, 0, 0}, 1, 4, func(_ int, img *image.RGBA) error { assertImagesEqual(t, img, want); return nil })
		if e != nil {
			t.Fatal(e)
		}
	}
	o := bdf.NewObject()
	o.Font(o.AddFont(bdf.SystemFont("unavailable-font", 400, bdf.StyleNormal)), 12).FillText("Missing", 0, 12, 0)
	doc, v := onePage(o, 12, 12)
	r := New(doc, &Options{NoSystemFonts: true})
	if e := r.Pages(stdcontext.Background(), v, []int{0, 0, 0}, 1, 3, func(int, *image.RGBA) error { return nil }); e != nil {
		t.Fatal(e)
	}
	before := slices.Clone(r.Warnings())
	if len(before) == 0 {
		t.Fatal("worker warning lost")
	}
	if e := r.Pages(stdcontext.Background(), v, []int{0, 0}, 1, 2, func(int, *image.RGBA) error { return nil }); e != nil {
		t.Fatal(e)
	}
	if !slices.Equal(before, r.Warnings()) {
		t.Fatal("warning duplicated on cache hit")
	}
}
