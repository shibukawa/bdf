package print

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/pdf"
	"github.com/shibukawa/bdf/print/fakeprinter"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/raster"
)

// a4 is the size of A4 in points.
const a4w, a4h = 595.276, 841.89

// testDoc is a document of n A4 pages, each with a background layer (a
// gray frame) and a body (a black bar whose width tells the page).
func testDoc(n int) *bdf.Document {
	d := bdf.NewDocument()
	bg := bdf.NewObject()
	bg.FillColor(0xc0c0c0ff).FillRect(20, 20, a4w-40, 10)
	bgh, _ := d.AddObject(bg)
	v := d.NewView("pages", bdf.ViewFixed, "")
	for i := 0; i < n; i++ {
		o := bdf.NewObject()
		o.FillColor(0x000000ff).FillRect(50, 100, float32(50*(i+1)), 40)
		h, _ := d.AddObject(o)
		v.AddPage(a4w, a4h, bdf.Layer{Role: bdf.RoleBackground, Obj: bgh}, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	}
	return d
}

func serve(t *testing.T, p *fakeprinter.Printer) *ipp.Client {
	t.Helper()
	return &ipp.Client{URL: fakeprinter.Start(t, p), User: "tester"}
}

func attr(j *fakeprinter.Job, name string) string {
	for _, a := range j.Attrs {
		if a.Name == name {
			return a.String()
		}
	}
	return ""
}

func TestPrintPDF(t *testing.T) {
	p := fakeprinter.New()
	c := serve(t, p)
	var events []string
	res, err := Print(context.Background(), c, testDoc(3), &Options{
		Ticket: Ticket{JobName: "invoice.bdf", Copies: 2, Sides: "two-sided-long-edge", ColorMode: "monochrome", Pages: []int{1, 3}},
		Wait:   time.Minute, Progress: func(e Event) { events = append(events, e.Kind) }})
	if err != nil {
		t.Fatal(err)
	}
	if res.Format != FormatPDF || res.Media != "iso_a4_210x297mm" || res.Sheets != 2 || res.Job.State != ipp.JobCompleted || len(res.Unsupported) != 0 {
		t.Fatalf("result: format %s media %s sheets %d job %v unsupported %v", res.Format, res.Media, res.Sheets, res.Job, res.Unsupported)
	}
	if strings.Join(events, ",") != "created,sent" {
		t.Errorf("events: %v", events)
	}
	j := p.Jobs()[0]
	if j.Name != "invoice.bdf" || j.Format != FormatPDF || attr(j, "copies") != "2" || attr(j, "sides") != "two-sided-long-edge" ||
		attr(j, "media") != "iso_a4_210x297mm" || attr(j, "print-color-mode") != "monochrome" {
		t.Fatalf("job %q %s: %v", j.Name, j.Format, j.Attrs)
	}
	doc := j.Document(0)
	back, err := pdf.Convert(bytes.NewReader(doc), nil)
	if err != nil {
		t.Fatal(err)
	}
	if n := len(back.Doc.Views[0].Pages); n != 2 {
		t.Fatalf("the PDF has %d pages", n)
	}
	if res.Bytes != int64(len(doc)) {
		t.Errorf("bytes: %d sent, %d received", res.Bytes, len(doc))
	}
}

func rasterPrinter() *fakeprinter.Printer {
	p := fakeprinter.New()
	p.Formats = []string{FormatPWGRaster, "image/jpeg"}
	return p
}

func readRaster(t *testing.T, j *fakeprinter.Job) []*raster.ReadPage {
	t.Helper()
	pages, err := raster.Read(bytes.NewReader(j.Document(0)), 0)
	if err != nil {
		t.Fatal(err)
	}
	return pages
}

func TestPrintRaster(t *testing.T) {
	p := rasterPrinter()
	c := serve(t, p)
	res, err := Print(context.Background(), c, testDoc(2), &Options{Ticket: Ticket{ColorMode: "monochrome"}, Wait: time.Minute, BandRows: 64})
	if err != nil {
		t.Fatal(err)
	}
	if res.Format != FormatPWGRaster || res.Sheets != 2 {
		t.Fatalf("result: %s %d", res.Format, res.Sheets)
	}
	pages := readRaster(t, p.Jobs()[0])
	if len(pages) != 2 {
		t.Fatalf("%d pages", len(pages))
	}
	pg := pages[0]
	// monochrome: gray at 600 dpi on A4
	if pg.Type != raster.SGray8 || pg.XDPI != 600 || pg.Width != 4960 || pg.Height != 7015 || pg.PageSizeName != "iso_a4_210x297mm" ||
		pg.PageWidth != 595 || pg.PageHeight != 842 || pg.TotalPages != 2 {
		t.Fatalf("header: %+v", pg.Page)
	}
	// the bar of page 2 is twice as wide as page 1's: 100 points at 600 dpi
	at := func(rp *raster.ReadPage, xpt, ypt float64) uint8 {
		x, y := int(xpt*600/72), int(ypt*600/72)
		r, _, _, _ := rp.Image.At(x, y).RGBA()
		return uint8(r >> 8)
	}
	if at(pages[0], 120, 120) != 255 || at(pages[1], 120, 120) != 0 || at(pages[0], 75, 120) != 0 {
		t.Errorf("the bars: page 1 at 120pt %d, page 2 %d, page 1 at 75pt %d", at(pages[0], 120, 120), at(pages[1], 120, 120), at(pages[0], 75, 120))
	}
	if at(pages[0], 100, 25) == 255 {
		t.Errorf("the background layer is not drawn")
	}
}

func TestPrintRasterOptions(t *testing.T) {
	p := rasterPrinter()
	p.Copies = ipp.Range{Lo: 1, Hi: 1} // the printer makes no copies
	p.SheetBack = "flipped"
	c := serve(t, p)
	res, err := Print(context.Background(), c, testDoc(3), &Options{
		Ticket:    Ticket{Copies: 2, Sides: "two-sided-long-edge", ColorMode: "monochrome", Bilevel: true, Resolution: 300},
		SkipRoles: []string{bdf.RoleBackground}, Wait: time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	j := p.Jobs()[0]
	if attr(j, "copies") != "" {
		t.Errorf("copies were asked of a printer that makes none")
	}
	pages := readRaster(t, j)
	// 3 pages twice, each copy starting on a front: 3 + blank + 3
	if len(pages) != 7 || res.Sheets != 7 {
		t.Fatalf("%d pages, %d sheets", len(pages), res.Sheets)
	}
	for i, pg := range pages {
		back := i%2 == 1
		if pg.Type != raster.Black1 || pg.XDPI != 300 || !pg.Duplex || pg.Tumble || (pg.Feed == -1) != back || pg.CrossFeed != 1 {
			t.Errorf("page %d: %+v", i+1, pg.Page)
		}
	}
	at := func(rp *raster.ReadPage, xpt, ypt float64) uint8 {
		x, y := int(xpt*300/72), int(ypt*300/72)
		r, _, _, _ := rp.Image.At(x, y).RGBA()
		return uint8(r >> 8)
	}
	// the back of a sheet is upside down (flipped), the background left out
	if at(pages[0], 75, 120) != 0 || at(pages[1], 75, a4h-120) != 0 || at(pages[1], 75, 120) != 255 || at(pages[0], 100, 25) != 255 {
		t.Errorf("front %d, back flipped %d, back unflipped %d, background %d", at(pages[0], 75, 120), at(pages[1], 75, a4h-120), at(pages[1], 75, 120), at(pages[0], 100, 25))
	}
	if blank := pages[3]; at(blank, 75, 120) != 255 {
		t.Errorf("the sheet between the copies is not blank")
	}
}

func TestPrintFallsBackToRaster(t *testing.T) {
	p := fakeprinter.New()
	p.Reject = []string{FormatPDF}
	c := serve(t, p)
	res, err := Print(context.Background(), c, testDoc(1), &Options{Wait: time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	if res.Format != FormatPWGRaster || p.Jobs()[0].Format != FormatPWGRaster {
		t.Fatalf("format %s", res.Format)
	}
	if !strings.Contains(strings.Join(res.Warnings, "\n"), "refuses PDF") {
		t.Errorf("warnings: %v", res.Warnings)
	}
	pg := readRaster(t, p.Jobs()[0])[0]
	if pg.Type != raster.SRGB8 || pg.XDPI != 300 {
		t.Errorf("colour by default: %+v", pg.Page)
	}
}

func TestPrintUnsupported(t *testing.T) {
	p := fakeprinter.New()
	c := serve(t, p)
	tk := Ticket{Sides: "two-sided-long-edge", MediaSource: "tray-9", Finishings: []int{4}}
	res, err := Print(context.Background(), c, testDoc(1), &Options{Ticket: tk})
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(res.Unsupported, ";"); got != "media-source=tray-9;finishings=4" {
		t.Errorf("unsupported: %q", got)
	}
	j := p.Jobs()[0]
	if attr(j, "finishings") != "" || attr(j, "sides") != "two-sided-long-edge" || !strings.Contains(attr(j, "media-col"), "x-dimension=21000") {
		t.Errorf("attributes: %v", j.Attrs)
	}
	if _, err := Print(context.Background(), c, testDoc(1), &Options{Ticket: tk, Strict: true}); !errors.Is(err, ErrUnsupported) {
		t.Errorf("strict: %v", err)
	}
	if _, err := Print(context.Background(), c, testDoc(1), &Options{Ticket: Ticket{Sides: "both"}}); err == nil {
		t.Error("a ticket of no sides")
	}
	if _, err := Print(context.Background(), c, testDoc(1), &Options{Ticket: Ticket{Pages: []int{5}}}); !errors.Is(err, ErrNoPages) {
		t.Errorf("no page selected: %v", err)
	}
	if _, err := Print(context.Background(), c, testDoc(1), &Options{Format: "urf"}); err == nil {
		t.Error("URF")
	}
}

func TestPrintStoppedAbortedCanceled(t *testing.T) {
	p := fakeprinter.New()
	c := serve(t, p)
	p.State, p.Reasons = ipp.PrinterStopped, []string{"media-empty-error"}
	if _, err := Print(context.Background(), c, testDoc(1), nil); !errors.Is(err, ErrPrinterStopped) || !strings.Contains(err.Error(), "media-empty") {
		t.Fatalf("stopped: %v", err)
	}
	p.State, p.Reasons = ipp.PrinterIdle, []string{"none"}
	p.Processing = time.Hour
	_, err := Print(context.Background(), c, testDoc(1), &Options{Wait: time.Minute, Progress: func(e Event) {
		if e.Kind == "sent" {
			p.Abort(e.Job.ID)
		}
	}})
	if !errors.Is(err, ErrAborted) {
		t.Fatalf("aborted: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	_, err = Print(ctx, c, testDoc(1), &Options{Wait: -1, Progress: func(e Event) {
		if e.Kind == "sent" {
			go func() { time.Sleep(100 * time.Millisecond); cancel() }()
		}
	}})
	if !errors.Is(err, ErrCanceled) {
		t.Fatalf("canceled: %v", err)
	}
	last := p.Jobs()[len(p.Jobs())-1]
	if p.JobState(last) != ipp.JobCanceled {
		t.Fatalf("the job at the printer is %v", p.JobState(last))
	}
	_, err = Print(context.Background(), c, testDoc(1), &Options{Wait: 300 * time.Millisecond})
	if !errors.Is(err, ErrTimeout) {
		t.Fatalf("timeout: %v", err)
	}
}

func TestOrder(t *testing.T) {
	p := &ipp.Printer{Copies: ipp.Range{Lo: 1, Hi: 1}}
	for _, c := range []struct {
		t    Ticket
		want []int
	}{
		{Ticket{Copies: 2}, []int{0, 1, 2, 0, 1, 2}},
		{Ticket{Copies: 2, Collate: "uncollated"}, []int{0, 0, 1, 1, 2, 2}},
		{Ticket{Copies: 2, Sides: "two-sided-long-edge"}, []int{0, 1, 2, -1, 0, 1, 2}},
		{Ticket{Copies: 1, Sides: "two-sided-long-edge"}, []int{0, 1, 2}},
	} {
		j := &job{p: p, t: &c.t, pages: []int{0, 1, 2}}
		if got := j.order(); !equal(got, c.want) {
			t.Errorf("%+v: %v, want %v", c.t, got, c.want)
		}
	}
	if got := (&job{p: &ipp.Printer{Copies: ipp.Range{Lo: 1, Hi: 99}}, t: &Ticket{Copies: 5}, pages: []int{1}}).order(); !equal(got, []int{1}) {
		t.Errorf("the printer makes the copies: %v", got)
	}
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
