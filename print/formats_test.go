package print

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/ipp/ipptest"
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
	if fallbackRaster(printer(FormatPDF, FormatURF)) != FormatURF || fallbackRaster(printer(FormatPDF)) != "" {
		t.Error("the raster of a printer that refuses PDF")
	}
}

// TestPrintURF prints on a printer of Apple Raster only: gray at its
// resolution, the back of a sheet as urf-supported says (DM3: rotated).
func TestPrintURF(t *testing.T) {
	p := ipptest.New()
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
	p := ipptest.New()
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
