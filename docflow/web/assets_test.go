package web

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// The embedded tree is public as it is written, served at AssetsPath
// with an ETag of each file.
func TestAssets(t *testing.T) {
	a, err := newAssets("")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"app.css", "tablet.js", "admin.js", "i18n.js", "ja.js", "topology.js", "canvas.js"} {
		if _, err := fs.Stat(a.fsys, name); err != nil {
			t.Errorf("%s is not in the embedded tree: %v", name, err)
		}
	}
	get := func(path string, header http.Header) *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodGet, path, nil)
		r.Header = header
		w := httptest.NewRecorder()
		a.ServeHTTP(w, r)
		return w
	}
	w := get(AssetsPath+"/app.css", nil)
	if w.Code != http.StatusOK || !strings.HasPrefix(w.Header().Get("Content-Type"), "text/css") || w.Header().Get("ETag") == "" || w.Header().Get("Cache-Control") != "public, no-cache" {
		t.Errorf("app.css: %d %v", w.Code, w.Header())
	}
	if w2 := get(AssetsPath+"/app.css", http.Header{"If-None-Match": {w.Header().Get("ETag")}}); w2.Code != http.StatusNotModified {
		t.Errorf("app.css again: %d", w2.Code)
	}
	if w := get(AssetsPath+"/topology.js", nil); w.Code != http.StatusOK || !strings.HasPrefix(w.Header().Get("Content-Type"), "text/javascript") {
		t.Errorf("topology.js: %d %s", w.Code, w.Header().Get("Content-Type"))
	}
	for _, path := range []string{AssetsPath + "/", AssetsPath + "/nothing.js", AssetsPath + "/../assets.go", AssetsPath + "/.hidden"} {
		if w := get(path, nil); w.Code != http.StatusNotFound {
			t.Errorf("%s: %d", path, w.Code)
		}
	}
	// a directory, while the files are edited
	d, err := newAssets("public")
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest(http.MethodGet, AssetsPath+"/i18n.js", nil)
	w = httptest.NewRecorder()
	d.ServeHTTP(w, r)
	if w.Code != http.StatusOK || w.Header().Get("ETag") == "" {
		t.Errorf("from the directory: %d %v", w.Code, w.Header())
	}
	if _, err := newAssets("nowhere"); err == nil {
		t.Error("a directory that is not there is taken")
	}
}

// The scripts' texts are written in English in _("…") (i18n.js), and
// public/ja.js has their Japanese: every text a script says has its
// translation there, with the same {name}s, every translation there is
// of a text a script says, and _ is always called with the text written
// in the call, which is what lets this see every one.
func TestScriptTexts(t *testing.T) {
	ja := japanese(t)
	call := regexp.MustCompile(`(?:^|[^\w$.])_\(`)
	text := regexp.MustCompile(`^"((?:[^"\\\n]|\\.)*)"`)
	said := map[string]string{} // a text, and the first script that says it
	scripts, err := filepath.Glob(filepath.Join("public", "*.js"))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range scripts {
		if base := filepath.Base(path); base == "ja.js" || base == "i18n.js" {
			continue
		}
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for i, line := range strings.Split(string(src), "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "//") {
				continue
			}
			for _, at := range call.FindAllStringIndex(line, -1) {
				m := text.FindStringSubmatch(line[at[1]:])
				if m == nil {
					t.Errorf("%s:%d: _ is called without its text written in the call", path, i+1)
					continue
				}
				s, err := strconv.Unquote(`"` + m[1] + `"`)
				if err != nil {
					t.Errorf("%s:%d: %v", path, i+1, err)
					continue
				}
				if _, ok := said[s]; !ok {
					said[s] = path
				}
			}
		}
	}
	if len(said) == 0 {
		t.Fatal("no texts in the scripts")
	}
	for s, path := range said {
		tr, ok := ja[s]
		switch {
		case !ok:
			t.Errorf("%s: %q has no Japanese in public/ja.js", path, s)
		case !slices.Equal(placeholders(s), placeholders(tr)):
			t.Errorf("public/ja.js: %q says %v, its English %v", tr, placeholders(tr), placeholders(s))
		}
	}
	for s := range ja {
		if _, ok := said[s]; !ok {
			t.Errorf("public/ja.js: %q is said by no script", s)
		}
	}
}

// japanese reads public/ja.js: a JSON object, but for its comment lines,
// after export default. A text twice is an error.
func japanese(t *testing.T) map[string]string {
	t.Helper()
	src, err := os.ReadFile(filepath.Join("public", "ja.js"))
	if err != nil {
		t.Fatal(err)
	}
	var lines []string
	for _, line := range strings.Split(string(src), "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "//") {
			lines = append(lines, line)
		}
	}
	obj, ok := strings.CutPrefix(strings.TrimSpace(strings.Join(lines, "\n")), "export default ")
	if !ok {
		t.Fatal("public/ja.js: no export default")
	}
	dec := json.NewDecoder(strings.NewReader(strings.TrimSuffix(obj, ";")))
	if _, err := dec.Token(); err != nil {
		t.Fatalf("public/ja.js: %v", err)
	}
	out := map[string]string{}
	for dec.More() {
		k, err := dec.Token()
		if err != nil {
			t.Fatalf("public/ja.js: %v", err)
		}
		var v string
		if err := dec.Decode(&v); err != nil {
			t.Fatalf("public/ja.js: %q: %v", k, err)
		}
		if _, dup := out[k.(string)]; dup {
			t.Errorf("public/ja.js: %q twice", k)
		}
		out[k.(string)] = v
	}
	return out
}

// placeholders are the {name}s of a text, each once, in order of name.
func placeholders(s string) []string {
	var out []string
	for _, m := range regexp.MustCompile(`\{(\w+)\}`).FindAllStringSubmatch(s, -1) {
		out = append(out, m[1])
	}
	slices.Sort(out)
	return slices.Compact(out)
}
