# テキスト帳票テンプレート 実装案（form-m6）

2026-10-20 の実装案。同日 S1〜S6 を実装した（`.knowledge/` の `milestone:form-m6`。計画から変えた所はそこの `changed_from_plan`）。要件は `text-template.ja.md` と、`.knowledge/` の `requirement:text-template`（決定 6 つは 2026-10-19 に利用者が決めた）。印刷の print-m5（ラベル・レシート・連帳。ZPL / SBPL / TPCL、ESC/POS、ESC/P の書き手）と並行して、別のセッションで作った。2026-10-20 利用者が「form-m6 を別セッションで並行して実装する」と決め、print-m5 はラスタの出力と BARCODE の MARK を先に作り、テキストモードはこの Part が入ってから足す。

## 0. 範囲

form-m6 は「プレーンテキストのテンプレートを、穴埋めから `.bdf` と `.txt` まで通す」。ESC/P と ESC/POS のテキストモードの書き手は作らない（print-m5）。その書き手が読む印字像の Part（`t: "text"`）の形と読み手は、ここで決めて固定する。

1. **テンプレートの読み込み**: 前書き（`---` の YAML）と本文、タグの位置（行）、区間の木。
2. **格子への穴埋め**: 東アジアの文字幅、`width` / `right` / `center` と `cut` / `wrap`、区間（同じ行・枠・`rows`・`pages`）、`title`、FF、`{{#preprinted}}` と倍角。
3. **BDF のページ**: 等幅フォントの fixed ページ、run の advance、MARK、用紙の層（`background`）。
4. **印字像の Part**: `t: "text"` の JSON。型と読み手は `form` も変換器も import しない小さなパッケージに置く。
5. **プレーンテキスト**: `bdf fill … out.txt`（FF 区切り、`-param encoding=`）。
6. **`form` への組み込み**: 判別、Fill・Fields・Compile（`.bdft`）、スキーマ、CLI。

**入れないもの**: ESC/P・ESC/POS の書き手（print-m5）、画像とバーコード（決定どおり）、縦書き、外字の解決、Office のテンプレートからの格子の推定、ページ番号を本文に入れること（`{{_page}}` は `| pages` の区間の中だけ。`rows` の改ページで生まれたページの番号は、穴埋めの後でしか決まらない）。

## 1. 決めたこと（暫定で推奨どおり進めた）

