# BDF Print Manager

A software-defined document network. Documents enter at virtual printers
(every node is an IPP printer) and at API endpoints (a document, or a
template and the data to fill it with), flow through the nodes an
administrator composes — queues that fan out, classes that send to one
member in turn, inboxes that hold a job until a person sends it on — and
end on display boards (the devices enrolled at a location, such as the
tablets of a warehouse), on printers, and in archives that keep a copy
nobody changes. Paper is one destination among others. The topology is
laid out on the admin page, as the printers and classes of a CUPS server
are, with edges drawn between nodes (see *Topology* below).

The manager is [print/manager](../../print/manager); this is its web
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
(`[dev.idp]` in `popcornweb.toml`): it gives the Print Manager its issuer
and a client of the run (`AUTH_OIDC_*`), so nothing about the IdP is typed.
Signing in is picking one of the people of
[devidp.toml](devidp.toml), no password checked: Ada Admin (`pm-admins`,
the administrators by `config.dev.toml`), Mina Member (`warehouse-a`) and
Otto Outsider (no groups). Editing the roster takes effect without a
restart; `TestDevRoster` checks that its people still hit the rules. Or

```sh
go run ./cmd/printmanager
```

with an IdP of your own in `.env.local` (copy `.env.example`; git ignores
it). Neither needs a secret in development.

In development (no `APP_ENV`), it reads `config.dev.toml` at the
repository's root (the `[oidc]` rules above, plain-http cookies), keeps its
store in `printmanager.db` and makes a random secret. Configure it with `config.{APP_ENV}.toml` (the
`[manager]` table; `printmanager --generate-config toml` prints every key)
or with variables:

| Variable | |
|---|---|
| `MANAGER_RECORDS` | store of the records: `bbolt:path`, `dynamodb://table?region=…`, `firestore://project/database`, `s3://bucket?region=…&endpoint=…`, `r2://bucket?endpoint=…`, `gcs://bucket` |
| `MANAGER_DOCUMENTS` | store of the slips' documents (`bbolt`, `s3`, `r2`, `gcs`); empty for the records' store |
| `MANAGER_SECRET` | signs the devices' sessions, the apps' tokens and the pages' sessions: 32 bytes or more, the same for every instance |
| `MANAGER_PREVIOUS_SECRETS` | retired secrets that still open what was signed with them, while a rotation goes on (set the new secret, keep the old one here for a day, then drop it) |
| `MANAGER_ADMIN_TOKEN` | bearer token of the API (`/api`), for scripts (people sign in) |
| `MANAGER_TICK_TOKEN` | bearer token a scheduler calls `POST /api/tick` with |
| `MANAGER_SUBMIT_TOKENS` | `name:secret` pairs that may send jobs (IPP Basic, or bearer `name:secret`); none lets anyone send |
| `MANAGER_SERVERLESS` | refuse a store one process keeps, and tick only when called |
| `MANAGER_TICK_INTERVAL` | how often a process ticks (default 30s) |
| `MANAGER_JOB_KEEP` | how long a job that ended is kept, with its document once nothing shows it (default 336h, 14 days) |

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

The rest of `[auth]` and `[session]` is the Print Manager's
(`handlers.ConfigureAuth`): the state of a login on its way through the
IdP is kept in the records' store, so that any instance finishes it, and
the session is a sealed cookie (`session.keyring.secret`, or made of
`MANAGER_SECRET`). The groups are copied at each login; who is an admin
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
`print.template`, `jobs.cancel`, `tick` (a scheduler's) and `admin` (a
script's: the API as an administrator). An app's jobs are its `group`'s
when the destination gives them none. Disabling an app refuses its
credentials and the tokens out.

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

A person signs in at the IdP (see Sign-in) and sees `/jobs`: their jobs
and their groups', and a form that prints a document to a node they may
print to. `/topology` shows everyone the network as the administrator
laid it out (the nodes, their kinds and edges; not their settings), with
the activity of the jobs the person may see: how many are open, held or
waiting at each node, and the jobs of the last minutes running along the
edges; a node lists the person's jobs there and prints there when they
may. An administrator sees every job, and goes on to the
administrator's pages, one item each: `/admin/topology` (the editor of
the topology: every node is placed, connected and set up there, with its
templates, key or devices in the sidebar), `/admin/inboxes` (release and
discard), `/admin/jobs` (every job and its route), `/admin/deliveries`
(what waits for a printer or an archive, and the tick) and `/audit` (who
did what: sign-ins, changes, keys, devices, refused documents, releases,
cancellations). The API takes the session's cookie (from the manager's
own pages) as it takes the tokens; scripts use the tokens.

