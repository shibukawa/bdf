# BDF Docflow

A document workflow engine: a software-defined document network. Documents enter at virtual printers
(every node is an IPP printer) and at API endpoints (a document, or a
template and the data to fill it with), flow through the nodes an
administrator composes — queues that fan out, classes that send to one
member in turn, inboxes that hold a job until a person sends it on — and
end on display boards (the devices enrolled at a location, such as the
tablets of a warehouse), on printers, and in archives that keep a copy
nobody changes. Paper is one destination among others. The topology is
laid out on the admin page, as the printers and classes of a CUPS server
are, with edges drawn between nodes (see *Topology* below).

The engine is [docflow](../../docflow); this is its web
application on [Popcorn Web](https://github.com/shibukawa/popcornweb). It
keeps no state in the process: the records and the documents are in a
[yuruidb](https://github.com/shibukawa/yuruidb) store, so the same binary
runs as one process on a site's server or as a function of AWS Lambda.

## Run it

From the repository's root, either

```sh
pw dev
```

which regenerates the pages and the public tree, rebuilds and restarts on
every change, and signs people in at pw's development OpenID Provider
(`[dev.idp]` in `popcornweb.toml`): it gives Docflow its issuer
and a client of the run (`AUTH_OIDC_*`), so nothing about the IdP is typed.
Signing in is picking one of the people of
[devidp.toml](devidp.toml), no password checked: Ada Admin (`docflow-admins`,
the administrators by `config.dev.toml`), Mina Member (`warehouse-a`) and
Otto Outsider (no groups). Editing the roster takes effect without a
restart; `TestDevRoster` checks that its people still hit the rules. Or

```sh
go run ./cmd/docflow
```

with an IdP of your own in `.env.local` (copy `.env.example`; git ignores
it). Neither needs a secret in development.

In development (no `APP_ENV`), it reads `config.dev.toml` at the
repository's root (the `[oidc]` rules above, plain-http cookies), keeps its
store in `docflow.db` and makes a random secret. Configure it with `config.{APP_ENV}.toml` (the
`[docflow]` table; `docflow --generate-config toml` prints every key)
or with variables:

| Variable | |
|---|---|
| `DOCFLOW_RECORDS` | store of the records: `bbolt:path`, `dynamodb://table?region=…`, `firestore://project/database`, `s3://bucket?region=…&endpoint=…`, `r2://bucket?endpoint=…`, `gcs://bucket` |
| `DOCFLOW_DOCUMENTS` | store of the slips' documents (`bbolt`, `s3`, `r2`, `gcs`); empty for the records' store |
| `DOCFLOW_SECRET` | signs the devices' sessions, the apps' tokens and the pages' sessions: 32 bytes or more, the same for every instance |
| `DOCFLOW_PREVIOUS_SECRETS` | retired secrets that still open what was signed with them, while a rotation goes on (set the new secret, keep the old one here for a day, then drop it) |
| `DOCFLOW_ADMIN_TOKEN` | bearer token of the API (`/api`), for scripts (people sign in) |
| `DOCFLOW_TICK_TOKEN` | bearer token a scheduler calls `POST /api/tick` with |
| `DOCFLOW_SUBMIT_TOKENS` | `name:secret` pairs that may send jobs (IPP Basic, or bearer `name:secret`); none lets anyone send |
| `DOCFLOW_SERVERLESS` | refuse a store one process keeps, and tick only when called |
| `DOCFLOW_TICK_INTERVAL` | how often a process ticks (default 30s) |
| `DOCFLOW_JOB_KEEP` | how long a job that ended is kept, with its document once nothing shows it (default 336h, 14 days) |

## Sign-in

Every page but `/tablet` is for a person signed in at an OIDC IdP
(authorization code with PKCE, through popcornweb's `plugin/auth`).
Anyone the IdP signs in gets in, with no registration here, and sees
their own jobs and their groups' jobs, and the nodes they may print to;
the admins see everything and get `/admin` and `/audit`, the others get
403 there. The IdP is popcornweb's `[auth.oidc]`; who gets in
and who is an admin is `[oidc]`:

| Variable | |
|---|---|
| `AUTH_OIDC_ISSUER`, `AUTH_OIDC_CLIENT_ID`, `AUTH_OIDC_CLIENT_SECRET` | the IdP and this client (required outside dev) |
| `AUTH_OIDC_REDIRECT_URL` | `https://<host>/auth/callback`, as registered at the IdP (outside dev) |
| `OIDC_GROUPS_CLAIM` | claim of the groups: `groups`, `roles`, or a dotted path such as Keycloak's `realm_access.roles` (it must be in the ID token) |
| `OIDC_GROUPS_PREFIX` | take only the groups of this prefix |
| `OIDC_ADMINS` | admins: `issuer#subject`, or an email the IdP verified |
| `OIDC_ADMIN_GROUPS` | groups whose members are admins (one admin rule is required outside dev) |
| `OIDC_TOPOLOGY_ADMINS`, `OIDC_TOPOLOGY_ADMIN_GROUPS` | the role `topology`: the admin pages and the API of the topology (nodes, templates, keys, apps, devices, inboxes, deliveries), short of an admin |
| `OIDC_AUDITORS`, `OIDC_AUDITOR_GROUPS` | the role `audit`: `/audit`, `GET /api/audit`, and every job (reading); an admin has both roles |
| `OIDC_ALLOWED_ISSUERS`, `OIDC_ALLOWED_DOMAINS`, `OIDC_ALLOWED_GROUPS` | narrow who gets in (email domain or Google's `hd`; admins always get in) |

The rest of `[auth]` and `[session]` is Docflow's
(`handlers.ConfigureAuth`): the state of a login on its way through the
IdP is kept in the records' store, so that any instance finishes it, and
the session is a sealed cookie (`session.keyring.secret`, or made of
`DOCFLOW_SECRET`). The groups are copied at each login; who is an admin
is decided on each request. A sign-in after a sign-out asks the IdP to
let the person choose who they are: popcornweb (0.5.14 or later) marks
`/auth/login` to carry `prompt=select_account login` until a sign-in
completes, and sends only requests for a page to the sign-in (the icons
Safari asks for on its own, `/favicon.ico` and `/apple-touch-icon*.png`,
get 401), so the IdP does not sign the same person in again without a
word.

**Apps** are the applications that send with credentials of their own (a
business system, a WMS): an administrator registers one on `/admin/apps`
(`PUT /api/apps/<id>` with its `scopes`; the secret comes back once,
`POST /api/apps/<id>/secret` rotates it), and an endpoint lists the apps
it takes (`apps`). The app takes a token at `POST /oauth/token` (the
client credentials grant: `-u id:secret -d grant_type=client_credentials`,
ten minutes) and sends with it as a bearer to `/in/<endpoint>`, or
prints over IPP with the same id and secret as Basic; it sees and, with
the scope `jobs.cancel`, cancels its own jobs. The scopes are `print`,
`print.template`, `jobs.cancel`, `export` (reads the exports that list
it: [help/en/exports.md](help/en/exports.md)), `tick` (a scheduler's) and `admin` (a
script's: the API as an administrator). An app's jobs are its `group`'s
when the destination gives them none. Disabling an app refuses its
credentials and the tokens out.

A **public client** is a frontend of one's own in people's browsers (a
customer portal, a single-page app on another origin): registered with
`"public": true` and its `redirect_uris` (https, or http on localhost),
it has no secret, and gets a token of the person using it with the
authorization code grant and PKCE (S256): the frontend sends the person
to `GET /oauth/authorize?response_type=code&client_id=…&redirect_uri=…&code_challenge=…&code_challenge_method=S256&state=…`,
where they sign in at the IdP as for the pages (once a session; the
registration is the consent, no screen asks), and comes back with a
`code`, which it exchanges at `POST /oauth/token`
(`grant_type=authorization_code`, `client_id`, `code`, `code_verifier`,
`redirect_uri`) for a token of an hour. The token is the person, with
the client's scopes: it opens the exports and `/api/me`, nothing else of
the API. A page of the client's origin may call the token endpoint and
the exports (CORS).

Every node (a display, a printer, a queue, a class …) has a scope and an
owner (`"scope":"group", "owner":"warehouse-a"`, a value of the groups
claim; `personal` and a person's email; `shared`, the default): who may
print to it from the pages, and the group whose members see the jobs
sent to it (of a display with none, its fallback printer's). A job
records who sent it as verified (`submitted_by`: the person signed in,
the token, the endpoint, or anonymous) apart from the name the sender
gave (`claimed_user`, shown only). The sender and the admins cancel a
job; its group only sees it.

`pw dev`'s provider is the local IdP. The tests also use
`print/fakeidp` (`fakeidp.New(...).Listen("127.0.0.1:9781")` serves it by
hand), for what the roster cannot say: a nested groups claim
(`realm_access.roles`), an email the IdP did not verify, a token tampered
with.

AWS credentials come from the environment (`AWS_ACCESS_KEY_ID` …), web
identity, the container's or the instance's role.

## Use it

The manual is in the application: **Help** in the bar of every page, at
`/help`, in English and Japanese by the browser's language. Its chapters
are the Markdown files of [help/en](help/en) and [help/ja](help/ja),
embedded by [package help](help/help.go) and rendered when served
(goldmark, no raw HTML): the pages, the topology and its editor, every
kind of node and its settings, what flows along the edges and the
metadata, templates and data, the entries (endpoints, CUPS server nodes,
apps, IPP), display boards and devices, inboxes and deliveries, sign-in,
roles and audit, the API, and the configuration. A link between chapters
is written as `jobs.md`, so that the files read on GitHub too; every
chapter has the same section ids in every language, so that the pages
link to the manual without a language (the editor's sidebar links
`/help/nodes#printer`), and `help_test.go` checks the links, the ids,
and that every kind of node has its section. Changing the manual is
editing the Markdown: no generation.

In short: a person signs in at the IdP (see Sign-in) and sees `/jobs`,
their jobs and their groups' and a form that prints to a node they may
send to, and `/topology`, the network with the activity of their jobs.
The role `topology` gets the administrator's pages under `/admin` (the
editor of the topology, apps, inboxes, jobs, deliveries) and the role
`audit` gets `/audit`. Every node is `PUT /api/nodes/<name>` with its
`kind` and the kind's settings — `endpoint`, `cups`, `queue`, `class`,
`filter`, `template`, `convert`, `inbox`, `display`, `printer`,
`archive`, `sink`, `export` ([help/en/nodes.md](help/en/nodes.md)) —
and an IPP printer at `/ipp/print/<name>`; a display's tablets are
enrolled with a one-time code from the display's sidebar
([help/en/displays.md](help/en/displays.md)); an export after an
archive lets backends and people's frontends fetch the documents at
`/out/<name>` ([help/en/exports.md](help/en/exports.md)); and the API is
[help/en/api.md](help/en/api.md). For example: an audit copy, a site's
printers in turn, both warehouses at once, a review before printing, and
an application that prints invoices from data:

```sh
curl -X PUT -H "Authorization: Bearer $TOKEN" localhost:8080/api/nodes/audit -d '{"kind":"archive","store":"s3://audit?region=ap-northeast-1"}'
curl -X PUT -H "Authorization: Bearer $TOKEN" localhost:8080/api/nodes/invoices -d '{"kind":"queue","to":["audit"]}'
curl -X PUT -H "Authorization: Bearer $TOKEN" localhost:8080/api/nodes/tokyo -d '{"kind":"class","members":["mfp-3f","mfp-4f"]}'
curl -X PUT -H "Authorization: Bearer $TOKEN" localhost:8080/api/nodes/both -d '{"kind":"queue","to":["wh-a","wh-b"]}'
curl -X PUT -H "Authorization: Bearer $TOKEN" localhost:8080/api/nodes/review -d '{"kind":"inbox","to":["wh-a","tokyo"]}'
curl -X PUT -H "Authorization: Bearer $TOKEN" localhost:8080/api/nodes/erp -d '{"kind":"endpoint","to":["invoices"],"formats":["pdf"]}'
# → {"name":"erp", …, "key":"bdfk_…"}: the application sends with that key
curl -X POST -H "Authorization: Bearer bdfk_…" --data-binary @statement.pdf "localhost:8080/in/erp?name=statement&meta.dept=sales"
# data for a template: a data endpoint, and a template node that fills the invoice
curl -X PUT -H "Authorization: Bearer $TOKEN" localhost:8080/api/nodes/forms -d '{"kind":"template","template":"invoice","to":["invoices"]}'
curl -X PUT -H "Authorization: Bearer $TOKEN" "localhost:8080/api/nodes/forms/templates/invoice?name=invoice.xlsx" --data-binary @invoice.xlsx
curl -X PUT -H "Authorization: Bearer $TOKEN" localhost:8080/api/nodes/erp-data -d '{"kind":"endpoint","accept":"data","to":["forms"]}'
curl -X POST -H "Authorization: Bearer bdfk_…" -H "Content-Type: application/json" "localhost:8080/in/erp-data?name=INV-1" -d '{"customer":"ACME","total":12000}'
```

## On AWS Lambda

From the repository's root:

```sh
pw build --target=lambda
docker build .pw/build/lambda/nethttp
```

builds the function (linux/amd64) and its container image with the AWS
Lambda Web Adapter. No config file is shipped: configure the function with
variables, `DOCFLOW_SERVERLESS=true` among them. Keep the records in
DynamoDB (a table with the string keys `pk` and `sk`) and the documents in
S3, and let EventBridge Scheduler call `POST /api/tick` every minute. To
print on a site's printers, put the function in a VPC that reaches them
(Direct Connect or a VPN).

## Generated files

pw reads `popcornweb.toml` at the repository's root (it wants it beside
`go.mod`): `[generate]` names only `cmd/docflow/handlers`, and
`[assets]` places the public tree in `cmd/docflow` (`public`, built
into `dist/public` and embedded by `public.go`). After changing the pages
(`.pw.html`), `handlers/config.go` or a file in `public` (the manual in
`help/` needs nothing: the Markdown is embedded as it is), run from the
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
embedded in the binary: no CDN); the manual's HTML is styled by Tailwind's
typography plugin (`prose`), which the standalone CLI bundles. No Node or npm: the standalone
`tailwindcss` CLI loads daisyUI's release files in `tailwind/` (versions
and hashes in `tailwind/VERSIONS`). After changing the pages' markup or
classes, with that CLI on the PATH (`brew install tailwindcss`):

```sh
go generate ./cmd/docflow
pw generate
```

CI builds both again and fails on a difference.
