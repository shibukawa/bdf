package music

import (
	"math"
	"slices"
	"sort"
)

// Standard tuning, from the first (highest) string to the sixth.
var standardGuitar = [6]int{64, 59, 55, 50, 45, 40}

type tabPreset struct {
	id, name, explanation  string
	low, shift, open, high float64
	spanLimit              int
}

// Each view is an alternative fingering, not a claim about how the source
// performance was actually played.
var tabPresets = []tabPreset{
	{"low", "Low positions", "Favors notes near the nut", 1.8, 0.65, 0.2, 0.6, 5},
	{"smooth", "Less shifting", "Favors a stable hand position across the phrase", 0.25, 4.0, 0.3, 0.25, 5},
	{"open", "Open strings", "Favors ringing open strings", 0.4, 1.0, 5.0, 0.35, 5},
	{"easy", "Beginner friendly", "Favors low frets, short stretches and open strings", 2.0, 2.0, 2.0, 2.0, 4},
}

type tabInput struct {
	note  *Note
	pitch int
	end   int // absolute tick at which this string can be played again
}

type tabOnset struct {
	tick  int
	notes []tabInput
}

type tabCandidate struct {
	pos     TabPosition
	penalty float64 // leaving an authored string/fret is a last resort
}

type tabShape struct {
	positions []TabPosition
	cost      float64
	base      int
}

type tabState struct {
	cost  float64
	prev  int
	shape int
	held  [6]int // end tick of the note sounding on each string
	frets [6]int // frets of the held notes (for physical span checks)
}

// tabCandidates returns playable positions for a sounding MIDI pitch.
// An authored position is used alone when possible, even when the written
// pitch is octave transposed. A relaxed retry adds alternatives for conflicts.
func tabCandidates(in tabInput, relaxed bool) []tabCandidate {
	var out []tabCandidate
	p := in.note.Tab
	authored := p != nil && p.String >= 1 && p.String <= 6 && p.Fret >= 0 && p.Fret <= 24
	if authored {
		out = append(out, tabCandidate{pos: *p})
		if !relaxed {
			return out
		}
	}
	for i, open := range standardGuitar {
		fret := in.pitch - open
		if fret < 0 || fret > 24 {
			continue
		}
		pos := TabPosition{String: i + 1, Fret: fret}
		if authored && pos == *p {
			continue
		}
		penalty := 0.0
		if authored {
			penalty = 100
		}
		out = append(out, tabCandidate{pos: pos, penalty: penalty})
	}
	return out
}

// tabOnsets collects all voices and staves of one part by absolute onset,
// so simultaneous notes are assigned as a chord rather than independently.
func tabOnsets(s *Score, part *Part) []tabOnset {
	byTick := map[int][]tabInput{}
	start := 0
	for mi, m := range s.Measures {
		if mi < len(part.Measures) && part.Measures[mi] != nil {
			for _, ev := range part.Measures[mi].Events {
				if ev == nil || ev.Hidden || ev.Rest {
					continue
				}
				for _, note := range ev.Notes {
					if note != nil {
						tick := start + ev.Offset
						pitch := note.Pitch.MIDI()
						if note.SoundingMIDI != nil {
							pitch = *note.SoundingMIDI
						}
						byTick[tick] = append(byTick[tick], tabInput{note: note, pitch: pitch, end: tick + max(ev.Duration, 1)})
					}
				}
			}
		}
		start += m.Length
	}
	ticks := make([]int, 0, len(byTick))
	for tick := range byTick {
		ticks = append(ticks, tick)
	}
	slices.Sort(ticks)
	out := make([]tabOnset, 0, len(ticks))
	for _, tick := range ticks {
		out = append(out, tabOnset{tick, byTick[tick]})
	}
	return out
}

// tabShapes enumerates chord positions with one note per string and a hand
// span of at most four or five frets (open strings do not stretch the hand).
func tabShapes(notes []tabInput, preset tabPreset, relaxed bool) []tabShape {
	if len(notes) == 0 || len(notes) > 6 {
		return nil
	}
	type choice struct {
		index int
		list  []tabCandidate
	}
	choices := make([]choice, len(notes))
	for i, note := range notes {
		choices[i] = choice{i, tabCandidates(note, relaxed)}
		if len(choices[i].list) == 0 {
			return nil
		}
	}
	sort.SliceStable(choices, func(i, j int) bool { return len(choices[i].list) < len(choices[j].list) })
	positions := make([]TabPosition, len(notes))
	var shapes []tabShape
	var visit func(int, uint8, int, int, float64)
	visit = func(index int, used uint8, minFret, maxFret int, penalty float64) {
		if index == len(choices) {
			base := 0
			if minFret <= 24 {
				base = minFret
			}
			cost := penalty + preset.low*float64(base) + preset.high*math.Pow(float64(max(0, maxFret-12)), 2)
			for _, pos := range positions {
				if pos.Fret == 0 {
					cost -= preset.open
				}
			}
			shapes = append(shapes, tabShape{slices.Clone(positions), cost, base})
			if len(shapes) >= 128 {
				sort.SliceStable(shapes, func(i, j int) bool { return shapes[i].cost < shapes[j].cost })
				shapes = shapes[:64]
			}
			return
		}
		for _, candidate := range choices[index].list {
			p := candidate.pos
			bit := uint8(1 << (p.String - 1))
			if used&bit != 0 {
				continue
			}
			lo, hi := minFret, maxFret
			if p.Fret > 0 {
				lo, hi = min(lo, p.Fret), max(hi, p.Fret)
			}
			if lo <= 24 && hi-lo > preset.spanLimit {
				continue
			}
			positions[choices[index].index] = p
			visit(index+1, used|bit, lo, hi, penalty+candidate.penalty)
		}
	}
	visit(0, 0, 25, 0, 0)
	sort.SliceStable(shapes, func(i, j int) bool {
		if shapes[i].cost != shapes[j].cost {
			return shapes[i].cost < shapes[j].cost
		}
		for k := range shapes[i].positions {
			a, b := shapes[i].positions[k], shapes[j].positions[k]
			if a.String != b.String {
				return a.String < b.String
			}
			if a.Fret != b.Fret {
				return a.Fret < b.Fret
			}
		}
		return false
	})
	if len(shapes) > 64 {
		shapes = shapes[:64]
	}
	return shapes
}

