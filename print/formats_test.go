package print

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/print/fakeprinter"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/raster"
)

func TestChooseFormat(t *testing.T) {
	printer := func(formats ...string) *ipp.Printer { return &ipp.Printer{Formats: formats} }
	for _, c := range []struct {
		formats []string
		asked   string
		want    string
	}{
		{[]string{FormatPWGRaster, FormatPDF, FormatURF}, "", FormatPDF},
		{[]string{FormatPCLm, FormatURF, FormatPWGRaster}, "auto", FormatPWGRaster},
		{[]string{"application/pclm", FormatURF}, "", FormatURF},
		{[]string{"application/pclm"}, "", FormatPCLm},
		{nil, "", FormatPDF},
		{[]string{FormatPDF, FormatURF}, "urf", FormatURF},
		{[]string{FormatPDF, FormatPCLm}, "PCLm", FormatPCLm},
		{nil, "image/urf", FormatURF},
		{[]string{FormatPWGRaster, FormatPostScript}, "", FormatPostScript},
		{[]string{FormatPDF, FormatPostScript}, "ps", FormatPostScript},
		{nil, "PostScript", FormatPostScript},
	} {
		got, err := chooseFormat(printer(c.formats...), c.asked)
		if err != nil || got != c.want {
			t.Errorf("%v asked %q: %q %v, want %q", c.formats, c.asked, got, err, c.want)
		}
	}
	for _, c := range []struct {
		formats []string
		asked   string
	}{{[]string{"image/jpeg"}, ""}, {[]string{FormatPDF}, "urf"}, {nil, "postscript-level-1"}} {
		if got, err := chooseFormat(printer(c.formats...), c.asked); err == nil {
			t.Errorf("%v asked %q: %q", c.formats, c.asked, got)
		}
	}
	if fallbackRaster(printer(FormatPDF, FormatPostScript, FormatURF)) != FormatURF || fallbackRaster(printer(FormatPDF, FormatPostScript)) != "" {
		t.Error("the raster of a printer, for what PDF leaves out")
	}
	if fallbackFormat(printer(FormatPDF, FormatURF, FormatPostScript)) != FormatPostScript || fallbackFormat(printer(FormatPDF)) != "" {
		t.Error("the format of a printer that refuses PDF")
	}
}

