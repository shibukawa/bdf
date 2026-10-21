# BDF Print Manager 設計案（保管・検索・ビューア・プリンタの束ね）

2026-10-20 の設計案。2026-10-20 利用者「BDF Print Manager 的な CUPS サーバーとしても振る舞う Web サービスがあっても良いと思った。電子帳簿保存法対応で、印刷されたものをオブジェクトストレージに BDF で保存。vault lock で上書き・削除できないモードでファイルを書き出す。日付やファイル名、その他メタデータや本文の検索もできる、Google Drive 的なビューアを持つ。また、CUPS サーバーの先には、その保存用の仮想プリンタ以外に、実プリンタを何台もぶら下げられる。スプール機能を持ち、クライアントから任意のプリンタに直接送ることもできる。グループになっているプリンタにラウンドロビンで送るのもできる。キャンセルなどもできる。仮想プリンタを複数作ることもできて、特定のユーザーに見せる（ペーパーレス）のも可能」。

`.knowledge/` の `requirement:print-manager`、`system:worm-object-storage`、`system:denchoho`。土台は `requirement:print`（print-m1〜m5 は main に入った）と `decision:thin-store`（yuruidb）。

## 0. 既にあるものと、足すもの

利用者の挙げた機能を、今の `bdf printd` と並べるとこうなる。

| 利用者の挙げた機能 | 今（print-m2〜m5） | 足すもの |
|---|---|---|
| スプール、キャンセル | ある（`print/spool`、`bdf jobs cancel`、IPP の Cancel-Job） | — |
| 実プリンタを何台もぶら下げる | ある（`printers:` に ipp・socket・lpd・serial・usb・winspool） | — |
| グループのプリンタに送る | ある（キューに複数のプリンタ、`policy: any` / `failover`） | `policy: round-robin`（§4） |
| 任意のプリンタに直接送る | キューを 1 台ずつ書けばできる | プリンタごとの暗黙のキュー（CUPS と同じ。§4） |
| IPP プリンタとして見せる | ある（m4。`/ipp/print/<キュー>`、DNS-SD、AirPrint） | — |
| CUPS のサーバーとして振る舞う | 無い（クライアント側の CUPS-Get-Printers だけ） | CUPS の名前の空間と操作（§5） |
| 保存用の仮想プリンタ | 無い | **保管庫（vault）と仮想プリンタ**（§1、§2） |
| 上書き・削除できない保存 | 無い | オブジェクトストレージの WORM（§2） |
| 日付・ファイル名・メタデータ・本文の検索 | 部品はある（`Document.SearchText`、`examples/search`） | 検索項目の取り出しと索引（§3） |
| Google Drive 的なビューア | 状態のページだけ（`print/server/web`。読むだけ） | 保管した文書の画面（§6） |
| 特定のユーザーに見せる（ペーパーレス） | 無い（トークンはキューへの投入を絞るだけ） | 倉庫ごとの表示先と、登録した端末だけに見せるデバイス認証（§7） |

つまり、**新しく作る芯は「現場の表示とデバイス認証」と「保管庫」と「検索」と「画面」**で、プリンタを束ねる側は小さな足し算で済む。

## 1. 仮想プリンタ（保存用の出力先）

プリンタの設定に、紙に出さない出力先を足す。

```yaml
vaults:
  books:                          # 帳簿書類の保管庫
    store: s3://bdf-archive/books/?region=ap-northeast-1
    lock: compliance              # compliance | governance | none
    retain: 10y                   # 保存の期間（法人の欠損金の年度は 10 年）
printers:
  archive-books: {url: "archive://books"}          # 保存だけの仮想プリンタ
  mfp-3f: {url: "ipps://mfp-3f.local/ipp/print"}
queues:
  invoices:                       # 紙に出して、控えを保存する
    printers: [mfp-3f]
    archive: books                # 刷ったジョブを保管庫にも入れる（tee）
  paperless-invoices:             # 紙に出さない
    printers: [archive-books]
```

使い方は 2 つ。

- **保存だけ**（ペーパーレス）: 出力先が `archive://` だけのキュー。IPP プリンタとして見せれば、Windows・Mac・iPhone の「印刷」がそのまま保存になる。
- **刷って控えを保存する**: キューの `archive:`。実プリンタに送ったジョブを、同じ .bdf で保管庫にも入れる。電子帳簿保存法の「自己がコンピュータを使用して作成して取引相手に交付する書類の写し（控え）」はこちら（`system:denchoho`）。

保存するのは**スプールが刷った .bdf**（受けた PDF・ラスタ・Office をスプールが変換したもの）。受けたままの入力（PDF など）も並べて残すかはキューの設定で選ぶ（既定は残す。真実性の説明が「受けたもの」と「刷ったもの」の両方でできる）。

