package hpgl

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// FuzzConvert feeds damaged plots to the converter: it may fail, but must
// not panic. The seeds are the test plots and some instructions.
func FuzzConvert(f *testing.F) {
	files, _ := filepath.Glob("testdata/*.plt")
	for _, name := range files {
		if b, err := os.ReadFile(name); err == nil && len(b) < 64<<10 {
			f.Add(b)
		}
	}
	for _, s := range []string{
		"IN;SP1;PM0;CI100;PM1;AA0,0,90;PM2;FT4,10,30;FP1;EP;",
		"IN;SC0,10,0,10,1,0,100;RO270;IW1,1,9,9;SP1;PR1,1;RT1,1,2,0;BR1,1,2,2,3,0;",
		"IN;SP1;DV1,1;LO19;ES1,-1;SL1;SI-1,-1;LBa\rb\nc\x08d\te\x0ef\x0fg\x03CP-2,3;SM#;PD1,1;",
		"\x1bE\x1b*v18W\x00\x01\x04\x08\x08\x08\x00\xff\x00\xff\x00\xff\x00\x00\x00\x00\x00\x00\x1b*r1A\x1b*b5M\x1b*b9W\x00\x00\x02\x12\x34\x05\x00\x03\x1b*rC",
	} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		Convert(bytes.NewReader(b), int64(len(b)), &Options{NoSystemFonts: true, NoTextIndex: true})
	})
}
