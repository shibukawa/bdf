# 帳票テンプレート 実装案（form-m1）

2026-10-13 時点の実装案。設計の経緯と決定は `.knowledge/` の `requirement:form-template`、`decision:fill-at-converter-seams`、`taxonomy:form-layout-axes`、`milestone:form-m1` にある。ここは「何をどの順で作るか」を書く。

## 0. 決まっていること

- 出力は `.bdf` だけ。穴埋め済みの Office ファイルは作らない。
- 穴埋めは各変換器の内部モデルの継ぎ目で行う。Excel なら行とセルの構造体を書き換えて、既存の組版を走らせる。
- 第 1 段階（form-m1）は Excel 入力で、固定ページ式、固定行数明細のページ複製、行挿入式の 3 形態。
- 繰越・小計・合計はデータ側で計算する。固定行数明細のデータは呼び出し側がページ単位の配列で渡す。
- 前提として Excel の印刷ビューを作る。穴埋めの出力は print View だけ。
- 後続: PowerPoint（m2）、Word の表（m3）、PDF 書き出しとフォント規則（別トラック）、ラベル・連帳・レシート・フロー文書（未定）。

## 1. 成果物

次のコマンドが動くこと。

```
bdf generate -param views=print invoice.xlsx invoice.bdf          # S1: 印刷ビューだけ
bdf fill invoice.xlsx -data order.json out.bdf                      # S3〜S5: 直接穴埋め
bdf fill invoice.xlsx -data order.json -data lines=lines.csv out.bdf
bdf template compile invoice.xlsx invoice.bdft -sample sample.json  # S6
bdf template schema invoice.bdft
bdf fill invoice.bdft -data order.json out.bdf
```

`out.bdf` は flow View（`print`）だけを持ち、ヘッダー・フッター・本文のレイヤーと本文矩形を持つ。ビューアの紙の表示、連続表示、サムネイル、検索テキストがそのまま効く。

## 2. 構成

```
converter/xlsx/
  print.go         印刷設定の読み取りとページ分割、ページの描画先         (S1)
  headerfooter.go  &L &C &R &P &N &D &T &F &A と書式コードの解釈          (S1)
  template.go      テンプレートモード: タグの走査、行の挿入、参照のずらし、
                   ブロックの複製、式の再計算、値の型付け                  (S3〜S5)
  template_test.go
form/
  tag.go           タグの字句解析と構文（Mustache 系）                   (S2)
  data.go          JSON / YAML / CSV / XLSX の読み込みと名前付き配列、-each (S2)
  schema.go        フィールドのスキーマ（名前・型・繰り返し・書式・必須）    (S2, S6)
  fill.go          テンプレートの種類を見て変換器の穴埋めを呼ぶ            (S3)
  bdft.go          .bdft の読み書き（BDF コンテナ＋tpl / src Part）         (S6)
cmd/bdf/
  fill.go          bdf fill                                             (S3)
  template.go      bdf template compile | schema | fonts                 (S6)
bdf/ (ルート)
  manifest.go      PartTemplate = "tpl"、PartSource = "src"              (S6)
docs/spec.md       parts[].t に tpl / src を追記（任意の Part。読み手は無視してよい）(S6)
test/xlsx/gen.py   印刷設定つきブックとテンプレートブックの生成            (S1, S3〜S5)
testdata/form/     穴埋めの golden                                       (S3〜S5)
```

Excel の穴埋めを `converter/xlsx` の中に置く理由は、ワークシートのモデル（`row`、`cell`、`cellRange`、`sheetCtx`）と数式の評価器（`parseFormula`、`fenv`）と表示形式（`parseNumFormat`）がすべて非公開で、公開すると API の表面が大きくなるため。`form` パッケージは 3 形式に共通の部分だけを持つ。

## 3. データの契約

- テンプレートから見えるのは JSON の木 1 つ。`-data file.json|yaml` が根、`-data name=file.csv|xlsx` は根の `name` に配列を足す。CSV は 1 行目が見出し、ブックはシート名ごとの配列。
- 値の型は JSON のまま使う。数値は数値（セルの表示形式と右揃えが効く）、`"2026-10-13"` や `"2026-10-13T09:00:00+09:00"` の文字列は日付として読んでシリアル値にする、真偽はチェックボックス、null は空。
- 固定行数明細のデータはページ単位の配列。

