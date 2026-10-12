// Package conformance writes a machine-readable conformance suite: inputs
// and what a reader of BDF v1 must do with them, as this implementation
// does it. Other implementations run the suite without Go (libbdf's
// test_conformance does), and `bdf conformance` writes it for a release.
//
// The suite is a directory:
//
//	suite.json    the cases
//	inputs/...    the bytes of each case
//
// suite.json holds {"format": "bdf-conformance/1", "bdf": 1, "opset": 1,
// "generator": "...", "cases": [...]}. A case has an id, a kind, the input's
// path and an expectation:
//
//   - kind "single": a single-file document. expect.open is "ok" or "error";
//     for "ok", expect.locked says whether it is encrypted and locked,
//     expect.views lists {id, kind, pages} of the views, expect.parts counts
//     the parts, expect.reads lists parts to read with "read" "ok" or "error"
//     and the decoded "size", and expect.allParts true says every part reads.
//   - kind "object": an Object part. expect.decode is "ok" or "error"; for
//     "ok", expect.walk says whether every instruction reads, and
//     expect.instructions counts the instructions read before the end or
//     the error.
//   - kind "paths": a path collection part. expect.decode, and expect.paths
//     counts the paths.
//
// expect.message is this implementation's error, for people; readers need
// not reproduce it. Expectations are what this implementation does with the
// input, checked against what the case was written to show.
package conformance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/shibukawa/bdf"
)

// Format names the shape of suite.json.
const Format = "bdf-conformance/1"

// Suite is a set of cases and their inputs.
type Suite struct {
	Format    string `json:"format"`
	BDF       int    `json:"bdf"`
	Opset     int    `json:"opset"`
	Generator string `json:"generator"`
	Cases     []Case `json:"cases"`

	inputs map[string][]byte
}

// Case is one input and what a reader must do with it.
type Case struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	Input  string `json:"input"`
	Note   string `json:"note,omitempty"`
	Expect Expect `json:"expect"`
}

// Expect is what happens to a case's input.
type Expect struct {
	// single
	Open     string       `json:"open,omitempty"`
	Locked   bool         `json:"locked,omitempty"`
	Views    []ViewExpect `json:"views,omitempty"`
	Parts    *int         `json:"parts,omitempty"`
	Reads    []ReadExpect `json:"reads,omitempty"`
	AllParts bool         `json:"allParts,omitempty"`
	// object, paths
	Decode       string `json:"decode,omitempty"`
	Walk         string `json:"walk,omitempty"`
	Instructions *int   `json:"instructions,omitempty"`
	Paths        *int   `json:"paths,omitempty"`
	// for people
	Message string `json:"message,omitempty"`
}

// ViewExpect summarises a view.
type ViewExpect struct {
	ID    string `json:"id"`
	Kind  string `json:"kind"`
	Pages int    `json:"pages"`
}

// ReadExpect is the outcome of reading one part.
type ReadExpect struct {
	H    string `json:"h"`
	Read string `json:"read"`
	Size *int   `json:"size,omitempty"`
}

const (
	OK  = "ok"
	Err = "error"
)

func intp(n int) *int { return &n }

// Observe runs this implementation on an input of a kind and reports what
// it does: the expectation of a case.
func Observe(kind string, data []byte) (Expect, error) {
	switch kind {
	case "single":
		return observeSingle(data), nil
	case "object":
		return observeObject(data), nil
	case "paths":
		return observePaths(data), nil
	}
	return Expect{}, fmt.Errorf("conformance: unknown kind %q", kind)
}

func observeSingle(data []byte) Expect {
	r, err := bdf.OpenSingle(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return Expect{Open: Err, Message: err.Error()}
	}
	e := Expect{Open: OK, Locked: r.Locked()}
	for _, v := range r.Manifest.Views {
		e.Views = append(e.Views, ViewExpect{ID: v.ID, Kind: v.Kind, Pages: len(v.Pages)})
	}
	e.Parts = intp(len(r.Manifest.Parts))
	// a few parts are listed one by one; many are summed up
	listAll := len(r.Manifest.Parts) <= 8
	all := true
	for _, p := range r.Manifest.Parts {
		b, err := r.Part(p.H)
		switch {
		case err != nil:
			all = false
			e.Reads = append(e.Reads, ReadExpect{H: p.H.String(), Read: Err})
		case listAll:
			e.Reads = append(e.Reads, ReadExpect{H: p.H.String(), Read: OK, Size: intp(len(b))})
		}
	}
	e.AllParts = all && !listAll
	return e
}

