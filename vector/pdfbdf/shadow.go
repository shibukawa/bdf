package pdfbdf

import (
	"crypto/sha256"
	"image"
	"image/color"
	"math"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/gauss"
	"github.com/shibukawa/bdf/raster/imagebdf"
	"github.com/shibukawa/bdf/vector/internal/export"
)

// A PDF has no shadows: the shadow of what is drawn is written as a
// picture of it. What casts the shadow is drawn by itself with the image
// renderer (raster/imagebdf), as the viewer would draw it; its alpha is
// blurred as the renderer blurs shadows, and written as the soft mask of an
// image of the color of the shadow, which is drawn under what casts it: at
// the offset of the shadow, in the directions of the page (a shadow is not
// turned with what casts it).

// Limits of the shadows of a document.
const (
	// maxShadowBlur bounds the blur of a shadow, in units of the page.
	maxShadowBlur = 1024
	// maxShadowPixels bounds the picture of a shadow, and
	// maxShadowsPixels those of all the shadows of a document.
	maxShadowPixels  = 4 << 20
	maxShadowsPixels = 256 << 20
	// shadowDPI is the most a picture of a shadow has: that of a shadow
	// without a blur, whose edges are those of what casts it.
	shadowDPI = 300
)

// reach is how far a line reaches beyond its path, in units of the state.
func (st *gstate) reach() float64 {
	return st.width / 2 * math.Max(math.Sqrt2, st.miter)
}

// shadows draws what casts the shadows of a document.
type shadows struct {
	doc    *bdf.Document
	r      *imagebdf.Renderer
	have   map[bdf.Hash]bool
	images map[[32]byte]*pdfImage
	pixels int
}

// part gives the document of the shadows a part of the document: a font
// or an image of what casts a shadow.
func (sh *shadows) part(ex *exporter, h bdf.Hash) {
	if sh.have[h] {
		return
	}
	sh.have[h] = true
	if p := ex.Doc.Part(h); p != nil {
		sh.doc.AddPart(p.Type, p.Data)
	}
}

// deps gives the document of the shadows the parts that an object refers
// to: its fonts, images and paths, and the objects it uses with theirs.
func (sh *shadows) deps(ex *exporter, o *bdf.ObjectPart) {
	for _, h := range o.Deps() {
		if sh.have[h] {
			continue
		}
		sh.part(ex, h)
		if p := ex.Doc.Part(h); p != nil && p.Type == bdf.PartObject {
			if child, err := bdf.DecodeObject(p.Data); err == nil {
				sh.deps(ex, child)
			}
		}
	}
}

// styleOf writes a fill or a stroke style of the state for an object of
// the shadows.
func (sh *shadows) styleOf(ex *exporter, o *bdf.Object, sty style, stroke bool) {
	if sty.paint == nil {
		if stroke {
			o.StrokeColor(sty.color)
		} else {
			o.FillColor(sty.color)
		}
		return
	}
	p := *sty.paint
	if p.Kind == bdf.PaintPattern {
		if sty.obj == nil || int(p.Image) >= len(sty.obj.Images) {
			return
		}
		h := sty.obj.Images[p.Image]
		sh.part(ex, h)
		p.Image = o.AddImage(h)
	}
	if stroke {
		o.StrokePaint(o.AddPaint(p))
	} else {
		o.FillPaint(o.AddPaint(p))
	}
}

