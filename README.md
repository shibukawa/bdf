# bdf

ブラウザの Canvas 2D にそのまま描画できる、プレビュー用の文書フォーマット（ドラフト）。

Office 系ファイルをサーバーで変換しておき、フロントエンドの Worker で描画する用途を想定しています。PDF が持つ機能のうちブラウザが標準 API で代替できるもの（フォントラスタライズ、画像デコード、圧縮）はブラウザに任せ、デコーダを最小にします。

- 固定サイズページ（スライド）、無限平面（シート）、ページ分割かつ連続表示可能な文書（ワープロ）の 3 モデル
- 命令セットは `CanvasRenderingContext2D` に 1:1 対応
- `DecompressionStream` で展開、Worker + `OffscreenCanvas` で描画
- 内容アドレスの Part によりマスターや繰り返し部品を自動共有
- テキスト索引 Part と Worker 内の全文検索（行またぎ、NFKC・かな正規化、ヒット矩形）
- 1 ファイル形式と分割ファイル形式を相互変換可能

PDF からの変換（`pdf2bdf`）は、埋め込みフォントをブラウザが読める形に組み直し、フォーム XObject を共有オブジェクトに、テキストを検索可能な run に変換します。詳細は design.md の §3.1 を参照してください。

ドキュメント:

- [docs/spec.md](docs/spec.md) — フォーマット仕様ドラフト
- [docs/design.md](docs/design.md) — 設計判断の理由と実装方針

## 構成

| 場所 | 内容 |
|---|---|
| `*.go`, `cmd/bdf` | Go のエンコーダ・デコーダ・コンテナ I/O と CLI |
| `fixture/` | サンプル文書の生成（埋め込みフォント付き） |
| `imgconv/` | 画像の格納方針（そのまま / WebP に変換）。純 Go の libwebp を同梱 |
| `pdf2bdf/`, `cmd/pdf2bdf` | PDF → BDF 変換器と CLI |
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

# PDF → BDF
go run ./cmd/pdf2bdf in.pdf out.bdf     # 1 ファイル形式
go run ./cmd/pdf2bdf in.pdf out/        # 分割形式
go run ./cmd/pdf2bdf -pages 1-3 -kind flow in.pdf out.bdf
go run ./cmd/pdf2bdf -images keep in.pdf out.bdf   # 画像を変換しない
go build -tags bdf_noconv ./...                    # コーデックを含めないビルド（ブラウザ向け）

# TypeScript: ビルドとテスト
npm ci
npm test                             # デコーダのテスト（Node）
npm run test:golden                  # Chromium で描画して golden 画像と比較
npm run test:golden:update           # golden 画像を更新
npm run fixtures                     # fixtures/ を再生成（Go が必要）
node test/render.mjs out.bdf pngdir/  # 任意の .bdf を Chromium で PNG に描画

# デモビューア
npm run demo                         # http://127.0.0.1:8765/examples/viewer/.out/
```

Go は 1.26 以上が必要です。golden テストは `playwright-core` を使います。Chromium は `CHROMIUM_PATH` で指定するか、`npx playwright-core install chromium` で入れてください。

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
```
