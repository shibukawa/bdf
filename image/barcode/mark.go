package barcode

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/shibukawa/bdf"
)

// Mark is what a BARCODE mark of a drawing says (bdf.MarkBarcode,
// docs/spec.md §7.8): the symbol that the drawing after it draws, and the
// box it is drawn in. A printer of labels prints such a symbol with its
// own barcode commands, and a raster of few dots an inch draws its
// modules on whole dots; anything else draws the drawing as it is.
type Mark struct {
	// Kind is the symbology, and Text what the symbol says (with the
	// check digit Encode adds).
	Kind string `json:"kind"`
	Text string `json:"text"`
	// Level is the level of error correction of a QR code.
	Level string `json:"level,omitempty"`
	// Quiet is Options.Quiet: 0 the quiet zone of the symbology, -1
	// none, else its modules.
	Quiet int `json:"quiet,omitempty"`
	// NoText is set when the drawing has no digits under the bars.
	NoText bool `json:"noText,omitempty"`
	// Box is the box the symbol and its quiet zone are drawn in (x, y,
	// width, height), in the units of the drawing after the mark: the box
	// of Symbol.Path.
	Box [4]float32 `json:"box"`
}

// Mark returns the payload of the BARCODE mark of the symbol drawn in the
// box at (x, y) of the size w × h, as Path draws it.
func (s *Symbol) Mark(x, y, w, h float32) string {
	m := Mark{Kind: s.Kind, Text: s.Text, Level: s.opts.Level, Quiet: s.opts.Quiet, NoText: s.noText, Box: [4]float32{x, y, w, h}}
	if s.Kind == QR && m.Level == "" {
		m.Level = "M"
	}
	if m.Quiet < 0 {
		m.Quiet = NoQuiet
	}
	b, _ := json.Marshal(m)
	return string(b)
}

// Fill draws the symbol in the box at (x, y) of the size w × h on an
// object, in its fill colour, after the BARCODE mark that says what it
// is.
func (s *Symbol) Fill(o *bdf.Object, x, y, w, h float32) {
	o.Mark(bdf.MarkBarcode, s.Mark(x, y, w, h))
	o.FillPath(o.AddPath(s.Path(x, y, w, h)), bdf.NonZero)
}

// ParseMark reads the payload of a BARCODE mark.
func ParseMark(payload string) (*Mark, error) {
	if len(payload) > MaxText*2+256 {
		return nil, fmt.Errorf("barcode: a mark of %d bytes", len(payload))
	}
	var m Mark
	if err := json.Unmarshal([]byte(payload), &m); err != nil {
		return nil, fmt.Errorf("barcode: the mark: %w", err)
	}
	if Kind(m.Kind) == "" {
		return nil, fmt.Errorf("barcode: the mark: %q is not a symbology", m.Kind)
	}
	m.Kind = Kind(m.Kind)
	for _, v := range m.Box {
		if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
			return nil, fmt.Errorf("barcode: the mark: a box of %v", m.Box)
		}
	}
	if m.Box[2] <= 0 || m.Box[3] <= 0 {
		return nil, fmt.Errorf("barcode: the mark: an empty box")
	}
	return &m, nil
}

// Symbol makes the symbol of a mark again.
func (m *Mark) Symbol() (*Symbol, error) {
	return Encode(m.Kind, m.Text, &Options{Level: m.Level, Quiet: m.Quiet, NoText: m.NoText})
}

// Bars returns the columns of the symbol's first and last dark modules
// (x1 is past the last): the symbol without its quiet zone.
func (s *Symbol) Bars() (x0, x1 int) {
	x0, x1 = s.W, 0
	for y := 0; y < s.H; y++ {
		for x := 0; x < s.W; x++ {
			if s.mods[y*s.W+x] {
				x0, x1 = min(x0, x), max(x1, x+1)
			}
		}
	}
	if x1 <= x0 {
		return 0, 0
	}
	return x0, x1
}

// HasText reports whether Path draws digits under the bars.
func (s *Symbol) HasText() bool {
	return !s.noText && (s.Kind == EAN13 || s.Kind == EAN8 || s.Kind == UPCA) && len(s.Text) > 0
}

// BarHeight returns the height of the bars of the symbol drawn in a box
// w × h by Path: h less the room of the digits.
func (s *Symbol) BarHeight(w, h float32) float32 {
	if !s.HasText() || w <= 0 || h <= 0 {
		return h
	}
	digitScale := min(h*0.22/750, w/float32(len(s.Text)*620))
	return h - float32(float32(digitScale*750)*1.25)
}
