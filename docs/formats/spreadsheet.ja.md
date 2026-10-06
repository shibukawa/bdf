# Excel・CSV・Parquet

Excel のブック、CSV・TSV の表、Apache Parquet のファイルは、BDF の中では同じものになる。表形式のデータのための無限平面のレイアウト、`sheet` View だ（[spec.md §4.1](../spec.md#41-view-の種類)）。CSV・TSV と Parquet は自分の描画コードを持たない。`converter/csv` と `converter/parquet` は値だけの単純な格子を組み立て、それをワークシートのセルを組んで描いている `converter/xlsx` のレイアウト・描画コードにそのまま渡す。だから CSV や Parquet のファイルは、ビューアの中で Excel のシートとまったく同じに見え、同じに振る舞う。太字で固定した見出し行、グリッド線と固定ペイン、表計算ソフトと同じ選び方（ドラッグ、Shift、行・列の見出し、矢印キー）で選べるセル。コピーすれば画像ではなくタブ区切りの値と HTML の表になる。

## 試す

ファイルを[ビューア](https://shibukawa.github.io/bdf/viewer/)にドロップするか、サンプルを試せる: [Excel のブック](https://shibukawa.github.io/bdf/viewer/?file=samples/features.xlsx)、[数式オブジェクトのあるシート](https://shibukawa.github.io/bdf/viewer/?file=samples/math.xlsx)、[Shift_JIS の TSV ファイル](https://shibukawa.github.io/bdf/viewer/?file=samples/japanese.tsv)、[Parquet の表](https://shibukawa.github.io/bdf/viewer/?file=samples/basic.parquet)。

## Excel（.xlsx）

`converter/xlsx` は OPC の zip パッケージ（.xlsx、.xlsm、.xltx、.xltm）を読み、ワークシートごとに `sheet` View を作る。セルの配置・書式・表示形式の解釈はすべて変換器自身が行う。BDF に再レイアウトはないので、ここが解釈の唯一の場所になり、シートは Tile に描かれる。セルの値はファイルに保存されているもの（数式の計算結果のキャッシュ）で、数式そのものは計算しない。

読み込むもの:
- 表示形式（日付、和暦、分数、会計など Excel の書式コードが表せるもの一式）、フォントとリッチテキスト、塗りと罫線、配置（和文の禁則付きの折り返し、空きセルへのはみ出し、回転、縮小して全体を表示）、セルの結合。
- 条件付き書式: カラースケール、データバー、アイコンセット、数式を使う規則。これらは BDF が内蔵する小さな数式評価器（四則演算、セル参照と範囲、論理、よく使う検索・文字列・日付の関数）で評価する。開いた日時に依存する規則は評価せずに警告を出す。
- テーブルとその組み込み・独自のスタイル、画像・図形・グラフ。いずれも 1 回だけ描き、重なる Tile から使い回す。グラフシートはシートではなく、独立した `fixed` View のページになる。
- 列幅と行の高さは、Excel 自身と同じ規則（既定のフォントの数字の幅が基準）で計算する。固定ペインとグリッド線の有無は manifest に書き、ビューアが描く。
- 読み上げのためのセル単位の構造: 値のあるセルはそれぞれ `MARK CELL` になり、固定した見出し行やテーブルの見出し行のセルは列見出しとして印を付ける。

読み込まないもの: 右から左のシート（左から右に描く）、ピボットテーブルのスタイル（値は普通のセルとして描く）、スパークライン、フォームコントロール、旧形式（VML）の図形、セル内の画像、古いバイナリ形式の .xls。

オプション（`-param`。`bdf generate -h` の内容）:

| オプション | 値 | 既定 |
|---|---|---|
| `-param hidden=` | `true` | 非表示シートを除く |

`-hidden`（`-param hidden=true` の短縮形）で非表示シートも含められる。マクロシートとダイアログシートは常に読まず、警告を出す。詳細は [design.md §3.6](../design.md#36-excel--bdf-変換器converterxlsxの構造) を、シートのセルがどのように Tile に分割されるかは [design.md §6](../design.md#6-excel-シートの-tile-化) を参照。

### 数式オブジェクト

シートに「挿入 → 数式」で置いた数式（テキストボックスや図形の中の Office Math）は、PowerPoint と同じく `a14:m` の Office Math から組む。サンプルは [`converter/xlsx/testdata/math.xlsx`](https://github.com/shibukawa/bdf/blob/main/converter/xlsx/testdata/math.xlsx)。セルの計算式（`=SUM(A1:A3)`）はこれとは別で、保存されている計算結果を表示する。組み方は [Word・HTML・Markdown の数式](document.ja.md#数式)を参照。

## CSV・TSV

`converter/csv` はテキストの表を読み、Excel で開いたときと同じ見た目の `sheet` View を 1 枚作る。実際に描いているのが Excel の変換器そのものだからだ。このパッケージが受け持つのは、CSV ファイルが明示していないことの推定である: 文字コード、方言、見出し行があるかどうか。どの推定も `-param` で上書きでき、実際にした推定は変換の要約（`Result`、CLI の最終行）に出る。

推定するもの:
- **文字コード**: バイト順マーク（BOM）があればそれに従う（UTF-8、UTF-16LE/BE、UTF-32LE/BE）。なければ先頭 64 KiB から推定する。BOM のない UTF-16 や ISO-2022-JP のエスケープシーケンスに特有のバイトパターン、正しい UTF-8、そのいずれでもなければ Shift_JIS・EUC-JP・Windows-1252 のうち「自然な文章らしさ」の点数が最も高いものを選ぶ。中国語・韓国語の旧来の文字コード（GB18030、Big5、EUC-KR）は推定せず、`-param charset=` で指定する。
- **方言**: 区切り文字（カンマ・タブ・セミコロン・縦棒）とクオートの書き方（ダブル・シングル・なし、二重化かバックスラッシュのエスケープ）の組を先頭 64 KiB に試し、フィールド数が最もそろう組を選ぶ。
- **見出し行**: 列ごとに点数を付ける。数値・日付・真偽値・メールアドレス・URL の列の上にある文字列は見出しとして加点し、同じ列の下にも現れる値はデータとして加点する。

値は書かれたままの表記で示す（Excel は `1.50` を `1.5` に変えてしまう）。先頭が 0 の数値（`007`、郵便番号）は数値にならず文字列のままだが、本物の数値と同じく右揃えになる。列幅は内容に合わせ、改行を含むセルは折り返して行が高くなる。そして Excel の変換器で描いているので、見出し行と判定された行は太字にして固定し、読み上げ用のテキスト層では列見出しとして印を付ける。固定ペインのある Excel のシートと何一つ変わらない。

オプション（`-param`。`bdf generate -h` の内容）:

| オプション | 値 | 既定 |
|---|---|---|
| `-param charset=` | `utf-8`、`shift_jis`、`euc-jp`、`euc-kr`、`gb18030`、`big5`、`windows-1252` など | 推定 |
| `-param delimiter=` | カンマ・タブ・セミコロン・縦棒、または他の文字 | 推定 |
| `-param quote=` | `double`、`single`、`none` | 推定 |
| `-param header=` | `true`、`false` | 推定 |
| `-param table=` | 組み込みの Excel テーブルスタイル（例: `TableStyleMedium2`） | なし |

詳細は [design.md §3.10](../design.md#310-csvtsv--bdf-変換器convertercsvの構造) を参照。

## Apache Parquet

`converter/parquet` は Parquet のファイル（pandas・Polars・DuckDB・Spark が書き、Hugging Face のデータセットもこの形で配られる列指向の形式）を読み、表を 1 枚の `sheet` View にする。描画はここでも `converter/xlsx` に任せる。CSV と違って推定するものはない。Parquet のスキーマに列の型が書かれているので、あとは値をどう文字にするかだけの問題になる。

読み込むもの:
- 読み込みは Arrow を使わず、Parquet 自体の仕様から書いた。apache/arrow-go や parquet-go の parquet パッケージをリンクすると、Office 系フォーマット用の WebAssembly モジュールがほぼ倍の大きさになるため、表示に要る部分だけを実装している。Snappy・gzip・Zstandard（klauspost の実装。速く、Polars と Databricks の既定の圧縮でもある）・LZ4 はどこでも伸張できるが、Brotli の伸張はブラウザ版のビルドでは外し、モジュールを小さく保っている。
- 辞書ページとデータページ（v1、v2）、実際に使われているすべてのエンコーディング: PLAIN、RLE とビットパックのハイブリッド、PLAIN/RLE の辞書エンコーディング、DELTA 系のエンコーディング、BYTE_STREAM_SPLIT。
- 入れ子の値（リスト、マップ、構造体、リストのリスト）は、葉の列の繰り返しレベルと定義レベルから組み立て直す。古い書き手が使っていた標準以前のエンコーディングも読む。
- Variant 型の列（Spark 4、Iceberg v3、DuckDB の variant 型）。分解保存（shredding）された値も含む。
- 表示する行に必要な行グループとページだけを読む。ファイル全体をメモリに読み込むことはない。

値はデータツールが表示するのと同じ書き方にする: decimal は全桁、タイムスタンプは UTC で列が必要とする桁数の小数秒、UUID と interval はそれぞれの標準的な表記、リスト・マップ・構造体は JSON、ジオメトリの列（`GEOMETRY`、`GEOGRAPHY`、GeoParquet の WKB 列）は WKT にする。列名は 1 行目に太字で固定し、CSV の見出し行と同じ扱いにする。その下の 2 行目には、各列の Parquet の型（`int64`、`decimal(10, 2)`、`timestamp[ms, UTC]`、`list<string>` など）を灰色の行で固定して示す。

既定では先頭の 10,000 行だけを表示する（`-param rows=` で別の数を、`all` でシートの上限の 1,048,574 行まで指定できる）。読み込み自体は速い（100 万行・8 列で約 0.6 秒）が、シートの組版はそうではない: 10 万行で約 3 秒・690 MB、100 万行で約 45 秒・2.7 GB かかる。既定の行数が小さいのはこのためだ。`rows=` の指定にかかわらずセル数には 400 万の上限もあり、列の多い表は行数が減る。列数は Excel・CSV と同じくシートの上限の 16,384 列まで。

読み込まないもの: Parquet Modular Encryption（フッターが暗号化されたファイルは判別はできるが開けない）、ブラウザ版での Brotli 圧縮の列、LZO、まだ実験的な扱いの ALP エンコーディング。

オプション（`-param`。`bdf generate -h` の内容）:

| オプション | 値 | 既定 |
|---|---|---|
| `-param rows=` | 数値、または `all` | `10000` |
| `-param types=` | `true`、`false` | `true` |
| `-param table=` | 組み込みの Excel テーブルスタイル（例: `TableStyleMedium2`） | なし |

詳細は [design.md §3.26](../design.md#326-parquet--bdf-変換器converterparquetの構造) を参照。
