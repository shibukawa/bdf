package formula_test

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/formula"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

const quadratic = `x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}`

func typesetter(t *testing.T, opts *formula.Options) *formula.Typesetter {
	t.Helper()
	ts, err := formula.New(opts)
	if err != nil {
		t.Fatal(err)
	}
	return ts
}

// ops counts the instructions of an object by opcode.
func ops(t *testing.T, o *bdf.Object) map[byte]int {
	t.Helper()
	part, err := bdf.DecodeObject(o.Encode())
	if err != nil {
		t.Fatal(err)
	}
	n := map[byte]int{}
	if err := part.Walk(func(in bdf.Instr) { n[in.Op]++ }); err != nil {
		t.Fatal(err)
	}
	return n
}

// inked counts the pixels of an image that something is drawn on.
func inked(img *image.RGBA) int {
	n := 0
	for i := 3; i < len(img.Pix); i += 4 {
		if img.Pix[i] > 127 {
			n++
		}
	}
	return n
}

// save writes an image to look at when FORMULA_SAMPLES names a directory.
func save(t *testing.T, name string, img image.Image) {
	dir := os.Getenv("FORMULA_SAMPLES")
	if dir == "" {
		return
	}
	f, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
}

func TestLayout(t *testing.T) {
	ts := typesetter(t, nil)
	f := formula.ParseTeX(quadratic)
	if got, want := f.Text(), "x=(−b±√(b^2−4ac))/(2a)"; got != want {
		t.Errorf("text %q, want %q", got, want)
	}
	inline := ts.Layout(f, formula.Style{Size: 20})
	disp := ts.Layout(f, formula.Style{Size: 20, Display: true})
	if inline.Width() <= 0 || inline.Height() <= 0 || inline.Depth() <= 0 {
		t.Fatalf("inline: %v × %v + %v", inline.Width(), inline.Height(), inline.Depth())
	}
	if disp.Height()+disp.Depth() <= inline.Height()+inline.Depth() {
		t.Errorf("display style %v is not taller than text style %v", disp.Height()+disp.Depth(), inline.Height()+inline.Depth())
	}
	// the size scales everything
	big := ts.Layout(f, formula.Style{Size: 40})
	if r := big.Width() / inline.Width(); r < 1.999 || r > 2.001 {
		t.Errorf("twice the size is %v times as wide", r)
	}
	if disp.Text() != f.Text() {
		t.Errorf("text of the layout %q", disp.Text())
	}
}

// The object is paths only (no text, no font, no other part), with the
// baseline's left end at its origin, and draws as it is.
func TestObject(t *testing.T) {
	ts := typesetter(t, nil)
	l := ts.Layout(formula.ParseTeX(quadratic), formula.Style{Size: 20, Display: true})
	o := l.Object()
	if deps := o.Deps(); len(deps) != 0 {
		t.Errorf("the object refers to %d parts", len(deps))
	}
	n := ops(t, o)
	if n[bdf.OpFillPathRun] == 0 || n[bdf.OpFillRect] != 2 || n[bdf.OpFillText] != 0 || n[bdf.OpFont] != 0 {
		t.Errorf("instructions: %d glyph runs, %d rules (the fraction and the radical), %d texts, %d fonts",
			n[bdf.OpFillPathRun], n[bdf.OpFillRect], n[bdf.OpFillText], n[bdf.OpFont])
	}
	if ink := l.Ink(); o.BBox != ink || ink.Y >= 0 || ink.Y+ink.H <= 0 {
		t.Errorf("bounding box %v, ink %v", o.BBox, ink)
	}
	img, err := imagebdf.Object(o, 2, nil)
	if err != nil {
		t.Fatal(err)
	}
	save(t, "quadratic.png", img)
	b := img.Bounds()
	if b.Min.Y >= 0 || b.Max.Y <= 0 || b.Dx() < int(2*l.Width())-2 {
		t.Errorf("image bounds %v for a formula %v wide, %v high, %v deep", b, l.Width(), l.Height(), l.Depth())
	}
	if k := inked(img); k < b.Dx()*b.Dy()/50 || k > b.Dx()*b.Dy()/2 {
		t.Errorf("%d of %d pixels drawn", k, b.Dx()*b.Dy())
	}
	if c := img.RGBAAt(b.Min.X, b.Min.Y); c.A != 0 {
		t.Errorf("the corner is %v, want transparent", c)
	}
}

// Draw places the formula in an object of the caller's, such as a page.
func TestDraw(t *testing.T) {
	ts := typesetter(t, nil)
	l := ts.Layout(formula.ParseTeX(`\textcolor{red}{E} = mc^2 \quad \boxed{\cancel{x}}`), formula.Style{Size: 20, Color: bdf.RGBA(0, 0, 255, 255)})
	page := bdf.NewObject()
	l.Draw(page, 10, 30)
	if b := page.BBox; b.X < 9 || b.Y > 30 || b.Y+b.H < 30 || b.W < float32(l.Width())-1 {
		t.Errorf("bounding box %v of a formula %v wide drawn at (10, 30)", b, l.Width())
	}
	n := ops(t, page)
	if n[bdf.OpStrokePath] != 1 || n[bdf.OpFillRect] != 4 || n[bdf.OpSave] != 0 {
		t.Errorf("%d strokes (the strike), %d rules (the box), %d saves", n[bdf.OpStrokePath], n[bdf.OpFillRect], n[bdf.OpSave])
	}
	img, err := imagebdf.Object(page, 2, &imagebdf.Options{Background: color.White})
	if err != nil {
		t.Fatal(err)
	}
	save(t, "draw.png", img)
	red, blue := 0, 0
	for i := 0; i < len(img.Pix); i += 4 {
		switch {
		case img.Pix[i] > 200 && img.Pix[i+2] < 60:
			red++
		case img.Pix[i+2] > 200 && img.Pix[i] < 60:
			blue++
		}
	}
	if red == 0 || blue <= red {
		t.Errorf("%d red pixels (the E) and %d blue ones (the rest)", red, blue)
	}
}

