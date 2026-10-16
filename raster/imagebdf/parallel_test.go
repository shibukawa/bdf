package imagebdf

import (
	"bytes"
	stdcontext "context"
	"image"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// TestParallelPassesChangeNoPixel draws documents with the passes over
// the canvas split among goroutines (even small ones) and alone, and
// checks that the images are the same.
func TestParallelPassesChangeNoPixel(t *testing.T) {
	defer func(n int) { parallelMinPixels = n }(parallelMinPixels)
	parallelMinPixels = 1 << 10
	for _, name := range []string{"pdf/chrome-doc", "pptx/features", "drawio/showcase", "docx/basic", "gerber/board", "xlsx/features", "demo"} {
		doc := openDoc(t, "../../testdata/"+name+".bdf")
		for _, mode := range []AntialiasMode{AntialiasOn, AntialiasOff} {
			draw := func(threads int) (out [][]byte) {
				r := New(doc, &Options{FontDirs: testFonts, NoSystemFonts: true, Antialias: mode, ImageCacheBytes: -1, CanvasCacheBytes: 256 << 20})
				r.threads = threads
				add := func(img image.Image, err error) {
					if err != nil {
						t.Fatal(name, err)
					}
					switch img := img.(type) {
					case *image.RGBA:
						out = append(out, img.Pix)
					case *image.Gray:
						out = append(out, img.Pix)
					case *image.Paletted:
						out = append(out, img.Pix)
					}
				}
				for _, v := range doc.Views {
					switch v.Kind {
					case bdf.ViewSheet:
						add(r.Region(v, 0, bdf.Rect{X: 300.5, Y: 200.25, W: 500, H: 400}, 375, 300))
						add(r.RegionMono(v, 0, bdf.Rect{W: 500, H: 400}, 250, 200, Gray))
					case bdf.ViewScroll:
						add(r.Region(v, 0, bdf.Rect{X: 100, Y: 300, W: 400, H: 400}, 300, 300))
					default:
						if len(v.Pages) == 0 {
							continue
						}
						add(r.Page(v, 0, 1))
						add(r.Page(v, 0, 2.5)) // the same canvas again, larger
						add(r.PageMono(v, 0, 1.5, Bilevel))
						p := v.Pages[len(v.Pages)-1]
						add(r.Region(v, len(v.Pages)-1, bdf.Rect{X: p.W / 3, Y: p.H / 4, W: p.W / 2, H: p.H / 2}, 200, 200))
					}
				}
				return out
			}
			alone, split := draw(1), draw(7)
			for i := range alone {
				if !bytes.Equal(alone[i], split[i]) {
					t.Errorf("%s (antialias %d): image %d differs", name, mode, i)
				}
			}
		}
	}
}

// TestSplitRowsPanics checks that a panic of a band is the drawing's.
func TestSplitRowsPanics(t *testing.T) {
	defer func(n int) { parallelMinPixels = n }(parallelMinPixels)
	parallelMinPixels = 1
	defer func() {
		if p := recover(); p == nil || !strings.Contains(p.(string), "band 2") {
			t.Fatalf("recovered %v", p)
		}
	}()
	splitRows(4, 0, 8, 1, func(y0, y1, k int) {
		if k == 2 {
			panic("band 2 failed")
		}
	})
}

// TestCanvasBudget checks that the canvases and masks kept stay within
// CanvasCacheBytes, and that a batch keeps one canvas a worker.
func TestCanvasBudget(t *testing.T) {
	doc := openDoc(t, "../../testdata/pdf/chrome-doc.bdf")
	v := doc.Views[0]
	r := New(doc, &Options{FontDirs: testFonts, NoSystemFonts: true, ImageCacheBytes: -1, CanvasCacheBytes: -1})
	if _, err := r.Page(v, 0, 2); err != nil {
		t.Fatal(err)
	}
	if r.freeBytes != 0 || r.freeMaskBytes != 0 || len(r.free) != 0 || len(r.freeMasks) != 0 {
		t.Fatalf("kept %d bytes of canvases and %d of masks with a negative budget", r.freeBytes, r.freeMaskBytes)
	}
	r = New(doc, &Options{FontDirs: testFonts, NoSystemFonts: true, ImageCacheBytes: -1, CanvasCacheBytes: 1 << 20})
	for i := 0; i < 3; i++ {
		if _, err := r.Page(v, 0, 2); err != nil {
			t.Fatal(err)
		}
	}
	if r.freeBytes > 1<<20 || r.freeMaskBytes > 1<<20 {
		t.Fatalf("kept %d bytes of canvases and %d of masks past the budget", r.freeBytes, r.freeMaskBytes)
	}
	// a batch keeps the canvas of a page while it lasts, and trims after
	r = New(doc, &Options{FontDirs: testFonts, NoSystemFonts: true, ImageCacheBytes: -1})
	pages := []int{0, 0, 0}
	var pix []byte
	err := r.Pages(stdcontext.Background(), v, pages, 2, 1, func(page int, img *image.RGBA) error {
		if pix == nil {
			pix = img.Pix
		} else if !bytes.Equal(pix, img.Pix) {
			t.Errorf("page %d differs", page)
		}
		w := r.pageWorkers[0]
		if n := 16 * img.Rect.Dx() * img.Rect.Dy(); w.freeBudget() < n {
			t.Errorf("a worker's budget is %d bytes in a batch of pages of %d", w.freeBudget(), n)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if w := r.pageWorkers[0]; w.freeBytes > DefaultCanvasCacheBytes || w.freeBudget() != DefaultCanvasCacheBytes {
		t.Fatalf("after the batch a worker keeps %d bytes within %d", w.freeBytes, w.freeBudget())
	}
}

// TestLazyCanvasesReadNothingUnwritten draws documents on reused canvases
// poisoned with values no drawing leaves, with the shortcuts and without
// (where the canvases are written whole), and checks that the images are
// the same: a pixel read before it is written would show.
func TestLazyCanvasesReadNothingUnwritten(t *testing.T) {
	defer func() { poisonSurfaces = false }()
	defer func(n int) { parallelMinPixels = n }(parallelMinPixels)
	parallelMinPixels = 1 << 10
	for _, name := range []string{"pdf/chrome-doc", "pptx/features", "drawio/showcase", "xlsx/features", "pdf/chrome-masks", "demo"} {
		doc := openDoc(t, "../../testdata/"+name+".bdf")
		for _, mode := range []AntialiasMode{AntialiasOn, AntialiasOff} {
			draw := func(pl bool) (out [][]byte) {
				plain, poisonSurfaces = pl, !pl
				defer func() { plain, poisonSurfaces = false, false }()
				r := New(doc, &Options{FontDirs: testFonts, NoSystemFonts: true, Antialias: mode, ImageCacheBytes: -1, CanvasCacheBytes: 1 << 30, BandRows: 40})
				add := func(img image.Image, err error) {
					if err != nil {
						t.Fatal(name, err)
					}
					switch img := img.(type) {
					case *image.RGBA:
						out = append(out, img.Pix)
					case *image.Gray:
						out = append(out, img.Pix)
					case *image.Paletted:
						out = append(out, img.Pix)
					}
				}
				for _, v := range doc.Views {
					switch v.Kind {
					case bdf.ViewSheet:
						// twice: the second drawing reuses the canvas
						add(r.Region(v, 0, bdf.Rect{X: 300.5, Y: 200.25, W: 500, H: 400}, 375, 300))
						add(r.RegionMono(v, 0, bdf.Rect{X: 1, Y: 2, W: 500, H: 400}, 250, 200, Gray))
					case bdf.ViewScroll:
						add(r.Region(v, 0, bdf.Rect{X: 100, Y: 300, W: 400, H: 400}, 300, 300))
						add(r.Region(v, 0, bdf.Rect{X: 0, Y: 0, W: 400, H: 400}, 200, 200))
					default:
						if len(v.Pages) == 0 {
							continue
						}
						add(r.Page(v, 0, 1))
						add(r.PageMono(v, 0, 1.5, Bilevel))
						add(r.Page(v, 0, 0.7))
						var bands []byte
						if err := r.PageRows(v, 0, 1.2, Gray, func(y0 int, rows image.Image) error {
							bands = append(bands, rows.(*image.Gray).Pix...)
							return nil
						}); err != nil {
							t.Fatal(name, err)
						}
						out = append(out, bands)
					}
				}
				return out
			}
			fast, slow := draw(false), draw(true)
			for i := range fast {
				if !bytes.Equal(fast[i], slow[i]) {
					n := 0
					for k := range fast[i] {
						if k < len(slow[i]) && fast[i][k] != slow[i][k] {
							n++
						}
					}
					t.Errorf("%s (antialias %d): image %d differs in %d bytes", name, mode, i, n)
				}
			}
		}
	}
}
