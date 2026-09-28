# Sample architectures

Under [`examples/`](https://github.com/shibukawa/bdf/tree/main/examples) sit four small, runnable Go servers. The first three pick different points on the same tradeoff: how much conversion work the server does ahead of time, and how much waits for the browser when a document is actually opened. They serve the same six files from [`examples/sample-files/`](https://github.com/shibukawa/bdf/tree/main/examples/sample-files) and share the small embeddable viewer, [`examples/miniviewer`](https://github.com/shibukawa/bdf/tree/main/examples/miniviewer), described in [In the browser](../architecture/browser.html).

| | [light-server](light-server.html) | [preview-server](preview-server.html) | [search](search.html) |
|---|---|---|---|
| Server converts to bdf | never | once, cached | once, cached |
| Server makes a thumbnail | once, cached | once, cached | once, cached |
| Where a document converts | in the browser, on open | — (already bdf) | — (already bdf) |
| Front end ships | wasm converters + renderer | renderer only | renderer only |
| Extra service | — | — | Meilisearch |
| What it's for | many documents, light server, no conversion capacity to provision | documents opened often, fastest and most consistent open | full-text search across every document |

The fourth, [secure-reader](secure-reader.html), is about who may read: readers log in, and a book is sent ten pages at a time, each segment sealed for a key made for its request.

None of these is "the" recommended architecture. Read [Why bdf](../why.html) and pick based on what you're actually building — how often a given document is opened relative to how often it changes, whether the server has spare CPU for conversion, and whether documents need to be searchable.

## light-server

The server never converts anything. It only draws a thumbnail once per document (converting in memory just for that, then discarding the result) and serves the original files as they are; opening one converts it to bdf inside the browser, with the same wasm converter modules the [demo site](https://shibukawa.github.io/bdf/) itself ships. This is the architecture that costs the server the least — no conversion code runs on a request — at the cost of a wasm module the browser fetches once and a moment of local conversion on each open. → [Read more](light-server.html)

## preview-server

The opposite tradeoff: the server converts every document to bdf once, ahead of time, and caches it. Opening a document only fetches the cached bdf (`net/http`'s `ServeContent` answers Range requests on its own, so only the parts of the pages in view are ever transferred) and draws it — the browser front end carries no converter, only the renderer. Documents open faster and more predictably, at the cost of conversion capacity on the server and a place to cache the result. → [Read more](preview-server.html)

## search

Builds on preview-server's server-side conversion and adds a search index: each page's text (`Document.SearchText`) goes into [Meilisearch](https://www.meilisearch.com/), one entry per page, and `/search` proxies a query to it. A hit opens the document at the exact page it was found on. → [Read more](search.html)

## secure-reader

Readers log in and read the PDF books they own; a book is never sent as a file. The server converts each book once and hands it out ten pages at a time (package `segment`), each segment sealed for a key pair the browser made for that request alone, so a recording of the traffic stays sealed even if a key or a session is taken later. A reader who does not own a book may read its first ten pages, and nobody may read faster than a person does. → [Read more](secure-reader.html)
