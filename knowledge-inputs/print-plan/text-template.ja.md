# テキスト帳票テンプレート 要件整理（form: text）

2026-10-19 利用者の提案「いまは docx / xlsx / pptx / pdfme のテンプレートしかないが、ESC/P 用にプレーンテキストのテンプレートも用意した方がよさそう」を受けた整理。構造化した記録は `requirement:text-template`、最初の実装のまとまりは `milestone:form-m6`。印刷側の要件は `print-requirements.ja.md` §3.9（ESC/P のテキストモード）。

## 0. 位置づけ

連続帳票（複写伝票、納品書、受領書、給与明細）は「行 × 桁の格子」に文字を置くもので、罫線と項目名はプレ印刷用紙が持つことが多い（`pattern:continuous-text-form`）。Office のテンプレートでこれを作るのは遠回りで、Excel の等幅フォント＋固定行高のシートを格子とみなす案は、列幅がプリンタ解像度で丸まる Excel の癖に引きずられる。**プレーンテキストのテンプレート**なら、エディタで見たとおりの行と桁がそのまま印字像になり、ドットインパクトのテキストモード（速く、打点が濃く、複写に強い）に 1:1 で写せる。

同じテンプレートが 4 つの出口を持つ。

1. **BDF**（画面のプレビュー、PDF、IPP のプリンタへのラスタ / PDF）。等幅フォントで格子に置く。
2. **ESC/P のテキストモード**（ドットインパクト、連帳）。文字と制御コードを送る。
3. **ESC/POS のテキストモード**（レシート）。幅固定・長さ可変の同じ格子。
4. **プレーンテキスト**（固定長の印字像そのもの。旧来のシステムへの受け渡し、メール、ログ）。

既存の仕組み（タグ `form/tag`、データ `form.LoadData`、`-each` のバッチ、`.bdft` のコンパイル、スキーマ、CLI）はそのまま使い、テンプレートの形式を 1 つ足すだけにする。

## 1. ゴール

- **G1 同じタグ・同じデータ**: `{{name}}`、`{{#items}} … {{/items}}`、`{{^}}`、`fmt:` の書式、全角の波括弧、`\{{`、`-data` / `-each`、`-strict`、`.bdft`、`bdf template schema` がそのまま効く。
- **G2 見たまま**: テンプレートの行と桁が印字の行と桁。全角は 2 桁。エディタの等幅表示が最終形のプレビューになる。
- **G3 格子の幾何をテンプレートが持つ**: 桁ピッチ（CPI）、行送り（LPI）、用紙長（行数かインチ）、桁数、左と上の余白、ミシン目スキップ。ESC/P の命令（`ESC C`、`ESC N`、`ESC 2` / `ESC 0`、`ESC P` / `M` / `g`、`ESC l` / `Q`）に 1:1 で写る。
- **G4 BDF に格子が残る**: できた `.bdf` から ESC/P のテキストモードを損失なしに作れるよう、印字像（行 × 桁の文字）と幾何を Part に持つ。ビューアと PDF は fixed ページとして見せる。
- **G5 プレ印刷用紙**: 罫線・項目名を書いた行や範囲を「用紙の層」と印せば、プレビューには出て、追い刷り（`SkipRoles` / 印刷側の設定）では抜ける。位置合わせのオフセット（上下左右）は印刷時のチケットで。
- **G6 レシートも同じ型**: 用紙長なし（長さ可変）、桁数 42 / 48、で ESC/POS のテキストモードに写る。

## 2. 非ゴール

- プロポーショナルフォントの組版、折り返し、箱の中の収め直し（それは Office のテンプレートの仕事）。
- 式・計算（既存の方針どおりデータ側）。
- 画像とバーコード（テキストモードには出せない。BDF の出力でだけ `{{%}}` を許すかは論点）。
- 外字の自動解決（私用領域の字は警告。`decision:font-rules` の外字規則に従う）。
- 縦書き。

## 3. テンプレートの形式

UTF-8 のテキストファイル。先頭に任意の前書き（`---` で囲んだ YAML）、続いて本文。本文がページの印字像で、1 行が 1 行、1 文字が半角 1 桁、全角 2 桁。

