//go:build js && wasm

// Command bdfwasm is the converters built for the browser: a page hands it
// the bytes of a file and gets a BDF document back, so that the file never
// leaves the browser. Build it with
//
//	GOOS=js GOARCH=wasm go build -tags bdf_noconv -o bdf.wasm ./cmd/bdfwasm
//
// and run it with Go's wasm_exec.js (in a Worker: a conversion keeps the
// thread busy). -tags pdfonly or officeonly leaves out the Office
// converters (Word, PowerPoint, Excel, CSV, Visio, draw.io, DXF, metafiles) or
// the PDF one, for smaller modules. bdf_noconv leaves out the image and
// WOFF2 encoders: a document drawn where it is converted gains nothing from
// them.
//
// The program sets globalThis.bdfConverter and waits for calls:
//
//	bdfConverter.formats: [{name, description, extensions}]
//	bdfConverter.convert(data: Uint8Array, options?: {
//	  format?: string,   // a format name; detected from the content when absent
//	  password?: string, // the open password of an encrypted input
//	  fonts?: string,    // URL of a font directory (see below)
//	}): Promise<{bdf: Uint8Array, format: string, summary: string, warnings: string[], protected: boolean}>
//	bdfConverter.open(data, options?): Promise<{bdf, format, pages: number, warnings, stream?}>
//
// convert converts the input whole into a single-file bdf. open starts a
// conversion done a page at a time (converter.OpenStream), for a viewer that
// draws the pages as they are converted: with pages > 0, bdf is the outline
// (every page sized, the pages of the first view without layers) and
// stream converts the rest:
//
//	stream.page(i: number): Promise<{bdf: Uint8Array, warnings: string[]}>
//	stream.finish(): Promise<{bdf, format, summary, warnings, protected}>
//	stream.close(): void
//
// page converts page i of the first view and returns a single-file bdf
// whose only page is that page, with the parts no earlier call returned;
// warnings are the new ones. finish returns the whole document, as convert
// makes it, with every warning. close lets the stream go (call it after
// finish too). A format that is not converted a page at a time comes back
// whole from open, with pages 0, a summary and no stream.
//
// A failed conversion rejects with an Error whose code is
// "password-required", "wrong-password" or "unknown-format" when it is one
// of those. The document is not encrypted even when the input was.
//
// The Office converters lay text out with the fonts of the font directory:
// index.json there lists its files as [{"name", "size", "scan"}], where
// scan holds the [offset, length] ranges of the tables the font scan reads
// (the table directories, name, OS/2 and post), fetched ahead in parallel.
// A font's whole file is fetched when a document uses it.
package main

import (
	"bytes"
	"compress/flate"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"syscall/js"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
)

func main() {
	var formats []any
	for _, f := range converter.Formats() {
		exts := make([]any, len(f.Extensions))
		for i, e := range f.Extensions {
			exts[i] = e
		}
		formats = append(formats, map[string]any{"name": f.Name, "description": f.Description, "extensions": exts})
	}
	api := js.ValueOf(map[string]any{"formats": formats})
	api.Set("convert", js.FuncOf(convert))
	api.Set("open", js.FuncOf(open))
	js.Global().Set("bdfConverter", api)
	select {}
}

var (
	fontMu sync.Mutex
	fontFS = map[string]*httpFS{}
)

// fonts returns the font directory at url, which keeps what it fetched for
// the conversions that follow.
func fonts(url string) (*httpFS, error) {
	fontMu.Lock()
	defer fontMu.Unlock()
	if f, ok := fontFS[url]; ok {
		return f, nil
	}
	f, err := newHTTPFS(url)
	if err != nil {
		return nil, err
	}
	fontFS[url] = f
	return f, nil
}

// request is the input of convert and open.
type request struct {
	data                      []byte
	format, password, fontURL string
}

func readRequest(args []js.Value) (*request, error) {
	if len(args) == 0 || args[0].Type() != js.TypeObject {
		return nil, errors.New("want the input bytes as a Uint8Array")
	}
	req := &request{data: make([]byte, args[0].Get("length").Int())}
	js.CopyBytesToGo(req.data, args[0])
	str := func(name string) string {
		if len(args) < 2 || args[1].Type() != js.TypeObject {
			return ""
		}
		if v := args[1].Get(name); v.Type() == js.TypeString {
			return v.String()
		}
		return ""
	}
	req.format, req.password, req.fontURL = str("format"), str("password"), str("fonts")
	return req, nil
}

func (req *request) options() (*converter.Options, error) {
	opts := &converter.Options{Password: req.password, NoSystemFonts: true}
	if req.fontURL != "" {
		fsys, err := fonts(req.fontURL)
		if err != nil {
			return nil, fmt.Errorf("fonts: %w", err)
		}
		opts.FontFS = fsys
	}
	return opts, nil
}

func convert(_ js.Value, args []js.Value) any {
	req, err := readRequest(args)
	if err != nil {
		return reject(fmt.Errorf("convert: %w", err))
	}
	return promise(func() (any, error) {
		opts, err := req.options()
		if err != nil {
			return nil, err
		}
		res, err := converter.Convert(bytes.NewReader(req.data), int64(len(req.data)), req.format, opts)
		if err != nil {
			return nil, err
		}
		return result(res, res.Warnings)
	})
}

