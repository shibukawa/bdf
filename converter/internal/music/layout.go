package music

import (
	"math"
	"slices"
	"sort"
	"strconv"
)

// The engraver lays a score out in staff spaces (sp): x from the left end
// of the staff lines of a system, y down from the top line of each staff
// (the bottom line is at 4). Staff positions count half spaces down from
// the top line (the middle line is 4, the bottom line 8).

// staffRef is one staff of the score, top to bottom.
type staffRef struct {
	part, staff int
}

// staffMeasure is what one staff has in one measure, resolved.
type staffMeasure struct {
	clef        Clef // at the start of the measure
	clefChange  bool // the clef changes at the start of the measure
	clefChanges []ClefChange
	key         KeySig
	keyChange   bool // a key signature starts here (drawn unless the system starts here)
	prevKey     KeySig
	events      []*Event
	acc         map[*Note]Accidental // the accidentals drawn
	voices      int                  // the number of voices with notes drawn
	directions  []*Direction
}

// slot is an onset of a measure: where the events that start together at
// an offset stand.
type slot struct {
	offset int
	dur    int     // to the next onset or the end of the measure
	pre    float64 // room before the noteheads (accidentals, grace notes, clef changes)
	head   float64 // the width of the noteheads
	post   float64 // room after the noteheads (dots, flags, noteheads beside the stem)
	lyric  float64 // half the width of the widest syllable centred on it
	ideal  float64 // the ideal distance to the next onset
	min    float64 // the least distance to the next onset
	x      float64 // after justification: from the start of the measure
}

// column is the horizontal layout of a measure, shared by every staff.
type column struct {
	index     int
	slots     []*slot
	header    float64 // key and time signature changes at the start
	timeSig   *TimeSig
	lead      float64 // room between the bar line (or the header) and the first slot
	width     float64 // after justification
	x         float64 // after justification: from the start of the system
	bodyStart float64 // where the first slot's space starts, from the start of the measure
}

// system is a line of measures.
type system struct {
	first, last int     // measures
	header      float64 // clef, key and time signature at the start
	width       float64
	staves      []*staffLayout
	y           []float64 // top line of each staff from the top line of the first
	above       float64   // extent above the first staff (negative)
	below       float64   // extent below the last staff's top line
	page        int
	top         float64 // top line of the first staff on the page (pt)
	left        float64 // room left of the staves for names and brackets
}

// staffLayout is the drawing of one staff of a system, in staff-local
// coordinates.
type staffLayout struct {
	prims        []prim
	above, below float64 // extents (above is negative)
	lyricsY      []float64
}

// engraver lays out and draws one score.
type engraver struct {
	s   *Score
	o   *Options
	sp  float64 // pt per staff space
	def engravingDefaults

	staves []staffRef
	sm     [][]*staffMeasure // [measure][staff]
	cols   []*column
	starts []int // start tick of each measure

	text  *texter
	warn  func(string)
	notes int
}

// engravingDefaults are the thicknesses the engraver draws with, in staff
// spaces.
type engravingDefaults struct {
	staffLine, stem, leger, legerExt, beam, beamSpacing, thinBar, thickBar, barSep  float64
	tieEnd, tieMid, slurEnd, slurMid, tupletBracket, hairpin, octaveLine, pedalLine float64
}

