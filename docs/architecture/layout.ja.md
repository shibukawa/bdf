# フォルダ構成

## コア

| 場所 | 内容 |
|---|---|
| `*.go`, `cmd/bdf` | Go のエンコーダ・デコーダ・コンテナ I/O と CLI |
| `converter/` | 入力形式の登録と共通のオプション、形式の判別、ページ指定 |
| `converter/internal/` | Office 系・CAD の変換器が共有するもの（下記） |
| `raster/imagebdf/` | ページと単体の Object を純 Go で画像に描く（ビューアと同じ描き方）。ソフトウェアのラスタライザと、SVG の画像を描く SVG レンダラ |
| `raster/ebitenginebdf/` | パスだけの Object（数式など）を Ebitengine で描く。Ebitengine に依存するので独自の `go.mod` を持つ別モジュール。`raster/internal/shapes` は Object を塗りと線の図形の並びに展開する共通部分 |
| `formula/` | LaTeX・MathML の数式を組んでパスだけの Object にする公開パッケージ。数式フォント（STIX Two Math）を同梱 |
| `thumbnail/` | 文書のサムネイル。文書の種類によるレイアウトと、PNG・JPEG・WebP の書き出し |
| `internal/` | フォントの探索・計測・サブセット化（`fontdb`）、TrueType/OpenType の読み書きとグリフの輪郭（`sfnt`）、CFF の読み取りとサブセット化（`cff`）、OpenType のレイアウトテーブル GSUB・GPOS・GDEF の読み取り（`otlayout`）。数式の木・読み手（LaTeX、MathML、Office Math）・組版（`mathlayout`）と、Office の XML の木（`xmltree`）。変換器と `imagebdf` が共有する |
| `woff2/` | TrueType/OpenType ↔ WOFF2（glyf 変換と Brotli） |
| `imgconv/` | 画像の格納方針（そのまま／WebP に変換）。純 Go の libwebp を同梱 |
| `fixture/` | テストと `bdf demo` が使うサンプル文書の生成（埋め込みフォント付き） |

## 変換器（入力形式ごとに 1 パッケージ）

| 場所 | 変換する形式 |
|---|---|
| `converter/pdf` | PDF |
| `converter/ai` | Illustrator（.ai）。PDF 変換器の上に |
| `converter/psd` | Photoshop（.psd, .psb） |
| `converter/pptx` | PowerPoint（.pptx） |
| `converter/xlsx` | Excel（.xlsx） |
| `converter/csv` | CSV・TSV（描画は `converter/xlsx`） |
| `converter/parquet` | Apache Parquet（描画は `converter/xlsx`） |
| `converter/docx` | Word（.docx） |
| `converter/visio` | Visio（.vsdx） |
| `converter/drawio` | draw.io（.drawio / .drawio.svg / .drawio.png） |
| `converter/dxf` | AutoCAD DXF |
| `converter/jww` | Jw_cad（.jww） |
| `converter/sxf` | SXF（.p21, .p2z, .sfc） |
| `converter/cgm` | CGM（.cgm, .cgz） |
| `converter/hpgl` | HP-GL/2 のプロットファイル（.plt。HP-GL、PJL、PCL、HP RTL） |
| `converter/gerber` | Gerber・Excellon（と基板のファイルをまとめた ZIP） |
| `converter/kicad` | KiCad の回路図・基板・プロジェクト（とプロジェクトをまとめた ZIP） |
| `converter/emf` | Windows メタファイル（.emf, .wmf） |
| `converter/tiff` | TIFF（.tif, .tiff） |
| `converter/html` | HTML（.html, .xhtml, .mhtml）。リーダー表示 |
| `converter/markdown` | Markdown（HTML を経由） |
| `converter/epub` | EPUB（リフロー型の本はリーダー表示、固定レイアウトの本は画像のページ） |
| `converter/mml` | MML（.mml。汎用、マビノギ、PPMCK の方言）。楽譜にする |
| `converter/midi` | Standard MIDI File（.mid, .kar, .rmi）。楽譜にする |
| `converter/musicxml` | MusicXML（.musicxml, .mxl） |
| `converter/font` | フォントファイル（.ttf, .otf, .ttc, .woff, .woff2）。文字・グリフ・OpenType フィーチャーのプレビュー |
| `converter/image` | PNG・JPEG・GIF・WebP・AVIF・BMP・ICO・SVG（そのまま格納し、メタデータを読む） |
| `converter/audio` | MP3・M4A・FLAC・Ogg・WAV・AIFF（カバーアートとタグのカード。音声は格納しない） |
| `converter/all` | すべての入力形式を登録する（副作用のために import する） |

