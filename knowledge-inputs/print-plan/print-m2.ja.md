# 印刷 実装案（print-m2）

2026-10-19 時点の実装案（2026-10-20 に保存を yuruidb に改めた）。要件と決定は `print-requirements.ja.md` §4.5〜§4.10 と `.knowledge/` の `requirement:print`（facts の spool・server・budgets）、`milestone:print-m2` にある。前段の print-m1（`bdf print`、`bdf printers`、`print`・`print/ipp`・`print/raster`）は実装して main に入った（#104）。ここは「何をどの順で作るか」を書く。

## 0. 決まっていること（と、暫定で進めること）

決まっていること。

- スプールはジョブを受け付けて永続化し、順に流し、失敗を扱い、状態を見せる。ジョブの状態は IPP の `job-state` と `job-state-reasons` の語彙をそのまま使う（m4 で IPP サーバーとして見せるとき、写すだけで済む）。
- 再起動の後に同じジョブを二重に刷らない。送る前にプリンタの `job-id` を記録し、再開時は `Get-Jobs` で確かめてから再送する。
- プリンタ側で `aborted` になったジョブは自動で再送しない。
- 既定は `localhost` に束縛する。外に開けるときは TLS とトークンとキューごとの許可。
- ログに文書の中身・パスワード・トークンを書かない。
- 暗号化した `.bdf` はスプールに暗号化のまま置く。描いた結果（Raster）は既定で残さない。

2026-10-19 に利用者が決めたこと。

- **サーバーは `bdf printd`**（`cmd/bdf` のサブコマンド）。`cmd/bdfprintd` は作らない。
- **永続化は組み込みの KV**。`go.etcd.io/bbolt`（MIT、Go だけ、etcd が使う B+ 木）の **1 ファイル**（`spool.db`）に、ジョブの記録・文書・連番・実績・プリンタの控え・証明書の指紋を全部入れる（何を入れるかは §3.2）。メモリだけの KV ではない。書き込みのたびに `fsync` するので、電源が落ちても受け付けたジョブは残る。人が書く設定（YAML）だけが別のファイル。
- **サーバーは `converter` の入力を受ける**。PDF、Office、画像などを受けて、スプールの中で `.bdf` に変換してから印刷する。
- **プリンタの中に溜めるジョブの数（`max-inflight`）の既定は 2**。1 つを刷っている間に次を送っておき、業務機が続けて刷れるようにする。まだ送っていないジョブは並べ替えも別のプリンタへの移しも効く。送った後のジョブは `Cancel-Job` で取り消す。
- **サービスとしての登録は持たない**（Windows の `x/sys/windows/svc` も作らない）。`bdf printd` は前面のプロセスとして正しく動くこと（シグナルと Ctrl-C で止まる、ログは標準エラー）だけを受け持ち、常駐は systemd・launchd・NSSM・WinSW などの外の道具に任せる。docs に例を書く。
- **m2 の後（m3 と m4 のどちらを先にするか）は m2 の後で決める**。

決めたことから出てくること。

- 1 つの `spool.db` を開けるのは 1 プロセス（bbolt がファイルをロックする）。並行はプロセスの中で行う（プリンタごと、描画の枠ごと）。2 つ目の `bdf printd` が同じファイルで起動しようとすると、ロックを待たずにエラーで終わる。
- **保存は yuruidb（`github.com/shibukawa/yuruidb`）を通す**（2026-10-20 利用者「作れそうということなので実装の開始をお願いしたい」。利用者が作った、SDK なしで bbolt・S3・R2・GCS・DynamoDB・RDB を同じ CRUD に載せ、楽観ロックで競合を検出するライブラリ。評価は `store.ja.md` §9）。レコードは yuruidb の bbolt ドライバ、文書は同じ `spool.db` の bbolt の Blob、写しは `Store.Backup`。yuruidb に触るのは `print/spool/store.go` だけにする。bbolt は yuruidb の下で使い、bdf が直接触るのは止めた状態の詰め直し（`bdf printd -compact`）だけ。m2 は 1 プロセスなので、ジョブの「借り」は作らない。
- 1 ジョブの処理は短い「段」（変換・送信・見届け）の繰り返しで、次に動かす時刻を記録に持つ（§3.4）。再起動の後の照合がどの段からでも同じ形で書ける。

暫定で推奨どおりに進めること（異論があれば変える。どれも後から差し替えられる位置に置く）。

- **暗号化文書のパスワードはメモリにだけ持つ**。再起動の後は `pending-held`（理由 `job-password-wait`）になり、API でパスワードを渡し直すと流れる。
- **既定の束縛は `127.0.0.1:8640`**。loopback だけに束縛し、トークンを 1 つも書いていなければ認証なし。
- **`bdf print -spool dir`（スプールに直接置いて帰る）は作らない**。スプールの持ち主は 1 つのプロセス（`bdf printd`）だけにし、置くのは常に API を通す。要件からの変更。

## 1. 成果物

次のことができること。

```
bdf printd -config /etc/bdf/printd.yaml          # 前面で動く。キューとプリンタは設定ファイルで
bdf printd -config printd.yaml -check            # 設定の検査（プリンタに Get-Printer-Attributes まで）
bdf printd -config printd.yaml -compact          # 止めた状態で spool.db を詰め直す

bdf print -server http://127.0.0.1:8640 -queue office -sides long invoice.bdf     # スプールに置く（既定は置いて帰る）
bdf print -server http://127.0.0.1:8640 -queue office -wait 10m report.xlsx       # 変換と印刷の完了まで待つ
bdf jobs -server http://127.0.0.1:8640                       # 終わっていないジョブ
bdf jobs -server http://127.0.0.1:8640 -all -queue office    # 終わったものも
bdf jobs -server … cancel 42 | hold 42 | release 42 | retry 42 | show 42
bdf jobs -server … accounting -from 2026-10-01 -to 2026-10-31 -csv     # 実績の書き出し
bdf printers -server http://127.0.0.1:8640                   # サーバーのキューとプリンタと状態

curl -F document=@invoice.bdf -F 'ticket={"sides":"two-sided-long-edge","copies":2}' \
     -H "Authorization: Bearer $TOKEN" https://print.example.local:8640/api/v1/queues/office/jobs
```