// prepare resolves the clefs, keys, events and accidentals of every staff
// in every measure.
func (e *engraver) prepare() {
	for pi, p := range e.s.Parts {
		for st := 0; st < max(p.Staves, 1); st++ {
			e.staves = append(e.staves, staffRef{pi, st})
		}
	}
	e.starts = make([]int, len(e.s.Measures)+1)
	for i, m := range e.s.Measures {
		e.starts[i+1] = e.starts[i] + m.Length
	}
	e.sm = make([][]*staffMeasure, len(e.s.Measures))
	clefs := make([]Clef, len(e.staves))
	for i := range clefs {
		clefs[i] = Clef{Sign: ClefG}
		if e.staves[i].staff == 1 {
			clefs[i] = Clef{Sign: ClefF}
		}
		if p := e.s.Parts[e.staves[i].part]; p.Percussion {
			clefs[i] = Clef{Sign: ClefPercussion}
		}
	}
	endClefs := slices.Clone(clefs)
	keys := make([]KeySig, len(e.s.Parts))
	for mi := range e.s.Measures {
		row := make([]*staffMeasure, len(e.staves))
		prevKeys := slices.Clone(keys)
		keyChanged := make([]bool, len(e.s.Parts))
		for pi, p := range e.s.Parts {
			if mi < len(p.Measures) && p.Measures[mi] != nil && p.Measures[mi].Key != nil {
				keys[pi] = *p.Measures[mi].Key
				keyChanged[pi] = mi > 0 && keys[pi] != prevKeys[pi]
			}
		}
		for si, ref := range e.staves {
			p := e.s.Parts[ref.part]
			sm := &staffMeasure{acc: map[*Note]Accidental{}, key: keys[ref.part], prevKey: prevKeys[ref.part],
				keyChange: keyChanged[ref.part]}
			var pm *PartMeasure
			if mi < len(p.Measures) {
				pm = p.Measures[mi]
			}
			if pm != nil {
				for _, c := range pm.Clefs {
					if c.Staff != ref.staff {
						continue
					}
					if c.Offset <= 0 {
						clefs[si] = c.Clef
					} else {
						sm.clefChanges = append(sm.clefChanges, c)
					}
				}
				for _, ev := range pm.Events {
					if ev.Staff == ref.staff || ev.Staff >= max(p.Staves, 1) && ref.staff == 0 {
						sm.events = append(sm.events, ev)
					}
				}
				for _, d := range pm.Directions {
					if d.Staff == ref.staff || d.Staff >= max(p.Staves, 1) && ref.staff == 0 {
						sm.directions = append(sm.directions, d)
					}
				}
			}
			sm.clef = clefs[si]
			if mi > 0 {
				sm.clefChange = sm.clef != endClefs[si]
			}
			if p.Percussion {
				sm.key, sm.prevKey = KeySig{}, KeySig{}
				sm.keyChange = false
			}
			for _, c := range sm.clefChanges {
				clefs[si] = c.Clef
			}
			endClefs[si] = clefs[si]
			slices.SortStableFunc(sm.events, func(a, b *Event) int {
				if a.Offset != b.Offset {
					return a.Offset - b.Offset
				}
				if a.Grace != b.Grace {
					if a.Grace {
						return -1
					}
					return 1
				}
				return a.Voice - b.Voice
			})
			voices := map[int]bool{}
			for _, ev := range sm.events {
				if !ev.Hidden && !(ev.Rest && ev.MeasureRest && len(voices) > 0) {
					voices[ev.Voice] = true
				}
			}
			sm.voices = len(voices)
			e.accidentals(sm)
			row[si] = sm
		}
		e.sm[mi] = row
	}
}

// accidentals works out the accidentals drawn in a measure of a staff:
// those the score gives, and for AccAuto the ones the key signature and
// the notes before in the measure call for.
func (e *engraver) accidentals(sm *staffMeasure) {
	type pk struct{ step, oct int }
	state := map[pk]int{}
	keyAlter := keyAlters(sm.key)
	for _, ev := range sm.events {
		for _, n := range ev.Notes {
			k := pk{n.Pitch.Step, n.Pitch.Octave}
			cur, ok := state[k]
			if !ok {
				cur = keyAlter[(n.Pitch.Step%7+7)%7]
			}
			switch n.Accidental {
			case AccAuto:
				if n.TieStop {
					continue // a tied note keeps its accidental without drawing it
				}
				if n.Pitch.Alter != cur || n.Cautionary {
					sm.acc[n] = accidentalFor(n.Pitch.Alter)
				}
				state[k] = n.Pitch.Alter
			case AccNone:
				state[k] = n.Pitch.Alter
			default:
				sm.acc[n] = n.Accidental
				state[k] = n.Pitch.Alter
			}
		}
	}
}

func accidentalFor(alter int) Accidental {
	switch alter {
	case 1:
		return AccSharp
	case -1:
		return AccFlat
	case 2:
		return AccDoubleSharp
	case -2:
		return AccDoubleFlat
	}
	return AccNatural
}

// sharpOrder and flatOrder are the steps of a key signature's sharps and
// flats in order.
var (
	sharpOrder = [7]int{3, 0, 4, 1, 5, 2, 6} // F C G D A E B
	flatOrder  = [7]int{6, 2, 5, 1, 4, 0, 3} // B E A D G C F
)

// keyAlters returns the alteration the key signature gives each step.
func keyAlters(k KeySig) [7]int {
	var a [7]int
	for i := 0; i < min(k.Fifths, 7); i++ {
		a[sharpOrder[i]] = 1
	}
	for i := 0; i < min(-k.Fifths, 7); i++ {
		a[flatOrder[i]] = -1
	}
	return a
}

