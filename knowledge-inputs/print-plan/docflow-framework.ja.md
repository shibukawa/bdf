# Docflow をフレームワークにする案（2026-10-26）

2026-10-26 の利用者の依頼を `/knowledge-memory` で整理した叙述。要件の概念は `.knowledge/` の `requirement:docflow-framework`（問いと暫定の選択はそこ）、最初の段は `milestone:fw-1`。土台は `requirement:docflow`、`milestone:df-1`（種類の登録）、`requirement:docflow-node-editing`（エディタの必須の欄）、`system:popcornweb`。

利用者「docflow はひととおりできた。docflow をベースに、static import でちょっとアレンジして独自ノードを追加するなどして、アレンジできるようなフレームワーク化をしたい。また、popcornweb ベースなので、構造体を登録するだけで設定は自由に足せるはず。S3 の接続先のような設定も、そこで定義した設定をプルダウンで選べたりする、トポロジーもプリセットのを組み込んでスタンドアロンアプリとできる、独自管理画面やユーザー向け画面を足せる、みたいなことが実現できると良いと思う」。

## 0. 今の形（2026-10-26 の main、#165 まで。コードで確かめた）

- **種類の登録**は df-1 で `docflow.RegisterKind` になっている。芯の種類は `kinds_core.go`、入出力は `docflow/io/{ipp,http,smtp,s3,notify}` と `docflow/render` が `init` で登録し、`cmd/docflow/handlers` が blank import する。外のモジュールが同じ口で登録すること自体は今でもできる。
- **種類ごとの設定は平らな `Node`**（df-1 の利用者の答え）。`url`・`credential`・`path`・`text`・`headers` … を使い回す。外で作る種類は自分の欄を持てない。
- **エディタ**（`topology.js` 2078 行、`canvas.js` 1122 行）は種類ごとの `switch` で専用の欄を作る。アダプタの種類は JSON の textarea（`ADAPTER_SETTINGS` の説明文と `SETTINGS` の鍵の白名簿）。**パレットは `pages.pw.html` に固定**（`kind-xxx` のメッセージ）、アイコンと役割の色は `canvas.js` の表（`ICONS`・`ROLE_OF`・`PASSES_ON`・`MAKES_BDF`）。`GET /api/kinds` は名前・題・役割・受ける種類・`required` までを返す（df-1 の「エディタはまだ自分の表を使う」の積み残し）。つまり外で登録した種類は、パレットに出ず、アイコンは queue の代用で、設定は JSON の欄になる。
- **設定**は `handlers.RegisterConfig` が `docflow` と `oidc` の構造体を pwconfig に登録する。アプリは `pw.RegisterConfig[T]` で足せる（生成の束縛は自分のディレクトリに書かれる）。
- **画面**は `handlers.App`（`Open`・`Handler`・`StartTimer`・`StartHTTPS`・`Outermost`・`Close`）。`Handler()` は閉じた mux で、`AppBar` のリンクは固定、`page()`・`viewer()`・`account()` は非公開、`Shell`・`AppBar` は export されていない部品。公開の木（JS・CSS）と pw の manifest は `package main`（`cmd/docflow`）。メッセージの catalog は `cmd/docflow/messages`、マニュアルは `cmd/docflow/help`（章は固定の一覧）。
- **保存**の外向きの口は `Store.PutState/TakeState/PruneStates` だけ。種類が自分のコレクションを持つ口は無い。
- **試験**の土台 `docflow/internal/flowtest`（Env の Put・Submit・Tick・Delivery・Required）は外から使えない。
- **芯（`docflow`・`docflow/io/...`・`docflow/render`）は popcornweb に依存しない**（`go list -deps` で 0）。
- **popcornweb（v0.5.14）の制約**（コードで確かめた）：`middlewares.RegisterPublicManifest` は 1 つの表だけ（2 度目は置き換え）。`pwruntime.RegisterLocales` は同じ集合なら no-op、違えば panic。`pw.WriteHTML` は登録された document で葉を包む。`document.pw.html` は 1 プロジェクトに 1 つ。`[]T` の設定は配列の表（TOML だけ、環境変数は `${NAME}` 参照）。`pw` はプロジェクトの直下に `go.mod` を求める。

## 1. ねらいと組み立ての形

「static import で組む」を、**Go のモジュールが bdf を import し、`init` の登録と `main` の数行で、独自の種類・設定・カタログ・プリセット・画面を足した 1 つのバイナリを `pw build` で作る**こととする。プラグインの動的読み込み・Wasm・別プロセスはしない（別プロセスの変換は `requirement:document-flow` の外部サービスのまま）。

