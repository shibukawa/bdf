package dxf

import (
	"math"
	"sort"
	"strings"

	"github.com/shibukawa/bdf/converter/internal/cad"
)

// entity is an entity, a table entry or an object: its type and the tags
// that follow it.
type entity struct {
	typ  string
	tags []tag
	// kids are the VERTEX entities of a POLYLINE and the ATTRIB entities of
	// an INSERT.
	kids []*entity
}

// get returns the first tag with a code, before the extended data.
func (e *entity) get(code int) (tag, bool) {
	for _, t := range e.tags {
		if t.code >= 1000 {
			break
		}
		if t.code == code {
			return t, true
		}
	}
	return tag{}, false
}

func (e *entity) has(code int) bool {
	_, ok := e.get(code)
	return ok
}

func (e *entity) str(code int) string {
	t, _ := e.get(code)
	return t.s
}

func (e *entity) num(code int, def float64) float64 {
	if t, ok := e.get(code); ok {
		return t.f
	}
	return def
}

func (e *entity) int(code int, def int) int {
	if t, ok := e.get(code); ok {
		return int(t.f)
	}
	return def
}

// pt returns the point whose x has the code (y has code+10).
func (e *entity) pt(code int) cad.Point {
	return cad.Point{X: e.num(code, 0), Y: e.num(code+10, 0)}
}

// vec3 returns the 3D point whose x has the code.
func (e *entity) vec3(code int, def [3]float64) [3]float64 {
	if !e.has(code) {
		return def
	}
	return [3]float64{e.num(code, 0), e.num(code+10, 0), e.num(code+20, 0)}
}

// sub returns the tags of a subclass (after its 100 marker), as an entity
// of the same type.
func (e *entity) sub(name string) *entity {
	out := &entity{typ: e.typ}
	in := false
	for _, t := range e.tags {
		if t.code == 100 {
			if in {
				break
			}
			in = strings.EqualFold(t.s, name)
			continue
		}
		if t.code >= 1000 {
			break
		}
		if in {
			out.tags = append(out.tags, t)
		}
	}
	return out
}

// xdata returns the extended data of an application.
func (e *entity) xdata(app string) []tag {
	var out []tag
	in := false
	for _, t := range e.tags {
		if t.code == 1001 {
			in = strings.EqualFold(t.s, app)
			continue
		}
		if in {
			out = append(out, t)
		}
	}
	return out
}

func (e *entity) handle() string {
	if e.typ == "DIMSTYLE" {
		if h := e.str(105); h != "" {
			return strings.ToUpper(h)
		}
	}
	return strings.ToUpper(e.str(5))
}

type layer struct {
	name       string
	color      int   // ACI; negative: the layer is off
	trueColor  int32 // -1 when not set
	ltype      string
	lineweight int
	frozen     bool
	noPlot     bool
	alpha      float64 // 1 = opaque
}

type ltype struct {
	name  string
	elems []float64 // positive: dash, negative: gap, 0: dot
}

type style struct {
	name            string
	font, bigFont   string
	family          string // TrueType family from the extended data
	bold, italic    bool
	height, width   float64
	oblique         float64
	backward, upset bool
	vertical        bool
}

type block struct {
	name     string
	base     cad.Point
	flags    int
	xref     string
	entities []*entity
}

type blockRecord struct {
	name   string
	layout string // handle of the layout of a paper space block
}

type layout struct {
	name     string
	order    int
	block    string // handle of the block record
	handle   string
	plot     *entity // AcDbPlotSettings
	lay      *entity // AcDbLayout
	blockKey string  // the block holding its entities (upper case)
}

// document is the content of a DXF file.
type document struct {
	version    string
	encoding   string
	header     map[string][]tag
	layers     map[string]*layer
	ltypes     map[string]*ltype
	styles     map[string]*style
	dimstyles  map[string]*entity
	records    map[string]*blockRecord // by handle
	blocks     map[string]*block       // by upper-case name
	entities   []*entity
	layouts    []*layout
	objects    map[string]*entity // by handle
	byHandle   map[string]*entity // entities by handle
	dictionary map[string]string  // named object dictionary: name → handle
	// table entries by handle
	layerHandles map[string]*layer
	styleHandles map[string]*style
}