// clefRef returns the staff position and the diatonic index of the pitch
// the clef names (G4 on the second line for a treble clef).
func clefRef(c Clef) (pos, diatonic int) {
	switch c.Sign {
	case ClefF:
		line := c.Line
		if line == 0 {
			line = 4
		}
		return (5 - line) * 2, 3*7 + 3 + 7*c.Octave // F3
	case ClefC:
		line := c.Line
		if line == 0 {
			line = 3
		}
		return (5 - line) * 2, 4*7 + 7*c.Octave // C4
	default:
		line := c.Line
		if line == 0 || c.Sign != ClefG {
			line = 2
		}
		return (5 - line) * 2, 4*7 + 4 + 7*c.Octave // G4
	}
}

// staffPos returns the staff position of a pitch under a clef.
func staffPos(c Clef, p Pitch) int {
	pos, d := clefRef(c)
	return pos - (p.Octave*7 + p.Step - d)
}

// keyPositions returns the staff positions of a key signature's
// accidentals under a clef.
func keyPositions(c Clef, k KeySig) []int {
	// the treble clef's positions, moved with the clef
	sharps := []int{0, 3, -1, 2, 5, 1, 4}
	flats := []int{4, 1, 5, 2, 6, 3, 7}
	shift := 0
	switch c.Sign {
	case ClefF:
		p, _ := clefRef(c)
		shift = 2 + (p - 2)
	case ClefC:
		p, _ := clefRef(c)
		shift = 1 + (p - 4)
	case ClefG:
		p, _ := clefRef(c)
		shift = p - 6
	}
	var out []int
	n := k.Fifths
	src := sharps
	if n < 0 {
		n, src = -n, flats
	}
	for i := 0; i < min(n, 7); i++ {
		v := src[i] + shift
		for v > 8 {
			v -= 7
		}
		for v < -2 {
			v += 7
		}
		out = append(out, v)
	}
	return out
}

// keyWidth is the width of a key signature.
func keyWidth(k KeySig) float64 {
	n := k.Fifths
	g := sym("accidentalSharp")
	if n < 0 {
		n, g = -n, sym("accidentalFlat")
	}
	if n == 0 {
		return 0
	}
	return float64(n)*(g.advance+0.1) + 0.6
}

// cancelWidth is the width of the naturals that cancel a key signature.
func cancelWidth(prev, next KeySig) float64 {
	n := cancelled(prev, next)
	if n == 0 {
		return 0
	}
	return float64(n)*(sym("accidentalNatural").advance+0.1) + 0.4
}

