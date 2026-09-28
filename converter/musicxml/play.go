package musicxml

import (
	"sort"

	"github.com/shibukawa/bdf/converter/internal/music"
)

// played is a measure as played: the pass of its repeat.
type played struct {
	measure, pass int
}

// section is a repeated section being played.
type section struct {
	start, pass int
	explicit    bool // it starts at a forward repeat
}

// maxNesting caps the repeats inside repeats.
const maxNesting = 16

// playOrder expands the repeats, voltas and jumps into the order the
// measures are played in.
//
// A backward repeat goes back to its forward repeat, or without one to
// the start of the music or the measure after the last repeat that ended;
// its times (default 2, or the highest volta number around it) is the
// number of passes. Repeats may nest. On pass k only the voltas whose
// numbers hold k are played (voltas without numbers are numbered in their
// order); the volta that holds the last pass ends the section. D.C. and
// D.S. jump at the end of their measure, each once; after a jump repeats
// are not taken, the last volta of each group is played, Fine ends the
// music at the end of its measure and To Coda jumps to the coda (only one
// of them when the words of the jump say "al Fine" or "al Coda"). Jumps
// are those of the sound elements, or, without one, of the words "D.C.",
// "D.S.", "Fine" and "To Coda" and the segno and coda signs.
func (r *reader) playOrder() []played {
	n := len(r.ms)
	explicitToCoda := false
	for _, mi := range r.ms {
		if len(mi.tocoda) > 0 {
			explicitToCoda = true
		}
	}
	stack := []section{{start: 0, pass: 1}}
	top := func() *section { return &stack[len(stack)-1] }
	// done ends the section at the top after measure i: the section around
	// it goes on, unless it started only because nothing did before (a
	// backward repeat without a forward one goes back to the measure after
	// this repeat), or the ended one is inside its volta
	done := func(i int) {
		inVolta := top().start < n && r.ms[top().start].ending != nil
		stack = stack[:len(stack)-1]
		switch {
		case len(stack) == 0:
			stack = append(stack, section{start: i + 1, pass: 1})
		case !top().explicit && !inVolta:
			top().start, top().pass = i+1, 1
		}
	}
	var out []played
	i := 0
	tick := 0 // where the music has come to
	// the measures the jumps go to, found once
	segnos, codas := map[int]int{}, map[int]int{}
	segnoTarget := func(i int, name string) int {
		j, ok := segnos[i]
		if !ok {
			j = r.segnoTarget(i, name)
			segnos[i] = j
		}
		return j
	}
	codaTarget := func(i int, name string) int {
		j, ok := codas[i]
		if !ok {
			j = r.codaTarget(i, name)
			codas[i] = j
		}
		return j
	}
	jumped, codaTaken, back := false, false, false
	alFine, alCoda := false, false // what the last jump goes on to
	taken := map[int]bool{}        // the jumps taken
	jump := func(from, to int) {
		mi := r.ms[from]
		jumped, codaTaken, taken[from] = true, false, true
		alFine, alCoda = !mi.alCoda || mi.alFine, !mi.alFine || mi.alCoda
		i = to
	}
	for i < n {
		if len(out) >= maxPlays {
			r.warn("the repeats play more than %d measures; the performance stops there", maxPlays)
			break
		}
		mi := r.ms[i]
		// a volta is chosen when the music reaches it, not when a repeat
		// that starts at it goes back
		if b := mi.ending; b != nil && i == b.first && !back {
			play := b.final
			if !jumped {
				play = contains(b.numbers, top().pass)
			}
			if !play {
				i = b.last + 1
				if !jumped && b.repeats && (i >= n || r.ms[i].ending == nil) {
					done(i - 1)
				}
				continue
			}
		}
		back = false
		if mi.left && !jumped {
			switch {
			case i == top().start:
				top().explicit = true
			case len(stack) < maxNesting:
				stack = append(stack, section{start: i, pass: 1, explicit: true})
			}
		}
		if tick += mi.m.Length; tick > maxPlayQuarters*music.PPQ {
			r.warn("the music plays longer than %d quarter notes; the performance stops there", maxPlayQuarters)
			break
		}
		out = append(out, played{i, top().pass})

		if jumped && alFine && mi.fine {
			break
		}
		if jumped && alCoda && !codaTaken {
			target := -1
			switch {
			case len(mi.tocoda) > 0:
				target = codaTarget(i, mi.tocoda[0])
			case !explicitToCoda && len(mi.coda) > 0:
				// a coda sign where the music leaves for the coda: the
				// next coda sign is the coda
				target = codaTarget(i, "")
			}
			if target >= 0 {
				codaTaken = true
				i = target
				continue
			}
		}
		if !jumped {
			switch b := mi.ending; {
			case mi.right:
				if t := top(); t.pass < r.repeatTimes(i) {
					t.pass++
					i, back = t.start, true
					continue
				}
				done(i)
			case b != nil && i == b.last && !mi.left && top().pass > 1 && top().pass >= b.groupMax:
				done(i)
			}
		}
		if !taken[i] && mi.dacapo {
			jump(i, 0)
			continue
		}
		if !taken[i] && mi.dalsegno != nil {
			if j := segnoTarget(i, *mi.dalsegno); j >= 0 {
				jump(i, j)
				continue
			}
			r.warn("D.S. without a segno is not played")
		}
		i++
	}
	return out
}

