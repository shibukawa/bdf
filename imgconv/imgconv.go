// Package imgconv decides how raster images are stored in a BDF document.
//
// In Keep mode images are stored as they are (decoded pixels become PNG).
// In Convert mode, images that browsers already decode efficiently (WebP,
// AVIF, SVG) are kept, and PNG/JPEG/GIF/BMP are re-encoded as WebP when
// that makes them smaller; otherwise the original bytes are kept. Lossless sources stay lossless unless they look like
// photographs. The codecs are pure Go (see internal/README.md); building
// with -tags bdf_noconv leaves them out and Convert behaves like Keep.
//
// Raster inputs whose size on the page is known (a scanned page) are also
// capped in resolution, in either mode: Options.FitSize says how far to scale
// one down (at most MaxDPI pixels per inch and MaxPixels pixels) and Resize
// scales it.
package imgconv

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"

	"golang.org/x/image/bmp"
)

// Mode selects between keeping images and converting them.
type Mode int

const (
	// Keep stores images unchanged; decoded pixel data is written as PNG.
	Keep Mode = iota
	// Convert tries lossless and lossy WebP and keeps the smallest.
	Convert
)

// Options configures conversion.
type Options struct {
	Mode Mode
	// Quality for lossy WebP (and for JPEG, which EncodePixels writes in
	// Keep mode), 1-100 (default 80).
	Quality int
	// Method is the WebP effort, 0 (fast) to 6 (small); default 6.
	Method int
	// NoLossyForLossless keeps lossless sources lossless even when they look like photos.
	NoLossyForLossless bool
	// PhotoColors is the number of distinct sampled colours above which an
	// image counts as a photograph (default 4096).
	PhotoColors int

	// MaxDPI caps the resolution of the raster inputs that are scaled down
	// to fit (see FitSize), in pixels per inch of the page: 0 is
	// DefaultMaxDPI, a negative value removes the cap.
	MaxDPI float64
	// MaxPixels caps the pixel count (width × height) of those inputs: 0 is
	// DefaultMaxPixels, a negative value removes the cap.
	MaxPixels int
}

func (o Options) quality() int {
	if o.Quality <= 0 || o.Quality > 100 {
		return 80
	}
	return o.Quality
}

func (o Options) method() int {
	if o.Method < 0 || o.Method > 6 {
		return 6
	}
	if o.Method == 0 {
		return 6
	}
	return o.Method
}

func (o Options) photoColors() int {
	if o.PhotoColors <= 0 {
		return 4096
	}
	return o.PhotoColors
}

// Result is a stored image.
type Result struct {
	Data []byte
	// Format is "png", "jpeg", "webp", "avif", "gif", "bmp" or "svg".
	Format string
	// Converted reports whether Data differs from the input.
	Converted bool
	// Lossless reports whether Data reproduces the input pixels exactly.
	Lossless bool
}

// ErrNotAvailable is returned in Convert mode when the codecs were left out
// of the build (-tags bdf_noconv).
var ErrNotAvailable = errors.New("imgconv: codecs not compiled in (bdf_noconv)")

// Sniff returns the format of encoded image bytes, or "".
func Sniff(b []byte) string {
	switch {
	case len(b) >= 8 && bytes.Equal(b[:8], []byte("\x89PNG\r\n\x1a\n")):
		return "png"
	case len(b) >= 3 && b[0] == 0xff && b[1] == 0xd8 && b[2] == 0xff:
		return "jpeg"
	case len(b) >= 12 && bytes.Equal(b[:4], []byte("RIFF")) && bytes.Equal(b[8:12], []byte("WEBP")):
		return "webp"
	case len(b) >= 12 && bytes.Equal(b[4:8], []byte("ftyp")) && (bytes.Equal(b[8:12], []byte("avif")) || bytes.Equal(b[8:12], []byte("avis"))):
		return "avif"
	case len(b) >= 6 && (bytes.Equal(b[:6], []byte("GIF87a")) || bytes.Equal(b[:6], []byte("GIF89a"))):
		return "gif"
	case len(b) >= 2 && b[0] == 'B' && b[1] == 'M':
		return "bmp"
	case bytes.Contains(b[:min(len(b), 512)], []byte("<svg")):
		return "svg"
	}
	return ""
}

