# 印刷 実装案（print-m4）

2026-10-20 の実装案。要件は `print-requirements.ja.md` §3.5・§3.6・§4.6・§6 の 4・§7 と、`.knowledge/` の `requirement:print`（facts の os・ios・android_chromeos）。print-m1（IPP で直接、`bdf print`）、print-m2（スプールとプリントサーバー、`bdf printd`）、print-m3（PostScript・Apple Raster・PCLm、9100 と LPD）は main に入った（#104、#105、#106）。2026-10-20 利用者「印刷の milestone3 まで終わりました。次の 4 の実装をお願いします」。

## 0. 範囲

m4 は「サーバーを IPP プリンタとして見せる」。`bdf printd` のキューを、どの OS からもドライバなしで「プリンタの追加」で入れられる IPP Everywhere のプリンタにし、AirPrint にも見せる。拠点の代理（NAT の向こうの拠点が上流のサーバーからジョブを取りに行く）も入れる。

1. **IPP の受け口**（`print/server` の `ipp*.go`）: キューごとに `/ipp/print/<queue>` の IPP Everywhere のプリンタ。操作は PWG 5100.14 が求めるもの（Print-Job、Validate-Job、Create-Job、Send-Document、Close-Job、Cancel-Job、Cancel-My-Jobs、Get-Job-Attributes、Get-Jobs、Get-Printer-Attributes、Identify-Printer）と Hold-Job / Release-Job。受けたジョブは `Spool.Submit` に渡し、ジョブの状態はスプールのもの（IPP の語彙のまま）を写す。圧縮（`gzip`、`deflate`）。
2. **受ける形式**: PDF・JPEG・PNG と `converter` の形式は今のサーバーの変換で、PWG Raster と Apple Raster（URF）はページの絵の `.bdf` にする（`print/raster` の読み手から）。`application/octet-stream` は中身で見分ける。`.bdf` そのもの（`application/x-bdf`）も受ける。
3. **DNS-SD の広告**（`print/ipp` の応答者）: mDNS で `_ipp._tcp` と `_ipps._tcp`、サブタイプ `_print`（IPP Everywhere）と `_universal`（AirPrint）、TXT（`txtvers`、`rp`、`ty`、`pdl`、`URF`、`UUID`、`TLS`、`Color`、`Duplex`、`kind`、`PaperMax`、`air`、`adminurl`、`note`、`priority`、`qtotal`、`product`）。5353 を Bonjour・Avahi と共有する。別サブネット向けに、企業の DNS に入れる広域 DNS-SD のレコードを書き出す（`bdf printd -dns-sd`）。
4. **AirPrint 互換と ipps**: 1 つのポートで ipp と ipps を受ける（最初のバイトで TLS を見分ける。CUPS と同じ）。証明書は設定のものか、無ければ自分で作った自己署名（スプールの隣に保存し、作り直さない）。`printer-icons`（自分で描いた PNG）。MDM の AirPrint ペイロード（IP か名前、`ResourcePath` = 固定の `rp`、`Port`、`ForceTLS`）と自己署名の証明書を入れた構成プロファイル（`.mobileconfig`）を書き出す（`bdf printd -airprint-profile`）。
5. **拠点の代理**（IPP INFRA、PWG 5100.18 の部分集合）: 上流のサーバーのキューに「出力装置」（`proxy`）を置き、拠点の `bdf printd` が `upstreams` の設定で取りに行く（Update-Output-Device-Attributes、Get-Jobs `which-jobs=fetchable`、Fetch-Job、Acknowledge-Job、Fetch-Document、Acknowledge-Document、Update-Job-Status、Update-Active-Jobs、Deregister-Output-Device）。拠点は受けたジョブを自分のスプールのキューに入れ、状態を上流へ返す。上流から下流へは接続しないので NAT を越える。
6. **確かめ方**: `ipptool` の `ipp-everywhere.test`（`ipp-1.1.test`・`ipp-2.0.test` を含む。PWG の自己認証の IPP の部分と同じ検査）を通す。macOS の CUPS に `lpadmin -m everywhere` でキューを作って印刷する（macOS からの印刷）。Windows・iOS・Android・ChromeOS は実機で利用者が確かめる（手順を docs に書く）。