func contains(s []int, v int) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

// repeatTimes returns the passes of the repeat that ends at a measure: its
// times, or the highest volta number of the voltas around it (the voltas
// one after the other are a group).
func (r *reader) repeatTimes(i int) int {
	times := r.ms[i].m.Times
	if times <= 0 {
		times = 2
	}
	if b := r.ms[i].ending; b != nil {
		times = min(max(times, b.groupMax), 100)
	}
	return times
}

// segnoTarget returns the measure a D.S. at measure i goes back to: the
// segno of that name, else the last segno before it, else the first.
func (r *reader) segnoTarget(i int, name string) int {
	if name != "" && name != "yes" {
		for j, mi := range r.ms {
			for _, s := range mi.segno {
				if s == name {
					return j
				}
			}
		}
	}
	for j := i; j >= 0; j-- {
		if len(r.ms[j].segno) > 0 {
			return j
		}
	}
	for j, mi := range r.ms {
		if len(mi.segno) > 0 {
			return j
		}
	}
	return -1
}

// codaTarget returns the measure To Coda at measure i goes to: the coda
// of that name after it, else the next coda.
func (r *reader) codaTarget(i int, name string) int {
	if name != "" && name != "yes" {
		for j := i + 1; j < len(r.ms); j++ {
			for _, c := range r.ms[j].coda {
				if c == name {
					return j
				}
			}
		}
	}
	for j := i + 1; j < len(r.ms); j++ {
		if len(r.ms[j].coda) > 0 {
			return j
		}
	}
	return -1
}

// perform makes the performance of the score: the measures in their
// play order, a track for each part.
func (r *reader) perform() {
	s := r.score
	order := r.playOrder()
	ticks := make([]int, len(order))
	tick := 0
	inOrder := len(order) == len(r.ms)
	for k, pl := range order {
		ticks[k] = tick
		tick += r.ms[pl.measure].m.Length
		if pl.measure != k {
			inOrder = false
		}
	}
	if !inOrder {
		s.PlayOrder = make([]music.PlayedMeasure, len(order))
		for k, pl := range order {
			s.PlayOrder[k] = music.PlayedMeasure{Tick: ticks[k], Measure: pl.measure}
		}
	}
	perf := &music.Performance{Title: s.Title, Subtitle: s.Subtitle, Composer: s.Composer, Lyricist: s.Lyricist,
		Arranger: s.Arranger, Copyright: s.Rights, End: tick}
	s.Play = perf
	perf.Tempo = r.tempos(order, ticks)

	// time and key signatures where they are first played
	var keyPart *music.Part
	for _, want := range []bool{true, false} {
		for _, p := range r.parts {
			if p.seen && !p.part.Percussion && (!p.transposed || !want) {
				keyPart = p.part
				break
			}
		}
		if keyPart != nil {
			break
		}
	}
	seen := make([]bool, len(r.ms))
	for k, pl := range order {
		if seen[pl.measure] {
			continue
		}
		seen[pl.measure] = true
		if t := r.ms[pl.measure].m.Time; t != nil {
			perf.Time = append(perf.Time, music.TimeChange{Tick: ticks[k], Time: *t})
		}
		if keyPart != nil && pl.measure < len(keyPart.Measures) {
			if ks := keyPart.Measures[pl.measure].Key; ks != nil {
				perf.Key = append(perf.Key, music.KeyChange{Tick: ticks[k], Key: *ks})
			}
		}
	}

	r.tracks(perf, order, ticks)
}

