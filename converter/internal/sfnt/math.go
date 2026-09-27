package sfnt

// Math is the OpenType MATH table of a font: the constants, the per-glyph
// information and the glyph constructions that lay formulas out. Values
// are in font units.
type Math struct {
	Constants MathConstants
	// MinConnectorOverlap is the overlap of the parts of a glyph assembly.
	MinConnectorOverlap int

	italics   map[uint16]int
	topAccent map[uint16]int
	extended  map[uint16]bool
	kerns     map[uint16]*[4]mathKern
	vert      map[uint16]*GlyphConstruction
	horiz     map[uint16]*GlyphConstruction
}

// MathConstants are the global layout parameters of the MATH table, in the
// order of the table (font units, percentages for the three that say so).
type MathConstants struct {
	ScriptPercentScaleDown                   int
	ScriptScriptPercentScaleDown             int
	DelimitedSubFormulaMinHeight             int
	DisplayOperatorMinHeight                 int
	MathLeading                              int
	AxisHeight                               int
	AccentBaseHeight                         int
	FlattenedAccentBaseHeight                int
	SubscriptShiftDown                       int
	SubscriptTopMax                          int
	SubscriptBaselineDropMin                 int
	SuperscriptShiftUp                       int
	SuperscriptShiftUpCramped                int
	SuperscriptBottomMin                     int
	SuperscriptBaselineDropMax               int
	SubSuperscriptGapMin                     int
	SuperscriptBottomMaxWithSubscript        int
	SpaceAfterScript                         int
	UpperLimitGapMin                         int
	UpperLimitBaselineRiseMin                int
	LowerLimitGapMin                         int
	LowerLimitBaselineDropMin                int
	StackTopShiftUp                          int
	StackTopDisplayStyleShiftUp              int
	StackBottomShiftDown                     int
	StackBottomDisplayStyleShiftDown         int
	StackGapMin                              int
	StackDisplayStyleGapMin                  int
	StretchStackTopShiftUp                   int
	StretchStackBottomShiftDown              int
	StretchStackGapAboveMin                  int
	StretchStackGapBelowMin                  int
	FractionNumeratorShiftUp                 int
	FractionNumeratorDisplayStyleShiftUp     int
	FractionDenominatorShiftDown             int
	FractionDenominatorDisplayStyleShiftDown int
	FractionNumeratorGapMin                  int
	FractionNumDisplayStyleGapMin            int
	FractionRuleThickness                    int
	FractionDenominatorGapMin                int
	FractionDenomDisplayStyleGapMin          int
	SkewedFractionHorizontalGap              int
	SkewedFractionVerticalGap                int
	OverbarVerticalGap                       int
	OverbarRuleThickness                     int
	OverbarExtraAscender                     int
	UnderbarVerticalGap                      int
	UnderbarRuleThickness                    int
	UnderbarExtraDescender                   int
	RadicalVerticalGap                       int
	RadicalDisplayStyleVerticalGap           int
	RadicalRuleThickness                     int
	RadicalExtraAscender                     int
	RadicalKernBeforeDegree                  int
	RadicalKernAfterDegree                   int
	RadicalDegreeBottomRaisePercent          int
}

// GlyphConstruction is how a glyph grows (MathGlyphConstruction): larger
// variants of it, then an assembly of parts for any larger size.
type GlyphConstruction struct {
	Variants []GlyphVariant
	Parts    []GlyphPart // bottom to top, or left to right; nil when there is no assembly
	// ItalicsCorrection is that of the assembly.
	ItalicsCorrection int
}

// GlyphVariant is a larger form of a glyph and its size in the direction
// it grows.
type GlyphVariant struct {
	Glyph   uint16
	Advance int
}

// GlyphPart is a part of a glyph assembly.
type GlyphPart struct {
	Glyph                        uint16
	StartConnector, EndConnector int
	FullAdvance                  int
	Extender                     bool // repeated as many times as needed (or not at all)
}

// Math kern corners.
const (
	KernTopRight = iota
	KernTopLeft
	KernBottomRight
	KernBottomLeft
)

type mathKern struct {
	heights []int
	kerns   []int // one more than heights
}

// ParseMath reads the MATH table of f; it returns nil when the font has
// none or it cannot be read.
func ParseMath(f *Font) *Math {
	b := f.Tables["MATH"]
	if len(b) < 10 || be16(b, 0) != 1 {
		return nil
	}
	m := &Math{italics: map[uint16]int{}, topAccent: map[uint16]int{}, extended: map[uint16]bool{},
		kerns: map[uint16]*[4]mathKern{}, vert: map[uint16]*GlyphConstruction{}, horiz: map[uint16]*GlyphConstruction{}}
	if !m.parseConstants(sub(b, int(be16(b, 4)))) {
		return nil
	}
	m.parseGlyphInfo(sub(b, int(be16(b, 6))))
	m.parseVariants(sub(b, int(be16(b, 8))))
	return m
}