利用者側の `main` の見本（暫定の名前）：

```go
package main

import (
	_ "github.com/shibukawa/bdf/docflow/io/all" // 標準の入出力（要るものだけでもよい）
	"github.com/shibukawa/bdf/docflow/web"
	"github.com/shibukawa/popcornweb/pw"

	_ "example.com/acme/docflow/kinds/folder" // 独自の種類（init で RegisterKind）
	"example.com/acme/docflow/pages"
)

func main() {
	web.RegisterConfig()                        // [docflow] と [oidc]
	pw.RegisterConfig[pages.AcmeConfig]("acme") // 自分の設定
	web.Main(web.Options{
		Preset: web.PresetFS(presets, "topology.yaml", docflow.PresetManaged),
		Routes: pages.Mount, // 自分の画面と API
		Nav:    pages.Nav,   // バーの項目
		Public: pages.Assets, // 自分の公開の木（pw の /public）
	})
}
```

パッケージの配置（Claude の暫定、問い 3）：

| パッケージ | 役 |
|---|---|
| `docflow` | 芯。足すのは `Settings`・`Fields`・`Text`・カタログ・プリセット。popcornweb に依存しないまま |
| `docflow/web`（`cmd/docflow/handlers`・`messages`・`help`・`public` を移す） | Web アプリの部品：`Config`・`RegisterConfig`・`Main`・`Run`・`App`・`Page`・`API`・`BindShell`/`BindAppBar`・`RegisterNav`・`RegisterRole`・`RegisterHelp` |
| `docflow/flowtest`（`internal` から出す） | 独自の種類の試験 |
| `docflow/io/all` | 標準の入出力をまとめて blank import する便宜 |
| `cmd/docflow` | 参照アプリ。`web.RegisterConfig(); web.Main(web.Options{})` だけ |
| `examples/docflow-custom` | 独自のアプリの見本（入れ子のモジュール、`replace ../..`。pw は直下に go.mod を求めるので） |

## 2. 独自のノードの種類

### 2.1 設定の持ち方：`Node.Settings`（問い 1）

- `Node` に `Settings json.RawMessage`（JSON の `settings`）を足す：**その種類だけの設定のオブジェクト**。芯は中の意味を知らず、`Kind.Fields` に照らして形（型・必須・カタログの値）だけ確かめ、それから `Check` に渡す。
- 芯の種類と今のアダプタ（smtp・s3・http・notify・render・s3-trigger・cups・printer）の平らな欄は当面そのまま（df-1 の答えを覆さない）。`settings` に移すかは問い 2。
- Go の側：`docflow.SettingsOf[T any](n *Node) (T, error)` で `Check`・`Send`・`Step` が自分の構造体に読む。
- 辺：独自の種類も次のノードは芯の欄 `to`・`on_error` で指す（`PassesOn`・`Makes`・`Takes` は今のまま）。branch の case のような独自の口は後（fw-4）。

平らな `Node` に `extra map[string]any` を足す案（B）もあるが、種類ごとの欄が芯の欄と同じ名前空間に混ざり、YAML と OpenAPI の形が崩れる。A を推す。

### 2.2 `Kind` に足すもの

```go
type Kind struct {
	// … 今のもの（Name・Title・Role・Takes・PassesOn・Stocks・Makes・Fixed・Check・Required・Ready・Send・Serve・Step・Keyed・Fetch・FetchSchema・Discover）
	Titles  Text   // 題（言語ごと。en は必須、無い言語は en。Title は en の略記）
	Help    Text   // 説明（パレットとサイドバー）
	Group   string // パレットの見出し。既定は Role から（entry→source、route→control、step→convert、exit→sink）
	Icon    []string // 24x24・stroke の path。無ければ役割の既定
	Fields  []Field  // Settings の欄（エディタのフォームと PutNode の確かめ）
	Doc     Text     // マニュアルの節（Markdown、言語ごと）
	Handler func(m *Manager) http.Handler // 任意。/api/<kind>/{name}/… に載る自分の API
	Access  string   // Handler の門：既定は役割 topology。"viewer"（サインインかトークン）か "public" で緩める
}

// Text は言語の札 → 文。API は全部を返し、エディタが lang で選ぶ。
type Text map[string]string

type Field struct {
	Key         string
	Type        FieldType // string・text・int・bool・duration・enum・catalog・credential・nodes・path・json
	Title, Help Text
	Placeholder Text
	Required bool
	If       map[string]string // 条件つきの必須（今の Required.If と同じ）
	Choices  []Choice          // enum
	Catalog  string            // catalog：登録したカタログの名前
	Open     bool              // catalog でも自由入力を許す
	Default  any
	Expands  bool // {{meta.<key>}} を展開する欄（説明に出す）
}
```

