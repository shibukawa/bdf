package ipp

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"math"
	"net/http"
	"strconv"
	"sync"

	"github.com/shibukawa/bdf/image/rowenc"
)

// The icon of the printers bdf serves (bdf printd's queues, the Docflow's
// CUPS server nodes): printer-icons, at /icons/<size>.png of the server.

// IconSizes are the sizes of printer-icons: small, large, and the one
// AirPrint shows.
var IconSizes = []int{48, 128, 512}

var (
	iconOnce sync.Once
	icons    map[int][]byte
)

// printerIcon draws the icon of the printers: a printer with a sheet in
// it, in the colours of the viewer.
func printerIcon(n int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, n, n))
	s := float64(n) / 64
	// signed distance to a rounded rectangle, for edges smoothed over a
	// pixel
	box := func(px, py, x0, y0, x1, y1, r float64) float64 {
		cx, cy := (x0+x1)/2, (y0+y1)/2
		hx, hy := (x1-x0)/2-r, (y1-y0)/2-r
		dx, dy := math.Max(math.Abs(px-cx)-hx, 0), math.Max(math.Abs(py-cy)-hy, 0)
		return math.Hypot(dx, dy) - r
	}
	paint := func(x, y int, c color.RGBA, d float64) {
		a := math.Min(math.Max(0.5-d/s, 0), 1)
		if a == 0 {
			return
		}
		o := img.PixOffset(x, y)
		for i, v := range []uint8{c.R, c.G, c.B, 255} {
			old := float64(img.Pix[o+i])
			img.Pix[o+i] = uint8(old*(1-a) + float64(v)*a + 0.5)
		}
	}
	body := color.RGBA{0x37, 0x47, 0x5a, 255}
	paper := color.RGBA{0xff, 0xff, 0xff, 255}
	line := color.RGBA{0x2b, 0x7c, 0xd3, 255}
	tray := color.RGBA{0x24, 0x2f, 0x3c, 255}
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			px, py := (float64(x)+0.5)/s, (float64(y)+0.5)/s
			paint(x, y, paper, box(px, py, 18, 6, 46, 26, 1.5))
			paint(x, y, body, box(px, py, 6, 22, 58, 48, 5))
			paint(x, y, tray, box(px, py, 14, 40, 50, 44, 1))
			paint(x, y, paper, box(px, py, 18, 42, 46, 60, 1.5))
			for _, ly := range []float64{48, 52, 56} {
				paint(x, y, line, box(px, py, 22, ly, 42, ly+1.5, 0.7))
			}
			paint(x, y, color.RGBA{0x5c, 0xd6, 0x8a, 255}, math.Hypot(px-50, py-29)-2.2)
		}
	}
	var b bytes.Buffer
	if err := rowenc.Encode(&b, "png", img, rowenc.RGBA, nil); err != nil {
		panic(err)
	}
	return b.Bytes()
}

// ServeIcon serves the icon file names ("128.png"), one of IconSizes.
func ServeIcon(w http.ResponseWriter, r *http.Request, file string) {
	iconOnce.Do(func() {
		icons = map[int][]byte{}
		for _, n := range IconSizes {
			icons[n] = printerIcon(n)
		}
	})
	var n int
	if _, err := fmt.Sscanf(file, "%d.png", &n); err != nil || icons[n] == nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(icons[n])))
	w.Header().Set("Cache-Control", "max-age=86400")
	w.Write(icons[n])
}
