# 印刷 実装案（print-m2）

2026-10-19 時点の実装案。要件と決定は `print-requirements.ja.md` §4.5〜§4.10 と `.knowledge/` の `requirement:print`（facts の spool・server・budgets）、`milestone:print-m2` にある。前段の print-m1（`bdf print`、`bdf printers`、`print`・`print/ipp`・`print/raster`）は実装して main に入った（#104）。ここは「何をどの順で作るか」を書く。

## 0. 決まっていること（と、暫定で進めること）

決まっていること。

- スプールはジョブを受け付けて永続化し、順に流し、失敗を扱い、状態を見せる。ジョブの状態は IPP の `job-state` と `job-state-reasons` の語彙をそのまま使う（m4 で IPP サーバーとして見せるとき、写すだけで済む）。
- 置き場所はディレクトリとファイル。制御ファイルは JSON、文書は内容アドレスで 1 回。一時ファイル＋改名。
- 再起動の後に同じジョブを二重に刷らない。送る前にプリンタの `job-id` を記録し、再開時は `Get-Jobs` で確かめてから再送する。
- プリンタ側で `aborted` になったジョブは自動で再送しない。
- 既定は `localhost` に束縛する。外に開けるときは TLS とトークンとキューごとの許可。
- ログに文書の中身・パスワード・トークンを書かない。
- 暗号化した `.bdf` はスプールに暗号化のまま置く。描いた結果（Raster）は既定で残さない。

暫定で推奨どおりに進めること（`requirement:print` の open_questions。異論があれば変える。どれも後から差し替えられる位置に置く）。

- **サーバーは別コマンド `cmd/bdfprintd`**。`cmd/bdf` を常駐の依存（設定ファイル、HTTP サーバー、全変換器）で重くしない。`cmd/bdf` はクライアント（`bdf print -server`、`bdf jobs`）だけを持つ。
- **永続化はファイルだけ**。データベースも組み込みの KV も使わない。件数の規模は 1 サーバーで数万ジョブ。
- **サーバーは `converter` の入力を受ける**。PDF、Office、画像などを受けて、スプールの中で `.bdf` に変換してから印刷する。「何でも紙にする」入口にする。
- **暗号化文書のパスワードはメモリにだけ持つ**。再起動の後は `pending-held`（理由 `job-password-wait`）になり、API でパスワードを渡し直すと流れる。鍵の管理を m2 に持ち込まない。
- **`bdf print -spool dir`（ディスクのスプールに直接置いて帰る）は作らない**。スプールの持ち主は 1 つのプロセス（`bdfprintd`）だけにし、置くのは常に API を通す。要件からの変更。

## 1. 成果物

次のことができること。

```
bdfprintd -config /etc/bdf/bdfprintd.yaml          # 常駐。キューとプリンタは設定ファイルで
bdfprintd -config bdfprintd.yaml -check            # 設定の検査（プリンタに Get-Printer-Attributes まで）

bdf print -server http://127.0.0.1:8640 -queue office -sides long invoice.bdf     # スプールに置く（既定は置いて帰る）
bdf print -server http://127.0.0.1:8640 -queue office -wait 10m report.xlsx       # 変換と印刷の完了まで待つ
bdf jobs -server http://127.0.0.1:8640                       # 終わっていないジョブ
bdf jobs -server http://127.0.0.1:8640 -all -queue office    # 終わったものも
bdf jobs -server … cancel 42 | hold 42 | release 42 | retry 42 | show 42
bdf printers -server http://127.0.0.1:8640                   # サーバーのキューとプリンタと状態

curl -F document=@invoice.bdf -F 'ticket={"sides":"two-sided-long-edge","copies":2}' \
     -H "Authorization: Bearer $TOKEN" https://print.example.local:8640/api/v1/queues/office/jobs
```

`BDF_SERVER` と `BDF_TOKEN`（または `-token-file`）が `-server` と認証の既定になる。ブラウザで `http://127.0.0.1:8640/` を開くと、キュー・プリンタ・ジョブの状態のページ（読むだけ。更新はサーバーからの通知）が見える。

## 2. 構成