// castShadow writes the shadow of what shape draws: box is around it in
// the space of the state (reach more on every side, for a line), and shape
// draws it on an object that has the state.
func (ex *exporter) castShadow(box [4]float64, reach float64, shape func(o *bdf.Object)) {
	s := ex.s
	st := &s.st
	if !st.shadowOn() {
		return
	}
	to := s.toPage.Mul(st.ctm)
	back, ok := to.Invert()
	if !ok || !to.Finite() {
		return
	}
	if ex.shadows == nil {
		sh := &shadows{doc: bdf.NewDocument(), have: map[bdf.Hash]bool{}, images: map[[32]byte]*pdfImage{}}
		sh.r = imagebdf.New(sh.doc, &imagebdf.Options{FontFS: ex.opts.FontFS, FontDirs: ex.opts.FontDirs, NoSystemFonts: ex.opts.NoSystemFonts, Background: color.Transparent})
		ex.shadows = sh
	}
	sh := ex.shadows
	// where the shadow may show, on the page: around what casts it, as
	// far as the blur reaches
	sigma := st.shadowBlur / 2
	pad := math.Ceil(3 * sigma)
	b := to.Bounds(box[0]-reach, box[1]-reach, box[2]+reach, box[3]+reach)
	x0, y0, x1, y1 := b[0]-pad, b[1]-pad, b[2]+pad, b[3]+pad
	if ex.pageW > 0 {
		// (nothing of a shadow shows outside the page)
		x0, y0 = math.Max(x0, -st.shadowDX-pad), math.Max(y0, -st.shadowDY-pad)
		x1, y1 = math.Min(x1, ex.pageW-st.shadowDX+pad), math.Min(y1, ex.pageH-st.shadowDY+pad)
	}
	if !(x1 > x0) || !(y1 > y0) || math.IsInf(x1-x0, 0) || math.IsInf(y1-y0, 0) {
		return
	}
	// the pixels of a unit: four for the width of the blur, which is as
	// much as a blurred edge shows
	res := float64(shadowDPI) / 72
	if sigma > 0 {
		res = math.Min(math.Max(4/sigma, 1), res)
	}
	if area := (x1 - x0) * (y1 - y0) * res * res; area > maxShadowPixels {
		res *= math.Sqrt(maxShadowPixels / area)
	}
	w, h := int(math.Ceil((x1-x0)*res)), int(math.Ceil((y1-y0)*res))
	if w < 1 || h < 1 {
		return
	}
	if sh.pixels+w*h > maxShadowsPixels {
		ex.Warnf("the shadows of a document have more than %d pixels in all: the rest of them are not drawn", maxShadowsPixels)
		return
	}
	sh.pixels += w * h
	x1, y1 = x0+float64(w)/res, y0+float64(h)/res

	// what casts the shadow, by itself, on a page of the place of the shadow
	o := bdf.NewObject()
	m := matrix{1, 0, 0, 1, -x0, -y0}.Mul(to)
	o.Transform(float32(m[0]), float32(m[1]), float32(m[2]), float32(m[3]), float32(m[4]), float32(m[5]))
	sh.styleOf(ex, o, st.fill, false)
	sh.styleOf(ex, o, st.stroke, true)
	o.Line(float32(st.width), st.cap, st.join, float32(st.miter))
	if len(st.dash) > 0 {
		o.Dash(st.dash, st.dashOff)
	}
	o.Alpha(float32(st.alpha))
	o.Smoothing(st.smooth, 0)
	if st.font != nil {
		if st.font.Kind == bdf.FontEmbedded {
			sh.part(ex, st.font.Hash)
		}
		o.Font(o.AddFont(*st.font), float32(st.size))
	}
	o.TextStyle(st.align, st.baseline, st.dir, float32(st.spacing))
	shape(o)
	pw, ph := float32(x1-x0), float32(y1-y0)
	o.SetBBox(0, 0, pw, ph)
	obj, _ := sh.doc.AddObject(o)
	view := &bdf.View{ID: "shadow", Kind: bdf.ViewFixed, Pages: []*bdf.Page{{W: pw, H: ph, Layers: []bdf.Layer{{Role: bdf.RoleBody, Obj: obj}}}}}
	img, err := sh.r.PageSize(view, 0, w, h)
	if err != nil {
		ex.Warnf("a shadow is not drawn: %v", err)
		return
	}
	// its alpha, blurred
	plane := make([]float32, w*h)
	any := false
	for y := 0; y < h; y++ {
		row := img.Pix[y*img.Stride : y*img.Stride+4*w]
		for x := 0; x < w; x++ {
			a := row[4*x+3]
			plane[y*w+x] = float32(a) / 255
			any = any || a != 0
		}
	}
	if !any {
		return
	}
	if sigma > 0 {
		gauss.Blur(plane, w, h, sigma*res)
	}
	c := uint32(st.shadow)
	pix := image.NewNRGBA(image.Rect(0, 0, w, h))
	for i, a := range plane {
		pix.Pix[4*i], pix.Pix[4*i+1], pix.Pix[4*i+2] = byte(c>>24), byte(c>>16), byte(c>>8)
		pix.Pix[4*i+3] = byte(math.Round(float64(min(max(a, 0), 1)) * 255))
	}
	// (the shadows of what is drawn again, on every page, are one image)
	key := sha256.Sum256(pix.Pix)
	key[0], key[1] = key[0]^byte(w), key[1]^byte(h)
	out := sh.images[key]
	if out == nil {
		out = &pdfImage{name: ex.newName("Im"), ref: ex.pixels(pix, ex.interpolate()), w: float64(w), h: float64(h)}
		sh.images[key] = out
	}
	// (in a tagged file a shadow is an artifact: what casts it is content
	// again after it)
	var content *structElem
	if t := ex.tags; t != nil && t.off == 0 {
		if s.mc.open {
			content = s.mc.elem
		}
		ex.tagArtifact()
	}
	// under what casts it, in the space of the page: with the alpha of
	// the color of the shadow and the blend mode of the state
	ex.blendMode()
	ex.scoped(func() {
		if a := alphaOf(st.shadow); a != st.pCa {
			s.op("%s gs", ex.extGState("/ca "+export.Fine(a)))
		}
		place := matrix{x1 - x0, 0, 0, -(y1 - y0), x0 + st.shadowDX, y1 + st.shadowDY}
		s.op("%s cm %s cm %s Do", back.Exact(), export.Mat(place), s.use("XObject", out.name, out.ref))
	})
	if content != nil {
		s.setMC(ex.tags, content)
	}
}

