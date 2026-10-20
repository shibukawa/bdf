# 印刷 実装案（print-m5）

2026-10-20 の実装案。同日 S1〜S5 を実装した（`.knowledge/` の `milestone:print-m5`。計画から変えた所はそこの `changed_from_plan`）。各言語の命令は同日に各社のリファレンスを Web で読んで確かめた（実機は無い）。要件は `print-requirements.ja.md` §3.9・§4.1（機器ごとの拡張）・§6 の 5 と、`.knowledge/` の `requirement:print`、`system:label-printer-languages`、`system:dot-impact-escp`、`pattern:label-sheet`、`pattern:roll-receipt`、`pattern:continuous-text-form`。print-m1〜m4 は main に入った（#104〜#107）。2026-10-20 利用者「印刷 milestone 5 の実装をお願いします」。

## 0. 範囲

m5 は「ラベル・レシート・連帳」。IPP を話さず、PDF も PostScript も受けない機器に、BDF のページをその機器の言語で届ける。

1. **ラベルの言語の書き手**（`print/pdl`）: ZPL（Zebra）、SBPL（SATO）、TPCL（東芝テック）。ラベル 1 枚を 1 ビットのラスタで送り、BARCODE の印のある記号は機器のバーコード命令に置き換える。
2. **レシートの書き手**: ESC/POS（Epson TM 系と互換機）。帯のラスタ（`GS v 0`）、カット（`GS V`）、ドロワ（`ESC p`）。
3. **ドットインパクトの書き手**: ESC/P（Epson VP 系、ESC/P-J と互換機）。24 ピンのビットイメージ（`ESC *`）、用紙長（`ESC C`）、ミシン目スキップ（`ESC N`）。
4. **2 値化の追加**: 画像だけを組織的ディザ（`imagebdf.Options.DitherImages`）、記号のモジュールを画素に揃える（`AlignBarcodes`）。
5. **機種のプロファイル**（dpi、印字幅、最大長、言語、印刷側バーコードの可否）と**チケットの機器ごとの拡張**。
6. **搬送路**: 9100（既存）、シリアル（Bluetooth の SPP を含む）、USB の装置ファイル（Linux / BSD）、Windows のメーカードライバのキューへ `winspool` の `RAW`、CUPS の raw キュー（IPP で `application/vnd.cups-raw`。macOS と Linux の USB）。

**入れないもの**:

- **ESC/P・ESC/POS のテキストモード**（最初の S1〜S5 には入れなかった）。テキスト帳票テンプレート（form-m6）が書く印字像の Part（`t: text`）から作る決まり（`requirement:text-template`）で、2026-10-20 利用者の判断で form-m6 は別のセッションで並行して実装した。両方が main に入った同日に §6 のとおり足した。
- Bluetooth Low Energy（機種ごとの GATT。SPP の機種はシリアルとして届く）、P-touch Template・Dymo（`system:label-printer-languages` にあるが要件の対象外）、プリンタへのフォーマットの保存（`^DF` / `^XF`。送信側で組版を済ませるので要らない）、状態の取得（`~HS`、`DLE EOT`。9100 の見張りと同じく後で）。

## 1. 決めたこと（暫定で推奨どおり進める）

