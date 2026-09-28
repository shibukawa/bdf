# Font files

A font file dropped on bdf isn't used to render something else — it's rendered *as itself*: the converter is a font previewer, turning `.ttf`, `.otf`, `.ttc`, `.woff` and `.woff2` files into pages that show what characters and glyphs a font has and what its OpenType features actually do. bdf runs on Canvas 2D, and `fillText` only applies a browser's own default features (`liga`, `kern`, …); there's no way for a document to ask for others. So a feature's effect gets drawn glyph by glyph, rather than by shaping real text.

## Try it

Drop a file on the [viewer](../../viewer/), or try a sample: [an OpenType font](../../viewer/?file=samples/stix.otf) (STIX Two Text).

## Font files (.ttf, .otf, .ttc, .woff, .woff2)

`converter/font` reads TrueType, OpenType, their collections (`.ttc`/`.otc`), and WOFF/WOFF2. It builds four scroll views per font — overview, characters, glyphs, and features (features is left out if the font has neither GSUB nor GPOS) — and a collection gets a full set of four views for each font it holds, numbered `-1`, `-2`, … and titled with that font's own name. `-param font=` picks which fonts of a collection to show; left unset, it shows all of them, or the first ones up to a total of 262,144 glyphs (some CJK collections bundle ten 65,535-glyph fonts).

- **Overview** — the font's name shown large, a specimen and its character set, sample text of each script it covers at sizes from 8 to 72 pt (that script's pangram, or `-param text=`), the `name` table in every language it carries, the license and OS/2 embedding permission, metrics, variation axes and named instances, the scripts and features of its layout tables, color palettes, and the table list.
- **Characters** — Unicode code charts by block, with each block's coverage.
- **Glyphs** — every glyph by ID, with its name, the characters that map to it, and its GDEF class.
- **Features** — every GSUB and GPOS feature, with the scripts and languages it applies to, whether a browser applies it without being asked, and what its lookups do, drawn glyph by glyph: substitutions (ligatures, alternates, small capitals, the lookups a contextual rule applies), kerning pairs before and after with their values, single adjustments with their advance boxes, and marks attached to base glyphs. `-param examples=` caps how many substitutions and pairs are shown per feature.

The GSUB/GPOS reader behind the features view is bdf's own (`internal/otlayout`), checked lookup by lookup against fontTools.

Glyphs are drawn by embedding the font itself twice. One embedding keeps the outlines but drops the layout tables and remaps `cmap` so each glyph sits at its own private-use codepoint — a shaper could otherwise still substitute something else for a lone character, and this is the only way to guarantee glyph *N* draws as glyph *N* (color glyphs included). The other embedding keeps the layout tables but prunes the glyph set down to what the sample text can reach, so the browser shapes that text itself (Arabic joining, Indic conjuncts, default kerning and ligatures). Variable fonts draw at their default instance, and a font whose OS/2 `fsType` forbids embedding is drawn as outlines instead — unless `-ignore-fstype` is given, and only with the rights to do so. The browser build carries no Brotli, so it can't read WOFF2 files.

| Option | Values | Default |
|---|---|---|
| `-param font=` | fonts of a collection to show, 1-based and comma-separated, or `all` | all, or the first ones of a collection of more than 262,144 glyphs |
| `-param text=` | sample text shown at several sizes instead of the pangram of the font's script | — |
| `-param examples=` | substitutions and pairs shown per feature, or `all` | 200 |
| `-ignore-fstype` | embed a font whose OS/2 `fsType` forbids embedding or subsetting | off |

See [design.md §3.28](../design.html#328-フォントファイル--bdf-変換器converterfontの構造) for the internals.
