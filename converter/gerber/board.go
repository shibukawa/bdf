package gerber

import (
	"fmt"
	"math"
	"path"
	"strings"

	"github.com/shibukawa/bdf"
)

// magnification returns how many times larger than the board its pages
// are: the largest whole number that keeps the long side within that of
// A3 (420 mm); larger boards are drawn smaller, to fit it, as DXF drawings
// are.
func magnification(long float64) float64 {
	if long <= 0 {
		return 1
	}
	if long > 420 {
		return 420 / long
	}
	return math.Min(math.Floor(420/long), 10000)
}

// build makes the document of the layers read.
func (c *conv) build() (*Result, error) {
	pal, err := c.palette()
	if err != nil {
		return nil, err
	}
	doc := bdf.NewDocument()
	doc.Meta.Source = "gerber"
	c.meta(doc)

	// the extent of the board: what its layers draw (drawings beside it,
	// such as drill maps, make their own pages larger)
	var extent box
	var outlines []*layer
	board := false
	for _, l := range c.layers {
		switch l.fn.kind {
		case kOther:
			continue
		case kOutline:
			outlines = append(outlines, l)
		case kCopper, kMask, kSilk:
			board = true
		}
		extent = extent.union(l.img.bounds)
	}
	if !extent.ok {
		for _, l := range c.layers {
			extent = extent.union(l.img.bounds)
		}
	}
	if !extent.ok {
		return nil, fmt.Errorf("gerber: the files draw nothing")
	}
	loops, closed := boardShape(outlines)
	if len(outlines) > 0 && !closed {
		c.warnf("the outline does not close; the board is drawn as its bounding box")
	}
	if !closed {
		loops = []contour{rect(extent.Min.add(extent.Max).mul(0.5), extent.w(), extent.h())}
	}
	long := math.Max(extent.w(), extent.h())
	margin := math.Max(1, 0.03*long)
	mag := magnification(long + 2*margin)
	m := mapping{k: mag * 72 / 25.4, origin: vec{extent.Min.X - margin, extent.Max.Y + margin}}
	pw, ph := f32((extent.w()+2*margin)*m.k), f32((extent.h()+2*margin)*m.k)
	enc := &encoder{doc: doc, m: m}
	objs := make([]layerObject, len(c.layers))
	for i, l := range c.layers {
		objs[i] = enc.object(l.img)
	}
	bp := &bdf.Path{}
	w := &pathWriter{m: m, p: bp}
	for i := range loops {
		w.contour(&loops[i], true)
	}
	var lb box
	for i := range loops {
		lb = lb.union(loops[i].bounds())
	}
	br := m.rect(lb)
	// groups cover the board (and a little more for the lines along its edge)
	br = bdf.Rect{X: br.X - 2, Y: br.Y - 2, W: br.W + 4, H: br.H + 4}
	b := &builder{c: c, doc: doc, m: m, w: pw, h: ph, objs: objs, boardPath: bp, boardRect: br, clip: closed, pal: pal}

	res := &Result{Doc: doc, Layers: len(c.layers), W: extent.w(), H: extent.h(), Scale: mag}
	if board && len(c.layers) > 1 && c.opts.Views != ViewsLayers {
		b.boardView("top", "Top", sTop)
		b.boardView("bottom", "Bottom", sBottom)
		res.Board = true
	}
	if c.opts.Views != ViewsBoard || !res.Board {
		b.layerViews()
	}
	return res, nil
}

// meta sets the document's metadata: the project name the files or the job
// file give as the title, and the date they were made.
func (c *conv) meta(doc *bdf.Document) {
	title, date := "", ""
	if c.job != nil {
		title, date = c.job.GeneralSpecs.ProjectId.Name, c.job.Header.CreationDate
	}
	for _, l := range c.layers {
		if title == "" {
			if v, _, _ := strings.Cut(l.attrs["ProjectId"], ","); v != "" {
				title = v
			}
		}
		if date == "" {
			date = l.attrs["CreationDate"]
		}
	}
	if c.opts.Title != "" {
		title = c.opts.Title
	}
	if title != "" {
		doc.Meta.DC.Title = bdf.DCValues{title}
	}
	if date != "" {
		doc.Meta.DC.Created = bdf.DCValues{date}
	}
}

// builder makes the pages.
type builder struct {
	c         *conv
	doc       *bdf.Document
	m         mapping
	w, h      float32
	objs      []layerObject
	boardPath *bdf.Path
	boardRect bdf.Rect
	// clip: the board path is the outline, which clips the board views
	clip bool
	pal  palette
}

// page is a page object being drawn: the object, the rectangle groups
// cover (in the coordinates drawn in), the board's path and the layers'
// objects it refers to.
type page struct {
	o     *bdf.Object
	r     bdf.Rect
	board bdf.PathRef
	added bool
	path  *bdf.Path
	refs  map[int]bdf.ObjRef
}

func newPage(o *bdf.Object, r bdf.Rect, board *bdf.Path) *page {
	return &page{o: o, r: r, path: board, refs: map[int]bdf.ObjRef{}}
}

