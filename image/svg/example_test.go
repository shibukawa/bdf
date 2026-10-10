package svg_test

import (
	"fmt"

	"github.com/shibukawa/bdf/image/svg"
)

func Example() {
	doc, err := svg.Parse([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 50">
  <style>.warm { fill: #f80 }</style>
  <g transform="translate(10 5) scale(2)">
    <rect class="warm" width="20" height="1em"/>
    <path id="mark" d="M0 0h10v10z" style="fill: rgb(0 128 255 / 50%)"/>
  </g>
</svg>`), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("view box:", svg.ParseNumbers(doc.Root.Attr["viewBox"]))
	g := doc.Root.Children[1]
	fmt.Println("transform:", svg.ParseTransform(g.Attr["transform"]))
	for _, n := range g.Children {
		fill, _ := n.Declared("fill")
		c, _ := svg.ParseColor(fill, svg.Color{A: 1})
		fmt.Printf("%s: fill %.2f %.2f %.2f %.2f\n", n.Name, c.R, c.G, c.B, c.A)
	}
	h, _ := svg.ParseLength(g.Children[0].Attr["height"], 0, 16)
	fmt.Println("height of the rect:", h)
	// path data is the string of the document
	fmt.Println("path data:", doc.IDs["mark"].Attr["d"])
	// Output:
	// view box: [0 0 100 50]
	// transform: [2 0 0 2 10 5]
	// rect: fill 1.00 0.53 0.00 1.00
	// path: fill 0.00 0.50 1.00 0.50
	// height of the rect: 16
	// path data: M0 0h10v10z
}
