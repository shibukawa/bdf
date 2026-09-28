package music

import (
	"fmt"
	"math"
	"slices"
	"sort"
	"strings"
)

// NotateOptions control how a performance is written as a score.
type NotateOptions struct {
	// Time overrides the time signature (MML has none: 4/4 by default).
	Time *TimeSig
	// Key overrides the key signature; nil takes the performance's or
	// estimates it from the notes.
	Key *KeySig
}

// Notate writes a performance as a score: it divides the tracks into
// measures, quantizes the notes to note values (with ties across beats and
// bar lines, dots and triplets), separates voices, chooses clefs and a
// grand staff for wide parts, spells the notes in the key, and keeps the
// performance as what is played.
func Notate(p *Performance, o NotateOptions) *Score {
	return notate(p, &o)
}

// maxMeasures bounds the measures of a notated performance.
const maxMeasures = 20000

// maxStaffMeasures bounds the measures of all the staves of a score
// together (its staves times its measures): every staff has at least a
// rest in every measure. It is a variable for the tests.
var maxStaffMeasures = 1_000_000

// errStaves is the error for a score of more than maxStaves.
func errStaves() error {
	return fmt.Errorf("music: the score has more than %d staves", maxStaves)
}

// errStaffMeasures is the error for a score of more than maxStaffMeasures:
// n of what (staves, tracks) in its measures.
func errStaffMeasures(n int, what string, measures int) error {
	return fmt.Errorf("music: the score has more than %d measures on its staves (%d %s of %d measures)", maxStaffMeasures, n, what, measures)
}

// measureSpan is a measure of a performance on its timeline.
type measureSpan struct {
	start, length int
	time          TimeSig
}

// beat returns the length of the beat of a time signature and whether the
// meter is compound (beats of three eighths).
func (t TimeSig) beat() (int, bool) {
	if t.BeatType == 8 && t.Beats%3 == 0 && t.Beats > 3 {
		return 3 * PPQ / 2, true
	}
	if t.BeatType <= 0 {
		return PPQ, false
	}
	return 4 * PPQ / t.BeatType, false
}