```
print/
  print.go         Options.OnCreated（Create-Job の直後、文書を送る前に呼ぶ。エラーならジョブを取り消して返す）、
                   Options.Printer（能力を聞き直さない）、Wait（既存のプリンタのジョブを待つ。m1 の wait を公開）  (S1)
  ticket.go        Ticket に JSON のタグ（キーは IPP の属性の名前: copies、sides、media、media-source …）と、
                   キューの既定に重ねる Merge                                                                 (S1)
print/ipp/ipptest/
  ipptest.go       プリンタの停止と再開、ジョブの履歴の消去（忘れたプリンタ）、文書の途中で接続を切る、
                   which-jobs=all、ジョブを受けて保留する（job-incoming のまま）                               (S1)
print/spool/
  spool.go         Spool: Open（ロック、読み込み、復旧）、Submit、Get、List、Cancel、Hold、Release、Retry、Move、
                   SetPriority、Close。変更の購読（Subscribe）                                                 (S2)
  job.go           Job（制御ファイルの中身）と状態機械（IPP の job-state と理由）、遷移の検査                     (S2)
  store.go         ディレクトリの配置、一時ファイル＋改名、文書の内容アドレス、次の id、ロックファイル、GC          (S2)
  hold.go          job-hold-until の時間帯（indefinite、day-time、evening、night、weekend、second-shift、third-shift）
                   と job-hold-until-time                                                                      (S2)
  account.go       アカウンティング（日ごとの JSON Lines）                                                       (S2)
  sched.go         スケジューラ: キューごとの順（優先度、投入順、保留）、プリンタごとのワーカー、組（any / failover）、
                   同時に描く数の上限                                                                          (S3)
  worker.go        1 ジョブの実行: 変換 → print.Print（OnCreated で記録）→ 結果の写し、再送の間隔                (S3)
  monitor.go       プリンタの見張り: printer-state と理由、能力の更新（printer-config-change-time）              (S3)
  recover.go       再起動の後の照合（二重印刷の防止）                                                           (S3)
  convert.go       受けた入力の変換（converter を呼ぶのは bdfprintd が渡す関数。spool は converter に依らない）     (S5)
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
cmd/bdfprintd/
  main.go          -config、-check、-listen、-spool。converter/all を import して変換を渡す。SIGINT / SIGTERM で
                   受け付けを止め、送信中のジョブを待って閉じる                                                  (S4, S5)
cmd/bdf/
  print.go         -server、-queue、-token-file（送る先がサーバーのとき）                                        (S5)
  jobs.go          bdf jobs                                                                                    (S5)
  printers.go      -server（サーバーのキューとプリンタ）                                                         (S5)
examples/print-server/
  main.go          帳票を穴埋めしてスプールに置くサンプル。server をライブラリとして埋め込む                        (S6)
docs/
  architecture/cli（bdfprintd、bdf jobs、-server）、examples/print-server、design §3.36 に追記              (S6)
test/print/ci.sh   bdfprintd を ippeveprinter につないで投入と完了、再起動をまたぐ                                (S6)
```

`spool` は `print` と `print/ipp` に依り、`server` は `spool` に依り、`client` はどれにも依らない（API の JSON の型だけを持つ）。`converter` に依るのは `cmd/bdfprintd` と `examples` だけ。

## 3. 各部の設計

### 3.1 print パッケージに足すもの（S1）

m1 の `print.Print` は 1 回の呼び出しで聞く・作る・送る・待つ。スプールにはその途中の 2 か所が要る。

- `Options.OnCreated func(*ipp.JobStatus) error`: `Create-Job` の応答を受けた直後、`Send-Document` の前に同じ goroutine で呼ぶ。スプールはここでプリンタの `job-id` と `job-uri` を制御ファイルに書く（書けなければエラーを返し、`Print` はプリンタのジョブを取り消して返す）。`Progress` の `created` は今どおり残す。
- `Options.OnSent func(*ipp.JobStatus) error`: `Send-Document` が成功した直後。スプールは「文書を送り終えた」印を書く（再開時の判断に使う。§3.5）。
- `Options.Printer *ipp.Printer`: 与えられれば `Get-Printer-Attributes` を省く（見張りが持っている能力を使う）。停止の判定は今どおり行う。
- `Wait(ctx, c, jobID, wait, progress) (*ipp.JobStatus, error)`: 既存のプリンタのジョブを終わるまで待つ。m1 の `job.wait` をこの形で公開する。再起動の後と、待つ途中で接続が切れたときに使う。
- `Ticket` に JSON のタグ。キーは IPP の属性の名前（`copies`、`multiple-document-handling` ではなく分かりやすい `collate`、`sides`、`media`、`media-source`、`media-type`、`print-color-mode`、`bilevel`、`print-quality`、`printer-resolution`、`pages`、`print-scaling`、`orientation-requested`、`finishings`、`output-bin`、`job-hold-until`、`job-priority`、`job-name`）。`Merge(defaults Ticket) Ticket` で、ジョブが言わない項目をキューの既定で埋める。