// runShadow writes the shadow of a FILL_PATH_RUN.
func (ex *exporter) runShadow(rule byte, glyphs []bdf.Glyph) {
	box := [4]float64{math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)}
	type placed struct {
		p    *bdf.Path
		x, y float32
	}
	var run []placed
	for _, g := range glyphs {
		bp := ex.Path(ex.obj, int(g.Path))
		if bp == nil {
			continue
		}
		b, ok := pathBounds(bp)
		if !ok {
			continue
		}
		run = append(run, placed{bp, g.X, g.Y})
		x, y := float64(g.X), float64(g.Y)
		box[0], box[1], box[2], box[3] = math.Min(box[0], b[0]+x), math.Min(box[1], b[1]+y), math.Max(box[2], b[2]+x), math.Max(box[3], b[3]+y)
	}
	if len(run) == 0 {
		return
	}
	ex.castShadow(box, 0, func(o *bdf.Object) {
		refs := map[*bdf.Path]bdf.PathRef{}
		out := make([]bdf.Glyph, len(run))
		for i, g := range run {
			ref, ok := refs[g.p]
			if !ok {
				ref = o.AddPath(g.p)
				refs[g.p] = ref
			}
			out[i] = bdf.Glyph{Path: ref, X: g.x, Y: g.y}
		}
		o.FillPathRun(rule, out)
	})
}

// pathBounds returns a box around a path: around its points and the
// points that steer its curves. It has none without points.
func pathBounds(bp *bdf.Path) (box [4]float64, ok bool) {
	box = [4]float64{math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)}
	add := func(x, y float64) {
		if math.IsNaN(x) || math.IsNaN(y) || math.IsInf(x, 0) || math.IsInf(y, 0) {
			return
		}
		box[0], box[1], box[2], box[3] = math.Min(box[0], x), math.Min(box[1], y), math.Max(box[2], x), math.Max(box[3], y)
		ok = true
	}
	a := bp.Args
	i := 0
	f := func(k int) float64 {
		if i+k >= len(a) {
			return 0
		}
		return float64(a[i+k])
	}
	for _, v := range bp.Verbs {
		switch v {
		case bdf.VerbMove, bdf.VerbLine:
			add(f(0), f(1))
			i += 2
		case bdf.VerbQuad:
			add(f(0), f(1))
			add(f(2), f(3))
			i += 4
		case bdf.VerbCubic:
			add(f(0), f(1))
			add(f(2), f(3))
			add(f(4), f(5))
			i += 6
		case bdf.VerbRect:
			add(f(0), f(1))
			add(f(0)+f(2), f(1)+f(3))
			i += 4
		case bdf.VerbEllipse:
			// (turned or not: within the circle of its longer radius)
			r := math.Max(math.Abs(f(2)), math.Abs(f(3)))
			add(f(0)-r, f(1)-r)
			add(f(0)+r, f(1)+r)
			i += 8
		case bdf.VerbArcTo:
			add(f(0), f(1))
			add(f(2), f(3))
			i += 5
		case bdf.VerbRoundRect:
			add(f(0), f(1))
			add(f(0)+f(2), f(1)+f(3))
			i += 5
		}
	}
	return box, ok
}