- **置き場所**: 埋め手は `converter/textform`（計画の `form/text` から変えた）。フォントの選択と埋め込み（`converter/internal/fontset`）とオブジェクトの組み立て（`converter/internal/canvas`）は `converter/internal` にあり、Go の決まりで `converter/` の下からしか import できない。pdfme（`converter/pdfme`）と PDF フォーム（`converter/pdfform`）と同じ並び。印字像の型と読み手は `form/textpart` で、import するのは `bdf` と `golang.org/x/text/width` だけ（`print` が読むため）。
- **前書きの読み手**: YAML のうち、キーと値（素の値と引用符の値）、`margin` の流れ形式（`{left: 2, top: 1}`）と字下げのブロック、コメントだけを読む自前の小さな読み手にした。`go.yaml.in/yaml/v3` は `form` のデータの読み込みで TinyGo のビルドから外しているもので、テンプレートの読み込みをそれに頼るとブラウザの wasm で前書きが読めなくなる。知らないキーは警告。
- **キーを 1 つ足した**: `ambiguous: narrow | wide`。東アジアの文字幅が「あいまい」な文字（罫線素片、`○`、`※`、ギリシャ文字）を 1 桁と数えるか 2 桁と数えるか。叙述の例（`┌──┬` と `│ 品名 │`）は 1 桁で数えたときだけ揃い、多くのエディタもそう見せるので既定は `narrow`。一方、ESC/P-J の漢字モードは罫線素片を全角で印字するので、そのプリンタで揃えたいテンプレートは `wide` と書く。Part の `ambiguousWide` に残り、テキストモードの書き手はそれで桁を数える。
- **幅のない値の改行**: 幅を書かない値の中の改行は、`| wrap` と同じく下の行の同じ桁に続ける（住所を複数行で渡したデータ）。幅を書いて `wrap` のない値では、改行を空白にして 1 行にする。
- **タブ**: 8 桁ごとの位置まで空白にする（ESC/P の既定のタブ位置）。制御文字は捨てる。
- **枠の行**: Excel と同じく、タグのない枠の行は最初の枠の行と同じものとして埋める（罫線の文字はそのまま）。
- **はみ出し**: `| rows` の区間のない帳票（固定のページ）は、ページに収まらなければエラー。ただしページの外にこぼれた空行だけは捨てる（ファイルの末尾の空行で落ちないように）。`| rows` があれば自動で改ページして `title` を繰り返す。`title` は「帳票の最初の n 行」（テンプレートの行ではなく、穴埋めの後の行）。
- **`| pages`**: 2 つ目以降の項目の前で改ページする（最初の項目は、その前の行と同じページ）。Word の `| pages` と同じ。
- **倍高の行**: その行は 2 行を占め、行の文字はすべて下の行のベースラインに立つ（ESC/POS が倍高の文字を含む行を印字するのと同じ）。Part では、その行の次の要素を空の行にして、`lines[i]` が常に格子の i 行目になるようにした。
- **描き方**: フォントの大きさは行の高さ、ベースラインは行の上から 0.88（CJK の仮想ボディの下端が行の下端に来るので、罫線素片が上下の行とつながる）。文字は run の advance（桁の幅）で格子に合わせる。最初は「半角と全角で run を分ける」だけにしたが、BIZ UDゴシックの罫線素片は全角の幅を持つ（あいまいを 1 桁にすると半角の文字と比が違う）ので、`│ 品名` の `│` と空白を同じ run にすると run の中で縦線がずれた。run は「フォントでの幅 ÷ 桁数」が等しい文字だけでまとめることにした。等幅のフォントなら半角と全角が 1 つの run になり、比例のフォントでは 1 文字ずつになる。どのフォントにもない文字は 1 文字で 1 run（advance がほかの文字を引き伸ばさないように）。
- **倍角の描き方**: 倍幅は advance を 2 倍、倍高はフォントを 2 倍にして advance はそのまま、倍角は両方。変換の命令（`TRANSFORM`）を足さずに済む。
- **MARK**: 行ごとに PARAGRAPH。`rows` と枠の区間の行は TABLE と CELL（1 列。`A1`、`A2` …）。列の区切りはテンプレートからは決まらない（罫線のない帳票もある）ので 1 列にした。用紙の層のレイヤーは PARAGRAPH だけ。
- **フォント**: マシンのフォントから、BIZ UDゴシック、IPAゴシック、Noto Sans Mono CJK JP、各 OS の等幅の順に最初にあるものを選び、字のないものは残りの順に探す。前書きの `font:` は、拡張子が `.ttf` / `.otf` / `.ttc` / `.otc` ならファイル、そうでなければファミリー名。ファイルは `TemplateData.Files`（CLI では `-files`、既定はテンプレートのディレクトリ）から読み、その外は読まない（テンプレートは信頼できない入力）。
- **プレーンテキスト**: `bdf fill … out.txt` は、文書の Part を `textform.Plain` でテキストにする（描いたページと同じものから書く）。余白は空白と空行で書く（固定長の印字像そのもの）。`-param encoding=` は WHATWG の名前（`golang.org/x/text/encoding/htmlindex`。`shift_jis` は Windows-31J）、ない文字は `〓`（それもなければ `?`）で警告。ほかに `newline=crlf` と `paper=preprinted`（用紙の層を空白に）を足した。
- **予約した名前**: `preprinted` は区間の名前としてデータを引かない（叙述どおり。`_` で始まらないので、データに同じ名前があれば引けなくなる）。
- **判別**: UTF-8 のテキスト（NUL を含まない）で、`---` の前書きで始まるか、タグを持つ。最初の 1 KiB にタグがなければ 1 MiB まで読んで探す。タグのない、前書きもないテキストはテンプレートではない。

## 2. 印字像の Part（`t: "text"`、版 1）

```jsonc
{
  "version": 1,
  "cpi": 10, "lpi": 6,              // 1 インチの桁数と行数
  "columns": 80, "lines": 33,       // 用紙（余白を含む）。lines 0 はロール紙
  "margin": {"left": 2, "top": 1},  // 桁と行
  "skip": 2,                        // ミシン目スキップ（省略は 0）
  "ambiguousWide": false,           // あいまいな幅の文字を 2 桁と数えたか（省略は false）
  "font": "BIZ UDGothic",           // ページを描いたフォント
  "pages": [                        // 文書の最初の View のページの順
    {"lines": [                     // lines[i] は上の余白の下の i 行目
      {"text": "納品書", "size": [{"col": 18, "cols": 12, "size": "big"}]},
      {"text": ""},                 // 倍高の行の下の行
      {"text": "┌──┐", "paper": [{"col": 0, "cols": 4}]}
    ]}
  ]
}
```

- `text` は行の文字（左の余白を含まず、行末の空白なし）。倍幅の文字も 1 回だけ入る。桁は `textpart.Width`（全角・幅広 2、結合文字と制御文字 0、あいまいは `ambiguousWide`）で数え、`wide` と `big` の範囲では 2 倍。`Image.Cells(line)` が文字ごとの桁と属性を返す。
- `paper` はプレ印刷の範囲、`size` は倍角の範囲（`wide` / `tall` / `big`）。どちらも重ならず、桁の順。`col` と `cols` は印字の桁（倍幅の分を含む）。
- 読み手（`textpart.Decode`）は版、ピッチ、用紙、余白、ページと行の数、範囲の重なりを確かめ、外れた Part を断る。
- ESC/P の書き手（print-m5）は、`cpi` → `ESC P` / `ESC M` / `ESC g`、`lpi` → `ESC 2` / `ESC 0`（ほかは `ESC 3 n`）、`lines` → `ESC C n`（0 は付けない）、`skip` → `ESC N n`、`margin.left` → `ESC l`、`margin.top` → 空行、`size` → `ESC W` / `ESC w`（ESC/POS は `GS !`）、`paper` → チケットが「プレ印刷」なら空白、ページの終わり → `FF` と写す。

