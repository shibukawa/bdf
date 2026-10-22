# Print Manager のデータ取り出し口（export）— 保管から認証で取り出す、メタデータの持ち主で絞る、バックエンドとフロントエンド

2026-10-22 の実装案と記録。同日に実装した（`.knowledge/` の `requirement:print-manager-export` と `milestone:pm-export-1`）。土台は `milestone:pm-4`（archive）、`milestone:pm-auth-2`（アプリと client credentials）、`milestone:fl-1`（メタデータ）。

利用者「print manager にデータ取り出しのエンドポイントという sink を作りたい。認証で archive から取り出す。メタデータでフィルタをして、送信したユーザーの ID にマッチしなければ 404 みたいなこともしたい。バックエンドサーバーから client credential でアクセス以外に、oauth でフロントエンドからアクセスとかも考えたい」。

## 0. 今の形（2026-10-22 の main、#130 まで）

- archive は終点で、辺を持たない。書いたものは Manager から読み出す口が無い（保管庫 pm-2 の読み出しは未）。
- sink は「最新の文書を画面で見る」終点で、ジョブの文書（`job-documents`）を Live で守り、押し出されたら手放す。外から取り出す口ではない。
- メタデータは送り手の名乗り（`X-Meta-*`、`?meta.*`、IPP の属性）と入口の `meta` と `$` の確かめた鍵（fl-1）。
- アプリは client credentials で 10 分のトークン（HMAC の封）。人は OIDC のセッション（クッキー）だけで、API のトークンを持たない。

## 1. 決めたこと（Claude の暫定。確認待ち）

### 1.1 新しい種類 `export`（sink を広げない）

利用者は「sink」と呼んだが、今の sink は「画面で最新を見る」もので、文書をジョブと一緒に手放す。取り出し口は保管と同じ寿命が要るので別の種類にした。名前は `export`（出口。外から読まれる）。画面のパレットでは「Export: lists what an archive wrote」。

### 1.2 archive の後ろに置く（archive に `to` を足す）

「archive から取り出す」をそのまま図にした：archive は書き終えたら（Delivery が done）`to` のノードへ続ける。export はその後ろにだけ置ける。

- PUT の確かめ：入口・何も送ってこないノード（自分が宛先）・template・convert・inbox（新しいペイロードか、流し直しで作る経路）から辺を辿り、archive を通らずに export に着く道があれば 400（理由に道：`api → out`）。
- 実行時：経路（routing）が「最後の archive が書いたもの（archive の名前と鍵）」を運び、export はそれが無ければ error の足跡。queue・class・filter は経路の中で運ぶので、archive → filter → export（種類で振り分け）は通る。filter に直接送られたジョブは error。export そのものに直接送るジョブは 400 で断り、export は宛先の一覧（画面の Print here）に出さない。
- export は文書をジョブから持たず、**archive から読む**。ジョブの記録が `JobKeep`（14 日）で消えても、項目と文書は残る（保管の寿命）。

### 1.3 項目（entry）と一覧

export は着いたジョブを項目として記録する（`export-entries`、ID は `<ノード>.<ジョブ>`）：ジョブの番号・名前・種類・形式・ページ・大きさ・sha256、archive の名前と鍵とオブジェクト、**持ち主**、写したメタデータ（`fields`）、送った主体とグループ。一覧（`export-lists`）は全員の分と持ち主ごとの分で、最新 `entries`（既定 1000、上限 5000）件の ID。一覧から押し出された項目もジョブの番号で読める。

### 1.4 持ち主と 404

`by` はメタデータの鍵で、その値が持ち主（送り手が `X-Meta-Customer: alice@example.com` と付ける、など）。人が読むとき、持ち主と人の `subject` を比べる：

| subject | 比べるもの |
|---|---|
| `email`（既定） | IdP が確かめたメールアドレス（大文字小文字を区別しない） |
| `sub` | IdP の subject |
| `id` | Manager の人の ID（issuer と subject のハッシュ） |
| `group` | 人のグループのどれか（groups claim の値） |

自分のものでない項目、無いジョブ、別の export のジョブ、`?owner=` で他人を指したものは全部 **404**（403 にしない。export が何を持っているかを言わない）。`by` の無い export は、人がそのジョブを見られるか（送った本人・ジョブのグループ。`CanSee` と同じ）で絞る。

