package tiff

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/shibukawa/bdf/image/imgconv"
)

// Documents embed TIFF pictures (PowerPoint, Excel and Word blips, Visio
// foreign data). They show the first page, as it is stored: the
// Orientation tag is not applied, and the picture keeps its resolution
// (the document gives its size).

// FirstPage opens a TIFF file held in memory and returns its first page.
func FirstPage(data []byte) (*IFD, error) {
	f, err := Open(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	pages := f.Pages()
	if len(pages) == 0 {
		return nil, errors.New("tiff: the file has no pages")
	}
	return pages[0], nil
}

// MaxPicturePixels bounds the pixels of the pictures Picture returns: they
// are stored as large as they are, for the viewer to decode, and the
// converters and the rasterizer decode no larger ones.
var MaxPicturePixels int64 = imgconv.MaxDecodePixels

// Picture returns the first page of a TIFF file as an image file that
// browsers decode: the page's JPEG strips joined into one JPEG when they
// allow it (IFD.JPEG), otherwise its pixels, decoded and stored in the
// smallest form: PNG, or JPEG at opts.Quality when the page is JPEG
// (imgconv.EncodePixels in Keep mode). The caller stores it through
// imgconv.Optimize like any other picture, which converts it under
// Convert mode. damaged reports pixel data that did not all decode.
func Picture(data []byte, opts imgconv.Options) (pic []byte, damaged bool, err error) {
	d, err := FirstPage(data)
	if err != nil {
		return nil, false, err
	}
	if w, h := d.Size(); int64(w)*int64(h) > MaxPicturePixels {
		return nil, false, fmt.Errorf("tiff: a picture of %d×%d pixels is too large", w, h)
	}
	if b, ok := d.JPEG(); ok {
		return b, false, nil
	}
	img, damaged, err := d.Decode()
	if err != nil {
		return nil, false, err
	}
	opts.Mode = imgconv.Keep
	r, _ := imgconv.EncodePixels(img, d.Compression() != CompressionJPEG, opts) // Keep mode does not fail
	return r.Data, damaged, nil
}
