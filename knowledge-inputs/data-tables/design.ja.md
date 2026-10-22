# データ表（描かれない表）の設計

2026-10-22。グラフの系列、CSV・Parquet の入力、Excel のシートといった「表としての値」を、描画とは別に BDF に持たせる。RAG・LLM・データツールが値を型付きで取り出せるようにし、Markdown にはグラフの表を書く。状態: 提案（`requirement:data-tables`）。

## 1. 要件

- BDF は描画済みの表示リストで、描かれない内容を持たない。グラフは図形と文字になり、系列の値は失われる（§3.37 では代替テキストに数値を書いて補っている）。CSV・Parquet・Excel の値も、表示用の文字列になる。
- 持ちたいもの: 文書に付いた **表の列** 。列名と型、行。どこの図の表かが分かること。
- できること: 一覧（デコードなしで列名・型・行数・場所が分かる）、取り出し（CSV・JSON・Parquet）、Markdown での表の書き出し、将来はビューアで図からの表示。
- wire format 1.0 の互換を保つ。知らない読み手は何も壊れない。

## 2. 格納

### Part `table`

- 表 1 つ = Part 1 つ。中身は **Apache Parquet ファイル** そのもの（メディアタイプ `application/vnd.apache.parquet`）。
- BDF の書き手が作る表は、平坦なスカラー列、PLAIN 符号化、UNCOMPRESSED、row group 1 つ、データページは 1 MiB 程度ごと。Part 自体が deflate されるので Parquet 側の圧縮は要らない。
- 入力が Parquet のときは **入力ファイルをそのまま** 格納する（符号化・圧縮・入れ子の型も入力のまま。再符号化しない。§3.2 の「そのまま格納」と同じ考え）。
- Parquet にした理由: 列の型がある、データツール（pandas、DuckDB、Arrow）がそのまま読む、Excel の 100 万行でも小さい、リポジトリに仕様準拠の読み手が既にある（`converter/parquet`）。CSV は型がなく、JSON は大きい。

### manifest の `tables`

任意メンバー（§2「同じ意味を保つ任意の manifest メンバーは format version を変えずに追加できる」）。

```jsonc
"tables": [
  { "id": "chart-1",                   // 文書内で一意。書き手が付ける（chart-1、sheet1、data、table-3 …）
    "title": "Sales by quarter",       // 任意。グラフの題、シート名、Excel のテーブル名
    "part": "<hash>",                  // t: table の Part
    "rows": 4,                         // 行数
    "columns": [                       // 列。一覧のためで、Parquet のスキーマと一致する
      { "name": "Quarter", "type": "string" },
      { "name": "East",    "type": "number" } ],
    "view": "slides", "page": 3 }      // 任意: 表が表す図のある場所（page は 1 から。sheet の表は view だけ）
]
```

`columns[].type` の語彙（Parquet の型との対応。null はどの列でも可）:

| type | Parquet | Go の値 | 用途 |
|---|---|---|---|
| `string` | BYTE_ARRAY (STRING) | string | 文字列、分類 |
| `number` | DOUBLE | float64 | 数値 |
| `integer` | INT64 | int64 | 整数（番号、件数） |
| `boolean` | BOOLEAN | bool | 真偽 |
| `date` | INT32 (DATE) | time.Time（UTC 0:00） | 日付 |
| `datetime` | INT64 (TIMESTAMP, millis, UTC) | time.Time | 日時 |
| `other` | それ以外（入力の Parquet の入れ子・decimal など） | 読み手が JSON や文字列にする | そのまま格納した入力 |

- 上限（§4 の Limits と同じ扱い）: `tables` は 1 文書 65536 まで、`columns` は 1 表 16384 まで。行数は Parquet の footer が言う値を書き、読み手は取り出すときに行数と大きさを制限できる。
- 暗号化した文書では Part として封印され、manifest も封印される。

### 図との結び付け: MARK kind 14 `DATA`

| kind | 名前 | payload | 意味 |
|---|---|---|---|
| 14 | DATA | 表の id | 開いている FIGURE（直前の FIGURE の中）が表すデータ表。区切りに影響しない。知らない読み手は無視する |

