// Command scene writes shared path workloads for native renderer comparisons.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/shibukawa/bdf"
)

type operation struct {
	Color    [4]uint8    `json:"color"`
	Rule     string      `json:"rule"`
	Stroke   float32     `json:"stroke"`
	Segments [][]float32 `json:"segments"`
}
type scene struct {
	Version       int    `json:"version"`
	Name          string `json:"name"`
	Width, Height int
	Operations    []operation `json:"operations"`
}

func main() {
	output := flag.String("output", "/tmp/bdf-native-scenes", "output directory")
	flag.Parse()
	must(os.MkdirAll(*output, 0755))
	for _, name := range []string{"opaque-rectangles", "cubic-paths", "translucent-overlap", "stroked-curves"} {
		s := scene{Version: 1, Name: name, Width: 1024, Height: 1024}
		count := 1000
		if name == "translucent-overlap" {
			count = 300
		}
		for i := 0; i < count; i++ {
			x := float32((i*37)%960) + .25
			y := float32((i*61)%960) + .125
			w := float32(16 + i%48)
			h := float32(12 + i%53)
			op := operation{Color: [4]uint8{uint8(30 + i%200), uint8(20 + (i*3)%220), uint8(10 + (i*7)%230), 255}, Rule: "nonzero"}
			switch name {
			case "opaque-rectangles":
				op.Segments = [][]float32{{0, x, y}, {1, x + w, y}, {1, x + w, y + h}, {1, x, y + h}, {4}}
			case "cubic-paths":
				op.Segments = [][]float32{{0, x, y}, {3, x + w*.2, y - h*.5, x + w*1.3, y + h*.3, x + w, y + h}, {3, x + w*.7, y + h*.3, x - w*.2, y + h*.8, x, y}, {4}}
			case "translucent-overlap":
				x = 512 + float32(math.Sin(float64(i)*.37)*240)
				y = 512 + float32(math.Cos(float64(i)*.43)*240)
				w = 180
				h = 150
				op.Color[3] = uint8(30 + i%160)
				if i%2 == 0 {
					op.Rule = "evenodd"
				}
				op.Segments = [][]float32{{0, x - w, y - h}, {1, x + w, y + h}, {1, x - w, y + h}, {1, x + w, y - h}, {4}}
			case "stroked-curves":
				op.Stroke = 1.25 + float32(i%4)
				op.Segments = [][]float32{{0, x, y}, {3, x + w*.2, y - h*.5, x + w*1.3, y + h*.3, x + w, y + h}}
			}
			s.Operations = append(s.Operations, op)
		}
		data, err := json.Marshal(s)
		must(err)
		must(os.WriteFile(filepath.Join(*output, name+".json"), append(data, '\n'), 0644))
		f, err := os.Create(filepath.Join(*output, name+".txt"))
		must(err)
		fmt.Fprintln(f, s.Width, s.Height, len(s.Operations))
		obj := bdf.NewObject()
		for _, op := range s.Operations {
			fmt.Fprintln(f, op.Color[0], op.Color[1], op.Color[2], op.Color[3], op.Rule, op.Stroke, len(op.Segments))
			p := &bdf.Path{}
			for _, seg := range op.Segments {
				for _, v := range seg {
					fmt.Fprint(f, v, " ")
				}
				fmt.Fprintln(f)
				switch int(seg[0]) {
				case 0:
					p.MoveTo(seg[1], seg[2])
				case 1:
					p.LineTo(seg[1], seg[2])
				case 2:
					p.QuadTo(seg[1], seg[2], seg[3], seg[4])
				case 3:
					p.CubicTo(seg[1], seg[2], seg[3], seg[4], seg[5], seg[6])
				case 4:
					p.Close()
				}
			}
			color := bdf.RGBA(op.Color[0], op.Color[1], op.Color[2], op.Color[3])
			ref := obj.AddPath(p)
			if op.Stroke > 0 {
				obj.StrokeColor(color).Line(op.Stroke, bdf.CapButt, bdf.JoinMiter, 10).StrokePath(ref)
			} else {
				rule := byte(bdf.NonZero)
				if op.Rule == "evenodd" {
					rule = bdf.EvenOdd
				}
				obj.FillColor(color).FillPath(ref, rule)
			}
		}
		must(f.Close())
		doc := bdf.NewDocument()
		hash, _ := doc.AddObject(obj)
		v := doc.NewView("scene", bdf.ViewFixed, name)
		v.AddPage(float32(s.Width), float32(s.Height), bdf.Layer{Role: bdf.RoleBody, Obj: hash})
		file, err := os.Create(filepath.Join(*output, name+".bdf"))
		must(err)
		must(doc.WriteSingle(file))
		must(file.Close())
	}
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
