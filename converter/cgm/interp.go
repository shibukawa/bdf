package cgm

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"golang.org/x/text/encoding"
)

func itoa(v int) string { return strconv.Itoa(v) }

// Specification modes of line widths, marker sizes, edge widths and the
// sizes of interiors (ISO/IEC 8632-1 §7.3).
const (
	modeAbsolute   = 0 // VDC
	modeScaled     = 1 // a factor of the nominal size
	modeFractional = 2 // a fraction of the VDC extent
	modeMM         = 3 // millimetres on paper
)

// Interior styles.
const (
	styleHollow = iota
	styleSolid
	stylePattern
	styleHatch
	styleEmpty
	styleGeometric
	styleInterpolated
)

// Nominal sizes on paper (mm) of scaled line widths and marker sizes, and
// the device-dependent sizes of dashes, hatching and patterns.
const (
	nominalLine    = 0.25
	nominalMarker  = 2.5
	hatchSpacing   = 2.0
	defaultPattern = 4.0
	mm             = 72 / 25.4 // pt
)

// colour is a colour as an element gives it: an index into the colour
// table (looked up when it is drawn), or a direct colour.
type colour struct {
	index  int
	direct bool
	c      bdf.Color
}

// size is a width or a size in the specification mode it was given in.
type size struct {
	v    float64
	mode int
}

type lineRep struct {
	typ    int
	width  size
	colour colour
}

type textRep struct {
	font, prec         int
	spacing, expansion float64
	colour             colour
}

type fillRep struct {
	style          int
	colour         colour
	hatch, pattern int
}

// lineTypeDef is a line type of LINE AND EDGE TYPE DEFINITION: dashes and
// gaps in proportion, repeated every repeat.
type lineTypeDef struct {
	dash   []float64
	repeat size
}

// hatchDef is a hatch style of HATCH STYLE DEFINITION.
type hatchDef struct {
	cross      bool
	dir1, dir2 cad.Point
	cycle      size
	gaps       []float64
	types      []int
}

// patternDef is a pattern of PATTERN TABLE.
type patternDef struct {
	nx, ny int
	cells  []colour
}

// Aspect source flags (ISO/IEC 8632-1 §7.7.35): which attributes come from
// the bundles.
const (
	asfLineType = iota
	asfLineWidth
	asfLineColour
	asfMarkerType
	asfMarkerSize
	asfMarkerColour
	asfFont
	asfTextPrecision
	asfExpansion
	asfSpacing
	asfTextColour
	asfInteriorStyle
	asfFillColour
	asfHatchIndex
	asfPatternIndex
	asfEdgeType
	asfEdgeWidth
	asfEdgeColour
	asfCount
)

// attrs are the primitive attributes, which SAVE PRIMITIVE CONTEXT saves.
type attrs struct {
	lineBundle, lineType int
	lineWidth            size
	lineColour           colour
	lineCap, lineJoin    int
	lineTypeOffset       float64

	markerBundle, markerType int
	markerSize               size
	markerColour             colour

	textBundle, font     int
	expansion, spacing   float64
	textColour           colour
	charHeight           float64 // negative: the default
	up, base             cad.Point
	path, hAlign, vAlign int
	contH, contV         float64
	charSet, altCharSet  int
	restrictedType       int

	fillBundle, interior int
	fillColour           colour
	hatch, pattern       int
	fillRef              cad.Point
	fillRefSet           bool
	patternSize          [4]float64 // height vector, width vector
	patternSizeSet       bool
	interp               *interpolated

	edgeBundle, edgeType int
	edgeWidth            size
	edgeColour           colour
	edgeVisible          bool
	edgeCap, edgeJoin    int
	edgeTypeOffset       float64

	asf [asfCount]bool // bundled
}

// interpolated is an INTERPOLATED INTERIOR.
type interpolated struct {
	style   int
	geom    []float64
	stages  []float64
	colours []colour
}

// state is the state of a picture: its descriptor, the control elements
// and the attributes.
type state struct {
	scaling                      int // 0 abstract, 1 metric
	metric                       float64
	colourMode                   int // 0 indexed, 1 direct
	lineMode, markerMode         int
	edgeMode, interiorMode       int
	extent                       [2]cad.Point
	extentSet                    bool
	background                   bdf.Color
	viewportMode                 int
	viewportScale                float64
	viewport                     [2]cad.Point
	viewportSet                  bool
	lineReps, markerReps, edgeRs map[int]lineRep
	textReps                     map[int]textRep
	fillReps                     map[int]fillRep
	lineTypes                    map[int]lineTypeDef
	hatchStyles                  map[int]hatchDef
	patterns                     map[int]patternDef
	table                        map[int]bdf.Color
	aux                          colour
	transparency                 bool
	clip                         [2]cad.Point
	clipSet, clipOn              bool
	transparentCell              *colour
	a                            attrs
	saved                        map[int]attrs
}

