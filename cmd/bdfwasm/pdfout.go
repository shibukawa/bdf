//go:build js && wasm && (previewonly || formonly)

package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/vector/pdfbdf"
)

func init() {
	methods["pdf"] = pdfCall
}

// bdfInput copies the single-file bdf of a call (data, options?).
func bdfInput(args []js.Value) ([]byte, error) { return input(args, "bdf bytes") }

// readDocument reads a single-file bdf, unlocking an encrypted one with
// password.
func readDocument(data []byte, password string) (*bdf.Document, error) {
	r, err := bdf.OpenSingle(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	if r.Locked() && password != "" {
		if err := r.Unlock(password); err != nil {
			return nil, err
		}
	}
	return r.ToDocument() // bdf.ErrLocked without the password
}

// strs returns an option that is an array of strings, without what is none.
func strs(args []js.Value, name string) []string {
	v := option(args, name)
	if v.Type() != js.TypeObject || !js.Global().Get("Array").Call("isArray", v).Bool() {
		return nil
	}
	var out []string
	for i := 0; i < v.Length() && i < 64; i++ {
		if s := v.Index(i); s.Type() == js.TypeString {
			out = append(out, s.String())
		}
	}
	return out
}

// pdfCall writes a single-file bdf as a PDF file, as bdf pdf does
// (vector/pdfbdf).
func pdfCall(_ js.Value, args []js.Value) any {
	data, err := bdfInput(args)
	if err != nil {
		return reject(fmt.Errorf("pdf: %w", err))
	}
	password, fontURL := str(args, "password"), str(args, "fonts")
	// no system fonts in a browser: fonts referred to by name come from the font directory
	opts := &pdfbdf.Options{Views: strs(args, "views"), SkipRoles: strs(args, "skip"), NoSystemFonts: true,
		Tagged: option(args, "tagged").Truthy(), PDFA: str(args, "pdfa")}
	return promise(func() (any, error) {
		doc, err := readDocument(data, password)
		if err != nil {
			return nil, err
		}
		if fontURL != "" {
			fsys, err := fonts(fontURL)
			if err != nil {
				return nil, fmt.Errorf("fonts: %w", err)
			}
			opts.FontFS = fsys
		}
		var b bytes.Buffer
		res, err := pdfbdf.Write(&b, doc, opts)
		if err != nil {
			return nil, err
		}
		out := js.Global().Get("Uint8Array").New(b.Len())
		js.CopyBytesToJS(out, b.Bytes())
		return map[string]any{"pdf": out, "pages": res.Pages, "fonts": res.Fonts, "warnings": anys(res.Warnings)}, nil
	})
}
