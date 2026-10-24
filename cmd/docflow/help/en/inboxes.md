# Inboxes and deliveries

Two pages are about jobs on their way: Inboxes, where a person sends jobs on by hand (an application does the same through the API), and the administrator's Deliveries, where what the manager sends on its own is watched and retried.

## Inboxes {#inboxes}

An inbox is a node that holds every job that reaches it until someone releases it or discards it: a filter worked by hand — an approval before printing, a batch sent on at the end of the day. The Inboxes page (`/inboxes`; in the bar of those who work an inbox, and in Admin's tabs) shows every inbox you work with the items it holds, oldest first, each with a thumbnail of its first page:

- **Release** sends the item on to nodes: those of the inbox's `to`, each checked by default, or, when the inbox names none, the names typed in the box, comma-separated (the role topology's alone).
- **Discard** ends the item: the job is not sent on from here.
- The name in **Worker** is written on the record and in the audit log, as who released or discarded; your own name when it is empty.
- The thumbnail downloads the document the item holds.

Released and discarded items are shown under the list for the inbox's keep (14 days by default), with where each went and who sent it. Edit the node opens the inbox on the topology.

The topology works an inbox too: click it, and its sidebar (under Usage, for the role topology) lists what it holds, each with its first page, the nodes it goes to, Release and Discard, and the Worker; Discard asks first there. A click on an item's page shows every page large, to zoom in, with Release and Discard beside the download ([Looking at a document](topology.md#lightbox)). A job canceled on the Jobs pages leaves its inboxes.

## Who works an inbox {#workers}

| Who | Works |
|---|---|
| the role topology (an administrator, the admin token, an app with the scope `admin`) | every inbox, and releases to any node of an inbox that names none |
| the people of the inbox's scope | of a `group` inbox, its members; of a `personal` one, its owner; a `shared` inbox has none: the role topology works it |
| an app the inbox lists (Apps allowed), with the scope `inbox` | that inbox, with its token |
| a person's token from a frontend (a public client with the scope `inbox`) | the inboxes the person works |

All but the role topology release to the inbox's `to` alone: an inbox that names no node releases where the role topology says. An inbox you do not work is refused (403), with why.

## The inboxes' API {#api}

What the page does is the API, for an application that approves documents itself — a workflow system, a portal where a manager signs off:

| | |
|---|---|
| `GET /api/inbox` | the inboxes the caller works: `name`, `title`, `description`, `to` (the nodes it releases to), `held` (how many items) and `rev` |
| `GET /api/inbox/<name>` | the inbox: `held`, the items held oldest first, `ended`, those released or discarded still shown, latest first, and `to`. With `If-None-Match` and the ETag of the last answer, 304 while nothing changed |
| `GET /api/inbox/<name>/<job>` | an item, by its job's number: who sent it, its name, its pages, its state |
| `GET /api/inbox/<name>/<job>/thumbnail?size=&mode=` | its first page as PNG: the whole page in `size` pixels on its longer side (256; 32 to 1024), or with `mode=crop` the square at its top |
| `GET /api/inbox/<name>/<job>/pages/<p>.png?w=` | a page, `w` pixels wide (800) |
| `GET /api/inbox/<name>/<job>/document` | the document it holds, as it is held; each download is in the audit log |
| `POST /api/inbox/<name>/<job>/release` | `{"to": ["<node>", …], "worker": "…"}`: sends it on; no `to`, to every node the inbox releases to. 409 once it is no longer held |
| `POST /api/inbox/<name>/<job>/discard` | `{"worker": "…"}`: ends it |

The thumbnail and the pages are of a `.bdf` document: the item of a file kept as received (an uploader with `output: file`) or of data has none (404), only its document. The OpenAPI document describes an inbox's API, with the nodes it releases to ([OpenAPI](api.md#openapi)).

```sh
TOKEN=$(curl -s -u 'approvals:<secret>' -d grant_type=client_credentials https://<host>/oauth/token | jq -r .access_token)
curl -H "Authorization: Bearer $TOKEN" https://<host>/api/inbox/review
curl -H "Authorization: Bearer $TOKEN" -o 12.png "https://<host>/api/inbox/review/12/thumbnail?size=256"
curl -X POST -H "Authorization: Bearer $TOKEN" -d '{"to": ["laser"], "worker": "Sato"}' https://<host>/api/inbox/review/12/release
```

## Deliveries {#deliveries}

A delivery is what the manager sends on its own: a job to a printer or to an archive, and the step of a template or a convert node. It is made when the job is routed, and done right after by the process, or by the tick. Admin → Deliveries lists the deliveries waiting or failed, with the job, the node, the printer (of a class, the member tried), the state, how many tries, when the next is due, and the last error.

- A delivery that fails — a printer out of reach, a store that does not answer — is tried again after a wait that grows: 30 seconds, then a minute, two, four, eight, then every 15 minutes, up to 20 tries, about four hours; a class tries its next member. After the last try the delivery fails for good, and the job is `aborted`.
- A failure no retry mends — a template missing, data that cannot be read, a format with no converter — fails for good at once.
- **Retry now** sends a delivery again now: one waiting for its time, or one that failed for good.

## The tick {#tick}

The tick does what is due: it sends the deliveries that wait or are to be tried again, prints on paper the slips nobody opened in time (and every unopened slip of a display none of whose devices is connected), removes the slips expired (their display's keep after they came) and the items ended longer ago than their inbox keeps them, the jobs that ended and nothing shows any more, and the enrollment codes expired. A process ticks every 30 seconds (the operator's `tick_interval`); on a serverless platform a scheduler calls `POST /api/tick` every minute ([Serverless](operate.md#serverless)). **Tick now** on Admin → Deliveries runs it and reports what it did.