```text
---
cpi: 10            # 桁ピッチ（10 / 12 / 15 / 17.1。既定 10）
lpi: 6             # 行送り（6 / 8。既定 6）
columns: 136       # 1 行の桁数（既定 80。15 インチ幅の連帳は 136）
lines: 66          # 用紙長を行数で（11 インチ × 6 LPI = 66。0 は長さ可変 = レシート）
margin: {left: 2, top: 1}   # 桁と行で
skip: 3            # ミシン目スキップの行数（ESC N）
title: 6           # 先頭の行数を各ページに繰り返す（印刷タイトル）
font: "BIZ UDゴシック"       # BDF の表示に使う等幅フォント（既定は font-rules の等幅の連鎖）
---
                  納　品　書                      {{date | fmt:"yyyy年m月d日"}}
{{customer | width:30}} 御中                       No. {{no | right:8}}
{{#preprinted}}
┌──────────────────────┬──────┬────────────┐
│ 品名                 │ 数量 │ 金額       │
{{/preprinted}}
{{#rows}}
│ {{name | width:20}} │ {{qty | right:4}} │ {{amount | fmt:"#,##0" | right:10}} │
{{/rows}}
{{#preprinted}}
└──────────────────────┴──────┴────────────┘
{{/preprinted}}
                                     合計 {{total | fmt:"#,##0" | right:10}}
```

- **タグだけの行**（`{{#rows}}`、`{{/rows}}`、`{{#preprinted}}` とその閉じ）は印字されない（行を占めない）。
- **幅と揃え**: `width:n`（左揃え、n 桁に収め、はみ出しは切って警告）、`right:n`、`center:n`。桁は東アジア幅（全角 2）で数える。幅を書かないタグは値の長さで行が伸び縮みする（Mustache と同じ）ので、桁を揃える帳票では幅を書く。
- **区間**: 開始と終了が同じ行なら行の中の文字の繰り返し。別の行なら、間の行が項目を 1 つずつ入れる**枠**（`slots`。既定。余った行は空行のまま）。`{{#items | rows}}` は行を増やす繰り返し（用紙長を超えたらページが増え、`title` の行が繰り返される）、`{{#pages | pages}}` はページ単位の複製（10 行の伝票を 25 行分で 3 枚）。規則は Excel と同じ語彙（`docs/forms.ja.md`）。`_index` / `_count` / `_first` / `_last` / `_page` / `_pages` / `_blank` も同じ。
- **用紙の層**: `{{#preprinted}} … {{/preprinted}}` の中の行は用紙の絵（プレ印刷）で、BDF では `background` の役割の層に描かれ、追い刷りでは抜ける。罫線は罫線素片（JIS X 0208 にあり、ESC/P-J の漢字モードで印字できる）で書く。
- **手動の改ページ**: 本文の `\f`（FF）。
- **倍角**: 伝票の表題で使う倍幅・倍高（ESC/P の `ESC W` / `ESC w`、ESC/POS の `GS !`）は `{{#_wide}} … {{/_wide}}` の区間で囲む案。幅は 2 倍で数える。第 1 段階では見送り、論点に残す。

## 4. 出力

### 4.1 BDF

- fixed View。ページの大きさは `columns / cpi` × `lines / lpi` インチ（余白を含む）。長さ可変（`lines: 0`）は 1 枚の高さ可変ページ（`pattern:roll-receipt` と同じ）。
- 各文字をセルに置く。x = `(left + 桁) / cpi` インチ、ベースラインは行の高さの約 0.8。BDF の文字は run の advance で横に伸縮できるので、どのフォントでも格子は正確に揃う（半角 = 1 桁、全角 = 2 桁）。
- フォントは等幅の連鎖（`decision:font-rules`。BIZ UDゴシック、IPAゴシック、Noto Sans Mono CJK JP）。`font` で上書き。埋め込みは使った字だけ。
- MARK は行ごとに PARAGRAPH、`{{#rows}}` の枠は TABLE / CELL に写す（検索・コピー・タグ付き PDF のため）。
- **印字像の Part**: `t: "text"` の Part に、格子の幾何（cpi、lpi、columns、lines、margin、skip）と、ページごとの行の文字列と、どの行が用紙の層かを JSON で持つ。読み手は知らない Part を無視してよい（spec §4）。ESC/P と ESC/POS のテキストモード、プレーンテキストの出力はこの Part から作る。表現（Part か MARK か）は論点。

### 4.2 ESC/P のテキストモード（`print` 側。`print-requirements.ja.md` §3.9）

