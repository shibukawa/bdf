# 薄い保存層（store）— bbolt・オブジェクトストレージ・DynamoDB・RDB を同じ約束で

2026-10-19 時点の案。**後に回した**（同日の利用者「store レイヤーは一旦頑張らずに bbolt 前提で良い」）。print-m2 のスプールは bbolt を直接使い、bbolt に触るのを `print/spool/db.go` だけにしておく（`print-m2.ja.md` §0、§6）。保存先を替える・複数台で回す・Lambda で動かすときに、この案で `db.go` を置き換える。調べた結果（§1、§4.1）はそのときの出発点。

案の中身: スプール（と BDF のほかのサーバー機能）が使う保存先を、1 つの小さな約束の上に置き、汎用の層として `store/` に置く。

## 0. 問いと答え

利用者の問い（2026-10-19）: bbolt、DynamoDB、その他の KVS（と RDB）を包む仕組みは無いか。SDK を同梱すると重いので、HTTP の API を呼ぶ薄い仕組みがよい。外部キーなしの insert・update・select（更新順の並びだけ）・delete ができれば多くの保存先に対応できるはず。楽観ロックのオブジェクトストレージに対応できれば、費用も運用の楽さもだいぶよくなる。

答え。

- 既にある Go の仕組みは、どれも「SDK を同梱する」か「楽観ロックか更新順の一覧が無い」かのどちらかで、そのままは使えない（§1）。
- 約束は問いのとおり小さくできる。操作は **Get・Insert・Update・Delete・Scan（更新順）** と GC 用の Purge、それに大きなバイト列（文書）の **Blob** だけ。楽観ロックは「版（version）が同じときだけ書く」で、Insert は「無いときだけ書く」（§2）。
- この約束は、bbolt、S3 互換のオブジェクトストレージ、GCS、Azure Blob、DynamoDB、RDB のどれにも素直に写る（§3）。
- **オブジェクトストレージは楽観ロックに対応できる**。AWS S3 は 2024 年に `If-None-Match`（8 月）と `If-Match`（11 月）の条件付き PUT を、2025 年 9 月に条件付き DELETE を足した。GCS（generation）と Azure Blob（ETag）は前から持つ。R2 や Scaleway なども対応している。ただし対応しない・黙って無視する保存先もある（Backblaze B2 は拒むとの報告）ので、起動時に実際に確かめてから使う（§4）。
- SDK は入れない。AWS の署名（SigV4）、資格情報の取得、各サービスの REST を自前で薄く書く（§5）。依存は増えない。
- print-m2 では作らない（bbolt を直接使う）。作るときの順は §7。

## 1. 既にある仕組み

| 仕組み | 何ができる | 足りないもの |
|---|---|---|
| Go CDK（`gocloud.dev`）の `blob` と `docstore` | `blob` は S3・GCS・Azure・ファイル・メモリ。`docstore` は DynamoDB・Firestore・MongoDB・メモリで、版（`DocstoreRevision`）による楽観ロックがある | 各社の公式 SDK を同梱する。`blob` の条件付き書き込みは「無いときだけ」（`IfNotExist`）で、`If-Match` での置き換えは無い（削除は各ドライバの前提条件を渡す口だけ）。`docstore` はオブジェクトストレージと RDB を持たない |
| `philippgille/gokv` | Set・Get・Delete の 1 つの interface に bbolt、DynamoDB、S3、PostgreSQL、MySQL、Redis など 20 余り | 版も比較交換も一覧も無い（一覧は「予定」）。借りやキューが作れない |
| `valkeyrie`（docker/libkv の後継） | Consul・etcd・ZooKeeper・Redis・BoltDB・DynamoDB に AtomicPut・AtomicDelete（比較交換） | オブジェクトストレージと RDB が無い。更新順の一覧が無い。SDK を同梱 |
| `k3s-io/kine` | etcd の API を SQLite・PostgreSQL・MySQL・NATS の上に作る。版つきの比較交換と、版の順の監視 | etcd のサーバーとして動く重い仕組み。考え方（1 つの表と版の順）は参考になる |
| AWS SDK for Go v2、`minio-go` などの単体の SDK | 1 つのサービスを完全に話す | 1 つのサービスだけ。束ねる層は別に要る |

結論: 約束を自分で決め、各保存先を HTTP で薄く話すドライバを書く。各ドライバは数百行で済む（§5）。

## 2. 約束

