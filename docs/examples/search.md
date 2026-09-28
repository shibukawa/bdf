# search: with a search engine

[`examples/search`](https://github.com/shibukawa/bdf/tree/main/examples/search) — builds on [preview-server](preview-server.html)'s server-side conversion and adds full-text search across every document, backed by [Meilisearch](https://www.meilisearch.com/).

```sh
docker run --rm -p 7700:7700 getmeili/meilisearch:v1.12 \
  meilisearch --master-key=dev-key --no-analytics

cd examples/search
node web/build.mjs
go run . -meili-key dev-key
open http://127.0.0.1:8083/
```

Any HTTP search engine would do the same job here — Meilisearch's plain REST API is small enough that this example calls it directly with `net/http`, adding no client library dependency.

## What the server does

On top of what [preview-server](preview-server.html) does (convert and cache every document as bdf, and a thumbnail), it opens each cached bdf back up (`bdf.OpenSingleFile`, `(*Reader).ToDocument`) and calls `(*bdf.Document).SearchText()` — the same call a server would use to make its own converted document searchable without reconverting. `SearchText` returns the document's metadata and the text of each page (a sheet is one entry, without a page number); each page becomes one Meilisearch document, with an id, the file and view it came from, its page number, and the document's title.

```go
type meiliDoc struct {
	ID    string `json:"id"`
	File  string `json:"file"`
	BDF   string `json:"bdf"`
	View  string `json:"view"`
	Page  int    `json:"page,omitempty"`
	Title string `json:"title"`
	Text  string `json:"text"`
}
```

These are pushed to Meilisearch in one batch at startup (`POST /indexes/{index}/documents`), and the server waits for the indexing task to finish before it starts answering requests — so the first search right after startup already finds everything, not just whatever Meilisearch got to first (indexing is normally asynchronous).

`GET /search?q=` proxies a query to Meilisearch's own `/search` endpoint, asking for a cropped, `<mark>`-highlighted snippet of the matching text (`attributesToCrop`, `attributesToHighlight`), and returns a small JSON array — file, bdf, view, page, title, snippet — that says nothing about how the search itself works. The API key stays on the server; the browser never talks to Meilisearch directly.

## What the browser does

The home page's search box (`web/search.ts`) calls `/search`, lists what comes back, and links each hit to `/view/?src=/cache/NAME.bdf#view=ID&page=N`. `web/main.ts` — otherwise identical to preview-server's — reads `view` and `page` off the URL fragment once the document is open and calls `viewer.show(view, page - 1)` to jump straight there, the same [`MiniViewer`](https://github.com/shibukawa/bdf/tree/main/examples/miniviewer) API any viewer built on it would use for a deep link.

## Why this shape

Search only needs the document's text, not its rendering, so it composes cleanly with whichever conversion architecture a system already uses — this example builds it on preview-server's for convenience, but light-server's would work too: SearchText only needs a converted bdf, made once when a document is indexed (or re-indexed), not on every request. Keeping the search API key server-side and proxying the query, rather than querying Meilisearch straight from the browser, means the browser never needs credentials for the search engine at all.