func open(_ js.Value, args []js.Value) any {
	req, err := readRequest(args)
	if err != nil {
		return reject(fmt.Errorf("open: %w", err))
	}
	return promise(func() (any, error) {
		opts, err := req.options()
		if err != nil {
			return nil, err
		}
		w := &warnings{}
		opts.Warn = w.add
		s, err := converter.OpenStream(bytes.NewReader(req.data), int64(len(req.data)), req.format, opts)
		if err != nil {
			return nil, err
		}
		if s.Pages() == 0 {
			res, err := s.Finish()
			if err != nil {
				return nil, err
			}
			out, err := result(res, w.all())
			if err != nil {
				return nil, err
			}
			out["pages"] = 0
			return out, nil
		}
		outline := s.Outline()
		b, err := single(outline)
		if err != nil {
			return nil, err
		}
		return map[string]any{"bdf": b, "format": outline.Meta.Source, "pages": s.Pages(), "warnings": w.fresh(), "stream": streamValue(s, w)}, nil
	})
}

// streamValue is the JavaScript side of a stream: page, finish and close.
func streamValue(s converter.Stream, w *warnings) js.Value {
	var mu sync.Mutex // the calls run in goroutines of their own
	var funcs []js.Func
	method := func(fn func(args []js.Value) any) js.Func {
		f := js.FuncOf(func(_ js.Value, args []js.Value) any { return fn(args) })
		funcs = append(funcs, f)
		return f
	}
	v := js.Global().Get("Object").New()
	v.Set("page", method(func(args []js.Value) any {
		if len(args) == 0 || args[0].Type() != js.TypeNumber {
			return reject(errors.New("page: want a page index"))
		}
		i := args[0].Int()
		return promise(func() (any, error) {
			mu.Lock()
			defer mu.Unlock()
			doc, err := s.Page(i)
			if err != nil {
				return nil, err
			}
			b, err := single(doc)
			if err != nil {
				return nil, err
			}
			return map[string]any{"bdf": b, "warnings": w.fresh()}, nil
		})
	}))
	v.Set("finish", method(func([]js.Value) any {
		return promise(func() (any, error) {
			mu.Lock()
			defer mu.Unlock()
			res, err := s.Finish()
			if err != nil {
				return nil, err
			}
			return result(res, w.all())
		})
	}))
	v.Set("close", method(func([]js.Value) any {
		for _, f := range funcs {
			f.Release()
		}
		return nil
	}))
	return v
}

// warnings collects what a conversion warns about as it goes.
type warnings struct {
	mu   sync.Mutex
	list []string
	seen int
}

func (w *warnings) add(msg string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.list = append(w.list, msg)
}

// fresh returns the warnings added since the last call.
func (w *warnings) fresh() []any {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := anys(w.list[w.seen:])
	w.seen = len(w.list)
	return out
}

func (w *warnings) all() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.list
}

func anys(list []string) []any {
	out := make([]any, len(list))
	for i, s := range list {
		out[i] = s
	}
	return out
}

// single writes a document in the single-file form. Parts are compressed
// quickly: the document is drawn where it is made, and only a download
// would gain from smaller parts.
func single(doc *bdf.Document) (js.Value, error) {
	doc.CompressionLevel = flate.BestSpeed
	var out bytes.Buffer
	if err := doc.WriteSingle(&out); err != nil {
		return js.Value{}, err
	}
	b := js.Global().Get("Uint8Array").New(out.Len())
	js.CopyBytesToJS(b, out.Bytes())
	return b, nil
}

func result(res *converter.Result, warnings []string) (map[string]any, error) {
	b, err := single(res.Doc)
	if err != nil {
		return nil, err
	}
	return map[string]any{"bdf": b, "format": res.Doc.Meta.Source, "summary": res.Summary,
		"warnings": anys(warnings), "protected": res.Protected}, nil
}

// promise runs fn in a goroutine, which may block (on fetches), and settles
// a Promise with its result. A panic rejects it instead of ending the
// program.
func promise(fn func() (any, error)) js.Value {
	var executor js.Func
	executor = js.FuncOf(func(_ js.Value, args []js.Value) any {
		resolve, rejectFn := args[0], args[1]
		go func() {
			defer executor.Release()
			defer func() {
				if r := recover(); r != nil {
					rejectFn.Invoke(jsError(fmt.Errorf("panic: %v\n%s", r, debug.Stack())))
				}
			}()
			v, err := fn()
			if err != nil {
				rejectFn.Invoke(jsError(err))
				return
			}
			resolve.Invoke(v)
		}()
		return nil
	})
	return js.Global().Get("Promise").New(executor)
}

// reject returns a Promise rejected with err.
func reject(err error) js.Value {
	return js.Global().Get("Promise").Call("reject", jsError(err))
}

// jsError makes an Error of err, with a code for the errors a page acts on.
func jsError(err error) js.Value {
	e := js.Global().Get("Error").New(err.Error())
	switch {
	case errors.Is(err, converter.ErrPasswordRequired):
		e.Set("code", "password-required")
	case errors.Is(err, converter.ErrWrongPassword):
		e.Set("code", "wrong-password")
	case errors.Is(err, converter.ErrUnknownFormat):
		e.Set("code", "unknown-format")
	}
	return e
}
