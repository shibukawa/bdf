//go:build js && wasm

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"syscall/js"
	"testing"
)

// TestHTTPFSFontURLs exercises the same JS fetch interface as the browser,
// without relying on a network connection or a particular CDN response.
func TestHTTPFSFontURLs(t *testing.T) {
	content := make([]byte, blockSize*2+7)
	for i := range content {
		content[i] = byte(i % 251)
	}
	for _, tc := range []struct {
		name, base, ref, source string
		ignoreRange             bool
	}{
		{name: "beside manifest", base: "https://app.example/fonts/", source: "https://app.example/fonts/Font.ttf"},
		{name: "relative directory without URL", base: "fonts/", source: "fonts/Font.ttf"},
		{name: "relative URL", base: "https://app.example/fonts/", ref: "../shared/remote.ttf", source: "https://app.example/shared/remote.ttf"},
		{name: "relative directory with URL", base: "fonts/", ref: "../shared/remote.ttf", source: "https://app.example/app/shared/remote.ttf"},
		{name: "parent directory with URL", base: "../fonts/", ref: "../shared/remote.ttf", source: "https://app.example/shared/remote.ttf"},
		{name: "absolute HTTPS URL", base: "https://app.example/fonts/", ref: "https://fonts.gstatic.com/s/example/v1/remote.ttf?version=1", source: "https://fonts.gstatic.com/s/example/v1/remote.ttf?version=1"},
		{name: "server ignores Range", base: "https://app.example/fonts/", ref: "https://fonts.gstatic.com/s/example/v1/remote.ttf", source: "https://fonts.gstatic.com/s/example/v1/remote.ttf", ignoreRange: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			entry := map[string]any{"name": "Font.ttf", "size": len(content), "scan": [][2]int64{{0, 4}}}
			if tc.ref != "" {
				entry["url"] = tc.ref
			}
			index, err := json.Marshal([]map[string]any{entry})
			if err != nil {
				t.Fatal(err)
			}
			previousLocation := js.Global().Get("location")
			js.Global().Set("location", map[string]any{"href": "https://app.example/app/index.html"})
			t.Cleanup(func() { js.Global().Set("location", previousLocation) })
			type request struct{ url, rangeHeader string }
			var requests []request
			mockFontFetch(t, func(location, rangeHeader string) ([]byte, int) {
				requests = append(requests, request{location, rangeHeader})
				if location == tc.base+"index.json" {
					return index, 200
				}
				if location != tc.source {
					return nil, 404
				}
				if rangeHeader == "" || tc.ignoreRange {
					return content, 200
				}
				var start, end int64
				if _, err := fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end); err != nil || start < 0 || end >= int64(len(content)) || end < start {
					return nil, 416
				}
				return content[start : end+1], 206
			})
			h, err := newHTTPFS(tc.base)
			if err != nil {
				t.Fatal(err)
			}
			if len(requests) != 2 || requests[1] != (request{tc.source, "bytes=0-3"}) {
				t.Fatalf("prefetch requests = %v", requests)
			}
			entries, err := fs.ReadDir(h, ".")
			if err != nil || len(entries) != 1 || entries[0].Name() != "Font.ttf" {
				t.Fatalf("directory entries = %v, %v", entries, err)
			}
			file, err := h.Open("Font.ttf")
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()
			info, err := file.Stat()
			if err != nil || info.Name() != "Font.ttf" || info.Size() != int64(len(content)) {
				t.Fatalf("file info = %v, %v", info, err)
			}
			if _, err := h.Open(tc.source); !errors.Is(err, fs.ErrNotExist) {
				t.Errorf("opening a URL as a file name: %v", err)
			}
			reader := file.(io.ReaderAt)
			first := make([]byte, 4)
			if n, err := reader.ReadAt(first, 0); err != nil || n != len(first) || !bytes.Equal(first, content[:4]) || len(requests) != 2 {
				t.Fatalf("reading prefetched bytes: %x, %d, %v; requests %v", first, n, err, requests)
			}
			middle := make([]byte, 7)
			const offset = blockSize + 11
			if n, err := reader.ReadAt(middle, offset); err != nil || n != len(middle) || !bytes.Equal(middle, content[offset:offset+len(middle)]) {
				t.Fatalf("reading another block: %x, %d, %v", middle, n, err)
			}
			wantRequests := 3
			if tc.ignoreRange {
				wantRequests = 2 // the prefetched request already returned the whole font
			} else if requests[len(requests)-1] != (request{tc.source, "bytes=65536-131071"}) {
				t.Errorf("block request = %v", requests[len(requests)-1])
			}
			if _, err := reader.ReadAt(middle, offset); err != nil || len(requests) != wantRequests {
				t.Fatalf("reading cached block: %v; requests %v", err, requests)
			}
			whole, err := fs.ReadFile(h, "Font.ttf")
			if err != nil || !bytes.Equal(whole, content) {
				t.Fatalf("reading whole font: %v", err)
			}
			if !tc.ignoreRange {
				wantRequests++
			}
			whole[0] ^= 0xff // ReadFile returns a copy, not the cache's backing array.
			again, err := fs.ReadFile(h, "Font.ttf")
			if err != nil || !bytes.Equal(again, content) || len(requests) != wantRequests {
				t.Fatalf("reading cached whole font: %v; requests %v", err, requests)
			}
		})
	}
}

func TestHTTPFSRejectsInvalidFontURLs(t *testing.T) {
	for _, ref := range []string{"https://", "https:font.ttf", "file:///font.ttf", "data:font/ttf;base64,AA==", ":bad", "https://fonts.example/%zz"} {
		t.Run(ref, func(t *testing.T) {
			index, err := json.Marshal([]map[string]any{{"name": "Font.ttf", "size": 10, "url": ref}})
			if err != nil {
				t.Fatal(err)
			}
			var unexpected bool
			mockFontFetch(t, func(location, rangeHeader string) ([]byte, int) {
				if location != "https://app.example/fonts/index.json" || rangeHeader != "" {
					unexpected = true
				}
				return index, 200
			})
			if _, err := newHTTPFS("https://app.example/fonts/"); err == nil || !strings.Contains(err.Error(), "font \"Font.ttf\"") {
				t.Fatalf("invalid URL %q: %v", ref, err)
			}
			if unexpected {
				t.Error("invalid font URL reached fetch")
			}
		})
	}
}

// mockFontFetch replaces only the HTTP boundary. Response.arrayBuffer and
// Promise scheduling remain the JavaScript runtime's real implementations.
func mockFontFetch(t *testing.T, handle func(location, rangeHeader string) ([]byte, int)) {
	t.Helper()
	previous := js.Global().Get("fetch")
	mock := js.FuncOf(func(_ js.Value, args []js.Value) any {
		rangeHeader := ""
		if len(args) > 1 {
			if headers := args[1].Get("headers"); !headers.IsUndefined() {
				rangeHeader = headers.Get("Range").String()
			}
		}
		data, status := handle(args[0].String(), rangeHeader)
		body := js.Global().Get("Uint8Array").New(len(data))
		js.CopyBytesToJS(body, data)
		response := js.Global().Get("Response").New(body, map[string]any{"status": status})
		return js.Global().Get("Promise").Call("resolve", response)
	})
	js.Global().Set("fetch", mock)
	t.Cleanup(func() {
		js.Global().Set("fetch", previous)
		mock.Release()
	})
}