`BDF_SERVER` と `BDF_TOKEN`（または `-token-file`）が `-server` と認証の既定になる。ブラウザで `http://127.0.0.1:8640/` を開くと、キュー・プリンタ・ジョブの状態のページ（読むだけ。更新はサーバーからの通知）が見える。

## 2. 構成

```
print/
  print.go         Options.OnCreated（Create-Job の直後、文書を送る前に呼ぶ。エラーならジョブを取り消して返す）、
                   Options.OnSent、Options.Printer（能力を聞き直さない）、
                   Check（プリンタのジョブを 1 回聞いて、終わっていれば m1 と同じ誤りに写す）                  (S1)
  ticket.go        Ticket に JSON のタグ（キーは IPP の属性の名前: copies、sides、media、media-source …）と、
                   キューの既定に重ねる Merge                                                                 (S1)
print/ipp/
  trust.go         TrustStore の保存先の差し替え（TrustBackend。nil なら今どおり Path の JSON）                  (S1)
print/ipp/ipptest/
  ipptest.go       プリンタの停止と再開、ジョブの履歴の消去（忘れたプリンタ）、文書の途中で接続を切る、
                   Create-Job の応答の前に接続を切る、which-jobs=all                                           (S1)
print/spool/
  spool.go         Spool: Open（spool.db のパス、時計）、Submit、Get、List、Cancel、Hold、Release、Retry、Move、
                   SetPriority、SetPassword、Subscribe（変更の購読）、Close                                    (S2)
  job.go           Job（記録の中身）と状態機械（IPP の job-state と理由）、遷移の検査                           (S2)
  store.go         yuruidb に触る唯一のファイル: Bind（jobs・done・meta・printers・trust）、ジョブの記録の
                   読み書き（版つき）、文書の Blob、連番、写し                                                   (S2)
  hold.go          job-hold-until の時間帯（indefinite、day-time、evening、night、weekend、second-shift、third-shift）
                   と job-hold-until-time                                                                      (S2)
  account.go       実績（終わったジョブごとに 1 件）と書き出し（JSON Lines、CSV）                               (S2)
  step.go          段: convert、send、watch。Step(ctx, id) が 1 段を行い、次の段と時刻を書く                      (S3)
  recover.go       再起動の後の照合（二重印刷の防止）                                                           (S3)
  sched.go         回し手（Runner）: キューの順、プリンタごとの送り手と見届け、max-inflight、組、描画の枠        (S3)
  monitor.go       プリンタの見張り: printer-state と理由、能力の更新（printer-config-change-time）              (S3)
  gc.go            保持期間・件数・バイトで消す                                                                 (S2)
print/server/
  config.go        設定ファイル（YAML）と環境変数、検査                                                         (S4)
  server.go        Server: New（設定 → Spool とプリンタとキュー）、Handler、Run、Shutdown                        (S4)
  api.go           HTTP API（/api/v1）                                                                         (S4)
  auth.go          トークンとキューごとの許可、利用者の名前                                                     (S4)
  events.go        変更の通知（Server-Sent Events）                                                             (S4)
  metrics.go       /metrics（Prometheus のテキスト。依存なしで書く）、/healthz                                   (S4)
  web/index.html   状態のページ（埋め込み。ビルドなしの HTML と小さな JS）                                       (S4)
print/server/client/
  client.go        API の Go のクライアント。cmd/bdf と examples が使う。spool と server に依らない                 (S5)
cmd/bdf/
  printd.go        bdf printd: -config、-check、-compact、-listen、-log。converter/all で変換を渡す。
                   SIGINT / SIGTERM（Windows は Ctrl-C と Ctrl-Break）で受け付けを止め、送信中のジョブを待って閉じる (S4, S5)
  print.go         -server、-queue、-token-file（送る先がサーバーのとき）                                        (S5)
  jobs.go          bdf jobs                                                                                    (S5)
  printers.go      -server（サーバーのキューとプリンタ）                                                         (S5)
examples/print-server/
  main.go          帳票を穴埋めしてスプールに置くサンプル。server をライブラリとして埋め込む                        (S6)
docs/
  architecture/cli（bdf printd、bdf jobs、-server）、常駐のさせ方、examples/print-server、design §3.36 に追記 (S6)
test/print/ci.sh   bdf printd を ippeveprinter につないで投入と完了、再起動をまたぐ                               (S6)
```

`spool` は `print`、`print/ipp`、yuruidb に依る（yuruidb を import するのは `store.go` だけ）。`server` は `spool` に依り、`client` はどれにも依らない（API の JSON の型だけを持つ）。`converter` に依るのは `cmd/bdf` と `examples` だけ。新しい依存は `github.com/shibukawa/yuruidb`（まだ版が無いので疑似版。その下の `go.etcd.io/bbolt` と `golang.org/x/sys`）。

## 3. 各部の設計

### 3.1 print パッケージに足すもの（S1）

m1 の `print.Print` は 1 回の呼び出しで聞く・作る・送る・待つ。スプールはこれを「送るまで」と「見届け」に分けて使う（`Options.Wait` が 0 なら送った時点で返る。m1 のまま）。