- 今の `Required []Required` は `Fields` から作る（残す）。`KindInfo` に `titles`・`help`・`group`・`icon`・`fields` と、`canvas.js` が表で持つもの（`passes_on`・`makes`・`stocks`・`takes` は既にある）を足し、エディタは `GET /api/kinds` だけで種類を知る。

### 2.3 エディタ

- **パレット**：`pages.pw.html` の固定の一覧を `GET /api/kinds` から作る（役割 → 4 つの見出し、`Group` で上書き）。芯の種類の文（`kind-xxx`）は `Titles`・`Help` に写す。
- **アイコンと表**：`Kind.Icon` → `canvas.js` の `ICONS` の代わりに API から。`ROLE_OF`・`PASSES_ON`・`MAKES_BDF`・`TAKES_STOCKED`・`PRINTS` も API の値から。
- **設定の欄**：`Fields` から作る汎用のフォーム（`fieldsForm`）。型ごとの部品は、文字列・複数行・数・真偽・期間・選択・カタログの pulldown（`GET /api/catalogs/<name>`）・資格の名前（カタログ `credentials` の pulldown と自由入力）・ノードの picker（既存）・path。必須の印は今の仕組み（`missingOf` が `fields` を見る。`requirement:docflow-node-editing`）。
- 芯の種類の専用の欄（Bonjour、Find queues、テンプレートの一覧…）はそのまま。

### 2.4 実行

- 役割ごとの口は今のもの：入口（`Keyed`+`Fetch`、`Serve`）、経路（`PassesOn`、`Step` で通過しつつ何かする＝notify の形）、段（`Step`）、出口（`Send`）。**独自の分岐**（`to` の一部を選ぶ）は `Kind.Pick func(ctx, *Out) ([]string, error)`（fw-4）。
- `Kind.Handler` の門は既定で役割 topology（`Access` で viewer・public に緩める）。中では `docflow.ViewerOf(r)` で誰かを知る。
- 保存済みの `settings` は PUT のときだけ確かめる。種類の欄が後から変わっても古いノードはそのまま読め、`SettingsOf` が無い鍵を既定で埋める。
- 設定の読み：`docflow.SettingsOf[T](o.Node)`。popcornweb の設定は種類の側で `pw.ConfigContext[T](ctx)`（`ParseConfig` の後は背景の goroutine からも読める。`main` が nil の ctx で読んでいる）。芯は popcornweb を import しない。
- 自分の記録（`RegisterCollection[T]`）と定期の仕事（`RegisterTicker`、tick に載る）は fw-4。それまでは `PutState/TakeState` と自分の保存先。

### 2.5 マニュアルと言葉

- `Kind.Doc` を `/help/nodes` の末尾（登録した種類の節）に出す。`Text` は en 必須、ja 任意。
- JS の文は `_()` と `ja.js` のまま（`requirement:docflow-ui-language-decisions`）。独自の種類の文は `Text` で届くので JS の catalog は要らない。

## 3. 設定（popcornweb の構造体）

- そのまま使える：アプリは `pw.RegisterConfig[T]("acme")` を自分のディレクトリで呼び、自分の `popcornweb.toml` の `generate.config` に載せる。`docflow/web` の `Config`（今の `DocflowConfig`）と `OIDCConfig` は `web.RegisterConfig()` が登録し、生成の束縛は `docflow/web` にコミット済みなので、アプリの `pw generate` は docflow の分を作らない。
- **制約**：アプリの `[i18n]` は `docflow/web` と同じ `locales`・`default_locale`（`["en", "ja"]`、`en`）と `cookie_routes = ["/"]` を宣言する。違うと `RegisterLocales` が起動で panic する。3 つ目の言語は先に `docflow/web` に足す。

## 4. カタログ（設定から pulldown）

