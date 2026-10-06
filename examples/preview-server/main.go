// Command preview-server is a sample architecture (docs/examples/preview-server):
// the server converts every document to bdf itself, once, and caches both
// the bdf and its thumbnail on disk (makePreviews). The browser never
// converts anything — it only ever fetches an already-made bdf file and
// draws it, so the front end needs none of the wasm converter modules
// light-server ships (web/build.mjs here is a few hundred KB, not tens of
// megabytes).
//
// This is the heavier of the two ends of docs/examples/index's tradeoff:
// the server needs the same Go converter code as any other bdf tool (still
// no cgo, no external process — see docs/why), and every document it might
// ever serve is converted up front rather than on first view. In return,
// a document opens as fast as its bdf can be range-fetched, on any device,
// without downloading a converter.
//
//	cd examples/preview-server && node web/build.mjs && go run .
//	open http://127.0.0.1:8082/
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

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
	_ "github.com/shibukawa/bdf/converter/all"
	"github.com/shibukawa/bdf/thumbnail"
)

func main() {
	dir := flag.String("dir", "../sample-files", "directory of documents to convert and serve")
	cacheDir := flag.String("cache", "./cache", "directory the converted .bdf and thumbnail files are cached in")
	webDir := flag.String("web", "./web", "directory of the built front end (see web/build.mjs)")
	addr := flag.String("listen", "127.0.0.1:8082", "address to listen on")
	flag.Parse()

	if err := os.MkdirAll(*cacheDir, 0o755); err != nil {
		log.Fatal(err)
	}
	files, err := listFiles(*dir)
	if err != nil {
		log.Fatal(err)
	}
	// A real server would do this as each file is uploaded (or lazily, the
	// first time it's requested) instead of scanning the directory and
	// blocking startup on every conversion; kept simple and synchronous
	// here so the whole pipeline reads top to bottom in one function.
	docs := makePreviews(*dir, *cacheDir, files)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", indexHandler(docs))
	// http.FileServer answers Range requests on its own (net/http's
	// ServeContent): the viewer fetches only the parts of a .bdf it needs
	// for the pages in view, exactly as it would from a CDN.
	mux.Handle("GET /cache/", http.StripPrefix("/cache/", http.FileServer(http.Dir(*cacheDir))))
	mux.Handle("GET /view/", http.StripPrefix("/view/", http.FileServer(http.Dir(*webDir))))

	log.Printf("preview-server: %d of %d document(s) converted, listening on http://%s/", len(docs), len(files), *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

// doc is one converted document, ready to serve.
type doc struct {
	Name    string // the source file's name
	BDF     string // its .bdf file's name under -cache
	Summary string // what the converter said of it
}

func listFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("preview-server: %w", err)
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

// makePreviews converts each file that has no cached bdf yet, and its
// thumbnail alongside it, and returns the documents ready to serve (in the
// order of files). A file that fails to convert is logged and left out.
func makePreviews(dir, cacheDir string, files []string) []doc {
	var docs []doc
	for _, name := range files {
		bdfName := name + ".bdf"
		bdfPath := filepath.Join(cacheDir, bdfName)
		pngPath := filepath.Join(cacheDir, name+".png")
		if st, err := os.Stat(bdfPath); err == nil && st.Size() > 0 {
			docs = append(docs, doc{Name: name, BDF: bdfName, Summary: "cached"})
			continue
		}
		res, err := converter.ConvertFile(filepath.Join(dir, name), "", &converter.Options{Params: map[string]string{"remote": "false"}})
		if err != nil {
			log.Printf("preview-server: %s: not converted: %v", name, err)
			continue
		}
		if err := writeBdf(bdfPath, res.Doc); err != nil {
			log.Printf("preview-server: %s: writing bdf: %v", name, err)
			continue
		}
		if err := writeThumbnail(pngPath, res.Doc); err != nil {
			// the document itself converted fine; a thumbnail failure is not fatal
			log.Printf("preview-server: %s: thumbnail: %v", name, err)
		}
		docs = append(docs, doc{Name: name, BDF: bdfName, Summary: res.Summary})
	}
	return docs
}

func writeBdf(path string, d *bdf.Document) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()
	return d.WriteSingle(f)
}

func writeThumbnail(path string, d *bdf.Document) (err error) {
	th, err := thumbnail.Make(d, &thumbnail.Options{Size: 320})
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()
	return thumbnail.Encode(f, th.Image, thumbnail.PNG)
}

var indexTemplate = template.Must(template.New("index").Parse(`<!doctype html>
<meta charset="utf-8">
<title>preview-server: BDF sample architecture</title>
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
  .files small { display: block; padding: 0 .75rem .5rem; color: #666; font-size: .75rem; }
</style>
<h1>preview-server</h1>
<p class="lead">A sample architecture (<a href="https://github.com/shibukawa/bdf/blob/main/docs/examples/preview-server.md">docs/examples/preview-server</a>): the server converted every document below to BDF already. Opening one only fetches and draws it.</p>
<ul class="files">
{{range .}}<li><a href="/view/?src=/cache/{{.BDF}}"><img src="/cache/{{.Name}}.png" alt="" loading="lazy"><span>{{.Name}}</span><small>{{.Summary}}</small></a></li>
{{end}}</ul>
`))

func indexHandler(docs []doc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := indexTemplate.Execute(w, docs); err != nil {
			log.Print(err)
		}
	}
}
