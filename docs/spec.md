# BDF (Browser Document Format) 仕様ドラフト v0.1

> 状態: ドラフト。名称 "BDF" は仮称（X11 のビットマップフォント形式 BDF と衝突するため、公開時には要再考）。

## 1. 目的と非目的

**目的**

- Office 系ファイル（PowerPoint / Excel / Word 等）や図（draw.io）を**サーバー側で変換**し、**ブラウザ（Worker）で Canvas 2D にそのまま描画**できる中間フォーマット。
- 命令セットは `CanvasRenderingContext2D` の API に 1:1 に近い形で対応させ、デコーダが「バイト列を読んで ctx のメソッドを呼ぶだけ」になるようにする。
- 圧縮・展開は `CompressionStream` / `DecompressionStream`（`deflate-raw`）に任せ、独自実装を持たない。
- **1ファイル**（配布・キャッシュ向け）と**分割ファイル群**（オンデマンド取得・CDN 向け）を同じ論理構造で表現できる。
- スライドマスターやページ背景など**同一内容の繰り返しは共有オブジェクトとして1回だけ格納**し、サイズを抑える。
- 3種類のレイアウトモデルを表現できる。
  - 固定サイズページの列（PowerPoint、PDF）
  - 無限サイズの平面（Excel のシート）
  - ページ分割されつつ「縦一枚」としても閲覧できる文書（Word）。ページを持たずにレイアウトした 1 枚の長い面も持てる
- パスワードで保護された入力からは、同じパスワードで暗号化した BDF を作れる。暗号化しても Part 単位の取得はそのまま使える（§3.5）。

**非目的**

- 編集可能な文書モデル（DOM のような意味構造）を持つこと。BDF は「描画済み表示リスト」であり、再レイアウトはしない。ただし読み上げに必要な最小限の構造（見出し・リスト・表・図の代替テキスト・言語、§7.8）は描画に影響しない `MARK` として持つ。
- 印刷用の高度な色管理（CMYK、ICC）。sRGB のみ。
- 署名。必要なら配布層で行う。暗号化（§3.5）は保護された入力の保護を保ち続けるためのもので、配布層のアクセス制御の代わりではない。

## 2. 用語

| 用語 | 意味 |
|---|---|
| Part | コンテナに格納される1つのバイト列。内容のハッシュで識別される。 |
| Object | 描画命令列を含む Part。ページ本体・マスター・シートのタイルなど、描けるものはすべて Object。 |
| View | 閲覧単位。スライド集、シート、文書、図のページなど。1 BDF に複数の View を持てる（Excel の複数シート、draw.io の複数ページ等）。 |
| Page | 固定サイズの矩形。`fixed` / `flow` View を構成する。 |
| Tile | `sheet` View を格子状に分割した矩形。1 Tile = 1 Object。 |
| Resource | Object から参照される Font / Image / Path / Paint。 |

## 3. 全体構造

```
BDF
├── manifest (JSON)               … View 一覧、ページ寸法、Part 一覧
└── parts/<hash>                  … 内容アドレスの Part 群
    ├── Object（描画命令列）
    ├── Font（WOFF2 等、そのままのバイト列）
    ├── Image（PNG/JPEG/WebP/AVIF、そのままのバイト列）
    ├── Path collection（パス群のバイナリ）
    └── Index（大きなページ表・タイル表）
```

### 3.1 内容アドレス（content addressing）

- Part 名は **SHA-256 の先頭 16 バイトを小文字 hex にした 32 文字**。
- 同一内容の Part は自動的に 1 つに畳まれる。マスタースライド、ヘッダー/フッター、繰り返し画像などはエンコーダ側で「同じバイト列を作る」だけで共有される。
- 読み手はハッシュ検証を**してもよい**（`crypto.subtle.digest`）が必須ではない。

### 3.2 圧縮

- Part ごとに独立して圧縮する。方式は `identity`（無圧縮）または `deflate-raw`。
- 読み手は `new DecompressionStream("deflate-raw")` で展開する。
- 独立圧縮にするのは Range 取得・部分取得・タイル単位キャッシュを可能にするため。Part をまたぐ冗長性は 3.1 の共有で吸収する。
- 目安: 展開後 512 バイト未満、または画像・WOFF2 のような既圧縮データは `identity`。

### 3.3 1ファイル形式（single）

先頭にディレクトリを置く。ストリーミング読み（`fetch().body` を先頭から読み進めながら描画開始）と Range 取得の両方に対応する。

```
offset 0   : magic          u8[4]   "bdf\0"（62 64 66 00）
offset 4   : version        u16     フォーマットバージョン（1）
offset 6   : flags          u16     bit 0: 暗号化（§3.5）。他のビットは予約（0）
offset 8   : manifestOff    u64     manifest の先頭オフセット（通常 32）
offset 16  : manifestLen    u64     manifest の圧縮後長
offset 24  : manifestEnc    u8      0=identity, 1=deflate-raw
offset 25  : reserved       u8[7]
offset 32  : manifest
offset ... : parts 領域（manifest 直後から始まり、manifest.parts に記載された順に連続配置）
```

- すべてリトルエンディアン。
- マジックは小文字の `bdf` に NUL を 1 バイト続けたもの。NUL があるのでテキストファイルと取り違えない。バージョンは `version` で表し、マジックには含めない。
- `manifest.parts[i]` は `off`（**parts 領域の先頭**、すなわち `manifestOff + manifestLen` からのオフセット）と `len` を持つ。manifest の内容が自身の長さに依存しないようにするため。
- **推奨配置順**: manifest → Font → 共有 Object（マスター等）→ View の先頭ページから順。先頭から読むだけで 1 ページ目が描けるようにする。
- 書き手は全 Part を確定してから書き出す（サーバー変換なので問題ない）。

### 3.4 分割形式（split）

```
<base>/manifest.json
<base>/parts/<hash>
```

