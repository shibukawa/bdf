package music

import (
	"strconv"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/music/smufl"
)

// fontDefaults returns Bravura's engraving defaults.
func fontDefaults() engravingDefaults {
	d := smufl.Engraving()
	return engravingDefaults{
		staffLine: d.StaffLineThickness, stem: d.StemThickness, leger: d.LegerLineThickness,
		legerExt: d.LegerLineExtension, beam: d.BeamThickness, beamSpacing: d.BeamSpacing,
		thinBar: d.ThinBarlineThickness, thickBar: d.ThickBarlineThickness, barSep: d.ThinThickBarlineSeparation,
		tieEnd: d.TieEndpointThickness, tieMid: d.TieMidpointThickness, slurEnd: d.SlurEndpointThickness,
		slurMid: d.SlurMidpointThickness, tupletBracket: d.TupletBracketThickness, hairpin: d.HairpinThickness,
		octaveLine: d.OctaveLineThickness, pedalLine: d.PedalLineThickness,
	}
}

// glyphKey is a glyph at a size (pt per staff space, horizontally and
// vertically).
type glyphKey struct {
	name   string
	sx, sy float64
}

// glyphSet is the glyph outlines of a document: one path collection part
// that every page refers to.
type glyphSet struct {
	paths []*bdf.Path
	index map[glyphKey]int
}

// glyphCollection gathers the glyphs the systems draw, at their sizes.
func (e *engraver) glyphCollection(systems []*system) *glyphSet {
	gs := &glyphSet{index: map[glyphKey]int{}}
	add := func(p prim) {
		if p.kind != primGlyph {
			return
		}
		k := glyphKey{p.g.name, p.scale * p.sx * e.sp, p.scale * p.sy * e.sp}
		if _, ok := gs.index[k]; ok {
			return
		}
		gs.index[k] = len(gs.paths)
		gs.paths = append(gs.paths, p.g.outline.pathXY(k.sx, k.sy))
	}
	for _, sys := range systems {
		for _, st := range sys.staves {
			for _, p := range st.prims {
				add(p)
			}
		}
		for _, p := range e.systemPrims(sys) {
			add(p)
		}
	}
	return gs
}

// pathXY returns the outline scaled by sx horizontally and sy vertically.
func (o *outline) pathXY(sx, sy float64) *bdf.Path {
	p := &bdf.Path{Verbs: append([]byte(nil), o.verbs...), Args: make([]float32, len(o.args))}
	for i, v := range o.args {
		if i%2 == 0 {
			p.Args[i] = f32(v * sx)
		} else {
			p.Args[i] = f32(v * sy)
		}
	}
	return p
}

// writer writes primitives into the object of a page.
type writer struct {
	e      *engraver
	cv     *canvas.Canvas
	glyphs *glyphSet
	ext    bdf.Hash
	refs   map[int]bdf.PathRef
	font   *fontset.Use
	size   float64
	lineW  float64
	dashed bool
}

func (w *writer) obj() *bdf.Object { return w.cv.Obj }

// draw writes primitives placed at (ox, oy) pt.
func (w *writer) draw(prims []prim, ox, oy float64) {
	sp := w.e.sp
	o := w.obj()
	for _, p := range prims {
		switch p.kind {
		case primGlyph:
			k := glyphKey{p.g.name, p.scale * p.sx * sp, p.scale * p.sy * sp}
			i := w.glyphs.index[k]
			ref, ok := w.refs[i]
			if !ok {
				ref = o.AddExtPath(w.ext, uint32(i))
				w.refs[i] = ref
			}
			o.FillPathAt(ref, 0, f32(ox+p.x*sp), f32(oy+p.y*sp))
		case primRect:
			o.FillRect(f32(ox+p.x*sp), f32(oy+p.y*sp), f32(p.w*sp), f32(p.h*sp))
		case primFill:
			path := &bdf.Path{Verbs: p.out.verbs, Args: make([]float32, len(p.out.args))}
			for i, v := range p.out.args {
				if i%2 == 0 {
					path.Args[i] = f32(ox + v*sp)
				} else {
					path.Args[i] = f32(oy + v*sp)
				}
			}
			o.FillPath(o.AddPath(path), 0)
		case primLine:
			if lw := p.width * sp; lw != w.lineW {
				o.Line(f32(lw), 0, 0, 10)
				w.lineW = lw
			}
			if p.dash > 0 {
				o.Dash([]float32{f32(p.dash * sp), f32(p.dash * sp * 0.6)}, 0)
				w.dashed = true
			} else if w.dashed {
				o.Dash(nil, 0)
				w.dashed = false
			}
			path := (&bdf.Path{}).MoveTo(f32(ox+p.x*sp), f32(oy+p.y*sp)).LineTo(f32(ox+p.w*sp), f32(oy+p.h*sp))
			o.StrokePath(o.AddPath(path))
		case primText:
			w.text(p.text, ox+p.x*sp, oy+p.y*sp, p.align)
		case primMark:
			o.Mark(p.mark, p.alt)
		}
	}
}

