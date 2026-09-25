// Package fixture builds the sample document used by tests and demos.
package fixture

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"

	"github.com/shibukawa/bdf"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

//go:embed fonts/*.ttf
var fontFS embed.FS

// Measurer computes text advances from an sfnt font.
type Measurer struct {
	font *sfnt.Font
	buf  sfnt.Buffer
}

// NewMeasurer parses a TTF/OTF.
func NewMeasurer(data []byte) (*Measurer, error) {
	f, err := sfnt.Parse(data)
	if err != nil {
		return nil, err
	}
	return &Measurer{font: f}, nil
}

// Advance returns the width of s at the given size in units (no kerning).
func (m *Measurer) Advance(s string, size float32) float32 {
	ppem := fixed.Int26_6(math.Round(float64(size) * 64))
	var total fixed.Int26_6
	for _, r := range s {
		g, err := m.font.GlyphIndex(&m.buf, r)
		if err != nil {
			continue
		}
		a, err := m.font.GlyphAdvance(&m.buf, g, ppem, 0)
		if err != nil {
			continue
		}
		total += a
	}
	return float32(total) / 64
}

// Fonts holds the embedded fixture fonts.
type Fonts struct {
	Regular, Bold       bdf.Hash
	RegularM, BoldM     *Measurer
	RegularRef, BoldRef bdf.Font
}

// LoadFonts adds the fixture fonts to d.
func LoadFonts(d *bdf.Document) (*Fonts, error) {
	reg, err := fontFS.ReadFile("fonts/DejaVuSans-sub.ttf")
	if err != nil {
		return nil, err
	}
	bold, err := fontFS.ReadFile("fonts/DejaVuSans-Bold-sub.ttf")
	if err != nil {
		return nil, err
	}
	f := &Fonts{Regular: d.AddFont(reg), Bold: d.AddFont(bold)}
	if f.RegularM, err = NewMeasurer(reg); err != nil {
		return nil, err
	}
	if f.BoldM, err = NewMeasurer(bold); err != nil {
		return nil, err
	}
	f.RegularRef = bdf.EmbeddedFont(f.Regular, 400, bdf.StyleNormal)
	f.RegularRef.Family = "sans-serif"
	f.BoldRef = bdf.EmbeddedFont(f.Bold, 700, bdf.StyleNormal)
	f.BoldRef.Family = "sans-serif"
	return f, nil
}

// checkerPNG renders a small test image.
func checkerPNG(size int) []byte {
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			c := color.NRGBA{uint8(255 * x / size), uint8(255 * y / size), 160, 255}
			if (x/8+y/8)%2 == 0 {
				c = color.NRGBA{c.R / 2, c.G / 2, 90, 255}
			}
			img.SetNRGBA(x, y, c)
		}
	}
	var b bytes.Buffer
	_ = png.Encode(&b, img)
	return b.Bytes()
}

// Demo builds a document with a slide deck, a flow document and a sheet.
func Demo() (*bdf.Document, error) {
	d := bdf.NewDocument()
	d.Meta.Title = "BDF fixture"
	d.Meta.Source = "fixture"
	fonts, err := LoadFonts(d)
	if err != nil {
		return nil, err
	}
	img := d.AddImage(checkerPNG(64))

	// A small icon reused everywhere (USE_AT).
	icon := bdf.NewObject()
	star := &bdf.Path{}
	for i := 0; i < 10; i++ {
		r := float32(12)
		if i%2 == 1 {
			r = 5
		}
		a := float64(i) * math.Pi / 5
		x, y := 12+r*float32(math.Sin(a)), 12-r*float32(math.Cos(a))
		if i == 0 {
			star.MoveTo(x, y)
		} else {
			star.LineTo(x, y)
		}
	}
	star.Close()
	sp := icon.AddPath(star)
	icon.FillColor(bdf.RGB(0xf5, 0xa6, 0x23)).FillPath(sp, bdf.NonZero)
	icon.StrokeColor(bdf.RGB(0x7a, 0x4a, 0x00)).Line(1, bdf.CapRound, bdf.JoinRound, 10).StrokePath(sp)
	iconH, iconBB := d.AddObject(icon)

	buildSlides(d, fonts, img, iconH, iconBB)
	buildFlow(d, fonts, iconH, iconBB)
	buildSheet(d, fonts, img)
	for _, v := range d.Views {
		if _, err := d.BuildTextIndex(v); err != nil {
			return nil, err
		}
	}
	return d, nil
}

