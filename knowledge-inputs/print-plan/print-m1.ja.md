# 印刷 実装案（print-m1）

2026-10-19 時点の実装案。要件と決定は `print-requirements.ja.md` と `.knowledge/` の `requirement:print`、`system:ipp-driverless`、`system:vector-pdls`、`milestone:print-m1` にある。ここは「何をどの順で作るか」を書く。利用者は 2026-10-19 に「print-m1 を最初に」と決めた。

## 0. 決まっていること（と、暫定で進めること）

決まっていること。

- 出力先の第 1 は IPP 直結（IPP/2.0 の操作、`ipps`、DNS-SD）。macOS / Linux の CUPS のキューも同じクライアントで。Windows は IPP 直結だけ。
- 文書形式の第 1 は PDF（ベクタのまま。`pdfbdf`）、プリンタが PDF を受けないときだけ PWG Raster（`imagebdf` の帯から流す）。URF / PCLm は m3。
- チケットの名前は IPP の属性。ジョブの状態は IPP の `job-state`。
- cgo なし。wasm / TinyGo のビルドに入れない。
- m1 の範囲は「手元／LAN の 1 台に `bdf print` で印刷し、完了まで待つ」。スプール・サーバーは m2。

暫定で推奨どおりに進めること（`requirement:print` の open_questions。異論があれば変える。どれも後から差し替えられる位置に置く）。

- IPP の符号化は自前（`print/ipp`。RFC 8010 は小さい。`goipp` はテストの参照に使うだけで依存にしない）。
- DNS-SD は自前の最小実装（mDNS の問い合わせと PTR / SRV / TXT / A / AAAA の解釈。`golang.org/x/net/dns/dnsmessage` は既に依存にある）。
- パッケージ名は `print/`。

## 1. 成果物

次のコマンドが動くこと。

```
bdf printers                                   # LAN の IPP プリンタを DNS-SD で列挙
bdf printers -url ipp://192.0.2.10/ipp/print    # 1 台の能力（用紙、両面、色、解像度、形式、仕上げ、認証）
bdf print -printer ipp://192.0.2.10/ipp/print invoice.bdf            # PDF で送り、完了まで待つ
bdf print -printer ipps://mfp.example.local/ipp/print -sides two-sided-long-edge \
          -media jis_b5_182x257mm -copies 3 -mono -pages 1-3 invoice.bdf
bdf print -printer ipp://localhost:631/printers/office -format pwg-raster -resolution 300 invoice.bdf   # CUPS のキュー、Raster を明示
bdf print -printer ipp://192.0.2.10/ipp/print -skip background slip.bdf      # プレ印刷用紙（用紙の層を抜く）
bdf print -printer ipp://192.0.2.10/ipp/print report.xlsx                    # converter の入力はその場で変換
```

`bdf print` は、プリンタの能力を取り、チケットを IPP の属性に写し（受けない属性は警告か `-strict` でエラー）、形式を選び、`Create-Job` → `Send-Document` で描きながら送り、`job-state` が `completed` / `canceled` / `aborted` になるまで `Get-Job-Attributes` で待って、結果に応じた終了コードで終わる。Ctrl-C で `Cancel-Job`。

## 2. 構成