`ESC @` → `ESC C n`（用紙長、行数）→ `ESC N n`（ミシン目スキップ）→ `ESC 2` / `ESC 0`（LPI）→ `ESC P` / `ESC M` / `ESC g`（CPI）→ `ESC l` / `ESC Q`（余白）→ 行ごとに文字と `CR LF`、ページの終わりに `FF`。日本語は ESC/P-J の漢字モード（`FS &` / `FS .`）で JIS X 0208 の符号に変換（`golang.org/x/text/encoding/japanese`。既に依存にある）。変換できない字（私用領域、JIS X 0213 だけの字）は警告して `〓`。用紙の層の行は、チケットが「プレ印刷」なら空行にする。

### 4.3 ESC/POS のテキストモード

`ESC @` → `ESC t`（文字コード表）→ 漢字モード `FS &` → 行ごとに文字と `LF` → `GS V`（カット）。桁数は 42（80 mm、Font A）/ 48（Font B）など機種のプロファイルから。

### 4.4 プレーンテキスト

`bdf fill invoice.txt -data … out.txt` で印字像そのものを書く（ページは `FF` 区切り、改行は `LF`、文字コードは `-param encoding=`）。固定長の帳票を取り込む旧来のシステムやメールに。

## 5. 検出と CLI

- `form` の判別に「UTF-8 のテキストで、前書き `---` で始まるか、タグ `{{` を含む」を足す（`FormatText = "text"`）。拡張子は `.txt`。`.bdft` の `format` は `"text"`、`src` は元のテキスト。
- `bdf fill slip.txt -data order.json out.bdf`、`bdf fill slip.txt -data order.json out.txt`、`bdf template compile slip.txt slip.bdft -sample sample.json`、`bdf template schema slip.bdft` がそのまま。
- `bdf print out.bdf -printer … -mode text` で ESC/P のテキストモード（print-m5）。

## 6. 構成案

```
form/text/              テンプレートの読み込み（前書き、本文、タグの位置）、幅と揃え、区間の展開（slots / rows / pages）、格子への配置
form/text/layout.go     格子 → BDF の fixed ページ（等幅フォント、advance、MARK、用紙の層）
form/text/image.go      格子 → t: "text" の Part、プレーンテキストの出力
form/form.go            FormatText の判別と Fill / Compile / Fields の分岐
docs/spec.md            parts[].t に text を追記（任意の Part。読み手は無視してよい）
docs/forms.ja.md        「テキストのテンプレート」の節
print/escp（m5）        t: "text" の Part から ESC/P のテキストモード、ESC/POS のテキストモード
```

## 7. 検証

- 単体: 前書きと本文の読み込み、東アジア幅の桁数（半角・全角・結合文字・絵文字は警告）、`width` / `right` / `center`、区間の 3 方式、タグだけの行の扱い、用紙長を超える `rows` のページ分割と `title` の繰り返し、fuzz。
- golden: 納品書（固定ページ、用紙の層つき）、明細が増える一覧（`rows`）、ページ複製（`pages`）、レシート（長さ可変）の 4 つを BDF にして描画を固定。`t: "text"` の Part をプレーンテキストに戻して本文と比較。
- ESC/P: 書き手の出力を、JIS 符号と制御コードの期待値と比較。実機（Epson VP 系）で複写伝票に印字して位置を確かめる（m5）。

## 8. 未決の論点

1. **印字像の表現**: `t: "text"` の Part（推奨）か、MARK の種類を足すか。
2. **幅の既定**: 幅を書かないタグを「値の長さで伸び縮み」（Mustache）にするか、「タグの文字数を幅にする」にするか。推奨は前者＋コンパイル時に桁数の超過を警告。
3. **倍角**: `{{#_wide}}` の区間を第 1 段階に入れるか。
4. **画像・バーコード**: BDF の出力でだけ `{{%}}` を許すか（テキストモードでは落ちる）。推奨は許さない（形式の約束を単純に保つ）。
5. **等幅の既定フォント**: BIZ UDゴシック（Windows 10 以降に同梱、OFL）を第 1 に、IPAゴシック、Noto Sans Mono CJK JP の連鎖でよいか。
6. **区間の既定**: 別の行にまたがる区間の既定を `slots` にする（Excel と同じ）か、`rows` にするか。推奨は Excel と同じ。