// cancelled is the number of the previous key's accidentals the next key
// cancels with naturals: all of them when the new key is C or turns from
// sharps to flats.
func cancelled(prev, next KeySig) int {
	if prev.Fifths == 0 {
		return 0
	}
	if next.Fifths == 0 || prev.Fifths > 0 != (next.Fifths > 0) {
		return abs(prev.Fifths)
	}
	return 0
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// timeWidth is the width of a time signature.
func timeWidth(t *TimeSig) float64 {
	if t == nil {
		return 0
	}
	if t.Symbol == "common" || t.Symbol == "cut" {
		return sym("timeSigCommon").advance + 0.8
	}
	w := 0.0
	for _, s := range []string{strconv.Itoa(t.Beats), strconv.Itoa(t.BeatType)} {
		ww := 0.0
		for _, c := range s {
			ww += sym("timeSig" + string(c)).advance
		}
		w = max(w, ww)
	}
	return w + 0.8
}

// clefGlyph returns the glyph of a clef (the smaller form for a change).
func clefGlyph(c Clef, change bool) *glyph {
	name := "gClef"
	switch c.Sign {
	case ClefF:
		name = "fClef"
		if c.Octave == -1 {
			name = "fClef8vb"
		}
	case ClefC:
		name = "cClef"
	case ClefPercussion:
		name = "unpitchedPercussionClef1"
	case ClefTab:
		name = "6stringTabClef"
	case ClefNone:
		return nil
	default:
		switch c.Octave {
		case -1:
			name = "gClef8vb"
		case 1:
			name = "gClef8va"
		}
	}
	if change {
		switch name {
		case "gClef":
			name = "gClefChange"
		case "fClef":
			name = "fClefChange"
		case "cClef":
			name = "cClefChange"
		}
	}
	return sym(name)
}

// noteheadGlyph returns the notehead of an event.
func noteheadGlyph(ev *Event, n *Note) *glyph {
	if n != nil && n.Head == HeadX {
		switch {
		case ev.Type <= Whole:
			return sym("noteheadXWhole")
		case ev.Type == Half:
			return sym("noteheadXHalf")
		}
		return sym("noteheadXBlack")
	}
	switch ev.Type {
	case Breve:
		return sym("noteheadDoubleWhole")
	case Whole:
		return sym("noteheadWhole")
	case Half:
		return sym("noteheadHalf")
	}
	return sym("noteheadBlack")
}

// restGlyph returns the rest of a note value.
func restGlyph(t NoteType) *glyph {
	names := []string{"restDoubleWhole", "restWhole", "restHalf", "restQuarter", "rest8th", "rest16th",
		"rest32nd", "rest64th", "rest128th"}
	return sym(names[min(max(int(t), 0), len(names)-1)])
}

// graceScale is the size of grace and cue notes.
const graceScale = 0.6

// columns works out the slots of every measure and their ideal and least
// distances.
func (e *engraver) columns() {
	e.cols = make([]*column, len(e.s.Measures))
	for mi, m := range e.s.Measures {
		col := &column{index: mi, timeSig: m.Time}
		// clef, key and time signature changes at the start of the measure
		for _, sm := range e.sm[mi] {
			w := 0.0
			if sm.clefChange {
				if g := clefGlyph(sm.clef, true); g != nil {
					w += g.advance + 0.6
				}
			}
			if sm.keyChange {
				w += cancelWidth(sm.prevKey, sm.key) + keyWidth(sm.key)
			}
			col.header = max(col.header, w)
		}
		if mi > 0 && m.Time != nil {
			col.header += timeWidth(m.Time)
		}
		// the onsets of every staff
		offs := map[int]*slot{}
		for _, row := range e.sm[mi] {
			for _, ev := range row.events {
				if ev.Grace {
					continue
				}
				if offs[ev.Offset] == nil {
					offs[ev.Offset] = &slot{offset: ev.Offset}
				}
			}
		}
		if len(offs) == 0 {
			offs[0] = &slot{offset: 0}
		}
		for _, sl := range offs {
			col.slots = append(col.slots, sl)
		}
		sort.Slice(col.slots, func(i, j int) bool { return col.slots[i].offset < col.slots[j].offset })
		length := max(m.Length, 1)
		for i, sl := range col.slots {
			next := length
			if i+1 < len(col.slots) {
				next = col.slots[i+1].offset
			}
			sl.dur = max(next-sl.offset, 1)
		}
		// what stands at each slot
		for _, sm := range e.sm[mi] {
			for _, ev := range sm.events {
				sl := offs[ev.Offset]
				if sl == nil {
					continue
				}
				pre, head, post := e.eventWidths(sm, ev)
				if ev.Grace {
					sl.pre = max(sl.pre, pre+head+post+0.3)
					continue
				}
				sl.pre = max(sl.pre, pre)
				sl.head = max(sl.head, head)
				sl.post = max(sl.post, post)
				for _, l := range ev.Lyrics {
					sl.lyric = max(sl.lyric, e.text.width(l.Text, e.lyricSize(), false, false)/e.sp/2)
				}
			}
			for _, c := range sm.clefChanges {
				if sl := offs[c.Offset]; sl != nil {
					sl.pre += clefGlyph(c.Clef, true).advance + 0.5
				}
			}
		}
		// the shortest onset spacing decides the ideal distances
		shortest := PPQ / 2
		for _, sl := range col.slots {
			shortest = min(shortest, sl.dur)
		}
		for i, sl := range col.slots {
			d := float64(sl.dur)
			sl.ideal = 2.6 + 1.6*math.Log2(d/float64(shortest))
			if sl.head == 0 {
				sl.head = 1.18
			}
			sl.min = sl.head + sl.post + 0.5
			if i+1 < len(col.slots) {
				nx := col.slots[i+1]
				sl.min += nx.pre
				if sl.lyric > 0 || nx.lyric > 0 {
					sl.min = max(sl.min, sl.lyric+nx.lyric+0.9)
				}
			} else {
				sl.min += 0.6 // before the bar line
				if sl.lyric > 0 {
					sl.min = max(sl.min, sl.lyric+0.4)
				}
			}
			sl.ideal = max(sl.ideal, sl.min)
		}
		col.lead = 1.0 + col.slots[0].pre
		if col.slots[0].lyric > col.lead+col.slots[0].head/2 {
			col.lead = col.slots[0].lyric - col.slots[0].head/2
		}
		// a whole-measure rest alone gets room of its own
		e.cols[mi] = col
	}
}

// eventWidths returns the room an event needs before, at and after its
// noteheads.
func (e *engraver) eventWidths(sm *staffMeasure, ev *Event) (pre, head, post float64) {
	scale := 1.0
	if ev.Grace || ev.Cue {
		scale = graceScale
	}
	if ev.Rest {
		g := restGlyph(ev.Type)
		if ev.MeasureRest {
			g = restGlyph(Whole)
		}
		return 0, g.advance * scale, float64(ev.Dots) * 0.5 * scale
	}
	if len(ev.Notes) == 0 {
		return 0, 0, 0
	}
	g := noteheadGlyph(ev, ev.Notes[0])
	head = g.advance * scale
	// accidentals
	accs := 0.0
	for _, n := range ev.Notes {
		if a, ok := sm.acc[n]; ok && a != AccNone {
			accs = max(accs, accidentalGlyph(a).advance)
		}
	}
	if accs > 0 {
		// stacked accidentals of a chord take more columns
		cols := e.accidentalColumns(sm, ev)
		pre = (accs + 0.2) * float64(max(cols, 1)) * scale
	}
	// seconds in a chord put noteheads on both sides of the stem
	if hasSecond(sm.clef, ev) {
		head *= 2
	}
	if ev.Dots > 0 {
		post += (0.3 + 0.5*float64(ev.Dots)) * scale
	}
	if ev.Type >= Eighth && ev.Beams == nil && len(ev.Notes) > 0 {
		// a flag on an up stem may reach into the next slot
		post = max(post, 0.8*scale)
	}
	return pre, head, post
}

// hasSecond reports whether a chord has two notes a second apart.
func hasSecond(c Clef, ev *Event) bool {
	if len(ev.Notes) < 2 {
		return false
	}
	pos := make([]int, len(ev.Notes))
	for i, n := range ev.Notes {
		pos[i] = staffPos(c, n.Pitch)
	}
	slices.Sort(pos)
	for i := 1; i < len(pos); i++ {
		if pos[i]-pos[i-1] == 1 {
			return true
		}
	}
	return false
}

// accidentalGlyph returns the glyph of an accidental.
func accidentalGlyph(a Accidental) *glyph {
	switch a {
	case AccSharp:
		return sym("accidentalSharp")
	case AccFlat:
		return sym("accidentalFlat")
	case AccDoubleSharp:
		return sym("accidentalDoubleSharp")
	case AccDoubleFlat:
		return sym("accidentalDoubleFlat")
	}
	return sym("accidentalNatural")
}

// accidentalColumns counts the columns the accidentals of a chord take.
func (e *engraver) accidentalColumns(sm *staffMeasure, ev *Event) int {
	cols := placeAccidentals(sm, ev)
	n := 0
	for _, c := range cols {
		n = max(n, c.col+1)
	}
	return n
}

// placedAccidental is an accidental of a chord in its column (0 is the
// nearest the noteheads).
type placedAccidental struct {
	note *Note
	acc  Accidental
	pos  int
	col  int
}

// placeAccidentals puts the accidentals of a chord in columns, top to
// bottom, each in the first column where it clears the others by three
// spaces.
func placeAccidentals(sm *staffMeasure, ev *Event) []placedAccidental {
	var out []placedAccidental
	for _, n := range ev.Notes {
		if a, ok := sm.acc[n]; ok && a != AccNone {
			out = append(out, placedAccidental{note: n, acc: a, pos: staffPos(sm.clef, n.Pitch)})
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].pos < out[j].pos })
	for i := range out {
		for c := 0; ; c++ {
			ok := true
			for j := 0; j < i; j++ {
				if out[j].col == c && out[i].pos-out[j].pos < 6 {
					ok = false
					break
				}
			}
			if ok {
				out[i].col = c
				break
			}
		}
	}
	return out
}

