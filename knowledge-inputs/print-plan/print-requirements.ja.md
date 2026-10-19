# 印刷サブシステム 要件整理（print）

2026-10-19 時点の要件の整理。構造化した記録は `.knowledge/` の `requirement:print`（要件と未決の論点）、`system:ipp-driverless`（プロトコルの調査）、`milestone:print-m1`（最初の実装のまとまり）にある。ここは「何を、なぜ、どの順で」を叙述で書く。ファイル単位の実装案は、論点が決まってから `print-m1.ja.md` として別に書く。

## 0. 位置づけ

BDF の出口は 3 つ揃った。画面（ビューア）、PDF（`vector/pdfbdf`。タグ付き PDF と PDF/A、プレ印刷用紙向けのレイヤー省略）、画像（`raster/imagebdf`。帯に分けた描画、モノクロ 2 値、ページ並列、600 dpi のページを帯の記憶量で流せる）。帳票（`form`）が `.bdf` を量産できるようになり、残る出口は紙である。

印刷は 3 つの層に分けて考える。

1. **印刷指示（チケット）**: 何を、何部、どの用紙に、両面か、モノクロか、どのレイヤーを省くか。
2. **スプール（キュー）**: ジョブを受け付けて永続化し、順に流し、失敗を扱い、状態を見せる。
3. **接続（出力先）**: プリンタに文書を渡す。プロトコル（IPP、OS のスプーラ、生ソケット、USB …）と、プリンタが受け取る文書形式（PWG Raster、PDF、Apple Raster、PCLm、各社の PDL）。

ローカル（手元の 1 台に `bdf print`）とプリントサーバー（複数台のリモートプリンタを束ねる常駐プロセス）は、1 と 3 を共有し、2 の有無と大きさが違うだけにする。`bdf print` はメモリ上の小さなスプール（1 ジョブを送って終わる）、サーバーはディスクのスプール（再起動をまたぐ）。

BDF 側の強みをそのまま使う。文書はベクタの命令列なので、PDF を受けるプリンタには **ベクタのまま**渡してプリンタの解像度で描かせ（`pdfbdf`）、受けないプリンタには自前のラスタライザで描いて**ビューアと同じ絵**を渡す（`imagebdf`）。どちらも OS のドライバに頼らない。同じ文書は内容アドレスで 1 回だけ持てるので、1000 枚の請求書をまとめた `.bdf` も、スプールに置くのは 1 つで足りる。

## 1. ゴール

- **G1 手元のプリンタに `.bdf` を印刷できる**: `bdf print invoice.bdf -printer ipp://…` が、OS のドライバを使わずに Go だけで紙を出す。ベクタのまま渡せるプリンタにはベクタのまま渡す。`converter` が読める入力（PDF、Office、画像 …）はその場で変換して印刷する。
- **G2 Windows、macOS、Linux で同じコード**: プロトコルは OS に寄らないものを主にし、OS のスプーラは補助（後述の層 2）。cgo なし。
- **G3 プリントサーバー**: 複数のキュー（1 キュー = 1 台、または複数台の組）を持ち、ジョブを受けて順に流し、再起動をまたいで続きから流す。クライアントは CLI と HTTP API。後の段階で、サーバー自身が IPP のプリンタとして見える（各 OS からドライバなしで投げられる）。
- **G4 帳票の印刷に要る指示が通る**: 部数と丁合、両面（長辺・短辺）、用紙（サイズ・トレイ・種類。JIS B と はがき と 封筒を含む）、モノクロ、ページ範囲、用紙への合わせ方（原寸・収める・自動回転）、プレ印刷用紙向けのレイヤー省略（`SkipRoles`）、ホチキスなどの仕上げの素通し。
- **G5 失敗の扱い**: プリンタの停止・用紙切れ・ネットワーク断・プリンタ側の中止を区別し、再送できるものは再送し、できないものは保留して理由を見せる。再起動の後に同じジョブを二重に刷らない。
- **G6 新しいプロトコルから始める**: 最初の出力先は IPP（driverless）。古いプロトコル（LPD、9100 の生ソケット、SMB）と OS 固有の経路は後の段階で、IPP の無いプリンタのためだけに足す。
- **G7 ラベル・レシート・連帳への道を塞がない**: 出力先の抽象は「文書形式 × 搬送路」で、PDL の書き手（ZPL、ESC/POS、ESC/P …）と搬送路（USB、シリアル、Bluetooth）を後から足せる形にする（`system:label-printer-languages`、`system:dot-impact-escp`）。

## 2. 非ゴール（初期）

- プリンタドライバを書くこと（Windows の v4 ドライバや Print Support App、CUPS のフィルタや Printer Application）。
- 各社固有の PDL（PCL 5/6、PostScript、各社の ESC/Page 系）の生成。IPP の無い古い機種は OS のスプーラか生ソケットで PDF を渡す。
- ラベル・レシート・ドットインパクト向けの PDL（ZPL、SBPL、TPCL、ESC/POS、ESC/P）。形態の整理は `taxonomy:form-layout-axes` にあり、後の段階（print-m5）。
- クラウド印刷サービス（Microsoft Universal Print、各社のクラウド）との連携。IPP INFRA の代理（プロキシ）も後。
- ブラウザ（wasm）からの印刷。ブラウザは `window.print()` の領分で、BDF のビューアからは「サーバーに投げる」で足りる。
- 課金・部門管理・認証印刷（IC カードで釈放）。履歴は出すが、その上の運用は外の仕事。IPP の属性（`job-password`、`job-account-id`）で渡せる分は渡す。
- Windows 共有プリンタ（SMB、`\\server\queue`）経由の印刷。Windows の印刷サーバーを義務にしている組織では BDF のサーバーをその代わりに置く。
- スキャン（eSCL）や FAX。
- プリンタの管理（ファームウェア、消耗品の発注、SNMP による詳細な監視）。状態は IPP が返す範囲。

## 3. 出力先の技術選定: なぜ IPP から始めるか

### 3.1 2026-10 の状況

**IPP が唯一の共通語になった。** IPP/1.1 は RFC 8010（符号化）と RFC 8011（意味論）、IPP/2.x は PWG 5100.12-2024（第 4 版、2024-11 承認）。ドライバなし印刷（driverless）の 3 つのブランドは、いずれも IPP の上に文書形式と DNS-SD の探索を決めたものである。

| ブランド | 主体 | 要求する文書形式 | 探索 |
|---|---|---|---|
| IPP Everywhere 1.1（PWG 5100.14-2020） | PWG | PWG Raster（`image/pwg-raster`）必須。PDF は推奨 | DNS-SD `_ipp._tcp` / `_ipps._tcp` |
| AirPrint | Apple | Apple Raster（`image/urf`）。PDF も多い | DNS-SD（TXT の `URF` キー） |
| Mopria | Mopria Alliance | PWG Raster、PCLm、PDF のいずれか | DNS-SD、Wi-Fi Direct |

2015 年ごろ以降の事務機・複合機・家庭用プリンタは、ほぼどれかを名乗る。1 台の実機は 3 つを同時に満たすことが多く、`document-format-supported` に PWG Raster・URF・PDF（＋ PCLm）が並ぶ。

**OS もそこへ寄せている。**

