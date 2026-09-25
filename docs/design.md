# 設計メモ: 判断の理由と実装方針

`spec.md` が「何を作るか」、このファイルが「なぜそうしたか」と「どう作るか」。

## 1. PDF と何が違うのか

| | PDF | BDF |
|---|---|---|
| 描画モデル | PostScript 由来。独自のグラフィックス状態、フォント機構、色空間 | Canvas 2D そのもの。ブラウザが持っていない概念は入れない |
| フォント | 埋め込み形式が多数（Type1/CFF/TrueType/Type3）。ラスタライザが必要 | WOFF2 を `FontFace` に渡すだけ。テキスト描画はブラウザ任せ |
| 画像 | 独自フィルタ（DCT/Flate/CCITT/JBIG2…） | ブラウザがデコードできる形式をそのまま格納 |
| 圧縮 | Flate 等を自前実装 | `DecompressionStream` |
| 構造 | xref テーブル、オブジェクトグラフ | 内容アドレスの Part 群 + JSON manifest |
| 共有 | Form XObject | Object（内容アドレスで自動 dedup） |
| 無限平面 | なし | `sheet` View + Tile |
| 部分取得 | Linearized PDF（複雑、対応まちまち） | Part 単位の Range / 個別ファイル |

pdf.js が数万行かけてやっていること（フォントラスタライズ、フィルタ、パーサ）をブラウザの標準 API に肩代わりさせるのが狙い。デコーダは「バイト列を読んで ctx を呼ぶループ」だけになる。

## 2. Go と wasm について

**結論: サーバー側エンコーダ・変換 CLI は Go、ブラウザ側ビューアは TypeScript。ビューアで Go/wasm は使わない（少なくとも v1 では）。**

理由:

- **デコーダは軽い**: フォーマットを Canvas API に 1:1 で寄せたので、デコーダは DataView を読むループで数百行に収まる。wasm で速くする対象がない。
- **境界コストが支配的**: 描画のホットパスは Canvas API 呼び出し。Go の wasm から `syscall/js` 経由で `ctx.fillText` などを呼ぶのは、TS から直接呼ぶより遅く、アロケーションも増える。命令 1 つごとに JS↔wasm を往復する構造は最悪のパターン。
- **バイナリサイズ**: Go 標準ツールチェーンの wasm は最小でも 2MB 級。プレビュー用ビューアの初期ロードとして重い。TinyGo で縮むが、GC・リフレクション・goroutine の制約と引き換えになる。
- **ブラウザ API への依存**: `DecompressionStream`、`OffscreenCanvas`、`FontFace`、`createImageBitmap` はすべて JS 側の非同期 API。wasm で包んでも JS の薄い層が必ず残る。

wasm が意味を持つケース:

- **クライアントで変換したい**とき（例: ユーザーがドロップした xlsx をブラウザ内で BDF に落として描く）。エンコーダ本体を Go で書いておけば、そのまま wasm にして使い回せる。これが「Go で書く」ことの本当の利点で、ビューア側の高速化ではない。
- Path2D 生成やテキスト抽出のような、純粋なデータ処理で JS より速くしたい部分が見つかった場合。プロファイルしてから判断する。

したがって:

- Go: `bdf` パッケージ（Writer / Object builder / 各 Part のエンコード）と変換 CLI。エンコーダは `io.Writer` に対して決定的に出力する。
- TypeScript: `@bdf/core`（デコード・型定義）、`@bdf/render`（Canvas バックエンド、Worker）、`@bdf/viewer`（UI）。
- 両者の契約は**ワイヤフォーマットとフィクスチャ**。Go でエンコードしたテストファイルを TS がデコードし、Playwright でスクリーンショットを golden 比較する。Go 側に描画は持たない。

## 3. 変換パイプライン（現実的な順序）

Office ファイルを直接 BDF にするには Word 相当のレイアウトエンジンが要る。ここが一番重いので、段階を踏む。

1. **PDF → BDF**（最初のターゲット）
   - PDF の内容ストリーム演算子（`m l c re f S q Q cm Tf Tj Do …`）は BDF の命令にほぼ 1:1 で落ちる。
   - LibreOffice headless で Office → PDF にしてから BDF 化すれば、それだけで PowerPoint / Word のプレビューが高忠実度で成立する。
   - フォントは PDF 埋め込みフォントを抽出して WOFF2 に再パッケージ（CFF/TrueType ならほぼ機械的。Type3 はアウトライン化）。
   - 難点: PDF ではマスター共有が失われている（各ページに同じ命令が展開される）。同一バイト列になる Object は自動 dedup されるが、微妙に違うページでは効かない。ページ内の「先頭の共通プレフィックス」を切り出して共有 Object にするヒューリスティックで補う（§3.1）。