// systemHeader is the width of the clef, key and time signature at the
// start of a system that starts with measure mi.
func (e *engraver) systemHeader(mi int, first bool) float64 {
	w := 0.0
	for _, sm := range e.sm[mi] {
		g := clefGlyph(sm.clef, false)
		cw := 0.8
		if g != nil {
			cw += g.advance + 0.8
		}
		w = max(w, cw+keyWidth(sm.key))
	}
	if t := e.timeAt(mi); first || e.s.Measures[mi].Time != nil {
		w += timeWidth(t)
	}
	return w
}

// timeAt returns the time signature in effect in a measure.
func (e *engraver) timeAt(mi int) *TimeSig {
	for i := mi; i >= 0; i-- {
		if t := e.s.Measures[i].Time; t != nil {
			return t
		}
	}
	return &TimeSig{Beats: 4, BeatType: 4}
}

// measureWidth returns the fixed and stretchable widths of a measure in a
// system (the header of a measure that starts a system is the system's).
func (e *engraver) measureWidth(mi int, startsSystem bool) (fixed, ideal float64) {
	col := e.cols[mi]
	fixed = col.lead
	if !startsSystem {
		// at the start of a system the key or time signature is in its header
		fixed += col.header
	}
	for _, sl := range col.slots {
		ideal += sl.ideal
	}
	fixed += barlineWidth(e.rightBar(mi)) + barlineWidth(e.leftBar(mi))
	return fixed, ideal
}

