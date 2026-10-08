// Command light-server is a sample architecture (docs/examples/light-server):
// the server never converts a document to bdf. It serves the documents as
// they are and a thumbnail image of each (made once, ahead of time, by
// converting in memory and throwing the bdf away — see makeThumbnails), and
// the browser does the actual conversion and rendering, using the bdf
// converters built as WebAssembly (web/build.mjs builds them, reusing
// examples/common/build.mjs, the same helper the demo site is built with).
//
// This is the lightest server of the three sample architectures
// (docs/examples/index): no Go conversion code runs on every view, only
// once per document to draw its thumbnail, and the server holds no bdf at
// all. It costs the browser a WebAssembly module to fetch and a moment to
// convert the file locally — the same tradeoff the demo site itself makes.
// See docs/examples/preview-server for the opposite tradeoff.
//
//	cd examples/light-server && node web/build.mjs && go run .
//	open http://127.0.0.1:8081/
package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/shibukawa/bdf/converter"
	_ "github.com/shibukawa/bdf/converter/all"
	"github.com/shibukawa/bdf/thumbnail"
)

func main() {
	dir := flag.String("dir", "../sample-files", "directory of documents to serve")
	thumbDir := flag.String("thumbs", "./thumbs", "directory thumbnail images are cached in")
	webDir := flag.String("web", "./web", "directory of the built front end (see web/build.mjs)")
	addr := flag.String("listen", "127.0.0.1:8081", "address to listen on")
	flag.Parse()

	if err := os.MkdirAll(*thumbDir, 0o755); err != nil {
		log.Fatal(err)
	}
	files, err := listFiles(*dir)
	if err != nil {
		log.Fatal(err)
	}
	// Thumbnails are made once, up front: a real server would do this when
	// a file is uploaded (or lazily, on the first request for it) instead
	// of scanning the directory at startup.
	makeThumbnails(*dir, *thumbDir, files)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", indexHandler(files, *thumbDir))
	// The originals: fetched by the browser's own converter, which runs
	// entirely client-side (web/main.ts). http.FileServer answers Range
	// requests on its own, though this server never needs them itself.
	mux.Handle("GET /files/", http.StripPrefix("/files/", http.FileServer(http.Dir(*dir))))
	mux.Handle("GET /thumbs/", http.StripPrefix("/thumbs/", http.FileServer(http.Dir(*thumbDir))))
	mux.Handle("GET /view/", http.StripPrefix("/view/", http.FileServer(http.Dir(*webDir))))

	log.Printf("light-server: %d document(s) from %s, listening on http://%s/", len(files), *dir, *addr)
	// a connection that sends its request slowly, or none, is not kept for ever
	srv := &http.Server{Addr: *addr, Handler: mux, ReadHeaderTimeout: 10 * time.Second, ReadTimeout: time.Minute, IdleTimeout: 2 * time.Minute}
	log.Fatal(srv.ListenAndServe())
}

// listFiles returns the regular files directly inside dir, sorted by name.
func listFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("light-server: %w", err)
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)
	return files, nil
}

// makeThumbnails converts each file that has no cached thumbnail yet, draws
// its thumbnail, and discards the converted bdf: this server never keeps
// one, so only the first page, all the thumbnail shows, is converted. A file whose format isn't registered, or that fails to convert, is
// skipped with a log line rather than stopping the server.
func makeThumbnails(dir, thumbDir string, files []string) {
	for _, name := range files {
		out := filepath.Join(thumbDir, name+".png")
		if _, err := os.Stat(out); err == nil {
			continue // already made
		}
		res, err := converter.ConvertFile(filepath.Join(dir, name), "", &converter.Options{Pages: converter.PageList(1), Params: map[string]string{"remote": "false"}})
		if err != nil {
			log.Printf("light-server: %s: not converted for its thumbnail: %v", name, err)
			continue
		}
		th, err := thumbnail.Make(res.Doc, &thumbnail.Options{Size: 320})
		if err != nil {
			log.Printf("light-server: %s: thumbnail: %v", name, err)
			continue
		}
		f, err := os.Create(out)
		if err != nil {
			log.Fatal(err)
		}
		err = thumbnail.Encode(f, th.Image, thumbnail.PNG)
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
		if err != nil {
			log.Printf("light-server: %s: writing thumbnail: %v", name, err)
			os.Remove(out)
		}
	}
}

var indexTemplate = template.Must(template.New("index").Parse(`<!doctype html>
<meta charset="utf-8">
<title>light-server: BDF sample architecture</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
  :root { color-scheme: light dark; font: 16px system-ui, sans-serif; }
  body { max-width: 60rem; margin: 2rem auto; padding: 0 1rem; }
  h1 { font-size: 1.375rem; }
  p.lead { color: #666; }
  .files { display: grid; grid-template-columns: repeat(auto-fill, minmax(11rem, 1fr)); gap: 1rem; list-style: none; margin: 1.5rem 0; padding: 0; }
  .files li { border: 1px solid #ccc; border-radius: 8px; overflow: hidden; }
  .files a { display: block; color: inherit; text-decoration: none; }
  .files img { display: block; width: 100%; aspect-ratio: 4 / 3; object-fit: contain; background: #eee; }
  .files span { display: block; padding: .5rem .75rem; font-size: .875rem; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
</style>
<h1>light-server</h1>
<p class="lead">A sample architecture (<a href="https://github.com/shibukawa/bdf/blob/main/docs/examples/light-server.md">docs/examples/light-server</a>): the server only ever made these thumbnails. Opening a document converts it to BDF inside your browser.</p>
<ul class="files">
{{range .}}<li><a href="/view/?file={{.}}"><img src="/thumbs/{{.}}.png" alt="" loading="lazy"><span>{{.}}</span></a></li>
{{end}}</ul>
`))

func indexHandler(files []string, thumbDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only list files whose thumbnail actually exists: one that failed
		// to convert (see makeThumbnails) is left out rather than shown
		// broken.
		shown := make([]string, 0, len(files))
		for _, name := range files {
			if _, err := os.Stat(filepath.Join(thumbDir, name+".png")); err == nil {
				shown = append(shown, name)
			}
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := indexTemplate.Execute(w, shown); err != nil {
			log.Print(err)
		}
	}
}
