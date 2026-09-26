# BDF (Browser Document Format) 仕様ドラフト v0.1

> 状態: ドラフト。名称 "BDF" は仮称（X11 のビットマップフォント形式 BDF と衝突するため、公開時には要再考）。

## 1. 目的と非目的

**目的**

- Office 系ファイル（PowerPoint / Excel / Word 等）を**サーバー側で変換**し、**ブラウザ（Worker）で Canvas 2D にそのまま描画**できる中間フォーマット。
- 命令セットは `CanvasRenderingContext2D` の API に 1:1 に近い形で対応させ、デコーダが「バイト列を読んで ctx のメソッドを呼ぶだけ」になるようにする。
- 圧縮・展開は `CompressionStream` / `DecompressionStream`（`deflate-raw`）に任せ、独自実装を持たない。
- **1ファイル**（配布・キャッシュ向け）と**分割ファイル群**（オンデマンド取得・CDN 向け）を同じ論理構造で表現できる。
- スライドマスターやページ背景など**同一内容の繰り返しは共有オブジェクトとして1回だけ格納**し、サイズを抑える。
- 3種類のレイアウトモデルを表現できる。
  - 固定サイズページの列（PowerPoint、PDF）
  - 無限サイズの平面（Excel のシート）
  - ページ分割されつつ「縦一枚」としても閲覧できる文書（Word）

**非目的**

- 編集可能な文書モデル（DOM のような意味構造）を持つこと。BDF は「描画済み表示リスト」であり、再レイアウトはしない。
- 印刷用の高度な色管理（CMYK、ICC）。sRGB のみ。
- 暗号化・署名。必要なら配布層（HTTPS、ストレージ）で行う。

## 2. 用語