- `Options.OnCreated func(*ipp.JobStatus) error`: `Create-Job` の応答を受けた直後、`Send-Document` の前に同じ goroutine で呼ぶ。スプールはここでプリンタの `job-id` と `job-uri` を記録する（書けなければエラーを返し、`Print` はプリンタのジョブを取り消して返す）。`Progress` の `created` は今どおり残す。
- `Options.OnSent func(*ipp.JobStatus) error`: `Send-Document` が成功した直後。スプールは「文書を送り終えた」印を書く（引き継ぎの判断に使う。§3.6）。
- `Options.Printer *ipp.Printer`: 与えられれば `Get-Printer-Attributes` を省く（見張りが持っている能力を使う）。停止の判定は今どおり行う。
- `Check(ctx, c, last *ipp.JobStatus) (*ipp.JobStatus, error)`: プリンタのジョブを `Get-Job-Attributes` で 1 回聞き、前の答えに重ねる。終わっていれば m1 の `wait` と同じく `ErrCanceled`・`ErrAborted` に写す。m1 の `wait` はこれの繰り返しに書き直す（動きは同じ）。見届けの段はこれを時刻ごとに 1 回呼ぶ。
- `Ticket` に JSON のタグ。キーは IPP の属性の名前（`copies`、`multiple-document-handling` ではなく分かりやすい `collate`、`sides`、`media`、`media-source`、`media-type`、`print-color-mode`、`bilevel`、`print-quality`、`printer-resolution`、`pages`、`print-scaling`、`orientation-requested`、`finishings`、`output-bin`、`job-hold-until`、`job-priority`、`job-name`）。`Merge(defaults Ticket) Ticket` で、ジョブが言わない項目をキューの既定で埋める。
- `ipp.TrustStore` に保存先の差し替え（`Backend TrustBackend`。`Load() (map[string]string, error)` と `Save(map[string]string) error`）。`bdf printd` は指紋を `spool.db` に置く。設定に `fingerprint:` を書いたプリンタは最初から固定（初回の指紋を使わない）。

`ipptest` に足すもの: プリンタを止めて再開する（`Stop(reasons)`、`Resume()`）、ジョブの履歴を忘れる（`Forget()`。履歴を持たないプリンタの再現）、`which-jobs=all`、`Send-Document` の本文を n バイトで切る（`FailAfter(n)`）、`Create-Job` の応答を返す前に接続を切る（`DropCreateResponse()`。記録の前に落ちた場合の再現）。

### 3.2 何を保存するか、どこに（S2）

保存するのはジョブの番号だけではなく、ジョブを再起動の後も流し続けるのに要るもの全部。

| 中身 | 大きさの目安 | いつまで |
|---|---|---|
| ジョブの記録: id、uuid、キュー、利用者、名前、チケット、状態と理由、優先度と保留、段と次に動かす時刻、試行の回数、プリンタの名前とプリンタ側の `job-id` / `job-uri` と送り終えた印（二重印刷の防止の要）、時刻、枚数、警告 | 1 件 1〜2 KB | 終わってから保持期間（既定 7 日） |
| 文書: 受けたままのバイト列（`.bdf`、PDF、docx …）と、変換した `.bdf` | ジョブの大部分（KB〜数百 MB） | 終わってから保持期間（既定 7 日） |
| 連番（次の `job-id`）、形式の版 | 数バイト | ずっと |
| 実績: 終わったジョブごとに時刻・利用者・キュー・プリンタ・枚数・印象数・用紙・両面・色 | 1 件 300 B | 既定 400 日 |
| プリンタの控え: 最後に聞いた能力と状態（`printer-config-change-time` つき） | 1 台 数十〜数百 KB | 設定にある間 |
| ipps の証明書の指紋 | 1 台 80 B | ずっと |

保存しないもの: 描いた Raster（送るたびに描く）、パスワード（メモリだけ）、トークン（設定ファイルとその参照先）、設定（YAML のまま）。

置き場所は yuruidb の bbolt ドライバの 1 ファイル（`spool.db`）。Go の型ごとにコレクションを `Bind` する。

```
spool.db（yuruidb の bbolt ドライバ）
  コレクション jobs       ID は 10 桁の 0 埋めの job-id（"0000000042"）→ 終わっていないジョブの記録
  コレクション done       同じ ID → 終わったジョブの記録（実績を兼ねる。type finishedJob Job で別の型にする）
  コレクション meta       "next-id" → 次の job-id
  コレクション printers   プリンタの名前（ID に使える文字に符号化）→ 控え
  コレクション trust      ホスト（同じく符号化）→ sha256 の指紋
  Blob docs              "0000000042.original"（受けたまま）、"0000000042.bdf"（変換した .bdf）
```