func notate(p *Performance, o *NotateOptions) *Score {
	s := &Score{Title: p.Title, Subtitle: p.Subtitle, Composer: p.Composer, Lyricist: p.Lyricist,
		Arranger: p.Arranger, Rights: p.Copyright, Play: p}

	end := 0
	for _, t := range p.Tracks {
		for _, n := range t.Notes {
			end = max(end, n.Tick+max(n.Dur, 1))
		}
	}
	if end == 0 {
		end = p.End
	}

	// measures from the time signatures
	times := slices.Clone(p.Time)
	if o.Time != nil {
		times = []TimeChange{{Tick: 0, Time: *o.Time}}
	}
	sort.SliceStable(times, func(i, j int) bool { return times[i].Tick < times[j].Tick })
	cur := TimeSig{Beats: 4, BeatType: 4}
	var spans []measureSpan
	ti := 0
	t := 0
	for (t < end || len(spans) == 0) && len(spans) < maxMeasures {
		changed := len(spans) == 0
		for ti < len(times) && times[ti].Tick <= t {
			if times[ti].Time.Beats > 0 && times[ti].Time.BeatType > 0 {
				changed = changed || times[ti].Time != cur
				cur = times[ti].Time
			}
			ti++
		}
		length := cur.Length()
		if ti < len(times) && times[ti].Tick < t+length {
			length = times[ti].Tick - t // a time signature that changes inside a measure cuts it short
		}
		spans = append(spans, measureSpan{t, length, cur})
		m := &Measure{Length: length}
		if changed {
			tm := cur
			m.Time = &tm
		}
		s.Measures = append(s.Measures, m)
		t += length
	}
	s.Measures[len(s.Measures)-1].Right = BarFinal
	if t < end {
		s.warnings = append(s.warnings, fmt.Sprintf("the music is longer than %d measures; the rest is left out", maxMeasures))
	}
	// a track is a part of one or two staves, written in every measure
	tracks := 0
	for _, t := range p.Tracks {
		if len(t.Notes) > 0 {
			tracks++
		}
	}
	switch {
	case tracks > maxStaves:
		s.err = errStaves()
		return s
	case tracks*len(spans) > maxStaffMeasures:
		s.err = errStaffMeasures(tracks, "tracks", len(spans))
		return s
	}

	measureAt := func(tick int) int {
		i := sort.Search(len(spans), func(i int) bool { return spans[i].start > tick }) - 1
		return max(i, 0)
	}

	// key signatures
	var keys []KeyChange
	switch {
	case o.Key != nil:
		keys = []KeyChange{{Tick: 0, Key: *o.Key}}
	case len(p.Key) > 1 || len(p.Key) == 1 && (p.Key[0].Key != KeySig{} || p.Key[0].Tick != 0):
		keys = slices.Clone(p.Key)
		sort.SliceStable(keys, func(i, j int) bool { return keys[i].Tick < keys[j].Tick })
	default:
		// no key, or the C major that sequencers write by default
		var hist [12]float64
		flats := 0
		for _, t := range p.Tracks {
			if t.Channel == 9 {
				continue
			}
			for _, n := range t.Notes {
				hist[((n.Key%12)+12)%12] += float64(max(n.Dur, 1))
				flats -= n.Spell
			}
		}
		keys = []KeyChange{{Tick: 0, Key: estimateKey(hist, flats)}}
	}
	keyAt := make([]*KeySig, len(spans)) // the key that starts at each measure
	keyOfMeasure := make([]KeySig, len(spans))
	for _, k := range keys {
		i := measureAt(k.Tick)
		if spans[i].start < k.Tick && i+1 < len(spans) {
			i++ // a change inside a measure starts at the next one
		}
		kk := k.Key
		kk.Fifths = min(max(kk.Fifths, -7), 7)
		keyAt[i] = &kk
	}
	cur2 := KeySig{}
	for i := range spans {
		if keyAt[i] != nil {
			if *keyAt[i] == cur2 && i > 0 {
				keyAt[i] = nil
			} else {
				cur2 = *keyAt[i]
			}
		}
		keyOfMeasure[i] = cur2
	}

	s.tempoMarks(p, spans, measureAt)

	for _, t := range p.Tracks {
		if len(t.Notes) == 0 {
			continue
		}
		part := notateTrack(t, spans, keyOfMeasure, measureAt)
		if part == nil {
			continue // its notes are all in what is left out
		}
		for i, pm := range part.Measures {
			if keyAt[i] != nil && !part.Percussion {
				k := *keyAt[i]
				pm.Key = &k
			}
		}
		s.Parts = append(s.Parts, part)
	}
	return s
}

// tempoMarks shows the tempo changes of a performance that last: the first
// tempo and each one that differs from the last shown by 3% or more and
// holds for a beat, at most one per measure.
func (s *Score) tempoMarks(p *Performance, spans []measureSpan, measureAt func(int) int) {
	tempos := slices.Clone(p.Tempo)
	sort.SliceStable(tempos, func(i, j int) bool { return tempos[i].Tick < tempos[j].Tick })
	last := 0.0
	lastMeasure := -1
	for i, t := range tempos {
		if t.BPM <= 0 {
			continue
		}
		if i+1 < len(tempos) && tempos[i+1].Tick-t.Tick < PPQ {
			continue
		}
		if last > 0 && math.Abs(t.BPM-last)/last < 0.03 {
			continue
		}
		m := measureAt(t.Tick)
		if m == lastMeasure {
			continue
		}
		sp := spans[m]
		beat, compound := sp.time.beat()
		mark := TempoMark{Offset: roundTo(t.Tick-sp.start, beat), Beat: Quarter, PerMinute: math.Round(t.BPM)}
		if compound {
			mark.Dots = 1
			mark.PerMinute = math.Round(t.BPM / 1.5)
		}
		if mark.Offset >= sp.length {
			mark.Offset = 0
		}
		s.Measures[m].Tempo = append(s.Measures[m].Tempo, mark)
		last = t.BPM
		lastMeasure = m
	}
}

func roundTo(v, grid int) int {
	if grid <= 0 {
		return v
	}
	return (v + grid/2) / grid * grid
}

// qnote is a note of a staff after quantization.
type qnote struct {
	on, off   int
	key       int
	vel       int
	spell     int
	head      Notehead
	pitch     *Pitch // a fixed display pitch (percussion)
	voiceHint int    // percussion: 2 for the feet
}