// newState returns the defaults of a picture (ISO/IEC 8632-1 Annex B);
// device-dependent ones are those of a white sheet of paper.
func newState() state {
	return state{
		lineMode: modeScaled, markerMode: modeScaled, edgeMode: modeScaled, interiorMode: modeAbsolute,
		background: bdf.RGB(255, 255, 255), viewportScale: 1,
		lineReps: map[int]lineRep{}, markerReps: map[int]lineRep{}, edgeRs: map[int]lineRep{},
		textReps: map[int]textRep{}, fillReps: map[int]fillRep{}, lineTypes: map[int]lineTypeDef{},
		hatchStyles: map[int]hatchDef{}, patterns: map[int]patternDef{}, table: map[int]bdf.Color{},
		transparency: true, clipOn: true, saved: map[int]attrs{},
		aux: colour{index: 0},
		a: attrs{
			lineBundle: 1, lineType: 1, lineWidth: size{1, modeScaled}, lineColour: colour{index: 1}, lineCap: 1, lineJoin: 1,
			markerBundle: 1, markerType: 3, markerSize: size{1, modeScaled}, markerColour: colour{index: 1},
			textBundle: 1, font: 1, expansion: 1, textColour: colour{index: 1}, charHeight: -1,
			up: cad.Point{Y: 1}, base: cad.Point{X: 1}, charSet: 1, altCharSet: 1, restrictedType: 1,
			fillBundle: 1, interior: styleHollow, fillColour: colour{index: 1}, hatch: 1, pattern: 1,
			edgeBundle: 1, edgeType: 1, edgeWidth: size{1, modeScaled}, edgeColour: colour{index: 1}, edgeCap: 1, edgeJoin: 1,
		},
	}
}

// picture is a picture being drawn, which becomes a page.
type picture struct {
	number int // 1-based
	// the page in mm, and the transform from VDC to it (pt, y down)
	w, h float64
	m    canvas.Matrix
	k    float64 // pt per VDC unit
	bg   bdf.Color
	d    *cad.Drawing
	// long is the longer side of the VDC extent
	long float64
}

// figure is a closed figure being built (BEGIN FIGURE): its boundary and
// the edges drawn with it.
type figure struct {
	path *cad.Path
}

// segment is a segment being recorded or recorded.
type segment struct {
	id     int
	d      *cad.Drawing
	points int // in.points when it began
}

// interp interprets the elements of a metafile.
type interp struct {
	c  *converter
	pr precisions
	// fonts measures text; fallback reads the strings that the character
	// sets do not say are multibyte (nil: Latin-1)
	fonts    *cad.Fonts
	fallback encoding.Encoding
	// scan only collects the strings of text elements (strs)
	scan bool
	strs [][]byte
	// clear text
	text bool

	title        string
	colourExt    [2][4]float64
	colourExtSet bool
	fontNames    []string
	fontSpecs    []cad.FontSpec
	charsets     []charset
	colourModel  int
	defaults     []element

	st      state
	pic     *picture
	body    bool
	out     *cad.Drawing
	clipped bool // a clip group is open in the picture's drawing
	fig     *figure
	// compound is the path of a compound line (BEGIN COMPOUND LINE)
	compound *cad.Path
	pending  *pendingText
	segs     map[int]*cad.Drawing
	open     []*segment
	// hidden counts the enclosing application structures that are
	// invisible, skip whether the picture is left out
	aps    []bool
	hidden int
	skip   bool
	tiles  *tileArray

	count    int // pictures begun
	pictures []*picture
	selected func(n int) bool
	// points and cellsUsed count against maxPoints and maxTotalCells;
	// weights holds the points a recorded segment adds with each copy
	points, cellsUsed int
	weights           map[int]int
	// patternImages caches the images of the patterns, which change with
	// the pattern and colour tables (generation)
	patternImages map[[2]int]bdf.Hash
	generation    int
}

// Limits against malformed metafiles (variables so that tests can lower
// them): the points drawn, counting what hatching and copies of segments
// add (the DXF converter draws at most 5 million entities a view, §3.12);
// the cells of all cell arrays and tiles, and of one; the pictures; the
// elements of METAFILE DEFAULTS REPLACEMENT, which every picture replays.
var (
	maxPoints     = 10_000_000
	maxTotalCells = 64 << 20
	maxCells      = 16 << 20
	maxPictures   = 10_000
	maxDefaults   = 1024
)

// maxPattern bounds the cells of a pattern.
const maxPattern = 1 << 16

func (in *interp) warnOnce(key, format string, args ...any) { in.c.warnOnce(key, format, args...) }

// drawing reports whether primitives are drawn: in a picture that is
// converted and outside hidden application structures.
func (in *interp) drawing() bool { return in.pic != nil && !in.skip && in.hidden == 0 }

// run interprets the elements.
func (in *interp) run(next func() (element, bool)) {
	for {
		e, ok := next()
		if !ok {
			break
		}
		if e.code == eEndMetafile {
			break
		}
		in.element(&e)
		if in.points > maxPoints {
			in.warnOnce("points", "the metafile has more than %d points; the rest is left out", maxPoints)
			break
		}
		if in.count > maxPictures {
			in.warnOnce("pictures", "the metafile has more than %d pictures; the rest are left out", maxPictures)
			break
		}
	}
	if in.pic != nil {
		in.endPicture()
	}
}