- **スプールは jobs をメモリに写して使う**。起動時に jobs を全部読み、以後は自分の書き込みをそのまま写す（書くのは 1 プロセスだけ）。キューの順、時刻の来たジョブ、一覧はメモリで作る。yuruidb の一覧は ID の順だけなので、索引は持たない。
- **記録の書き込み**: yuruidb の Update に読んだときの Revision を渡す。1 プロセスなので競合は起きないはずで、起きたらバグとして誤りにする。書けたか分からない（`ErrOutcomeUnknown`）ときは、メモリの写しを読み直しで合わせる。
- **終わったジョブ**: done に Insert してから jobs を Delete する（コレクションをまたぐトランザクションは無い）。途中で落ちても、起動時に「done にもあるのに jobs に残る」ものを消せば揃う。jobs は小さく保たれる。
- **文書**: 受けた本文はまず一時ファイルに落とし、大きさと SHA-256 を数えてから Blob に入れる（受け付けの速さが相手の回線に左右されない。yuruidb の Blob は 2 MiB ずつ短い書き込みで入れるので、ほかの書き込みを止めない）。文書はジョブごとに持ち、ジョブをまたいで共有しない（参照の数を数えずに済み、GC が単純になる）。読むのは Blob の `io.ReaderAt`（短い読み取りのトランザクション）で、記録の SHA-256 と大きさで確かめる。
- **連番**: meta の `next-id` を Update で足してから使う（同じ id を 2 度使わない。飛ぶのは構わない）。IPP の `job-id` と同じ正の整数。
- **GC**（1 時間ごとと起動時）: done の文書を保持期間（既定 7 日）の後に Blob から消し、done の記録を実績の保持期間（既定 400 日）の後に消す。jobs の件数（既定 1 万）と文書の合計バイト（既定 10 GiB）を超えたら受け付けを断る（507）。どのジョブにも属さない Blob（受け付けの途中で落ちた）を Blob の一覧から探して消す。
- **ファイルの大きさ**: bbolt は消した分のページを使い回すが、ファイルは縮まない。`bdf printd -compact` が止めた状態で `bbolt.Compact` を使い、隣に新しいファイルを書いて置き換える。
- **写し**: `GET /api/v1/backup`（admin）が `Store.Backup` でいったんローカルの一時ファイルに一貫した写しを書いてから相手に流す（遅い相手に直接流すと、bbolt の読み取りのトランザクションが長く続き、ファイルを広げる書き込みを待たせるため）。
- **置き場所の既定**: 設定ファイルと同じディレクトリの `spool.db`（設定の `spool:` で変える）。NFS・SMB の共有には置かない（bbolt のロックと mmap が当てにならない）。

### 3.3 ジョブ（S2）

記録の中身（JSON）。

```jsonc
{
  "id": 42, "uuid": "urn:uuid:…",                 // uuid は m4 で job-uuid として見せる
  "queue": "office", "name": "invoice.pdf", "user": "erp",
  "document": { "sha256": "…", "size": 123456, "type": "pdf",      // docs/42/original
                "bdf": { "sha256": "…", "size": 80211 } },       // 変換した .bdf（docs/42/bdf）。.bdf を受けたときは無い
  "ticket": { "sides": "two-sided-long-edge", "copies": 2 },   // キューの既定を重ねる前の、受けたままのもの
  "options": { "format": "auto", "view": "", "skip": ["background"], "strict": false, "raster-if-lossy": false },
  "state": "processing", "reasons": ["job-printing"], "message": "",
  "priority": 50, "hold-until": "", "hold-until-time": null,
  "step": "watch",                                 // convert | send | watch | ""（終わった、保留）
  "due": "2026-10-19T10:00:12Z",                   // 次に動かす時刻（再送の待ち、保留の終わり）。見届けの間隔はメモリだけ
  "attempts": 1, "retry-since": null,
  "printer": { "name": "mfp1", "url": "ipps://…/ipp/print", "attempt-started": "…",
               "job-id": 12, "job-uri": "…", "sent": true },
  "submitted": "…", "started": "…", "finished": null,
  "sheets": 3, "impressions": 0, "bytes": 14623, "format": "application/pdf", "media": "iso_a4_210x297mm",
  "unsupported": [], "warnings": []
}
```

状態機械（IPP の `job-state`）。理由は IPP の語彙で、無いものだけ `bdf-` を付ける。

| 状態 | いつ | 主な理由 |
|---|---|---|
| pending（3） | 受けた、順番待ち、再送の待ち | `job-incoming`、`job-queued`、`printer-stopped`、`bdf-retry-wait` |
| pending-held（4） | 保留（API、`job-hold-until`、パスワード待ち、照合できない） | `job-hold-until-specified`、`job-password-wait`、`bdf-duplicate-check` |
| processing（5） | 変換・描画・送信・プリンタが処理中 | `job-transforming`、`job-printing`、`job-queued-for-marker` |
| processing-stopped（6） | プリンタが止まった（送った後） | `printer-stopped`、プリンタの理由（`media-empty-error` …） |
| canceled（7） | 取り消し（API、プリンタ側） | `job-canceled-by-user`、`job-canceled-at-device` |
| aborted（8） | プリンタが中止、文書が壊れている、受けない（strict）、再送の上限 | `aborted-by-system`、`document-format-error`、`bdf-unsupported`、`bdf-retry-exhausted` |
| completed（9） | 終わった | `job-completed-successfully`、`job-completed-with-warnings` |

遷移は `job.go` の表で検査し、表に無い遷移はエラー（バグを早く見つける）。終わった状態から出られるのは `retry`（aborted と canceled から pending へ。新しい試行として `attempts` を足し、`printer` を空にする）だけ。

`job-hold-until` は CUPS と同じ時間帯（day-time 6〜18 時、evening 18 時〜、night 18〜6 時、second-shift 16〜24 時、third-shift 0〜8 時、weekend 土日）と `indefinite`、`job-hold-until-time`（時刻）。時間帯はサーバーの地方時。時間帯の保留は `due` に解ける時刻を置くだけで、回し手が時刻に起こす。

### 3.4 段（S3）

1 ジョブの処理を、どれも短く、途中で落ちてもやり直せる 3 つの段に分ける。

