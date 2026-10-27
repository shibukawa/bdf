# Docflow の FTP / SFTP 対応（2026-10-27）

2026-10-27 の利用者の依頼を `/knowledge-memory` で整理した叙述。要件の概念は `.knowledge/` の `requirement:docflow-file-transfer`（問いと暫定の選択はそこ）、他の転送プロトコルの調査は `survey:transfer-protocols`、最初の段は `milestone:fl-4a`。

利用者「ftp, sftp で scp なしでお願いします。/knowledge-memory で要件整理をお願いします。他に対応すべき転送プロトコルはありますか？ hulft はプロプライエタリなのでなしでよいかと」。直前に Claude が「Docflow は FTP / SFTP / SCP のどれにも対応していない。SFTP が業務の本命、SCP は SFTP があれば省ける、入口は見回りか常駐かが論点」と答えていた。

## 0. 今の形（2026-10-27、v0.2.0 の後。コードで確かめた）

- **入口**は upload（HTTP の API）、cups（IPP、自分のポート）、s3-trigger（外の Lambda / EventBridge が POST する）。**出口**は http（push）、smtp、s3、cups-out、printer、download（取りに来る）。`docflow/io/{http,smtp,s3,ipp,notify}` が `init` で登録する。
- 雛形の `kinds/folder`（`docflow/web/scaffold`）がサーバーのディレクトリに書く出口の例。一時名に書いて rename、`overwrite=false` なら既にあるファイルは「書いた扱い」（結果の分からなかった再試行のため）。FTP / SFTP の出口はこの形をそのまま遠隔に持っていくもの。
- 出口の口は `Kind.Send`（delivery、at-least-once、冪等の鍵 = delivery の ID）。入口の口は `Keyed` + `Fetch`（参照を POST で受けて読む）と `Serve`（自分のポートの `http.Handler`。HTTP だけ）。定期の仕事は `Tick`（POST /api/tick をスケジューラが叩くか、プロセスのタイマー）で、tickDeliveries・tickJobs・tickInbox がある。**種類ごとの tick の口はまだ無い**（`requirement:docflow-framework` の kind-collections-tickers、fw-4 の予定）。
- 設定は fw-1 で `Node.Settings` + `Kind.Fields`（string・path・bool・enum・credential・duration …）。既存のアダプタ（smtp・s3・http）は平らな `Node` の欄のまま（framework の open_questions.adapters-to-settings）。
- 資格は名前だけを Node に置き、値は `DOCFLOW_CRED_<name>` かファイル `CredentialDir/<name>`（複数行でもよい）。
- 依存：`golang.org/x/crypto v0.58.0` が go.mod にある（`ssh` が使える）。FTP / SFTP の部品は無い。
- ftp は fl-1・fl-2 の non_goals で「fl-4」に送られていた。fl-4 は人の操作のノード・人の入力の入口・split・dedupe・流量と同じ束で、まだ概念が無い。

## 1. 決めたこと（利用者、2026-10-27）

- **FTP と SFTP** を入れる。**SCP は入れない**。**HULFT は入れない**（プロプライエタリ）。
- Claude の読み：FTP には **FTPS**（RFC 4217 の明示的 TLS、`AUTH TLS` + `PROT P`）を含む。`requirement:document-flow` は最初から「FTP(S)」と書いていた。平文の FTP も残す（複合機・NAS・古い取引先のサーバー）が、エディタで注意を出す。

## 2. 出口（sink）：`ftp` と `sftp`

folder の遠隔版。届いたペイロードを相手のサーバーのパスに置く。

- **設定**（`Node.Settings` + `Kind.Fields`。新しい種類は最初から fw-1 の形で、平らな Node の欄は使わない）
  - `url`：`ftp://user@host:21/base/dir`、`ftps://user@host:990/base`（暗黙の TLS）、`sftp://user@host:22/base`。ユーザー名は URL に、パスワードは入れない（入っていたら PUT で断る）。
  - `credential`：資格の名前。値は FTP ならパスワード、SFTP ならパスワードか **PEM の秘密鍵**（`-----BEGIN` で見分ける。鍵はファイルの資格に置く。暗号化された鍵は当面断る）。
  - `host_key`（sftp だけ、必須）：サーバーのホスト鍵（OpenSSH の公開鍵の 1 行か `SHA256:` の指紋）。**`Kind.Discover` が接続して取ってくる**ので、エディタで「取得」して確かめてから保存する（人の目を通した TOFU）。known_hosts のファイルは後。鍵の検証を外す設定は作らない。
  - `tls`（ftp だけ）：`explicit`（既定。`AUTH TLS` が無ければ失敗）・`none`（平文。エディタに注意）。`ftps://` は暗黙。自己署名の相手は `tls_fingerprint`（証明書の SHA-256）で留める。
  - `path`：`{{meta.<key>}}` を展開。既定 `{{meta.$job}}/{{meta.$name}}`（s3・folder と同じ）。`..` と頭の `/` は落とす。
  - `mkdirs`：親のディレクトリを作る（既定 on）。`overwrite`：既にあるファイルに上書き（既定 off：あるものは「書いた扱い」）。
