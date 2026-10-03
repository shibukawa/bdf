//go:build js && wasm && previewonly

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"syscall/js"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/raster/imagebdf"
	"github.com/shibukawa/bdf/thumbnail"
)

func init() {
	methods["thumbnail"] = thumbnailCall
	methods["text"] = textCall
}

// maxThumbnail is the largest size a thumbnail call takes: its canvas holds
// 16 bytes a pixel.
const maxThumbnail = 2048

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

func thumbnailCall(_ js.Value, args []js.Value) any {
	data, err := bdfInput(args)
	if err != nil {
		return reject(fmt.Errorf("thumbnail: %w", err))
	}
	opts := &thumbnail.Options{Size: thumbnail.DefaultSize, View: str(args, "view")}
	if v := option(args, "size"); v.Type() == js.TypeNumber {
		opts.Size = v.Int()
	}
	if opts.Size < 1 || opts.Size > maxThumbnail {
		return reject(fmt.Errorf("thumbnail: size must be from 1 to %d", maxThumbnail))
	}
	if v := option(args, "sheetDpi"); v.Type() == js.TypeNumber {
		opts.SheetDPI = v.Float()
	}
	mode, err := thumbnail.ParseMode(str(args, "mode"))
	if err != nil {
		return reject(err)
	}
	opts.Mode = mode
	format := str(args, "format")
	switch format {
	case "":
		format = thumbnail.PNG
	case thumbnail.PNG, thumbnail.JPEG, thumbnail.WebP: // WebP fails in a build without encoders (bdf_noconv)
	default:
		return reject(errors.New("thumbnail: format must be png, jpeg or webp"))
	}
	password, fontURL := str(args, "password"), str(args, "fonts")
	return promise(func() (any, error) {
		doc, err := readDocument(data, password)
		if err != nil {
			return nil, err
		}
		// no system fonts in a browser: fonts referred to by name come from the font directory
		opts.Raster = imagebdf.Options{NoSystemFonts: true}
		if fontURL != "" {
			fsys, err := fonts(fontURL)
			if err != nil {
				return nil, fmt.Errorf("fonts: %w", err)
			}
			opts.Raster.FontFS = fsys
		}
		res, err := thumbnail.Make(doc, opts)
		if err != nil {
			return nil, err
		}
		var b bytes.Buffer
		if err := thumbnail.Encode(&b, res.Image, format); err != nil {
			return nil, err
		}
		img := js.Global().Get("Uint8Array").New(b.Len())
		js.CopyBytesToJS(img, b.Bytes())
		return map[string]any{"image": img, "format": format, "width": res.Image.Rect.Dx(), "height": res.Image.Rect.Dy(),
			"mode": res.Mode.String(), "warnings": anys(res.Warnings)}, nil
	})
}

func textCall(_ js.Value, args []js.Value) any {
	data, err := bdfInput(args)
	if err != nil {
		return reject(fmt.Errorf("text: %w", err))
	}
	password := str(args, "password")
	return promise(func() (any, error) {
		doc, err := readDocument(data, password)
		if err != nil {
			return nil, err
		}
		st, err := doc.SearchText()
		if err != nil {
			return nil, err
		}
		// as bdf text writes it
		var b bytes.Buffer
		enc := json.NewEncoder(&b)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		if err := enc.Encode(st); err != nil {
			return nil, err
		}
		return map[string]any{"json": b.String()}, nil
	})
}