## 3. 構成

```
form/textpart/textpart.go      印字像の型（Image、Page、Line、Range）、Width、Cells、Text、Encode・Decode、Read・FromDocument・AddTo
form/tag/tag.go                ParseText（値の width・right・center、cut・wrap。Token の Width・Align・Overflow）
converter/textform/frontmatter.go  前書きの読み手と格子
converter/textform/parse.go    本文の行、タグ、区間の木、フィールド、IsText・IsTemplate
converter/textform/fill.go     穴埋め（値の幅、区間、続きの行）、ページ分け、印字像
converter/textform/layout.go   印字像 → fixed ページ（フォント、run、MARK、用紙の層）
converter/textform/textform.go Fill・Show・Inspect、フォントファイル、Plain（文字コード）
form/form.go                   FormatText の判別と Fill・Fields・Compile の分岐
cmd/bdf/fill.go                out.txt、使い方の文
manifest.go                    bdf.PartText
docs/                          spec §4 の一覧、forms（テキストのテンプレート）、cli、api、design §3.34
conformance/inputs.go          参照されない Part（text と未知の型）の事例
```

## 4. 順番（スライス）

1. **S1** 前書きと本文の読み込み、タグの位置、区間の木（`frontmatter.go`、`parse.go`、`tag.ParseText`）。
2. **S2** 幅と揃え（東アジアの文字幅、`width` / `right` / `center`、`cut` / `wrap`）。
3. **S3** 区間の展開（枠・`rows`・`pages`、`title` の繰り返し、FF、倍角と用紙の層）。
4. **S4** BDF の fixed ページ（`layout.go`）。テストのフォントは BIZ UDゴシックのサブセット。
5. **S5** 印字像の Part（`form/textpart`）とプレーンテキスト、spec §4 の追記、conformance。
6. **S6** `form` の判別と分岐、`.bdft`、スキーマ、CLI、docs、CI のテストデータ、`.knowledge`。

## 5. テスト

- **単体**（`converter/textform`）: 前書き（既定、流れ形式とブロック、コメント、誤り）、判別、値の幅（半角・全角・揃え・書式）、超えた値（警告・cut・wrap、続きの行の桁）、改行とタブ、欠けた名前と `-strict`、区間（文字・逆・枠と `_blank_first`・`rows` の入れ子・同じ行の `| rows`・誤りの文言）、ページ（`rows` の改ページと `title`、FF、`| pages`、固定のページのはみ出し、ページの外の空行、ロール紙、`-pages`）、倍角と用紙の層（Part の範囲、レイヤーの役割）、run の位置・advance・フォントの大きさと MARK、フォントファイル、タグのままのページ（`Show`）、プレーンテキスト（Shift_JIS、CR LF、プレ印刷、`〓`）、手数の上限、`FuzzFill`。
- **帳票**: 納品書（固定のページ＋用紙の層＋倍角＋cut / wrap）、出荷明細（`rows` で 3 ページ、`title`）、請求明細書（`pages` と枠、`_page` / `_pages`、倍幅・倍高）、レシート（ロール紙）の 4 つを、`test/text/data` のデータで穴埋めして、印字像のテキストを `testdata/want/*.txt` と比べる。`testdata/form/text-*.bdf` と `text-slip.bdft`、Shift_JIS の `text-slip.txt` を CI が作り直してバイトで比べる。
- **`form`**: 4 つの帳票の判別・フィールド・スキーマ、データで穴埋めした `.bdft` が同じ帳票を作ること、データなしの `.bdft` がタグのままのページを見せること、最初の 1 KiB より後にタグがあるテキスト。
- **読み手**: `form/textpart` の往復と誤り、`bdf.OpenSingle` が Part のある文書を開くこと。conformance に「参照されない `text` と未知の型の Part」の事例、TS の `packages/core` のテストで `text-slip.bdf` を開く。
- 画素の golden は #99 で無くなったので作らない。目では、`bdf render` の絵で罫線のつながり（上の run の決め方）を確かめた。

## 6. print-m5 への引き継ぎ

- テキストモードの書き手は `textpart.Read(rd)`（なければ `textpart.ErrNone`）で印字像を取り、`Image.Cells(line)` で文字・桁・倍角・用紙を歩く。プレーンテキストの `Image.Text` も、書き手のテストの期待値に使える。
- 漢字モード（ESC/P-J の `FS &`、ESC/POS の `FS &`）では、全角の文字は 2 桁、`ambiguousWide` が false のテンプレートの罫線素片は 1 桁で書かれている。1 桁の罫線素片は JIS X 0208 の全角では幅が合わないので、PC437 などの 1 バイトの罫線にするか、`ambiguous: wide` を勧める警告を出すかを m5 で決める。
- `ESC N`（スキップ）を送るなら、Part の `lines[]` はスキップの分を含まない（ページの行は `lines - margin.top - skip` まで）。
