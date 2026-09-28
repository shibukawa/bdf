package musicxml

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/converter/internal/music"
)

// The reader keeps what it used to look up again for every note and every
// measure played. These tests compare it with the code it had, kept here.

// numberBefore gives a name its number from the numbers given so far.
func numberBefore(given map[string]int, name string) int {
	if name == "" {
		name = "1"
	}
	if v, ok := given[name]; ok {
		return v
	}
	v, err := strconv.Atoi(name)
	if err != nil || v < 1 || v > 99 {
		taken := map[int]bool{}
		for _, t := range given {
			taken[t] = true
		}
		for v = 1; taken[v]; v++ {
		}
	}
	given[name] = v
	return v
}

func TestNumbersSame(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	names := []string{"", "1", "2", "3", "5", "10", "98", "99", "100", "0", "-1", "007", "a", "b", "soprano", "part1verse1", "x y", "1.5"}
	for i := 0; i < 500; i++ {
		before, ns := map[string]int{}, newNumbers(0)
		for k := 0; k < 150; k++ {
			name := names[r.Intn(len(names))]
			if r.Intn(3) == 0 {
				name = fmt.Sprint("v", r.Intn(120))
			}
			if got, want := ns.get(name), numberBefore(before, name); got != want {
				t.Fatalf("run %d: %q is %d, want %d", i, name, got, want)
			}
		}
	}
}

func TestNumbersTime(t *testing.T) {
	// 100000 names: as many steps, where they were five billion
	ns := newNumbers(0)
	start := time.Now()
	for i := 1; i <= 100000; i++ {
		if got := ns.get(fmt.Sprint("voice", i)); got != i {
			t.Fatalf("name %d is %d", i, got)
		}
	}
	if d := time.Since(start); d > 2*time.Second {
		t.Errorf("took %v", d)
	}
}