// element interprets an element.
func (in *interp) element(e *element) {
	p := e.params(&in.pr)
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(shortError); ok {
				n := e.code.name()
				if e.text {
					n = e.word
				}
				in.warnOnce("short:"+n, "%s elements whose parameters could not be read are ignored", n)
				return
			}
			panic(r)
		}
	}()
	if in.scan {
		in.scanElement(e, p)
		return
	}
	switch e.class {
	case 0:
		in.delimiter(e, p)
	case 1:
		in.descriptor(e, p)
	case 2:
		in.pictureDescriptor(e.code, p)
	case 3:
		in.control(e.code, p)
	case 4:
		in.primitive(e.code, p)
	case 5:
		in.attribute(e.code, p)
	case 6, 7:
		// escapes, messages and application data draw nothing
	case 8:
		in.segmentElement(e.code, p)
	case 9:
		in.apsAttribute(p)
	default:
		if e.code == eUnknownKeyword {
			in.warnOnce("keyword:"+e.word, "unknown element %s is ignored", e.word)
		} else if e.code != eEndMFDefaults {
			in.warnOnce("class", "elements of unknown classes are ignored")
		}
	}
}

func (in *interp) delimiter(e *element, p params) {
	// a text that is not final ends at any structure
	in.flushText()
	switch e.code {
	case eBeginMetafile:
		if p.more() {
			in.title = strings.TrimSpace(in.decodeName(p.str()))
		}
	case eBeginPicture:
		if in.pic != nil {
			in.endPicture()
		}
		in.beginPicture()
	case eBeginPictureBody:
		if in.pic == nil {
			in.beginPicture()
		}
		in.beginBody()
	case eEndPicture:
		if in.pic != nil {
			in.endPicture()
		}
	case eBeginSegment:
		id := p.name()
		in.ensureBody()
		s := &segment{id: id, d: &cad.Drawing{}, points: in.points}
		in.open = append(in.open, s)
		in.out = s.d
	case eEndSegment:
		if n := len(in.open); n > 0 {
			s := in.open[n-1]
			in.open = in.open[:n-1]
			in.out = in.target()
			in.segs[s.id] = s.d
			in.weights[s.id] = in.points - s.points + 1
			if !s.d.Empty() {
				in.out.Append(s.d, canvas.Identity, nil)
			}
		}
	case eBeginFigure:
		in.ensureBody()
		in.fig = &figure{path: &cad.Path{}}
	case eEndFigure:
		if f := in.fig; f != nil {
			in.fig = nil
			f.path.Close()
			in.area(f.path, nil)
		}
	case eBeginCompoundLine:
		in.ensureBody()
		in.compound = &cad.Path{}
	case eEndCompoundLine:
		if c := in.compound; c != nil {
			in.compound = nil
			if in.drawing() {
				in.out.Stroke(c, in.linePen())
			}
		}
	case eBeginTileArray:
		in.ensureBody()
		in.beginTileArray(p)
	case eEndTileArray:
		in.endTileArray()
	case eBeginAPS:
		// its identifier, type and inheritance do not change the drawing
		in.aps = append(in.aps, false)
	case eEndAPS:
		if n := len(in.aps); n > 0 {
			if in.aps[n-1] {
				in.hidden--
			}
			in.aps = in.aps[:n-1]
		}
	case eBeginCompoundText, eEndCompoundText, eBeginProtection, eEndProtection, eBeginAPSBody:
	case eBeginMetafile, eEndMetafile:
	}
}

// apsAttribute reads an APPLICATION STRUCTURE ATTRIBUTE: a WebCGM
// application structure whose visibility is off is not drawn.
func (in *interp) apsAttribute(p params) {
	name := strings.ToLower(strings.TrimSpace(string(p.str())))
	if name != "visibility" || len(in.aps) == 0 {
		return
	}
	off := false
	if sdr := p.sdr(); len(sdr) > 0 && len(sdr[0].values) > 0 {
		switch v := sdr[0].values[0].(type) {
		case int:
			off = v == 0
		case float64:
			off = v == 0
		case string:
			off = v == "OFF"
		case []byte:
			off = strings.EqualFold(string(v), "off")
		}
	}
	n := len(in.aps)
	if off != in.aps[n-1] {
		in.aps[n-1] = off
		if off {
			in.hidden++
		} else {
			in.hidden--
		}
	}
}

