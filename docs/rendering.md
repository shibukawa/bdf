# Rendering

How to draw converted documents, and formulas written in LaTeX or MathML, on three kinds of surface.

| What | Go `image.Image` | Ebitengine | Browser canvas |
|---|---|---|---|
| A page of a document | [`raster/imagebdf`](#a-document-into-an-imageimage) | [an image drawn by `imagebdf`, as an `ebiten.Image`](#a-document-with-ebitengine) | [`@bdfkit/viewer`, `@bdfkit/render`](#a-document-on-a-browser-canvas) |
| A formula | [`formula` → `imagebdf.Object`](#a-formula-into-an-imageimage) | [`formula` → `raster/ebitenginebdf`](#a-formula-with-ebitengine), as vector paths | [`formula` → a one-page BDF → `@bdfkit/render`](#a-formula-on-a-browser-canvas) |

All three renderers run the same instructions ([spec §7](spec.md)).

- **`raster/imagebdf`** — a software rasterizer in pure Go. Draws pages, regions and single objects into an `*image.RGBA`
- **`raster/ebitenginebdf`** — draws objects made of paths with Ebitengine's `vector` package. A module of its own, since it depends on Ebitengine
- **`@bdfkit/render`** — the TypeScript package that draws on a browser's Canvas 2D; `@bdfkit/viewer` is the viewer built on it

## Documents

### A document into an image.Image

```go
import (
	"image/png"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

r, err := bdf.OpenSingleFile("report.bdf") // after a conversion, use res.Doc as it is
if err != nil {
	return err
}
defer r.Close()
doc, err := r.ToDocument()
if err != nil {
	return err
}

ren := imagebdf.New(doc, nil)
img, err := ren.Page(doc.Views[0], 0, 2) // the first page at 2 pixels per unit: an *image.RGBA
if err != nil {
	return err
}
err = png.Encode(w, img)
```

- `Page` draws a page of a fixed or flow view; `Region` draws part of a page, a stretch of a scroll view or a range of a sheet.
- Fonts referred to by name (HTML, Markdown, EPUB) and characters an embedded font lacks are looked up in `Options.FontFS`, `Options.FontDirs` and the system's fonts.
- What could not be drawn is listed by `ren.Warnings()`. For thumbnails, [`thumbnail`](api.md#go-サムネイルとページの画像thumbnailimagebdf) also decides the crop.

From the command line: `bdf render -page 1 -scale 2 report.bdf page1.png`.

### Raster quality and speed

`imagebdf.Options.Antialias` explicitly controls the edges of shapes, text and geometric clips. It applies to `Page`, `PageDPI`, `PageSize`, `Region`, `Object` and `SVGImage`; neither dpi nor scale changes the mode automatically.

| Mode | Behavior |
|---|---|
| `imagebdf.AntialiasDefault` (zero value) | The existing antialiased rendering at every resolution |
| `imagebdf.AntialiasOn` | Explicitly request the same antialiased rendering |
| `imagebdf.AntialiasOff` | Test pixel centers once per row and fill whole pixel spans, omitting edge coverage calculations for faster rendering |

For a high-resolution print image, select `AntialiasOff` and compare small text and thin rules at the final output size:

```go
ren := imagebdf.New(doc, &imagebdf.Options{Antialias: imagebdf.AntialiasOff})
img, err := ren.PageDPI(doc.Views[0], 0, 600)
```

The CLI equivalent is `bdf render -antialias=off -scale 8.333333 report.bdf print.png` (600 dpi). `-antialias=default|on|off` is available on `render`; `thumbnail` and generated previews retain their existing smoothing. In Go, a thumbnail can opt in with `thumbnail.Options.Raster.Antialias`.

The passes over a whole canvas (its background, fills of large areas, the products of large clips, the scanning of large paths and the conversion of the drawn page to an image) are split among up to eight goroutines, as many as the machine has; `Pages` divides them among its workers. The pixels are the same as when drawn alone.

### Monochrome output

For printing in black, `PageMono`, `RegionMono` and `PagesMono` take a `Mono` in place of the RGBA image, which is a quarter the size or less and faster to make:

| `Mono` | Image | Pixels |
|---|---|---|
| `imagebdf.Gray` | `*image.Gray` | the luminance (Rec. 601 weights, as `image/color` converts to gray) of each pixel over white, rounded |
| `imagebdf.Bilevel` | `*image.Paletted` of black and white (black first), which `image/png` writes at 1 bit a pixel | black where the gray value is below 128 |

```go
ren := imagebdf.New(doc, &imagebdf.Options{Antialias: imagebdf.AntialiasOff})
img, err := ren.PageMono(doc.Views[0], 0, 600.0/72, imagebdf.Bilevel) // an *image.Paletted
```

The page is drawn as for `Page` and converted at the end, so the antialiasing mode decides the edges: with `AntialiasOff` the bilevel image is exactly the drawing, with antialiasing the coverage of the edges is cut at the half, not dithered. The images are kept in the completed-image cache like the RGBA ones, apart from them. The CLI equivalent is `bdf render -mono=gray|bilevel`, for a page or a batch (`-pages`), as PNG or TIFF (8 bits or 1 bit a pixel), JPEG (in gray) or WebP (in colour).

Source transparency, blending, image interpolation, soft masks and blur remain enabled. Effects still depend on the geometry rendered with the selected mode, so their pixels can change too. Lines narrower than one pixel still use the existing one-pixel width and reduced alpha, so disabling edge antialiasing does not make every output pixel opaque. Low-resolution output can show jagged edges or lose small details; choose the mode for the intended output instead of relying on a fixed dpi threshold. An unknown Go `AntialiasMode` value keeps antialiasing enabled and adds a renderer warning; the CLI rejects unknown values.

### Drawing in bands

A page at 600 dpi is 35 million pixels: its canvas takes 560 MiB (16 bytes a pixel) and the RGBA image 140 MiB more. `PageRows` draws the page in bands of rows from the top instead, and hands each to a callback as it is drawn, so that a page is written to a file with the memory of a band:

```go
f, _ := os.Create("print.png")
w := bufio.NewWriter(f)
width, height, _ := imagebdf.PagePixels(v, 0, 600.0/72)
enc, _ := rowenc.NewPNG(w, width, height, rowenc.Bilevel, &rowenc.Options{DPI: 600})
err := ren.PageRows(v, 0, 600.0/72, imagebdf.Bilevel, func(y0 int, rows image.Image) error {
	return enc.WriteRows(rows) // an *image.Paletted of the rows from y0
})
// then enc.Close(), w.Flush(), f.Close()
```

The callback takes the first row of the band and the image of its rows, whose bounds are those rows of the page (`Min.Y` is `y0`): an `*image.RGBA` for a `Mono` of 0, else the image of the `Mono`. The image is reused for a later band, so copy what you keep. The bands after the first are drawn in parallel, up to eight at a time, by workers that each decode what they draw; the callback is called in order. `PagesRows` does the same for a selection of pages, one after the other. The pixels are those of `Page` and `PageMono`: each band is drawn with the transforms of the whole page, on a canvas of its own rows and of as many rows above and below as the shadows of the page reach (a band whose shadows reach past its margin is drawn again with as much), and a text run that crosses a band is rasterized whole and kept for the next band.

`Options.BandRows` sets the rows of a band (default `DefaultBandRows`, 128: a canvas of 10 MiB for an A4 page at 600 dpi). The instructions of the page are run once a band; the first band notes the extent of each draw, and the bands after it skip the draws that lie outside them, so a page is drawn in bands about as fast as whole. The bands after the first are drawn in parallel, each by a worker with a canvas of its own: `Options.BandWorkers` sets how many at a time (zero: the threads of the machine, up to eight, within `MaxBandCanvases`, 96 MiB, of canvases; one: the bands one after the other, with the memory of one). The workers share the text masks and each decodes the images it draws, so a page of large images takes a decoded copy a worker. A page whose shadows are blurred widely takes bands as tall as the shadows reach above and below, hundreds of rows at 600 dpi. The groups, soft masks and clips of a band take the memory of the band, not of the page, and the images are not kept in the completed-image cache.

A canvas, whole or a band, is lazy: it begins as its background without the background being written, and its pixels are written in blocks of 64 as the draws reach them. What is not drawn on (most of a page) is neither written nor read, and the conversion to the image writes the background for it, so a page of text costs a fraction of the memory passes of one filled first.

`image/rowenc` writes PNG and TIFF files from the rows as they come, for the three kinds of image (`RGBA`, `Gray`, `Bilevel`): the PNG compressed with zlib in IDAT chunks of a megabyte, the TIFF uncompressed in one strip with its directory first, so that neither seeks. The rows of a PNG are filtered and compressed in chunks of a fixed number of rows (about 256 KiB of them), each by a goroutine (`Options.Threads`, zero for those of the machine, up to eight): the stream is the streams of the chunks one after the other, each ended with a sync flush, which decodes as one, and the file is the same whatever the number of goroutines. At zlib level 9 a 600 dpi gray page compresses in a quarter of the time of one goroutine. `bdf render -compression=fast|normal|best` picks the level (1, 6 or 9); the default is `fast`, for an image that is printed or converted further: levels 1 and 6 make files of about the same size, and level 9 files a half to a third of that in one and a half to two times the time. The PNG of an RGBA image is written in colour type RGBA whether or not the page is opaque, since the rows to come are not known; `image/png` would write an opaque page as RGB. Both decode to the pixels of the whole image, and the file of a page written in bands is byte for byte that of the page written whole.

The CLI draws a page in bands by default when the output is a PNG or TIFF file (`.png`, `.tif`): `bdf render -mono=bilevel -scale 8.333333 report.bdf print.tif` writes a 600 dpi page with about 70 MiB of memory where drawing it whole took 1.1 GiB. `-band=N` sets the rows of a band and `-workers` how many bands are drawn at a time; `-band=0` draws the page whole, as JPEG and WebP output does.

### Cached redraws and parallel pages

Reuse a `Renderer` to reuse completed `Page`, `PageDPI`, `PageSize`, `Region` and `Renderer.Object` images and whole text-run coverage masks. Returned images are independent copies, so editing them cannot corrupt the cache. Retaining an image adds a copy on its first render; set `ImageCacheBytes: -1` for one-shot output.

`Options.ImageCacheBytes` defaults to 64MiB and `TextCacheBytes` to 32MiB. Zero selects the default, a negative value disables the cache, and a positive value sets a byte budget. LRU eviction accounts backing arrays, keys and conservative entry overhead. Decoded fonts, working canvases and Go allocator overhead are additional. Text masks cache coverage only; paint, alpha, shadow and compositing state are applied on each draw.

`Options.CanvasCacheBytes` (default 64MiB; a negative value keeps none) bounds the working canvases kept from one drawing for the next, 16 bytes a pixel of the image drawn and of each open group, and as much again the coverage of large masks kept from one fill for the next. A page at 600 dpi takes a canvas of hundreds of megabytes, which costs the system as much again to give back and fault in for each page: pages drawn one after the other with `Page` reuse it within this budget (set it to the size of a page, say `1 << 30`), and `Pages` keeps one canvas a worker for the length of a batch whatever the budget, trimming to the budget after.

Image keys include layout, region, pixel size, background and AA mode. Call `ren.Invalidate()` after changing part bytes or `FontFS`; it releases warnings, decoded resources, both caches and page workers. `FontDirs` and system fonts keep the existing process-wide scan/load cache; use `FontFS` for replaceable fonts. The package-level `imagebdf.Object` is a one-shot helper. To reuse object images, add the object to a document and call `ren.Object(hash, scale)`.

Use `Pages` / `PagesDPI` for multiple pages. Page numbers are zero-based; nil selects all pages. Independent workers draw pages, while callbacks run in selection order on the calling goroutine.

```go
err := ren.PagesDPI(ctx, doc.Views[0], nil, 300, 0,
    func(page int, img *image.RGBA) error {
        return savePage(page, img)
    })
```

Zero workers selects up to four; 1..16 sets a maximum. Concurrency also respects an aggregate 32M-pixel base-canvas budget, drawing a larger individual page alone (an A4 page at 600 dpi is one: it is drawn alone, with its passes split among the goroutines instead). Workers persist between calls, share the text-mask budget and divide the retained working-canvas allowance. Layers and output images retained by callbacks are additional. Cancellation and callback errors take effect between drawings/batches; already-started draws finish before return. Keep the renderer, document and font resources unchanged until return, including inside callbacks. `FontFS` must support concurrent reads for multiple workers.

The CLI supports `bdf render -pages=all -workers=4 -antialias=off report.bdf 'page-{page}.png'` and selections such as `-pages=1-3,5`. `{page}` expands to a one-based page number. `-page` and `-pages` are mutually exclusive; sheets are excluded. CLI output is consumed once, so it skips completed-image retention while reusing text masks.

### A document with Ebitengine

A page holds text, images, gradients and clips, so draw it into an image with `imagebdf` and hand that to Ebitengine.

```go
import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/shibukawa/bdf/raster/imagebdf"
)

// at start-up, or when the page changes
v := doc.Views[0]
scale := float64(screenHeight) / float64(v.Pages[0].H) // fit the height of the screen
img, err := imagebdf.New(doc, nil).Page(v, 0, scale)
if err != nil {
	return err
}
page := ebiten.NewImageFromImage(img)

// in Draw
var op ebiten.DrawImageOptions
op.GeoM.Translate(20, 20)
screen.DrawImage(page, &op)
```

It is an image, so it blurs when enlarged: draw it again at the new scale when the size it is shown at changes.

### A document on a browser canvas

For page layout, scrolling, zoom and text selection, use `@bdfkit/viewer` ([Integrating into a frontend](integrating-to-frontend.md)). To draw one page on a canvas of your own, use `PageRenderer` of `@bdfkit/render`.

```ts
import { BdfDocument, fetchSingle } from "@bdfkit/core";
import { PageRenderer } from "@bdfkit/render";

const doc = await BdfDocument.open(await fetchSingle("/documents/report.bdf"));
const renderer = new PageRenderer(doc, {}, document.fonts);

const page = doc.manifest.views[0].pages![0];
const scale = devicePixelRatio; // device pixels per unit
const canvas = document.querySelector<HTMLCanvasElement>("#page")!;
canvas.width = Math.ceil(page.w * scale);
canvas.height = Math.ceil(page.h * scale);
canvas.style.width = `${page.w}px`;
await renderer.renderPage(canvas.getContext("2d")!, page, { scale });
```

To keep the main thread free, draw in a worker with `BdfWorkerClient` ([API reference](api.md#typescript-bdfkitrender), in Japanese).

## Formulas

`github.com/shibukawa/bdf/image/formula` lays out a formula written in LaTeX or MathML and returns it as **an object of paths only**: the outlines of its glyphs and the rules of its fractions and radicals. It refers to no font, so any renderer that draws paths draws it. It is the engine the converters lay out the formulas of Word, PowerPoint, Excel, HTML, EPUB, Markdown and draw.io documents with.

```go
import (
	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/image/formula"
)

ts, err := formula.New(nil) // lays out with the embedded STIX Two Math
if err != nil {
	return err
}

f := formula.ParseTeX(`x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}`)
l := ts.Layout(f, formula.Style{
	Size:    24,               // the font size; the result is in the same units
	Display: true,             // display style: a formula on a line of its own
	Color:   bdf.RGB(0, 0, 0), // black when left out; \textcolor and the like override it
})

obj := l.Object() // a *bdf.Object: the origin is the left end of the baseline, y grows downwards
```

| Name | What it does |
|---|---|
| `formula.New(*Options) (*Typesetter, error)` | A typesetter. `Options.Math` is the formula font (TrueType or OpenType with a MATH table; the embedded STIX Two Math when left out), `Options.Text` are fonts for the characters the formula font lacks (Japanese in `\text{}`, say). Safe for concurrent use |
| `formula.ParseTeX(src) *Formula` | LaTeX math mode, with the commands of amsmath and amssymb. Never fails: an unknown command stays as text |
| `formula.ParseMathML(src) (*Formula, error)` | The first `math` element of `src` (presentation MathML). `display="block"` is reported by `Formula.Display()`, and `Layout` then uses display style |
| `(*Formula).Text()` / `(*Layout).Text()` | The linear notation (`x=(−b±√(b^2−4ac))/(2a)`), for alternative text, search and copy |
| `(*Typesetter).Layout(f, Style) *Layout` | Lays the formula out. `Style` has `Size`, `Display`, `Color` and `Bold` |
| `(*Layout).Width()` / `Height()` / `Depth()` | The advance, the height above the baseline and the depth below it: what placing a formula in a line of text needs |
| `(*Layout).Ink()` | The extent of what is drawn, which may reach past the box (italic overhangs, accents) |
| `(*Layout).Object()` | The object of paths; its bounding box is `Ink()` |
| `(*Layout).Draw(obj, x, y)` | Draws the formula into an existing object with the left end of its baseline at (x, y) |
| `(*Layout).Document(margin)` | A document of one page with a margin; the text of the page is the linear notation |

The embedded STIX Two Math is under the SIL Open Font License 1.1, and a program that imports `formula` holds it: see [Fonts and licenses](licenses.md).

A character that no font has is not drawn (its width is left empty). Pass a font in `Options.Text` for Japanese and other text.

### A formula into an image.Image

```go
import "github.com/shibukawa/bdf/raster/imagebdf"

img, err := imagebdf.Object(obj, 2, nil) // 2 pixels per unit: an *image.RGBA on a transparent background
if err != nil {
	return err
}
// img.Bounds() is in pixels from the origin (the left end of the baseline):
// with the baseline of the line at (x, y), the top left corner goes to (x+Min.X, y+Min.Y)
b := img.Bounds()
draw.Draw(dst, b.Add(image.Pt(x, y)), img, b.Min, draw.Over)
```

Pass `&imagebdf.Options{Background: color.White}` to fill the background.

### A formula with Ebitengine

`raster/ebitenginebdf` draws the object as paths, so its outlines stay sharp at every size.

```sh
go get github.com/shibukawa/bdf/raster/ebitenginebdf
```

```go
import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/image/formula"
	"github.com/shibukawa/bdf/raster/ebitenginebdf"
)

// once, at start-up: lay out and compile
ts, _ := formula.New(nil)
l := ts.Layout(formula.ParseTeX(`e^{i\pi} + 1 = 0`), formula.Style{Size: 48, Color: bdf.RGB(255, 255, 255)})
drawing, err := ebitenginebdf.Compile(l.Object())
if err != nil {
	return err
}

// in Draw
func (g *Game) Draw(screen *ebiten.Image) {
	var opts ebitenginebdf.DrawOptions
	opts.GeoM.Scale(g.zoom, g.zoom) // outlines stay smooth when scaled or rotated
	opts.GeoM.Translate(40, 120)    // where the origin of the object (the left end of the baseline) goes
	opts.ColorScale.ScaleAlpha(0.8) // scales the colors and the alpha
	opts.AntiAlias = true
	drawing.Draw(screen, &opts)
}
```

- `Compile` turns the object into `vector.Path`s; only `Draw` runs every frame, and it keeps the transformed paths while the `GeoM` stays the same.
- It draws fills and strokes of paths in plain colors and nothing else: an object with text, images, gradients, clips, groups, masks or dashed lines is an `*ebitenginebdf.UnsupportedError`. The objects of formulas always draw.
- To follow the text color of a game, lay the formula out in white and color it with `ColorScale`.
- Drawing paths costs more than drawing an image: for many formulas whose size does not change, draw them with `imagebdf.Object` and pass the images to `ebiten.NewImageFromImage`.

A runnable example is in [`raster/ebitenginebdf/example`](https://github.com/shibukawa/bdf/tree/main/raster/ebitenginebdf/example): the page of a document beside formulas.

```sh
cd raster/ebitenginebdf
go run ./example
```

### A formula on a browser canvas

The server returns the formula as a BDF of one page, and the browser draws it the way it draws a document. A formula takes 2–3 KB.

```go
// GET /formula?tex=...
func handler(w http.ResponseWriter, r *http.Request) {
	l := ts.Layout(formula.ParseTeX(r.URL.Query().Get("tex")), formula.Style{Size: 20, Display: true})
	w.Header().Set("Content-Type", "application/octet-stream")
	l.Document(8).WriteSingle(w) // a page with a margin of 8 units
}
```

```ts
import { BdfDocument, fetchSingle } from "@bdfkit/core";
import { PageRenderer } from "@bdfkit/render";

const tex = String.raw`x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}`;
const doc = await BdfDocument.open(await fetchSingle(`/formula?tex=${encodeURIComponent(tex)}`));
const page = doc.manifest.views[0].pages![0];
const scale = devicePixelRatio;
const canvas = document.querySelector<HTMLCanvasElement>("#formula")!;
canvas.width = Math.ceil(page.w * scale);
canvas.height = Math.ceil(page.h * scale);
canvas.style.width = `${page.w}px`;
await new PageRenderer(doc, {}, document.fonts).renderPage(canvas.getContext("2d")!, page, { scale });
```

The page is drawn on white paper. Its text is the linear notation, so opened with `@bdfkit/viewer` the formula can be searched, copied and read aloud.

## Formulas inside documents

Office Math in Word, PowerPoint and Excel, MathML in HTML and EPUB, and LaTeX in Markdown and draw.io are laid out by the converters with the same engine and become part of the document. There they are drawn as text in the embedded formula font rather than as paths, so that the whole document shares its glyphs. Each page under [Supported formats](formats/index.md) shows how to write them; [Conversion and output → Formulas](architecture/conversion.md#formulas) explains how it works.