// tempos returns the tempo changes of the measures played from ticks: the
// changes in score order, restated where a jump lands.
func (r *reader) tempos(order []played, ticks []int) []music.Tempo {
	var out []music.Tempo
	points := make([][]tempoPoint, len(r.ms))
	startTempo := make([]float64, len(r.ms))
	cur := 120.0
	for i, mi := range r.ms {
		startTempo[i] = cur
		pts := append([]tempoPoint(nil), mi.tempo...)
		sort.SliceStable(pts, func(a, b int) bool {
			if pts[a].offset != pts[b].offset {
				return pts[a].offset < pts[b].offset
			}
			return !pts[a].metronome && pts[b].metronome
		})
		for _, p := range pts {
			if len(points[i]) > 0 && points[i][len(points[i])-1].offset == p.offset {
				continue
			}
			points[i] = append(points[i], p)
			cur = p.bpm
		}
	}
	// a measure starts at the tempo of a change at its start, and changes
	// it where a change differs from the tempo before it: the same every
	// time the measure is played
	changes := make([][]tempoPoint, len(r.ms))
	for i := range r.ms {
		t := startTempo[i]
		for _, p := range points[i] {
			switch {
			case p.offset <= 0:
				t = p.bpm
				startTempo[i] = t
			case p.bpm != t:
				t = p.bpm
				changes[i] = append(changes[i], p)
			}
		}
	}
	last := -1.0
tempo:
	for k, pl := range order {
		if t0 := startTempo[pl.measure]; t0 != last {
			out = append(out, music.Tempo{Tick: ticks[k], BPM: t0})
			last = t0
		}
		for _, p := range changes[pl.measure] {
			if len(out) >= maxPlayChanges {
				r.warn("more than %d tempo changes are played; the rest are not", maxPlayChanges)
				break tempo
			}
			out = append(out, music.Tempo{Tick: ticks[k] + p.offset, BPM: p.bpm})
			last = p.bpm
		}
	}
	if len(out) == 0 {
		out = []music.Tempo{{Tick: 0, BPM: 120}}
	}
	return out
}