```go
package store

// Table は 1 つの表。キーと値（不透明なバイト列）の行を持つ。外部キー、表をまたぐトランザクション、
// 副索引、更新順以外の問い合わせは無い。
type Table interface {
	Get(ctx context.Context, key string) (Row, error)                      // 無ければ ErrNotFound
	Insert(ctx context.Context, key string, value []byte) (Row, error)     // 無いときだけ書く。あれば ErrExists
	Update(ctx context.Context, key string, value []byte, ver string) (Row, error) // 版が同じときだけ。違えば ErrConflict
	Delete(ctx context.Context, key string, ver string) error              // 版が同じときだけ、墓標にする
	Scan(ctx context.Context, since Cursor, limit int) ([]Row, Cursor, error) // 更新順。墓標を含む
	Purge(ctx context.Context, key string) error                           // 墓標を本当に消す（GC だけが使う）
}

type Row struct {
	Key     string
	Value   []byte
	Ver     string    // 不透明。ETag、generation、連番
	Updated time.Time // 保存先が付けた更新の時刻（並びに使う）
	Deleted bool      // 墓標
}

// Blobs は書いた後に変えない大きなバイト列（文書）。
type Blobs interface {
	Put(ctx context.Context, key string, r io.Reader, size int64) error   // 新しいキーにだけ書く
	Open(ctx context.Context, key string) (Blob, error)                   // ReadAt（範囲の読み取り）と Size
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) iter.Seq2[string, error]     // GC が持ち主の無い Blob を探す
}
```

約束の細目。

- **キー**は `[a-z0-9._/-]` の 1〜200 文字。**値**は 256 KiB まで（DynamoDB の 1 項目 400 KB に収まる）。大きいものは Blob に置く。
- **1 つのキーの上では線形化可能**: Insert・Update・Delete の条件は保存先が原子的に判定し、条件付きの書き込みが成功した後の Get はそれを返す。
- **Scan は遅れてよい**: `since` より後に変わった行を（更新の時刻、キー）の順に返す。保存先によっては直近の変更が少し遅れて現れる（DynamoDB の GSI、同じ秒の LastModified）。使う側は `since` を少し戻して読み直し（重なりの窓。既定 30 秒）、`(Key, Ver)` で重複を捨てる。墓標は `Purge` まで残るので、窓の中の削除も見える。窓より遅れた読み手（長く止まっていた）は全件を読み直す。
- **Delete は墓標**: 他の読み手が Scan で削除を知るため。GC が十分古い墓標を `Purge` する。
- **版**は保存先が決める不透明な文字列。使う側は比べるだけで、中身を解釈しない。
- **競合**: 条件が外れた（HTTP 412）と、同時の条件付き書き込みにぶつかった（409）のどちらも `ErrConflict`。使う側は読み直してやり直す。

使う側（スプール）が約束の外で受け持つこと。

- **並べ替えと絞り込み**: スプールは表を Scan で読んでメモリに持ち、キューの順、時刻の来たジョブ、利用者ごとの一覧はメモリで作る。1 台のサーバーで数万行なので足りる。
- **1 行で閉じる**: 1 つのジョブは 1 行。表をまたぐ変更（終わったジョブを `jobs` から `done` に移す）は「`done` に Insert → `jobs` を Delete」の順にし、途中で落ちても次の起動で同じ順にやり直せば揃う。
- **連番**（IPP の `job-id`）: `meta` 表の 1 行を Update（版つき）で 1 つずつ足す。競合したら読み直す。
- **借り**: ジョブの行の中に持ち主と期限を書き、Update（版つき）で取る（`print-m2.ja.md` §3.4）。

## 3. 各保存先への写し方