func (in *interp) descriptor(e *element, p params) {
	switch e.code {
	case eVDCType:
		in.pr.vdcReal = p.enum("INTEGER REAL") == 1
	case eIntegerPrecision:
		if !in.text {
			in.pr.intBits = bitsOf(p.int(), in.pr.intBits)
		}
	case eRealPrecision:
		if !in.text {
			in.pr.real = in.realPrecision(p, in.pr.real)
		}
	case eIndexPrecision:
		if !in.text {
			in.pr.indexBits = bitsOf(p.int(), in.pr.indexBits)
		}
	case eColourPrecision:
		if in.text {
			if v := p.int(); v > 0 {
				in.pr.colourBits = v // the largest component value
			}
		} else {
			in.pr.colourBits = bitsOf(p.int(), in.pr.colourBits)
		}
	case eColourIndexPrec:
		if !in.text {
			in.pr.colourIndexBits = bitsOf(p.int(), in.pr.colourIndexBits)
		}
	case eNamePrecision:
		if !in.text {
			in.pr.nameBits = bitsOf(p.int(), in.pr.nameBits)
		}
	case eColourValueExtent:
		if in.colourModel == 1 || in.colourModel == 4 || in.colourModel == 0 {
			n := p.components()
			for i := 0; i < n; i++ {
				in.colourExt[0][i] = p.component()
			}
			for i := 0; i < n; i++ {
				in.colourExt[1][i] = p.component()
			}
			in.colourExtSet = true
		}
	case eDefaultsReplacement:
		add := func(sub element) {
			if len(in.defaults) >= maxDefaults {
				in.warnOnce("defaults", "metafile defaults replacements of more than %d elements are cut short", maxDefaults)
				return
			}
			in.defaults = append(in.defaults, sub)
		}
		if e.text {
			for _, sub := range e.sub {
				add(sub)
			}
		} else {
			r := &binReader{data: e.data}
			for {
				sub, ok := r.next()
				if !ok {
					break
				}
				add(sub)
			}
		}
	case eFontList:
		in.fontNames, in.fontSpecs = nil, nil
		for p.more() {
			n := in.decodeName(p.str())
			in.fontNames = append(in.fontNames, n)
			in.fontSpecs = append(in.fontSpecs, fontSpec(n))
		}
	case eCharacterSetList:
		in.charsets = in.charsets[:0]
		for p.more() {
			typ := p.enum("STD94 STD96 STD94MULTIBYTE STD96MULTIBYTE COMPLETECODE")
			in.charsets = append(in.charsets, charset{typ: typ, tail: string(p.str())})
		}
	case eColourModel:
		in.colourModel = p.index()
		switch in.colourModel {
		case 1:
			in.pr.components = 3
		case 4:
			in.pr.components = 4
		default:
			in.pr.components = 3
			in.warnOnce("colour model", "colour model %d is read as RGB", in.colourModel)
		}
	}
}

// bitsOf returns a precision in bits, or old when v is not one.
func bitsOf(v, old int) int {
	switch v {
	case 8, 16, 24, 32:
		return v
	}
	return old
}

// realPrecision reads REAL PRECISION or VDC REAL PRECISION: the form
// (floating or fixed), and the sizes of the exponent and fraction (or of
// the whole and fractional parts).
func (in *interp) realPrecision(p params, old realFormat) realFormat {
	form := p.enum("FLOATING FIXED")
	a, b := p.int(), p.int()
	switch {
	case form == 0 && a == 9 && b == 23:
		return realFormat{bits: 32}
	case form == 0 && a == 12 && b == 52:
		return realFormat{bits: 64}
	case form == 1 && a == 16 && b == 16:
		return realFormat{fixed: true, bits: 32}
	case form == 1 && a == 32 && b == 32:
		return realFormat{fixed: true, bits: 64}
	}
	in.warnOnce("real precision", "unsupported real precision (%d, %d, %d) is ignored", form, a, b)
	return old
}

// readColour reads a colour in the colour selection mode.
func (in *interp) readColour(p params) colour {
	if in.st.colourMode == 1 {
		return colour{direct: true, c: in.readDirect(p)}
	}
	return colour{index: p.colourIndex()}
}

// readDirect reads a direct colour.
func (in *interp) readDirect(p params) bdf.Color {
	n := p.components()
	var v [4]float64
	for i := 0; i < n; i++ {
		v[i] = p.component()
	}
	return in.rgb(v, n, p.colourMax())
}

// rgb converts the components of a direct colour, scaled by the colour
// value extent.
func (in *interp) rgb(v [4]float64, n int, max float64) bdf.Color {
	var lo, hi [4]float64
	if in.colourExtSet {
		lo, hi = in.colourExt[0], in.colourExt[1]
	} else {
		for i := range hi {
			hi[i] = max
		}
	}
	var f [4]float64
	for i := 0; i < n; i++ {
		if d := hi[i] - lo[i]; d != 0 {
			f[i] = (v[i] - lo[i]) / d
		}
	}
	return rgbUnit(f, n)
}

// rgbUnit converts colour components from 0 to 1: red, green and blue, or
// cyan, magenta, yellow and black.
func rgbUnit(f [4]float64, n int) bdf.Color {
	for i := range f {
		f[i] = math.Min(math.Max(f[i], 0), 1)
	}
	if n == 4 {
		k := f[3]
		f[0], f[1], f[2] = (1-f[0])*(1-k), (1-f[1])*(1-k), (1-f[2])*(1-k)
	}
	return bdf.RGB(uint8(math.Round(f[0]*255)), uint8(math.Round(f[1]*255)), uint8(math.Round(f[2]*255)))
}