// sub returns b from off, or nil when off is 0 or out of range.
func sub(b []byte, off int) []byte {
	if off <= 0 || off >= len(b) {
		return nil
	}
	return b[off:]
}

func (m *Math) parseConstants(b []byte) bool {
	c := &m.Constants
	// the first four are plain values, the last a percentage; the others
	// MathValueRecords (a value and a device table offset)
	plain := []*int{&c.ScriptPercentScaleDown, &c.ScriptScriptPercentScaleDown, &c.DelimitedSubFormulaMinHeight, &c.DisplayOperatorMinHeight}
	records := []*int{&c.MathLeading, &c.AxisHeight, &c.AccentBaseHeight, &c.FlattenedAccentBaseHeight,
		&c.SubscriptShiftDown, &c.SubscriptTopMax, &c.SubscriptBaselineDropMin,
		&c.SuperscriptShiftUp, &c.SuperscriptShiftUpCramped, &c.SuperscriptBottomMin, &c.SuperscriptBaselineDropMax,
		&c.SubSuperscriptGapMin, &c.SuperscriptBottomMaxWithSubscript, &c.SpaceAfterScript,
		&c.UpperLimitGapMin, &c.UpperLimitBaselineRiseMin, &c.LowerLimitGapMin, &c.LowerLimitBaselineDropMin,
		&c.StackTopShiftUp, &c.StackTopDisplayStyleShiftUp, &c.StackBottomShiftDown, &c.StackBottomDisplayStyleShiftDown,
		&c.StackGapMin, &c.StackDisplayStyleGapMin,
		&c.StretchStackTopShiftUp, &c.StretchStackBottomShiftDown, &c.StretchStackGapAboveMin, &c.StretchStackGapBelowMin,
		&c.FractionNumeratorShiftUp, &c.FractionNumeratorDisplayStyleShiftUp, &c.FractionDenominatorShiftDown,
		&c.FractionDenominatorDisplayStyleShiftDown, &c.FractionNumeratorGapMin, &c.FractionNumDisplayStyleGapMin,
		&c.FractionRuleThickness, &c.FractionDenominatorGapMin, &c.FractionDenomDisplayStyleGapMin,
		&c.SkewedFractionHorizontalGap, &c.SkewedFractionVerticalGap,
		&c.OverbarVerticalGap, &c.OverbarRuleThickness, &c.OverbarExtraAscender,
		&c.UnderbarVerticalGap, &c.UnderbarRuleThickness, &c.UnderbarExtraDescender,
		&c.RadicalVerticalGap, &c.RadicalDisplayStyleVerticalGap, &c.RadicalRuleThickness, &c.RadicalExtraAscender,
		&c.RadicalKernBeforeDegree, &c.RadicalKernAfterDegree}
	if len(b) < 2*len(plain)+4*len(records)+2 {
		return false
	}
	p := 0
	for _, v := range plain {
		*v = int(int16(be16(b, p)))
		p += 2
	}
	for _, v := range records {
		*v = int(int16(be16(b, p)))
		p += 4
	}
	c.RadicalDegreeBottomRaisePercent = int(int16(be16(b, p)))
	return true
}

// coverage reads a Coverage table into glyph → coverage index.
func coverage(b []byte) map[uint16]int {
	out := map[uint16]int{}
	if len(b) < 4 {
		return out
	}
	switch be16(b, 0) {
	case 1:
		n := int(be16(b, 2))
		for i := 0; i < n && 4+i*2+2 <= len(b); i++ {
			out[be16(b, 4+i*2)] = i
		}
	case 2:
		n := int(be16(b, 2))
		for i := 0; i < n && 4+i*6+6 <= len(b); i++ {
			r := 4 + i*6
			start, end, idx := int(be16(b, r)), int(be16(b, r+2)), int(be16(b, r+4))
			for g := start; g <= end && g-start < 0x10000; g++ {
				out[uint16(g)] = idx + g - start
			}
		}
	}
	return out
}

// valueTable reads a table of a coverage and a MathValueRecord per glyph.
func valueTable(b []byte, into map[uint16]int) {
	if len(b) < 4 {
		return
	}
	cov := coverage(sub(b, int(be16(b, 0))))
	n := int(be16(b, 2))
	for g, i := range cov {
		if i < n && 4+i*4+2 <= len(b) {
			into[g] = int(int16(be16(b, 4+i*4)))
		}
	}
}

