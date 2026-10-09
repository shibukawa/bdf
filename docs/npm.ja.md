# npm パッケージで使う

ブラウザ向けの BDF は、サーバーで変換済みの BDF を表示するだけなら `@bdfkit/viewer`、入力ファイルをブラウザ内で BDF に変換するならプリセットを 1 つ追加します。これらはリポジトリ内の npm workspace パッケージで、下の `npm install` は公開後の利用例です。公開前に `npm run build && npm run build:wasm` で配布物を生成します。wasm は TinyGo で生成し、各変換パッケージの tarball に含めます。インストール先で TinyGo や Go を動かす必要はありません。

用途別の手順は[フロントエンドに組み込む](integrating-to-frontend.ja.md)、[サーバーに組み込む](integrate-to-server.ja.md)、[Wasm ランタイムをビルドする](build-wasm-runtime.ja.md)を参照してください。

| 構成 | インストールするパッケージ | 入力形式 |
|---|---|---|
| レンダラーのみ | `@bdfkit/viewer` | サーバー変換済みの BDF |
| PDF + EPUB | `@bdfkit/viewer`、`@bdfkit/convert-pdf-epub` | PDF、EPUB |
| Office + PDF + EPUB | `@bdfkit/viewer`、`@bdfkit/convert-office` | Word、Excel、PowerPoint、Visio、InDesign、CSV、TSV、Parquet、PDF、EPUB |
| All | `@bdfkit/viewer`、`@bdfkit/convert-all` | このリポジトリのブラウザ用変換器すべて |