**持ち主は名乗り**（fl-1 の制約のまま）：入口に送れる者が、誰の文書かを決める。入口をバックエンド（自分の利用者を知っているもの）にだけ渡すか、入口の `meta` で固定する。README に書いた。

### 1.5 誰が読めるか

| 読み手 | 資格 | 読めるもの |
|---|---|---|
| バックエンド | アプリ（client credentials、10 分）。export の `apps` に載り、scope `export` | 全部。`?owner=` で 1 人の分 |
| 管理者・監査 | セッション、管理トークン、scope admin のアプリ | 全部 |
| 人（画面） | OIDC のセッション | 自分が持ち主の項目 |
| 人（フロントエンド） | 公開クライアントが人のために取ったトークン（1 時間） | 自分が持ち主の項目 |
| それ以外 | submit トークン、載っていないアプリ、scope の無いアプリ | 403。資格なしは 401 |

export の `scope` と `owner`（group・personal）は「誰が読んでよいか」の外枠（送り先の範囲と同じ判断）。

### 1.6 口

| | |
|---|---|
| `GET /out/<名前>` | 読める項目、新しい順（`?owner=`、`?limit=` 既定 100） |
| `GET /out/<名前>/<ジョブ>` | 1 項目 |
| `GET /out/<名前>/<ジョブ>/record` | archive が書いた記録の JSON（メタデータ・ペイロード・足跡） |
| `GET /out/<名前>/<ジョブ>/document` | archive が書いた文書（監査に `export-document`） |

`/out/**` は popcornweb の保護の外（自分で確かめる）。セッションのクッキーも効く（`Options.Viewer`）。

### 1.7 フロントエンド：公開クライアントと認可コード + PKCE

pm-auth の `authorization_server: hybrid`（人は外の IdP、機械には Manager が小さな認可サーバー）を、人のトークンにも広げた。

- アプリに `public`（秘密なし）と `redirect_uris`（https、localhost だけ http、fragment なし、8 つまで）を足した。公開クライアントは client credentials を取れず、秘密の回転も断る。
- `GET /oauth/authorize`（popcornweb の保護の中のページ）：登録したクライアントと redirect URI でなければ JSON の 400（誤ったところへ飛ばさない）。それ以外の誤り（`response_type`、S256 の PKCE が無い、クライアントに無い scope）は RFC 6749 §4.1.2.1 のとおり redirect URI に `error` で返す。ログインしていなければ IdP へ行き、戻ってから code を返す。**同意の画面は無い**：管理者がクライアントを登録したことが同意（社内のフロントエンドの想定）。
- code は 32 バイトの乱数、5 分、1 度だけ（`PutState`・`TakeState`。ログインの状態と同じ置き場。間違った verifier でも消える）。
- `POST /oauth/token` の `grant_type=authorization_code`：`client_id`・`code`・`code_verifier`・`redirect_uri` が code と合えば、1 時間のトークン（HMAC の封、種類 `user`）。refresh token は無い（期限が来たら `/oauth/authorize` をもう一度。セッションがあれば IdP を経ず即座に戻る）。
- トークンは人（`Via: oauth`）で、役割を持たない。効くのは `/out/**` と `/api/me` だけ（他の API は 401）。要求ごとに人の記録（最後のログインのグループ、sub）とクライアント（無効化・scope の縮小）を読み直す。
- CORS：登録した redirect URI の origin からの `/oauth/token` と `/out/**` に `Access-Control-Allow-Origin`（`Allow-Credentials` は付けない。クッキーは別 origin へ行かない）。
- 監査：`oauth-code`、`user-token`、`code-refused`。

### 1.8 フロントエンドのもう 1 つの形（BFF）

フロントエンドに自分のバックエンドがあるなら、トークンをブラウザに置かない形の方が簡単で強い：バックエンドが client credentials で `?owner=<その人>` を読み、自分のセッションでブラウザに渡す。Manager 側の追加は要らない（1.5 のバックエンドの行）。README の例はこの形も含む。

## 2. 構成