```
print/
  ticket.go        Ticket（部数、丁合、両面、用紙、トレイ、種類、色、品質、解像度、ページ範囲、合わせ方、向き、仕上げ、名前、利用者、保留、優先、形式、
                   SkipRoles、ビュー、モノクロの方式）と IPP の属性への写し（*-supported との照合、unsupported の列挙）           (S4)
  media.go         PWG 5101.1 の用紙の名前の表（A 系、B 系（JIS）、Letter / Legal、はがき、封筒 長形 3 / 4、角形 2、custom_）と
                   ページの大きさ（pt）からの選択、±2 mm の許容                                                                 (S4)
  select.go        形式の選択（明示 > PDF（受けるなら）> PWG Raster > URF / PCLm（m3 までは「未対応」のエラー））と、
                   Raster の解像度・型（srgb_8 / sgray_8 / black_1）・裏面の向きの決定                                           (S4)
  doc.go           文書の用意: .bdf を開く（暗号化は $BDF_PASSWORD）、converter の入力は変換、ビューとページの選択              (S4)
  print.go         1 ジョブの流れ（Printer を取る → 検証 → 形式 → Create-Job → Send-Document → 待つ）と Job の状態、終了の理由  (S4, S5)
  errors.go        ErrConnect、ErrAuth、ErrCertificateChanged、ErrUnsupported（属性の列）、ErrFormat、ErrPrinterStopped、
                   ErrAborted、ErrCanceled                                                                                      (S5)
print/ipp/
  message.go       Message、Group、Attribute、Value（RFC 8010 の全ての型と out-of-band）、Encode / Decode、限界            (S1)
  tags.go          value-tag と operation-id、status-code、job-state、printer-state の定数と名前                             (S1)
  attrs.go         取り出しの補助（Keywords、Int、Bool、Resolution、Range、Collection、"media-col" の読み書き）               (S1)
  client.go        Client（HTTP、URL の解釈 ipp / ipps → http / https と 631、Basic 認証、TLS の方針）、
                   Get-Printer-Attributes、Validate-Job、Create-Job、Send-Document（流しながら）、Print-Job、
                   Get-Job-Attributes、Get-Jobs、Cancel-Job、Identify-Printer、CUPS-Get-Printers                              (S2)
  printer.go       Printer: 能力の構造体（URI、形式、用紙、トレイ、種類、両面、色、解像度、部数、pwg-raster-*、urf-supported、
                   仕上げ、print-scaling、認証、状態と理由、printer-config-change-time）と Get-Printer-Attributes からの読み取り (S2)
  trust.go         ipps の証明書の方針: システムの CA → 初回の指紋の記録（TOFU）、-insecure。記録は
                   $XDG_CONFIG_HOME/bdf/printers.json（os.UserConfigDir）                                                      (S5)
  dnssd.go         mDNS の問い合わせ（_ipp._tcp / _ipps._tcp、QU ビットで一時ポートに単一応答を受ける）、PTR / SRV / TXT / A / AAAA、
                   TXT の rp / pdl / URF / TLS / Color / Duplex / UUID / adminurl / ty / note / kind / air / priority → URL      (S2)
  ipptest/         偽の IPP プリンタ（http.Handler）: 能力を設定でき、受け取った文書を保存し、job-state の遷移を再現する。
                   cmd/bdf と print のテストが使う                                                                             (S2)
  testdata/        RFC 8010 の例と、実機から採った Get-Printer-Attributes の応答（個体を特定する属性は伏せる）                   (S1, S2)
print/raster/
  pwg.go           PWG Raster の書き手: 同期語 RaS2、ページごとの 1796 バイトのヘッダー、行の繰り返しと画素の連の RLE。
                   imagebdf の帯（RegionRows）を受けて io.Writer に流す。1 / 8 / 24 ビット                                   (S3)
  pwg_read.go      読み手（テスト用。ヘッダーと RLE を戻す）                                                                   (S3)
  place.go         ページを用紙に置く: fit / actual / shrink、自動回転、中央寄せ → imagebdf の領域と回転                        (S3)
raster/imagebdf/
  rows.go          RegionRows（領域 × 画素数を帯で描く。PageRows はその特殊形）と Options.Rotate（0 / 90 / 180 / 270。
                   出力の行が帯になるので帯の描画がそのまま使える）                                                            (S3)
vector/pdfbdf/
  pdfbdf.go        Options.Pages（書くページの選択。page-ranges を受けないプリンタのため）                                     (S4)
cmd/bdf/
  print.go         bdf print、bdf printers                                                                                      (S4, S5)
  main.go          usage の追記
.github/workflows/ci.yml
                   print ジョブ: ippeveprinter（cups-ipp-utils）と CUPS のキューに対して bdf print、受け取ったファイルの検査。
                   Windows / macOS の runner で go vet と go test ./print/...（ネットワークなし）                              (S5)
docs/
  architecture/cli.ja.md（bdf print / printers）、architecture/layout.ja.md（print/）、design.md §3.36                       (S5)
```

`print` は Go の組み込み関数と同名だが、パッケージ名が組み込みを隠すだけで害はない（`print.New` と書く）。

## 3. 各部の設計

### 3.1 IPP の符号化（`print/ipp`、S1）

RFC 8010 のまま。

