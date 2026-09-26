//go:build js && wasm

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"sort"
	"strings"
	"sync"
	"syscall/js"
	"time"
)

// httpFS is the read-only fs.FS of a font directory on the web (see the
// package documentation). ReadAt is served from the ranges fetched ahead
// and otherwise fetches the 64 KiB blocks it touches with Range requests;
// ReadFile fetches the whole file. Everything fetched is kept.
type httpFS struct {
	base  string
	files map[string]*httpFile
	names []string
}

type httpFile struct {
	fs   *httpFS
	name string
	size int64

	mu    sync.Mutex
	whole []byte
	segs  []segment // fetched ranges
}

type segment struct {
	off  int64
	data []byte
}

const blockSize = 64 << 10

func newHTTPFS(base string) (*httpFS, error) {
	if !strings.HasSuffix(base, "/") {
		base += "/"
	}
	data, status, err := fetch(base+"index.json", -1, 0)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("%sindex.json: HTTP %d", base, status)
	}
	var index []struct {
		Name string     `json:"name"`
		Size int64      `json:"size"`
		Scan [][2]int64 `json:"scan"`
	}
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("%sindex.json: %w", base, err)
	}
	h := &httpFS{base: base, files: map[string]*httpFile{}}
	var wg sync.WaitGroup
	for _, e := range index {
		if !fs.ValidPath(e.Name) || strings.Contains(e.Name, "/") || e.Size <= 0 {
			continue
		}
		f := &httpFile{fs: h, name: e.Name, size: e.Size}
		h.files[e.Name] = f
		h.names = append(h.names, e.Name)
		for _, r := range e.Scan {
			if r[0] < 0 || r[1] <= 0 || r[0]+r[1] > e.Size {
				continue
			}
			wg.Go(func() { f.fetchRange(r[0], r[1]) })
		}
	}
	sort.Strings(h.names)
	wg.Wait()
	return h, nil
}

// fetchRange fetches [off, off+n) and keeps it; a server that ignores the
// Range header sends the whole file, which is kept instead. It reports
// whether it got the bytes.
func (f *httpFile) fetchRange(off, n int64) bool {
	data, status, err := fetch(f.fs.base+f.name, off, off+n)
	if err != nil {
		return false
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	switch {
	case status == 206 && int64(len(data)) == n:
		f.segs = append(f.segs, segment{off, data})
		return true
	case status == 200 && int64(len(data)) == f.size:
		f.whole = data
		return true
	}
	return false
}

// cached copies [off, off+len(p)) from what was fetched, if it is there.
func (f *httpFile) cached(p []byte, off int64) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.whole != nil {
		copy(p, f.whole[off:])
		return true
	}
	for _, s := range f.segs {
		if off >= s.off && off+int64(len(p)) <= s.off+int64(len(s.data)) {
			copy(p, s.data[off-s.off:])
			return true
		}
	}
	return false
}

func (f *httpFile) ReadAt(p []byte, off int64) (int, error) {
	if off < 0 {
		return 0, errors.New("negative offset")
	}
	if off >= f.size {
		return 0, io.EOF
	}
	n := len(p)
	var err error
	if off+int64(n) > f.size {
		n, err = int(f.size-off), io.EOF
	}
	if f.cached(p[:n], off) {
		return n, err
	}
	start := off / blockSize * blockSize
	end := min((off+int64(n)+blockSize-1)/blockSize*blockSize, f.size)
	if !f.fetchRange(start, end-start) || !f.cached(p[:n], off) {
		return 0, fmt.Errorf("%s: cannot fetch bytes %d-%d", f.name, start, end-1)
	}
	return n, err
}

// readAll fetches the whole file (once).
func (f *httpFile) readAll() ([]byte, error) {
	f.mu.Lock()
	whole := f.whole
	f.mu.Unlock()
	if whole != nil {
		return whole, nil
	}
	data, status, err := fetch(f.fs.base+f.name, -1, 0)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("%s: HTTP %d", f.name, status)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.whole, f.segs = data, nil
	return data, nil
}

func (h *httpFS) file(op, name string) (*httpFile, error) {
	f, ok := h.files[name]
	if !ok {
		return nil, &fs.PathError{Op: op, Path: name, Err: fs.ErrNotExist}
	}
	return f, nil
}

