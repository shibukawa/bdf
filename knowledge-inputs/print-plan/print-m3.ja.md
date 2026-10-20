# 印刷 実装案（print-m3）

2026-10-20 の実装案。同日 S1〜S5 を実装した（`.knowledge/` の `milestone:print-m3`。計画から変えた所はそこの `changed_from_plan`）。要件は `print-requirements.ja.md` §3.3・§3.4・§6 の 3 と、`.knowledge/` の `requirement:print`（facts の formats・vector_formats・enterprise_printers）。print-m1（IPP で直接、`bdf print`）と print-m2（スプールとプリントサーバー、`bdf printd`）は main に入った（#104、#105）。2026-10-20 利用者「次の実装をお願いします。m3。」。

## 0. 範囲

m3 は「第 2 の形式と古い搬送路」。PDF を受けないプリンタ、IPP を話さないプリンタに届かせる。

1. **PostScript の書き手 `vector/psbdf`**（第 2 のベクタの経路）。PDF を受けず PostScript を受ける業務機（2010 年前後までの機種、Kyocera KPDL・Brother BR-Script などの互換機、大判の PS 機）に、文書をベクタのまま渡す。PostScript には透明が無いので、透明な所だけを画像にする平坦化（`vector/internal/flatten`）と一緒に作る。
2. **Apple Raster（`image/urf`）と PCLm（`application/PCLm`）の書き手**（`print/raster`）。PDF も PWG Raster も受けない AirPrint 専用機と Mopria 専用機のため。
3. **9100（生のソケット）と LPD の搬送路**（`print/raw`）。IPP を持たないプリンタに PostScript か PDF をそのまま流す。`bdf print -printer socket://host:9100`、`lpd://host/queue`。スプールとサーバー（`bdf printd`）のプリンタにも書ける。
4. **IPP USB**: Linux の `ipp-usb` と macOS は USB の IPP 機を `localhost` の IPP として見せ、ループバックで DNS-SD に広告する。`bdf printers` の探索でループバックも聞く。
5. **形式の選び方**を要件どおり **PDF → PostScript → PWG Raster → Apple Raster → PCLm** にする。

`CUPS-Get-Printers` の一覧は m1 で入った（`bdf printers -cups`）。

**入れないもの**: HP-GL/2 の書き手（要件は「CAD の現場があれば」。現場の声を待つ。m3 の後で足せる位置に置く）、PCL XL と各社の PDL（要件どおり作らない）、PCL 5 のラスタ（要件の最後の後ろ盾。PS も Raster も受けない機種が出たら）、Windows の OS 経由（作らない）、PJL による状態の取得（9100 の双方向。後で）。

## 1. 決めたこと（暫定で推奨どおり進める）

- **PostScript は Level 3**（`%!PS-Adobe-3.0`、DSC 準拠）。2000 年代以降の PS 機は全部 Level 3。滑らかなシェーディング（`shfill`）、`/FlateDecode` を使う。画像と圧縮した本文は ASCII85 で 7 ビットに収める（LPD やシリアルを通っても壊れない）。
- **文字は Type 3 のフォント**（グリフのアウトラインを手続きにし、`setcachedevice` でプリンタのフォントキャッシュに載せる）。Type 42（TrueType）と CFF（CIDFontType 0C）は互換機の解釈器に差が大きく、日本語は CID の組みが要る。Type 3 はどの解釈器でも同じに出る。ヒントは無いが 600 dpi で問題にならない。1 つのフォントは 256 グリフまでで、日本語のように多いときは分けて持つ。
- **平坦化は「透明な所の矩形だけを画像に」**: ページを 2 回なめる。1 回目で透明なもの（α が 1 でない色と GLOBAL_ALPHA、合成、α か合成のあるグループ、ソフトマスク、影、α を持つ画像、α の止まりのあるグラデーション、パターン、円錐のグラデーション、FILTER）が描かれる矩形（ページの単位、クリップと交わる所）を集めて重なりをまとめる。2 回目でベクタを、それらの矩形を除いたクリップ（ページの矩形＋矩形の偶奇）の下に書き、透明なものは書かない。最後に各矩形を `imagebdf.Region` で描いた画像（既定 300 dpi。ページの中で後から描くものも含めて、その矩形の正しい絵）を置く。矩形の外は全部ベクタのまま。Acrobat の平坦化の「原子領域」を矩形に粗くしたもの。
- **用紙・部数・両面**は PostScript では `setpagedevice`（`/PageSize`、`/NumCopies`、`/Duplex`、`/Tumble`、`/MediaType`、`/MediaPosition`）で文書の中に書く（9100 と LPD に属性の通り道が無いため）。IPP では今どおり属性で。
- **9100 と LPD のプリンタには能力を聞けない**ので、送る形式（`ps` か `pdf`）は URL（`socket://host:9100?format=pdf`）、フラグ（`-format`）、設定ファイルで決める。既定は PostScript（9100 の古い機種の共通語）。PDF を送るときは PJL の見出し（`@PJL ENTER LANGUAGE=PDF`、部数・両面・用紙）を付ける（`pjl=false` で外す）。
- **状態**: 9100 と LPD はジョブの状態を返さない。送り終えたら「送った」で終わる（`-wait` は送り終えるまで）。LPD は受け取りの確認（各段の `\0`）まで。
- **スプールで二重に刷らない**: LPD は制御ファイルを受けて初めて刷る（RFC 1179）ので、制御ファイルを送る前に落ちたら送り直してよい。9100 は途中まで送った分が刷られることがあるので、送り終えた印が無いまま落ちたジョブは人に任せて保留する（`bdf-duplicate-check`）。送り終えた印があれば終わったとみなす（見届ける手段が無い）。
- **URF と PCLm**: 解像度と色はプリンタの `urf-supported`（`RS300`、`W8`、`SRGB24`、`DM1`〜`DM4`）と `pclm-source-resolution-supported`・`pclm-strip-height-preferred`・`pclm-compression-method-preferred`・`pclm-raster-back-side` から選ぶ。URF の行の符号は PWG Raster と同じ（CUPS の `raster-stream.c` と同じ）なので、書き手の符号化を共有する。PCLm は帯ごとの画像（Flate か JPEG）を並べた PDF 1.7 で、Mopria の PCLm 1.0 の形（2 行目の `%PCLm 1.0`、ページごとの帯）に従う。