// defaultTable holds the colours of the indexes a picture has not set.
var defaultTable = []bdf.Color{
	bdf.RGB(255, 255, 255), bdf.RGB(0, 0, 0), bdf.RGB(255, 0, 0), bdf.RGB(0, 255, 0),
	bdf.RGB(0, 0, 255), bdf.RGB(255, 255, 0), bdf.RGB(0, 255, 255), bdf.RGB(255, 0, 255),
}

// colourOf resolves a colour: index 0 is the background colour unless the
// colour table sets it.
func (in *interp) colourOf(c colour) bdf.Color {
	if c.direct {
		return c.c
	}
	if v, ok := in.st.table[c.index]; ok {
		return v
	}
	if c.index == 0 {
		return in.st.background
	}
	if c.index > 0 && c.index < len(defaultTable) {
		return defaultTable[c.index]
	}
	return bdf.RGB(0, 0, 0)
}

// readSize reads a width or size (VDC in absolute mode, a real otherwise).
func (in *interp) readSize(p params, mode int) size {
	if mode == modeAbsolute {
		return size{p.vdc(), mode}
	}
	return size{p.real(), mode}
}

// modeEnum lists the clear text words of the specification modes.
const modeEnum = "ABS|ABSOLUTE SCALED FRACTIONAL|FRACTION MM|METRIC"

func (in *interp) pictureDescriptor(c code, p params) {
	st := &in.st
	switch c {
	case eScalingMode:
		st.scaling = p.enum("ABSTRACT METRIC")
		if st.scaling == 1 && p.more() {
			st.metric = p.scaleFactor()
		}
	case eColourSelectionMode:
		st.colourMode = p.enum("INDEXED DIRECT")
	case eLineWidthMode:
		st.lineMode = p.enum(modeEnum)
	case eMarkerSizeMode:
		st.markerMode = p.enum(modeEnum)
	case eEdgeWidthMode:
		st.edgeMode = p.enum(modeEnum)
	case eInteriorStyleMode:
		st.interiorMode = p.enum(modeEnum)
	case eVDCExtent:
		st.extent[0], st.extent[1] = p.point(), p.point()
		st.extentSet = true
	case eBackgroundColour:
		st.background = in.readDirect(p)
	case eDeviceViewportMode:
		st.viewportMode = p.enum("FRACTION MM PHYDEVCOORD")
		if p.more() {
			st.viewportScale = p.real()
		}
	case eDeviceViewport:
		if st.viewportMode == 2 {
			st.viewport[0] = cad.Point{X: float64(p.int()), Y: float64(p.int())}
			st.viewport[1] = cad.Point{X: float64(p.int()), Y: float64(p.int())}
		} else {
			st.viewport[0] = cad.Point{X: p.real(), Y: p.real()}
			st.viewport[1] = cad.Point{X: p.real(), Y: p.real()}
		}
		st.viewportSet = true
	case eLineRep, eMarkerRep, eEdgeRep:
		i := p.index()
		mode := st.lineMode
		reps := st.lineReps
		switch c {
		case eMarkerRep:
			mode, reps = st.markerMode, st.markerReps
		case eEdgeRep:
			mode, reps = st.edgeMode, st.edgeRs
		}
		r := lineRep{typ: p.index()}
		r.width = in.readSize(p, mode)
		r.colour = in.readColour(p)
		reps[i] = r
	case eTextRep:
		i := p.index()
		r := textRep{font: p.index(), prec: p.enum("STRING CHAR STROKE")}
		r.spacing, r.expansion = p.real(), p.real()
		r.colour = in.readColour(p)
		st.textReps[i] = r
	case eFillRep:
		i := p.index()
		r := fillRep{style: p.enum(interiorEnum)}
		r.colour = in.readColour(p)
		r.hatch, r.pattern = p.index(), p.index()
		st.fillReps[i] = r
	case eLineEdgeTypeDef:
		i := p.index()
		d := lineTypeDef{repeat: in.readSize(p, st.lineMode)}
		for p.more() {
			d.dash = append(d.dash, math.Abs(float64(p.int())))
		}
		if i < 0 {
			st.lineTypes[i] = d
		}
	case eHatchStyleDef:
		i := p.index()
		d := hatchDef{cross: p.enum("PARALLEL CROSSHATCH") == 1}
		m := st.interiorMode
		d.dir1 = cad.Point{X: in.readSize(p, m).v, Y: in.readSize(p, m).v}
		d.dir2 = cad.Point{X: in.readSize(p, m).v, Y: in.readSize(p, m).v}
		d.cycle = in.readSize(p, m)
		n := p.int()
		for j := 0; j < n && p.more(); j++ {
			d.gaps = append(d.gaps, math.Abs(float64(p.int())))
		}
		for j := 0; j < n && p.more(); j++ {
			d.types = append(d.types, p.index())
		}
		if i < 0 {
			st.hatchStyles[i] = d
		}
	case eGeoPatternDef:
		in.warnOnce("geopat", "geometric patterns are filled with the fill colour")
	}
}