**入れないもの**: IPP の通知（Create-*-Subscriptions、Get-Notifications。IPP Everywhere では任意。代理も Get-Jobs のポーリングで足りる）、複数文書のジョブ（`multiple-document-jobs-supported` は false）、Print-URI / Send-URI（サーバーが外を取りに行く穴になる）、PDF の素通し（受けた PDF を変換せずに PDF のプリンタへ渡す。下の §5）、IPP の System Service（PWG 5100.22）、OAuth。

## 1. 決めたこと（暫定で推奨どおり進める）

- **受け口の置き場所**: `bdf printd` の設定に `ipp:` を足し、HTTP API とは別のポート（既定 8631。631 は CUPS と重なり、1024 未満は権限が要る）で聞く。API は `localhost` が既定のまま、IPP は LAN に開くのが普通なので、束縛と認証を別にする。キューは全部見せる（`ipp.queues` で絞れる）。`ipp.default` のキューは `/ipp/print` でも見せる。
- **認証**: `ipp.auth` が `none`（LAN の普通のプリンタと同じ。`requesting-user-name` をそのまま使う）か `basic`（トークンの名前と秘密。`submit` の役割と、キューの許しを API と同じに見る）。トークンがあれば既定は `basic`。`basic` のときは平文の ipp を断り（HTTP 426）、`_ipps._tcp` だけを広告する。他人のジョブの取り消しは `manage` の役割だけ（`none` では同じ `requesting-user-name` のジョブだけ）。
- **ジョブの番号**: IPP の `job-id` はスプールの id、`job-uuid` はスプールの uuid（m2 で揃えた）。Create-Job は文書の前に番号を返すので、スプールの番号を先に取り（`Spool.ReserveID`、`Submission.ID`）、Send-Document が来るまでサーバーのメモリに「来る途中」のジョブとして持つ（`job-incoming`）。`multiple-operation-time-out`（300 秒）を過ぎたら捨てる（`aborted`、Get-Job-Attributes には `client-error-not-found` ではなく中止として答える）。再起動で消えるのは来る途中のジョブだけ。
- **IPP の属性 → チケット**: `copies`、`multiple-document-handling`、`sides`、`media`・`media-col`（`media-size` から PWG の名前、`media-source`、`media-type`）、`print-color-mode`、`print-quality`、`printer-resolution`、`page-ranges`、`print-scaling`、`orientation-requested`、`finishings`、`output-bin`、`job-hold-until`、`job-priority`、`job-name`。`document-password`（PWG 5100.13）は暗号化した PDF・Office・`.bdf` のパスワードで、スプールの決まりどおりメモリにだけ持つ。知らない属性と出せない値は `unsupported` の組で返し、`ipp-attribute-fidelity` が true なら断る（RFC 8011 §5.2.4）。
- **見せる能力**: キューの最初のプリンタの能力（スプールが取ったもの）から、用紙・両面・色・仕上げ・トレイ・排紙を写す。能力が分からない（9100 / LPD、まだ届かない）ときは事務機の既定（A4・A3・B5・A5・Letter・Legal・はがき、両面、カラー）。部数・ページ範囲・合わせ方・向きはサーバーが描くので、プリンタに関わらず受ける。解像度は 300 と 600 dpi（PWG Raster・URF を送ってくる側の解像度で、送り先のものではない）。
- **ラスタの受け取り**: PWG Raster と URF は 1 ページずつ読み（ページ全体の記憶は 1 ページ分）、ページを PNG（`black_1` は 1 ビット、`sgray_8` はグレー、`srgb_8` は RGB）にした絵だけのページの `.bdf` にする。ページの大きさはラスタの見出しの用紙の大きさ。16 ビットと CMYK は 8 ビットの RGB に落とす。`print/raster` に `NewPageReader`（PWG と URF を流して読む）と `Document`（`.bdf` を作る）を足し、`cmd/bdf` の変換器が `pwg-raster` と `urf` の形式として使う（`converter/all` は `print` を import しない決まりなので、変換器の側で包む）。
- **mDNS の応答者**: 自前（`golang.org/x/net/ipv4`・`ipv6` と `dnsmessage`。依存は足さない）。5353 を `SO_REUSEADDR`（BSD では `SO_REUSEPORT` も）で束縛してマルチキャストの群に入り（Go の `ListenMulticastUDP` が既にそうする）、Bonjour・Avahi と共存する。名前の衝突を確かめる問い（probe、RFC 6762 §8.1）を 3 回、告知（announce）を 2 回、終わるときは TTL 0 の告知（goodbye）。衝突したらインスタンス名に ` (2)` を足す。ホスト名は OS の名前（`<hostname>.local`）で、A / AAAA は自分のインターフェースのアドレス（同じ名前で OS の応答者も答えるが、同じ値なので衝突ではない）。
- **TXT**: AirPrint の求める `URF` は必ず入れる（無いと iOS に出ない）。`pdl` は `application/pdf,image/urf,image/pwg-raster,image/jpeg,image/png`（PDF を先頭に出して、送る側にベクタを選ばせる）。`kind=document`、`PaperMax` は能力の最大の用紙、`Color` と `Duplex` は能力から、`TLS=1.2` は ipps のとき、`air=none` か `username,password`（`basic`）。
- **ipps**: 証明書は `ipp.tls.cert` / `key`、無ければ自己署名（ECDSA P-256、10 年、SAN にホスト名・`.local` の名前・アドレス）を `ipp-tls.pem`（鍵と証明書、0600）にスプールの隣で作り、以後それを使う（指紋が変わると TOFU のクライアントが拒むので作り直さない）。
- **アイコン**: `printer-icons` は 48・128・512 px の PNG を Go で描き（プリンタの形）、`/ipp/icon-<n>.png` で返す。
- **代理の型**: 要件の論点 11 は IPP INFRA の部分集合にする（推奨）。独自の取りに行く API より、IPP の語彙のままのスプールに合い、PAPPL の `ippproxy` のような他の代理とも話せる。通知は使わず、代理は Get-Jobs（`which-jobs=fetchable`、`output-device-uuid`）を 10 秒ごとに聞く。文書は `.bdf`（`application/x-bdf`）を先に、受けない代理には PDF（上流で `pdfbdf`）を渡す。暗号化した文書は鍵のまま渡し、パスワードは Fetch-Document の答えの `document-password` で渡す（拠点のスプールもメモリにだけ持つ）。
- **代理で二重に刷らない**: 拠点は上流のジョブを自分のスプールに入れる前に、上流のジョブの uuid を `Job.Origin` に持たせて探し、既にあれば入れない。上流は Acknowledge-Job で出力装置にジョブを結び（他の装置に渡さない）、Acknowledge-Document で「渡した」を記録する（`PrinterJob.Sent`）。上流が落ちて戻ったときは、渡したジョブは状態の報告を待ち、渡す前のジョブはもう一度取らせる。拠点が落ちて戻ったときは Update-Active-Jobs で上流と突き合わせる（上流で取り消されたジョブは拠点でも取り消す）。