- `Message{Version [2]uint8; Op uint16（要求では operation-id、応答では status-code）; RequestID uint32; Groups []Group}`、`Group{Tag Tag; Attrs []Attribute}`、`Attribute{Name string; Values []Value}`、`Value{Tag Tag; …}`。`Value` は型つきの union で、integer / boolean / enum は `int32`、octetString / 文字列系（keyword、name、text、uri、uriScheme、charset、naturalLanguage、mimeMediaType）は `string`、textWithLanguage / nameWithLanguage は `{Lang, Text}`、dateTime は `time.Time`、resolution は `{X, Y, Units}`、rangeOfInteger は `{Lo, Hi}`、collection は `[]Attribute`、out-of-band（unsupported / unknown / no-value / not-settable / delete-attribute / admin-define）は値なし。
- `Encode(w io.Writer) error` と `Decode(r io.Reader, lim *Limits) (*Message, error)`。復号の限界: 全体 4 MiB、属性 10,000 個、1 つの値 64 KiB、collection の深さ 8、1 属性の値 1,000 個。超えたら `ErrLimit`。
- 1setOf は name-length 0 の続きの値、collection は begCollection / memberAttrName / endCollection。memberAttrName の値が名前、次の値が中身。
- 文字列は UTF-8 のまま持つ（`attributes-charset` は `utf-8` しか送らず、応答が別の charset を名乗ったら警告して bytes のまま）。
- `FuzzDecode`（panic しない、復号 → 符号化 → 復号で同じ）。`goipp` で同じバイト列を復号してテストで突き合わせる（テストの依存は `go.mod` の `// indirect` にもしない: 実機の応答を `goipp` で一度 JSON にしておき、その JSON を testdata に置く）。

### 3.2 クライアントと能力（`print/ipp`、S2）

- URL: `ipp://host[:631]/path` → `http://host:631/path`、`ipps://` → `https://`（既定 631）。`http(s)://` もそのまま受ける。
- 要求の共通属性: `attributes-charset=utf-8`、`attributes-natural-language=ja`（`-lang` で変更）、`printer-uri`、`requesting-user-name`（`-user`、既定は OS の利用者名）。
- 送信は `http.Request` の `Body` に `io.Pipe` をつなぎ、IPP のメッセージの後に文書を書く。`Content-Type: application/ipp`、`Content-Length` なし（chunked）、`Expect: 100-continue`（`http.Transport.ExpectContinueTimeout`）。応答の `Content-Type` が `application/ipp` でなければ `ErrProtocol`。
- 401 は Basic で 1 回だけ再試行（`-password-file` か環境変数 `BDF_PRINTER_PASSWORD`）。
- `Printer` は `Get-Printer-Attributes`（`requested-attributes` は `all` ではなく必要な名前の列。`media-col-database` は大きい機種があるので別の要求で、`-verbose` のときだけ）から読む。`printer-config-change-time` と一緒に保存できる（m2 のキューの設定のため。m1 は毎回取る）。
- `Get-Jobs` は `which-jobs=completed` と `my-jobs=true`、`requested-attributes=job-id,job-name,job-state,job-state-reasons,job-impressions-completed`。
- CUPS 固有の `CUPS-Get-Printers`（0x4002）で `localhost:631` のキューを列挙（`bdf printers -cups`）。
- DNS-SD（`dnssd.go`）: 224.0.0.251:5353 と `[ff02::fb]:5353` に `_ipp._tcp.local` と `_ipps._tcp.local` の PTR を、QU ビットを立てて一時ポートから送る（RFC 6762 §5.4。応答は一時ポートに単一で返るので 5353 を束縛せず、Avahi / Bonjour と衝突しない）。`-timeout`（既定 2 秒）まで集め、同じ `UUID` / 同じ `ty`＋`rp` をまとめる。TXT の `rp` から `ipp(s)://host:port/rp`。インターフェースごとに送る（`net.Interfaces`）。単一 DNS の DNS-SD（企業の網）は m2 以降。m1 は URL の直接指定で足りる。

### 3.3 PWG Raster（`print/raster`、S3）