- **Windows**: 同梱の IPP クラスドライバ（Mopria ベース。Windows 10 21H2 からネットワークと USB の Mopria 機を扱う）で、ドライバなしに IPP プリンタを使う。Microsoft は 2023-09 に他社製ドライバの提供終了計画を出し、2026-01-15 以降は Windows Update に新しい他社製ドライバを載せず、2027-07-01 以降は他社製ドライバの更新をセキュリティ修正に限る。**Windows Ready Print（WRP）**: 2026-06-12 に発表（旧称 Modern Print Platform）、6 月 23 日のプレビュー更新から展開、7 月以降は**新しく追加するプリンタの既定のインストール方法**が IPP クラスドライバになった（窓の杜 2026-07-02）。導入済みのキューは変わらず、メーカー製ドライバも選べ、グループポリシー「Configure Windows Ready Print driver ranking」で従来の順位に戻せる。**Windows Protected Print（WPP）** はその上の厳格モードで、有効にすると他社製ドライバ・他社製のポートモニタと印刷プロバイダが外れ、WRP で入らないプリンタは追加できない（Windows 11 24H2 以降、既定では無効。Microsoft は将来の既定化を表明）。メーカー機能は Print Support App で足す。
- **macOS**: AirPrint が標準。PPD のメーカードライバは非推奨で（CUPS の `lpadmin` / `lpinfo` は `everywhere` 以外のモデルを非推奨とし、ドライバとバックエンドは将来の版で外すと明記。Apple 独自の拡張 `APDialogExtension` は 10.14 から非推奨）、macOS 26 Tahoe では古いメーカーのドライバが入らず「AirPrint で」と案内される機種が増えた。中の CUPS は Apple の分岐で 2.3.6（2022-05）のまま凍結、上流は OpenPrinting。アプリの印刷は PDF を `cupsd`（`localhost:631`、IPP）に渡す経路で、PostScript 機にも CUPS が PDF → PS にする。
- **Linux**: CUPS 2.4 / 2.5 が driverless を標準にし（`driverless` で IPP Everywhere の PPD を生成）、CUPS 3.0 は PAPPL ベースの `cups-local`（利用者ごと、一時キュー、IPP over ドメインソケット）と `cups-sharing`（共有、永続キュー、アカウンティング）に分かれ、ドライバをなくして Printer Application（それ自体が IPP プリンタ）に置き換える。libcups 3.0.2 は 2026-06、PAPPL 1.4.12 は 2026-08。`cups-local` / `cups-sharing` 3.0.0 の正式版が出たかは未確認（計画は 2025 年 8〜9 月だった）。USB のプリンタは `ipp-usb` が IPP over USB（USB クラス 7/1/4）を `localhost` の IPP に見せる。

つまり **CUPS のある macOS と Linux では、OS のスプーラ自体が IPP のサーバー**であり、プリンタに直接話す IPP クライアントが、そのまま OS のキューにも使える。Windows だけが IPP のサーバーを持たない（共有は SMB / RPC）が、Windows Ready Print の後は Windows が扱うプリンタ自体が IPP プリンタなので、Windows では OS を通さず直結だけにする（§3.5）。

**古いプロトコルは残るが、主ではない。** LPD（RFC 1179、1990）、AppSocket / JetDirect（TCP 9100 の生ストリーム）、SMB 共有、WSD（Microsoft の Web Services for Devices、2006）、状態を取る SNMP（Printer MIB）。IPP の無いプリンタは、ラベル・レシート・ドットインパクトと、古い事務機に残る。9100 は PDL を素通しするだけなので、ラベル・レシートの段階で「搬送路」として足す。

**次に来るもの。** IPP Everywhere 2.0 は安定版草案（2026-04-30 更新）で、未承認。IPP INFRA（PWG 5100.18-2015。プリンタ側の代理がサーバーからジョブを取りに来る。NAT の向こうの拠点に向く）、IPP System Service（PWG 5100.22-2019。複数のプリンタを束ねるシステムの操作）、IPP Encrypted Jobs and Documents（S/MIME で文書を封印したまま送る。草案）。OAuth 2.0 による認証は IPP/2.x 第 4 版と INFRA の範囲。

### 3.2 結論

1. **第 1 の出力先は IPP**（IPP/2.0 の操作、`ipps` の TLS、DNS-SD の探索）。プリンタ直結と、macOS / Linux の CUPS のキューを同じクライアントで扱う。
2. **第 1 の文書形式は PDF（ベクタのまま）**。文書はベクタの命令列で、`pdfbdf` が既に書ける。プリンタの `document-format-supported` に `application/pdf` があればそれを送り、文字と線はプリンタの解像度で出て、ジョブは小さい。部数・丁合・両面・用紙への合わせ方は IPP の属性でプリンタに任せられる。**PWG Raster は PDF を受けないプリンタのための後ろ盾**（IPP Everywhere が必須とする形式）で、`PagesRows` の帯をそのまま行の符号化に流すので、600 dpi でもページ全体の画像を持たない（2026-10-19 利用者「ベクトルデータもあるのでそちらで印刷できるならそれを優先したい」）。
3. **Apple Raster（URF）と PCLm は第 2 の形式**。URF だけを名乗る機種（AirPrint 専用の家庭用）と PCLm だけの機種（Mopria の一部）のため。どちらも Raster の書き手の変種で、PWG Raster ができれば小さい。
4. **OS のスプーラ経由は macOS / Linux の CUPS だけ**（IPP で話せる）。**Windows は IPP 直結だけにし、OS 経由（GDI / EMF / XPS）は作らない**（§3.5）。
5. **IPP USB は IPP の搬送路の 1 つ**として後で足す（Linux は `ipp-usb` が、Windows と macOS は OS が面倒を見るので、優先度は低い）。

### 3.3 プリンタに何を送るか（文書形式の選び方）

| 形式 | MIME | 要求する規格 | 性質 |
|---|---|---|---|
| PDF | `application/pdf` | IPP Everywhere 推奨。業務機はほぼ対応 | **ベクタのまま**。`pdfbdf` の出力そのまま。ジョブが小さく、文字と線はプリンタの解像度で出る。絵はプリンタの解釈器次第 |
| PWG Raster | `image/pwg-raster` | IPP Everywhere 必須、Mopria | ラスタ。ページごとに 1796 バイトのヘッダー＋行単位の RLE（行の繰り返し数と画素の連）。色空間は `sgray_1`（2 値）、`sgray_8`、`srgb_8` など。`pwg-raster-document-resolution-supported`、`-type-supported`、`-sheet-back`（両面の裏面の向き）をプリンタが名乗る |
| Apple Raster | `image/urf` | AirPrint | ラスタ。`UNIRAST` ヘッダー＋ページごとに 32 バイト。`urf-supported`（`W8`、`SRGB24`、`DM1`〜`DM4`、`RS300`、`IS1` …）で能力を名乗る |
| PCLm | `application/PCLm` | Mopria、Wi-Fi Direct | ラスタ。帯（strip）の画像を並べた制限付き PDF。`pclm-strip-height-preferred`、`pclm-raster-back-side`、`pclm-source-resolution-supported` |
| 生の PDL | `application/octet-stream` | — | 9100 と同じ素通し。後の段階 |

既定の選び方（`auto`）は **PDF があればそれ（ベクタのまま）、なければ PWG Raster、どちらも無ければ URF → PCLm**。利用者はキューごと・ジョブごとに `-format pwg-raster` などで上書きできる。PDF を優先する理由は、文書がもともとベクタであること（ラスタ化はプリンタの解像度で、プリンタがする方がよい）、ジョブが小さいこと（帳票 1 ページが 100 KB 程度。Raster の 600 dpi カラーは生で 1 ページ 100 MB 超、RLE で縮んでも桁が違う）、部数・両面・合わせ方・仕上げ（面付け・製本）をプリンタに任せられること。Raster に落とす場面は 3 つ。プリンタが PDF を受けないとき（AirPrint 専用の家庭用、Raster 専用の IPP Everywhere 機）、`pdfbdf` が書けない命令（`FILTER`、Canvas 固有の合成）を文書が使っていてビューアの絵を優先したいとき（キューの設定 `raster-if-lossy`。既定は警告のみ）、プリンタの PDF 解釈器に癖があって絵が違うとき（機種ごとの表で `-format pwg-raster` を既定にする）。

driverless の範囲でベクタのまま渡せる形式は PDF だけ（URF・PWG Raster・PCLm はすべてラスタ）。PDF 以外のベクタの経路は §3.4 にまとめた。