// beatGrid is how the notes of one beat are quantized.
type beatGrid struct {
	start, length int
	compound      bool
	grid          int
	triplet       bool
}

// grids works out the quantization of every beat of the measures from the
// onsets of the notes: the coarsest grid (eighths, triplet eighths,
// sixteenths, triplet sixteenths, thirty-seconds of a quarter beat) that
// holds every onset of the beat closely enough, or the closest one.
func grids(spans []measureSpan, onsets []int) []beatGrid {
	var beats []beatGrid
	for _, sp := range spans {
		b, compound := sp.time.beat()
		for t := 0; t < sp.length; t += b {
			beats = append(beats, beatGrid{start: sp.start + t, length: min(b, sp.length-t), compound: compound})
		}
	}
	slices.Sort(onsets)
	j := 0
	for i := range beats {
		bg := &beats[i]
		for j < len(onsets) && onsets[j] < bg.start {
			j++
		}
		k := j
		for k < len(onsets) && onsets[k] < bg.start+bg.length {
			k++
		}
		in := onsets[j:k]
		type cand struct {
			grid    int
			triplet bool
		}
		var cands []cand
		l := bg.length
		if bg.compound {
			// a compound beat: eighths, sixteenths, thirty-seconds
			cands = []cand{{l / 3, false}, {l / 6, false}, {l / 12, false}}
		} else {
			cands = []cand{{l / 2, false}, {l / 3, true}, {l / 4, false}, {l / 6, true}, {l / 8, false}, {l / 12, true}}
		}
		best, bestErr := cands[len(cands)-1], math.MaxInt
		for _, c := range cands {
			if c.grid <= 0 {
				continue
			}
			maxErr, sum := 0, 0
			for _, on := range in {
				r := (on - bg.start) % c.grid
				e := min(r, c.grid-r)
				maxErr = max(maxErr, e)
				sum += e
			}
			if maxErr <= min(c.grid/4, PPQ/24) {
				best = c
				break
			}
			if sum < bestErr {
				best, bestErr = c, sum
			}
		}
		// a beat too short for a grid (a time signature that changes a tick
		// after a bar line) is written as it is
		bg.grid, bg.triplet = best.grid, best.triplet && best.grid > 0
		j = k
	}
	return beats
}

// gridAt returns the beat that holds tick.
func gridAt(beats []beatGrid, tick int) *beatGrid {
	i := sort.Search(len(beats), func(i int) bool { return beats[i].start > tick }) - 1
	if i < 0 {
		i = 0
	}
	if i >= len(beats) {
		i = len(beats) - 1
	}
	return &beats[i]
}

func (bg *beatGrid) snap(tick int) int {
	if bg.grid <= 0 {
		return tick
	}
	return bg.start + roundTo(tick-bg.start, bg.grid)
}