func key(name string) string { return strings.ToUpper(strings.TrimSpace(name)) }

func (d *document) hdr(name string) (tag, bool) {
	ts := d.header[name]
	if len(ts) == 0 {
		return tag{}, false
	}
	return ts[0], true
}

func (d *document) hnum(name string, def float64) float64 {
	if t, ok := d.hdr(name); ok {
		return t.f
	}
	return def
}

func (d *document) hpt(name string) (cad.Point, bool) {
	ts := d.header[name]
	var p cad.Point
	var okx, oky bool
	for _, t := range ts {
		switch t.code {
		case 10:
			p.X, okx = t.f, true
		case 20:
			p.Y, oky = t.f, true
		}
	}
	return p, okx && oky && !math.IsInf(p.X, 0) && math.Abs(p.X) < 1e19 && math.Abs(p.Y) < 1e19
}

// split cuts tags into entities at each group code 0.
func split(tags []tag) []*entity {
	var out []*entity
	var cur *entity
	for _, t := range tags {
		if t.code == 0 {
			cur = &entity{typ: strings.ToUpper(strings.TrimSpace(t.s))}
			out = append(out, cur)
			continue
		}
		if cur != nil {
			cur.tags = append(cur.tags, t)
		}
	}
	return out
}

// group attaches the VERTEX entities to their POLYLINE and the ATTRIB
// entities to their INSERT, dropping the SEQEND markers.
func group(list []*entity) []*entity {
	var out []*entity
	for i := 0; i < len(list); i++ {
		e := list[i]
		kid := ""
		switch e.typ {
		case "POLYLINE":
			kid = "VERTEX"
		case "INSERT":
			kid = "ATTRIB"
		}
		out = append(out, e)
		if kid == "" {
			continue
		}
		for i+1 < len(list) && list[i+1].typ == kid {
			e.kids = append(e.kids, list[i+1])
			i++
		}
		if i+1 < len(list) && list[i+1].typ == "SEQEND" {
			i++
		}
	}
	return out
}

// parse reads the sections of a DXF file.
func parse(tags []tag) *document {
	d := &document{header: map[string][]tag{}, layers: map[string]*layer{}, ltypes: map[string]*ltype{},
		styles: map[string]*style{}, dimstyles: map[string]*entity{}, records: map[string]*blockRecord{},
		blocks: map[string]*block{}, objects: map[string]*entity{}, byHandle: map[string]*entity{},
		dictionary: map[string]string{}, layerHandles: map[string]*layer{}, styleHandles: map[string]*style{}}
	d.encoding = decodeStrings(tags)
	for i := 0; i < len(tags); i++ {
		t := tags[i]
		if t.code != 0 || t.s != "SECTION" {
			continue
		}
		name := ""
		if i+1 < len(tags) && tags[i+1].code == 2 {
			name = strings.ToUpper(tags[i+1].s)
			i++
		}
		j := i + 1
		for j < len(tags) && !(tags[j].code == 0 && (tags[j].s == "ENDSEC" || tags[j].s == "EOF")) {
			j++
		}
		body := tags[i+1 : j]
		i = j
		switch name {
		case "HEADER":
			d.readHeader(body)
		case "TABLES":
			d.readTables(split(body))
		case "BLOCKS":
			d.readBlocks(split(body))
		case "ENTITIES":
			d.entities = group(split(body))
			for _, e := range d.entities {
				if h := e.handle(); h != "" {
					d.byHandle[h] = e
				}
			}
		case "OBJECTS":
			d.readObjects(split(body))
		}
	}
	if t, ok := d.hdr("$ACADVER"); ok {
		d.version = strings.ToUpper(strings.TrimSpace(t.s))
	}
	d.resolveLayouts()
	return d
}

func (d *document) readHeader(tags []tag) {
	var name string
	for _, t := range tags {
		if t.code == 9 {
			name = strings.ToUpper(t.s)
			continue
		}
		if name != "" {
			d.header[name] = append(d.header[name], t)
		}
	}
}

