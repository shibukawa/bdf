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

Copy `.env.example` to `.env.local` (git ignores it) and fill in the admin
token and the secret. Then, from the repository's root, either

```sh
pw dev
```

which regenerates the pages, rebuilds and restarts on every change (it
reports that `public.go` is missing: the pages' files are embedded by
cmd/printmanager itself, so that step is not needed), or

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
| `MANAGER_ADMIN_TOKEN` | bearer token of the API (`/api`) and of the page `/admin` |
| `MANAGER_TICK_TOKEN` | bearer token a scheduler calls `POST /api/tick` with |
| `MANAGER_SUBMIT_TOKENS` | `name:secret` pairs that may send jobs (IPP Basic, or bearer `name:secret`); none lets anyone send |
| `MANAGER_SERVERLESS` | refuse a store one process keeps, and tick only when called |
| `MANAGER_TICK_INTERVAL` | how often a process ticks (default 30s) |

AWS credentials come from the environment (`AWS_ACCESS_KEY_ID` …), web
identity, the container's or the instance's role.

## Use it

1. Add a printer and a location (on `/admin`, or the API):

   ```sh
   curl -X PUT -H "Authorization: Bearer $TOKEN" localhost:8080/api/printers/wh-a-laser -d '{"url":"ipps://wh-a-laser.local/ipp/print"}'
   curl -X PUT -H "Authorization: Bearer $TOKEN" localhost:8080/api/locations/wh-a -d '{"title":"第 1 倉庫","fallback":{"printer":"wh-a-laser","after":"10m"}}'
   ```

2. Enroll a tablet: "端末を足す" on `/admin` shows a one-time code and its QR
   code; the tablet opens `/tablet?code=…`. Its page makes a key that cannot
   be exported and keeps it; only that browser signs in as the device.
   Add the page to the home screen.
3. Print to `ipp://<host>:<port>/ipp/print/wh-a` (or `POST /api/locations/wh-a/slips`
   with the document as the body). PDF, Office documents, images, PWG and
   Apple rasters and `.bdf` are taken.

## On AWS Lambda

[deploy/lambda/Dockerfile](deploy/lambda/Dockerfile) builds the container
image (with the AWS Lambda Web Adapter). Keep the records in DynamoDB
(a table with the string keys `pk` and `sk`) and the documents in S3, and
let EventBridge Scheduler call `POST /api/tick` every minute. To print on a
site's printers, put the function in a VPC that reaches them (Direct
Connect or a VPN).

## Generated code

The pages are `.pw.html` templates. After changing them, or
`handlers/config.go`, run from the repository's root:

```sh
pw generate --code-only
```

and commit the `*_pw_gen.go` files (the build of the module does not run
pw). `popcornweb.toml` is at the root because pw wants it beside `go.mod`;
it reads only `cmd/printmanager/handlers`.

## Style

The pages use [Tailwind CSS](https://tailwindcss.com) and
[daisyUI](https://daisyui.com), built into `public/app.css` (committed, and
embedded in the binary: no CDN). No Node or npm: the standalone
`tailwindcss` CLI loads daisyUI's release files in `tailwind/` (versions
and hashes in `tailwind/VERSIONS`). After changing the pages' markup or
classes, with that CLI on the PATH (`brew install tailwindcss`):

```sh
go generate ./cmd/printmanager
```

CI builds it again and fails on a difference.

