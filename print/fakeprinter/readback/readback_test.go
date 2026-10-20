package readback_test

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/draw"
	"os/exec"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/fakeprinter"
	"github.com/shibukawa/bdf/print/fakeprinter/readback"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/pdl"
	"github.com/shibukawa/bdf/print/raster"
	"github.com/shibukawa/bdf/print/raw"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

const a4w, a4h = 595.2756, 841.8898

// testDoc is pages of size w × h of black and grey bars, a bar wider on
// each page.
func testDoc(n int, w, h float32) (*bdf.Document, *bdf.View) {
	d := bdf.NewDocument()
	v := d.NewView("pages", bdf.ViewFixed, "")
	for i := 0; i < n; i++ {
		o := bdf.NewObject()
		o.FillColor(0x000000ff).FillRect(w/10, h/5, w/10*float32(i+2), h/8)
		o.FillColor(0x808080ff).FillRect(w/10, h/2, w/2, h/10)
		hh, _ := d.AddObject(o)
		v.AddPage(w, h, bdf.Layer{Role: bdf.RoleBody, Obj: hh})
	}
	return d, v
}

func gray(c color.Color) int {
	return int(color.GrayModel.Convert(c).(color.Gray).Y)
}

// differ counts the pixels of two images that differ by more than tol in
// grey, over the smaller of the two; the sizes must be within a pixel.
func differ(t *testing.T, name string, got, want image.Image, tol int) int {
	t.Helper()
	g, w := got.Bounds(), want.Bounds()
	if abs(g.Dx()-w.Dx()) > 1 || abs(g.Dy()-w.Dy()) > 1 {
		t.Fatalf("%s: %v, want %v", name, g, w)
	}
	n := 0
	for y := 0; y < min(g.Dy(), w.Dy()); y++ {
		for x := 0; x < min(g.Dx(), w.Dx()); x++ {
			if abs(gray(got.At(g.Min.X+x, g.Min.Y+y))-gray(want.At(w.Min.X+x, w.Min.Y+y))) > tol {
				n++
			}
		}
	}
	return n
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// printIPP prints a document on a preset reached over IPP and reads back
// its job.
func printIPP(t *testing.T, preset string, d *bdf.Document, o *print.Options) (*readback.Result, *fakeprinter.Job) {
	t.Helper()
	ps, _ := fakeprinter.LookupPreset(preset)
	p := ps.New(preset)
	c := &ipp.Client{URL: fakeprinter.Start(t, p), User: "tester"}
	o.Wait = time.Minute
	o.NoSystemFonts = true
	if _, err := print.Print(context.Background(), c, d, o); err != nil {
		t.Fatal(err)
	}
	j := p.Jobs()[0]
	res, err := readback.Job(p, j, nil)
	if err != nil {
		t.Fatal(err)
	}
	return res, j
}

// TestRasters prints PWG Raster and Apple Raster: the pages come back
// as imagebdf places them on the sheet, pixel for pixel.
func TestRasters(t *testing.T) {
	d, v := testDoc(2, a4w, a4h)
	r := imagebdf.New(d, &imagebdf.Options{NoSystemFonts: true, ImageCacheBytes: -1})
	m, _ := print.ParseMedia("iso_a4_210x297mm")
	sheet := raster.Sheet{Width: m.Width, Height: m.Height, DPI: 300}
	for _, c := range []struct{ preset, format string }{{"pwg-only", "pwg-raster"}, {"airprint-only", "urf"}} {
		res, j := printIPP(t, c.preset, d, &print.Options{Format: c.format, Ticket: print.Ticket{Resolution: 300, ColorMode: "monochrome", Media: "iso_a4_210x297mm"}})
		if res.Format != j.Format || len(res.Pages) != 2 {
			t.Fatalf("%s: %s, %d pages", c.preset, res.Format, len(res.Pages))
		}
		for i, pg := range res.Pages {
			if pg.DPI != 300 {
				t.Errorf("%s page %d: %d dpi", c.preset, i+1, pg.DPI)
			}
			pl, err := raster.Fit(float64(v.Pages[i].W), float64(v.Pages[i].H), sheet, "", 0)
			if err != nil {
				t.Fatal(err)
			}
			pw, ph := sheet.Pixels()
			want := image.NewGray(image.Rect(0, 0, pw, ph))
			if err := r.PlaceRows(v, i, pl, imagebdf.Gray, func(y0 int, rows image.Image) error {
				b := rows.Bounds()
				draw.Draw(want, image.Rect(0, y0, b.Dx(), y0+b.Dy()), rows, b.Min, draw.Src)
				return nil
			}); err != nil {
				t.Fatal(err)
			}
			if n := differ(t, c.preset, pg.Image, want, 0); n != 0 {
				t.Errorf("%s page %d: %d pixels differ", c.preset, i+1, n)
			}
		}
	}
}

// TestVectors prints PDF, PostScript and PCLm: drawn back at 150 dpi they
// are imagebdf's pages but for the edges of their shapes.
func TestVectors(t *testing.T) {
	d, v := testDoc(2, a4w, a4h)
	r := imagebdf.New(d, &imagebdf.Options{NoSystemFonts: true})
	cases := []struct{ preset, format string }{{"office-mfp", "pdf"}, {"pclm-inkjet", "pclm"}}
	if _, err := exec.LookPath("gs"); err == nil {
		cases = append(cases, struct{ preset, format string }{"ps-office", "ps"})
	}
	for _, c := range cases {
		res, _ := printIPP(t, c.preset, d, &print.Options{Format: c.format, Ticket: print.Ticket{Media: "iso_a4_210x297mm"}})
		if len(res.Pages) != 2 {
			t.Fatalf("%s: %d pages", c.format, len(res.Pages))
		}
		for i, pg := range res.Pages {
			want, _ := r.PageDPI(v, i, float64(pg.DPI))
			all := want.Bounds().Dx() * want.Bounds().Dy()
			// the edges of the bars, a pixel wide, may differ
			if n := differ(t, c.format, pg.Image, want, 48); n > all/200 {
				t.Errorf("%s page %d: %d pixels of %d differ", c.format, i+1, n, all)
			}
		}
	}
}

// TestLabels prints labels, receipts and forms over 9100, LPD and a
// serial line: what the fake reads back is what print/pdl wrote for the
// same printer, dot for dot, and its barcodes are the printer's own.
func TestLabels(t *testing.T) {
	for _, c := range []struct {
		preset, profile, transport string
		w, h                       float32
	}{
		{"zebra-zd421-203", "zd421-203", "socket", 4 * 72, 3 * 72},
		{"sato-cl4nx", "cl4nx-203", "socket", 4 * 72, 3 * 72},
		{"tec-bex4t1", "b-ex4t1-203", "lpd", 4 * 72, 3 * 72},
		{"epson-tm-t88", "tm-t88", "socket", 72 * 72 / 25.4, 200},
		{"epson-vp-f2000", "vp-f2000", "serial", 8 * 72, 5.5 * 72},
	} {
		t.Run(c.preset, func(t *testing.T) {
			d, v := testDoc(2, c.w, c.h)
			ps, _ := fakeprinter.LookupPreset(c.preset)
			p := ps.New(c.preset)
			var url string
			switch c.transport {
			case "socket":
				url = fakeprinter.StartSocket(t, p).URL + "?profile=" + c.profile
			case "lpd":
				url = fakeprinter.StartLPD(t, p, "labels").URL + "?profile=" + c.profile
			case "serial":
				url = "serial://" + fakeprinter.StartSerial(t, p) + "?profile=" + c.profile
			}
			rp, err := raw.Parse(url)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := print.PrintRaw(context.Background(), rp, d, &print.Options{NoSystemFonts: true}); err != nil {
				t.Fatal(err)
			}
			for range 100 {
				if len(p.Jobs()) > 0 {
					break
				}
				time.Sleep(20 * time.Millisecond)
			}
			jobs := p.Jobs()
			if len(jobs) != 1 || jobs[0].Transport != c.transport || jobs[0].Format != p.Language {
				t.Fatalf("%d jobs: %+v", len(jobs), jobs)
			}
			got, err := readback.Job(p, jobs[0], nil)
			if err != nil {
				t.Fatal(err)
			}
			prof, _ := pdl.LookupProfile(c.profile)
			var buf bytes.Buffer
			if _, err := pdl.Write(&buf, d, v, &pdl.Options{Profile: prof, Copies: 1}); err != nil {
				t.Fatal(err)
			}
			want, err := readback.Read(buf.Bytes(), &readback.Options{Format: prof.Language, Width: prof.Width, DPI: prof.DPI})
			if err != nil {
				t.Fatal(err)
			}
			if len(got.Pages) != len(want.Pages) || len(got.Pages) == 0 {
				t.Fatalf("%d pages, want %d", len(got.Pages), len(want.Pages))
			}
			for i := range got.Pages {
				if n := differ(t, c.preset, got.Pages[i].Image, want.Pages[i].Image, 0); n != 0 {
					t.Errorf("page %d: %d dots differ", i+1, n)
				}
				if black(got.Pages[i].Image) == 0 {
					t.Errorf("page %d is blank", i+1)
				}
			}
		})
	}
}

func black(img image.Image) int {
	n := 0
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if gray(img.At(x, y)) < 128 {
				n++
			}
		}
	}
	return n
}
