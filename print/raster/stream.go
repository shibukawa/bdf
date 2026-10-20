package raster

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"io"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/image/rowenc"
)

// Formats of the streams a PageReader reads.
const (
	FormatPWG = "image/pwg-raster"
	FormatURF = "image/urf"
)

// ErrNotRaster is returned by NewPageReader for a stream that is neither
// PWG Raster nor Apple Raster.
var ErrNotRaster = errors.New("raster: neither a PWG Raster nor an Apple Raster stream")

// Sniff names the raster format a stream begins with (FormatPWG or
// FormatURF), or returns "".
func Sniff(head []byte) string {
	switch {
	case bytes.HasPrefix(head, []byte(Sync)):
		return FormatPWG
	case bytes.HasPrefix(head, []byte(URFSync)):
		return FormatURF
	}
	return ""
}

// PageReader reads the pages of a PWG Raster or Apple Raster stream one
// at a time, and the rows of each a band at a time: a page is never whole
// in memory. It is what a print server reads the raster documents it is
// sent with.
type PageReader struct {
	br        *bufio.Reader
	format    string
	maxPixels int
	count     int // the pages an Apple Raster stream says it has
	pages     int // the pages begun
	page      *Page
	lines     *lineDecoder
	row       int // the rows of the page read
}

// NewPageReader begins to read a raster stream. It refuses pages of more
// than maxPixels pixels (64 Mi when 0).
func NewPageReader(r io.Reader, maxPixels int) (*PageReader, error) {
	if maxPixels <= 0 {
		maxPixels = 64 << 20
	}
	br := bufio.NewReaderSize(r, 64<<10)
	head, _ := br.Peek(len(URFSync))
	pr := &PageReader{br: br, format: Sniff(head), maxPixels: maxPixels}
	switch pr.format {
	case FormatPWG:
		br.Discard(len(Sync))
	case FormatURF:
		var h [12]byte
		if _, err := io.ReadFull(br, h[:]); err != nil {
			return nil, ErrNotRaster
		}
		pr.count = int(binary.BigEndian.Uint32(h[8:]))
	default:
		return nil, ErrNotRaster
	}
	return pr, nil
}

// Format is the stream's format: FormatPWG or FormatURF.
func (r *PageReader) Format() string { return r.format }

// Next reads the header of the next page, skipping what is left of the
// rows of the page before; it returns io.EOF after the last page. The
// pixels of a page of a bit a pixel are read as an *image.Gray of 0 (ink)
// and 255.
func (r *PageReader) Next() (*Page, error) {
	if r.page != nil {
		for r.row < r.page.Height {
			if err := r.lines.next(); err != nil {
				return nil, r.pageError(err)
			}
			r.row++
		}
	}
	r.page, r.lines, r.row = nil, nil, 0
	var p *Page
	if r.format == FormatPWG {
		var h [HeaderSize]byte
		if _, err := io.ReadFull(r.br, h[:]); err != nil {
			if err == io.EOF {
				return nil, io.EOF
			}
			return nil, fmt.Errorf("raster: page %d: %w", r.pages+1, err)
		}
		var err error
		if p, err = parseHeader(h); err != nil {
			return nil, fmt.Errorf("raster: page %d: %w", r.pages+1, err)
		}
	} else {
		var h [URFHeaderSize]byte
		if _, err := io.ReadFull(r.br, h[:]); err != nil {
			if err == io.EOF {
				if r.pages != r.count {
					return nil, fmt.Errorf("raster: %d pages of %d", r.pages, r.count)
				}
				return nil, io.EOF
			}
			return nil, fmt.Errorf("raster: page %d: %w", r.pages+1, err)
		}
		if r.pages >= r.count {
			return nil, fmt.Errorf("raster: more than the %d pages the stream has", r.count)
		}
		var err error
		if p, err = parseURFHeader(h, r.count); err != nil {
			return nil, fmt.Errorf("raster: page %d: %w", r.pages+1, err)
		}
	}
	if p.Width > 1<<20 || p.Height > 1<<20 || p.Width*p.Height > r.maxPixels {
		return nil, fmt.Errorf("raster: page %d is %d × %d pixels", r.pages+1, p.Width, p.Height)
	}
	r.pages++
	r.page = p
	r.lines = newLineDecoder(r.br, p.Type, p.Width)
	return p, nil
}

func (r *PageReader) pageError(err error) error {
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return fmt.Errorf("raster: page %d: %w", r.pages, err)
}