// notateTrack writes one track as a part, nil when none of its notes is in
// the measures.
func notateTrack(t *Track, spans []measureSpan, keys []KeySig, measureAt func(int) int) *Part {
	part := &Part{Name: t.Name, Staves: 1, Percussion: t.Channel == 9}
	last := spans[len(spans)-1]
	end := last.start + last.length

	var onsets []int
	for _, n := range t.Notes {
		onsets = append(onsets, n.Tick)
	}
	beats := grids(spans, onsets)

	// quantize
	var notes []qnote
	for _, n := range t.Notes {
		if n.Tick >= end || n.Key < 0 || n.Key > 127 {
			continue
		}
		on := gridAt(beats, n.Tick).snap(n.Tick)
		off := n.Tick + max(n.Dur, 1)
		off = gridAt(beats, off).snap(off)
		if g := gridAt(beats, on).grid; off <= on {
			off = on + max(g, 1)
		}
		off = min(off, end)
		if on >= end {
			continue
		}
		q := qnote{on: on, off: off, key: n.Key, vel: n.Vel, spell: n.Spell}
		if part.Percussion {
			d := drumOf(n.Key)
			q.pitch, q.head, q.voiceHint = &d.pitch, d.head, d.voice
		}
		notes = append(notes, q)
	}
	sort.SliceStable(notes, func(i, j int) bool {
		if notes[i].on != notes[j].on {
			return notes[i].on < notes[j].on
		}
		return notes[i].key > notes[j].key
	})
	// the same key twice at an onset is one note
	notes = slices.CompactFunc(notes, func(a, b qnote) bool { return a.on == b.on && a.key == b.key })
	if len(notes) == 0 {
		return nil
	}

	// staves and clefs
	var clefs []Clef
	var staffOf func(q qnote) int
	switch {
	case part.Percussion:
		clefs = []Clef{{Sign: ClefPercussion}}
		staffOf = func(qnote) int { return 0 }
		// a drum hit lasts until the next hit of its voice, within its beat
		for i := range notes {
			q := &notes[i]
			next := gridAt(beats, q.on)
			q.off = next.start + next.length
			for k := i + 1; k < len(notes); k++ {
				if notes[k].on > q.on && notes[k].voiceHint == q.voiceHint {
					q.off = min(q.off, notes[k].on)
					break
				}
			}
		}
	default:
		keysOf := make([]int, len(notes))
		below := 0
		lo, hi := 127, 0
		for i, q := range notes {
			keysOf[i] = q.key
			if q.key < 60 {
				below++
			}
			lo, hi = min(lo, q.key), max(hi, q.key)
		}
		slices.Sort(keysOf)
		median := keysOf[len(keysOf)/2]
		frac := float64(below) / float64(len(notes))
		// keyboards (pianos, chromatic percussion, organs, the harp) take a
		// grand staff when they use both clefs; other parts only when their
		// range is too wide for one
		keyboard := t.Program <= 23 || t.Program == 46
		guitar := t.Program >= 24 && t.Program <= 31
		switch {
		case !guitar && (keyboard && frac >= 0.1 && frac <= 0.9 && lo < 55 && hi >= 60 ||
			frac >= 0.25 && frac <= 0.75 && hi-lo >= 36):
			part.Staves = 2
			clefs = []Clef{{Sign: ClefG}, {Sign: ClefF}}
			staffOf = func(q qnote) int {
				if q.key >= 60 {
					return 0
				}
				return 1
			}
		case guitar && median < 67:
			clefs = []Clef{{Sign: ClefG, Octave: -1}} // guitar: sounds an octave lower
			staffOf = func(qnote) int { return 0 }
		case median < 57:
			clefs = []Clef{{Sign: ClefF}}
			staffOf = func(qnote) int { return 0 }
		default:
			clefs = []Clef{{Sign: ClefG}}
			staffOf = func(qnote) int { return 0 }
		}
	}

	part.Measures = make([]*PartMeasure, len(spans))
	for i := range part.Measures {
		part.Measures[i] = &PartMeasure{}
	}
	for st, c := range clefs {
		part.Measures[0].Clefs = append(part.Measures[0].Clefs, ClefChange{Offset: 0, Staff: st, Clef: c})
	}

	lyrics := map[int]string{}
	for _, l := range t.Lyrics {
		// a syllable belongs to the note that starts nearest to it
		on := gridAt(beats, l.Tick).snap(l.Tick)
		if s := lyrics[on]; s != "" {
			lyrics[on] = s + l.Text
		} else {
			lyrics[on] = l.Text
		}
	}

	for st := range clefs {
		var sn []qnote
		for _, q := range notes {
			if staffOf(q) == st {
				sn = append(sn, q)
			}
		}
		n := &staffNotator{part: part, staff: st, spans: spans, beats: beats, keys: keys, measureAt: measureAt,
			lyrics: lyrics, percussion: part.Percussion}
		n.run(sn, end)
	}
	return part
}

// staffNotator writes the notes of one staff into the measures of a part.
type staffNotator struct {
	part       *Part
	staff      int
	spans      []measureSpan
	beats      []beatGrid
	keys       []KeySig
	measureAt  func(int) int
	lyrics     map[int]string
	percussion bool
}

// chord is notes of one voice that start and end together.
type chord struct {
	on, off int
	notes   []qnote
}