`ipptest` に足すもの: プリンタを止めて再開する（`Stop(reasons)`、`Resume()`）、ジョブの履歴を忘れる（`Forget()`。履歴を持たないプリンタの再現）、`which-jobs=all`、`Send-Document` の本文を n バイトで切る（`FailAfter(n)`）、`Create-Job` の応答を返す前に接続を切る（`DropCreateResponse()`。記録の前に落ちた場合の再現）。

### 3.2 スプールの置き場所（S2）

```
<spool>/
  lock              ロックファイル（pid、ホスト名、開始時刻、bdfprintd の版）
  next-id           次のジョブの id（整数。IPP の job-id と同じ形）
  jobs/00000042.json        制御ファイル（1 ジョブ 1 ファイル）
  docs/sha256-<hex>         受けた文書（.bdf か変換前の入力。内容アドレス）
  docs/sha256-<hex>.bdf     変換した .bdf（変換前の入力の隣。再送で変換し直さない）
  accounting/2026-10-19.jsonl
  printers/<name>.json      プリンタの能力の控え（printer-config-change-time つき。起動を速く、プリンタが落ちていても見せられる）
```

- 書き込みは同じディレクトリの一時ファイル（`.tmp-<rand>`）に書いて `fsync` してから改名する。起動時に残った一時ファイルは消す。
- 文書は受けながら SHA-256 を計算して一時ファイルに書き、改名で `docs/` に入れる。同じ文書は 1 回（1000 枚の請求書の束を 3 回送っても 1 つ）。
- **ロック**: `lock` を `O_CREATE|O_EXCL` で作り、pid とホスト名を書く。既にあれば、同じホストでその pid のプロセスがいなければ（Unix は `kill(pid, 0)`、Windows は `os.FindProcess` が失敗する）古いものとして取り直す。別のホストの名前なら取らずにエラー（NFS の共有を誤って 2 台で使うのを止める）。`golang.org/x/sys` を依存に足さないための形。
- **id**: `next-id` を読み、足して書き直してから使う（落ちても同じ id を 2 度使わない。飛ぶのは構わない）。
- **GC**（1 時間ごとと起動時）: 終わったジョブ（canceled・aborted・completed）の制御ファイルを保持期間（既定 7 日）の後に消す。件数（既定 1 万）と文書の合計バイト（既定 10 GiB）を超えたら古い順に消す。どのジョブからも参照されない文書を消す。

### 3.3 ジョブ（S2）

制御ファイルの中身（JSON）。

