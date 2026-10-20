# プリントサーバー

`bdf printd` はプリントサーバーです。HTTP でジョブを受け、1 つのファイルに保存し、順に IPP のプリンタ（と 9100・LPD のプリンタ）へ送り、プリンタが刷り終えるまで見届け、届かないプリンタには送り直します。業務システムは文書を渡したら先へ進み、待つのはサーバーです。プリンタとは [`bdf print`](architecture/cli.ja.md#print-printers) と同じく、OS のドライバを使わずに IPP で話します。

```sh
bdf printd -config printd.yaml                           # 前面で動く。Ctrl-C か SIGTERM で止まる
bdf print -server http://127.0.0.1:8640 -queue office invoice.bdf     # 投入して、ジョブの id を出して帰る
bdf print -server http://127.0.0.1:8640 -queue office -wait 10m report.xlsx   # 刷り終えるまで見届ける
bdf jobs -server http://127.0.0.1:8640                    # 終わっていないジョブ
bdf printers -server http://127.0.0.1:8640                # キューとプリンタ
```

`$BDF_SERVER` と `$BDF_TOKEN` が `-server` とトークンの代わりになります。ブラウザで `http://127.0.0.1:8640/` を開くと、キュー・プリンタ・ジョブの状態が変わるたびに見えます。

## すること

- **再起動してもジョブは残る。** 投入に答える前に、ジョブを文書ごと記録します。ジョブ、文書、実績、プリンタの証明書の指紋は、すべて 1 つのファイル `spool.db` に入ります（bbolt の上の [yuruidb](https://github.com/shibukawa/yuruidb)。書くたびに `fsync`）。開けるのは 1 プロセスで、同じファイルで 2 つ目の `bdf printd` を起動するとエラーで終わります。
- **同じジョブを二重に刷らない。** 文書を送る前に試行を記録し、プリンタが `Create-Job` に答えたらすぐにプリンタのジョブ id を、プリンタが文書を全部受けたらすぐにそのことを記録します。どの時点で止まっても（落ちても、電源が切れても）、プリンタにあったジョブは何かを送り直す前にプリンタと照らし合わせます。プリンタにあれば見届け、なければ送り直し、このジョブかもしれないジョブがプリンタに 2 つ以上あるときだけ人に任せて保留します（`bdf-duplicate-check`）。プリンタでのジョブ名は試行ごとに `#<ジョブ>.<試行>` で始まり、照合はこれを探します。9100 と LPD のプリンタはジョブの記録を持たないので、プリンタが刷り始めうるもの（9100 は最初のバイト、LPD は制御ファイル）を送る直前にそのことを記録します。それより前に止まったジョブは送り直し、後で止まったジョブは人に任せて保留し（`bdf-duplicate-check`）、プリンタが受け終えたジョブは刷ったとみなします。
- **プリンタを待ち、ジョブは失敗にしない。** 電源が切れている、ほかの仕事で手がふさがっている、届かないプリンタには、10 秒、30 秒、1 分、2 分、5 分、その後は 10 分ごとに 24 時間（`retry-for`）送り直します。その間の理由は `bdf-retry-wait`、諦めて中止したら `bdf-retry-exhausted` です。止まったプリンタ（用紙切れ、紙詰まり）は受けたジョブを持ったままで、次のジョブは待ちます。プリンタが中止したジョブは自動では送り直しません。`bdf jobs retry` で送り直します。
- **プリンタには 2 つまで。** キューのジョブは、1 台のプリンタに同時に 2 つまで渡します（`max-inflight`）。1 つを刷っている間に次が届いているためです。残りはキューで待ち、プリンタに触れずに並べ替え、保留、別のキューへの移し、取り消しができます。
- **複数のプリンタのキュー。** キューのプリンタは順に試します。`failover`（既定）は動いている最初のプリンタへ、`any` は空きのある最初のプリンタへ送ります。
- **ほかの形式は変換する。** PDF、Office 文書、画像など `bdf generate` の形式は、スプールの中で `.bdf` に変換します（1 回だけ。送り直しで変換し直さない）。刷るページとして変換し、Excel のブックは Excel が印刷するとおりに、HTML と Markdown は A4 のページにします。HTML・Markdown・EPUB の画像をネットワークから取りに行くことはしません。
- **状態は IPP のもの。** ジョブは `pending`、`pending-held`、`processing`、`processing-stopped`、`canceled`、`aborted`、`completed` のどれかで、理由は IPP の `job-state-reasons`（`job-queued`、`job-printing`、`job-password-wait` など）と、bdf 独自のいくつか（`bdf-retry-wait`、`bdf-retry-exhausted`、`bdf-unsupported`、`bdf-duplicate-check`）です。
- **順番**: 優先度の高いもの（`job-priority`、1〜100、既定 50）から、同じなら古いものから。`job-hold-until` で保留します（`indefinite`、`day-time`、`evening`、`night`、`weekend`、`second-shift`、`third-shift`、または時刻まで）。

## 設定

```yaml
listen: 127.0.0.1:8640              # 既定
spool: /var/lib/bdf/spool.db        # 既定はこのファイルと同じディレクトリの spool.db
user: bdf                           # プリンタに見せる requesting-user-name
tls: { cert: /etc/bdf/tls.crt, key: /etc/bdf/tls.key }   # loopback 以外では必須
tokens:
  - { name: erp, file: /etc/bdf/erp.token, queues: [office, warehouse], roles: [submit], act-as: true }
  - { name: ops, env: BDF_OPS_TOKEN, roles: [submit, manage, admin] }
printers:
  mfp1: { url: "ipps://mfp1.example.local/ipp/print", fingerprint: "sha256:…" }
  mfp2: { url: "ipp://192.0.2.20/ipp/print", username: print, password-env: MFP2_PASSWORD }
  old-laser: { url: "socket://192.0.2.30:9100" }                     # IPP の無いプリンタ: 生のポートに PostScript
  plotter: { url: "lpd://192.0.2.31/lp1", format: pdf, pjl: true }   # LPD で PDF、PJL の見出し付き
  labels: { url: "socket://192.0.2.40:9100", profile: zd421-203 }    # ZPL のラベル
  register: { url: "serial:///dev/ttyUSB0?baud=19200", format: escpos }   # ESC/POS のレシート
  zebra-usb: { url: "ipp://localhost:631/printers/zebra", profile: zd421-300 }   # CUPS の raw キュー
queues:
  office:
    printers: [mfp1, mfp2]
    policy: failover                # か any
    defaults: { sides: two-sided-long-edge, print-color-mode: monochrome }
    format: auto                    # auto（PDF、PostScript、PWG Raster、Apple Raster、PCLm の順に受けるもの）、pdf、ps、pwg-raster、urf、pclm
    raster-if-lossy: true
    max-inflight: 2
  warehouse: { printers: [mfp2] }
retention: { documents: 7d, records: 400d }
limits: { max-document: 256MiB, max-jobs: 10000, max-bytes: 10GiB, max-impressions: 10000,
          renders: 0, retry-for: 24h, uploads: 8, events: 64 }
fonts: { dirs: [/usr/share/fonts/bdf], no-system: false }
temp-dir: /var/tmp/bdf
```

- **printers**: URL（`ipp://`、`ipps://`、`http://`、`https://`。CUPS のキューは `ipp://host:631/printers/<名前>`）、Basic 認証を求めるプリンタの利用者とパスワード（`password-env` か `password-file`）、`ipps` のプリンタの証明書。`fingerprint` で固定し、書かなければ最初に見た証明書をスプールに記録して以後はそれを求めます（CUPS と同じ）。`insecure: true` はどの証明書も受けます。IPP の無いプリンタは `socket://host:9100` か `lpd://host/queue` で、受ける形式 `format`（`ps` が既定、か `pdf`）と `pjl`（文書の前の PJL の見出し。既定は PDF のときだけ）を書きます。チケットをファイルの中に書いた PostScript か PDF を送り、ジョブの状態は返らず、動いていると分かっている間は接続を試しません。ラベル・レシート・連帳のプリンタ（[下](#ラベルレシート連帳)）は、`format` に言語（`zpl`、`sbpl`、`tpcl`、`escpos`、`escp`）か `profile` に機種を書き、`socket://`、`lpd://`、`serial://`、`usb://`、`winspool://` で、または IPP で CUPS の raw キュー（データをそのまま渡す）に送ります。
- **queues**: プリンタの順、方針、チケットの `defaults`（IPP の名前: `copies`、`sides`、`media`、`media-source`、`media-type`、`print-color-mode`、`print-quality`、`printer-resolution`、`print-scaling`、`orientation-requested`、`finishings`、`output-bin` など。ジョブのチケットがこの上に重なる）、送り方。
- **retention**: 終わったジョブの文書を残す期間（7 日）と、記録（実績）を残す期間（400 日）。
- **limits**: 文書の大きさ（256 MiB。超えると 413）、終わっていないジョブの数とその文書の量（1 万件と 10 GiB。超えると 507）、ジョブのページ数、同時に変換・描画するジョブの数（0 は CPU 数）、プリンタに送り直す期間、同時に受ける投入とイベントの接続の数。
- 期間は `s`、`m`、`h`、`d`、大きさは `KiB`、`MiB`、`GiB`（`KB` なども）。知らないキーは誤りにします（打ち間違いを黙って無視しない）。相対パスは設定ファイルのディレクトリからです。設定は起動時に読みます。

`bdf printd -config printd.yaml -check` は設定を読み、各プリンタにできることを聞いて表にします。スプールは開きません（答えないプリンタがあれば終了コード 1）。

`ipp:` はキューを IPP プリンタとして見せ（[後述](#ipp-プリンタとして見せる)）、`upstreams:` と `proxy: true` のプリンタは NAT を越えてジョブを取りに行かせます（[後述](#nat-の向こうの拠点)）。

## 誰が何をできるか

- **loopback で、トークンが無いとき**（既定）は、その機械の誰でも何でもできます。ジョブの利用者は `X-BDF-User`（`bdf print` はプロセスの利用者を入れる）です。
- **loopback 以外**では、TLS とトークンが無ければ起動しません（`insecure: true` で明示したときを除く）。
- **トークン**は名前（そのジョブの利用者）、16 文字以上の秘密（`file`、`env`、または設定に直接書く `value`。勧めません）、投入できるキュー（無ければ全部）、役割を持ちます。役割は `submit`（投入と自分のジョブの操作）、`manage`（他人のジョブの操作と実績）、`admin`（キューの停止・再開・受け付けの停止と再開、プリンタの探索、写し）、`proxy`（拠点のサーバーとしてキューのジョブを取りに行く）。`Authorization: Bearer …` で送り、比べるのは定数時間です。
- **`act-as`** があるトークンは `X-BDF-User` でジョブの利用者を名乗れます。利用者に代わって印刷する業務システムのためです。`manage` が無ければ、名乗った利用者のジョブを見て操作します。
- ログには文書、パスワード、トークン、プリンタの URL の利用者とパスワードを書きません。

## API

JSON で `/api/v1` の下。エラーは `{"error": {"code": "…", "message": "…"}}` と HTTP の状態です: 400 `invalid` か `bad-document`、401 `unauthorized`、403 `forbidden` か `wrong-password`、404 `not-found` か `no-queue`、409 `not-possible`（ジョブの状態ではできない）、413 `too-large`、415 `unknown-format`、503 `not-accepting` か `busy`、507 `full`。

| メソッドと道 | |
|---|---|
| `POST /queues/{q}/jobs` | 投入。`multipart/form-data` で `document`（ファイル）と、順不同で `ticket` と `options`（JSON）、`name`、`password`、`priority`、`hold-until`、`hold-until-time`。201 とジョブ |
| `GET /jobs?queue=&user=&state=&after=&limit=` | 終わっていないジョブ。`state=all` で終わったものも、または状態で絞る。`after` と `next` で id 順に続きを取る |
| `GET /jobs/{id}` | ジョブ |
| `POST /jobs/{id}/cancel`・`hold`・`release`・`retry`・`password`・`move` | 取り消し（プリンタでも）、保留、解除（`{"password"}` も）、送り直し、パスワードを渡す、移す（`{"queue"}`） |
| `PATCH /jobs/{id}` | `{"priority": n}` |
| `GET /queues`、`GET /queues/{q}` | キューと、状態ごとのジョブの数 |
| `POST /queues/{q}/pause`・`resume`・`reject`・`accept` | 送るのを止める・再開する・受け付けを止める・再開する（admin） |
| `GET /printers`、`GET /printers/{p}` | プリンタと状態、1 台のできること |
| `GET /discover?timeout=` | DNS-SD でネットワークの IPP プリンタを探す（admin） |
| `GET /accounting?from=&to=&format=jsonl\|csv` | 終わったジョブ: 誰が何をどこで何枚刷ったか（manage） |
| `GET /backup` | `spool.db` の一貫した写し（admin） |
| `GET /events?job=` | Server-Sent Events。`job`、`queue`、`printer` とそれぞれの JSON |

```sh
curl -F document=@invoice.bdf -F 'ticket={"sides":"two-sided-long-edge","copies":2}' \
     -H "Authorization: Bearer $TOKEN" https://print.example.local:8640/api/v1/queues/office/jobs
```

`/api/v1` の外に、`/healthz`（動いている間は 200）、`/metrics`（Prometheus のテキスト。キューと状態ごとのジョブ、終わったジョブ、枚数、プリンタの状態、送ったバイト、送り直し、スプールの大きさ）、`/`（状態のページ。読むだけで、トークンはページで渡せます）があります。

Go からは [`print/server/client`](https://github.com/shibukawa/bdf/tree/main/print/server/client) が API のクライアントです（`Submit`、`Job`、`Jobs`、`Cancel` など、`Wait`、`Events`、`Accounting`）。bdf のほかの部分に依りません。プログラムにサーバーを組み込むこともできます。[`print/server`](https://github.com/shibukawa/bdf/tree/main/print/server) が設定のスプールを開き、`http.Handler` を返します（[サンプル](examples/print-server.ja.md)）。

## 暗号化した文書

暗号化した `.bdf` はスプールでも暗号化のままです。ジョブと一緒に渡したパスワード（`bdf print` は `$BDF_PASSWORD` を、先に確かめてから送る）はメモリにだけ持ちます。再起動の後のジョブは、渡し直す（`bdf jobs password <id>`）まで待ちます（`pending-held`、`job-password-wait`）。パスワードの要る PDF や Office 文書はそのパスワードで変換し、変換した `.bdf` も同じパスワードで暗号化して置きます。`bdf print` は、TLS の無い接続で別の機械へパスワードを送りません（`-allow-plaintext` のときを除く）。

## IPP プリンタとして見せる

設定に `ipp:` があると、各キューは IPP Everywhere のプリンタにもなり、どの OS もドライバなしで追加できます。Windows（Windows Ready Print と同じ IPP クラスドライバ）、macOS と iOS（AirPrint）、Android（既定の印刷サービスの Mopria）、ChromeOS と Linux（CUPS）。送ってくるのは PDF（先頭に出すので、ベクタのまま送られる）、Apple Raster か PWG Raster、JPEG か PNG で、サーバーは受けたものを投入と同じく `.bdf` にしてキューのプリンタで刷ります。ラスタのページは、描かれた用紙の大きさの絵になります。

```yaml
ipp:
  listen: ":8631"              # 既定。ipp と ipps を同じポートで
  host: print.example.com      # プリンタの URI と証明書の名前。既定は機械の名前の .local
  auth: none                   # none（LAN のプリンタと同じ）か basic（トークンの名前と秘密）。トークンがあれば basic
  require-tls: false           # ipp を断る（HTTP 426）。basic のときはいつも
  tls: { cert: /etc/bdf/ipp.crt, key: /etc/bdf/ipp.key }   # 無ければ自己署名を作る（スプールの隣の ipp-tls.pem）
  queues: [office, warehouse]  # 見せるキュー。既定は全部
  default: office              # /ipp/print にも
  location: "2F"               # printer-location。キューの location が優先
  organization: Example Inc.
  dns-sd: true                 # LAN に広告する（mDNS）。既定
queues:
  office: { printers: [mfp1], info: "Office printer", location: "2F north" }
```

- **プリンタ**: `ipp://host:8631/ipp/print/<キュー>` と `ipps://…`（同じポート。接続の最初のバイトで TLS を見分けます）。操作は IPP Everywhere が求めるもの（Print-Job、Validate-Job、Create-Job、Send-Document、Close-Job、Cancel-Job、Cancel-My-Jobs、Get-Job-Attributes、Get-Jobs、Get-Printer-Attributes、Identify-Printer）と Hold-Job・Release-Job。ジョブの id はスプールのもの（`bdf jobs` と同じ番号）、状態もスプールのものです。CUPS の `ipptool` の IPP/1.1・IPP/2.0・IPP Everywhere の属性の検査が通ります（`test/print/ipp-everywhere-server.test`）。
- **ジョブが頼むこと**はチケットになります。`copies`、`sides`、`media` か `media-col`（大きさはキューの用紙に合わせる）、`print-color-mode`、`print-quality`、`printer-resolution`、`page-ranges`、`print-scaling`、`orientation-requested`、`finishings`、`output-bin`、`job-hold-until`、`job-priority`。キューにできないことは返事で伝え、`ipp-attribute-fidelity` があればジョブを断ります。`document-password` は暗号化した PDF・Office 文書・`.bdf` を開き、投入と同じくメモリにだけ持ちます。
- **できると名乗ること**はキューの最初のプリンタのもの（用紙、トレイ、仕上げ、カラー、両面）に、サーバーが自分でするもの（部数、ページ範囲、合わせ方）を足したものです。まだ答えていないプリンタと、9100・LPD のプリンタは事務機の用紙で見せます。
- **誰が**: `auth: none` ではネットワークの誰もが、クライアントの名乗る `requesting-user-name` で印刷し、その名前のジョブを取り消せます。他人のジョブは名前を隠して一覧に出ます。`auth: basic` ではトークンの名前と秘密を求め（OS が 1 回聞いて覚えます）、秘密が TLS を通るように別の機械からの ipp を断り、`_ipps._tcp` だけを広告します。トークンのキューと役割は API と同じに効きます。
- **DNS-SD**: 各キューを `_ipp._tcp` と `_ipps._tcp` で、サブタイプ `_print`（IPP Everywhere）と `_universal`（AirPrint）を付けて広告します。ホスト名はサーバー自身のもの（`<機械>-bdf.local`。機械の名前は OS の応答者のもの）で、TXT にはクライアントが選ぶ材料（`pdl`、`URF`、`Color`、`Duplex`、`PaperMax`、`TLS`、`air` …）を入れ、後ろのプリンタが答えたら更新します。mDNS のポートは Bonjour・Avahi・Windows と共有します。
- **別のサブネット**: mDNS はリンクの中だけです。`bdf printd -config printd.yaml -dns-sd example.com` は同じレコードを社内の DNS のゾーン向けに書き出します（広域 DNS-SD。`ipp.host` はゾーンでのサーバーの名前に）。`bdf printd -config printd.yaml -airprint-profile print.example.com > printers.mobileconfig` は iPhone・iPad 向けの構成プロファイル（各キューの AirPrint の項目と、自己署名の証明書を信頼するルート）を書き出し、MDM で配れます。プリンタのページ（`https://host:8631/`）にも、開いたアドレス向けの同じプロファイルへのリンクがあります（https だけ。平文の HTTP では別の機械が自分の証明書に差し替えられるため）。サーバーが作る証明書は自分で署名した末端の証明書で、他の証明書には署名しません。信頼した機器が信じるのはこのサーバーだけです。

プリンタの追加:

| OS | |
|---|---|
| Windows 10・11 | 設定 → Bluetooth とデバイス → プリンターとスキャナー → デバイスの追加。キューが一覧に出ます（mDNS）。別のサブネットは「IP アドレスまたはホスト名を使ってプリンターを追加する」→ IPP デバイス → `https://host:8631/ipp/print/office`（か `http://`）。IPP クラスドライバが使われ、ドライバは入りません |
| macOS | システム設定 → プリンタとスキャナ → プリンタを追加。一覧に出ます（AirPrint）。別のサブネットは IP のタブで、プロトコル IPP、アドレス `host:8631`、キュー `ipp/print/office`。シェルからは `lpadmin -p office -E -v ipps://host:8631/ipp/print/office -m everywhere` |
| iPhone・iPad | 同じサブネットなら印刷のときに一覧に出ます（AirPrint）。別のサブネットは上のプロファイル |
| Android | 既定の印刷サービス（Mopria）の一覧に出ます。別のサブネットはサービスの設定でアドレスから |
| ChromeOS | 設定 → プリンタ → プリンタを追加。一覧に出るか、アドレスとプロトコル IPP / IPPS、パス `ipp/print/office` で |
| Linux | `lpadmin -p office -E -v ipps://host:8631/ipp/print/office -m everywhere`、またはデスクトップのプリンタ設定（一覧に出ます） |

## NAT の向こうの拠点

本社のサーバーから届かない拠点（NAT の内側で VPN が無い）は、拠点に自前の `bdf printd` を置き、本社のキューのジョブを取りに行って拠点のプリンタで刷らせます。つなぐのは拠点からで、本社からは拠点へつなぎません。両者が話すのは IPP の共有インフラ拡張（PWG 5100.18）なので、同じ規格の他の代理も取りに来られます。

```yaml
# 本社: プリンタが拠点の代理であるキュー
printers:
  osaka: { proxy: true }
queues:
  osaka: { printers: [osaka] }
ipp: { listen: ":8631", host: print.example.com }
tokens:
  - { name: osaka-proxy, file: /etc/bdf/osaka.token, roles: [proxy], queues: [osaka] }
```

```yaml
# 拠点
printers:
  mfp: { url: "ipp://192.0.2.10/ipp/print" }
queues:
  local: { printers: [mfp] }
upstreams:
  hq:
    url: ipps://print.example.com:8631/ipp/print/osaka
    queue: local
    username: osaka-proxy
    password-file: /etc/bdf/osaka.token
    poll: 10s              # 既定
    # fingerprint: "sha256:…"   本社の証明書を固定する。無ければ最初に見たものを覚える
```

- 本社のキューに投入されたジョブ（API、`bdf print -server`、IPP のどれでも）は本社で `job-fetchable` として待ちます。拠点は `poll` ごとに聞き、ジョブを受け（`Fetch-Job`、`Acknowledge-Job`）、文書（本社が持つ変換済みの `.bdf`。パスワードがあれば一緒に受け、拠点でもメモリにだけ持つ）を取り、自分のキューに入れて、持ったことを本社に伝えます（`Acknowledge-Document`）。以後はジョブの状態を報告し（`Update-Job-Status`）、本社の `bdf jobs` に拠点での印刷中・完了・印象数が出ます。
- **二重に刷らない。** ジョブは受けた拠点に結びつき、他の拠点には渡しません。拠点は本社のジョブの uuid を自分のジョブに持たせ、既に持っているジョブ（受けてから伝えるまでの間に止まった）は受け直さずに伝えるだけにします。本社が渡したのに拠点が持っていないジョブ（拠点がスプールを失った）は本社で中止にし、人の再送を待ちます（自動では送り直しません）。
- **本社で取り消したジョブ**は拠点でも取り消します（次の報告か、`poll` ごとに互いのジョブを突き合わせる `Update-Active-Jobs` で）。
- 本社は拠点のプリンタをそのキューのプリンタとして見せます（状態・用紙・カラーは拠点の報告から。2 分報告が無ければ止まっている扱い）。そのキューの IPP プリンタも同じです。
- 拠点のトークンの役割は `proxy`（キューの文書を読むので）。`queues` のキューだけに効き、ipps の Basic 認証でだけ使えます。代理のプリンタは 1 つのキューだけのプリンタです。拠点が受けられないジョブ（読めない文書、大きすぎる）は中止として報告し、本社で retry したジョブは取り直します。

## ラベル・レシート・連帳

ラベル・レシート・連帳のプリンタのキューも、ほかのキューと同じもの（穴埋めした帳票、PDF、画像 …）を受け、そのプリンタの言語で刷ります（細部は [`bdf print`](architecture/cli.ja.md#print-printers)）。

| 言語 | プリンタ | 送るもの |
|---|---|---|
| ZPL | Zebra と ZPL II の互換機 | ラベルごとに 1 つのフォーマット。ドットは `^GFA`（ZPL の ASCII 圧縮）、QR コードとバーコードは `^BQ`・`^BC`・`^BE` …、部数は `^PQ` |
| SBPL | SATO | ラベルごとに STX と ETX で囲んだ `<A>` … `<Z>`。`<G>H` のグラフィック、`<2D30>`・`<B>`・`<BG>` の記号、`<Q>` の部数 |
| TPCL | 東芝テック | `{D` の大きさ、ニブルモードの `{SG;`、`{XB`・`{RB` の記号（Code 128 は除く）、`{XS;I` の発行 |
| ESC/POS | Epson TM、レシートプリンタ | 描いた帯から順に `GS v 0`、`GS V` のカット、`ESC p` のドロワ |
| ESC/P | Epson VP、ドットインパクト | 24 ピンのビットイメージ（`ESC *`、180 dpi）、`ESC C` の用紙長、`ESC N` のミシン目スキップ、改ページ |

キューの `defaults` には、こうしたプリンタへの指示もチケットの名前で書けます: `print-darkness`、`print-speed`、`media-tracking`（IPP の名前）と、`bdf-label-gap`、`bdf-cut`、`bdf-cash-drawer`、`bdf-perforation-skip`、`bdf-no-printer-barcodes`、`bdf-print-mode`（`raster` か `text`。テキストのテンプレートの帳票の文字をプリンタのフォントで刷る。ESC/P はそうした帳票を既定でテキストモードで刷る）。部数はプリンタが作り（`^PQ`、`<Q>`、`{XS`）、レシートと連帳は繰り返して送ります。

命令は各社のプログラミングリファレンス（Zebra の ZPL II、SATO の CL4NX Plus と CT4-LX、東芝テックの B-EX4・B-FV4・BX400 の外部インタフェース仕様、Epson の ESC/POS と ESC/P）に従っています。確かめたのは自前の読み手でドットに戻すところまでで、実機では確かめていません。大量に刷る前に 1 枚試し、記号がおかしければ `bdf-no-printer-barcodes: true` でドットとして描かせてください（203 dpi でも読めます）。SBPL は SATO のステータスのプロトコルが求めるとおり STX と ETX で囲みます。USB のプリンタには、Windows はドライバのキュー（`winspool:///名前`、`RAW` で渡す）、macOS は CUPS（`ipp://localhost:631/printers/<キュー>`）、Linux はどちらかか `usb:///dev/usb/lp0` で届きます。

## 常駐させる

`bdf printd` は前面で動き、ログを標準エラーに書き（JSON。人が読むなら `-log text`）、SIGINT か SIGTERM（Windows は Ctrl-C か Ctrl-Break）で止まります。止まるときは受け付けをやめ、送っている文書を 30 秒まで待ってからスプールを閉じます。途中で打ち切ったものは次の起動で直します。サービスとしての登録は持たず、常駐は OS のサービス管理に任せます。

**Linux（systemd）**: `/etc/systemd/system/bdf-printd.service` を置いて `systemctl enable --now bdf-printd`。

```ini
[Unit]
Description=bdf print server
After=network-online.target
Wants=network-online.target

[Service]
# 専用の利用者（useradd --system bdf）。トークンのファイルを読めること
User=bdf
ExecStart=/usr/local/bin/bdf printd -config /etc/bdf/printd.yaml
Restart=on-failure
# printd.yaml の spool: /var/lib/bdf/spool.db
StateDirectory=bdf
TimeoutStopSec=45

[Install]
WantedBy=multi-user.target
```

**macOS（launchd）**: `/Library/LaunchDaemons/io.github.shibukawa.bdf-printd.plist` を置いて `sudo launchctl bootstrap system /Library/LaunchDaemons/io.github.shibukawa.bdf-printd.plist`。

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key><string>io.github.shibukawa.bdf-printd</string>
  <key>ProgramArguments</key>
  <array>
    <string>/usr/local/bin/bdf</string><string>printd</string>
    <string>-config</string><string>/usr/local/etc/bdf/printd.yaml</string>
  </array>
  <key>RunAtLoad</key><true/>
  <key>KeepAlive</key><dict><key>SuccessfulExit</key><false/></dict>
  <key>StandardErrorPath</key><string>/usr/local/var/log/bdf-printd.log</string>
</dict>
</plist>
```

**Windows**: [NSSM](https://nssm.cc/) や [WinSW](https://github.com/winsw/winsw) のようなサービスの包みで動かします。

```bat
nssm install bdf-printd "C:\Program Files\bdf\bdf.exe" printd -config "C:\ProgramData\bdf\printd.yaml"
nssm set bdf-printd AppStopMethodConsole 30000
nssm set bdf-printd AppStderr "C:\ProgramData\bdf\printd.log"
nssm start bdf-printd
```

```xml
<!-- bdf-printd.xml。bdf-printd.exe に改名した WinSW の隣に置き、bdf-printd.exe install、start -->
<service>
  <id>bdf-printd</id>
  <name>bdf print server</name>
  <executable>C:\Program Files\bdf\bdf.exe</executable>
  <arguments>printd -config "C:\ProgramData\bdf\printd.yaml"</arguments>
  <stoptimeout>30 sec</stoptimeout>
  <onfailure action="restart" delay="10 sec"/>
  <log mode="roll"/>
</service>
```

包みを使わないなら、タスク スケジューラで起動時に動かせます: `schtasks /create /tn bdf-printd /sc onstart /ru SYSTEM /tr "\"C:\Program Files\bdf\bdf.exe\" printd -config C:\ProgramData\bdf\printd.yaml"`（落ちたときに起動し直しはしません）。

## スプールの手入れ

- **写し**: `GET /api/v1/backup`（admin）が、動いているサーバーの `spool.db` の一貫した写しを返します。写しはそのままスプールとして開けます。まず一時ファイル（`temp-dir`）に作るので、遅いダウンロードがスプールを止めることはありません。
- **大きさ**: 保持期間を過ぎた文書を消しても、ファイルは縮みません（空いた所は使い回します）。止めた状態で `bdf printd -config printd.yaml -compact` を実行すると、その分を詰めて書き直します。
- **共有しない**: 1 つの `spool.db` は 1 つのサーバーのものです。複数のサーバーが開くネットワークのファイルシステムに置かないでください。
