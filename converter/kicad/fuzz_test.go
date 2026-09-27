package kicad

import (
	"bytes"
	"testing"
	"time"
)

// FuzzConvert converts arbitrary schematics and boards: each converts or
// fails, soon, without an internal error.
func FuzzConvert(f *testing.F) {
	for _, name := range []string{"demo/demo.kicad_sch", "demo/sub.kicad_sch", "demo/demo.kicad_pcb", "frame/frame.kicad_sch"} {
		f.Add(readTestdata(f, name))
	}
	f.Add([]byte(schematic(`(label "x" (at 0 0 0) (effects (font (size 1.27 1.27))))`)))
	f.Fuzz(func(t *testing.T, b []byte) {
		start := time.Now()
		opts := &Options{NoSystemFonts: true, NoTextIndex: true, Warn: func(m string) {
			if bytes.Contains([]byte(m), []byte("internal error")) {
				t.Error(m)
			}
		}}
		Convert(bytes.NewReader(b), int64(len(b)), opts)
		if d := time.Since(start); d > 2*time.Second {
			t.Errorf("%v to convert %d bytes", d, len(b))
		}
	})
}
