// Command search is a sample architecture (docs/examples/search): on top of
// what preview-server does (convert and cache every document as bdf, and a
// thumbnail), it also extracts each document's text with
// bdf.Document.SearchText and indexes it in Meilisearch, one document per
// page (a sheet is one entry, without a page number — see
// bdf.Document.SearchText's doc comment). /search proxies a query to
// Meilisearch and the browser opens a hit at its exact page, no server-side
// rendering involved: only the bdf and a JSON search response ever cross
// the wire.
//
// Needs a running Meilisearch (any search engine with an HTTP API would do
// the same job; Meilisearch's is small enough to call with net/http alone,
// so this example adds no dependency):
//
//	docker run --rm -p 7700:7700 getmeili/meilisearch:v1.12 \
//	  meilisearch --master-key=dev-key --no-analytics
//	cd examples/search && node web/build.mjs && go run . -meili-key dev-key
//	open http://127.0.0.1:8083/
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
	_ "github.com/shibukawa/bdf/converter/all"
	"github.com/shibukawa/bdf/thumbnail"
)

func main() {
	dir := flag.String("dir", "../sample-files", "directory of documents to convert, serve and index")
	cacheDir := flag.String("cache", "./cache", "directory the converted .bdf and thumbnail files are cached in")
	webDir := flag.String("web", "./web", "directory of the built front end (see web/build.mjs)")
	addr := flag.String("listen", "127.0.0.1:8083", "address to listen on")
	meiliHost := flag.String("meili", "http://127.0.0.1:7700", "Meilisearch URL")
	meiliKey := flag.String("meili-key", "", "Meilisearch API key (its --master-key, or a key with search and document rights)")
	meiliIndex := flag.String("meili-index", "bdf-search-example", "Meilisearch index to use")
	flag.Parse()

	if err := os.MkdirAll(*cacheDir, 0o755); err != nil {
		log.Fatal(err)
	}
	files, err := listFiles(*dir)
	if err != nil {
		log.Fatal(err)
	}
	m := &meili{host: *meiliHost, key: *meiliKey, index: *meiliIndex}

	// As in preview-server: converted and cached up front for this example,
	// where a real server would do it per upload. The text index is rebuilt
	// from the cached bdf too, so re-running this program re-indexes
	// without reconverting.
	docs := makePreviews(*dir, *cacheDir, files)
	if err := m.indexAll(*cacheDir, docs); err != nil {
		log.Printf("search: Meilisearch: %v (is it running? see -meili, -meili-key)", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", indexHandler(docs))
	mux.HandleFunc("GET /search", searchHandler(m))
	mux.Handle("GET /cache/", http.StripPrefix("/cache/", http.FileServer(http.Dir(*cacheDir))))
	mux.Handle("GET /view/", http.StripPrefix("/view/", http.FileServer(http.Dir(*webDir))))

	log.Printf("search: %d of %d document(s) indexed, listening on http://%s/", len(docs), len(files), *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

// doc is one converted, indexed document.
type doc struct {
	Name    string // the source file's name
	BDF     string // its .bdf file's name under -cache
	Summary string
}

func listFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
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
// thumbnail alongside it — the same shape as examples/preview-server's.
func makePreviews(dir, cacheDir string, files []string) []doc {
	var docs []doc
	for _, name := range files {
		bdfName := name + ".bdf"
		bdfPath := filepath.Join(cacheDir, bdfName)
		pngPath := filepath.Join(cacheDir, name+".png")
		if st, err := os.Stat(bdfPath); err != nil || st.Size() == 0 {
			res, err := converter.ConvertFile(filepath.Join(dir, name), "", &converter.Options{Params: map[string]string{"remote": "false"}})
			if err != nil {
				log.Printf("search: %s: not converted: %v", name, err)
				continue
			}
			if err := writeFile(bdfPath, res.Doc.WriteSingle); err != nil {
				log.Printf("search: %s: writing bdf: %v", name, err)
				continue
			}
			if err := writeThumbnail(pngPath, res.Doc); err != nil {
				log.Printf("search: %s: thumbnail: %v", name, err) // not fatal: the document itself converted fine
			}
			docs = append(docs, doc{Name: name, BDF: bdfName, Summary: res.Summary})
			continue
		}
		docs = append(docs, doc{Name: name, BDF: bdfName, Summary: "cached"})
	}
	return docs
}

func writeFile(path string, write func(io.Writer) error) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()
	return write(f)
}

func writeThumbnail(path string, d *bdf.Document) error {
	th, err := thumbnail.Make(d, &thumbnail.Options{Size: 320})
	if err != nil {
		return err
	}
	return writeFile(path, func(w io.Writer) error { return thumbnail.Encode(w, th.Image, thumbnail.PNG) })
}

var indexTemplate = template.Must(template.New("index").Parse(`<!doctype html>
<meta charset="utf-8">
<title>search: bdf sample architecture</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
  :root { color-scheme: light dark; font: 16px system-ui, sans-serif; }
  body { max-width: 60rem; margin: 2rem auto; padding: 0 1rem; }
  h1 { font-size: 1.375rem; }
  p.lead { color: #666; }
  form { display: flex; gap: .5rem; margin: 1.5rem 0; }
  input[type=search] { flex: 1; padding: .5rem .75rem; font: inherit; }
  button { padding: .5rem 1rem; font: inherit; }
  #hits { margin: 0; padding: 0; list-style: none; }
  #hits li { padding: .75rem 0; border-top: 1px solid #ddd; }
  #hits a { font-weight: 600; text-decoration: none; color: inherit; }
  #hits .snippet { color: #444; margin: .25rem 0 0; }
  #hits mark { background: #ff0; }
  .files { display: grid; grid-template-columns: repeat(auto-fill, minmax(11rem, 1fr)); gap: 1rem; list-style: none; margin: 1.5rem 0; padding: 0; }
  .files li { border: 1px solid #ccc; border-radius: 8px; overflow: hidden; }
  .files a { display: block; color: inherit; text-decoration: none; }
  .files img { display: block; width: 100%; aspect-ratio: 4 / 3; object-fit: contain; background: #eee; }
  .files span { display: block; padding: .5rem .75rem; font-size: .875rem; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
</style>
<h1>search</h1>
<p class="lead">A sample architecture (<a href="https://github.com/shibukawa/bdf/blob/main/docs/examples/search.md">docs/examples/search</a>): every page of every document below is indexed in Meilisearch. A hit opens the document at that page.</p>
<form id="searchForm" role="search">
  <input type="search" id="q" placeholder="search the documents…" autofocus>
  <button type="submit">Search</button>
</form>
<ul id="hits"></ul>
<ul class="files">
{{range .}}<li><a href="/view/?src=/cache/{{.BDF}}"><img src="/cache/{{.Name}}.png" alt="" loading="lazy"><span>{{.Name}}</span></a></li>
{{end}}</ul>
<script type="module" src="/view/search.js"></script>
`))

func indexHandler(docs []doc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := indexTemplate.Execute(w, docs); err != nil {
			log.Print(err)
		}
	}
}

// hit is what /search returns for one page: enough for the front end to
// show a result and open it at the right place.
type hit struct {
	File    string `json:"file"`
	BDF     string `json:"bdf"`
	View    string `json:"view"`
	Page    int    `json:"page,omitempty"`
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
}

func searchHandler(m *meili) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if q == "" {
			json.NewEncoder(w).Encode([]hit{})
			return
		}
		hits, err := m.search(q)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(hits)
	}
}

