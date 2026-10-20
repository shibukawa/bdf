package pdl

import (
	"fmt"
	"image"
	"math"
	"sync"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/image/barcode"
	"github.com/shibukawa/bdf/print/raster"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

// Bitmap is rows of dots, black or white: a row of W dots is Stride
// bytes, the first dot in the high bit of the first byte, 1 for black.
// The dots past W in the last byte of a row are 0.
type Bitmap struct {
	W, H   int
	Stride int
	Bits   []byte
}

// NewBitmap returns a white bitmap of w × h dots.
func NewBitmap(w, h int) *Bitmap {
	s := (w + 7) / 8
	return &Bitmap{W: w, H: h, Stride: s, Bits: make([]byte, s*h)}
}

// Black reports whether the dot at x, y is black.
func (b *Bitmap) Black(x, y int) bool {
	if x < 0 || y < 0 || x >= b.W || y >= b.H {
		return false
	}
	return b.Bits[y*b.Stride+x/8]&(0x80>>(x&7)) != 0
}

// Row returns the bytes of a row.
func (b *Bitmap) Row(y int) []byte { return b.Bits[y*b.Stride : (y+1)*b.Stride] }

// Rows returns the bitmap of rows y0 to y1, sharing the bytes.
func (b *Bitmap) Rows(y0, y1 int) *Bitmap {
	return &Bitmap{W: b.W, H: y1 - y0, Stride: b.Stride, Bits: b.Bits[y0*b.Stride : y1*b.Stride]}
}

// blank reports whether rows y0 to y1 have no black dot.
func (b *Bitmap) blank(y0, y1 int) bool {
	for _, c := range b.Bits[y0*b.Stride : y1*b.Stride] {
		if c != 0 {
			return false
		}
	}
	return true
}

// pack writes the rows of a bilevel image (black at index 0) into the
// bitmap from row y.
func (b *Bitmap) pack(y int, rows image.Image) {
	p, ok := rows.(*image.Paletted)
	r := rows.Bounds()
	for j := 0; j < r.Dy() && y+j < b.H; j++ {
		row := b.Row(y + j)
		for x := 0; x < r.Dx() && x < b.W; x++ {
			var black bool
			if ok {
				black = p.Pix[j*p.Stride+x] == 0
			} else {
				c, _, _, _ := rows.At(r.Min.X+x, r.Min.Y+j).RGBA()
				black = c < 0x8000
			}
			if black {
				row[x/8] |= 0x80 >> (x & 7)
			}
		}
	}
}

// Symbol is a QR code or a barcode that the printer prints with its own
// commands in place of the drawing of a page.
type Symbol struct {
	Mark   *barcode.Mark
	Symbol *barcode.Symbol
	// X, Y, W and H are the box of the symbol and its quiet zone on the
	// label, in dots, as imagebdf.Barcode has it, and Turns its quarter
	// turns clockwise.
	X, Y, W, H float64
	Turns      int
}

// Module returns the width of a module in whole dots: as many as fit in
// the box (and, for a QR code, in its height), 0 when not even one does.
func (s *Symbol) Module() int {
	m := math.Floor(s.W / float64(s.Symbol.W))
	if !s.Symbol.Linear {
		m = math.Min(m, math.Floor(s.H/float64(s.Symbol.H)))
	}
	return int(m)
}

// Origin returns the top left corner of the symbol's bars (or of the QR
// code without its quiet zone), upright, in dots of the label, placed as
// AlignBarcodes places them: the modules of Module() dots in the middle
// of the box. Turned symbols are placed in their box turned; the corner is
// the one the printer's commands place (the top left of the turned
// symbol on the label), and w, h its size on the label.
func (s *Symbol) Origin() (x, y, w, h int) {
	m := s.Module()
	sym := s.Symbol
	sw := float64(m * sym.W)
	sh := s.H
	if !sym.Linear {
		sh = float64(m * sym.H)
	}
	ox, oy := math.Floor((s.W-sw)/2), math.Floor((s.H-sh)/2)
	b0, b1 := sym.Bars()
	q0 := 0
	if !sym.Linear {
		q0 = b0 // the quiet zone is as wide at the top
	}
	// the bars in the upright box
	ux, uy := ox+float64(b0*m), oy+float64(q0*m)
	uw := float64((b1 - b0) * m)
	uh := sh - float64(2*q0*m)
	if sym.Linear {
		uh = float64(sym.BarHeight(float32(sw), float32(s.H)))
	}
	// to the label: the box's corner is at X, Y and turns about it
	bx, by := math.Round(s.X), math.Round(s.Y)
	switch s.Turns & 3 {
	case 1:
		return int(bx - uy - uh), int(by + ux), int(uh), int(uw)
	case 2:
		return int(bx - ux - uw), int(by - uy - uh), int(uw), int(uh)
	case 3:
		return int(bx + uy), int(by - ux - uw), int(uh), int(uw)
	}
	return int(bx + ux), int(by + uy), int(uw), int(uh)
}

// render draws pages for a printer of few dots.
type render struct {
	r     *imagebdf.Renderer
	view  *bdf.View
	o     *Options
	p     Profile
	warns []string

	// symbols is set while a page is drawn whose symbols the printer
	// prints, with those it takes
	mu      sync.Mutex
	symbols map[[5]float64]*Symbol
	order   []*Symbol
}

func newRender(doc *bdf.Document, view *bdf.View, o *Options, p Profile) *render {
	rd := &render{view: view, o: o, p: p}
	ro := &imagebdf.Options{FontFS: o.FontFS, FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, ImageCacheBytes: -1,
		Antialias: imagebdf.AntialiasOff, DitherImages: true, AlignBarcodes: true,
		BandRows: o.BandRows, BandWorkers: o.BandWorkers}
	if p.Barcodes && !o.Settings.NoPrinterBarcodes {
		ro.Barcode = rd.barcode
	}
	rd.r = imagebdf.New(doc, ro)
	return rd
}

// barcode takes the symbols the printer prints (imagebdf.Options.Barcode):
// those it has, of modules a dot wide or more.
func (rd *render) barcode(b imagebdf.Barcode) bool {
	rd.mu.Lock()
	defer rd.mu.Unlock()
	if rd.symbols == nil {
		return false
	}
	key := [5]float64{b.X, b.Y, b.W, b.H, float64(b.Turns)}
	if _, ok := rd.symbols[key]; ok {
		return true
	}
	sym, err := b.Mark.Symbol()
	if err != nil {
		return false
	}
	s := &Symbol{Mark: b.Mark, Symbol: sym, X: b.X, Y: b.Y, W: b.W, H: b.H, Turns: b.Turns}
	if s.Module() < 1 || !printable(rd.p.Language, s) || rd.p.Language == TPCL && len(rd.order) >= 32 {
		return false // (TPCL has 32 fields)
	}
	rd.symbols[key] = s
	rd.order = append(rd.order, s)
	return true
}

// Label is a page drawn on a label: where it goes, in dots.
type label struct {
	W, H  int // the label
	place imagebdf.Place
}

// labelOf places a page on a label of the profile's width (or the media's)
// and of the media's length, or of the page's own when there is none.
func (rd *render) labelOf(page int) (label, error) {
	p := rd.view.Pages[page]
	o := rd.o
	dpi := rd.p.DPI
	width := float64(rd.p.Width) * 72 / float64(dpi)
	if o.MediaWidth > 0 {
		width = math.Min(width, o.MediaWidth)
	}
	length := o.MediaLength
	pw, ph := float64(p.W), float64(p.H)
	if length <= 0 {
		// as long as the page, as wide as it is fitted
		tw, th := pw, ph
		if o.Orientation == raster.Landscape || o.Orientation == raster.ReverseLandscape {
			tw, th = ph, pw
		}
		f := 1.0
		switch o.Scaling {
		case raster.ScaleFit, raster.ScaleFill:
			f = width / tw
		case raster.ScaleNone:
		default:
			f = math.Min(1, width/tw)
		}
		length = th * f
	}
	sheet := raster.Sheet{Width: width, Height: length, DPI: dpi}
	orient := o.Orientation
	if o.MediaLength <= 0 && orient == raster.OrientAuto {
		orient = raster.Portrait // a roll is as long as a page needs
	}
	pl, err := raster.Fit(pw, ph, sheet, o.Scaling, orient)
	if err != nil {
		return label{}, err
	}
	if rd.p.YDPI != 0 && rd.p.YDPI != dpi {
		return label{}, fmt.Errorf("pdl: dots of %d × %d dpi", dpi, rd.p.YDPI)
	}
	if max := rd.p.MaxLength; max > 0 && pl.H > max {
		return label{}, fmt.Errorf("pdl: page %d is %d dots long, more than the %d the printer takes", page+1, pl.H, max)
	}
	return label{W: pl.W, H: pl.H, place: pl}, nil
}

// page draws a page as a whole bitmap, with the symbols the printer
// prints left out of it when symbols is set.
func (rd *render) page(page int, l label, symbols bool) (*Bitmap, []*Symbol, error) {
	bm := NewBitmap(l.W, l.H)
	if symbols {
		rd.mu.Lock()
		rd.symbols, rd.order = map[[5]float64]*Symbol{}, nil
		rd.mu.Unlock()
	}
	err := rd.r.PlaceRows(rd.view, page, l.place, imagebdf.Bilevel, func(y0 int, rows image.Image) error {
		bm.pack(y0, rows)
		return nil
	})
	rd.mu.Lock()
	syms := rd.order
	rd.symbols, rd.order = nil, nil
	rd.mu.Unlock()
	if err != nil {
		return nil, nil, fmt.Errorf("pdl: page %d: %w", page+1, err)
	}
	return bm, syms, nil
}

// rows draws a page in bands of rows, for a printer that prints as it is
// sent.
func (rd *render) rows(page int, l label, yield func(y0 int, rows *Bitmap) error) error {
	err := rd.r.PlaceRows(rd.view, page, l.place, imagebdf.Bilevel, func(y0 int, rows image.Image) error {
		bm := NewBitmap(l.W, rows.Bounds().Dy())
		bm.pack(0, rows)
		return yield(y0, bm)
	})
	if err != nil {
		return fmt.Errorf("pdl: page %d: %w", page+1, err)
	}
	return nil
}
