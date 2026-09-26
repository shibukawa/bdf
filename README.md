# bdf

ブラウザの Canvas 2D にそのまま描画できる、プレビュー用の文書フォーマット（ドラフト）。

Office 系ファイルをサーバーで変換しておき、フロントエンドの Worker で描画する用途を想定しています。PDF が持つ機能のうちブラウザが標準 API で代替できるもの（フォントラスタライズ、画像デコード、圧縮）はブラウザに任せ、デコーダを最小にします。

- 固定サイズページ（スライド）、無限平面（シート）、ページ分割かつ連続表示可能な文書（ワープロ）の 3 モデル
- 命令セットは `CanvasRenderingContext2D` に 1:1 対応
- `DecompressionStream` で展開、Worker + `OffscreenCanvas` で描画
- 内容アドレスの Part によりマスターや繰り返し部品を自動共有
- テキスト索引 Part と Worker 内の全文検索（行またぎ、NFKC・かな正規化、ヒット矩形）
- 透明 DOM のテキスト選択層とコピー（空白・改行は MARK 境界から復元、ページまたぎ、連続モード対応）
- 1 ファイル形式と分割ファイル形式を相互変換可能

変換は `bdf generate` サブコマンドで行い、入力の形式（PDF / PowerPoint）は中身から判別します。

- **PDF**（`converter/pdf`）: 埋め込みフォント（TrueType、CFF、OpenType、Type1）を使うグリフだけの WOFF2 に組み直し（OS/2 の埋め込み許諾 `fsType` を確認し、著作権表示は引き継ぐ）、フォーム XObject を共有オブジェクトに、テキストを検索可能な run に変換し、各ページ先頭の共通部分（マスター）を共有 Object に切り出します。詳細は design.md の §3.1。
- **PowerPoint .pptx**（`converter/pptx`）: DrawingML を直接描画します。スライドマスターとレイアウトの図形はスライド間で共有されるレイヤー Object になり、プリセット図形は ECMA-376 の図形定義式から、テキストは変換側で折り返し（和文の禁則・縦書き・箇条書き・段落書式）、表・グラフ・SmartArt・EMF/WMF の図も描きます。レイアウトに使ったフォントはサブセットの WOFF2 にして埋め込むので、閲覧環境のフォントに依存しません。詳細は design.md の §3.4。

ドキュメント:

- [docs/spec.md](docs/spec.md) — フォーマット仕様ドラフト
- [docs/design.md](docs/design.md) — 設計判断の理由と実装方針

## 構成

| 場所 | 内容 |
|---|---|
| `*.go`, `cmd/bdf` | Go のエンコーダ・デコーダ・コンテナ I/O と CLI（`bdf generate` で変換） |
| `fixture/` | サンプル文書の生成（埋め込みフォント付き） |
| `imgconv/` | 画像の格納方針（そのまま / WebP に変換）。純 Go の libwebp を同梱 |
| `converter/` | 変換器の共通部分（入力形式の判別、ページ指定） |
| `converter/pdf` | PDF → BDF 変換器 |
| `converter/pptx` | PowerPoint (.pptx) → BDF 変換器 |
| `converter/internal/` | フォントの探索・計測・サブセット化（`fontdb`）、TrueType/OpenType の読み書き（`sfnt`） |
| `woff2/` | TrueType/OpenType → WOFF2（glyf 変換と Brotli） |
| `packages/core` | `@bdf/core`: TypeScript のデコーダ、コンテナ読み込み、テキスト抽出 |
| `packages/render` | `@bdf/render`: Canvas レンダラ、ページ/連続/シート描画、Worker |
| `examples/viewer` | デモビューア |
| `fixtures/` | 生成済みサンプルと golden 画像 |

## 使い方