## 2. 保管庫（vault）— 上書き・削除できない書き出し

### 2.1 何を書くか

1 つのジョブにつき、次の 3 つを**書いた後に変えない**オブジェクトとして書く。

```
<prefix>/docs/2026/10/20/<uuid>.bdf         刷った .bdf
<prefix>/docs/2026/10/20/<uuid>.src         受けたままの入力（任意）
<prefix>/docs/2026/10/20/<uuid>.json        記録（下）
<prefix>/manifests/2026-10-20.json          その日の一覧と前の日の一覧の SHA-256（鎖）
```

記録（`<uuid>.json`）: uuid、.bdf と入力の SHA-256 と大きさ、ページ数、受けた時刻（サーバーの時計）、キュー・出力先・利用者・ジョブ名・部数・用紙、出どころ（IPP・API・帳票）、**検索項目**（§3。取引年月日、金額、取引先、書類の種類、書類番号）、取り出した方法とその版、`supersedes`（訂正の前の uuid）、保存の期限。

**1 日の一覧（manifest）**は、その日に入れたすべての uuid とハッシュと、前の日の一覧のハッシュを持つ。WORM は「消せない」を、鎖は「黙って抜いていない」を示す。任意で RFC 3161 のタイムスタンプを一覧に付ける（§2.4）。

### 2.2 訂正と削除

- **上書きしない**。訂正は新しいジョブとして入れ、記録の `supersedes` で前のものを指す（前のものも残る）。画面は最新を出し、履歴をたどれる。
- **削除**は期限が過ぎたものだけ、管理者の操作で（WORM の期限の前はストレージが断る）。訴訟などで止めるときは legal hold。
- **変えられる層は別に持つ**: タグ、フォルダ、配達先、人が直した検索項目。これらはスプールの DB（yuruidb）に、誰がいつ何を何から何に変えたかの履歴つきで置き、保管庫のオブジェクトは変えない。検索項目を人が直したときも、元の値と直した履歴が残る。

### 2.3 WORM の写し方（`system:worm-object-storage`）

| 保存先 | 仕組み | 写し方 |
|---|---|---|
| Amazon S3 | Object Lock（バケットの作成時に有効化、版の管理が要る） | PUT に `x-amz-object-lock-mode`（COMPLIANCE / GOVERNANCE）と `x-amz-object-lock-retain-until-date`、内容のチェックサム。legal hold は `x-amz-object-lock-legal-hold` |
| MinIO | S3 互換の Object Lock | 同じ。ただしコミュニティ版は 2026-02 に凍結（`system:conditional-storage`）。試験は自前の S3 の偽物（Object Lock の見出しを守る）を主にし、MinIO の最後の版と実際の S3 で確かめる |
| Google Cloud Storage | Object Retention Lock（オブジェクトごとの保持期限）と Bucket Lock | オブジェクトの `retention`（mode Locked、retainUntilTime） |
| Azure Blob | 不変ストレージ（版ごとの時間ベースの保持） | 版ごとの immutability policy |
| Cloudflare R2 | バケットのロックの規則（接頭辞と期間） | オブジェクトごとの指定は無く、バケットの規則に頼る |
| ローカルのディレクトリ | 無し | `lock: none` のときだけ。**法令の用途には使えない**と起動時と画面に出す |

「vault lock」という名前の AWS の機能は S3 Glacier の Vault Lock と AWS Backup の Vault Lock で、どちらもオブジェクトを 1 つずつ書く用途ではない。ここでの相当は **S3 Object Lock の compliance モード**。

保存の層は yuruidb（`decision:thin-store`）の Blob を使う。yuruidb の Blob の書き込みに「保持の期限とモード、legal hold」を渡す口と、起動時に「そのバケットで本当にロックが効くか」を確かめる口が要る（**yuruidb への追加。利用者のリポジトリ**）。確かめは、使い捨てのキーに短い期限で書いて、消せないことを見る（compliance のときは消せないものが残るので、governance で確かめるか、確かめ用の接頭辞を使う。§9 の問い）。

`compliance` は**誰も**（クラウドの管理者も）期限まで消せない。設定を誤ると費用ごと戻せないので、既定は `governance`、`compliance` は明示のときだけ、起動時に期限と対象を表示する。

### 2.4 真実性のための追加（任意）

- RFC 3161 のタイムスタンプ（認定タイムスタンプ局）を 1 日の一覧に付ける。電子取引の受け手側の要件の 1 つ（`system:denchoho`）。WORM が「訂正削除のできないシステム」の要件を満たすなら必須ではない。
- .bdf への電子署名は今の形式に無いので、ここでは使わない。

