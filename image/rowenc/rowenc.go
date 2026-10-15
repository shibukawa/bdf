// Package rowenc writes PNG and TIFF files from the rows of an image as
// they come, for an image drawn in bands (imagebdf.PageRows) that is not
// all in memory at once: a page at 600 dpi is written with the memory of
// a band of it and of the compression.
//
// A Writer is made for an image of a width, height and Kind, takes the
// rows in order with WriteRows, and is closed. The files decode (with
// image/png and a TIFF reader) to the pixels given: the PNG of an RGBA
// image is written as image/png writes an *image.RGBA, in colour type
// RGBA with the alpha unmultiplied, whether or not the image is opaque,
// since the rows to come are not known.
package rowenc

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
)

// Kind is the pixels of an image written.
type Kind int

const (
	// RGBA is 8 bits a channel with premultiplied alpha, as *image.RGBA
	// holds them: the PNG is of colour type RGBA (the alpha unmultiplied,
	// as image/png writes it), the TIFF of RGB with an associated alpha
	// sample.
	RGBA Kind = iota
	// Gray is 8 bits a pixel, as *image.Gray.
	Gray
	// Bilevel is a bit a pixel, black and white: the rows are an
	// *image.Paletted of two colours, black first (as imagebdf.Bilevel
	// draws them), or an *image.Gray whose values below 128 are black.
	// The PNG is of a palette of the two colours at a bit a pixel, the
	// TIFF of a bit a pixel with 0 black (min-is-black).
	Bilevel
)

// Options configures a Writer.
type Options struct {
	// Level is the zlib compression level of a PNG, from 1 (fastest) to
	// 9 (smallest), as compress/flate numbers them; 0 is the default
	// level (6). A TIFF is not compressed.
	Level int
	// DPI is the resolution noted in the file: the pHYs chunk of a PNG
	// and the resolution tags of a TIFF. Zero notes none in a PNG and 72
	// in a TIFF.
	DPI float64
}

// Writer writes an image from its rows, from the top.
type Writer interface {
	// WriteRows writes the next rows: an image as wide as the one
	// written, of the kind of the writer (see Kind), whose rows follow
	// those written before. Its bounds may start anywhere.
	WriteRows(rows image.Image) error
	// Close writes what ends the file and reports an error if the rows
	// written are not all of the image.
	Close() error
}

// rower converts the rows of an image to the bytes of each, packed as
// the file holds them, and counts them.
type rower struct {
	width, height int
	kind          Kind
	written       int
	row           []byte // the bytes of a row, for packing
	bits          []byte // a byte a pixel of a bilevel row from gray pixels
}

func newRower(width, height int, kind Kind) (*rower, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("rowenc: the image is empty")
	}
	if width > 1<<30 || height > 1<<30 {
		return nil, errors.New("rowenc: the image is too large")
	}
	if kind < RGBA || kind > Bilevel {
		return nil, fmt.Errorf("rowenc: unknown kind %d", kind)
	}
	r := &rower{width: width, height: height, kind: kind}
	r.row = make([]byte, r.rowBytes())
	return r, nil
}

// rowBytes is the bytes of a row in the file.
func (r *rower) rowBytes() int {
	switch r.kind {
	case RGBA:
		return 4 * r.width
	case Gray:
		return r.width
	}
	return (r.width + 7) / 8
}

