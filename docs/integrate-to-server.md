# Integrating into a server

The Go `converter` registry contains only formats whose packages your program imports. Blank import the formats your service accepts, then serve the resulting BDF to the [frontend viewer](integrating-to-frontend.md). Go 1.27 or later is required.

## Choose the input formats

For example, a service accepting PDF, EPUB, Word, and CSV/TSV imports these converters. The `converter/csv` package handles TSV too.

```go
import (
    "github.com/shibukawa/bdf/converter"
    _ "github.com/shibukawa/bdf/converter/pdf"
    _ "github.com/shibukawa/bdf/converter/epub"
    _ "github.com/shibukawa/bdf/converter/docx"
    _ "github.com/shibukawa/bdf/converter/csv"
)
```

Add `converter/pptx` for PowerPoint, `converter/xlsx` for Excel, `converter/visio` for Visio, or `converter/parquet` for Parquet. To accept every format, import `converter/all` instead of the individual packages. See [Formats](formats/index.md) and the [`converter/`](../converter/) directories for the available names.

To omit music conversion, leave out `converter/mml`, `converter/midi`, and `converter/musicxml`. A server that only converts to BDF does not need to import the `imagebdf` or `thumbnail` packages for drawing.

## Convert an HTTP upload

Save this example as `cmd/my-bdf-server/main.go` inside this repository. Send source bytes to `POST /convert?name=report.pdf`; it returns a single-file BDF. The `name` query parameter also gives the converter an extension for inputs such as CSV and Markdown that may be hard to identify from their bytes alone.

```go
package main

import (
    "bytes"
    "io"
    "log"
    "net/http"
    "path/filepath"

    "github.com/shibukawa/bdf/converter"
    _ "github.com/shibukawa/bdf/converter/csv"
    _ "github.com/shibukawa/bdf/converter/docx"
    _ "github.com/shibukawa/bdf/converter/epub"
    _ "github.com/shibukawa/bdf/converter/pdf"
)

func main() {
    http.HandleFunc("POST /convert", func(w http.ResponseWriter, r *http.Request) {
        input, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 32<<20))
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        name := filepath.Base(r.URL.Query().Get("name"))
        result, err := converter.Convert(bytes.NewReader(input), int64(len(input)), "",
            &converter.Options{FileName: name})
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnprocessableEntity)
            return
        }
        var output bytes.Buffer
        if err := result.Doc.WriteSingle(&output); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/octet-stream")
        _, _ = w.Write(output.Bytes())
    })
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

```sh
go run ./cmd/my-bdf-server
curl --data-binary @report.pdf 'http://localhost:8080/convert?name=report.pdf' -o report.bdf
```

This compact example keeps the input and output in memory. For a service, choose upload limits, authentication, conversion concurrency, and a cache appropriate to your workload.

## Store and serve BDF

To reuse a conversion, write `result.Doc.WriteSingle(file)` to a `.bdf` file and serve it through `http.FileServer` or equivalent. Serving single-file BDF with HTTP Range support lets the viewer fetch the parts it needs. For a CDN or object storage, use `result.Doc.WriteSplit(dir)` and open it in the frontend with `{ kind: "split", base: "..." }`.

For a file already on the server, use `converter.ConvertFile(path, "", &converter.Options{})`. Formats with references to other files, such as HTML images or other KiCad sheets, can read them through `Options.Dir`, which `ConvertFile` sets from the input path. For multi-file uploads, provide `Options.Files`. See [preview-server](examples/preview-server.md) for a working cache and serving example.
