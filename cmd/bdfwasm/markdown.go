//go:build js && wasm && previewonly

package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path"
	"syscall/js"

	"github.com/klauspost/compress/flate"
	"github.com/shibukawa/bdf/text/markdownbdf"
	base64 "github.com/shibukawa/tinygodriver/encoding/simdbase64"
)

func init() {
	methods["markdown"] = markdownCall
}

// markdownCall writes a single-file bdf as Markdown, as bdf markdown does
// (text/markdownbdf): the text, and its images as files beside it, as data
// URLs in it, or left out; and, when asked, a zip file of the text and the
// files.
func markdownCall(_ js.Value, args []js.Value) any {
	data, err := bdfInput(args)
	if err != nil {
		return reject(fmt.Errorf("markdown: %w", err))
	}
	password, name := str(args, "password"), str(args, "name")
	if name == "" {
		name = "document"
	}
	images := str(args, "images")
	switch images {
	case "", "files":
		images = "files"
	case "embed", "none":
	default:
		return reject(fmt.Errorf("markdown: images must be files, embed or none, not %q", images))
	}
	wantZip := option(args, "zip").Truthy()
	opts := &markdownbdf.Options{Views: strs(args, "views"), InternalLinks: option(args, "internalLinks").Truthy(),
		KeepRunning: option(args, "keepRunning").Truthy(), Tables: str(args, "tables")}
	if skip := strs(args, "skip"); skip != nil {
		opts.SkipRoles = skip
	}
	return promise(func() (any, error) {
		doc, err := readDocument(data, password)
		if err != nil {
			return nil, err
		}
		type file struct {
			name string
			data []byte
		}
		var files []file
		switch images {
		case "files":
			opts.Image = func(name string, data []byte) (string, error) {
				files = append(files, file{name, data})
				return name, nil
			}
			opts.File = opts.Image
		case "embed":
			opts.Image = func(file string, data []byte) (string, error) {
				return "data:" + mimeOf(file) + ";base64," + base64.StdEncoding.EncodeToString(data), nil
			}
		}
		var md bytes.Buffer
		res, err := markdownbdf.Write(&md, doc, opts)
		if err != nil {
			return nil, err
		}
		out := map[string]any{"markdown": bytesToJS(md.Bytes()), "views": res.Views, "tables": res.Tables, "warnings": anys(res.Warnings)}
		list := make([]any, 0, len(files))
		for _, f := range files {
			list = append(list, map[string]any{"name": f.name, "data": bytesToJS(f.data)})
		}
		out["images"] = list
		if wantZip {
			var zb bytes.Buffer
			zw := zip.NewWriter(&zb)
			zw.RegisterCompressor(zip.Deflate, func(w io.Writer) (io.WriteCloser, error) { return flate.NewWriter(w, flate.BestSpeed) })
			add := func(file string, data []byte) error {
				w, err := zw.Create(file)
				if err != nil {
					return err
				}
				_, err = w.Write(data)
				return err
			}
			if err := add(name+".md", md.Bytes()); err != nil {
				return nil, err
			}
			for _, f := range files {
				if err := add(f.name, f.data); err != nil {
					return nil, err
				}
			}
			if err := zw.Close(); err != nil {
				return nil, err
			}
			out["zip"] = bytesToJS(zb.Bytes())
		}
		return out, nil
	})
}

// bytesToJS copies bytes into a Uint8Array.
func bytesToJS(b []byte) js.Value {
	out := js.Global().Get("Uint8Array").New(len(b))
	js.CopyBytesToJS(out, b)
	return out
}

// mimeOf is the media type of an image file, by its extension.
func mimeOf(file string) string {
	switch path.Ext(file) {
	case ".png":
		return "image/png"
	case ".jpg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".avif":
		return "image/avif"
	case ".bmp":
		return "image/bmp"
	case ".svg":
		return "image/svg+xml"
	}
	return "application/octet-stream"
}