- ファイル: 先頭 4 バイト `RaS2`、ページごとに 1796 バイトのヘッダー（PWG 5102.4 §4.3 の順に: `PwgRaster`、MediaColor、MediaType、PrintContentOptimize、CutMedia、Duplex、HWResolution、InsertSheet、Jog、LeadingEdge、MediaPosition、MediaWeightMetric、NumCopies、Orientation、PageSize（pt）、Tumble、Width、Height、BitsPerColor、BitsPerPixel、BytesPerLine、ColorOrder、ColorSpace、NumColors、TotalPageCount、CrossFeedTransform、FeedTransform、ImageBox、AlternatePrimary、PrintQuality、VendorIdentifier / Length / Data、RenderingIntent、PageSizeName）。ColorSpace は PWG 5102.4 の表（`sgray` = 18、`srgb` = 19、`black` = 3）。ビッグエンディアン。
- 画素: 行ごとに「繰り返し数 − 1」の 1 バイト、続いて連: 0〜127 は次の画素を（値＋1）回、128〜255 は（257 − 値）個の画素がそのまま。1 ビットは 1 バイト（8 画素）を画素の単位にする。`black_1` は 1 = 黒、`sgray_8` / `srgb_8` は 0 = 黒。
- 入力は `imagebdf.RegionRows` の帯（`*image.Gray`、`*image.Paletted`（2 値）、`*image.RGBA`）。帯の中で同じ行が続けば繰り返し数にまとめ、行を RLE にして `io.Writer` に流す。ページ全体の画像は持たない。ヘッダーの値（幅、高さ、解像度、用紙）は描く前に全部分かる。
- 用紙への配置（`place.go`）: 文書のページ（pt）と用紙（PWG の名前 → mm）から、`fit`（収める）/ `actual`（原寸、中央）/ `shrink`（大きいときだけ縮める。既定）と自動回転（縦横が逆なら 90 度）で、用紙の画素格子に対するページの変換を決め、`RegionRows` に「用紙の矩形をページ座標で表した領域」と `Rotate` を渡す。用紙の外の部分は白。
- 両面の裏面: `pwg-raster-document-sheet-back` が `flipped` なら裏面を上下反転、`rotated` なら 180 度、`manual-tumble` なら綴じが短辺のときだけ 180 度、`normal` はそのまま。ヘッダーの CrossFeedTransform / FeedTransform も合わせる（CUPS の `rastertopwg` と同じ規則。実機で確かめる）。
- 部数: `copies-supported` の上限が足りればヘッダーの NumCopies と属性の `copies`、足りなければ送る側でページを繰り返す（`sheet-collate` も）。
- 解像度: `-resolution` か、`printer-resolution-supported` から既定（モノクロ 600、カラー 300。無ければ最も近い値）。型: `print-color-mode=monochrome` なら `sgray_8`、`-mono bilevel` で `black_1`（無ければ `sgray_1`）、カラーは `srgb_8`。
- `imagebdf` に足すもの: `RegionRows`（`Region` の帯版。`PageRows` はその特殊形になる）と `Options.Rotate`。どちらも既存の描画の行列に回転と平行移動を前置するだけで、帯の分割は出力の行で行うので並列の帯描画がそのまま使える。
- テスト: `pwg_read.go` で戻して `imagebdf.PageDPI`（回転なし、原寸）と画素が一致すること。testdata の全文書を 300 dpi の `sgray_8` と `black_1` で書いて、バイト列のハッシュを golden に（確定性）。両面の裏面の変換、1 ビットの詰め方、行の繰り返しの境界。

### 3.4 チケットと形式（`print`、S4）

- `Ticket` の項目は `print-requirements.ja.md` §4.1 の表。`ToIPP(p *ipp.Printer) (job []ipp.Attribute, unsupported []string)` が `*-supported` と照合する（`media` は `media-supported`、`media-col` は `media-col-database` があればそれ、`sides-supported`、`print-color-mode-supported`、`print-quality-supported`、`printer-resolution-supported`、`copies-supported`、`finishings-supported`、`print-scaling-supported`、`page-ranges-supported`）。受けない項目は `unsupported` に入れ、`bdf print` は警告（`-strict` でエラー）。
- 用紙（`media.go`）: ページの大きさから候補（JIS B5 = 182 × 257 mm = 515.9 × 728.5 pt …）を ±2 mm で探し、`media-ready` → `media-supported` の順に、無ければ `custom_<w>x<h>mm`（`media-supported` に `custom_min_…` / `custom_max_…` があるとき）、それも無ければプリンタの `media-default` に `print-scaling=auto`。
- 形式（`select.go`）: `-format` の明示 > `application/pdf`（`document-format-supported` にあり、`pdfbdf` が警告なしで書ける見込みのとき。`FILTER` や Canvas 固有の合成を持つ文書は警告し、`-raster-if-lossy` で Raster に落とす）> `image/pwg-raster` > `image/urf` / `application/PCLm`（m1 では「m3 で対応」のエラー）> それも無ければ受ける形式の一覧を添えたエラー。
- PDF のとき: `pdfbdf.Write` に `SkipRoles`、`Views`、`Pages`（`page-ranges-supported` が無ければ）を渡し、部数・丁合・両面・合わせ方・ページ範囲は IPP の属性で。暗号化した文書は平文の PDF になる（`ipps` でない接続なら警告）。
- Raster のとき: `SkipRoles` は `imagebdf` のレイヤー省略（既存の仕組みがあればそれ、無ければ描くレイヤーを選ぶ小さな追加）、ページ範囲と部数は送る側で。
- 文書（`doc.go`）: `.bdf`（単一 / 分割、`$BDF_PASSWORD`）か、`converter` の入力（`converter/all` を import するのは `cmd/bdf` だけ。`print` パッケージは `*bdf.Document` を受ける）。ビューは `-view`、既定はページのある最初のビュー（`fixed` / `flow`）。`scroll` は帯を積んだ 1 ページ、`sheet` は印刷ビューがあればそれ、無ければ m1 ではエラー（PDF 書き出しと同じ扱いは m2 で）。

