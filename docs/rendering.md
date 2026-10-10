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