レンダラーを外す場合は上表の `@bdfkit/viewer` をインストールせず、変換パッケージだけを使います。音楽機能も通常のビューアーから分離してあり、ピアノロールと再生が必要な場合だけ `@bdfkit/viewer/music` を import します。All プリセットから MML・MIDI・MusicXML の変換器も外すには、[カスタム Wasm ビルド](build-wasm-runtime.ja.md#tinygo-好きな形式を選ぶ)で `--formats all --without-music` を指定します。

`@bdfkit/viewer` は `@bdfkit/render` と `@bdfkit/core` に依存します。変換パッケージは `@bdfkit/convert` に依存します。アプリのバンドラーはパッケージ内の `new URL(..., import.meta.url)` が参照する `worker.js` と `converter.wasm` を配信してください。TinyGo の実行コードは変換 Worker に含まれます。CSP で Worker と WebAssembly を許可する必要もあります。

```sh
npm install @bdfkit/viewer @bdfkit/convert-office
```

## レンダラーのみ（サーバー変換）

サーバー側で `bdf generate` などを使って BDF を作る場合、ブラウザに wasm を送る必要はありません。

```sh
npm install @bdfkit/viewer
```

```html
<div id="preview" style="width:100%;height:70vh"></div>
```

```ts
import { Viewer } from "@bdfkit/viewer";

const viewer = new Viewer(document.querySelector("#preview")!, { layout: "single" });
await viewer.open({ kind: "single", url: "/documents/report.bdf" });
// サーバーが分割 BDF を配る場合:
// await viewer.open({ kind: "split", base: "/documents/report/" });
viewer.setLayout("spread");
// 画面の破棄時に viewer.destroy()
```

ページ配置を自作する場合は低レベルの `@bdfkit/render` と `@bdfkit/core` を直接使えます。`@bdfkit/render` は描画 Worker を、`@bdfkit/viewer` はスクロール面と表示モードを提供します。

## ブラウザ内変換を追加する

```ts
import { Viewer } from "@bdfkit/viewer";
import { createOfficeConverter } from "@bdfkit/convert-office";

const viewer = new Viewer(document.querySelector("#preview")!, {
  mode: "embedded",       // または "lightbox"
  layout: "single",       // "spread"、"continuous" も選択可能
});
const converter = createOfficeConverter();
const file = (document.querySelector("input[type=file]") as HTMLInputElement).files![0];
await viewer.openFile(converter, await file.arrayBuffer(), { name: file.name, fonts: "/fonts/" });

// ホストアプリ側のボタンなどから呼ぶ。
viewer.setLayout("spread");
viewer.setView(viewer.document?.views[0]?.id);
viewer.scrollToPage(3);             // 0 始まり
// mode: "lightbox" の場合は viewer.show() / viewer.hide()

// 画面を破棄するとき
viewer.destroy();
converter.terminate();
```

埋め込み表示では `#preview` に幅と高さを指定します。ライトボックスは `show()` まで非表示です。パッケージは文書面、ページ配置、スクロール、シートのビューポート描画を担当します。ページ単位で変換できる入力は、輪郭を先に表示し、変換済みページを順に反映します。ファイル選択、ツールバー、タブ、ズームボタン、ライトボックスの閉じるボタンはアプリ側で作り、`setView`、`setLayout`、`setZoom`、`show`、`hide` を呼びます。`documentchange`、`viewchange`、`conversionprogress`、`conversiondone`、`error` イベントを監視できます。

## 配信サイズ

下表は 2026 年 9 月 30 日、TinyGo 0.42.0 で `npm run build && npm run build:wasm && npm run size:npm` を実行した結果です。各構成で依存パッケージを一度だけ数え、`dist` の実行用 `.js` と `.wasm` の合計を示します。gzip は各ファイルを個別に gzip レベル 9 で圧縮した合計です。型定義、source map、フォント、入力文書、HTML/CSS は含みません。レンダラーなしの構成は `@bdfkit/viewer` をインストールしません。実際のバンドラーが未使用コードを除けば JS は小さくなります。

<!-- npm-size-table -->
| 構成 | 無圧縮 | gzip | 対象ファイル |
|---|---:|---:|---:|
| core + render（低レベル） | 356.7 KiB | 108.5 KiB | 29 |
| レンダラーのみ | 371.4 KiB | 112.7 KiB | 31 |
| PDF + EPUB（レンダラーなし） | 13.79 MiB | 6.09 MiB | 4 |
| PDF + EPUB | 14.15 MiB | 6.20 MiB | 35 |
| Office + PDF + EPUB（レンダラーなし） | 14.74 MiB | 6.44 MiB | 4 |
| Office + PDF + EPUB | 15.10 MiB | 6.55 MiB | 35 |
| All（レンダラーなし） | 20.11 MiB | 8.81 MiB | 4 |
| All | 20.48 MiB | 8.92 MiB | 35 |
<!-- /npm-size-table -->

変換用 Wasm は TinyGo 0.42.0 でビルドしています。`tools/build-npm-converters.mjs` は `tinygo build -target wasm -no-debug` を呼び出し、TinyGo の既定値である `-opt=z` のサイズ最適化を使います。これは小さな言語ランタイムだけではなく、選択した形式のパーサー、レイアウト、ラスタライザー、検索テーブルなどを静的リンクした変換モジュールです。そのため `All` は全変換器を含み、フロントエンドの JavaScript の tree shaking では Wasm 内のコードを削れません。

レンダラーのみの構成なら変換用 Wasm は不要です。ブラウザで変換する場合は小さいプリセットを選ぶか、`--formats` で独自の Wasm を作ります。独自ビルドでは `--without-music` を指定して音楽系変換器を外せます。転送時は `.wasm` を Brotli などの HTTP 圧縮で配信すると、gzip より小さくできます。上表は環境をまたいで再現しやすいよう gzip レベル 9 で計測しています。

上表のファイル数には未使用のエントリも含みます。アプリ側の JS を minify してバンドルした場合、音楽機能を import しない構成は 11.7 KiB（gzip 4.4 KiB）、`@bdfkit/viewer/music` からピアノロールと `MusicPlayer` を import する構成は 29.5 KiB（gzip 10.9 KiB）でした。この比較には描画 Worker と Wasm は含みません。

## 好きな形式だけを入れた wasm

プリセットにない組み合わせはリポジトリから生成できます。形式名は `converter/` のパッケージ名を使い、`tsv` は `csv` と指定しても同じ変換器にまとめます。

```sh
npm ci
node tools/build-npm-converters.mjs --formats pdf,epub,docx,csv --out ./public/my-converter.wasm
```

生成した wasm を `@bdfkit/convert` の Worker で動かします。

```ts
import { createConverter } from "@bdfkit/convert";
const converter = createConverter(new URL("/my-converter.wasm", location.href));
```

`await converter.ready` で実際に組み込まれた形式の一覧を確認できます。フォントが必要な入力では、`convert` または `openFile` のオプションに `fonts` としてフォントディレクトリの URL を渡します。ディレクトリには `index.json` が必要です。生成方法は `examples/common/build.mjs` の `copyFonts` を参照してください。

## 配布物の確認

```sh
npm run build
npm run build:wasm
npm pack --dry-run --workspace @bdfkit/convert-office
```

`@bdfkit/convert` と使用するプリセット、`@bdfkit/viewer`、`@bdfkit/render`、`@bdfkit/core` を同じバージョンで公開します。

## GitHub Actions から公開する

`.github/workflows/publish.yml` は `v*` タグを push すると 7 つのパッケージを公開します。タグと全パッケージのバージョンが一致することを確認し、TypeScript と Wasm のビルド・テスト、tarball の確認を行ってから依存順に公開します。

初回リリースは手元から公開し、その後 npm の Trusted Publisher を 7 パッケージに登録します。公開するマシンで npm にログインし、依存順に実行します。

```sh
npm login
npm ci
npm test
npm run build:wasm

for package in \
  @bdfkit/core \
  @bdfkit/render \
  @bdfkit/convert \
  @bdfkit/viewer \
  @bdfkit/convert-pdf-epub \
  @bdfkit/convert-office \
  @bdfkit/convert-all; do
  npm publish --workspace="$package" --access=public
done
```

初回公開後、GitHub の Settings → Environments で `release` Environment を作成します。必要なレビュアーを設定し、デプロイ対象のブランチまたはタグをリリース用（例: `v*`）に限定してください。npm の各パッケージの Trusted Publisher には、このリポジトリ、`publish.yml`、Environment name `release` を登録し、直接の `npm publish` を許可します。ワークフローには GitHub OIDC 用の `id-token: write` を設定済みです。以後は `release` の承認後だけ短期 OIDC credential で公開され、`NPM_TOKEN` secret は不要です。

```sh
git tag v0.1.0
git push origin main v0.1.0
```
