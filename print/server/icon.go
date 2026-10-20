package server

import (
	"bytes"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"math"
	"net/http"
	"strconv"
	"sync"

	"github.com/shibukawa/bdf/image/rowenc"
	"github.com/shibukawa/bdf/print/ipp"
)

// iconSizes are the sizes of printer-icons: small, large, and the one
// AirPrint shows.
var iconSizes = []int{48, 128, 512}

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

// icon serves an icon of printer-icons.
func (x *ippServer) icon(w http.ResponseWriter, r *http.Request) {
	iconOnce.Do(func() {
		icons = map[int][]byte{}
		for _, n := range iconSizes {
			icons[n] = printerIcon(n)
		}
	})
	var n int
	if _, err := fmt.Sscanf(r.PathValue("file"), "%d.png", &n); err != nil || icons[n] == nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(icons[n])))
	w.Header().Set("Cache-Control", "max-age=86400")
	w.Write(icons[n])
}

var ippPage = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.Title}}</title>
<style>
:root{color-scheme:light dark;--fg:#1d2733;--bg:#fff;--muted:#5b6876;--line:#d5dbe1}
@media (prefers-color-scheme:dark){:root{--fg:#e6ebf0;--bg:#14191f;--muted:#9aa7b4;--line:#2e3946}}
body{font:15px/1.5 system-ui,sans-serif;color:var(--fg);background:var(--bg);margin:0;padding:24px 16px;max-width:760px;margin:auto}
h1{font-size:20px;margin:0 0 16px;display:flex;gap:12px;align-items:center}
table{border-collapse:collapse;width:100%}td,th{text-align:left;padding:8px 6px;border-bottom:1px solid var(--line)}
th{font-weight:600;color:var(--muted);font-size:13px}code{font-size:13px;word-break:break-all}
</style></head><body>
<h1><img src="/icons/48.png" width="32" height="32" alt="">{{.Title}}</h1>
<table><tr><th>Printer</th><th>State</th><th>Jobs</th><th>Address</th></tr>
{{range .Queues}}<tr><td>{{.Info}}{{if .Location}}<br><small>{{.Location}}</small>{{end}}</td><td>{{.State}}{{if .Message}}<br><small>{{.Message}}</small>{{end}}</td><td>{{.Jobs}}</td><td><code>{{.URI}}</code></td></tr>
{{end}}</table>
{{if .Secure}}<p><a href="/airprint.mobileconfig">Add these printers to an iPhone or iPad</a> that does not list them (a configuration profile: open it in Safari, then install it in Settings).</p>{{else}}<p>To add these printers to an iPhone or iPad that does not list them, open this page over <a href="https://{{.Host}}/">https</a>.</p>{{end}}
</body></html>
`))

// page shows the printers and their states: printer-more-info.
func (x *ippServer) page(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Info, Location, State, Message, URI string
		Jobs                                int
	}
	data := struct {
		Title, Host string
		Secure      bool
		Queues      []row
	}{Title: "bdf printd on " + x.host, Host: r.Host, Secure: isSecure(r)}
	scheme := "ipp"
	if isSecure(r) {
		scheme = "ipps"
	}
	for _, q := range x.order {
		if name := r.PathValue("queue"); name != "" && name != q.name {
			continue
		}
		st, _, msg, accepting := x.queueState(q.name)
		state := map[ipp.PrinterState]string{ipp.PrinterIdle: "idle", ipp.PrinterProcessing: "printing", ipp.PrinterStopped: "stopped"}[st]
		if !accepting {
			state += ", not accepting jobs"
		}
		n := 0
		if info, err := x.s.sp.Queue(q.name); err == nil {
			for _, c := range info.Jobs {
				n += c
			}
		}
		data.Queues = append(data.Queues, row{Info: q.info, Location: q.location, State: state, Message: msg, Jobs: n, URI: scheme + "://" + r.Host + q.path})
	}
	if len(data.Queues) == 0 {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; img-src 'self'; style-src 'unsafe-inline'")
	ippPage.Execute(w, data)
}
