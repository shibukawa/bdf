# Fonts and licenses

BDF itself is under the [MIT License](https://github.com/shibukawa/bdf/blob/main/LICENSE). The fonts it holds or serves have licenses of their own.

## Fonts inside the library

| Font | Where it is | Source | License |
|---|---|---|---|
| **STIX Two Math** 2.13 b171 | Embedded in the package [`formula`](rendering.md#formulas) (`go:embed`, the whole font, unmodified): the formula font when `Options` name no other | [stipub/stixfonts](https://github.com/stipub/stixfonts), tag `v2.13b171`, `fonts/static_otf/STIXTwoMath-Regular.otf` | [SIL Open Font License 1.1](https://github.com/shibukawa/bdf/blob/main/formula/fonts/OFL.txt) |
| **Bravura** 1.482 | The outlines and metrics of music symbols, as a generated table in `converter/internal/music/smufl`, for the scores of MML, MIDI and MusicXML | [steinbergmedia/bravura](https://github.com/steinbergmedia/bravura) | [SIL Open Font License 1.1](https://github.com/shibukawa/bdf/blob/main/converter/internal/music/smufl/OFL.txt) |
| **NewStroke** | The stroke font KiCad draws text with, as a table in `converter/kicad` | vovanium's NewStroke release; [`NEWSTROKE.txt`](https://github.com/shibukawa/bdf/blob/main/converter/kicad/NEWSTROKE.txt) names the file | [CC0 1.0](https://creativecommons.org/publicdomain/zero/1.0/) |
| **DejaVu Sans**, Regular and Bold (subsets: ASCII) | In `fixture`, the package that builds the sample document (`bdf demo`) | [dejavu-fonts](https://dejavu-fonts.github.io/) | [DejaVu Fonts License](https://github.com/shibukawa/bdf/blob/main/fixture/testdata/fonts/LICENSE.txt) |

STIX Two Math: Copyright 2001-2021 The STIX Fonts Project Authors (https://github.com/stipub/stixfonts), with Reserved Font Name "TM Math". STIX Fonts™ is a trademark of The Institute of Electrical and Electronics Engineers, Inc.

Bravura: Copyright © 2015, Steinberg Media Technologies GmbH (http://www.steinberg.net/), with Reserved Font Name "Bravura".

The file of STIX Two Math in the repository is the one of the upstream tag, byte for byte; [`formula/fonts/README.md`](https://github.com/shibukawa/bdf/blob/main/formula/fonts/README.md) gives its SHA-256 and the command that checks it.

A program that imports `formula` holds STIX Two Math. The Open Font License allows bundling and embedding a font with any software, commercial or not, and asks that copies keep the copyright notice and the license — the font carries both in its name table, and `OFL.txt` is beside it. A program that lays formulas out with a font of its own (`formula.Options.Math`) still links the embedded one in.

The converters hold no font for the text of documents. They lay text out with the fonts they are given (`-font-dir`, the fonts of the system) and embed subsets of those in the documents they write; the licenses of those fonts are theirs.

## Fonts on this website

The demo pages convert files in the browser and fetch the fonts a document asks for from `fonts/`. They are published with their license files beside them:

| Font | Stands in for | License |
|---|---|---|
| Liberation Sans, Serif, Mono | Arial, Times New Roman, Courier New | SIL Open Font License 1.1 |
| Carlito | Calibri | SIL Open Font License 1.1 |
| Caladea | Cambria | the license in its copyright file (`fonts/fonts-crosextra-caladea.copyright.txt`) |
| IPAex Gothic, IPAex Mincho | Japanese text | IPA Font License Agreement v1.0 |
| DejaVu Sans, Serif, Mono | symbols | Bitstream Vera Fonts license; the DejaVu changes are in the public domain |
| STIX Two Math | Cambria Math, and formulas | SIL Open Font License 1.1 |

The files come from the Debian packages `fonts-liberation`, `fonts-crosextra-carlito`, `fonts-crosextra-caladea`, `fonts-ipaexfont-gothic`, `fonts-ipaexfont-mincho` and `fonts-dejavu-core`, and from the STIX release above; the [workflow that builds the site](https://github.com/shibukawa/bdf/blob/main/.github/workflows/pages.yml) lists them.

## Fonts of the tests

The test fonts in the repository — subsets of M PLUS 1p, STIX Two Math, STIX Two Text and DejaVu Sans — sit in `testdata` directories with their licenses.