### 3.5 `bdf print` の流れと失敗の扱い（S4, S5）

1. フラグを読む。`-printer` は URL か、`bdf printers` の名前（DNS-SD で引く）。
2. 文書を用意する。
3. `Get-Printer-Attributes`。接続できなければ `ErrConnect`（終了 3）、証明書の指紋が違えば `ErrCertificateChanged`（終了 3、`-insecure` か記録の削除を案内）。`printer-state` が `stopped` なら理由を出して、`-wait-printer` が無ければ `ErrPrinterStopped`（終了 5）。
4. チケットを写し、`unsupported` を警告（`-strict` なら終了 4）。
5. 形式を選ぶ。
6. `Validate-Job`。`server-error-operation-not-supported` は無視、`client-error-*` は属性の名前を添えて終了 4。
7. `Create-Job` → `job-id`、`job-uri`。`Send-Document`（`last-document=true`、`document-format`、`document-name`）で描きながら送る。`client-error-document-format-not-supported` なら次の形式で 1 回だけやり直す（PDF → Raster）。
8. `Get-Job-Attributes` を 1 秒から 5 秒まで伸ばしながら呼び、`job-state` と `job-state-reasons`、`job-impressions-completed` を表示（`-q` で黙る）。`completed`（終了 0）、`canceled`（終了 1）、`aborted`（理由を添えて終了 1）。`-wait=0` なら送った時点で終わる（終了 0、`job-id` を表示）。
9. Ctrl-C（`os/signal`）で `Cancel-Job` を送ってから終了 1。
10. 送信の途中で接続が切れたら、`job-id` を持っていれば `Get-Jobs` で状態を確かめ、プリンタが知っていれば再送しない（m2 の二重印刷の防止の元になる）。

`bdf printers`: `-timeout`、`-cups`（`localhost:631` のキュー）、`-url`（1 台の能力の表）、`-json`。

### 3.6 TLS の方針（S5）

`ipps` は、システムの CA で検証できればそれ、できなければ葉の証明書の SHA-256 を `$XDG_CONFIG_HOME/bdf/printers.json`（`os.UserConfigDir()`）に `printer-uuid` か host をキーに記録し、以後は一致を求める（`tls.Config.InsecureSkipVerify=true` ＋ `VerifyPeerCertificate` で自前の検証。CUPS の `TrustOnFirstUse` と同じ）。`-insecure` で検証なし（警告）。記録の形式は JSON で、m2 のサーバーの設定と共有する。

### 3.7 予算と限界

- IPP の応答: 4 MiB、属性 10,000、値 64 KiB、深さ 8。
- DNS-SD: 応答 1,000 件まで、1 件 9 KB まで。
- 文書: 大きさに上限なし（流す）。印象数（部数 × ページ）は `-max-impressions`（既定 10,000）で止める。
- 待ち: `-wait`（既定 10 分。0 で待たない）。
- `imagebdf` の限界と警告はそのまま。警告は `bdf print` が出す。

## 4. テスト