## 2. 構成

```
vector/psbdf/              PostScript 3 の書き手。Write(w, doc, *Options)。pdfbdf と同じ命令の解釈で、
                           パス・クリップ・線・破線・色・線形と放射のグラデーション（shfill）・画像・文字（Type 3）・
                           USE。透明な所は flatten の矩形の画像。setpagedevice で用紙・部数・両面
vector/internal/flatten/   透明な所の矩形を集め（Collector）、まとめ、画像にする（imagebdf.Region）
print/raster/              urf.go（Apple Raster の書き手）、pclm.go（PCLm の書き手）、読み手（テスト用）
print/                     select.go（PDF → PS → PWG → URF → PCLm）、print.go（PS・URF・PCLm の送り方）
print/raw/                 9100 と LPD: Send(ctx, url, doc, *print.Options)。PJL の見出し
print/ipp/                 DNS-SD でループバックも聞く（ipp-usb）、能力に pclm-* と application/postscript
print/ipp/ipptest/         URF・PCLm・PostScript を受ける
print/spool/               9100 と LPD のプリンタ（Printer.Client の代わりに Printer.Raw）。照合の規則
print/server/              設定の printers に socket:// と lpd://、format、pjl
cmd/bdf/                   bdf print -printer socket://… lpd://…、-format ps|urf|pclm
docs/                      cli、print-server、design §3.36、psbdf の API
test/print/ci.sh           ippeveprinter に PostScript、9100 の偽の受け手、Ghostscript で PS を描いて確かめる
```

## 3. 順番（スライス）

1. **S1** Apple Raster と PCLm の書き手、読み手、形式の選び方、ipptest、`-format urf|pclm`。
2. **S2** `vector/psbdf` と `vector/internal/flatten`。Ghostscript で描いて `imagebdf` の絵と比べる。
3. **S3** PostScript を IPP で送る（`application/postscript`）、選び方の順を PDF → PS → PWG にする。
4. **S4** `print/raw`（9100、LPD、PJL）、`bdf print` の URL、スプールとサーバーのプリンタ、照合の規則。
5. **S5** IPP USB（ループバックの探索）、docs、CI、`.knowledge`。

## 4. テスト

- **psbdf**: testdata の文書（帳票、PowerPoint、Word、Excel の印刷ビュー、PDF 由来のソフトマスク、draw.io、写真、縦書き）を PostScript にして Ghostscript（`gs -sDEVICE=png16m`）で描き、`imagebdf` で描いた絵との差を pdfbdf の往復テストと同じ尺度で比べる。Ghostscript が無ければ飛ばす（CI の print ジョブで入れる）。透明な所が画像になり、それ以外がベクタのまま残ること（ページの中の画像の数と大きさ）。DSC の構造と `setpagedevice`。限界（入れ子、再利用、巨大な画像）と fuzz。
- **URF と PCLm**: 自前の読み手で戻して描いた画素と一致。PCLm は converter/pdf と Ghostscript でも開ける。
- **print**: ipptest に PostScript・URF・PCLm だけを名乗らせ、選び方と送り方。
- **raw**: Go の偽の 9100 の受け手と LPD の受け手（制御ファイルとデータファイル、確認の `\0`）。PJL の見出し。スプールで 9100 と LPD のプリンタに流し、落ちた時点ごとの照合（LPD は送り直し、9100 は保留）。
- **CI**: `ippeveprinter -f application/postscript,…` に `bdf print -format ps`、受けた PS を Ghostscript で描けること。9100 は `nc -l` か Go の受け手。

## 5. 後で変えられる位置に置くもの

- HP-GL/2 の書き手は `vector/hpglbdf` として `flatten` を共有する（現場の要望で）。
- Type 42・CFF の埋め込みは、Type 3 で困る解釈器が出たら `psbdf.Options.Fonts` で足す。
- PJL の状態の取得（`@PJL INFO STATUS`、`USTATUS`）は 9100 の見張りとして後で。
