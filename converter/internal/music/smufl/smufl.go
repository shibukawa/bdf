// Package smufl holds the music symbols the score engraver draws: glyphs
// of Bravura, the reference font of the Standard Music Font Layout
// (SMuFL), with their metrics and anchors, and the font's engraving
// defaults. The table (bravura.go) is written by tools/gen-smufl; Bravura
// is licensed under the SIL Open Font License 1.1 (OFL.txt).
//
// Glyphs are in font units: 1000 per em, 250 per staff space, y up, the
// origin on the glyph's reference point (a notehead's is its left edge on
// the staff position it sits on).
package smufl

// Space is the size of a staff space in font units.
const Space = 250

// Glyph is an embedded glyph.
type Glyph struct {
	Rune    rune
	Advance int
	// BBox is the bounding box: x0, y0 (south-west), x1, y1 (north-east).
	BBox [4]float64
	// Anchors are SMuFL's named attachment points (stemUpSE, stemDownNW,
	// cutOutNE, …).
	Anchors map[string][2]float64
	// Path is the outline as SVG path data (M, L, Q, C, Z; absolute).
	Path string
}

// Defaults are the engraving defaults of the font, in staff spaces.
type Defaults struct {
	ArrowShaftThickness        float64
	BarlineSeparation          float64
	BeamSpacing                float64
	BeamThickness              float64
	BracketThickness           float64
	DashedBarlineDashLength    float64
	DashedBarlineGapLength     float64
	DashedBarlineThickness     float64
	HBarThickness              float64
	HairpinThickness           float64
	LegerLineExtension         float64
	LegerLineThickness         float64
	LyricLineThickness         float64
	OctaveLineThickness        float64
	PedalLineThickness         float64
	RepeatBarlineDotSeparation float64
	RepeatEndingLineThickness  float64
	SlurEndpointThickness      float64
	SlurMidpointThickness      float64
	StaffLineThickness         float64
	StemThickness              float64
	SubBracketThickness        float64
	TextEnclosureThickness     float64
	ThickBarlineThickness      float64
	ThinBarlineThickness       float64
	ThinThickBarlineSeparation float64
	TieEndpointThickness       float64
	TieMidpointThickness       float64
	TupletBracketThickness     float64
}

// Lookup returns the glyph named name, or nil.
func Lookup(name string) *Glyph { return glyphs[name] }

// Engraving returns the engraving defaults.
func Engraving() Defaults { return defaults }