```json
{
  "customer": "株式会社例",
  "pages": [
    {"rows": [{"name": "A", "qty": 2, "amount": 1000}, "..."], "carry_in": 0,     "carry_out": 12000},
    {"rows": ["..."],                                           "carry_in": 12000, "carry_out": 25000}
  ],
  "total": 25000
}
```

- エンジンは配列を切らない。1 ページの `rows` が明細域の行数を超えたらエラー。`-chunk lines=10` のような CLI の補助は後で足してよい。
- 画像はファイル名か base64。ファイルは `Options.Files` の FS から読む。ネットワークは読まない。
- `-each rows.csv` はレコードごとに 1 文書。出力先に `{{code}}` を含めるとファイルごと、含めなければ 1 つの BDF に全部。

## 4. タグ（m1 の範囲）

| 書き方 | 意味 |
|---|---|
| `{{name}}`、`{{a.b.c}}` | 値。セルがタグ 1 つだけなら値の型をセルに移す |
| `{{#items}} … {{/items}}` | 区間。配列なら繰り返し、真偽や null なら条件 |
| `{{^paid}} … {{/paid}}` | 否定の区間 |
| `{{total \| fmt:"#,##0"}}` | 書式。Excel の表示形式コード。日付は `yyyy年m月d日`、和暦は `ggge年m月d日` |
| `{{%photo}}` | 画像。セル（結合範囲）に contain で置く |
| `{{_page}}`、`{{_pages}}`、`{{_first}}`、`{{_last}}`、`{{_index}}`、`{{_count}}`、`{{_blank}}`、`{{_blank_first}}` | 文脈の変数 |

- 全角の `｛｛ ｝｝` も受け付ける。
- Excel での区間の範囲は、開始と終了のタグの位置で決まる。同じセル → セルの中の文字列。同じ行の別のセル → その行を行テンプレートにして挿入（行挿入式）。別の行 → その間の行の組が明細域で、行数が固定（固定行数式）。明細域を含む印刷範囲は、外側の区間（`{{#pages}}`）があればページ要素ごとに複製される。
- 区間はこの 2 段（`pages` の中に `rows`）までを m1 の対象にする。それより深い入れ子はコンパイルでエラー。
- タグの構文エラーと、行や印刷範囲をまたぐ区間は、セル参照つきのエラーにする。

## 5. Excel の穴埋めパイプライン

1. **読む**。`readWorksheet` を今のまま使う。テンプレートモードでは、セルに式の文字列（`f` 要素）を残す。今は捨てている。
2. **タグを集める**。`cellStr` のセルの文字列を走査する。共有文字列のセルはテンプレートモードでは行ごとに独立した文字列として扱い、同じ共有文字列を参照する他のセルに影響させない。
3. **テンプレートの構造を作る**。値のタグ、区間の開始と終了の対、区間の範囲（セル・行・明細域）、印刷範囲との対応、画像のタグ。ここでエラーを出す。`.bdft` のスキーマはこの構造から作る。
4. **穴埋め**。
   - 値のタグ: 置換して型を付ける。書式フィルタは `parseNumFormat` と `formatNumber` で文字列にする。
   - 行挿入: 行テンプレートを N−1 行複製して下へずらす。ずらす対象は、結合セル、条件付き書式の `sqref`、ハイパーリンク、定義名（印刷範囲と印刷タイトル）、手動改ページ、テーブルの範囲、全シートの式の参照。`SUM(D10:D12)` は挿入分だけ広がる（Excel の「範囲の中に挿入したら伸びる」規則）。
   - 固定行数ブロック: 印刷範囲の行をページ要素ごとに複製し、下に並べる。各複製の明細域に `rows` を入れ、余った行はタグを消して静的部分（罫線）を残す。複製の間に手動改ページを入れ、印刷範囲を複製ぶんだけ伸ばす。ページ番号と総ページ数は複製の数から入れる。
   - 画像: セルの矩形に `drawingml` の経路を使わず、`canvas` に直接 IMAGE を置く（図形の読み込みを経由しない方が単純）。
