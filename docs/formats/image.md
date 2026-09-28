# Images, Photoshop

Four converters share this page because they all end up as a picture on a page rather than a laid-out document: plain images that browsers can already decode, Photoshop's composited artwork, scanned or faxed TIFF, and the pictures held in Windows metafiles. Each becomes one page (or one page per artboard, for Photoshop) sized to the image itself.

## Try it

Drop a file on the [viewer](../../viewer/), or try a sample: [an SVG image](../../viewer/?file=samples/drawing.svg), [a JPEG photo with EXIF](../../viewer/?file=samples/photo.jpg), or [a Photoshop file with three artboards](../../viewer/?file=samples/artboards.psd).

## Images (PNG, JPEG, GIF, WebP, AVIF, BMP, ICO, SVG)

`converter/image` passes through any image a browser's `<img>` element can already show. The file is stored as it is — no decoding, no re-encoding — in one page the size of the image, drawn by a single `IMAGE` instruction; opening the result looks exactly like opening the file in a browser. The converter reads only the image's size (from its own header — IHDR, the first SOF, the logical screen descriptor, VP8X, the DIB header, the largest ICO entry, or the HEIF `ispe`, honoring EXIF/`irot` orientation for JPEG, PNG and AVIF) and its metadata, which becomes the document's Dublin Core the same way it does for every other format: XMP first when present, then EXIF, then IPTC, then the format's own text fields (PNG text chunks, GIF comments, an SVG's `title`/`desc`). An animated image (APNG, animated GIF or WebP, an AVIF image sequence) shows the frame the browser decodes first, and converting one prints a warning.

SVG needs a workaround: `createImageBitmap` can't decode an SVG Blob inside a Worker, so the renderer keeps an SVG image as vector data and asks the page (which does have an `<img>` element) to rasterize it at whatever resolution the current zoom needs, caching a handful of rasters per image. That keeps SVG sharp at any zoom without running its scripts or fetching its external resources — the same restrictions a browser applies to an `<img>`.

There is no `-param` for this format; an image given directly as input is always stored as it is, regardless of the `-images` flag that controls how images *embedded in other documents* are stored.

See [design.md §3.19](../design.html#319-画像--bdf-変換器converterimageの構造) for the internals.

## Photoshop (.psd, .psb)

`converter/psd` reads the image Photoshop composited, not the layer structure — what a viewer sees is pixels Photoshop already rendered, since text and vector layers are stored as composited pixels too. A document without artboards becomes one page for the whole canvas; a document with artboards gets one page per visible artboard, cut from the composite in the order the layers panel lists them. Page size comes from the document's resolution. Every color mode and bit depth is read — RGB, CMYK, grayscale, Lab, indexed, and bitmap, at 8, 16 and 32 bits — as well as the large-document format (`.psb`). A file saved without "Maximize Compatibility" carries no composite at all, so the converter composites the layers itself: blend modes, masks, clipping masks, groups and fill layers, but not adjustment layers or layer effects (skipped with a warning). Pages finer than the shared resolution cap for image inputs are scaled down to it, the same as TIFF.

| Option | Values | Default |
|---|---|---|
| `-param artboards=` | `false`: the whole canvas on one page instead of a page per artboard | one page per artboard |
| `-max-dpi` | image inputs (TIFF, Photoshop): scale pages down to at most this many pixels per inch (`0`: no limit) | 192 |
| `-max-pixels` | image inputs (TIFF, Photoshop): scale pages down to at most this many pixels, width × height (`0`: no limit) | 14745600 |

See [design.md §3.18](../design.html#318-photoshop--bdf-変換器converterpsd) for the internals.

## TIFF (.tif, .tiff)

`converter/tiff` gives each page of the file a page sized to its resolution, drawn by one image; a multi-page TIFF is usually a scanner's or fax's output, image only, with no text. The reader is bdf's own — golang.org/x/image/tiff only reads the first IFD, and reads neither JPEG compression nor CCITT's two-dimensional coding — and it handles classic TIFF and BigTIFF, strips and tiles, uncompressed/PackBits/LZW/Deflate/JPEG, and CCITT fax coding (Group 3 one- and two-dimensional, with or without fill bits, and Group 4), at 1 to 16 bits in bilevel, gray, palette, RGB and CMYK. The `Orientation` tag rotates the page. Pages finer than the resolution cap (192 dpi and 3840 × 3840 pixels by default) are scaled down to it — a bilevel page stays bilevel rather than turning gray — and a JPEG page that needs no scaling is stored by joining its strips into one JPEG without re-encoding.

| Option | Values | Default |
|---|---|---|
| `-param dpi=` | resolution in pixels per inch of pages that do not give one | 96 |
| `-max-dpi` | image inputs (TIFF, Photoshop): scale pages down to at most this many pixels per inch (`0`: no limit) | 192 |
| `-max-pixels` | image inputs (TIFF, Photoshop): scale pages down to at most this many pixels, width × height (`0`: no limit) | 14745600 |

See [design.md §3.15](../design.html#315-tiff--bdf-変換器convertertiffの構造) for the internals.

## Windows metafiles (.emf, .wmf)

`converter/emf` turns an EMF or WMF file into one page the size of the picture, replaying its records with the same replayer bdf uses for EMF/WMF pictures embedded inside Office documents. Text is laid out and its fonts embedded the same way as for PowerPoint, so the font flags that apply there (`-font-dir`, `-fonts`, `-no-system-fonts`, `-no-subset`, `-no-woff2`) apply here too. There is no format-specific `-param` for metafiles.

See [design.md §3.5](../design.html#35-入力形式の登録と-emfwmf-変換器converteremf) for the internals.