// graceTimingBefore looks for the principal note of each group among all
// the notes, and for each grace note among those of its group.
func graceTimingBefore(notes []*playNote) []*playNote {
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

// playedNotes copies notes (their events are shared) and writes them out.
func playedNotes(notes []*playNote) ([]*playNote, string) {
	var b strings.Builder
	out := make([]*playNote, len(notes))
	for i, n := range notes {
		c := *n
		out[i] = &c
		fmt.Fprintf(&b, "%d+%d k%d v%d g%v; ", n.offset, n.dur, n.key, n.voice, n.grace)
	}
	return out, b.String()
}

func TestGraceTimingSame(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 1000; i++ {
		// notes in the order of the file: grace notes (chords of them, in
		// voices of their own) before principal notes, at a few offsets
		var notes []*playNote
		for k := r.Intn(30); k > 0; k-- {
			ev := &playEvent{shorten: 1}
			n := playNote{offset: r.Intn(5) * 240, voice: 1 + r.Intn(3), key: 60 + r.Intn(12), ev: ev, grace: r.Intn(2) == 0}
			if !n.grace {
				n.dur = r.Intn(8) * 120
			}
			for c := r.Intn(3); c >= 0; c-- {
				m := n
				if r.Intn(6) == 0 {
					m.voice = 1 + r.Intn(3)
				}
				notes = append(notes, &m)
			}
		}
		a, _ := playedNotes(notes)
		b, _ := playedNotes(notes)
		_, got := playedNotes(graceTiming(a))
		_, want := playedNotes(graceTimingBefore(b))
		if got != want {
			t.Fatalf("run %d:\n%s\nwant\n%s", i, got, want)
		}
	}
}

func TestGraceTimingTime(t *testing.T) {
	// 50000 grace notes before one note: as many steps, where they were
	// more than a billion
	var notes []*playNote
	for i := 0; i < 50000; i++ {
		notes = append(notes, &playNote{voice: 1, key: 60, grace: true, ev: &playEvent{}})
	}
	notes = append(notes, &playNote{voice: 1, key: 62, dur: 4096 * music.PPQ, ev: &playEvent{}})
	start := time.Now()
	out := graceTiming(notes)
	if d := time.Since(start); d > 2*time.Second {
		t.Errorf("took %v", d)
	}
	if len(out) != 50001 || out[1].offset != 39 || out[49999].offset != 49999*39 || out[50000].offset != 50000*39 {
		t.Errorf("%d notes, the last at %d", len(out), out[len(out)-1].offset)
	}
}

// read reads a document as far as Parse does before it makes the
// performance.
func read(t *testing.T, doc string) *reader {
	t.Helper()
	rd := newReader()
	data, _ := toUTF8([]byte(doc))
	if err := rd.read(newDecoder(bytes.NewReader(data))); err != nil {
		t.Fatal(err)
	}
	rd.finish()
	return rd
}

// repeatTimesBefore goes through the voltas around a measure.
func repeatTimesBefore(r *reader, i int) int {
	times := r.ms[i].m.Times
	if times <= 0 {
		times = 2
	}
	if r.ms[i].ending == nil {
		return times
	}
	lo, hi := i, i
	for lo > 0 && r.ms[lo-1].ending != nil {
		lo--
	}
	for hi+1 < len(r.ms) && r.ms[hi+1].ending != nil {
		hi++
	}
	for j := lo; j <= hi; j++ {
		for _, v := range r.ms[j].ending.numbers {
			times = max(times, v)
		}
	}
	return min(times, 100)
}

// temposBefore goes through the tempo changes of a measure every time it
// is played.
func temposBefore(r *reader, order []played, ticks []int) []music.Tempo {
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
	last := -1.0
	for k, pl := range order {
		t0 := startTempo[pl.measure]
		for _, p := range points[pl.measure] {
			if p.offset <= 0 {
				t0 = p.bpm
			}
		}
		if t0 != last {
			out = append(out, music.Tempo{Tick: ticks[k], BPM: t0})
			last = t0
		}
		for _, p := range points[pl.measure] {
			if p.offset > 0 && p.bpm != last {
				out = append(out, music.Tempo{Tick: ticks[k] + p.offset, BPM: p.bpm})
				last = p.bpm
			}
		}
	}
	if len(out) == 0 {
		out = []music.Tempo{{Tick: 0, BPM: 120}}
	}
	return out
}

// repeats makes a score of measures with repeats, voltas, jumps and tempo
// changes.
func repeats(r *rand.Rand) string {
	pick := func(s ...string) string { return s[r.Intn(len(s))] }
	ms := make([]string, 1+r.Intn(16))
	for i := range ms {
		var b strings.Builder
		if i == 0 {
			b.WriteString(attrs(4))
		}
		if r.Intn(5) == 0 {
			b.WriteString(`<barline location="left"><repeat direction="forward"/></barline>`)
		}
		if r.Intn(3) == 0 {
			fmt.Fprintf(&b, `<barline location="left"><ending number="%s" type="start"/></barline>`,
				pick("1", "2", "3", "1, 2", "1-3", "", "2-5", "7", "1 3", "100", "0", "x"))
		}
		for k := r.Intn(4); k > 0; k-- {
			switch r.Intn(3) {
			case 0:
				fmt.Fprintf(&b, `<sound tempo="%d"/>`, 60+20*r.Intn(3))
			case 1:
				fmt.Fprintf(&b, `<direction><direction-type><metronome><beat-unit>%s</beat-unit><per-minute>%d</per-minute></metronome></direction-type>%s</direction>`,
					pick("quarter", "half"), 60+20*r.Intn(3), pick("", "", `<sound tempo="100"/>`))
			}
			if r.Intn(2) == 0 {
				b.WriteString(n("C5", 1+r.Intn(8), ""))
			}
		}
		b.WriteString(n("C5", 4, "quarter"))
		if r.Intn(12) == 0 {
			fmt.Fprintf(&b, `<sound %s/>`, pick(`dacapo="yes"`, `dalsegno="s"`, `segno="s"`, `coda="c"`, `tocoda="c"`, `fine="yes"`))
		}
		if r.Intn(3) == 0 {
			fmt.Fprintf(&b, `<barline location="right"><ending number="1" type="%s"/></barline>`, pick("stop", "discontinue"))
		}
		if r.Intn(4) == 0 {
			fmt.Fprintf(&b, `<barline location="right"><repeat direction="backward"%s/></barline>`, pick("", "", ` times="3"`, ` times="1"`, ` times="200"`))
		}
		ms[i] = b.String()
	}
	return score(ms...)
}

func TestRepeatsSame(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	voltas, changes := 0, 0
	for i := 0; i < 1000; i++ {
		rd := read(t, repeats(r))
		for k, mi := range rd.ms {
			if got, want := rd.repeatTimes(k), repeatTimesBefore(rd, k); got != want {
				t.Fatalf("score %d measure %d: %d passes, want %d", i, k, got, want)
			}
			if mi.ending != nil {
				voltas++
			}
		}
		order := rd.playOrder()
		ticks := make([]int, len(order))
		tick := 0
		for k, pl := range order {
			ticks[k] = tick
			tick += rd.ms[pl.measure].m.Length
		}
		got, want := rd.tempos(order, ticks), temposBefore(rd, order, ticks)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("score %d: tempo %v, want %v", i, got, want)
		}
		changes += len(got)
	}
	if voltas < 1000 || changes < 3000 {
		t.Errorf("the scores have %d measures under voltas and %d tempo changes", voltas, changes)
	}
}