func TestMathML(t *testing.T) {
	f, err := formula.ParseMathML(`<p>Euler: <math display="block"><msup><mi>e</mi><mrow><mi>i</mi><mi>π</mi></mrow></msup><mo>+</mo><mn>1</mn><mo>=</mo><mn>0</mn></math></p>`)
	if err != nil {
		t.Fatal(err)
	}
	if !f.Display() || f.Text() != "e^(iπ)+1=0" {
		t.Errorf("display %v, text %q", f.Display(), f.Text())
	}
	if l := typesetter(t, nil).Layout(f, formula.Style{}); l.Width() <= 0 {
		t.Errorf("width %v", l.Width())
	}
	if _, err := formula.ParseMathML(`<p>no formula</p>`); err == nil {
		t.Error("no math element: no error")
	}
}

// Characters the formula font lacks are drawn with the text fonts.
func TestTextFonts(t *testing.T) {
	jp, err := os.ReadFile("../converter/docx/testdata/fonts/MPLUS1p-Regular-subset.ttf")
	if err != nil {
		t.Fatal(err)
	}
	f := formula.ParseTeX(`\text{速さ} = \frac{\text{距離}}{\text{時間}}`)
	without := typesetter(t, nil).Layout(f, formula.Style{Size: 20})
	with := typesetter(t, &formula.Options{Text: [][]byte{jp}}).Layout(f, formula.Style{Size: 20})
	a, err := imagebdf.Object(without.Object(), 2, nil)
	if err != nil {
		t.Fatal(err)
	}
	b, err := imagebdf.Object(with.Object(), 2, nil)
	if err != nil {
		t.Fatal(err)
	}
	save(t, "text.png", b)
	if inked(b) < 2*inked(a) {
		t.Errorf("%d pixels with the text font, %d without", inked(b), inked(a))
	}

	if _, err := formula.New(&formula.Options{Math: jp}); err == nil {
		t.Error("a formula font without a MATH table: no error")
	}
	if _, err := formula.New(&formula.Options{Text: [][]byte{[]byte("not a font")}}); err == nil {
		t.Error("a text font that is no font: no error")
	}
	stix, err := os.ReadFile("../converter/docx/testdata/fonts/STIXTwoMath-subset.ttf")
	if err != nil {
		t.Fatal(err)
	}
	own := typesetter(t, &formula.Options{Math: stix}).Layout(formula.ParseTeX(quadratic), formula.Style{Size: 20})
	if own.Width() <= 0 {
		t.Errorf("a formula font of the caller's: width %v", own.Width())
	}
}

func TestConcurrent(t *testing.T) {
	ts := typesetter(t, nil)
	f := formula.ParseTeX(quadratic)
	want := ts.Layout(f, formula.Style{Size: 20}).Object().Encode()
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 20 {
				if got := ts.Layout(f, formula.Style{Size: 20}).Object().Encode(); string(got) != string(want) {
					t.Error("the object differs")
					return
				}
			}
		}()
	}
	wg.Wait()
}

// Document is a page of the formula, with the formula as its text.
func TestDocument(t *testing.T) {
	ts := typesetter(t, nil)
	l := ts.Layout(formula.ParseTeX(quadratic), formula.Style{Size: 20, Display: true})
	doc := l.Document(10)
	v := doc.Views[0]
	ink := l.Ink()
	if p := v.Pages[0]; p.W != ink.W+20 || p.H != ink.H+20 {
		t.Errorf("page %v × %v for an ink of %v × %v", p.W, p.H, ink.W, ink.H)
	}
	img, err := imagebdf.New(doc, &imagebdf.Options{NoSystemFonts: true}).Page(v, 0, 2)
	if err != nil {
		t.Fatal(err)
	}
	save(t, "document.png", img)
	if dir := os.Getenv("FORMULA_SAMPLES"); dir != "" {
		// to draw with the browser's renderer: node test/render.mjs document.bdf out
		f, err := os.Create(filepath.Join(dir, "document.bdf"))
		if err != nil {
			t.Fatal(err)
		}
		if err := doc.WriteSingle(f); err != nil {
			t.Fatal(err)
		}
		f.Close()
	}
	// white paper, the margin empty, the formula inside it
	dark := image.Rectangle{}
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			if c := img.RGBAAt(x, y); c.R < 128 {
				dark = dark.Union(image.Rect(x, y, x+1, y+1))
			}
		}
	}
	if in := img.Bounds().Inset(18); dark.Empty() || !dark.In(in) || dark.Dx() < in.Dx()*8/10 {
		t.Errorf("the formula is drawn in %v of %v", dark, img.Bounds())
	}
	st, err := doc.SearchText()
	if err != nil {
		t.Fatal(err)
	}
	if len(st.Views) != 1 || len(st.Views[0].Pages) != 1 || st.Views[0].Pages[0].Text != l.Text() {
		t.Errorf("text of the document %+v, want %q", st.Views, l.Text())
	}
}
