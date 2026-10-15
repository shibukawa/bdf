package htmlro

// Parse must build the tree golang.org/x/net/html.Parse builds, for every
// input, including the HTML5 test suites and malformed inputs.

import (
	"bytes"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

// compareParse checks that Parse gives the tree html.Parse gives for input,
// or fails as it fails.
func compareParse(t *testing.T, name string, input []byte, scripting bool) {
	t.Helper()
	want, wantErr := html.ParseWithOptions(bytes.NewReader(input), html.ParseOptionEnableScripting(scripting))
	for _, f := range feeds {
		var got *html.Node
		var err error
		opts := Options{NoScripting: !scripting}
		switch f.name {
		case "bytes":
			got, err = ParseBytes(input, opts)
		case "stream":
			got, err = Parse(bytes.NewReader(input), opts)
		default:
			opts.BufferSize = 1
			got, err = Parse(&oneByteReader{s: string(input)}, opts)
		}
		if (err != nil) != (wantErr != nil) {
			t.Errorf("%s (%s): error %v, x/net %v", name, f.name, err, wantErr)
			continue
		}
		if err != nil {
			continue
		}
		if err := checkTreeConsistency(got); err != nil {
			t.Errorf("%s (%s): %v", name, f.name, err)
		}
		g, err := dump(got)
		if err != nil {
			t.Fatal(err)
		}
		w, err := dump(want)
		if err != nil {
			t.Fatal(err)
		}
		if g != w {
			t.Errorf("%s (%s):\n%s\nx/net:\n%s", name, f.name, g, w)
		}
		// Render agrees too, which is what a consumer that writes the tree
		// back sees.
		var gr, wr bytes.Buffer
		html.Render(&gr, got)
		html.Render(&wr, want)
		if gr.String() != wr.String() {
			t.Errorf("%s (%s): renders differently:\n%s\nx/net:\n%s", name, f.name, gr.String(), wr.String())
		}
	}
}

func TestParseMatchesXNet(t *testing.T) {
	for _, tt := range tokenTests {
		compareParse(t, tt.desc, []byte(tt.html), true)
	}
	for i, in := range html5libTokenizerInputs(t) {
		compareParse(t, "tokenizer test "+strings.TrimSpace(strings.SplitN(in, "\n", 2)[0]), []byte(in), i%2 == 0)
	}
	parseTestCases(t, func(file string, i int, ta *testAttrs) {
		if ta.context != "" {
			return
		}
		compareParse(t, file+"/"+ta.text, []byte(ta.text), ta.scripting)
	})
}
