# なぜ bdf か

## オフィススイートを動かさなくてよい

Office のファイルをブラウザでプレビューするには、サーバーで LibreOffice や OpenOffice をヘッドレスで動かして PDF にするのが定番です。これはインストールだけで 1 GB を超え、何もしていない状態でこの大きさです。プロセスの起動・維持・隔離も必要で、壊れたファイル 1 つでプロセスが詰まることもあります。

bdf の変換器は cgo も外部プログラムも使わない Go のパッケージです。PDF、Word、PowerPoint、Excel、CSV、Parquet、Visio、draw.io、DXF、Jw_cad、SXF、CGM、HP-GL/2、Gerber、Excellon、KiCad、TIFF、Illustrator、Photoshop、メタファイル、HTML、Markdown、EPUB、MML、MIDI、MusicXML、フォントファイルを 1 つのバイナリで変換します。同じコードを WebAssembly にすればブラウザの中でも変換でき（PDF 用が gzip で約 7.3 MB、Office 系・draw.io・CAD・KiCad・フォント用が約 7.8 MB）、ファイルをアップロードする必要すらありません。サムネイルと検索用のテキストも同じプロセスで作れます。ページを描くのは純 Go のラスタライザで、サーバー側にもブラウザは要りません。

## 内容に合った形で見せる

PDF はすべてを紙に切り分けます。スプレッドシートを印刷したページでは、横に長い表がページをまたいで分断されて行を追えず、目当てのセルも見つけにくくなります。固定した見出しや枠線は消え、ブックの中で切り替えていたシートは一続きのページになります。Word の文書もページ単位でしか読めません。

bdf は内容の種類ごとにレイアウトのモデルを持ちます。スライド・図面・PDF には固定サイズのページ、ワークシートにはシートごとの無限平面（タイルで描画し、ウィンドウ枠の固定、行・列見出し、枠線つき）、ワープロ文書にはページでも一続きのスクロールでも読めるフローと、ページなしで 1 本の長い列に組み直した表示を用意しています。Illustrator と Photoshop のアートボードはページになります。ブックのシート、draw.io の図のページ、DXF のモデル空間とレイアウト、プリント基板の表・裏と各層はそれぞれ 1 つの表示になり、ビューアのタブで切り替えます。

## ブラウザ表示に特化している

命令セットは Canvas 2D と 1 対 1 に対応します。フォントは `FontFace` に渡す WOFF2、画像はブラウザがデコードできる形式で、Part の圧縮は `DecompressionStream` で展開できる形式です。pdf.js のような PDF ビューアが数万行かけて実装しているフォントのラスタライズ、画像のデコード、展開はブラウザに任せ、bdf のデコーダとレンダラは TypeScript で約 3,400 行です（レンダラの Worker は gzip で 21 KB）。描画は Worker の `OffscreenCanvas` で行い、メインスレッドはビットマップを置くだけです。Part は内容アドレスなので、マスターや繰り返し現れる要素は 1 度だけ格納され、ビューアは表示中のページに要る Part だけを Range リクエストや CDN 上の分割形式から取得します。ブラウザ内で変換する PDF は、表示中のページを優先して変換できたページから表示します。

## 多くの形式を 1 つのレンダラで

PDF、Word（.docx）、PowerPoint（.pptx）、Excel（.xlsx）、CSV・TSV、Apache Parquet、Visio（.vsdx、.vdx）、draw.io（.drawio と、図を埋め込んだ SVG・PNG の書き出し）、AutoCAD DXF、Jw_cad（.jww）、SXF（.p21、.p2z、.sfc）、CGM（.cgm）、HP-GL/2 のプロットファイル（.plt）、Gerber（RS-274X）と Excellon の穴あけファイル（1 つずつでも、基板のファイルをまとめた ZIP でも）、KiCad の回路図と基板（.kicad_sch、.kicad_pcb、プロジェクトの .kicad_pro、プロジェクトの ZIP）、TIFF、Illustrator（.ai）、Photoshop（.psd、.psb）、Windows メタファイル（.emf、.wmf）、HTML（リーダー表示）、Markdown、EPUB（リフロー型の本は和文の縦書きも、固定レイアウトのマンガも）、楽譜に組む MML（.mml）・MIDI（.mid、.kar）・MusicXML（.musicxml、.mxl）、文字・グリフ・OpenType フィーチャーを見せるフォントファイル（.ttf、.otf、.ttc、.woff、.woff2）を、パスワード付きの Office 文書や PDF も含めて同じフォーマットにします。どの形式も同じレンダラで描き、検索、テキスト選択、読み上げ用のテキスト層（見出し、リスト、表、代替テキスト）も共通です。

各変換器が何を保ち、どうレイアウトするかは[対応形式](formats/index.ja.html)を、変換と描画がどうつながっているかは[アーキテクチャ](architecture/index.ja.html)を参照してください。