- `manifest.json` と Part の内容は single と同一。`off` は無視される（存在しなくてよい）。
- 各 Part ファイルは `enc` に従って**圧縮済みのまま置く**。HTTP の `Content-Encoding` に依存しないので、S3 等のオブジェクトストレージや CDN にそのまま置ける。
- 読み手は `fetch(base + "/parts/" + hash)` で必要な Part だけを取得する。
- single ↔ split は Part の再圧縮なしに相互変換できる（ディレクトリの組み替えだけ）。

### 3.5 暗号化

パスワードで保護された入力（Office の読み取りパスワード、PDF のユーザーパスワード）から作った文書は、同じパスワードで暗号化する。Part を 1 つずつ封印するので、Range 取得・split 形式・CDN への配置はそのまま使える。

```
外側の manifest（平文）  … 暗号方式、鍵スロット、封印された Part の一覧
封印された Part
├── manifest             … 文書の manifest（§4。views、meta、parts）
└── 各 Part              … Part の格納バイト列（圧縮済み）を暗号化したもの
```

**外側の manifest**

```jsonc
{
  "bdf": 1,
  "encryption": {
    "cipher": "A256GCM",
    "keys": [
      { "type": "password", "kdf": "PBKDF2-SHA256", "iter": 600000,
        "salt": "<base64, 16 バイト>",
        "key": "<base64: AES-KW でラップしたコンテンツ鍵, 40 バイト>" }
    ],
    "manifest": { "part": "<hash>", "enc": "deflate-raw" }   // 封印された manifest とその符号化
  },
  "parts": [
    { "h": "<hash>", "t": "sealed", "enc": "identity", "len": 1474, "size": 1474, "off": 0 },
    …
  ]
}
```

- views・meta は持たない。題名などの Dublin Core も封印された manifest の中にある。
- single 形式ではヘッダの `flags` の bit 0 を立てる（§3.3）。bit 0 と `encryption` の有無が食い違うファイルは不正。
- 封印された Part の名前は、**格納バイト列（暗号文）** の SHA-256 の先頭 16 バイト。平文のハッシュを外に出すと、よく使われるフォントのサブセットや画像が含まれているかを外から照合できてしまう。格納バイト列のハッシュなので、鍵がなくても名前を検証でき、single ↔ split も鍵なしで組み替えられる（§3.4）。
- 並び順は、封印された manifest が先頭で、その後は §3.3 の推奨順。

**鍵**

- コンテンツ鍵: 文書ごとにランダムな 256 ビットの鍵。manifest とすべての Part をこの鍵で封印する。
- 鍵スロット（`keys`）: コンテンツ鍵を AES-KW（RFC 3394）でラップしたもの。`password` スロットのラップ鍵は PBKDF2-HMAC-SHA-256 で導く。入力はパスワードを **NFC に正規化して UTF-8 にしたもの**、salt は 16 バイト、反復回数は `iter`（書き手の既定は 600,000）。
- 読み手はスロットを順に試し、アンラップの完全性検査を通ったものを使う。どれも通らなければパスワード違い。
- 読み手は `iter` が 1〜10,000,000 の範囲外なら拒否する。開くだけで長く計算させられるのを防ぐため。
- スロットを複数持てるのは、パスワードの変更（スロットの差し替えだけで済み、Part は暗号化し直さない）や、回復用の鍵などを後から足すため。

**封印**

- 封印した Part の格納バイト列は `nonce（12 バイト）‖ 暗号文 ‖ tag（16 バイト）`。AES-256-GCM で、nonce はランダム。
- 追加認証データ（AAD）は、Part なら**平文の Part 名（16 バイト）**、manifest なら ASCII の `manifest`。Part の差し替えや manifest との取り違えは復号で失敗する。
- 圧縮してから暗号化する。封印を解いた中身は、封印された manifest の `parts[].enc` に従って展開する。

**封印された manifest**

§4 の manifest と同じ形で、`parts[]` の各エントリに `sealed`（その Part を封印した外側の Part の名前）を持つ。`off` は使わず、外側のエントリの `off`/`len` で取得する。

読み手の手順は、外側の manifest を読む → パスワードからラップ鍵を導いてコンテンツ鍵をアンラップ → `encryption.manifest.part` を復号・展開 → 以降は Part 名で引き、`sealed` が指す外側の Part を取得して復号・展開、となる。必要なのは WebCrypto の `PBKDF2`・`AES-KW`・`AES-GCM` だけ。コンテンツ鍵は抽出できない `CryptoKey` として Worker の中に置く。

**守るもの・守らないもの**

- 守る: 配布層（ストレージ、CDN、キャッシュ）からの漏洩。パスワードを知らない人にわかるのは、Part の数と大きさだけ。
- 守らない: パスワードの総当たり。ファイルを手に入れればオフラインで試せるのは元の Office ファイルや PDF と同じ。PBKDF2-SHA-256 の 60 万回は、元のファイル（Office の Agile 暗号化は SHA-512 を 10 万回）より弱い経路を作らないために選んだ。
- 変換はパスワードと平文を扱う。変換の後でパスワードを残さないのは運用側の責務。

## 4. Manifest

JSON。読みやすさとツールでの扱いやすさを優先する。巨大になる部分（ページ表・タイル表）は Index Part に逃がす。