PDF に決まる細部:

- **部数・丁合・両面・合わせ方・ページ範囲**は IPP の属性（`copies`、`sheet-collate`、`sides`、`print-scaling`、`page-ranges`）でプリンタに任せる。受けない属性（`unsupported`）は送信側で肩代わりする（ページ範囲は `pdfbdf` にページの選択を足す、部数は文書を繰り返す）。
- **レイヤー省略**（プレ印刷用紙）は `pdfbdf` の `SkipRoles`。
- **フォント**は `pdfbdf` が埋め込む（名前で参照するフォントは送信側のマシンで探す）。埋め込みを許さないフォントはアウトラインになる。
- **暗号化した文書**は、`pdfbdf` がまだ PDF の暗号化を書かないので、`ipps` の上を平文の PDF で送る（スプールには暗号化のまま置く）。
- **同じ絵**: `pdfbdf` は書いた PDF を読み直して元の文書と同じ絵になることをラスタライザで確かめている（Ghostscript と Quartz でも）。実機の解釈器での差は機種ごとに記録する。

Raster に落ちるときの細部:

- **解像度**はプリンタの `printer-resolution-supported` から選ぶ（既定 600 dpi。モノクロ文書は 600、カラーは 300 を推奨値にして、チケットで上書き）。
- **モノクロ**は `print-color-mode=monochrome` のとき `sgray_8`、プリンタが `black_1` / `sgray_1` を名乗り、文書が線画なら 2 値（`imagebdf.Bilevel`＋`AntialiasOff`）。既定は 8 ビットのグレー（アンチエイリアスを保つ）。
- **両面**は `pwg-raster-document-sheet-back`（`normal` / `flipped` / `rotated` / `manual-tumble`）に従って裏面のページを描くときに回す（`imagebdf` の変換で済む）。
- **部数**は `copies-supported` が 1 しかない機種（Raster 専用機に多い）では送信側が繰り返す。丁合もその中で。
- **用紙へ合わせる**: ページの大きさと用紙が違うとき、`fit`（収める）/ `actual`（原寸、中央）/ `shrink`（大きいときだけ縮める）と自動回転（横長のページを縦の用紙に）。既定は `shrink`＋自動回転。
- **確定的**: 同じ文書・同じチケット・同じフォントなら同じバイト列（golden テストと再送の照合のため）。


### 3.4 PDF 以外のベクタの経路

PDF は「プリンタが直接受けるベクタ形式」の 1 つに過ぎない。ベクタの経路は 3 つの層にある（2026-10-19 利用者「PDF 以外のベクトルデータを使う仕組みも多い。Microsoft のやつとか」）。

**(a) OS の印刷パイプライン**: OS のドライバがプリンタの言語に変える。アプリはベクタで渡す。

- **Windows** には 2 つの経路がある。**GDI**（Win32 の `CreateDC` → `StartDoc` → `BeginPath` / `PolyBezier` / `ExtTextOut` / `StretchDIBits`）は、v3 ドライバがそのまま描き、v4 ドライバと IPP クラスドライバでは MXDC が GDI → XPS に変換するので、どちらでもベクタのまま届く。**XPS** は v4 ドライバの母語で、アプリは `IPrintDocumentPackageTarget`（COM）か XPS Print API で渡す。`winspool` に `RAW` で PDL を書く経路は v3 ドライバ用（PS / PCL の素通し）で、v4 ドライバでは 0 バイトのスプールになる。`XPS_PASS` は**フィルタのパイプラインを飛ばしてプリンタに直送する**型（ヘッド調整のコマンドなど向け）であって、「XPS をドライバに描かせる」経路ではない（前の版の案は誤りなので改めた）。GDI には α もグラデーションも無い（`AlphaBlend` と `GradientFill` だけ）が、GDI+（`gdiplus.dll` の平坦な C API。syscall で呼べる）には α・線形 / パスグラデーション・アンチエイリアスがあり、プリンタの DC に描ける。ブレンドモードはどちらにも無い。チケットは Print Schema の PrintTicket（XML）か `DEVMODE` に写す。
- **macOS**: Core Printing（`PMPrinter`）は cgo が要るので使わない。CUPS に PDF を渡すのが OS 経由の経路で、PostScript のプリンタには `pdftops`（Ghostscript の `ps2write`）がベクタのまま PostScript にする（透明の部分だけ画像になる）。
- **Linux**: 同じ。CUPS の driverless と Printer Application の多くはラスタ（PWG Raster）に落とすので、ベクタで届くのは PS / PCL のプリンタに Ghostscript 経由のとき。

**(b) プリンタが直接受けるベクタ PDL**: IPP の `document-format-supported` に並ぶか、9100 で送る。

| PDL | MIME | 誰の・仕様の公開 | どこで受ける | ベクタの表現力 | フォント | 透明 |
|---|---|---|---|---|---|---|
| PDF 1.7 | `application/pdf` | ISO 32000。公開 | 業務機のほぼ全部、大判の PS 機 | パス・クリップ・破線・グラデーション・パターン・画像・透明グループ・ソフトマスク・16 のブレンド | CID（TrueType / CFF）の埋め込み、ToUnicode | あり |
| PostScript 3 | `application/postscript` | Adobe。公開（PLRM 第 3 版） | 業務機（純正 PS か互換: Kyocera KPDL、Brother BR-Script、Ricoh・Canon・Xerox・Konica Minolta の PS）、大判の PS 機 | パス・クリップ・破線・shading（タイプ 1〜7）・パターン・画像 | Type 42（TrueType）、CIDFontType0C（CFF）、Type 1 | なし。平坦化が要る |
| PCL XL（PCL 6） | `application/vnd.hp-PCLXL` | HP。Class 3.0 の文書は会員向けで非公開扱い。Ghostscript の pxlmono / pxlcolor が参考実装 | レーザー機のほぼ全部（HP、Kyocera、Brother、Canon の業務機、Konica Minolta、Ricoh、Fujifilm BI） | パス（直線・ベジエ）・クリップ・破線・ROP・画像（RLE / JPEG）・パターン | ビットマップか TrueType のグリフのダウンロード | なし。平坦化 |
| PCL 5 | `application/vnd.hp-PCL` | HP。公開 | 古い機種まで | 文字＋ラスタ＋HP-GL/2（線画は HP-GL/2 のみ） | ビットマップのソフトフォント | なし |
| HP-GL/2（＋HP RTL） | `application/vnd.hp-HPGL` | HP。公開 | 大判（HP DesignJet は常に HP-GL/2。Canon imagePROGRAF と Epson SureColor は HP-GL/2 エミュレーション）。PCL 5 の中にも | 線・円弧・ベジエ・多角形の塗り（ハッチ）・ペンの幅と色・クリップ窓。画像は HP RTL | ストロークフォント（`LB`）。文字はパスにする | なし |
| OpenXPS | `application/oxps`（XPS は `application/vnd.ms-xpsdocument`） | Microsoft → ECMA-388。公開 | Windows の v4 ドライバ（`IPrintDocumentPackageTarget`）、Universal Print（PDF か XPS。XPS → PDF の変換あり）、XPS 直接印刷を持つ少数の機種 | パス・クリップ・破線・不透明度・OpacityMask・線形 / 放射グラデーション・画像 / タイル / Visual のブラシ | OpenType（TrueType / CFF）の埋め込み（難読化） | α はある。ブレンドモードは無い |
| 各社の PDL（Canon LIPS、Epson ESC/Page、Fujifilm BI ART、Ricoh RPCS / RPDL、NEC NPDL、Kyocera PRESCRIBE） | `vnd.*` | 非公開か部分公開（IANA の登録に名前だけ。ART は「ART IV Reference Manual」の参照のみ） | 国内の業務機。ただし同じ機種が PS か PCL のエミュレーションを持つことが多い | — | — | — |
| PDF/VT、JDF / JMF | `application/pdf`、`application/vnd.cip4-jdf+xml` | ISO 16612-2、CIP4。公開 | 商業印刷・大量発行の機械 | PDF と同じ。可変データの構造（DPart）とジョブの指示 | PDF と同じ | あり |

