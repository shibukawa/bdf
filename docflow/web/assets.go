package web

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// The pages' scripts and style are written in public (the style built
// there by Tailwind CSS with daisyUI, from tailwind/app.css, by go
// generate with the standalone CLI of tailwind/VERSIONS), embedded here
// and served at AssetsPath by this package itself, not at Popcorn Web's
// /public: that is the one tree of the program's own public assets (its
// manifest, which a program that composes Docflow builds of its own
// public directory; middlewares.RegisterPublicManifest keeps one), and
// Docflow's are the same in every program.

//go:generate tailwindcss --input tailwind/app.css --output public/app.css --minify

//go:embed public
var embedded embed.FS

// AssetsPath is where Docflow's scripts and style are served:
// /_docflow/app.css, /_docflow/topology.js ….
const AssetsPath = "/_docflow"

// assets serves the scripts and the style: the embedded copy, or a
// directory (Options.AssetsDir), each file with an ETag of its content
// and asked for again on each page (no-cache).
type assets struct {
	fsys fs.FS
	// etags are the embedded files', by name; a directory's are computed
	// as they are served (they change).
	etags map[string]string
}

func newAssets(dir string) (*assets, error) {
	if dir != "" {
		st, err := os.Stat(dir)
		if err != nil || !st.IsDir() {
			return nil, fmt.Errorf("docflow: assets dir %s is not a directory", dir)
		}
		return &assets{fsys: os.DirFS(dir)}, nil
	}
	sub, err := fs.Sub(embedded, "public")
	if err != nil {
		return nil, err
	}
	a := &assets{fsys: sub, etags: map[string]string{}}
	err = fs.WalkDir(sub, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := fs.ReadFile(sub, p)
		if err != nil {
			return err
		}
		a.etags[p] = etagOf(b)
		return nil
	})
	return a, err
}

func etagOf(b []byte) string {
	sum := sha256.Sum256(b)
	return `"` + hex.EncodeToString(sum[:8]) + `"`
}

func (a *assets) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := strings.TrimPrefix(r.URL.Path, AssetsPath+"/")
	if name == "" || name != path.Clean(name) || strings.HasPrefix(name, ".") || strings.Contains(name, "/.") || strings.HasPrefix(name, "/") {
		http.NotFound(w, r)
		return
	}
	b, err := fs.ReadFile(a.fsys, name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	tag := a.etags[name]
	if tag == "" {
		tag = etagOf(b)
	}
	w.Header().Set("ETag", tag)
	w.Header().Set("Cache-Control", "public, no-cache")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(b))
}