// tracks makes a track of each part: its channel and program, and its
// notes and lyrics as played.
func (r *reader) tracks(perf *music.Performance, order []played, ticks []int) {
	var parts []*partState
	for _, p := range r.parts {
		if p.seen || len(r.score.Parts) == 1 && r.score.Parts[0] == p.part {
			parts = append(parts, p)
		}
	}
	// channels: the instrument's, else percussion on 9 and the others on
	// the free ones
	channel := make([]int, len(parts))
	used := map[int]bool{}
	for k, p := range parts {
		channel[k] = -1
		if len(p.instruments) > 0 && p.instruments[0].channel >= 0 {
			channel[k] = p.instruments[0].channel
			used[channel[k]] = true
		}
	}
	free := func(k int) int {
		for c := 0; c < 16; c++ {
			if c != 9 && !used[c] {
				return c
			}
		}
		c := k % 15 // more parts than channels: they share
		if c >= 9 {
			c++
		}
		return c
	}
	for k, p := range parts {
		switch {
		case channel[k] >= 0:
		case p.part.Percussion || p.unpitched:
			channel[k] = 9
		default:
			channel[k] = free(k)
			used[channel[k]] = true
		}
	}

	total, looked, sung := 0, 0, 0
	for k, p := range parts {
		t := &music.Track{Name: p.trackName, Channel: channel[k], Volume: -1, Pan: -1}
		if t.Name == "" {
			t.Name = p.part.Name
		}
		if len(p.instruments) > 0 {
			in := p.instruments[0]
			t.Program = max(in.program, 0)
			t.Volume, t.Pan = in.volume, in.pan
		}
		perf.Tracks = append(perf.Tracks, t)

		// the dynamics in effect at the start of each measure, and the
		// marks in it
		startVel := make([]int, len(p.play))
		marks := make([]*dynamics, len(p.play))
		vel, d := 80, 0
		for i := range p.play {
			startVel[i] = vel
			for ; d < len(p.dyn) && p.dyn[d].measure <= i; d++ {
				if p.dyn[d].measure == i {
					if marks[i] == nil {
						marks[i] = &dynamics{}
					}
					marks[i].add(p.dyn[d])
				}
				if p.dyn[d].vel > 0 {
					vel = p.dyn[d].vel
				}
			}
		}

		type open struct{ idx, end int }
		ties := map[int]*open{}
		for k, pl := range order {
			if pl.measure >= len(p.play) {
				continue
			}
			pm := p.play[pl.measure]
			base := ticks[k]
			dyn := marks[pl.measure]
			for _, n := range pm.notes {
				// the notes looked at count too: tied notes and notes out
				// of range add none to the track
				if total >= maxPlayNotes || looked >= 4*maxPlayNotes {
					r.warn("more than %d notes are played; the rest are not", maxPlayNotes)
					break
				}
				looked++
				tick := base + n.offset
				if n.tieStop {
					if o := ties[n.key]; o != nil && abs(o.end-tick) <= 2 {
						t.Notes[o.idx].Dur = tick + n.dur - t.Notes[o.idx].Tick
						o.end = tick + n.dur
						if !n.tieStart {
							delete(ties, n.key)
						}
						continue
					}
				}
				v, accent := startVel[pl.measure], n.ev.accent
				level, a := dyn.at(n.offset)
				if level > 0 {
					v = level
				}
				accent += a
				if n.vel > 0 {
					v = n.vel
				}
				dur := n.dur
				if !n.tieStart && n.ev.shorten > 1 {
					dur = max(dur/n.ev.shorten, 1)
				}
				if dur <= 0 || n.key < 0 || n.key > 127 {
					continue
				}
				t.Notes = append(t.Notes, music.PlayNote{Tick: tick, Dur: dur, Key: n.key, Vel: min(max(v+accent, 1), 127)})
				total++
				if n.tieStart {
					ties[n.key] = &open{len(t.Notes) - 1, tick + n.dur}
				} else {
					delete(ties, n.key)
				}
			}
			// verse k on the k-th pass when the measure has it
			if pm.verses == nil && len(pm.lyrics) > 0 {
				pm.verses = map[int][]playLyric{}
				for _, l := range pm.lyrics {
					pm.verses[l.verse] = append(pm.verses[l.verse], l)
				}
			}
			verse := 1
			if len(pm.verses[pl.pass]) > 0 {
				verse = pl.pass
			}
			for _, l := range pm.verses[verse] {
				if sung >= maxPlayChanges {
					r.warn("more than %d syllables are sung; the rest are not", maxPlayChanges)
					break
				}
				t.Lyrics = append(t.Lyrics, music.Lyric{Tick: base + l.offset, Text: l.text})
				sung++
			}
		}
		sort.SliceStable(t.Notes, func(a, b int) bool { return t.Notes[a].Tick < t.Notes[b].Tick })
		sort.SliceStable(t.Lyrics, func(a, b int) bool { return t.Lyrics[a].Tick < t.Lyrics[b].Tick })
		for _, n := range t.Notes {
			perf.End = max(perf.End, n.Tick+n.Dur)
		}
	}
}

// dynamics are the dynamics marks of a part in a measure, in the order of
// their offsets.
type dynamics struct {
	offsets []int
	levels  []int // the dynamics set by the marks up to each one; 0 when none sets them
	accents []int // the accents of the marks up to each one, added up
}

// add adds the next mark of the measure.
func (d *dynamics) add(e dynPoint) {
	level, accents := 0, 0
	if n := len(d.offsets); n > 0 {
		level, accents = d.levels[n-1], d.accents[n-1]
	}
	if e.vel > 0 {
		level = e.vel
	}
	d.offsets = append(d.offsets, e.offset)
	d.levels = append(d.levels, level)
	d.accents = append(d.accents, accents+e.accent)
}

// at returns the dynamics set by the marks up to an offset (0 when none
// sets them) and the accents of the marks at the offset. d may be nil.
func (d *dynamics) at(offset int) (level, accent int) {
	if d == nil {
		return 0, 0
	}
	before := sort.SearchInts(d.offsets, offset) // the marks before the offset
	upTo := sort.SearchInts(d.offsets, offset+1) // and with those at it
	if upTo == 0 {
		return 0, 0
	}
	accent = d.accents[upTo-1]
	if before > 0 {
		accent -= d.accents[before-1]
	}
	return d.levels[upTo-1], accent
}