// TestPrintURF prints on a printer of Apple Raster only: gray at its
// resolution, the back of a sheet as urf-supported says (DM3: rotated).
func TestPrintURF(t *testing.T) {
	p := fakeprinter.New()
	p.Formats = []string{FormatURF}
	p.URF = []string{"V1.4", "CP1", "W8", "SRGB24", "RS300", "DM3"}
	c := serve(t, p)
	res, err := Print(context.Background(), c, testDoc(3), &Options{Ticket: Ticket{ColorMode: "monochrome", Sides: "two-sided-long-edge"}, Wait: time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	j := p.Jobs()[0]
	if res.Format != FormatURF || j.Format != FormatURF {
		t.Fatalf("sent as %s (%s)", res.Format, j.Format)
	}
	pages, err := raster.ReadURF(bytes.NewReader(j.Document(0)), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != 3 {
		t.Fatalf("%d pages", len(pages))
	}
	pg := pages[1]
	if pg.Type != raster.SGray8 || pg.XDPI != 300 || pg.Width != 2480 || pg.Height != 3507 || !pg.Duplex || pg.Tumble {
		t.Fatalf("header %+v", pg.Page)
	}
	at := func(rp *raster.ReadPage, xpt, ypt float64) uint8 {
		r, _, _, _ := rp.Image.At(int(xpt*300/72), int(ypt*300/72)).RGBA()
		return uint8(r >> 8)
	}
	// the back is turned a half turn: the bar at the left and top of the
	// page is at the right and bottom of the sheet
	if at(pages[0], 75, 120) != 0 || at(pages[1], a4w-75, a4h-120) != 0 || at(pages[1], 75, 120) != 255 {
		t.Errorf("front %d, back rotated %d, back unrotated %d", at(pages[0], 75, 120), at(pages[1], a4w-75, a4h-120), at(pages[1], 75, 120))
	}
}

// TestPrintPCLm prints on a printer of PCLm only: in colour, in strips of
// its height, compressed as it prefers.
func TestPrintPCLm(t *testing.T) {
	p := fakeprinter.New()
	p.Formats = []string{FormatPCLm}
	p.PCLmStripHeight = 32
	p.PCLmBackSide = "normal"
	p.PCLmResolutions = []ipp.Resolution{{X: 300, Y: 300, Units: ipp.UnitsDPI}}
	p.PCLmCompression = []string{"jpeg", "flate"}
	c := serve(t, p)
	res, err := Print(context.Background(), c, testDoc(2), &Options{Wait: time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	j := p.Jobs()[0]
	doc := j.Document(0)
	if res.Format != FormatPCLm || j.Format != FormatPCLm || !bytes.HasPrefix(doc, []byte("%PDF-1.7\n%PCLm 1.0\n")) {
		t.Fatalf("sent as %s (%s): %q", res.Format, j.Format, doc[:min(len(doc), 20)])
	}
	// an A4 page at 300 dpi is 3507 rows: 110 strips of 32 a page
	if n := bytes.Count(doc, []byte("/Subtype /Image")); n != 2*110 {
		t.Errorf("%d strips", n)
	}
	if !bytes.Contains(doc, []byte("/ColorSpace /DeviceRGB /BitsPerComponent 8 /Filter /DCTDecode")) || !bytes.Contains(doc, []byte("/MediaBox [0 0 595.2 841.68]")) {
		t.Errorf("strips or page: %q", doc[:min(len(doc), 600)])
	}
}

// TestPrintPostScript prints on a printer of PostScript: the pages placed
// on its paper, gray for a monochrome job, the sides asked in the file;
// a printer that refuses PDF is sent PostScript before raster pages.
func TestPrintPostScript(t *testing.T) {
	p := fakeprinter.New()
	p.Formats = []string{FormatPostScript, FormatPWGRaster}
	c := serve(t, p)
	doc := testDoc(3)
	res, err := Print(context.Background(), c, doc, &Options{Ticket: Ticket{ColorMode: "monochrome", Sides: "two-sided-long-edge", Media: "na_letter_8.5x11in", Resolution: 600}, Wait: time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	j := p.Jobs()[0]
	ps := string(j.Document(0))
	if res.Format != FormatPostScript || j.Format != FormatPostScript || res.Sheets != 3 || !strings.HasPrefix(ps, "%!PS-Adobe-3.0\n") {
		t.Fatalf("sent as %s (%s), %d sheets: %q", res.Format, j.Format, res.Sheets, ps[:min(len(ps), 20)])
	}
	for _, want := range []string{"%%Pages: 3\n", "<< /PageSize [612 792] >>", "<< /Duplex true /Tumble false >>", "\n0 g\n"} {
		if !strings.Contains(ps, want) {
			t.Errorf("the file lacks %q", want)
		}
	}
	if strings.Contains(ps, " rg\n") || attr(j, "printer-resolution") != "600x600dpi" || attr(j, "print-scaling") != "" {
		t.Errorf("colours, or attributes %v", j.Attrs)
	}
	// an A4 page on Letter paper: scaled to fit, centred
	if pages := ghostscript(t, ps); len(pages) == 3 {
		k := math.Min(612/a4w, 792/a4h)
		x, y := (612-a4w*k)/2+75*k, (792-a4h*k)/2+120*k
		if c := pages[0].RGBAAt(int(x), int(y)); c.R != 0 {
			t.Errorf("the bar of the first page is not at its place: %v", c)
		}
		if c := pages[0].RGBAAt(int(x)+60, int(y)); c.R != 255 {
			t.Errorf("the bar of the first page is too long: %v", c)
		}
		if c := pages[1].RGBAAt(int(x)+60, int(y)); c.R != 0 {
			t.Errorf("the bar of the second page is too short: %v", c)
		}
	}

	p = fakeprinter.New()
	p.Formats = []string{FormatPDF, FormatPostScript, FormatPWGRaster}
	p.Reject = []string{FormatPDF}
	c = serve(t, p)
	if res, err = Print(context.Background(), c, testDoc(1), &Options{Wait: time.Minute}); err != nil {
		t.Fatal(err)
	}
	if res.Format != FormatPostScript || p.Jobs()[0].Format != FormatPostScript || !strings.Contains(strings.Join(res.Warnings, "\n"), "refuses PDF: sent as PostScript") {
		t.Errorf("sent as %s: %v", res.Format, res.Warnings)
	}
}

// ghostscript draws the pages of a PostScript file at 72 dpi; without
// Ghostscript, it draws none.
func ghostscript(t *testing.T, ps string) []*image.RGBA {
	t.Helper()
	gs, err := exec.LookPath("gs")
	if err != nil {
		t.Log("Ghostscript (gs) is not installed: the pages are not drawn")
		return nil
	}
	dir := t.TempDir()
	cmd := exec.Command(gs, "-q", "-dSAFER", "-dBATCH", "-dNOPAUSE", "-sDEVICE=png16m", "-r72", "-sOutputFile="+filepath.Join(dir, "p%d.png"), "-")
	cmd.Stdin = strings.NewReader(ps)
	if out, err := cmd.CombinedOutput(); err != nil || len(out) > 0 {
		t.Fatalf("Ghostscript: %v\n%s", err, out)
	}
	var pages []*image.RGBA
	for n := 1; ; n++ {
		f, err := os.Open(filepath.Join(dir, fmt.Sprintf("p%d.png", n)))
		if err != nil {
			return pages
		}
		img, err := png.Decode(f)
		f.Close()
		if err != nil {
			t.Fatal(err)
		}
		rgba := image.NewRGBA(img.Bounds())
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
				rgba.Set(x, y, img.At(x, y))
			}
		}
		pages = append(pages, rgba)
	}
}

// TestPlacePS places pages on paper as raster pages are placed: scaled to
// fit, centred, landscape pages turned onto portrait paper.
func TestPlacePS(t *testing.T) {
	if PlacePS("no-such-paper", "", 0) != nil {
		t.Error("a placement on unknown paper")
	}
	apply := func(m [6]float64, x, y float64) [2]float64 {
		return [2]float64{math.Round(m[0]*x + m[2]*y + m[4]), math.Round(m[1]*x + m[3]*y + m[5])}
	}
	for _, c := range []struct {
		media       string
		scaling     raster.Scaling
		orientation raster.Orientation
		w, h        float64
		// where the top left and the bottom right corners of the page go
		tl, br [2]float64
	}{
		{"iso_a4_210x297mm", "", 0, 595.2756, 841.8898, [2]float64{0, 842}, [2]float64{595, 0}},
		{"iso_a4_210x297mm", "", 0, 841.8898, 595.2756, [2]float64{0, 0}, [2]float64{595, 842}}, // turned counterclockwise
		{"iso_a4_210x297mm", "", raster.ReversePortrait, 595.2756, 841.8898, [2]float64{595, 0}, [2]float64{0, 842}},
		{"iso_a4_210x297mm", raster.ScaleNone, raster.Portrait, 100, 100, [2]float64{248, 471}, [2]float64{348, 371}}, // centred
		{"iso_a4_210x297mm", raster.ScaleFit, raster.Portrait, 100, 100, [2]float64{0, 719}, [2]float64{595, 123}},
	} {
		w, h, m := PlacePS(c.media, c.scaling, c.orientation)(c.w, c.h)
		if math.Round(w) != 595 || math.Round(h) != 842 || apply(m, 0, 0) != c.tl || apply(m, c.w, c.h) != c.br {
			t.Errorf("%v %s %d: %v × %v, %v %v", c.media, c.scaling, c.orientation, w, h, apply(m, 0, 0), apply(m, c.w, c.h))
		}
	}
}