```jsonc
{
  "id": 42, "uuid": "urn:uuid:…",            // m4 で job-uuid として見せる
  "queue": "office", "name": "invoice.pdf", "user": "erp",
  "document": { "hash": "sha256-…", "size": 123456, "type": "bdf" /* or the converter's format name */, "bdf": "sha256-….bdf" },
  "ticket": { "sides": "two-sided-long-edge", "copies": 2 },   // キューの既定を重ねる前の、受けたままのもの
  "options": { "format": "auto", "view": "", "skip": ["background"], "strict": false, "raster-if-lossy": false },
  "state": "processing", "reasons": ["job-printing"], "message": "",
  "priority": 50, "hold-until": "", "hold-until-time": null,
  "attempts": 1, "next-attempt": null,
  "printer": { "name": "mfp1", "url": "ipps://…/ipp/print", "job-id": 12, "job-uri": "…", "created": "…", "sent": true },
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

`job-hold-until` は CUPS と同じ時間帯（day-time 6〜18 時、evening 18 時〜、night 18〜6 時、second-shift 16〜24 時、third-shift 0〜8 時、weekend 土日）と `indefinite`、`job-hold-until-time`（時刻）。時間帯はサーバーの地方時。

### 3.4 スケジューラ（S3）

- **キュー**: 設定のキューごとに、出られるジョブ（pending で、保留でなく、`next-attempt` を過ぎた）を `priority` の降順、`submitted` の昇順で並べる。キューは一時停止できる（受け付けは続け、流さない）。受け付け自体を止める（`accepting: false`）こともできる。
- **プリンタのワーカー**: 設定のプリンタごとに 1 本の goroutine。そのプリンタで「まだ終わっていないジョブ」が `max-inflight`（既定 1）未満で、見張りが止まっていない・届くと言っているときに、そのプリンタを使えるキューから次のジョブを取る。プリンタの中に溜めずに BDF の側で並べておくのは、取り消し・並べ替え・別のプリンタへの移しを効かせるため。
- **組（クラス）**: キューに複数のプリンタを書けば組になる。`policy: any`（空いているプリンタから。ワーカーが取り合う）か `failover`（並べた順。前のプリンタが止まっている・届かないときだけ次）。
- **描く数の上限**: 変換と PWG Raster の描画は CPU を使う。全体で同時に `renders`（既定は CPU 数）まで。PDF を送るジョブは書き出しが軽いので同じ枠で数える。
- **移す**: `Move` で別のキューへ。pending か pending-held のときだけ。

### 3.5 1 ジョブの実行と、失敗と、二重印刷の防止（S3）

ワーカーがジョブを取ると、

1. `processing`（`job-transforming`）。文書が `.bdf` でなければ変換する（§3.7）。変換に失敗したら `aborted`（`document-format-error`）。
2. 暗号化した `.bdf` なら、メモリのパスワードで開く。無ければ `pending-held`（`job-password-wait`）に戻す。
3. キューの既定にジョブのチケットを重ね、`print.Print` を呼ぶ（`Printer` は見張りの控え、`Wait` は負＝文脈が続く限り、`OnCreated` と `OnSent` で記録）。
4. `Print` の結果を写す: 終わった → `completed`（警告や受けない属性があれば `job-completed-with-warnings`）、`ErrAborted` → `aborted`、プリンタ側の取り消し → `canceled`、`ErrUnsupported`（strict）と `client-error-*` → `aborted`。
5. 接続の失敗・`server-error-busy`・`server-error-temporary-error`・`server-error-not-accepting-jobs`・`ErrPrinterStopped` で、まだ `Create-Job` が通っていなければ `pending` に戻して `next-attempt` を置く（10 秒、30 秒、1 分、2 分、5 分、以後 10 分ごと）。止まっているならプリンタを「止まった」として見張りに渡し、キューの他のプリンタ（組）に回す。再送の上限は時間で持つ（既定 24 時間で `aborted`、`bdf-retry-exhausted`）。
6. `Create-Job` の後で接続が切れたら、ジョブは `processing` のまま、見張りが届くと言うまで待って `print.Wait` で続ける（作り直さない）。

**再起動の後の照合**（`recover.go`）。`processing` か `processing-stopped` のジョブそれぞれについて、

| 記録 | プリンタに聞いて | すること |
|---|---|---|
| `job-id` あり、`sent` あり | `Get-Job-Attributes` で見つかる | 終わっていれば写す、終わっていなければ `print.Wait` で続ける |
| 〃 | 見つからない（履歴を持たない・消した） | `completed`（警告「プリンタが履歴を持たない」）。送り終えているので再送しない |
| `job-id` あり、`sent` なし | 見つかって pending（文書待ち） | プリンタのジョブを取り消して、新しい試行で送り直す |
| 〃 | 見つからない、または終わっている | 新しい試行で送り直す |
| `job-id` なし（`Create-Job` の応答を記録する前に落ちた） | `Get-Jobs`（which-jobs=all、my-jobs）で、同じ `job-name` と利用者で試行の開始より後のもの | 1 つなら上の表に入れて扱う。無ければ送り直す。2 つ以上なら `pending-held`（`bdf-duplicate-check`）にして人に任せる |

`job-name` はプリンタによっては切られるので前方一致で比べる。`requesting-user-name` は常にスプールの利用者名（設定の `user`、既定 `bdf`）にし、`job-name` にジョブの名前を入れる（プリンタのパネルに見える名前）。

### 3.6 プリンタの見張り（S3）

プリンタごとに 1 本。`Get-Printer-Attributes` の状態の属性（`printer-state`、`printer-state-reasons`、`printer-state-message`、`printer-is-accepting-jobs`、`printer-config-change-time`）を、何もないとき 60 秒、ジョブを待たせているときと止まっているとき 10 秒ごとに聞く。届かないときは間隔を延ばす（最大 5 分）。`printer-config-change-time` が変わったら全部の能力を聞き直して `printers/<name>.json` に控える。`Create-Printer-Subscriptions` と `Get-Notifications` は m2 では使わない（機種の差が大きい）。

### 3.7 受けた入力の変換（S5）

`spool` は変換器に依らず、`spool.Options.Convert func(ctx, r io.ReaderAt, size int64, name string) (*bdf.Document, []string, error)` を受ける。`cmd/bdfprintd` が `converter/all` で作って渡す。変換はワーカーの中（`job-transforming`）で行い、HTTP の受け付けは速く返す。HTML と Markdown の外部の画像は取りに行かない（`remote=false`。サーバーが内側の網に届く入口にならないため）。Excel は `views=print`。変換の予算は変換器のまま。結果の `.bdf` を `docs/<hash>.bdf` に控え、再送で変換し直さない。

### 3.8 HTTP API（S4）

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
| `GET /discover` | DNS-SD でネットワークのプリンタを探す（管理者だけ） |
| `GET /events` | Server-Sent Events。ジョブ・キュー・プリンタの変化（`event: job`、`data: {ジョブ}`） |

`/healthz`（スプールが開いていれば 200）、`/metrics`、`/`（状態のページ）は `/api/v1` の外。

**認証と許可**:

- 設定の `tokens` に名前・トークン（ファイルか環境変数から）・使えるキュー・役割（`submit`、`manage`（自分以外のジョブの取り消しや保留）、`admin`（キューと発見））を書く。`Authorization: Bearer …`。比べるのは定数時間。
- 束縛が loopback だけで `tokens` が空なら、認証なし（手元の 1 台向け）。loopback 以外に束縛するときは、TLS（`tls.cert`・`tls.key`）とトークンが無ければ起動しない（`insecure: true` で明示したときだけ許す）。
- ジョブの利用者はトークンの名前。トークンに `act-as: true` があるときだけ `X-BDF-User` を利用者として受ける（ERP が利用者に代わって出すとき）。
- 自分のジョブは `submit` で取り消せる。他人のジョブは `manage`。

**限界**: 文書の大きさ（`max-document`、既定 256 MiB。超えたら 413）、同時の受け付け（既定 8）、スプールの合計容量（超えたら 507）、`limit` は最大 1000、SSE の同時接続（既定 64）。ヘッダーと本文の読み込みに時間の上限。

### 3.9 設定（S4）

```yaml
listen: 127.0.0.1:8640
spool: /var/spool/bdf
user: bdf                         # プリンタへの requesting-user-name
tls: { cert: /etc/bdf/tls.crt, key: /etc/bdf/tls.key }   # loopback 以外では必須
tokens:
  - { name: erp, file: /etc/bdf/erp.token, queues: [office, warehouse], roles: [submit], act-as: true }
  - { name: admin, env: BDF_ADMIN_TOKEN, roles: [submit, manage, admin] }