```jsonc
{
  "bdf": 1,                     // フォーマットバージョン
  "opset": 1,                   // 命令セットバージョン（§7）
  "unit": "pt",                 // 1 unit = 1/72 inch。座標はすべてこの単位
  "meta": {                     // §4.3
    "dc": { "title": "…", "creator": ["…", "…"], "language": "ja", "created": "2026-09-25T07:21:16+09:00" },
    "source": "pptx", "generator": "bdf-go/0.1" },

  "views": [
    { "id": "slides", "kind": "fixed", "title": "スライド",
      "pages": [ { "w": 960, "h": 540, "layers": [
          { "role": "master", "obj": "3f9a…" },
          { "role": "body",   "obj": "b1c2…" } ] },
        … ] },

    { "id": "doc", "kind": "flow", "title": "本文",
      "pages": [ { "w": 595.3, "h": 841.9,
                   "body": { "x": 72, "y": 72, "w": 451.3, "h": 697.9 },
                   "layers": [
          { "role": "header", "obj": "…" },
          { "role": "body",   "obj": "…" },
          { "role": "footer", "obj": "…" } ] }, … ],
      "continuous": { "gap": 24 },
      "textIndex": "<hash>" },                       // 任意: テキスト索引 Part（§7.9）

    { "id": "scroll", "kind": "scroll", "title": "本文",
      "pages": [ { "w": 523.3, "h": 1014.2, "layers": [ { "role": "body", "obj": "…" } ] },
                 { "w": 523.3, "h": 988.5,  "layers": [ { "role": "body", "obj": "…" } ] }, … ] },

    { "id": "sheet1", "kind": "sheet", "title": "Sheet1",
      "tile": 2048,
      "cols": [[26, 64], [1, 120], [16358, 64]],   // [count, size] のランレングス
      "rows": [[1048576, 20]],
      "freeze": { "cols": 1, "rows": 1 },
      "gridlines": true,
      "tiles": { "0,0": "…", "1,0": "…" }           // 小さければインライン
      // "tilesRef": "<hash>"                        // 大きければ Index Part
    }
  ],

  "parts": [
    { "h": "3f9a…", "t": "obj",  "enc": "deflate-raw", "len": 1234, "size": 5678, "off": 0 },
    { "h": "c0ff…", "t": "font", "enc": "identity",    "len": 40210, "size": 40210, "off": 1234 },
    …
  ]
}
```

`parts[].t` の値: `obj` / `font` / `img` / `path` / `idx`。暗号化した文書の外側の manifest では `sealed`（§3.5）。

### 4.1 View の種類

**`fixed`** — 固定サイズページの列。ページごとに `w`,`h` が異なってよい。ビューアはページを縦または横に並べる、あるいは 1 ページずつ表示する。

**`flow`** — `fixed` に加え、各ページが `body`（余白を除いた本文矩形）を持ち、レイヤーに役割（`role`）が付く。ビューアは 2 つの表示モードを提供できる。

- ページモード: 紙の形で全レイヤーを描く。
- 連続モード: `role: body` のレイヤーだけを `body` 矩形で切り出し、`continuous.gap` を挟んで縦に積む。ヘッダー・フッター・背景は描かない。

連続モードは**再レイアウトではなく本文矩形の積み上げ**である（BDF は表示リストであり、テキストの折り返し位置は変えられない）。「1 枚の長いページとして再レイアウトしたもの」も欲しい場合、変換側が別 View（`kind: scroll`）を追加で生成する。共有 Object により画像やフォントは二重に格納されない。

**`scroll`** — ページを持たない縦に長い 1 枚の面（Word の下書き・Web レイアウト表示のように、用紙の高さを無限にしてレイアウトした文書）。`pages` はこの面を上から順に水平に切った**帯**（strip）で、形は `fixed` のページと同じ（`w`、`h`、`layers`）。

- ビューアは帯を隙間なく縦に積み、それぞれ帯の矩形で切り抜いて描く（`flow` の連続モードで `continuous.gap` が 0、`body` が帯全体の場合と同じ）。帯はページではないので、ページとして 1 枚ずつ見せたりページ番号を示したりしない。
- 帯に切るのは、Canvas の大きさの上限を超えないため、また見えている範囲の Part だけを読み込むためである。書き手は帯の境界を行の間に置き（テキストの行を横切らない）、境界をまたぐ図・塗り・罫線は両方の帯に描く。帯の高さはそろっていなくてよい。帯の幅はすべて同じ。
- 開いている構造（§7.8）は帯ごとに閉じるので、帯をまたぐ表やリストは次の帯で開き直す。LINK の `#page=N` は N 番目の帯を指す。

**`sheet`** — 原点 (0,0) から右下に伸びる無限平面。

- `cols` / `rows` はランレングス `[count, size]` の配列。ビューアは行・列ヘッダー、グリッド線、固定ペイン（`freeze`）をこの情報から自前で描く。
- 内容は `tile` unit 四方の Tile に分割され、各 Tile は 1 Object。Tile 内の座標は **Tile 原点を基準**にする（f32 の精度限界を超えないため）。
- 複数 Tile にまたがるオブジェクト（結合セル、画像、グラフ）は共有 Object にして、触れる各 Tile から `USE` で参照する。ビューアは Tile 矩形でクリップして描くので重複描画は正しく処理される。小さな描画（セルの塗りや罫線、隣のセルにはみ出したテキスト）は共有せず、触れる各 Tile にその Tile の原点で描いてもよい。
- テキストの抽出（テキスト索引、テキスト層、検索）では、run はアンカー（`FILL_TEXT` などの x, y を変換したもの）が含まれる Tile に属し、他の Tile が重ねて描いた同じ run は除く。Tile の範囲は左上の辺を含み、右下の辺を含まない。
- 空 Tile は表に載せない。
- `tiles` の Index Part 形式（`t: idx`）: `[u32 tx, u32 ty, u8[16] hash]` を (ty, tx) 昇順に並べた固定長レコード列。二分探索で引く。

**複数の View** — `views` の順が表示の順。View が 2 つ以上あるとき、ビューアは表計算ソフトのシート見出しのように View をタブで並べて切り替える（Excel のシート、draw.io のページ）。View の `id` は文書内で一意な任意の文字列で、LINK の `#view=` で参照される（§7.7）。

### 4.2 レイヤーの役割（role）

