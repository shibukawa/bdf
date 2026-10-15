package imagebdf

import (
	"bytes"
	stdcontext "context"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf"
)

// assembled draws a page in bands and puts them together, checking the
// bands as they come: in order, as wide as the page, and of the kind of
// the mono. The image returned is of the kind of the mono.
func assembled(t testing.TB, r *Renderer, v *bdf.View, page int, scale float64, mono Mono) (image.Image, int) {
	t.Helper()
	var out draw.Image
	next, bands := 0, 0
	var last []byte
	err := r.PageRows(v, page, scale, mono, func(y0 int, rows image.Image) error {
		b := rows.Bounds()
		if y0 != next || b.Min.Y != y0 || b.Min.X != 0 || b.Dy() <= 0 {
			t.Fatalf("page %d: band at %d, bounds %v, after %d rows", page, y0, b, next)
		}
		var pix []byte
		switch img := rows.(type) {
		case *image.RGBA:
			if mono != 0 {
				t.Fatalf("an RGBA band for mono %d", mono)
			}
			pix = img.Pix
			if out == nil {
				out = image.NewRGBA(image.Rect(0, 0, b.Dx(), 1<<20))
			}
		case *image.Gray:
			if mono != Gray {
				t.Fatalf("a gray band for mono %d", mono)
			}
			pix = img.Pix
			if out == nil {
				out = image.NewGray(image.Rect(0, 0, b.Dx(), 1<<20))
			}
		case *image.Paletted:
			if mono != Bilevel || len(img.Palette) != 2 {
				t.Fatalf("a %T of %d colours for mono %d", rows, len(img.Palette), mono)
			}
			pix = img.Pix
			if out == nil {
				out = image.NewPaletted(image.Rect(0, 0, b.Dx(), 1<<20), img.Palette)
			}
		default:
			t.Fatalf("a %T band", rows)
		}
		if bands > 0 && last != nil && &last[0] != &pix[0] && cap(last) >= len(pix) {
			t.Errorf("page %d: the band at %d is not drawn into the memory of the one before", page, y0)
		}
		last = pix
		draw.Draw(out, b, rows, b.Min, draw.Src)
		next = b.Max.Y
		bands++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// cut to the rows drawn
	switch img := out.(type) {
	case *image.RGBA:
		return img.SubImage(image.Rect(0, 0, img.Rect.Dx(), next)), bands
	case *image.Gray:
		return img.SubImage(image.Rect(0, 0, img.Rect.Dx(), next)), bands
	case *image.Paletted:
		return img.SubImage(image.Rect(0, 0, img.Rect.Dx(), next)), bands
	}
	return nil, bands
}

// samePixels compares two images of the same kind pixel by pixel and
// returns how many differ, and the first that does.
func samePixels(a, b image.Image) (n int, first image.Point) {
	ab, bb := a.Bounds(), b.Bounds()
	if ab.Dx() != bb.Dx() || ab.Dy() != bb.Dy() {
		return -1, image.Point{}
	}
	for y := 0; y < ab.Dy(); y++ {
		for x := 0; x < ab.Dx(); x++ {
			if a.At(ab.Min.X+x, ab.Min.Y+y) != b.At(bb.Min.X+x, bb.Min.Y+y) {
				if n == 0 {
					first = image.Pt(x, y)
				}
				n++
			}
		}
	}
	return n, first
}

// TestRowsChangeNoPixel draws the pages of the documents of testdata in
// bands of a few rows and whole, and checks that the images are the
// same.
func TestRowsChangeNoPixel(t *testing.T) {
	if testing.Short() {
		t.Skip("slow")
	}
	files, _ := filepath.Glob("../../testdata/*/*.bdf")
	if len(files) < 50 {
		t.Fatalf("%d documents in testdata", len(files))
	}
	files = append(files, "../../testdata/demo.bdf")
	for _, f := range files {
		rd, err := bdf.OpenSingleFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if rd.Encrypted() {
			rd.Close()
			continue
		}
		d, err := rd.ToDocument()
		rd.Close()
		if err != nil {
			t.Fatal(f, err)
		}
		for _, v := range d.Views {
			if v.Kind == bdf.ViewSheet || len(v.Pages) == 0 {
				continue
			}
			for _, mode := range []AntialiasMode{AntialiasOn, AntialiasOff} {
				r := New(d, &Options{FontDirs: testFonts, NoSystemFonts: true, Antialias: mode, BandRows: 48, ImageCacheBytes: -1})
				for _, mono := range []Mono{0, Gray, Bilevel} {
					if mode == AntialiasOff && mono != Bilevel || mode == AntialiasOn && mono == Bilevel {
						continue // the three images are of one drawing
					}
					whole, err := r.PageMono(v, 0, 1, mono)
					if err != nil {
						t.Fatal(f, err)
					}
					rows, bands := assembled(t, r, v, 0, 1, mono)
					if h := whole.Bounds().Dy(); bands != (h+47)/48 {
						t.Errorf("%s %s: %d bands of %d rows", f, v.ID, bands, h)
					}
					if n, at := samePixels(whole, rows); n != 0 {
						t.Errorf("%s %s (antialias %d, mono %d): %d pixels differ, the first at %v", f, v.ID, mode, mono, n, at)
					}
				}
			}
		}
	}
}

// TestRowsOfShadows draws shadows that reach across the bands, with
// their blur and offsets, in bands and whole, and checks that the
// images are the same: the bands take the margin the shadows need.
func TestRowsOfShadows(t *testing.T) {
	font, err := os.ReadFile("../../converter/pptx/testdata/fonts/MPLUS1p-Bold-subset.ttf")
	if err != nil {
		t.Fatal(err)
	}
	d := bdf.NewDocument()
	o := bdf.NewObject()
	fontRef := o.AddFont(bdf.EmbeddedFont(d.AddFont(font), 700, bdf.StyleNormal))
	o.FillColor(0xf0f0f0ff).FillRect(0, 0, 200, 300)
	// blurred, offset up and down, under a clip
	o.Save().ClipRect(5, 5, 190, 290)
	o.Save().Shadow(0x00000080, 12, 6, 9).FillColor(0xff0000ff).FillRect(20, 20, 60, 40).Restore()
	o.Save().Shadow(0x0000ff80, 20, -4, -30).FillColor(0x00ff00ff).FillRect(100, 60, 70, 50).Restore()
	o.Save().Shadow(0x000000ff, 0, 0, 40).FillColor(0x0000ffff).FillRect(30, 120, 20, 20).Restore()
	// a group with a shadow, holding a draw with its own shadow
	o.Save().Shadow(0x00000060, 10, 5, 12).GroupBegin(0.8, bdf.BlendSourceOver, 90, 130, 100, 100)
	o.FillColor(0xffff00ff).FillRect(100, 140, 60, 60)
	o.Save().Shadow(0xff000080, 6, 0, 20).FillColor(0x00ffffff).FillRect(120, 150, 30, 20).Restore()
	o.GroupEnd().Restore()
	// text (synthesized italic too) and a stroke with shadows
	o.Save().Shadow(0x000000a0, 5, 3, 3).Font(fontRef, 24).FillColor(0xff00ffff).FillText("Shadow", 20, 230, 0).Restore()
	o.Save().Shadow(0x008000a0, 4, 2, -6).StrokeColor(0x000000ff).Line(3, bdf.CapRound, bdf.JoinRound, 10)
	o.StrokePath(o.AddPath((&bdf.Path{}).MoveTo(20, 250).LineTo(120, 280).LineTo(60, 295))).Restore()
	o.Restore()
	// a blend mode that changes what it does not cover, after the bands
	// have been noted: it must not be skipped
	o.Save().ClipRect(150, 240, 40, 40).Blend(bdf.BlendSourceIn).FillColor(0x00ff00ff).FillRect(160, 250, 10, 10).Restore()
	h, _ := d.AddObject(o)
	v := d.NewView("v", bdf.ViewFixed, "")
	v.AddPage(200, 300, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	for _, rows := range []int{8, 50, 1000} {
		for _, mode := range []AntialiasMode{AntialiasOn, AntialiasOff} {
			var recs []*bandRecord
			bandHook = func(b *bandRecord) { recs = append(recs, b) }
			r := New(d, &Options{NoSystemFonts: true, Antialias: mode, BandRows: rows, ImageCacheBytes: -1})
			whole, err := r.Page(v, 0, 2)
			if err != nil {
				t.Fatal(err)
			}
			got, bands := assembled(t, r, v, 0, 2, 0)
			bandHook = nil
			if n, at := samePixels(whole, got); n != 0 {
				t.Errorf("bands of %d rows (antialias %d): %d pixels differ, the first at %v", rows, mode, n, at)
			}
			if len(recs) == 0 || recs[len(recs)-1].margin == 0 {
				t.Errorf("bands of %d rows: no margin was noted", rows)
			}
			if rows == 1000 && bands != 1 {
				t.Errorf("%d bands of 1000 rows", bands)
			}
			if len(r.Warnings()) != 0 {
				t.Errorf("warnings: %v", r.Warnings())
			}
		}
	}
}

// TestRowsSkipDraws checks that the bands after the first skip the draws
// that lie outside them, and none that reach into them.
func TestRowsSkipDraws(t *testing.T) {
	o := bdf.NewObject()
	for y := 0; y < 100; y++ {
		o.FillColor(bdf.Color(uint32(y)<<24|0xff)).FillRect(float32(y), float32(10*y), 20, 5)
	}
	o.Save().Blend(bdf.BlendSourceIn).FillRect(0, 0, 10, 10).Restore()
	d, v := onePage(o, 200, 1000)
	var recs []*bandRecord
	bandHook = func(b *bandRecord) { recs = append(recs, b) }
	defer func() { bandHook = nil }()
	r := New(d, &Options{NoSystemFonts: true, BandRows: 100, ImageCacheBytes: -1})
	whole, err := r.Page(v, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	got, bands := assembled(t, r, v, 0, 1, 0)
	if n, at := samePixels(whole, got); n != 0 || bands != 10 {
		t.Fatalf("%d bands, %d pixels differ, the first at %v", bands, n, at)
	}
	rec := recs[len(recs)-1]
	known, unknown := 0, 0
	for _, b := range rec.boxes {
		if b.known {
			known++
		} else {
			unknown++
		}
	}
	// the fills are known, the colours and the fill of the blend that
	// changes what it does not cover are not
	if known != 100 || unknown < 103 {
		t.Errorf("%d draws noted, %d instructions not", known, unknown)
	}
	// each band draws its 10 rectangles (one reaching from the band
	// before, a few rows tall), none of the rest
	for i, b := range rec.boxes {
		if !b.known || i >= 200 {
			continue
		}
		for band := 0; band < 10; band++ {
			canvas := image.Rect(0, 100*band, 200, 100*band+100)
			y := (i - 1) / 2 // the fills come after their colours
			want := 10*y+5 > canvas.Min.Y-1 && 10*y < canvas.Max.Y+1
			if skips := rec.skips(i, canvas); skips == want {
				t.Errorf("instruction %d (rectangle at %d) in band %d: skipped %v", i, 10*y, band, skips)
			}
		}
	}
}

// TestRowsEntryPoints checks PagesRows and the errors of PageRows.
func TestRowsEntryPoints(t *testing.T) {
	doc := openDoc(t, "../../testdata/pdf/chrome-doc.bdf")
	v := doc.Views[0]
	r := New(doc, &Options{FontDirs: testFonts, NoSystemFonts: true, BandRows: 64, ImageCacheBytes: -1})
	if err := r.PageRows(v, 0, 1, Mono(9), func(int, image.Image) error { return nil }); err == nil {
		t.Error("an unknown mono")
	}
	if err := r.PageRows(v, 0, 1, 0, nil); err == nil {
		t.Error("no callback")
	}
	if err := r.PageRows(v, len(v.Pages), 1, 0, func(int, image.Image) error { return nil }); !errors.Is(err, ErrNoPage) {
		t.Errorf("a page past the end: %v", err)
	}
	if err := r.PageRows(v, 0, 1e6, 0, func(int, image.Image) error { return nil }); !errors.Is(err, ErrTooLarge) {
		t.Errorf("a huge page: %v", err)
	}
	stop := errors.New("stop")
	n := 0
	if err := r.PageRows(v, 0, 1, Gray, func(int, image.Image) error { n++; return stop }); !errors.Is(err, stop) || n != 1 {
		t.Errorf("the error of the callback: %v after %d bands", err, n)
	}
	// PagesRows: the pages in selection order, each in bands from the top
	var seen []int
	last := -1
	err := r.PagesRows(stdcontext.Background(), v, []int{1, 0, 1}, 0.5, Bilevel, func(page, y0 int, rows image.Image) error {
		if y0 == 0 {
			seen = append(seen, page)
			last = 0
		} else if y0 != last {
			t.Errorf("page %d: band at %d after %d", page, y0, last)
		}
		last = rows.Bounds().Max.Y
		if _, ok := rows.(*image.Paletted); !ok {
			t.Errorf("a %T band", rows)
		}
		return nil
	})
	if err != nil || len(seen) != 3 || seen[0] != 1 || seen[1] != 0 || seen[2] != 1 {
		t.Errorf("pages %v: %v", seen, err)
	}
	if err := r.PagesRows(stdcontext.Background(), v, []int{0, 99}, 1, 0, func(int, int, image.Image) error { t.Fatal("drawn"); return nil }); !errors.Is(err, ErrNoPage) {
		t.Errorf("a selection past the end: %v", err)
	}
	ctx, cancel := stdcontext.WithCancel(stdcontext.Background())
	n = 0
	err = r.PagesRows(ctx, v, nil, 1, 0, func(int, int, image.Image) error { n++; cancel(); return nil })
	if !errors.Is(err, stdcontext.Canceled) || n != 1 {
		t.Errorf("cancelled: %v after %d bands", err, n)
	}
	if v.Kind != bdf.ViewSheet {
		sheet := &bdf.View{Kind: bdf.ViewSheet}
		if err := r.PagesRows(stdcontext.Background(), sheet, nil, 1, 0, func(int, int, image.Image) error { return nil }); err == nil {
			t.Error("a sheet")
		}
	}
	// a transparent background, and a region of the page outside it
	r = New(doc, &Options{FontDirs: testFonts, NoSystemFonts: true, BandRows: 48, Background: color.Transparent, ImageCacheBytes: -1})
	whole, err := r.Page(v, 0, 0.7)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := assembled(t, r, v, 0, 0.7, 0)
	if n, at := samePixels(whole, got); n != 0 {
		t.Errorf("transparent: %d pixels differ, the first at %v", n, at)
	}
}

// TestRowsInParallel draws bands with the passes split among goroutines
// and alone (see TestParallelPassesChangeNoPixel): the same pixels, and
// no race.
func TestRowsInParallel(t *testing.T) {
	defer func(n int) { parallelMinPixels = n }(parallelMinPixels)
	parallelMinPixels = 1 << 10
	for _, name := range []string{"pptx/features", "pdf/chrome-doc", "drawio/showcase"} {
		doc := openDoc(t, "../../testdata/"+name+".bdf")
		for _, v := range doc.Views {
			if v.Kind == bdf.ViewSheet || len(v.Pages) == 0 {
				continue
			}
			var out [][]byte
			for _, threads := range []int{1, 7} {
				r := New(doc, &Options{FontDirs: testFonts, NoSystemFonts: true, BandRows: 50, ImageCacheBytes: -1})
				r.threads = threads
				var buf bytes.Buffer
				if err := r.PageRows(v, 0, 1.2, 0, func(y0 int, rows image.Image) error {
					buf.Write(rows.(*image.RGBA).Pix)
					return nil
				}); err != nil {
					t.Fatal(name, err)
				}
				out = append(out, buf.Bytes())
			}
			if !bytes.Equal(out[0], out[1]) {
				t.Errorf("%s %s: the bands differ when the passes are split", name, v.ID)
			}
		}
	}
}