2. **XLSX → BDF**（直接変換）
   - PDF 経由だと無限シートが失われるので、こちらは直接。Go には excelize があり、セル値・書式・列幅・結合・条件付き書式（一部）が取れる。セルのレイアウトは行列の格子なので、文書レイアウトほど難しくない。
   - グラフや図形は当面画像化（LibreOffice に描かせて PNG）で逃げる。
3. **PPTX / DOCX → BDF**（直接変換）
   - PPTX は絶対配置なので DOCX より先に手が届く（テキストボックス内の折り返しは必要）。実装済み（§3.3）。
   - マスター・レイアウト・スライドの継承構造が BDF の共有 Object にそのまま対応するので、直接変換するとサイズ面の効果が最も大きい。
   - DOCX は Word 相当のレイアウトエンジンが要るので長期課題。

## 3.1 PDF → BDF 変換器（converter/pdf）の構造

`converter/pdf` パッケージは PDF オブジェクト層に pdfcpu を使い、内容ストリームの解釈は自前で行う。CLI は `bdf generate`（入力の形式は中身から判別する）。

- **命令の対応**: `q`/`Q` → SAVE/RESTORE、`cm` → TRANSFORM、パス演算子 → Path、`W n` → CLIP_PATH、`Do`（Form）→ 共有 Object + USE、`Do`（Image）→ IMAGE、`sh` と shading パターン → Paint、tiling パターン → セル Object を USE_AT で敷き詰め、ExtGState の `ca`/`CA`/`BM` → ALPHA/BLEND、透明グループ → GROUP。座標は PDF のユーザー空間をそのまま TRANSFORM で写す（ページ先頭で y 反転と回転を 1 回かける）。
- **状態の遅延出力**: 色・線・アルファなどは描画命令の直前に、前回出力した値と違うときだけ書く。SAVE/RESTORE で出力済み状態のスタックも巻き戻す。
- **テキスト**: 行列（Tm/Td/T*）ごとに SAVE + TRANSFORM のブロックを開き、その中で x オフセットだけで FILL_TEXT を並べる。`advance` には PDF の幅を入れる。Skia のように 1 グリフずつ `Td` で位置決めする PDF は、同じ行で筆記位置が連続していれば同じ run に結合する（カーニング分は run の advance に吸収）。`Tz`/`Ts` は変換行列、`Tc` は letterSpacing、`Tw` はスペースを独立 run にして advance 補正で広げる。不可視テキスト（Tr 3）は透明色で描いて検索可能にする。`/ActualText` は run の文字列に使う。
- **フォント**: 埋め込み TrueType/CFF/OpenType は、使われたコードごとに (Unicode, GID) を集め、Unicode → GID の cmap を合成して TTF/OTF に組み直す（name、OS/2、post も生成し、ブラウザのサニタイザを通す）。同じ Unicode に別のグリフが割り当たる場合や合字（ToUnicode が複数文字）は私用領域 U+E000〜 に逃がし、`ALT_TEXT` で本来の文字列を持つ。TrueType は使ったグリフ（合成グリフの構成要素を含む）以外のアウトラインを空にして `glyf` を縮める（GID は付け替えないので `loca`/`hmtx` の長さはグリフ数のまま。フルフォント埋め込みの PDF で 611KB → 56KB 程度になる。`-no-subset` で無効化）。CFF は CharStrings INDEX と各オフセットの書き換えが必要なので現状は縮めない。Type1 (FontFile) は変換せずシステムフォントに落とす。非埋め込みフォントは名前とフラグから serif/sans-serif/monospace と太さ・斜体を決める。Type3 はグリフ手続きを Object にして USE する。
- **画像**: DCT はそのまま JPEG、それ以外はデコードして PNG（SMask/ステンシルマスクはアルファに合成、ImageMask は塗り色で PNG 化）。JPX と JBIG2 は未対応。格納前に `imgconv`（§3.2）を通す。
- **注釈**: リンクは LINK 命令、外観ストリームは Form として描く。
- **ページをまたぐ共通プレフィックスの共有**: PDF はマスター（ヘッダ・ロゴ・フッタ）を各ページの内容ストリームに展開してしまうので、変換後の各ページ Object の先頭から一致するバイト列を切り出して共有 Object にする（`bdf.SharePrefixes`）。切れる位置は「深さ 0 の命令境界」に限る（SAVE/RESTORE と GROUP が釣り合っていて、それより前の深さ 0 に CLIP/SHADOW/FILTER がなく、直前が MARK でない）。共有部分は USE で呼ぶが USE は暗黙の save/restore を持つので、残り部分の先頭で切断時点の状態（正味の変換行列、最後に出力した塗り・線・アルファ・フォントなど）を書き直してから続きを出す。候補の鍵は命令列と参照リソース（Path、Paint、フォント、画像、子 Object）の内容の累積ハッシュで、同じ鍵を持つページの組を「(ページ数 − 1) × 切り出すバイト数」の大きい順に貪欲に採用する（3 ページで共有できる短い接頭辞を、2 ページだけで共有できる長い接頭辞より優先する）。512 バイト未満の接頭辞は Part のオーバーヘッドの方が大きいので共有しない。`-no-share` で無効化。
- **未対応（警告を出して無視）**: ExtGState のソフトマスク、メッシュ系シェーディング（平均色で代用）、Type 4 関数（中央値で代用）、埋め込みでない定義済み CMap。

