# API {#the-api}

API は `/api` の下の JSON です。リクエストは次のどれかを添えます。

- マネージャー自身のページからのセッションのクッキー。人なら見られるもの・できることの範囲で、管理者ならすべてです。
- 運用者の管理トークンをベアラーとして（`Authorization: Bearer …`）。スクリプト向けに、すべてです。
- アプリのトークン（[アプリ](entries.md#apps)）。自分のジョブと、そのアプリを挙げたエンドポイントです。
- サブミットトークンの `name:secret` をベアラーとして。ジョブを送ります。
- ティックトークン。`POST /api/tick` だけです。

エラーは `{"error": "…"}` とステータスです。400 はマネージャーが受け付けられないリクエスト（循環を作るノード、ノードが受け取れないペイロードの種類）、401 は資格情報がない、または間違っている、403 は自分のものではない、404 は見つからない、409 は衝突（ほかの種類の名前、終わったジョブ）、413 は文書が大きすぎる、415 は入口が受け取らない形式、502 はプリンタの失敗です。

## 誰か {#me}

| | |
|---|---|
| `GET /api/me` | リクエストの主。種類（人、アプリ、トークン）、id、名前、メールアドレス、グループ、管理者かどうか、役割（`topology`、`audit`） |
| `GET /auth/login`、`POST /auth/logout` | プロバイダでサインインする。サインアウトする |
| `POST /oauth/token` | トークン。アプリのものはクライアントクレデンシャルグラントで、公開クライアントのための人のものは PKCE 付きの認可コードグラントで |
| `GET /oauth/authorize` | サインインした人をコード付きで公開クライアントへ戻す（[フロントエンド](exports.md#frontends)） |

## ジョブ {#jobs}

| | |
|---|---|
| `GET /api/my/jobs?limit=` | 呼び手が見られるジョブ。自分のものと自分のグループのもの（管理者はすべてのジョブ） |
| `GET /api/destinations` | 呼び手が送れるノード |
| `POST /api/destinations/<node>/jobs?name=&user=&meta.<key>=` | ボディ（文書、またはテンプレートノード向けのデータ）をノードに送る。`X-Meta-*` ヘッダも使えます。エンドポイントと同じく `multipart/form-data` のボディも受け取ります。201 とジョブ |
| `POST /api/locations/<display>/slips` | 同じものをディスプレイに。古い名前です |
| `GET /api/destinations/<node>/jobs?limit=` | ノードに送られたジョブ（管理者） |
| `GET /api/jobs/<n>` | ジョブ。ペイロード、メタデータ、レグ |
| `GET /api/jobs/<n>/document?payload=` | その `.bdf` 文書（最後に作られたもの、なければ受け取ったもの）、または id で指定したペイロード |
| `GET /api/jobs/<n>/pages/<p>.png` | 文書のページを描いたもの |
| `POST /api/jobs/<n>/cancel` | 取り消す |
| `POST /api/jobs/<n>/reprint` | `{"to": "<node>"}`。呼び手の新しいジョブとして、送れるノードへもう一度送る。既定はジョブのノード |
| `POST /in/<endpoint>` | エンドポイントの入口。キーかアプリのトークンを添えます（[エンドポイント](entries.md#endpoints)） |

## トポロジー {#topology}

トポロジーはサインインした全員が読めます。変えられるのは役割 `topology` と管理トークンです。

| | |
|---|---|
| `GET /api/topology` | ノードとエッジ。`takes`、`emits`、`accepts`（いまノードに送れるもの）、`broken`（そこから出る、もうつながっていないエッジ。[壊れたエッジ](flow.md#broken)）付き。管理者には設定と数も |
| `GET /api/topology/live` | 呼び手が見られるジョブの動き。各ノードの数と直近 10 分のレグ。ポーリング用の ETag（変化がなければ 304） |
| `GET`、`PUT /api/layout` | ノードを描く位置 |
| `GET /api/nodes`、`GET /api/nodes/<name>` | すべてのノード。1 つのノード |
| `PUT /api/nodes/<name>` | ノードを保存する。`kind` と種類ごとの設定（[ノードの種類](nodes.md)）。新しいエンドポイントの応答には `key` が 1 度だけ入ります |
| `DELETE /api/nodes/<name>` | 削除する。ほかのノードが送り先にしている間はできません |
| `POST /api/nodes/<name>/key` | エンドポイントのキーを回す。新しいキーを 1 度だけ |
| `GET /api/formats` | 入口が受け取れる形式 |
| `GET /api/openapi.json?node=` | 呼び出し元が使えるエンドポイントとエクスポートの OpenAPI 文書（[OpenAPI](#openapi)） |

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

## インボックスと配送 {#deliveries}

| | |
|---|---|
| `GET /api/inboxes/<name>` | インボックス。留まっている項目と、まだ表示している終わった項目 |
| `POST /api/items/<id>/release` | `{"to": ["<node>", …], "worker": "…"}` |
| `POST /api/items/<id>/discard` | `{"worker": "…"}` |
| `GET /api/deliveries` | 待っている配送と失敗した配送 |
| `POST /api/deliveries/<id>/retry` | 今すぐ送り直す |
| `POST /api/tick` | やるべきことを行う（[ティック](inboxes.md#tick)）。ティックトークン、スコープ `tick` のアプリ、管理トークン、または役割 `topology` の人 |

## HTTP (pull) {#exports}

バックエンドのトークン、フロントエンドが得た本人のトークン、またはページのセッションで読みます（[HTTP (pull)](exports.md)）。

| | |
|---|---|
| `GET /out/<name>?owner=&limit=` | 読み手が読めるエントリ。新しい順 |
| `GET /out/<name>/<job>` | 1 つのエントリ |
| `GET /out/<name>/<job>/record` | アーカイブのジョブの記録 |
| `GET /out/<name>/<job>/document` | アーカイブが書いたとおりの文書 |

## アプリと監査 {#apps}

| | |
|---|---|
| `GET /api/apps`、`GET /api/apps/<id>` | アプリ。1 つのアプリ |
| `PUT /api/apps/<id>` | 登録または変更する。`title`、`scopes`、`group`、`disabled`、公開クライアントなら `public` と `redirect_uris`。新しいアプリの応答には `secret` が 1 度だけ入ります（公開クライアントにはありません） |
| `DELETE /api/apps/<id>` | 削除する |
| `POST /api/apps/<id>/secret` | シークレットを回す |
| `GET /api/audit?limit=` | 監査ログ（役割 `audit`） |

## OpenAPI {#openapi}

`GET /api/openapi.json` は、トポロジーが公開しているものを、いまのトポロジーから作った OpenAPI 3.1 の文書です。アプリケーションの開発者はこれからクライアントを生成したり、OpenAPI を読むツールでエンドポイントを試したりします。

- すべてのエンドポイントの `POST /in/<name>`。受け取る形式のメディアタイプが入ります。`accept: data` のエンドポイントなら、ボディのスキーマはデータが届くテンプレートのデータの JSON Schema です（[データ](templates.md#data)）。`?meta.<key>=` パラメータは後ろのノードが見るキーです。フィルタのルール、テンプレートノードの `template_from`（値はテンプレートの ID）、エクスポートの `by` と `fields`、出口の設定の `{{meta.<key>}}` です。エンドポイント自身が設定するキーは載りません。そちらの値が勝つからです。
- すべてのエクスポートの `GET /out/<name>` とそのエントリ。エントリにはエクスポートの fields が入ります。
- 送ったジョブを追うアプリ向けの `GET /api/jobs/<n>`、`…/document`、`POST …/cancel`。
- 入り方。エンドポイントのキー、管理トークン、アプリのトークン（`/oauth/token` のクライアントクレデンシャル）、人のトークン（PKCE 付きの認可コード）です。

載る範囲は読み手のものです。役割 `topology` と管理トークンにはすべてのエンドポイントとエクスポートが、アプリには自分を載せていて送るスコープ（文書なら `print`、データなら `print.template`）のあるエンドポイントと、自分を載せていてスコープ `export` のあるエクスポートが、サインインした人にはエクスポートが載ります。`?node=<name>`（繰り返すか、カンマ区切り）でそのノードだけに絞ります。読み手に何も公開していないノードは見つかりません（404）。`info.version` は残りの部分のハッシュで、公開しているものが変わると変わります。トポロジーのページの OpenAPI で文書全体を、エンドポイントやエクスポートのサイドバーにある同じ名前のリンクでそのノードの分をダウンロードできます。

```sh
curl -H "Authorization: Bearer $TOKEN" "https://<host>/api/openapi.json?node=erp" > erp.openapi.json
```

## IPP {#ipp}

`POST /ipp/print/<node>` はすべてのノードの IPP プリンタで、`POST /ipp/print` はサーバーです。その CUPS-Get-Printers がノードを一覧します（[IPP](entries.md#ipp)）。
