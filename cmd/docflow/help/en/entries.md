# Where documents come in

A job enters the network at an entry: the form on the Jobs page, an uploader an application POSTs to, a CUPS server node computers print to, an S3 trigger a bucket's events reach, or the API. Who may send where is [Scopes of nodes](roles.md#scopes).

## From the pages {#pages}

The form on [Jobs](jobs.md#print), and the form of a node on the [Topology](topology.md#view) page, send a document as yourself to a source you may send to — an uploader, a CUPS server or an S3 trigger, the nodes where documents come in — so that what a source takes and sets holds for people as for applications: the job is yours, and its group the node's. A person with the role `topology` sends to any node of theirs from the Jobs page, to try a part of the topology out.

## IPP {#ipp}

IPP enters at CUPS server nodes alone ([CUPS server nodes](#cups)): a computer adds one as an IPP printer, at its own port (`ipp://<host>:<port>/printers/<name>`) or at the manager's own server (`ipp://<host>:<port>/ipp/print/<name>`; `ipps://` when the manager is behind TLS), and the job follows the node's edges. `bdf printers -url ipp://<host>:<port>/ipp/print` lists them, as a CUPS client's `lpstat` does (CUPS-Get-Printers). The other nodes are no IPP printers: a queue, a printer or a display is reached along the edges. The job attributes of the request (`media`, `copies`, `job-name` …) and its `document-metadata` go with the job as metadata, `ipp.<name>`; the account (`job-account-id`, `job-accounting-user-id`), `document-metadata` and attributes of one's own (`x-…`) go on to the printers the job reaches ([What goes on to the printer](flow.md#ipp-attributes)).

Who may print over IPP is the operator's `submit_tokens`: with any set, the client gives a token's name and secret as Basic authentication (an app gives its id and secret the same way); with none, anyone who reaches the manager may send. The requesting-user-name the client gives is shown as "(as …)", never checked.

## Uploaders {#uploaders}

An uploader is a REST entry for an application: a business system, a WMS. It is a node on the topology ([Uploader](nodes.md#upload)); when it is saved for the first time its key is shown once, and Rotate key in its sidebar makes a new one and shows it once (the old one stops working). The manager keeps only the key's hash.

The application sends with the key as a bearer token:

```sh
curl -X POST -H "Authorization: Bearer bdfk_…" \
     --data-binary @statement.pdf \
     "https://<host>/api/upload/erp?name=statement&meta.dept=sales"
```

or, as a form:

```sh
curl -H "Authorization: Bearer bdfk_…" \
     -F 'meta={"dept": "sales"};type=application/json' \
     -F file=@statement.pdf \
     https://<host>/api/upload/erp
```

- The body is the document, or the data for an uploader with `accept: data`; the name (`?name=`) and the `Content-Type` tell its format when the content does not.
- `?name=` names the job; `?user=` names who sent it, shown as "(as …)".
- `X-Meta-<Key>` headers and `?meta.<key>=` parameters are the job's metadata ([Metadata](flow.md#metadata)).
- Or the body is `multipart/form-data`, told by its `Content-Type`: the document in the part `file`, whose name and `Content-Type` are the document's (`?name=` over the name), and the metadata as a JSON object in the part `meta`, over the headers' and the parameters'. An [HTTP (push)](nodes.md#http) node of another Docflow sends either.
- Without `?name=`, a body as it is takes its name from `Content-Disposition` (`attachment; filename="statement.pdf"`).
- The answer is the job (201), with its number and its legs so far. A document of a format the uploader does not take is refused (415), and the refusal is in the audit log; a wrong key is refused (401).
- An [assert](nodes.md#assert) right after the uploader checks the metadata before the answer: a job that fails is refused (422) with `{"error": "…", "job": <number>, "failures": [{"node", "key", "error"}]}`, and kept, aborted, for its legs to show what failed. A CUPS server answers the same as IPP's `client-error-bad-request`, with what failed as the status message.

The uploader's OpenAPI document — the link OpenAPI in its sidebar, or `GET /api/openapi.json?node=<name>` — says what it takes: the media types, the schema of the data its templates fill, and the metadata the nodes after it look at ([OpenAPI](api.md#openapi)).

Everything an application calls is the API of the manager's own server, by the kind of the node and its name: an uploader at `/api/upload/<name>`, an S3 trigger at `/api/s3-trigger/<name>`, an inbox at `/api/inbox/<name>` ([The inboxes' API](inboxes.md#api)), a downloader at `/api/download/<name>` ([Downloaders](downloads.md)). An uploader has no port of its own: only a CUPS server listens on one, for IPP.

## CUPS server nodes {#cups}

A CUPS server node ([CUPS server](nodes.md#cups)) is for computers and business systems that print: a CUPS server on a port of the process's own, with one printer, the node itself. A computer adds the printer with `ipp://<host>:<port>/printers/<name>` (`ipps://` when the node has a certificate), as it adds any CUPS printer, and what it prints there goes on to the node's `to`: a printer makes the node a relay, a sink shows what was printed, a filter keeps only what fits. `lpstat -h <host>:<port> -p` lists it. A new CUPS server gets the first port no node has, from 6310. To print the other way, from this system to a CUPS server's queue, put a [CUPS server (sink)](nodes.md#cups-out) at the end of the edges. The node's Usage tab shows the URL, how macOS, Windows and CUPS add the printer, and the commands.

The printer answers as AirPrint and IPP Everywhere printers do (PDF preferred, Apple and PWG rasters, paper from A3 to postcards), so a Mac makes an AirPrint queue of it, which sends PDF. When adding it by its address, choose AirPrint as the driver: a Generic PostScript Printer sends PostScript, which is not taken, and the Mac keeps retrying the job (each refusal is logged as `ipp: print-job refused`). It is advertised with Bonjour (`dnssd`, on unless it is `false`): Macs, iPhones and iPads on the local network find it as `<title> @ <computer>`, on a host name of the process's own (`<computer>-docflow.local`), while its port is open.

## Apps {#apps}

An app is an application that sends, works inboxes or reads downloaders with credentials of its own rather than an uploader's key: the credentials are the app's, and the audit log and the jobs say which app sent. An administrator registers one on Admin → Apps: an id, a title, the group whose members see its jobs (optional), and its scopes:

| Scope | |
|---|---|
| `print` | send documents |
| `print.template` | fill templates: send data |
| `jobs.cancel` | cancel its own jobs |
| `inbox` | work the inboxes that list it: list what they hold, their thumbnails and documents, release and discard ([The inboxes' API](inboxes.md#api)) |
| `download` | read the downloaders that list it ([Backends](downloads.md#backends)) |
| `tick` | call `POST /api/tick`: a scheduler |
| `admin` | the API as an administrator: a script |

The secret is shown once, when the app is added, and Rotate secret makes a new one (the old one stops working). Disable refuses the app's credentials and the tokens out; Delete removes it.

The app takes a token at `POST /oauth/token`, the OAuth 2.0 client credentials grant, and sends with it as a bearer to the uploaders that list the app among their apps (Apps allowed, in the uploader's sidebar), or prints over IPP with its id and secret as Basic authentication to the CUPS server nodes that list it the same way:

```sh
curl -u 'erp:<secret>' -d grant_type=client_credentials https://<host>/oauth/token
# {"access_token": "…", "token_type": "Bearer", "expires_in": 600, "scope": "print"}
curl -X POST -H "Authorization: Bearer <access_token>" --data-binary @statement.pdf "https://<host>/api/upload/erp?name=statement"
```

A token lasts ten minutes. An app sees its own jobs (`GET /api/my/jobs`) and cancels them with the scope `jobs.cancel`; its jobs are its group's when the node gives them none.

## Scripts {#scripts}

A script sends with the operator's admin token, or a submit token's `name:secret`, as a bearer: `POST /api/destinations/<node>/jobs` with the document as the body takes `?name=`, `?user=` and the metadata as an uploader does ([Jobs](api.md#jobs)).
