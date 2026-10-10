# CPU 効率と安全性の再調査

調査日: 2026-10-08。[前回の調査](security-performance-review.ja.md)と[通常処理の改善](performance-review.ja.md)の後に残っているものを探した。対象は Go の BDF コンテナ（読み書き・暗号化・区間）、ラスタライザ、変換器と共通部品、サーバー例、TypeScript の core/render/viewer と依存パッケージ。

コードを読むほかに、入力全体を通すファザーを 2 本書いて計 20 分ほど回し、CPU とメモリ確保のプロファイルから改善点を探した。前回と同じく、あらゆる入力で問題が起きないことを保証する監査ではない。

## 見つかった問題と修正

| 優先度 | 問題と影響 | 修正 |
|---|---|---|
| 中（サーバー利用時） | HTML/Markdown のネットワーク画像の取得は、ライブラリの既定では接続先を制限しない。変換するサーバーの内側（ループバック、プライベートネットワーク、クラウドのメタデータサービスが答えるリンクローカル）に、名前解決やリダイレクトの先としても届く | 既定の HTTP クライアントは、接続する直前の IP アドレスが公開アドレスでなければ接続を拒否する（`Dialer.Control`。名前解決後・リダイレクト後も同じ検査を通る）。`Options.AllowPrivate` / `-param private=true` で許可。環境変数のプロキシ設定があるときはプロキシに任せる。TinyGo にはダイヤラがないのでブラウザの規則に任せる |
| 中（サーバー利用時） | OOXML の部品は 1 GiB まで展開して木にする。要素 1 つが数百バイトの木になるので、`<a/>` を詰めた 1 MiB の pptx が数十 GiB を要求できる | 木の要素数に上限を付けた。1 つの木は 8M 要素まで（`xmltree.MaxElements`）、パッケージが同時に保持する木は 16M 要素まで（`DiscardXML` で返却）。ストリーミングで読む行や段落も 1 要素あたり同じ上限 |
| 低 | サンプルサーバー 4 つがタイムアウトなしの `http.ListenAndServe` で待つ。ヘッダーを送らない接続が溜まる | `ReadHeaderTimeout` 10 秒、`ReadTimeout` 1 分、`IdleTimeout` 2 分の `http.Server` に変更（Range の送信を切らないよう `WriteTimeout` は付けない） |
| 低 | secure-reader がクライアント由来の URL パスをそのままログに書く | `%q` で引用 |

依存パッケージ: `govulncheck`（最新、2026-10-08）と `npm audit` はどちらも 0 件。`gosec` の G115（整数変換）と G304（CLI の引数のファイル）はこの用途では問題ではなく、G602 の配列添字 2 件は入力側で範囲が確かめられていた。`staticcheck` の未使用値 2 件（dxf/hatch.go、hpgl/gl.go）は動作に影響しない。

## 調べて変えなかったこと

- BDF の部品の展開上限は manifest が書く `size` で、deflate の比率（最大 1032 倍）までは入力より大きく膨らむ。Go と TypeScript の両方とも設計どおりで、信頼しない .bdf をそのまま受け取るサービスは入力サイズに上限を置く。
- デモビューアの `?src=` はどのオリジンの文書でも開く（CORS が許す範囲）。描くのはキャンバスで、LINK は http/https/mailto と内部リンクだけ、FILTER はフィルタ関数だけなので、スクリプトの注入経路はない。
- light-server の `http.FileServer` はディレクトリ一覧を返すが、並ぶのは最初から一覧に出すファイル。
- 暗号化（AES-256-GCM、AES-KW、PBKDF2、ECDH + HKDF）と区間の実装は仕様どおりで、部品の AAD にハッシュを使うため部品の入れ替えもできない。
- ビューアの requestAnimationFrame と setInterval はすべて操作中・再生中だけ動き、待機中に回り続けるループはない。

## CPU とメモリ確保の改善

### ラスタライザ（raster/imagebdf）

プロファイルでは、描画ごとのキャンバス確保（16 バイト/画素）を OS に返して取り直す `madvise` と、走査線の交点のソートが目立った。

- キャンバスを Renderer が 64 MiB まで保持して次の描画に使い回す（グループ・ソフトマスクの一時キャンバスも同じ）。
- クリップ矩形の積: 一方が全面の矩形でもう一方を含むときは、積を作らずその一方をそのまま使う。
- 交点のソート: 前の走査線の順に並んだ辺は挿入ソートが速い。16 個を超える走査線も、1 交点あたり 4 回までの移動で済むなら挿入ソートで済ませ、超えたら pdqsort に切り替える（順序は変わらないので画素は同じ）。
- 塗り: 不透明のクリップなし・source-over という最も多い場合の内側のループから、画素ごとの分岐を外した。浮動小数点の式は元と同じ形にして、`TestShortcutsChangeNoPixel` と `TestFillWithinTheCoverage`（ビット単位の比較）を通す。