func buildSlides(d *bdf.Document, fonts *Fonts, img bdf.Hash, iconH bdf.Hash, iconBB bdf.Rect) {
	const w, h = 960, 540
	v := d.NewView("slides", bdf.ViewFixed, "Slides")

	master := bdf.NewObject()
	bg := master.AddPaint(bdf.LinearGradient(0, 0, 0, h,
		bdf.Stop{Offset: 0, Color: bdf.RGB(0xf7, 0xf9, 0xfc)},
		bdf.Stop{Offset: 1, Color: bdf.RGB(0xdc, 0xe6, 0xf2)}))
	master.FillPaint(bg).FillRect(0, 0, w, h)
	master.FillColor(bdf.RGB(0x1f, 0x3a, 0x5f)).FillRect(0, h-36, w, 36)
	ic := master.AddObject(iconH, iconBB)
	for i := 0; i < 3; i++ {
		master.UseAt(ic, w-120+float32(i)*32, h-30)
	}
	fr := master.AddFont(fonts.RegularRef)
	master.Font(fr, 12).FillColor(bdf.RGB(0xff, 0xff, 0xff)).
		FillText("BDF fixture deck", 24, h-14, fonts.RegularM.Advance("BDF fixture deck", 12))
	masterH, masterBB := d.AddObject(master)

	titles := []string{"Shapes and paths", "Images, clips and blending", "Text styles"}
	for i, title := range titles {
		body := bdf.NewObject()
		fb := body.AddFont(fonts.BoldRef)
		fr := body.AddFont(fonts.RegularRef)
		body.Mark(bdf.MarkBox, "title").Font(fb, 36).FillColor(bdf.RGB(0x1f, 0x3a, 0x5f)).
			FillText(title, 48, 80, fonts.BoldM.Advance(title, 36))
		num := fmt.Sprintf("%d / %d", i+1, len(titles))
		body.Mark(bdf.MarkBox, "page-number").Font(fr, 12).FillColor(bdf.RGB(0xff, 0xff, 0xff)).TextStyle(bdf.AlignRight, bdf.BaselineAlphabetic, bdf.DirInherit, 0).
			FillText(num, w-140, h-14, fonts.RegularM.Advance(num, 12)).
			TextStyle(bdf.AlignLeft, bdf.BaselineAlphabetic, bdf.DirInherit, 0)

		switch i {
		case 0:
			rr := body.AddPath((&bdf.Path{}).RoundRect(60, 130, 260, 160, 18))
			body.FillColor(bdf.RGB(0x4c, 0x9b, 0xe8)).FillPath(rr, bdf.NonZero)
			body.StrokeColor(bdf.RGB(0x1f, 0x3a, 0x5f)).Line(4, bdf.CapRound, bdf.JoinRound, 10).
				Dash([]float32{12, 8}, 0).StrokePath(rr).Dash(nil, 0)
			rg := body.AddPaint(bdf.RadialGradient(470, 200, 10, 470, 210, 90,
				bdf.Stop{Offset: 0, Color: bdf.RGB(0xff, 0xf2, 0xb0)},
				bdf.Stop{Offset: 1, Color: bdf.RGB(0xe0, 0x7a, 0x1f)}))
			circle := body.AddPath((&bdf.Path{}).Circle(470, 210, 80))
			body.FillPaint(rg).FillPath(circle, bdf.NonZero)
			wave := &bdf.Path{}
			wave.MoveTo(600, 280).CubicTo(680, 120, 760, 360, 880, 160).LineTo(880, 300).Close()
			wp := body.AddPath(wave)
			body.FillColor(bdf.RGBA(0x2e, 0xa0, 0x6c, 0xa0)).FillPath(wp, bdf.NonZero)
			body.StrokeColor(bdf.RGB(0x1c, 0x60, 0x40)).Line(3, bdf.CapButt, bdf.JoinMiter, 10).StrokePath(wp)
			// even-odd ring
			ring := &bdf.Path{}
			ring.Circle(200, 400, 70).Circle(200, 400, 40)
			rp := body.AddPath(ring)
			body.FillColor(bdf.RGB(0x9b, 0x59, 0xb6)).FillPath(rp, bdf.EvenOdd)
			body.Save().Translate(500, 400).Transform(1, 0, 0.5, 1, 0, 0).
				FillColor(bdf.RGB(0xe7, 0x4c, 0x3c)).FillRect(-60, -40, 120, 80).
				StrokeColor(bdf.RGB(0x40, 0x10, 0x10)).Line(2, bdf.CapButt, bdf.JoinBevel, 10).StrokeRect(-60, -40, 120, 80).Restore()
			body.Font(fr, 16).FillColor(bdf.RGB(0x33, 0x33, 0x33)).
				FillText("dashed round rect, radial circle, bezier, even-odd ring, skewed rect", 48, 500,
					fonts.RegularM.Advance("dashed round rect, radial circle, bezier, even-odd ring, skewed rect", 16))
		case 1:
			im := body.AddImage(img)
			body.Image(im, 60, 130, 200, 200)
			body.Smoothing(false, 0).Image(im, 290, 130, 200, 200).Smoothing(true, 2)
			clip := body.AddPath((&bdf.Path{}).Circle(640, 230, 100))
			body.Save().ClipPath(clip, bdf.NonZero).Image(im, 540, 130, 200, 200).Restore()
			body.Save().ClipRect(780, 130, 120, 200).ImageSub(im, 16, 16, 32, 32, 760, 130, 200, 200).Restore()
			body.Save().Alpha(0.5).FillColor(bdf.RGB(0xff, 0x00, 0x00)).FillRect(60, 360, 160, 120).
				Blend(bdf.BlendMultiply).FillColor(bdf.RGB(0x00, 0x80, 0xff)).FillRect(140, 400, 160, 120).Restore()
			body.Save().Shadow(bdf.RGBA(0, 0, 0, 0x80), 12, 6, 6).FillColor(bdf.RGB(0xff, 0xff, 0xff)).
				FillRect(400, 370, 200, 120).Restore()
			body.GroupBegin(0.6, bdf.BlendSourceOver, 640, 360, 260, 150).
				FillColor(bdf.RGB(0x2e, 0xa0, 0x6c)).FillRect(660, 380, 140, 100).
				FillColor(bdf.RGB(0xe7, 0x4c, 0x3c)).FillRect(740, 410, 140, 90).GroupEnd()
			pat := body.AddPaint(bdf.Pattern(im, bdf.RepeatBoth))
			body.FillPaint(pat).FillRect(60, 500, 840, 20)
		case 2:
			y := float32(140)
			for _, sz := range []float32{12, 18, 24, 36} {
				s := fmt.Sprintf("Regular %gpt — The quick brown fox", sz)
				body.Font(fr, sz).FillColor(bdf.RGB(0x22, 0x22, 0x22)).FillText(s, 60, y, fonts.RegularM.Advance(s, sz))
				y += sz * 1.4
			}
			s := "Bold 28pt with stroke"
			body.Font(fb, 28).FillColor(bdf.RGB(0x4c, 0x9b, 0xe8)).FillText(s, 60, y+20, fonts.BoldM.Advance(s, 28))
			body.StrokeColor(bdf.RGB(0x1f, 0x3a, 0x5f)).Line(1, bdf.CapButt, bdf.JoinMiter, 10).StrokeText(s, 60, y+20, fonts.BoldM.Advance(s, 28))
			// Deliberately wrong advance to show the correction (rendered narrower).
			s = "advance correction: squeezed to 60%"
			body.Font(fr, 18).FillColor(bdf.RGB(0xc0, 0x39, 0x2b)).FillText(s, 60, y+70, fonts.RegularM.Advance(s, 18)*0.6)
			s = "letterSpacing 4pt"
			body.TextStyle(bdf.AlignLeft, bdf.BaselineAlphabetic, bdf.DirInherit, 4).FillText(s, 60, y+110, 0).
				TextStyle(bdf.AlignLeft, bdf.BaselineAlphabetic, bdf.DirInherit, 0)
			sysf := body.AddFont(bdf.SystemFont("serif", 400, bdf.StyleItalic))
			s = "system serif italic (no correction)"
			body.Font(sysf, 18).FillColor(bdf.RGB(0x22, 0x22, 0x22)).FillText(s, 60, y+150, 0)
			body.Link(60, y+130, 300, 24, "https://example.com/")
		}
		bodyH, _ := d.AddObject(body)
		v.AddPage(w, h, bdf.Layer{Role: bdf.RoleMaster, Obj: masterH}, bdf.Layer{Role: bdf.RoleBody, Obj: bodyH})
	}
	_ = masterBB
}

