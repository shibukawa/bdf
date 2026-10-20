// Package raster writes the raster formats of driverless printing: PWG
// Raster (PWG 5102.4, image/pwg-raster), which IPP Everywhere printers
// take, Apple Raster (image/urf), which AirPrint printers take, and PCLm
// (application/PCLm), which Mopria printers take. A page is written a
// band of rows at a time, as raster/imagebdf draws it, so a page at 600
// dpi is never whole in memory.
//
// The package also places pages on sheets of paper (Fit): the scale, the
// quarter turn and the offset that print a page of a document on a sheet
// as print-scaling and orientation ask.
package raster

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"io"
)

// Type is the colour space and depth of the pixels of a page, as
// pwg-raster-document-type-supported names them.
type Type string

// Types of pixels.
const (
	// Black1 is a bit a pixel, 1 for black (ink): for black and white.
	Black1 Type = "black_1"
	// SGray1 is a bit a pixel, 1 for white.
	SGray1 Type = "sgray_1"
	// SGray8 is a byte a pixel of the luminance of sRGB, 0 for black.
	SGray8 Type = "sgray_8"
	// SRGB8 is three bytes a pixel, red, green and blue of sRGB.
	SRGB8 Type = "srgb_8"
)

// colour spaces of the header (cups_cspace_t)
const (
	csBlack = 3
	csSGray = 18
	csSRGB  = 19
)

func (t Type) layout() (colorSpace, bitsPerColor, bitsPerPixel, numColors int, ok bool) {
	switch t {
	case Black1:
		return csBlack, 1, 1, 1, true
	case SGray1:
		return csSGray, 1, 1, 1, true
	case SGray8:
		return csSGray, 8, 8, 1, true
	case SRGB8:
		return csSRGB, 8, 24, 3, true
	}
	return 0, 0, 0, 0, false
}

// Page describes a page of PWG Raster: its header.
type Page struct {
	// Width and Height are the pixels of the page (of the sheet).
	Width, Height int
	// XDPI and YDPI are the resolution.
	XDPI, YDPI int
	Type       Type
	// PageWidth and PageHeight are the size of the sheet in points.
	PageWidth, PageHeight int
	// PageSizeName is the PWG name of the media ("iso_a4_210x297mm").
	PageSizeName string
	// MediaType ("stationery") and MediaColor, when known.
	MediaType, MediaColor string
	// Duplex and Tumble tell the sides: Duplex for two-sided, Tumble
	// for the short edge.
	Duplex, Tumble bool
	// CrossFeed and Feed tell how the image of the page is transformed:
	// -1 for mirrored across (left to right) or along (top to bottom)
	// the feed, 1 (or 0) for not.
	CrossFeed, Feed int
	// Copies is NumCopies (1 when 0).
	Copies int
	// TotalPages is the pages of the document, 0 when not known.
	TotalPages int
	// Quality is the print-quality (3 draft, 4 normal, 5 high; 0 for the
	// printer's default).
	Quality int
	// MediaPosition is the input tray (0 for automatic).
	MediaPosition int
}

// HeaderSize is the bytes of the header of a page.
const HeaderSize = 1796

// offsets of the fields of the header (PWG 5102.4 §4.3)
const (
	offMediaClass   = 0 // "PwgRaster"
	offMediaColor   = 64
	offMediaType    = 128
	offCutMedia     = 268
	offDuplex       = 272
	offHWResolution = 276
	offMediaPos     = 324
	offNumCopies    = 340
	offPageSize     = 352
	offTumble       = 368
	offWidth        = 372
	offHeight       = 376
	offBitsPerColor = 384
	offBitsPerPixel = 388
	offBytesPerLine = 392
	offColorOrder   = 396
	offColorSpace   = 400
	offNumColors    = 420
	offTotalPages   = 452
	offCrossFeed    = 456
	offFeed         = 460
	offImageBox     = 464
	offQuality      = 484
	offPageSizeName = 1732
)

// Sync is the first four bytes of a PWG Raster stream.
const Sync = "RaS2"

