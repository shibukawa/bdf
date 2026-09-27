package musicxml

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// FuzzParse feeds damaged scores to the reader: it may fail, but must not
// panic. The seeds are the test scores and some fragments.
func FuzzParse(f *testing.F) {
	files, _ := filepath.Glob("testdata/*.musicxml")
	mxl, _ := filepath.Glob("testdata/*.mxl")
	for _, name := range append(files, mxl...) {
		if b, err := os.ReadFile(name); err == nil && len(b) < 64<<10 {
			f.Add(b)
		}
	}
	for _, s := range []string{
		`<score-partwise><part-list><score-part id="P1"/></part-list><part id="P1"><measure><attributes><divisions>0</divisions></attributes><note><chord/><grace/><pitch><step>H</step><alter>1e308</alter><octave>99</octave></pitch><duration>-5</duration><time-modification><actual-notes>0</actual-notes></time-modification><notations><tuplet type="stop"/><tied type="stop"/></notations></note><backup><duration>1e300</duration></backup></measure></part></score-partwise>`,
		`<score-timewise><measure><part id="A"><note><rest measure="yes"/><duration>4</duration></note><barline><repeat direction="backward" times="1000000"/><ending number="1-99999" type="start"/></barline></part></measure><measure><part id="A"><sound dacapo="yes" dalsegno="x" tocoda="y" fine="yes"/><direction><direction-type><octave-shift type="down" size="99"/><wedge type="stop" number="-1"/></direction-type><offset>-99999</offset></direction></part></measure></score-timewise>`,
		`<score-partwise><part id="P"><measure><attributes><time><beats>3+x</beats><beat-type>0</beat-type></time><time><senza-misura/></time><clef number="99"><sign>TAB</sign></clef><staves>1000</staves></attributes><harmony><kind text="7"/><degree><degree-type>add</degree-type></degree></harmony></measure></part></score-partwise>`,
	} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		Detect(b[:min(len(b), 1024)], bytes.NewReader(b), int64(len(b)))
		Parse(bytes.NewReader(b), int64(len(b)), nil)
	})
}