| 操作 | bbolt（ファイル） | S3 互換 | GCS | Azure Blob | DynamoDB | RDB |
|---|---|---|---|---|---|---|
| 行の置き方 | 表ごとのバケット | `<prefix>/<表>/<キー>` のオブジェクト | 同左 | 同左 | 1 つの表。パーティションキー `表#キー` | 1 つの表 `bdf_rows(tbl, k, ver, upd, del, val)`、主キー `(tbl, k)` |
| Get | 読み取りのトランザクション | GET（ETag） | GET（generation） | GET（ETag） | GetItem（`ConsistentRead`） | SELECT |
| Insert | 書き込みのトランザクションで有無を見て書く | PUT `If-None-Match: *` | PUT `x-goog-if-generation-match: 0` | PUT `If-None-Match: *` | PutItem `attribute_not_exists(k)` | INSERT（主キーの重複で失敗） |
| Update | 版を比べて書く | PUT `If-Match: <ETag>` | PUT `x-goog-if-generation-match: <gen>` | PUT `If-Match: <ETag>` | PutItem `ver = :v` | `UPDATE … WHERE tbl=? AND k=? AND ver=?`（0 行なら競合） |
| Delete | 墓標 | 墓標を PUT `If-Match` | 墓標を PUT（generation） | 墓標を PUT `If-Match` | 墓標を PutItem `ver = :v` | `UPDATE … SET del=1 …` |
| Scan | 更新の連番の索引バケット | ListObjectsV2 で一覧（キー、ETag、LastModified）→ 時刻で並べ、変わったものだけ GET | objects.list（`updated`、`generation`） | List Blobs（`Last-Modified`、`Etag`） | GSI（パーティション `tbl`、ソート `upd`）を Query | `WHERE tbl=? AND upd >= ? ORDER BY upd, k LIMIT ?`、索引 `(tbl, upd, k)` |
| Purge | 削除 | DELETE | DELETE | DELETE | DeleteItem | DELETE |
| 版 | 行ごとの連番 | ETag | generation | ETag | 数の属性 `ver` | 数の列 `ver` |
| 更新の時刻 | プロセスの単調な連番 | LastModified（秒の精度） | `updated`（ミリ秒） | `Last-Modified`（秒） | 書き手の時計（µs） | 書き手の時計（µs） |
| Blob | 1 MiB の塊のバケット | オブジェクト（大きいものは multipart） | オブジェクト | ブロック BLOB | 持てない（S3 と組む） | 塊の行の表 `bdf_blobs(k, n, data)` |

- 行と Blob は別の保存先にできる（DynamoDB の行と S3 の Blob、RDB の行と S3 の Blob）。既定は同じ保存先。
- RDB の表は 1 つ（`bdf_rows`）で、論理的な表は列 `tbl` で分ける。外部キーも、トリガーも、方言ごとの機能も使わない（型の名前と `LIMIT` / `TOP` だけが方言）。DBA に渡すのは CREATE 文 2 つ。
- 書き手の時計を使う保存先（DynamoDB、RDB）では、台ごとの時計のずれを重なりの窓（30 秒）が吸収する。

## 4. オブジェクトストレージ

### 4.1 条件付き書き込みの対応

| 保存先 | 無いときだけ PUT | 版が同じときだけ PUT | 条件付き DELETE | 備考 |
|---|---|---|---|---|
| AWS S3 | ○ `If-None-Match: *`（2024-08） | ○ `If-Match`（2024-11） | ○ `If-Match`（2025-09） | 署名（SigV4）か HTTPS が要る。同時の書き込みは 409 `ConditionalRequestConflict` もある。LIST も含めて強い一貫性（2020-12） |
| Google Cloud Storage | ○ `x-goog-if-generation-match: 0` | ○ generation | ○ generation | 前から。一覧も強い一貫性。S3 互換の XML API でなく専用のドライバにする |
| Azure Blob Storage | ○ `If-None-Match: *` | ○ `If-Match`（ETag） | ○ `If-Match` | 前から。外れると 412 |
| Cloudflare R2 | ○ | ○ | 要確認 | 同時に `If-None-Match: *` を送ると成功は 1 つだけ、と文書にある |
| Scaleway、OVHcloud、CoreWeave など | ○ | ○ | 要確認 | 各社の文書に `If-None-Match: *` と `If-Match` |
| MinIO | ○ | ○ | ○ | 2023 年から。ただしコミュニティ版は 2025 年に保守のみになり 2026-02 に保存庫が凍結。2026-03 に「読み取りの定足数が無いと条件を確かめずに書く」不具合の報告。手元の試験用にとどめる |
| Backblaze B2 | ×（`If-None-Match` を拒む、2025-07 の報告） | 不明 | 不明 | 使えない |
| Wasabi、さくらのクラウドなど | 不明 | 不明 | 不明 | 起動時の確かめで判断 |

### 4.2 確かめてから使う

ヘッダーを黙って無視する保存先では「最後に書いた者が勝つ」になり、借りが効かず、二重に刷る。だから `s3` ドライバは起動時（と `bdf printd -check`）に、使い捨てのキー（`<prefix>/_probe/<乱数>`）で次を確かめ、1 つでも外れたら起動しない。