**(c) プリンタ常駐フォーム**: PCL のマクロ、PostScript の `execform`、AFP のオーバーレイ。静的部分を 1 回だけ送り、可変部分だけ毎回送る。BDF の共有 Object と同型（`system:printer-resident-forms`）。帳票の大量発行で効く。後の段階。

**結論**:

1. 直結の第 1 は PDF（§3.2）。
2. **PostScript の書き手（`vector/psbdf`）を第 2 のベクタ経路にする**。PDF を受けず PS を受ける業務機は国内外に多く（2010 年前後までの機種）、大判の PS 機、CUPS を通る経路、Windows の v3 ドライバへの `RAW` でも使える。`pdfbdf` の構造（命令列 → ページの内容、`internal/textrun` の文字の配置、影のラスタ化）をそのまま写せる。透明（グループの α・ソフトマスク・ブレンド）は平坦化が要る。
3. **Windows の OS 経由（GDI / EMF / XPS）は作らない**。Windows Ready Print の後は、Windows が扱うプリンタは IPP プリンタであり、IPP 直結で届く（§3.5）。XPS は Windows の印刷パイプラインの内部形式で、パイプラインを通らない BDF には要らない。Universal Print は PDF を受ける。GDI は IPP の無い機器（メーカードライバの USB 機、ラベル・レシートのドライバ）のためだけに残る話で、需要が出てから考える。
4. **HP-GL/2 の書き手（`vector/hpglbdf`）は CAD の出口として価値が高い**。DXF / JWW / SXF / CGM の変換器を既に持ち、大判プロッタの共通語で、読み手 `converter/hpgl` があるので往復テストができる。線画だけ（文字はパス、画像は HP RTL）。順は現場の要望次第。
5. **PCL XL は当面作らない**。仕様が非公開扱いで、機種ごとの実装の癖が大きい。PCL XL だけの機種は、同じ機種の PS エミュレーション、PWG Raster、または PCL 5 のラスタ（小さな書き手）で逃げる。
6. **各社の PDL は作らない**。仕様が無い。同じ機種の PS / PCL のエミュレーションか Raster で。
7. **平坦化の共通部（`vector/internal/flatten`）**: 透明グループ・ソフトマスク・ブレンド・`FILTER` の部分だけをラスタライザで画像にし、残りをベクタのまま残す。`pdfbdf` の影の手法の一般化で、PS・HP-GL/2・EMF（・PCL）が共有する。

形式は搬送路と独立に選ぶ。IPP でも 9100 でも CUPS でも、プリンタが名乗る形式の中から同じ順（PDF → PS → PWG Raster → URF → PCLm → PCL 5 ラスタ）で選ぶ。

### 3.5 Windows は IPP だけでよいか

よい（2026-10-19 利用者の問い。窓の杜 2026-07-02 の記事を受けて）。理由と、残る穴を書く。

- **Windows Ready Print が示すこと**: Windows 自身が、新しく追加するプリンタを IPP クラスドライバで入れ、IPP で話すようになった。WRP で入るプリンタ = IPP プリンタ。BDF は Windows のスプーラを通らずに同じプリンタへ IPP で直接話せるので、**Windows で印刷に要るのは IPP 直結だけ**になる。WPP（厳格モード）の環境では IPP プリンタしか無いので、なおさら。
- **XPS は要らない**: XPS はアプリ → スプーラ → ドライバの内部の形式で、パイプラインを通らない BDF には出番がない。Universal Print は PDF を受ける（XPS → PDF に変換するのは Microsoft の側）。OpenXPS の書き手は計画から外す。
- **GDI / EMF も要らない**: GDI はスプーラに描くための API で、IPP 直結には不要。残る用途は IPP の無い機器（メーカードライバで入れた USB のホストベース機、ラベル・レシート・ドットインパクトのメーカードライバ、`\\server\queue` の Windows 共有プリンタ）だけで、WPP の環境ではそもそも入らない機器。計画から外し、需要が出たら「OS のスプーラ経由」の枠（`print/internal/osspool`）に足す。
- **残る穴は USB 接続**: WRP は USB の Mopria 機も IPP クラスドライバ（IPP over USB）で扱うが、その USB インターフェースは Windows のドライバが握るので、BDF が自前の IPP USB で同じ機器に話すのは衝突する。Windows で USB 直結のプリンタに出したいときは、ネットワーク（同じ機種はほぼ LAN / Wi-Fi を持つ）に替えてもらうのが現実的。どうしても要るなら、IPP クラスドライバのキューに `XPS_PASS` で PDF を直送できるか（IPP ポートが `application/octet-stream` で渡すか）を 1 日の実験で確かめる価値はある。通れば GDI なしで済む。
- **強まること**: プリントサーバーを IPP Everywhere のプリンタとして見せる段階（m4）は、WRP の Windows がまさに期待する相手になる。Windows の PC からは「プリンタの追加」で BDF のサーバーのキューが入り、ドライバもアプリも要らない。WPP の環境向けに `ipps`（証明書）を持つ。

### 3.6 macOS と iOS はどうなっているか

Apple の側は Windows より先に「IPP だけ」になっている（2026-10-19 利用者の問い）。

- **macOS は 2 つとも IPP**。プリンタ直結は AirPrint（IPP。URF か PDF）。OS 経由は `cupsd` が `localhost:631` で IPP を話すので、BDF の同じ IPP クライアントで届く。USB 接続のプリンタも macOS がキューにして `localhost:631` に見せるので、Windows と違って穴にならない。メーカーの PPD ドライバは非推奨で消えていく途中（§3.1）だが、残っているキューにも CUPS 経由（IPP）で出せる。PostScript 機には CUPS が PDF → PS にするので、ベクタのまま届く。Core Printing（`PMPrinter`）の C API は cgo が要るうえに、CUPS に IPP で話せば同じことができるので使わない。macOS 向けに足すものは無い。
- **iOS / iPadOS は AirPrint しか無い**。アプリは `UIPrintInteractionController` に PDF か画像を渡し、OS が Bonjour で見つけた AirPrint プリンタに IPP で送る。ドライバも PPD も無く、CUPS は外に出ていない。BDF は Go なので iOS のアプリは作らない（gomobile の経路は非ゴール）。iOS からの印刷は 2 つの経路になる。
  1. **BDF のサーバーが AirPrint のプリンタに見える**（m4）。iPhone / iPad の共有シートからどのアプリでも BDF のサーバーに印刷でき、サーバーが最終のプリンタへ流す。AirPrint の条件: DNS-SD の TXT に `URF` キーが**必ず**要る（無いと一覧に出ない。`pdl` に `image/urf` が無くてもよい）、`image/urf` を受ける、`pdl` に `application/pdf` を出すと iOS は PDF（ベクタ）を送る（要検証）、`ipps`、`air`（認証）と `TLS` の TXT。別サブネットの iOS 端末には MDM の AirPrint ペイロード（IP か FQDN、`ResourcePath` = `rp`、`Port`、`Use TLS`。iOS 11 以降）か、広域 DNS-SD で配る。サーバーの `rp` は固定の値にしておく。
  2. **ビューアから印刷**: Safari の BDF ビューアは `window.print()`（ブラウザが描いて AirPrint）か、ビューアの「印刷」ボタンでサーバーの HTTP API にジョブを投げる（サーバーが .bdf を最終のプリンタへ。帳票ならこちらが本命）。