| 段 | すること | 次 |
|---|---|---|
| `convert` | 文書が `.bdf` でなければ変換して `docs/<id>/bdf` に置く（§3.8）。キューのプリンタが決まる前に行う（プリンタが前のジョブを刷っている間に変換が済む） | `send`（時刻は今） |
| `send` | パスワードの確認、キューの既定とチケットを重ね、`print.Print`（`Wait` 0、`OnCreated` と `OnSent` で記録）。送り終えたところで返る | `watch`（時刻は 2 秒後） |
| `watch` | `print.Check` で 1 回聞き、状態と理由を写す。終わっていれば実績を書いて終わり | 終わっていなければ `watch`（2 秒から 10 秒に延ばす。止まっていれば 10 秒） |

- **1 ジョブを動かすのは 1 つの goroutine**: 回し手がジョブを段に渡すとき、メモリの「動かしている」印を付ける（1 プロセスなので、記録に持ち主を書く必要は無い）。段を終えたら次の段と時刻を記録に書いて印を外す（1 つの書き込み）。
- **書く回数**: 記録に書くのは状態・理由・段が変わったときだけ。見届けの次の時刻はメモリに持つ。
- **回し手**（`Runner`）: 段を goroutine で回す。時刻の来たジョブをメモリの写しから取り、`convert` は描画の枠で、`send` はプリンタごとの送り手で、`watch` はプリンタごとの見届けで行う。変更の購読と「次の時刻」の時計で起きる。
- **再起動の後**: bbolt のファイルのロックが取れた時点で前のプロセスは死んでいるので、`send` と `watch` の段で残ったジョブを §3.6 の照合に回す。

### 3.5 スケジューラ（S3）

- **キュー**: 設定のキューごとに、出られるジョブ（pending で、保留でなく、`due` を過ぎた）を `priority` の降順、`submitted` の昇順、id の昇順で並べる。キューは一時停止できる（受け付けは続け、流さない）。受け付け自体を止める（`accepting: false`）こともできる。
- **プリンタの送り手と見届け**: 設定のプリンタごとに送り手 1 本と見届け 1 本。送り手は、そのプリンタで「送ったがまだ終わっていないジョブ」が `max-inflight`（既定 2）未満で、見張りが止まっていない・届くと言っているときに、そのプリンタを使えるキューから次のジョブを取って送る。送るのは 1 ジョブずつ順番どおり（プリンタの中でも順番が保たれる）。見届けは、そのプリンタに送ったジョブを時刻ごとに `Check` する。
- **組（クラス）**: キューに複数のプリンタを書けば組になる。`policy: any`（空いているプリンタから。送り手が取り合う）か `failover`（並べた順。前のプリンタが止まっている・届かないときだけ次）。送った後で止まったプリンタの中のジョブは `processing-stopped` のままそのプリンタに残る（別のプリンタに回すのは、取り消して `retry` してから `move`）。
- **描く数の上限**: 変換と PWG Raster の描画は CPU を使う。全体で同時に `renders`（既定は CPU 数）まで。PDF を送るジョブは書き出しが軽いので同じ枠で数える。
- **移す**: `Move` で別のキューへ。pending か pending-held のときだけ。

### 3.6 失敗と、二重印刷の防止（S3）

段の中で起きたこと。

1. 変換に失敗したら `aborted`（`document-format-error`）。
2. 暗号化した `.bdf` でメモリにパスワードが無ければ `pending-held`（`job-password-wait`）。
3. `Print` と `Check` の結果を写す: 終わった → `completed`（警告や受けない属性があれば `job-completed-with-warnings`）、`ErrAborted` → `aborted`、プリンタ側の取り消し → `canceled`、`ErrUnsupported`（strict）と `client-error-*` → `aborted`。
4. 接続の失敗・`server-error-busy`・`server-error-temporary-error`・`server-error-not-accepting-jobs`・`ErrPrinterStopped` で、まだ `Create-Job` が通っていなければ `pending` に戻して `due` を置く（10 秒、30 秒、1 分、2 分、5 分、以後 10 分ごと）。止まっているならプリンタを「止まった」として見張りに渡し、キューの他のプリンタ（組）に回す。再送の上限は時間で持つ（既定 24 時間で `aborted`、`bdf-retry-exhausted`）。
5. `Create-Job` の後で接続が切れたら、段は `watch` のまま（作り直さない）。見張りが届くと言うまで `due` を延ばす。

**再起動の後の照合**（`recover.go`）。`send` と `watch` の段で残ったジョブについて、

| 記録 | プリンタに聞いて | すること |
|---|---|---|
| `job-id` あり、`sent` あり | `Get-Job-Attributes` で見つかる | 終わっていれば写す、終わっていなければ `watch` を続ける |
| 〃 | 見つからない（履歴を持たない・消した） | `completed`（警告「プリンタが履歴を持たない」）。送り終えているので再送しない |
| `job-id` あり、`sent` なし | 見つかって pending（文書待ち） | プリンタのジョブを取り消して、新しい試行で送り直す |
| 〃 | 見つからない、または終わっている | 新しい試行で送り直す |
| `job-id` なし、`attempt-started` あり（`Create-Job` の応答を記録する前に落ちた） | `Get-Jobs`（which-jobs=all、my-jobs）で、同じ `job-name` と利用者で `attempt-started` より後のもの | 1 つなら上の表に入れて扱う。無ければ送り直す。2 つ以上なら `pending-held`（`bdf-duplicate-check`）にして人に任せる |
| `attempt-started` なし | 聞かない | 送り直す（まだ何も送っていない） |

`job-name` はプリンタによっては切られるので前方一致で比べる。`requesting-user-name` は常にスプールの利用者名（設定の `user`、既定 `bdf`）にし、`job-name` にジョブの名前を入れる（プリンタのパネルに見える名前）。

### 3.7 プリンタの見張り（S3）

