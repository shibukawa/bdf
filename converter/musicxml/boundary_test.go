package musicxml

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/shibukawa/bdf/converter/internal/music"
)

func TestNumberAssignments(t *testing.T) {
	ns := newNumbers(0)
	for _, c := range []struct {
		name string
		want int
	}{
		{"", 1}, {"1", 1}, {"alpha", 2}, {"99", 99}, {"100", 3},
		{"alpha", 2}, {"007", 7}, {"7", 7}, {"0", 4},
	} {
		if got := ns.get(c.name); got != c.want {
			t.Errorf("name %q is %d, want %d", c.name, got, c.want)
		}
	}
}

func TestManyNumberAssignments(t *testing.T) {
	ns := newNumbers(0)
	for i := 1; i <= 100000; i++ {
		if got := ns.get(fmt.Sprint("voice", i)); got != i {
			t.Fatalf("name %d is %d", i, got)
		}
	}
}

func TestGraceTimingManyNotes(t *testing.T) {
	var notes []*playNote
	for range 50000 {
		notes = append(notes, &playNote{voice: 1, key: 60, grace: true, ev: &playEvent{}})
	}
	notes = append(notes, &playNote{voice: 1, key: 62, dur: 4096 * music.PPQ, ev: &playEvent{}})
	out := graceTiming(notes)
	if len(out) != 50001 {
		t.Fatalf("%d notes, want 50001", len(out))
	}
	if out[1].offset != 39 || out[49999].offset != 49999*39 || out[50000].offset != 50000*39 {
		t.Errorf("grace offsets: second %d, last %d, principal %d", out[1].offset, out[49999].offset, out[50000].offset)
	}
}

func TestRepeatsManyMeasures(t *testing.T) {
	rd := newReader()
	for i := 0; i < 20000; i++ {
		mi := rd.measureAt(i)
		mi.right, mi.m.Length = true, 4*music.PPQ
	}
	rd.ms[0].endingStart, rd.ms[19999].endStop = &endingStart{number: "1-100"}, true
	rd.endings()
	order := rd.playOrder()
	if len(order) != maxPlays || order[99] != (played{0, 100}) || order[100] != (played{1, 1}) {
		t.Errorf("%d measures played", len(order))
	}
}

func TestDynamicsPlayed(t *testing.T) {
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
