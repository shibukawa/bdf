package musicxml

import (
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/converter/internal/music"
)

// measureInfo is what the parts share about a measure, with what the
// performance needs of it.
type measureInfo struct {
	m        *music.Measure
	numbered bool
	time     *music.TimeSig // the first time signature read in the measure
	free     bool           // senza misura
	content  int            // the furthest position a part reaches, in ticks

	// the volta read from one part
	endingPart   *partState
	endingStart  *endingStart
	endStop      bool
	endOpen      bool
	ending       *bracket
	left, right  bool // forward and backward repeats
	tempo        []tempoPoint
	segno, coda  []string
	tocoda       []string
	fine, dacapo bool
	dalsegno     *string
	// alFine and alCoda are what the words of a jump say it goes on to
	alFine, alCoda bool
}

// endingStart is the start of a volta as the file has it.
type endingStart struct {
	number, text string
}

// bracket is a volta over one or more measures.
type bracket struct {
	numbers     []int // the passes it is played on
	first, last int   // measures
	groupMax    int   // the last pass of the voltas it is one of
	final       bool  // it holds the last pass of its group
	repeats     bool  // it has a backward repeat
}

// tempoPoint is a tempo change in a measure.
type tempoPoint struct {
	offset    int
	bpm       float64
	metronome bool // from a metronome mark, not a sound tempo
}

// partState is a part being read.
type partState struct {
	id        string
	trackName string
	part      *music.Part
	play      []*playMeasure // aligned with part.Measures
	count     int            // measures read
	seen      bool

	instruments []*instrument
	instByID    map[string]*instrument

	div        float64 // divisions per quarter note
	time       music.TimeSig
	key        *music.KeySig
	clefs      map[int]music.Clef
	transpose  int
	transposed bool
	voices     map[string]int
	wedges     map[int]music.DirectionKind
	shifts     []*octaveSpan
	openShifts map[[2]int]*octaveSpan
	dyn        []dynPoint
	unpitched  bool
}

// instrument is a MIDI instrument of a part; -1 is unknown.
type instrument struct {
	id                                       string
	channel, program, unpitched, volume, pan int
}

// octaveSpan is an octave line: the notes of a staff under it are
// written delta octaves from their pitch.
type octaveSpan struct {
	staff, delta     int
	startM, startOff int
	endM, endOff     int // endM -1: to the end
}

// dynPoint is a change of the playing dynamics of a part, or an accent
// (sf) on the notes at its position.
type dynPoint struct {
	measure, offset, seq int
	vel                  int // 0: unchanged
	accent               int
}

// playMeasure is what a part plays in a measure.
type playMeasure struct {
	notes  []*playNote
	lyrics []playLyric
}

// playNote is a note as the part plays it in a measure.
type playNote struct {
	offset, dur       int
	key               int // sounding MIDI note number
	tieStart, tieStop bool
	staff, voice      int
	grace             bool
	vel               int // from the note's dynamics attribute; 0: the part's
	ev                *playEvent
}

// playEvent is what the notes of an event share in playing.
type playEvent struct {
	shorten int // staccato: the duration is divided by it
	accent  int
}

// playLyric is a syllable sung from an offset.
type playLyric struct {
	offset, verse int
	text          string
}

func newPart(id, name string) *partState {
	return &partState{id: id, trackName: name, part: &music.Part{Staves: 1}, instByID: map[string]*instrument{},
		div: 1, time: music.TimeSig{Beats: 4, BeatType: 4}, clefs: map[int]music.Clef{}, voices: map[string]int{},
		wedges: map[int]music.DirectionKind{}, openShifts: map[[2]int]*octaveSpan{}}
}

// mstate is the state of reading a measure of a part.
type mstate struct {
	idx       int
	mi        *measureInfo
	pm        *music.PartMeasure
	play      *playMeasure
	pos, end  float64 // quarter notes
	last      *music.Event
	lastStart float64
	lastPlay  *playEvent
	tuplets   map[int]*tupletState
	seq       int
}

// tupletState is a tuplet group being read in a voice.
type tupletState struct {
	t              *music.Tuplet
	marker         bool   // started by a tuplet element
	number         string // of the tuplet element
	actual, normal int
	target, sum    int // notated ticks the group fills, and has
}

// measureAt returns what the parts share of a measure, adding measures
// up to it; nil past the budget.
func (r *reader) measureAt(idx int) *measureInfo {
	if idx >= maxMeasures {
		r.warn("more than %d measures; the rest are left out", maxMeasures)
		return nil
	}
	for len(r.ms) <= idx {
		m := &music.Measure{}
		r.ms = append(r.ms, &measureInfo{m: m})
		r.score.Measures = append(r.score.Measures, m)
	}
	return r.ms[idx]
}

// partMeasure returns the part's measure, adding measures up to it.
func (p *partState) measureAt(idx int) (*music.PartMeasure, *playMeasure) {
	for len(p.part.Measures) <= idx {
		p.part.Measures = append(p.part.Measures, &music.PartMeasure{})
		p.play = append(p.play, &playMeasure{})
	}
	return p.part.Measures[idx], p.play[idx]
}

// quarters converts a duration in divisions into quarter notes.
func (p *partState) quarters(s string) float64 {
	v, ok := parseFloat(s)
	if !ok {
		return 0
	}
	return min(max(v/p.div, -maxQuarters), maxQuarters)
}

