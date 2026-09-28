# bdf

[English](README.md) | 日本語

**bdf**（Browser-specific Document Format）は、ブラウザの Canvas 2D にそのまま描画できる、プレビュー用の文書フォーマット（ドラフト）です。

PDF、Office のファイル（Word、PowerPoint、Excel、CSV、Parquet、Visio）、draw.io の図、CAD の図面（DXF、Jw_cad、SXF、CGM、HP-GL/2）、プリント基板（Gerber、Excellon、KiCad）、Illustrator、Photoshop、HTML、Markdown、EPUB、音楽（MML、MIDI、MusicXML。ビューアで演奏可能）、フォントファイルを bdf に変換し、同じレンダラで Web Worker 内に描画します。オフィススイートもサーバー上のヘッドレスブラウザも不要です。変換器は素の Go で、WebAssembly にすればブラウザの中でも変換できます。

**[試してみる](https://shibukawa.github.io/bdf/)** — デモサイトにファイルをドロップしてください。アップロードせずにブラウザの中で変換して描画します。**[ドキュメントを読む](https://shibukawa.github.io/bdf/docs/)** — なぜ bdf があるのか、対応するすべての形式、サーバー構成のサンプル 3 つ、ビルドの方法まで。

## クイックスタート

```sh
go run ./cmd/bdf generate report.pptx report.bdf   # 変換
go run ./cmd/bdf ls report.bdf                     # 中身を見る

npm ci && npm run demo                             # ブラウザで開く
# http://127.0.0.1:8765/examples/viewer/.out/?src=/report.bdf
```

ソースからビルドするには Go 1.27 以上と Node.js が要ります。続きは[はじめかた](https://shibukawa.github.io/bdf/docs/getting-started.ja.html)を参照してください（Go のコードから直接変換する方法も含みます）。

## 構成

| 場所 | 内容 |
|---|---|
| `*.go`、`converter/`、`raster/`、`internal/` | Go のエンコーダ・デコーダと、入力形式ごとの変換器パッケージ |
| `packages/core`、`packages/render` | `@bdf/core`・`@bdf/render`: TypeScript のデコーダと Canvas レンダラ |
| `examples/`、`site/` | デモビューア、サーバー構成のサンプル 3 つ、GitHub Pages で公開しているサイト |
| `docs/` | このドキュメント |

全体の内訳は[フォルダ構成](https://shibukawa.github.io/bdf/docs/architecture/layout.ja.html)を、`go test`・golden テスト・フィクスチャの再生成は[ビルドとテスト](https://shibukawa.github.io/bdf/docs/architecture/development.ja.html)を参照してください。

## ライセンス

[MIT](LICENSE)。
