package web

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// What gen-project writes is examples/docflow-custom, file for file, with
// the example's own name and module path: the example is the scaffold's
// output, kept and tested in the repository (its go.mod has what go mod
// tidy added, and its Go version is the module's).
func TestGenProjectIsTheExample(t *testing.T) {
	dir := t.TempDir()
	var out bytes.Buffer
	if err := GenerateProject(GenProject{Dir: dir, Module: "github.com/shibukawa/bdf/examples/docflow-custom", Name: "docflow-custom", Bdf: "../.."}, &out); err != nil {
		t.Fatal(err)
	}
	example := filepath.Join("..", "..", "examples", "docflow-custom")
	written := 0
	err := filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(dir, p)
		got, _ := os.ReadFile(p)
		want, err := os.ReadFile(filepath.Join(example, rel))
		if err != nil {
			t.Errorf("%s: written, but not in the example: %v", rel, err)
			return nil
		}
		written++
		if rel == "go.mod" {
			// the module's lines are the example's, which has more
			// (what go mod tidy added) and the go version of the module
			for _, line := range strings.Split(string(got), "\n") {
				if strings.HasPrefix(line, "go ") || strings.TrimSpace(line) == "" {
					continue
				}
				if !strings.Contains(string(want), line) {
					t.Errorf("go.mod: %q is not in the example's", line)
				}
			}
			return nil
		}
		if !bytes.Equal(got, want) {
			t.Errorf("%s differs from the example's (docflow gen-project --force writes the example again)", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if written < 10 {
		t.Errorf("%d files written", written)
	}
	// the example has nothing the scaffold does not write, but what the
	// tools made: go.sum, the generated manifest, what pw built (dist),
	// pw dev's logs and store, the binary
	filepath.WalkDir(example, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(example, p)
		if d.IsDir() {
			if rel == "dist" || rel == "bin" || rel == ".log" {
				return filepath.SkipDir
			}
			return nil
		}
		if rel == "go.sum" || strings.HasSuffix(rel, "_pw_gen.go") || rel == "docflow.db" {
			return nil
		}
		if _, err := os.Stat(filepath.Join(dir, rel)); err != nil {
			t.Errorf("%s is in the example, but gen-project does not write it", rel)
		}
		return nil
	})
	for _, s := range []string{"wrote " + filepath.Join(dir, "main.go"), "go mod tidy", "pw dev"} {
		if !strings.Contains(out.String(), s) {
			t.Errorf("the report lacks %q:\n%s", s, out.String())
		}
	}
}

// A file that is there is not written over without --force; the module
// path and the name have defaults, and bad ones are refused.
func TestGenProjectChecks(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "acme-docflow")
	var out bytes.Buffer
	if err := GenerateProject(GenProject{Dir: dir}, &out); err != nil {
		t.Fatal(err)
	}
	gomod, _ := os.ReadFile(filepath.Join(dir, "go.mod"))
	toml, _ := os.ReadFile(filepath.Join(dir, "popcornweb.toml"))
	if !strings.HasPrefix(string(gomod), "module acme-docflow\n") || strings.Contains(string(gomod), "replace ") || !strings.Contains(string(toml), `name = "acme-docflow"`) {
		t.Errorf("defaults:\n%s\n%s", gomod, toml)
	}
	if !strings.Contains(out.String(), "v0.0.0") {
		t.Errorf("a checkout's docflow does not say so:\n%s", out.String())
	}
	if err := GenerateProject(GenProject{Dir: dir}, &out); err == nil || !strings.Contains(err.Error(), "--force") {
		t.Errorf("written over: %v", err)
	}
	if err := GenerateProject(GenProject{Dir: dir, Force: true, Module: "example.com/acme/flow", Bdf: "../bdf"}, &out); err != nil {
		t.Fatal(err)
	}
	gomod, _ = os.ReadFile(filepath.Join(dir, "go.mod"))
	main, _ := os.ReadFile(filepath.Join(dir, "main.go"))
	if !strings.HasPrefix(string(gomod), "module example.com/acme/flow\n") || !strings.Contains(string(gomod), "replace github.com/shibukawa/bdf => ../bdf\n") || !strings.Contains(string(main), `"example.com/acme/flow/kinds/folder"`) || !strings.Contains(string(main), "Command flow ") {
		t.Errorf("--force with a module:\n%s\n%s", gomod, main)
	}
	for _, c := range []GenProject{{Dir: ""}, {Dir: dir, Force: true, Module: "a b"}, {Dir: dir, Force: true, Name: "no spaces"}, {Dir: dir, Force: true, Module: "../x"}} {
		if err := GenerateProject(c, &out); err == nil {
			t.Errorf("%+v: written", c)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "kinds", "folder", "folder.go")); err != nil {
		t.Error(err)
	}
}
