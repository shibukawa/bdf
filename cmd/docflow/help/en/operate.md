# Running the manager

The manager is one binary, `docflow`, on Popcorn Web. It keeps no state in the process: the records and the documents are in a store, so the same binary runs as one process on a site's server, as several behind a load balancer, or as a function of AWS Lambda. It is configured with a file, `config.<APP_ENV>.toml` (`docflow --generate-config toml` prints every key), or with environment variables, named below.

## The store and the secrets {#store}

| Variable | |
|---|---|
| `DOCFLOW_RECORDS` | the store of the records: `bbolt:path` (one process), `dynamodb://table?region=…`, `firestore://project/database`, `s3://bucket?region=…&endpoint=…`, `r2://bucket?endpoint=…`, `gcs://bucket` |
| `DOCFLOW_DOCUMENTS` | the store of the documents (`bbolt`, `s3`, `r2`, `gcs`); empty: the records' store |
| `DOCFLOW_SECRET` | signs the devices' sessions, the apps' tokens and the pages' sessions: 32 bytes or more, the same for every instance; required outside development |
| `DOCFLOW_PREVIOUS_SECRETS` | retired secrets that still open what was signed with them, while a rotation goes on: set the new secret, keep the old one here for a day, then drop it |
| `DOCFLOW_ADMIN_TOKEN` | the bearer token of the API, for scripts; people sign in |
| `DOCFLOW_TICK_TOKEN` | the bearer token a scheduler calls `POST /api/tick` with |
| `DOCFLOW_SUBMIT_TOKENS` | `name:secret` pairs that may send jobs over IPP and the API; none lets anyone send |
| `DOCFLOW_TICK_INTERVAL` | how often a process ticks (`30s`); `0` leaves it to a scheduler |
| `DOCFLOW_JOB_KEEP` | how long a job that ended is kept (`336h`, 14 days) |
| `DOCFLOW_MAX_DOCUMENT_MB` | the largest document taken (64) |
| `DOCFLOW_FONT_DIRS`, `DOCFLOW_NO_SYSTEM_FONTS` | the fonts the documents are drawn with: directories of fonts, and whether the machine's own are used |
| `DOCFLOW_SECURE_COOKIES` | `Secure` on the cookies (true; false only for plain http in development) |
| `DOCFLOW_SERVERLESS` | refuse a store one process keeps, tick only when called, and open no ports |

## Sign-in {#oidc}

People sign in at an OpenID Connect provider, Popcorn Web's `[auth.oidc]`:

| Variable | |
|---|---|
| `AUTH_OIDC_ISSUER`, `AUTH_OIDC_CLIENT_ID`, `AUTH_OIDC_CLIENT_SECRET` | the provider, and this client |
| `AUTH_OIDC_REDIRECT_URL` | `https://<host>/auth/callback`, as registered at the provider |

Who gets in, and who has which role ([Roles](roles.md#roles)), is `[oidc]`:

| Variable | |
|---|---|
| `OIDC_GROUPS_CLAIM` | the claim of the groups: `groups`, `roles`, or a dotted path such as Keycloak's `realm_access.roles` (it must be in the ID token) |
| `OIDC_GROUPS_PREFIX` | take only the groups of this prefix |
| `OIDC_ADMINS`, `OIDC_ADMIN_GROUPS` | the administrators: `issuer#subject` or a verified email, and groups; one rule is required outside development |
| `OIDC_TOPOLOGY_ADMINS`, `OIDC_TOPOLOGY_ADMIN_GROUPS` | the role `topology` |
| `OIDC_AUDITORS`, `OIDC_AUDITOR_GROUPS` | the role `audit` |
| `OIDC_ALLOWED_ISSUERS`, `OIDC_ALLOWED_DOMAINS`, `OIDC_ALLOWED_GROUPS` | narrow who gets in: an issuer, an email domain (or Google's `hd`), a group; administrators always get in |

## Serverless {#serverless}

On AWS Lambda (through the Lambda Web Adapter) the function is configured with variables, `DOCFLOW_SERVERLESS=true` among them; the records are in DynamoDB (a table with the string keys `pk` and `sk`) and the documents in S3; EventBridge Scheduler calls `POST /api/tick` every minute with the tick token. There are no ports of the nodes' own: computers print to `/ipp/print/<name>`, and applications call the API (`/api/upload/<name>`, `/api/inbox/<name>`, `/api/download/<name>`) as anywhere. To print on a site's printers, the function is in a VPC that reaches them. Building the function and its image is in the repository's README.

## Printers and fonts {#printers}

A printer node reaches its printer over the network: `ipp`, `ipps`, `socket` (port 9100) or `lpd`, or a queue of a CUPS server. The documents are drawn with the machine's fonts and those of `DOCFLOW_FONT_DIRS`. A fake printer farm (`tools/fakeprinter`) stands in for printers while a topology is set up: a printer node's Preview shows what the fake printer received.