1. 無いキーへの `If-None-Match: *` の PUT が成功する。
2. 同じキーへの 2 回目の `If-None-Match: *` が 412 か 409 で断られる。
3. 古い ETag の `If-Match` の PUT が断られ、今の ETag なら成功する。
4. 書いた直後の LIST に、そのキーと新しい ETag が出る（一覧の一貫性）。
5. 条件付き DELETE（古い ETag で断られる）。対応しなければ Purge を条件なしの DELETE にする（墓標のキーは使い回さないので安全）と記録する。

### 4.3 一覧と費用

- Scan は ListObjectsV2 で表の一覧（1 回 1000 件）を取り、LastModified と ETag で変わった行だけを GET する。一覧の費用は表の件数に比例するので、生きている表を小さく保つ。スプールは終わったジョブを `jobs` から `done` に移し、見回りで一覧するのは `jobs` だけ。`done`（実績）は API で問われたときだけ読む。
- LastModified は秒の精度なので、同じ秒の変更は重なりの窓と ETag の比較で拾う。
- 見回りの間隔は、1 台のサーバーなら要らない（同じプロセスの変更はメモリで通知する。安全網として 5 分ごとに読み直す）。複数台なら 5〜10 秒（後の段階）。
- 書き込みを減らす: ジョブの行に書くのは状態が変わったときと借りを延ばすときだけ。見届けの次の時刻はメモリに持ち、行に書かない。
- **費用の目安**（S3 Standard の公表価格での概算。1 日 1000 ジョブ、1 ジョブ 1 MB、文書を 7 日残す）: 書き込みは 1 ジョブ十数回で月 40 万回強、月 2 ドル前後。保存は 7 GB 程度で月 0.2 ドル前後。複数台で 5 秒ごとに見回ると一覧が月 50 万回増えて、さらに 2〜3 ドル。サーバーも容量の計画もバックアップも要らない。RDB は最小の管理型インスタンスでも月十数ドルからで、運用が付く。

### 4.4 文書の置き方

- 文書はジョブごとに `docs/<job-id>/original`（受けたまま）と `docs/<job-id>/bdf`（変換した .bdf）に置く。ジョブをまたいで共有しない。内容アドレスで共有すると、複数の持ち主が同時に「参照が 0 になった Blob を消す」と「同じ Blob を参照する新しいジョブを作る」をしたときに、消した Blob を指すジョブができる。同じ文書を別のジョブで送ったときに容量が重なるだけで、正しさを取る。
- 記録に SHA-256 と大きさを持ち、読むときに確かめる。
- 範囲の GET で `io.ReaderAt` にする（1 MiB ずつ、直前の塊を覚える）。

## 5. SDK を入れない

- **AWS**（S3、DynamoDB）: Signature Version 4 を自前で書く（HMAC-SHA256、正規化した要求の署名。S3 の本文は HTTPS の上で `UNSIGNED-PAYLOAD`）。資格情報は次の順に探す。
  1. 環境変数（`AWS_ACCESS_KEY_ID`、`AWS_SECRET_ACCESS_KEY`、`AWS_SESSION_TOKEN`。Lambda もこれ）
  2. 共有の設定の `credential_process`（外の道具に任せる。AWS IAM Identity Center（SSO）は `aws configure export-credentials` をここから呼ぶ）
  3. 共有の資格情報ファイルの静的なキー（`AWS_PROFILE`）
  4. Web identity（`AWS_WEB_IDENTITY_TOKEN_FILE` と `AWS_ROLE_ARN`。EKS の IRSA。STS の `AssumeRoleWithWebIdentity` は署名の要らない HTTP）
  5. コンテナの資格情報（`AWS_CONTAINER_CREDENTIALS_FULL_URI` / `_RELATIVE_URI`。ECS、EKS Pod Identity）
  6. EC2 のインスタンスメタデータ（IMDSv2）
