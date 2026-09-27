# API 一覧

bdf の公開 API をパッケージごとにまとめる。引数や細かい挙動は Go のパッケージコメント（`go doc`）と TypeScript の型定義に書いてある。フォーマットそのものは [spec.md](spec.md)、設計の理由は [design.md](design.md) を参照。

使う場所ごとの入口は次のとおり。

| やりたいこと | 使うもの |
|---|---|
| サーバーやバッチでファイルを bdf に変換する | Go の [`converter`](#go-変換converter) と形式ごとのパッケージ、または [`bdf` コマンド](#コマンドbdf) |
| ブラウザの中でファイルを bdf に変換する | [wasm の変換器](#ブラウザ内変換cmdbdfwasm)（`cmd/bdfwasm`） |
| bdf を自分で組み立てる・読む | Go の [`bdf`](#go-文書の組み立てと読み込みbdf) パッケージ |
| bdf をブラウザに表示する | [`@bdf/render`](#typescript-bdfrender) の Worker とテキスト層 |
| bdf を読んでテキストや構造を取り出す（Node でも） | [`@bdf/core`](#typescript-bdfcore) |
| サーバーでサムネイル・ページの画像・検索用のテキストを作る | Go の [`thumbnail` と `raster`](#go-サムネイルとページの画像thumbnailraster)、[`Document.SearchText`](#読み込み)、または [`bdf` コマンド](#コマンドbdf) |

## Go: 変換（converter）

`github.com/shibukawa/bdf/converter` は入力形式の登録簿と、形式に共通する変換の入口である。形式ごとの変換器はサブパッケージで、import したものだけが登録される（`converter/all` は全形式）。

```go
import (
	"github.com/shibukawa/bdf/converter"
	_ "github.com/shibukawa/bdf/converter/all" // 全形式を登録する
)

res, err := converter.ConvertFile("in.xlsx", "", &converter.Options{}) // "" は形式を内容から判定する
if err != nil {
	return err
}
f, _ := os.Create("out.bdf")
defer f.Close()
err = res.Doc.WriteSingle(f) // res.Doc.WriteSplit("out/") なら分割形式
```

| 名前 | 内容 |
|---|---|
| `Convert(r io.ReaderAt, size int64, name string, opts *Options) (*Result, error)` | 入力を変換する。`name` は形式名（`"pdf"`、`"pptx"` など）、`""` なら内容から判定し、だめなら `Options.FileName` の拡張子で決める。パスワード付きの Office 文書は `Options.Password` で復号してから判定する |
| `ConvertFile(path, name string, opts *Options) (*Result, error)` | ファイルを変換する（`FileName` も設定する） |
| `OpenStream(r, size, name, opts) (Stream, error)` | ページ単位の変換を始める。対応していない形式は丸ごと変換し、`Pages()` が 0 のストリームとして返す |
| `Stream` | `Outline()`（全ページの大きさだけでレイヤーが空の文書）、`Pages()`、`Page(i)`（そのページと、まだ返していない Part だけを持つ文書）、`Finish()`（残りを変換して `Convert` と同じ結果を返す）。詳細は design.md の「ページ単位の変換」 |
| `CheckPassword(r, size, password) (protected bool, err error)` | 変換せずにパスワードが要るか・合っているかを調べる。アップロード時にすぐ答えたいサーバー向け |
| `Detect(r, size) *Format` / `DetectFile(path)` | 形式を判定する（`nil` は不明） |
| `Formats() []*Format` / `Lookup(name) *Format` | 登録されている形式 |
| `Register(f *Format)` | 形式を登録する（変換器パッケージが `init` で呼ぶ） |
| `Pages` / `ParsePages(spec)` / `PageList(n...)` | ページ（スライド、シート）の選択。`"1-3,5,8-"` のような範囲を指定順に持ち、末尾までの範囲は総数を知らなくてよい。`Numbers(count)` で 1 始まりの番号にする |
| `ErrUnknownFormat` / `ErrPasswordRequired` / `ErrWrongPassword` | `errors.Is` で調べるエラー |

`Result` は `Doc *bdf.Document`、`Warnings []string`（`Options.Warn` が無いとき）、`Summary`（1 行の要約）、`Protected`（パスワードで開いた入力。同じパスワードで `bdf.NewPasswordLock` して暗号化するとよい）を持つ。

`Options` の項目（形式ごとに使うものだけを読む）:

| 項目 | 内容 |
|---|---|
| `Title` | 文書のタイトルを上書きする |
| `Pages Pages` | 変換するページ・スライド・シート（`ParsePages("1-3,5")`、`PageList(2)` など。`nil` は全部） |
| `Images imgconv.Options` | ラスター画像を再エンコードするか（ゼロ値はそのまま格納）と、解像度の上限 |
| `FontFS fs.FS` / `FontDirs []string` / `NoSystemFonts` | テキストをレイアウトする形式（Office 系、メタファイル）のフォントの探し先。`FontFS` が最初、次に `FontDirs`、最後にシステムのフォント |
| `SystemFonts` | フォントを埋め込まず名前で参照する |
| `NoSubset` / `NoWOFF2` / `IgnoreFSType` | フォントを subset しない / WOFF2 にしない / OS/2 の fsType が禁じていても埋め込む（権利がある場合だけ） |
| `NoTextIndex` | テキスト索引の Part を作らない |
| `Params map[string]string` | 形式ごとの設定（`Format.Params` に一覧。`bdf generate -h` でも表示される） |
| `Password` | パスワード付き入力を開くパスワード |
| `FileName` | 入力のファイル名（内容で判定できない CSV の判定、CSV と Parquet のシート名） |
| `Warn func(string)` | 警告を受け取る（無ければ `Result.Warnings` に集める） |

### 形式ごとのパッケージ

直接呼ぶと、`converter.Options` に無い形式固有の設定も使える。どのパッケージも `Convert` と `ConvertFile` を持ち、`Options` には `Title` と `Warn`、テキストをレイアウトする形式ならフォントまわり（`FontFS`、`FontDirs` など）と `NoTextIndex`、画像を持つ形式なら `Images` がある。ページの選択は `converter.Pages` 型である。

| パッケージ | 形式 | 固有の入口と設定 |
|---|---|---|
| `converter/pdf` | PDF | `Convert(rs io.ReadSeeker, opts)`、`NewStream(rs, opts)`（ページ単位の変換）。`Kind`（`fixed` か `flow`）、`NoSharePrefix`（共通プレフィックスを共有しない）、`NoAnnotations`、`Password` |
| `converter/pptx` | PowerPoint | `Slides`（`converter.Pages`）、`Hidden`（非表示スライドも含める） |
| `converter/xlsx` | Excel | `Sheets`（`converter.Pages`）、`Hidden`。`ConvertGrid(*Grid, opts)` は書式のない値の表（`Grid`）を Excel の新しいブックのように見せる |
| `converter/csv` | CSV・TSV | `Charset`、`Delimiter`、`Quote`、`Header`（`HeaderAuto` / `HeaderYes` / `HeaderNo`）、`TableStyle`、`Name`（シート名）。指定しなかったものは推測する |
| `converter/parquet` | Apache Parquet | `Rows`（表示する行数。0 は `DefaultRows`、-1 は全行）、`NoTypes`（型の行を付けない）、`TableStyle`、`Name`（シート名）。フッターが暗号化されたファイルは `ErrEncrypted` |
| `converter/docx` | Word | `Pages`、`Views`（`ViewsBoth` / `ViewsPages` / `ViewsScroll`） |
| `converter/visio` | Visio（.vsdx、.vdx） | `Pages` |
| `converter/drawio` | draw.io（.drawio、図を埋め込んだ .drawio.svg・.drawio.png） | `Convert(data []byte, opts)`（入力をバイト列で渡す）。`Pages`、`Border`（図の周りの余白） |
| `converter/dxf` | AutoCAD DXF | `Pages`、`Views`（`all` / `model` / `layouts`）、`Light`（モデル空間を白い紙に描く）。`Detect(head)` |
| `converter/gerber` | Gerber（RS-274X）・Excellon と、基板のファイルをまとめた ZIP | `Views`（`ViewsAll` / `ViewsBoard` / `ViewsLayers`）、`Mask`・`Silkscreen`・`Finish`（基板の表と裏の色）、`FileName`（1 つのファイルの層を名前から見分ける）。`Detect`、`IsGerber(head)`、`IsExcellon(head)` |
| `converter/tiff` | TIFF | `Pages`、`DPI`（解像度の無いページに仮定する値）、`Images`（解像度の上限もここ） |
| `converter/emf` | EMF・WMF | — |
| `converter/html` | HTML・XHTML・MHTML（リーダー表示） | `ConvertBytes(data, opts)`、`ConvertNode(*html.Node, opts)`。`Views`、`Extract`（`auto` / `article` / `none`）、`Dir`・`BaseURL`（参照の基準）、`NoRemote`・`Fetch`（ネットワークの画像）、`Width`・`FontSize`・`Font`・`MonoFont`、`EmbedFonts` |
| `converter/markdown` | Markdown | `ConvertBytes(data, opts)`、`ToHTML(data)`。`Options` は html と同じ |
| `converter/epub` | EPUB | `Views`（既定は `ViewsPages`）、`Paper`（`ParsePaper("b6")` など。既定は A5）、`Width`・`FontSize`・`Font`・`MonoFont`、`EmbedFonts`。暗号化された本は `ErrDRM`。`Result.FixedLayout`・`Vertical` |
| `converter/all` | 全形式を登録するだけ（`import _`） | — |

ページ単位の変換の使い方:

```go
s, err := converter.OpenStream(r, size, "", &converter.Options{})
if err != nil {
	return err
}
outline := s.Outline() // すぐ表示できる: 全ページの大きさ
send(outline)
for i := range s.Pages() { // 読み手が見ているページから順に呼んでよい
	page, err := s.Page(i) // ページ 1 枚と、新しい Part だけ
	if err != nil {
		return err
	}
	send(page)
}
res, err := s.Finish() // Convert と同じ完成した文書
```

## Go: 文書の組み立てと読み込み（bdf）

`github.com/shibukawa/bdf` はフォーマットそのものを扱う。変換器もこれで文書を作る。

### 書き出し

| 名前 | 内容 |
|---|---|
| `NewDocument() *Document` | 空の文書。`Meta`（Dublin Core の `DC`、`Source`）、`Views`、`CompressionLevel`、`MinCompress`、`Lock`（暗号化）を持つ |
| `(*Document).NewView(id, kind, title) *View` | View を足す。`kind` は `ViewFixed`、`ViewFlow`、`ViewSheet`、`ViewScroll` |
| `(*View).AddPage(w, h, layers...) *Page` | ページを足す（fixed・flow・scroll）。シートは `Tiles`、`Cols`、`Rows`、`Freeze` などの項目を直接設定する。右綴じの本は `Direction` を `DirectionRTL` にする |
| `(*Document).AddObject(*Object) (Hash, Rect)` | Object を格納し、ハッシュと外接矩形を返す（参照先の Part が先に要る） |
| `AddFont` / `AddImage` / `AddPaths` / `AddPart` | フォント、画像、パス集合、任意の Part を格納する。同じ内容は 1 度だけ格納される |
| `(*Document).BuildTextIndex(*View) (Hash, error)` | テキスト索引の Part を作って View に設定する |
| `(*Document).WriteSingle(io.Writer)` / `WriteSplit(dir)` | 1 ファイル形式・分割形式で書く |
| `NewPasswordLock(password, iter) (*Lock, error)` | `Document.Lock` に入れると、Part ごとに AES-256-GCM で封をして書く |

### 描画命令（Object）

`NewObject()` で作り、メソッドを連ねて命令を足す（どれも Canvas 2D の同名の操作に対応する）。

| 種類 | メソッド |
|---|---|
| 状態 | `Save`、`Restore`、`Transform`、`Translate`、`Scale`、`ClipRect`、`ClipPath` |
| スタイル | `FillColor`、`StrokeColor`、`FillPaint`、`StrokePaint`、`Line`、`Dash`、`Alpha`、`Blend`、`Shadow`、`Filter`、`Smoothing` |
| 図形 | `FillRect`、`StrokeRect`、`ClearRect`、`FillPath`、`FillPathAt`、`FillPathRun`、`StrokePath` |
| テキスト | `Font`、`TextStyle`、`FillText`、`StrokeText` |
| 画像 | `Image`、`ImageSub` |
| 合成 | `Use`、`UseAt`（子 Object）、`GroupBegin` / `GroupEnd`、`MaskBegin` / `MaskEnd` |
| メタ | `Mark`（構造、代替テキスト、言語）、`Link`、`Ext` |
| 資源の登録 | `AddPath`、`AddExtPath`、`AddPaint`（`LinearGradient`、`RadialGradient`、`ConicGradient`、`Pattern`）、`AddFont`（`EmbeddedFont`、`SystemFont`）、`AddImage`、`AddObject`、`UpdateFont`、`UpdateObject` |
| その他 | `SetBBox`、`Encode`、`Deps`、`RGB` / `RGBA`（色） |

### 読み込み

| 名前 | 内容 |
|---|---|
| `OpenSingle(r, size)` / `OpenSingleFile(path)` / `OpenSplit(dir)` | 文書を開く（`*Reader`）。`Manifest` を持つ |
| `(*Reader).Unlock(password)` | 暗号化された文書を開く（`Encrypted`、`Locked` で状態を調べる。`ErrLocked`、`ErrWrongPassword`） |
| `(*Reader).Part(h)` / `Object(h)` / `Entry(h)` | Part の中身、デコードした Object、Part の表の項目 |
| `(*Reader).ToDocument()` / `WriteSingle` / `WriteSplit` | 再エンコードせずに書き出す（1 ファイル形式と分割形式の相互変換） |
| `DecodeObject(data) (*ObjectPart, error)` | Object Part をデコードする。`Walk` / `Instructions` で命令を読む |
| `Disassemble(data) (string, error)` | Object を人が読める命令列にする |
| `ExtractText(o, resolve) ([]TextRun, error)` | Object（と子 Object）のテキストを取り出す |
| `DecodeTextIndex` / `EncodeTextIndex` / `PlainText` | テキスト索引の読み書きと、索引からの平文 |
| `(*Document).SearchText() (*SearchText, error)` | 検索エンジンに入れるテキスト。`Meta`（Dublin Core と入力形式）と、View ごとの `Pages`（`Page` は 1 始まりのページ番号、シートは 0 で丸ごと、`Text` は平文）。テキスト索引の Part があればそれを、なければ Object から取り出す。テキストのないページは除き、flow View のある文書の scroll View（同じ本文の別レイアウト）も除く。JSON にできる |
| `DecodePathCollection` / `EncodePathCollection` | パス集合の Part |
| `SharePrefixes(objs, minBytes)` | 複数の Object に共通する先頭部分を共有 Object に切り出す |
| `HashOf` / `ParseHash` | Part のハッシュ |

### その他のパッケージ

| パッケージ | 内容 |
|---|---|
| `imgconv` | 画像の格納方法。`Options{Mode: imgconv.Convert, Quality: 80}` で WebP を試して小さい方を残す（`Keep` はそのまま）。`MaxDPI`（既定 `DefaultMaxDPI` = 192）と `MaxPixels` はページ上の大きさが分かるラスター入力の解像度の上限。`Optimize(data, opts)`、`EncodePixels`、`Resize`、`Available()`（`bdf_noconv` ビルドでは false） |
| `woff2` | TrueType・OpenType を WOFF2 にする、WOFF2 を戻す。`Encode(font)`、`Decode(data)`（glyf・loca・hmtx の変換を戻す。フォントコレクションは扱わない）、`Available()`、`IsWOFF(data)` |

## Go: サムネイルとページの画像（thumbnail、raster）

`github.com/shibukawa/bdf/thumbnail` は文書のサムネイルを、`github.com/shibukawa/bdf/raster` は任意のページや範囲の画像を、ブラウザを使わずに Go で描く（design.md §3.25）。

```go
d, _ := r.ToDocument() // 変換した結果なら res.Doc
th, err := thumbnail.Make(d, &thumbnail.Options{Size: 256, Raster: raster.Options{FontDirs: dirs}})
if err != nil {
	return err
}
err = thumbnail.Encode(w, th.Image, thumbnail.PNG) // JPEG、WebP も
```

| 名前 | 内容 |
|---|---|
| `thumbnail.Make(doc, *Options) (*Result, error)` | サムネイルを描く。`Options` は `Size`（既定 256 px。切り抜きなら一辺、全体なら長辺）、`Mode`（`Auto`・`Crop`・`Fit`）、`View`（既定は最初の View）、`Raster`（フォントと背景）。`Result` は `Image`、選んだ `Mode`、`Warnings` |
| `Auto` の選び方 | Word・HTML・Markdown・Excel・CSV と、縦長のページの PDF・TIFF は `Crop`: 1 ページ目の左上から、ページの幅（横長ならページの高さ）の正方形。scroll View は先頭から、シートは A1 から、シートの短い辺（最大 `MaxSheetSide` = 480 単位）の正方形を枠線つきで。それ以外（PowerPoint、Visio、draw.io、CAD、プリント基板、Illustrator、Photoshop、画像、EPUB の表紙、横長の PDF・TIFF）は `Fit`: 1 ページ目の全体。判定は `Meta.Source` と View の種類による |
| `thumbnail.Encode(w, img, format)` / `FormatOf(name)` | `PNG`・`JPEG`（品質 85）・`WebP`（非可逆、品質 80）で書く / ファイル名の拡張子から形式を決める |
| `raster.New(doc, *Options) *Renderer` | 文書を描くレンダラ。`Options` は `FontFS`・`FontDirs`・`NoSystemFonts`（名前で参照するフォントと、埋め込みフォントにない字の探し先）と `Background`（既定は白）。デコードした Object・画像・フォントを保持する。並行には使えない |
| `(*Renderer).Page(v, page, scale)` | fixed・flow View のページ（scroll View の帯）を、1 単位 `scale` 画素で描く |
| `(*Renderer).Region(v, page, rect, w, h)` | 範囲を `w` × `h` 画素に描く。fixed・flow はそのページの座標、scroll View は帯を積んだ連続の座標、シートはシートの座標（枠線も描く） |
| `(*Renderer).Warnings()` | ビューアどおりに描けなかったもの（AVIF の画像、FILTER、フォントが見つからないテキストなど） |
| `raster.MaxPixels` / `ErrTooLarge` / `ErrNoPage` / `SheetSize(v)` | 描く画像の大きさの上限（64 M 画素）とそのエラー、ないページ、シートの大きさ |

描けるもの: パス（nonzero・evenodd、アンチエイリアス）、線（端・結合・破線。1 画素より細い線はビューアと同じく 1 画素幅で薄く）、クリップ、単色・線形・放射・扇形のグラデーションとパターン、画像（PNG、JPEG と EXIF の向き、GIF、BMP、WebP、SVG）、埋め込みフォントと名前で参照するフォントのテキスト（`advance` 補正、揃え、ベースライン、字間、太字・斜体の合成、字ごとのフォールバック、右から左の文字の並べ替え）、グループの透明度とブレンド、ソフトマスク、影。描かないもの: AVIF の画像、FILTER、ヒンティング、カーニング、合字、アラビア文字の字形の変化。

## コマンド（bdf）

`go run ./cmd/bdf <サブコマンド>`。暗号化された文書は `$BDF_PASSWORD` のパスワードで読む。

| サブコマンド | 内容 |
|---|---|
| `generate [flags] <input> <out.bdf \| dir/>` | 変換する（`/` で終わる出力は分割形式）。フラグと形式ごとの `-param` は `bdf generate -h` |
| `ls` / `manifest` | View と Part の一覧 / manifest の JSON |
| `disasm <file> <hash>` / `extract <file> <hash> <out>` | Object の逆アセンブル / Part の取り出し |
| `split` / `join` | 1 ファイル形式と分割形式の変換（パスワード不要） |
| `encrypt` / `decrypt` | パスワードで暗号化する / 暗号化を外す |
| `demo` | サンプル文書（testdata/demo.bdf と同じもの）を書く |
| `thumbnail [flags] <file> <out.png \| .jpg \| .webp>` | サムネイルを描く。`-size`（既定 256）、`-mode auto\|crop\|fit`、`-view`、`-font-dir`、`-no-system-fonts` |
| `text [flags] <file> [out.json]` | メタデータとページごとのテキストを JSON で書く（既定は標準出力） |
| `render [flags] <file> <out.png \| .jpg \| .webp>` | ページを描く。`-view`、`-page`（1 始まり）、`-scale`、シートは A1 からの `-width`・`-height` |

`generate` の `-thumbnail <file>`（`-thumbnail-size`、`-thumbnail-mode`）と `-text <file>`（`-` で標準出力）は、変換と同時にサムネイルとテキストを書く。サムネイルとテキストは暗号化されないので、暗号化した文書（`generate` では暗号化して書き出す文書。既定ではパスワード付きの入力のもの）については書かない。`-allow-plaintext` を付けたときだけ書く（`thumbnail`・`text`・`render` も同じ）。

## ブラウザ内変換（cmd/bdfwasm）

`cmd/bdfwasm` は変換器を `GOOS=js GOARCH=wasm` でビルドしたもので、Go の `wasm_exec.js` とともに Worker で動かす（デモサイトは `examples/viewer/convert-worker.ts`）。起動すると `globalThis.bdfConverter` を設定する。`-tags pdfonly`、`officeonly`、`webonly`、`imageonly` で PDF 用、Office 系用、HTML・Markdown・EPUB 用、画像用に分けてビルドできる。`previewonly` は変換器を持たず、bdf のサムネイルと検索用テキストを作るモジュールになる（`thumbnail` と `text`）。

| 名前 | 内容 |
|---|---|
| `formats` | `[{name, description, extensions}]` |
| `convert(data: Uint8Array, options?)` | 丸ごと変換する。`{bdf, format, summary, warnings, protected}` を返す Promise |
| `open(data, options?)` | ページ単位の変換を始める。`{bdf, format, pages, warnings, stream?}` を返す。`pages > 0` なら `bdf` は outline で、`stream` が残りを変換する。対応していない形式は丸ごと変換した結果（`pages: 0`、`summary` つき） |
| `stream.page(i)` | ページ i を変換し、`{bdf, warnings}`（そのページだけのページ文書と新しい警告）を返す |
| `stream.finish()` | 残りを変換し、`convert` と同じ形の完成した文書を返す |
| `stream.close()` | ストリームを手放す（`finish` のあとも呼ぶ） |
| `thumbnail(bdf: Uint8Array, options?)` | `previewonly` のみ。単一ファイル形式の bdf のサムネイルを描き、`{image, format, width, height, mode, warnings}` を返す。`options` は `{size?, mode?, format?, view?, password?, fonts?}`（`size` は既定 256・最大 2048、`mode` は `auto`・`crop`・`fit`、`format` は `png`・`jpeg`。`bdf_noconv` でなければ `webp` も） |
| `text(bdf: Uint8Array, options?)` | `previewonly` のみ。`bdf text` と同じ JSON（メタデータとページごとのテキスト）を `{json}` で返す。`options` は `{password?}` |

`options` は `{format?, password?, fonts?}`。`fonts` はフォントのディレクトリの URL で、`index.json` にファイルとフォントの走査が読む範囲を並べておく（`examples/viewer/site.mjs` が作る）。失敗した Promise の Error は、パスワードが要る・違う・形式が分からないときに `code` が `"password-required"`、`"wrong-password"`、`"unknown-format"` になる。返す文書は入力が暗号化されていても暗号化しない。

## TypeScript: @bdf/core

bdf を読むためのパッケージ（`packages/core`）。DOM に依存しないので Node でも動く。

### 文書を開く

| 名前 | 内容 |
|---|---|
| `BdfDocument.open(source, {password?})` | 文書を開く。暗号化された文書でパスワードが無い・違うときは `BdfPasswordError`（`reason` が `"required"` / `"wrong"`） |
| `BufferSource(bytes)` | メモリ上の 1 ファイル形式 |
| `RangeSource(url, init?)` | 1 ファイル形式を HTTP Range で必要な Part だけ取得する |
| `SplitSource(base, init?)` | 分割形式（`manifest.json` と `parts/<hash>`） |
| `fetchSingle(url, init?)` | 1 ファイル形式を丸ごと取得して `BufferSource` にする |
| `PartSource` | 上のソースの共通インターフェース（`manifest()`、`stored(entry)`）。自前の取得方法も実装できる |

`BdfDocument` のメソッド:

| 名前 | 内容 |
|---|---|
| `manifest` / `view(id)` | manifest と View |
| `ensure(hash, onResource?)` | Object と、それが参照する Object を再帰的に読み込む。フォント・画像・パス集合の Part は `onResource` に渡す |
| `object(hash)` / `objectSync(hash)` | Object（`objectSync` は読み込み済みのもの） |
| `part(hash)` / `entry(hash)` / `pathCollection(hash)` | 展開した Part、Part の表の項目、パス集合 |
| `textIndex(view)` | View のテキスト索引（索引の Part が無ければ全ページから作る） |
| `addPage(viewId, index, pageDoc)` | ページ単位の変換が返したページ文書のページを置き、その Part を加える |

### Object とテキスト

| 名前 | 内容 |
|---|---|
| `decodeObject(bytes)` / `walk(obj, sink)` / `NoopSink` | Object のデコードと、命令を `OpSink` に流すループ |
| `objectDeps(obj)` / `opHistogram(obj)` | 参照している Part / 命令ごとの数 |
| `decodePath` / `decodePathCollection` | パス |
| `extractText(obj, resolve, matrix?)` | テキストの run（位置、フォント、区切り）を取り出す |
| `extractContent(obj, resolve, matrix?)` | run に構造（見出し、リスト、表、図）とリンクを付けた `TextContent` |
| `parseCellRef` / `guessSep` / `Mark` / `Sep` | セル参照の解釈、run の区切りの推測、MARK と区切りの定数 |
| `TextSearch(runs)` / `decodeTextIndex(bytes)` | テキスト索引から検索する（NFKC、大文字小文字、小書きのかな、数式のマイナス記号とハイフンマイナスを同一視し、行をまたいで一致する）。`search(query, {limit, caseSensitive, context})` は `SearchHit[]` |
| `normalizeQuery` / `normalizeChar` | 検索と同じ正規化 |
| `dcValues(value)` | Dublin Core の値（文字列または配列）を配列にする |
| `parseHeader` / `decode` / `MAGIC` / `HEADER_SIZE` | 1 ファイル形式のヘッダと Part の展開 |
| `BdfFormatError` / `BdfPasswordError` / `SealedSource` | 形式の誤り、パスワードの誤り、暗号化された文書のソース |
| `Manifest`、`View`、`Page`、`Layer`、`PartEntry` など | manifest の型（spec.md §4） |

## TypeScript: @bdf/render

ブラウザで描くためのパッケージ（`packages/render`）。デコードと描画は Worker（`@bdf/render/worker` をバンドルしたもの）で行い、メインスレッドはビットマップとテキスト層を置く。

```ts
import { BdfWorkerClient, buildTextLayer, installCopyHandler, TEXT_LAYER_CSS } from "@bdf/render";

const client = new BdfWorkerClient(new Worker("./worker.js", { type: "module" })); // @bdf/render/worker をバンドルしたもの
const manifest = await client.open({ kind: "single", url: "doc.bdf", range: true });
const view = manifest.views[0];
const zoom = 1;

// ページ 1: Worker が描いたビットマップを canvas に置き、その上にテキスト層を重ねる
const bmp = await client.page(view.id, 0, zoom * devicePixelRatio);
canvas.width = bmp.width;
canvas.height = bmp.height;
canvas.getContext("2d")!.drawImage(bmp, 0, 0);
pageElement.append(buildTextLayer(await client.content(view.id, 0), zoom, { lang: "ja" }));

// テキスト層の CSS と、ページをまたいだコピー
document.head.append(Object.assign(document.createElement("style"), { textContent: TEXT_LAYER_CSS }));
installCopyHandler(container);
```

### BdfWorkerClient

| メソッド | 内容 |
|---|---|
| `open(source, password?, options?)` | 文書を開く。`source` は `{kind: "single", url, range?}`、`{kind: "split", base}`、`{kind: "buffer", buffer}`（転送される）。`options.imageBudget` はデコードした画像を保持するバイト数 |
| `unlock(password)` | パスワードが要る・違うと `open` が `BdfWorkerError`（`code` が `"password-required"` / `"wrong-password"`）で失敗したあと、同じ文書を別のパスワードで開く |
| `page(view, index, scale, roles?)` | ページを `ImageBitmap` に描く。`scale` は 1 単位あたりのデバイスピクセル |
| `continuous(view, viewport, scale)` | flow・scroll View を縦に続けた配置の矩形を描く |
| `sheet(view, viewport, scale)` | シートの矩形を描く（行・列見出しとウィンドウ枠の固定はビューアの仕事） |
| `content(view, index)` / `continuousContent(view, viewport)` / `sheetContent(view, viewports)` | テキスト層用の `TextContent`（run、構造、リンク） |
| `text(view, index)` / `continuousText(view, viewport)` | run だけ |
| `search(view, query, options?)` / `locate(view, hits)` | 検索と、ヒットの矩形 |
| `addPage(view, index, buffer)` | ページ単位の変換が返したページ文書を、開いている文書に置く |
| `replace(source)` | 同じ View とページを持つ文書（ページ単位の変換の完成品）に差し替える。実行中の要求は古い文書で終わる |
| `close()` / `terminate()` | 文書を閉じる / Worker を止める |

### テキスト層

| 名前 | 内容 |
|---|---|
| `buildTextLayer(content, scale, options?)` | 透明なテキスト層の要素を作る。選択、コピー、検索のハイライト位置合わせ、読み上げ（見出し、リスト、表、代替テキスト、リンク、言語）に使う。`options` は `lang`、`sheet`（シートを表として並べる）、`selectable`（`false` でテキストを選択できなくする。セルを自前で選ぶシート用）、`altText`、`measure`、`linkLabel`、`onSpan`、`className`。表のあるセルから別のセルへ届いた選択は、その間の矩形のセルの選択になる |
| `TEXT_LAYER_CSS` | テキスト層の CSS。run の span は z-index 1（他の run より上に置く span は 2）。選択をドラッグしてテキストのない所へ出ても選択が飛ばないように、層の末尾の要素がドラッグの間だけスクロール領域を覆う |
| `installCopyHandler(container)` | ページをまたいだ選択を、文書の空白と改行を戻してコピーする。表のセルの選択はタブ区切りと HTML の表にする |
| `selectionText` / `selectedRuns` / `joinRuns` | 選択範囲のテキスト |
| `selectionCells(selection?, root?)` | 選択が表のセルを選んでいれば、その矩形の `{text, html}`（タブ区切りと HTML の表）。テキストの選択なら `undefined` |
| `tableCells(content, table?, keep?)` / `cellClipboard(cells, range, options?)` | `TextContent` の表（`table` は表のノードの番号。`-1` で表の外のセル、つまりシートのセル）のセルとそのテキスト（`CellText`）/ セルの矩形（`CellRange`、0 始まりで両端を含む）を表計算ソフトが貼り付けられるタブ区切りと HTML の表にする。`options` は `trim`（列・行全体の選択を、テキストのある最後の行・列までにする）、`skipRow` / `skipCol`（非表示の行・列を除く） |
| `RUN_ATTR` / `linkHref` / `internalLink` | run の要素に付く属性名 / リンク先として安全な URL / 文書内リンク（`#page=N`、`#view=ID&page=N`）の解釈 |

### 同じスレッドで描く

Worker を使わない場合や、自前の Worker に組み込む場合の部品。

| 名前 | 内容 |
|---|---|
| `PageRenderer(doc, options?, fontSet?, resources?)` | `renderPage(ctx, page, {scale, roles?, background?})`、`renderContinuous`、`renderSheet`、`continuousLayout`、`sheetSize`、`preparePage`、`preparePageText`、`dispose` |
| `CanvasRenderer(res, options?)` | Object 1 つを 2D コンテキストに描く（`draw(ctx, obj, reset?, visible?)`。`visible` は Object の座標で見える矩形で、そこから十分に離れたテキストを飛ばす） |
| `ResourceCache(doc, fontSet?, {imageBudget?, decodeImage?})` | Part から作る `Path2D`、`ImageBitmap`、`FontFace` のキャッシュ。デコードした画像は `imageBudget`（既定 `DEFAULT_IMAGE_BUDGET` = 256 MiB）まで保持し、描画中のものは `hold()` / `release()`（`ImageHold`）で守る |
| `DocumentSearch(doc, measure)` | 検索とヒットの矩形（`search`、`locate`、`text`、`forget`） |
| `fontString` / `embeddedFamily` / `buildPath2D` / `cssColor` / `resetState` | 描画の小さな部品 |