| 用語 | 意味 |
|---|---|
| Part | コンテナに格納される1つのバイト列。内容のハッシュで識別される。 |
| Object | 描画命令列を含む Part。ページ本体・マスター・シートのタイルなど、描けるものはすべて Object。 |
| View | 閲覧単位。スライド集、シート、文書など。1 BDF に複数の View を持てる（Excel の複数シート等）。 |
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
offset 0   : magic          u8[4]   "BDF1"
offset 4   : version        u16     フォーマットバージョン（1）
offset 6   : flags          u16     予約（0）
offset 8   : manifestOff    u64     manifest の先頭オフセット（通常 32）
offset 16  : manifestLen    u64     manifest の圧縮後長
offset 24  : manifestEnc    u8      0=identity, 1=deflate-raw
offset 25  : reserved       u8[7]
offset 32  : manifest
offset ... : parts 領域（manifest 直後から始まり、manifest.parts に記載された順に連続配置）
```

- すべてリトルエンディアン。
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

## 4. Manifest

JSON。読みやすさとツールでの扱いやすさを優先する。巨大になる部分（ページ表・タイル表）は Index Part に逃がす。

```jsonc
{
  "bdf": 1,                     // フォーマットバージョン
  "opset": 1,                   // 命令セットバージョン（§7）
  "unit": "pt",                 // 1 unit = 1/72 inch。座標はすべてこの単位
  "meta": { "title": "…", "source": "pptx", "generator": "bdf-go/0.1" },

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

`parts[].t` の値: `obj` / `font` / `img` / `path` / `idx`。

### 4.1 View の種類

**`fixed`** — 固定サイズページの列。ページごとに `w`,`h` が異なってよい。ビューアはページを縦または横に並べる、あるいは 1 ページずつ表示する。

**`flow`** — `fixed` に加え、各ページが `body`（余白を除いた本文矩形）を持ち、レイヤーに役割（`role`）が付く。ビューアは 2 つの表示モードを提供できる。

- ページモード: 紙の形で全レイヤーを描く。
- 連続モード: `role: body` のレイヤーだけを `body` 矩形で切り出し、`continuous.gap` を挟んで縦に積む。ヘッダー・フッター・背景は描かない。

連続モードは**再レイアウトではなく本文矩形の積み上げ**である（BDF は表示リストであり、テキストの折り返し位置は変えられない）。「1 枚の長いページとして再レイアウトしたもの」も欲しい場合、変換側が別 View（`kind: fixed`、ページ 1 枚）を追加で生成する。共有 Object により画像やフォントは二重に格納されない。

**`sheet`** — 原点 (0,0) から右下に伸びる無限平面。

- `cols` / `rows` はランレングス `[count, size]` の配列。ビューアは行・列ヘッダー、グリッド線、固定ペイン（`freeze`）をこの情報から自前で描く。
- 内容は `tile` unit 四方の Tile に分割され、各 Tile は 1 Object。Tile 内の座標は **Tile 原点を基準**にする（f32 の精度限界を超えないため）。
- 複数 Tile にまたがるオブジェクト（結合セル、画像、グラフ）は共有 Object にして、触れる各 Tile から `USE` で参照する。ビューアは Tile 矩形でクリップして描くので重複描画は正しく処理される。
- 空 Tile は表に載せない。
- `tiles` の Index Part 形式（`t: idx`）: `[u32 tx, u32 ty, u8[16] hash]` を (ty, tx) 昇順に並べた固定長レコード列。二分探索で引く。

### 4.2 レイヤーの役割（role）

`background` / `master` / `header` / `footer` / `body` / `notes` / `annotation`。ビューアは役割ごとの表示切替に使える。未知の role は `body` として扱う。

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

- `USE` の対象 Object は読み手が `(hash, 現在の拡大率)` をキーにビットマップキャッシュしてよい（マスタースライド等の高速化）。
- `GROUP_BEGIN`/`END` は入れ子可。`x y w h` は一時キャンバスの範囲。

### 7.7 メタ・拡張

| op | 名前 | オペランド | 意味 |
|---|---|---|---|
| 0x70 | LINK | f32 x y w h, str url | リンク領域。描画には影響しない |
| 0x71 | MARK | u8 kind, str payload | 構造マーク（§7.8）。描画には影響しない |
| 0xFF | EXT | u32 len, u8[len] | 拡張。未知なら読み飛ばす |

### 7.8 MARK の種類とテキストの構造

`MARK` は描画に影響しない境界情報で、検索・選択・コピーのために命令列に読み順を与える。

| kind | 名前 | payload | 意味 |
|---|---|---|---|
| 0 | PARAGRAPH | 任意 | 段落の開始。以降の run は前の run と結合されない |
| 1 | LINE | 任意 | 同じ段落内の行の開始。前の run との間に空白 1 つがあるものとして扱う |
| 2 | CELL | セル参照（例 `B12`） | 表・シートのセルの開始。PARAGRAPH と同じ境界 |
| 3 | BOX | 任意 | テキストボックス・図形内テキストの開始。PARAGRAPH と同じ境界 |
| 4 | ALT_TEXT | 文字列 | 直後の描画命令 1 つ（`FILL_TEXT` / `STROKE_TEXT` / `FILL_PATH_AT` / `FILL_PATH_RUN` / `USE` / `USE_AT`）が表す文字列。テキスト抽出ではその命令の文字列の代わりにこの文字列を、位置はその命令の位置を使う。合字や私用領域の文字で描いた run の本来の文字列、アウトライン化した文字、1 文字ずつ描いた縦書きなどに使う。描画命令が続かない場合は位置なしの run になる |
| 5 | WRAP | 任意 | 同じ段落内の行の折り返し。前の run との間に区切りを入れずに結合する。和文など語を空白で区切らない文字の間で折り返した行に使い、行をまたぐ検索を可能にする |

- `MARK` の効果は次の `FILL_TEXT` / `STROKE_TEXT` / `ALT_TEXT` に及び、run どうしの結合規則を決める。同じ行で書式だけが変わった run の間には `MARK` を置かない（結合される）。
- `MARK` を出さないエンコーダのために、読み手は位置に基づく推定（y が変われば行、同じ行で字送りの 0.2 倍以上の隙間があれば空白）で補ってよい。ただし推定は不正確なので、エンコーダは `MARK` を出すことが強く推奨される。

### 7.9 テキスト索引 Part

View ごとに任意で持てる索引。`View.textIndex` にハッシュを置く（`t: idx`）。全 Object をデコードせずに全文検索するため、またサーバー側の検索エンジンに同じデータを流すためのもの。

```
u8[4]   "BTXT"
u16     version (1)
varuint nRuns
run ×n:
  varuint a        fixed/flow: ページ番号   sheet: タイル x
  varuint b        fixed/flow: レイヤー番号 sheet: タイル y
  varuint ordinal  その Object のテキスト抽出（§8 のテキストバックエンド、USE の子を含む走査順）における run の通し番号
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
- 命令列は「実行できるデータ」なので、バックエンドを差し替えられる: Canvas レンダラ、テキスト抽出（検索・選択用の DOM レイヤー）、ヒットテスト（リンク）、将来的な WebGL/WebGPU レンダラ。テキストレイヤーのために別の Part を持つ必要はない。
- 検索はビューアではなくライブラリ（Worker）が提供する。索引 Part（§7.9）があればそれを、なければ Object を走査して同じ形の run 列を作り、正規化して検索し、ヒットの矩形を返す。ビューアはヒット一覧とハイライトの描画だけを担当する。
- 大きな文書ではページ表を Index Part にして、可視範囲のページだけ Part を取得する。single 形式でも `off`/`len` により Range 取得できる。

必要なブラウザ機能（いずれも 2023 年時点の主要ブラウザで利用可能）: `DecompressionStream("deflate-raw")`、Worker 内 `OffscreenCanvas`、Worker 内 `FontFace` / `self.fonts`、`createImageBitmap`、`Path2D`、`roundRect`。`letterSpacing`、`filter` は任意機能とし、非対応環境では無視または代替描画する。

## 9. 書き手（エンコーダ）の責務

- レイアウトは完了させてから出力する（BDF に再レイアウトはない）。
- 同一内容の描画は**同じバイト列**になるよう決定的に出力する（浮動小数の丸め、文字列表の順序、リソースの並びを正規化）。共有はハッシュ一致で自動的に起きる。
- フォントはサブセット化して WOFF2 で埋め込む。埋め込めない場合は `system` かアウトライン化。
- `FILL_TEXT` の `advance` を必ず埋める。
- シートは Tile 原点相対で出力し、またがるオブジェクトは共有 Object + `USE` にする。
- Part の推奨順（§3.3）に従って並べる。

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
