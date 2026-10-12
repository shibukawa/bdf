package conformance

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestBuild builds the suite with the test documents, which checks every
// synthetic case against what it was written to show, writes it out and
// reads it back.
func TestBuild(t *testing.T) {
	s, err := Build("../testdata")
	if err != nil {
		t.Fatal(err)
	}
	kinds := map[string]int{}
	verdicts := map[string]int{}
	for _, c := range s.Cases {
		kinds[c.Kind]++
		verdicts[c.Expect.Verdict()]++
		if c.Expect.Verdict() == "" {
			t.Errorf("%s: no verdict", c.ID)
		}
		if !strings.HasPrefix(c.Input, "inputs/") {
			t.Errorf("%s: input %s", c.ID, c.Input)
		}
	}
	t.Logf("%d cases: %v, %v", len(s.Cases), kinds, verdicts)
	if kinds["single"] < 50 || kinds["object"] < 50 || kinds["paths"] < 3 || verdicts[OK] < 100 || verdicts[Err] < 50 {
		t.Errorf("the suite is thin: %v %v", kinds, verdicts)
	}

	dir := t.TempDir()
	if err := s.Write(dir); err != nil {
		t.Fatal(err)
	}
	back, err := Read(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(back.Cases) != len(s.Cases) {
		t.Fatalf("%d cases read back, %d written", len(back.Cases), len(s.Cases))
	}
	// what was observed holds for what was read back
	for _, c := range back.Cases {
		e, err := Observe(c.Kind, back.Input(c))
		if err != nil {
			t.Fatal(err)
		}
		if e.Verdict() != c.Expect.Verdict() || e.Walk != c.Expect.Walk || len(e.Views) != len(c.Expect.Views) {
			t.Errorf("%s: %+v read back, %+v observed", c.ID, c.Expect, e)
		}
	}
	if _, err := Read(filepath.Join(dir, "nowhere")); err == nil {
		t.Error("a missing suite read")
	}
}

// TestBuildWithoutTestdata builds the synthetic cases alone.
func TestBuildWithoutTestdata(t *testing.T) {
	s, err := Build("")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range s.Cases {
		if strings.HasPrefix(c.ID, "testdata/") {
			t.Errorf("%s in a suite without testdata", c.ID)
		}
	}
}