func (m *Math) parseGlyphInfo(b []byte) {
	if len(b) < 8 {
		return
	}
	valueTable(sub(b, int(be16(b, 0))), m.italics)
	valueTable(sub(b, int(be16(b, 2))), m.topAccent)
	for g := range coverage(sub(b, int(be16(b, 4)))) {
		m.extended[g] = true
	}
	if k := sub(b, int(be16(b, 6))); len(k) >= 4 {
		cov := coverage(sub(k, int(be16(k, 0))))
		n := int(be16(k, 2))
		for g, i := range cov {
			rec := 4 + i*8
			if i >= n || rec+8 > len(k) {
				continue
			}
			var corners [4]mathKern
			for c := 0; c < 4; c++ {
				corners[c] = parseKern(sub(k, int(be16(k, rec+c*2))))
			}
			m.kerns[g] = &corners
		}
	}
}

func parseKern(b []byte) mathKern {
	if len(b) < 2 {
		return mathKern{}
	}
	n := int(be16(b, 0))
	if len(b) < 2+(2*n+1)*4 {
		return mathKern{}
	}
	k := mathKern{heights: make([]int, n), kerns: make([]int, n+1)}
	for i := 0; i < n; i++ {
		k.heights[i] = int(int16(be16(b, 2+i*4)))
	}
	for i := 0; i <= n; i++ {
		k.kerns[i] = int(int16(be16(b, 2+(n+i)*4)))
	}
	return k
}

func (m *Math) parseVariants(b []byte) {
	if len(b) < 10 {
		return
	}
	m.MinConnectorOverlap = int(be16(b, 0))
	vcov := coverage(sub(b, int(be16(b, 2))))
	hcov := coverage(sub(b, int(be16(b, 4))))
	nv, nh := int(be16(b, 6)), int(be16(b, 8))
	read := func(cov map[uint16]int, n, base int, into map[uint16]*GlyphConstruction) {
		for g, i := range cov {
			if i >= n || base+i*2+2 > len(b) {
				continue
			}
			if c := parseConstruction(sub(b, int(be16(b, base+i*2)))); c != nil {
				into[g] = c
			}
		}
	}
	read(vcov, nv, 10, m.vert)
	read(hcov, nh, 10+nv*2, m.horiz)
}

func parseConstruction(b []byte) *GlyphConstruction {
	if len(b) < 4 {
		return nil
	}
	c := &GlyphConstruction{}
	n := int(be16(b, 2))
	for i := 0; i < n && 4+i*4+4 <= len(b); i++ {
		c.Variants = append(c.Variants, GlyphVariant{Glyph: be16(b, 4+i*4), Advance: int(be16(b, 4+i*4+2))})
	}
	if a := sub(b, int(be16(b, 0))); len(a) >= 6 {
		c.ItalicsCorrection = int(int16(be16(a, 0)))
		np := int(be16(a, 4))
		for i := 0; i < np && 6+i*10+10 <= len(a); i++ {
			r := 6 + i*10
			c.Parts = append(c.Parts, GlyphPart{Glyph: be16(a, r), StartConnector: int(be16(a, r+2)),
				EndConnector: int(be16(a, r+4)), FullAdvance: int(be16(a, r+6)), Extender: be16(a, r+8)&1 != 0})
		}
	}
	return c
}

// Italic returns the italics correction of a glyph.
func (m *Math) Italic(g uint16) (int, bool) {
	v, ok := m.italics[g]
	return v, ok
}

// TopAccent returns the horizontal position where an accent attaches
// above a glyph.
func (m *Math) TopAccent(g uint16) (int, bool) {
	v, ok := m.topAccent[g]
	return v, ok
}

// Extended reports whether a glyph is an extended shape (a tall glyph
// whose scripts are placed as for a box, not for a character).
func (m *Math) Extended(g uint16) bool { return m.extended[g] }

// Kern returns the kerning at a corner of a glyph for a script whose edge
// is at height (font units).
func (m *Math) Kern(g uint16, corner int, height int) int {
	ks := m.kerns[g]
	if ks == nil || corner < 0 || corner > 3 {
		return 0
	}
	k := ks[corner]
	if len(k.kerns) == 0 {
		return 0
	}
	i := 0
	for i < len(k.heights) && height >= k.heights[i] {
		i++
	}
	return k.kerns[i]
}

// Vertical returns how a glyph grows vertically (nil when it does not).
func (m *Math) Vertical(g uint16) *GlyphConstruction { return m.vert[g] }

// Horizontal returns how a glyph grows horizontally (nil when it does not).
func (m *Math) Horizontal(g uint16) *GlyphConstruction { return m.horiz[g] }
