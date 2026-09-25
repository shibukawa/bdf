module github.com/shibukawa/bdf/tools/gen-webp

go 1.26

require github.com/goccy/wasm2go v0.0.0

// The wasm2go fork is consumed from a checkout (tools/gen-codecs.sh sets this).
replace github.com/goccy/wasm2go => ../../../bdf-codec-work/wasm2go-fork
