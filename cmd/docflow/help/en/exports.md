# HTTP (pull) {#exports}

An HTTP (pull) node — the kind `export`, an export for short — is where documents leave the network for the applications and people that fetch them: a customer portal that shows each customer their invoices, a backend that files what was printed. It lists a stocked file: what an archive right before it wrote, for good, or what an inbox right before it holds, while it holds it. Applications and people read the list and the documents at `/out/<name>`. For the manager to send documents to a web service itself, use an [HTTP (push)](nodes.md#http) node instead.

## Setting one up {#setup}

Put an archive or an inbox, and an export after it: drag from its right port onto the export on the Topology page, or give it `"to": ["<export>"]`. The edge carries a stocked file ([Kinds of payload](flow.md#kinds)), drawn cyan. The settings are in [HTTP (pull)](nodes.md#export).

| Before it | Listed | Until |
|---|---|---|
| an archive | once the archive wrote the job | for good: the archive keeps the document, and the entry outlives the job's record ([States](jobs.md#states)) |
| an inbox | as the inbox holds the job | a person releases or discards it ([Inboxes](inboxes.md#inboxes)): the entry goes, and its URLs are not found |

An entry says where the document is kept: `archive`, `key` and `object` of an archive's, or `inbox` and `item` of an inbox's.

- An edge to an export from any other node — an entry, a queue, a filter, a node after an archive — is refused when the node is saved (400). To list documents by their metadata, put a filter in front of the archives, an archive and an export for each. A job sent straight to an export is refused, and one that reaches it otherwise ends there with an error leg. An export is not a node to print to, nor one an inbox releases to.
- An export lists the latest `entries` (1000) for everyone and for each owner; an older entry is still read by its job's number.

## Whose a document is {#owners}

Each entry is listed under its owner: the value of the metadata key the export's `by` names — a customer's id or email the sender gave, `X-Meta-Customer: alice@example.com` or `?meta.customer=…` ([Metadata](flow.md#metadata)). A person reads the entries whose owner is them, comparing the owner with what the export's `subject` says:

| Subject | Compared with |
|---|---|
| `email` (the default) | the person's email, as the identity provider verified it, without case |
| `sub` | the person's subject at the identity provider |
| `id` | the person's ID in the manager |
| `group` | any of the person's groups |

An export without `by` shows a person what they could see of the job anyway: their own, their groups'. The export's scope and owner say who may read it at all, as they say who may send to a node ([Scopes of nodes](roles.md#scopes)).

The owner is what the sender claimed: whoever may send to the entry decides whose a document is. Give the entry only to the backends that know their users, or fix the owner with the entry's `meta`.

## Reading {#reading}

| | |
|---|---|
| `GET /out/<name>` | the entries the reader may read, newest first: the job's number, name, kind, format, pages, size, sha256, owner and fields. `?owner=` one owner's, `?limit=` (100) |
| `GET /out/<name>/<job>` | one entry, with who sent it and where the archive wrote it |
| `GET /out/<name>/<job>/record` | the archive's record of the job: its metadata, payloads and legs when it was written; of an inbox's entry, the job as it is now |
| `GET /out/<name>/<job>/document` | the document as the archive wrote it; each download is in the audit log |
| `GET /out/<name>/<job>/pages/<n>.png?w=` | a page of a `.bdf` document drawn as PNG, `w` pixels wide (800), to look at before the download; not in the audit log |

The export's OpenAPI document — the link OpenAPI in its sidebar, or `GET /api/openapi.json?node=<name>` — describes these, with the fields of its entries ([OpenAPI](api.md#openapi)).

What the reader may not read is **not found** (404), not forbidden: an entry of another owner, of a job that is not there, of a job sent elsewhere. An export does not say what it holds.

## Backends {#backends}

A backend reads every entry of the exports that list it: an app ([Apps](entries.md#apps)) with the scope `export`, ticked under Apps allowed in the export's sidebar. It takes a token with its id and secret (client credentials) and narrows to one of its users with `?owner=`. The administrators and the auditors read every entry too.

```sh
TOKEN=$(curl -s -u 'portal-api:<secret>' -d grant_type=client_credentials https://<host>/oauth/token | jq -r .access_token)
curl -H "Authorization: Bearer $TOKEN" "https://<host>/out/statements?owner=alice@example.com"
curl -H "Authorization: Bearer $TOKEN" -O -J https://<host>/out/statements/1/document
```

## Frontends {#frontends}

A frontend in people's browsers — a customer portal, a single-page app on another origin — reads as the person using it. Register it on Admin → Apps as a **public client**: tick Public client, give its redirect URIs (https, or http on localhost) and the scope `export`. It has no secret. A public client needs no listing under an export's Apps allowed: the person's own entries are what it reads.

1. The frontend sends the person to `GET /oauth/authorize?response_type=code&client_id=…&redirect_uri=…&code_challenge=…&code_challenge_method=S256&state=…`, the authorization code grant with PKCE (S256).
2. The person signs in at the identity provider as for the pages, once a session; the registration is the consent, and no screen asks.
3. The manager sends them back to the redirect URI with a `code`, which the frontend exchanges at `POST /oauth/token` (`grant_type=authorization_code`, `client_id`, `code`, `code_verifier`, `redirect_uri`) for a token of an hour.

The token is the person, with the client's scopes: it opens the exports and `/api/me`, nothing else of the API. A page of a registered origin may call the token endpoint and the exports (CORS). The codes given and the tokens made are in the audit log ([Audit](roles.md#audit)).