## 2. 構成

```
print/server/ipp.go          IPP の受け口: ハンドラ、操作の振り分け、認証、ジョブの属性
print/server/ippattrs.go     プリンタの属性（能力から）、ジョブの属性 → チケット、unsupported
print/server/ipplisten.go    1 つのポートで ipp と ipps（TLS の見分け）、自己署名の証明書
print/server/icon.go         printer-icons の PNG
print/server/advertise.go    キューの DNS-SD のサービス（TXT）、広域 DNS-SD のレコード、.mobileconfig
print/server/infra.go        上流: 出力装置、Fetch-Job …（PWG 5100.18）
print/server/upstream.go     下流（拠点）: 取りに行く、状態を返す
print/ipp/responder.go       mDNS の応答者（probe、announce、問いに答える、goodbye）
print/ipp/tags.go            INFRA の操作、output-device の属性
print/ipp/client.go          答えの後ろの文書を読む Stream（Fetch-Document）
print/raster/stream.go       PWG と URF を 1 ページずつ読む、ページの絵の .bdf
print/spool/                 ReserveID、Submission.ID・Type、Job.Origin、出力装置のプリンタ（Printer.Proxy）
cmd/bdf/printd.go            -dns-sd、-airprint-profile
cmd/bdf/printconv.go         pwg-raster と urf の変換
docs/                        print-server（IPP の受け口、各 OS からの追加、AirPrint、MDM、代理）、design §3.37
test/print/ci.sh             ipptool の ipp-everywhere.test、CUPS の lpadmin -m everywhere
```