5. **再計算**。値が変わったセルに依存する式を、参照の依存関係の順に `fenv.eval` で評価する。循環は警告して止める。未対応の関数は警告して空。評価できる関数は条件付き書式のもの（四則、SUM、AVERAGE、MIN、MAX、COUNT、COUNTIF、IF、IFERROR、ROUND 類、LEN、LEFT、RIGHT、MID、TEXT の一部、YEAR、MONTH、DAY など）。
6. **組む**。`worksheet()` の流れを `views=print` で通す。印刷設定の読み取りとページ分割は S1 の機能。
7. **書き出す**。print View だけ。sheet View は作らない。

## 6. 印刷ビュー（S1）

- **読む**。`workbook.xml` の `definedName` から `_xlnm.Print_Area` と `_xlnm.Print_Titles`（`localSheetId` でシートに結ぶ）。ワークシートの `pageSetup`、`pageMargins`、`printOptions`、`rowBreaks`、`colBreaks`、`headerFooter`、`sheetPr/pageSetUpPr`。`sheet.go` の `element()` に case を足し、`worksheet` に `print` 構造体を持たせる。
- **ページに切る**。用紙（`paperSize` の表、または `paperWidth` / `paperHeight`）から余白を引いた印刷可能域に、拡大率をかけた列と行を順に詰める。手動改ページは必ず切り、1 行が 1 ページより高ければ途中で切る。非表示は入れない。「ページに合わせる」は横 n 縦 m に収まる最大の拡大率を 10% まで 1% 刻みで探す（`fitToWidth` / `fitToHeight` の 0 は自動）。印刷範囲が複数ならその順。なければ使用範囲。ページの順序は `downThenOver` / `overThenDown`。
- **描く**。ページごとに `tileCv` と同じ型の描画先を作り、Object の先頭に `Transform(s, 0, 0, s, tx, ty)` を置いてシート座標のまま描く。`tilesIn(x0, y0, x1, y1)` が「矩形に触れる描画先」を返すところを、sheet View ではタイル、print View ではページを返すようにする（`sheetCtx` に描画先の選び方を持たせる）。印刷タイトルの行と列は、タイトルを含まないページで、ページの上と左の領域に同じ方法で描く。m1 ではページごとに描き直し、子 Object の共有は後の最適化。枠線と行列番号は `printOptions` に従ってページに描く。
- **ヘッダーとフッター**。`&L` `&C` `&R` の 3 区画、`&P` `&N` `&D` `&T` `&F` `&A`、`&"フォント,スタイル"`、`&サイズ`、`&B` `&I` `&U`、`&&`。奇数・偶数・先頭ページの別。`&D` と `&T` は `Options` で与えた日時だけを使い、なければ空（決定性のため）。
- **View**。`flow` View `print`。各ページは用紙の大きさ、`body` は印刷可能域、レイヤーは `background`（なし）、`header`、`footer`、`body`。テキスト索引も作る。
- **オプション**。`Options.Views`（`sheet` / `print` / `both`、既定 `sheet`）と `-param views=`。`thumbnail` は print があればその 1 ページ目の全体にする（後で）。
- **コンパイルの報告**（S6）。印刷範囲がない、縦 1 ページ固定のページ合わせ、自動改ページが行の境の 3 pt 以内、は警告にする。

## 7. スライスと受け入れ条件

