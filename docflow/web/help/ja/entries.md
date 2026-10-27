# 文書の入口 {#where-documents-come-in}

ジョブは入口からネットワークに入ります。ジョブのページのフォーム、アプリケーションが POST するアップローダー、コンピュータが印刷する CUPS サーバーノード、バケットのイベントが届く S3 トリガー、API です。誰がどこへ送れるかは[ノードのスコープ](roles.md#scopes)にあります。

## ページから {#pages}

[ジョブのページ](jobs.md#print)のフォームと、[トポロジーのページ](topology.md#view)のノードのフォームは、自分が送れるソースに、自分として文書を送ります。ソースは文書が入るノード、つまりアップローダー、CUPS サーバー、S3 トリガーです。ソースが受け取るものと付けるものが、アプリケーションと同じく人にも効くようにするためです。ジョブは自分のもので、グループはノードのものです。役割 `topology` を持つ人は、トポロジーの一部を試すために、ジョブのページから自分が送れるどのノードにも送れます。

## IPP {#ipp}

IPP が入ってくるのは CUPS サーバーノードだけです（[CUPS サーバーノード](#cups)）。コンピュータはそれを IPP プリンタとして追加します。ノード自身のポート（`ipp://<host>:<port>/printers/<name>`）でも、マネージャー自身のサーバー（`ipp://<host>:<port>/ipp/print/<name>`。マネージャーが TLS の後ろにあれば `ipps://`）でも構いません。ジョブはそこからノードのエッジをたどります。`bdf printers -url ipp://<host>:<port>/ipp/print` はそれらを一覧します。CUPS クライアントの `lpstat` も同じです（CUPS-Get-Printers）。ほかのノードは IPP プリンタではありません。キュー、プリンタ、ディスプレイにはエッジを通って届きます。リクエストのジョブ属性（`media`、`copies`、`job-name` ……）と `document-metadata` はメタデータ `ipp.<name>` としてジョブとともに流れます。アカウント（`job-account-id`、`job-accounting-user-id`）、`document-metadata`、独自の属性（`x-…`）は、ジョブが届くプリンタへ渡ります（[プリンタへ渡るもの](flow.md#ipp-attributes)）。

IPP で印刷できる人は運用者の `submit_tokens` で決まります。1 つでも設定されていれば、クライアントはトークンの名前とシークレットを Basic 認証で渡します（アプリも id とシークレットを同じように渡します）。設定がなければ、マネージャーに届く誰もが送れます。クライアントが渡す requesting-user-name は「（… として）」と表示されるだけで、確認はしません。

## アップローダー {#uploaders}

アップローダーはアプリケーション（業務システム、WMS）向けの REST の入口です。トポロジー上のノードです（[アップローダー](nodes.md#upload)）。最初に保存したときにキーが 1 度だけ表示され、サイドバーの「キーを回す」で新しいキーを作って 1 度だけ表示します（古いキーは使えなくなります）。マネージャーが保存するのはキーのハッシュだけです。

アプリケーションはキーをベアラートークンとして送ります。

```sh
curl -X POST -H "Authorization: Bearer bdfk_…" \
     --data-binary @statement.pdf \
     "https://<host>/api/upload/erp?name=statement&meta.dept=sales"
```

フォームとして送るなら次のとおりです。

```sh
curl -H "Authorization: Bearer bdfk_…" \
     -F 'meta={"dept": "sales"};type=application/json' \
     -F file=@statement.pdf \
     https://<host>/api/upload/erp
```

- ボディは文書、`accept: data` のアップローダーならデータです。内容から形式が分からないときは、名前（`?name=`）と `Content-Type` で判別します。
- `?name=` はジョブの名前、`?user=` は送った人の名前で、「（… として）」と表示されます。
- `X-Meta-<Key>` ヘッダと `?meta.<key>=` パラメータはジョブのメタデータです（[メタデータ](flow.md#metadata)）。
- ボディは `multipart/form-data` でも構いません。`Content-Type` で見分けます。文書はパート `file` に入れ、その名前と `Content-Type` が文書のものになります（名前は `?name=` が優先）。メタデータはパート `meta` に JSON オブジェクトで入れ、ヘッダとパラメータのものより優先されます。別の Docflow の [HTTP (push)](nodes.md#http) ノードはどちらの形でも送ります。
- `?name=` がなければ、そのままのボディは `Content-Disposition`（`attachment; filename="statement.pdf"`）から名前を取ります。
- 応答はジョブ（201）で、番号とそこまでのレグが入ります。アップローダーが受け取らない形式の文書は拒否され（415）、拒否は監査ログに残ります。キーが違えば拒否されます（401）。
- アップローダーの直後の[アサート](nodes.md#assert)は、応答の前にメタデータを検査します。通らなかったジョブは `{"error": "…", "job": <番号>, "failures": [{"node", "key", "error"}]}` とともに拒否され（422）、何が通らなかったかをレグで見られるよう、中止として残ります。CUPS サーバーは同じことを IPP の `client-error-bad-request` で答え、何が通らなかったかをステータスメッセージに入れます。

アップローダーが何を受け取るか（メディアタイプ、テンプレートが埋めるデータのスキーマ、後ろのノードが見るメタデータ）は、その OpenAPI 文書に載ります。サイドバーのリンク OpenAPI か、`GET /api/openapi.json?node=<name>` です（[OpenAPI](api.md#openapi)）。

アプリケーションが呼ぶものはすべて、マネージャー自身のサーバーの API にあり、ノードの種類と名前で決まります。アップローダーは `/api/upload/<name>`、S3 トリガーは `/api/s3-trigger/<name>`、インボックスは `/api/inbox/<name>`（[インボックスの API](inboxes.md#api)）、ダウンローダーは `/api/download/<name>`（[ダウンローダー](downloads.md)）です。アップローダーには専用のポートがありません。ポートで待ち受けるのは CUPS サーバーだけで、IPP のためです。

## CUPS サーバーノード {#cups}

CUPS サーバーノード（[CUPS サーバー](nodes.md#cups)）は、印刷するコンピュータや業務システム向けです。プロセス自身のポートで動く CUPS サーバーで、プリンタはノード自身の 1 台です。コンピュータは、ほかの CUPS プリンタと同じように `ipp://<host>:<port>/printers/<name>`（ノードに証明書があれば `ipps://`）でプリンタを追加し、そこに印刷したものはノードの `to` へ進みます。プリンタならノードは中継になり、シンクなら印刷されたものが見え、フィルタなら条件に合うものだけが残ります。`lpstat -h <host>:<port> -p` で一覧できます。新しい CUPS サーバーには、6310 から、どのノードも使っていない最初のポートが付きます。逆に、このシステムから CUPS サーバーのキューへ印刷するには、エッジの先に [CUPS サーバー（シンク）](nodes.md#cups-out)を置きます。ノードの使い方タブに、URL、macOS・iOS・Android・Windows・CUPS での追加のしかた、コマンドが表示されます。

プリンタは AirPrint、IPP Everywhere、Mopria のプリンタと同じように答えます（PDF 優先、それぞれのラスター（Apple、PWG、PCLm）、A3 からはがきまでの用紙）。そのため Mac は AirPrint のキューを作って PDF を送り、Android の印刷サービス（デフォルト印刷サービス、Mopria Print Service）は Mopria のプリンタとして、Windows（IPP クラスドライバ）と CUPS は IPP Everywhere のプリンタとして扱います。Mac でアドレスから追加するときは、ドライバに AirPrint を選びます。一般的な PostScript プリンタは PostScript を送りますが、これは受け付けず、Mac はジョブを再送し続けます（拒否のたびに `ipp: print-job refused` がログに出ます）。Bonjour（`dnssd`。`false` にしない限り入っています）で告知するので、ポートが開いている間、ローカルネットワークの Mac、iPhone、iPad、Android のスマートフォンとタブレット、Windows と Linux のコンピュータが「<タイトル> @ <コンピュータ>」として見つけます。ホスト名はプロセス自身のもの（`<コンピュータ>-docflow.local`）です。

ポートでは、IPP Everywhere が求め、これらのクライアントが送るとおり、ジョブが 2 つのリクエストで来ることもあります。ジョブの番号を返す Create-Job と、文書を運ぶ Send-Document です（文書はジョブに 1 つ。5 分以内に文書が来ないジョブは中止します）。マネージャー自身のサーバー（`/ipp/print/<name>`）は、どのインスタンスも文書を待たずに済むよう、1 つのリクエストで来るジョブ（Print-Job）だけを受けます。文書はそのまま先へ進みます。印刷するページはクライアントが選んで送り、`page-ranges` は受けません。

## FTP と SFTP（見回り） {#poll}

[FTP（見回り）](nodes.md#ftp-poll)と [SFTP（見回り）](nodes.md#sftp-poll)のノードは、API を呼ぶのではなくサーバーにファイルを置くシステム向けです。EDI の取引先の送信箱、バッチの書き出し、SFTP で公開したディレクトリに書くスキャナなどです。ノードは待ち受けません。tick のたび（`POST /api/tick`。スケジューラか、プロセスのタイマー）、**間隔**より密にはならずにディレクトリを一覧し、ファイルを改名で取り置き（`<name>.docflow-<node>.taking`）、アップローダーと同じように取り込み（入口の accept・output・formats が効き、ファイルのパス・大きさ・時刻がメタデータに入ります）、消すか済みのディレクトリへ移します。逆にサーバーへ置くのは [FTP](nodes.md#ftp) と [SFTP](nodes.md#sftp) の出口です。Docflow 自身は FTP や SFTP のサーバーにはなりません。押すことしかできない送り手は自前のサーバーに置き、見回りのノードがそれを見ます。

## アプリ {#apps}

アプリは、アップローダーのキーではなく自前の資格情報で送ったり、インボックスを扱ったり、ダウンローダーを読んだりするアプリケーションです。資格情報はアプリのもので、監査ログとジョブにはどのアプリが送ったかが残ります。管理 → アプリで、管理者が登録します。id、題名、そのジョブを見られるメンバーのグループ（任意）、そしてスコープです。

| スコープ | |
|---|---|
| `print` | 文書を送る |
| `print.template` | テンプレートを埋める。データを送る |
| `jobs.cancel` | 自分のジョブを取り消す |
| `inbox` | 自分を挙げたインボックスを扱う。留まっているものを一覧し、サムネイルと文書を見て、送り出しと破棄をする（[インボックスの API](inboxes.md#api)） |
| `download` | 自分を挙げたダウンローダーを読む（[バックエンド](downloads.md#backends)） |
| `tick` | `POST /api/tick` を呼ぶ。スケジューラ向け |
| `admin` | 管理者として API を使う。スクリプト向け |

シークレットはアプリを追加したときに 1 度だけ表示され、「シークレットを回す」で新しいものを作ります（古いものは使えなくなります）。「無効にする」はアプリの資格情報と発行済みのトークンを拒否し、「削除」はアプリを消します。

アプリは `POST /oauth/token`（OAuth 2.0 のクライアントクレデンシャルグラント）でトークンを取り、ベアラーとして、そのアプリを「許可するアプリ」（アップローダーのサイドバー）に挙げたアップローダーへ送ります。または、同じようにそのアプリを挙げた CUPS サーバーノードへ、id とシークレットを Basic 認証で渡して IPP で印刷します。

```sh
curl -u 'erp:<secret>' -d grant_type=client_credentials https://<host>/oauth/token
# {"access_token": "…", "token_type": "Bearer", "expires_in": 600, "scope": "print"}
curl -X POST -H "Authorization: Bearer <access_token>" --data-binary @statement.pdf "https://<host>/api/upload/erp?name=statement"
```

トークンは 10 分間有効です。アプリは自分のジョブを見られ（`GET /api/my/jobs`）、スコープ `jobs.cancel` があれば取り消せます。ノードがグループを与えないとき、アプリのジョブはアプリのグループのものです。

## スクリプト {#scripts}

スクリプトは運用者の管理トークン、またはサブミットトークンの `name:secret` をベアラーとして送ります。`POST /api/destinations/<node>/jobs` に文書をボディとして送り、アップローダーと同じく `?name=`、`?user=`、メタデータを受け取ります（[ジョブ](api.md#jobs)）。