プリンタごとに 1 本。`Get-Printer-Attributes` の状態の属性（`printer-state`、`printer-state-reasons`、`printer-state-message`、`printer-is-accepting-jobs`、`printer-config-change-time`）を、何もないとき 60 秒、ジョブを待たせているときと止まっているとき 10 秒ごとに聞く。届かないときは間隔を延ばす（最大 5 分）。`printer-config-change-time` が変わったら全部の能力を聞き直して `printers` に控える。`Create-Printer-Subscriptions` と `Get-Notifications` は m2 では使わない（機種の差が大きい）。

### 3.8 受けた入力の変換（S5）

`spool` は変換器に依らず、`spool.Options.Convert func(ctx, r io.ReaderAt, size int64, name string) (*bdf.Document, []string, error)` を受ける。`bdf printd` が `converter/all` で作って渡す。変換は `convert` の段（`job-transforming`）で行い、HTTP の受け付けは速く返す。HTML と Markdown の外部の画像は取りに行かない（`remote=false`。サーバーが内側の網に届く入口にならないため）。Excel は `views=print`。変換の予算は変換器のまま。結果の `.bdf` を `docs` に置いて記録の `document.bdf` から指し、再送で変換し直さない。

### 3.9 HTTP API（S4）

JSON。`/api/v1` の下。エラーは `{"error": {"code": "…", "message": "…"}}` と HTTP の状態。

| メソッドと道 | すること |
|---|---|
| `POST /queues/{q}/jobs` | 投入。`multipart/form-data` で `document`（ファイル、必須）、`ticket`（JSON）、`options`（JSON）、`name`、`password`（暗号化した文書）。201 とジョブ。受け付けを止めたキュー・容量の超過は 503 / 507 |
| `GET /jobs?queue=&state=&user=&limit=&after=` | 一覧。既定は終わっていないもの、`state=all` で全部。`after` は id で続きを取る |
| `GET /jobs/{id}` | 詳細 |
| `POST /jobs/{id}/cancel`・`hold`・`release`・`retry` | 取り消し・保留・保留の解除（`{"password": "…"}` も受ける）・再送 |
| `POST /jobs/{id}/move` | `{"queue": "…"}` |
| `PATCH /jobs/{id}` | `{"priority": n}` |
| `GET /queues`、`GET /queues/{q}` | キューと、その順番、プリンタ |
| `POST /queues/{q}/pause`・`resume`・`reject`・`accept` | 流すのを止める・再開・受け付けを止める・再開 |
| `GET /printers`、`GET /printers/{p}` | 設定のプリンタと状態、能力 |
| `GET /discover` | DNS-SD でネットワークのプリンタを探す（admin） |
| `GET /accounting?from=&to=&format=jsonl\|csv` | 実績の書き出し（manage） |
| `GET /backup` | `spool.db` の一貫した写し（admin） |
| `GET /events` | Server-Sent Events。ジョブ・キュー・プリンタの変化（`event: job`、`data: {ジョブ}`） |

`/healthz`（スプールが開いていれば 200）、`/metrics`、`/`（状態のページ）は `/api/v1` の外。

**認証と許可**:

- 設定の `tokens` に名前・トークン（ファイルか環境変数から）・使えるキュー・役割（`submit`、`manage`（自分以外のジョブの取り消しや保留、実績）、`admin`（キュー、発見、写し））を書く。`Authorization: Bearer …`。比べるのは定数時間。
- 束縛が loopback だけで `tokens` が空なら、認証なし（手元の 1 台向け）。loopback 以外に束縛するときは、TLS（`tls.cert`・`tls.key`）とトークンが無ければ起動しない（`insecure: true` で明示したときだけ許す）。
- ジョブの利用者はトークンの名前。トークンに `act-as: true` があるときだけ `X-BDF-User` を利用者として受ける（ERP が利用者に代わって出すとき）。
- 自分のジョブは `submit` で取り消せる。他人のジョブは `manage`。

**限界**: 文書の大きさ（`max-document`、既定 256 MiB。超えたら 413）、同時の受け付け（既定 8）、スプールの合計容量（超えたら 507）、`limit` は最大 1000、SSE の同時接続（既定 64）。ヘッダーと本文の読み込みに時間の上限。

### 3.10 設定（S4）

```yaml
listen: 127.0.0.1:8640
spool: /var/lib/bdf/spool.db      # 既定は設定ファイルと同じディレクトリの spool.db
compact-on-start: true
user: bdf                         # プリンタへの requesting-user-name
tls: { cert: /etc/bdf/tls.crt, key: /etc/bdf/tls.key }   # loopback 以外では必須
tokens:
  - { name: erp, file: /etc/bdf/erp.token, queues: [office, warehouse], roles: [submit], act-as: true }
  - { name: admin, env: BDF_ADMIN_TOKEN, roles: [submit, manage, admin] }
printers:
  mfp1: { url: ipps://mfp1.example.local/ipp/print, fingerprint: "sha256:…" }   # 書けば固定、無ければ初回の指紋
  mfp2: { url: ipp://192.0.2.20/ipp/print, username: print, password-env: MFP2_PASSWORD }
queues:
  office:
    printers: [mfp1, mfp2]
    policy: failover                # any | failover
    defaults: { sides: two-sided-long-edge, print-color-mode: monochrome }
    format: auto                    # auto | pdf | pwg-raster
    raster-if-lossy: true
    max-inflight: 2
retention: { done: 168h, max-jobs: 10000, max-bytes: 10GiB, accounting: 400d }
limits: { max-document: 256MiB, max-impressions: 10000, renders: 0, retry-for: 24h }
fonts: { dirs: [/usr/share/fonts/bdf], no-system: false }
```

