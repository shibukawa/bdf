# otf: font files mapped into memory

`github.com/shibukawa/bdf/contrib/otf` opens OpenType font files, TrueType-flavored (.ttf, .ttc) and CFF-flavored (.otf, .otc) alike, for reading with little memory. It depends on the standard library only.

A font file is mostly glyph data. A Japanese font holds twenty thousand glyphs or more, of which a document uses a few hundred, and a collection (.ttc) holds several fonts of which one is wanted. Reading such a file with `os.ReadFile` to lay a page out keeps all of it on the heap for as long as the font is held: one face of a 19 MB collection costs 19 MB. `otf.Open` maps the file into memory instead where the platform allows. The result looks like a `[]byte` of the whole file, so any parser that takes a byte slice works unchanged, but only the pages that are read take memory (the table directory, the small tables, the glyphs in use), and the kernel may drop them again under pressure since the file backs them.

```go
f, err := otf.Open("/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc")
if err != nil {
	return err
}
// f.Bytes() is the whole file; parse it as usual.
font, err := sfnt.ParseIndex(f.Bytes(), 0)
// Keep f (or let it live for the process) while the parsed font is in use;
// f.Close() unmaps the file and invalidates every slice into it.
```

`OpenFS(fsys, name)` is `Open` for an `fs.FS`: the file is mapped when the file system opens it as an `*os.File` (`os.DirFS` does) and read whole otherwise (fonts embedded in the program with `embed`, say).

## What is mapped, what is read

| Platform | `Open` of a regular, non-empty file |
|---|---|
| Unix (Linux, macOS, the BSDs, Solaris, AIX) | mapped with `mmap(2)`, `PROT_READ`, `MAP_SHARED` |
| Windows | mapped with `CreateFileMapping` and `MapViewOfFile`, read-only |
| js/wasm, wasip1 and anything else | read whole, in one allocation of the file's size |

Empty files, files that are not regular files, files larger than the platform can address, and files the mapping call refuses (some network and FUSE file systems) are read whole. `Mapped()` says which happened. The file descriptor is closed as soon as the mapping exists; the mapping outlives it.

## API

| | |
|---|---|
| `Open(path string) (*File, error)` | opens a file by path |
| `OpenFS(fsys fs.FS, name string) (*File, error)` | opens a file of an `fs.FS` |
| `(*File).Bytes() []byte` | the whole file, read-only; `nil` after `Close` |
| `(*File).Size() int` | the size of the file |
| `(*File).Mapped() bool` | whether the file is mapped rather than read |
| `(*File).Close() error` | unmaps or drops the bytes; closing twice is fine |

The bytes are read-only in both senses: the mapping is read-only, and writing to a read fallback would change nothing on disk. A mapped file that another process truncates while it is mapped faults when the truncated pages are read, which Go cannot recover from; fonts in system and font directories are not rewritten in place, which is the use this package is for. Do not map files that are being written.

## How it differs from existing packages

| Package | Shape of the result | Platforms without mapping | `fs.FS` | Notes |
|---|---|---|---|---|
| `otf` (this package) | `[]byte` of the whole file | read whole | yes | read-only; opened by path; falls back when mapping fails |
| `golang.org/x/exp/mmap` | `*ReaderAt` (`io.ReaderAt`, `Len`, `At`) | read whole | no | every `ReadAt` copies; a `[]byte` parser needs the file copied once more |
| `github.com/edsrzf/mmap-go` | `MMap` (`[]byte`) | not supported | no | the caller opens the `*os.File` and picks `PROT`/`MAP` flags; read-write and copy-on-write mappings; `Flush`, `Lock` |
| `os.ReadFile`, `fs.ReadFile` | `[]byte` | — | `fs.ReadFile` | the whole file on the heap |
| `golang.org/x/image/font/sfnt.ParseReaderAt` | a parsed font | — | — | lazy at the parser: reads tables through an `io.ReaderAt` on demand; a different approach to the same problem, tied to that parser |

`otf` returns a slice rather than an `io.ReaderAt` because the parsers of bdf (`internal/sfnt`, `internal/cff`, `woff2`), like most font parsers, index a byte slice: a table is a sub-slice of the file, and a glyph is a sub-slice of a table, without copying. With a `ReaderAt`, each table or glyph would be copied out, or the parsers rewritten to read through the reader. Mapping keeps the slice API and still pays only for the pages touched.

## How it relates to bdf's own packages

`internal/fontdb` (font lookup, measurement and subsetting for the converters and the rasterizer) opened the fonts it loaded with `os.ReadFile` and kept the bytes for the life of the process; it now opens them with `otf.Open` or `otf.OpenFS` and keeps the `*otf.File` instead. Everything downstream (`internal/sfnt`, `internal/cff`, `woff2`, `raster/imagebdf`) takes the same `[]byte` and is unchanged. Converting the sample documents with `bdf generate -thumbnail` on macOS, with the system fonts (Hiragino collections of 7–19 MB):

| Input | Peak RSS, before → after | Peak footprint, before → after |
|---|---:|---:|
| basic.docx | 73 → 45 MiB | 60 → 28 MiB |
| basic.pptx | 53 → 41 MiB | 39 → 26 MiB |
| basic.xlsx | 54 → 38 MiB | 40 → 23 MiB |

The outputs are byte-identical. The package lives in the bdf module, under `contrib/`, because `fontdb` imports it: a module of its own would need a published version before each change to it could be used here. Importing it alone (`go get github.com/shibukawa/bdf/contrib/otf`) adds no dependency beyond the bdf module itself, since it uses only the standard library.
