# 印刷 実装案（print-m2）

2026-10-19 時点の実装案。要件と決定は `print-requirements.ja.md` §4.5〜§4.10 と `.knowledge/` の `requirement:print`（facts の spool・server・budgets）、`milestone:print-m2` にある。前段の print-m1（`bdf print`、`bdf printers`、`print`・`print/ipp`・`print/raster`）は実装して main に入った（#104）。ここは「何をどの順で作るか」を書く。

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
- ただし、後で複数台や Lambda に広げられるように、スプールの中を 3 つの継ぎ目で書く（§3.4、§6）。(1) 保存は `Store` のインターフェースで、ジョブの記録は版つきの条件付き書き込み（版が変わっていなければ書く）。(2) ジョブを取るのは期限つきの「借り」（持ち主と期限の時刻）。(3) 1 ジョブの処理は短い「段」（変換・送信・見届け）の繰り返しで、次に動かす時刻を記録に持つ。m2 で作る `Store` は bbolt とメモリ（テスト用）の 2 つだけ。

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
  spool.go         Spool: Open（Store、時計、持ち主の名前）、Submit、Get、List、Cancel、Hold、Release、Retry、Move、
                   SetPriority、SetPassword、Subscribe（変更の購読）、Close                                    (S2)
  job.go           Job（記録の中身）と状態機械（IPP の job-state と理由）、遷移の検査                           (S2)
  store.go         Store のインターフェース（記録の版つき書き込み、索引での一覧、連番、文書、小さな記録）と
                   ErrConflict                                                                                 (S2)
  memstore.go      メモリの Store（テストと、2 つの Spool が 1 つの Store を取り合うテスト）                     (S2)
  hold.go          job-hold-until の時間帯（indefinite、day-time、evening、night、weekend、second-shift、third-shift）
                   と job-hold-until-time                                                                      (S2)
  account.go       実績（終わったジョブごとに 1 件）と書き出し（JSON Lines、CSV）                               (S2)
  lease.go         借り: 取る（条件付き書き込み）、延ばす、返す、期限切れの判定                                  (S3)
  step.go          段: convert、send、watch。Step(ctx, id) が 1 段を行い、次の段と時刻を書く                      (S3)
  takeover.go      持ち主が替わったジョブの照合（再起動もこれ。二重印刷の防止）                                   (S3)
  sched.go         回し手（Runner）: キューの順、プリンタごとの送り手と見届け、max-inflight、組、描画の枠        (S3)
  monitor.go       プリンタの見張り: printer-state と理由、能力の更新（printer-config-change-time）              (S3)
  gc.go            保持期間・件数・バイトで消す、参照されない文書を消す                                         (S2)
print/spool/storetest/
  storetest.go     Store の約束のテスト（どの実装にも同じものを流す）                                           (S2)
print/spool/boltstore/
  boltstore.go     bbolt の Store（バケットの配置、文書の分割、索引、詰め直し、写しの書き出し）                   (S2)
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

`spool` は `print` と `print/ipp` に依り、bbolt には依らない（依るのは `boltstore` だけ）。`server` は `spool` に依り、`client` はどれにも依らない（API の JSON の型だけを持つ）。`converter` に依るのは `cmd/bdf` と `examples` だけ。新しい依存は `go.etcd.io/bbolt`（v1.5.0。実行時に引くのは `golang.org/x/sys` だけ）。

## 3. 各部の設計

### 3.1 print パッケージに足すもの（S1）

m1 の `print.Print` は 1 回の呼び出しで聞く・作る・送る・待つ。スプールはこれを「送るまで」と「見届け」に分けて使う（`Options.Wait` が 0 なら送った時点で返る。m1 のまま）。

- `Options.OnCreated func(*ipp.JobStatus) error`: `Create-Job` の応答を受けた直後、`Send-Document` の前に同じ goroutine で呼ぶ。スプールはここでプリンタの `job-id` と `job-uri` を記録する（書けなければ、あるいは借りを失っていれば、エラーを返し、`Print` はプリンタのジョブを取り消して返す）。`Progress` の `created` は今どおり残す。
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
| ジョブの記録: id、uuid、キュー、利用者、名前、チケット、状態と理由、優先度と保留、段と次に動かす時刻、借り、試行の回数、プリンタの名前とプリンタ側の `job-id` / `job-uri` と送り終えた印（二重印刷の防止の要）、時刻、枚数、警告 | 1 件 1〜2 KB | 終わってから保持期間（既定 7 日） |
| 文書: 受けたままのバイト列（`.bdf`、PDF、docx …）と、変換した `.bdf` | ジョブの大部分（KB〜数百 MB） | 参照するジョブが無くなるまで |
| 連番（次の `job-id`）、形式の版 | 数バイト | ずっと |
| 実績: 終わったジョブごとに時刻・利用者・キュー・プリンタ・枚数・印象数・用紙・両面・色 | 1 件 300 B | 既定 400 日 |
| プリンタの控え: 最後に聞いた能力と状態（`printer-config-change-time` つき） | 1 台 数十〜数百 KB | 設定にある間 |
| ipps の証明書の指紋 | 1 台 80 B | ずっと |