- `docflow.RegisterCatalog(name string, list func(ctx context.Context) ([]Choice, error))`。`Choice{Value, Label Text, Group}`。`GET /api/catalogs/<name>`（役割 topology）。`Field.Type = catalog` で `Catalog` に名前。`PutNode` は値がカタログにあるかを確かめる（`Open` なら自由入力も可）。
- `docflow/web` が持つカタログ：
  - **`connections`**：`[[docflow.connections]]`（`name`・`url`・`credential`・`note`）。配列の表は TOML だけ、秘密は `${NAME}` で環境変数を指す。s3・smtp・http・s3-trigger・cups-out・printer に `connection` の欄（カタログ `connections`、URL の scheme で絞る）を足し、`url`+`credential` の代わりに**名前で**選ぶ。保存は名前、送る時に解決する（設定の変更がそのまま効く）。芯の欄 `Node.Connection` と `Manager.Connection(name) (URL, credential)`。
  - **`credentials`**：`DOCFLOW_CRED_*` と `credential_dir` のファイルの**名前**（値は出さない）。資格の欄の pulldown。
  - `nodes` は今の picker のまま。
- アプリの設定からのカタログは `RegisterCatalog` の中で `pw.ConfigContext[AcmeConfig](ctx)` を読めば足りる（見本に書く）。
- 画面で接続を足す（保存先に封をして置く）形は後（`requirement:document-flow` の `credentials-store` の問いのまま。問い 4）。

## 5. プリセットのトポロジーとスタンドアロン

- 文書の形（YAML。JSON も同じ木）。鍵・資格の値は含めない（`requirement:docflow` の `yaml-secrets` の暫定）：

```yaml
version: 1
nodes:
  - {name: intake, kind: upload, to: [convert]}
  - {name: convert, kind: convert, to: [out]}
  - name: out
    kind: folder
    settings: {dir: /var/spool/acme, path: "{{meta.$job}}/{{meta.$name}}"}
layout:
  intake: {x: 80, y: 120}
```

- 読み込み：`Manager.ImportTopology(ctx, doc, mode)`。辺の先から順に `PutNode`（閉路は無いので並べられる）。レイアウトも。
- **モード**（問い 6）：`seed`（ノードが 1 つも無いときだけ）、`managed`（毎起動で揃え、プリセットのノードは `Node.Managed` の印で画面も API も読み取り専用。他のノードは自由。違いがあるときだけ書くので、サーバーレスの cold start で毎回は書かない）、`locked`（全部読み取り専用、パレット無し）。既定は `managed`。
- 書き出しと読み込み：`GET /api/topology/export`（YAML）、`POST /api/topology/import?mode=replace|merge`、サブコマンド `docflow topology export|import`（`pw.RegisterSubCommand`）。= `requirement:docflow` の `topology-yaml`。
- スタンドアロン＝プリセット（`locked`）+ 自分の画面 + 入口への送りは自分のコードから（`Manager.Submit` を `web` が公開）か apps。

## 6. 画面（管理・利用者向け）

- **骨**：`docflow/web` の `Shell`・`AppBar` を export（`BindShell`・`BindAppBar`）。アプリの `.pw.html` は自分のパッケージで葉を書き、`web.Page(role, func(w, r, ac web.Account) { pw.WriteHTMLChain(w, r, []pw.HTMLWrapper{web.BindShell(...), web.BindAppBar(...)}, Leaf(...)) })` で包む（popcornweb の「別の骨は WriteHTMLChain で」の形）。
- **制約**：`docflow/web` の画面は `WriteHTMLChain(w, r, nil, leaf)` で、プロジェクトの document で包まれない。アプリは自分の `document.pw.html` を持ってよい。
- **バー**：`web.RegisterNav(NavItem{Path, Title Text, Role, Place})`（`Place` は bar・admin）。`AppBar`・`AdminNav` は `{for}` で表から描く（今の固定のリンクも同じ表に載せる）。
- **ロール**：`web.RegisterRole(name, Title Text)` と `[[oidc.roles]]`（`name`・`people`・`groups`）。`Viewer.Can(name)`。`topology`・`audit` も同じ仕組みに載せ、今の `topology_admins` … は残す。
- **API**：`web.API(role, fn)`（JSON、セッションかトークン）。Viewer は `docflow.ViewerOf(r)`。アプリの API は `/api/<prefix>/…` を自分で選ぶ（`/api/<kind>/<name>` とぶつからない名前。`web` が確かめる）。
- **公開の木**（問い 5）：pw の `/public` はアプリのもの（アプリの `[assets]` と manifest）。`docflow/web` は自分の JS・CSS を **`/_docflow/`（名前は暫定）で自分で配る**（embed、ETag、`no-cache`）。理由は `RegisterPublicManifest` が 1 つの表だけなこと。参照アプリ `cmd/docflow` も同じ道にし、bdf の `popcornweb.toml` から `[assets]` を外し、`dist/public` と manifest を消す。開発は `web.Options.AssetsDir`（か `DOCFLOW_ASSETS_DIR`）でディスクから読む。
- **CSS**：docflow の `app.css` は `docflow/web` が持ち、`Shell` が読む。アプリの画面の追加のクラスは、アプリの Tailwind で自分の CSS を作り `web.Options.Stylesheets` で `Shell` に足す。1 つにまとめたいなら `@source` に `docflow/web/*.pw.html` を足して自分で作り、`Options.Stylesheet` で置き換える。
- **help**：`web.RegisterHelp(chapter, lang→Markdown)`（目次の末尾）。`Kind.Doc` は nodes の章へ。