const interiorEnum = "HOLLOW SOLID PAT|PATTERN HATCH EMPTY GEOPAT|GEOPATTERN INTERP|INTERPOLATED"

func (in *interp) control(c code, p params) {
	st := &in.st
	switch c {
	case eVDCIntegerPrecision:
		if !in.text {
			in.pr.vdcIntBits = bitsOf(p.int(), in.pr.vdcIntBits)
		}
	case eVDCRealPrecision:
		if !in.text {
			in.pr.vdcRealFormat = in.realPrecision(p, in.pr.vdcRealFormat)
		}
	case eAuxiliaryColour:
		st.aux = in.readColour(p)
	case eTransparency:
		st.transparency = p.enum("OFF ON") == 1
	case eClipRectangle:
		st.clip[0], st.clip[1] = p.point(), p.point()
		st.clipSet = true
		in.updateClip()
	case eClipIndicator:
		st.clipOn = p.enum("OFF ON") == 1
		in.updateClip()
	case eNewRegion:
		if in.fig != nil {
			in.fig.path.Close()
		}
	case eSaveContext:
		st.saved[p.name()] = st.a
	case eRestoreContext:
		if a, ok := st.saved[p.name()]; ok {
			st.a = a
		}
	case eTransparentCell:
		if p.enum("OFF ON") == 1 {
			c := in.readColour(p)
			st.transparentCell = &c
		} else {
			st.transparentCell = nil
		}
	}
}

func (in *interp) attribute(c code, p params) {
	a := &in.st.a
	st := &in.st
	switch c {
	case eLineBundleIndex:
		a.lineBundle = p.index()
	case eLineType:
		a.lineType = p.index()
	case eLineWidth:
		a.lineWidth = in.readSize(p, st.lineMode)
	case eLineColour:
		a.lineColour = in.readColour(p)
	case eMarkerBundleIndex:
		a.markerBundle = p.index()
	case eMarkerType:
		a.markerType = p.index()
	case eMarkerSize:
		a.markerSize = in.readSize(p, st.markerMode)
	case eMarkerColour:
		a.markerColour = in.readColour(p)
	case eTextBundleIndex:
		a.textBundle = p.index()
	case eTextFontIndex:
		a.font = p.index()
	case eCharExpansion:
		a.expansion = p.real()
	case eCharSpacing:
		a.spacing = p.real()
	case eTextColour:
		a.textColour = in.readColour(p)
	case eCharHeight:
		a.charHeight = math.Abs(p.vdc())
	case eCharOrientation:
		a.up = cad.Point{X: p.vdc(), Y: p.vdc()}
		a.base = cad.Point{X: p.vdc(), Y: p.vdc()}
	case eTextPath:
		a.path = p.enum("RIGHT LEFT UP DOWN")
	case eTextAlignment:
		a.hAlign = p.enum("NORMHORIZ LEFT CTR|CENTRE|CENTER RIGHT CONTHORIZ")
		a.vAlign = p.enum("NORMVERT TOP CAP HALF BASE BOTTOM CONTVERT")
		if p.more() {
			a.contH = p.real()
			a.contV = p.real()
		}
	case eCharSetIndex:
		a.charSet = p.index()
	case eAltCharSetIndex:
		a.altCharSet = p.index()
	case eFillBundleIndex:
		a.fillBundle = p.index()
	case eInteriorStyle:
		a.interior = p.enum(interiorEnum)
	case eFillColour:
		a.fillColour = in.readColour(p)
	case eHatchIndex:
		a.hatch = p.index()
	case ePatternIndex:
		a.pattern = p.index()
	case eEdgeBundleIndex:
		a.edgeBundle = p.index()
	case eEdgeType:
		a.edgeType = p.index()
	case eEdgeWidth:
		a.edgeWidth = in.readSize(p, st.edgeMode)
	case eEdgeColour:
		a.edgeColour = in.readColour(p)
	case eEdgeVisibility:
		a.edgeVisible = p.enum("OFF ON") == 1
	case eFillReferencePoint:
		a.fillRef, a.fillRefSet = p.point(), true
	case ePatternTable:
		i := p.index()
		nx, ny := p.int(), p.int()
		prec := p.int()
		d := patternDef{nx: nx, ny: ny}
		if in.cells(p, nx, ny, prec, 1, maxPattern, func(_, _ int, c colour) { d.cells = append(d.cells, c) }) && i > 0 {
			st.patterns[i] = d
			in.generation++
		}
	case ePatternSize:
		for i := range a.patternSize {
			a.patternSize[i] = in.readSize(p, st.interiorMode).v
		}
		a.patternSizeSet = true
	case eColourTable:
		i := p.colourIndex()
		for p.more() && i < 1<<20 {
			st.table[i] = in.readDirect(p)
			i++
		}
		in.generation++
	case eASF:
		for p.more() {
			typ := p.enum("LINETYPE LINEWIDTH LINECOLR MARKERTYPE MARKERSIZE MARKERCOLR TEXTFONTINDEX TEXTPREC " +
				"CHAREXP|CHAREXPAN CHARSPACE TEXTCOLR INTSTYLE FILLCOLR HATCHINDEX PATINDEX EDGETYPE EDGEWIDTH EDGECOLR")
			bundled := p.enum("INDIV BUNDLED") == 1
			switch {
			case typ >= 0 && typ < asfCount:
				a.asf[typ] = bundled
			case typ >= 506 && typ <= 511:
				// all edge, line, marker, text, fill, or all attributes
				groups := map[int][]int{
					506: {asfEdgeType, asfEdgeWidth, asfEdgeColour},
					507: {asfLineType, asfLineWidth, asfLineColour},
					508: {asfMarkerType, asfMarkerSize, asfMarkerColour},
					509: {asfFont, asfTextPrecision, asfExpansion, asfSpacing, asfTextColour},
					510: {asfInteriorStyle, asfFillColour, asfHatchIndex, asfPatternIndex},
				}
				if typ == 511 {
					for j := range a.asf {
						a.asf[j] = bundled
					}
				}
				for _, j := range groups[typ] {
					a.asf[j] = bundled
				}
			}
		}
	case eLineCap:
		a.lineCap = p.index()
	case eLineJoin:
		a.lineJoin = p.index()
	case eLineTypeOffset:
		a.lineTypeOffset = p.real()
	case eEdgeCap:
		a.edgeCap = p.index()
	case eEdgeJoin:
		a.edgeJoin = p.index()
	case eEdgeTypeOffset:
		a.edgeTypeOffset = p.real()
	case eRestrictedTextType:
		a.restrictedType = p.index()
	case eInterpolatedInt:
		ii := &interpolated{style: p.index()}
		n := map[int]int{1: 2, 2: 4}[ii.style]
		for j := 0; j < n; j++ {
			ii.geom = append(ii.geom, in.readSize(p, st.interiorMode).v)
		}
		k := p.int()
		for j := 0; j < k && p.more(); j++ {
			ii.stages = append(ii.stages, p.real())
		}
		for p.more() {
			ii.colours = append(ii.colours, in.readColour(p))
		}
		a.interp = ii
	}
}