保存しないもの: 描いた Raster（送るたびに描く）、パスワード（メモリだけ）、トークン（設定ファイルとその参照先）、設定（YAML のまま）。

置き場所は bbolt の 1 ファイル。

```
spool.db（bbolt）
  meta        format（形式の版）、next-id、owner（今開いている持ち主の名前）
  jobs        id（8 バイトの大端）→ ジョブの記録（JSON）
  idx-due     次に動かす時刻＋id → 空      回し手が「時刻が来たジョブ」を順に取る
  idx-queue   キュー＋状態＋id → 空        一覧と、キューの順番の表示
  idx-done    終わった時刻＋id → 空        保持期間の GC
  docs        sha256 → 文書のバケット { info: {size、type、refs}、0000000000: 1 MiB、0000000001: 1 MiB、… }
  uploads     受け付け中の文書（起動時に残っていれば消す）
  acct        終わった時刻＋id → 実績の 1 件（JSON）
  printers    名前 → 控え（JSON）
  trust       ホスト → sha256 の指紋
```

- **記録の書き込み**: ジョブの記録と索引は 1 つの書き込み（bbolt のトランザクション）で変える。記録には `rev` があり、`PutJob` は読んだときの `rev` と今の `rev` が同じときだけ書いて 1 足す（違えば `ErrConflict`）。1 プロセスの bbolt では書き込みは元から 1 つずつだが、同じ約束で書いておけば、複数台の Store（§6）に替えても `spool` は変わらない。
- **文書**: 受けながら SHA-256 を計算し、1 MiB ずつ `uploads` に書く（数 MiB ごとに 1 つの書き込みにして、他の書き込みを長く止めない）。受け終わったら 1 つの書き込みで `docs/<sha256>` に移し（`Tx.MoveBucket`）、ジョブの記録を作る。同じ文書が既にあれば `refs` を足して、受けた分は捨てる（1000 枚の請求書の束を 3 回送っても 1 つ）。
- **文書を読む**: `io.ReaderAt`。1 MiB の塊ごとに短い読み取りで取り出し、直前の塊を 1 つだけ覚える。長い読み取りのトランザクションを開いたままにしない（bbolt では長い読み取りが空きページの再利用とファイルの伸長を止めるため）。文書は書いた後は変わらず、使っているジョブがある間は GC が消さないので、塊ごとに読んでも一貫する。
- **連番**: `next-id` を同じ書き込みで足してから使う（同じ id を 2 度使わない。飛ぶのは構わない）。IPP の `job-id` と同じ正の整数。
- **GC**（1 時間ごとと起動時）: 終わったジョブ（canceled・aborted・completed）の記録を保持期間の後に消す。件数（既定 1 万）と文書の合計バイト（既定 10 GiB）を超えたら古い順に消す。`refs` が 0 の文書を消す。実績は別の保持期間（既定 400 日）。
- **ファイルの大きさ**: bbolt は消した分のページを使い回すが、ファイルは縮まない。起動時に空きが全体の半分以上かつ 1 GiB 以上なら詰め直す（`bbolt.Compact` で隣に新しいファイルを書き、`fsync` して置き換える。設定 `compact-on-start` で切れる）。`bdf printd -compact` は止めた状態で同じことをする。
- **写し**: `GET /api/v1/backup`（admin）が、動いたまま一貫した `spool.db` の写しを流す（読み取りのトランザクションの `WriteTo`）。1 ファイルにした利点。
- **置き場所の既定**: 設定ファイルと同じディレクトリの `spool.db`（設定の `spool:` で変える）。NFS・SMB の共有には置かない（bbolt のロックと mmap が当てにならない。docs に書き、起動時にネットワークのファイルシステムと分かれば警告する）。

### 3.3 ジョブ（S2）

記録の中身（JSON）。