テスト用 PDF は Chromium（Skia）と reportlab で生成し（`npm run test:pdf:gen`）、変換結果は `fixtures/pdf/` に置いて golden テストで描画を比較する。

## 3.2 画像の格納と変換（imgconv）

画像の扱いは「そのまま」と「変換あり」の 2 モードで、`pdf.Options.Images` / `pptx.Options.Images`（CLI は `-images keep|convert`）で選ぶ。今後の Office 変換器も同じパッケージを使う。

- **Keep**: 元ファイルにある JPEG/PNG などをそのまま Part にする。デコードした画素は PNG にする。ブラウザ内で変換を走らせる場合はこちらを使う。
- **Convert**: サーバー側の事前変換用。時間は掛けてよい前提で、次の順に決める。
  1. WebP / AVIF / SVG はそのまま。
  2. JPEG は WebP 非可逆（既定 q80, method 6）を試す。
  3. PNG / GIF / BMP と、PDF からデコードした画素は WebP 可逆を試す。色数が多く写真らしい画像（標本で 4096 色超）は非可逆も試す。
  4. 元より小さくなった候補のうち最小のものを採用し、小さくならなければ元のまま。

コーデックは libwebp（エンコーダのみ）を wasi-sdk で wasm にし、[shibukawa/wasm2go-fork](https://github.com/shibukawa/wasm2go-fork)（pgmem ブランチ）で純 Go に変換したもので、cgo も wasm ランタイムも使わない。生成物は `imgconv/internal/webpw`（スカラー、約 5MB、44 ファイル）と `imgconv/internal/webpwsimd`（SIMD、`GOEXPERIMENT=simd` 専用、後述）で、`tools/gen-codecs.sh` で再生成する。フォークの `-symbol-names`（関数名を wasm の name セクションから付ける）、`-group-files`（`vp8_enc.go` のように主題ごとのファイルに分ける）、`-addr-consts`（静的データのアドレスを名前付き定数にする）を使い、libwebp を更新しても差分が小さく収まるようにしている。gen2brain 同梱の wasm は name セクションが落とされているので、自前でビルドしている。

**ビルドタグ**

| タグ | 効果 |
|---|---|
| （なし） | WebP 可逆・非可逆を同梱 |
| `bdf_noconv` | コーデックを一切リンクしない。Convert を指定しても Keep として動き、`ErrNotAvailable` を返す。tinygo でブラウザ向けに wasm 化するときはこれを使う |

**測定（1024×768、純 Go、4 コアのコンテナ）**

| 画像 | PNG | JPEG q80 | WebP 可逆 | WebP q80 m6 | AVIF q60 s10 | AVIF q60 s6 | AVIF 可逆 |
|---|---|---|---|---|---|---|---|
| 図版（少色数） | 6.7 KB | 64 KB | 0.3 KB / 87 ms | 14 KB / 150 ms | 28 KB / 155 ms | 3.2 KB / 1.0 s | 63 KB / 0.2 s |
| 写真風 | 1463 KB | 90 KB | 1405 KB / 0.5 s | 43 KB / 0.35 s | 41 KB / 0.16 s | 39 KB / 1.35 s | 1422 KB / 0.4 s |

フォークで変換したコードは、gen2brain 同梱の Go 化コードより WebP で 1.6〜2 倍、wazero 実行の AVIF より 3〜4 倍速く、出力はバイト単位で同一だった。

**同梱の wasm はスカラーでビルドしている**（測定の結果）。wasm SIMD を有効にした 2 種類のビルド（clang の自動ベクトル化のみ、および emscripten の SSE 互換ヘッダで libwebp 自身の SSE2/SSE4.1 経路を有効化したもの）を、wasm2go の純 Go モードと asm バックエンド（`-fuse-loops`、SIMD 命令 10,866 箇所をインライン展開）で変換して比べたが、いずれもスカラーより遅かった。v128 を `[2]uint64` の組で運ぶこれらのバックエンドでは、4×4〜16×16 ブロック単位の短い SIMD カーネルの命令 1 つごとに汎用レジスタとの往復が入り、ベクトル化の利得を上回る。

| 1024×768 写真風 q80 m6 | スカラー | 自動ベクトル化 | SSE2/SSE4.1 経路 |
|---|---|---|---|
| 純 Go モード（GOAMD64=v2） | 0.35 s | 0.75 s | 2.4 s |
| asm バックエンド | 0.35 s | 0.39 s | 0.52 s |

この往復をなくすため、フォークに `-simd=go127` を追加した。Go 1.27 の `GOEXPERIMENT=simd` で入る `simd/archsimd` の 128 bit ベクトル型で v128 を運び、SIMD 命令ごとにインライン展開される `base.Simd_g_*` メソッド連鎖を生成する（amd64 は AVX/AVX2、arm64 は NEON）。v128 に触れる関数は `[2]uint64` 版と archsimd 版の 2 本が出力され、`goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)` のビルドタグで切り替わる。archsimd は互換性保証がなく 1.26 と 1.27 で API が変わったため、対象の Go リリースを明示し、生成コード本体は `base.V128` と `base.Simd_g_*` だけを参照する。移植性を重視した `simd` パッケージも検討したが、ベクトル長が実行時に決まり（AVX2 機で 256 bit）、shuffle や narrow などの命令が無く、simd 型を使う関数が 4 本にクローンされるため、128 bit 固定の v128 には使えなかった。

SSE 経路の libwebp を `-simd=go127` で変換し、`GOEXPERIMENT=simd` でビルドして測った結果（各 5 回、出力はスカラー版とバイト単位で同一）:

| | スカラー（同梱） | SSE 経路 + archsimd | ネイティブ C（SIMD なし → あり） |
|---|---|---|---|
| 写真風 q80 m6 | 0.33〜0.35 s | 0.30〜0.31 s | 226 → 185 ms |
| 図版 q80 m6 | 0.15 s | 0.11〜0.12 s | 88 → 62 ms |
| 図版 可逆 | 38〜41 ms | 32〜41 ms | 25 → 23 ms |
| 写真風 可逆 | 0.48〜0.53 s | 0.46〜0.50 s | 305 → 248 ms |

非可逆は 1.1〜1.3 倍速くなり、可逆は同等である。最初の計測では可逆が 1.5 倍遅かったが、原因は libwebp の可逆 SSE 経路ではなく（ネイティブ C では SSE ありの方が可逆も速い）、Go 1.27 側の 2 つの挙動の組み合わせだった。gc はベクトル型のゼロ値をレガシー SSE の `MOVUPS X15, Xn` で作り、ランタイムの非同期プリエンプションは AVX-512 機でレジスタを `VMOVDQU64 Z0..Z31` で復元して `VZEROUPPER` を実行しないため、最初のプリエンプション以降はレガシー SSE 命令 1 つごとに SSE/AVX 状態遷移（実測で約 300 サイクル）が発生する。可逆の予測器はピクセルごとに `i8x16.narrow_i16x8_u` を呼び、そのエミュレーション中のゼロ値がこれに当たっていた（`GODEBUG=asyncpreemptoff=1` で 30 倍速くなることで確認）。ヘルパーからゼロ値ベクトルをなくして解消したが、gc がスピル/リロードに使う `MOVUPS` は残る（生成コード中に約 1 万箇所）。この libwebp では影響は小さかったが、AVX-512 機で `GOEXPERIMENT=simd` を使う場合は `//go:debug asyncpreemptoff=1` を検討すること。なお `i8x16.narrow_i16x8_u`（VPACKUSWB）自体は Go 1.27 の archsimd が AVX の範囲で公開しておらず 7 命令でエミュレートしているので、ここも将来の改善余地である。

同じ wasm から出る `[2]uint64` 版は 2.4 s のもので、`GOEXPERIMENT=simd` を付けない通常のビルドが 7 倍遅くなる。そのため 2 つのパッケージを同梱し、ビルドタグで切り替えている。`imgconv/internal/webpw` はスカラーの libwebp を `[2]uint64` で運ぶ版（約 5MB）、`imgconv/internal/webpwsimd` は SSE 経路の libwebp を archsimd で運ぶ版で、`goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)` のときだけコンパイルされる（`[2]uint64` 側の関数は落としてあり約 8.6MB、git 上は gzip で 1.3MB）。imgconv の `webp_scalar.go` / `webp_simd.go` が同じタグで束ね直し、`imgconv.SIMD()` がどちらが入ったかを返す。利用者は Go 1.27 で `GOEXPERIMENT=simd go build` するだけで SIMD 版になる（amd64 は実行時に AVX2 が必要で、無ければ init で panic する）。再生成は `tools/gen-codecs.sh`（スカラー）と `SIMD=1 tools/gen-codecs.sh`（SIMD）。AVIF も評価したが採用しなかった。可逆は図版で WebP の 200 倍大きく、非可逆は speed 6 で図版に効くものの 3〜4 倍遅く、生成コードが 47MB になるためである。

## 3.3 PowerPoint → BDF 変換器（converter/pptx）の構造

`converter/pptx` は .pptx（OPC の zip）を読み、DrawingML を直接 BDF の命令に落とす。Office から PDF を経由する経路と比べると、マスター・レイアウトの共有とテキストの構造（段落・箇条書き・行の継ぎ目）が残る。XML は構造体にデコードせず汎用の要素木として読む（DrawingML は省略可能な要素と多段の継承が多いため）。`mc:AlternateContent` は Fallback を使う。

- **レイヤー構成**: 各ページは `background`（スライド → レイアウト → マスターの順で最初に見つかった背景）、`master`（マスターのプレースホルダー以外の図形）、`master`（レイアウトの同様の図形）、`body`（スライドの図形）の 4 レイヤー。マスターとレイアウトのレイヤーはスライドの配色マップで描くので、同じマスター・レイアウトを使うスライドどうしでは同一の Object になり、内容アドレスで 1 回だけ格納される。`showMasterSp="0"` で隠されたものは載せず、描くもののないレイヤーは省く。
- **プレースホルダーの継承**: スライドのプレースホルダーはレイアウトのものに idx → type の順で、レイアウトのものはマスターのものに type（`ctrTitle` → `title`、`subTitle`/`obj` → `body`）で対応付け、位置（xfrm の off と ext は別々に）、形状、塗り・線・効果、`bodyPr`（属性単位）、`lstStyle`（レベル単位）を継承する。マスターとレイアウトのプレースホルダーそのものは描かない。
- **テキストの書式**: 段落の `pPr` → 図形の `lstStyle` →（図形スタイル `p:style` の `fontRef`、表スタイルの `tcTxStyle`）→ 継承元プレースホルダーの `lstStyle` → マスターの `txStyles`（タイトル／本文／その他）またはプレゼンテーションの `defaultTextStyle` の順に、属性ごとに最初に見つかったものを使う。`+mj-lt` などのテーマフォントと `schemeClr` はスライドの配色マップで解決し、色の修飾（`lumMod`/`lumOff`、`tint`/`shade` は線形 RGB で、`alpha` など）を順に適用する。
- **図形**: プリセット図形 187 種は ECMA-376 Part 1 の `presetShapeDefinitions.xml` を圧縮して埋め込み（`presets.xml.gz`、`tools/gen-presets.py` で生成、30 KB）、ガイド式を評価してパスにする。`arcTo` の角度は楕円の見かけの角度なので Canvas の媒介変数角に直す。`custGeom` も同じ評価器を通す。塗りは単色・グラデーション（線形は角度と `scaled`、パスは放射状で近似）・画像（伸縮、`srcRect` の切り抜き、タイル）・パターン（8×8 の PNG を作ってパターン Paint に）、パスごとの `lighten`/`darken`。線は幅・端・結合・破線・矢印（三角・ステルス・菱形・楕円・開いた矢印）、外側の影は SHADOW。グループは子の配置をスライド座標に畳み込み（変換行列を入れ子にしない）、回転・反転は図形ごとの TRANSFORM で表す。テキストは反転させない（上下反転は 180° 回転）。
- **テキストレイアウト**: BDF に再レイアウトはないので行分割は変換側で行う。空白の後・和文の文字間・和文と欧文の間・語中のハイフンの後で改行でき、閉じ括弧・句読点・小書きの仮名・長音の前と開き括弧の後では改行しない（禁則）。インデントと箇条書き（記号、Wingdings/Symbol の文字は Unicode に置き換え、自動番号）、タブ、行間（割合・固定）、段落前後の間隔（先頭段落の前は取らない）、左・中央・右・両端揃え（欧文は語間、和文は字間を TEXT_STYLE の letterSpacing で広げる）、上・中央・下の配置、自動調整（保存されている `fontScale`/`lnSpcReduction`）、上付き・下付き、下線・取り消し線・ハイライト、縦書き（`vert`/`vert270` は枠の回転、`eaVert` は漢字・仮名を 1 文字ずつ正立させ、句読点は縦書き用字形 U+FE10〜 に、括弧・長音・欧文は回転したまま）を扱う。行の高さは Office と同じくフォントの Windows 用メトリクス（OS/2 winAscent + winDescent）に行間を掛けたもの。
- **フォント**: `converter/internal/fontdb` がフォントディレクトリ（既定でシステムのもの、`-font-dir` で追加、`-no-system-fonts` で限定）の name テーブルから索引を作り、要求されたファミリーを実物 → 計量互換の代替（Calibri → Carlito、Arial → Liberation Sans、游ゴシック → Noto Sans CJK JP / IPAexGothic など）→ 同系統の汎用フォント（serif / sans-serif / monospace、和文は明朝／ゴシック）→ 任意のフォント の順に解決する。文字ごとに欧文は latin、和文は ea のフォントを使い、グリフがなければもう一方、さらにフォールバック列から探す（どこにもなければ警告）。計測したフォントは使った文字だけのサブセット（TrueType はグリフを詰め直し、CFF は 2 MB までならそのまま）にして埋め込むので、閲覧側でも計測と同じ字幅で描かれる。ブラウザに登録する FontFace は太さ・斜体の記述子を持たないので、太字のフェイスそのものを埋め込んだときは FONT の weight を 400 にし、フェイスがなく合成が必要なときだけ 700 や italic を指定する。`-fonts system` では埋め込まずにファミリー名（要求名、代替名、汎用名の順）で参照し、`advance` 補正で行幅を保つ。ライセンス（OS/2 fsType）が埋め込みを禁じるフォントも名前で参照する。
- **テキストの構造**: テキスト枠ごとに MARK BOX、段落と `a:br` に PARAGRAPH、折り返しに LINE（和文どうしの折り返しは区切りを入れない WRAP、spec §7.8）、箇条書きの記号の後に LINE を置く。書式が同じで分割されているだけの run（スペルチェックや言語タグ）は 1 つの FILL_TEXT にまとめる。縦書きの行は子 Object に描いて ALT_TEXT に行の文字列を持たせるので、検索では横書きと同じく 1 行 1 run になる。ハイパーリンク（外部 URL、スライドへのジャンプは `#page=N`）は LINK。
- **表**: `tblGrid` と行から格子を作り、結合セル、セルの余白・配置、塗り・罫線（セルの `tcPr` が表スタイルより優先）を描く。表スタイルは `tableStyles.xml` から全体・縞模様の行／列・先頭／末尾の行／列・角のセルの順に重ね、ファイルにない既定スタイル（Medium Style 2 - Accent 1）は内蔵する。行の高さはセルのテキストに合わせて伸ばす。
- **グラフ**: グラフパートにキャッシュされた値（`numCache`/`strCache`）から、縦棒・横棒（集合・積み上げ・100%）、折れ線、面、円・ドーナツ、散布図を描く。軸の範囲と目盛の単位は Office と同じ規則（データが 0 から遠くなければ 0 を含め、端に 5% の余白を取り、1・2・5×10^n の単位で目盛の数を図の大きさに合わせる）。系列の色はアクセント色の循環で、凡例、タイトル（単一系列なら系列名）、データラベル、数値の書式コード（%、桁区切り、小数桁）を扱う。3-D グラフは平面で描く。
- **SmartArt**: PowerPoint がデータモデルと一緒に保存している描画パート（`diagrams/drawingN.xml`）の図形を、グラフィックフレームを枠とするグループとして描く（テキストは `txXfrm` の枠に置き、その回転は図形の回転に足す）。
- **その他**: 非表示スライドは既定で除く（`-hidden` で含める）。OLE オブジェクトはプレビュー画像を描く。画像は `imgconv`（§3.2）を通し、TIFF は PNG にデコードする。構造の壊れたスライドで描画が失敗した場合は空のページにして警告する。
- **未対応（警告を出す）**: EMF/WMF 画像、レーダー・バブル・等高線グラフ、段組み（1 段で配置）、画像の色変更効果（複色・二値化など）、光彩・反射・ぼかしなどの効果、インク、旧形式（VML のみ）の OLE プレビュー。

テスト用のデッキは python-pptx で生成し（`npm run test:pptx:gen`、`test/pptx/gen.py`）、変換結果は `fixtures/pptx/` に置いて golden テストで描画を比較する。フォントは `converter/pptx/testdata/fonts` の M PLUS 1p のサブセットだけを使うので、出力は実行環境に依存しない。開発中は Apache POI のテストデータ（PowerPoint で作られた実ファイル約 90 本）でも変換を確かめ、LibreOffice の描画（PPTX → PDF → BDF）と見比べた。`lumMod`/`lumOff` と `alpha` を併用した色、グラデーションの線、縦書きなどでは LibreOffice の方が崩れる。

## 4. テキストの扱い

一番忠実度を左右する部分。3 段階を用意する。

| 方法 | 忠実度 | サイズ | 検索/選択 | 用途 |
|---|---|---|---|---|
| `FILL_TEXT` + 埋め込み WOFF2 + `advance` 補正 | 高 | 小 | 可 | 標準 |
| `FILL_TEXT` + システムフォント + `advance` 補正 | 中 | 最小 | 可 | 埋め込み不可のフォント、和文のフォールバック |
| アウトライン化（`FILL_PATH_RUN` + Path collection） | 完全 | 大 | 不可（`MARK` で文字列を別途持てば可） | ロゴ、特殊フォント、Type3 |

`advance` 補正は pdf.js が長年使っている手法で、フォント代替時の行崩れを実用上ほぼ消せる。文字単位の位置がずれても行の右端は揃う。

Canvas にはグリフ ID で描く API がないので、「サブセットフォントを埋め込んでも、文字列→グリフの対応はブラウザのシェーピングに任せる」ことになる。合字やコンテキスト依存の字形が原文と違うことはありうる。それが許容できない箇所だけアウトライン化する。

### 4.1 検索と選択

検索と選択は別の問題として扱い、検索はライブラリ（Worker）が提供する。

- **検索はデータの問題**。`FILL_TEXT` の文字列と `MARK` の境界（§7.8）から読み順の run 列が得られる。View ごとのテキスト索引 Part（§7.9）があれば Object を一切デコードせずに全文検索でき、同じ Part をサーバー側の検索エンジンにも流せる。索引がなければ Object を走査して同じ run 列を作る（Go と TypeScript の抽出器は同じ順序・通し番号を返すことをテストで保証している）。
- **正規化は検索時に行う**。NFKC、大文字小文字、カタカナ→ひらがなを 1 文字ずつ正規化し、正規化後の各コード単位から元の (run, オフセット) への対応表を持つ。合字や全角半角で長さが変わっても、ヒットを元の文字範囲に戻せる。
- **ヒットの矩形は Worker で計算する**。該当 Object をデコードして run の位置・フォント・`advance` を取り、接頭辞幅を `measureText` で測って部分一致の矩形を出す。フォントは Worker に読み込まれているので、埋め込みフォントでも正しく測れる。行をまたぐヒットは複数の矩形になる。
- **選択・コピーは透明 DOM**（pdf.js 方式）。`@bdf/render` の `buildTextLayer(runs, scale)` が run ごとに透明な `span` を絶対配置した層を作り、ページ型では可視ページだけに置く。span の幅は Worker が埋め込みフォントで測った `advance` に合わせて `scaleX` で伸縮するので、メインスレッドにフォントが無くても選択範囲が描画と一致する。回転・斜体の run は `matrix` をそのまま CSS transform に渡す。連続モードは `continuousText(view, viewport)` が帯座標に移した run を返し、同じ層を帯ごとに置く。シートは DOM にせず、セル範囲の選択モデルをビューアが持つ方が自然（10 万セルの span は重すぎる）。
- **コピーはブラウザの直列化に任せない**。無関係な span を選択したときのブラウザのテキスト化は区切りが落ちる（単語がくっつく）ので、`installCopyHandler(container)` が copy イベントを取り、選択範囲に交差する span を文書順に集めて `selectionText()` で組み立てる。run の `sep`（MARK 由来、無ければ位置からの推定）を使い、LINE は空白、PARAGRAPH/CELL/BOX は改行、ページをまたいだら空行にする。span の途中で始まる・終わる選択は Range のオフセットで切る。
- ビューアの責務はヒット一覧の表示とハイライトの描画だけで、フォーマットの内部事情（run の分断、MARK、正規化）を知らなくてよい。

## 5. 共有オブジェクトの効き方

内容アドレスなので、エンコーダが「同じものを同じバイト列で出す」だけで共有される。効く条件は決定的な出力:

- f32 の丸め方を固定（Go 側で `float32()` に落としてから書く）。
- 文字列表・パス表の順序は出現順に固定。
- Object の bbox はコンテンツから計算し、ページ位置に依存させない（ローカル座標）。
- 「ほぼ同じだが少し違う」ものは共有されない。マスター＋差分という構造は変換元の情報がある場合（PPTX 直接変換）にのみ作れる。

さらに読み手側では `USE` 対象を `(hash, scale)` でビットマップキャッシュできるので、共有はサイズだけでなくスクロール時の描画コストも下げる。

## 6. Excel シートの Tile 化

- Tile サイズは 2048 unit（≈ 28 インチ）を既定にする。標準行高 20pt なら 100 行、標準列幅 64pt なら 32 列がおおむね 1 Tile。
- 空 Tile は表に載せないので、疎なシートは小さくなる。
- 座標を Tile 原点相対にする理由: Excel の最大行数 1,048,576 × 20pt ≈ 2,000 万 unit で f32 の整数精度（約 1,677 万）を超える。
- 行列ヘッダー・グリッド線・固定ペインはビューアが `cols`/`rows`/`freeze` から描く。内容にグリッド線を焼き込まないのは、表示切替とズーム時の線幅制御のため。
- 1 Tile を 1 Object にして「クリップして描くだけ」にしたのは、レンダラの Tile キャッシュを単純にするため。またがるオブジェクトの重複は `USE` 参照の重複（16 バイト + 数バイト）だけなので許容。

## 7. Word の「縦一枚」表示

再レイアウトはしないと決めた。理由は BDF が表示リストであること、そして再レイアウトを許すとフォーマットが「文書モデル」になって PDF との差別化（軽いデコーダ）が消えること。

代わりに 2 つの手段を用意する。

1. `flow` View の連続モード: `body` レイヤーを本文矩形で切り出して積む。ページ境界で段落が割れるのは残るが、ヘッダー・フッター・余白は消える。多くのプレビュー用途はこれで足りる。
2. 変換側で「用紙高さ無限」の設定でもう一度レイアウトし、別 View として同梱する。画像・フォント・繰り返し部品は共有されるので、増えるのは本文の命令列だけ。

## 8. ロードマップ

1. **仕様固め**: `spec.md` の opset 1 を確定。Go のエンコーダ + TS のデコーダを最小実装し、手書きの Object で Canvas に描けるところまで。
2. **フィクスチャと golden テスト**: Go でフィクスチャ生成 → Playwright でスクリーンショット比較。
3. **PDF → BDF**: 最初の実用変換（実装済み、§3.1）。
4. **ビューア**: Worker + OffscreenCanvas、ページ/連続/シートの 3 モード、テキストレイヤー、検索。
5. **XLSX → BDF**: 直接変換、Tile 化、固定ペイン。
6. **PPTX 直接変換**: マスター共有の本領（実装済み、§3.3）。

## 9. リポジトリ構成（案）

```
bdf/
├── docs/              spec.md, design.md
├── *.go               Go パッケージ bdf（module github.com/shibukawa/bdf）: Object builder、Part エンコード、コンテナ I/O、デコーダ
├── cmd/bdf/           CLI: generate / ls / manifest / disasm / extract / split / join / demo
├── imgconv/           画像の格納方針と WebP/AVIF 変換（internal/ は wasm2go で生成した純 Go コーデック）
├── converter/         変換器の共通部分（入力形式の判別、ページ指定）
│   ├── pdf/           PDF → BDF 変換器（testdata/ にテスト用 PDF）
│   ├── pptx/          PowerPoint → BDF 変換器（testdata/ にテスト用デッキとフォント）
│   └── internal/      fontdb（フォントの探索・解決・計測・サブセット）、sfnt（TrueType/OpenType の読み書き）
├── fixture/           フィクスチャ生成（埋め込みフォント、計測、サンプル文書）
├── packages/
│   ├── core/          @bdf/core  デコーダ・コンテナ読み込み・テキスト抽出（依存なし）
│   └── render/        @bdf/render Canvas バックエンド、ページ/連続/シート描画、Worker とクライアント
├── examples/viewer/   デモビューア（Worker 描画、テキストレイヤー）
├── fixtures/          Go が生成した demo.bdf / demo-split と golden PNG
└── test/              Playwright による golden テスト
```

将来の変換器（XLSX など）も `converter/` のサブパッケージとして追加し、`bdf generate` から形式を判別して呼ぶ。
