//go:build js && wasm && formonly

package main

import (
	"strings"
	"syscall/js"
	"testing"

	"github.com/shibukawa/bdf/form/tag"
)

func TestFormJSON(t *testing.T) {
	for _, text := range []string{`{} {}`, `{}]`, `{} trailing`, `[]`, `null`} {
		args := []js.Value{js.Undefined(), js.ValueOf(map[string]any{"data": text})}
		if _, err := jsonData(args, "data"); err == nil {
			t.Errorf("accepted data %q", text)
		}
	}
	args := []js.Value{js.Undefined(), js.ValueOf(map[string]any{"data": `{"items":[{"price":3.5}],"__proto__":"ordinary data"} `})}
	data, err := jsonData(args, "data")
	if err != nil || data["__proto__"] != "ordinary data" || data["items"].([]any)[0].(map[string]any)["price"] != 3.5 {
		t.Fatalf("data = %v, %v", data, err)
	}
}

func TestFormFilesTotalLimit(t *testing.T) {
	// Reuse one JS backing array; every individual picture is within its
	// limit, while their sum cannot be copied into the module's memory.
	file := js.Global().Get("Uint8Array").New(tag.MaxPicture)
	files := map[string]any{}
	for i := 0; i <= maxFormFileData/tag.MaxPicture; i++ {
		files[string(rune('a'+i))+".png"] = file
	}
	args := []js.Value{js.Undefined(), js.ValueOf(map[string]any{"files": files})}
	if _, err := fileOption(args); err == nil || !strings.Contains(err.Error(), "in total") {
		t.Fatalf("files total limit: %v", err)
	}
}