// Header returns the header of a page.
func (p *Page) Header() ([HeaderSize]byte, error) {
	var h [HeaderSize]byte
	cs, bpc, bpp, nc, ok := p.Type.layout()
	if !ok {
		return h, fmt.Errorf("raster: no PWG Raster type %q", p.Type)
	}
	if p.Width <= 0 || p.Height <= 0 || p.XDPI <= 0 || p.YDPI <= 0 {
		return h, errors.New("raster: an empty page")
	}
	str := func(off int, s string) error {
		if len(s) > 63 {
			return fmt.Errorf("raster: %q is longer than 63 bytes", s)
		}
		copy(h[off:off+64], s)
		return nil
	}
	u32 := func(off int, v int) { binary.BigEndian.PutUint32(h[off:], uint32(int32(v))) }
	b2i := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}
	sign := func(v int) int {
		if v < 0 {
			return -1
		}
		return 1
	}
	for off, s := range map[int]string{offMediaClass: "PwgRaster", offMediaColor: p.MediaColor, offMediaType: p.MediaType, offPageSizeName: p.PageSizeName} {
		if err := str(off, s); err != nil {
			return h, err
		}
	}
	copies := p.Copies
	if copies <= 0 {
		copies = 1
	}
	u32(offDuplex, b2i(p.Duplex))
	u32(offHWResolution, p.XDPI)
	u32(offHWResolution+4, p.YDPI)
	u32(offMediaPos, p.MediaPosition)
	u32(offNumCopies, copies)
	u32(offPageSize, p.PageWidth)
	u32(offPageSize+4, p.PageHeight)
	u32(offTumble, b2i(p.Tumble))
	u32(offWidth, p.Width)
	u32(offHeight, p.Height)
	u32(offBitsPerColor, bpc)
	u32(offBitsPerPixel, bpp)
	u32(offBytesPerLine, (p.Width*bpp+7)/8)
	u32(offColorOrder, 0) // chunky
	u32(offColorSpace, cs)
	u32(offNumColors, nc)
	u32(offTotalPages, p.TotalPages)
	u32(offCrossFeed, sign(p.CrossFeed))
	u32(offFeed, sign(p.Feed))
	u32(offImageBox, 0)
	u32(offImageBox+4, 0)
	u32(offImageBox+8, p.Width)
	u32(offImageBox+12, p.Height)
	u32(offQuality, p.Quality)
	return h, nil
}

// Writer writes a PWG Raster stream: the sync word, then for each page
// its header and its rows, compressed as PWG 5102.4 §4.4 says (a line
// repeated up to 256 times, and in a line, runs of a pixel repeated or
// of pixels as they are).
type Writer struct {
	w       *bufio.Writer
	started bool
	// urf writes Apple Raster (NewURFWriter) of urfPages pages
	urf      bool
	urfPages int
	page     *Page
	bpp      int // bytes a pixel unit: 1 for 1 bit (8 pixels), 1 for sgray_8, 3 for srgb_8
	line     []byte
	prev     []byte
	repeat   int // times prev is repeated, beyond the first
	rows     int // rows written of the page
	held     bool
	enc      []byte
}

// NewWriter returns a writer of a PWG Raster stream to w.
func NewWriter(w io.Writer) *Writer { return &Writer{w: bufio.NewWriterSize(w, 64<<10)} }

// ErrRows is returned when a page is given more or fewer rows than its
// height.
var ErrRows = errors.New("raster: the rows do not fill the page")

// StartPage writes the header of a page; its rows follow with WriteRows.
// The page before it must have had all its rows.
func (w *Writer) StartPage(p Page) error {
	if w.page != nil && w.rows != w.page.Height {
		return fmt.Errorf("%w: %d of %d rows", ErrRows, w.rows, w.page.Height)
	}
	if err := w.flushLine(); err != nil {
		return err
	}
	var h []byte
	if w.urf {
		uh, err := p.urfHeader()
		if err != nil {
			return err
		}
		h = uh[:]
	} else {
		ph, err := p.Header()
		if err != nil {
			return err
		}
		h = ph[:]
	}
	if err := w.start(); err != nil {
		return err
	}
	if _, err := w.w.Write(h); err != nil {
		return err
	}
	_, _, bits, _, _ := p.Type.layout()
	w.page, w.rows = &p, 0
	w.bpp = max(bits/8, 1)
	n := (p.Width*bits + 7) / 8
	w.line = make([]byte, n)
	w.prev = make([]byte, n)
	w.held = false
	return nil
}

// WriteRows writes rows of the page from an image: an *image.Gray, an
// *image.RGBA, or an *image.Paletted of black and white (black first),
// as raster/imagebdf draws them, whose width is the page's and whose
// rows follow the ones written. Other images are converted pixel by
// pixel. Colour is put over white where it is not opaque.
func (w *Writer) WriteRows(img image.Image) error {
	if w.page == nil {
		return errors.New("raster: WriteRows before StartPage")
	}
	b := img.Bounds()
	if b.Dx() != w.page.Width {
		return fmt.Errorf("raster: rows of %d pixels for a page of %d", b.Dx(), w.page.Width)
	}
	if w.rows+b.Dy() > w.page.Height {
		return fmt.Errorf("%w: %d rows past %d", ErrRows, w.rows+b.Dy(), w.page.Height)
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		w.convert(img, y)
		if err := w.addLine(); err != nil {
			return err
		}
	}
	return nil
}

// Close ends the stream: the last page must have had all its rows.
func (w *Writer) Close() error {
	if w.page != nil && w.rows != w.page.Height {
		return fmt.Errorf("%w: %d of %d rows", ErrRows, w.rows, w.page.Height)
	}
	if err := w.flushLine(); err != nil {
		return err
	}
	if err := w.start(); err != nil {
		return err
	}
	return w.w.Flush()
}

// start writes the beginning of the stream, once.
func (w *Writer) start() error {
	if w.started {
		return nil
	}
	w.started = true
	if w.urf {
		var b [12]byte
		copy(b[:], URFSync)
		binary.BigEndian.PutUint32(b[8:], uint32(w.urfPages))
		_, err := w.w.Write(b[:])
		return err
	}
	_, err := w.w.WriteString(Sync)
	return err
}