func (d *document) readTables(list []*entity) {
	for _, e := range list {
		switch e.typ {
		case "LAYER":
			l := &layer{name: e.str(2), color: e.int(62, 7), trueColor: -1, ltype: e.str(6),
				lineweight: e.int(370, -3), frozen: e.int(70, 0)&1 != 0, noPlot: e.has(290) && e.int(290, 1) == 0, alpha: 1}
			if e.has(420) {
				l.trueColor = int32(e.int(420, 0))
			}
			for _, x := range e.xdata("AcCmTransparency") {
				if x.code == 1071 {
					l.alpha = transparency(int(x.f))
				}
			}
			if strings.EqualFold(l.name, "Defpoints") {
				l.noPlot = true
			}
			d.layers[key(l.name)] = l
			d.layerHandles[e.handle()] = l
		case "LTYPE":
			lt := &ltype{name: e.str(2)}
			for _, t := range e.tags {
				if t.code == 49 {
					lt.elems = append(lt.elems, t.f)
				}
			}
			d.ltypes[key(lt.name)] = lt
		case "STYLE":
			s := &style{name: e.str(2), font: e.str(3), bigFont: e.str(4), height: e.num(40, 0), width: e.num(41, 1),
				oblique: e.num(50, 0)}
			flags := e.int(70, 0)
			gen := e.int(71, 0)
			s.backward, s.upset, s.vertical = gen&2 != 0, gen&4 != 0, flags&4 != 0
			for _, x := range e.xdata("ACAD") {
				switch x.code {
				case 1000:
					if s.family == "" {
						s.family = x.s
					}
				case 1071:
					v := int(x.f)
					s.italic, s.bold = v&0x1000000 != 0, v&0x2000000 != 0
				}
			}
			d.styles[key(s.name)] = s
			d.styleHandles[e.handle()] = s
		case "DIMSTYLE":
			d.dimstyles[key(e.str(2))] = e
		case "BLOCK_RECORD":
			d.records[e.handle()] = &blockRecord{name: e.str(2), layout: strings.ToUpper(e.str(340))}
		}
	}
}

func (d *document) readBlocks(list []*entity) {
	var cur *block
	var ents []*entity
	for _, e := range list {
		switch e.typ {
		case "BLOCK":
			cur = &block{name: e.str(2), base: e.pt(10), flags: e.int(70, 0), xref: e.str(1)}
			if cur.name == "" {
				cur.name = e.str(3)
			}
			ents = nil
		case "ENDBLK":
			if cur != nil {
				cur.entities = group(ents)
				for _, x := range cur.entities {
					if h := x.handle(); h != "" {
						d.byHandle[h] = x
					}
				}
				d.blocks[key(cur.name)] = cur
			}
			cur = nil
		default:
			if cur != nil {
				ents = append(ents, e)
			}
		}
	}
}

func (d *document) readObjects(list []*entity) {
	first := true
	for _, e := range list {
		if h := e.handle(); h != "" {
			d.objects[h] = e
		}
		if first && e.typ == "DICTIONARY" {
			// the named object dictionary
			first = false
			var name string
			for _, t := range e.tags {
				switch t.code {
				case 3:
					name = strings.ToUpper(t.s)
				case 350, 360:
					if name != "" {
						d.dictionary[name] = strings.ToUpper(t.s)
					}
				}
			}
		}
		if e.typ == "LAYOUT" {
			lay := e.sub("AcDbLayout")
			l := &layout{name: lay.str(1), order: lay.int(71, 0), block: strings.ToUpper(lay.str(330)), handle: e.handle(),
				plot: e.sub("AcDbPlotSettings"), lay: lay}
			d.layouts = append(d.layouts, l)
		}
	}
}

// resolveLayouts finds the block of each layout and orders them by tab.
func (d *document) resolveLayouts() {
	for _, l := range d.layouts {
		if r := d.records[l.block]; r != nil {
			l.blockKey = key(r.name)
		}
		if l.blockKey == "" {
			if strings.EqualFold(l.name, "Model") {
				l.blockKey = "*MODEL_SPACE"
			}
		}
	}
	sort.SliceStable(d.layouts, func(i, j int) bool { return d.layouts[i].order < d.layouts[j].order })
}

// transparency converts a DXF transparency value (0x020000TT) into an
// opacity.
func transparency(v int) float64 {
	if v&0x02000000 == 0 {
		return 1
	}
	return float64(v&0xff) / 255
}
