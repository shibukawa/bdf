# Generated codec package

`webpw` is libwebp (v1.6.0, encoder only) compiled to WebAssembly with
wasi-sdk and translated to pure Go with
[shibukawa/wasm2go-fork](https://github.com/shibukawa/wasm2go-fork) (branch
`pgmem`, v0.5.15-fork.7). No cgo and no wasm runtime are needed; building with
`-tags bdf_noconv` leaves it out entirely.

The translation uses the fork's options that keep a regenerated tree
diff-friendly: `-symbol-names` (functions are named after the wasm name
section, not their index), `-group-files` (one file per subject such as
`vp8_enc.go` instead of one multi-megabyte file), `-addr-consts` (static-data
addresses as named constants so a relinked data layout does not rewrite every
function body), and DCE rooted at `_initialize`, `encode`, `malloc` and `free`.

Regenerate with `tools/gen-codecs.sh`; the C entry point is `tools/webp/webp.c`
and the translator driver is `tools/gen-webp`.
