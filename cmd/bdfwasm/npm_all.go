//go:build js && wasm && npm_all

package main

// Keep the npm preset in sync with the registry of every converter.
import _ "github.com/shibukawa/bdf/converter/all"
