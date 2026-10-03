package music

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf/raster/imagebdf"
)

// samples are performances engraved by TestSamples.
func samples() map[string]*Performance {
	out := map[string]*Performance{}

	// Twinkle, twinkle, little star, with its words
	tw := melody(60, q, 60, q, 67, q, 67, q, 69, q, 69, q, 67, h, 65, q, 65, q, 64, q, 64, q, 62, q, 62, q, 60, h,
		67, q, 67, q, 65, q, 65, q, 64, q, 64, q, 62, h, 67, q, 67, q, 65, q, 65, q, 64, q, 64, q, 62, h)
	words := []string{"Twin-", "kle", "twin-", "kle", "lit-", "tle", "star,", "how", "I", "won-", "der", "what", "you", "are."}
	for i, w := range words {
		tw.Tracks[0].Lyrics = append(tw.Tracks[0].Lyrics, Lyric{Tick: tw.Tracks[0].Notes[i].Tick, Text: w})
	}
	tw.Title = "Twinkle, Twinkle, Little Star"
	tw.Composer = "Traditional"
	tw.Tempo = []Tempo{{0, 100}}
	out["twinkle"] = tw

	// a piano piece: a melody over broken chords, in G
	var right, left []PlayNote
	mel := []int{67, 71, 74, 72, 71, 69, 67, 66, 67, 69, 71, 72, 74, 72, 71, 69}
	for i, k := range mel {
		right = append(right, PlayNote{Tick: i * e8, Dur: e8, Key: k + 12, Vel: 90})
	}
	right = append(right, PlayNote{Tick: 16 * e8, Dur: 3 * q, Key: 79, Vel: 90}, PlayNote{Tick: 16 * e8, Dur: 3 * q, Key: 83, Vel: 90})
	right = append(right, PlayNote{Tick: 16*e8 + 3*q, Dur: t3, Key: 81, Vel: 90}, PlayNote{Tick: 16*e8 + 3*q + t3, Dur: t3, Key: 79, Vel: 90},
		PlayNote{Tick: 16*e8 + 3*q + 2*t3, Dur: t3, Key: 78, Vel: 90})
	for i := 0; i < 5; i++ {
		base := []int{43, 50, 55, 50}
		if i%2 == 1 {
			base = []int{48, 52, 55, 52}
		}
		for j, k := range base {
			left = append(left, PlayNote{Tick: i*4*q + j*q, Dur: q, Key: k, Vel: 70})
		}
	}
	right = append(right, PlayNote{Tick: 5 * 4 * q, Dur: 4 * q, Key: 79, Vel: 80})
	left = append(left, PlayNote{Tick: 5 * 4 * q, Dur: 4 * q, Key: 43, Vel: 80}, PlayNote{Tick: 5 * 4 * q, Dur: 4 * q, Key: 31, Vel: 80})
	notes := append(append([]PlayNote{}, right...), left...)
	sortNotes(notes)
	out["piano"] = &Performance{Title: "Minuet in G", Composer: "after J. S. Bach", Tempo: []Tempo{{0, 96}},
		Tracks: []*Track{{Name: "Piano", Program: 0, Volume: -1, Pan: -1, Notes: notes}}}

	// three parts
	tr := func(name string, prog int, keys ...int) *Track {
		t := &Track{Name: name, Program: prog, Volume: -1, Pan: -1}
		for i := 0; i+1 < len(keys); i += 2 {
			tick := 0
			for _, n := range t.Notes {
				tick = max(tick, n.Tick+n.Dur)
			}
			if keys[i] > 0 {
				t.Notes = append(t.Notes, PlayNote{Tick: tick, Dur: keys[i+1], Key: keys[i], Vel: 90})
			} else {
				t.Notes = append(t.Notes, PlayNote{Tick: tick, Dur: keys[i+1], Key: 0, Vel: 0})
			}
		}
		// rests were written as key 0: drop them
		var kept []PlayNote
		for _, n := range t.Notes {
			if n.Key > 0 {
				kept = append(kept, n)
			}
		}
		t.Notes = kept
		return t
	}
	out["trio"] = &Performance{Title: "Ode to Joy", Composer: "L. van Beethoven", Tempo: []Tempo{{0, 120}},
		Time: []TimeChange{{0, TimeSig{4, 4, ""}}}, Key: []KeyChange{{0, KeySig{Fifths: 2}}},
		Tracks: []*Track{
			tr("Flute", 73, 78, q, 78, q, 79, q, 81, q, 81, q, 79, q, 78, q, 76, q, 74, q, 74, q, 76, q, 78, q, 78, 3*e8, 76, e8, 76, h),
			tr("Clarinet", 71, 74, h, 74, h, 71, h, 69, h, 69, h, 71, h, 74, 3*e8, 73, e8, 73, h),
			tr("Cello", 42, 50, q, 57, q, 50, q, 57, q, 47, q, 54, q, 45, q, 52, q, 43, q, 50, q, 45, q, 52, q, 45, q, 57, q, 45, h),
		}}
	return out
}

func sortNotes(ns []PlayNote) {
	for i := 1; i < len(ns); i++ {
		for j := i; j > 0 && ns[j].Tick < ns[j-1].Tick; j-- {
			ns[j], ns[j-1] = ns[j-1], ns[j]
		}
	}
}

// TestSamples engraves the samples; with MUSIC_SAMPLES=dir it writes their
// pages there as PNG images.
func TestSamples(t *testing.T) {
	dir := os.Getenv("MUSIC_SAMPLES")
	for name, p := range samples() {
		sc := Notate(p, NotateOptions{})
		res, err := Build(sc, Options{Source: "test", FontDirs: []string{"../../pptx/testdata/fonts"}, NoSystemFonts: true})
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if res.Pages == 0 || res.Notes == 0 {
			t.Errorf("%s: %s", name, res.Summary())
		}
		if dir == "" {
			continue
		}
		t.Logf("%s: %s %v", name, res.Summary(), res.Warnings)
		r := imagebdf.New(res.Doc, &imagebdf.Options{FontDirs: []string{"../../pptx/testdata/fonts"}, NoSystemFonts: true})
		v := res.Doc.Views[0]
		for i := range v.Pages {
			img, err := r.Page(v, i, 2)
			if err != nil {
				t.Fatal(err)
			}
			f, err := os.Create(filepath.Join(dir, name+"-"+string(rune('1'+i))+".png"))
			if err != nil {
				t.Fatal(err)
			}
			png.Encode(f, img)
			f.Close()
		}
	}
}