func TestRepeatsTime(t *testing.T) {
	// backward repeats in the 20000 measures of a volta of 100 numbers:
	// each of the 10000 measures played looked at all of them
	rd := newReader()
	for i := 0; i < 20000; i++ {
		mi := rd.measureAt(i)
		mi.right, mi.m.Length = true, 4*music.PPQ
	}
	rd.ms[0].endingStart, rd.ms[19999].endStop = &endingStart{number: "1-100"}, true
	rd.endings()
	start := time.Now()
	order := rd.playOrder()
	if d := time.Since(start); d > 2*time.Second {
		t.Errorf("took %v", d)
	}
	if len(order) != maxPlays || order[99] != (played{0, 100}) || order[100] != (played{1, 1}) {
		t.Errorf("%d measures played", len(order))
	}
}

// levelBefore goes through the marks of the measure for a note.
func levelBefore(dyn []dynPoint, start, offset int) (v, accent int) {
	v = start
	for _, e := range dyn {
		if e.offset > offset {
			break
		}
		if e.vel > 0 {
			v = e.vel
		}
		if e.offset == offset {
			accent += e.accent
		}
	}
	return v, accent
}

func TestDynamicsSame(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 1000; i++ {
		// the marks of a measure in the order of their offsets
		var dyn []dynPoint
		var marks *dynamics
		off := 0
		for k := r.Intn(12); k > 0; k-- {
			off += r.Intn(3) * 240
			e := dynPoint{offset: off, vel: r.Intn(3) * 40, accent: r.Intn(2) * 20}
			dyn = append(dyn, e)
			if marks == nil {
				marks = &dynamics{}
			}
			marks.add(e)
		}
		for offset := -240; offset <= off+240; offset += 120 {
			v, accent := 80, 0
			level, a := marks.at(offset)
			if level > 0 {
				v = level
			}
			accent += a
			if wv, wa := levelBefore(dyn, 80, offset); v != wv || accent != wa {
				t.Fatalf("run %d at %d: %d and %d, want %d and %d (marks %+v)", i, offset, v, accent, wv, wa, dyn)
			}
		}
	}
}

func TestDynamicsPlayed(t *testing.T) {
	// the marks at a note and before it, played twice
	dynamics := func(mark string) string {
		return `<direction><direction-type><dynamics><` + mark + `/></dynamics></direction-type></direction>`
	}
	m := attrs(1) + dynamics("p") + n("C5", 1, "quarter") + dynamics("sf") + dynamics("f") + n("D5", 1, "quarter") +
		n("E5", 1, "quarter") + dynamics("fp") + n("F5", 1, "quarter") + `<barline location="right"><repeat direction="backward"/></barline>`
	s, _ := parse(t, score(m))
	var vel []int
	for _, n := range s.Play.Tracks[0].Notes {
		vel = append(vel, n.Vel)
	}
	if want := []int{49, 116, 96, 69, 49, 116, 96, 69}; !reflect.DeepEqual(vel, want) {
		t.Errorf("played %v, want %v", vel, want)
	}
}
