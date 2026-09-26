package jww

import (
	"fmt"
	"strings"
)

// base is what every figure has (CData).
type base struct {
	group  uint32
	style  int // line type
	color  int // line color
	width  int // line width (0: the color's)
	layer  int
	glayer int
	flag   int
}

type point struct{ x, y float64 }

type sen struct {
	base
	a, b point
}

type enko struct {
	base
	c                   point
	r, start, arc, tilt float64 // radians
	flat                float64
	full                bool
}

type ten struct {
	base
	p      point
	kari   bool // a temporary point, not printed
	code   int
	rot    float64
	scale  float64
	hasSub bool
}

type moji struct {
	base
	a, b          point
	kind          int // 文字種 (+10000 italic, +20000 bold)
	w, h, spacing float64
	angle         float64 // degrees
	font, s       string
}

type sunpou struct {
	base
	line   sen
	text   moji
	sxf    bool
	aux    [2]*sen
	points [4]*ten
	hasAux bool
}

type solid struct {
	base
	p1, p4, p2, p3 point
	rgb            uint32
}

type block struct {
	base
	p      point
	sx, sy float64
	rot    float64 // radians
	number int
}

type blockDef struct {
	base
	number int
	name   string
	list   []any
}

// data is the figures of a file and its block definitions.
type data struct {
	figures []any
	blocks  map[int]*blockDef
	images  map[string][]byte
}

// readData reads the list of figures and the list of block definitions
// (and the images of version 7). The size of CTime in the definitions
// depends on the version of MFC that wrote the file: timeSize is 4 or 8.
func readData(a *archive, h *header, timeSize int) (*data, error) {
	d := &data{blocks: map[int]*blockDef{}}
	r := &dataReader{a: a, h: h, timeSize: timeSize}
	d.figures = r.list()
	if a.err != nil {
		return d, a.err
	}
	for _, o := range r.list() {
		if b, ok := o.(*blockDef); ok {
			d.blocks[b.number] = b
		}
	}
	if a.err != nil {
		return d, a.err
	}
	if h.version >= 700 && a.i < len(a.b) {
		n := int(a.u32())
		for range n {
			name := a.str()
			size := int(a.u32())
			if a.err != nil || !a.need(size) {
				break
			}
			if d.images == nil {
				d.images = map[string][]byte{}
			}
			d.images[name] = a.b[a.i : a.i+size]
			a.i += size
		}
		// images are optional; a short tail is not an error
		a.err = nil
	}
	return d, nil
}

type dataReader struct {
	a        *archive
	h        *header
	timeSize int
	depth    int
}

func (r *dataReader) list() []any {
	n := r.a.count()
	out := make([]any, 0, min(n, 1<<16))
	for range n {
		if r.a.err != nil {
			break
		}
		if o := r.object(); o != nil {
			out = append(out, o)
		}
	}
	return out
}

func (r *dataReader) object() any {
	class, idx, old := r.a.objectTag()
	if class == "" {
		return old
	}
	var o any
	switch class {
	case "CDataSen":
		o = r.sen()
	case "CDataEnko":
		o = r.enko()
	case "CDataTen":
		o = r.ten()
	case "CDataMoji":
		o = r.moji()
	case "CDataSunpou":
		o = r.sunpou()
	case "CDataSolid":
		o = r.solid()
	case "CDataBlock":
		o = r.block()
	case "CDataList":
		o = r.blockDef(idx)
	default:
		r.a.err = fmt.Errorf("figures of class %s cannot be read", class)
		return nil
	}
	r.a.objects[idx] = o
	return o
}

func (r *dataReader) base() base {
	a := r.a
	b := base{group: a.u32(), style: int(a.u8()), color: int(a.u16())}
	if r.h.version >= 351 {
		b.width = int(a.u16())
	}
	b.layer = int(a.u16())
	b.glayer = int(a.u16())
	b.flag = int(a.u16())
	return b
}

func (r *dataReader) pt() point { return point{r.a.f64(), r.a.f64()} }

func (r *dataReader) sen() *sen {
	s := &sen{base: r.base()}
	s.a, s.b = r.pt(), r.pt()
	return s
}

func (r *dataReader) enko() *enko {
	e := &enko{base: r.base()}
	e.c = r.pt()
	e.r, e.start, e.arc, e.tilt, e.flat = r.a.f64(), r.a.f64(), r.a.f64(), r.a.f64(), r.a.f64()
	e.full = r.a.u32() != 0
	return e
}

func (r *dataReader) ten() *ten {
	t := &ten{base: r.base()}
	t.p = r.pt()
	t.kari = r.a.u32() != 0
	if t.style == 100 {
		t.hasSub = true
		t.code = int(r.a.u32())
		t.rot = r.a.f64()
		t.scale = r.a.f64()
	}
	return t
}

func (r *dataReader) moji() *moji {
	m := &moji{base: r.base()}
	m.a, m.b = r.pt(), r.pt()
	m.kind = int(r.a.u32())
	m.w, m.h, m.spacing = r.a.f64(), r.a.f64(), r.a.f64()
	m.angle = r.a.f64()
	m.font = r.a.str()
	m.s = r.a.str()
	return m
}

func (r *dataReader) sunpou() *sunpou {
	s := &sunpou{base: r.base()}
	s.line = *r.sen()
	s.text = *r.moji()
	if r.h.version >= 420 {
		s.hasAux = true
		s.sxf = r.a.u16() != 0
		s.aux[0], s.aux[1] = r.sen(), r.sen()
		for i := range s.points {
			s.points[i] = r.ten()
		}
	}
	return s
}

func (r *dataReader) solid() *solid {
	s := &solid{base: r.base()}
	s.p1, s.p4, s.p2, s.p3 = r.pt(), r.pt(), r.pt(), r.pt()
	if s.color == 10 {
		s.rgb = r.a.u32()
	}
	return s
}

func (r *dataReader) block() *block {
	b := &block{base: r.base()}
	b.p = r.pt()
	b.sx, b.sy, b.rot = r.a.f64(), r.a.f64(), r.a.f64()
	b.number = int(r.a.u32())
	return b
}

func (r *dataReader) blockDef(idx uint32) *blockDef {
	b := &blockDef{base: r.base()}
	// register the definition before its content, which may refer to it
	r.a.objects[idx] = b
	b.number = int(r.a.u32())
	r.a.u32() // referred to
	r.a.skip(r.timeSize)
	b.name = r.a.str()
	if i := strings.Index(b.name, "@@SfigorgFlag@@"); i >= 0 {
		b.name = b.name[:i]
	}
	r.depth++
	if r.depth > 32 {
		r.a.err = fmt.Errorf("block definitions nested deeper than 32 levels")
		return b
	}
	b.list = r.list()
	r.depth--
	return b
}
