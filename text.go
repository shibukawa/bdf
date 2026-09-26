package bdf

import "math"

// MARK kinds (docs/spec.md §7.8).
const (
	MarkParagraph byte = 0
	MarkLine      byte = 1
	MarkCell      byte = 2
	MarkBox       byte = 3
	MarkAltText   byte = 4
	MarkWrap      byte = 5
	MarkHeading   byte = 6  // payload: level "1"-"6"
	MarkList      byte = 7  // opens a list, closed by MarkEnd
	MarkListItem  byte = 8  // starts an item of the innermost list
	MarkTable     byte = 9  // opens a table, closed by MarkEnd
	MarkFigure    byte = 10 // payload: alternative text; opens a figure, closed by MarkEnd
	MarkEnd       byte = 11 // closes the innermost list, table or figure
	MarkLang      byte = 12 // payload: BCP 47 language of the runs that follow ("" = document default)
)

// Separators between consecutive text runs (docs/spec.md §7.9).
const (
	SepNone  byte = 0
	SepSpace byte = 1
	SepBreak byte = 2
)

// TextRun is a text drawing instruction (or ALT_TEXT mark) located in the
// coordinate space of the object being walked.
type TextRun struct {
	Text    string
	X, Y    float32 // anchor after applying the transform
	Advance float32
	Size    float32
	Font    *Font
	Align   byte
	Matrix  [6]float32
	Sep     byte // separator from the previous run, from MARKs or position heuristics
	Ordinal int  // position in the extraction order of the top-level object
	AltText bool
}

type matrix [6]float32

func mul(m, n matrix) matrix {
	return matrix{
		m[0]*n[0] + m[2]*n[1],
		m[1]*n[0] + m[3]*n[1],
		m[0]*n[2] + m[2]*n[3],
		m[1]*n[2] + m[3]*n[3],
		m[0]*n[4] + m[2]*n[5] + m[4],
		m[1]*n[4] + m[3]*n[5] + m[5],
	}
}

type textState struct {
	m     matrix
	font  *Font
	size  float32
	align byte
}

type textExtractor struct {
	runs    []TextRun
	resolve func(Hash) *ObjectPart
	pending byte // separator requested by a MARK for the next run
	hasMark bool
	alt     *string // ALT_TEXT waiting for the drawing op it describes
}