- **S3 互換**: 道の形（path-style）と仮想ホスト形、`endpoint` の指定で R2、Scaleway、MinIO などに向ける。PUT・GET（範囲）・HEAD・DELETE・ListObjectsV2（XML）・multipart（大きな Blob）。
- **GCS**: JSON API。OAuth2 のトークンはメタデータサーバー（GCE、Cloud Run、GKE）か、サービスアカウントの鍵（`GOOGLE_APPLICATION_CREDENTIALS`。JWT を標準ライブラリの RS256 で署名してトークンに替える）。
- **Azure Blob**: REST。SAS、共有キー（HMAC）、マネージド ID（IMDS）、クライアントシークレット。
- **DynamoDB**: JSON の API（`X-Amz-Target` ヘッダー）＋SigV4。表と GSI は利用者が作る（docs に CLI と CloudFormation の例）か、`-init` で CreateTable。
- **RDB**: HTTP ではなく各 DB の通信規約なので、どこかのドライバが要る。案は 2 つ（§8 の論点）。(a) `database/sql` と、ビルドタグで選ぶ純 Go のドライバ（PostgreSQL は pgx、MySQL、SQLite）。(b) HTTP で SQL を受けるもの（Amazon RDS Data API は SigV4 の HTTP、PostgREST、libSQL / Turso、Neon）を薄く話す。
- **規模の目安**: 署名と資格情報で 1 千行弱、各ドライバ 300〜600 行。`go.mod` に足すものは無い（bbolt は print-m2 で別に足す）。

## 6. 置き場所と選び方

```
store/              Table、Blobs、Row、Cursor、誤り（ErrNotFound、ErrExists、ErrConflict）、Open（URL で選ぶ）
store/storetest/    約束のテスト（どのドライバにも同じものを流す）
store/mem/          メモリ（テスト、2 つの持ち主の取り合いのテスト）
store/bolt/         bbolt の 1 ファイル（既定。詰め直し、写し）
store/s3/           S3 互換（条件付き書き込みの確かめ、multipart）
store/s3/s3test/    S3 の偽物（httptest。ETag、If-Match、If-None-Match、409、LIST）
store/gcs/          後
store/azblob/       後
store/dynamodb/     後
store/sql/          後（§8）
store/internal/awsauth/  SigV4 と資格情報の順
```

選ぶのは URL。

```yaml
store:
  records: file:///var/lib/bdf/spool.db                          # 既定は設定ファイルと同じディレクトリの spool.db
  # records: s3://my-bucket/bdf?region=ap-northeast-1
  # records: s3://my-bucket/bdf?endpoint=https://<account>.r2.cloudflarestorage.com&region=auto
  # records: gs://my-bucket/bdf           azblob://account/container/bdf
  # records: dynamodb://bdf-rows?region=ap-northeast-1&blobs=s3://my-bucket/bdf
  # records: postgres://…
  blobs: ""                                                      # 空なら records と同じ保存先
  overlap: 30s                                                   # Scan の重なりの窓
```

## 7. 段階

print-m2 では作らない（2026-10-19 利用者）。作るときの案。

1. `store`（約束、`storetest`）、`mem`、`bolt`（print-m2 の `db.go` をこの約束に載せ替える）、`s3`（S3 互換。AWS S3 と R2 を手で確かめる）、`awsauth`。スプールのジョブの記録に持ち主と期限の「借り」を足し、2 つの持ち主の取り合いを `mem` と `s3test` のテストで確かめる。
2. 続けて `gcs`、`azblob`、`dynamodb`、`sql`。複数台で 1 つの保存先を回すこと（プリンタごとの借り、台をまたぐ通知、パスワードの封印）と Lambda。

## 8. 論点

1. **パッケージの名前と置き場所**: 汎用の `store/` にする（推奨。印刷以外のサーバー機能も使える）か、`print/spool/store` に閉じるか。
2. **いつ作るか**: m2 の後に決める（`requirement:print` の `scale-out`）。
3. **文書をジョブごとに置き、内容アドレスで共有しない**（§4.4）。print-m2 の bbolt でも同じ置き方にしてある。
4. **RDB の話し方**: `database/sql`＋ビルドタグのドライバか、HTTP で SQL を受けるもの（RDS Data API、PostgREST、libSQL）か、両方か。
5. **次に足すドライバの順**: GCS・Azure Blob（オブジェクトストレージをそろえる）か、DynamoDB か、RDB か。

## 9. yuruidb の評価（2026-10-20）

利用者が別に作った `github.com/shibukawa/yuruidb`（SDK なしの HTTP ドライバで bbolt・SQL・S3・R2・DynamoDB・GCS・Firestore・Realtime Database・D1・Cloudflare KV を同じ CRUD に載せ、楽観ロックで競合を検出するライブラリ。コミット 276310b）を読み、テストを走らせ、手元で確かめた。