- **Android と ChromeOS も同じ側**（参考）。Android の既定の印刷サービスは Mopria（IPP。PDF / PWG Raster / PCLm）、ChromeOS は CUPS の driverless（IPP）。m4 のサーバーはどちらからも「プリンタの追加」で入る。
- **帰結**: m4 の IPP サーバーは、IPP Everywhere に加えて **AirPrint 互換を要件にする**（`URF` の TXT、`image/urf` の読み手、`ipps`、MDM 向けの固定 `rp`）。これで Windows（WRP）、macOS、iOS、Android、ChromeOS のすべてがドライバなしで BDF のサーバーに投げられる。

### 3.7 まとめ: 3 つの OS も業務用プリンタも IPP が基本

2026-10-19 利用者の確認。**そのとおり**。OS → プリンタの脚は Windows（Ready Print）、macOS（AirPrint と CUPS）、Linux（CUPS）のすべてが IPP になり、BDF → プリンタの脚も IPP 直結にする。業務用の複合機・レーザー（Canon imageRUNNER / Satera、Ricoh IM、Fujifilm Apeos、Kyocera TASKalfa、Konica Minolta bizhub、Sharp、Toshiba e-STUDIO、HP、Xerox、Brother、Epson）は 2010 年代半ば以降の機種がほぼ AirPrint / Mopria を名乗るので、IPP で届き、PDF をそのまま受けるものが多い。サーバーからの大量印刷も同じ IPP で、CUPS や出力管理製品がしていることと同じになる。

IPP が「基本」であって「すべて」ではない点は 4 つ。

1. **古い業務機**（2012 年ごろ以前、低価格のレーザー、IPP を持たない国内メーカーの機種）は 9100 / LPD に PS か PCL か各社の PDL しか受けない。ここは 9100 ＋ PostScript の書き手（m3）で拾い、PS も PCL も無い機種（LIPS / ART / RPCS のみ）は対象外（§3.4）。
2. **業務機の運用機能**は IPP の属性で届く範囲が機種で違う。暗証番号で留め置く認証印刷は `job-password` / `job-password-encryption`、部門 ID は `job-account-id` / `job-accounting-user-id`（いずれも IPP Enterprise Printing Extensions、PWG 5100.11）、留め置きと解放は `job-hold-until` / `Release-Job`。メーカードライバが PJL で埋めていた部門コードを IPP で受けるかは実機で確かめ、機種ごとの表に書く。IC カードで解放する認証印刷（出力管理製品）はサーバーの側の仕事で、BDF のサーバーはそこへ IPP で渡す。
3. **組織が Windows の印刷サーバー経由を義務にしている**場合、Windows 共有プリンタ（SMB）には IPP で届かない。BDF のサーバー（m2 / m4）をその代わりに置くか、対象外とする（非ゴール）。
4. **ラベル・レシート・ドットインパクト・大判プロッタ**は IPP でないことが多い（9100 / USB / シリアルに ZPL / ESC/POS / ESC/P / HP-GL/2）。m5 と HP-GL/2 の書き手で拾う。商業印刷の機械（JDF / ホットフォルダ）は対象外。

出力先の抽象は「搬送路 × 文書形式」のままにして、IPP 以外の搬送路（9100、LPD、USB、シリアル）を後から足せる形を保つ。

## 4. 要件

### 4.1 印刷指示（チケット）

名前は IPP の属性に合わせる。IPP が共通語なので、Windows の `DEVMODE` も CUPS のオプションも、ここから写せる。

| 項目 | IPP の属性 | 備考 |
|---|---|---|
| 部数、丁合 | `copies`、`multiple-document-handling` / `sheet-collate` | 送信側で繰り返す場合あり |
| 両面 | `sides`（`one-sided`、`two-sided-long-edge`、`two-sided-short-edge`） | |
| 用紙 | `media`（PWG 5101.1 の名前: `iso_a4_210x297mm`、`jis_b5_182x257mm`、`jpn_hagaki_100x148mm`、`jpn_chou3_120x235mm`、`jpn_kaku2_240x332mm`、`custom_…`）、`media-col`（`media-source` のトレイ、`media-type` の種類） | `media-ready` / `media-col-database` から選ぶ。ページの大きさからの自動選択 |
| 色 | `print-color-mode`（`monochrome`、`color`、`auto`） | |
| 品質と解像度 | `print-quality`、`printer-resolution` | |
| ページ範囲 | `page-ranges` | BDF 側で切る |
| 向きと合わせ方 | `orientation-requested`、`print-scaling`（`auto`、`fit`、`none`、`auto-fit`、`fill`） | BDF 側で描く |
| 面付け | `number-up` | BDF 側で描く（2 in 1 など）。後の段階 |
| 仕上げ | `finishings`（ホチキス、パンチ、折り）、`output-bin` | 素通し。`finishings-supported` に無ければ警告 |
| ジョブの名前と利用者 | `job-name`、`requesting-user-name` | 履歴に残る |
| 保留と優先 | `job-hold-until`、`job-priority` | スプールが扱う |
| 認証印刷と部門 ID | `job-password`（＋ `job-password-encryption`）、`job-account-id`、`job-accounting-user-id` | 業務機の留め置き印刷と部門管理（PWG 5100.11）。対応は機種次第で、`printer-attributes` に無ければ警告 |
| BDF 固有 | ビュー（`-view`）、省くレイヤーの役割（`SkipRoles`。プレ印刷用紙）、モノクロの方式（グレー / 2 値）、フォント（`FontFS` / `FontDirs`）、暗号化文書のパスワード、文書の中の区切り（`-each` のバッチを 1 ジョブにするか文書ごとにするか） | |

### 4.2 文書の描画と形式

- 入力は `.bdf`（単一ファイルと分割形式）と、`converter` が読める形式（変換してから）。サーバーは「何でも紙にする」入口。
- 既定は PDF。`pdfbdf.Write` に `SkipRoles`・`Views` を渡し、書きながら送る。ページ範囲と部数と両面は IPP の属性で。
- Raster に落ちるときは `raster/imagebdf` の `PagesRows`（帯）を Raster の行の符号化に直結し、ページ全体の画像を持たない。ページ並列はプリンタへの送信が直列なので帯の並列だけで足りる（`BandWorkers`）。
- 描画の予算（`imagebdf` の限界と警告）はそのまま。警告はジョブの記録に残す。

### 4.3 プリンタの発見と能力

- **DNS-SD**: `_ipp._tcp` と `_ipps._tcp` を mDNS（LAN）と単一 DNS（企業の DNS-SD）で引く。TXT の `rp`（パス）、`pdl`、`URF`、`TLS`、`Color`、`Duplex`、`UUID`、`adminurl`、`ty`、`kind`、`air`（認証の種類）を読む。`bdf printers` で一覧。
- **URL の直接指定**: `ipp://host[:631]/ipp/print`、`ipps://…`。CUPS のキューは `ipp://localhost:631/printers/<名前>`、CUPS 3 の `cups-local` はドメインソケット（`ipp+unix://` 相当の指定が要る）。
- **能力**: `Get-Printer-Attributes` を取り、キューの設定として保存する（起動のたびに取り直さない。`printer-config-change-time` で更新を検知）。見せる項目は用紙・トレイ・両面・色・解像度・文書形式・仕上げ・認証・TLS。
- **状態**: `printer-state`（`idle` / `processing` / `stopped`）と `printer-state-reasons`（`media-empty`、`media-jam`、`toner-low`、`door-open` …）を定期に取り、キューの表示に出す。
- CUPS の一覧は `CUPS-Get-Printers`（CUPS 固有の操作）で。Windows の一覧は `EnumPrinters`（後の段階）。

### 4.4 ジョブの送信と監視