## 3. 検索 — 項目の取り出しと索引

### 3.1 検索項目をどこから得るか

電子帳簿保存法の書類の検索は「取引年月日その他の日付」が中心で、電子取引は「日付・金額・取引先」（`system:denchoho`）。印刷のストリームにはこれらの項目が無いので、**どこから取るか**が設計の要。優先の順：

1. **投げる側が渡す**: API・`bdf print -meta date=2026-10-20 -meta amount=12000 -meta partner=…`。IPP では標準の属性に入らないので、IPP の利用者は 2〜4 に頼る。
2. **帳票テンプレート**: bdf の帳票（Word・Excel・PowerPoint・pdfme・テキスト帳票）は差し込んだデータを知っている。テンプレートに「検索項目への写し」（`archive: {date: invoiceDate, amount: total, partner: customer.name}`）を書き、できた .bdf に項目を入れる。**これが一番確か**。
3. **キューの規則**: 本文（`Document.SearchText`）に対する正規表現や位置の規則（「請求日」の右の日付、「合計」の右の金額）。ベンダーの請求書の書式ごとに。
4. **人が入れる**: 項目の取れなかったものは画面の「未分類」に並び、人が入れる（履歴つき。§2.2）。

ファイル名（ジョブ名）、利用者、キュー、受けた日時、ページ数、本文は常に取れる。

### 3.2 索引

索引は**保管庫から作り直せる派生物**にする（正は保管庫のオブジェクト）。`bdf printd archive reindex` で保管庫を全部読んで作り直せる。

- 項目（日付・金額の範囲、取引先、種類、利用者、キュー、組み合わせ）と本文の全文検索。
- 暫定の実装：自前の索引を別の bbolt のファイル（`index.db`）に。本文は 2-gram（日本語に分かち書きは要らない）。外の検索エンジン（Meilisearch。`examples/search` に例がある）にも出せる口を作る。
- 規模の見積もり（暫定）：1 日 1000 ジョブ × 10 年で約 365 万件。本文 2 KB / 件なら 7 GB 程度。2-gram の転置索引を bbolt に置くのは、試して決める（§9 の問い `search-engine`）。

## 4. プリンタの束ね方の足し算

- **`policy: round-robin`**: 空いているかによらず順に回す（今の `any` は「最初に空いている 1 台」、`failover` は「順に、生きている最初の 1 台」）。止まっているプリンタは飛ばす。
- **プリンタごとの暗黙のキュー**: すべてのプリンタに同じ名前のキューがあるものとして、`/ipp/print/<プリンタ>` と API で直接送れる（CUPS ではプリンタ＝キューで、クラスが組）。設定の `queues:` に同じ名前があればそちらが勝つ。見せるかは `ipp.queues` で選べる。
- **ジョブの移動**: 別のキューへ（API、`bdf jobs move`、CUPS-Move-Job）。スプールの要件（`requirement:print` の spool.queue）に既にあり、未実装。

## 5. CUPS のサーバーとしての振る舞い

IPP Everywhere のプリンタとして見せる（m4）ことで、OS は 1 台ずつ足せる。CUPS のサーバーとして振る舞うと、**CUPS のクライアントがサーバーごと**見られる：

- 名前の空間 `/printers/<キュー>`、`/classes/<キュー>`（複数のプリンタのキュー）、`/jobs/<id>`
- 操作 `CUPS-Get-Printers`、`CUPS-Get-Classes`、`CUPS-Get-Default`、`CUPS-Move-Job`、`CUPS-Accept-Jobs`・`CUPS-Reject-Jobs`、`CUPS-Get-Document`（保管した文書を CUPS の流儀で取り出す）
- 使い道：Linux / macOS の `client.conf` の `ServerName`（`lpstat`、`lp`、アプリの印刷の一覧がサーバーのキューになる）、`lpstat -h`、cups-browsed

`lpadmin` などの管理の操作（CUPS-Add-Modify-Printer …）は持たない（設定は YAML と API）。ポートは m4 の 8631 のまま（631 で動かすかは配置の判断）。偽物の `cups` 搬送路（`fake-printers.ja.md` §3）で先に往復を固める。

## 6. ビューア（Google Drive 的な画面）

`bdf printd` の Web（今の状態のページ）を広げる。