`converter/internal/` には複数の変換器が共有するものが入っています: OOXML のパッケージと XML（`ooxml`）、DrawingML の図形・テキスト・表・グラフ（`ooxml/drawingml`）、テキストレイアウト用のフォント選択・計測・埋め込み（`fontset`。draw.io も使う）、組み立て中の Object（`canvas`。draw.io も使う）、EMF/WMF の再生（`metafile`）、ISO base media file format の箱（`isobmff`。AVIF の画像と M4A の音声）、CAD 図面のページへの描画（`cad`）、行分割の規則（`linebreak`）、複合ファイル（`cfb`）とパスワード付き Office 文書の復号（`offcrypto`）。PDF 用の Adobe の定義済み CJK CMap（`cjkcmap`）と JPEG 2000・JBIG2 のデコーダ（`jpx`、`jbig2`）。TIFF の読み取りと CCITT の FAX 符号のデコーダ（`tiff`）。Word・HTML・Markdown・EPUB の組版エンジン（`wordproc`: 段落・表・ページ・scroll View）と、それらが共有する HTML・XHTML の読み込み（`webdoc`）。XMP メタデータの Dublin Core（`xmp`）。数式エンジン（`internal/mathlayout`）を文書のフォントで動かし、埋め込みフォントのテキストとして描く部分（`equation`）。MML・MIDI・MusicXML の楽譜の組版と演奏する音楽（`music`。Bravura の SMuFL の字形は `music/smufl`）。

## TypeScript のパッケージ

| 場所 | 内容 |
|---|---|
| `packages/core` | `@bdfkit/core` — デコーダ、コンテナ読み込み、テキスト抽出（依存なし） |
| `packages/render` | `@bdfkit/render` — Canvas レンダラ、ページ/連続/シート描画（scroll View は連続描画）、Worker とそのクライアント（SVG の画像は Worker がデコードできないのでメインスレッドが描く） |
| `cmd/bdfwasm` | ブラウザ内変換用に wasm にした変換器（PDF 用、Office 系用、HTML・Markdown・EPUB 用、画像用、変換器を持たないサムネイル・検索テキスト用のモジュール） |

## デモサイトとそのページ

| 場所 | 内容 |
|---|---|
| `examples/viewer` | フル機能のデモビューア（Worker による描画、読み上げ用のテキスト層、ページめくり） |
| `examples/common` | ビューアとサイトの残りが共有するもの（変換 Worker のクライアント、サイトの esbuild 補助） |
| `examples/miniviewer` | 埋め込み用の小さなビューア（下記のサンプルプロジェクトが使う） |
| `site/` | GitHub Pages で公開するサイト本体: トップページ、サムネイル・検索テキストのページ、ドキュメントの描画（`site/docs.mjs` がこのページを `docs/*.md` から作る） |
| `docs/` | このドキュメント（Markdown） |

## サンプルプロジェクト、CLI、テスト

| 場所 | 内容 |
|---|---|
| `examples/light-server`、`examples/preview-server`、`examples/search`、`examples/secure-reader` | BDF をシステムに組み込む方法ごとの小さな Go サーバー 4 つ。[構成のサンプル](../examples/index.ja.md)を参照 |
| `cmd/bdf` | `bdf` コマンド: `generate`、`thumbnail`、`text`、`render`、`ls`、`manifest`、`disasm`、`extract`、`split`、`join`、`encrypt`、`decrypt`、`demo`。[BDF コマンド](cli.ja.md)を参照 |
| `tools/` | ソースが変わったときに手で実行する生成ツール群: Unicode のデータ、SMuFL の字形パス、draw.io のステンシルと AWS アイコンの対応表、テーブルスタイル、cmap、純 Go の WebP コーデックの生成テーブル |
| `testdata/` | 生成済みサンプルと golden 画像。CI でバイト単位（または画素単位）に照合する |
| `test/` | Playwright による golden テストと、各形式の `testdata/` を作るフィクスチャ生成スクリプト |

各パッケージが何を公開しているかは [API 一覧](../api.md)を、これらのビルド・テスト・再生成の方法は[ビルドとテスト](development.ja.md)を参照してください。