func (n *staffNotator) run(notes []qnote, end int) {
	// legato: a note that overlaps the next onset a little ends there
	for i := range notes {
		if n.percussion {
			break
		}
		q := &notes[i]
		for k := i + 1; k < len(notes); k++ {
			if notes[k].on > q.on {
				if over := q.off - notes[k].on; over > 0 && over <= max(gridAt(n.beats, q.on).grid, (q.off-q.on)/4) {
					q.off = notes[k].on
				}
				break
			}
		}
	}
	// chords of notes with the same onset and end
	var chords []*chord
	for _, q := range notes {
		if len(chords) > 0 {
			c := chords[len(chords)-1]
			if c.on == q.on && c.off == q.off && c.notes[0].voiceHint == q.voiceHint {
				c.notes = append(c.notes, q)
				continue
			}
		}
		chords = append(chords, &chord{on: q.on, off: q.off, notes: []qnote{q}})
	}
	// two voices: a chord goes to the first voice that is free
	var voices [2][]*chord
	var busy [2]int
	for _, c := range chords {
		pref := []int{0, 1}
		if c.notes[0].voiceHint == 2 {
			pref = []int{1, 0}
		}
		placed := false
		for _, v := range pref {
			if busy[v] <= c.on {
				voices[v] = append(voices[v], c)
				busy[v] = c.off
				placed = true
				break
			}
		}
		if placed {
			continue
		}
		// both voices busy: join a chord that starts together, or cut the
		// first voice short
		joined := false
		for v := range voices {
			if k := len(voices[v]); k > 0 && voices[v][k-1].on == c.on {
				last := voices[v][k-1]
				last.notes = append(last.notes, c.notes...)
				last.off = max(last.off, c.off)
				busy[v] = last.off
				joined = true
				break
			}
		}
		if joined {
			continue
		}
		v := pref[0]
		k := len(voices[v])
		voices[v][k-1].off = c.on
		voices[v] = append(voices[v], c)
		busy[v] = c.off
	}
	// the higher line is the first voice in each measure
	n.swapVoices(&voices)

	for v := range voices {
		n.voice(voices[v], v+1, end)
	}
}

// swapVoices makes the first voice the higher one in each measure where
// the second voice sounds.
func (n *staffNotator) swapVoices(voices *[2][]*chord) {
	if len(voices[1]) == 0 || n.percussion {
		return
	}
	type acc struct{ sum, cnt [2]float64 }
	per := map[int]*acc{}
	for v := range voices {
		for _, c := range voices[v] {
			m := n.measureAt(c.on)
			a := per[m]
			if a == nil {
				a = &acc{}
				per[m] = a
			}
			for _, q := range c.notes {
				a.sum[v] += float64(q.key)
				a.cnt[v]++
			}
		}
	}
	swap := map[int]bool{}
	for m, a := range per {
		if a.cnt[0] > 0 && a.cnt[1] > 0 && a.sum[1]/a.cnt[1] > a.sum[0]/a.cnt[0] {
			swap[m] = true
		}
	}
	if len(swap) == 0 {
		return
	}
	var out [2][]*chord
	for v := range voices {
		for _, c := range voices[v] {
			w := v
			if swap[n.measureAt(c.on)] {
				w = 1 - v
			}
			out[w] = append(out[w], c)
		}
	}
	for v := range out {
		sort.SliceStable(out[v], func(i, j int) bool { return out[v][i].on < out[v][j].on })
	}
	*voices = out
}

// voice writes the chords of one voice, with rests in the gaps: in every
// measure for the first voice, in the measures it sounds in for the
// second.
func (n *staffNotator) voice(chords []*chord, voice, end int) {
	used := map[int]bool{}
	if voice == 1 {
		for i := range n.spans {
			used[i] = true
		}
	} else {
		for _, c := range chords {
			for m := n.measureAt(c.on); m < len(n.spans) && n.spans[m].start < c.off; m++ {
				used[m] = true
			}
		}
	}
	t := 0
	emitRest := func(a, b int) {
		for a < b {
			m := n.measureAt(a)
			sp := n.spans[m]
			e := min(b, sp.start+sp.length)
			if used[m] {
				n.segment(m, a, e, nil, voice, false, false)
			}
			a = e
		}
	}
	for _, c := range chords {
		if c.on > t {
			emitRest(t, c.on)
		}
		on, off := max(c.on, t), c.off
		if off <= on {
			continue
		}
		first := true
		for a := on; a < off; {
			m := n.measureAt(a)
			sp := n.spans[m]
			e := min(off, sp.start+sp.length)
			n.segment(m, a, e, c.notes, voice, !first, e < off)
			first = false
			a = e
		}
		t = off
	}
	last := n.spans[len(n.spans)-1]
	if e := last.start + last.length; t < e {
		emitRest(t, e)
	}
	_ = end
}