- **送り**：一時名 `.<delivery の ID>.part` に上げて、`RNFR`/`RNTO`（FTP）か `SSH_FXP_RENAME`（SFTP。`posix-rename@openssh.com` があればそれ）で最終の名前へ。読み手が途中のファイルを見ない。再試行は、最終の名前があれば done、無ければ一時名から上げ直す。FTP はバイナリ（`TYPE I`）だけ、受動（`EPSV`、無ければ `PASV`）だけ。
- 記録するのは置いたパス（delivery の result）。資格の値はログにも API にも出さない（fl-2 の制約のまま）。
- 入力は any の多相（smtp・s3・http と同じ。.bdf のままも、描画を前に置いた PDF も送れる）。

## 3. 入口（source）：`ftp-poll` と `sftp-poll`

FTP / SFTP のサーバーは何も知らせてこないので、S3 のように外からの POST で起こす形は取れない。**tick で見回る**（Claude の推奨。利用者の答え待ち）。pm-1 の「常駐に頼らない」はそのまま：見回りは Tick の中で行い、サーバーレスでもスケジューラの POST /api/tick で動く。fl-2 の「S3 の見回りはしない」は S3 にイベントがあるからで、矛盾しない。

- **設定**：`url`・`credential`・`host_key`・`tls`・`tls_fingerprint` は出口と同じ。`pattern`（取る名前の glob、既定 `*`）。`every`（見回りの間隔、既定 1 分。Tick がそれより密でも間隔は守る）。`after`：取ったファイルをどうするか（`delete`（既定）・`move`（`done_dir` へ。相手の了解がいる））。`settle`（この時間の間に大きさか更新時刻が変わったファイルは次に回す、既定 10 秒。`.part`・`.tmp`・`.filepart` の名前も飛ばす）。`max_per_tick`（1 回に取る数、既定 20）。入口の共通の欄（accept・output・formats・meta）も。
- **取り方（at-least-once、2 つのインスタンスが同じファイルを取らない）**：まず **rename で取り置き**（`<name>.docflow-<node>.taking`）。rename は両プロトコルで原子的なので、負けた方は飛ばす。それから読んで（MaxDocument を超えるものは取らずに TickReport に出す）ジョブにし（upload と同じ道。`$source` = `ftp` / `sftp`、`ftp.path` / `sftp.path` と名前・大きさ・更新時刻をメタデータに）、最後に `after` のとおり消すか移す。取り置きのまま残ったファイル（途中で落ちた）は次の tick で同じノードのものだけ拾い直す。
- **種類ごとの tick の口**が要る：`Kind.Tick(ctx, m, n, rep)`。fw-4 の kind-collections-tickers の tick の方だけを前に出す。`every` を守るための「最後に見回った時刻」はノードの小さな状態（`Store.PutState`）に置く。
- Docflow 自身が FTP / SFTP の **サーバーになる**形（複合機の「スキャン to FTP」、押すしかできない古い送り手）は別の口（生の TCP の Listen、FTP の受動のポートの範囲、SSH のホスト鍵の持ち物）が要るので、この要件では問いに留める（open_questions.server-mode）。

## 4. 実装の部品（Claude の推奨。確認待ち）

- **SFTP**：SSH の搬送・認証・チャネルは `golang.org/x/crypto/ssh`（既にある）。SFTP そのもの（version 3、draft-ietf-secsh-filexfer-02：OPEN・WRITE・READ・CLOSE・RENAME・MKDIR・OPENDIR・READDIR・STAT・REMOVE と `posix-rename@openssh.com`）は小さなパケットのプロトコルなので**自前**（数百行）。`github.com/pkg/sftp` は代案（BSD-2。サーバーも持つ）。
- **FTP**：RFC 959 + 2428（EPSV）+ 3659（MLSD・SIZE・MDTM）+ 4217（AUTH TLS、PROT P）を `net` と `crypto/tls` で**自前**。データ接続の TLS は制御接続のセッションを再利用する（多くの FTPS サーバーが要求する。`tls.Config.ClientSessionCache` を共有）。`github.com/jlaffaye/ftp` は代案（MIT）。
- 他社の SDK を入れない方針（S3 は自前の PUT/GET、SMTP は net/smtp）に揃える。新しい依存は 0。
- 試験は同じプロセスの中の偽のサーバー：SFTP は `x/crypto/ssh` のサーバー + 自前の最小の SFTP サーバー（試験用）、FTP は最小の FTP サーバー（試験用、TLS あり）。実物の確認は macOS の `sftp-server`（OpenSSH）と、Docker の vsftpd / pure-ftpd で手で。