// ReadRows reads the next rows of the page into dst, from its top: an
// *image.Gray for the types of gray and of a bit a pixel, an *image.RGBA
// for srgb_8, as wide as the page. It returns the rows read, fewer than
// dst has at the end of the page, and io.EOF once all are read.
func (r *PageReader) ReadRows(dst image.Image) (int, error) {
	if r.page == nil {
		return 0, errors.New("raster: no page")
	}
	b := dst.Bounds()
	if b.Dx() != r.page.Width {
		return 0, fmt.Errorf("raster: rows %d wide for a page %d wide", b.Dx(), r.page.Width)
	}
	n := 0
	for y := b.Min.Y; y < b.Max.Y && r.row < r.page.Height; y++ {
		if err := r.lines.next(); err != nil {
			return n, r.pageError(err)
		}
		if err := r.lines.put(dst, y); err != nil {
			return n, err
		}
		r.row++
		n++
	}
	if n == 0 && r.row >= r.page.Height {
		return 0, io.EOF
	}
	return n, nil
}

// lineDecoder decodes the lines of a page, compressed as PWG Raster and
// Apple Raster compress them: a line with its repeat count, then runs of
// a pixel repeated and of pixels as they are.
type lineDecoder struct {
	br     *bufio.Reader
	typ    Type
	width  int
	unit   int // bytes of a pixel, 1 at least
	line   []byte
	repeat int // the times the line is still to be given again
}

func newLineDecoder(br *bufio.Reader, typ Type, width int) *lineDecoder {
	_, _, bits, _, _ := typ.layout()
	return &lineDecoder{br: br, typ: typ, width: width, unit: max(bits/8, 1), line: make([]byte, (width*bits+7)/8), repeat: -1}
}

// next makes the next line of the page the current one.
func (d *lineDecoder) next() error {
	if d.repeat > 0 {
		d.repeat--
		return nil
	}
	rep, err := d.br.ReadByte()
	if err != nil {
		return err
	}
	n, unit, line := len(d.line), d.unit, d.line
	for i := 0; i < n; {
		c, err := d.br.ReadByte()
		if err != nil {
			return err
		}
		if c < 128 {
			count := (int(c) + 1) * unit
			if i+count > n {
				return errors.New("a run past the end of the line")
			}
			if _, err := io.ReadFull(d.br, line[i:i+unit]); err != nil {
				return err
			}
			for k := unit; k < count; k++ {
				line[i+k] = line[i+k-unit]
			}
			i += count
		} else {
			count := (257 - int(c)) * unit
			if i+count > n {
				return errors.New("pixels past the end of the line")
			}
			if _, err := io.ReadFull(d.br, line[i:i+count]); err != nil {
				return err
			}
			i += count
		}
	}
	d.repeat = int(rep)
	return nil
}

// put writes the current line into the row y of dst.
func (d *lineDecoder) put(dst image.Image, y int) error {
	line := d.line
	switch d.typ {
	case SRGB8:
		rgba, ok := dst.(*image.RGBA)
		if !ok {
			return fmt.Errorf("raster: srgb_8 rows go in an *image.RGBA, not %T", dst)
		}
		row := rgba.Pix[rgba.PixOffset(rgba.Rect.Min.X, y):]
		for x := 0; x < d.width; x++ {
			row[4*x], row[4*x+1], row[4*x+2], row[4*x+3] = line[3*x], line[3*x+1], line[3*x+2], 255
		}
		return nil
	}
	gray, ok := dst.(*image.Gray)
	if !ok {
		return fmt.Errorf("raster: %s rows go in an *image.Gray, not %T", d.typ, dst)
	}
	row := gray.Pix[gray.PixOffset(gray.Rect.Min.X, y):]
	if d.typ == SGray8 {
		copy(row[:d.width], line)
		return nil
	}
	for x := 0; x < d.width; x++ {
		bit := line[x>>3]&(0x80>>(x&7)) != 0
		row[x] = 255
		if bit == (d.typ == Black1) {
			row[x] = 0
		}
	}
	return nil
}

// parseURFHeader reads the header of a page of Apple Raster.
func parseURFHeader(h [URFHeaderSize]byte, count int) (*Page, error) {
	p := &Page{Width: int(binary.BigEndian.Uint32(h[12:])), Height: int(binary.BigEndian.Uint32(h[16:])),
		Quality: int(h[3]), MediaPosition: int(h[5]), CrossFeed: 1, Feed: 1, Copies: 1, TotalPages: count}
	p.XDPI = int(binary.BigEndian.Uint32(h[20:]))
	p.YDPI = p.XDPI
	switch {
	case h[0] == 8 && h[1] == urfSGray:
		p.Type = SGray8
	case h[0] == 24 && h[1] == urfSRGB:
		p.Type = SRGB8
	default:
		return nil, fmt.Errorf("%d bits a pixel in colour space %d", h[0], h[1])
	}
	p.Duplex, p.Tumble = h[2] == urfShort || h[2] == urfLong, h[2] == urfShort
	if p.Width <= 0 || p.Height <= 0 || p.XDPI <= 0 || p.XDPI > 10000 {
		return nil, fmt.Errorf("a page of %d × %d pixels at %d dpi", p.Width, p.Height, p.XDPI)
	}
	// an Apple Raster page has no size of paper but its pixels
	p.PageWidth, p.PageHeight = (p.Width*72+p.XDPI/2)/p.XDPI, (p.Height*72+p.YDPI/2)/p.YDPI
	return p, nil
}