// voice returns the number of a voice of the part: its own when it is a
// small number, else one not taken.
func (p *partState) voice(s string) int {
	if s == "" {
		s = "1"
	}
	if v, ok := p.voices[s]; ok {
		return v
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 1 || v > 99 {
		taken := map[int]bool{}
		for _, t := range p.voices {
			taken[t] = true
		}
		for v = 1; taken[v]; v++ {
		}
	}
	p.voices[s] = v
	return v
}

// staff returns the 0-based staff of an element with a <staff> child,
// growing the part to hold it.
func (p *partState) staff(n *node) int {
	s, ok := atoi(n.textOf("staff"))
	if !ok || s < 1 {
		return 0
	}
	s = min(s, maxStaves) - 1
	p.part.Staves = max(p.part.Staves, s+1)
	return s
}

// measure reads the measure of a part: mn carries the measure's
// attributes, cn its music (the same element in a partwise score).
func (r *reader) measure(p *partState, idx int, mn, cn *node) {
	mi := r.measureAt(idx)
	if mi == nil {
		return
	}
	r.gotMusic = true
	m := mi.m
	if !mi.numbered {
		mi.numbered = true
		if t, ok := mn.attrOK("text"); ok && t != "" {
			m.Number = t
		} else {
			m.Number = mn.attr("number")
		}
	}
	if mn.yes("implicit") {
		m.Implicit = true
	}
	pm, play := p.measureAt(idx)
	st := &mstate{idx: idx, mi: mi, pm: pm, play: play, tuplets: map[int]*tupletState{}}
	for _, k := range cn.kids {
		switch k.name {
		case "attributes":
			r.attributes(p, st, k)
		case "note":
			r.note(p, st, k)
		case "backup":
			st.pos = max(st.pos-p.quarters(k.textOf("duration")), 0)
		case "forward":
			st.pos = min(st.pos+max(p.quarters(k.textOf("duration")), 0), maxQuarters)
			st.end = max(st.end, st.pos)
		case "direction":
			r.direction(p, st, k)
		case "harmony":
			r.harmony(p, st, k)
		case "barline":
			r.barline(p, st, k)
		case "print":
			if idx > 0 {
				switch {
				case k.yes("new-page"):
					m.Break = music.BreakPage
				case k.yes("new-system") && m.Break == music.BreakNone:
					m.Break = music.BreakSystem
				}
			}
		case "sound":
			r.sound(p, st, k, st.pos)
		}
	}
	mi.content = max(mi.content, tickOf(st.end))
	play.notes = graceTiming(play.notes)
}

// attributes reads the attributes of a part at a position: divisions,
// key, time, staves, clefs and transposition.
func (r *reader) attributes(p *partState, st *mstate, n *node) {
	for _, k := range n.kids {
		switch k.name {
		case "divisions":
			if v, ok := parseFloat(k.trim()); ok && v > 0 {
				p.div = v
			}
		case "key":
			if num := k.attr("number"); num != "" && num != "1" || k.child("fifths") == nil {
				continue
			}
			f, _ := atoi(k.textOf("fifths"))
			ks := music.KeySig{Fifths: min(max(f, -7), 7), Minor: k.textOf("mode") == "minor"}
			if p.key == nil || *p.key != ks {
				p.key = &ks
				c := ks
				st.pm.Key = &c
			}
		case "time":
			if num := k.attr("number"); num != "" && num != "1" {
				continue
			}
			if k.child("senza-misura") != nil {
				if st.mi.time == nil {
					st.mi.free = true
				}
				continue
			}
			if ts, ok := timeSig(k); ok {
				p.time = ts
				if st.mi.time == nil {
					st.mi.time = &ts
					st.mi.free = false
				}
			}
		case "staves":
			if v, ok := atoi(k.trim()); ok && v >= 1 {
				p.part.Staves = max(p.part.Staves, min(v, maxStaves))
			}
		case "clef":
			if k.yes("additional") {
				continue
			}
			staff := 0
			if v, ok := atoi(k.attr("number")); ok && v >= 1 {
				staff = min(v, maxStaves) - 1
				p.part.Staves = max(p.part.Staves, staff+1)
			}
			c := clef(k)
			if c.Sign == music.ClefPercussion {
				p.part.Percussion = true
			}
			if cur, ok := p.clefs[staff]; ok && cur == c {
				continue
			}
			p.clefs[staff] = c
			off := tickOf(st.pos)
			cs := st.pm.Clefs[:0]
			for _, e := range st.pm.Clefs {
				if e.Offset != off || e.Staff != staff {
					cs = append(cs, e)
				}
			}
			st.pm.Clefs = append(cs, music.ClefChange{Offset: off, Staff: staff, Clef: c})
		case "transpose":
			if num := k.attr("number"); num != "" && num != "1" {
				continue
			}
			ch, _ := parseFloat(k.textOf("chromatic"))
			oc, _ := atoi(k.textOf("octave-change"))
			p.transpose = min(max(int(math.Round(ch))+12*oc, -48), 48)
			p.transposed = p.transposed || p.transpose != 0
		case "directive":
			if t := collapse(k.text); t != "" {
				above := true
				st.pm.Directions = append(st.pm.Directions, &music.Direction{Offset: tickOf(st.pos), Kind: music.DirWords, Text: t, Above: &above})
			}
		}
	}
}

// timeSig reads a time signature: beats may be sums ("3+2"), and several
// signatures (3/8+2/4) add up over their common beat type.
func timeSig(n *node) (music.TimeSig, bool) {
	var beats, types []int
	for _, k := range n.kids {
		switch k.name {
		case "beats":
			sum := 0
			for _, s := range strings.Split(k.trim(), "+") {
				v, ok := atoi(s)
				if !ok || v < 1 {
					return music.TimeSig{}, false
				}
				sum += min(v, 1000)
			}
			beats = append(beats, sum)
		case "beat-type":
			v, ok := atoi(k.trim())
			if !ok || v < 1 {
				return music.TimeSig{}, false
			}
			types = append(types, min(v, 1024))
		}
	}
	if len(beats) == 0 || len(beats) != len(types) {
		return music.TimeSig{}, false
	}
	ts := music.TimeSig{Beats: beats[0], BeatType: types[0]}
	for i := 1; i < len(beats); i++ {
		l := lcm(ts.BeatType, types[i])
		if l > 1024 {
			break
		}
		ts.Beats = ts.Beats*(l/ts.BeatType) + beats[i]*(l/types[i])
		ts.BeatType = l
	}
	ts.Beats = min(ts.Beats, 1000)
	switch n.attr("symbol") {
	case "common", "cut":
		ts.Symbol = n.attr("symbol")
	}
	return ts, true
}

func lcm(a, b int) int {
	x, y := a, b
	for y != 0 {
		x, y = y, x%y
	}
	return a / x * b
}

// clef reads a clef; the usual line of its sign is 0.
func clef(n *node) music.Clef {
	var c music.Clef
	line, _ := atoi(n.textOf("line"))
	usual := 0
	switch strings.ToUpper(n.textOf("sign")) {
	case "F":
		c.Sign, usual = music.ClefF, 4
	case "C":
		c.Sign, usual = music.ClefC, 3
	case "PERCUSSION":
		c.Sign, line = music.ClefPercussion, 0
	case "TAB":
		c.Sign, line = music.ClefTab, 0
	case "NONE":
		c.Sign, line = music.ClefNone, 0
	default: // G, jianpu
		c.Sign, usual = music.ClefG, 2
	}
	if line >= 1 && line <= 5 && line != usual {
		c.Line = line
	}
	if v, ok := atoi(n.textOf("clef-octave-change")); ok {
		c.Octave = min(max(v, -2), 2)
	}
	return c
}

// noteTypes are the note values of <type>.
var noteTypes = map[string]music.NoteType{
	"maxima": music.Breve, "long": music.Breve, "breve": music.Breve, "whole": music.Whole, "half": music.Half,
	"quarter": music.Quarter, "eighth": music.Eighth, "16th": music.N16th, "32nd": music.N32nd, "64th": music.N64th,
	"128th": music.N128th, "256th": music.N128th, "512th": music.N128th, "1024th": music.N128th,
}

// inferType returns the note value and dots nearest a notated length.
func inferType(ticks int) (music.NoteType, int) {
	if ticks <= 0 {
		return music.Quarter, 0
	}
	for t := music.Breve; t <= music.N128th; t++ {
		for dots := 0; dots <= 3; dots++ {
			if t.Ticks(dots) == ticks {
				return t, dots
			}
		}
	}
	t := music.Breve
	for t < music.N128th && t.Ticks(0) > ticks {
		t++
	}
	dots := 0
	for dots < 3 && t.Ticks(dots+1) <= ticks {
		dots++
	}
	return t, dots
}

var steps = map[string]int{"C": 0, "D": 1, "E": 2, "F": 3, "G": 4, "A": 5, "B": 6}

// note reads a note, rest or chord tone.
func (r *reader) note(p *partState, st *mstate, n *node) {
	if r.notes >= maxNotes {
		r.warn("more than %d notes; the rest are left out", maxNotes)
		return
	}
	r.notes++
	graceN := n.child("grace")
	grace := graceN != nil
	chord := n.child("chord") != nil
	restN := n.child("rest")
	pitchN, unpN := n.child("pitch"), n.child("unpitched")
	typeN := n.child("type")
	cue := n.child("cue") != nil || typeN.attr("size") == "cue"

	dur, hasDur := 0.0, false
	if d := n.child("duration"); d != nil && !grace {
		dur, hasDur = max(p.quarters(d.trim()), 0), true
	}
	actual, normal := 1, 1
	tm := n.child("time-modification")
	if tm != nil {
		a, _ := atoi(tm.textOf("actual-notes"))
		b, _ := atoi(tm.textOf("normal-notes"))
		if a >= 1 && b >= 1 && a <= 1000 && b <= 1000 {
			actual, normal = a, b
		}
	}
	typ, typed := noteTypes[typeN.trim()]
	dots := 0
	for _, k := range n.kids {
		if k.name == "dot" && dots < 4 {
			dots++
		}
	}
	if !typed {
		switch {
		case grace:
			typ = music.Eighth
		case hasDur:
			typ, dots = inferType(tickOf(dur) * actual / normal)
		default:
			typ = music.Quarter
		}
	}
	if !hasDur && !grace {
		dur = float64(typ.Ticks(dots)*normal) / float64(actual) / music.PPQ
	}
	voice := p.voice(n.textOf("voice"))
	staff := p.staff(n)

	// the tuplet elements of the note
	var starts, stops []*node
	for _, nt := range n.kids {
		if nt.name != "notations" {
			continue
		}
		for _, k := range nt.kids {
			if k.name == "tuplet" {
				switch k.attr("type") {
				case "start":
					starts = append(starts, k)
				case "stop":
					stops = append(stops, k)
				}
			}
		}
	}

	var ev *music.Event
	var pe *playEvent
	start := st.pos
	if chord && st.last != nil {
		start = st.lastStart
		if !st.last.Rest && restN == nil && st.last.Grace == grace {
			ev, pe = st.last, st.lastPlay
		}
	}
	if ev == nil {
		ev = &music.Event{Offset: tickOf(start), Staff: staff, Voice: voice, Type: typ, Dots: dots, Grace: grace, Cue: cue}
		if !grace {
			ev.Duration = tickOf(start+dur) - ev.Offset
		} else {
			ev.Slash = graceN.yes("slash")
		}
		ev.Hidden = n.attr("print-object") == "no"
		switch n.textOf("stem") {
		case "up":
			ev.Stem = music.StemUp
		case "down":
			ev.Stem = music.StemDown
		case "none":
			ev.Stem = music.StemNone
		}
		if restN != nil || pitchN == nil && unpN == nil {
			ev.Rest = true
			ev.MeasureRest = restN.yes("measure") || !typed && start == 0 && abs(ev.Duration-p.time.Length()) <= 2
		}
		st.pm.Events = append(st.pm.Events, ev)
		pe = &playEvent{shorten: 1}
		st.last, st.lastStart, st.lastPlay = ev, start, pe
		if !chord && !grace {
			st.pos = min(start+dur, maxQuarters)
			st.end = max(st.end, st.pos)
		}
		if !grace {
			st.tuplet(ev, voice, tm, actual, normal, starts)
		}
	}
	if !grace {
		st.tupletStop(voice, stops)
	}

	// beams: those of the first note of a chord that has any
	var beams []music.Beam
	for _, k := range n.kids {
		if k.name != "beam" {
			continue
		}
		r.hasBeams = true
		num := atoiDef(k.attr("number"), 1)
		if num < 1 || num > 6 {
			continue
		}
		for len(beams) < num {
			beams = append(beams, 0)
		}
		switch k.trim() {
		case "begin":
			beams[num-1] = music.BeamBegin
		case "continue":
			beams[num-1] = music.BeamContinue
		case "end":
			beams[num-1] = music.BeamEnd
		case "forward hook":
			beams[num-1] = music.BeamForwardHook
		case "backward hook":
			beams[num-1] = music.BeamBackwardHook
		}
	}
	if len(beams) > 0 && len(ev.Beams) == 0 {
		ev.Beams = beams
	} else if ev.Beams == nil {
		ev.Beams = []music.Beam{}
	}

	r.notations(p, st, n, ev, pe, start)
	r.lyrics(st, n, ev, start)
	if ev.Rest || pitchN == nil && unpN == nil {
		return
	}

	note := &music.Note{Accidental: music.AccNone}
	key := 0
	if pitchN != nil {
		step, ok := steps[strings.ToUpper(pitchN.textOf("step"))]
		if !ok {
			step = 0
		}
		alterF, _ := parseFloat(pitchN.textOf("alter"))
		alter := int(math.Round(min(max(alterF, -4), 4)))
		oct := atoiDef(pitchN.textOf("octave"), 4)
		oct = min(max(oct, -1), 9)
		note.Pitch = music.Pitch{Step: step, Alter: min(max(alter, -2), 2), Octave: oct}
		if alter != 0 {
			r.hasAltered = true
		}
		key = music.Pitch{Step: step, Alter: alter, Octave: oct}.MIDI() + p.transpose
	} else {
		step, ok := steps[strings.ToUpper(unpN.textOf("display-step"))]
		oct, ok2 := atoi(unpN.textOf("display-octave"))
		if !ok || !ok2 {
			step, oct = 6, 4 // the middle line
		}
		note.Pitch = music.Pitch{Step: step, Octave: min(max(oct, -1), 9)}
		key = p.unpitchedKey(n, note.Pitch)
		p.unpitched = true
	}
	if acc := n.child("accidental"); acc != nil {
		r.hasAccidentals = true
		note.Accidental = accidental(acc.trim(), note.Pitch.Alter)
		note.Cautionary = acc.yes("cautionary") || acc.yes("parentheses") || acc.yes("editorial") || acc.yes("bracket")
	}
	switch n.textOf("notehead") {
	case "x", "cross":
		note.Head = music.HeadX
	}
	// ties: <tie> is what sounds, <tied> what is drawn
	var soundStart, soundStop, drawStart, drawStop, hasTie bool
	for _, k := range n.kids {
		if k.name == "tie" {
			hasTie = true
			switch k.attr("type") {
			case "start":
				soundStart = true
			case "stop":
				soundStop = true
			}
		}
	}
	for _, nt := range n.kids {
		if nt.name != "notations" {
			continue
		}
		for _, k := range nt.kids {
			if k.name == "tied" {
				switch k.attr("type") {
				case "start":
					drawStart = true
				case "stop":
					drawStop = true
				case "continue":
					drawStart, drawStop = true, true
				}
			}
		}
	}
	note.TieStart, note.TieStop = soundStart || drawStart, soundStop || drawStop
	if !hasTie {
		soundStart, soundStop = drawStart, drawStop
	}
	ev.Notes = append(ev.Notes, note)

	if cue {
		return
	}
	nd := dur
	if !hasDur || chord && nd <= 0 {
		nd = float64(ev.Duration) / music.PPQ
	}
	pn := &playNote{offset: tickOf(start), key: key, tieStart: soundStart, tieStop: soundStop,
		staff: staff, voice: voice, grace: grace, ev: pe}
	if !grace {
		pn.dur = tickOf(start+nd) - pn.offset
	}
	if v, ok := parseFloat(n.attr("dynamics")); ok {
		pn.vel = min(max(int(math.Round(v*90/100)), 1), 127)
	}
	st.play.notes = append(st.play.notes, pn)
}

// tuplet puts an event in a tuplet group of its voice: the one its tuplet
// elements start or continue, or without them the group of the events
// before with the same ratio until it fills its time.
func (st *mstate) tuplet(ev *music.Event, voice int, tm *node, actual, normal int, starts []*node) {
	ts := st.tuplets[voice]
	if tm == nil || actual == normal {
		delete(st.tuplets, voice)
		return
	}
	if len(starts) > 0 && (ts == nil || !ts.marker || starts[0].attr("number") == ts.number) {
		s := starts[0]
		t := &music.Tuplet{Actual: actual, Normal: normal, ShowNumber: s.attr("show-number") != "none"}
		if a, ok := atoi(s.child("tuplet-actual").textOf("tuplet-number")); ok && a >= 1 {
			t.Actual = a
		}
		if b, ok := atoi(s.child("tuplet-normal").textOf("tuplet-number")); ok && b >= 1 {
			t.Normal = b
		}
		switch s.attr("bracket") {
		case "yes":
			b := true
			t.Bracket = &b
		case "no":
			b := false
			t.Bracket = &b
		}
		ts = &tupletState{t: t, marker: true, number: s.attr("number")}
		st.tuplets[voice] = ts
		ev.Tuplet = t
		return
	}
	if ts != nil && ts.marker {
		ev.Tuplet = ts.t
		return
	}
	if ts == nil || ts.actual != actual || ts.normal != normal || ts.sum >= ts.target {
		unit := ev.Type.Ticks(ev.Dots)
		if nt, ok := noteTypes[tm.textOf("normal-type")]; ok {
			dots := 0
			for _, k := range tm.kids {
				if k.name == "normal-dot" {
					dots++
				}
			}
			unit = nt.Ticks(min(dots, 3))
		}
		ts = &tupletState{t: &music.Tuplet{Actual: actual, Normal: normal, ShowNumber: true}, actual: actual, normal: normal, target: actual * unit}
		st.tuplets[voice] = ts
	}
	ev.Tuplet = ts.t
	ts.sum += ev.Type.Ticks(ev.Dots)
}

// tupletStop ends the tuplet group of a voice at its stop element.
func (st *mstate) tupletStop(voice int, stops []*node) {
	ts := st.tuplets[voice]
	if ts == nil || !ts.marker {
		return
	}
	for _, s := range stops {
		if s.attr("number") == ts.number || s.attr("number") == "" || ts.number == "" {
			delete(st.tuplets, voice)
			return
		}
	}
}

// accidental maps an accidental element to the signs the engine draws:
// microtonal and other kinds to the nearest.
func accidental(s string, alter int) music.Accidental {
	switch s {
	case "sharp":
		return music.AccSharp
	case "flat":
		return music.AccFlat
	case "natural":
		return music.AccNatural
	case "double-sharp", "sharp-sharp", "triple-sharp", "double-sharp-down", "double-sharp-up":
		return music.AccDoubleSharp
	case "flat-flat", "double-flat", "triple-flat", "double-flat-down", "double-flat-up":
		return music.AccDoubleFlat
	}
	switch {
	case strings.Contains(s, "sharp") || s == "sori":
		return music.AccSharp
	case strings.Contains(s, "flat") || s == "koron":
		return music.AccFlat
	case strings.Contains(s, "natural"):
		return music.AccNatural
	}
	switch {
	case alter >= 2:
		return music.AccDoubleSharp
	case alter == 1:
		return music.AccSharp
	case alter == -1:
		return music.AccFlat
	case alter <= -2:
		return music.AccDoubleFlat
	}
	return music.AccNatural
}

// articulations maps the articulations and ornaments to SMuFL names.
var articulations = map[string][]string{
	"staccato": {"articStaccato"}, "accent": {"articAccent"}, "tenuto": {"articTenuto"},
	"strong-accent": {"articMarcato"}, "staccatissimo": {"articStaccatissimo"}, "spiccato": {"articStaccatissimo"},
	"detached-legato": {"articTenuto", "articStaccato"}, "breath-mark": {"breathMarkComma"}, "caesura": {"caesura"},
	"trill-mark": {"ornamentTrill"}, "turn": {"ornamentTurn"}, "delayed-turn": {"ornamentTurn"},
	"mordent": {"ornamentMordent"}, "inverted-mordent": {"ornamentShortTrill"},
}

// notations reads the marks of a note onto its event: articulations,
// fermatas, ornaments, tremolos, slurs, and dynamics under the note.
func (r *reader) notations(p *partState, st *mstate, n *node, ev *music.Event, pe *playEvent, start float64) {
	add := func(names ...string) {
		for _, name := range names {
			found := false
			for _, a := range ev.Artics {
				if a == name {
					found = true
				}
			}
			if !found {
				ev.Artics = append(ev.Artics, name)
			}
		}
	}
	for _, nt := range n.kids {
		if nt.name != "notations" || nt.attr("print-object") == "no" {
			continue
		}
		for _, k := range nt.kids {
			switch k.name {
			case "articulations", "ornaments":
				for _, a := range k.kids {
					add(articulations[a.name]...)
					switch a.name {
					case "staccato", "spiccato", "detached-legato":
						pe.shorten = max(pe.shorten, 2)
					case "staccatissimo":
						pe.shorten = max(pe.shorten, 4)
					case "accent":
						pe.accent = max(pe.accent, 10)
					case "strong-accent":
						pe.accent = max(pe.accent, 20)
					case "tremolo":
						if typ := a.attr("type"); typ == "" || typ == "single" {
							ev.Tremolo = min(max(atoiDef(a.trim(), 3), 1), 8)
						}
					}
				}
			case "fermata":
				add("fermata")
			case "slur":
				mark := music.SlurMark{Number: atoiDef(k.attr("number"), 1)}
				switch k.attr("type") {
				case "start":
					mark.Start = true
				case "stop":
				default:
					continue
				}
				mark.Above = placement(k.attr("placement"))
				ev.Slurs = append(ev.Slurs, mark)
			case "dynamics":
				if t := dynamicsText(k); t != "" {
					off := tickOf(start)
					st.pm.Directions = append(st.pm.Directions, &music.Direction{Offset: off, Staff: ev.Staff, Kind: music.DirDynamic,
						Text: t, Above: placement(k.attr("placement"))})
					p.dynamic(st, off, t)
				}
			}
		}
	}
}

// lyrics reads the lyrics of a note onto its event (a verse once).
func (r *reader) lyrics(st *mstate, n *node, ev *music.Event, start float64) {
	for _, l := range n.kids {
		if l.name != "lyric" || l.attr("print-object") == "no" {
			continue
		}
		verse := r.verse(l.attr("number"))
		var text strings.Builder
		first, last := "", ""
		extend := false
		for _, k := range l.kids {
			switch k.name {
			case "syllabic":
				if first == "" {
					first = k.trim()
				}
				last = k.trim()
			case "text":
				text.WriteString(k.text)
			case "elision":
				if t := k.text; strings.TrimSpace(t) != "" {
					text.WriteString(t)
				} else {
					text.WriteString("‿")
				}
			case "extend":
				if t := k.attr("type"); t == "" || t == "start" {
					extend = true
				}
			}
		}
		t := strings.TrimSpace(text.String())
		if t == "" {
			continue
		}
		dup := false
		for _, s := range ev.Lyrics {
			if s.Verse == verse {
				dup = true
			}
		}
		if dup {
			continue
		}
		starts := first == "" || first == "single" || first == "begin"
		ends := last == "" || last == "single" || last == "end"
		syl := music.SyllableSingle
		switch {
		case starts && !ends:
			syl = music.SyllableBegin
		case !starts && ends:
			syl = music.SyllableEnd
		case !starts && !ends:
			syl = music.SyllableMiddle
		}
		ev.Lyrics = append(ev.Lyrics, &music.LyricSyllable{Verse: verse, Text: t, Syllabic: syl, Extend: extend})
		if !ev.Grace {
			sung := t
			if !ends {
				sung += "-"
			}
			st.play.lyrics = append(st.play.lyrics, playLyric{offset: tickOf(start), verse: verse, text: sung})
		}
	}
}

// verse returns the verse of a lyric number: its own when it is a small
// number, else the next one not taken, in the order they appear.
func (r *reader) verse(num string) int {
	if num == "" {
		num = "1"
	}
	if v, ok := r.verses[num]; ok {
		return v
	}
	v, err := strconv.Atoi(num)
	if err != nil || v < 1 || v > 99 {
		taken := map[int]bool{}
		for _, t := range r.verses {
			taken[t] = true
		}
		for v = 1; taken[v]; v++ {
		}
	}
	r.verses[num] = v
	return v
}

// placement returns the Above of a placement attribute.
func placement(s string) *bool {
	switch s {
	case "above":
		b := true
		return &b
	case "below":
		b := false
		return &b
	}
	return nil
}

// dynamicsText writes the marks of a dynamics element.
func dynamicsText(n *node) string {
	var b strings.Builder
	for _, k := range n.kids {
		if k.name == "other-dynamics" {
			b.WriteString(collapse(k.text))
		} else {
			b.WriteString(k.name)
		}
	}
	return b.String()
}

// dynamicLevels are the velocities of the dynamics marks.
var dynamicLevels = map[string]int{
	"pppppp": 3, "ppppp": 5, "pppp": 10, "ppp": 23, "pp": 36, "p": 49, "mp": 64,
	"mf": 80, "f": 96, "ff": 112, "fff": 120, "ffff": 127, "fffff": 127, "ffffff": 127,
}

// dynamic records the playing dynamics a mark sets at an offset: a level,
// or an accent on the notes there (sf, sfz, fz) with the level after it
// (fp, sfp).
func (p *partState) dynamic(st *mstate, off int, text string) {
	t := strings.ToLower(strings.TrimSpace(text))
	d := dynPoint{measure: st.idx, offset: off, seq: st.nextSeq()}
	if v, ok := dynamicLevels[t]; ok {
		d.vel = v
	} else {
		for _, pre := range []string{"sffz", "sfz", "sf", "rfz", "rf", "fz", "f"} {
			if rest, ok := strings.CutPrefix(t, pre); ok && (rest == "" || dynamicLevels[rest] > 0 && rest[0] == 'p') {
				d.accent = 20
				d.vel = dynamicLevels[rest]
				break
			}
		}
		if d.accent == 0 {
			return
		}
	}
	p.dyn = append(p.dyn, d)
}

func (st *mstate) nextSeq() int {
	st.seq++
	return st.idx<<20 | st.seq
}

// unpitchedKey returns the MIDI note of an unpitched note: that of its
// instrument, else a General MIDI drum for its place on the staff.
func (p *partState) unpitchedKey(n *node, pitch music.Pitch) int {
	if in := p.instByID[n.child("instrument").attr("id")]; in != nil && in.unpitched >= 0 {
		return in.unpitched
	}
	if len(p.instruments) == 1 && p.instruments[0].unpitched >= 0 {
		return p.instruments[0].unpitched
	}
	switch pitch.Step + 7*pitch.Octave {
	case 1 + 7*4: // D4 hi-hat pedal
		return 44
	case 2 + 7*4: // E4
		return 35
	case 3 + 7*4: // F4 bass drum
		return 36
	case 4 + 7*4: // G4
		return 41
	case 5 + 7*4: // A4 floor tom
		return 43
	case 6 + 7*4: // B4
		return 45
	case 0 + 7*5: // C5 snare
		return 38
	case 1 + 7*5: // D5
		return 47
	case 2 + 7*5: // E5
		return 50
	case 3 + 7*5: // F5 ride
		return 51
	case 4 + 7*5: // G5 hi-hat
		return 42
	case 5 + 7*5: // A5 crash
		return 49
	}
	return 38
}

// tempoWords are the starts of words that are tempo marks.
var tempoWords = []string{
	"allegr", "andant", "moderat", "adagi", "prest", "lento", "largo", "larghetto", "larghissimo", "grave", "vivace",
	"vivacissimo", "vivo", "maestoso", "tempo ", "a tempo", "tempo i", "tempo primo", "tempo di", "animato", "sostenuto",
	"lent", "modéré", "vif", "rasch", "langsam", "mäßig", "massig", "schnell", "lebhaft", "moderately", "slowly",
	"fast", "brightly", "lively", "freely", "swing", "ballad",
}

// isTempoWords reports whether words are a tempo mark.
func isTempoWords(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	for _, w := range tempoWords {
		if strings.HasPrefix(s, w) || s == strings.TrimSpace(w) {
			return true
		}
	}
	return false
}

// direction reads a direction: dynamics, words, hairpins, pedal marks,
// octave lines, segno and coda signs, rehearsal and tempo marks, and its
// sound.
func (r *reader) direction(p *partState, st *mstate, n *node) {
	pos := st.pos + p.quarters(n.textOf("offset"))
	off := max(tickOf(pos), 0)
	staff := p.staff(n)
	above := placement(n.attr("placement"))
	snd := n.child("sound")
	_, sndTempo := snd.attrOK("tempo")
	dir := func(d *music.Direction) {
		d.Offset, d.Staff = off, staff
		if d.Above == nil {
			d.Above = above
		}
		st.pm.Directions = append(st.pm.Directions, d)
	}
	type words struct {
		text   string
		italic bool
	}
	var ws []words
	var metro *music.TempoMark
	var segno, coda bool
	for _, dt := range n.kids {
		if dt.name != "direction-type" {
			continue
		}
		// the words of a direction type are runs of one text
		var run strings.Builder
		italic, hasWords := false, false
		for _, k := range dt.kids {
			if k.name == "words" {
				if !hasWords {
					italic, hasWords = k.attr("font-style") == "italic", true
				}
				run.WriteString(k.text)
			}
		}
		if t := collapse(run.String()); t != "" {
			ws = append(ws, words{t, italic})
		}
		for _, k := range dt.kids {
			switch k.name {
			case "dynamics":
				if t := dynamicsText(k); t != "" {
					dir(&music.Direction{Kind: music.DirDynamic, Text: t})
					p.dynamic(st, off, t)
				}
			case "wedge":
				num := atoiDef(k.attr("number"), 1)
				switch k.attr("type") {
				case "crescendo":
					p.wedges[num] = music.DirCrescendo
					dir(&music.Direction{Kind: music.DirCrescendo, Number: num})
				case "diminuendo":
					p.wedges[num] = music.DirDiminuendo
					dir(&music.Direction{Kind: music.DirDiminuendo, Number: num})
				case "stop":
					kind, ok := p.wedges[num]
					if !ok {
						kind = music.DirCrescendo
					}
					delete(p.wedges, num)
					dir(&music.Direction{Kind: kind, Stop: true, Number: num})
				}
			case "pedal":
				num := atoiDef(k.attr("number"), 1)
				switch k.attr("type") {
				case "start", "sostenuto", "resume":
					dir(&music.Direction{Kind: music.DirPedal, Number: num})
				case "stop", "discontinue":
					dir(&music.Direction{Kind: music.DirPedal, Stop: true, Number: num})
				case "change":
					dir(&music.Direction{Kind: music.DirPedal, Stop: true, Number: num})
					dir(&music.Direction{Kind: music.DirPedal, Number: num})
				}
			case "octave-shift":
				r.octaveShift(p, st, k, staff, off, dir)
			case "segno":
				segno = true
				dir(&music.Direction{Kind: music.DirSegno})
			case "coda":
				coda = true
				dir(&music.Direction{Kind: music.DirCoda})
			case "rehearsal":
				if t := collapse(k.text); t != "" && st.mi.m.Rehearsal == "" {
					st.mi.m.Rehearsal = t
				}
			case "metronome":
				metro = metronome(k)
			}
		}
	}

	var texts []string
	for _, w := range ws {
		texts = append(texts, w.text)
	}
	if len(ws) > 0 && (metro != nil || sndTempo || isTempoWords(ws[0].text)) || metro != nil {
		mark := music.TempoMark{Offset: off, Text: strings.Join(texts, " ")}
		if metro != nil {
			mark.Beat, mark.Dots, mark.PerMinute = metro.Beat, metro.Dots, metro.PerMinute
		}
		if mark.Text != "" || mark.PerMinute > 0 {
			st.mi.addTempoMark(mark)
		}
	} else if len(ws) > 0 {
		dir(&music.Direction{Kind: music.DirWords, Text: strings.Join(texts, " "), Italic: ws[0].italic})
	}

	jumps := false
	if snd != nil {
		jumps = r.sound(p, st, snd, pos)
	}
	for _, t := range texts {
		l := strings.ToLower(t)
		st.mi.alFine = st.mi.alFine || strings.Contains(l, "al fine")
		st.mi.alCoda = st.mi.alCoda || strings.Contains(l, "al coda")
	}
	if metro != nil && metro.PerMinute > 0 && !sndTempo {
		bpm := metro.PerMinute * float64(metro.Beat.Ticks(metro.Dots)) / music.PPQ
		if bpm > 0 && bpm < 10000 {
			st.mi.tempo = append(st.mi.tempo, tempoPoint{offset: off, bpm: bpm, metronome: true})
		}
	}
	if jumps {
		return
	}
	// the jumps of a direction without a sound, from its words and signs
	toCoda := false
	for _, t := range texts {
		switch l := strings.ToLower(t); {
		case strings.HasPrefix(l, "d.c.") || strings.HasPrefix(l, "da capo") || strings.HasPrefix(l, "d. c."):
			st.mi.dacapo = true
		case strings.HasPrefix(l, "d.s.") || strings.HasPrefix(l, "dal segno") || strings.HasPrefix(l, "d. s."):
			s := ""
			st.mi.dalsegno = &s
		case strings.TrimRight(l, ".") == "fine":
			st.mi.fine = true
		case l == "segno":
			segno = true
		case l == "coda":
			coda = true
		case strings.HasPrefix(l, "to coda") || strings.HasPrefix(l, "al coda"):
			toCoda = true
			st.mi.tocoda = append(st.mi.tocoda, "")
		}
	}
	if segno {
		st.mi.segno = append(st.mi.segno, "")
	}
	if coda && !toCoda {
		st.mi.coda = append(st.mi.coda, "")
	}
}

// addTempoMark adds a tempo mark, once for the parts that repeat it and
// merged with a mark at the same place that has only words or only a
// metronome mark.
func (mi *measureInfo) addTempoMark(t music.TempoMark) {
	for i := range mi.m.Tempo {
		e := &mi.m.Tempo[i]
		if e.Offset != t.Offset {
			continue
		}
		switch {
		case e.Text == t.Text && e.PerMinute == t.PerMinute:
			return
		case e.Text == "" && t.PerMinute == 0:
			e.Text = t.Text
			return
		case e.PerMinute == 0 && t.Text == "":
			e.Beat, e.Dots, e.PerMinute = t.Beat, t.Dots, t.PerMinute
			return
		}
	}
	mi.m.Tempo = append(mi.m.Tempo, t)
}

// metronome reads a metronome mark: a note value and a count per minute
// ("c. 120" and "120-132" count 120). Marks that equate two note values
// have no count.
func metronome(n *node) *music.TempoMark {
	t := &music.TempoMark{Beat: music.Quarter}
	units := 0
	for _, k := range n.kids {
		switch k.name {
		case "beat-unit":
			units++
			if units == 1 {
				if nt, ok := noteTypes[k.trim()]; ok {
					t.Beat = nt
				}
			}
		case "beat-unit-dot":
			if units == 1 {
				t.Dots = min(t.Dots+1, 3)
			}
		case "per-minute":
			s := k.trim()
			start := strings.IndexAny(s, "0123456789")
			if start < 0 {
				continue
			}
			end := start
			for end < len(s) && (s[end] >= '0' && s[end] <= '9' || s[end] == '.') {
				end++
			}
			if v, ok := parseFloat(s[start:end]); ok && v > 0 && v < 10000 {
				t.PerMinute = v
			}
		}
	}
	if units != 1 {
		t.PerMinute = 0
	}
	return t
}

// octaveShift reads an octave line: the Direction, and the span whose
// notes are written shifted. MusicXML notes carry their sounding pitch;
// "down" (8va) shows them lower, "up" (8vb) higher.
func (r *reader) octaveShift(p *partState, st *mstate, k *node, staff, off int, dir func(*music.Direction)) {
	num := atoiDef(k.attr("number"), 1)
	size := atoiDef(k.attr("size"), 8)
	if size != 15 && size != 22 {
		size = 8
	}
	octaves := (size + 1) / 7
	key := [2]int{staff, num}
	closeOpen := func() {
		if s := p.openShifts[key]; s != nil {
			s.endM, s.endOff = st.idx, off
			delete(p.openShifts, key)
		}
	}
	switch k.attr("type") {
	case "down", "up":
		closeOpen()
		d := &music.Direction{Kind: music.DirOctave, Number: num, Size: size}
		s := &octaveSpan{staff: staff, delta: -octaves, startM: st.idx, startOff: off, endM: -1}
		if k.attr("type") == "up" {
			d.Size, s.delta = -size, octaves
		}
		dir(d)
		p.shifts = append(p.shifts, s)
		p.openShifts[key] = s
	case "stop":
		closeOpen()
		dir(&music.Direction{Kind: music.DirOctave, Stop: true, Number: num})
	}
}

// harmony reads a chord symbol.
func (r *reader) harmony(p *partState, st *mstate, n *node) {
	if n.attr("print-object") == "no" {
		return
	}
	text := harmonyText(n)
	if text == "" {
		return
	}
	pos := st.pos + p.quarters(n.textOf("offset"))
	st.pm.Directions = append(st.pm.Directions, &music.Direction{Offset: max(tickOf(pos), 0), Staff: p.staff(n),
		Kind: music.DirChord, Text: text, Above: placement(n.attr("placement"))})
}

// chordKinds are the suffixes of the chord kinds.
var chordKinds = map[string]string{
	"major": "", "minor": "m", "augmented": "+", "diminished": "dim", "dominant": "7", "major-seventh": "maj7",
	"minor-seventh": "m7", "diminished-seventh": "dim7", "augmented-seventh": "+7", "half-diminished": "m7b5",
	"major-minor": "m(maj7)", "major-sixth": "6", "minor-sixth": "m6", "dominant-ninth": "9", "major-ninth": "maj9",
	"minor-ninth": "m9", "dominant-11th": "11", "major-11th": "maj11", "minor-11th": "m11", "dominant-13th": "13",
	"major-13th": "maj13", "minor-13th": "m13", "suspended-second": "sus2", "suspended-fourth": "sus4",
	"power": "5", "augmented-ninth": "+9", "other": "", "none": "N.C.",
}

// harmonyText writes a chord symbol: root, kind, degrees and bass, with
// "#" and "b" for sharps and flats.
func harmonyText(n *node) string {
	var b strings.Builder
	kind := n.child("kind")
	if kind.trim() == "none" {
		if t, ok := kind.attrOK("text"); ok && t != "" {
			return t
		}
		return "N.C."
	}
	root := n.child("root")
	switch {
	case root != nil:
		step := root.child("root-step")
		if t, ok := step.attrOK("text"); ok {
			b.WriteString(t)
		} else {
			b.WriteString(strings.ToUpper(step.trim()))
			b.WriteString(alterText(root.textOf("root-alter")))
		}
	case n.child("function") != nil:
		b.WriteString(n.textOf("function"))
	case n.child("numeral") != nil:
		b.WriteString(n.child("numeral").textOf("numeral-root"))
	}
	if t, ok := kind.attrOK("text"); ok {
		b.WriteString(t)
	} else {
		b.WriteString(chordKinds[kind.trim()])
	}
	for _, d := range n.kids {
		if d.name != "degree" || d.attr("print-object") == "no" {
			continue
		}
		v, alter := d.textOf("degree-value"), alterText(d.textOf("degree-alter"))
		typ := d.child("degree-type")
		if t, ok := typ.attrOK("text"); ok {
			b.WriteString(t + alter + v)
			continue
		}
		switch typ.trim() {
		case "add":
			b.WriteString("add" + alter + v)
		case "alter":
			b.WriteString(alter + v)
		case "subtract":
			b.WriteString("no" + v)
		}
	}
	if bass := n.child("bass"); bass != nil {
		b.WriteString("/")
		step := bass.child("bass-step")
		if t, ok := step.attrOK("text"); ok {
			b.WriteString(t)
		} else {
			b.WriteString(strings.ToUpper(step.trim()))
			b.WriteString(alterText(bass.textOf("bass-alter")))
		}
	}
	return strings.TrimSpace(b.String())
}

// alterText writes a chromatic alteration as "#" or "b".
func alterText(s string) string {
	v, _ := parseFloat(s)
	switch a := int(math.Round(v)); {
	case a > 0:
		return strings.Repeat("#", min(a, 2))
	case a < 0:
		return strings.Repeat("b", min(-a, 2))
	}
	return ""
}

// barStyles map the bar styles.
var barStyles = map[string]music.Barline{
	"regular": music.BarRegular, "dotted": music.BarDashed, "dashed": music.BarDashed, "heavy": music.BarFinal,
	"light-light": music.BarDouble, "light-heavy": music.BarFinal, "heavy-light": music.BarHeavyLight,
	"heavy-heavy": music.BarFinal, "tick": music.BarRegular, "short": music.BarRegular, "none": music.BarNone,
}

// barline reads a bar line: its style, repeat, volta, and the segno or
// coda sign on it.
func (r *reader) barline(p *partState, st *mstate, n *node) {
	mi, m := st.mi, st.mi.m
	left := n.attr("location") == "left"
	if n.attr("location") == "middle" {
		return
	}
	if style, ok := barStyles[n.textOf("bar-style")]; ok && style != music.BarRegular {
		if left && m.Left == music.BarRegular {
			m.Left = style
		} else if !left && m.Right == music.BarRegular {
			m.Right = style
		}
	}
	if rep := n.child("repeat"); rep != nil {
		switch rep.attr("direction") {
		case "forward":
			m.Left, mi.left = music.BarRepeatStart, true
		case "backward":
			m.Right, mi.right = music.BarRepeatEnd, true
			if t, ok := atoi(rep.attr("times")); ok && t >= 1 {
				m.Times = max(m.Times, min(t, 100))
			}
		}
	}
	if e := n.child("ending"); e != nil && (mi.endingPart == nil || mi.endingPart == p) {
		mi.endingPart = p
		switch e.attr("type") {
		case "start":
			mi.endingStart = &endingStart{number: e.attr("number"), text: collapse(e.text)}
		case "stop":
			mi.endStop = true
		case "discontinue":
			mi.endStop, mi.endOpen = true, true
		}
	}
	off := 0
	if !left {
		off = tickOf(st.pos)
	}
	for _, k := range n.kids {
		switch k.name {
		case "segno":
			st.pm.Directions = append(st.pm.Directions, &music.Direction{Offset: off, Kind: music.DirSegno})
		case "coda":
			st.pm.Directions = append(st.pm.Directions, &music.Direction{Offset: off, Kind: music.DirCoda})
		}
	}
	if v, ok := n.attrOK("segno"); ok {
		mi.segno = append(mi.segno, v)
	}
	if v, ok := n.attrOK("coda"); ok {
		mi.coda = append(mi.coda, v)
	}
}

// sound reads a sound at a position: tempo, dynamics and jumps. It
// reports whether it has jumps.
func (r *reader) sound(p *partState, st *mstate, n *node, pos float64) bool {
	mi := st.mi
	off := max(tickOf(pos), 0)
	if v, ok := parseFloat(n.attr("tempo")); ok && v > 0 && v < 10000 {
		mi.tempo = append(mi.tempo, tempoPoint{offset: off, bpm: v})
	}
	if v, ok := parseFloat(n.attr("dynamics")); ok && v >= 0 {
		p.dyn = append(p.dyn, dynPoint{measure: st.idx, offset: off, seq: st.nextSeq(), vel: min(max(int(math.Round(v*90/100)), 1), 127)})
	}
	jumps := false
	if n.yes("dacapo") {
		mi.dacapo, jumps = true, true
	}
	if v, ok := n.attrOK("dalsegno"); ok {
		mi.dalsegno, jumps = &v, true
	}
	if v, ok := n.attrOK("segno"); ok {
		mi.segno, jumps = append(mi.segno, v), true
	}
	if v, ok := n.attrOK("coda"); ok {
		mi.coda, jumps = append(mi.coda, v), true
	}
	if v, ok := n.attrOK("tocoda"); ok {
		mi.tocoda, jumps = append(mi.tocoda, v), true
	}
	if _, ok := n.attrOK("fine"); ok {
		mi.fine, jumps = true, true
	}
	return jumps
}

// endings makes the voltas: each bracket from the measure that starts it
// to the one that stops it (or only its first measure when nothing does).
// Voltas without numbers are numbered in the order of their group (the
// voltas one after the other).
func (r *reader) endings() {
	var group []*bracket
	closeGroup := func() {
		top := 0
		for _, b := range group {
			for _, n := range b.numbers {
				top = max(top, n)
			}
		}
		for _, b := range group {
			b.groupMax = top
			b.final = contains(b.numbers, top)
		}
		group = group[:0]
	}
	for i := 0; i < len(r.ms); i++ {
		es := r.ms[i].endingStart
		if len(group) > 0 && group[len(group)-1].last != i-1 {
			closeGroup()
		}
		if es == nil {
			continue
		}
		end := i
		for j := i; j < len(r.ms); j++ {
			if j > i && r.ms[j].endingStart != nil {
				break
			}
			if r.ms[j].endStop {
				end = j
				break
			}
		}
		b := &bracket{numbers: endingNumbers(es.number), first: i, last: end}
		if len(b.numbers) == 0 {
			b.numbers = []int{len(group) + 1}
		}
		text := es.text
		if text == "" {
			var parts []string
			for _, n := range b.numbers {
				parts = append(parts, strconv.Itoa(n))
			}
			text = strings.Join(parts, ", ") + "."
		}
		group = append(group, b)
		for j := i; j <= end; j++ {
			r.ms[j].ending = b
			b.repeats = b.repeats || r.ms[j].right
			r.ms[j].m.Ending = &music.Ending{Text: text, Start: j == i, Stop: j == end && r.ms[j].endStop,
				Open: j == end && r.ms[j].endOpen}
		}
		i = end
	}
	if len(group) > 0 {
		closeGroup()
	}
}

// endingNumbers reads the passes of a volta: "1", "1, 2", "1 2", "1-3".
func endingNumbers(s string) []int {
	var out []int
	for _, f := range strings.FieldsFunc(s, func(c rune) bool { return c == ',' || c == ' ' || c == ';' || c == '.' }) {
		if a, b, ok := strings.Cut(f, "-"); ok {
			x, ok1 := atoi(a)
			y, ok2 := atoi(b)
			if ok1 && ok2 && x >= 1 && y >= x && y-x < 100 {
				for v := x; v <= y; v++ {
					out = append(out, v)
				}
			}
			continue
		}
		if v, ok := atoi(f); ok && v >= 1 && v <= 100 {
			out = append(out, v)
		}
	}
	return out
}

// finish completes a part: a measure for each of the score, the clefs of
// the first measure, the notes under octave lines written shifted, and
// the dynamics in score order.
func (p *partState) finish(r *reader) {
	if len(r.ms) > 0 {
		p.measureAt(len(r.ms) - 1)
	}
	p.part.Measures = p.part.Measures[:len(r.ms)]
	p.play = p.play[:len(r.ms)]
	if len(p.part.Measures) > 0 {
		first := p.part.Measures[0]
		var add []music.ClefChange
		for s := 0; s < p.part.Staves; s++ {
			found := false
			for _, c := range first.Clefs {
				if c.Staff == s && c.Offset == 0 {
					found = true
				}
			}
			if !found {
				c := music.Clef{Sign: music.ClefG}
				switch {
				case p.part.Percussion:
					c.Sign = music.ClefPercussion
				case s == 1 && p.part.Staves == 2:
					c.Sign = music.ClefF
				}
				add = append(add, music.ClefChange{Staff: s, Clef: c})
			}
		}
		first.Clefs = append(add, first.Clefs...)
	}
	for i, pm := range p.part.Measures {
		length := r.ms[i].m.Length
		sort.SliceStable(pm.Clefs, func(a, b int) bool { return pm.Clefs[a].Offset < pm.Clefs[b].Offset })
		for j := range pm.Clefs {
			pm.Clefs[j].Offset = min(pm.Clefs[j].Offset, length)
		}
		for _, d := range pm.Directions {
			d.Offset = min(max(d.Offset, 0), length)
		}
	}
	// a line ends where the next one on its staff starts: a stop placed
	// past the last note must not reach into it
	sort.SliceStable(p.shifts, func(a, b int) bool {
		x, y := p.shifts[a], p.shifts[b]
		return x.staff < y.staff || x.staff == y.staff && before(x.startM, x.startOff, y.startM, y.startOff)
	})
	for k := 0; k+1 < len(p.shifts); k++ {
		s, next := p.shifts[k], p.shifts[k+1]
		if s.staff == next.staff && (s.endM < 0 || before(next.startM, next.startOff, s.endM, s.endOff)) {
			s.endM, s.endOff = next.startM, next.startOff
		}
	}
	for _, s := range p.shifts {
		for i := s.startM; i < len(p.part.Measures) && (s.endM < 0 || i <= s.endM); i++ {
			for _, ev := range p.part.Measures[i].Events {
				if ev.Staff != s.staff || i == s.startM && ev.Offset < s.startOff || i == s.endM && ev.Offset >= s.endOff {
					continue
				}
				for _, n := range ev.Notes {
					n.Pitch.Octave += s.delta
				}
			}
		}
	}
	sort.SliceStable(p.dyn, func(a, b int) bool {
		x, y := p.dyn[a], p.dyn[b]
		if x.measure != y.measure {
			return x.measure < y.measure
		}
		if x.offset != y.offset {
			return x.offset < y.offset
		}
		return x.seq < y.seq
	})
}

// before reports whether a position (measure, offset) is before another.
func before(m1, off1, m2, off2 int) bool {
	return m1 < m2 || m1 == m2 && off1 < off2
}

// graceTiming gives the grace notes of a measure their time before their
// principal note: a sixteenth each at most, taken from the principal.
// Grace notes without one are not played.
func graceTiming(notes []*playNote) []*playNote {
	// a grace note may be on another staff than its principal note
	type key struct{ voice, offset int }
	graces := map[key][]*playEvent{}
	for _, n := range notes {
		if n.grace {
			k := key{n.voice, n.offset}
			if evs := graces[k]; len(evs) == 0 || evs[len(evs)-1] != n.ev {
				graces[k] = append(evs, n.ev)
			}
		}
	}
	if len(graces) == 0 {
		return notes
	}
	step := map[key]int{}
	for k, evs := range graces {
		principal := 0
		for _, n := range notes {
			if !n.grace && n.voice == k.voice && n.offset == k.offset && n.dur > 0 {
				if principal == 0 || n.dur < principal {
					principal = n.dur
				}
			}
		}
		if g := min(music.PPQ/4, principal/(2*len(evs))); g >= 1 {
			step[k] = g
		}
	}
	out := notes[:0]
	for _, n := range notes {
		k := key{n.voice, n.offset}
		g, ok := step[k]
		switch {
		case n.grace && !ok:
			continue
		case n.grace:
			i := 0
			for i < len(graces[k]) && graces[k][i] != n.ev {
				i++
			}
			n.offset += i * g
			n.dur, n.grace = g, false
		case ok:
			shift := g * len(graces[k])
			n.offset += shift
			n.dur = max(n.dur-shift, 1)
		}
		out = append(out, n)
	}
	sort.SliceStable(out, func(a, b int) bool { return out[a].offset < out[b].offset })
	return out
}
