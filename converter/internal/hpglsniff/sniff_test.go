package hpglsniff

import "testing"

func TestIs(t *testing.T) {
	for _, c := range []struct {
		name string
		data string
		want bool
	}{
		{"bare", "IN;SP1;PA0,0;PD1000,0,1000,1000;PU;", true},
		{"lower case, no semicolons", "in sp1 pu0,0 pd100,100 pu", true},
		{"lines", "IN;\r\nSP1;\r\nPU0,0;\r\nPD1000,0;\r\n", true},
		{"a label first", "IN;DT@;LBTitle@SP1;PA0,0;", true},
		{"device control", "\x1b.(\x1b.I81;;17:\x1b.N;19:IN;IP;SC0,100,0,100;SP1;", true},
		{"a comment and a picture name", `CO"made by hand";BP1,"sheet";PS1000,1000;`, true},
		{"cut short", "IN;SP1;PA10", true},
		{"PJL HP-GL/2", "\x1b%-12345X@PJL JOB\r\n@PJL ENTER LANGUAGE = HPGL2\r\nBP;IN;", true},
		{"PCL with HP-GL/2", "\x1bE\x1b&l2A\x1b%0BIN;SP1;", true},
		{"HP RTL raster", "\x1bE\x1b*v6W\x00\x03\x08\x08\x08\x08\x1b*r100S", true},
		{"PJL PostScript", "\x1b%-12345X@PJL ENTER LANGUAGE=POSTSCRIPT\r\n%!PS-Adobe-3.0", false},
		{"PCL text", "\x1bE\x1b&l2A\x1b(8U\x1b(s0p10h12v0s0b3THello, world\r\n\x0c\x1bE", false},
		{"English", "IN THE BEGINNING was the plot", false},
		{"instruction-like words", "INSPECT", false},
		{"CSV with semicolons", "id;name;price\n1;apple;120\n2;pear;90\n", false},
		{"Markdown", "# Plot\n\nIN;SP1 is HP-GL.", false},
		{"numbers", "1,2,3\n4,5,6\n", false},
	} {
		if got := Is([]byte(c.data)); got != c.want {
			t.Errorf("%s: %v, want %v", c.name, got, c.want)
		}
	}
}