## 7. 試験と組み立て

- bdf の `popcornweb.toml` は `[generate]` と `[i18n].catalog` を `docflow/web` に向け、`[assets]` を外す。CI の `pw generate` の差分の確認も `docflow/web` に移る。
- `docflow/flowtest` を公開する（`Env.Put`・`Submit`・`Tick`・`Delivery`・`Required`・`Payload`）。独自の種類の試験はこれで書く。
- `examples/docflow-custom`：`go.mod`（`replace ../..`）、`popcornweb.toml`、`main.go`、`kinds/folder`（出口：ディレクトリに書く。`Fields` は `dir`・`path`（`{{meta.<key>}}`）・`overwrite`）、`[acme]` と `[[docflow.connections]]` の設定、`presets/topology.yaml`、`pages`（`/acme/dashboard` の管理画面と `/acme/portal` の利用者向け）。CI で `go build`・`go test`・`pw generate` の差分なしを見る。

## 8. マイルストーン

| 段 | 中身 | 受け入れ |
|---|---|---|
| **fw-1 独自の種類** | `Node.Settings`、`Kind.Titles/Help/Fields/Icon/Group/Doc`、`KindInfo` の拡張、汎用のフォーム、API からのパレット・アイコン・役割の表、`docflow/web` への移動（`Main`/`Run`/`Options` の最小）、`/_docflow/` の自前の配信、`flowtest` の公開、`examples/docflow-custom` の `folder` | 外のモジュールから `pw build` が通る。パレットに folder が出て、置いて設定して流すとディレクトリに書かれる。見本の `go test` が flowtest で通る |
| **fw-2 カタログとプリセット** | `RegisterCatalog`、`connections`・`credentials`、`connection` の欄、プリセットの 3 モード、export/import、サブコマンド | s3 ノードを接続の名前で作れる。プリセット付きのバイナリが空の保存先で起動してトポロジーが揃い、managed のノードは変えられない |
| **fw-3 画面** | `BindShell`/`BindAppBar`、`Page`/`API`、`RegisterNav`、`RegisterRole`、`Stylesheets`、`RegisterHelp`、見本の 2 画面 | 見本の管理画面がバーに出て役割で守られ、利用者向けの画面が自分の文書で描ける |
| fw-4（後） | `Kind.Pick`、`RegisterCollection[T]`、`RegisterTicker`、独自の口と at-node の節 | — |

fw-2 と fw-3 は fw-1 の後なら独立。

## 9. 問い（番号つき。推しと結果）

必須（fw-1 の形を決める）：

1. **独自の種類の設定の持ち方**：A. `settings` のオブジェクト（推し：種類の欄と芯の欄を分け、YAML・OpenAPI・フォームが素直）／B. 平らな `Node` に `extra` の map（名前空間が混ざる）。
3. **フレームワークの置き場**：A. bdf の `docflow/web`（推し：今の go.mod・CI・生成の束縛をそのまま使え、切り出しは後からでも機械的）／B. 別リポジトリ（bdf を依存に。版の同期が要る）。

任意（既定で進める）：

2. 今のアダプタを `settings` に移すか：既定は fw-1 では移さない（保存済みのノードの互換を保つ）。移すなら互換なし。
4. 接続（`connections`）の持ち場所：A. 設定 TOML（推し：運用の持ち物、秘密は環境変数の参照）／B. 保存先と画面（封をした秘密が要る）／C. 両方。
5. docflow の資産の配り方：A. `/_docflow/` を自分で（推し）／B. アプリの public に写してから `pw generate`（手順が増える）。
6. プリセットの既定のモード：`managed`（推し）／`seed`／`locked`。
7. 見本の置き場：A. `examples/docflow-custom`（入れ子のモジュール、推し）／B. 別リポジトリ。