`bdf printd -check` は設定を読み、プリンタに `Get-Printer-Attributes` を送って届くかと能力を表にして終わる（`spool.db` は開かない）。設定の変更は再起動で反映する（m2 では再読み込みしない）。

### 3.11 クライアントと CLI（S5）

- `print/server/client`: `Submit(ctx, queue, name, r io.Reader, ticket, opts) (*Job, error)`（multipart を流しながら送る）、`Job`、`Jobs`、`Cancel` など、`Accounting`、`Watch(ctx, id)`（SSE。切れたらポーリングに落ちる）。JSON の型（`Job`、`Queue`、`Printer`）はこのパッケージに置き、`server` もそれを使う。
- `bdf print -server URL -queue Q`: チケットのフラグは m1 と同じ。既定は投入して帰る（ジョブの id を出す）。`-wait` で終わるまで待ち、終了コードは m1 と同じ意味（1 取り消し・中止、3 サーバーに届かない、4 拒まれた、5 キューが受け付けない、6 時間切れ）。`-server` のときは文書を変換せずにそのまま送る（変換はサーバー）。
- `bdf jobs`: 一覧（id、状態、理由、キュー、名前、利用者、時刻）と `show`・`cancel`・`hold`・`release`・`retry`・`accounting`。`-json`。
- `bdf printers -server`: サーバーのキューとプリンタと状態。

### 3.12 観測（S4）

- ログは `log/slog` の JSON を標準エラーへ（`-log text` で人向け）。ジョブの受け付け・状態の変化・再送・引き継ぎ・プリンタの状態の変化。文書の中身、パスワード、トークン、URL の利用者情報は書かない。
- `/metrics`（Prometheus のテキスト形式を自前で書く。依存を足さない）: `bdf_print_jobs_total{queue,state}`、`bdf_print_jobs{queue,state}`（今の数）、`bdf_print_printer_up{printer}`、`bdf_print_printer_state{printer}`、`bdf_print_sheets_total{queue}`、`bdf_print_sent_bytes_total{printer}`、`bdf_print_retries_total{printer}`、`bdf_print_spool_bytes`（生きている文書）、`bdf_print_spool_file_bytes`（`spool.db` の大きさ）。
- 状態のページ: キューごとの順番とプリンタの状態とジョブ。`/events` で更新。読むだけ（操作は API と CLI）。

### 3.13 常駐のさせ方（S6、文書）

`bdf printd` は前面で動くプロセスで、サービスとしての登録は持たない。止めるのは SIGINT・SIGTERM（Windows は Ctrl-C と Ctrl-Break。`os/signal` の `os.Interrupt`）で、受け付けを止め、送信中の段を終えてから閉じる（既定 30 秒で打ち切る。打ち切ったジョブは次の起動で引き継ぎの照合を通る）。docs に例を書く: systemd の unit（`Restart=on-failure`、`StateDirectory=bdf`）、macOS の launchd の plist、Windows は NSSM か WinSW でサービスにする例とタスク スケジューラの「起動時」の例。

### 3.14 サンプル（S6）

`examples/print-server`: `print/server` をライブラリとして埋め込んだ小さなサーバー。請求書のテンプレート（`converter/xlsx/testdata/form-pages.xlsx`）とデータの入力欄のページを出し、送られたデータで `form.Fill` して `.bdf` を作り、スプールに置き、ジョブの状態を `/events` で見せる。プリンタが無い手元でも、`ipptest` の偽のプリンタを同じプロセスで立てる `-fake-printer` で動く。`docs/examples/print-server(.ja).md` と構成のサンプルの一覧に足す。

## 4. テスト

- **store.go**（`go test ./print/spool/ -run Store`。一時ディレクトリの `spool.db`）: 記録の往復、連番が戻らない、done への移しと途中で落ちたときの揃え直し、文書の往復（0 バイト、範囲の読み取り、SHA-256 の確かめ）、持ち主の無い文書の掃除、GC。閉じて開き直しても同じ、2 つ目の起動がロックで断られる、詰め直しと写しで中身が変わらない。
- **単体**（`go test ./print/spool/`）: 状態機械の遷移の表（許される遷移と許されない遷移）、保留の時間帯（偽の時計で）、優先度と投入順、再送の間隔（偽の時計）。
- **偽のプリンタで**（`ipptest`）: 2 台の組（any と failover、1 台を止める・再開する）、`max-inflight` 2 でプリンタの中に 2 つまで、プリンタ側の中止で再送しない、取り消し（pending・processing・プリンタに送った後）、保留と解除、別のキューへの移し、変換して印刷（PDF と TIFF の入力）、暗号化した文書とパスワード待ち。
- **二重印刷の防止**: テストのバイナリを子プロセスで `bdf printd` の代わりに動かし、決めた時点で `os.Exit` させてから、同じ `spool.db` を開き直す（本当に落ちた状態から始める）。落ちる時点は `Create-Job` の前、`Create-Job` の応答の前（`DropCreateResponse`）、記録の後で送る前、送る途中（`FailAfter`）、送った後で見届ける間、履歴を持たないプリンタ（`Forget`）。どの場合もプリンタが刷るのは 1 回（取り消した不完全なジョブは数えない）。
- **HTTP**（`httptest`）: 認証と許可（トークンなし、別のキュー、manage でない他人のジョブ、act-as）、loopback 以外で TLS が無いと起動しない、大きさの上限（413）、容量の超過（507）、SSE、写し、エラーの形。
- **クライアントと CLI**: `bdf print -server`（帰る・待つ・終了コード）、`bdf jobs` の各操作、`bdf printd -check`。
- **競合**: `go test -race`。投入・取り消し・状態の変化・GC を同時に。
- **3 OS**: CI の `print-os`（Windows と macOS）で `go test ./print/...`。bbolt のロックと詰め直しの置き換え（Windows は開いているファイルを改名できない）を確かめる。
- **CI**: `test/print/ci.sh` に、`bdf printd` を `ippeveprinter` につないで起動し、`bdf print -server` で 2 つ投入し、1 つを取り消し、`bdf printd` を止めて再起動し、残りが完了して `ippeveprinter` が受けた文書が 1 つ（二重でない）ことを確かめる段を足す。

