package main

import (
	"github.com/shibukawa/bdf"
	"os"
)

func main() {
	o := bdf.NewObject()
	o.Save().Shadow(0x000000c0, 12, 18, 7).FillColor(0xf03a21ff).FillRect(90.25, 110.5, 40, 35).Restore()
	o.GroupBegin(.7, bdf.BlendMultiply, 0, 0, 256, 256).FillColor(0x32acdcff).FillRect(100.5, 100.25, 65, 65).MaskBegin(bdf.MaskAlpha, 0, nil).FillColor(0xffffffbf).FillRect(105.25, 105.5, 60, 60).MaskEnd().GroupEnd()
	o.Save().ClipRect(110.5, 0, 60, 256).Alpha(.6).FillColor(0x61b725ff).FillRect(100, 125.5, 70, 10).Restore()
	doc := bdf.NewDocument()
	h, _ := doc.AddObject(o)
	v := doc.NewView("effects", bdf.ViewFixed, "effects")
	v.AddPage(256, 256, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	f, e := os.Create(os.Args[1])
	if e != nil {
		panic(e)
	}
	if e = doc.WriteSingle(f); e != nil {
		panic(e)
	}
	if e = f.Close(); e != nil {
		panic(e)
	}
}