func buildFlow(d *bdf.Document, fonts *Fonts, iconH bdf.Hash, iconBB bdf.Rect) {
	const w, h = 595.3, 841.9
	const mx, my = 72, 72
	v := d.NewView("doc", bdf.ViewFlow, "Document")
	v.Continuous = &bdf.Continuous{Gap: 24}

	header := bdf.NewObject()
	hf := header.AddFont(fonts.BoldRef)
	header.Font(hf, 10).FillColor(bdf.RGB(0x66, 0x66, 0x66)).FillText("BDF fixture — flow view", mx, 48, fonts.BoldM.Advance("BDF fixture — flow view", 10))
	header.StrokeColor(bdf.RGB(0xcc, 0xcc, 0xcc)).Line(0.5, bdf.CapButt, bdf.JoinMiter, 10)
	rule := header.AddPath((&bdf.Path{}).MoveTo(mx, 56).LineTo(w-mx, 56))
	header.StrokePath(rule)
	headerH, _ := d.AddObject(header)

	footer := bdf.NewObject()
	ic := footer.AddObject(iconH, iconBB)
	footer.UseAt(ic, w/2-12, h-52)
	footerH, _ := d.AddObject(footer)

	lorem := []string{
		"BDF is a display-list format for browsers. Each page is a list of",
		"objects, and each object is a stream of instructions that map directly",
		"onto the Canvas 2D API. Fonts and images are stored as the browser",
		"already understands them, so the decoder stays small.",
		"",
		"Repeated content such as this header and footer is stored once and",
		"referenced by hash from every page. The body of each page is its own",
		"object, laid out inside the body rectangle so that a viewer can stack",
		"bodies vertically for a continuous reading mode.",
	}
	for page := 0; page < 2; page++ {
		body := bdf.NewObject()
		fb := body.AddFont(fonts.BoldRef)
		fr := body.AddFont(fonts.RegularRef)
		y := float32(my + 20)
		title := fmt.Sprintf("Section %d", page+1)
		body.Mark(bdf.MarkParagraph, "heading").Font(fb, 20).FillColor(bdf.RGB(0x1f, 0x3a, 0x5f)).FillText(title, mx, y, fonts.BoldM.Advance(title, 20))
		y += 32
		body.Font(fr, 11).FillColor(bdf.RGB(0x22, 0x22, 0x22))
		for rep := 0; rep < 3; rep++ {
			body.Mark(bdf.MarkParagraph, "")
			for _, line := range lorem {
				if line == "" {
					body.Mark(bdf.MarkParagraph, "")
				} else {
					body.Mark(bdf.MarkLine, "").FillText(line, mx, y, fonts.RegularM.Advance(line, 11))
				}
				y += 15
			}
			y += 8
		}
		// simple table
		body.StrokeColor(bdf.RGB(0x99, 0x99, 0x99)).Line(0.75, bdf.CapButt, bdf.JoinMiter, 10)
		tx, ty := float32(mx), y
		cw, rh := float32(120), float32(18)
		for r := 0; r < 4; r++ {
			for c := 0; c < 3; c++ {
				if r == 0 {
					body.FillColor(bdf.RGB(0xe8, 0xee, 0xf6)).FillRect(tx+float32(c)*cw, ty+float32(r)*rh, cw, rh)
				}
				body.StrokeRect(tx+float32(c)*cw, ty+float32(r)*rh, cw, rh)
				cell := fmt.Sprintf("R%dC%d", r+1, c+1)
				body.Mark(bdf.MarkCell, fmt.Sprintf("%c%d", 'A'+c, r+1)).FillColor(bdf.RGB(0x22, 0x22, 0x22)).FillText(cell, tx+float32(c)*cw+6, ty+float32(r)*rh+13, fonts.RegularM.Advance(cell, 11))
			}
		}
		bodyH, _ := d.AddObject(body)
		p := v.AddPage(w, h, bdf.Layer{Role: bdf.RoleHeader, Obj: headerH}, bdf.Layer{Role: bdf.RoleBody, Obj: bodyH}, bdf.Layer{Role: bdf.RoleFooter, Obj: footerH})
		p.Body = &bdf.RectDef{X: mx, Y: my, W: w - 2*mx, H: h - 2*my}
	}
}