// rightBar is the bar line drawn at the end of a measure: the score's,
// and a final bar line at the end of the music when the score has a plain
// one there.
func (e *engraver) rightBar(mi int) Barline {
	b := e.s.Measures[mi].Right
	if mi == len(e.s.Measures)-1 && b == BarRegular {
		return BarFinal
	}
	return b
}

// leftBar is the bar line drawn at the start of a measure (a start repeat).
func (e *engraver) leftBar(mi int) Barline {
	switch e.s.Measures[mi].Left {
	case BarRepeatStart, BarHeavyLight:
		return e.s.Measures[mi].Left
	}
	return BarNone
}

func barlineWidth(b Barline) float64 {
	switch b {
	case BarDouble:
		return 0.6
	case BarFinal, BarHeavyLight:
		return 0.9
	case BarRepeatStart, BarRepeatEnd:
		return 1.5
	case BarRepeatBoth:
		return 2.4
	}
	return 0.2
}

// breakSystems divides the measures into systems of width w (staff
// spaces) and justifies them.
func (e *engraver) breakSystems(width, firstIndent float64) []*system {
	var out []*system
	mi := 0
	for mi < len(e.s.Measures) {
		first := len(out) == 0
		sys := &system{first: mi}
		avail := width
		if first {
			avail -= firstIndent
		}
		used := e.systemHeader(mi, first)
		sys.header = used
		last := mi
		for k := mi; k < len(e.s.Measures); k++ {
			if k > mi && e.s.Measures[k].Break != BreakNone {
				break
			}
			fixed, ideal := e.measureWidth(k, k == mi)
			if k > mi && used+fixed+ideal > avail {
				break
			}
			used += fixed + ideal
			last = k
		}
		sys.last = last
		sys.width = avail
		out = append(out, sys)
		mi = last + 1
	}
	for i, sys := range out {
		e.justify(sys, i == len(out)-1)
	}
	return out
}

// justify stretches the measures of a system to its width: the ideal
// distances grow (or shrink down to the least ones) by one factor. The
// last system stays at its ideal width unless it is nearly full.
func (e *engraver) justify(sys *system, last bool) {
	fixed := sys.header
	var springs []*slot
	for mi := sys.first; mi <= sys.last; mi++ {
		f, _ := e.measureWidth(mi, mi == sys.first)
		fixed += f
		springs = append(springs, e.cols[mi].slots...)
	}
	total := func(f float64) float64 {
		t := 0.0
		for _, sl := range springs {
			t += max(sl.min, sl.ideal*f)
		}
		return t
	}
	target := sys.width - fixed
	f := 1.0
	natural := total(1)
	if !last || natural > target*0.8 {
		lo, hi := 0.0, 1.0
		for total(hi) < target && hi < 64 {
			hi *= 2
		}
		for i := 0; i < 50; i++ {
			mid := (lo + hi) / 2
			if total(mid) < target {
				lo = mid
			} else {
				hi = mid
			}
		}
		f = hi
	}
	if last && natural <= target*0.8 {
		sys.width = fixed + natural
	}
	x := sys.header
	for mi := sys.first; mi <= sys.last; mi++ {
		col := e.cols[mi]
		col.x = x
		fx := barlineWidth(e.leftBar(mi))
		if mi != sys.first {
			fx += col.header
		}
		col.bodyStart = fx
		sx := fx + col.lead
		for _, sl := range col.slots {
			sl.x = sx
			sx += max(sl.min, sl.ideal*f)
		}
		col.width = sx + barlineWidth(e.rightBar(mi))
		x += col.width
	}
	sys.width = x
}