- **左**: 保管庫、仮想プリンタ（キュー）、「自分の文書」「自分あて」、日付の木（年・月・日）、未分類
- **一覧**: サムネイル（`thumbnail`）か表。列はジョブ名、日付、取引先、金額、種類、利用者、ページ数
- **検索の欄**: 本文と、日付・金額の範囲、取引先、組み合わせ。条件を URL に持つ
- **開く**: bdf のビューア（サイトの viewer の部品）でそのまま。ページ・検索の当たり・ダウンロード（.bdf、PDF に書き出し）
- **詳細**: ハッシュ、保存の期限とモード、legal hold、訂正の履歴（supersedes の鎖）、項目の変更の履歴
- **操作**: 実プリンタに刷り直す（保管した .bdf を参照する新しいジョブ）、タグ、フォルダ、項目の修正、配達先の追加

保管した文書の読み出しは、権限の確認と監査の記録（誰がいつ何を見たか）を通す。

## 7. 現場の表示 — 倉庫ごとの仮想プリンタとタブレット（ペーパーレスの本来の用途）

2026-10-20 利用者「基本的にはペーパーレスで、伝票を見せる、という用途を考えていた。現場の倉庫ごとの仮想プリンタを作り、タブレットからアクセス。デバイス認証とかで認証してその端末からのみ見せる、といったことを想定していた」。

つまり「見せる」相手は**人ではなく場所（倉庫）と、そこに置いた端末**。基幹システムや WMS は今までどおり「倉庫 A のプリンタ」に刷り、紙の代わりに倉庫 A のタブレットに出る。

### 7.1 表示先（display）の仮想プリンタ

```yaml
locations:
  warehouse-a:
    name: 第 1 倉庫
    fallback: {printer: wh-a-laser, after: 10m}   # 10 分だれも開かなければ紙に出す（任意）
    keep: 14d                                      # 完了した伝票を端末に出しておく日数
    archive: books                                 # 控えを保管庫にも入れる（任意。§2）
printers:
  wh-a-display: {url: "display://warehouse-a"}
  wh-a-laser:   {url: "ipps://wh-a-laser.local/ipp/print"}
queues:
  warehouse-a:  {printers: [wh-a-display]}         # /ipp/print/warehouse-a として見せる
```

- 出力先 `display://<場所>` は、ジョブを .bdf にしてスプールに置き、その場所の端末に知らせる（今の `/events` の SSE を場所ごとに絞る）。ジョブは置いた時点で `completed`。
- **伝票の状態**は、ジョブとは別に持つ：新着 → 開いた（どの端末が・いつ）→ 完了（端末で「完了」を押す）。送った側が完了の前に Cancel-Job したら端末から消える（取り消しの印を残す）。
- **紙への逃げ道**（任意）：決めた時間だれも開かない、またはその場所の端末がどれもつながっていないとき、その倉庫の実プリンタに刷る。ペーパーレスにしても出荷が止まらないように。
- 控え（電子帳簿保存法の書類の写し）にもしたいなら `archive:` で保管庫にも入れる（§2）。表示と保管は別の設定。

### 7.2 デバイス認証 — 登録した端末だけに見せる

共有のタブレットには個人のアカウントが無いので、**端末そのもの**を認証する。

**既定：登録コード＋端末に閉じた鍵（MDM 不要）**

1. 管理者が画面で「倉庫 A に端末を足す」を押すと、短い期限・1 回きりの登録コード（QR）が出る。
2. タブレットのブラウザで QR を読むと、ページが WebCrypto で**取り出せない（non-extractable）** ECDSA P-256 の鍵を作り、IndexedDB に置き、公開鍵とコードを送る。
3. サーバーは端末 ID を場所に結びつける（管理者の承認を挟むかは設定）。
4. 以後、開くたびにサーバーの challenge に鍵で署名し、短い期限のセッション（HttpOnly、SameSite=Strict の Cookie）を得る。鍵はページの外に出せないので、URL や Cookie を他の端末に写しても通らない。
5. 端末の一覧（名前、場所、最後の接続、IP）から管理者がいつでも外せる。

弱点：ブラウザのデータを消すと登録が消える（登録し直す）。iPadOS の Safari はしばらく開かないサイトのストレージを消すことがあるので、**ホーム画面に追加した Web アプリ**で使う（未確認、実機で確かめる）。同じ端末の上で悪意のあるスクリプトが動けば鍵を使われる（取り出せはしない）。

**選べる強化**

- **MDM のある端末群**：MDM でクライアント証明書を入れ、mTLS で端末を見分ける（iPadOS は Managed Device Attestation で鍵を Secure Enclave に閉じられる、Android は Android Enterprise と鍵の証明。未確認）。ブラウザの鍵より強く、消えない。
- **網の制限**：場所ごとに、つながってよい IP の範囲（倉庫の LAN、VPN）。
- passkey（WebAuthn）は使わない：同期されて端末に閉じず、共有の端末に個人のアカウントが無い。

**作業者**：共有の端末で「だれが完了したか」を残したいときは、端末の上で名前を選ぶか短い PIN を入れる（記録のためで、認証ではない）。