func (h *httpFS) Open(name string) (fs.File, error) {
	if name == "." {
		return &dirHandle{fs: h}, nil
	}
	f, err := h.file("open", name)
	if err != nil {
		return nil, err
	}
	return &fileHandle{f: f}, nil
}

// ReadFile makes fs.ReadFile fetch a file in one request.
func (h *httpFS) ReadFile(name string) ([]byte, error) {
	f, err := h.file("read", name)
	if err != nil {
		return nil, err
	}
	data, err := f.readAll()
	if err != nil {
		return nil, err
	}
	return append([]byte(nil), data...), nil
}

func (h *httpFS) Stat(name string) (fs.FileInfo, error) {
	if name == "." {
		return fileInfo{name: ".", dir: true}, nil
	}
	f, err := h.file("stat", name)
	if err != nil {
		return nil, err
	}
	return fileInfo{name: f.name, size: f.size}, nil
}

func (h *httpFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if name != "." {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrNotExist}
	}
	out := make([]fs.DirEntry, len(h.names))
	for i, n := range h.names {
		out[i] = fs.FileInfoToDirEntry(fileInfo{name: n, size: h.files[n].size})
	}
	return out, nil
}

type fileHandle struct {
	f   *httpFile
	off int64
}

func (h *fileHandle) Stat() (fs.FileInfo, error) {
	return fileInfo{name: h.f.name, size: h.f.size}, nil
}
func (h *fileHandle) Close() error { return nil }
func (h *fileHandle) ReadAt(p []byte, off int64) (int, error) {
	return h.f.ReadAt(p, off)
}
func (h *fileHandle) Read(p []byte) (int, error) {
	n, err := h.f.ReadAt(p, h.off)
	h.off += int64(n)
	return n, err
}

type dirHandle struct {
	fs   *httpFS
	read bool
}

func (d *dirHandle) Stat() (fs.FileInfo, error) { return fileInfo{name: ".", dir: true}, nil }
func (d *dirHandle) Close() error               { return nil }
func (d *dirHandle) Read([]byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: ".", Err: errors.New("is a directory")}
}
func (d *dirHandle) ReadDir(n int) ([]fs.DirEntry, error) {
	if d.read {
		if n > 0 {
			return nil, io.EOF
		}
		return nil, nil
	}
	d.read = true
	return d.fs.ReadDir(".")
}

type fileInfo struct {
	name string
	size int64
	dir  bool
}

func (i fileInfo) Name() string { return i.name }
func (i fileInfo) Size() int64  { return i.size }
func (i fileInfo) Mode() fs.FileMode {
	if i.dir {
		return fs.ModeDir | 0o555
	}
	return 0o444
}
func (i fileInfo) ModTime() time.Time { return time.Time{} }
func (i fileInfo) IsDir() bool        { return i.dir }
func (i fileInfo) Sys() any           { return nil }

// fetch gets url, or the bytes [start, end) of it when start is not
// negative, and waits for the answer: the calling goroutine blocks while
// the browser goes on. It returns the HTTP status with the body.
func fetch(url string, start, end int64) ([]byte, int, error) {
	init := map[string]any{}
	if start >= 0 {
		init["headers"] = map[string]any{"Range": fmt.Sprintf("bytes=%d-%d", start, end-1)}
	}
	res, err := await(js.Global().Call("fetch", url, init))
	if err != nil {
		return nil, 0, fmt.Errorf("fetch %s: %w", url, err)
	}
	status := res.Get("status").Int()
	buf, err := await(res.Call("arrayBuffer"))
	if err != nil {
		return nil, 0, fmt.Errorf("fetch %s: %w", url, err)
	}
	u8 := js.Global().Get("Uint8Array").New(buf)
	data := make([]byte, u8.Get("length").Int())
	js.CopyBytesToGo(data, u8)
	return data, status, nil
}

// await waits for a Promise to settle.
func await(p js.Value) (js.Value, error) {
	done := make(chan struct{})
	var v js.Value
	var err error
	then := js.FuncOf(func(_ js.Value, args []js.Value) any {
		v = args[0]
		close(done)
		return nil
	})
	catch := js.FuncOf(func(_ js.Value, args []js.Value) any {
		err = js.Error{Value: args[0]}
		close(done)
		return nil
	})
	defer then.Release()
	defer catch.Release()
	p.Call("then", then, catch)
	<-done
	return v, err
}
