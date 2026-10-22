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
| `POST /api/destinations/<node>/jobs?name=&user=&meta.<key>=` | send the body — a document, or data for a template node — to a node; `X-Meta-*` headers as well. 201 and the job |
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

## Exports {#exports}

A backend's token, a frontend's token of the person, or the pages' session ([Exports](exports.md)):

| | |
|---|---|
| `GET /out/<name>?owner=&limit=` | the entries the reader may read, newest first |
| `GET /out/<name>/<job>` | one entry |
| `GET /out/<name>/<job>/record` | the archive's record of the job |
| `GET /out/<name>/<job>/document` | the document, as the archive wrote it |

## Apps and audit {#apps}

| | |
|---|---|
| `GET /api/apps`, `GET /api/apps/<id>` | the apps; one app |
| `PUT /api/apps/<id>` | register or change one: `title`, `scopes`, `group`, `disabled`, and `public` with `redirect_uris` for a public client; the answer of a new app carries its `secret`, once (a public client has none) |
| `DELETE /api/apps/<id>` | delete it |
| `POST /api/apps/<id>/secret` | rotate its secret |
| `GET /api/audit?limit=` | the audit log (the role `audit`) |

## IPP {#ipp}

`POST /ipp/print/<node>` is every node's IPP printer, and `POST /ipp/print` the server, whose CUPS-Get-Printers lists the nodes ([IPP](entries.md#ipp)).
