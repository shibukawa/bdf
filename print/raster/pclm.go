package raster

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"strconv"
	"strings"
)

// PCLmOptions are how PCLm is written, as the printer prefers it.
type PCLmOptions struct {
	// StripHeight is the rows of a strip (pclm-strip-height-preferred):
	// 16 when 0.
	StripHeight int
	// Compression is the compression of the strips: "flate" (the
	// default), "jpeg" or "rle" (pclm-compression-method-preferred).
	Compression string
	// JPEGQuality is the quality of "jpeg": 90 when 0.
	JPEGQuality int
}

// PCLmWriter writes PCLm (application/PCLm, the raster of Mopria): a PDF
// 1.7 file, whose second line says it is PCLm 1.0, of pages that are
// strips of image drawn one under the other, at the resolution of the
// printer. Its pages are of SGray8 (DeviceGray) or SRGB8 (DeviceRGB).
// Strips are written as their rows come, and the file's objects are
// listed at its end.
type PCLmWriter struct {
	w    *countWriter
	opts PCLmOptions

	offsets []int64 // of the objects, by number (0 is not one)
	pages   []int   // the objects of the pages
	page    *Page
	rows    int // rows written of the page
	unit    int // bytes a pixel
	strip   []byte
	inStrip int   // rows in strip
	strips  []int // the objects of the strips of the page, from its top
	heights []int
	err     error
}

type countWriter struct {
	w *bufio.Writer
	n int64
}

func (c *countWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.n += int64(n)
	return n, err
}

// NewPCLmWriter returns a writer of PCLm to w.
func NewPCLmWriter(w io.Writer, opts PCLmOptions) (*PCLmWriter, error) {
	if opts.StripHeight <= 0 {
		opts.StripHeight = 16
	}
	switch opts.Compression {
	case "":
		opts.Compression = "flate"
	case "flate", "jpeg", "rle":
	default:
		return nil, fmt.Errorf("raster: PCLm compression %q: flate, jpeg or rle", opts.Compression)
	}
	if opts.JPEGQuality <= 0 {
		opts.JPEGQuality = 90
	}
	pw := &PCLmWriter{w: &countWriter{w: bufio.NewWriterSize(w, 64<<10)}, opts: opts, offsets: []int64{0, 0, 0}}
	// the catalog (1) and the page tree (2) are written at the end
	_, pw.err = io.WriteString(pw.w, "%PDF-1.7\n%PCLm 1.0\n")
	return pw, pw.err
}

// object starts the next object of the file and returns its number.
func (p *PCLmWriter) object() int {
	p.offsets = append(p.offsets, p.w.n)
	n := len(p.offsets) - 1
	fmt.Fprintf(p.w, "%d 0 obj\n", n)
	return n
}

// StartPage starts a page; its rows follow with WriteRows.
func (p *PCLmWriter) StartPage(pg Page) error {
	if p.err != nil {
		return p.err
	}
	if err := p.endPage(); err != nil {
		return err
	}
	if pg.Width <= 0 || pg.Height <= 0 || pg.XDPI <= 0 {
		return errors.New("raster: an empty page")
	}
	if pg.YDPI <= 0 {
		pg.YDPI = pg.XDPI
	}
	switch pg.Type {
	case SGray8:
		p.unit = 1
	case SRGB8:
		p.unit = 3
	default:
		return fmt.Errorf("raster: PCLm has no type %q (sgray_8 or srgb_8)", pg.Type)
	}
	p.page, p.rows, p.inStrip = &pg, 0, 0
	p.strips, p.heights = p.strips[:0], p.heights[:0]
	p.strip = make([]byte, 0, pg.Width*p.unit*p.opts.StripHeight)
	return nil
}

// WriteRows writes rows of the page, as Writer.WriteRows does.
func (p *PCLmWriter) WriteRows(img image.Image) error {
	if p.err != nil {
		return p.err
	}
	if p.page == nil {
		return errors.New("raster: WriteRows before StartPage")
	}
	b := img.Bounds()
	if b.Dx() != p.page.Width {
		return fmt.Errorf("raster: rows of %d pixels for a page of %d", b.Dx(), p.page.Width)
	}
	if p.rows+b.Dy() > p.page.Height {
		return fmt.Errorf("%w: %d rows past %d", ErrRows, p.rows+b.Dy(), p.page.Height)
	}
	line := make([]byte, p.page.Width*p.unit)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		convertRow(line, img, y, p.page.Type)
		p.strip = append(p.strip, line...)
		p.rows++
		p.inStrip++
		if p.inStrip == p.opts.StripHeight || p.rows == p.page.Height {
			if err := p.writeStrip(); err != nil {
				return err
			}
		}
	}
	return nil
}