```sh
# Go: テストと CLI
go test ./...
go run ./cmd/bdf demo out.bdf        # サンプル文書を生成
go run ./cmd/bdf ls out.bdf          # Part 一覧
go run ./cmd/bdf disasm out.bdf <hash>
go run ./cmd/bdf split out.bdf out/  # 分割形式へ

# PDF / PowerPoint → BDF（形式は中身から判別。-format pdf|pptx で指定も可）
go run ./cmd/bdf generate in.pdf out.bdf      # 1 ファイル形式
go run ./cmd/bdf generate in.pptx out/        # 分割形式
go run ./cmd/bdf generate -pages 1-3 in.pptx out.bdf   # ページ（スライド）を選ぶ
go run ./cmd/bdf generate -images keep in.pdf out.bdf  # 画像を変換しない
go run ./cmd/bdf generate -kind flow in.pdf out.bdf    # PDF: flow View にする
go run ./cmd/bdf generate -no-share in.pdf out.bdf     # PDF: ページ共通の先頭部分（マスター）を共有 Object にしない
go run ./cmd/bdf generate -font-dir fonts/ in.pptx out.bdf         # PowerPoint: フォントを探すディレクトリを追加
go run ./cmd/bdf generate -fonts system in.pptx out.bdf            # PowerPoint: フォントを埋め込まず名前で参照
go run ./cmd/bdf generate -hidden in.pptx out.bdf                  # PowerPoint: 非表示スライドも含める
go run ./cmd/bdf generate -no-woff2 in.pdf out.bdf     # フォントを WOFF2 にせず TTF/OTF のまま格納
go run ./cmd/bdf generate -ignore-fstype in.pdf out.bdf # fsType が埋め込みやサブセット化を禁じるフォントも埋め込む（権利がある場合のみ）
go build -tags bdf_noconv ./...                    # コーデック（WebP、WOFF2 の Brotli）を含めないビルド（ブラウザ向け）
GOEXPERIMENT=simd go build ./...                   # Go 1.27 amd64/arm64: SIMD 版コーデック（amd64 は AVX2 必須）

# TypeScript: ビルドとテスト
npm ci
npm test                             # デコーダのテスト（Node）
npm run test:golden                  # Chromium で描画して golden 画像と比較
npm run test:golden:update           # golden 画像を更新
npm run fixtures                     # fixtures/ を再生成（Go が必要）
npm run test:pptx:gen                # PowerPoint のテスト用デッキを再生成（python-pptx が必要）
node test/render.mjs out.bdf pngdir/  # 任意の .bdf を Chromium で PNG に描画

# デモビューア
npm run demo                         # http://127.0.0.1:8765/examples/viewer/.out/
```

Go は 1.27 以上が必要です。golden テストは `playwright-core`（固定バージョン。golden 画像はその Chromium ビルドの headless shell で描いたもの）を使います。`npx playwright-core install chromium` で入れるか、同じビルドの headless shell を `CHROMIUM_PATH` で指定してください。

## エンコーダ API の雰囲気（Go）

```go
d := bdf.NewDocument()
font := d.AddFont(woff2Bytes)

master := bdf.NewObject()
bg := master.AddPaint(bdf.LinearGradient(0, 0, 0, 540,
    bdf.Stop{Offset: 0, Color: bdf.RGB(0xf7, 0xf9, 0xfc)},
    bdf.Stop{Offset: 1, Color: bdf.RGB(0xdc, 0xe6, 0xf2)}))
master.FillPaint(bg).FillRect(0, 0, 960, 540)
masterHash, _ := d.AddObject(master)

body := bdf.NewObject()
f := body.AddFont(bdf.EmbeddedFont(font, 700, bdf.StyleNormal))
body.Font(f, 36).FillColor(bdf.RGB(0x1f, 0x3a, 0x5f)).FillText("Hello", 48, 80, advance)
bodyHash, _ := d.AddObject(body)

v := d.NewView("slides", bdf.ViewFixed, "Slides")
v.AddPage(960, 540,
    bdf.Layer{Role: bdf.RoleMaster, Obj: masterHash},
    bdf.Layer{Role: bdf.RoleBody, Obj: bodyHash})
d.WriteSingle(w)
```

## ビューア API の雰囲気（TypeScript）

```ts
import { BdfWorkerClient } from "@bdf/render";
const client = new BdfWorkerClient(new Worker(workerUrl, { type: "module" }));
const manifest = await client.open({ kind: "single", url: "/doc.bdf", range: true });
const bitmap = await client.page("slides", 0, devicePixelRatio); // ImageBitmap
const runs = await client.text("slides", 0);                      // 選択・コピー用テキスト
const hits = await client.search("doc", "list of objects");         // 全文検索（正規化つき）
const rects = await client.locate("doc", hits);                    // ハイライト矩形（ページ座標）

// 選択層: ページの上に透明な span を置き、コピーは run の区切りから組み立てる
import { buildTextLayer, installCopyHandler, TEXT_LAYER_CSS } from "@bdf/render";
pageElement.append(canvas, buildTextLayer(runs, zoom));          // TEXT_LAYER_CSS を読み込んでおく
installCopyHandler(stage);                                         // copy で選択範囲のテキストを整形
```