// --- Meilisearch, over its plain HTTP API (no client library: a handful of JSON calls) ---

type meili struct {
	host, key, index string
}

// call makes one Meilisearch API request and decodes its JSON response into
// out (nil to ignore the body); errCode, when not nil, receives the
// response's "code" field on a non-2xx status (Meilisearch's way of saying
// what went wrong, such as "index_already_exists"), which callers can check
// without having to parse the error text.
func (m *meili) call(method, path string, body, out any) error {
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, m.host+path, r)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if m.key != "" {
		req.Header.Set("Authorization", "Bearer "+m.key)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		var apiErr struct {
			Code, Message string
		}
		b, _ := io.ReadAll(res.Body)
		json.Unmarshal(b, &apiErr)
		if apiErr.Code != "" {
			return &meiliError{code: apiErr.Code, message: apiErr.Message}
		}
		return fmt.Errorf("meilisearch: %s: %s", res.Status, b)
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(res.Body).Decode(out)
}

// meiliError is a Meilisearch API error, with the machine-readable code
// (e.g. "index_already_exists") callers can switch on.
type meiliError struct{ code, message string }

func (e *meiliError) Error() string { return fmt.Sprintf("meilisearch: %s: %s", e.code, e.message) }

// meiliDoc is one page of one document, as Meilisearch stores it. ID must
// be Meilisearch-safe (letters, digits, -, _), so it isn't a file name.
type meiliDoc struct {
	ID    string `json:"id"`
	File  string `json:"file"`
	BDF   string `json:"bdf"`
	View  string `json:"view"`
	Page  int    `json:"page,omitempty"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

// indexAll converts each document's cached bdf into meiliDocs with
// bdf.Document.SearchText and pushes them to Meilisearch, waiting for the
// indexing tasks so that a search right after startup already finds them.
func (m *meili) indexAll(cacheDir string, docs []doc) error {
	err := m.call("POST", "/indexes", map[string]string{"uid": m.index, "primaryKey": "id"}, nil)
	// "index_already_exists" is expected on a second run; anything else is a real problem.
	if apiErr, ok := err.(*meiliError); err != nil && (!ok || apiErr.code != "index_already_exists") {
		return err
	}

	var out []meiliDoc
	for i, d := range docs {
		r, err := bdf.OpenSingleFile(filepath.Join(cacheDir, d.BDF))
		if err != nil {
			log.Printf("search: %s: %v", d.Name, err)
			continue
		}
		document, err := r.ToDocument()
		if err != nil {
			log.Printf("search: %s: %v", d.Name, err)
			continue
		}
		st, err := document.SearchText()
		if err != nil {
			log.Printf("search: %s: %v", d.Name, err)
			continue
		}
		title := d.Name
		if t := st.Meta.DC.Title.First(); t != "" {
			title = t
		}
		for v, view := range st.Views {
			for _, page := range view.Pages {
				out = append(out, meiliDoc{
					ID: fmt.Sprintf("d%d-v%d-p%d", i, v, page.Page), File: d.Name, BDF: d.BDF,
					View: view.ID, Page: page.Page, Title: title, Text: page.Text,
				})
			}
		}
	}
	if len(out) == 0 {
		return nil
	}
	var pushed struct {
		TaskUID int `json:"taskUid"`
	}
	if err := m.call("POST", "/indexes/"+m.index+"/documents", out, &pushed); err != nil {
		return err
	}
	return m.awaitTask(pushed.TaskUID)
}

// awaitTask polls a Meilisearch task until it finishes or 10 seconds pass:
// indexing runs asynchronously, and a demo reads better when the first
// search after startup already sees every document.
func (m *meili) awaitTask(id int) error {
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		var t struct {
			Status string `json:"status"`
			Error  any    `json:"error"`
		}
		if err := m.call("GET", fmt.Sprintf("/tasks/%d", id), nil, &t); err != nil {
			return err
		}
		switch t.Status {
		case "succeeded":
			return nil
		case "failed", "canceled":
			return fmt.Errorf("meilisearch: indexing task %d %s: %v", id, t.Status, t.Error)
		}
		time.Sleep(150 * time.Millisecond)
	}
	return fmt.Errorf("meilisearch: indexing task %d did not finish in time", id)
}

// search queries Meilisearch and turns its hits into the JSON /search
// returns: a cropped, highlighted snippet of the page's text around the
// match, ready to show as a result.
func (m *meili) search(q string) ([]hit, error) {
	body := map[string]any{
		"q": q, "limit": 20,
		"attributesToCrop": []string{"text"}, "cropLength": 24,
		"attributesToHighlight": []string{"text"},
		"highlightPreTag":       "<mark>", "highlightPostTag": "</mark>",
	}
	var res struct {
		Hits []struct {
			meiliDoc
			Formatted struct {
				Text string `json:"text"`
			} `json:"_formatted"`
		} `json:"hits"`
	}
	if err := m.call("POST", "/indexes/"+m.index+"/search", body, &res); err != nil {
		return nil, err
	}
	hits := make([]hit, len(res.Hits))
	for i, h := range res.Hits {
		hits[i] = hit{File: h.File, BDF: h.BDF, View: h.View, Page: h.Page, Title: h.Title, Snippet: h.Formatted.Text}
	}
	return hits, nil
}