func (t *textExtractor) walk(o *ObjectPart, m matrix) error {
	st := textState{m: m, size: 10}
	var stack []textState
	return o.Walk(func(in Instr) {
		switch in.Op {
		case OpSave:
			stack = append(stack, st)
		case OpRestore:
			if n := len(stack); n > 0 {
				st = stack[n-1]
				stack = stack[:n-1]
			}
		case OpTransform:
			st.m = mul(st.m, matrix{f(in, 0), f(in, 1), f(in, 2), f(in, 3), f(in, 4), f(in, 5)})
		case OpTranslate:
			st.m = mul(st.m, matrix{1, 0, 0, 1, f(in, 0), f(in, 1)})
		case OpScale:
			st.m = mul(st.m, matrix{f(in, 0), 0, 0, f(in, 1), 0, 0})
		case OpFont:
			if i := int(in.Args[0].(uint64)); i < len(o.Fonts) {
				st.font = &o.Fonts[i]
			}
			st.size = f(in, 1)
		case OpTextStyle:
			st.align = byte(in.Args[0].(uint64))
		case OpMark:
			kind := byte(in.Args[0].(uint64))
			payload := in.Args[1].(string)
			switch kind {
			case MarkLine:
				t.flushAlt(st)
				t.mark(SepSpace)
			case MarkParagraph, MarkCell, MarkBox, MarkHeading, MarkList, MarkListItem, MarkTable, MarkFigure, MarkEnd:
				t.flushAlt(st)
				t.mark(SepBreak)
			case MarkWrap:
				t.flushAlt(st)
				t.mark(SepNone)
			case MarkAltText:
				t.flushAlt(st)
				p := payload
				t.alt = &p
			}
			// MarkLang and unknown kinds leave the separator to the next MARK
			// or to the position guess.
		case OpFillText, OpStrokeText:
			if t.alt != nil {
				// The drawing op right after ALT_TEXT renders that text.
				t.emit(*t.alt, f(in, 1), f(in, 2), f(in, 3), st, true)
				t.alt = nil
				return
			}
			t.emit(in.Args[0].(string), f(in, 1), f(in, 2), f(in, 3), st, false)
		case OpFillPathAt:
			if t.alt != nil {
				t.emit(*t.alt, f(in, 2), f(in, 3), 0, st, true)
				t.alt = nil
			}
		case OpFillPathRun:
			if t.alt != nil {
				if gl := in.Args[1].([]Glyph); len(gl) > 0 {
					t.emit(*t.alt, gl[0].X, gl[0].Y, 0, st, true)
				} else {
					t.emit(*t.alt, 0, 0, 0, st, true)
				}
				t.alt = nil
			}
		case OpUse, OpUseAt:
			if t.alt != nil {
				x, y := float32(0), float32(0)
				if in.Op == OpUseAt {
					x, y = f(in, 1), f(in, 2)
				}
				t.emit(*t.alt, x, y, 0, st, true)
				t.alt = nil
				return
			}
			child := t.resolve(o.Objects[int(in.Args[0].(uint64))])
			if child == nil {
				return
			}
			cm := st.m
			if in.Op == OpUseAt {
				cm = mul(cm, matrix{1, 0, 0, 1, f(in, 1), f(in, 2)})
			}
			_ = t.walk(child, cm)
		}
	})
}

func f(in Instr, i int) float32 { return in.Args[i].(float32) }

// flushAlt emits an ALT_TEXT that no drawing op consumed, without a position.
func (t *textExtractor) flushAlt(st textState) {
	if t.alt != nil {
		t.emit(*t.alt, 0, 0, 0, st, true)
		t.alt = nil
	}
}

func (t *textExtractor) mark(sep byte) {
	if !t.hasMark || sep > t.pending {
		t.pending = sep
	}
	t.hasMark = true
}

func (t *textExtractor) emit(text string, x, y, advance float32, st textState, alt bool) {
	m := st.m
	run := TextRun{
		Text: text, Advance: advance, Size: st.size, Font: st.font, Align: st.align, Matrix: m,
		X: m[0]*x + m[2]*y + m[4], Y: m[1]*x + m[3]*y + m[5], Ordinal: len(t.runs), AltText: alt,
	}
	if t.hasMark {
		run.Sep = t.pending
	} else if len(t.runs) > 0 {
		run.Sep = guessSep(t.runs[len(t.runs)-1], run)
	}
	t.pending, t.hasMark = SepNone, false
	t.runs = append(t.runs, run)
}

// guessSep estimates the separator between runs that carry no MARK information.
func guessSep(prev, cur TextRun) byte {
	size := prev.Size
	if size <= 0 {
		size = 10
	}
	if math.Abs(float64(cur.Y-prev.Y)) > float64(size)*0.5 {
		return SepSpace
	}
	end := prev.X + prev.Advance
	if prev.Advance == 0 || cur.X-end > float64To32(0.2)*size {
		return SepSpace
	}
	return SepNone
}

func float64To32(v float64) float32 { return float32(v) }

// ExtractText returns the text runs of an object and the objects it USEs, in
// walk order. resolve must return decoded child objects (nil skips the child).
func ExtractText(o *ObjectPart, resolve func(Hash) *ObjectPart) ([]TextRun, error) {
	t := &textExtractor{resolve: resolve}
	err := t.walk(o, matrix{1, 0, 0, 1, 0, 0})
	t.flushAlt(textState{m: matrix{1, 0, 0, 1, 0, 0}, size: 10})
	return t.runs, err
}
