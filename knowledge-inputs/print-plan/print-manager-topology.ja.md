# Print Manager のトポロジー（pm-4）— ノードを組み合わせる仮想プリンタ、クラス、インボックス、保管

2026-10-21 の実装案と記録。同日に実装した（`.knowledge/` の `milestone:pm-4`。計画から変えた所と、決めた所はそこの `facts.delivery`）。要件は `requirement:print-manager` と `print-manager.ja.md`（§4「束ね方の足し算」、§10 の pm-4）。pm-1（現場の表示。`milestone:pm-1`）は main に入った（#114）。

2026-10-21 利用者「print manager 上で CUPS の仕組みを使って、トポロジー構造を構成できるようにしたい。いまも、ラウンドロビンで複数プリンタ出力、仮想プリンタなどがあるが、それらを任意に組変えられるようにしたい。たとえば、外から見える仮想プリンタを 3 つつくり、1 つは監査用書き込み専用ストレージにそのまま入れる。もう 1 つは特定の拠点の実プリンタにラウンドロビンで流す。もう 1 つは 2 つに分岐して、それぞれ仮想プリンタに入れて閲覧できるようにする、など。あとは、流れてきたら勝手に保存し、手動でピックアップしたものを手動で後続に流すフィルターというかインボックスみたいな仕組みとか」。

§10 の順では pm-4（束ね方と CUPS）は保管庫（pm-2）と事務所の画面（pm-3）の後だったが、利用者の求めで前に出した。pm-4 の中身（round-robin、プリンタごとのキュー、CUPS の名前の空間）に、分岐・インボックス・保管のノードを足した形。

同日、続けて利用者「入口のノードとしては API エンドポイントも欲しいです。これも任意の数作れると良さそう。API エンドポイントには pdf などを送ることでその後のフローが実行できるとか、bdft ファイルが複数登録でき、データとテンプレート ID を渡すとプリントが行えるとかできても良さそう。UI は全部英語にしたいのと、ユースケースが倉庫のペーパーレスというだけでそれに限定するわけではないです。SDN の様なソフトウェア定義の帳票ネットワークの様なイメージです。必ずしも紙に限定するわけではなく」「GUI でレイアウトできると良いなと思います」。これを受けて、エンドポイント（§1.11）、テンプレート（§1.12）、英語の画面とキャンバス（§1.13）を同じ milestone に入れ、枠組みを「ソフトウェア定義の帳票ネットワーク」に広げた（README と Home の文）。

## 0. 範囲

- **ノードの種類**と**辺**を管理者が組み替えられる（API と `/admin`）。pm-1 の場所（display）とプリンタはそのままノードになる。
- **どのノードも IPP の宛先**（`/ipp/print/<名前>`）。CUPS でプリンタもクラスも宛先であるのと同じ。
- **ジョブ**は 1 つの文書と、通ったノードの足跡（legs）。
- **CUPS のクライアントから一覧**（CUPS-Get-Printers、CUPS-Get-Classes）。

入れないもの：ジョブの移動（CUPS-Move-Job。インボックスの「流す」がその役）、`job-hold-until`（インボックスが明示の仕組み）、規則で振り分けるフィルター（名前や利用者の正規表現で宛先を変える。§5）、保管庫の Object Lock の見出し（pm-2。§1.6）、拠点の printd が取りに来る branch の route（pm-1 のまま 501）。

## 1. 決めたこと

### 1.1 1 つの名前の空間

場所（display）、プリンタ、経路のノード（queue・class・inbox・archive）は 1 つの名前の空間を分け合う。辺は名前だけで相手を指す（種類を書かない）。CUPS の `/printers/<name>` と `/classes/<name>` が同じ名前を取り合うのと同じ。別の種類に同じ名前を付けると 409。

### 1.2 ノードの種類