- **印刷側バーコードの印**: 2026-10-20 利用者が決定。spec §7.8 の MARK に kind 13 `BARCODE` を足す（payload は JSON: `kind`・`text`・`level`・`quiet`・`noText`・`box`）。帳票（Word・Excel・PowerPoint・pdfme）が記号を描くときに付ける（`barcode.Symbol.Fill`）。読み手は知らない kind を無視するので互換は崩れない。
- **ラスタの描き方**: `imagebdf` の `PlaceRows` で 1 ビット（`Bilevel`、`AntialiasOff`）。画像だけディザ（帯をまたいで同じ絵になる組織的ディザ。誤差拡散は帯の境で切れるので使わない）、文字と線はしきい値のまま。印刷側に置き換えない記号はモジュールを画素の整数倍に揃えて描き直す（203 dpi でも読める）。
- **ページとラベルの写し方**: ラベルの幅はプロファイルの印字幅（チケットの `media` があればその幅）、長さは `media` の長さ、無ければページの高さ。ページはラベルに `raster.Fit`（既定は大きければ縮める）。レシートと連帳は長さ可変で、ページの高さをそのまま（幅に縮めた分だけ縮む）。
- **印刷側バーコード**: ZPL・SBPL・TPCL で、記号が正立か 90° 単位の回転で、モジュールが 1 ドット以上になるとき。ZPL は `^BY` のモジュール幅を整数ドットに丸め、記号を箱の中央に置く。置き換えられない記号（傾いている、細すぎる、言語に無い種類）はラスタに揃えて描く。チケットの `bdf-printer-barcodes=false` で全部ラスタにできる。
- **チケットの拡張**: 名前は IPP にあるものは IPP（`print-darkness`（-100〜100）、`print-speed`（1/100 mm/s）、`media-tracking`（`continuous`・`gap`・`mark`）。PWG 5100.13 と LPrint の語彙）、無いものは `bdf-` を付ける（`bdf-cut`（`none`・`page`・`job`）、`bdf-cash-drawer`（1 か 2）、`bdf-label-gap`（1/100 mm）、`bdf-perforation-skip`（1/100 mm）、`bdf-print-mode`（`raster`・`text`）、`bdf-printer-barcodes`）。言語に通り道の無い値は `unsupported` に入れる（`-strict` で断る）。複写の枚数は ESC/P に命令が無いので扱わない（機種のパネルの設定）。
- **プロファイル**: 組み込みの表（Zebra ZD421 203/300 dpi、GK420d、SATO CL4NX、東芝テック B-EX4T1、Epson TM-T88VII、TM-m30III、58 mm 機、Epson VP-F2000、VP-D1300 など）と、URL の `profile=` / `lang=`・`dpi=`・`width=`（`104mm`、`4in`、ドット数）で上書き。
- **搬送路の URL**: `socket://host:9100?lang=zpl&profile=zd421-203`（既存の 9100 に言語を足す）、`serial:///dev/ttyUSB0?baud=9600&lang=escpos`（Windows は `serial://COM3`）、`usb:///dev/usb/lp0?lang=zpl`（Linux / BSD の装置ファイル）、`winspool:///Zebra%20ZD421?lang=zpl`（Windows のキューに `RAW`。共有のキューは `winspool://server/queue`）。CUPS の raw キューは `ipp://localhost:631/printers/<名前>` に `-format zpl` などで、IPP の `application/vnd.cups-raw`（無ければ `application/octet-stream`）として送る。シリアルの設定は `baud`・`bits`・`parity`・`stop`・`flow`（`none`・`rtscts`・`xonxoff`）。
- **二重に刷らない**: 9100 と同じく、送り始めた後に落ちたジョブは人に任せて保留する（シリアル・USB・winspool も同じ扱い。どれも状態を返さない）。

## 2. 構成

```
image/barcode/mark.go        BARCODE の印（Mark、ParseMark、Symbol.Fill、Bars、BarHeight）
raster/imagebdf/barcode.go   Options.Barcode（印の記号を呼び出し側へ）、AlignBarcodes、DitherImages
print/pdl/                   書き手: zpl.go、sbpl.go、tpcl.go、escpos.go、escp.go、
                             render.go（ページ → 1 ビットの行と記号）、profile.go、settings.go、
                             読み手（テスト用。出力を絵に戻す）
print/raw/                   serial（unix は termios、Windows は DCB）、usb（装置ファイル）、
                             winspool（Windows の RAW）、URL の lang・profile
print/                       Ticket の拡張、PrintRaw が言語を書く、IPP に raw で送る
print/spool/、print/server/  設定の printers に言語とプロファイル、キューの既定のチケット
cmd/bdf/                     bdf print -printer serial://… usb://… winspool://…、-format zpl…、-darkness …
docs/                        cli、print-server、design §3.38、spec §7.8
```

## 3. 順番（スライス）

1. **S1** BARCODE の印（spec、`image/barcode`、帳票の描き手）、`imagebdf` の記号の扱いとディザ。
2. **S2** `print/pdl` の共通部と ZPL・SBPL・TPCL、読み手、プロファイル。
3. **S3** ESC/POS と ESC/P。
4. **S4** 搬送路（シリアル、USB、winspool、IPP raw）、`print` のチケットと送り方、スプールとサーバー、CLI。
5. **S5** docs、CI、`.knowledge`。

## 4. テスト

- **印**: 記号ごとに `Mark` → `ParseMark` → `Symbol` で同じ格子。帳票（docx・xlsx・pptx・pdfme）の出力に印があり、その後の描画が記号であること。
- **imagebdf**: 揃えた記号の走査線の黒白の長さがモジュールの整数倍で、格子のとおり（帯の大きさによらない）。呼び出し側に渡した記号の箱と回転、取ったときに描かないこと。ディザの黒の割合と帯での一致、塗りはディザしないこと。
- **書き手**: 各言語の出力を自前の読み手で絵に戻し、`imagebdf` の 1 ビットの絵と画素で一致（ZPL の Z64 は zlib を戻し CRC を確かめる）。印刷側バーコードの命令の位置・モジュール・高さ・回転・データ。チケットの拡張の命令。限界（長いラベル、幅を超えるページ）と fuzz（読み手）。
- **搬送路**: 偽の 9100、擬似端末（unix のシリアル）、一時ファイル（usb の装置ファイル）、ipptest の raw 受け。winspool は Windows の CI で関数の結び付けまで。
- **実機**: 手元に無い。docs に各機種での確かめ方を書き、利用者が確かめる。

## 5. 後で変えられる位置に置くもの

- **テキストモード**（ESC/P・ESC/POS）: 2026-10-20 に足した（§6）。
- 状態の取得（ZPL の `~HS`、ESC/POS の `DLE EOT`、SBPL の ENQ）を 9100 の見張りに。
- プリンタに保存したフォーマット（`^DF`）へのデータだけの送信（大量のラベルで速くする）。