## 5. 順番

1. **S1** `print` の継ぎ目（`OnCreated`、`OnSent`、`Printer`、`Check`、`Ticket` の JSON と `Merge`）、`TrustBackend`、`ipptest` の足し算。小さく、ほかの全部の前提。
2. **S2** `store.go`（yuruidb）、`job`、`hold`、`account`、`gc` と `Spool` の受け付け・一覧・操作（まだ流さない）。
3. **S3** 段、回し手、見張り、再起動の後の照合。ここで「置けば刷られる」になる。二重印刷の防止のテストはここ。
4. **S4** `print/server`（設定、API、認証、SSE、/metrics、状態のページ）と `bdf printd`。
5. **S5** `client`、`bdf print -server`、`bdf jobs`、`bdf printers -server`、受けた入力の変換、暗号化した文書。
6. **S6** `examples/print-server`、docs（常駐のさせ方を含む）、CI、`.knowledge` の記録。

S2 と S4 の HTTP の骨組みは並行できる。S3 が山。

## 6. 後で変えられる位置に置くもの

**プロセスの中の並行**は m2 で持つ。プリンタごとに送り手と見届け、変換と描画は CPU 数の枠。bbolt は読み取りが並行、書き込みは 1 つずつだが、1 回の書き込みは記録 1 件か文書の数 MiB で短い。並びの順は（優先度、受けた時刻、id）。

**保存先を替える・複数台で回す・Lambda で動かす**のは m2 の後。調べたことと案は `store.ja.md` にある。

- 1 つの `spool.db` は 1 プロセスだけが開ける。共有のファイルシステム（NFS、SMB、EFS）に置いて複数台で開くのも不可。
- 保存先は yuruidb のドライバ（S3・R2・GCS・DynamoDB・RDB、文書は S3・R2・GCS・RDB の Blob）を設定で選べば替えられる。`store.go` が yuruidb の `Store` と `blob.Store` を受け取る形にしておく。
- 複数台で 1 つの保存先を回すには、ジョブの記録に持ち主と期限の「借り」を足す。時刻だけでは 2 台の取り合いを防げないので、勝つのは条件付きの書き込み（yuruidb の Revision つきの Update）の 1 つにし、期限は落ちた持ち主から引き継ぐために使う。
- 段（§3.4）と再起動の後の照合（§3.6）は、そのまま「借りが切れたジョブの引き継ぎ」になる。
- Lambda は、受け付けと変換をクラウド（Lambda＋オブジェクトストレージ）に、印刷を拠点の `bdf printd` に分ける形が現実的（事務所の LAN のプリンタにクラウドから届かないため）。

ほかに替えられる位置。

- パスワードの持ち方（メモリ → サーバーの鍵で封印）は `spool.Options.Secrets`（インターフェース）。
- m4 の IPP サーバーは、`server` に `/ipp/print/<queue>` の IPP の受け口を足し、`Spool.Submit` に渡す。ジョブの状態は IPP の語彙のままなので写すだけ。`spool` の id は IPP の `job-id` の形（正の整数）、`uuid` は `job-uuid` になる。

## 7. m1 から変わるもの

- `print.Options` に `OnCreated`、`OnSent`、`Printer` が増え、`print.Check` が公開される（m1 の `wait` はその繰り返しに書き直す）。m1 の `bdf print`（サーバーなし）の動きは変わらない。
- `Ticket` が JSON になる（タグを足すだけ。Go の利用者には影響なし）。
- `ipp.TrustStore` に `Backend` が増える（nil なら今どおり）。
- `bdf print` に `-server`、`-queue`、`-token-file` が増える。`-server` が無ければ m1 のまま。`bdf printd` と `bdf jobs` が増える。
- `go.mod` に `github.com/shibukawa/yuruidb` が増える（と、それが引く `go.etcd.io/bbolt`、`golang.org/x/sys`）。

## 8. 未決の論点（決めてほしいこと）

決まったもの（2026-10-19〜20 利用者）: サーバーは `bdf printd`、永続化は組み込みの KV（bbolt の 1 ファイル）、保存は yuruidb を通す（bbolt ドライバと bbolt の Blob）、変換の入力を受ける、`max-inflight` の既定は 2、サービスとしての登録は持たない、m2 の後は後で決める。

残るもの。どれも暫定で推奨どおり進める。

1. **文書も `spool.db` の中に置く**（推奨。yuruidb の bbolt の Blob。1 ファイルで写しと移しが楽）か、文書だけ隣のディレクトリに置くか。中に置くと、ファイルは消した分だけ縮まず、詰め直しが要る（§3.2）。文書はジョブごとに持ち、ジョブをまたいで共有しない（要件の「内容アドレスで同じ文書は 1 回」を改める。GC が単純になる）。
2. **パスワードの持ち方**: メモリだけで、再起動の後は `job-password-wait`（推奨）か、サーバーの鍵で封印して記録に置くか。
3. **既定の束縛と認証**: `127.0.0.1:8640` で、loopback かつトークンが無ければ認証なし、でよいか。
4. **`bdf print -spool dir` を作らない**こと（要件からの変更。スプールの持ち主を 1 プロセスにする）。
