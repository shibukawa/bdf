package printconv

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/shibukawa/bdf/converter/markdown"
	_ "github.com/shibukawa/bdf/converter/pdf"
	_ "github.com/shibukawa/bdf/converter/pptx"
	"github.com/shibukawa/bdf/print/raster"
	"github.com/shibukawa/bdf/print/spool"
)

func TestConverter(t *testing.T) {
	c := Converter{FontDirs: []string{filepath.Join("..", "..", "converter", "pptx", "testdata", "fonts")}, NoSystemFonts: true}
	ctx := context.Background()
	read := func(path string) *bytes.Reader {
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		return bytes.NewReader(b)
	}

	// an encrypted deck: its name tells it, or the encryption does
	deck := read(filepath.Join("..", "..", "converter", "internal", "offcrypto", "testdata", "agile.pptx"))
	if f, ok := c.Detect(deck, deck.Size(), "deck.pptx"); !ok || f != "pptx" {
		t.Fatalf("by name: %q %v", f, ok)
	}
	if f, ok := c.Detect(deck, deck.Size(), "upload"); !ok || f != EncryptedOffice {
		t.Fatalf("by its encryption: %q %v", f, ok)
	}
	const password = "パスワード🔑bdf"
	for _, format := range []string{"pptx", EncryptedOffice} {
		if _, _, err := c.Convert(ctx, deck, deck.Size(), format, "deck.pptx", ""); !errors.Is(err, spool.ErrPasswordRequired) {
			t.Fatalf("%s without the password: %v", format, err)
		}
		if _, _, err := c.Convert(ctx, deck, deck.Size(), format, "deck.pptx", "wrong"); !errors.Is(err, spool.ErrWrongPassword) {
			t.Fatalf("%s with a wrong password: %v", format, err)
		}
		doc, _, err := c.Convert(ctx, deck, deck.Size(), format, "deck.pptx", password)
		if err != nil {
			t.Fatal(err)
		}
		// what needed the password stays locked in the spool
		if doc.Lock == nil {
			t.Fatalf("%s: the converted deck is not locked", format)
		}
	}

	// Markdown: by its name, as pages
	md := bytes.NewReader([]byte("# Notes\n\nA line.\n"))
	f, ok := c.Detect(md, md.Size(), "notes.md")
	if !ok || f != "markdown" {
		t.Fatalf("markdown: %q %v", f, ok)
	}
	doc, _, err := c.Convert(ctx, md, md.Size(), f, "notes.md", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Views) == 0 || len(doc.Views[0].Pages) == 0 {
		t.Fatalf("markdown views %+v", doc.Views)
	}
	if _, ok := c.Detect(bytes.NewReader([]byte{0, 1, 2, 3}), 4, "blob"); ok {
		t.Fatal("bytes detected")
	}
	// a conversion stops waiting when asked
	cctx, cancel := context.WithCancel(ctx)
	cancel()
	if _, _, err := c.Convert(cctx, md, md.Size(), f, "notes.md", ""); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled: %v", err)
	}
}

// PCLm, the raster of Mopria, is a PDF of strips of image: the PDF
// converter reads it, in each of its compressions, as pages of the size
// the strips make.
func TestPCLm(t *testing.T) {
	c := Converter{NoSystemFonts: true}
	for _, comp := range []string{"flate", "jpeg", "rle"} {
		t.Run(comp, func(t *testing.T) {
			var b bytes.Buffer
			w, err := raster.NewPCLmWriter(&b, raster.PCLmOptions{Compression: comp})
			if err != nil {
				t.Fatal(err)
			}
			// two A6 pages at 150 dpi, a grey band across each
			for range 2 {
				pg := raster.Page{Width: 620, Height: 874, XDPI: 150, YDPI: 150, Type: raster.SRGB8, PageWidth: 298, PageHeight: 420}
				if err := w.StartPage(pg); err != nil {
					t.Fatal(err)
				}
				img := image.NewRGBA(image.Rect(0, 0, pg.Width, pg.Height))
				for y := range pg.Height {
					for x := range pg.Width {
						v := uint8(255)
						if y > 400 && y < 480 {
							v = 128
						}
						img.Set(x, y, color.RGBA{v, v, v, 255})
					}
				}
				if err := w.WriteRows(img); err != nil {
					t.Fatal(err)
				}
			}
			if err := w.Close(); err != nil {
				t.Fatal(err)
			}
			r := bytes.NewReader(b.Bytes())
			f, ok := c.Detect(r, r.Size(), "job")
			if !ok || f != "pdf" {
				t.Fatalf("detected %q %v", f, ok)
			}
			doc, _, err := c.Convert(context.Background(), r, r.Size(), f, "job", "")
			if err != nil {
				t.Fatal(err)
			}
			v := doc.Views[0]
			if len(v.Pages) != 2 {
				t.Fatalf("%d pages", len(v.Pages))
			}
			if w, h := float64(v.Pages[0].W), float64(v.Pages[0].H); math.Abs(w-298) > 1 || math.Abs(h-420) > 1 {
				t.Errorf("page %vx%v, want 298x420", w, h)
			}
		})
	}
}