| スライス | 作るもの | 受け入れ |
|---|---|---|
| S1 印刷ビュー | `print.go`、`headerfooter.go`、定義名、`views` パラメータ、`gen.py` に印刷設定つきブック | 印刷範囲・タイトル・ページ設定・余白・手動改ページ・ヘッダーとフッターのブックが Excel と同じ規則でページに切れ、golden が固定される。sheet View の golden は不変 |
| S2 タグとデータ | `form/tag.go`、`form/data.go`、`form/schema.go` | 構文の単体テスト。全角波括弧。開閉の不一致とまたぎのエラーにセル参照 |
| S3 固定ページ | `template.go` の走査と値の置換、`form/fill.go`、`bdf fill` | 1 ページのテンプレートと JSON から型と書式の効いた .bdf。欠けたフィールドは警告して空、`-strict` でエラー |
| S4 行挿入 | 行の複製と参照のずらし、式の再計算 | ずらしと再計算の単体テスト。見出し行の繰り返しつきで複数ページに割れる golden |
| S5 固定行数 | ブロックの複製、ページ変数、余り行の処理 | `pages` 配列の要素ごとに 1 ページ、`_page` / `_pages`、余った行の罫線、繰越額の配置の golden |
| S6 .bdft | `form/bdft.go`、`bdf template`、Part 型、spec の追記、読み手の確認 | compile と fill の往復で同じ .bdf。Go / TS / libbdf が `tpl` と `src` を無視して開ける。プレビューがビューアで見える |

S1 と S2 は独立。S3 は S1 と S2 に依存。S4 と S5 は S3 に依存し、互いに独立。S6 は S3 の後ならいつでも。

## 8. テストとフィクスチャ

- `test/xlsx/gen.py` に、印刷設定つきのブック（`print.xlsx`: 複数の印刷範囲、印刷タイトル、A4 縦横、拡大率とページ合わせ、手動改ページ、ヘッダーとフッターのコード）と、テンプレートのブック（`form-fixed.xlsx`、`form-rows.xlsx`、`form-pages.xlsx`）を足す。openpyxl は定義名、`print_title_rows`、`page_setup`、`row_breaks`、`oddHeader` を書ける。
- 穴埋めの golden は `testdata/form/` に、テンプレートとデータと出力の組で置く。描画の比較は既存の golden の仕組み（ブラウザと imagebdf）をそのまま使う。
- ずらしと再計算は表駆動の単体テスト。`SUM(D10:D12)` に 3 行挿入して `SUM(D10:D15)`、結合 `B10:C10` が `B13:C13` へ、条件付き書式の `sqref`、別シートからの参照。
- 開発中の見た目の確認は LibreOffice headless で同じブックを PDF にして比べる。CI には入れない（ページ分割の丸めが違う）。
- fuzz: タグの字句解析と、テンプレートの構造の組み立てに `FuzzTemplate` を足す。テンプレートは信頼できない入力。

## 9. 暫定の選択

- 欠けたフィールドは警告して空。`-strict` でエラー。帳票の現場はエラー寄りかもしれないので運用で見直す。
- 画像はセルに contain で置くだけ。既存の図形の差し替えは後。
- 式の文字列はテンプレートモードのときだけ cell に残す。通常の変換のメモリは増やさない。
- 表示形式のコードは m1 では xlsx の内部のまま使う。Word と PowerPoint で要るときに `converter/internal/numfmt` に切り出す。
- 印刷タイトルは m1 ではページごとに描き直す。子 Object の共有は最適化として後。
- `.bdft` は BDF コンテナ。`tpl`（テンプレート定義の JSON）と `src`（元の .xlsx）を Part に持ち、ページはサンプルデータで穴埋めしたプレビュー。Go の読み手は Part の型を見ずに開ける。TS の `@bdfkit/core` と libbdf は S6 で確かめ、無視しなければ spec に型を足す。
- フォントは m1 では既存の `-font-dir` と `FontFS` で組む。`.bdft` へのフォントの固定と規則ファイルは別トラック。

## 10. 大きさの目安

| 部分 | 行数の見込み |
|---|---|
| 印刷ビュー（読み取り、ページ分割、ページの描画先、ヘッダーとフッター） | 1,200〜1,800 |
| テンプレートモード（走査、挿入、ずらし、複製、再計算、画像） | 1,500〜2,200 |
| form パッケージ（タグ、データ、スキーマ、fill、bdft） | 1,200〜1,600 |
| CLI とドキュメント | 300〜500 |
| テストとフィクスチャの生成 | 1,000〜1,500 |

先に S1 を終えると、帳票と無関係に「Excel の印刷の見た目で開ける」という価値が出るので、ここで一度 PR を切るのがよい。