FIGURE の直後に置く（`MARK FIGURE "Column chart: Sales"` → `MARK DATA "chart-1"` → 描画 → `MARK END`）。FIGURE の外の DATA は最上位（ページ）のデータとして読み手が扱ってよい。代替テキストは自由文なので、そこに id を混ぜる案は採らない。シートの表は DATA を持たず `tables[].view` で結ぶ。

## 3. 書き手（最初に対応する入力）

| 入力 | 表 | id / title | 列 |
|---|---|---|---|
| グラフ（PowerPoint・Excel・Word の drawingml） | グラフ 1 つに 1 表。全系列が同じ分類を持てば横持ち（`Category, East, West`）、違えば縦持ち（`Series, Category, Value`）。散布図は縦持ち（`Series, X, Y`） | `chart-<n>`（文書の通し番号）/ グラフの題 | 分類は string、値は number |
| CSV・TSV | 全行（表示は既定 N 行でも） | `data` / ファイル名 | 1 行目をヘッダーと判定すれば列名に、しなければ `C1, C2…`。型は推定（整数・数値・日付・真偽・文字列） |
| Parquet | 入力ファイルそのまま | `data` / ファイル名 | footer のスキーマから。入れ子は `other` |
| Excel | シートごとに使用範囲を 1 表。Excel のテーブル（ListObject）はヘッダー付きの別の表 | `sheet<n>` / シート名、テーブルは `table-<n>` / テーブル名 | 列名は `A, B…`（テーブルは見出し行）。値は型付き（数値・文字列・真偽・日付は表示書式で判断・エラーは null） |

- 大きさの既定: 表の Parquet が 64 MiB、行が 1,048,576 を超えるものは持たない（警告）。`-param table=none` で外し、`table=all` で上限なし。
- 後の候補: Word・PowerPoint・HTML・Markdown の表（MARK TABLE から復元できるので優先度低）、タグ付き PDF の表、draw.io / Visio の図形データ。
- グラフの代替テキスト: 表があるときは「種類: 題」だけにし、数値は表に任せる（代替テキストが長くなりすぎないため）。表を持たない設定（`table=none`）では今のまま数値を書く。

## 4. 読み手と API

### Go

```go
// package bdf
type Table struct {
    ID, Title string
    Part      string        // Part のハッシュ
    Rows      int
    Columns   []TableColumn
    View      string
    Page      int
}
type TableColumn struct{ Name, Type string }

func (d *Document) AddTable(t *Table, parquet []byte) error   // Part を登録し t.Part を埋める。id の重複は error
func (d *Document) Tables() []*Table
func (r *Reader) Table(id string) (*Table, []byte, error)     // 項目と Parquet のバイト列（封印は解いて）
func (m *Manifest) Table(id string) *Table

// package encoding/parquet（新設。依存なし。書き手と、平坦なファイルの読み手）
type Column struct { Name string; Type Type }                 // Type: String, Number, Integer, Boolean, Date, DateTime
type Table struct { Columns []Column; Rows [][]any }          // any: string, float64, int64, bool, time.Time, nil
func Write(w io.Writer, t *Table) error
func Schema(data []byte) ([]Column, int64, error)             // footer だけ読む: 列と行数
func Read(data []byte, limits Limits) (*Table, error)         // 平坦なスカラー列（PLAIN / 辞書 / RLE の定義レベル、UNCOMPRESSED・Snappy）

// package converter/parquet（既存の読み手を公開）
func ReadTable(r io.ReaderAt, size int64, opts *ReadOptions) (*parquet.Table, error)  // 任意の Parquet。入れ子は JSON の文字列
```

- `bdf.Document.AddTable` に渡す Parquet は `encoding/parquet.Write` で作る。変換器はそれを使う。
- `bdf` パッケージは `encoding/parquet` に依存してよい（葉のパッケージ）。`converter/parquet` の読み手は符号化・コーデック・入れ子を全部持つので、取り出しの CLI と wasm はそちらを使い、`encoding/parquet.Read` は BDF の書き手が作った平坦なファイル向けに小さく保つ。後で `converter/parquet` の中核を `encoding/parquet` に寄せて 1 つにする。
- 検索用テキスト（`SearchText`）には入れない（描かれた文字と重なる。グラフの値は入らないが、検索の対象ではない）。