// each calls fn with the bytes of each row of rows, in order, after
// checking the rows against the image. png says that the alpha of RGBA
// pixels is unmultiplied, as PNG holds it.
func (r *rower) each(rows image.Image, png bool, fn func(row []byte) error) error {
	if rows == nil {
		return errors.New("rowenc: no rows")
	}
	b := rows.Bounds()
	if b.Dx() != r.width {
		return fmt.Errorf("rowenc: rows %d wide for an image %d wide", b.Dx(), r.width)
	}
	if b.Dy() <= 0 {
		return nil
	}
	if r.written+b.Dy() > r.height {
		return fmt.Errorf("rowenc: %d rows past the %d of the image", r.written+b.Dy()-r.height, r.height)
	}
	var pack func(y int) error
	switch r.kind {
	case RGBA:
		img, ok := rows.(*image.RGBA)
		if !ok {
			return fmt.Errorf("rowenc: RGBA rows are %T", rows)
		}
		pack = func(y int) error {
			src := img.Pix[(y-b.Min.Y)*img.Stride+4*(b.Min.X-img.Rect.Min.X):]
			src = src[:4*r.width]
			if !png {
				copy(r.row, src)
				return nil
			}
			for i := 0; i+3 < len(src); i += 4 {
				a := src[i+3]
				switch a {
				case 255:
					copy(r.row[i:i+4], src[i:i+4])
				case 0:
					r.row[i], r.row[i+1], r.row[i+2], r.row[i+3] = 0, 0, 0, 0
				default:
					// as color.NRGBAModel converts an RGBA colour: in 16
					// bits, then the high byte
					c := color.NRGBAModel.Convert(color.RGBA{src[i], src[i+1], src[i+2], a}).(color.NRGBA)
					r.row[i], r.row[i+1], r.row[i+2], r.row[i+3] = c.R, c.G, c.B, c.A
				}
			}
			return nil
		}
	case Gray:
		img, ok := rows.(*image.Gray)
		if !ok {
			return fmt.Errorf("rowenc: gray rows are %T", rows)
		}
		pack = func(y int) error {
			src := img.Pix[(y-b.Min.Y)*img.Stride+b.Min.X-img.Rect.Min.X:]
			copy(r.row, src[:r.width])
			return nil
		}
	case Bilevel:
		var white func(y int) []byte // a byte a pixel: 0 black, else white
		switch img := rows.(type) {
		case *image.Paletted:
			if len(img.Palette) != 2 || !black(img.Palette[0]) || !whiteColor(img.Palette[1]) {
				return errors.New("rowenc: bilevel rows must be of two colours, black then white")
			}
			white = func(y int) []byte {
				src := img.Pix[(y-b.Min.Y)*img.Stride+b.Min.X-img.Rect.Min.X:]
				return src[:r.width]
			}
		case *image.Gray:
			if len(r.bits) < r.width {
				r.bits = make([]byte, r.width)
			}
			white = func(y int) []byte {
				src := img.Pix[(y-b.Min.Y)*img.Stride+b.Min.X-img.Rect.Min.X:]
				src = src[:r.width]
				for i, v := range src {
					if v >= 128 {
						r.bits[i] = 1
					} else {
						r.bits[i] = 0
					}
				}
				return r.bits[:r.width]
			}
		default:
			return fmt.Errorf("rowenc: bilevel rows are %T", rows)
		}
		pack = func(y int) error {
			src := white(y)
			out := r.row[:r.rowBytes()]
			clear(out)
			for i, v := range src {
				if v != 0 {
					out[i>>3] |= 0x80 >> (i & 7)
				}
			}
			return nil
		}
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		if err := pack(y); err != nil {
			return err
		}
		if err := fn(r.row[:r.rowBytes()]); err != nil {
			return err
		}
		r.written++
	}
	return nil
}

// done reports an error if the rows written are not all of the image.
func (r *rower) done() error {
	if r.written != r.height {
		return fmt.Errorf("rowenc: %d of %d rows written", r.written, r.height)
	}
	return nil
}

func black(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return r == 0 && g == 0 && b == 0
}

func whiteColor(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return r == 0xffff && g == 0xffff && b == 0xffff
}

// New returns a writer of the format: "png" or "tiff".
func New(w io.Writer, format string, width, height int, kind Kind, opts *Options) (Writer, error) {
	switch format {
	case "png":
		return NewPNG(w, width, height, kind, opts)
	case "tiff":
		return NewTIFF(w, width, height, kind, opts)
	}
	return nil, fmt.Errorf("rowenc: unknown format %q", format)
}

// Encode writes a whole image in a format (see New): the image is the
// rows.
func Encode(w io.Writer, format string, img image.Image, kind Kind, opts *Options) error {
	b := img.Bounds()
	enc, err := New(w, format, b.Dx(), b.Dy(), kind, opts)
	if err != nil {
		return err
	}
	if err := enc.WriteRows(img); err != nil {
		return err
	}
	return enc.Close()
}

// KindOf returns the Kind of an image the renderer draws (see
// imagebdf.Mono): an *image.Gray is Gray, an *image.Paletted of black
// and white is Bilevel, the rest RGBA (converted when they are not an
// *image.RGBA).
func KindOf(img image.Image) Kind {
	switch img := img.(type) {
	case *image.Gray:
		return Gray
	case *image.Paletted:
		if len(img.Palette) == 2 && black(img.Palette[0]) && whiteColor(img.Palette[1]) {
			return Bilevel
		}
	}
	return RGBA
}