// beginPicture starts a picture: the picture state returns to its
// defaults and the metafile's defaults replacement.
func (in *interp) beginPicture() {
	in.count++
	in.pic = &picture{number: in.count, d: &cad.Drawing{}}
	in.skip = in.selected != nil && !in.selected(in.count)
	in.body = false
	in.st = newState()
	in.pr.vdcIntBits = 16
	in.pr.vdcRealFormat = realFormat{fixed: true, bits: 32}
	in.out = in.pic.d
	in.clipped = false
	in.fig, in.compound, in.pending, in.tiles = nil, nil, nil, nil
	in.segs = map[int]*cad.Drawing{}
	in.weights = map[int]int{}
	in.patternImages = map[[2]int]bdf.Hash{}
	in.open = nil
	in.aps, in.hidden = nil, 0
	in.replayDefaults()
}

// replayDefaults interprets the elements of METAFILE DEFAULTS
// REPLACEMENT: picture descriptor, control and attribute elements only.
func (in *interp) replayDefaults() {
	for i := range in.defaults {
		if c := in.defaults[i].class; c == 2 || c == 3 || c == 5 {
			in.element(&in.defaults[i])
		}
	}
}

// ensureBody begins a picture and its body when a file draws without
// them.
func (in *interp) ensureBody() {
	if in.pic == nil {
		in.beginPicture()
	}
	if !in.body {
		in.beginBody()
	}
}

// extent returns the VDC extent: the picture's, or the default one.
func (in *interp) extent() [2]cad.Point {
	if in.st.extentSet {
		return in.st.extent
	}
	if in.pr.vdcReal {
		return [2]cad.Point{{}, {X: 1, Y: 1}}
	}
	return [2]cad.Point{{}, {X: 32767, Y: 32767}}
}

// maxPage is the largest side of a page (mm): 200 inches, as PDF.
const maxPage = 5080