| 種類 | 何をするか | 設定 |
|---|---|---|
| `display`（場所、pm-1） | その場所の端末に伝票として見せる。紙への逃げ道 | `fallback`、`keep` |
| `printer`（pm-1） | 紙に刷る（direct。branch は未） | `url`、`route` |
| `queue` | `to` のノードすべてに流す（複数なら分岐、1 つなら別名） | `to` |
| `class` | `members` の 1 つに送る。`round-robin`（既定。クラスごとの数えで順に）か `failover`（先頭）。取らなかったプリンタは次のメンバーへ | `members`、`policy` |
| `inbox` | 人が流すまで持つ。流す先は `to` のうちから選ぶ（`to` が空ならどこへでも）。捨てることもできる。流した・捨てた項目は `keep`（14 日）残す | `to`、`keep` |
| `archive` | ジョブの `.bdf` と記録の JSON を、yuruidb の Blob の保存先（bbolt・S3・R2・GCS）にそれぞれ 1 回だけ書く。Manager は変えも消しもしない | `store`、`collection`（`archive-<名前>`） |

「キュー」の名は、送る側が選んで名前を付ける仮想プリンタという意味で付けた。分岐が要らない所でも、クラスやインボックスの前に表示名を持つ宛先を置ける。

### 1.3 どのノードも宛先

利用者の「外から見える仮想プリンタ」は、見せる相手を認証で絞る段（`requirement:print-manager-auth` のトークンごとの範囲）で決めるのが筋なので、ノードの側には「見える・見えない」を持たせなかった。クラスにも保管にも直接刷れる。プリンタに直接刷れば、そのプリンタだけに行く（§4 の「プリンタごとの暗黙のキュー」）。

### 1.4 ジョブと足跡

受けた文書は 1 つの **Job**（番号＝IPP の `job-id`、宛先、名前、利用者、ページ、ハッシュ）になり、文書はジョブの番号で 1 つだけ置く。宛先から辺をたどり、着いたノードごとに **Leg**（ノード、種類、来た所、状態、作った記録の ID、クラスが選んだメンバー、誤り）を足す。

| ノード | 作る記録 | Leg の状態 |
|---|---|---|
| display | **Slip**（pm-1。ID は `<番号>.<場所>`。2 つの場所に分岐すれば 2 枚） | shown |
| inbox | **Item**（`<番号>.<インボックス>`） | held → released / discarded |
| printer・archive・class（メンバーがプリンタ） | **Delivery**（`<番号>.<Leg の番地>`） | pending → done / failed / canceled |
| queue・class（メンバーがプリンタでない） | — | routed |

ジョブの状態は Leg から導く：pending があれば `processing`、held があれば `held`、failed があれば `aborted`、取り消しは `canceled`、それ以外は `completed`。IPP では processing / pending-held（`job-hold-until-specified`）/ aborted / canceled / completed。

**取り消し**（Cancel-Job、`POST /api/jobs/<n>/cancel`）は、開いている伝票をボードから下げ、持っている項目を捨て、送っていない Delivery を止める。完了したジョブでも、伝票がまだ開いていれば取り消せる（pm-1 の振る舞いを保つ）。

### 1.5 Delivery — 送りの記録と再試行

プリンタと保管への送りは要求の中で行わない（送りに数秒〜数分かかる。Lambda は応答の後に止まる）。routing が Delivery を pending で作り、

- **プロセス**（`Options.Background`、サーバーレスでなければ既定）では、受け付けた要求の後に goroutine が送る。
- **tick** が pending をすべて見て、期限（`RetryAt`）の来たものを送る。サーバーレスはこれだけ（EventBridge の 1 分）。

送りは pm-1 の伝票の逃げ道と同じ**借り**（lease、5 分）の下で行い、同時の tick と背景の送りが 2 度送らない。失敗は `Attempts` と `LastError` と `RetryAt`（30 秒から 15 分まで伸びる）に残し、20 回で `failed`（ジョブは aborted）。管理者は `POST /api/deliveries/<id>/retry` で今すぐ、または failed からやり直せる。