printers:
  mfp1: { url: ipps://mfp1.example.local/ipp/print }
  mfp2: { url: ipp://192.0.2.20/ipp/print, username: print, password-env: MFP2_PASSWORD }
queues:
  office:
    printers: [mfp1, mfp2]
    policy: failover                # any | failover
    defaults: { sides: two-sided-long-edge, print-color-mode: monochrome }
    format: auto                    # auto | pdf | pwg-raster
    raster-if-lossy: true
    max-inflight: 1
retention: { done: 168h, max-jobs: 10000, max-bytes: 10GiB, accounting: 400d }
limits: { max-document: 256MiB, max-impressions: 10000, renders: 0, retry-for: 24h }
fonts: { dirs: [/usr/share/fonts/bdf], no-system: false }
trust-store: /var/lib/bdf/printers.json   # ipps の証明書の指紋
```

`bdfprintd -check` は設定を読み、プリンタに `Get-Printer-Attributes` を送って届くかと能力を表にして終わる。設定の変更は再起動で反映する（m2 では再読み込みしない）。

### 3.10 クライアントと CLI（S5）

- `print/server/client`: `Submit(ctx, queue, name, r io.Reader, ticket, opts) (*Job, error)`（multipart を流しながら送る）、`Job`、`Jobs`、`Cancel` など、`Watch(ctx, id)`（SSE。切れたらポーリングに落ちる）。JSON の型（`Job`、`Queue`、`Printer`）はこのパッケージに置き、`server` もそれを使う。
- `bdf print -server URL -queue Q`: チケットのフラグは m1 と同じ。既定は投入して帰る（ジョブの id を出す）。`-wait` で終わるまで待ち、終了コードは m1 と同じ意味（1 取り消し・中止、3 サーバーに届かない、4 拒まれた、5 キューが受け付けない、6 時間切れ）。`-server` のときは文書を変換せずにそのまま送る（変換はサーバー）。
- `bdf jobs`: 一覧（id、状態、理由、キュー、名前、利用者、時刻）と `show`・`cancel`・`hold`・`release`・`retry`。`-json`。
- `bdf printers -server`: サーバーのキューとプリンタと状態。

### 3.11 観測（S4）

- ログは `log/slog` の JSON（`-log text` で人向け）。ジョブの受け付け・状態の変化・再送・照合・プリンタの状態の変化。文書の中身、パスワード、トークン、URL の利用者情報は書かない。
- `/metrics`（Prometheus のテキスト形式を自前で書く。依存を足さない）: `bdf_print_jobs_total{queue,state}`、`bdf_print_jobs{queue,state}`（今の数）、`bdf_print_printer_up{printer}`、`bdf_print_printer_state{printer}`、`bdf_print_sheets_total{queue}`、`bdf_print_sent_bytes_total{printer}`、`bdf_print_retries_total{printer}`、`bdf_print_spool_bytes`。
- 状態のページ: キューごとの順番とプリンタの状態とジョブ。`/events` で更新。読むだけ（操作は API と CLI）。
- アカウンティング: 終わったジョブごとに 1 行（時刻、id、キュー、プリンタ、利用者、名前、状態、枚数、印象数、用紙、両面、色、形式、バイト）。集計は外。

### 3.12 常駐のさせ方（S6、文書）

m2 の `bdfprintd` は前面で動くプロセス。systemd の unit と macOS の launchd の plist の例を docs に書く。Windows はサービスとして登録する仕組み（`golang.org/x/sys/windows/svc`）を m2 には入れず、タスク スケジューラか NSSM で動かす例を書く（論点 §8）。

### 3.13 サンプル（S6）

`examples/print-server`: `print/server` をライブラリとして埋め込んだ小さなサーバー。請求書のテンプレート（`converter/xlsx/testdata/form-pages.xlsx`）とデータの入力欄のページを出し、送られたデータで `form.Fill` して `.bdf` を作り、スプールに置き、ジョブの状態を `/events` で見せる。プリンタが無い手元でも、`ipptest` の偽のプリンタを同じプロセスで立てる `-fake-printer` で動く。`docs/examples/print-server(.ja).md` と構成のサンプルの一覧に足す。

## 4. テスト

- **単体**（`go test ./print/spool/`）: 状態機械の遷移の表（許される遷移と許されない遷移）、制御ファイルの往復、一時ファイルと改名（途中で落ちたファイルの掃除）、ロック（生きている pid、死んだ pid、別のホスト）、id が戻らない、GC（期間・件数・バイト・参照されない文書）、保留の時間帯（偽の時計で）、優先度と投入順、再送の間隔（偽の時計）。
- **偽のプリンタで**（`ipptest`）: 2 台の組（any と failover、1 台を止める・再開する）、プリンタ側の中止で再送しない、取り消し（pending・processing・プリンタに送った後）、保留と解除、別のキューへの移し、変換して印刷（PDF と TIFF の入力）、暗号化した文書とパスワード待ち。
- **二重印刷の防止**: 落ちた時点ごとに、スプールを閉じて開き直し、プリンタの受けたジョブの数を数える。`Create-Job` の前、`Create-Job` の応答の前（`DropCreateResponse`）、記録の後で送る前、送る途中（`FailAfter`）、送った後で待つ間、履歴を持たないプリンタ（`Forget`）。どの場合もプリンタが刷るのは 1 回（取り消した不完全なジョブは数えない）。
- **HTTP**（`httptest`）: 認証と許可（トークンなし、別のキュー、manage でない他人のジョブ、act-as）、loopback 以外で TLS が無いと起動しない、大きさの上限（413）、容量の超過（507）、SSE、エラーの形。
- **クライアントと CLI**: `bdf print -server`（帰る・待つ・終了コード）、`bdf jobs` の各操作。
- **競合**: `go test -race`。投入・取り消し・状態の変化・GC を同時に。
- **CI**: `test/print/ci.sh` に、`bdfprintd` を `ippeveprinter` につないで起動し、`bdf print -server` で 2 つ投入し、1 つを取り消し、`bdfprintd` を止めて再起動し、残りが完了して `ippeveprinter` が受けた文書が 1 つ（二重でない）ことを確かめる段を足す。

## 5. 順番

1. **S1** `print` の継ぎ目（`OnCreated`、`OnSent`、`Printer`、`Wait`、`Ticket` の JSON と `Merge`）と `ipptest` の足し算。小さく、ほかの全部の前提。
2. **S2** `print/spool` の保存（`store`、`job`、`hold`、`account`）と `Spool` の API の受け付け・一覧・操作（まだ流さない）。
3. **S3** スケジューラ、ワーカー、見張り、照合。ここで「置けば刷られる」になる。二重印刷の防止のテストはここ。
4. **S4** `print/server`（設定、API、認証、SSE、/metrics、状態のページ）と `cmd/bdfprintd`。
5. **S5** `client`、`bdf print -server`、`bdf jobs`、`bdf printers -server`、受けた入力の変換、暗号化した文書。
6. **S6** `examples/print-server`、docs、CI、`.knowledge` の記録。

S2 と S4 の HTTP の骨組みは並行できる。S3 が山。

## 6. 後で変えられる位置に置くもの

- 永続化をファイルから KV に替えるなら `store.go` だけ。`Spool` は `store` のインターフェース（制御ファイルの読み書き、文書、id、ロック）にだけ依る。
- パスワードの持ち方（メモリ → サーバーの鍵で封印）は `spool.Options.Secrets`（インターフェース）にしておく。
- サーバーの形（`cmd/bdfprintd` → `bdf printd`）は `cmd` の置き場所だけ。
- m4 の IPP サーバーは、`server` に `/ipp/print/<queue>` の IPP の受け口を足し、`Spool.Submit` に渡す。ジョブの状態は IPP の語彙のままなので写すだけ。`spool` の id は IPP の `job-id` の形（正の整数）、`uuid` は `job-uuid` になる。
- 拠点の代理（`remote-proxy`）は、代理側の `bdfprintd` が上流の API（`GET /jobs?state=pending&queue=…` と文書の取得）で取りに行く形が `client` の上に作れる。IPP INFRA に合わせるかは m4 で決める。

## 7. m1 から変わるもの

- `print.Options` に `OnCreated`、`OnSent`、`Printer` が増え、`print.Wait` が公開される。m1 の `bdf print`（サーバーなし）の動きは変わらない。
- `Ticket` が JSON になる（タグを足すだけ。Go の利用者には影響なし）。
- `bdf print` に `-server`、`-queue`、`-token-file` が増える。`-server` が無ければ m1 のまま。

## 8. 未決の論点（決めてほしいこと）

1. **サーバーの形**: 別コマンド `cmd/bdfprintd`（推奨、暫定で進める）か `bdf printd` か。
2. **永続化**: ファイルだけ（推奨、暫定で進める）か、組み込みの KV か。
3. **変換の入力を受けるか**: 受ける（推奨、暫定で進める）。
4. **パスワードの持ち方**: メモリだけで、再起動の後は `job-password-wait`（推奨、暫定で進める）か、サーバーの鍵で封印して制御ファイルに置くか。
5. **既定の束縛と認証**: `127.0.0.1:8640` で、loopback かつトークンが無ければ認証なし、でよいか。
6. **`bdf print -spool dir` を作らない**こと（要件からの変更。スプールの持ち主を 1 プロセスにする）。
7. **プリンタの中に溜めるジョブの数**（`max-inflight`）の既定を 1 にするか。1 は取り消しと並べ替えが効く代わりに、ジョブの合間にプリンタが待つ。業務機の連続印刷の速さを優先するなら 2。
8. **Windows のサービス化**: m2 では前面のプロセスとタスク スケジューラ・NSSM の例だけにし、`x/sys/windows/svc` による登録は後にするか。
9. **m4（IPP サーバー）の前倒し**: m2 の次に m4 を置くか（`requirement:print` の `ipp-server-priority`）。m2 はその受け口を足すだけで済む形にしてある（§6）。
