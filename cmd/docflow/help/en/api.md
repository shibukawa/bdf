# The API

The API is JSON under `/api`. A request carries one of:

- the session's cookie, from the manager's own pages: a person, with what they may see and do; an administrator, everything;
- the operator's admin token as a bearer (`Authorization: Bearer …`): everything, for scripts;
- an app's token ([Apps](entries.md#apps)): its own jobs, and the endpoints that list it;
- a submit token's `name:secret` as a bearer: sending jobs;
- the tick token: `POST /api/tick` alone.

An error is `{"error": "…"}` with its status: 400 a request the manager cannot take (a node that would make a cycle, a kind of payload a node does not take), 401 no credentials or wrong ones, 403 not yours, 404 not found, 409 a conflict (a name of another kind, a job that ended), 413 a document too large, 415 a format the entry does not take, 502 a printer that failed.

## Who {#me}

| | |
|---|---|
| `GET /api/me` | who the request is from: its kind (a person, an app, a token), id, name, email and groups, whether an administrator, and its roles (`topology`, `audit`) |
| `GET /auth/login`, `POST /auth/logout` | sign in at the provider, sign out |
| `POST /oauth/token` | a token: an app's by the client credentials grant, or a person's for a public client by the authorization code grant with PKCE |
| `GET /oauth/authorize` | sends a person signed in back to a public client with a code ([Frontends](exports.md#frontends)) |

## Jobs {#jobs}

| | |
|---|---|
| `GET /api/my/jobs?limit=` | the jobs the caller sees: theirs and their groups' (an administrator: every job) |
| `GET /api/destinations` | the nodes the caller may send to |
| `POST /api/destinations/<node>/jobs?name=&user=&meta.<key>=` | send the body — a document, or data for a template node — to a node; `X-Meta-*` headers as well, or a `multipart/form-data` body as an endpoint takes it. 201 and the job |
| `POST /api/locations/<display>/slips` | the same, to a display: the older name |
| `GET /api/destinations/<node>/jobs?limit=` | the jobs sent to a node (admin) |
| `GET /api/jobs/<n>` | a job: its payloads, metadata and legs |
| `GET /api/jobs/<n>/document?payload=` | its `.bdf` document — the last made, or the one received — or a payload by id |
| `GET /api/jobs/<n>/pages/<p>.png` | a page of its document, drawn |
| `POST /api/jobs/<n>/cancel` | cancel it |
| `POST /api/jobs/<n>/reprint` | `{"to": "<node>"}`: send it again as a new job of the caller's, to a node they may send to; the job's node by default |
| `POST /in/<endpoint>` | an endpoint's entry, with its key or an app's token ([Endpoints](entries.md#endpoints)) |

## The topology {#topology}

Everyone signed in reads the topology; changing it is the role `topology`'s and the admin token's.

| | |
|---|---|
| `GET /api/topology` | the nodes and their edges, with `takes`, `emits`, `accepts` (what may be sent to the node now) and `broken` (the edges from it that no longer connect: [Broken edges](flow.md#broken)); an administrator gets the settings and the counts too |
| `GET /api/topology/live` | the activity of the jobs the caller sees: counts at each node, and the legs of the last ten minutes; an ETag to poll with (304 while nothing changed) |
| `GET`, `PUT /api/layout` | where the nodes are drawn |
| `GET /api/nodes`, `GET /api/nodes/<name>` | every node; one node |
| `PUT /api/nodes/<name>` | save a node: its `kind` and the kind's settings ([Node kinds](nodes.md)); the answer of a new endpoint carries its `key`, once |
| `DELETE /api/nodes/<name>` | delete it; not while another node sends to it |
| `POST /api/nodes/<name>/key` | rotate an endpoint's key: the new one, once |
| `GET /api/formats` | the formats an entry may take |
| `GET /api/openapi.json?node=` | the OpenAPI document of the endpoints and exports the caller may use ([OpenAPI](#openapi)) |

## Templates {#templates}

| | |
|---|---|
| `GET /api/nodes/<node>/templates` | the templates of a template node |
| `GET /api/nodes/<node>/templates/<id>` | one: its format, the fields it looks up, a sample of its data |
| `GET /api/nodes/<node>/templates/<id>/file` | its file |
| `PUT /api/nodes/<node>/templates/<id>?name=` | upload one, the file as the body, with the file's name |
| `DELETE /api/nodes/<node>/templates/<id>` | delete it |

## Displays, sinks and printers {#devices}

| | |
|---|---|
| `GET /api/nodes/<display>/board` | the slips of a display |
| `GET /api/nodes/<display>/devices` | its devices |
| `POST /api/nodes/<display>/enrollments` | a one-time code to enroll a device, with its QR code and URL |
| `DELETE /api/devices/<id>` | remove a device |
| `GET /api/nodes/<sink>/latest` | the latest documents a sink keeps |
| `GET /api/nodes/<printer>/preview`, `…/preview/<job>/pages/<p>.png` | what a printer's fake printer received, and its pages |

## Inboxes and deliveries {#deliveries}

| | |
|---|---|
| `GET /api/inboxes/<name>` | an inbox: the items held, and the ended ones still shown |
| `POST /api/items/<id>/release` | `{"to": ["<node>", …], "worker": "…"}` |
| `POST /api/items/<id>/discard` | `{"worker": "…"}` |
| `GET /api/deliveries` | the deliveries waiting or failed |
| `POST /api/deliveries/<id>/retry` | send one again now |
| `POST /api/tick` | do what is due ([The tick](inboxes.md#tick)): the tick token, an app with the scope `tick`, the admin token, or a person with the role `topology` |

## HTTP (pull) {#exports}

A backend's token, a frontend's token of the person, or the pages' session ([HTTP (pull)](exports.md)):

| | |
|---|---|
| `GET /out/<name>?owner=&limit=` | the entries the reader may read, newest first |
| `GET /out/<name>/<job>` | one entry |
| `GET /out/<name>/<job>/record` | the archive's record of the job (of an inbox's entry, the job) |
| `GET /out/<name>/<job>/document` | the document, as the archive wrote it |

## Apps and audit {#apps}

| | |
|---|---|
| `GET /api/apps`, `GET /api/apps/<id>` | the apps; one app |
| `PUT /api/apps/<id>` | register or change one: `title`, `scopes`, `group`, `disabled`, and `public` with `redirect_uris` for a public client; the answer of a new app carries its `secret`, once (a public client has none) |
| `DELETE /api/apps/<id>` | delete it |
| `POST /api/apps/<id>/secret` | rotate its secret |
| `GET /api/audit?limit=` | the audit log (the role `audit`) |

## OpenAPI {#openapi}

`GET /api/openapi.json` is an OpenAPI 3.1 document of what the topology publishes, made of it as it is now: an application's developer generates a client of it, or tries the endpoints in a tool that reads OpenAPI.

- Every endpoint's `POST /in/<name>`, with the media types of the formats it takes. Of an endpoint with `accept: data`, the body's schema is the JSON Schema of the data of the templates its data reaches ([The data](templates.md#data)), and the `?meta.<key>=` parameters are the keys the nodes after it look at: a filter's rules, a template node's `template_from` (its templates' ids as the values), an export's `by` and `fields`, and the `{{meta.<key>}}` of an exit's settings. A key the endpoint sets itself is left out: its value wins.
- Every export's `GET /out/<name>` and its entries, which carry the export's fields.
- `GET /api/jobs/<n>`, `…/document` and `POST …/cancel`, for an app that follows the jobs it sent.
- The ways in: an endpoint's key, the admin token, an app's token (client credentials at `/oauth/token`), a person's token (the authorization code with PKCE).

What it describes is the reader's: the role `topology` and the admin token get every endpoint and export; an app, the endpoints that list it with the scope to send (`print` for documents, `print.template` for data) and the exports that list it with the scope `export`; a person signed in, the exports. `?node=<name>` (again, or names with commas) narrows it to those nodes; a node that publishes nothing to the reader is not found (404). Its `info.version` is a hash of the rest, which changes when what is published changes. On the topology's page, OpenAPI downloads the whole document, and the link of the same name in an endpoint's or an export's sidebar the node's.

```sh
curl -H "Authorization: Bearer $TOKEN" "https://<host>/api/openapi.json?node=erp" > erp.openapi.json
```

## IPP {#ipp}

`POST /ipp/print/<node>` is every node's IPP printer, and `POST /ipp/print` the server, whose CUPS-Get-Printers lists the nodes ([IPP](entries.md#ipp)).