`BenchmarkRenderPage`（Apple M3、`-cpu 1`、1 秒 × 3 回の中央値）:

| 文書 | 描画 | 修正前 | 修正後 | 時間 | 確保量 |
|---|---|---:|---:|---:|---:|
| PDF `chrome-doc` | 通常ページ | 15.5 ms | 13.2 ms | -15% | 23.6 MB → 11.7 MB |
| PDF `chrome-doc` | サムネイル | 4.82 ms | 4.66 ms | -3% | 1.84 MB → 1.08 MB |
| PowerPoint `features` | 通常ページ | 8.25 ms | 6.2 ms | -25% | 11.1 MB → 2.7 MB |
| PowerPoint `features` | サムネイル | 1.07 ms | 0.94 ms | -12% | 0.72 MB → 0.20 MB |
| draw.io `showcase` | 通常ページ | 2.52 ms | 2.22 ms | -12% | 2.7 MB → 1.0 MB |
| draw.io `showcase` | サムネイル | 0.71 ms | 0.67 ms | -6% | 0.30 MB → 0.18 MB |

### 変換

変換の CPU 時間は、変換器そのものより GC とメモリ管理（`madvise`、GC ワーカーの起動）が大きい。確保量の 47% は WOFF2 の Brotli エンコーダがフォントごとに作る 2 MiB 余りの検索テーブルで、ほかに `fontdb` が変換のたびに作業ディレクトリを `stat` していた。

- Brotli の writer をウィンドウ幅ごとに `sync.Pool` で使い回す（`Reset` で検索テーブルは残り、初期化し直される）。出力は新規の writer と同じバイト列になることをテストで確かめた（`TestEncodeIsDeterministic`）。
- `fontdb`: 一度走査したディレクトリは、渡された文字列からそのまま引く（相対パスの絶対化で `os.Getwd` を呼んでいた）。
- CSV/Excel の値の判別: コロンのない値を `:` で分けない、英字と数字のない値を月名の日付として分けない。判定結果は同じ。

`BenchmarkConvert`（新設。converter/all で各形式のテスト文書を丸ごと変換して書き出す。1 秒 × 3 回の中央値）:

| 入力 | 修正前 | 修正後 | 時間 | 確保量 |
|---|---:|---:|---:|---:|
| PDF `chrome-doc` | 6.30 ms | 5.05 ms | -20% | 16.7 MB → 10.5 MB |
| PowerPoint `features` | 4.05 ms | 3.77 ms | -7% | 7.7 MB → 6.8 MB |
| Excel `basic` | 5.92 ms | 5.82 ms | -2% | 10.5 MB → 10.3 MB |
| Word `basic` | 5.67 ms | 5.37 ms | -5% | 10.2 MB → 9.9 MB |
| EPUB `basic` | 1.86 ms | 1.73 ms | -7% | 4.0 MB → 3.2 MB |
| draw.io `showcase` | 1.31 ms | 1.10 ms | -16% | 4.05 MB → 1.11 MB |
| Markdown `basic` | 1.56 ms | 1.50 ms | -4% | 1.4 MB → 1.3 MB |

残っている大きなものは wordproc の行の配置（行ごとの要素のコピー）で、構造の変更になるので見送った。参考値として、`GOGC=400` で変換は 8〜13% 速くなるが常駐メモリが増えるので、プロセスの設定として使う側が決める。

### Office の XML の読み手を xmlro に（2026-10-09）

大きな文書では encoding/xml が確保するオブジェクトの 39%、バイトの 17% を占めていた（4 シート × 3,000 行の xlsx と 100 スライドの pptx のプロファイル）。トークンごとに名前の文字列 2 つ、属性のスライス、文字データのコピーを作るためで、木にする `xmltree` の確保はその上に乗る。

