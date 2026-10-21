# BDF Print Manager

Virtual printers that show slips on the tablets of a location (a warehouse)
instead of printing them. A WMS or a business system prints to the
location's IPP printer as it prints to any printer; the slip appears on the
tablets enrolled there, and a printer of the location prints what nobody
opened in time, or everything while no tablet is connected.

The manager is [print/manager](../../print/manager); this is its web
application on [Popcorn Web](https://github.com/shibukawa/popcornweb). It
keeps no state in the process: the records and the documents are in a
[yuruidb](https://github.com/shibukawa/yuruidb) store, so the same binary
runs as one process on a site's server or as a function of AWS Lambda.

## Run it

Copy `.env.example` to `.env.local` (git ignores it) and fill in the
secret and the IdP. Then, from the repository's root, either

```sh
pw dev
```

which regenerates the pages and the public tree, rebuilds and restarts on
every change, or

```sh
go run ./cmd/printmanager
```

In development (no `APP_ENV`), it keeps its store in `printmanager.db` and
makes a random secret. Configure it with `config.{APP_ENV}.toml` (the
`[manager]` table; `printmanager --generate-config toml` prints every key)
or with variables:

| Variable | |
|---|---|
| `MANAGER_RECORDS` | store of the records: `bbolt:path`, `dynamodb://table?region=…`, `firestore://project/database`, `s3://bucket?region=…&endpoint=…`, `r2://bucket?endpoint=…`, `gcs://bucket` |
| `MANAGER_DOCUMENTS` | store of the slips' documents (`bbolt`, `s3`, `r2`, `gcs`); empty for the records' store |
| `MANAGER_SECRET` | signs the devices' sessions: 32 bytes or more, the same for every instance |
| `MANAGER_ADMIN_TOKEN` | bearer token of the API (`/api`), for scripts (people sign in) |
| `MANAGER_TICK_TOKEN` | bearer token a scheduler calls `POST /api/tick` with |
| `MANAGER_SUBMIT_TOKENS` | `name:secret` pairs that may send jobs (IPP Basic, or bearer `name:secret`); none lets anyone send |
| `MANAGER_SERVERLESS` | refuse a store one process keeps, and tick only when called |
| `MANAGER_TICK_INTERVAL` | how often a process ticks (default 30s) |

## Sign-in

Every page but `/tablet` is for a person signed in at an OIDC IdP
(authorization code with PKCE, through popcornweb's `plugin/auth`).
Anyone the IdP signs in gets in, with no registration here, and sees
their own jobs and their groups' jobs, and the locations and printers
they may print to; the admins see everything and get `/admin`, the
others get 403 there. The IdP is popcornweb's `[auth.oidc]`; who gets in
and who is an admin is `[oidc]`:

| Variable | |
|---|---|
| `AUTH_OIDC_ISSUER`, `AUTH_OIDC_CLIENT_ID`, `AUTH_OIDC_CLIENT_SECRET` | the IdP and this client (required outside dev) |
| `AUTH_OIDC_REDIRECT_URL` | `https://<host>/auth/callback`, as registered at the IdP (outside dev) |
| `OIDC_GROUPS_CLAIM` | claim of the groups: `groups`, `roles`, or a dotted path such as Keycloak's `realm_access.roles` (it must be in the ID token) |
| `OIDC_GROUPS_PREFIX` | take only the groups of this prefix |
| `OIDC_ADMINS` | admins: `issuer#subject`, or an email the IdP verified |
| `OIDC_ADMIN_GROUPS` | groups whose members are admins (one admin rule is required outside dev) |
| `OIDC_ALLOWED_ISSUERS`, `OIDC_ALLOWED_DOMAINS`, `OIDC_ALLOWED_GROUPS` | narrow who gets in (email domain or Google's `hd`; admins always get in) |

The rest of `[auth]` and `[session]` is the Print Manager's
(`handlers.ConfigureAuth`): the state of a login on its way through the
IdP is kept in the records' store, so that any instance finishes it, and
the session is a sealed cookie (`session.keyring.secret`, or made of
`MANAGER_SECRET`). The groups are copied at each login; who is an admin
is decided on each request.

A location or a printer has a kind and an owner (`"kind":"group",
"owner":"warehouse-a"`, a value of the groups claim; `personal` and a
person's email; `shared`, the default): who may print to it, and the
group whose members see the jobs of a location (else of its fallback
printer). A job records who sent it as verified (`submitted_by`: the
person signed in, the submit token, or anonymous) apart from the name the
sender gave (`claimed_user`, shown only).

For a local IdP, `print/fakeidp` serves one with people of your choosing
(`fakeidp.New(...).Listen("127.0.0.1:9781")`).

AWS credentials come from the environment (`AWS_ACCESS_KEY_ID` …), web
identity, the container's or the instance's role.

## Use it

1. Add a printer and a location (on `/admin`, or the API):

   ```sh
   curl -X PUT -H "Authorization: Bearer $TOKEN" localhost:8080/api/printers/wh-a-laser -d '{"url":"ipps://wh-a-laser.local/ipp/print"}'
   curl -X PUT -H "Authorization: Bearer $TOKEN" localhost:8080/api/locations/wh-a -d '{"title":"第 1 倉庫","kind":"group","owner":"warehouse-a","fallback":{"printer":"wh-a-laser","after":"10m"}}'
   ```

2. Enroll a tablet: "端末を足す" on `/admin` shows a one-time code and its QR
   code; the tablet opens `/tablet?code=…`. Its page makes a key that cannot
   be exported and keeps it; only that browser signs in as the device.
   Add the page to the home screen.
3. Print to `ipp://<host>:<port>/ipp/print/wh-a` (or `POST /api/locations/wh-a/slips`
   with the document as the body). PDF, Office documents, images, PWG and
   Apple rasters and `.bdf` are taken.

## On AWS Lambda

From the repository's root:

```sh
pw build --target=lambda
docker build .pw/build/lambda/nethttp
```

builds the function (linux/amd64) and its container image with the AWS
Lambda Web Adapter. No config file is shipped: configure the function with
variables, `MANAGER_SERVERLESS=true` among them. Keep the records in
DynamoDB (a table with the string keys `pk` and `sk`) and the documents in
S3, and let EventBridge Scheduler call `POST /api/tick` every minute. To
print on a site's printers, put the function in a VPC that reaches them
(Direct Connect or a VPN).

## Generated files

pw reads `popcornweb.toml` at the repository's root (it wants it beside
`go.mod`): `[generate]` names only `cmd/printmanager/handlers`, and
`[assets]` places the public tree in `cmd/printmanager` (`public`, built
into `dist/public` and embedded by `public.go`). After changing the pages
(`.pw.html`), `handlers/config.go` or a file in `public`, run from the
repository's root:

```sh
pw generate
```

and commit what it wrote: the `*_pw_gen.go` files, `dist/public` and
`public_manifest_pw_gen.go`. The module's build runs no generator, so a
plain `go build` embeds what is committed.

## Style

The pages use [Tailwind CSS](https://tailwindcss.com) and
[daisyUI](https://daisyui.com), built into `public/app.css` (committed, and
embedded in the binary: no CDN). No Node or npm: the standalone
`tailwindcss` CLI loads daisyUI's release files in `tailwind/` (versions
and hashes in `tailwind/VERSIONS`). After changing the pages' markup or
classes, with that CLI on the PATH (`brew install tailwindcss`):

```sh
go generate ./cmd/printmanager
pw generate
```

CI builds both again and fails on a difference.
