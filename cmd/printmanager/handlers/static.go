package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"
)

// StaticPath is where the pages' scripts and style are served.
//
// They are not served at popcornweb's /public: under pw dev (the pwdev
// build tag) that mount reads only dist/public of the working directory,
// which pw's asset step writes and bdf does not use, so it answered 404
// for every file. This handler serves the embedded files in every mode.
const StaticPath = "/static/"

type staticFile struct {
	body        []byte
	etag, ctype string
}

// staticHandler serves the files of a tree, read once, with an ETag of
// their content (the tree is in the binary: it changes only with it).
func staticHandler(tree fs.FS) (http.Handler, error) {
	files := map[string]staticFile{}
	err := fs.WalkDir(tree, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := fs.ReadFile(tree, name)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(b)
		ctype := mime.TypeByExtension(path.Ext(name))
		if ctype == "" {
			ctype = "application/octet-stream"
		}
		files[name] = staticFile{body: b, etag: `"` + hex.EncodeToString(sum[:8]) + `"`, ctype: ctype}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, ok := files[strings.TrimPrefix(r.URL.Path, StaticPath)]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("ETag", f.etag)
		w.Header().Set("Cache-Control", "no-cache")
		if r.Header.Get("If-None-Match") == f.etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Content-Type", f.ctype)
		w.Write(f.body)
	}), nil
}