// segment writes a stretch of a note (or rest) that lies in one measure,
// divided into note values; tieIn and tieOut tie it to the stretches
// before and after.
func (n *staffNotator) segment(m, a, b int, notes []qnote, voice int, tieIn, tieOut bool) {
	sp := n.spans[m]
	pm := n.part.Measures[m]
	if notes == nil && a == sp.start && b == sp.start+sp.length {
		pm.Events = append(pm.Events, &Event{Offset: 0, Duration: sp.length, Staff: n.staff, Voice: voice,
			Rest: true, MeasureRest: true, Type: Whole})
		return
	}
	beat, compound := sp.time.beat()
	// divide at the beats whose grid changes between straight and triplet
	type piece struct {
		a, b    int
		triplet *beatGrid
	}
	var pieces []piece
	for x := a; x < b; {
		bg := gridAt(n.beats, x)
		e := min(b, bg.start+bg.length)
		// a stretch of a triplet beat that starts and ends on its straight
		// sixteenths is written straight (the other voice has the triplets)
		straight := max(bg.length/4, 1)
		if bg.triplet && ((x-bg.start)%straight != 0 || (e-bg.start)%straight != 0) {
			pieces = append(pieces, piece{x, e, bg})
		} else if k := len(pieces); k > 0 && pieces[k-1].triplet == nil && pieces[k-1].b == x {
			pieces[k-1].b = e
		} else {
			pieces = append(pieces, piece{x, e, nil})
		}
		x = e
	}
	var evs []*Event
	for _, pc := range pieces {
		if pc.triplet != nil {
			bg := pc.triplet
			tup := n.tuplet(m, voice, bg)
			// in the tuplet, time runs 3/2 as fast
			va, vb := (pc.a-bg.start)*3/2, (pc.b-bg.start)*3/2
			for _, v := range divide(va, vb, bg.length*3/2, bg.length*3/2, false, notes == nil) {
				evs = append(evs, &Event{Offset: bg.start + v.a*2/3 - sp.start, Duration: v.ticks * 2 / 3,
					Type: v.typ, Dots: v.dots, Tuplet: tup})
			}
			continue
		}
		for _, v := range divide(pc.a-sp.start, pc.b-sp.start, beat, sp.length, compound, notes == nil) {
			evs = append(evs, &Event{Offset: v.a, Duration: v.ticks, Type: v.typ, Dots: v.dots})
		}
	}
	for i, e := range evs {
		e.Staff, e.Voice = n.staff, voice
		if notes == nil {
			e.Rest = true
		} else {
			for _, q := range notes {
				nt := &Note{Head: q.head}
				if q.pitch != nil {
					nt.Pitch = *q.pitch
				} else {
					nt.Pitch = spell(q.key, n.keys[m], q.spell)
				}
				nt.TieStop = tieIn || i > 0
				nt.TieStart = tieOut || i < len(evs)-1
				e.Notes = append(e.Notes, nt)
			}
			slices.SortFunc(e.Notes, func(x, y *Note) int { return pitchOrder(x.Pitch) - pitchOrder(y.Pitch) })
			if !tieIn && i == 0 {
				if l := n.lyrics[sp.start+e.Offset]; l != "" {
					e.Lyrics = []*LyricSyllable{lyricSyllable(l)}
					delete(n.lyrics, sp.start+e.Offset)
				}
			}
		}
		pm.Events = append(pm.Events, e)
	}
}

func pitchOrder(p Pitch) int { return p.Octave*7 + p.Step }

// lyricSyllable reads a syllable as MIDI lyrics write it: a trailing
// hyphen joins it to the next.
func lyricSyllable(s string) *LyricSyllable {
	l := &LyricSyllable{Verse: 1}
	s = strings.TrimSpace(strings.Trim(s, "/\\\r\n"))
	if strings.HasSuffix(s, "-") {
		s = strings.TrimSuffix(s, "-")
		l.Syllabic = SyllableBegin
	}
	l.Text = s
	return l
}

