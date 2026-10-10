package formula_test

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"net/http"
	"net/http/httptest"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/image/formula"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

// A formula drawn into an image, on the baseline of a line of text.
func Example() {
	ts, err := formula.New(nil)
	if err != nil {
		panic(err)
	}
	f := formula.ParseTeX(`x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}`)
	l := ts.Layout(f, formula.Style{Size: 24, Display: true, Color: bdf.RGB(0, 0, 0)})
	obj := l.Object()

	img, err := imagebdf.Object(obj, 2, nil)
	if err != nil {
		panic(err)
	}
	dst := image.NewRGBA(image.Rect(0, 0, 480, 200))
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	x, y := 20, 100 // where the baseline starts
	b := img.Bounds()
	draw.Draw(dst, b.Add(image.Pt(x, y)), img, b.Min, draw.Over)

	fmt.Println(f.Text())
	fmt.Println(b.Min.Y < 0, b.Max.Y > 0)
	// Output:
	// x=(−b±√(b^2−4ac))/(2a)
	// true true
}

// A formula served as a document of one page, for a browser to draw.
func ExampleLayout_Document() {
	ts, err := formula.New(nil)
	if err != nil {
		panic(err)
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		l := ts.Layout(formula.ParseTeX(r.URL.Query().Get("tex")), formula.Style{Size: 20, Display: true})
		w.Header().Set("Content-Type", "application/octet-stream")
		l.Document(8).WriteSingle(w)
	}
	rec := httptest.NewRecorder()
	handler(rec, httptest.NewRequest("GET", `/formula?tex=e^{i\pi}%2B1=0`, nil))
	r, err := bdf.OpenSingle(bytes.NewReader(rec.Body.Bytes()), int64(rec.Body.Len()))
	if err != nil {
		panic(err)
	}
	doc, err := r.ToDocument()
	if err != nil {
		panic(err)
	}
	st, _ := doc.SearchText()
	fmt.Println(len(doc.Views[0].Pages), st.Views[0].Pages[0].Text)
	// Output:
	// 1 e^(iπ)+1=0
}