### CLI

```
bdf tables <file.bdf | dir>                 表の一覧: id、行数×列数、列名と型、題、場所
bdf table [-format csv|json|parquet] <file.bdf | dir> <id> [out]
                                            表を取り出す（既定: CSV を標準出力に。parquet は Part そのまま）
bdf generate -param table=none|all …        表を持たない / 上限なしで持つ（CSV・Parquet・Excel）
bdf markdown -tables inline|files|none …    図の表を GFM の表で書く（既定 inline: 100 行まで。超える表は <id>.csv を隣に書いてリンク）
```

### wasm（previewonly）と TypeScript

```
bdfConverter.tables(bdf) → [{id, title, rows, columns, view, page}]
bdfConverter.table(bdf, {id, format: "csv" | "json" | "parquet"}) → {data: Uint8Array, rows, columns}
```

- `@bdfkit/core`: `Manifest.tables` の型、`TextNode.data?: string`（FIGURE の DATA）。
- ビューア: 後で、DATA のある図に「Data」のボタン（表を開く・CSV で保存）。最初の段階では入れない。

### Markdown

- DATA のある FIGURE は、代替テキストの行の後に表を GFM で書く（列名が見出し行、`number` は書式なし、`date` は ISO 8601）。
- 100 行を超える表は 100 行まで書き、残りは `<id>.csv` を画像と同じ場所に書いてリンクする（`-tables files` で全部ファイルに、`none` で書かない）。
- シートの表（DATA のない `tables[].view`）は、今のセルからの表をそのまま使う（値は表示文字列）。`-tables files` のときは Parquet の値から CSV も書く。

## 5. spec の変更

- §3.3 の `parts[].t` に `table` を足す。
- §4 に「4.6 データ表（tables）」を足す: 上の JSON、型の語彙、上限、暗号化。
- §7.8 に kind 14 DATA を足す。区切りの規則に「DATA は区切りに影響しない」。
- 適合テスト: fixture の文書にグラフ風の図と表を 1 つ足し、conformance に「tables を持つ文書を開ける」「DATA を無視しても描画が同じ」を足す。

## 6. 互換

- `tables` は任意メンバー、`table` は知らない `t`、DATA は知らない kind。いずれも v1 の読み手が無視すると §2・§3.3・§7.8 で決めてある。
- ビューア（TypeScript）は `tables` を読んでも読まなくても描画は同じ。

## 7. 工数の見積もり

| 項目 | 目安 |
|---|---|
| `encoding/parquet` の書き手と平坦な読み手、テスト（`converter/parquet` のテスト用の書き手 `tenc` / `buildFile` を基に） | 1 日 |
| `bdf` の manifest・AddTable・Reader.Table・上限、spec、fixture と conformance | 半日 |
| CLI `tables` / `table`、wasm、TypeScript の型 | 半日 |
| グラフ（drawingml）の表と DATA、代替テキストの短縮 | 半日 |
| CSV・Parquet の入力 | 半日 |
| Excel のシートと ListObject | 1 日 |
| Markdown の表の書き出し | 半日 |

## 8. 決めてほしいこと

1. 列の型の語彙（`string` / `number` / `integer` / `boolean` / `date` / `datetime` / `other`）でよいか。
2. 表があるグラフの代替テキストから数値を外してよいか（「種類: 題」だけにする）。
3. CSV・Parquet・Excel の既定（64 MiB・1,048,576 行まで持つ。`-param table=none` で外す）でよいか。
4. Excel は「シートごとの使用範囲」に加えて ListObject のテーブルも出すか。
5. `encoding/parquet` を新設し、最初は書き手と平坦な読み手だけにして、`converter/parquet` の読み手は後で寄せる方針でよいか。
6. コマンド名 `bdf tables` / `bdf table` と、Markdown の既定（100 行までインライン、超えたら CSV を隣に）でよいか。
