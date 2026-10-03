// Command example shows the two ways of drawing BDF with Ebitengine: the
// page of a document, drawn into an image by raster/imagebdf, and formulas
// laid out by package formula, drawn as vector paths at any size.
//
//	go run ./example                       # with the page of testdata/docx/math.bdf
//	go run ./example -doc file.bdf -page 2
//	go run ./example -shot out.png         # save the first frame and quit
package main

import (
	"errors"
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/formula"
	"github.com/shibukawa/bdf/raster/ebitenginebdf"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

const width, height = 960, 600

type game struct {
	page      *ebiten.Image
	formulas  []*ebitenginebdf.Drawing
	zoomed    *ebitenginebdf.Drawing
	frame     int
	drawn     int // frames drawn: Update may run more than once before the first Draw
	shot      string
	offscreen *ebiten.Image
}

// page draws a page of a document into an image as high as the window.
func page(path string, n int) (*ebiten.Image, error) {
	r, err := bdf.OpenSingleFile(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	doc, err := r.ToDocument()
	if err != nil {
		return nil, err
	}
	if len(doc.Views) == 0 || n >= len(doc.Views[0].Pages) {
		return nil, errors.New("the document has no such page")
	}
	v := doc.Views[0]
	scale := float64(height-40) / float64(v.Pages[n].H)
	img, err := imagebdf.New(doc, nil).Page(v, n, scale)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

// formulas lays formulas out and compiles them into paths.
func formulas() (list []*ebitenginebdf.Drawing, zoomed *ebitenginebdf.Drawing, err error) {
	ts, err := formula.New(nil)
	if err != nil {
		return nil, nil, err
	}
	compile := func(tex string, st formula.Style) *ebitenginebdf.Drawing {
		if err != nil {
			return nil
		}
		var d *ebitenginebdf.Drawing
		d, err = ebitenginebdf.Compile(ts.Layout(formula.ParseTeX(tex), st).Object())
		return d
	}
	white := formula.Style{Size: 28, Display: true, Color: bdf.RGB(255, 255, 255)}
	for _, tex := range []string{
		`x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}`,
		`\int_0^\infty e^{-x^2}\,dx = \frac{\sqrt{\pi}}{2}`,
		`A = \begin{pmatrix} a_{11} & a_{12} \\ a_{21} & a_{22} \end{pmatrix}`,
		`\textcolor{orange}{\sum_{k=1}^{n} k^2} = \frac{n(n+1)(2n+1)}{6}`,
	} {
		list = append(list, compile(tex, white))
	}
	zoomed = compile(`e^{i\pi} + 1 = 0`, formula.Style{Size: 20, Color: bdf.RGB(120, 220, 255)})
	return list, zoomed, err
}

func (g *game) Update() error {
	g.frame++
	if g.shot != "" && g.drawn >= 2 {
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		g.offscreen.ReadPixels(img.Pix)
		f, err := os.Create(g.shot)
		if err != nil {
			return err
		}
		defer f.Close()
		if err := png.Encode(f, img); err != nil {
			return err
		}
		return ebiten.Termination
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	g.drawn++
	dst := g.offscreen
	dst.Fill(color.RGBA{24, 28, 40, 255})

	// the page: an image
	var op ebiten.DrawImageOptions
	op.GeoM.Translate(20, 20)
	dst.DrawImage(g.page, &op)

	// the formulas: paths, placed by the origin of each (the left end of
	// its baseline)
	x := float64(g.page.Bounds().Dx()) + 60
	y := 70.0
	for _, d := range g.formulas {
		var opts ebitenginebdf.DrawOptions
		opts.GeoM.Translate(x, y)
		opts.AntiAlias = true
		d.Draw(dst, &opts)
		y += 90
	}
	// one that grows, turns and fades: the same paths under a transform
	t := float64(g.frame) / 120
	var opts ebitenginebdf.DrawOptions
	b := g.zoomed.Bounds()
	opts.GeoM.Translate(-float64(b.X+b.W/2), -float64(b.Y+b.H/2))
	opts.GeoM.Scale(1.6+0.8*math.Sin(t), 1.6+0.8*math.Sin(t))
	opts.GeoM.Rotate(0.15 * math.Sin(t/2))
	opts.GeoM.Translate(x+170, y+80)
	opts.ColorScale.ScaleAlpha(float32(0.7 + 0.3*math.Cos(t)))
	opts.AntiAlias = true
	g.zoomed.Draw(dst, &opts)

	screen.DrawImage(dst, nil)
}

func (g *game) Layout(int, int) (int, int) { return width, height }

func main() {
	doc := flag.String("doc", "../../testdata/docx/math.bdf", "the document whose page is shown")
	n := flag.Int("page", 1, "the page shown (from 1)")
	shot := flag.String("shot", "", "save the first frame as a PNG file and quit")
	flag.Parse()

	g := &game{shot: *shot, offscreen: ebiten.NewImage(width, height)}
	var err error
	if g.page, err = page(*doc, *n-1); err != nil {
		log.Fatal(err)
	}
	if g.formulas, g.zoomed, err = formulas(); err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("BDF with Ebitengine")
	if err := ebiten.RunGame(g); err != nil && !errors.Is(err, ebiten.Termination) {
		log.Fatal(err)
	}
}
