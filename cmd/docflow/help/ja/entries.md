# 文書の入口 {#where-documents-come-in}

ジョブは入口からネットワークに入ります。Jobs ページのフォーム、アプリケーションが POST するエンドポイント、コンピュータが印刷する CUPS サーバーノード、バケットのイベントが届く S3 トリガー、API です。誰がどこへ送れるかは[ノードのスコープ](roles.md#scopes)にあります。

## ページから {#pages}

[Jobs](jobs.md#print) のフォームと、[Topology](topology.md#view) ページのノードのフォームは、自分が送れるノードに、自分として文書を送ります。ジョブは自分のもので、グループはノードのものです。

## IPP {#ipp}

IPP が入ってくるのは CUPS サーバーノードだけです（[CUPS サーバーノード](#cups)）。コンピュータはそれを IPP プリンタとして追加します。ノード自身のポート（`ipp://<host>:<port>/printers/<name>`）でも、マネージャー自身のサーバー（`ipp://<host>:<port>/ipp/print/<name>`。マネージャーが TLS の後ろにあれば `ipps://`）でも構いません。ジョブはそこからノードのエッジをたどります。`bdf printers -url ipp://<host>:<port>/ipp/print` はそれらを一覧します。CUPS クライアントの `lpstat` も同じです（CUPS-Get-Printers）。ほかのノードは IPP プリンタではありません。キュー、プリンタ、ディスプレイにはエッジを通って届きます。リクエストのジョブ属性（`media`、`copies`、`job-name` ……）と `document-metadata` はメタデータ `ipp.<name>` としてジョブとともに流れます。アカウント（`job-account-id`、`job-accounting-user-id`）、`document-metadata`、独自の属性（`x-…`）は、ジョブが届くプリンタへ渡ります（[プリンタへ渡るもの](flow.md#ipp-attributes)）。

IPP で印刷できる人は運用者の `submit_tokens` で決まります。1 つでも設定されていれば、クライアントはトークンの名前とシークレットを Basic 認証で渡します（アプリも id とシークレットを同じように渡します）。設定がなければ、マネージャーに届く誰もが送れます。クライアントが渡す requesting-user-name は「(as …)」として表示されるだけで、確認はしません。

## エンドポイント {#endpoints}

エンドポイントはアプリケーション（業務システム、WMS）向けの REST の入口です。トポロジー上のノードです（[エンドポイント](nodes.md#endpoint)）。最初に保存したときにキーが 1 度だけ表示され、サイドバーの Rotate key で新しいキーを作って 1 度だけ表示します（古いキーは使えなくなります）。マネージャーが保存するのはキーのハッシュだけです。

アプリケーションはキーをベアラートークンとして送ります。

```sh
curl -X POST -H "Authorization: Bearer bdfk_…" \
     --data-binary @statement.pdf \
     "https://<host>/in/erp?name=statement&meta.dept=sales"
```

フォームとして送るなら次のとおりです。

```sh
curl -H "Authorization: Bearer bdfk_…" \
     -F 'meta={"dept": "sales"};type=application/json' \
     -F file=@statement.pdf \
     https://<host>/in/erp
```

- ボディは文書、`accept: data` のエンドポイントならデータです。内容から形式が分からないときは、名前（`?name=`）と `Content-Type` で判別します。
- `?name=` はジョブの名前、`?user=` は送った人の名前で、「(as …)」として表示されます。
- `X-Meta-<Key>` ヘッダと `?meta.<key>=` パラメータはジョブのメタデータです（[メタデータ](flow.md#metadata)）。
- ボディは `multipart/form-data` でも構いません。`Content-Type` で見分けます。文書はパート `file` に入れ、その名前と `Content-Type` が文書のものになります（名前は `?name=` が優先）。メタデータはパート `meta` に JSON オブジェクトで入れ、ヘッダとパラメータのものより優先されます。別の Docflow の [HTTP (push)](nodes.md#http) ノードはどちらの形でも送ります。
- `?name=` がなければ、そのままのボディは `Content-Disposition`（`attachment; filename="statement.pdf"`）から名前を取ります。
- 応答はジョブ（201）で、番号とそこまでのレグが入ります。エンドポイントが受け取らない形式の文書は拒否され（415）、拒否は監査ログに残ります。キーが違えば拒否されます（401）。

エンドポイントが何を受け取るか（メディアタイプ、テンプレートが埋めるデータのスキーマ、後ろのノードが見るメタデータ）は、その OpenAPI 文書に載ります。サイドバーのリンク OpenAPI か、`GET /api/openapi.json?node=<name>` です（[OpenAPI](api.md#openapi)）。

`port` を設定すると、プロセスはそのポートでも待ち受け、そこへの POST はすべてこのエンドポイント宛てになります。パスを設定できないアプリケーション向けです。サーバーレスではポートを使えません。

## CUPS サーバーノード {#cups}

CUPS サーバーノード（[CUPS サーバー](nodes.md#cups)）は、印刷するコンピュータや業務システム向けです。プロセス自身のポートで動く CUPS サーバーで、プリンタはノード自身の 1 台です。コンピュータは、ほかの CUPS プリンタと同じように `ipp://<host>:<port>/printers/<name>`（ノードに証明書があれば `ipps://`）でプリンタを追加し、そこに印刷したものはノードの `to` へ進みます。プリンタならノードは中継になり、シンクなら印刷されたものが見え、フィルタなら条件に合うものだけが残ります。`lpstat -h <host>:<port> -p` で一覧できます。新しい CUPS サーバーには、6310 から、どのノードも使っていない最初のポートが付きます。逆に、このシステムから CUPS サーバーのキューへ印刷するには、エッジの先に [CUPS サーバー（シンク）](nodes.md#cups-out)を置きます。ノードの Usage タブに、URL、macOS・Windows・CUPS での追加のしかた、コマンドが表示されます。

## アプリ {#apps}

アプリは、エンドポイントのキーではなく自前の資格情報で送ったり、エクスポートを読んだりするアプリケーションです。資格情報はアプリのもので、監査ログとジョブにはどのアプリが送ったかが残ります。管理者は Admin → Apps で登録します。id、題名、そのジョブを見られるメンバーのグループ（任意）、そしてスコープです。

| スコープ | |
|---|---|
| `print` | 文書を送る |
| `print.template` | テンプレートを埋める。データを送る |
| `jobs.cancel` | 自分のジョブを取り消す |
| `export` | 自分を載せたエクスポートを読む（[バックエンド](exports.md#backends)） |
| `tick` | `POST /api/tick` を呼ぶ。スケジューラ向け |
| `admin` | 管理者として API を使う。スクリプト向け |

シークレットはアプリを追加したときに 1 度だけ表示され、Rotate secret で新しいものを作ります（古いものは使えなくなります）。Disable はアプリの資格情報と発行済みのトークンを拒否し、Delete は削除します。

アプリは `POST /oauth/token`（OAuth 2.0 のクライアントクレデンシャルグラント）でトークンを取り、ベアラーとして、そのアプリを Apps allowed（エンドポイントのサイドバー）に挙げたエンドポイントへ送ります。または、同じようにそのアプリを挙げた CUPS サーバーノードへ、id とシークレットを Basic 認証で渡して IPP で印刷します。

```sh
curl -u 'erp:<secret>' -d grant_type=client_credentials https://<host>/oauth/token
# {"access_token": "…", "token_type": "Bearer", "expires_in": 600, "scope": "print"}
curl -X POST -H "Authorization: Bearer <access_token>" --data-binary @statement.pdf "https://<host>/in/erp?name=statement"
```

トークンは 10 分間有効です。アプリは自分のジョブを見られ（`GET /api/my/jobs`）、スコープ `jobs.cancel` があれば取り消せます。ノードがグループを与えないとき、アプリのジョブはアプリのグループのものです。

## スクリプト {#scripts}

スクリプトは運用者の管理トークン、またはサブミットトークンの `name:secret` をベアラーとして送ります。`POST /api/destinations/<node>/jobs` に文書をボディとして送り、エンドポイントと同じく `?name=`、`?user=`、メタデータを受け取ります（[ジョブ](api.md#jobs)）。
