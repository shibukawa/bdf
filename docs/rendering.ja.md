# 描画する

変換した文書と、LaTeX・MathML で書いた数式を、3 つの描画先に描く方法です。

| 描くもの | Go の `image.Image` | Ebitengine | ブラウザの Canvas |
|---|---|---|---|
| 文書のページ | [`raster/imagebdf`](#文書を-imageimage-に描く) | [`imagebdf` で描いた画像を `ebiten.Image` にする](#文書を-ebitengine-で描く) | [`@bdfkit/viewer`、`@bdfkit/render`](#文書をブラウザの-canvas-に描く) |
| 数式 | [`formula` → `imagebdf.Object`](#数式を-imageimage-に描く) | [`formula` → `raster/ebitenginebdf`](#数式を-ebitengine-で描く)（ベクターのまま） | [`formula` → 1 ページの BDF → `@bdfkit/render`](#数式をブラウザの-canvas-に描く) |

描き手はどれも同じ命令列（[spec §7](spec.md)）を実行します。

- **`raster/imagebdf`** — 純 Go のソフトウェアラスタライザ。文書のページ、範囲、単体の Object を `*image.RGBA` に描く
- **`raster/ebitenginebdf`** — パスだけの Object を Ebitengine の `vector` パッケージで描く。Ebitengine に依存するので別モジュール
- **`@bdfkit/render`** — ブラウザの Canvas 2D に描く TypeScript のパッケージ。`@bdfkit/viewer` はその上のビューア

## 文書を描く

### 文書を image.Image に描く

```go
import (
	"image/png"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

r, err := bdf.OpenSingleFile("report.bdf") // 変換した結果なら res.Doc をそのまま使う
if err != nil {
	return err
}
defer r.Close()
doc, err := r.ToDocument()
if err != nil {
	return err
}

ren := imagebdf.New(doc, nil)
img, err := ren.Page(doc.Views[0], 0, 2) // 1 ページ目を 1 単位 2 画素で。*image.RGBA
if err != nil {
	return err
}
err = png.Encode(w, img)
```

- `Page` は fixed・flow View のページ、`Region` はページの一部、scroll View の連続した範囲、シートの範囲を描きます。
- 名前で参照するフォント（HTML・Markdown・EPUB）と、埋め込みフォントにない字は `Options` の `FontFS`・`FontDirs` とシステムのフォントから探します。
- 描けなかったものは `ren.Warnings()` に入ります。サムネイルなら [`thumbnail`](api.md#go-サムネイルとページの画像thumbnailimagebdf) が切り抜き方まで決めます。

コマンドなら `bdf render -page 1 -scale 2 report.bdf page1.png` です。

### 文書を Ebitengine で描く

文書のページはテキスト、画像、グラデーション、クリップを含むので、`imagebdf` で画像にしてから `ebiten.Image` にします。

```go
import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/shibukawa/bdf/raster/imagebdf"
)

// 起動時、またはページを切り替えるとき
v := doc.Views[0]
scale := float64(screenHeight) / float64(v.Pages[0].H) // 画面の高さに合わせる
img, err := imagebdf.New(doc, nil).Page(v, 0, scale)
if err != nil {
	return err
}
page := ebiten.NewImageFromImage(img)

// Draw の中
var op ebiten.DrawImageOptions
op.GeoM.Translate(20, 20)
screen.DrawImage(page, &op)
```

画像なので、拡大するとぼけます。表示する大きさが変わったら、その倍率で描き直してください。

### 文書をブラウザの Canvas に描く

ページの配置、スクロール、ズーム、テキストの選択まで要るなら `@bdfkit/viewer` を使います（[フロントエンドに組み込む](integrating-to-frontend.ja.md)）。自分の canvas に 1 ページだけ描くなら `@bdfkit/render` の `PageRenderer` です。

```ts
import { BdfDocument, fetchSingle } from "@bdfkit/core";
import { PageRenderer } from "@bdfkit/render";

const doc = await BdfDocument.open(await fetchSingle("/documents/report.bdf"));
const renderer = new PageRenderer(doc, {}, document.fonts);

const page = doc.manifest.views[0].pages![0];
const scale = devicePixelRatio; // 1 単位あたりのデバイスピクセル
const canvas = document.querySelector<HTMLCanvasElement>("#page")!;
canvas.width = Math.ceil(page.w * scale);
canvas.height = Math.ceil(page.h * scale);
canvas.style.width = `${page.w}px`;
await renderer.renderPage(canvas.getContext("2d")!, page, { scale });
```

メインスレッドを止めたくない場合は、Worker で描く `BdfWorkerClient` を使います（[API 一覧](api.md#typescript-bdfkitrender)）。

## 数式を描く

`github.com/shibukawa/bdf/formula` は、LaTeX か MathML の数式を組み、**パスだけの Object**（グリフの輪郭、分数や根号の線）にします。フォントを参照しないので、パスを描ける描き手ならどれでも描けます。変換器が Word・PowerPoint・Excel・HTML・EPUB・Markdown・draw.io の数式を組むのと同じエンジンです。

```go
import (
	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/formula"
)

ts, err := formula.New(nil) // 同梱の STIX Two Math で組む
if err != nil {
	return err
}

f := formula.ParseTeX(`x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}`)
l := ts.Layout(f, formula.Style{
	Size:    24,                 // 文字の大きさ。結果の座標も同じ単位
	Display: true,               // ディスプレイスタイル（独立した行の数式）
	Color:   bdf.RGB(0, 0, 0),   // 省略すると黒。\textcolor などが上書きする
})

obj := l.Object() // *bdf.Object。原点はベースラインの左端、y は下向き
```

| 名前 | 内容 |
|---|---|
| `formula.New(*Options) (*Typesetter, error)` | 組版器。`Options.Math` は数式フォント（MATH 表を持つ TrueType・OpenType。省略すると同梱の STIX Two Math）、`Options.Text` は数式フォントにない字（`\text{}` の中の日本語など）のためのフォント。並行して使える |
| `formula.ParseTeX(src) *Formula` | LaTeX の数式モード（amsmath・amssymb の命令）。失敗しない。知らない命令はテキストとして残す |
| `formula.ParseMathML(src) (*Formula, error)` | `src` の最初の `math` 要素（Presentation MathML）。`display="block"` は `Formula.Display()` でわかり、`Layout` もディスプレイスタイルで組む |
| `(*Formula).Text()` / `(*Layout).Text()` | 線形表記（`x=(−b±√(b^2−4ac))/(2a)`）。代替テキスト、検索、コピー用 |
| `(*Typesetter).Layout(f, Style) *Layout` | 組む。`Style` は `Size`、`Display`、`Color`、`Bold` |
| `(*Layout).Width()` / `Height()` / `Depth()` | 送り幅、ベースラインから上、ベースラインから下。行の中に置くときに使う |
| `(*Layout).Ink()` | 実際に描く範囲（斜体のはみ出しやアクセントで箱より広いことがある） |
| `(*Layout).Object()` | パスだけの Object。境界は `Ink()` |
| `(*Layout).Draw(obj, x, y)` | 既存の Object に、ベースラインの左端を (x, y) にして描き足す |
| `(*Layout).Document(margin)` | 余白つきの 1 ページの文書。ページのテキストは線形表記 |

同梱の STIX Two Math は SIL Open Font License 1.1 で、`formula` を import したプログラムはこのフォントを含みます（[フォントとライセンス](licenses.ja.md)）。

どのフォントにもない字は描かれません（幅だけ空きます）。日本語などを使うときは `Options.Text` にフォントを渡してください。

### 数式を image.Image に描く

```go
import "github.com/shibukawa/bdf/raster/imagebdf"

img, err := imagebdf.Object(obj, 2, nil) // 1 単位 2 画素。*image.RGBA、背景は透明
if err != nil {
	return err
}
// img.Bounds() は原点（ベースラインの左端）からの画素。
// 行のベースラインが (x, y) なら、画像の左上は (x+Min.X, y+Min.Y) に置く
b := img.Bounds()
draw.Draw(dst, b.Add(image.Pt(x, y)), img, b.Min, draw.Over)
```

背景を塗るなら `&imagebdf.Options{Background: color.White}` を渡します。

### 数式を Ebitengine で描く

`raster/ebitenginebdf` は Object をパスのまま描くので、どの大きさでも輪郭がぼけません。

```sh
go get github.com/shibukawa/bdf/raster/ebitenginebdf
```

```go
import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/formula"
	"github.com/shibukawa/bdf/raster/ebitenginebdf"
)

// 起動時に 1 度だけ組んでコンパイルする
ts, _ := formula.New(nil)
l := ts.Layout(formula.ParseTeX(`e^{i\pi} + 1 = 0`), formula.Style{Size: 48, Color: bdf.RGB(255, 255, 255)})
drawing, err := ebitenginebdf.Compile(l.Object())
if err != nil {
	return err
}

// Draw の中
func (g *Game) Draw(screen *ebiten.Image) {
	var opts ebitenginebdf.DrawOptions
	opts.GeoM.Scale(g.zoom, g.zoom)  // 拡大・回転しても輪郭は滑らか
	opts.GeoM.Translate(40, 120)     // Object の原点（ベースラインの左端）を置く場所
	opts.ColorScale.ScaleAlpha(0.8)  // 色と透明度を掛ける
	opts.AntiAlias = true
	drawing.Draw(screen, &opts)
}
```

- `Compile` は Object を `vector.Path` にします。毎フレーム呼ぶのは `Draw` だけです。同じ `GeoM` で描き続ける間は、変換したパスも使い回します。
- 描けるのは単色のパスの塗りと線だけです。テキスト、画像、グラデーション、クリップ、グループ、マスク、破線を含む Object は `*ebitenginebdf.UnsupportedError` になります。数式の Object は常に描けます。
- ゲームの文字色に合わせるなら、白で組んでおいて `ColorScale` で色を付けます。
- パスの描画は画像の描画より重いので、たくさんの数式を大きさを変えずに出すなら、`imagebdf.Object` で画像にして `ebiten.NewImageFromImage` に渡すほうが軽くなります。

動くサンプルは [`raster/ebitenginebdf/example`](https://github.com/shibukawa/bdf/tree/main/raster/ebitenginebdf/example) にあります（文書のページと数式を並べて描きます）。

```sh
cd raster/ebitenginebdf
go run ./example
```

### 数式をブラウザの Canvas に描く

サーバーで数式を 1 ページの BDF にして返し、ブラウザは文書と同じ方法で描きます。数式 1 つで 2〜3 KB ほどです。

```go
// GET /formula?tex=...
func handler(w http.ResponseWriter, r *http.Request) {
	l := ts.Layout(formula.ParseTeX(r.URL.Query().Get("tex")), formula.Style{Size: 20, Display: true})
	w.Header().Set("Content-Type", "application/octet-stream")
	l.Document(8).WriteSingle(w) // 余白 8 単位のページ
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

ページは白い紙に描かれます。ページのテキストは線形表記なので、`@bdfkit/viewer` で開けば検索、コピー、読み上げもできます。

## 文書の中の数式

Word・PowerPoint・Excel の Office Math、HTML・EPUB の MathML、Markdown・draw.io の LaTeX は、変換器が同じエンジンで組んで文書に入れます。こちらはパスではなく、埋め込んだ数式フォントのテキストとして描きます（文書全体でグリフを共有するため）。書き方のサンプルは[対応形式](formats/index.ja.md)の各ページ、仕組みは[変換と出力 → 数式](architecture/conversion.ja.md#数式)にあります。