- **単体**（`go test ./print/...`）: 3.1 の往復と fuzz、`Printer` の読み取り（testdata の応答）、`Ticket.ToIPP` の照合（`unsupported` の列）、用紙の選択（JIS B5、はがき、custom）、形式の選択、PWG Raster の往復と golden、配置（fit / actual / shrink、回転）、裏面の変換、DNS-SD の応答の解釈（testdata のパケット）。
- **偽のプリンタ**（`print/ipp/ipptest`）: `cmd/bdf` のテストが `bdf print` を通して走らせ、受け取った文書（PDF は `converter/pdf` で読み戻して絵を比べる、Raster は `pwg_read` で戻す）と属性（`copies`、`sides`、`media`）を検査。`copies-supported` が 1 の設定で送る側の繰り返しを確かめる。`document-format-not-supported` を返して Raster への切り替えを確かめる。`stopped` と `aborted` の遷移で終了コードを確かめる。
- **CI**（`.github/workflows/ci.yml` の `print` ジョブ、Ubuntu）: `apt-get install cups cups-ipp-utils`、`ippeveprinter -p 8631 -f application/pdf,image/pwg-raster -c /tmp/save test` を起動して `bdf print -printer ipp://localhost:8631/ipp/print` で testdata の文書を送り、保存されたファイルを検査。`lpadmin -p test -E -v ipp://localhost:8631/ipp/print -m everywhere` で CUPS にキューを作り、`bdf print -printer ipp://localhost:631/printers/test` と `bdf printers -cups` を確かめる。Windows と macOS の runner で `go vet ./print/... && go test ./print/...`（ネットワークなし。DNS-SD と TLS のコードがビルドできること）。
- **実機**: `print-requirements.ja.md` §7 の機種で、PDF と PWG Raster、両面（`sheet-back` の 4 値）、JIS B5、はがき、トレイ、`copies-supported`、`page-ranges`、認証の有無を確かめ、`.knowledge/concepts/system/` に機種ごとの表（`system:printer-quirks`）を作る。

## 5. 順番

1. **S1** `print/ipp` の符号化（message.go、tags.go、attrs.go、testdata、FuzzDecode）。ほかの何にも依らない。
2. **S2** クライアント、`Printer`、DNS-SD、`ipptest`。`bdf printers` が動く。
3. **S3** `imagebdf.RegionRows` と `Options.Rotate`、`print/raster` の書き手と読み手、配置。`bdf render` の `-band` の経路でも `RegionRows` が使えることを確かめる。
4. **S4** `Ticket`、`media`、`select`、`doc`、`print.go`、`pdfbdf.Options.Pages`、`bdf print`。偽のプリンタで通す。
5. **S5** 失敗の扱いと終了コード、TLS の指紋、Ctrl-C、`converter` の入力、CI の `print` ジョブ、docs。

S1 と S3 は独立なので並行できる。S4 が合流点。

## 6. 後で変えられる位置に置くもの

- IPP の符号化を `goipp` に替えるなら `print/ipp/message.go` だけ。`Client` と `Printer` は `Message` の型にだけ依る。
- DNS-SD を依存に替えるなら `dnssd.go` だけ。`bdf printers` と `-printer name` は `[]Service` を受ける。
- 形式の既定（PDF 優先）は `select.go` の順の 1 行。キューごとの既定は m2 の設定で上書き。
- `print/` → `printing/` の改名は import パスの置換だけ。

## 7. m2 への引き継ぎ

- `print.Job`（状態、`job-id`、時刻、警告）と `Ticket` はそのままスプールの制御ファイルの中身になる。
- `Printer` と TLS の指紋の記録は、キューの設定として保存される。
- `ipptest` は m2 のキュー・再送・二重印刷の防止のテストにそのまま使う。

## 8. 実装で変わったこと（2026-10-19）

S1〜S5 を実装した（`milestone:print-m1`、`docs/design.md` §3.36）。計画から変えたのは次の 4 つ。

- `imagebdf` には `RegionRows` と `Options.Rotate` ではなく、`Place`（用紙の画素数、倍率、四分の一回転、位置、左右・上下の鏡映）と `PlaceRows` を足した。両面の裏面を鏡映して受けるプリンタ（`pwg-raster-document-sheet-back` の `flipped`）があり、回転だけでは足りなかった。帯の描き方は `PageRows` と共有する。
- `page-ranges` は送らず、ページの選択は PDF の側（`pdfbdf.Options.Pages`）だけで行う。プリンタごとの扱いの差を避ける。
- 終了コードに 6（`-wait` のうちに終わらない）を足した。
- 確かめ方: 作業環境には IPv6 がなく `ippeveprinter` が起動しないので、手元では CUPS 2.4.7 の `cupsd`（raw のキューと IPP Everywhere のキュー）と偽のプリンタで確かめ、`ippeveprinter` は CI の `print` ジョブ（`test/print/ci.sh`）で確かめる。実機での確認（`manual-devices`）はまだ。
