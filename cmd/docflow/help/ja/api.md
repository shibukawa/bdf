# API {#the-api}

API は `/api` の下の JSON で、マネージャー自身のサーバーにあります。アプリケーションやスクリプトが外から呼ぶものはすべてここにあり、サーバーもポートも 1 つです（例外は CUPS サーバーノードの IPP だけで、そのノードのポートにあります）。外から届くノードは、種類と名前で `/api/<kind>/<name>` にあります。アップローダーは `/api/upload/<name>` で文書を受け取り、インボックスは `/api/inbox/<name>` で扱い、ダウンローダーは `/api/download/<name>` で読み、S3 トリガーは `/api/s3-trigger/<name>` でイベントの参照を受け取ります（[ノードの API](#nodes)）。リクエストは次のどれかを添えます。

- マネージャー自身のページからのセッションのクッキー。人なら見られるもの・できることの範囲で、管理者ならすべてです。
- 運用者の管理トークンをベアラーとして（`Authorization: Bearer …`）。スクリプト向けに、すべてです。
- アプリのトークン（[アプリ](entries.md#apps)）。自分のジョブと、そのアプリを挙げたアップローダー、インボックス、ダウンローダーです。
- フロントエンドがその人のために受け取った、人のトークン。その人が使えるダウンローダーとインボックス、そして `/api/me` です。
- アップローダーか S3 トリガーのキー。そのノードへ送ります。
- サブミットトークンの `name:secret` をベアラーとして。ジョブを送ります。
- ティックトークン。`POST /api/tick` だけです。

エラーは `{"error": "…"}` とステータスです。400 はマネージャーが受け付けられないリクエスト（循環を作るノード、ノードが受け取れないペイロードの種類）、401 は資格情報がない、または間違っている、403 は自分のものではない、404 は見つからない、409 は衝突（ほかの種類の名前、終わったジョブ）、413 は文書が大きすぎる、415 は入口が受け取らない形式、422 は入口の後ろのアサートが拒否したメタデータ（`error` のほかに `job` と `failures` が入ります。[アップローダー](entries.md#uploaders)）、502 はプリンタの失敗です。

## 誰か {#me}

| | |
|---|---|
| `GET /api/me` | リクエストの主。種類（人、アプリ、トークン）、id、名前、メールアドレス、グループ、管理者かどうか、役割（`topology`、`audit`） |
| `GET /auth/login`、`POST /auth/logout` | プロバイダでサインインする。サインアウトする |
| `POST /oauth/token` | トークン。アプリのものはクライアントクレデンシャルグラントで、公開クライアントのための人のものは PKCE 付きの認可コードグラントで |
| `GET /oauth/authorize` | サインインした人をコード付きで公開クライアントへ戻す（[フロントエンド](downloads.md#frontends)） |

## ジョブ {#jobs}

| | |
|---|---|
| `GET /api/my/jobs?limit=` | 呼び手が見られるジョブ。自分のものと自分のグループのもの（管理者はすべてのジョブ） |
| `GET /api/destinations` | 呼び手が送れるノード |
| `POST /api/destinations/<node>/jobs?name=&user=&meta.<key>=` | ボディ（文書、またはテンプレートノード向けのデータ）をノードに送る。`X-Meta-*` ヘッダも使えます。アップローダーと同じく `multipart/form-data` のボディも受け取ります。201 とジョブ |
| `POST /api/locations/<display>/slips` | 同じものをディスプレイに。古い名前です |
| `GET /api/destinations/<node>/jobs?limit=` | ノードに送られたジョブ（管理者） |
| `GET /api/jobs/<n>` | ジョブ。ペイロード、メタデータ、レグ |
| `GET /api/jobs/<n>/document?payload=` | その `.bdf` 文書（最後に作られたもの、なければ受け取ったもの）、または id で指定したペイロード |
| `GET /api/jobs/<n>/pages/<p>.png` | 文書のページを描いたもの |
| `POST /api/jobs/<n>/cancel` | 取り消す |
| `POST /api/jobs/<n>/reprint` | `{"to": "<node>"}`。呼び手の新しいジョブとして、送れるノードへもう一度送る。既定はジョブのノード |

## トポロジー {#topology}

トポロジーはサインインした全員が読めます。変えられるのは役割 `topology` と管理トークンです。

| | |
|---|---|
| `GET /api/topology` | ノードとエッジ。`takes`、`emits`、`accepts`（いまノードに送れるもの）、`broken`（そこから出る、もうつながっていないエッジ。[壊れたエッジ](flow.md#broken)）付き。管理者には設定と数も |
| `GET /api/topology/live` | 呼び手が見られるジョブの動き。各ノードの数と直近 10 分のレグ。ポーリング用の ETag（変化がなければ 304） |
| `GET`、`PUT /api/layout` | ノードを描く位置 |
| `GET /api/nodes`、`GET /api/nodes/<name>` | すべてのノード。1 つのノード |
| `PUT /api/nodes/<name>` | ノードを保存する。`kind` と種類ごとの設定（[ノードの種類](nodes.md)）。新しいアップローダー（または S3 トリガー）の応答には `key` が 1 度だけ入ります |
| `DELETE /api/nodes/<name>` | 削除する。ほかのノードが送り先にしている間はできません |
| `POST /api/nodes/<name>/key` | アップローダーか S3 トリガーのキーを回す。新しいキーを 1 度だけ |
| `POST /api/kinds/<kind>/discover` | 保存する前に、ノードの設定（ボディ。`PUT` と同じ）の向こう側に何があるかを尋ねる。CUPS サーバー（シンク）ならキューの一覧 `[{"name", "path", "info", "state", "accepting", "class"}]`。応答がなければ 502 |
| `POST /api/nodes/<name>/rename` | 名前を変える。`{"name": "<new>"}`。エッジとその名前で保存しているものが付いてきます。そこで処理中の仕事がある間は拒否します（409。[名前](topology.md#names)） |
| `GET /api/formats` | 入口が受け取れる形式 |
| `GET /api/openapi.json?node=` | 呼び手が使えるアップローダー、インボックス、ダウンローダーの OpenAPI 文書（[OpenAPI](#openapi)） |

## テンプレート {#templates}

| | |
|---|---|
| `GET /api/nodes/<node>/templates` | テンプレートノードのテンプレート |
| `GET /api/nodes/<node>/templates/<id>` | 1 つ。形式、引く項目、データのサンプル |
| `GET /api/nodes/<node>/templates/<id>/file` | そのファイル |
| `PUT /api/nodes/<node>/templates/<id>?name=` | アップロードする。ファイルをボディに、ファイル名を添えて |
| `DELETE /api/nodes/<node>/templates/<id>` | 削除する |

## ディスプレイ・シンク・プリンタ {#devices}

| | |
|---|---|
| `GET /api/nodes/<display>/board` | ディスプレイの伝票 |
| `GET /api/nodes/<display>/devices` | そのデバイス |
| `POST /api/nodes/<display>/enrollments` | デバイスを登録する 1 度きりのコード。QR コードと URL 付き |
| `DELETE /api/devices/<id>` | デバイスを外す |
| `GET /api/nodes/<sink>/latest` | シンクが保持する最新の文書 |
| `GET /api/nodes/<printer>/preview`、`…/preview/<job>/pages/<p>.png` | プリンタの偽プリンタが受け取ったものと、そのページ |

## ノードの API {#nodes}

アプリケーションが呼ぶノードを、種類と名前で並べます。

| | |
|---|---|
| `POST /api/upload/<name>?name=&user=&meta.<key>=` | アップローダーに文書かデータを送る。キーかアプリのトークンを添えます（[アップローダー](entries.md#uploaders)） |
| `POST /api/s3-trigger/<name>` | S3 のイベントの参照。ノードのキーを添えます（[S3 トリガー](nodes.md#s3-trigger)） |
| `GET /api/inbox` | 呼び手が扱うインボックス |
| `GET /api/inbox/<name>` | インボックスに留まっているものと、送り出せるノード |
| `GET /api/inbox/<name>/<job>`、`…/thumbnail`、`…/pages/<p>.png`、`…/document` | 項目、1 ページ目の縮小画像、ページ、文書 |
| `POST /api/inbox/<name>/<job>/release`、`…/discard` | `{"to": ["<node>", …], "worker": "…"}`、`{"worker": "…"}`。決定です（[インボックスの API](inboxes.md#api)） |
| `GET /api/download/<name>?owner=&limit=` | 読み手が読めるエントリ。新しい順（[ダウンローダー](downloads.md)） |
| `GET /api/download/<name>/<job>`、`…/record`、`…/document` | 1 つのエントリ、アーカイブのジョブの記録（インボックスのエントリならジョブ）、文書 |

インボックスとダウンローダーは、公開クライアントのオリジンのページにも応えます（CORS）。人のトークンで呼ぶフロントエンドのためです。

## 配送 {#deliveries}

| | |
|---|---|
| `GET /api/deliveries` | 待っている配送と失敗した配送 |
| `POST /api/deliveries/<id>/retry` | 今すぐ送り直す |
| `POST /api/tick` | やるべきことを行う（[ティック](inboxes.md#tick)）。ティックトークン、スコープ `tick` のアプリ、管理トークン、または役割 `topology` の人 |

## アプリと監査 {#apps}

| | |
|---|---|
| `GET /api/apps`、`GET /api/apps/<id>` | アプリ。1 つのアプリ |
| `PUT /api/apps/<id>` | 登録または変更する。`title`、`scopes`、`group`、`disabled`、公開クライアントなら `public` と `redirect_uris`。新しいアプリの応答には `secret` が 1 度だけ入ります（公開クライアントにはありません） |
| `DELETE /api/apps/<id>` | 削除する |
| `POST /api/apps/<id>/secret` | シークレットを回す |
| `GET /api/audit?limit=` | 監査ログ（役割 `audit`） |

## OpenAPI {#openapi}

`GET /api/openapi.json` は、トポロジーが公開しているものを、いまのトポロジーから作った OpenAPI 3.1 の文書です。アプリケーションの開発者はこれからクライアントを生成したり、OpenAPI を読むツールで API を試したりします。

- すべてのアップローダーの `POST /api/upload/<name>`（と S3 トリガーの `POST /api/s3-trigger/<name>`）。受け取る形式のメディアタイプが入ります。`accept: data` のアップローダーなら、ボディのスキーマはデータが届くテンプレートのデータの JSON Schema です（[データ](templates.md#data)）。`?meta.<key>=` パラメータは後ろのノードが見るキーです。フィルタのルール、テンプレートノードの `template_from`（値はテンプレートの ID）、ダウンローダーの `by` と `fields`、出口の設定の `{{meta.<key>}}` です。アップローダー自身が設定するキーは載りません。そちらの値が勝つからです。
- すべてのインボックスの `/api/inbox/<name>`。留まっているもの、項目、そのサムネイル・ページ・文書、そして送り出しです。送り出しの `to` は、インボックスが送り出せるノードのどれかです。`GET /api/inbox` も載ります。
- すべてのダウンローダーの `GET /api/download/<name>` とそのエントリ。エントリにはダウンローダーの fields が入ります。
- 送ったジョブを追うアプリ向けの `GET /api/jobs/<n>`、`…/document`、`POST …/cancel`。
- 入り方。入口のキー、管理トークン、アプリのトークン（`/oauth/token` のクライアントクレデンシャル）、人のトークン（PKCE 付きの認可コード）です。

載る範囲は読み手のものです。役割 `topology` と管理トークンにはすべてのアップローダー、インボックス、ダウンローダーが載ります。アプリには、自分を挙げていて送るスコープ（文書なら `print`、データなら `print.template`）のあるアップローダー、自分を挙げていてスコープ `inbox` のあるインボックス、自分を挙げていてスコープ `download` のあるダウンローダーが、サインインした人には、ダウンローダーと自分が扱うインボックスが載ります。`?node=<name>`（繰り返すか、カンマ区切り）でそのノードだけに絞ります。読み手に何も公開していないノードは見つかりません（404）。`info.version` は残りの部分のハッシュで、公開しているものが変わると変わります。トポロジーのページの OpenAPI で文書全体を、ノードのサイドバーにある同じ名前のリンクでそのノードの分をダウンロードできます。

```sh
curl -H "Authorization: Bearer $TOKEN" "https://<host>/api/openapi.json?node=erp" > erp.openapi.json
```

## IPP {#ipp}

`POST /ipp/print/<node>` は CUPS サーバーノードの IPP プリンタです（ほかのノードは IPP プリンタではありません）。`POST /ipp/print` はサーバーで、その CUPS-Get-Printers が CUPS サーバーノードを一覧します（[IPP](entries.md#ipp)）。
