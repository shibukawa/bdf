//go:build js && wasm && formonly

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"
	"syscall/js"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/form" // (registers the converters of the templates: Word, Excel, PowerPoint and PDF)
	"github.com/shibukawa/bdf/form/tag"
)

func init() {
	methods["fields"] = fieldsCall
	methods["fill"] = fillCall
}

// Bound the JSON, picture count, and total picture bytes of a call.
const (
	maxFormData     = 64 << 20
	maxFormFiles    = 256
	maxFormFileData = 64 << 20
)

// jsonData reads an option that is the JSON of an object ("" for none).
func jsonData(args []js.Value, name string) (map[string]any, error) {
	s := str(args, name)
	if s == "" {
		return nil, nil
	}
	if len(s) > maxFormData {
		return nil, fmt.Errorf("%s: more than %d MiB of JSON", name, maxFormData>>20)
	}
	var v any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	m, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: the data is an object of JSON", name)
	}
	return m, nil
}

// fieldsCall reads the tags of a template: the names they look up, and the
// JSON Schema of the data (form.Schema), which a page makes a form of.
func fieldsCall(_ js.Value, args []js.Value) any {
	data, err := input(args, "template bytes")
	if err != nil {
		return reject(fmt.Errorf("fields: %w", err))
	}
	sample, err := jsonData(args, "sample")
	if err != nil {
		return reject(fmt.Errorf("fields: %w", err))
	}
	return promise(func() (any, error) {
		fields, err := form.Fields(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return nil, err
		}
		if fields == nil {
			fields = []tag.Field{}
		}
		out := map[string]any{}
		// the format of the template, and the sample that it has for itself
		// (a compiled template, a template of pdfme)
		format, own, err := form.Describe(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return nil, err
		}
		out["format"] = format
		if bytes.HasPrefix(data, bdf.Magic[:]) {
			// a partly filled form: the names it leaves open, and what it holds
			def, _, err := form.ReadTemplate(bytes.NewReader(data), int64(len(data)))
			if err != nil {
				return nil, err
			}
			if def.IsPartial() {
				b, err := json.Marshal(def.Filled)
				if err != nil {
					return nil, err
				}
				out["open"], out["filled"] = anys(def.Open), string(b)
			}
		}
		if own != nil {
			b, err := json.Marshal(own)
			if err != nil {
				return nil, err
			}
			out["sample"] = string(b)
		}
		list, err := json.Marshal(fields)
		if err != nil {
			return nil, err
		}
		// Reuse the fields already read: SchemaOf would parse an Office or
		// PDF template a second time. A compiled template can carry a schema
		// with its own constraints, which SchemaOf preserves.
		var schema []byte
		if bytes.HasPrefix(data, bdf.Magic[:]) {
			schema, err = form.SchemaOf(bytes.NewReader(data), int64(len(data)), sample)
		} else {
			schema, err = form.Schema(fields, sample)
		}
		if err != nil {
			return nil, err
		}
		out["fields"], out["schema"] = string(list), string(schema)
		return out, nil
	})
}

// fillCall fills a template with data, as bdf fill does.
func fillCall(_ js.Value, args []js.Value) any {
	req, err := readRequest(args)
	if err != nil {
		return reject(fmt.Errorf("fill: %w", err))
	}
	data, err := jsonData(args, "data")
	if err != nil {
		return reject(fmt.Errorf("fill: %w", err))
	}
	if data == nil {
		data = map[string]any{}
	}
	files, err := fileOption(args)
	if err != nil {
		return reject(fmt.Errorf("fill: %w", err))
	}
	strict := option(args, "strict").Truthy()
	var open []string
	if v := option(args, "open"); v.Type() == js.TypeObject {
		if v.Length() > maxFormFiles {
			return reject(fmt.Errorf("fill: open: more than %d names", maxFormFiles))
		}
		for i := 0; i < v.Length(); i++ {
			if e := v.Index(i); e.Type() == js.TypeString {
				open = append(open, e.String())
			}
		}
	}
	return promise(func() (any, error) {
		opts, err := req.options()
		if err != nil {
			return nil, err
		}
		res, err := form.Fill(bytes.NewReader(req.data), int64(len(req.data)), &form.Options{Data: data, Strict: strict, Files: files, Convert: opts, Open: open})
		if err != nil {
			return nil, err
		}
		b, err := single(res.Doc)
		if err != nil {
			return nil, err
		}
		pages := 0
		for _, v := range res.Doc.Views {
			pages += len(v.Pages)
		}
		out := map[string]any{"bdf": b, "format": res.Format, "pages": pages, "warnings": anys(res.Warnings)}
		if res.Partial != nil {
			out["open"] = anys(res.Partial.Open)
		}
		return out, nil
	})
}

// fileOption reads the files option of a call: the pictures that the data
// names, by their paths.
func fileOption(args []js.Value) (memFS, error) {
	v := option(args, "files")
	if v.Type() != js.TypeObject {
		return nil, nil
	}
	keys := js.Global().Get("Object").Call("keys", v)
	if keys.Length() > maxFormFiles {
		return nil, fmt.Errorf("files: more than %d files", maxFormFiles)
	}
	out := memFS{}
	u8 := js.Global().Get("Uint8Array")
	total := 0
	for i := 0; i < keys.Length(); i++ {
		name := keys.Index(i).String()
		f := v.Get(name)
		if f.Type() != js.TypeObject || !f.InstanceOf(u8) {
			return nil, fmt.Errorf("files: %q is not a Uint8Array", name)
		}
		n := f.Get("length").Int()
		if n > tag.MaxPicture {
			return nil, fmt.Errorf("files: %q is larger than %d MiB", name, tag.MaxPicture>>20)
		}
		if n > maxFormFileData-total {
			return nil, fmt.Errorf("files: more than %d MiB in total", maxFormFileData>>20)
		}
		total += n
		name = strings.TrimPrefix(strings.ReplaceAll(name, "\\", "/"), "./")
		if !fs.ValidPath(name) {
			return nil, fmt.Errorf("files: %q is not the path of a file", name)
		}
		b := make([]byte, n)
		js.CopyBytesToGo(b, f)
		out[name] = b
	}
	return out, nil
}

// memFS holds the files of a call by their paths.
type memFS map[string][]byte

func (m memFS) Open(name string) (fs.File, error) {
	b, ok := m[name]
	if !ok {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	return &memFile{name: name, Reader: bytes.NewReader(b), size: int64(len(b))}, nil
}

func (m memFS) ReadFile(name string) ([]byte, error) {
	b, ok := m[name]
	if !ok {
		return nil, &fs.PathError{Op: "read", Path: name, Err: fs.ErrNotExist}
	}
	return b, nil
}

type memFile struct {
	*bytes.Reader
	name string
	size int64
}

func (f *memFile) Stat() (fs.FileInfo, error) { return f, nil }
func (f *memFile) Close() error               { return nil }
func (f *memFile) Name() string               { return f.name[strings.LastIndex(f.name, "/")+1:] }
func (f *memFile) Size() int64                { return f.size }
func (f *memFile) Mode() fs.FileMode          { return 0o444 }
func (f *memFile) ModTime() time.Time         { return time.Time{} }
func (f *memFile) IsDir() bool                { return false }
func (f *memFile) Sys() any                   { return nil }