- `internal/xmltree` の木は tinygodriver の `encoding/xmlro` で読む。xmlro はトークンに何も確保せず、名前や値はバッファへのスライスなので、木が持つものだけを写す。要素名と短い属性値（32 バイトまで）は部品ごとに intern して 1 つの文字列にし、名前空間は encoding/xml と同じに解決する（接頭辞は URI に、束縛のない接頭辞はそのまま、`xmlns:` の宣言は Space "xmlns" の属性）。`ooxml.Package` は読み手を 1 つ使い回し、最大のトークンまで育ったバッファを部品から部品へ持ち越す。
- xlsx の `sheetData` と共有文字列、pptx の非表示スライドの判定、docx の判別は、`ooxml.NewReader` で流し読みし、持つ要素だけ `Package.ReadFrom` で木にする。
- `TestTreesOfDocumentsMatch` が、リポジトリの Office 文書の全 XML 部品（12 文書、258 部品）と境界条件の文書で、encoding/xml から作った木との一致を確かめる（encoding/xml から木を作る関数はテストの中だけにある）。SVG、draw.io、MusicXML、EPUB、XHTML、XMP はこの時点では encoding/xml のまま（xmlro v1.3.3 は DOCTYPE の内部サブセットを拒み、`Strict = false` の寛容さがなかった）で、次の節で置き換えた。

合成した Word 本文（4,000 段落、3.7 MB）の木:

| | 時間 | 確保量 | 確保回数 |
|---|---:|---:|---:|
| encoding/xml（従来） | 57 ms | 44.5 MB | 1,061k |
| xmlro | 17 ms | 21.4 MB | 311k |

変換全体（`BenchmarkConvertFiles` は `tools/generate-memory-fixtures.py` の大きな文書、3 回の中央値）:

| 入力 | 修正前 | 修正後 | 時間 | 確保量 | 確保回数 |
|---|---:|---:|---:|---:|---:|
| 4 シート × 3,000 行の xlsx | 491 ms | 478 ms | -3% | 339 MB → 298 MB | 4.38M → 3.07M |
| 100 スライドの pptx | 131 ms | 99 ms | -24% | 139 MB → 116 MB | 1.74M → 1.01M |
| PowerPoint `features` | 3.77 ms | 3.07 ms | -19% | 6.8 MB → 6.2 MB | 34.8k → 20.0k |
| Excel `basic` | 5.82 ms | 5.27 ms | -9% | 10.3 MB → 10.0 MB | 49.6k → 32.6k |
| Word `basic` | 5.37 ms | 5.30 ms | -1% | 9.9 MB → 9.8 MB | 27.2k → 20.8k |

xlsx の変換時間の残りは、セルのレイアウト（確保量の 2 割）と出力の deflate 圧縮（CPU の 3 分の 1）が占める。

### 残りの XML の読み手も xmlro に（2026-10-09、tinygodriver v1.3.4）

xmlro が v1.3.4 で HTML 風の XML を読めるようになった（`Lenient` は encoding/xml の `Strict = false`、`htmlentity` の実体と空要素、`CharsetReader`、DOCTYPE の内部サブセットの実体、BOM 付きの UTF-16）。encoding/xml で読んでいた残りをすべて置き換え、変換器のコードに encoding/xml はなくなった。残っているのは、XML を書く開発用の `tools/gen-drawio-stencils`（xmlro は読むだけ）、encoding/xml を正解として比べるテスト、依存先の pdfcpu である（pdfcpu があるのでバイナリにはまだリンクされる）。