// convert puts row y of img into w.line in the page's type.
func (w *Writer) convert(img image.Image, y int) { convertRow(w.line, img, y, w.page.Type) }

// convertRow puts row y of img into line, in pixels of a type.
func convertRow(line []byte, img image.Image, y int, typ Type) {
	b := img.Bounds()
	switch typ {
	case SGray8:
		if g, ok := img.(*image.Gray); ok {
			copy(line, g.Pix[(y-g.Rect.Min.Y)*g.Stride:])
			return
		}
		for x := 0; x < b.Dx(); x++ {
			r, gg, bb := over(img, b.Min.X+x, y)
			line[x] = luma(r, gg, bb)
		}
	case SRGB8:
		for x := 0; x < b.Dx(); x++ {
			r, g, bb := over(img, b.Min.X+x, y)
			line[3*x], line[3*x+1], line[3*x+2] = r, g, bb
		}
	case Black1, SGray1:
		inkIsOne := typ == Black1
		for i := range line {
			line[i] = 0
		}
		if !inkIsOne {
			// padding bits are white: 1 in sgray_1
			for i := range line {
				line[i] = 0xff
			}
		}
		pal, isPal := img.(*image.Paletted)
		blackIndex := uint8(0)
		if isPal && len(pal.Palette) > 0 {
			if r, _, _, _ := pal.Palette[0].RGBA(); r > 0x8000 {
				blackIndex = 1
			}
		}
		for x := 0; x < b.Dx(); x++ {
			var dark bool
			if isPal {
				dark = pal.Pix[(y-pal.Rect.Min.Y)*pal.Stride+x] == blackIndex
			} else if g, ok := img.(*image.Gray); ok {
				dark = g.Pix[(y-g.Rect.Min.Y)*g.Stride+x] < 128
			} else {
				r, gg, bb := over(img, b.Min.X+x, y)
				dark = luma(r, gg, bb) < 128
			}
			bit := byte(0x80) >> (x & 7)
			if dark == inkIsOne {
				line[x>>3] |= bit
			} else {
				line[x>>3] &^= bit
			}
		}
	}
}

// over returns the colour of a pixel over white, in 8 bits.
func over(img image.Image, x, y int) (r, g, b uint8) {
	if rgba, ok := img.(*image.RGBA); ok {
		i := (y-rgba.Rect.Min.Y)*rgba.Stride + (x-rgba.Rect.Min.X)*4
		p := rgba.Pix[i : i+4 : i+4]
		inv := 255 - p[3] // premultiplied: add white for what is not covered
		return p[0] + inv, p[1] + inv, p[2] + inv
	}
	cr, cg, cb, ca := img.At(x, y).RGBA()
	inv := 0xffff - ca
	return uint8((cr + inv) >> 8), uint8((cg + inv) >> 8), uint8((cb + inv) >> 8)
}

// luma is the luminance of Rec. 601, as raster/imagebdf and image/color
// make gray.
func luma(r, g, b uint8) uint8 {
	return uint8((19595*uint32(r) + 38470*uint32(g) + 7471*uint32(b) + 1<<15) >> 16)
}

// addLine adds w.line to the page: as a repeat of the line before it,
// or as a line of its own once the one before is written.
func (w *Writer) addLine() error {
	w.rows++
	if w.held && w.repeat < 255 && string(w.line) == string(w.prev) {
		w.repeat++
		return nil
	}
	if err := w.flushLine(); err != nil {
		return err
	}
	w.prev, w.line = w.line, w.prev
	w.held, w.repeat = true, 0
	return nil
}

// flushLine writes the line held with its repeat count.
func (w *Writer) flushLine() error {
	if !w.held {
		return nil
	}
	w.held = false
	w.enc = encodeLine(w.enc[:0], byte(w.repeat), w.prev, w.bpp)
	_, err := w.w.Write(w.enc)
	return err
}

// encodeLine appends a line: its repeat byte, then runs of pixel units
// of bpp bytes: n-1 then a unit for n (1 to 128) equal units, 257-n then
// the units for n (2 to 128) units as they are.
func encodeLine(out []byte, repeat byte, line []byte, bpp int) []byte {
	out = append(out, repeat)
	n := len(line) / bpp
	unit := func(i int) []byte { return line[i*bpp : (i+1)*bpp] }
	eq := func(i, j int) bool { return string(unit(i)) == string(unit(j)) }
	for i := 0; i < n; {
		// a run of equal units
		j := i + 1
		for j < n && j-i < 128 && eq(i, j) {
			j++
		}
		if j-i > 1 || j == n {
			out = append(out, byte(j-i-1))
			out = append(out, unit(i)...)
			i = j
			continue
		}
		// units as they are, up to the next pair of equal ones
		j = i + 1
		for j < n && j-i < 128 && !(j+1 < n && eq(j, j+1)) {
			j++
		}
		if j-i == 1 {
			out = append(out, 0)
			out = append(out, unit(i)...)
			i = j
			continue
		}
		out = append(out, byte(257-(j-i)))
		out = append(out, line[i*bpp:j*bpp]...)
		i = j
	}
	return out
}