よいところ: 共通部は標準ライブラリだけ。記録の版（乱数）と保存先の版（ETag、generation など）の 2 段で、同じ内容を書き戻しても古い版の書き込みが通らない（ABA に強い）。CAS を持たない Cloudflare KV は条件付きの操作を断り、無条件の Set・Remove を別に明示する。送った後の通信断や 5xx を `ErrOutcomeUnknown` として区別する。共通の適合テスト、条件の無い書き込みを断る模擬サーバー、別プロセスの競合、SigV4 の公式の既知ベクトル。Go 1.26 で `go test -race ./...` と `go vet` が通る。

bdf のスプールに使う前に要るもの（手元で確かめた）。

1. **1 つの bbolt のファイルに 1 つのコレクションしか開けない**: 同じプロセスで同じパスの 2 つ目を開くと、ファイルのロックで 1 秒後に timeout になる。スプールは jobs・done・meta・printers・trust を 1 つの `spool.db` に持ちたい。パスごとに `*bolt.DB` を共有する必要がある。
2. **ID が UUIDv7 に限られる**: プリンタの名前、ホスト、連番の行のような名前のキーが置けない（IPP の `job-id` は整数）。
3. **墓標が消えない**: Delete は墓標を書き、物理的には消さない。生きている 10 件と消した 1990 件のコレクションを List すると、墓標ごとに 1 回、計 2000 回の読み取りになった（オブジェクトストレージでは GET）。終わったジョブを消すスプールでは、一覧が時間とともに重くなる。S3（2025-09 から条件付き DELETE）、GCS、DynamoDB、SQL、bbolt は条件付きの削除ができ、記録の版が乱数なので消して作り直しても ABA は起きない。
4. **`ErrOutcomeUnknown` を呼び手が照合できない**: 書けたのに応答が失われた Update の後、古い版で再試行すると競合になり、自分の書き込みかどうか分からない。記録の版は書き手が作るので、ライブラリが読み直して照合できる。
5. 一覧は Keys の後に 1 件ずつ読む（N+1）。Insert と Update は読んでから書く（2 往復）。bbolt では値が 2 重に base64 になり、同じ値の生の bbolt の約 2 倍のファイルになった。
6. AWS の資格情報は環境変数だけ（Lambda は動く。EC2・ECS・EKS は注入が要る）。S3 は path-style だけで、地域は既定の us-east-1（`AWS_REGION` を読まない）。一時的な 5xx・503 の再試行が無い。S3 互換のサービスが条件のヘッダーを本当に守るかの確かめが無い。大きな値（文書）を置く Blob の口が無い。

結論: 設計の芯（2 段の版、能力の明示、送った後の不明の区別）は bdf の案（§2）とよく合う。1〜3 が直れば、print-m2 の `db.go` の記録の部分を yuruidb に載せ替えられる。文書は別に置く（bbolt の塊か、オブジェクトストレージに直接）。

### 9.1 再評価（2026-10-20、コミット adb806b）

利用者が指摘を受けて更新した（Go 1.27 のメソッドの型引数で 1 つの `Store` に複数のコレクションを `Bind[T]("name")` で登録する形、名前のキー、条件付きの物理削除、書き込みの結果の照合、bbolt のバイナリ形式、AWS の資格情報の順と地域、virtual-hosted、一時的な誤りの再試行、S3 互換の確かめ、S3・R2・GCS の Blob）。両方の module で `go test -race` と `go vet` が通り、手元で次を確かめた。

- 1 つの bbolt のファイルに jobs・done・printers を同時に開け、同じファイルの 2 つ目の `Store` も開け、片方を閉じてももう片方が使える。
- 名前のキー（`0000000042`、`mfp1.example.local`）が置ける。使える文字は英数字と `.`・`_`・`-` なので、IPv6 のホスト（`:` を含む）やプリンタの任意の名前は bdf の側で符号化する。
- 2000 件作って 1990 件消した後の一覧は、キーの一覧 1 回と生きている 10 件の読み取りだけ（bbolt は物理削除）。Insert と Update は書き込み 1 回だけで、先読みが無い。
- bbolt の使用量は、10 KB の値で生の bbolt と同じ（1.00 倍）、300 B の値で 1.25 倍（ID と版の見出しの分）。
- 書けたのに応答が失われた Update は、ライブラリが読み直して成功を返す。DynamoDB で最初の PutItem が適用されて 500 になり再試行が条件で断られる場合も、正しく成功になった。

残る点。

