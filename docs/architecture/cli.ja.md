# bdf コマンド

```
bdf generate [flags] <input> <out.bdf | dir/>
                                   文書や図面を変換する（フラグと形式一覧は bdf generate -h）
bdf fill [flags] <template.xlsx | .docx | .pptx | .pdf | .json | .bdft> <out.bdf | dir/>
                                   帳票のテンプレートをデータで穴埋めする: タグを書いた Office の文書、
                                   入力欄のある PDF、pdfme のテンプレート（フラグは bdf fill -h）
bdf template compile [flags] <template> <out.bdft>
bdf template schema [-json-schema] <template> [out.json]
                                   テンプレートを検査してコンパイル済みにする、タグが引く名前を書き出す
bdf pdf [flags] <file.bdf | dir> <out.pdf>
                                   ページを PDF に書き出す（フラグは bdf pdf -h）
bdf thumbnail [flags] <file.bdf | dir | input> <out.png | .jpg | .webp>
                                   サムネイルを描く: 1 ページ目の上部、またはスライド全体
                                   （変換する入力なら 1 ページ目だけを変換する）
bdf text [flags] <file.bdf | dir> [out.json]
                                   検索用に、メタデータとページごとのテキストを書き出す
bdf render [flags] <file.bdf | dir> <out.png | .jpg | .webp>
                                   1 ページを描く（bdf thumbnail -h, text -h, render -h でフラグ一覧）
bdf ls <file.bdf | dir>            View と Part の一覧
bdf manifest <file.bdf | dir>      manifest を JSON で出力
bdf disasm <file.bdf | dir> <hash> Object Part を逆アセンブル
bdf extract <file.bdf | dir> <hash> <out>
bdf split <file.bdf> <dir>         分割形式で書き出す
bdf join <dir> <file.bdf>          1 ファイル形式で書き出す
bdf encrypt [-password-file f] <file.bdf | dir> <out.bdf | dir/>
                                   パスワードで文書を暗号化する
bdf decrypt [-password-file f] <file.bdf | dir> <out.bdf | dir/>
                                   暗号化を解いて書き出す
bdf demo <file.bdf | dir/>         フィクスチャ文書を書き出す（dir/ は末尾がスラッシュ）
bdf conformance [-testdata dir] <dir>
                                   他の読み手のための適合テストを書き出す（conformance パッケージ）
```

暗号化された文書は `$BDF_PASSWORD` のパスワードで読みます。`split` と `join` はパスワードなしで使えます。

## generate

1 つの入力を BDF に変換します。形式は入力の中身から判別するか、`-format` で明示します。出力は 1 ファイル、または末尾が `/` のディレクトリで分割形式になります。`-pages 1-3,5,8-` でページ・スライド・シート・アートボードを選び、`-thumbnail` と `-text` は出力と一緒にプレビューを書き出します（暗号化した出力には `-allow-plaintext` のときだけ）。`-dc 要素名=値` は Dublin Core の要素を上書きします（`-dc 要素名=` で削除）。入力形式ごとの `-param` オプションは[対応形式](../formats/index.ja.md)を、現在の完全な一覧は `bdf generate -h` を参照してください。

```sh
bdf generate report.pptx report.bdf              # 1 ファイル形式
bdf generate report.pptx out/                    # 分割形式
bdf generate -pages 1-3 book.pdf book.bdf         # 最初の 3 ページだけ
bdf generate -password-file pw.txt secret.xlsx out.bdf  # パスワード付きの入力
```

## fill, template

帳票のテンプレート（文字のところにタグを書いた Excel のブック、Word の文書、PowerPoint のプレゼンテーション）をデータで穴埋めします。形式は中身から判定します。書き方は[帳票を作る](../forms.ja.md)を参照してください。

```sh
bdf fill -data order.yaml -data items=lines.csv invoice.xlsx invoice.bdf   # 穴埋めして、印刷のページの文書に
bdf fill -data order.yaml notice.docx notice.bdf                           # Word の文書はページに、PowerPoint はスライドに
bdf fill -each orders.csv invoice.xlsx 'out/{{no}}.bdf'                    # レコードごとに 1 つずつ
bdf template compile -sample sample.yaml invoice.xlsx invoice.bdft         # 検査して、コンパイル済みテンプレートに
bdf template schema invoice.bdft                                           # タグが引く名前を JSON で
bdf template schema -json-schema invoice.xlsx                              # データの JSON Schema（入力フォームを作る元）
bdf fill -data applicant.json application.pdf application.bdf              # 入力欄のある PDF: 欄の場所に値を描く
bdf fill -data invoice.json pdfme-template.json invoice.bdf                # pdfme のテンプレート: pdfme の生成器と同じ規則で埋める
```