// writeStrip writes the rows held as an image object.
func (p *PCLmWriter) writeStrip() error {
	w, h := p.page.Width, p.inStrip
	var data []byte
	var filter string
	switch p.opts.Compression {
	case "flate":
		var buf bytes.Buffer
		zw := zlib.NewWriter(&buf)
		zw.Write(p.strip)
		zw.Close()
		data, filter = buf.Bytes(), "/FlateDecode"
	case "rle":
		data, filter = runLength(p.strip), "/RunLengthDecode"
	case "jpeg":
		var img image.Image
		if p.unit == 1 {
			img = &image.Gray{Pix: p.strip, Stride: w, Rect: image.Rect(0, 0, w, h)}
		} else {
			rgba := image.NewRGBA(image.Rect(0, 0, w, h))
			for i, j := 0, 0; i < len(p.strip); i, j = i+3, j+4 {
				rgba.Pix[j], rgba.Pix[j+1], rgba.Pix[j+2], rgba.Pix[j+3] = p.strip[i], p.strip[i+1], p.strip[i+2], 255
			}
			img = rgba
		}
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: p.opts.JPEGQuality}); err != nil {
			return err
		}
		data, filter = buf.Bytes(), "/DCTDecode"
	}
	space := "/DeviceGray"
	if p.unit == 3 {
		space = "/DeviceRGB"
	}
	n := p.object()
	fmt.Fprintf(p.w, "<< /Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace %s /BitsPerComponent 8 /Filter %s /Length %d >>\nstream\n",
		w, h, space, filter, len(data))
	p.w.Write(data)
	if _, err := io.WriteString(p.w, "\nendstream\nendobj\n"); err != nil {
		p.err = err
		return err
	}
	p.strips = append(p.strips, n)
	p.heights = append(p.heights, h)
	p.strip, p.inStrip = p.strip[:0], 0
	return nil
}

// endPage writes the content and the object of the page begun.
func (p *PCLmWriter) endPage() error {
	if p.page == nil {
		return nil
	}
	pg := p.page
	p.page = nil
	if p.rows != pg.Height {
		return fmt.Errorf("%w: %d of %d rows", ErrRows, p.rows, pg.Height)
	}
	// the strips in pixels, under a transform from pixels to points
	var c strings.Builder
	scale := 72 / float64(pg.XDPI)
	sy := 72 / float64(pg.YDPI)
	fmt.Fprintf(&c, "%s 0 0 %s 0 0 cm\n/P <</MCID 0>> BDC\n", fnum(scale), fnum(sy))
	y := pg.Height
	var xobjects strings.Builder
	for i, n := range p.strips {
		y -= p.heights[i]
		fmt.Fprintf(&c, "q %d 0 0 %d 0 %d cm /S%d Do Q\n", pg.Width, p.heights[i], y, i)
		fmt.Fprintf(&xobjects, " /S%d %d 0 R", i, n)
	}
	c.WriteString("EMC\n")
	content := p.object()
	fmt.Fprintf(p.w, "<< /Length %d >>\nstream\n%s\nendstream\nendobj\n", c.Len(), c.String())
	n := p.object()
	fmt.Fprintf(p.w, "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %s %s] /Resources << /XObject <<%s >> >> /Contents %d 0 R >>\nendobj\n",
		fnum(float64(pg.Width)*scale), fnum(float64(pg.Height)*sy), xobjects.String(), content)
	p.pages = append(p.pages, n)
	return nil
}

// Close ends the last page and the file.
func (p *PCLmWriter) Close() error {
	if p.err != nil {
		return p.err
	}
	if err := p.endPage(); err != nil {
		return err
	}
	p.offsets[1] = p.w.n
	io.WriteString(p.w, "1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	p.offsets[2] = p.w.n
	kids := make([]string, len(p.pages))
	for i, n := range p.pages {
		kids[i] = strconv.Itoa(n) + " 0 R"
	}
	fmt.Fprintf(p.w, "2 0 obj\n<< /Type /Pages /Kids [%s] /Count %d >>\nendobj\n", strings.Join(kids, " "), len(p.pages))
	xref := p.w.n
	fmt.Fprintf(p.w, "xref\n0 %d\n0000000000 65535 f \n", len(p.offsets))
	for _, off := range p.offsets[1:] {
		fmt.Fprintf(p.w, "%010d 00000 n \n", off)
	}
	fmt.Fprintf(p.w, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(p.offsets), xref)
	return p.w.w.Flush()
}

// fnum writes a number of a PDF: no exponent, no needless digits.
func fnum(f float64) string {
	s := strconv.FormatFloat(f, 'f', 6, 64)
	s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
	if s == "" || s == "-0" {
		return "0"
	}
	return s
}

// runLength compresses data as PDF's RunLengthDecode reads it: n (0 to
// 127) then n+1 bytes as they are, or 257-n (129 to 255) then a byte
// repeated, and 128 at the end.
func runLength(data []byte) []byte {
	out := make([]byte, 0, len(data)/2)
	for i := 0; i < len(data); {
		j := i + 1
		for j < len(data) && j-i < 128 && data[j] == data[i] {
			j++
		}
		if j-i > 1 {
			out = append(out, byte(257-(j-i)), data[i])
			i = j
			continue
		}
		j = i + 1
		for j < len(data) && j-i < 128 && !(j+1 < len(data) && data[j] == data[j+1]) {
			j++
		}
		out = append(out, byte(j-i-1))
		out = append(out, data[i:j]...)
		i = j
	}
	return append(out, 128)
}