1. 応答が失われ、照合の読み取りも失敗したときの誤りが `ErrOutcomeUnknown` と同時に `ErrConflict` にも当たる。「競合なら読み直してやり直す」と書いた呼び手は、自分の書き込みをもう一度当てうる。不明のときは `ErrConflict` に当てないか、先に `ErrOutcomeUnknown` を見るよう文書に書く。
2. S3 互換の確かめの使い捨てキーがコレクションの中（`<collection>/__yuruidb_probe_…`）に置かれる。条件付き DELETE の無い保存先で後片付けが失敗すると（権限、通信、取り消した context）、そのコレクションの一覧が「corrupt record」で止まった。コレクション名に使えない文字の接頭辞に置き、後片付けは `context.WithoutCancel` で行う。
3. Blob は S3・R2・GCS だけで、bbolt の Blob が無い。範囲の読み取り（`io.ReaderAt`）と一覧（持ち主の無い Blob の掃除）も無い。条件付き DELETE の無い保存先（確かめの結果による）では Blob を消せない。
4. bbolt の `*bolt.DB` を外に出さないので、bdf の写し（`/backup`）や詰め直しを動かしたまま行えず、同じファイルに bdf の独自のバケットも置けない。
5. DynamoDB の流量制限（HTTP 400 の ThrottlingException など）は再試行しない。

結論: 前回の 1〜3（複数のコレクション、名前のキー、墓標）は解消した。ジョブの記録は yuruidb の bbolt ドライバに載せられる。bdf の既定（bbolt の 1 ファイル）で使い切るには、3 と 4（bbolt の Blob、写しと詰め直し）が要る。それまでは文書だけを別のファイルに置く形でも始められる。

### 9.2 再評価（2026-10-20、コミット 3684337）

利用者がさらに更新した（不明の誤りを `ErrConflict`・`ErrAlreadyExists` に当てない、確かめの一時キーをコレクションの外の `.yuruidb-probes/` に置き後片付けを `context.WithoutCancel` で行う、bbolt・SQL の Blob、Blob の `io.ReaderAt`・一覧・無条件の `Remove`、bbolt の稼働中の `Backup`、DynamoDB の HTTP 400 の流量制限の再試行）。両方の module で `go test -race` と `go vet` が通り、手元で次を確かめた。

- 応答が失われ照合の読み取りも失敗したときの誤りは `ErrOutcomeUnknown` だけに当たり、`ErrConflict` には当たらない。
- 確かめの後片付けに失敗しても、コレクションの一覧は止まらない。
- 前回までの改善（複数のコレクション、名前のキー、物理削除、書き込み 1 回、照合）はそのまま。

残る点（確かめた）。

1. **`Store.Backup` の間、同じ `Store` のレコードの操作が止まる**: `Backup` が `Store` の読み取りロックを書き出しの間ずっと持ち、各操作のコレクションの解決が書き込みロックを取るため。遅い書き出し先への 0.55 秒のバックアップの間、Get が 0.55 秒待った。bdf の `GET /api/v1/backup` は遅い HTTP の相手に流すので、その間スプールが止まる。ロックは確認だけに使い、書き出しの間は放す（bbolt の Close は開いているトランザクションを待つ）。
2. **bbolt の Blob の Insert が入力を読み終えるまで bbolt の書き込みロックを持つ**: 1 MiB を 1.6 秒かけて受ける間、別のレコードの Update が 1.67 秒待った。入力全体もコミットまでメモリに残る。README に書いてはある。bdf は受けた文書を一時ファイルに落としてから入れれば避けられるが、yuruidb の側で数 MiB ごとの書き込みに分けて最後に見えるようにする（作業用のバケットから移す）と、ロックとメモリの両方が抑えられる。
3. HTTP の Blob の `Open` は本文の GET を始めて開いたままにするので、`ReadAt` だけを使う呼び手には無駄（小さい）。
4. bdf の側で扱うこと: コレクションをまたぐトランザクションは無い（jobs から done への移しは途中で落ちても揃え直せる順に書く）。ID の文字種に入らないプリンタの名前や IPv6 のホストは符号化する。

結論: bdf のスプールに要るもの（1 つの bbolt のファイルに記録・文書・バックアップ）が揃った。1 が直れば print-m2 は yuruidb の上に作れる。2 は bdf の側の一時ファイルでも避けられる。

### 9.3 再評価（2026-10-20、コミット b6e1529）