1. Add a printer and a display (on the topology, or the API; both are
   nodes):

   ```sh
   curl -X PUT -H "Authorization: Bearer $TOKEN" localhost:8080/api/nodes/wh-a-laser -d '{"kind":"printer","url":"ipps://wh-a-laser.local/ipp/print"}'
   curl -X PUT -H "Authorization: Bearer $TOKEN" localhost:8080/api/nodes/wh-a -d '{"kind":"display","title":"第 1 倉庫","scope":"group","owner":"warehouse-a","fallback":{"printer":"wh-a-laser","after":"10m"}}'
   ```

2. Enroll a tablet: "Add a device" in the display's sidebar on the topology
   shows a one-time code and its QR code; the tablet opens `/tablet?code=…`.
   Its page makes a key that cannot be exported and keeps it; only that
   browser signs in as the device. Add the page to the home screen.
3. Print to `ipp://<host>:<port>/ipp/print/wh-a` (or `POST /api/locations/wh-a/slips`
   with the document as the body, which answers with the slip). PDF, Office
   documents, images, PWG and Apple rasters and `.bdf` are taken. The
   metadata keys `ipp.job-account-id`, `ipp.job-accounting-user-id`,
   `ipp.document-metadata` and `ipp.x-…` go to the printer as those IPP
   attributes: an IPP request's own give them, or `?meta.ipp.job-account-id=…`
   of a `POST`.

## Topology