### 7.3 端末に残さない

- 既定はオフラインで持たない（開くたびに取りに行く）。ダウンロードのボタンも出さない。
- さらに守るなら、ebook で作った「要求ごとに封をする」仕組み（`examples/secure-reader`。エフェメラルの ECDH でページの束を封じる）を表示にも使える。
- 任意で、ページに端末名と時刻の透かし。

### 7.4 タブレットの画面

- 新着の伝票を大きく、新しい順に。着いたら音と表示（SSE）。今日・未完了・完了の切り替え。
- 開くと bdf のビューアで全画面。「完了」「この倉庫のプリンタで刷る」。
- つながっていないときは目立つ表示（紙への逃げ道が働くかも出す）。

### 7.5 事務所の人と管理者

事務所の画面（§6）と管理（端末の登録・場所・プリンタ）は人のログイン。まずは今のトークン、後で OIDC（Entra ID、Google Workspace）。IPP で投げる側（基幹、WMS）は今のトークン（Basic）。`requesting-user-name` は自己申告なので権限の根拠にしない。

## 8. 電子帳簿保存法との対応（`system:denchoho`。法令の判断ではない）

国税庁の資料（令和元年 7 月のパンフレット）では、自己がコンピュータで作成して交付する書類の写し（見積書・請求書・納品書・領収書の控え）を電子データのまま保存できる。要件のうち書類にかかるもの：

| 要件 | この設計で |
|---|---|
| システム関係書類（概要書・仕様書・操作説明書・事務処理マニュアル）の備付け | 概要書と操作説明書のひな形を docs に用意する。事務処理マニュアルは利用者の会社が作る |
| 見読可能装置の備付け、画面・書面に整然と明瞭に速やかに出せる | ビューアと刷り直し |
| 取引年月日その他の日付で検索できる | 検索項目の日付（§3） |
| 日付の範囲で検索できる | 同じ |

令和 4 年（2022）の改正で事前承認は無くなり、検索の要件にも緩和がある（税務職員のダウンロードの求めに応じる場合など）。**今の要件は国税庁の「電子帳簿保存法一問一答」の最新の版で確かめる**（未確認）。電子取引（取引先へ電子で渡す・受け取る）は別の条で、真実性（タイムスタンプ、訂正削除のできない・記録が残るシステム、事務処理規程のどれか）と、日付・金額・取引先の検索が要る。仮想プリンタで作った書類を取引先へ電子で渡す使い方はこちらに入る。

JIIMA の認証は目標にしない（第三者の審査が要る）。

## 9. 決めること（問い）

答えたもの：

- **順番**（2026-10-20 利用者「偽プリンタをまずきちんと作って、テストを全てそれをベースに行うのを先にしたい」）：偽プリンタ（print-m6、`fake-printers.ja.md`）が先。Print Manager はその後。
- **形**（2026-10-20 利用者「bdf printd とは別にサービスを作りたい」）：bdf printd を広げず、別のサービスにする。この設計書で「bdf printd の設定」として書いた YAML の例は、そのサービスの設定と読み替える。
- **ペーパーレスの意味**（2026-10-20 利用者）：倉庫ごとの仮想プリンタの伝票を、登録した端末（タブレット）だけに見せる（§7）。
- **画面のログイン**：タブレットはデバイス認証（§7.2）。事務所の人と管理者はまずトークン、OIDC は後（暫定）。
- **土台**（2026-10-20 利用者）：Web は popcornweb、電子帳簿保存の保管庫は S3、普段のスプールは yuruidb の使える保存先、サーバーレスでも動く（§11）。

pm-1 を始めるときに決めるもの：

- **bdf printd との関係**（`printd-relation`）。A：print/ の部品（スプール、IPP のサーバー、変換）をライブラリとして組み込み、実プリンタにも自分で送る。B：表示・保管・検索・画面の中央のサービスで、IPP の受け口と INFRA の Infrastructure Printer を持ち、倉庫の実プリンタへは各拠点の bdf printd が INFRA（m4 の upstreams）で取りに来る。C：bdf printd が受けて Manager の API に渡す。推奨は B（倉庫は NAT の向こうで m4 の拠点の代理がそのまま使え、タブレットは中央に繋ぐ）。サーバーレスでも VPC と Direct Connect・VPN なら A が動く（利用者の指摘）ので、プリンタごとに A（direct）と B（branch）を選ぶ（§11.1）。
- 置き場所（このリポジトリの `cmd/` か別のリポジトリか）。

暫定で決めて進めるもの（後で替えられる）：