func observeObject(data []byte) Expect {
	o, err := bdf.DecodeObject(data)
	if err != nil {
		return Expect{Decode: Err, Message: err.Error()}
	}
	ins, err := o.Instructions()
	e := Expect{Decode: OK, Walk: OK, Instructions: intp(len(ins))}
	if err != nil {
		e.Walk = Err
		e.Message = err.Error()
	}
	return e
}

func observePaths(data []byte) Expect {
	ps, err := bdf.DecodePathCollection(data)
	if err != nil {
		return Expect{Decode: Err, Message: err.Error()}
	}
	return Expect{Decode: OK, Paths: intp(len(ps))}
}

// Verdict is the one-word outcome of an expectation: open or decode.
func (e Expect) Verdict() string {
	if e.Open != "" {
		return e.Open
	}
	return e.Decode
}

// input is a case before its expectation is observed: what it shows, and
// the verdict it was written to have.
type input struct {
	id, kind, note string
	data           []byte
	want           string            // OK or Err: open or decode
	wantReads      map[string]string // part name → OK or Err, for single
	wantWalk       string            // for object: OK or Err, "" for either
}

// Build makes the suite: the synthetic cases, and the documents under
// testdataDir ("" for none) as accept cases.
func Build(testdataDir string) (*Suite, error) {
	s := &Suite{Format: Format, BDF: bdf.FormatVersion, Opset: bdf.OpsetVersion, Generator: "bdf-go", inputs: map[string][]byte{}}
	inputs := synthetic()
	if testdataDir != "" {
		docs, err := testdataInputs(testdataDir)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, docs...)
	}
	seen := map[string]bool{}
	for _, in := range inputs {
		if seen[in.id] {
			return nil, fmt.Errorf("conformance: two cases named %s", in.id)
		}
		seen[in.id] = true
		e, err := Observe(in.kind, in.data)
		if err != nil {
			return nil, err
		}
		if in.want != "" && e.Verdict() != in.want {
			return nil, fmt.Errorf("conformance: %s was written to %s but the reference says %s (%s)", in.id, in.want, e.Verdict(), e.Message)
		}
		if in.wantWalk != "" && e.Walk != in.wantWalk {
			return nil, fmt.Errorf("conformance: %s was written to walk %s but the reference says %s (%s)", in.id, in.wantWalk, e.Walk, e.Message)
		}
		for h, want := range in.wantReads {
			got := ""
			for _, rd := range e.Reads {
				if rd.H == h {
					got = rd.Read
				}
			}
			if got == "" && e.AllParts {
				got = OK
			}
			if got != want {
				return nil, fmt.Errorf("conformance: %s: part %s was written to read %s but the reference says %q", in.id, h, want, got)
			}
		}
		path := "inputs/" + strings.NewReplacer("/", "-", " ", "-").Replace(in.id) + ".bin"
		s.inputs[path] = in.data
		s.Cases = append(s.Cases, Case{ID: in.id, Kind: in.kind, Input: path, Note: in.note, Expect: e})
	}
	sort.SliceStable(s.Cases, func(i, j int) bool { return s.Cases[i].ID < s.Cases[j].ID })
	return s, nil
}

// Input returns the bytes of a case of a built suite.
func (s *Suite) Input(c Case) []byte { return s.inputs[c.Input] }

// Write writes the suite as a directory: suite.json and inputs/.
func (s *Suite) Write(dir string) error {
	if err := os.MkdirAll(filepath.Join(dir, "inputs"), 0o755); err != nil {
		return err
	}
	for path, data := range s.inputs {
		if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(path)), data, 0o644); err != nil {
			return err
		}
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "suite.json"), append(b, '\n'), 0o644)
}

// Read reads a suite directory back, for checking one.
func Read(dir string) (*Suite, error) {
	b, err := os.ReadFile(filepath.Join(dir, "suite.json"))
	if err != nil {
		return nil, err
	}
	var s Suite
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	if s.Format != Format {
		return nil, fmt.Errorf("conformance: %s is of format %q, not %q", dir, s.Format, Format)
	}
	s.inputs = map[string][]byte{}
	for _, c := range s.Cases {
		data, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(c.Input)))
		if err != nil {
			return nil, err
		}
		s.inputs[c.Input] = data
	}
	return &s, nil
}

// testdataInputs lists the single-file documents under dir as accept cases.
func testdataInputs(dir string) ([]input, error) {
	var files []string
	for _, pattern := range []string{"*.bdf", "*/*.bdf"} {
		m, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			return nil, err
		}
		files = append(files, m...)
	}
	sort.Strings(files)
	var out []input
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		rel, _ := filepath.Rel(dir, f)
		id := "testdata/" + strings.TrimSuffix(filepath.ToSlash(rel), ".bdf")
		out = append(out, input{id: id, kind: "single", data: data, want: OK, note: "a document a converter wrote"})
	}
	return out, nil
}