## 3. 順番（スライス）

1. **S1** IPP の受け口（操作、属性、チケット、圧縮、来る途中のジョブ）、ラスタの受け取り、`Spool.ReserveID`、設定の `ipp:`。`ipptool` の `ipp-everywhere.test` を手元で通す。
2. **S2** mDNS の応答者と広告、広域 DNS-SD のレコード。`bdf printers`（m1 の Browse）で自分のキューが見え、`dns-sd -B` / `ippfind` でも見える。
3. **S3** ipps（1 つのポート）、自己署名の証明書、アイコン、`.mobileconfig`。macOS の CUPS に `lpadmin -m everywhere` でキューを作り、`lp` で印刷して、サーバーの先の偽のプリンタ（ipptest）に届く。
4. **S4** 拠点の代理（INFRA の上流と下流）、二重に刷らない突き合わせ。
5. **S5** docs、CI（ipptool、CUPS）、`.knowledge`。

## 4. テスト

- **受け口**: Go の IPP クライアント（`print/ipp`）で、各操作、属性 → チケット、unsupported と fidelity、圧縮、Create-Job → Send-Document、時間切れ、認証（none の利用者名、basic、426）、取り消しの許し、ジョブの状態がスプールを写すこと。サーバーの先は ipptest の偽のプリンタで、受けた PDF を確かめる。
- **ラスタ**: `print/raster` の書き手で作った PWG Raster と URF を読み戻した `.bdf` のページの絵が、元の画素と一致する。限界（大きなページ、切れた流れ）と fuzz。
- **mDNS**: ループバックの応答者に Browse で問い、TXT・SRV・A を確かめる。probe の衝突で名前が変わる。goodbye。
- **ipptool**: `ipptool -t -V 2.0 ipp://127.0.0.1:<port>/ipp/print/<queue> ipp-everywhere.test` が通る（手元の macOS と CI の Ubuntu。無ければ飛ばす）。
- **代理**: 2 つのサーバー（上流と拠点）を同じテストの中で立て、上流に投げたジョブが拠点の偽のプリンタで刷られ、完了が上流に戻る。取り消しが下流に届く。上流・拠点を途中で落としても二重に刷らない。

## 5. 後で変えられる位置に置くもの

- **PDF の素通し**: 受けた PDF を、PDF を受けるプリンタへ変換せずに渡す（`converter/pdf` を通らない分、元の PDF の絵になる）。スプールが `.bdf` でない文書を送れるようにする変更で、キューの設定（`pass-through: pdf`）として足せる。
- IPP の通知（`ippget`）は、代理のポーリングの遅れが問題になったら。
- WSD（Windows の古い発見）と LPD の受け口は作らない（Windows は Ready Print の後は IPP）。