利用者が前回の 2 点を直した（`Store.Backup` は確認の後にロックを放す、bbolt の Blob の Insert はトランザクションの外で入力を読み、2 MiB ずつ作業用のバケットに書いて最後に移して見せる、途中で落ちた作業用のデータは次にファイルを最初に開いたときに消す、HTTP の Blob の `Open` は HEAD などで版と大きさだけを取り本文は最初の Read で取りに行く）。両方の module で `go test -race` と `go vet` が通り、手元で次を確かめた。

- `Backup` の間の Get は待たない（前回 0.55 秒 → 0 秒）。
- 遅い Blob の受け付けの間のレコードの Update は待たない（前回 1.67 秒 → 2 ミリ秒）。
- 途中で失敗した受け付けは何も残さず、同じ ID でもう一度入れられる。128 MiB の受け付けでヒープの増加は 12 MiB ほどで、100 MiB の位置の `ReadAt` も読める。

残る点: `Backup` の間の Update が、bbolt がファイルの対応づけ（mmap）を広げる必要があるときだけ、バックアップが終わるまで待つ（0.55 秒）。bbolt の性質で、長い読み取りのトランザクションの間は対応づけを広げられない。生の bbolt で確かめると、`InitialMmapSize` を 256 MiB にすれば同じ書き込みは 1 ミリ秒、既定では 515 ミリ秒だった。yuruidb の bbolt の Config に `InitialMmapSize` を出してもらうのがよい（Windows の bbolt はファイルを対応づけの大きさまで伸ばすので、既定を大きくはしない）。bdf の側では `GET /api/v1/backup` をいったんローカルの一時ファイルに書いてから相手に流せば、読み取りのトランザクションはディスクの速さで終わる。

結論: print-m2 を yuruidb の上に作れる（記録は bbolt のドライバ、文書は同じ `spool.db` の bbolt の Blob、写しは `Store.Backup`）。

## 10. 出典

- AWS S3 の条件付き書き込み（If-Match、2024-11）: https://aws.amazon.com/about-aws/whats-new/2024/11/amazon-s3-functionality-conditional-writes
- AWS S3 の条件付き書き込み（If-None-Match、2024-08）: https://www.amazonaws.cn/en/new/2024/amazon-s3-supports-conditional-writes/
- AWS S3 の条件付き削除（2025-09）: https://www.amazonaws.cn/en/new/2025/amazon-s3-supports-conditional-deletes-in-s3-general-purpose-buckets/
- AWS S3 の条件付き書き込みの手引き: https://docs.aws.amazon.com/AmazonS3/latest/userguide/conditional-writes.html
- AWS S3 の強い一貫性（LIST を含む、2020-12）: https://aws.amazon.com/about-aws/whats-new/2020/12/amazon-s3-now-delivers-strong-read-after-write-consistency-automatically-for-all-applications
- GCS の前提条件（generation、0 は「無いときだけ」）: https://docs.cloud.google.com/storage/docs/request-preconditions
- Azure Blob の並行性の管理（ETag と If-Match）: https://learn.microsoft.com/azure/storage/blobs/concurrency-manage
- Cloudflare R2 の S3 API の互換: https://developers.cloudflare.com/r2/api/s3/api/ （作業環境から開けず、検索の抜粋による）
- Scaleway の条件付き書き込み: https://www.scaleway.com/en/docs/object-storage/api-cli/using-conditional-writes/
- OVHcloud の条件付き書き込み: https://docs.ovhcloud.com/es/guides/storage-and-backup/object-storage/s3-conditional-writes.md
- CoreWeave の条件付き書き込み（409）: https://docs.coreweave.com/changelog/release-notes/ai-object-storage-conditional-writes
- MinIO の条件付き書き込み: https://www.min.io/blog/leading-the-way-minios-conditional-write-feature-for-modern-data-workloads 、保守のみへの移行（二次資料）: https://ayedo.de/en/posts/minio-im-maintenance-mode/
- Backblaze B2 が If-None-Match を拒む報告（Proxmox、2025-07）: https://lists.proxmox.com/pipermail/pbs-devel/2025-July/014234.html
- Go CDK の blob（`IfNotExist`）: https://pkg.go.dev/gocloud.dev/blob 、docstore（版による楽観ロック）: https://pkg.go.dev/gocloud.dev/docstore
- gokv: https://pkg.go.dev/github.com/philippgille/gokv
- valkeyrie: https://pkg.go.dev/github.com/kms9/valkeyrie
- kine: https://github.com/k3s-io/kine