// MIME returns the media type for a format name.
func MIME(format string) string {
	switch format {
	case "png":
		return "image/png"
	case "jpeg":
		return "image/jpeg"
	case "webp":
		return "image/webp"
	case "avif":
		return "image/avif"
	case "gif":
		return "image/gif"
	case "bmp":
		return "image/bmp"
	case "svg":
		return "image/svg+xml"
	}
	return "application/octet-stream"
}

// Optimize takes encoded image bytes and returns the representation to store.
func Optimize(data []byte, opts Options) (Result, error) {
	format := Sniff(data)
	keep := Result{Data: data, Format: format, Lossless: format != "jpeg"}
	if opts.Mode != Convert {
		return keep, nil
	}
	if !Available() {
		return keep, ErrNotAvailable
	}
	var img image.Image
	var err error
	lossless := true
	switch format {
	case "jpeg":
		img, err = jpeg.Decode(bytes.NewReader(data))
		lossless = false
	case "png":
		img, err = png.Decode(bytes.NewReader(data))
	case "gif":
		img, err = gif.Decode(bytes.NewReader(data))
	case "bmp":
		img, err = bmp.Decode(bytes.NewReader(data))
	default:
		// webp, avif, svg and unknown formats are kept as they are.
		return keep, nil
	}
	if err != nil {
		return keep, err
	}
	best, err := encodeBest(toNRGBA(img), lossless, opts, len(data))
	if err != nil {
		return keep, err
	}
	if best == nil {
		return keep, nil
	}
	best.Converted = true
	return *best, nil
}

// EncodeImage encodes decoded pixels. lossless says whether they come from a
// lossless source (so that lossy re-encoding is only used for photographs).
func EncodeImage(img *image.NRGBA, lossless bool, opts Options) (Result, error) {
	pngData := encodePNG(img)
	keep := Result{Data: pngData, Format: "png", Lossless: true}
	if opts.Mode != Convert {
		return keep, nil
	}
	if !Available() {
		return keep, ErrNotAvailable
	}
	best, err := encodeBest(img, lossless, opts, len(pngData))
	if err != nil {
		return keep, err
	}
	if best == nil {
		return keep, nil
	}
	best.Converted = true
	return *best, nil
}

// encodeBest returns the smallest candidate that beats limit bytes, or nil.
func encodeBest(img *image.NRGBA, lossless bool, opts Options, limit int) (*Result, error) {
	var best *Result
	consider := func(data []byte, format string, ll bool) {
		if len(data) == 0 || len(data) >= limit {
			return
		}
		if best == nil || len(data) < len(best.Data) {
			best = &Result{Data: data, Format: format, Lossless: ll}
		}
	}
	tryLossy := !lossless || (!opts.NoLossyForLossless && PhotoLike(img, opts.photoColors()))
	if lossless {
		data, err := encodeWebP(img, 75, opts.method(), true)
		if err != nil {
			return nil, err
		}
		consider(data, "webp", true)
	}
	if tryLossy {
		data, err := encodeWebP(img, opts.quality(), opts.method(), false)
		if err != nil {
			return nil, err
		}
		consider(data, "webp", false)
	}
	return best, nil
}

// PhotoLike reports whether an image has more than maxColors distinct colours
// in a sample of its pixels, which makes lossy encoding a reasonable choice.
func PhotoLike(img *image.NRGBA, maxColors int) bool {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w*h < 64*64 {
		return false
	}
	step := 1
	for (w/step)*(h/step) > 65536 {
		step++
	}
	seen := make(map[uint32]struct{}, maxColors+1)
	for y := 0; y < h; y += step {
		for x := 0; x < w; x += step {
			i := img.PixOffset(b.Min.X+x, b.Min.Y+y)
			c := uint32(img.Pix[i])<<16 | uint32(img.Pix[i+1])<<8 | uint32(img.Pix[i+2])
			seen[c] = struct{}{}
			if len(seen) > maxColors {
				return true
			}
		}
	}
	return false
}

func toNRGBA(img image.Image) *image.NRGBA {
	if n, ok := img.(*image.NRGBA); ok {
		return n
	}
	b := img.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(out, out.Bounds(), img, b.Min, draw.Src)
	return out
}

func encodePNG(img *image.NRGBA) []byte {
	var buf bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	_ = enc.Encode(&buf, img)
	return buf.Bytes()
}