クラスのメンバーがプリンタのときは、Leg はクラスのもので、Delivery が `Printer`（選んだメンバー）と `Class` を持つ。失敗したら `Tried` に積み、次の試みはクラスを読み直して「試していない次のメンバー」へ（全部試したら最初から）。round-robin も failover も、失敗の後は同じに次へ行く。違いは最初の 1 台の選び方だけ。

### 1.6 保管（archive）と pm-2 の保管庫

保管のオブジェクトの鍵は `<受けた日>.<番号>.<ハッシュの先頭 12 桁>.bdf` と `.json`（yuruidb の ID に `/` が使えないので、接頭辞の代わりにコレクションを使う）。Blob の Insert は「無いときだけ」なので、応答を失った後の再試行でも 2 度書かない（あれば書いたものとする）。記録の JSON はジョブの属性、ハッシュ、受けた時刻、保管した時刻、ここまでの足跡。

Manager から消す口は無い。S3 Object Lock の見出し（保持の期限とモード）は yuruidb の Blob に無いので（`requirement:print-manager` の `vault-client`）、**バケットの既定の保持**（Object Lock の default retention）で「誰も変えられない」を作る。pm-2 の保管庫は、この archive の上にオブジェクトごとの保持・legal hold・1 日の一覧の鎖・検索項目を足す。

### 1.7 組み替えの確かめ

PUT のときに、名指すノードが全部あること、名前が別の種類と重ならないこと、**自分へ戻る道が無いこと**（queue・inbox の `to` と class の `members` を辺として DFS）を確かめる。インボックスの `to` も辺に数える（人が流す辺でも、戻る道は混乱のもと）。送り先にされているノード・場所・プリンタ（`fallback` も）は消せない（409）。項目を持っているインボックスも消せない。

### 1.8 片付け

- 伝票は場所の `keep`、項目はインボックスの `keep` の後に消える（pm-1 と同じ tick）。
- ジョブは終わった日ごとの記録（`ended-days`、ID が日付なので一覧が古い順）に積み、`JobKeep`（既定 14 日。`MANAGER_JOB_KEEP`）の後、**文書を要る記録が無くなっていれば**（`Live`：伝票・項目・pending の Delivery の数。消えるたびに減る）、文書と Delivery とジョブを消す。
- 宛先ごとの最近のジョブ（`dest-states`、500 件）は Get-Jobs と画面の一覧のため。消えたジョブは読み飛ばす。

### 1.9 IPP と CUPS

- `POST /ipp/print/<名前>`：どのノードも。Get-Printer-Attributes に `printer-type`（クラスは 1）。Get-Jobs は RFC 8011 どおり既定が not-completed（pm-1 は全部だった）。
- `POST /ipp/print`：CUPS-Get-Printers（全ノード。`device-uri` はプリンタの URL か `<種類>://<名前>`、`printer-type` のクラスのビット、クラスには `member-names`・`member-uris`）、CUPS-Get-Classes（クラスだけ）。`bdf printers -url ipp://<host>:<port>/ipp/print` で一覧できる。CUPS-Get-Default は無い（既定の宛先を持たない）。

### 1.10 画面

`/admin` に、トポロジーの表（名前、種類、流れる先、設定、削除）とノードの追加、インボックスごとの保留の項目（流す先のチェックと「流す」「捨てる」、作業者の名前）、宛先ごとのジョブ（状態と足跡、取り消し）、送り待ち（試行、次、最後の誤り、「今すぐ」）。図は描かない（§5）。

### 1.11 エンドポイント（endpoint）— アプリのための入口

`endpoint` はキューと同じに `to` のすべてに流す入口で、**自分の鍵**を持つ。`POST /in/<名前>` に `Authorization: Bearer <鍵>` で：

- 文書（PDF・Office・画像・.bdf …。本文そのもの、`?name=`）→ 変換してジョブに。
- JSON（`Content-Type: application/json`、または本文が `{` で始まる）`{"template": "<id>", "data": {…}, "name": …, "strict": …}` → テンプレートを差し込んでジョブに（§1.12）。