- 索引は自前の bbolt の 2-gram、Meilisearch に出す口も（`search-engine`）
- 既定の WORM のモードは governance、compliance は明示（`default-lock-mode`）
- 試験は自前の S3 の偽物（Object Lock の見出し）と実際の S3。最初の本番の保存先は S3（`first-worm-store`。利用者の答え）
- 受けたままの入力も保存する（`keep-source`）
- 帳票テンプレートの「検索項目への写し」は form の側の足し算（`form-archive-fields`）
- デバイス認証の既定は登録コード＋WebCrypto の取り出せない鍵、MDM のある所は mTLS を選べる（`device-auth`）
- 端末に伝票を残さない（オフラインなし）。要求ごとの封は後で（`device-offline`）

## 10. 順番の案（Print Manager 側）

2026-10-20 に、利用者の主な用途（現場の表示）に合わせて並べ替えた。

1. **pm-1 現場の表示**: popcornweb の骨組み、`display://` と場所、伝票の状態（新着・開いた・完了・取り消し）、場所の版のポーリング（プロセスでは SSE も）、登録コードと端末の鍵の認証、端末の管理、タブレットの画面、INFRA で拠点の printd に渡す紙への逃げ道、tick。保存は yuruidb のスプール（WORM なし）。コンテナと Lambda の 2 つに出す（§11.6）
2. **pm-2 保管庫**: `archive://` とキュー・場所の `archive:`、記録と 1 日の一覧の鎖、S3 Object Lock（自前の偽物と実際の S3 で試験）と `lock: none`、検索項目の投げる側と規則、索引と検索の API、`bdf archive search / get / verify`、reindex
3. **pm-3 事務所のビューア**: §6 の Drive 的な画面、サムネイル、刷り直し、監査の記録、OIDC
4. **pm-4 束ね方と CUPS**: round-robin、プリンタごとのキュー、ジョブの移動、CUPS の名前の空間（偽物の S4 の後）。小さく独立なので前に入れてもよい
5. **後**: mTLS（MDM）の端末認証、要求ごとの封、帳票テンプレートの検索項目、RFC 3161、GCS・Azure・R2 のロック、期限の過ぎたものの削除、複数台

## 11. 土台 — popcornweb、保存の分け方、サーバーレス

2026-10-20 利用者「ウェブフレームワークには https://github.com/shibukawa/popcornweb これを利用して欲しい。電子帳簿保存用ストレージは S3 だが、普段のスプールのストレージは yuruidb が使えるものは利用できるように。print manager はサーバーレスでも使えると良い」。`system:popcornweb`（コミット 5af5b03）と `system:yuruidb`（コミット 004439a）を読んだ上での整理。

### 11.1 サーバーレスが決めること

サーバーレスの要件がいちばん形を決める。Lambda などの上の Manager は：

- **実プリンタへの道はプリンタごとに選ぶ**。最初の案では「サーバーレスは LAN に届かないので B だけ」としたが、利用者の指摘（「Lambda で動いていて、Direct Connect で拠点マシンなしで工場制御した事例があるので、不可能ではないと思う」）のとおり、VPC に置いた Lambda は Direct Connect・Site-to-Site VPN で拠点の LAN に届く。
  - **direct**：Manager が `print.Print`（ipp・ipps）と `print.PrintRaw`（socket・lpd）で直接送る。どちらも 1 回で終わる関数なので、呼び出しごとの送りに合う。送る前に記録に期限つきの借り（lease）を版で書いて 2 重の送りを防ぎ、失敗・busy は tick で再試行。呼び出しの時間の上限（Lambda は 15 分）を超える送りは branch に回す。VPC の Lambda から S3・DynamoDB へは VPC エンドポイントが要る。
  - **branch**：網が届かない拠点は、拠点の bdf printd が Manager の INFRA（Get-Jobs の fetchable、Fetch-Job、Fetch-Document、Acknowledge-Job、Update-Job-Status）を取りに来る。m4 の拠点の代理そのもの。
  - serial・usb・winspool は拠点の printd だけ。
- **DNS-SD・AirPrint の広告、8631 の振り分け、LPD・9100 の受け口を持てない**（マルチキャストは Direct Connect を越えない）。これらは拠点の printd の仕事。Manager の IPP は HTTPS の POST（`/ipp/print/<場所>`）だけ。
- **常駐の goroutine と時計を持てない**。期限は要求の中で遅れて評価する。branch のプリンタへの紙への逃げ道は「その場所の printd が fetchable を問うとき、期限を過ぎて開かれていない伝票を fetchable に変える」。printd の問い合わせそのものが時計になる。direct のプリンタは tick が期限を見て送る。保持・1 日の一覧・outbox の再試行は、外のスケジューラ（EventBridge Scheduler など）が叩く認証つきの `POST /_pm/tick`。プロセスで動くときは中の時計が同じ関数を呼ぶ。popcornweb は非 HTTP の trigger を持たないので、この形が合う。
- **インスタンスに状態を持てない**。伝票・端末・場所・期限はすべて保存先に。試験は「要求ごとにハンドラーを作り直しても通る」ことを確かめる。
- **要求の本文に上限がある**（popcornweb の既定 10 MiB、ホストによってはもっと小さい。数字は実装のときに確かめる）。大きな文書は API で S3 への署名つき PUT に上げてから確定する。IPP で上限を超えるものは拠点の printd を通す。