func tabTransition(a, b tabShape, preset tabPreset) float64 {
	move := math.Abs(float64(a.base - b.base))
	// The current shape may be an inversion of the previous one. Compare
	// each string to its nearest predecessor rather than pairing notes by
	// their input order.
	for _, current := range b.positions {
		nearest := math.MaxFloat64
		for _, previous := range a.positions {
			d := 0.22*math.Abs(float64(current.Fret-previous.Fret)) +
				0.35*math.Abs(float64(current.String-previous.String))
			nearest = min(nearest, d)
		}
		move += nearest / float64(len(b.positions))
	}
	return preset.shift * move
}

// chooseTab uses bounded dynamic programming over chord shapes for the whole part.
// It reports notes that cannot be played in standard tuning or an onset with
// more than six simultaneous pitches; those are left blank in the TAB.
func chooseTab(onsets []tabOnset, preset tabPreset) (map[*Note]TabPosition, int) {
	positions := map[*Note]TabPosition{}
	type step struct {
		notes  []tabInput
		shapes []tabShape
		states []tabState
	}
	var steps []step
	omitted := 0
	for _, onset := range onsets {
		var notes []tabInput
		for _, note := range onset.notes {
			if len(tabCandidates(note, true)) > 0 {
				notes = append(notes, note)
			} else {
				omitted++
			}
		}
		advance := func(shapes []tabShape) []tabState {
			var states []tabState
			for i, shape := range shapes {
				if len(steps) == 0 {
					state := tabState{cost: shape.cost, prev: -1, shape: i}
					if state.place(notes, shape, onset.tick, preset.spanLimit) {
						states = append(states, state)
					}
					continue
				}
				last := &steps[len(steps)-1]
				for j, prev := range last.states {
					state := tabState{cost: prev.cost + shape.cost + tabTransition(last.shapes[prev.shape], shape, preset),
						prev: j, shape: i, held: prev.held, frets: prev.frets}
					if state.place(notes, shape, onset.tick, preset.spanLimit) {
						states = append(states, state)
					}
				}
			}
			return states
		}
		shapes := tabShapes(notes, preset, false)
		states := advance(shapes)
		if len(states) == 0 {
			// Only a physical conflict permits overriding an authored position.
			shapes = tabShapes(notes, preset, true)
			states = advance(shapes)
		}
		if len(states) == 0 {
			omitted += len(notes)
			continue
		}
		sort.SliceStable(states, func(i, j int) bool { return states[i].cost < states[j].cost })
		if len(states) > 32 {
			states = states[:32]
		}
		steps = append(steps, step{notes, shapes, states})
	}
	if len(steps) == 0 {
		return positions, omitted
	}
	best := 0 // every step's states are sorted by cumulative cost
	for i := len(steps) - 1; i >= 0; i-- {
		st := steps[i]
		state := st.states[best]
		for j, note := range st.notes {
			positions[note.note] = st.shapes[state.shape].positions[j]
		}
		best = state.prev
	}
	return positions, omitted
}

// place checks strings held by earlier, overlapping notes as well as the
// new chord. A string can be reused exactly when its previous note ends.
func (s *tabState) place(notes []tabInput, shape tabShape, tick, spanLimit int) bool {
	for i, end := range s.held {
		if end <= tick {
			s.held[i], s.frets[i] = 0, 0
		}
	}
	for i, pos := range shape.positions {
		stringIndex := pos.String - 1
		if s.held[stringIndex] > tick {
			return false
		}
		end := notes[i].end
		if end <= tick {
			end = tick + 1
		}
		s.held[stringIndex], s.frets[stringIndex] = end, pos.Fret
	}
	minFret, maxFret := 25, 0
	for i, end := range s.held {
		if end > tick && s.frets[i] > 0 {
			minFret, maxFret = min(minFret, s.frets[i]), max(maxFret, s.frets[i])
		}
	}
	return minFret == 25 || maxFret-minFret <= spanLimit
}
