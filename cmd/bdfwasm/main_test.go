//go:build js && wasm

package main

import (
	"strings"
	"syscall/js"
	"testing"
)

// TestArguments calls every method with arguments that are not the bytes
// it wants: each call is a Promise that rejects. Run it with
//
//	GOOS=js GOARCH=wasm go test -exec "$(go env GOROOT)/lib/wasm/go_js_wasm_exec" ./cmd/bdfwasm
//
// (and -tags previewonly for the methods of the preview module).
func TestArguments(t *testing.T) {
	bytes := js.Global().Get("Uint8Array").New(100)
	for name, fn := range methods {
		call := guarded(name, fn)
		for what, args := range map[string][]js.Value{
			"no argument":         nil,
			"undefined":           {js.Undefined()},
			"null":                {js.Null()},
			"a string":            {js.ValueOf("bytes")},
			"an object":           {js.ValueOf(map[string]any{})},
			"an array":            {js.ValueOf([]any{1, 2, 3})},
			"an ArrayBuffer":      {bytes.Get("buffer")},
			"an object of length": {js.ValueOf(map[string]any{"length": 8})},
			"a Float64Array":      {js.Global().Get("Float64Array").New(8)},
			"bytes of no format":  {bytes, js.ValueOf(map[string]any{"size": 16, "format": "png"})},
			"options of no kind":  {bytes, js.ValueOf(map[string]any{"size": "large", "format": 1, "password": js.Null(), "view": js.Undefined()})},
		} {
			p, ok := call(js.Undefined(), args).(js.Value) // ends the program without the checks
			if !ok || !p.InstanceOf(js.Global().Get("Promise")) {
				t.Fatalf("%s with %s: no Promise", name, what)
			}
			if _, err := await(p); err == nil {
				t.Errorf("%s with %s: no error", name, what)
			} else if strings.Contains(err.Error(), "panic") {
				t.Errorf("%s with %s: %v", name, what, err)
			}
		}
	}
	// a panic of the part that reads the arguments rejects
	p := guarded("test", func(js.Value, []js.Value) any { return js.Undefined().Get("length") })(js.Undefined(), nil).(js.Value)
	if _, err := await(p); err == nil || !strings.Contains(err.Error(), "test: panic") {
		t.Errorf("a method that panics: %v", err)
	}
}