鍵は 32 バイトの乱数（`bdfk_` + base64url）で、作ったとき（PUT の答えの `key`）と `POST /api/nodes/<名前>/key` で回したときに 1 度だけ見せる。保存先には argon2id の PHC 文字列（m=19 MiB、t=2、p=1。OWASP の推奨）だけを置き、API はハッシュも返さない（`requirement:print-manager-auth` の「保存先に取り出せる秘密を置かない」「アプリの秘密は argon2id」に合わせた）。管理者トークンでも `/in/` に送れる。エンドポイントもノードなので IPP の宛先でもあり、`/api/destinations/<名前>/jobs` にも送れる。

### 1.12 テンプレート — データから刷る

テンプレートは `PUT /api/templates/<id>`（本文がファイル、`?name=` がファイル名）で登録する。コンパイル済み（`.bdft`）でも、Excel・Word・PowerPoint・PDF フォーム・pdfme・テキストの素のテンプレートでもよい（`form.Fill` が受けるもの）。登録時に `form.Describe`・`form.Fields` で形式と項目（タグの名前）を取り、記録に持つ（`GET /api/templates/<id>` で見える。データの形の手がかり）。ファイルは文書と同じ Blob の保存先のコレクション `templates`。

エンドポイントは `templates` で差し込めるテンプレートを絞れる（空なら全部）。差し込みは `form.Fill` に `Data` と `Strict`、変換の `FontDirs` を渡し、できた `bdf.Document` をそのまま `submitDoc` でジョブにする（`format` は `form:<形式>`）。テンプレートを名指すエンドポイントは消せない。

### 1.13 英語の画面とキャンバス

画面（Home・端末・管理）と README の文は英語にし、倉庫の語を「display board（登録した端末に見せる）」「documents」に広げた。API と記録の名前（location、slip）は pm-1 のまま。

`/admin` の **キャンバス**（SVG、依存なし）：ノードは種類ごとの色の箱、辺は矢印（クラスのメンバーは破線、場所の逃げ道は点線）。箱をドラッグして動かす（`PUT /api/layout` に位置を保存。`layout` の記録 1 つ、消えたノードの位置は落とす）。右の丸い口から別のノードへドラッグすると辺を足す（PUT /api/nodes でその `to`・`members` に足す。存在・閉路の確かめはサーバー）。辺をクリックで外す（確認つき）。箱をクリックで下のフォームに読み込んで編集。位置の無いノードは、入口からの深さで列、同じ列で行に並べる（autoLayout）。パン・ズームは無い（広い SVG を横スクロール）。

## 2. 構成

```
print/manager/model.go      Node、Job・Leg、Item、Delivery、状態、ID の形
print/manager/nodes.go      ノードの PUT・一覧・削除、名前の空間、参照の確かめ、閉路の検出、Topology
print/manager/route.go      routing（辺をたどる）、クラスの選択、Job の記録、取り消し、宛先ごとの一覧、ジョブの片付け
print/manager/inbox.go      項目、流す、捨てる、インボックスの状態と片付け
print/manager/deliver.go    Delivery、借りの下の送り、クラスの次のメンバー、保管の書き込み、再試行、tick
print/manager/submit.go     受け付け（Job を作り routing へ）、文書はジョブの番号で
print/manager/ipp.go        どのノードも宛先、CUPS-Get-Printers / Classes
print/manager/http.go       /api/nodes、/api/topology、/api/destinations、/api/jobs、/api/inboxes、/api/items、/api/deliveries
print/manager/endpoint.go   エンドポイントの鍵（argon2id）、/in/<名前>、テンプレートの登録・差し込み、キャンバスの位置
print/manager/store.go      新しいコレクション、テンプレートの Blob、保管の Blob の保存先の cache
cmd/printmanager            管理の画面（pages.pw.html、admin.js：キャンバス・テンプレート・鍵の表示）、英語の文、MANAGER_JOB_KEEP、Background
```

## 3. pm-1 から変えた所