## 6. テキストモード（2026-10-20 に足した）

利用者「milestone 5 と 5 に必要な 6 の機能の一部は終わりました。残りの作業を終わらせてください」。form-m6 の印字像の Part（`form/textpart`）から書く。ページの絵は見ない。`print/pdl` の `text.go`（モードの決め方と文字の符号化）、`escp_text.go`、`escpos_text.go`。

- **既定**: ESC/P は Part を持つ文書をテキストモードで刷る（§3.9 の表のとおり）。ESC/POS はドットが既定で、`bdf-print-mode=text`（`bdf print -print-mode text`）のときだけテキスト。`raster` を求めれば ESC/P もドット。ラベルの言語に `text` は誤り。
- **刷れないとき**: ESC/P に無いピッチ（10・12・15・17.1・20 CPI のほか）や行送り（6・8 LPI と 180/n のほか）、漢字が 2 桁に収まらないピッチ（17.1・20 CPI）、プリンタより広い帳票は、既定のモードなら警告を出してドットで刷り、`text` を求めたときは誤り。書き手はまずメモリに書き、通ってから送る。
- **ESC/P**: `ESC @`、`ESC R 0`、ピッチ（`ESC P`・`ESC M`・`ESC g`、`SI` の縮小）、行送り（`ESC 2`・`ESC 0`・`ESC 3 n`）、用紙長（`ESC C n`、128 行以上はインチ）、スキップ（`ESC N`、無ければ `ESC O`）、左余白（`ESC l`）、上の余白は LF、行は文字と CR LF、ページの終わりに FF。漢字は漢字モード（`FS &` … `FS .`）の JIS コードで、`FS S` の左右の空きで 2 桁に。あいまいな幅を 1 桁と数えたテンプレートの 1 桁の文字（罫線素片、○、※）は半角の漢字（`FS !` のビット 1。`FS SI` と同じ。`FS T` の空き）。倍角はキヤノンの ESC/P エミュレーションの資料に従う。`ESC W` は全部の文字の横倍で `FS !` のビット 2 と同じフラグなので、漢字のある帳票では横倍・半角（ビット 1）・縦横 2 倍（ビット 3、`FS W`）を `FS !` でまとめて決める。`ESC w`（1 バイトの文字の縦倍）はベースラインから上へ伸びて改行の量を変えないので、倍高の行は次の行に刷る。縦横 2 倍の漢字は上半分の後に 24/180 インチ送って下半分を刷るので、`ESC J` で（行の高さ − 24/180）下げてから刷る。漢字の縦倍だけは命令が無いので普通の大きさで刷る（警告）。
- **ESC/POS**: `GS P` で送りの単位をドットに、`ESC M 0`（Font A、12 ドット、漢字 24）に `ESC SP` と `FS S` の空きを足して 1 桁を `floor(dpi / cpi)` ドットにし、そこから格子の位置 `round(桁 × dpi / cpi)` にずれた文字は `ESC $` で置く（空白は送らない）。行は `ESC J` で格子の行の位置まで送る（倍高の行は 2 行分）。漢字は `FS C 0`（JIS）と `FS &` … `FS .`、半角カナは `ESC t 1`、1 桁の罫線素片は `ESC t 0`（コードページ 437。太線は細線にして警告）、`×` などは `ESC t 2`（850）、どれにも無い 1 桁の文字は `?` にして `ambiguous: wide` を勧める警告。倍角は `GS !` と、漢字のある帳票では `FS !`（同じ大きさ。重ねても倍にならない）。カットとドロワはドットと同じ。
- **ほか**: `¥` は `ESC R 8` の 0x5C、JIS X 0208 に無い文字（IBM 拡張、絵文字）は `〓` で警告。`-skip background` でプレ印刷の範囲を空白に。チケットの `bdf-perforation-skip` はテキストモードでは受けない（帳票の skip を使う）。結果の形式は「ESC/P (text)」「ESC/POS (text)」。
- **確かめ方**: 各言語の出力を自前の読み手で格子に戻し、`Image.Text`（余白込み）と比べる（`print/pdl/text_test.go`）。ESC/POS の命令は Epson の TM のリファレンスを内蔵ブラウザで読んだ（表示された本文だけ。利用規約には同意していない）。利用者が示した lrks の記事（https://telecotele.com/projects/escpos/）で、東芝テック B-EP2DL が Shift_JIS を受けず JIS を `FS &` … `FS .` で囲むこと、TM-T90 が `FS !` で漢字を倍幅にすることを確かめ、JIS と `FS !` を選んだ。ESC/P-J84 は Epson のリファレンスが読めなかったが、利用者の指摘でブラザー（PT-9700PC の ESC/P コマンドリファレンス）とキヤノン（imageRUNNER ADVANCE C350F の ESC/P エミュレーション）の資料を読み、漢字モードが JIS X 0208 のコード、`FS S` / `FS T` の単位が 1/180 インチ、`FS !` のビット、`ESC w` が上へ伸びること、`FS W` が 2 回に分けて刷ることを確かめた。Epson VP の実機では確かめていない。