// beginBody fixes the page of the picture: its size comes from the metric
// scaling factor, or a device viewport in millimetres; abstract pictures
// are fitted to 297 mm (the long side of A4).
func (in *interp) beginBody() {
	in.body = true
	pic := in.pic
	ext := in.extent()
	dx, dy := ext[1].X-ext[0].X, ext[1].Y-ext[0].Y
	if dx == 0 || !isFinite(dx) {
		dx = 1
	}
	if dy == 0 || !isFinite(dy) {
		dy = 1
	}
	adx, ady := math.Abs(dx), math.Abs(dy)
	pic.long = math.Max(adx, ady)
	st := &in.st
	var s float64 // mm per VDC unit
	var vw, vh float64
	switch {
	case st.scaling == 1 && st.metric > 0 && isFinite(st.metric):
		s = st.metric
	case st.viewportMode == 1 && st.viewportSet:
		vw = math.Abs(st.viewport[1].X-st.viewport[0].X) * st.viewportScale
		vh = math.Abs(st.viewport[1].Y-st.viewport[0].Y) * st.viewportScale
		if vw > 0 && vh > 0 && isFinite(vw) && isFinite(vh) {
			s = math.Min(vw/adx, vh/ady)
		} else {
			vw, vh = 0, 0
		}
	}
	w, h := adx*s, ady*s
	if !(math.Max(w, h) >= 1) || !isFinite(w) || !isFinite(h) {
		s = 297 / pic.long
		w, h, vw, vh = adx*s, ady*s, 0, 0
	}
	if l := math.Max(math.Max(w, h), math.Max(vw, vh)); l > maxPage {
		f := maxPage / l
		s, w, h, vw, vh = s*f, w*f, h*f, vw*f, vh*f
		in.warnOnce("page size", "pictures larger than %d mm are scaled down to fit", maxPage)
	}
	if vw > 0 {
		w, h = vw, vh
	}
	pic.w, pic.h = math.Max(w, 0.1), math.Max(h, 0.1)
	k := s * mm
	pic.k = k
	sx, sy := math.Copysign(k, dx), math.Copysign(k, dy)
	pic.m = canvas.Matrix{sx, 0, 0, -sy, -ext[0].X * sx, pic.h*mm + ext[0].Y*sy}
	pic.bg = st.background
	in.updateClip()
}

func isFinite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

// updateClip opens a clip group for the clip rectangle, when clipping is
// on and the rectangle is smaller than the page.
func (in *interp) updateClip() {
	if in.pic == nil || !in.body {
		return
	}
	d := in.pic.d
	if in.clipped {
		d.End()
		in.clipped = false
	}
	st := &in.st
	if !st.clipOn || !st.clipSet {
		return
	}
	ext := in.extent()
	r := cad.Rect{}.Add(st.clip[0]).Add(st.clip[1])
	e := cad.Rect{}.Add(ext[0]).Add(ext[1])
	eps := in.pic.long * 1e-9
	if r.Min.X <= e.Min.X+eps && r.Min.Y <= e.Min.Y+eps && r.Max.X >= e.Max.X-eps && r.Max.Y >= e.Max.Y-eps {
		return // the page clips as much
	}
	clip := (&cad.Path{}).Polyline([]cad.Point{r.Min, {X: r.Max.X, Y: r.Min.Y}, r.Max, {X: r.Min.X, Y: r.Max.Y}}, true)
	d.Begin(clip)
	in.clipped = true
}

// target returns where primitives go: the innermost segment being
// recorded, or the picture.
func (in *interp) target() *cad.Drawing {
	if n := len(in.open); n > 0 {
		return in.open[n-1].d
	}
	return in.pic.d
}

// endPicture finishes a picture: open structures are closed and the page
// is kept unless it was left out.
func (in *interp) endPicture() {
	if !in.body {
		in.beginBody()
	}
	in.flushText()
	for len(in.open) > 0 {
		in.element(&element{code: eEndSegment})
	}
	if in.fig != nil {
		in.element(&element{code: eEndFigure})
	}
	if in.compound != nil {
		in.element(&element{code: eEndCompoundLine})
	}
	if in.clipped {
		in.pic.d.End()
		in.clipped = false
	}
	if !in.skip {
		in.pictures = append(in.pictures, in.pic)
	}
	in.pic = nil
}

func (in *interp) segmentElement(c code, p params) {
	in.flushText()
	switch c {
	case eCopySegment:
		id := p.name()
		var m canvas.Matrix
		for i := 0; i < 4; i++ {
			m[i] = p.real()
		}
		m[4], m[5] = p.vdc(), p.vdc()
		d, ok := in.segs[id]
		if !ok {
			in.warnOnce("copyseg", "copies of segments that are not in the picture are left out")
			return
		}
		if in.drawing() && !d.Empty() {
			in.points += in.weights[id]
			if in.points <= maxPoints {
				in.out.Append(d, m, nil)
			}
		}
	case eSegmentTransf:
		in.warnOnce("segtran", "segment transformations are not applied")
	}
}

// scanElement interprets an element for the scan of strings: pictures,
// the metafile descriptor, VDC precisions and the strings of text.
func (in *interp) scanElement(e *element, p params) {
	switch e.code {
	case eBeginPicture:
		in.count++
		in.pr.vdcIntBits = 16
		in.pr.vdcRealFormat = realFormat{fixed: true, bits: 32}
		in.replayDefaults()
	case eVDCIntegerPrecision, eVDCRealPrecision:
		in.control(e.code, p)
	case eText:
		p.point()
		p.enum("NOTFINAL FINAL")
		in.strs = append(in.strs, p.str())
	case eRestrictedText:
		p.vdc()
		p.vdc()
		p.point()
		p.enum("NOTFINAL FINAL")
		in.strs = append(in.strs, p.str())
	case eAppendText:
		p.enum("NOTFINAL FINAL")
		in.strs = append(in.strs, p.str())
	default:
		if e.class == 1 {
			in.descriptor(e, p)
		}
	}
}