// ref returns the reference of layer i's object in the page object.
func (b *builder) ref(pg *page, i int) bdf.ObjRef {
	ref, ok := pg.refs[i]
	if !ok {
		ref = pg.o.AddObject(b.objs[i].hash, b.objs[i].bbox)
		pg.refs[i] = ref
	}
	return ref
}

// boardRef returns the board's path in the page object.
func (pg *page) boardRef() bdf.PathRef {
	if !pg.added {
		pg.board, pg.added = pg.o.AddPath(pg.path), true
	}
	return pg.board
}

func (pg *page) group(blend byte) {
	pg.o.GroupBegin(1, blend, pg.r.X, pg.r.Y, pg.r.W, pg.r.H)
}

// use draws layer i with the fill and stroke colors set: in a group of its
// own when it takes anything away, and as the board without it when it is
// negative.
func (b *builder) use(pg *page, i int) {
	o := pg.o
	if b.c.layers[i].negative {
		pg.group(bdf.BlendSourceOver)
		o.FillPath(pg.boardRef(), bdf.EvenOdd)
		b.subtract(pg, i)
		o.GroupEnd()
		return
	}
	ref := b.ref(pg, i)
	if b.objs[i].clear {
		pg.group(bdf.BlendSourceOver)
		o.Use(ref)
		o.GroupEnd()
		return
	}
	o.Use(ref)
}

// subtract takes layer i away from what is drawn.
func (b *builder) subtract(pg *page, i int) {
	o := pg.o
	ref := b.ref(pg, i)
	if b.objs[i].clear {
		pg.group(bdf.BlendDestinationOut)
		o.Use(ref)
		o.GroupEnd()
		return
	}
	o.Save()
	o.Blend(bdf.BlendDestinationOut)
	o.FillColor(bdf.RGB(0, 0, 0))
	o.StrokeColor(bdf.RGB(0, 0, 0))
	o.Use(ref)
	o.Restore()
}

func (b *builder) color(o *bdf.Object, c bdf.Color) {
	o.FillColor(c)
	o.StrokeColor(c)
}

// boardView adds a view of one side of the board as it looks: the
// substrate and the copper seen through the solder mask, what the openings
// of the mask show (the bare substrate and the finish of the copper), the
// silkscreen, and the holes. The bottom is seen from below, mirrored.
func (b *builder) boardView(id, title string, s side) {
	var cu, mask, silk, drill []int
	for i, l := range b.c.layers {
		f := l.fn
		onSide := f.side == s || s == sTop && f.side == sNone
		switch {
		case f.kind == kCopper && onSide:
			cu = append(cu, i)
		case f.kind == kMask && onSide:
			mask = append(mask, i)
		case f.kind == kSilk && onSide:
			silk = append(silk, i)
		case f.kind == kDrill:
			drill = append(drill, i)
		}
	}
	o := bdf.NewObject()
	o.SetBBox(0, 0, b.w, b.h)
	alt := "The top of the board"
	if s == sBottom {
		alt = "The bottom of the board, seen from below"
	}
	o.Mark(bdf.MarkFigure, alt)
	o.FillColor(background)
	o.FillRect(0, 0, b.w, b.h)
	if s == sBottom {
		o.Transform(-1, 0, 0, 1, b.w, 0)
	}
	pg := newPage(o, b.boardRect, b.boardPath)
	board := pg.boardRef()
	if b.clip {
		o.Save()
		o.ClipPath(board, bdf.EvenOdd)
	}
	if len(mask) == 0 {
		o.FillColor(b.pal.substrate)
		o.FillPath(board, bdf.EvenOdd)
		b.color(o, b.pal.finish)
		for _, i := range cu {
			b.use(pg, i)
		}
	} else {
		// the mask covers the substrate and the copper, which have one color
		// each: the colors they take under it are mixed here rather than
		// by drawing the mask translucent
		under := func(c bdf.Color) bdf.Color { return mix(b.pal.mask, c, b.pal.maskAlpha) }
		o.FillColor(under(b.pal.substrate))
		o.FillPath(board, bdf.EvenOdd)
		b.color(o, under(b.pal.copper))
		for _, i := range cu {
			b.use(pg, i)
		}
		// what the openings of the mask show: the bare substrate and the
		// finish of the copper, the openings' image the group's soft mask
		pg.group(bdf.BlendSourceOver)
		o.FillColor(b.pal.substrate)
		o.FillPath(board, bdf.EvenOdd)
		b.color(o, b.pal.finish)
		for _, i := range cu {
			b.use(pg, i)
		}
		o.MaskBegin(bdf.MaskAlpha, 0, nil)
		for _, i := range mask {
			o.Use(b.ref(pg, i))
		}
		o.MaskEnd()
		o.GroupEnd()
	}
	if len(silk) > 0 {
		// silkscreen is not printed on the openings of the mask
		pg.group(bdf.BlendSourceOver)
		b.color(o, b.pal.silk)
		for _, i := range silk {
			b.use(pg, i)
		}
		for _, i := range mask {
			b.subtract(pg, i)
		}
		o.GroupEnd()
	}
	b.color(o, background)
	for _, i := range drill {
		b.use(pg, i)
	}
	if b.clip {
		o.Restore()
	}
	o.Mark(bdf.MarkEnd, "")
	h, _ := b.doc.AddObject(o)
	v := b.doc.NewView(id, bdf.ViewFixed, title)
	v.AddPage(b.w, b.h, bdf.Layer{Role: bdf.RoleBody, Obj: h})
}