```
print/manager/export.go     KindExport の設定（checkExport）、exportsArchived（PUT の道の確かめ）、routeExport・listAt、Entry・Entries、dropExport、/out の口（exportViewer・sees・exportObject）
print/manager/oauth.go      公開クライアントの redirect URI、CORS、/oauth/authorize、RedeemCode、userTokenViewer、authorization_code の grant
print/manager/model.go      KindExport、Node に By・Subject・Fields・Entries、LegListed、archive の To
print/manager/scope.go      ViaOAuth、Viewer.Subject、SubjectEmail・Sub・ID・Group、identities・owns、Token()
print/manager/apps.go       ScopeExport、App.Public・RedirectURIs、PutApp と RotateAppSecret の公開クライアント、/oauth/token の grant の振り分け
print/manager/deliver.go    archive の done の後に continueFrom（archiveRef つき）
print/manager/transform.go  continueWith を continueFrom（made の有無と archiveRef）に一般化
print/manager/route.go      routing.archived、KindExport の route
print/manager/people.go     viewerOf に人のトークン（/out と /api/me だけ）
print/manager/store.go      export-entries・export-lists
cmd/printmanager            publicPaths に /out/** と /oauth/token（/oauth/authorize は保護の中）、viewer に Subject、パレットに Export、export のサイドバー（Owner's key・Subject・Fields・Entries・Apps allowed、Readers、Listed here）、archive の Forwards to と出力の口、/admin/apps の公開クライアントと redirect URI、README の Exports
```

## 3. 試験

`go test -race ./print/manager/`：TestExport（設定の 400、archive の無い道の 400 と道、template・convert・inbox からの道、filter を挟む道、持ち主のメールの大文字小文字、alice・bob・root・バックエンド・scope なし・載っていないアプリ・submit トークン・資格なし、group の export、1 項目の 200 と 404 の表、文書と記録、監査、topology と live の数、archive の段の 2 度目で 2 度載らない、直接送られた export は error、by なしの export、subject group、ジョブが消えた後も項目と文書、export の削除で項目と一覧が消え他の export は残る）、TestExportAfterFilter（archive → filter → 2 つの export、直接 filter は error）、TestOAuthCode（公開クライアントの登録の確かめ、秘密なし、client credentials は 401、authorize の 400・401・redirect の error、S256 以外、scope、code の 1 度きり、verifier と redirect URI の不一致、トークンの scope と期限、/api/me、他の API は 401、CORS と preflight、scope の縮小・無効化・期限、監査、機密クライアントに戻すと秘密）。`go test ./cmd/printmanager/...`：TestOAuthAuthorize（実の popcornweb と偽の IdP で、未ログインは /auth/login へ、ログイン後は code、RFC 7636 の例の verifier で交換、トークンで /api/me）。

## 4. 残り（問い）

- **名前**：`export`（利用者の言う sink とは別の種類にした）。sink を広げる方がよければ替える。
- **フロントエンドのトークンの出所**：今は Manager が出す（hybrid）。会社の IdP が出したアクセストークン（JWT）を Manager が JWKS で確かめる形（pm-auth の `external-only`・`later`）を足すか。IdP 側にこの API を登録する手間と、IdP ごとの aud・scope の違いがある。
- **持ち主の信頼**：送り手の名乗りのまま。「確かめた送り手（`$sender`）の分だけ」などの選択を足すか。
- **項目の寿命**：export の項目は消えない（archive と同じ）。件数が大きくなる。`dropExport` は全部を走査する。保管庫（pm-2）の保持の期限と合わせて刈るか。
- **ページ送り**：新しい順の `limit` だけ（カーソルが無い）。
- **大きな文書**：Manager が取り次ぐ。S3 の presigned URL で直接渡すか。
- **同意の画面**：無い（登録が同意）。社外のフロントエンドを許すなら要る。
- **refresh token**：無い（1 時間ごとに authorize をもう一度）。

確かめた：pw dev を内蔵ブラウザで、一般の人（Mina）のサイドバーに自分の 1 件だけ、管理者の編集のサイドバー、/admin/apps の公開クライアント。そこで `/out/` を Web アプリに mount し忘れていたこと（404）と、export に Print here が出ることを見つけて直した。

未確認：実の IdP と実のフロントエンド（別 origin の SPA からの PKCE と CORS）、S3・R2・GCS の archive から読むこと。
