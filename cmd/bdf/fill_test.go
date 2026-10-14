package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf/form/tag"
)

func TestFormPictureRoot(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "pictures")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	inside := []byte("inside")
	if err := os.WriteFile(filepath.Join(dir, "inside.svg"), inside, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "secret.svg"), []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	for name, target := range map[string]string{
		"inside-link.svg": "inside.svg",
		"escape.svg":      "../secret.svg",
		"absolute.svg":    filepath.Join(base, "secret.svg"),
		"escape-dir":      base,
	} {
		if err := os.Symlink(target, filepath.Join(dir, name)); err != nil {
			t.Skipf("symlinks unavailable: %v", err)
		}
	}
	for _, explicit := range []bool{false, true} {
		t.Run(map[bool]string{false: "template directory", true: "files flag"}[explicit], func(t *testing.T) {
			f := addFormFlags(flag.NewFlagSet("fill", flag.ContinueOnError), "fill", "data", "data")
			path := filepath.Join(dir, "template.xlsx")
			if explicit {
				*f.files = dir
				path = filepath.Join(base, "template.xlsx")
			}
			opts, root, err := f.options(path, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer root.Close()
			for _, name := range []string{"inside.svg", "inside-link.svg"} {
				data, err := tag.PictureData(opts.Files, name)
				if err != nil || !bytes.Equal(data, inside) {
					t.Errorf("PictureData(%q) = %q, %v", name, data, err)
				}
			}
			for _, name := range []string{"escape.svg", "absolute.svg", "escape-dir/secret.svg"} {
				if _, err := tag.PictureData(opts.Files, name); err == nil {
					t.Errorf("PictureData(%q) escaped the picture directory", name)
				}
			}
			if err := root.Close(); err != nil {
				t.Fatal(err)
			}
			if _, err := tag.PictureData(opts.Files, "inside.svg"); err == nil {
				t.Error("picture filesystem remained open after closing its root")
			}
		})
	}
}