- **開き方をひとつに**（`xmltree.Open`）: SVG（`converter/image`、`image/imgconv`、`raster/imagebdf`）、draw.io のファイルとステンシル、EPUB の package・container・encryption、MusicXML、XHTML（`webdoc`）、XMP は、読み手のオプションだけを決めて `xmltree.Open` で開く。名前空間を解決した名前は `ElementName` / `Attrs` / `AttrName`、接頭辞を除いた名前は `Local`、文字データ（テキストと CDATA）は `AppendText` で取る。`xmltree.Node` の属性は `xml.Attr` から `xmltree.Attr`（`ooxml.Attr`）になった。EPUB の 3 つのファイルは構造体へのデコード（リフレクション）をやめて小さな木にし、Agile 暗号の EncryptionInfo（`offcrypto`）は `xmltree.Parse` の木から読む。
- **実体の展開に上限**: v1.3.4 から xmlro は DOCTYPE が宣言した実体を置き換える。上限はトークンごと（`MaxBufferBytes`）で、文書全体にはない。50 KiB の実体を参照するだけのテキストを何千も並べれば、小さな文書が何百 MB もの文字列になる。`Open` は読む前に宣言と参照を数え（`BoundEntities`）、置き換えた合計が 1 MiB を超える文書は宣言を無効にして参照をそのまま残す（SVG が前からしていたことで、警告もそのまま出る）。数えるのは文字コードを変換した後の文書で、UTF-16 や変換される文字コードでも同じにかかる。
- **Office の部品は実体の宣言を拒む**（`xmltree.ErrEntities`）: Office の XML は実体を宣言しない。部品はストリームで読むので前もって数えられず、トークンの上限は図のデータのために 256 MiB ある。`<!ENTITY` を含む部品は、流し読みでも木でも読まない（v1.3.3 の xmlro が内部サブセットを拒んでいたのと同じ結果）。
- **文字コードの変換は 1 回**: encoding/xml も xmlro も、文字コードを書いた XML 宣言が現れるたびに残りを変換し直す。宣言を何万も並べた文書で残りが何万回も写されるので、`Open` は最初の変換のときに残りを全部変換し、後ろの宣言を別の名前の処理命令にする（コメントと CDATA の中は本文なので触らない）。それでも 2 度目の変換を求める文書は、そこで読むのをやめる。寛容な読み手が返すものは必ず UTF-8 で、そうでないバイトは U+FFFD にする（encoding/xml はそこで読むのをやめていた）。
- **XHTML の整形式の検査**: XHTML は厳密に読み、整形式でなければ HTML のパーサーに回す。xmlro は終了タグの対応と属性の書き方しか確かめないので、`webdoc` が残りを確かめる（UTF-8、制御文字、XHTML にない実体と単独の `&`、属性値の `<`、テキストの `]]>`、コメントの `--`）。HTML5 にしかない実体（`&check;`）やセミコロンのない `&nbsp` を、これまでどおり HTML のパーサーが読む。

**検証**: `TestOpenMatchesEncodingXML` が、名前空間、未束縛の接頭辞、HTML の実体と空要素、対応しない終了タグ、値のない属性、Latin-1 などの文書で、encoding/xml（`Strict = false`、`HTMLEntity`、`HTMLAutoClose`）と同じトークン列になることを確かめる。変更前のバイナリとの比較では、リポジトリの 406 ファイル（draw.io、EPUB、MusicXML、SVG、HTML、Markdown、Office、IDML、画像、PDF。うち 4 つは前後とも同じエラー）、手元の実ファイル 50（EPUB 3 冊、draw.io 28、Illustrator 12 など）、手元の SVG 451 の変換結果が警告も含めてすべてバイト一致し、SVG 451 のサムネイル（`raster/imagebdf` の読み手）も一致した。`FuzzConvert` を 4 分と 2 分半走らせてパニックはなく、1 秒を超えた入力は XML を読まない PDF がひとつ（ほかの計測と重なったときで、単独では変更前と同じ 0.6 秒）。

変換全体（`BenchmarkConvertFiles`、Apple M3、5 回の中央値）と読み手だけの時間:

| 入力 | 修正前 | 修正後 | 時間 | 確保回数 |
|---|---:|---:|---:|---:|
| 326 KB の SVG（画像として変換） | 3.04 ms | 1.14 ms | -63% | 11.5k → 8.8k |
| 同じ SVG を `raster/imagebdf` が読む | 2.15 ms | 0.38 ms | -82% | 4.3k → 2.5k |
| draw.io `showcase` | 0.86 ms | 0.73 ms | -15% | 7.7k → 6.6k |
| draw.io `aws` | 1.38 ms | 1.26 ms | -8% | 7.3k → 6.4k |
| EPUB `basic` | 1.25 ms | 1.15 ms | -7% | 9.6k → 8.1k |
| 実物の EPUB（約 360 ページ） | 138 ms | 135 ms | -2% | 609k → 535k |
| 541 小節の MusicXML（527 KB） | 159 ms | 151 ms | -5% | 955k → 788k |
| draw.io の最大のステンシル集（aws4）を読む | 44.2 ms | 42.3 ms | 差なし | 確保量 50.7 MB → 36.9 MB |

EPUB と MusicXML は組版と浄書が時間のほとんどを占め、XML を読む時間は小さい。

npm の変換器（TinyGo 0.42.0）も、変更前後のビルドでリポジトリの XML を含む 204 ファイルの変換結果がバイト一致した。大きさはほぼ変わらない。pdfcpu が encoding/xml を使い続けるので、All は 14.70 MB → 14.75 MB（+53 KB、gzip で +17 KB）。PDF と Illustrator を含まない組み合わせ（Office、IDML、EPUB、draw.io、画像、HTML、Markdown、MusicXML、CSV）では encoding/xml がなくなり、9.93 MB → 9.89 MB（-40 KB、gzip で -12 KB）。

