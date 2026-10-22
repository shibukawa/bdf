package main

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// The embedded tree is pw generate's build of public: the same files, so a
// change to public without pw generate is caught here as well as by CI.
func TestPublicTree(t *testing.T) {
	tree := publicTree()
	for _, name := range []string{"app.css", "tablet.js", "admin.js"} {
		built, err := fs.ReadFile(tree, name)
		if err != nil {
			t.Fatalf("%s is not in the embedded tree: %v", name, err)
		}
		authored, err := os.ReadFile(filepath.Join("public", name))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(built, authored) {
			t.Errorf("%s: the embedded tree is not public's: run pw generate", name)
		}
		for _, sidecar := range []string{".br", ".gz", ".zstd"} {
			if _, err := fs.Stat(tree, name+sidecar); err != nil {
				t.Errorf("%s%s: %v", name, sidecar, err)
			}
		}
	}
}