// text writes a line of text with its baseline at y pt.
func (w *writer) text(t *textLine, x, y float64, align int) {
	o := w.obj()
	x -= t.width * float64(align) / 2
	for _, r := range t.runs {
		if r.s == "" {
			continue
		}
		lang := ""
		if r.cjk {
			for _, c := range r.s {
				if lang = fontset.RuneLang(c); lang != "" {
					break
				}
			}
			if lang == "" && fontset.Script(w.cv.Lang()) != "" {
				lang = w.cv.Lang()
			}
		}
		w.cv.SetLang(lang)
		u := r.fc.Use
		if w.font == nil || *w.font != u || w.size != t.size {
			o.Font(w.cv.Font(u), f32(t.size))
			w.font, w.size = &u, t.size
		}
		o.FillText(r.s, f32(x+r.x), f32(y), f32(r.w))
	}
	w.cv.Drawn = true
}

// drawTitle draws the title, subtitle and credits above the first system.
func (e *engraver) drawTitle(w *writer, title string, pageW float64) {
	sp := e.sp
	y := margin
	if title != "" {
		t := e.text.line(title, 4.4*sp, true, false)
		y += t.asc
		w.obj().Mark(bdf.MarkHeading, "1")
		w.text(t, pageW/2, y, 1)
		y += t.desc + sp
	}
	if e.s.Subtitle != "" {
		t := e.text.line(e.s.Subtitle, 2.8*sp, false, false)
		y += t.asc
		w.obj().Mark(bdf.MarkParagraph, "")
		w.text(t, pageW/2, y, 1)
		y += t.desc + 0.5*sp
	}
	if e.s.Composer != "" || e.s.Lyricist != "" {
		size := 2.2 * sp
		y += 0.8 * size
		if e.s.Lyricist != "" {
			w.obj().Mark(bdf.MarkParagraph, "")
			w.text(e.text.line(e.s.Lyricist, size, false, false), margin, y, 0)
		}
		if e.s.Composer != "" {
			w.obj().Mark(bdf.MarkParagraph, "")
			w.text(e.text.line(e.s.Composer, size, false, false), pageW-margin, y, 2)
		}
		y += 0.5 * size
	}
	if e.s.Arranger != "" {
		size := 2.2 * sp
		y += 0.8 * size
		w.obj().Mark(bdf.MarkParagraph, "")
		w.text(e.text.line(e.s.Arranger, size, false, false), pageW-margin, y, 2)
	}
}

// drawPageNumber draws the number of a page at its foot.
func (e *engraver) drawPageNumber(w *writer, n int, pageW, pageH float64) {
	t := e.text.line(strconv.Itoa(n), 2.0*e.sp, false, false)
	w.obj().Mark(bdf.MarkParagraph, "")
	w.text(t, pageW/2, pageH-margin/2, 1)
}

// drawSystem draws a system whose first staff's top line starts at (ox,
// oy) pt.
func (e *engraver) drawSystem(w *writer, sys *system, ox, oy float64) {
	sp := e.sp
	w.draw(e.systemPrims(sys), ox, oy)
	for si, st := range sys.staves {
		w.draw(st.prims, ox, oy+sys.y[si]*sp)
	}
}