Every node, of every kind, is `PUT /api/nodes/<name>` with its `kind`
and the kind's settings, in one namespace, and each is an IPP printer at
`/ipp/print/<name>` (a job sent to a printer's name is printed there; one
sent to a display's name is shown there). A job sent to a node follows
the edges from it; its record (`/api/jobs/<n>`) keeps a leg for every
node it reached. The document is kept once for all of them.

| Kind | Does | Settings |
|---|---|---|
| `endpoint` | a REST entry for applications: `POST /in/<name>` with the endpoint's own key takes a document as the body, or data (`accept: data`: JSON, YAML, CSV, an Excel workbook or a Parquet file, for a template node after it), and forwards the job to every node of `to`. `output` is what goes on of a document: `bdf` (the default: the .bdf made of it) or `file` (as received, for a convert node or an archive after it). `formats` are the formats it takes (as the converter names them: `pdf`, `docx`, `image`, `csv`, `bdf` …, and `json`, `yaml`; empty: any); a file of another format is refused (415) and the refusal is in the audit log. `meta` is metadata the endpoint adds to every job it takes. The key is shown once, when the endpoint is made (`key` in the answer) or rotated (`POST /api/nodes/<name>/key`); the manager keeps only its argon2id hash. With `port`, a process listens there too: every `POST` to it is to the endpoint | `to`, `accept`, `output`, `formats`, `meta`, `port` |
| `cups` | an IPP entry for computers and business systems: a CUPS server on a `port` of the process's own, with one printer, the node itself (`lpstat -h <host>:<port> -p` lists it; computers print to `ipp://<host>:<port>/printers/<node>`), which forwards what it receives to every node of `to`; plain `ipp`, or `ipps` with `tls_cert` and `tls_key` (files). `accept`, `output`, `formats` and `meta` as an endpoint's; the job attributes of a Print-Job request are metadata (`ipp.media`, `ipp.copies` …) | `to`, `port`, `tls_cert`, `tls_key`, `accept`, `output`, `formats`, `meta` |

| `queue` | forwards the job to every node of `to` (a fan-out when there are several) | `to` |
| `class` | sends the job to one of `members`, in turn (`round-robin`, the default) or the first (`failover`); a printer that does not take it is skipped for the next member | `members`, `policy` |
| `filter` | sends the job on by `rules`: the first whose `match` fits (`name`, `user` and `format` as RE2 regular expressions, `min_pages`, `max_pages`, `meta` as regular expressions by key of the job's metadata, `kind` of the payload: `bdf`, `data` or `any`) forwards it to the rule's `to`; a job no rule fits goes to `otherwise`, or ends there (dropped) when there is none. What reaches it goes on as it is | `rules`, `otherwise` |
| `template` | fills one of its templates (`PUT /api/nodes/<name>/templates/<id>`; they go with the node) with the data that reaches it, and sends the form on as a .bdf document to every node of `to`. `template` names the template, or `template_from` the metadata key whose value names it (a sender's `X-Meta-Doc-Type`); `list_name` names data that is a list of records (a CSV file, a sheet, a JSON array; the file's name by default), an object is the data itself; `strict` refuses data that lacks a name the template looks up. The step is done like a delivery (the background, or the tick), and fails for good when the template is missing or the data cannot be read | `to`, `template`, `template_from`, `list_name`, `strict` |
| `convert` | makes a .bdf document of a file that reaches it as received (a PDF, an Office document, a picture …; `formats` limits which), and sends it on to every node of `to`; a .bdf goes on as it is. JSON and YAML have no document of their own (fill a template with them), and a CSV or Parquet file has no pages to print yet. The step is done like a delivery | `to`, `formats` |
| `inbox` | holds the job until a person releases it to nodes (those of `to`, or any when `to` is empty) or discards it; ended items stay `keep` (14 days) | `to`, `keep` |
| `display` | shows the job as a slip on the devices enrolled there, and prints on paper what nobody opened within `fallback.after` (or every unopened slip while no device is connected); slips done or canceled stay `keep` (14 days) | `fallback` (`printer`, `after`), `keep`, `offline_after` |
| `printer` | prints the job on paper, on a printer the manager reaches over the network (`ipp`, `ipps`, `socket`, `lpd`). With `preview`, the URL of the printer on a fake printer farm's API (`tools/fakeprinter`: `http://host:port/api/printers/<printer>`), what the fake printer received is shown beside the node, as the farm read it back (`GET /api/nodes/<name>/preview`, `…/preview/<job>/pages/<p>.png`) | `url`, `route` (`direct`), `preview` |
| `archive` | writes the payload that reaches it (a `.bdf`, or any file or data as received) and a JSON record of the job (its metadata, payloads and legs) to `store` (a documents URL: `bbolt`, `s3`, `r2`, `gcs`), each object once, under `collection` (`archive-<name>`); the manager never changes or deletes them. Give the S3 bucket a default retention (Object Lock) for an audit trail nobody can change | `store`, `collection` |
| `sink` | keeps the `latest` documents that reach it (10) to look at on the topology (the first page shows on the node, every page in the sidebar: `GET /api/nodes/<name>/latest`, `GET /api/jobs/<n>/pages/<p>.png`); an older one goes when a newer comes | `latest` |

Every node has a `title` and a `description` too (as a C4 element: the
name in the title bar, `[kind]` and the title under it, then the
description, a note of what it is for, up to 1000 characters, the first
lines on the canvas), and a `scope` and an `owner` (see Sign-in). A display, a
printer or a sink keeps its kind; the other kinds change into one
another, keeping the name.

What flows along an edge has a kind: a `bdf` document, `data` (JSON,
YAML, CSV, Excel, Parquet) or `any` other file as received. An entry
sends on what its `accept` and `output` say, a template or a convert
node the .bdf it makes, every other node what reaches it; a printer or
a display takes `bdf` only, a template node `data` only. The kinds that
reach every node are inferred from the entries along the edges when a
node is put, and a new edge that would bring a kind a node does not
take is refused (400: "laser (printer) takes bdf, but any reaches it:
raw → fan → laser"). A change of a node's settings alone (an entry
taking data now, a queue made a template node) is not refused for the
edges it leaves carrying what the node after them does not take: they
stay, broken, until mended (a convert or a template node put between,
the setting changed back, or the edge removed), and a job gets a leg of
the error where it reaches the node that does not take it. `GET
/api/topology` shows each node's `takes`, `emits`, `accepts` (what may
be sent to it now: what it takes, within what every node after it takes,
so that a queue in front of a printer takes `bdf` alone, as a wildcard
pin is typed by its first edge) and `broken` (the edges from it that no
longer connect, with the kinds), and the editor colors the ports and the
edges by them (blue bdf, green data, amber any, purple several; hollow
and dashed a port not typed yet), fades the nodes that would not take
the edge being dragged, and draws a broken edge dashed and red. A node's
title bar is colored by its role (source: endpoint, cups; control flow:
queue, class, inbox; filter; transform: template, convert; sink:
display, printer, archive, sink), with an icon of its kind. A job
carries its `payloads`: what was received
(`kind`, `format`, `pages`, `size`, `sha256`) and what the template and
convert nodes made of it (`made`), each leg saying which payload it got
(`payload`) and made (`made`); `GET /api/jobs/<n>/document` downloads
the .bdf document to show (the last made, or the one received), and
`?payload=<id>` any of them. A job's `meta` is its metadata: what the
sender gave (`X-Meta-<Key>` headers and `?meta.<key>=` parameters of
the API, the job attributes of an IPP request as `ipp.<name>`; claims,
routed by and never trusted), what the entry added (`meta`), and what
the manager verified, under keys that start with `$`: `$job`,
`$received`, `$entry`, `$source`, `$sender`, `$sender_kind`, `$file`. A
filter's rules and a template node's `template_from` see those, and the
payload's `$kind`, `$format`, `$pages`, `$size` and `$sha256`, and the
job's `name` and `user`.

A port of an entry node is the process's: it is opened when the node is
put and closed when the node loses it (and checked at every tick), on
every address of the host. A serverless deployment (`MANAGER_SERVERLESS`)
refuses ports: applications call the manager's own `/in/<name>`, and
there is no IPP port of a node's own there.

Three virtual printers and an endpoint, for example: an audit copy, a
site's printers in turn, both warehouses at once, and an application that
prints invoices from data.

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

A template belongs to its template node (`PUT /api/nodes/<name>/templates/<id>`;
it goes when the node does). It is a compiled template (`.bdft`, from
`bdf template compile`) or one of Excel, Word, PowerPoint, a PDF form,
pdfme or text, as `bdf fill` takes them; `GET /api/nodes/<name>/templates/<id>`
lists the fields its tags look up. The data that fills it is a JSON or
YAML object, a CSV file or an Excel workbook (a list of records gets the
node's `list_name`, or the file's name); a Parquet file is not filled into
a template yet.

A node is checked when it is put: the nodes it names exist, the name is
free of the other kinds, no path leads back to it, and every node takes
the kinds of payload that reach it. A node another one sends to is not
deleted. Deliveries to printers and archives, and the steps of template
and convert nodes, are done by a process right after the job is routed,
and by the tick (`POST /api/tick`; every minute on a serverless
platform); one that fails is tried again after a growing wait, up to 20
times (not at all after a failure no retry mends: a template missing,
data that cannot be read), and `POST /api/deliveries/<id>/retry` does it
again now. `/admin/topology` is the editor (a palette of kinds,
ports to connect, a sidebar for the node selected, pan and zoom; where the
nodes are drawn is `/api/layout`); a CUPS client lists every node with
CUPS-Get-Printers at `/ipp/print` (`bdf printers -url ipp://<host>:<port>/ipp/print`).

`GET /auth/login` signs in at the IdP and `POST /auth/logout` out;
`GET /api/me` tells who; `GET /api/my/jobs` lists a person's jobs and
`GET /api/destinations` the nodes to print to; `POST /api/jobs/<n>/reprint`
(`{"to": "<node>"}`, or the job's node) sends a job the person sees again,
as a new job of theirs, to a node they may send to; `GET /api/audit` is
the log (the role `audit`).
Everyone signed in reads `GET /api/topology` (the nodes and edges; an
administrator gets the settings and the counts too), `GET /api/layout`
(where they are drawn) and `GET /api/topology/live`: the activity of the
jobs they may see, as counts at each node (`held`, `open`, `pending`,
`failed`, `recent` in the last hour) and the legs of the last ten minutes
(`flows`), with an ETag to poll with (304 while nothing changed). The
rest of the topology's API is the admin token's: `GET`, `PUT` and `DELETE /api/nodes/<name>` and `POST
/api/nodes/<name>/key`; a display's `GET …/board`, `GET …/devices` and
`POST …/enrollments`, and `DELETE /api/devices/<id>`; a template node's `GET
…/templates` and `GET`, `PUT` and `DELETE …/templates/<id>` and its
`/file`; `GET /api/formats` (the formats an entry may take); `GET` and
`PUT /api/layout` (where the nodes are drawn); `GET` and `POST
/api/destinations/<name>/jobs`; `GET /api/jobs/<n>`, its `/document` and
`POST …/cancel`; `GET /api/inboxes/<name>`, `POST /api/items/<id>/release`
(`{"to": [...], "worker": "..."}`) and `…/discard`; `GET /api/deliveries`
and `POST /api/deliveries/<id>/retry`.

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