// layerViews adds a view for each layer, in the layer's color over the
// outline.
func (b *builder) layerViews() {
	titles := map[string]int{}
	for _, l := range b.c.layers {
		titles[b.title(l)]++
	}
	ids := map[string]bool{}
	outlines := []int{}
	for i, l := range b.c.layers {
		if l.fn.kind == kOutline {
			outlines = append(outlines, i)
		}
	}
	for i, l := range b.c.layers {
		title := b.title(l)
		base := path.Base(strings.ReplaceAll(l.name, "\\", "/"))
		if titles[title] > 1 && base != "" && base != "." && base != title {
			title += " (" + base + ")"
		}
		id := l.fn.id()
		if id == "" {
			id = "layer"
		}
		for n, base := 2, id; ids[id]; n++ {
			id = fmt.Sprintf("%s-%d", base, n)
		}
		ids[id] = true

		// the page holds the board's extent and whatever the layer draws
		// beside it (within reason: three times the board's page)
		x0, y0, x1, y1 := float32(0), float32(0), b.w, b.h
		if bb := l.img.bounds; bb.ok {
			r := b.m.rect(bb.grow(math.Max(1, 0.03*math.Max(bb.w(), bb.h()))))
			ux0, uy0 := min(x0, r.X), min(y0, r.Y)
			ux1, uy1 := max(x1, r.X+r.W), max(y1, r.Y+r.H)
			if max(ux1-ux0, uy1-uy0) <= 3*max(b.w, b.h) {
				x0, y0, x1, y1 = ux0, uy0, ux1, uy1
			} else if r.X < x0 || r.Y < y0 || r.X+r.W > x1 || r.Y+r.H > y1 {
				b.c.warnf("%s draws far outside the board; its page shows the board's extent only", b.title(l))
			}
		}
		w, h := x1-x0, y1-y0
		o := bdf.NewObject()
		o.SetBBox(0, 0, w, h)
		o.Mark(bdf.MarkFigure, title)
		o.FillColor(background)
		o.FillRect(0, 0, w, h)
		if x0 != 0 || y0 != 0 {
			o.Translate(-x0, -y0)
		}
		pg := newPage(o, bdf.Rect{X: x0, Y: y0, W: w, H: h}, b.boardPath)
		if l.fn.kind != kOutline {
			b.color(o, outlineDim)
			for _, j := range outlines {
				b.use(pg, j)
			}
		}
		b.color(o, layerColor(l.fn))
		b.use(pg, i)
		o.Mark(bdf.MarkEnd, "")
		hash, _ := b.doc.AddObject(o)
		v := b.doc.NewView(id, bdf.ViewFixed, title)
		v.AddPage(w, h, bdf.Layer{Role: bdf.RoleBody, Obj: hash})
	}
}

// title names a layer's view: what the layer is, or its file's name.
func (b *builder) title(l *layer) string {
	if t := l.fn.title(); t != "" {
		return t
	}
	if l.fn.other != "" {
		return otherTitle(l.fn.other)
	}
	if base := path.Base(strings.ReplaceAll(l.name, "\\", "/")); l.name != "" && base != "." {
		return base
	}
	return "Layer"
}

// otherNames are the names of the file functions this converter does not
// draw in the board views.
var otherNames = map[string]string{
	"drillmap": "Drill map", "vcut": "V-cut", "vcutmap": "V-cut map", "depthrout": "Depth route",
	"viafill": "Via fill", "carbonmask": "Carbon mask", "goldmask": "Gold mask", "heatsinkmask": "Heatsink mask",
	"peelablemask": "Peelable mask", "silvermask": "Silver mask", "tinmask": "Tin mask", "component": "Components",
}

// otherTitle names a file function that the board views do not draw
// ("AssemblyDrawing,Top" is Top assembly drawing); Other and OtherDrawing
// take the text the file gives.
func otherTitle(v string) string {
	fs := strings.Split(v, ",")
	name := strings.TrimSpace(fs[0])
	if (name == "Other" || name == "OtherDrawing") && len(fs) > 1 && strings.TrimSpace(fs[1]) != "" {
		return strings.TrimSpace(fs[1])
	}
	t, ok := otherNames[strings.ToLower(name)]
	if !ok {
		// AssemblyDrawing: Assembly drawing
		var b strings.Builder
		for i, r := range name {
			if i > 0 && r >= 'A' && r <= 'Z' {
				b.WriteByte(' ')
				r += 'a' - 'A'
			}
			b.WriteRune(r)
		}
		t = b.String()
	}
	for _, f := range fs[1:] {
		switch strings.TrimSpace(f) {
		case "Top":
			return "Top " + strings.ToLower(t[:1]) + t[1:]
		case "Bot":
			return "Bottom " + strings.ToLower(t[:1]) + t[1:]
		}
	}
	return t
}