- 送信は `Create-Job` → `Send-Document`（`last-document=true`）を基本にし、1 文書なら `Print-Job` でもよい。`Validate-Job` で属性の検証を先に。送る前に `job-id` をスプールに書き、送った後に完了を待つ（下の二重印刷の防止）。
- HTTP は chunked の送信で、Raster を描きながら送る（1 ジョブのメモリは帯の分）。`Expect: 100-continue` で属性の拒否を本文の前に受ける。
- 完了の検知は `Get-Job-Attributes` のポーリング（`job-state` が 7・8・9 になるまで）。`Create-Printer-Subscriptions` ＋ `Get-Notifications`（`ippget`）が使えれば待ち時間を減らす（任意。機種差が大きい）。
- `Cancel-Job`、`Hold-Job` / `Release-Job`（対応機のみ）。
- **二重印刷の防止**: 送信の途中で落ちたら、再起動後にまずプリンタの `Get-Jobs`（`which-jobs=completed` を含む）で、記録した `job-id` と `job-name` の状態を確かめる。プリンタが知っていれば再送しない。知らなければ（`job-id` を得る前に落ちた）再送する。
- **TLS**: `ipps` を優先（`TLS` の TXT キー、`printer-uri-supported` の `ipps`）。プリンタの証明書は自己署名が普通なので、既定は「初回に指紋を記録し、以後は一致を求める」（CUPS の `TrustOnFirstUse` と同じ）。システムの CA で検証できればそれを優先。`-insecure` は明示のときだけ。
- **認証**: なし / Basic（`uri-authentication-supported`）。OAuth は後。資格情報はキューの設定に（サーバーの設定ファイルは秘密を分けて持つ）。
- IPP の符号化は自前で書く（RFC 8010 は小さく、2〜3 千行）。属性の型（integer、boolean、enum、keyword、name、text、uri、dateTime、resolution、rangeOfInteger、collection、1setOf）、文字集合と言語、`unsupported` の応答の解釈。fuzz する。

### 4.5 スプール

- **ジョブの状態**は IPP の `job-state` をそのまま使う: `pending`（3）、`pending-held`（4）、`processing`（5）、`processing-stopped`（6）、`canceled`（7）、`aborted`（8）、`completed`（9）。理由は `job-state-reasons` の語彙（`job-incoming`、`printer-stopped`、`media-empty`、`connection-failed`（独自）…）。IPP サーバーとして見せるときも写すだけで済む。
- **置き場所**はディレクトリ（`spool/jobs/<id>/`）で、ジョブの制御ファイル（JSON。チケット、出力先、状態、プリンタ側の `job-id`、時刻、警告）と文書（`.bdf` をそのまま。内容アドレスで同じ文書は 1 回。`converter` で変換した結果の `.bdf` も置く）。書き込みは一時ファイル＋改名で壊れないように。データベースは使わない（依存を増やさない。件数は 1 サーバーで数万ジョブの規模）。
- **キュー**は出力先ごとに 1 本。FIFO に `job-priority` と `job-hold-until`。キューの一時停止と再開、ジョブの並べ替え、別のキューへの移動。
- **組（クラス）**: 複数の出力先を 1 キューにまとめ、空いている方に流す（負荷分散）か、順に試す（故障時の切り替え）。CUPS のクラスと同じ考え。
- **再送**: 接続の失敗は指数バックオフで再試行し、ジョブは `pending` のまま。プリンタが `stopped` ならキューを止めて状態を見張り、`idle` に戻れば再開。プリンタ側で `aborted` になったジョブは自動で再送しない（理由を付けて `aborted`。利用者が `retry` で戻す）。
- **描画の結果はキャッシュしない**（既定）。再送は描き直す。カラーの大きなジョブを何度も送るときのために、キューごとに「描いた Raster を残す」を選べる。
- **履歴と保持**: 完了・中止のジョブは制御ファイルを一定期間残し、文書は参照が無くなれば消す。上限（件数・バイト数・日数）を設定で。
- **アカウンティング**: ジョブごとに利用者・キュー・ページ数（`job-impressions-completed`）・用紙・時刻を JSON Lines に追記する。集計は外。
- `bdf print`（ローカル）は同じスプールの実装をメモリ上で使い、1 ジョブの完了まで待って終わる。`-spool dir` を付ければサーバーと同じディスクのスプールに置いて帰る（後で `bdf printd` が流す）。

### 4.6 プリントサーバー

- 常駐プロセス（`bdf printd`、または `cmd/bdfprintd`）。設定はファイル（YAML。既に依存にある）＋環境変数。
- **HTTP API**（JSON）: ジョブの投入（multipart で文書＋チケット、または URL / パスの参照）、一覧、詳細、取消、保留、再送、キューの一覧と状態、プリンタの発見。文書の形式は `.bdf` と `converter` の形式。既定は `localhost` に束縛し、外に開けるときは TLS とトークン。
- **CLI クライアント**: `bdf print -server http://… file.bdf`、`bdf jobs`、`bdf printers`。
- **IPP サーバー（後の段階）**: キューごとに IPP Everywhere のプリンタとして見せ、DNS-SD で広告する。Windows の IPP クラスドライバ（WRP）、macOS と iOS の AirPrint（`URF` の TXT キーと `image/urf` の受け入れが必須。§3.6）、Linux の CUPS、Android の Mopria からドライバなしで投げられる。受ける形式は PDF、PWG Raster、URF、JPEG、PNG と `converter` の形式。`pdl` に PDF を出して、送る側にベクタ（PDF）を選ばせる。これで「社内のどの OS からでも、BDF のサーバーに投げれば、どのプリンタにも出る」になり、Windows の他社製ドライバ終了（2027）への答えになる。
- **拠点をまたぐリモートプリンタ**: 直接届く（VPN か同じ LAN）なら IPP 直結。NAT の向こうは、拠点側に小さな代理（同じサーバーを「上流から取りに行く」設定で動かす）を置く。IPP INFRA の `Fetch-Job` / `Fetch-Document` の型に合わせるかは後の論点。
- **観測**: 構造化ログ（文書の中身は書かない）、`/metrics`（Prometheus のテキスト形式。任意）、キューとジョブの状態の HTTP ページ（小さな静的ページで足りる）。

### 4.7 OS ごとの対応

| | 直結（IPP） | OS のスプーラ経由 | 発見 |
|---|---|---|---|
| Windows | 唯一の経路。Windows Ready Print の後は Windows が扱うプリンタは IPP プリンタ（§3.5） | 作らない。IPP の無い機器（メーカードライバの USB 機、ラベル・レシートのドライバ、Windows 共有プリンタ）のために需要が出たら GDI を `print/internal/osspool` に足す。USB 直結の IPP 機は Windows のドライバが握るので対象外（ネットワークで） | mDNS（自前）。Windows のキューの一覧（`EnumPrinters`）は参考表示だけ |
| macOS | AirPrint のプリンタ（URF か PDF）。足すものは無い | CUPS に IPP（`localhost:631`）。USB 接続も macOS がキューにするので届く。キューの一覧は `CUPS-Get-Printers`。PS 機には CUPS が PDF → PS | mDNS（自前。Bonjour と共存） |
| iOS / iPadOS | BDF のアプリは作らない。端末からは AirPrint で BDF のサーバーへ（m4。`URF` の TXT、`image/urf`、`ipps`、MDM の AirPrint ペイロード） | — | Bonjour（サーバーが広告）か MDM |
| Linux | そのまま | CUPS 2.x に IPP（`localhost:631` か `/run/cups/cups.sock`）、CUPS 3 の `cups-local` はドメインソケット。USB は `ipp-usb` の `localhost` | mDNS（自前。Avahi と 5353 を共有するので `SO_REUSEADDR` / `SO_REUSEPORT` で束縛） |

- OS 固有のコードは `print/internal/osspool` に OS ごとの build tag で置く（`contrib/otf` と同じ流儀）。
- wasm / TinyGo のビルドに入れない。`cmd/bdfwasm` と `converter/all` から import しない。CI は 3 OS の `go vet` と `go test`（既に Windows と macOS のジョブがあれば乗る）。

### 4.8 セキュリティ

