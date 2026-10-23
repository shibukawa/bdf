# マネージャーを動かす {#running-the-manager}

マネージャーは Popcorn Web 上の 1 つのバイナリ `docflow` です。プロセスに状態を持ちません。記録と文書はストアにあるので、同じバイナリが拠点のサーバーで 1 プロセスとして、ロードバランサの後ろで複数として、あるいは AWS Lambda の関数として動きます。設定はファイル `config.<APP_ENV>.toml`（`docflow --generate-config toml` がすべてのキーを出力します）か、下に挙げる環境変数で行います。

## ストアとシークレット {#store}

| 変数 | |
|---|---|
| `DOCFLOW_RECORDS` | 記録のストア。`bbolt:path`（1 プロセス）、`dynamodb://table?region=…`、`firestore://project/database`、`s3://bucket?region=…&endpoint=…`、`r2://bucket?endpoint=…`、`gcs://bucket` |
| `DOCFLOW_DOCUMENTS` | 文書のストア（`bbolt`、`s3`、`r2`、`gcs`）。空なら記録のストア |
| `DOCFLOW_SECRET` | デバイスのセッション、アプリのトークン、ページのセッションに署名します。32 バイト以上で、すべてのインスタンスで同じもの。開発以外では必須 |
| `DOCFLOW_PREVIOUS_SECRETS` | 退役したシークレット。ローテーションの間、それで署名したものをまだ開けます。新しいシークレットを設定し、古いものをここに 1 日置き、その後に外します |
| `DOCFLOW_ADMIN_TOKEN` | API のベアラートークン。スクリプト向けです。人はサインインします |
| `DOCFLOW_TICK_TOKEN` | スケジューラが `POST /api/tick` を呼ぶときのベアラートークン |
| `DOCFLOW_SUBMIT_TOKENS` | IPP と API でジョブを送れる `name:secret` の組。なければ誰でも送れます |
| `DOCFLOW_TICK_INTERVAL` | プロセスがティックする間隔（`30s`）。`0` ならスケジューラに任せます |
| `DOCFLOW_JOB_KEEP` | 終わったジョブを保存する期間（`336h`、14 日） |
| `DOCFLOW_MAX_DOCUMENT_MB` | 受け取る文書の最大サイズ（64） |
| `DOCFLOW_FONT_DIRS`、`DOCFLOW_NO_SYSTEM_FONTS` | 文書を描くフォント。フォントのディレクトリと、マシン自身のフォントを使うかどうか |
| `DOCFLOW_SECURE_COOKIES` | クッキーの `Secure`（true。false は開発で平文の http を使うときだけ） |
| `DOCFLOW_DEVICE_URL` | ディスプレイのデバイスが開くオリジン（登録の URL と QR コード）。`https://docflow.example.com` など。空なら、管理者がアクセスしたときのマネージャー（[ネットワーク上のデバイス](displays.md#network)） |
| `DOCFLOW_HTTPS_PORT`、`DOCFLOW_TLS_CERT_FILE`、`DOCFLOW_TLS_KEY_FILE` | ネットワーク上のデバイスのための、HTTP と同じものを返す HTTPS の待ち受けと、その証明書（続けて中間証明書）と鍵（PEM）。開発では、なければ自己署名証明書 |
| `DOCFLOW_SERVERLESS` | 1 プロセスしか使えないストアを拒否し、呼ばれたときだけティックし、ポートを開きません |

## サインイン {#oidc}

人は OpenID Connect プロバイダでサインインします。Popcorn Web の `[auth.oidc]` です。

| 変数 | |
|---|---|
| `AUTH_OIDC_ISSUER`、`AUTH_OIDC_CLIENT_ID`、`AUTH_OIDC_CLIENT_SECRET` | プロバイダと、このクライアント |
| `AUTH_OIDC_REDIRECT_URL` | `https://<host>/auth/callback`。プロバイダに登録したもの |

誰が入れるか、誰がどの役割を持つか（[役割](roles.md#roles)）は `[oidc]` です。

| 変数 | |
|---|---|
| `OIDC_GROUPS_CLAIM` | グループのクレーム。`groups`、`roles`、または Keycloak の `realm_access.roles` のようなドット区切りのパス（ID トークンに含まれる必要があります） |
| `OIDC_GROUPS_PREFIX` | この接頭辞のグループだけを取ります |
| `OIDC_ADMINS`、`OIDC_ADMIN_GROUPS` | 管理者。`issuer#subject` か確認済みのメールアドレス、およびグループ。開発以外では 1 つは必須 |
| `OIDC_TOPOLOGY_ADMINS`、`OIDC_TOPOLOGY_ADMIN_GROUPS` | 役割 `topology` |
| `OIDC_AUDITORS`、`OIDC_AUDITOR_GROUPS` | 役割 `audit` |
| `OIDC_ALLOWED_ISSUERS`、`OIDC_ALLOWED_DOMAINS`、`OIDC_ALLOWED_GROUPS` | 入れる人を絞ります。発行者、メールのドメイン（または Google の `hd`）、グループ。管理者は常に入れます |

## サーバーレス {#serverless}

AWS Lambda（Lambda Web Adapter 経由）では、関数は環境変数で設定し、`DOCFLOW_SERVERLESS=true` を含めます。記録は DynamoDB（文字列キー `pk` と `sk` のテーブル）に、文書は S3 に置きます。EventBridge Scheduler が毎分、ティックトークンで `POST /api/tick` を呼びます。ノード自身のポートはありません。アプリケーションはマネージャーの `/in/<name>` を呼びます。拠点のプリンタに印刷するには、関数をそこに届く VPC に置きます。関数とそのイメージのビルドはリポジトリの README にあります。

## プリンタとフォント {#printers}

プリンタノードはネットワーク越しにプリンタへ届きます。`ipp`、`ipps`、`socket`（ポート 9100）、`lpd`、または CUPS サーバーのキューです。文書はマシンのフォントと `DOCFLOW_FONT_DIRS` のフォントで描かれます。トポロジーを設定している間は、偽プリンタファーム（`tools/fakeprinter`）がプリンタの代わりになります。プリンタノードの「プレビュー」に、偽プリンタが受け取ったものが表示されます。
