package drawio

import (
	"strings"
	"testing"
)

// TestCompactPath checks that the d attribute the stencil generator writes
// expands into the step elements of the original path.
func TestCompactPath(t *testing.T) {
	root, err := parseStencilXML([]byte(`<shape w="10" h="10"><foreground><path rounded="1" d="M1 2L3.5 -4Q1 2 3 4C1 2 3 4 5 6A1 1 0 0 1 5 5Z"/><fillstroke/></foreground></shape>`))
	if err != nil {
		t.Fatal(err)
	}
	path := root.find("path")
	var got []string
	for _, k := range path.kids {
		var b strings.Builder
		b.WriteString(k.name)
		for _, a := range k.attrs {
			b.WriteString(" " + a.Name.Local + "=" + a.Value)
		}
		got = append(got, b.String())
	}
	want := []string{"move x=1 y=2", "line x=3.5 y=-4", "quad x1=1 y1=2 x2=3 y2=4", "curve x1=1 y1=2 x2=3 y2=4 x3=5 y3=6",
		"arc rx=1 ry=1 x-axis-rotation=0 large-arc-flag=0 sweep-flag=1 x=5 y=5", "close"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("steps\n got %q\nwant %q", got, want)
	}
	if path.get("rounded") != "1" {
		t.Error("path attributes are kept")
	}
}