- プリンタへは `ipps` を優先し、証明書の扱いは 4.4 のとおり。
- サーバーの API は既定で `localhost`。外に開けるときは TLS とトークン、キューごとの許可。
- **暗号化した `.bdf`** はスプールに暗号化のまま置く。パスワードはジョブの投入時に受け、サーバーの鍵で暗号化して制御ファイルに置く（再起動の後も流せる）か、メモリにだけ持つ（再起動で `pending-held`、理由 `password-required`）かは論点。保護モード（`segment`）の文書は、読者の鍵を持つ側が平文の区間を投入する。
- 文書を描いた結果（Raster）は、残すと平文になる。既定で残さない。
- ログに文書の中身・パスワード・資格情報を書かない。

### 4.9 予算と信頼できない入力

- 1 ジョブの大きさ、ページ数、部数 × ページ数（印象数）に上限。描画は `imagebdf` の限界がそのまま効く。
- 同時に描くジョブの数（キューごとに 1、全体で CPU 数）と、同時に送る接続の数。
- IPP の応答の大きさ（属性の数とバイト数）に上限。悪意のあるプリンタ（偽の DNS-SD の応答）で止まらない。
- スプールのディレクトリの合計容量に上限。超えたら受け付けを `server-error-busy` 相当で断る。

### 4.10 観測と記録

- ジョブの記録: 投入・送信開始・プリンタの `job-id`・完了または中止の時刻と理由、描画の警告、送ったバイト数、印象数。
- プリンタの記録: 状態の変化（`stopped` ↔ `idle`、`media-empty` の出入り）。
- `bdf jobs` と API で見られ、JSON Lines でも追記する。

## 5. 構成案（パッケージ）

```
print/                  共通部: チケット（IPP の名前）、出力先のインターフェース（能力・送信・状態）、形式の選択、用紙の名前（PWG 5101.1）
print/ipp/              IPP のメッセージ（RFC 8010）の読み書き、クライアント（RFC 8011・PWG 5100.12）、DNS-SD の探索、CUPS 固有の操作
print/raster/           PWG Raster の書き手（imagebdf の帯から）。後で Apple Raster と PCLm
print/spool/            スプール: ジョブと制御ファイル、キュー、状態機械、再送、履歴、アカウンティング
print/server/           常駐サーバー: HTTP API、設定、プリンタの監視。後で IPP のサーバーと DNS-SD の広告
print/internal/osspool/ OS のスプーラ（CUPS の固有操作。Windows は無し。需要が出たら GDI をここに）
vector/psbdf/           PostScript 3 の書き手（pdfbdf の兄弟。第 2 のベクタ経路）
vector/hpglbdf/         HP-GL/2（＋HP RTL）の書き手（大判プロッタ。CAD の出口。converter/hpgl で往復テスト）
vector/internal/flatten 透明の部分だけを画像にする平坦化（psbdf・hpglbdf が共有）
cmd/bdf/print.go        bdf print、bdf printers、bdf jobs
cmd/bdfprintd/          サーバー（または bdf printd）
```

`print` は Go の組み込み関数 `print` と同じ名前だが、パッケージ名が組み込みを隠すだけで害はない。嫌なら `printing/`。

## 6. 段階

1. **print-m1 — IPP で手元／LAN の 1 台に印刷する（`bdf print`）**: IPP の符号化とクライアント、`Get-Printer-Attributes`、`Create-Job` / `Send-Document`、`Get-Jobs` / `Get-Job-Attributes` / `Cancel-Job`、`ipps`（初回の指紋）、DNS-SD の探索と `bdf printers`、PDF の送信（既定。ベクタのまま、`pdfbdf`）、PWG Raster の書き手（PDF を受けないプリンタの後ろ盾。帯から流す、モノクロ、両面の裏面、用紙への合わせ方）、形式の自動選択、チケットの主要項目（部数・丁合・両面・用紙・色・解像度・ページ範囲・`SkipRoles`）。CUPS のキュー（macOS / Linux）もこれで印刷できる。
2. **print-m2 — スプールとサーバー**: ディスクのスプール、キューと組、状態機械、再送と二重印刷の防止、履歴と保持、アカウンティング、HTTP API と CLI クライアント、プリンタの状態の監視、`converter` の入力の変換、`examples/print-server`（帳票を穴埋めして印刷するサンプル）。
3. **print-m3 — 第 2 の形式と古い搬送路**: PostScript の書き手（第 2 のベクタ経路。平坦化の共通部と一緒に）、Apple Raster、PCLm、`CUPS-Get-Printers` の一覧、9100 の生ソケットと LPD（PDF / PS の素通し）、IPP USB（Linux の `ipp-usb` と macOS。Windows は OS が握るので対象外）。CAD の現場があれば HP-GL/2 の書き手もここ（§3.4）。Windows の OS 経由は作らない（§3.5）。
4. **print-m4 — サーバーを IPP プリンタとして見せる**: IPP Everywhere のサーバー、DNS-SD の広告、受ける形式（PDF・PWG Raster・URF・JPEG・PNG）、AirPrint 互換（`URF` の TXT、`ipps`、MDM の AirPrint ペイロード向けの固定 `rp`。§3.6）、PWG の自己認証テスト（`ipptool`）を通す。Windows（WRP）・macOS・iOS・Android・ChromeOS からの印刷を確かめる。拠点の代理（上流から取りに行く）。
5. **print-m5 — ラベル・レシート・連帳**: ZPL / SBPL / TPCL、ESC/POS、ESC/P の書き手と、USB・シリアル・Bluetooth の搬送路。`.bdf` のページをラベルや行に写す規則は `pattern:label-sheet`、`pattern:roll-receipt`、`pattern:continuous-text-form` から。

m1 と m2 は独立に価値がある（m1 だけで手元の印刷、m2 でサーバー）。m3 以降は現場の機種で順を変える。

## 7. 検証

- **単体**: IPP の符号化の往復と fuzz、PWG Raster の書き手（自前の読み手で戻して `imagebdf` の画像と一致。golden）、チケット → IPP 属性の写し、スプールの状態機械（落ちた時点ごとの再開）。
- **結合（CI）**: `ippeveprinter`（libcups の道具。Ubuntu の `cups-ipp-utils`）を偽のプリンタとして起動し、受け取ったファイルを検査する。自前の偽プリンタ（Go）も持ち、依存なしのテストを回す。Ubuntu の runner の CUPS にキューを作って `localhost:631` 経由も確かめる。
- **実機**: 国内で多い機種（Canon、Epson、Brother、Kyocera、Ricoh、Fujifilm BI、Konica Minolta、HP）で、PDF（ベクタ）と PWG Raster、両面、JIS B5、はがき、トレイ指定を手で確かめ、PDF の絵が Raster の絵（＝ビューアの絵）と違う機種は癖として記録し、必要なら Raster を既定にする（機種ごとの表）。
- **サーバーが IPP プリンタになる段階**: PWG の IPP Everywhere 自己認証テスト（`ippeveselfcert`、`ipptool`）と、Windows の IPP クラスドライバ・macOS の AirPrint・CUPS からの印刷。

## 8. 未決の論点（決めてほしいこと）

