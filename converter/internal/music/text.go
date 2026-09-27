package music

import (
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/internal/fontdb"
)

// textFamily is the family the words of a score are set in: a serif, as
// in engraved music (Bravura's companion text font is Academico, a
// Century Schoolbook).
const textFamily = "Times New Roman"

// texter measures the words of a score (titles, lyrics, directions) with
// the fonts of the document.
type texter struct {
	set *fontset.Set
}

// textLine is a line of text split into runs of one face, measured.
type textLine struct {
	runs  []textRun
	width float64 // pt
	size  float64 // pt
	asc   float64 // pt above the baseline
	desc  float64 // pt below the baseline
}

// textRun is a run of one face within a line.
type textRun struct {
	fc   *fontset.Choice
	s    string
	cjk  bool
	x, w float64 // pt from the start of the line
}

// line measures s at size pt.
func (t *texter) line(s string, size float64, bold, italic bool) *textLine {
	l := &textLine{size: size, asc: 0.8 * size, desc: 0.25 * size}
	var cur *textRun
	x := 0.0
	for _, r := range s {
		fc := t.set.FaceFor(textFamily, "", bold, italic, r)
		adv := t.set.Advance(fc, r) * size
		cjk := fontdb.IsCJK(r)
		if cur == nil || cur.fc != fc || cur.cjk != cjk {
			l.runs = append(l.runs, textRun{fc: fc, x: x, cjk: cjk})
			cur = &l.runs[len(l.runs)-1]
			l.asc = max(l.asc, fc.Asc*size)
			l.desc = max(l.desc, fc.Desc*size)
		}
		cur.s += string(r)
		cur.w += adv
		x += adv
	}
	l.width = x
	return l
}

// width measures s at size pt.
func (t *texter) width(s string, size float64, bold, italic bool) float64 {
	return t.line(s, size, bold, italic).width
}

// lyricSize is the size of lyrics in pt.
func (e *engraver) lyricSize() float64 { return 2.2 * e.sp }

// wordSize is the size of directions (words, chord symbols) in pt.
func (e *engraver) wordSize() float64 { return 2.2 * e.sp }