func buildSheet(d *bdf.Document, fonts *Fonts, img bdf.Hash) {
	const tile = 2048
	const colW, rowH = 64, 20
	const nCols, nRows = 26, 200
	v := d.NewView("sheet1", bdf.ViewSheet, "Sheet1")
	v.Tile = tile
	v.Cols = []bdf.Run{{nCols, colW}}
	v.Rows = []bdf.Run{{nRows, rowH}}
	v.Freeze = &bdf.Freeze{Cols: 1, Rows: 1}
	v.Gridlines = true
	v.Tiles = map[string]string{}

	// A chart image spanning the tile boundary (rows 95..110 => y 1900..2200).
	chart := bdf.NewObject()
	im := chart.AddImage(img)
	chart.FillColor(bdf.RGB(0xff, 0xff, 0xff)).FillRect(0, 0, 320, 300)
	chart.StrokeColor(bdf.RGB(0x88, 0x88, 0x88)).Line(1, bdf.CapButt, bdf.JoinMiter, 10).StrokeRect(0.5, 0.5, 319, 299)
	chart.Image(im, 10, 10, 300, 280)
	chartH, chartBB := d.AddObject(chart)
	const chartX, chartY = 2*colW + 8, 95 * rowH

	tilesY := int(math.Ceil(float64(nRows*rowH) / tile))
	for ty := 0; ty < tilesY; ty++ {
		t := bdf.NewObject()
		fr := t.AddFont(fonts.RegularRef)
		fb := t.AddFont(fonts.BoldRef)
		oy := float32(ty * tile)
		r0, r1 := ty*tile/rowH, min(nRows, (ty+1)*tile/rowH+1)
		for r := r0; r < r1; r++ {
			y := float32(r*rowH) - oy
			for c := 0; c < nCols; c++ {
				x := float32(c * colW)
				var s string
				var f bdf.FontRef = fr
				switch {
				case r == 0:
					s = fmt.Sprintf("Col %c", 'A'+c)
					f = fb
					t.FillColor(bdf.RGB(0xe8, 0xee, 0xf6)).FillRect(x, y, colW, rowH)
				case c == 0:
					s = fmt.Sprintf("Row %d", r+1)
					f = fb
				default:
					s = fmt.Sprintf("%d", (r+1)*(c+1))
				}
				m := fonts.RegularM
				if f == fb {
					m = fonts.BoldM
				}
				adv := m.Advance(s, 10)
				t.Mark(bdf.MarkCell, fmt.Sprintf("%c%d", 'A'+c, r+1)).Font(f, 10).FillColor(bdf.RGB(0x22, 0x22, 0x22))
				if f == fr {
					t.TextStyle(bdf.AlignRight, bdf.BaselineAlphabetic, bdf.DirInherit, 0).FillText(s, x+colW-4, y+14, adv).
						TextStyle(bdf.AlignLeft, bdf.BaselineAlphabetic, bdf.DirInherit, 0)
				} else {
					t.FillText(s, x+4, y+14, adv)
				}
			}
		}
		if chartY < (ty+1)*tile && chartY+int(chartBB.H) > ty*tile {
			cr := t.AddObject(chartH, chartBB)
			t.UseAt(cr, chartX, float32(chartY)-oy)
		}
		th, _ := d.AddObject(t)
		v.Tiles[fmt.Sprintf("0,%d", ty)] = th.String()
	}
}