1. **IPP の符号化を自前で書くか、`OpenPrinting/goipp`（BSD-2、符号化のみ）を使うか**。推奨は自前（この repo の流儀。fuzz しやすい。goipp は参照実装としてテストで比べる）。
2. **DNS-SD を自前で書くか、依存（`grandcat/zeroconf`、`brutella/dnssd`、`hashicorp/mdns`）を使うか**。推奨は自前の最小実装（問い合わせと PTR / SRV / TXT / A / AAAA の解釈。`golang.org/x/net/dns/dnsmessage` が既に依存にある）。広告（m4）も同じ土台。
3. **既定の文書形式**: 決定済み（2026-10-19 利用者）。PDF（ベクタのまま）を優先し、PWG Raster は PDF を受けないプリンタの後ろ盾。残る点は、`pdfbdf` が書けない命令を含む文書を自動で Raster に落とすか（案: 既定は警告のみ、キューの設定 `raster-if-lossy` で落とす）。
4. **サーバーの形**: `bdf printd` のサブコマンドか、`cmd/bdfprintd` の別コマンドか。推奨は別コマンド（`cmd/bdf` を常駐の依存で重くしない）。
5. **スプールの永続化**: ファイルだけ（推奨）か、組み込みの KV。
6. **サーバーが `converter` の入力を受けるか**: 受ける（推奨。「何でも紙にする」入口になる）。wasm への影響は無い（サーバーはサーバー）。
7. **IPP サーバー（m4）の順**: m2 の直後に上げるか。Windows の他社製ドライバ終了が近いなら前倒しの価値がある。
8. **Windows の OS 経由の経路**: 決定済み（2026-10-19 利用者）。作らない。Windows Ready Print の後は IPP 直結だけで足りる（§3.5）。USB 直結の IPP 機だけが穴で、必要なら `XPS_PASS` で PDF を直送できるかの実験をしてから考える。
9. **暗号化文書のパスワードの持ち方**（4.8）。
10. **パッケージ名**: `print/` か `printing/` か。
11. **拠点の代理の型**: 独自の「上流から取りに行く」か、IPP INFRA の操作に合わせるか。
12. **PDF 以外のベクタの書き手の順**（§3.4）: 推奨は PostScript → HP-GL/2（CAD の現場があれば）。EMF と OpenXPS は Windows の OS 経由をやめたので外した。PCL XL と各社の PDL は作らない。

## 9. 用語

- **driverless / ドライバなし印刷**: OS にメーカーのドライバを入れず、IPP で能力を聞いて標準の形式（Raster / PDF）を送る印刷。IPP Everywhere、AirPrint、Mopria。
- **PDL**: Page Description Language。PCL、PostScript、ZPL、ESC/P など、プリンタが解釈する言語。
- **PWG**: Printer Working Group。IPP とその周辺の規格を作る団体。
- **スプール**: ジョブを一時的に貯めて順に流す仕組み。`spool/jobs/` のディレクトリ。
- **キュー**: 1 つの出力先（または組）に並ぶジョブの列。
- **組（クラス）**: 複数の出力先を 1 つのキューにまとめたもの。
- **印象（impression）**: 紙の面 1 つ。両面なら 1 枚 = 2 印象。IPP の `job-impressions`。
- **TOFU**: Trust On First Use。初回の証明書の指紋を記録して以後は一致を求める。

## 10. 参照

調査は 2026-10-19 の Web 検索。`pwg.org` と `learn.microsoft.com` は作業環境から直接開けなかったので、検索の要約と二次資料による。各項目の確度は `system:ipp-driverless` の `unverified` に書いた。

- PWG IPP ワークグループ（規格の一覧）: https://www.pwg.org/ipp/
- IPP/2.x 第 4 版の承認（PWG 5100.12-2024）: https://www.pwg.org/pipermail/pwg-announce/2024/004056.html
- IPP Everywhere: https://www.pwg.org/ipp/everywhere.html 、自己認証の道具: https://github.com/istopwg/ippeveselfcert
- IPP Everywhere 2.0 の草案の案内: https://www.pwg.org/archives/ipp/2023/021421.html
- IPP Encrypted Jobs and Documents の草案: https://pwg.org/archives/ipp/2021/020939.html
- PWG 5101.1 の用紙の名前の文法: https://ftp.pwg.org/pub/pwg/informational/pwg5101.1-media-name-abnf.txt
- CUPS 3.0 の構成と計画: https://github.com/OpenPrinting/cups/wiki/CUPS-3.0 、Debian の説明: https://wiki.debian.org/CUPSNewArchitecture
- PAPPL: https://www.msweet.org/pappl/ 、libcups 3.0.2: https://openprinting.github.io/libcups-3.0.2
- CUPS の driverless の説明（AirPrint / IPP Everywhere / Mopria と形式）: https://wiki.debian.org/CUPSDriverlessPrinting
- Mopria が求める PDL（Xerox の案内）: https://support.xerox.com/en-us/article/KB0037009
- Mopria と PCLm についての議論（Michael Sweet）: https://pwg.org/archives/ipp/2024/021565.html
- `ipp-usb`: https://man.archlinux.org/man/ipp-usb.8.en
- `goipp`: https://github.com/OpenPrinting/goipp
- Apple の CUPS が 2.3.6（2022-05）で止まっていること: https://mjtsai.com/blog/2025/03/26/apples-dormant-cups/ 、OpenPrinting が上流に: https://www.phoronix.com/news/Apple-No-More-CUPS 、CUPS の PPD 拡張の非推奨（`APDialogExtension`、10.14）: https://www.cups.org/doc/spec-ppd.html
- AirPrint と CUPS（`URF` の TXT キーが必須、PDF / URF を送る）: https://connor.zip/posts/2023-06-08-airprint-with-cups 、https://www.finnie.org/2010/11/13/airprint-and-linux/
- Apple の MDM の AirPrint ペイロード（IP / FQDN、ResourcePath、Port、Use TLS）: https://support.apple.com/en-gb/guide/mdm/mdm3b4cf515/web
- macOS 26 Tahoe で古いドライバが入らない例（HP のコミュニティ）: https://h30434.www3.hp.com/t5/Printer-Setup-Software-Drivers/P1102-MacOS-Tahoe-26-2/td-p/9583270
- Windows の `XPS_PASS`（パイプラインを飛ばす直送）: https://learn.microsoft.com/en-us/windows/win32/printdocs/sending-data-directly-to-an-xps-printer 、v4 ドライバで `RAW` が 0 バイトになる: https://learn.microsoft.com/en-us/previous-versions/troubleshoot/windows/win32/v4-print-driver-raw-mode-pcl-postscript
- Windows の印刷経路（GDI と XPS、MXDC）: https://learn.microsoft.com/en-us/windows-hardware/drivers/print/windows-print-path-overview 、v4 の描画: https://learn.microsoft.com/et-ee/windows-hardware/drivers/print/v4-driver-rendering-architecture 、OpenXPS: https://learn.microsoft.com/en-us/windows/desktop/printdocs/app-support-for-openxps-printing
- Universal Print の文書形式（XPS → PDF の変換）: https://learn.microsoft.com/en-us/graph/api/resources/printsettings?view=graph-rest-1.0
- HP PCL 6（PCL XL）: https://developers.hp.com/node/21581 、HP-GL/2（HP のドライバは常に HP-GL/2）: https://support.hp.com/sg-en/document/c01121263
- CUPS の `pdftops`（Ghostscript の ps2write、透明は画像）: https://github.com/OpenPrinting/cups-filters
- IANA のプリンタ言語の登録（LIPS、ART、RPCS、PRESCRIBE …）: https://shrubbery.net/rfc/IANA/printer-language-numbers 、PWG の文書形式の MIME 表: https://www.pwg.org/hypermail/ipp/att-4439/01-Print-Document-Format-MIME-media-types-table.doc
- Windows Protected Print の解説（二次資料）: https://printops.org/windows-protected-print/ 、https://www.myq-solution.com/windows-protected-print-mode
- Windows Ready Print（Canon の案内）: https://www.canon-europe.com/support/windows-ready-print/ 、Microsoft Learn（日本語）: https://learn.microsoft.com/ja-jp/windows/modern-print/windows-ready-print 、窓の杜 2026-07-02「新規のプリンターインストールは『Windows Ready Print』が既定に」: https://forest.watch.impress.co.jp/docs/news/2121908.html
- Microsoft Q&A（IPP クラスドライバと Print Support App）: https://learn.microsoft.com/en-gb/answers/questions/1690564/how-to-setup-the-microsoft-ipp-driver-for-the-prin
