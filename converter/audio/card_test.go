package audio

import (
	"slices"
	"strings"
	"testing"
)

func TestWrap(t *testing.T) {
	for _, c := range []struct {
		s     string
		width float32
		max   int
		want  []string
	}{
		{"MP3, 33 kbps (variable), 22.05 kHz, mono", 304, 12, []string{"MP3, 33 kbps (variable), 22.05 kHz, mono"}},
		{"Over the field the sun is high", 60, 8, []string{"Over the", "field the", "sun is high"}},
		// a word longer than a line is cut
		{"abcdefghijklmnopqrstuvwxyz", 40, 8, []string{"abcdef", "ghijkl", "mnopqr", "stuvwx", "yz"}},
		// wide characters break anywhere
		{"野原の上に太陽", 40, 8, []string{"野原の", "上に太", "陽"}},
		// more than max lines: an ellipsis
		{"one two three four five six", 44, 2, []string{"one two", "thre…"}},
		{"", 100, 3, nil},
	} {
		got := wrap(c.s, 11, c.width, c.max)
		if !slices.Equal(got, c.want) {
			t.Errorf("wrap(%q, %v, %d) = %q, want %q", c.s, c.width, c.max, got, c.want)
		}
	}
	// a line is never wider than estimated
	for _, l := range wrap(strings.Repeat("word ", 100), 11, 200, 100) {
		if w := estWidth(l, 11); w > 200 {
			t.Errorf("line %q is %v wide", l, w)
		}
	}
}

func TestClockAndRates(t *testing.T) {
	for sec, want := range map[float64]string{0: "0:00", 59.4: "0:59", 61: "1:01", 3600: "1:00:00", 3725.6: "1:02:06"} {
		if got := clock(sec); got != want {
			t.Errorf("clock(%v) = %q, want %q", sec, got, want)
		}
	}
	for rate, want := range map[int]string{44100: "44.1 kHz", 48000: "48 kHz", 22050: "22.05 kHz", 8000: "8 kHz", 11025: "11.025 kHz"} {
		if got := kHz(rate); got != want {
			t.Errorf("kHz(%d) = %q, want %q", rate, got, want)
		}
	}
	tr := &track{format: "flac", codec: "FLAC", bits: 16, bitrate: 900_000, sampleRate: 44100, channels: 2}
	if got := formatLine(tr); got != "FLAC, 16-bit, 900 kbps, 44.1 kHz, stereo" {
		t.Errorf("formatLine = %q", got)
	}
	tr = &track{format: "mp3", codec: "MP3", bitrate: 128_000, vbr: true, sampleRate: 44100, channels: 6, protected: true}
	if got := formatLine(tr); got != "MP3 (protected), 128 kbps (variable), 44.1 kHz, 6 channels" {
		t.Errorf("formatLine = %q", got)
	}
}
