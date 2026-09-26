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
| `FileName` | 入力のファイル名（内容で判定できない CSV の判定、CSV のシート名） |
| `Warn func(string)` | 警告を受け取る（無ければ `Result.Warnings` に集める） |

### 形式ごとのパッケージ

直接呼ぶと、`converter.Options` に無い形式固有の設定も使える。どのパッケージも `Convert` と `ConvertFile` を持ち、`Options` には `Title` と `Warn`、テキストをレイアウトする形式ならフォントまわり（`FontFS`、`FontDirs` など）と `NoTextIndex`、画像を持つ形式なら `Images` がある。ページの選択は `converter.Pages` 型である。

| パッケージ | 形式 | 固有の入口と設定 |
|---|---|---|
| `converter/pdf` | PDF | `Convert(rs io.ReadSeeker, opts)`、`NewStream(rs, opts)`（ページ単位の変換）。`Kind`（`fixed` か `flow`）、`NoSharePrefix`（共通プレフィックスを共有しない）、`NoAnnotations`、`Password` |
| `converter/pptx` | PowerPoint | `Slides`（`converter.Pages`）、`Hidden`（非表示スライドも含める） |
| `converter/xlsx` | Excel | `Sheets`（`converter.Pages`）、`Hidden`。`ConvertGrid(*Grid, opts)` は書式のない値の表（`Grid`）を Excel の新しいブックのように見せる |
| `converter/csv` | CSV・TSV | `Charset`、`Delimiter`、`Quote`、`Header`（`HeaderAuto` / `HeaderYes` / `HeaderNo`）、`TableStyle`、`Name`（シート名）。指定しなかったものは推測する |
| `converter/docx` | Word | `Pages`、`Views`（`ViewsBoth` / `ViewsPages` / `ViewsScroll`） |
| `converter/visio` | Visio（.vsdx、.vdx） | `Pages` |
| `converter/drawio` | draw.io（.drawio、図を埋め込んだ .drawio.svg・.drawio.png） | `Convert(data []byte, opts)`（入力をバイト列で渡す）。`Pages`、`Border`（図の周りの余白） |
| `converter/dxf` | AutoCAD DXF | `Pages`、`Views`（`all` / `model` / `layouts`）、`Light`（モデル空間を白い紙に描く）。`Detect(head)` |
| `converter/tiff` | TIFF | `Pages`、`DPI`（解像度の無いページに仮定する値）、`Images`（解像度の上限もここ） |
| `converter/emf` | EMF・WMF | — |
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
| `(*View).AddPage(w, h, layers...) *Page` | ページを足す（fixed・flow・scroll）。シートは `Tiles`、`Cols`、`Rows`、`Freeze` などの項目を直接設定する |
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
| `DecodePathCollection` / `EncodePathCollection` | パス集合の Part |
| `SharePrefixes(objs, minBytes)` | 複数の Object に共通する先頭部分を共有 Object に切り出す |
| `HashOf` / `ParseHash` | Part のハッシュ |

### その他のパッケージ

| パッケージ | 内容 |
|---|---|
| `imgconv` | 画像の格納方法。`Options{Mode: imgconv.Convert, Quality: 80}` で WebP を試して小さい方を残す（`Keep` はそのまま）。`MaxDPI`（既定 `DefaultMaxDPI` = 192）と `MaxPixels` はページ上の大きさが分かるラスター入力の解像度の上限。`Optimize(data, opts)`、`EncodePixels`、`Resize`、`Available()`（`bdf_noconv` ビルドでは false） |
| `woff2` | TrueType・OpenType を WOFF2 にする。`Encode(font)`、`Available()`、`IsWOFF(data)` |

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

## ブラウザ内変換（cmd/bdfwasm）

`cmd/bdfwasm` は変換器を `GOOS=js GOARCH=wasm` でビルドしたもので、Go の `wasm_exec.js` とともに Worker で動かす（デモサイトは `examples/viewer/convert-worker.ts`）。起動すると `globalThis.bdfConverter` を設定する。`-tags pdfonly` と `officeonly` で PDF 用と Office 系用に分けてビルドできる。

| 名前 | 内容 |
|---|---|
| `formats` | `[{name, description, extensions}]` |
| `convert(data: Uint8Array, options?)` | 丸ごと変換する。`{bdf, format, summary, warnings, protected}` を返す Promise |
| `open(data, options?)` | ページ単位の変換を始める。`{bdf, format, pages, warnings, stream?}` を返す。`pages > 0` なら `bdf` は outline で、`stream` が残りを変換する。対応していない形式は丸ごと変換した結果（`pages: 0`、`summary` つき） |
| `stream.page(i)` | ページ i を変換し、`{bdf, warnings}`（そのページだけのページ文書と新しい警告）を返す |
| `stream.finish()` | 残りを変換し、`convert` と同じ形の完成した文書を返す |
| `stream.close()` | ストリームを手放す（`finish` のあとも呼ぶ） |

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
| `TextSearch(runs)` / `decodeTextIndex(bytes)` | テキスト索引から検索する（NFKC、大文字小文字、小書きのかなを同一視し、行をまたいで一致する）。`search(query, {limit, caseSensitive, context})` は `SearchHit[]` |
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
| `buildTextLayer(content, scale, options?)` | 透明なテキスト層の要素を作る。選択、コピー、検索のハイライト位置合わせ、読み上げ（見出し、リスト、表、代替テキスト、リンク、言語）に使う。`options` は `lang`、`sheet`（シートを表として並べる）、`altText`、`measure`、`linkLabel`、`onSpan`、`className` |
| `TEXT_LAYER_CSS` | テキスト層の CSS |
| `installCopyHandler(container)` | ページをまたいだ選択を、文書の空白と改行を戻してコピーする |
| `selectionText` / `selectedRuns` / `joinRuns` | 選択範囲のテキスト |
| `RUN_ATTR` / `linkHref` / `internalLink` | run の要素に付く属性名 / リンク先として安全な URL / 文書内リンク（`#page=N`、`#view=ID&page=N`）の解釈 |

### 同じスレッドで描く

Worker を使わない場合や、自前の Worker に組み込む場合の部品。

| 名前 | 内容 |
|---|---|
| `PageRenderer(doc, options?, fontSet?, resources?)` | `renderPage(ctx, page, {scale, roles?, background?})`、`renderContinuous`、`renderSheet`、`continuousLayout`、`sheetSize`、`preparePage`、`preparePageText`、`dispose` |
| `CanvasRenderer(res, options?)` | Object 1 つを 2D コンテキストに描く（`draw(ctx, obj)`） |
| `ResourceCache(doc, fontSet?, {imageBudget?, decodeImage?})` | Part から作る `Path2D`、`ImageBitmap`、`FontFace` のキャッシュ。デコードした画像は `imageBudget`（既定 `DEFAULT_IMAGE_BUDGET` = 256 MiB）まで保持し、描画中のものは `hold()` / `release()`（`ImageHold`）で守る |
| `DocumentSearch(doc, measure)` | 検索とヒットの矩形（`search`、`locate`、`text`、`forget`） |
| `fontString` / `embeddedFamily` / `buildPath2D` / `cssColor` / `resetState` | 描画の小さな部品 |
