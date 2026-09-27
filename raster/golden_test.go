package raster

import (
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// goldenCase is a case of the Chromium golden test (test/page/harness.ts).
type goldenCase struct {
	name, src, kind, view string
	page                  int
	viewport              bdf.Rect
	scale                 float64
}

// goldenCases reads the cases the browser harness draws, so that both
// renderers are compared on the same pages.
func goldenCases(t *testing.T) []goldenCase {
	src, err := os.ReadFile("../test/page/harness.ts")
	if err != nil {
		t.Skip(err)
	}
	line := regexp.MustCompile(`(?m)^\s*\{ name: "([^"]+)"(.*)\},?\s*$`)
	field := func(s, name string) string {
		m := regexp.MustCompile(name + `: "?([^",}]+)"?`).FindStringSubmatch(s)
		if m == nil {
			return ""
		}
		return strings.TrimSpace(m[1])
	}
	vp := regexp.MustCompile(`viewport: \{ x: ([-\d.]+), y: ([-\d.]+), w: ([-\d.]+), h: ([-\d.]+) \}`)
	num := func(s string) float64 {
		v, _ := strconv.ParseFloat(s, 64)
		return v
	}
	var out []goldenCase
	for _, m := range line.FindAllStringSubmatch(string(src), -1) {
		rest := m[2]
		if field(rest, "roles") != "" {
			continue // drawing only some layers is the viewer's option
		}
		c := goldenCase{name: m[1], src: field(rest, "src"), kind: field(rest, "kind"), view: field(rest, "view"), scale: num(field(rest, "scale"))}
		if c.src == "" {
			c.src = "/testdata/demo.bdf"
		}
		c.page = int(num(field(rest, "page")))
		if v := vp.FindStringSubmatch(rest); v != nil {
			c.viewport = bdf.Rect{X: float32(num(v[1])), Y: float32(num(v[2])), W: float32(num(v[3])), H: float32(num(v[4]))}
		}
		out = append(out, c)
	}
	if len(out) < 50 {
		t.Fatalf("read %d cases from the harness", len(out))
	}
	return out
}

func openDoc(t testing.TB, path string) *bdf.Document {
	t.Helper()
	r, err := bdf.OpenSingleFile(path)
	if err != nil {
		t.Fatal(err)
	}
	d, err := r.ToDocument()
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func findView(d *bdf.Document, id string) *bdf.View {
	for _, v := range d.Views {
		if v.ID == id {
			return v
		}
	}
	return nil
}

// testFonts are the fonts the test documents were laid out with; the
// drawings use no others, so that they are the same on every machine.
var testFonts = []string{"../converter/pptx/testdata/fonts", "../converter/docx/testdata/fonts", "../fixture/testdata/fonts"}

// diff compares two images after scaling both down by k: the mean absolute
// difference of the channels (0–255) and the share of pixels that differ
// by more than 48 in a channel.
func diff(a, b image.Image, k int) (mean, bad float64) {
	w := min(a.Bounds().Dx(), b.Bounds().Dx()) / k
	h := min(a.Bounds().Dy(), b.Bounds().Dy()) / k
	if w == 0 || h == 0 {
		return 255, 1
	}
	avg := func(img image.Image, x, y int) [3]float64 {
		var s [3]float64
		o := img.Bounds().Min
		for dy := 0; dy < k; dy++ {
			for dx := 0; dx < k; dx++ {
				r, g, bb, al := img.At(o.X+x*k+dx, o.Y+y*k+dy).RGBA()
				// over white, as the pages are
				wt := float64(0xffff - al)
				s[0] += float64(r) + wt
				s[1] += float64(g) + wt
				s[2] += float64(bb) + wt
			}
		}
		for i := range s {
			s[i] /= float64(k*k) * 257
		}
		return s
	}
	total, nbad := 0.0, 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p, q := avg(a, x, y), avg(b, x, y)
			worst := 0.0
			for i := range p {
				d := math.Abs(p[i] - q[i])
				total += d
				worst = math.Max(worst, d)
			}
			if worst > 48 {
				nbad++
			}
		}
	}
	return total / float64(3*w*h), float64(nbad) / float64(w*h)
}