// systemPrims draws what spans the staves of a system: bar lines,
// brackets, braces and part names, in the coordinates of the first staff.
func (e *engraver) systemPrims(sys *system) []prim {
	c := newDrawing()
	def := e.def
	last := len(e.staves) - 1
	// the parts' extents
	type span struct{ top, bottom float64 }
	parts := make([]span, len(e.s.Parts))
	for i := range parts {
		parts[i] = span{1e9, -1e9}
	}
	for si, ref := range e.staves {
		p := &parts[ref.part]
		p.top = min(p.top, sys.y[si])
		p.bottom = max(p.bottom, sys.y[si]+4)
	}
	// the line joining the staves at the start
	if last > 0 {
		c.rect(0, 0, def.thinBar, sys.y[last]+4)
	}
	for pi, p := range e.s.Parts {
		sp := parts[pi]
		if p.Staves > 1 {
			g := sym("brace")
			h := sp.bottom - sp.top
			c.prims = append(c.prims, prim{kind: primGlyph, g: g, x: -0.4 - g.advance*h/4*1.2, y: sp.bottom,
				scale: 1, sx: h / 4 * 1.2, sy: h / 4})
		}
	}
	if len(e.s.Parts) > 1 && !e.grand {
		top, bottom := 0.0, sys.y[last]+4
		x := -1.0
		c.rect(x-def.thickBar, top-0.2, def.thickBar, bottom-top+0.4)
		c.glyph(sym("bracketTop"), x-def.thickBar, top-0.2, 1)
		c.glyph(sym("bracketBottom"), x-def.thickBar, bottom+0.2, 1)
	}
	// part names
	if len(e.s.Parts) > 1 {
		for pi, p := range e.s.Parts {
			name := p.Abbrev
			if sys.first == 0 {
				name = p.Name
			}
			if name == "" {
				continue
			}
			t := e.text.line(name, 2.2*e.sp, false, false)
			sp := parts[pi]
			y := (sp.top+sp.bottom)/2 + (t.asc-t.desc)/2/e.sp
			c.markText(bdf.MarkParagraph, "")
			c.text(t, -e.bracketWidth()-1.0, y, 2, e.sp)
		}
	}
	// bar lines, joined through the staves of a part
	tops := e.staffTops(sys)
	for mi := sys.first; mi <= sys.last; mi++ {
		col := e.cols[mi]
		for pi := range e.s.Parts {
			sp := parts[pi]
			if lb := e.leftBar(mi); lb != BarNone {
				e.barline(c, lb, col.x, sp.top, sp.bottom, tops[pi], true)
			}
			e.barline(c, e.rightBar(mi), col.x+col.width, sp.top, sp.bottom, tops[pi], false)
		}
	}
	return c.prims
}

// staffTops returns the top lines of the staves of each part in a system.
func (e *engraver) staffTops(sys *system) [][]float64 {
	out := make([][]float64, len(e.s.Parts))
	for si, ref := range e.staves {
		out[ref.part] = append(out[ref.part], sys.y[si])
	}
	return out
}

// barline draws a bar line at x (its right edge; its left edge for a bar
// line at the start of a measure) from top to bottom, with repeat dots on
// each staff.
func (e *engraver) barline(c *drawing, b Barline, x, top, bottom float64, staves []float64, left bool) {
	def := e.def
	h := bottom - top
	thin, thick, sep := def.thinBar, def.thickBar, def.barSep
	dots := func(x float64) {
		for _, y := range staves {
			c.glyph(sym("repeatDots"), x, y+4, 1)
		}
	}
	dotW := sym("repeatDots").advance
	if left {
		switch b {
		case BarHeavyLight, BarRepeatStart:
			c.rect(x, top, thick, h)
			c.rect(x+thick+sep, top, thin, h)
			if b == BarRepeatStart {
				dots(x + thick + sep + thin + 0.4)
			}
		}
		return
	}
	switch b {
	case BarNone:
	case BarDouble:
		c.rect(x-thin, top, thin, h)
		c.rect(x-2*thin-sep, top, thin, h)
	case BarFinal:
		c.rect(x-thick, top, thick, h)
		c.rect(x-thick-sep-thin, top, thin, h)
	case BarHeavyLight:
		c.rect(x-thin, top, thin, h)
		c.rect(x-thin-sep-thick, top, thick, h)
	case BarRepeatEnd, BarRepeatBoth:
		c.rect(x-thick, top, thick, h)
		c.rect(x-thick-sep-thin, top, thin, h)
		dots(x - thick - sep - thin - 0.4 - dotW)
	case BarDashed:
		for _, y := range staves {
			c.line(x-thin/2, y, x-thin/2, y+4, thin, 0.5)
		}
	default:
		c.rect(x-thin, top, thin, h)
	}
}