```jsonc
{
  "id": 42, "uuid": "urn:uuid:…", "rev": 17,      // uuid は m4 で job-uuid として見せる。rev は書くたびに 1 足す
  "queue": "office", "name": "invoice.pdf", "user": "erp",
  "document": { "hash": "sha256-…", "size": 123456, "type": "pdf", "bdf": "sha256-…" },  // bdf は変換した .bdf
  "ticket": { "sides": "two-sided-long-edge", "copies": 2 },   // キューの既定を重ねる前の、受けたままのもの
  "options": { "format": "auto", "view": "", "skip": ["background"], "strict": false, "raster-if-lossy": false },
  "state": "processing", "reasons": ["job-printing"], "message": "",
  "priority": 50, "hold-until": "", "hold-until-time": null,
  "step": "watch",                                 // convert | send | watch | ""（終わった、保留）
  "due": "2026-10-19T10:00:12Z",                   // 次に動かす時刻（再送の待ち、見届けの間隔、保留の終わり）
  "lease": { "owner": "print1/7f3c", "until": "2026-10-19T10:02:00Z" },   // 取っている者と期限
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

### 3.4 段と借り（S3）

1 ジョブの処理を、どれも短く、途中で落ちてもやり直せる 3 つの段に分ける。

| 段 | すること | 次 |
|---|---|---|
| `convert` | 文書が `.bdf` でなければ変換して `docs` に置く（§3.8）。キューのプリンタが決まる前に行う（プリンタが前のジョブを刷っている間に変換が済む） | `send`（時刻は今） |
| `send` | パスワードの確認、キューの既定とチケットを重ね、`print.Print`（`Wait` 0、`OnCreated` と `OnSent` で記録）。送り終えたところで返る | `watch`（時刻は 2 秒後） |
| `watch` | `print.Check` で 1 回聞き、状態と理由を写す。終わっていれば実績を書いて終わり | 終わっていなければ `watch`（2 秒から 10 秒に延ばす。止まっていれば 10 秒） |

- **借り**: 段を行う前に、記録の `lease` を自分（持ち主の名前 `ホスト名/起動ごとの乱数`）と期限（今から 2 分）にして条件付きで書く。書けた者だけがその段を行う。長い段（大きな文書の送信）は 30 秒ごとに期限を延ばす。段を終えたら借りを外して次の段と時刻を書く（1 つの書き込み）。
- **時刻だけでは足りない**: 2 つの回し手が同時に「時刻の来た 42 番」を見ても、条件付きの書き込みで勝つのは 1 つだけ。期限の時刻は「持ち主が落ちたとき、ほかの者が引き継げる」ためにある。期限が切れて引き継がれた後に古い持ち主が書こうとしても、`rev` が変わっているので書けない（`OnCreated` がエラーを返し、プリンタのジョブは取り消される）。
- **1 プロセスの `bdf printd`**: bbolt のロックで同じファイルを開いている者は自分しかいないので、起動時に前の持ち主の借りはすべて切れたものとして、すぐ引き継ぐ（期限を待たない）。引き継ぎは §3.6 の照合を通る。
- **回し手**（`Runner`）: 常駐の `bdf printd` では、段を goroutine で回す。時刻の来たジョブを `idx-due` から取り、`convert` は描画の枠で、`send` はプリンタごとの送り手で、`watch` はプリンタごとの見届けで行う。変更の購読と「次の時刻」の時計で起きる。Lambda なら同じ `Step(ctx, id)` を時刻起動で 1 回ずつ呼ぶ形になる（§6）。

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

**引き継ぎの照合**（`takeover.go`）。前の持ち主の借りのまま残った `send` と `watch` のジョブ（再起動の後はこれ）について、

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

- **Store の約束**（`storetest`。memstore と boltstore の両方に流す）: 版つきの書き込み（古い `rev` は `ErrConflict`）、索引の順、連番が戻らない、文書の往復（塊の境目、0 バイト、同じ文書の `refs`）、受け付け中に落ちた文書の掃除、GC。boltstore は閉じて開き直しても同じ、2 つ目の起動がロックで断られる、詰め直しの前後で中身が同じ、写しから開ける。
- **単体**（`go test ./print/spool/`）: 状態機械の遷移の表（許される遷移と許されない遷移）、借り（取る・延ばす・切れる・古い持ち主の書き込みが断られる）、保留の時間帯（偽の時計で）、優先度と投入順、再送の間隔（偽の時計）。
- **偽のプリンタで**（`ipptest`）: 2 台の組（any と failover、1 台を止める・再開する）、`max-inflight` 2 でプリンタの中に 2 つまで、プリンタ側の中止で再送しない、取り消し（pending・processing・プリンタに送った後）、保留と解除、別のキューへの移し、変換して印刷（PDF と TIFF の入力）、暗号化した文書とパスワード待ち。
- **二重印刷の防止**: テストのバイナリを子プロセスで `bdf printd` の代わりに動かし、決めた時点で `os.Exit` させてから、同じ `spool.db` を開き直す（本当に落ちた状態から始める）。落ちる時点は `Create-Job` の前、`Create-Job` の応答の前（`DropCreateResponse`）、記録の後で送る前、送る途中（`FailAfter`）、送った後で見届ける間、履歴を持たないプリンタ（`Forget`）。どの場合もプリンタが刷るのは 1 回（取り消した不完全なジョブは数えない）。
- **取り合い**: 1 つの memstore を 2 つの `Spool`（別の持ち主の名前）で同時に回し、100 のジョブを 2 台の偽のプリンタに流して、どのジョブも 1 回だけ送られること。片方を途中で止め、借りの期限の後にもう片方が引き継いで終わらせること（偽の時計）。m2 の `bdf printd` では起きないが、§6 の継ぎ目が本物であることを確かめる。
- **HTTP**（`httptest`）: 認証と許可（トークンなし、別のキュー、manage でない他人のジョブ、act-as）、loopback 以外で TLS が無いと起動しない、大きさの上限（413）、容量の超過（507）、SSE、写し、エラーの形。
- **クライアントと CLI**: `bdf print -server`（帰る・待つ・終了コード）、`bdf jobs` の各操作、`bdf printd -check`。
- **競合**: `go test -race`。投入・取り消し・状態の変化・GC を同時に。
- **3 OS**: CI の `print-os`（Windows と macOS）で `go test ./print/...`。bbolt のロックと詰め直しの置き換え（Windows は開いているファイルを改名できない）を確かめる。
- **CI**: `test/print/ci.sh` に、`bdf printd` を `ippeveprinter` につないで起動し、`bdf print -server` で 2 つ投入し、1 つを取り消し、`bdf printd` を止めて再起動し、残りが完了して `ippeveprinter` が受けた文書が 1 つ（二重でない）ことを確かめる段を足す。

## 5. 順番

1. **S1** `print` の継ぎ目（`OnCreated`、`OnSent`、`Printer`、`Check`、`Ticket` の JSON と `Merge`）、`TrustBackend`、`ipptest` の足し算。小さく、ほかの全部の前提。
2. **S2** `Store`（インターフェース、memstore、boltstore、約束のテスト）、`job`、`hold`、`account`、`gc` と `Spool` の受け付け・一覧・操作（まだ流さない）。
3. **S3** 段と借り、回し手、見張り、引き継ぎの照合。ここで「置けば刷られる」になる。二重印刷の防止と取り合いのテストはここ。
4. **S4** `print/server`（設定、API、認証、SSE、/metrics、状態のページ）と `bdf printd`。
5. **S5** `client`、`bdf print -server`、`bdf jobs`、`bdf printers -server`、受けた入力の変換、暗号化した文書。
6. **S6** `examples/print-server`、docs（常駐のさせ方を含む）、CI、`.knowledge` の記録。

S2 と S4 の HTTP の骨組みは並行できる。S3 が山。

## 6. 後で変えられる位置に置くもの（複数台・Lambda を含む）

**プロセスの中の並行**は m2 で持つ。プリンタごとに送り手と見届け、変換と描画は CPU 数の枠。bbolt は読み取りが並行、書き込みは 1 つずつだが、1 回の書き込みは記録 1 件か文書の数 MiB で短い。並びの順は（優先度、受けた時刻、id）で、時刻は順番を決めるのに使う。

**1 つのファイルを複数のプロセスで開くことはできない**。bbolt はファイルをロックするので、2 つ目のプロセスは開けない。共有のファイルシステム（NFS、SMB、EFS）に置いて複数台で開くのも、ロックと mmap が当てにならないので不可。複数台で回すには `Store` を、条件付きの書き込みを持つ共有の保存先に替える。

- **替える場所**: `print/spool/boltstore` の隣に別の `Store` を足すだけ。`spool` は `Store` のインターフェース（版つきの記録、索引での一覧、連番、文書）と借り・段の形にだけ依るので変わらない。m2 のテスト（取り合い）がそれを確かめている。
- **候補**: DynamoDB（条件付き書き込み）＋ S3（文書。署名つき URL で受ける）、または PostgreSQL（`SELECT … FOR UPDATE SKIP LOCKED` と BLOB か S3）。
- **時刻の注意**: 借りの期限は台ごとの時計に依る。期限は余裕をとり（2 分、30 秒ごとに延ばす）、最後の砦は `rev`（古い持ち主の書き込みを断る）と、送る前のプリンタへの照合（§3.6）。

**Lambda で動かすとき**（m2 では作らない。いつやるかは m2 の後で決める）。

- **そのまま載るもの**: 受け付け（API Gateway＋Lambda。大きな文書は S3 の署名つき URL に直接上げる。API Gateway の本文の上限が 10 MB のため）、変換と描画（Go は `provided.al2023` で動く。フォントはコンテナ イメージに入れる）、段の実行（`Step(ctx, id)` を EventBridge Scheduler や SQS の遅延で `due` の時刻に呼ぶ。1 段は 15 分の上限に収まる。待つのは段と段の間で、Lambda は待たない）。
- **そのままでは載らないもの**: (1) 事務所の LAN のプリンタにクラウドから届かない（私設アドレス、NAT）。VPN か Direct Connect と VPC の中の Lambda が要る。普通は拠点に `bdf printd` を置き、クラウドのスプールからジョブを取りに行かせる（拠点の代理。外向きの HTTPS だけで済む。`client` の上に作れる。IPP INFRA に合わせるかは後の論点 `remote-proxy`）。(2) SSE の常時接続 → ポーリングか別の通知。(3) メモリだけのパスワード → KMS で封印（`spool.Options.Secrets`）。(4) 状態のページと `/metrics` → CloudWatch など。
- だから現実的な形は「受け付けと変換をクラウド（Lambda＋S3＋DynamoDB）に、印刷を拠点の `bdf printd` に」。m2 はその土台（`Store`、借り、段、`client`）までを作る。

ほかに替えられる位置。

- パスワードの持ち方（メモリ → サーバーの鍵や KMS で封印）は `spool.Options.Secrets`（インターフェース）。
- m4 の IPP サーバーは、`server` に `/ipp/print/<queue>` の IPP の受け口を足し、`Spool.Submit` に渡す。ジョブの状態は IPP の語彙のままなので写すだけ。`spool` の id は IPP の `job-id` の形（正の整数）、`uuid` は `job-uuid` になる。

## 7. m1 から変わるもの

- `print.Options` に `OnCreated`、`OnSent`、`Printer` が増え、`print.Check` が公開される（m1 の `wait` はその繰り返しに書き直す）。m1 の `bdf print`（サーバーなし）の動きは変わらない。
- `Ticket` が JSON になる（タグを足すだけ。Go の利用者には影響なし）。
- `ipp.TrustStore` に `Backend` が増える（nil なら今どおり）。
- `bdf print` に `-server`、`-queue`、`-token-file` が増える。`-server` が無ければ m1 のまま。`bdf printd` と `bdf jobs` が増える。
- `go.mod` に `go.etcd.io/bbolt` が増える（と、それが引く `golang.org/x/sys`）。

## 8. 未決の論点（決めてほしいこと）

決まったもの（2026-10-19 利用者）: サーバーは `bdf printd`、永続化は組み込みの KV（bbolt の 1 ファイル）、変換の入力を受ける、`max-inflight` の既定は 2、サービスとしての登録は持たない、m2 の後は後で決める。

残るもの。どれも暫定で推奨どおり進める。

1. **文書も `spool.db` の中に置く**（推奨。1 ファイルで写しと移しが楽）か、文書だけ隣のディレクトリに置くか。中に置くと、ファイルは消した分だけ縮まず、詰め直しが要る（§3.2）。
2. **複数台・Lambda の継ぎ目**（`Store` のインターフェース、借り、段）を m2 に入れる（推奨。1 台だけでも引き継ぎの照合が同じ道を通るので、余計な手間は小さい）か、1 台に絞って単純に書くか。
3. **パスワードの持ち方**: メモリだけで、再起動の後は `job-password-wait`（推奨）か、サーバーの鍵で封印して記録に置くか。
4. **既定の束縛と認証**: `127.0.0.1:8640` で、loopback かつトークンが無ければ認証なし、でよいか。
5. **`bdf print -spool dir` を作らない**こと（要件からの変更。スプールの持ち主を 1 プロセスにする）。