// tuplet returns the tuplet of a triplet beat in a voice, shared by its
// events.
func (n *staffNotator) tuplet(m, voice int, bg *beatGrid) *Tuplet {
	pm := n.part.Measures[m]
	for _, e := range pm.Events {
		if e.Tuplet != nil && e.Voice == voice && e.Staff == n.staff {
			sp := n.spans[m]
			if abs := sp.start + e.Offset; abs >= bg.start && abs < bg.start+bg.length {
				return e.Tuplet
			}
		}
	}
	return &Tuplet{Actual: 3, Normal: 2, ShowNumber: true}
}

// value is a note value placed in a measure.
type value struct {
	a, ticks int
	typ      NoteType
	dots     int
}

// noteValues are the note values in order of length, with their dots.
var noteValues = func() []value {
	var vs []value
	for t := Breve; t <= N128th; t++ {
		for dots := 1; dots >= 0; dots-- {
			if t == N128th && dots > 0 {
				continue
			}
			vs = append(vs, value{ticks: t.Ticks(dots), typ: t, dots: dots})
		}
	}
	sort.SliceStable(vs, func(i, j int) bool { return vs[i].ticks > vs[j].ticks })
	return vs
}()

// divide writes the stretch [a, b) of a measure of length length and beat
// beat as note values: a value that starts on a beat may span whole beats,
// one that starts off the beat stays in its beat; rests are aligned to
// their own length and dotted only as a whole beat of a compound meter.
func divide(a, b, beat, length int, compound, rest bool) []value {
	var out []value
	for a < b {
		ok := false
		for _, v := range noteValues {
			if v.ticks > b-a || v.typ == Breve && length < v.ticks {
				continue
			}
			base := v.typ.Ticks(0)
			if rest {
				if v.dots > 0 && !(compound && v.ticks == beat && a%beat == 0) {
					continue
				}
				if a%base != 0 && !(v.dots > 0) {
					continue
				}
			}
			if a%beat == 0 {
				if v.ticks > beat {
					if compound && v.ticks%beat != 0 {
						continue
					}
					if !compound && base%beat != 0 {
						continue
					}
				}
			} else if a/beat != (a+v.ticks-1)/beat {
				continue
			}
			v.a = a
			out = append(out, v)
			a += v.ticks
			ok = true
			break
		}
		if !ok {
			// shorter than the shortest value: give it the shortest
			out = append(out, value{a: a, ticks: b - a, typ: N128th})
			break
		}
	}
	return out
}

// drum is how a General MIDI percussion key is written.
type drum struct {
	pitch Pitch
	head  Notehead
	voice int
}

// drumOf writes a General MIDI percussion key on a percussion staff, in
// the usual drum set notation: cymbals with x heads above the staff,
// drums in the spaces, the feet in the second voice.
func drumOf(key int) drum {
	p := func(step, oct int) Pitch { return Pitch{Step: step, Octave: oct} }
	switch key {
	case 35, 36: // bass drum
		return drum{p(3, 4), HeadNormal, 2}
	case 37: // side stick
		return drum{p(0, 5), HeadX, 1}
	case 38, 40: // snare
		return drum{p(0, 5), HeadNormal, 1}
	case 39: // hand clap
		return drum{p(0, 5), HeadX, 1}
	case 41, 43: // low floor toms
		return drum{p(5, 4), HeadNormal, 1}
	case 45, 47: // low and mid toms
		return drum{p(1, 5), HeadNormal, 1}
	case 48, 50: // high toms
		return drum{p(2, 5), HeadNormal, 1}
	case 42: // closed hi-hat
		return drum{p(4, 5), HeadX, 1}
	case 44: // pedal hi-hat
		return drum{p(1, 4), HeadX, 2}
	case 46: // open hi-hat
		return drum{p(4, 5), HeadX, 1}
	case 49, 57, 52, 55: // crash, chinese, splash
		return drum{p(5, 5), HeadX, 1}
	case 51, 59, 53: // ride
		return drum{p(3, 5), HeadX, 1}
	case 54, 56: // tambourine, cowbell
		return drum{p(2, 5), HeadX, 1}
	}
	return drum{p(0, 5), HeadNormal, 1}
}