## 5. 他の転送プロトコル（調査の結論。`survey:transfer-protocols`）

利用者の問い「他に対応すべき転送プロトコルはあるか」への Claude の答え（2026-10-27、Web 検索 3 回と Claude の知識）。

| 候補 | 判定 | 理由 |
| --- | --- | --- |
| FTPS | **FTP に含める** | 明示的 TLS は FTP の延長。実務の FTP はほぼこれ |
| SCP | 入れない（利用者） | OpenSSH 9.0 以降の scp は中身が SFTP。SFTP で足りる |
| **WebDAV** | **次の候補** | Nextcloud・ownCloud・SharePoint・Box・NAS・複合機の宛先。HTTP の上（PUT・MKCOL・MOVE・PROPFIND）なので出口は http の種類 + 2 命令、入口は PROPFIND の見回り。新しい依存なし |
| SMB / CIFS | 種類は作らず、**マウント + folder** | 日本の事務所のファイルサーバーと複合機の「スキャン to SMB」の既定。純 Go の SMB2 クライアント（hirochachacha/go-smb2、NTLMv2 だけ）はあるが、ホストでマウントして folder の種類で書けば新しいコードが要らない。コンテナやサーバーレスでマウントできない場面が出たら go-smb2 を考える |
| AS2（EDIINT） | **保留** | 小売・製薬・自動車の B2B EDI。流通 BMS も規定。S/MIME の署名と暗号、MDN の受領、証明書と取引先の台帳が要る大きな部品。具体的な取引先が現れたら |
| AS4 / ebXML MS 2.0 / JX 手順 | 保留 | ebMS と JX は流通 BMS の手順（JX は小売のサーバーへ取りに行く SOAP-RPC）、AS4 は EU の e-delivery（Peppol）。業界のメッセージの扱いと一緒でないと意味がない |
| 全銀協標準通信プロトコル | 入れない | ベーシック手順（ISDN）と旧 TCP/IP 手順は 2023 年末でサポート終了、広域 IP 網版（2017 年制定）が残る。仕様の入手と製品が前提の金融向けで、HULFT と同じ扱い |
| HULFT・JP1・EDI-Master 等 | 入れない（利用者） | プロプライエタリ |
| rsync | 入れない | プロトコルが実装の仕様。ssh + rsync は道具であってノードにしない |
| Azure Blob | Azure の配備が出たら | S3 と API が違う（SharedKey / SAS）。GCS は s3 の種類の endpoint で届くはず（HMAC 鍵。`If-None-Match: *` の PUT が通るかは未確認） |
| IMAP（メールの受信） | 別の要件 | 転送プロトコルではないが、複合機の「スキャン to メール」と取引先のメール添付は入口として多い。document-flow の inputs.smtp_receive のまま後 |
| NFS | マウント + folder | SMB と同じ |
| MQ（AMQP・Kafka・MQTT） | 範囲外 | ファイル転送ではない。トリガーの入口として別に |

結論：**FTP / SFTP の次は WebDAV**。SMB と NFS はマウントして folder。EDI 級（AS2・AS4・ebMS・JX・全銀）は具体的な相手が決まるまで保留。

## 6. 問い（利用者へ）

1. **入口の形（本質）**：A. 相手のサーバーを tick で見回る（推奨。サーバーレスでも動く、EDI・バッチの受け取り）／B. Docflow 自身が FTP / SFTP のサーバーになる（スキャン to FTP、押すしかできない送り手。自分のポートと受動のポートの範囲、ホスト鍵の持ち物）／C. 両方。暫定は A、B は open_questions.server-mode に引き金つきで残す。
2. 種類の数（任意）：A. `ftp`（ftp / ftps）と `sftp` の 2 つ（推奨。設定とホスト鍵が違う）／B. URL の scheme で分ける 1 つ。
3. 部品（任意）：A. 自前のクライアント（stdlib + x/crypto/ssh、依存 0。推奨）／B. jlaffaye/ftp + pkg/sftp（コードは減るが依存 2 つ）。
4. FTP の TLS の既定（任意）：A. AUTH TLS を要求し、平文はノードごとの opt-in（推奨）／B. あれば使う（日和見）。
5. 取ったファイルの後始末の既定（任意）：A. 消す（MFT の慣例。Docflow 側に残すのは archive の仕事。推奨）／B. done へ移す。