### deflate を klauspost/compress に（2026-10-09）

BDF の部品の圧縮と展開（document.go、reader.go、区間の manifest）、zip の部品の展開（`converter/internal/ziputil`。Office・EPUB・KiCad・Gerber・SXF・MusicXML の読み手が使い、展開器はプールする）、draw.io の deflate、PSD・TIFF・画像・WOFF の zlib、CGM と Parquet の gzip を、標準ライブラリから `github.com/klauspost/compress` に替えた（go.mod にあった依存）。埋め込みテーブルの gzip（起動時に 1 回）は標準ライブラリのまま。

全テスト文書の部品 366 個（2.1 MB）と、大きな xlsx・pptx を変換した部品（8.4 MB）で、klauspost の deflate は **どのレベルでも標準ライブラリと同じバイト列を書く**ことを確かめた（`TestFlateMatchesStandardLibrary` が毎回確かめる）。既定レベル（-1）だけは別のレベルを指すので、`compress` は -1 を標準ライブラリの 6 に写す。fixture は変わらない。TinyGo でもビルドでき、wasm モジュールは gzip 後で 29 KB 増える。

同じ部品での圧縮・展開（5 回の平均）:

| | 時間 | 速度 | 大きさ |
|---|---:|---:|---:|
| 標準ライブラリ レベル 9（既定） | 69.7 ms | 30 MB/s | 839,292（元の 39.8%） |
| klauspost レベル 9 | 60.9 ms | 35 MB/s | 同じ |
| klauspost レベル 7 | 29.5 ms | 72 MB/s | 841,558（+0.3%） |
| klauspost レベル 6 | 16.8 ms | 126 MB/s | 847,109（+0.9%） |
| 標準ライブラリの展開 | 10.9 ms | 194 MB/s | |
| klauspost の展開 | 8.3 ms | 253 MB/s | |

変換全体（3 回の中央値）: 文書の読み込みは 20〜30% 短縮（`BenchmarkDocumentIO` の read。確保回数は半分以下）、書き出しは 5〜13%、変換は PowerPoint `features` で 8%、Word と EPUB で 5〜6%、100 スライドの pptx で 4%、ほかは 0〜2%。同じレベルでは、圧縮器の違いより圧縮レベルの違いが大きい: レベル 6 なら圧縮が 4 倍速く大きさは 0.9% 増、レベル 7 なら 2.4 倍速く 0.3% 増。

この結果から、`Document.CompressionLevel` の既定をレベル 9 から 6 に変えた。`bdf generate -compression best` がレベル 9（`fast` はレベル 1）。split・join・encrypt が書き直す manifest と区間の manifest も 6。ブラウザの変換器は従来どおりレベル 1。testdata の fixture は amd64 ビルドで作り直した（部品のハッシュは展開後の内容のものなので変わらず、変わるのは圧縮後の長さとオフセット）。

## ファジング

- `FuzzReader`（ルート）: OpenSingle → ToDocument → SearchText → サムネイル → 先頭 2 ページの描画 → 区間 → 書き出し。testdata の .bdf を種にする。
- `FuzzConvert`（converter/all）: 登録済みの全形式を内容とファイル名で判別して変換し、書き出す。各変換器の testdata を種にする。
- `BDF_FUZZ_SLOW_DIR` を設定すると、1 秒より長くかかった入力を書き出す。

4〜6 分ずつ 2 回回して、パニック・ハング・1 秒を超える入力はなかった（進捗が止まって見えるのは、新しい入力を 60 秒かけて最小化しているとき）。

## 検証

- `go test ./...`、`go vet -unreachable=false ./...`、`go test -tags bdf_noconv ./image/imgconv/ ./converter/... ./font/woff2/`、raster/ebitenginebdf の vet とテストが成功。
- `go test -race` を woff2、fontdb、html、imagebdf、xmltree、ooxml、xlsx、pptx、docx、ziputil と zip を読む変換器、ルートで実行して成功（プールと走査キャッシュの並行利用）。
- `npm test` が成功。
- `bdf generate` で pptx・docx・xlsx・pdf・csv・parquet・drawio・epub・markdown・html の testdata と `bdf demo` を再生成し、既存の BDF とバイト単位で一致（docx/math.bdf だけは修正前のコードでも一致せず、arm64 と amd64 の浮動小数点の差）。
- `go vet -tags tinygo ./converter/html/` が通り、TinyGo 用のクライアントが選ばれる。