// DocumentOptions bound what Document makes.
type DocumentOptions struct {
	// MaxPixels bounds a page, as NewPageReader's; MaxPages the pages
	// (10,000 when 0); MaxTotalPixels the pages together (8 Gi when 0:
	// some 250 pages of A4 in colour at 600 dpi), as a small stream of
	// repeated lines can stand for many pixels to compress.
	MaxPixels, MaxPages int
	MaxTotalPixels      int64
	// Title is the title of the document.
	Title string
}

// RasterInfo tells what Document read.
type RasterInfo struct {
	Format string
	Pages  int
	// Duplex and Tumble are the first page's.
	Duplex, Tumble bool
	// Resolution is the first page's, in dots per inch.
	Resolution int
}

// Document makes a document of the pages of a PWG Raster or Apple Raster
// stream: each page an image of its pixels, as large as the paper the
// stream says (its pixels at their resolution for Apple Raster, which
// says none). The pixels are kept as PNG: gray, bilevel (for the types of
// a bit a pixel) or RGB, read a band at a time.
// It stops when ctx ends, between bands.
func Document(ctx context.Context, r io.Reader, o *DocumentOptions) (*bdf.Document, *RasterInfo, error) {
	if o == nil {
		o = &DocumentOptions{}
	}
	maxTotal := o.MaxTotalPixels
	if maxTotal <= 0 {
		maxTotal = 8 << 30
	}
	var total int64
	maxPages := o.MaxPages
	if maxPages <= 0 {
		maxPages = 10000
	}
	pr, err := NewPageReader(r, o.MaxPixels)
	if err != nil {
		return nil, nil, err
	}
	doc := bdf.NewDocument()
	doc.Meta.Source = map[string]string{FormatPWG: "pwg-raster", FormatURF: "urf"}[pr.format]
	if o.Title != "" {
		doc.Meta.DC.Title = bdf.DCValues{o.Title}
	}
	view := doc.NewView("pages", bdf.ViewFixed, o.Title)
	info := &RasterInfo{Format: pr.format}
	var png bytes.Buffer
	for {
		p, err := pr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		if info.Pages++; info.Pages > maxPages {
			return nil, nil, fmt.Errorf("raster: more than %d pages", maxPages)
		}
		if info.Pages == 1 {
			info.Duplex, info.Tumble, info.Resolution = p.Duplex, p.Tumble, p.XDPI
		}
		if total += int64(p.Width) * int64(p.Height); total > maxTotal {
			return nil, nil, fmt.Errorf("raster: more than %d pixels in all", maxTotal)
		}
		png.Reset()
		if err := encodePage(ctx, &png, pr, p); err != nil {
			return nil, nil, err
		}
		pw, ph := float32(p.PageWidth), float32(p.PageHeight)
		if pw <= 0 || ph <= 0 {
			pw, ph = float32(p.Width)*72/float32(p.XDPI), float32(p.Height)*72/float32(p.YDPI)
		}
		obj := bdf.NewObject()
		img := obj.AddImage(doc.AddImage(bytes.Clone(png.Bytes())))
		obj.Image(img, 0, 0, pw, ph)
		hash, _ := doc.AddObject(obj)
		view.AddPage(pw, ph, bdf.Layer{Role: bdf.RoleBody, Obj: hash})
	}
	if info.Pages == 0 {
		return nil, nil, errors.New("raster: no pages")
	}
	return doc, info, nil
}

// encodePage writes the pixels of the page as a PNG, a band of rows at a
// time.
func encodePage(ctx context.Context, w io.Writer, pr *PageReader, p *Page) error {
	const band = 64
	kind := rowenc.Gray
	var rows image.Image
	switch p.Type {
	case SRGB8:
		kind = rowenc.RGBA
		rows = image.NewRGBA(image.Rect(0, 0, p.Width, band))
	case Black1, SGray1:
		kind = rowenc.Bilevel
		rows = image.NewGray(image.Rect(0, 0, p.Width, band))
	default:
		rows = image.NewGray(image.Rect(0, 0, p.Width, band))
	}
	enc, err := rowenc.NewPNG(w, p.Width, p.Height, kind, &rowenc.Options{Level: 6, DPI: float64(p.XDPI)})
	if err != nil {
		return err
	}
	for {
		if err := ctx.Err(); err != nil {
			enc.Close()
			return err
		}
		n, err := pr.ReadRows(rows)
		if err == io.EOF {
			break
		}
		if err != nil {
			enc.Close()
			return err
		}
		var part image.Image
		switch img := rows.(type) {
		case *image.RGBA:
			part = img.SubImage(image.Rect(0, 0, p.Width, n))
		case *image.Gray:
			part = img.SubImage(image.Rect(0, 0, p.Width, n))
		}
		if err := enc.WriteRows(part); err != nil {
			enc.Close()
			return err
		}
	}
	return enc.Close()
}
