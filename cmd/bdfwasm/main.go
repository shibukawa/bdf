//go:build js && wasm

// Command bdfwasm is the converters built for the browser: a page hands it
// the bytes of a file and gets a BDF document back, so that the file never
// leaves the browser. Build it with
//
//	GOOS=js GOARCH=wasm go build -tags bdf_noconv -o bdf.wasm ./cmd/bdfwasm
//
// and run it with Go's wasm_exec.js (in a Worker: a conversion keeps the
// thread busy). -tags pdfonly or officeonly leaves out the Office
// converters (Word, PowerPoint, Excel, CSV, Visio, draw.io, DXF, Jw_cad, metafiles) or
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

func convert(_ js.Value, args []js.Value) any {
	if len(args) == 0 || args[0].Type() != js.TypeObject {
		return reject(errors.New("convert: want the input bytes as a Uint8Array"))
	}
	data := make([]byte, args[0].Get("length").Int())
	js.CopyBytesToGo(data, args[0])
	str := func(name string) string {
		if len(args) < 2 || args[1].Type() != js.TypeObject {
			return ""
		}
		if v := args[1].Get(name); v.Type() == js.TypeString {
			return v.String()
		}
		return ""
	}
	format, password, fontURL := str("format"), str("password"), str("fonts")
	return promise(func() (any, error) {
		opts := &converter.Options{Password: password, NoSystemFonts: true}
		if fontURL != "" {
			fsys, err := fonts(fontURL)
			if err != nil {
				return nil, fmt.Errorf("fonts: %w", err)
			}
			opts.FontFS = fsys
		}
		res, err := converter.Convert(bytes.NewReader(data), int64(len(data)), format, opts)
		if err != nil {
			return nil, err
		}
		// Parts are compressed quickly: the document is drawn where it is
		// made, and only a download would gain from smaller parts.
		res.Doc.CompressionLevel = flate.BestSpeed
		var out bytes.Buffer
		if err := res.Doc.WriteSingle(&out); err != nil {
			return nil, err
		}
		b := js.Global().Get("Uint8Array").New(out.Len())
		js.CopyBytesToJS(b, out.Bytes())
		warnings := make([]any, len(res.Warnings))
		for i, w := range res.Warnings {
			warnings[i] = w
		}
		return map[string]any{"bdf": b, "format": res.Doc.Meta.Source, "summary": res.Summary,
			"warnings": warnings, "protected": res.Protected}, nil
	})
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