## 7. 段

- **fl-4a**（`milestone:fl-4a`）：出口 `ftp`・`sftp`、入口 `ftp-poll`・`sftp-poll`、`Kind.Tick`、`Discover` のホスト鍵取得、エディタの欄、マニュアルの章。fl-4 の残り（人の操作・人の入力・split・dedupe・流量）はそのまま fl-4。
- 後：WebDAV（出口と見回り）、サーバーの形、known_hosts、fw-2 の接続先カタログ（connections）に FTP / SFTP の接続先を載せる。

## 8. 実装（2026-10-27、同日）

利用者「暫定の選択（入口は Tick の見回り、ftp と sftp の 2 種類、自前のクライアント、AUTH TLS を要求、取ったら消す）のまま milestone:fl-4a を実装してください」。§6 の 5 つの問いはこれで決まり（`requirement:docflow-file-transfer` の open_questions は answered）、サーバーの形と接続先カタログは deferred のまま。

- **芯**：`Kind.Tick`（`Manager.Tick` がノードごとに呼ぶ。`TickReport` に Ticked・Polled・Notes）、`KindState` / `PutKindState`（ノードの名で引く kind-states、64 KiB まで。`DeleteNode` が消す）。fw-4 の kind-collections-tickers の tick の方だけが前に出た。
- **docflow/io/transfer**：`FS`（List・Stat・Open・Create・Rename(replace)・Remove・Mkdir）、`Site`（`ParseURL`：scheme ごとの既定のポート、URL のパスワードは断る、`/~/dir` はログインのディレクトリの下、`..` は落とす）、`ExitFields` / `PollFields`、`Send`（`.<delivery>.part` → rename、あるものは書いた扱い、`MkdirAll`）、`Poll`（every → 一覧 → 取り置きの rename → 読む → `Submit` → delete / move。settle は「更新が新しければ残す、前回と同じなら取る」。refused は `.refused` に改名）。
- **docflow/io/sftp**：`client.go` が SFTP v3（READ/WRITE は 16 × 32 KiB の窓で pipeline、posix-rename）、`ssh.go` が x/crypto/ssh の `Dial`（パスワード + keyboard-interactive か PEM の鍵、`ParseHostKey` は OpenSSH の 1 行か SHA256:、`HostKey()` が署名なしで鍵だけ取る）、`sftp.go` が `sftp` と `sftp-poll`（`Discover` がホスト鍵を返す）。
- **docflow/io/ftp**：`client.go` が FTP(S)（AUTH TLS / 暗黙 / 平文、EPSV→PASV、MLSD か NLST+SIZE+MDTM、TLS セッションの再利用、`tls_fingerprint`）、`ftp.go` が `ftp` と `ftp-poll`。
- **エディタ・マニュアル**：Fields の種類で discovers のものに「サーバーに聞く」（答えの同名の欄を埋める）。4 つの章は `help/{en,ja}/nodes.md`（`Kind.Doc` ではない：web/help の TestLinks が静的な見出しを見るのと、他のアダプタと同じ置き場だから）、`entries.md` に見回りの節。
- **試験**：偽の SFTP サーバー（x/crypto/ssh のサーバー + 自前の最小の SFTP サーバー）と偽の FTP サーバー（TLS あり、MLSD の有無）、OpenSSH の `/usr/libexec/sftp-server` を stdio で、Docker の OpenSSH 9.9 sshd と Alpine の vsftpd 3.0.5（require_ssl_reuse=YES、MLSD なし）で実物の出口と見回り（`testdata/{sshd,vsftpd}` に Dockerfile と手順。環境変数で有効になる `TestSFTPReal` / `TestFTPReal`）。
- **作りながら決めたこと**：`$source` は種類の名（`sftp-poll`）、メタデータは `sftp.path`・`sftp.size`・`sftp.time`、名乗りのユーザーはログインの名。`max_per_tick` は試した数。`done_dir` の `..` は落とす。FTP の 550 は「無い」扱い。詳しくは `milestone:fl-4a` の facts.implementation。