- 伝票の ID は `<番号>.<場所>`（1 つのジョブが 2 つの場所に出る）。文書のコレクションは `job-documents`（番号が鍵）。
- `Submit` は宛先（どのノードでも）を取り Job を返す。`/api/locations/<場所>/slips` は残し、その場所の伝票を返す。`/api/destinations/<名前>/jobs` はジョブを返す。
- 伝票の取り消し（`CancelSlip`）とジョブの取り消し（`Cancel`）を分けた。IPP の Cancel-Job はジョブの取り消し。
- `removeSlip` は文書を消さず、ジョブの `Live` を減らす。
- Get-Jobs の既定が not-completed。
- tick の報告に delivered・undelivered・items・jobs。

## 4. 試験

`go test -race ./print/manager/`：TestQueueFanout（分岐、2 枚の伝票、取り消し、片付けで文書も消える）、TestClass（4 件が 2 台に交互、failover で落ちた 1 台を飛ばす、retry）、TestPrinterDirect（プリンタに直接、Background で tick なし）、TestInbox（保留、流す先の制限、流す、全部に流す、捨てる、取り消し、消せない、keep）、TestArchive（bbolt の保管に 2 つのオブジェクト、記録の中身、2 度書かない、悪い URL は 400）、TestNodes（種類・辺・名前・閉路・削除の確かめ、Topology）、TestCupsGetPrinters（`ipp.Client.CupsGetPrinters` で全ノード、クラスの印）、TestEndpoint（鍵は 1 度だけ、鍵なし・別の鍵は 401、文書、`pages.bdft` の登録と JSON での差し込み、無いテンプレート、テンプレートの制限と削除の拒否、鍵の回転、IPP の宛先、レイアウト）、TestKeys（argon2id の PHC）。pm-1 のテストはそのまま通る。

手で：偽プリンタの農場（`tools/fakeprinter`）の mfp-3f・postscript を使い、`/admin` でクラス・分岐・インボックス・保管を組み、API で 8 件送って、インボックスから画面で流した。クラスは mfp-3f と postscript に交互に、分岐は 2 つの倉庫のボードに、保管の bbolt に `.bdf` と `.json`。

## 5. 残り（Claude の暫定。確認待ち）

- **種類の名前**：queue（分岐）・class・inbox・archive。CUPS の語に寄せた。「tee」「group」の方が通りがよければ替える。
- **どのノードも宛先**（§1.3）。見せる相手の制限は認証（トークンごとの宛先）で。
- **クラスの failover は「先頭から」だけ**。bdf printd の `any`（空いている最初の 1 台）は、Manager がプリンタの状態を見張らないので作らなかった。
- **インボックスの流す先は `to` に限る**（`to` が空なら自由）。
- **規則で振り分けるフィルター**（名前・利用者・ページ数で宛先を変える）は作っていない。要るなら `filter` の種類を足す（辺は `rules: [{match, to}]`）。
- **画面に図を描かない**（表だけ）。SVG の図は後で。
- **branch の route** は pm-1 のまま未実装。
- **保管の保持**は S3 のバケットの既定の保持に頼る（§1.6）。pm-2 で yuruidb の Blob に見出しを足す。
- Lambda での背景の送りは無い（tick だけ）。1 分の遅れは EventBridge の刻みの通り。
- **エンドポイントの鍵は 1 つ**（アプリごとに 1 つのエンドポイント、またはエンドポイントを分ける）。複数の鍵・期限・鍵ごとの記録は認証の milestone で。`requesting-user-name` の代わりの `?user=` は自己申告のまま。
- **テンプレートの差し込みの入力の形**は JSON の `data` だけ（CSV・YAML は `form.LoadData` が読めるが口を作っていない）。画像（`Files`）も渡せない。
- **キャンバス**はパン・ズームなし、ノードの追加はフォームから。辺の種類は to と members だけで、インボックスの「人が選ぶ」辺も同じ矢印（点線にはしていない）。
- 画面の文は英語だけ（多言語の仕組みは無い）。