`-data` は JSON・YAML のオブジェクト（`file`）か、名前を付けた値（`name=file`。CSV・TSV やブックのシートはオブジェクトのリスト）で、繰り返し指定できます。`-strict` はデータにない名前のタグをエラーにし、`-files` はデータが名前で指す画像のディレクトリを決めます。フォントと画像のフラグ（`-font-dir`、`-no-system-fonts`、`-fonts`、`-images` など）と `-param`・`-pages` は `generate` と同じ意味です（`-pages` は、Excel ではシート、PowerPoint ではテンプレートのスライド、Word ではできた帳票のページを選びます）。`.bdft` は BDF の文書なので、`ls`・`thumbnail`・`render` でそのまま開けます。

## pdf

文書のページを PDF に書き出します。文字は埋め込んだフォントの文字として、図形はベクタとして入ります（[PDF にする](../forms.ja.md#pdf-にする)）。

```sh
bdf pdf invoice.bdf invoice.pdf
bdf pdf -view print1 book.bdf sheet1.pdf            # ビューを選ぶ（繰り返せる）
bdf pdf -skip background -skip master slip.bdf slip.pdf   # レイヤーを役割で省く
bdf pdf -tagged report.bdf report.pdf               # タグ付き PDF: 文書の構造を持ち、スクリーンリーダーで読める
bdf pdf -pdfa 2a invoice.bdf invoice.pdf            # 保存用の PDF/A（2b、2u、2a、3b、3u、3a）
```

既定では、ページのあるビュー（`fixed`・`flow`）をすべて書き出し、ページのない文書は全体を 1 ページにします。`-font-dir` と `-no-system-fonts` は、文書が名前で参照するフォントを探す場所です。`-ignore-fstype` は、埋め込みを許可していないシステムのフォントも埋め込みます（権利があるときだけ。既定ではその文字をアウトラインで描きます）。`-no-compression` はストリームを圧縮せずに書きます（中身を読むため）。暗号化した文書には `-allow-plaintext` が必要です（PDF は暗号化されません）。

## thumbnail, text, render

ブラウザなしで `.bdf` を見る 3 つの方法です。サーバーが使うのと同じ Go のコード（`imagebdf`、`thumbnail`、`Document.SearchText`）を使います。

```sh
bdf thumbnail -size 256 out.bdf thumb.png    # -mode auto|crop|fit, -view, -sheet-dpi, -font-dir, -no-system-fonts
bdf thumbnail -size 256 book.pdf thumb.png   # generate が変換できる入力なら 1 ページ目だけを変換する
bdf text out.bdf text.json                   # メタデータとページごとのテキストを JSON で（既定は標準出力）
bdf render -page 2 -scale 2 out.bdf page2.png  # シートは -width/-height も（A1 から）
bdf render -antialias=off -scale 8.333333 out.bdf print.png  # 600 dpi。境界の AA を省いて速度を優先
bdf render -pages=all -workers=4 -antialias=off out.bdf 'page-{page}.png'  # 独立ページを並列に描き、指定順に保存
bdf render -pages=all -antialias=off -mono=bilevel -scale 8.333333 out.bdf 'page-{page}.png'  # 600 dpi の 2 値（1 画素 1 bit）
```

`render -antialias=default|on|off` で図形・文字・幾何クリップの境界のアンチエイリアスを選べます。`default` と `on` は倍率によらず従来の描画を維持し、`off` は境界の被覆率計算を省きます。dpi による自動切り替えはありません。半透明・ブレンド・画像の補間・ソフトマスク・ぼかしと、1 画素より細い線の薄い alpha は維持します。`thumbnail` と変換時のプレビューは従来のアンチエイリアスを使います。Go の指定と出力上の注意は[描画する](../rendering.ja.md#ラスタライズの品質と速度)を参照してください。`render -mono=gray|bilevel` はカラーの代わりにモノクロの画像を書きます。1 画素 8 bit の輝度か、黒と白の 2 値（PNG では 1 画素 1 bit。グレー値 128 未満が黒）で、1 ページでも `-pages` のバッチでも使えます（[モノクロの出力](../rendering.ja.md#モノクロの出力)）。

`.bdf` でない入力（PDF、`.docx` など）を渡すと、`thumbnail` はサムネイルに必要な 1 ページ目だけを変換します（`-view` を付けたときは全体）。500 ページの PDF でも数秒ではなく数ミリ秒で済みます。パスワード付きの入力には、暗号化した文書と同じく `$BDF_PASSWORD` と `-allow-plaintext` が要ります。

暗号化した文書には `thumbnail` と `text` に `-allow-plaintext` が必要です（`$BDF_PASSWORD` で先に復号）。どちらの出力も暗号化されないため、明示的に許可する仕組みです。それぞれが何を描く・抽出するかは[変換と出力 → サーバー側のサムネイルと検索用テキスト](conversion.ja.md#サーバー側のサムネイルと検索用テキスト)を参照してください。

## 文書の中身を見る

```sh
bdf ls out.bdf              # すべての Part: ハッシュ、種類、エンコーディング、圧縮前後のサイズ
bdf manifest out.bdf        # manifest を JSON で: View、ページ、メタデータ
bdf disasm out.bdf <hash>   # 1 つの Object Part の描画命令を逆アセンブル
bdf extract out.bdf <hash> out.bin  # Part の生バイト列（フォント、画像など）
```

## 2 つのコンテナ形式の変換

```sh
bdf split out.bdf out/    # 1 ファイル形式 → 分割形式
bdf join out/ out.bdf     # 分割形式 → 1 ファイル形式
```

どちらもパスワードは要りません。暗号化されていてもいなくても、コンテナの Part をそのままコピーします。

## encrypt, decrypt

```sh
BDF_PASSWORD=… bdf encrypt in.bdf out.bdf   # 既存の bdf をパスワードで封印
BDF_PASSWORD=… bdf decrypt in.bdf out.bdf   # 暗号化を解除（開くのにパスワードが必要）
```

`generate` は `-encrypt auto`（既定: 入力にパスワードが要ったときだけ）・`always`・`never` で、この 2 段階を一度にできます。

## demo

```sh
bdf demo out.bdf      # golden テストとデモビューアの既定が使うフィクスチャ文書
bdf demo out/         # その分割形式
```

## conformance

```sh
bdf conformance out/                 # suite.json と inputs/。testdata/ の文書を受理の事例として含む
bdf conformance -testdata "" out/    # 合成した事例だけ
```

他の言語で書かれた読み手（libbdf）のためのもので、入力と、この実装がそれをどう扱うかを `suite.json` に記録します（`conformance` パッケージ）。CI は成果物 `bdf-conformance` として、リリースのワークフローはタグのリリースに `bdf-conformance-<tag>.tar.gz` として添付します。

## 形式を横断するフォント埋め込みのフラグ

テキストをレイアウトして埋め込むすべての形式（PowerPoint、Excel、Word、CSV、Parquet、メタファイル、HTML、Markdown、EPUB、フォントファイル）に共通のフラグです。

| フラグ | 効果 |
|---|---|
| `-font-dir DIR` | システムのフォントより先にこのディレクトリを探す（繰り返し指定可） |
| `-fonts embed \| system` | レイアウトに使ったフォントを埋め込む（サブセット化）か、名前で参照するか（ビューア自身のフォントで描く。HTML・Markdown・EPUB の既定） |
| `-no-system-fonts` | `-font-dir` のフォントだけを使い、システムのフォントは使わない |
| `-no-woff2` | 埋め込みフォントを WOFF2 ではなく TrueType/OpenType のまま格納 |
| `-no-subset` | 使う字形だけでなくフォントファイル全体を埋め込む |
| `-ignore-fstype` | OS/2 の `fsType` が埋め込みやサブセット化を禁じるフォントも埋め込む（権利がある場合のみ） |
| `-compression normal \| best \| fast` | 部品の deflate の強さ: `normal`（レベル 6。既定）、`best`（レベル 9。数% 小さくなるが圧縮に 4 倍かかる）、`fast`（レベル 1） |

対応する Go の呼び出しは [API 一覧](../api.md)を、`bdf generate` と `bdf demo` で `testdata/` のフィクスチャを再生成する方法は[ビルドとテスト](development.ja.md)を参照してください。
