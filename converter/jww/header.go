package jww

import "errors"

// pen is how a line color prints: its color (a Windows COLORREF, 0x00BBGGRR),
// the width of its lines and the radius of its points (mm).
type pen struct {
	color  uint32
	width  int
	radius float64
}

// lineType is a line type of dots: a bit pattern of unit bits, each bit
// pitch dots long on screen and prtPitch dots on paper.
type lineType struct {
	mask                  uint32
	unit, pitch, prtPitch int
}

// udLineType is a line type of SXF: dash and gap lengths in mm.
type udLineType struct {
	name     string
	segments int
	pitch    [10]float64
}

type layerGroup struct {
	state  int // 0 hidden, 1 shown only, 2 editable, 3 being written
	scale  float64
	layers [16]int
	name   string
	names  [16]string
}

// mojiShu is a predefined text size (文字種).
type mojiShu struct {
	w, h, spacing float64
	color         int
}

// header is the part of a JWW file before its data.
type header struct {
	version    int
	memo       string
	paper      int
	groups     [16]layerGroup
	dimension  [5]uint32
	maxWidth   int32
	screen     [10]uint32 // screen colors (COLORREF)
	prtPen     [10]pen
	lineTypes  map[int]lineType
	sxfPens    [257]pen
	sxfTypes   [33]udLineType
	hasSXF     bool
	drawPrtTen bool
	colorPrint bool
	dpi        int
	moji       [11]mojiShu
}

var errNotJWW = errors.New("not a JWW file")

// readHeader reads the header in the order of Jw_cad's data format (the
// version 7.02 description, jwdatafmt.txt), with the fields that appeared
// in later versions read only from files of those versions.
func readHeader(a *archive) (*header, error) {
	if len(a.b) < 12 || string(a.b[:8]) != "JwwData." {
		return nil, errNotJWW
	}
	a.i = 8
	h := &header{lineTypes: map[int]lineType{}, dpi: 300}
	h.version = int(a.u32())
	h.memo = a.str()
	h.paper = int(a.u32())
	a.u32() // the layer group being written
	for g := range h.groups {
		lg := &h.groups[g]
		lg.state = int(a.u32())
		a.u32() // the layer being written
		lg.scale = a.f64()
		a.u32() // protection
		for l := range lg.layers {
			lg.layers[l] = int(a.u32())
			a.u32() // protection
		}
	}
	a.skip(14 * 4) // dummies
	for i := range h.dimension {
		h.dimension[i] = a.u32()
	}
	a.u32() // dummy
	h.maxWidth = a.i32()
	a.skip(3 * 8) // print origin and magnification
	a.u32()       // print rotation and reference point
	a.u32()       // scale mode
	a.skip(5 * 8) // scale spacing and reference point
	for g := range h.groups {
		for l := range h.groups[g].names {
			h.groups[g].names[l] = a.str()
		}
	}
	for g := range h.groups {
		h.groups[g].name = a.str()
	}
	a.skip(2 * 8) // shadows: level, latitude
	a.u32()
	a.skip(8)
	if h.version >= 300 {
		a.skip(2 * 8) // sky map
	}
	a.u32()       // 2.5D units
	a.skip(3 * 8) // screen magnification and origin
	a.skip(3 * 8) // stored range
	if h.version >= 300 {
		for range 8 {
			a.skip(3 * 8)
			a.u32()
		}
		a.skip(3 * 8)
		a.u32()
		a.skip(3 * 8)
		a.u32()
	} else {
		a.skip(4 * 3 * 8)
	}
	a.skip(10 * 8) // parallel line spacings
	a.skip(8)
	for n := range h.screen {
		h.screen[n] = a.u32()
		a.u32() // the width on screen
	}
	for n := range h.prtPen {
		h.prtPen[n] = pen{color: a.u32(), width: int(a.u32()), radius: a.f64()}
	}
	readType := func(n int) {
		h.lineTypes[n] = lineType{mask: a.u32(), unit: int(a.u32()), pitch: int(a.u32()), prtPitch: int(a.u32())}
	}
	for n := 2; n <= 9; n++ {
		readType(n)
	}
	for n := 11; n <= 15; n++ {
		// random lines: pattern, amplitude and pitch on screen and paper
		lt := lineType{mask: a.u32()}
		a.u32()
		lt.pitch = int(a.u32())
		a.u32()
		lt.prtPitch = int(a.u32())
		h.lineTypes[n] = lt
	}
	for n := 16; n <= 19; n++ {
		readType(n)
	}
	a.u32() // points drawn at their radius on screen
	h.drawPrtTen = a.u32() != 0
	a.u32() // drawing order
	a.u32() // reverse drawing
	a.u32() // reverse search
	h.colorPrint = a.u32() != 0
	a.u32() // print by layer
	a.u32() // print by color
	a.u32() // continuous printing
	a.u32() // grey printing
	n1 := a.u32()
	if h.version >= 600 && n1/10%10 == 2 {
		h.dpi = 600
	}
	if h.version >= 223 {
		a.u32()       // drawing time
		a.u32()       // 2.5D flags
		a.skip(3 * 4) // eye angles
		a.skip(5 * 8) // eye positions
	}
	if h.version >= 225 {
		a.skip(4 * 8) // last lengths and radius
	}
	if h.version >= 230 {
		a.u32() // solids in any color
		a.u32() // the default color of solids
	}
	if h.version >= 420 {
		h.hasSXF = true
		a.skip(257 * 2 * 4) // screen colors and widths
		for n := range h.sxfPens {
			a.str() // the name of the color
			h.sxfPens[n] = pen{color: a.u32(), width: int(a.u32()), radius: a.f64()}
		}
		for n := 0; n <= 32; n++ {
			readType(n + 30)
		}
		for n := range h.sxfTypes {
			t := &h.sxfTypes[n]
			t.name = a.str()
			t.segments = int(a.u32())
			for j := range t.pitch {
				t.pitch[j] = a.f64()
			}
		}
	}
	for i := 1; i <= 10; i++ {
		h.moji[i] = mojiShu{w: a.f64(), h: a.f64(), spacing: a.f64(), color: int(a.u32())}
	}
	a.skip(3 * 8) // the text being written
	a.skip(2 * 4)
	a.skip(2 * 8) // text arrangement
	a.u32()       // offsets of the text reference point
	a.skip(6 * 8)
	return h, a.err
}