// TestGoldenPages draws the pages of the browser's golden test and checks
// that they look like the browser's: close once both are scaled down to
// thumbnail detail. BDF_RASTER_OUT=dir writes the drawings there.
func TestGoldenPages(t *testing.T) {
	if testing.Short() {
		t.Skip("slow")
	}
	outDir := os.Getenv("BDF_RASTER_OUT")
	docs := map[string]*bdf.Document{}
	for _, c := range goldenCases(t) {
		t.Run(c.name, func(t *testing.T) {
			if why, ok := goldenSkips[c.name]; ok {
				t.Skip(why)
			}
			golden, err := readPNG(filepath.Join("../testdata/golden", c.name+".png"))
			if err != nil {
				t.Skip(err)
			}
			d := docs[c.src]
			if d == nil {
				d = openDoc(t, ".."+c.src)
				docs[c.src] = d
			}
			v := findView(d, c.view)
			if v == nil {
				t.Fatalf("no view %q", c.view)
			}
			r := New(d, &Options{FontDirs: testFonts, NoSystemFonts: true})
			var img *image.RGBA
			gb := golden.Bounds()
			switch {
			case c.kind == "page":
				img, err = r.Page(v, c.page, c.scale)
			case c.kind == "sheet" || c.kind == "continuous" && v.Kind == bdf.ViewScroll:
				img, err = r.Region(v, 0, c.viewport, gb.Dx(), gb.Dy())
			default:
				t.Skip("continuous flow views are the viewer's layout")
			}
			if err != nil {
				t.Fatal(err)
			}
			if outDir != "" {
				// this drawing on the left, the browser's on the right
				side := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx()+gb.Dx(), max(img.Bounds().Dy(), gb.Dy())))
				draw.Draw(side, img.Bounds(), img, image.Point{}, draw.Src)
				draw.Draw(side, gb.Add(image.Pt(img.Bounds().Dx(), 0)), golden, gb.Min, draw.Src)
				writePNG(t, filepath.Join(outDir, c.name+".png"), side)
			}
			if d := img.Bounds().Size().Sub(gb.Size()); d.X < -1 || d.X > 1 || d.Y < -1 || d.Y > 1 {
				t.Fatalf("size %v, golden %v", img.Bounds().Size(), gb.Size())
			}
			mean, bad := diff(img, golden, 4)
			t.Logf("mean %.2f, bad %.2f%% %v", mean, bad*100, r.Warnings())
			if lim := goldenLimit(c.name); mean > lim.mean || bad > lim.bad {
				t.Errorf("differs from the browser: mean %.2f (limit %.2f), %.2f%% pixels off (limit %.2f%%)", mean, lim.mean, bad*100, lim.bad*100)
			}
		})
	}
}

type limit struct{ mean, bad float64 }

// goldenLimits are the cases that may differ more from the browser's
// drawing than the others, and why.
var goldenLimits = map[string]limit{
	// a photograph scaled down 2.7 times: mipmaps here, bilinear sampling in the browser
	"tiff-scan-2": {10, 0.02},
	// a checkerboard image magnified twice: its edges fall half a pixel apart
	"hpgl-job-2": {8, 0.02},
	// Chinese and Korean text in fonts referred to by name, which the test fonts lack
	"pdf-cjk-1": {8, 0.08},
}

// goldenSkips are the cases the renderer does not draw.
var goldenSkips = map[string]string{
	"image-rotated": "AVIF images are not decoded",
}

// goldenLimit is how far a case may be from the browser's drawing.
func goldenLimit(name string) limit {
	if l, ok := goldenLimits[name]; ok {
		return l
	}
	return limit{mean: 4, bad: 0.03}
}

func readPNG(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

func writePNG(t testing.TB, path string, img image.Image) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
}
