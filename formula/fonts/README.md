# The formula font

`STIXTwoMath-Regular.otf` is STIX Two Math 2.13 b171, the font that package
`formula` lays formulas out with when it is given no other. It is embedded
in the package with `go:embed`, unmodified.

| | |
|---|---|
| Source | https://github.com/stipub/stixfonts, tag `v2.13b171`, `fonts/static_otf/STIXTwoMath-Regular.otf` |
| Download | https://raw.githubusercontent.com/stipub/stixfonts/v2.13b171/fonts/static_otf/STIXTwoMath-Regular.otf |
| SHA-256 | `3a5f3f26f40d5698b3c62dd085d48d6663696a3f80825aab8b553d5097518e8c` |
| Size | 838,652 bytes |
| License | SIL Open Font License 1.1: `OFL.txt`, the file of the same tag |
| Copyright | Copyright 2001-2021 The STIX Fonts Project Authors (https://github.com/stipub/stixfonts), with Reserved Font Name "TM Math". STIX Fonts™ is a trademark of The Institute of Electrical and Electronics Engineers, Inc. |

To check the file against its source:

    curl -fsSL https://raw.githubusercontent.com/stipub/stixfonts/v2.13b171/fonts/static_otf/STIXTwoMath-Regular.otf | shasum -a 256

A program that imports `formula` holds the font. The license allows
bundling and embedding it with any software; it asks that copies keep the
copyright notice and the license, which the font also carries in its name
table.