### 11.2 popcornweb の上の形

- 画面（タブレット、事務所、管理）は `.pw.html`。API・IPP・INFRA は `http.Handler` として同じ mux に載せる。bdf の `print/ipp` はそのまま、`print/server` の INFRA の受け口は `print/spool` に結びついているので、記録の層を差し替えられるように切り出す。
- 端末のセッションは popcornweb の session の `cookie`（封をしたクッキー、ストア無し）。端末を外したときに効かせるため、セッションは短く（例 10 分）、切れたら端末の鍵で challenge に署名し直す（その時に端末の記録を読む）。
- 事務所の人のログインは popcornweb の OIDC がそのまま使える。pm-3 を待たずに入れる費用は小さい（`viewer-login`）。
- 設定は 2 つに分ける。動かし方（保存先、保管庫、鍵、tick の秘密）は popcornweb の TOML と環境変数。場所・端末・表示先・逃げ道の規則は管理の画面と API で変えるデータで、スプールの保存先に置く。§1・§7.1 の YAML の例はこの 2 つに読み替える。
- TinyGo と Cloudflare Workers は対象にしない（変換器とフォントの大きさ、ファイルシステム無し、Workers の storage が r2 だけで S3 の Object Lock に届かない。`workers-target`）。

### 11.3 保存の分け方

| 何 | 保存先 | 性質 |
|---|---|---|
| 保管庫（.bdf、入力、記録の JSON、1 日の一覧） | S3 の Object Lock | 書いた後に変えない。消せない |
| スプールの記録（伝票、端末、場所、tick の進み、outbox、監査） | yuruidb のレコード | 版で守る。消せる |
| スプールの文書（見せる間・刷る前の .bdf、受けたままの入力） | yuruidb の Blob | 消せる。保管が済むまで消さない |

配置ごとの組み合わせ（どれも yuruidb のドライバ）：

| 配置 | 記録 | 文書 |
|---|---|---|
| 1 台のプロセス（社内のサーバー、コンテナ 1 つ） | bbolt（既定）か SQLite | 同じファイル |
| オンプレのプロセス＋保管庫だけ S3 | bbolt か社内の PostgreSQL | 同じ（保管庫は回線越しの S3） |
| 複数台のコンテナ | PostgreSQL か S3 | PostgreSQL か S3 |
| AWS Lambda | DynamoDB | S3（保管庫とは別のバケット） |
| Cloud Run functions | Firestore | GCS |
| Cloudflare（後） | D1 | R2 |

- cfkv は条件付きの操作が無いので断る。サーバーレスで bbolt を選んだら起動時に断る（インスタンスをまたげない）。
- 保管庫とスプールが同じバケットの重なる接頭辞なら起動時に断る。Object Lock のバケットには yuruidb の `.yuruidb-probes/` の確かめを向けない（既定の保持があると確かめのキーが消せない）。
- yuruidb の一覧は ID の順だけで、場所ごと・状態ごとに絞れない。暫定で場所ごとの「開いている伝票」の記録（版で守る ID の集合、上限つき）を持つ。足りなければ yuruidb に接頭辞の一覧を足す（`list-by-location`）。

### 11.4 保管庫の書き込み（S3）

- yuruidb にも popcornweb（tinygodriver の S3）にも Object Lock の見出しを渡す口が無い（grep で確かめた）。yuruidb の Blob に足す（`vault-client`。利用者の答え）。保管庫は「無いときだけ作る」と版の照合が要り、Blob の約束に合う。大きな文書を受ける presign は popcornweb の storage を使う。
- 足すもの：保持の期限とモードと legal hold の見出し、保持つきの PUT のチェックサムの見出し（Content-MD5 か内容のチェックサム。未確認）、バケットのロックが効くかの確かめ、legal hold の付け外しと期限の延長、版の管理のあるバケットでの If-None-Match の確かめ。
- 保管の書き込みは記録（outbox）を先に置く。失敗しても tick と次の要求で再試行し、済むまでスプールの文書を消さない。
- 1 日の一覧は、その日の後の tick が保管庫の接頭辞を一覧にして作る。一覧のキーは無いときだけ作る（If-None-Match）ので、同時の tick でも 1 つ。一覧の後に届いたその日の項目は、次の日の一覧の `late` に入れる。