`background` / `master` / `header` / `footer` / `body` / `notes` / `annotation`。ビューアは役割ごとの表示切替に使える。未知の role は `body` として扱う。

### 4.3 メタデータ（meta）

| キー | 意味 |
|---|---|
| `dc` | 文書そのものの記述。Dublin Core（下記） |
| `source` | 変換元の形式（`pdf` / `pptx` / `xlsx` / `csv` / `vsdx` / `vdx` / `drawio` / `dxf` / `emf` / `wmf` / `tiff` / `fixture` …）。Dublin Core の `source` とは別物 |
| `generator` | 書き出したソフトウェア（例 `bdf-go/0.1`） |

`meta.dc` は [Dublin Core Metadata Element Set 1.1](https://www.dublincore.org/specifications/dublin-core/dces/) の 15 要素に、[DCMI Metadata Terms](https://www.dublincore.org/specifications/dublin-core/dcmi-terms/) の `created` と `modified` を加えたもの。キーは要素名（名前空間接頭辞なし）。

```jsonc
"dc": {
  "title": "四半期報告",
  "creator": ["山田太郎", "佐藤花子"],
  "subject": ["売上", "予算"],
  "language": "ja",
  "created": "2026-09-25T07:21:16+09:00"
}
```

- どの要素も省略でき、繰り返せる。値は文字列、または文字列の配列。書き手は値が 1 つなら文字列、2 つ以上なら配列で書く。読み手はどちらも受け付ける。空の配列は要素がないのと同じ。
- 未知のキーは無視する。

| 要素 | 内容 | 推奨する書き方 |
|---|---|---|
| `title` | 題名 | |
| `creator` | 作成者 | 人・組織ごとに 1 値 |
| `subject` | 主題・キーワード | キーワードごとに 1 値 |
| `description` | 概要 | |
| `publisher` | 発行者 | |
| `contributor` | 寄与者 | |
| `date` | 何らかの日付（発行日など） | W3CDTF（ISO 8601 の部分集合）。`2026`、`2026-09`、`2026-09-25`、`2026-09-25T07:21:16+09:00` |
| `type` | 種類 | DCMI Type Vocabulary（`Text`、`Image` など） |
| `format` | 形式 | MIME タイプ |
| `identifier` | 識別子 | URI、ISBN、DOI など |
| `source` | 派生元の資源 | |
| `language` | 言語 | BCP 47（`ja`、`en-US`）。最初の値が文書のテキストの既定の言語（§7.8 の LANG で空の run の言語） |
| `relation` | 関連する資源 | |
| `coverage` | 範囲（場所・期間） | |
| `rights` | 権利 | |
| `created` | 作成日時 | `date` と同じ書き方 |
| `modified` | 更新日時 | `date` と同じ書き方 |

変換器は入力文書のメタデータを次のように写す。

| 要素 | PDF（文書情報辞書） | PowerPoint、Excel、Word（コアプロパティ） | draw.io | TIFF（先頭のページのタグ） |
|---|---|---|---|---|
| `title` | `Title` | `dc:title` | – | `DocumentName` |
| `creator` | `Author` | `dc:creator` | – | `Artist`（`;` で分割） |
| `subject` | `Keywords`（`,` `;` `、` などで分割） | `dc:subject`、`cp:keywords`（同様に分割） | – | – |
| `description` | `Subject` | `dc:description` | – | `ImageDescription` |
| `identifier` | – | `dc:identifier` | – | – |
| `language` | – | `dc:language`（なければ PowerPoint は既定のテキストスタイルの、Word は既定の run の言語。Word は本文の多くが和文なら東アジアの言語） | – | – |
| `rights` | – | – | – | `Copyright` |
| `created` | `CreationDate`（W3CDTF に変換） | `dcterms:created` | – | – |
| `modified` | `ModDate`（W3CDTF に変換） | `dcterms:modified` | `mxfile` の `modified` | `DateTime`（W3CDTF に変換） |

PDF と TIFF の対応は、XMP が文書情報辞書と TIFF のタグを写す方法に合わせている（TIFF の `DocumentName` は XMP にないので題名にした）。`bdf generate` の `-dc 要素名=値`（繰り返し可）で要素を上書きでき、`-dc 要素名=` でその要素を消せる。

## 5. Object Part

描画命令列と、それが参照するリソース表を持つ自己完結した Part。ページ、マスター、Tile、部品はすべてこの形式。

```
u8[4]   "BOBJ"
u16     opset バージョン
u16     flags（予約）
f32×4   bbox (x, y, w, h)  … この Object が描く範囲。キャッシュ用。原点は Object 自身のローカル座標
varuint nStrings ; UTF-8 文字列 [ varuint len, bytes ]×n
varuint nPaths   ; パス [ u8 kind(0=inline,1=ext) , inline: パスデータ(§6.3) / ext: u8[16] hash, varuint index ]×n
varuint nPaints  ; グラデーション・パターン定義(§6.4)×n
varuint nFonts   ; フォント参照(§6.1)×n
varuint nImages  ; 画像参照 [ u8[16] hash ]×n
varuint nObjects ; 子 Object 参照 [ u8[16] hash ]×n
varuint opsLen   ; 命令列のバイト長
u8[opsLen]       ; 命令列（§7）
```

- 命令内のリソース参照（`pathRef`, `fontRef` …）は上記各表のローカルインデックス。
- 外部参照のハッシュは Manifest の `parts` に存在しなければならない。
- 座標は Object のローカル座標。`USE` する側が変換行列で配置する。

## 6. リソース

### 6.1 Font

```
u8      kind        0=embedded, 1=system
u8[16]  hash        embedded のときのみ（Font Part。WOFF2 推奨、TTF/OTF も可）
string  family      CSS フォントファミリー文字列（フォールバック列を含む）。embedded では空でよく、空でなければ埋め込みフォントの後ろにフォールバックとして付ける
u16     weight      100–900
u8      style       0=normal, 1=italic, 2=oblique
```

- Worker 内では `new FontFace("bdf-" + hash, buffer)` を生成し `self.fonts.add()` する。`ctx.font` にはこの合成名を使う。
- サブセット化はエンコーダの責務。ライセンス上埋め込めないフォントは `system` にするか、テキストをグリフのアウトライン（Path）に落とす。

### 6.2 Image

Part の生バイト列（PNG / JPEG / WebP / AVIF / SVG）。再エンコードしない。読み手は `createImageBitmap(new Blob([bytes]))` でデコードする。

### 6.3 Path

Path2D に変換されるバイナリ。

```
varuint nVerbs
u8[nVerbs]  verbs   M=0 L=1 Q=2 C=3 Z=4 R=5(rect x y w h) E=6(ellipse x y rx ry rot a0 a1 ccw) A=7(arcTo x1 y1 x2 y2 r) O=8(roundRect x y w h r)
f32[...]    verb ごとの引数を verbs と同じ順で連結
```

- `nPaths` の `ext` は Path collection Part（`t: path`）の `index` 番目を指す。Path collection は `varuint n` に続けてパスデータを n 個並べたもの。グリフアウトラインを全ページで共有する用途。
- 読み手は Path2D をハッシュ+index でキャッシュする。

### 6.4 Paint

```
u8  kind      0=linear 1=radial 2=conic 3=pattern
linear : f32 x0 y0 x1 y1
radial : f32 x0 y0 r0 x1 y1 r1
conic  : f32 angle x y
上記3種: varuint nStops ; [ f32 offset, u32 rgba ]×n
pattern: varuint imageRef ; u8 repeat(0=repeat 1=repeat-x 2=repeat-y 3=no-repeat) ; f32×6 matrix
```

## 7. 命令セット（opset 1）

- 命令は `u8 opcode` に続けてオペランドを並べる。オペランド型: `f32`, `u8`, `u32`（色は `0xRRGGBBAA`）, `varuint`（参照・個数）, `str`（文字列表への varuint インデックス）。
- 描画状態（変換・クリップ・スタイル）は Canvas と同じくスタックで管理する。`USE` は暗黙に `save`/`restore` で包まれる。
- 未知の opcode は**エラー**。前方互換は `opset` バージョンと `EXT` 命令で確保する。

### 7.1 状態

| op | 名前 | オペランド | Canvas 対応 |
|---|---|---|---|
| 0x01 | SAVE | – | `save()` |
| 0x02 | RESTORE | – | `restore()` |
| 0x03 | TRANSFORM | f32 a b c d e f | `transform(a,b,c,d,e,f)` |
| 0x04 | TRANSLATE | f32 x y | `translate` |
| 0x05 | SCALE | f32 sx sy | `scale` |
| 0x06 | CLIP_PATH | varuint pathRef, u8 rule | `clip(path2d, rule)` |
| 0x07 | CLIP_RECT | f32 x y w h | `beginPath(); rect(); clip()` |

`rule`: 0=nonzero, 1=evenodd（以下同じ）。

### 7.2 スタイル

| op | 名前 | オペランド | Canvas 対応 |
|---|---|---|---|
| 0x10 | FILL_COLOR | u32 rgba | `fillStyle = "#…"` |
| 0x11 | FILL_PAINT | varuint paintRef | `fillStyle = gradient/pattern` |
| 0x12 | STROKE_COLOR | u32 rgba | `strokeStyle` |
| 0x13 | STROKE_PAINT | varuint paintRef | `strokeStyle` |
| 0x14 | LINE | f32 width, u8 cap, u8 join, f32 miter | `lineWidth/lineCap/lineJoin/miterLimit` |
| 0x15 | DASH | varuint n, f32[n], f32 offset | `setLineDash / lineDashOffset` |
| 0x16 | ALPHA | f32 | `globalAlpha` |
| 0x17 | BLEND | u8 | `globalCompositeOperation`（表は付録 A） |
| 0x18 | SHADOW | u32 rgba, f32 blur dx dy | `shadowColor/Blur/OffsetX/OffsetY` |
| 0x19 | FILTER | str | `filter`（任意。非対応環境では無視） |
| 0x1A | FONT | varuint fontRef, f32 size | `font = "<style> <weight> <size>px <family>"` |
| 0x1B | TEXT_STYLE | u8 align, u8 baseline, u8 dir, f32 letterSpacing | `textAlign/textBaseline/direction/letterSpacing` |

SHADOW の `blur`・`dx`・`dy` は unit で表す。Canvas の影は変換行列の影響を受けないので、読み手は設定時の変換の拡大率（行列式の平方根）を掛けて渡す。オフセットの向きはページの向きで、回転には追従しない。

`cap`: 0=butt 1=round 2=square。`join`: 0=miter 1=round 2=bevel。`align`: 0=left 1=right 2=center 3=start 4=end。`baseline`: 0=alphabetic 1=top 2=middle 3=bottom 4=hanging 5=ideographic。

### 7.3 図形

| op | 名前 | オペランド | Canvas 対応 |
|---|---|---|---|
| 0x20 | FILL_RECT | f32 x y w h | `fillRect` |
| 0x21 | STROKE_RECT | f32 x y w h | `strokeRect` |
| 0x22 | FILL_PATH | varuint pathRef, u8 rule | `fill(path2d, rule)` |
| 0x23 | STROKE_PATH | varuint pathRef | `stroke(path2d)` |
| 0x24 | FILL_PATH_AT | varuint pathRef, u8 rule, f32 x y | translate して fill。グリフや繰り返し部品用 |
| 0x25 | FILL_PATH_RUN | u8 rule, varuint n, [varuint pathRef, f32 x y]×n | 連続するグリフ描画を 1 命令で |
| 0x26 | CLEAR_RECT | f32 x y w h | `clearRect` |

### 7.4 テキスト

| op | 名前 | オペランド | Canvas 対応 |
|---|---|---|---|
| 0x30 | FILL_TEXT | str s, f32 x y, f32 advance | `fillText(s, x, y)` |
| 0x31 | STROKE_TEXT | str s, f32 x y, f32 advance | `strokeText` |

`advance` はエンコーダが計算した**期待幅**（unit）。読み手は `measureText(s).width` と比較し、差が閾値（推奨 0.5%）を超える場合は `advance / measured` で x 方向にスケールして描く（pdf.js と同じ手法）。これでフォント代替や文字送りの差異による行のはみ出し・ずれを抑える。`advance` が 0 のときは補正しない。

文字単位の正確な配置が必要な場合（カーニング・詰めのある見出し、埋め込み不可フォント等）は、エンコーダが 1 文字ずつ FILL_TEXT を出すか、アウトライン化して FILL_PATH_RUN を使う。

### 7.5 画像

| op | 名前 | オペランド | Canvas 対応 |
|---|---|---|---|
| 0x40 | IMAGE | varuint imgRef, f32 x y w h | `drawImage(img, x, y, w, h)` |
| 0x41 | IMAGE_SUB | varuint imgRef, f32 sx sy sw sh dx dy dw dh | 9 引数版 `drawImage` |
| 0x42 | SMOOTHING | u8 enabled, u8 quality | `imageSmoothingEnabled/Quality` |

### 7.6 合成

| op | 名前 | オペランド | 意味 |
|---|---|---|---|
| 0x50 | USE | varuint objRef | 子 Object を現在の変換で描く。暗黙 save/restore |
| 0x51 | USE_AT | varuint objRef, f32 x y | translate(x,y) して USE |
| 0x52 | GROUP_BEGIN | f32 alpha, u8 blend, f32 x y w h | グループ透過。読み手は一時キャンバスに描いて合成する |
| 0x53 | GROUP_END | – | |
| 0x54 | MASK_BEGIN | u8 kind, u32 backdrop, varuint n, u8[n] transfer | ソフトマスクの描画を始める |
| 0x55 | MASK_END | – | グループの内容にソフトマスクを掛ける |

- `USE` の対象 Object は読み手が `(hash, 現在の拡大率)` をキーにビットマップキャッシュしてよい（マスタースライド等の高速化）。
- `GROUP_BEGIN`/`END` は入れ子可。`x y w h` は一時キャンバスの範囲。一時キャンバスは現在の変換を引き継ぎ、そのほかの状態（塗り・線の色、線幅、フォントを除く）は初期状態から始まる。
- `MASK_BEGIN` 〜 `MASK_END` はソフトマスク（PDF の ExtGState `/SMask`）で、グループの中の最後、`GROUP_END` の直前に 1 つだけ置く。読み手は最も内側のグループの一時キャンバスと同じ大きさのマスク用キャンバスを用意し、現在の変換だけを引き継いで初期状態から、間の命令をそこに描く。`MASK_END` でマスク用キャンバスの各画素から値 m（0〜255）を求め、グループの一時キャンバスの各画素のアルファに m/255 を掛ける（`source-in` で合成してよい）。
  - `kind` 0（alpha）: m はマスクのアルファ。マスク用キャンバスは透明から始まる。
  - `kind` 1（luminosity）: マスク用キャンバスを `backdrop` の色（`0xRRGGBBAA`、アルファは無視）で塗ってから描き、m はその輝度 `0.3 R + 0.59 G + 0.11 B`（PDF の非分離ブレンドモードと同じ重み）。
  - `transfer` は空か 256 バイトで、256 バイトなら m を `transfer[m]` に置き換える（PDF の `/TR`。反転したマスクなど）。
  - マスク用キャンバスのクリップと状態は `MASK_END` で捨てる。グループに残ったクリップはマスクの範囲を狭めない。
  - マスクの中の命令は内容ではないので、テキスト抽出（MARK・LINK を含む）は `MASK_BEGIN` 〜 `MASK_END` を読み飛ばす。

### 7.7 メタ・拡張

| op | 名前 | オペランド | 意味 |
|---|---|---|---|
| 0x70 | LINK | f32 x y w h, str url | リンク領域（現在の座標系の矩形）。描画には影響しない |
| 0x71 | MARK | u8 kind, str payload | 構造マーク（§7.8）。描画には影響しない |
| 0xFF | EXT | u32 len, u8[len] | 拡張。未知なら読み飛ばす |

`LINK` の矩形は他の命令と同じく現在の変換行列の座標系で表す。読み手は 4 隅を変換した外接矩形をリンク領域とし、その領域に重なる run をリンクの文字列とする。`url` は次のいずれか。読み手はそれ以外の `url` をリンクにしない。

- `http:` / `https:` / `mailto:` の絶対 URL
- `#page=N` — 同じ View のページ（1 始まり。`scroll` View では帯）
- `#view=ID` または `#view=ID&page=N` — 別の View（`ID` は View の `id` を URL エンコードしたもの）とそのページ。ビューアはその View に切り替える（draw.io のページ間リンクなど）

### 7.8 MARK の種類とテキストの構造

`MARK` は描画に影響しない境界情報で、検索・選択・コピーのために命令列に読み順を与え、読み上げのために最小限の構造（見出し・リスト・表・図・言語）を与える。

| kind | 名前 | payload | 意味 |
|---|---|---|---|
| 0 | PARAGRAPH | 任意 | 段落の開始。以降の run は前の run と結合されない |
| 1 | LINE | 任意 | 同じ段落内の行の開始。前の run との間に空白 1 つがあるものとして扱う |
| 2 | CELL | セル参照（例 `B12`） | 表・シートのセルの開始。PARAGRAPH と同じ境界 |
| 3 | BOX | 任意 | テキストボックス・図形内テキストの開始。PARAGRAPH と同じ境界 |
| 4 | ALT_TEXT | 文字列 | 直後の描画命令 1 つ（`FILL_TEXT` / `STROKE_TEXT` / `FILL_PATH_AT` / `FILL_PATH_RUN` / `USE` / `USE_AT`）が表す文字列。テキスト抽出ではその命令の文字列の代わりにこの文字列を、位置はその命令の位置を使う（`USE` / `USE_AT` では、子 Object の bbox の x 方向の範囲をこの文字列の範囲とし、位置を bbox の左端に、字送りを bbox の幅にする）。合字や私用領域の文字で描いた run の本来の文字列、アウトライン化した文字、1 文字ずつ描いた縦書きなどに使う。描画命令が続かない場合は位置なしの run になる |
| 5 | WRAP | 任意 | 同じ段落内の行の折り返し。前の run との間に区切りを入れずに結合する。和文など語を空白で区切らない文字の間で折り返した行に使い、行をまたぐ検索を可能にする |
| 6 | HEADING | レベル `1`〜`6` | 見出しの段落の開始。PARAGRAPH と同じ境界 |
| 7 | LIST | 任意 | リストを開く（END で閉じる） |
| 8 | LIST_ITEM | 任意 | 最も内側の LIST の項目の開始。次の LIST_ITEM か LIST の END までが項目の中身（入れ子の LIST や続きの段落を含められる）。最も内側が LIST でなければ PARAGRAPH と同じ |
| 9 | TABLE | 任意 | 表を開く（END で閉じる）。中のセルは CELL で始める |
| 10 | FIGURE | 代替テキスト | 図を開く（END で閉じる）。範囲は中の描画命令の外接矩形。代替テキストが空なら装飾で、読み上げない |
| 11 | END | 任意 | 開いている LIST / TABLE / FIGURE のうち最も内側を閉じる。開いていなければ何もしない |
| 12 | LANG | 言語（BCP 47）。空なら文書の既定（`meta.dc.language` の最初の値） | 以降の run の言語 |

- `MARK` の効果は次の `FILL_TEXT` / `STROKE_TEXT` / `ALT_TEXT` に及び、run どうしの結合規則を決める。同じ行で書式だけが変わった run の間には `MARK` を置かない（結合される）。
- `MARK` を出さないエンコーダのために、読み手は位置に基づく推定（y が変われば行、同じ行で字送りの 0.2 倍以上の隙間があれば空白）で補ってよい。ただし推定は不正確なので、エンコーダは `MARK` を出すことが強く推奨される。
- run どうしの区切りは kind だけで決まる。LINE は空白、WRAP は連結、PARAGRAPH / CELL / BOX / HEADING / LIST / LIST_ITEM / TABLE / FIGURE / END は段落境界。ALT_TEXT と LANG は区切りに影響しない（LANG の後の run は位置による推定のままで、保留中の ALT_TEXT もそのまま次の描画命令に付く）。

#### 構造

読み上げ用の構造は、MARK を走査順にたどって組み立てる。

- **葉の MARK**（PARAGRAPH、BOX、HEADING、LIST の外の LIST_ITEM、TABLE の外の CELL）は次の run が属する段落の種類を予約する。run を挟まずに続いた場合は最後のものが決める（BOX の直後の HEADING は見出しの段落）。
- **構造の MARK**（LIST、TABLE、FIGURE、END、最も内側が LIST のときの LIST_ITEM、最も内側が TABLE のときの CELL）はその場で効き、予約を取り消す。LIST や TABLE の中で項目・セルの前に現れた run は、それぞれ暗黙の項目、表の見出し（caption）になる。
- 最も内側が TABLE のときの CELL の payload は、表の左上を `A1` とする参照か範囲（`B2:C2`、結合セル）で、後ろに ` col`（列見出し）か ` row`（行見出し）を付けられる（例 `A1 col`）。セルは左上の位置の行優先の順に並べる。TABLE の外の CELL はシートのセル参照（例 `B12`、結合セルは範囲 `B2:C3`）で、同じく ` col` / ` row` を付けて見出しのセルにできる（シートの中の表の見出し行）。
- 開いている構造と LANG は、`USE` した子 Object や SAVE/RESTORE をまたいで走査順に一直線に続き、最上位の Object（ページのレイヤー、シートのタイル）ごとに初期状態（構造なし、言語は文書の既定）に戻る。最上位 Object の終わりで開いている構造は閉じる。共有プレフィックス（§5 の `USE`）で命令列が分かれても同じ結果になるための規則である。エンコーダは開いた構造を END で閉じ、言語を持つ子 Object を `USE` した後は必要なら LANG を出し直す。
- FIGURE の範囲は、中の描画命令（画像、矩形、パス、子 Object の中身、テキスト）を現在の変換で移した外接矩形を、その時点のクリップの外接矩形と交わらせたもの。中身が空なら大きさ 0 とする。FIGURE の中のテキストは検索・選択・コピーの対象のままだが、代替テキストがあれば読み上げには代替テキストを使う。

### 7.9 テキスト索引 Part

View ごとに任意で持てる索引。`View.textIndex` にハッシュを置く（`t: idx`）。全 Object をデコードせずに全文検索するため、またサーバー側の検索エンジンに同じデータを流すためのもの。

```
u8[4]   "BTXT"
u16     version (1)
varuint nRuns
run ×n:
  varuint a        fixed/flow/scroll: ページ（帯）番号   sheet: タイル x
  varuint b        fixed/flow/scroll: レイヤー番号       sheet: タイル y
  varuint ordinal  その Object のテキスト抽出（§8 のテキストバックエンド、USE の子を含む走査順）における run の通し番号。シートでは §4.1 の規則でその Tile に属さない run を載せないので、番号は飛ぶことがある
  u8      sep      前の run との結合: 0=連結 1=空白 2=段落境界（一致は境界をまたがない）
  str     text     run の文字列（ALT_TEXT の場合はその文字列）
```

- 順序は読み順。読み手は `text` を `sep` に従って連結した平文に対して検索し、ヒットを (a, b, ordinal, 文字範囲) に戻す。
- ヒットの矩形は該当 Object をデコードして run の位置・フォント・`advance` から計算する（部分一致は接頭辞幅の計測で求める）。
- 正規化（NFKC、大文字小文字、かな）は索引には施さず、検索時に読み手が行う。

## 8. 読み手（ビューア）アーキテクチャ

```
main thread                          worker
──────────────                       ────────────────────────────────────
Viewer UI                            Loader      fetch / DecompressionStream / Part キャッシュ
  スクロール・ズーム         ──msg──▶  Decoder     Object → 命令配列（またはその場実行）
  <canvas> or 転送済み        ◀─bitmap─ Renderer    OffscreenCanvas に描画。Path2D / ImageBitmap / FontFace / USE ビットマップのキャッシュ
  テキストレイヤー(DOM)      ◀─runs──  TextExtractor 同じ命令列を走査して文字列と矩形を返す
```

- ページやタイルなど最上位の Object は Canvas の初期状態（黒の塗り・線、線幅 1、`10px sans-serif`、不透明、`source-over`）から描き始める。`USE` された子 Object は親の状態を引き継ぐ（塗りを指定しないパターンセルなどに使う）。
- 描画は `(Object, 変換行列, クリップ矩形)` の純関数。Tile・ページ・高倍率時の部分再描画をすべて同じ経路で処理する。
- 命令列は「実行できるデータ」なので、バックエンドを差し替えられる: Canvas レンダラ、テキスト抽出（検索・選択・読み上げ用の DOM レイヤー。MARK の構造と LINK のリンクも含む）、将来的な WebGL/WebGPU レンダラ。テキストレイヤーのために別の Part を持つ必要はない。
- 検索はビューアではなくライブラリ（Worker）が提供する。索引 Part（§7.9）があればそれを、なければ Object を走査して同じ形の run 列を作り、正規化して検索し、ヒットの矩形を返す。ビューアはヒット一覧とハイライトの描画だけを担当する。
- 大きな文書ではページ表を Index Part にして、可視範囲のページだけ Part を取得する。single 形式でも `off`/`len` により Range 取得できる。

必要なブラウザ機能（いずれも 2023 年時点の主要ブラウザで利用可能）: `DecompressionStream("deflate-raw")`、Worker 内 `OffscreenCanvas`、Worker 内 `FontFace` / `self.fonts`、`createImageBitmap`、`Path2D`、`roundRect`。暗号化した文書（§3.5）には `crypto.subtle` が要る（HTTPS か localhost の安全なコンテキストでだけ使える）。`letterSpacing`、`filter` は任意機能とし、非対応環境では無視または代替描画する。

## 9. 書き手（エンコーダ）の責務

- レイアウトは完了させてから出力する（BDF に再レイアウトはない）。
- 同一内容の描画は**同じバイト列**になるよう決定的に出力する（浮動小数の丸め、文字列表の順序、リソースの並びを正規化）。共有はハッシュ一致で自動的に起きる。
- フォントはサブセット化して WOFF2 で埋め込む。埋め込めない場合は `system` かアウトライン化。
- `FILL_TEXT` の `advance` を必ず埋める。
- シートは Tile 原点相対で出力し、またがるオブジェクトは共有 Object + `USE` にする。
- Part の推奨順（§3.3）に従って並べる。
- パスワードで保護された入力は、同じパスワードで暗号化して出力する（§3.5）。

## 付録 A. BLEND 値

0=source-over 1=multiply 2=screen 3=overlay 4=darken 5=lighten 6=color-dodge 7=color-burn 8=hard-light 9=soft-light 10=difference 11=exclusion 12=hue 13=saturation 14=color 15=luminosity 16=destination-over 17=destination-in 18=destination-out 19=source-in 20=source-out 21=source-atop 22=destination-atop 23=xor 24=copy 25=lighter

## 付録 B. 命令列の例

幅 960×540 のスライド。マスター（Object #0）を敷き、タイトルを描く。

```
50 00                         USE obj[0]
1A 00 00 00 C0 41             FONT font[0], 24.0
10 20 20 20 FF                FILL_COLOR #202020FF
30 00 00 00 A0 42 00 00 F0 42 00 00 34 43
                              FILL_TEXT str[0] x=80 y=120 advance=180
```

## 付録 C. 検討したが採用しなかったもの

- **ZIP コンテナ**: ツールが豊富だが、セントラルディレクトリが末尾にあり先頭からのストリーミング描画と相性が悪い。ローカルヘッダの重複や ZIP64 も不要な複雑さ。single 形式は「ディレクトリ先頭の単純な連結」で十分。
- **ファイル全体の 1 本圧縮**: Range・Tile 単位取得ができなくなる。
- **可変長整数の固定小数点座標**: f32 より小さくなりうるが、デコードが複雑になる。deflate で差はかなり縮む。必要なら opset 2 で列指向（opcode 列 / f32 列 / 参照列を分離）を検討する。
- **テキスト専用 Part**: 命令列を別バックエンドで走査すれば得られるので不要。
- **ファイル全体の暗号化**: 配布層に任せれば仕様は変えずに済むが、Range 取得・split 形式・CDN への配置という BDF の中心の性質を失う。
- **Part 名を鍵付き HMAC にする**: 平文のハッシュは隠せるが、名前の検証にも single ↔ split の組み替えにも鍵が要る。格納バイト列のハッシュで足りる。
- **Argon2 などメモリハードな鍵導出**: 総当たりには強いが WebCrypto にないので、ビューアに wasm が要る。PBKDF2 は WebCrypto でそのまま使える。