### 11.4a オンプレのプロセス＋保管庫だけ S3

2026-10-20 利用者「オンプレでサーバープロセスを動かしながら、保存用ストレージだけ S3 というユースケースもあり得るかと思います」。保管庫は Manager の置き場所と独立に選べるので、この組み合わせを正式な配置の 1 つにする。オンプレならプリンタへは LAN で direct、tick は中の時計で、サーバーレスの制約は効かない。代わりに効くのは：

- **資格情報**：オンプレには IAM ロールも IMDS も無い。yuruidb は環境変数・web identity・コンテナ・IMDSv2 を探し、ほかは `auth/aws` の `Provider` の注入。まずは保管庫の接頭辞への操作だけを許した IAM ユーザーの鍵を環境変数で。IAM Roles Anywhere（X.509 の証明書で一時の資格情報）は `credential_process`（aws_signing_helper）を呼ぶ Provider を Manager に持たせて後で（`aws-credentials-on-prem`）。
- **回線の断**：受け付け・表示・刷ることは止めない。保管は outbox に積み、つながったら古い順に書く。スプールの文書は保管が済むまで消さないので、ディスクの残りと outbox の滞留（件数、最も古いものの時刻）を状態のページと監視に出す。ディスクが決めた量を下回ったら、保管庫の要るキューだけ受け付けを断り、表示だけのキューは続ける（暫定。`vault-backlog-limit`）。
- **WORM になるまでの遅れ**：受けてから S3 に書けるまでは、文書はローカルの消せる保存先にある。受け付けのときのハッシュを記録に、S3 に書けた時刻を記録の JSON と 1 日の一覧に残し、遅れの大きいもの（例 24 時間）に印を付ける。
- **時計と保持**：受けた時刻はオンプレの時計（NTP が前提）。S3 の LastModified も並べて残し、差が大きければ警告。保持の期限は受けた時刻から数える（遅れて書くときも retain-until は受けた時刻＋期間）。
- **網**：Direct Connect・VPN があればその上、無ければインターネットの HTTPS。`HTTPS_PROXY` を通す。暗号化は S3 の既定（SSE-S3）で足りる見込み。SSE-KMS の鍵を指定する見出しが yuruidb にあるかは未確認。

### 11.5 タブレットへの通知

基本は「場所の版」（場所の記録の Revision）を数秒おきに条件つき GET で問うポーリング。どのホストでも動き、インスタンスをまたぐ通知の仕組みが要らない。プロセスで動くときは popcornweb の型つきストリーム（SSE）でも押せる。費用は端末 1 台 5 秒おきで 1 日 1.7 万回（Lambda の要求の課金では小さいはず。pm-1 で見積もる。`notify-transport`）。

### 11.6 pm-1 で確かめること

- 同じソースを、コンテナ（bbolt のスプール）と Lambda（Web Adapter、DynamoDB＋S3 のスプール）に出して、同じ E2E が通る（`first-serverless-target`）。
- 要求ごとにハンドラーを作り直しても通る（インスタンスに状態が無い）。
- 紙への逃げ道（direct）：期限を過ぎた伝票を、Manager が偽プリンタ（`print/fakeprinter`）に送る。Lambda では偽プリンタを VPC の中に置いて試す。branch（拠点の printd が INFRA で取る）は続けて。
- `cmd/` の中で popcornweb のプロジェクト（popcornweb.toml）が `pw generate`・`pw build --target=lambda` で動くか、bdf の go.mod がどれだけ増えるか、tinygodriver が bdf の v1.3.5 に揃って popcornweb が動くか。
- 冷えた起動の時間と、イメージの大きさ（変換器とフォント込み）を測る。

### 11.7 新しい問い

答えたもの（2026-10-20 利用者）：

- `manager-location`：**bdf の `cmd/`**（go.mod は 1 つ。推奨の別 go.mod は採らなかった）。
- `first-serverless-target`：**AWS Lambda**。
- `vault-client`：**yuruidb の Blob** に Object Lock を足す（§11.4）。pm-2 の前に要る。
- `printd-relation`：「サーバーレスでは A が動かない」は誤り（Direct Connect の事例）。プリンタごとの route（direct・branch）に。pm-1 は direct を先に、branch は続けて（この順は Claude の案）。

残り：

- `spool-store-default`：暫定で bbolt。サーバーレスでは既定に頼らず設定で選ばせる。
- `list-by-location`、`notify-transport`、`tick-trigger`：上の暫定のとおり。
- `workers-target`：後。

