# Node kinds

Every node has a kind, and the settings of its kind. The entries — uploaders, CUPS servers and S3 triggers — take documents and data from outside; queues, classes, filters, asserts, branches and inboxes route; template and convert nodes make a document; sign nodes sign it; displays, printers, CUPS servers' queues, archives and sinks are where a job ends; a downloader after an inbox or an archive is where applications and people fetch it, and an HTTP (push) node sends it to a web service. IPP enters at CUPS server nodes alone ([IPP](entries.md#ipp)): computers print to them, applications post files to uploaders, and a bucket's events reach S3 triggers; every other node is reached along the edges. A node is put before the nodes after it: no edge below (`to`, a class's `members`) is needed to add it, and a job that reaches an entry, a queue or a step that sends to no node yet gets an error there.

## Every node {#common}

| Setting | |
|---|---|
| Name | the system's name of the node: letters, digits, `.`, `_` and `-`, no other node's. Changing it renames the node, and its edges follow ([Names, titles and descriptions](topology.md#names)) |
| Title | what people read, shown in the node's title bar |
| Description (`description`) | a note of what the node is for, up to 1000 characters: its first lines in the node on the canvas, all of it in the sidebar |
| Scope, Owner | who may send to the node from the pages, and whose its jobs are: `shared` (everyone signed in), `group` with the group as owner, or `personal` with the person's email as owner ([Scopes of nodes](roles.md#scopes)) |
| On error (`on_error`) | the nodes a job goes to, with what reached this node, when its send or step fails for good: a [notify](#notify) node, say. None: the manager's default (`docflow.on_error`), or the job stops here, aborted ([On error](#on-error)) |

A line in the sidebar tells what the node takes and what it sends on ([Kinds of payload](flow.md#kinds)).

## Uploader {#upload}

The kind `upload`. A REST entry for applications: `POST /api/upload/<name>` with the uploader's own key, or an app's token, takes a document as the body, or data for a template node, and forwards the job to every node of `to`. [Uploaders](entries.md#uploaders) has the requests.

| Setting | |
|---|---|
| Forwards to (`to`) | every node chosen gets the job; none yet: a job gets an error here |
| Accepts (`accept`) | `file`: documents (the default); `data`: JSON, YAML, CSV, an Excel workbook or a Parquet file, for a template node after it |
| Sends on (`output`) | of a document: `bdf`, the `.bdf` made of it (the default), or `file`, the file as received, for a convert node or an archive after it |
| Formats (`formats`) | the formats it takes, as the converter names them (`pdf`, `docx`, `image`, `csv`, `bdf`, `json`, `yaml` …); none takes any. A file of another format is refused (415), and the refusal is in the audit log |
| Metadata (`meta`) | keys the uploader adds to every job it takes; they win over the sender's ([Metadata](flow.md#metadata)) |
| Apps allowed (`apps`) | the apps that may send here with their own credentials, besides the key ([Apps](entries.md#apps)) |

The key is shown once, when the uploader is made or the key rotated; the manager keeps only its hash.

## CUPS server (source) {#cups}

The IPP entry, for computers and business systems, and the only one: a CUPS server on a port of the process's own, with one printer, the node itself. `lpstat -h <host>:<port> -p` lists it, computers print to `ipp://<host>:<port>/printers/<name>`, and what it receives goes on to every node of `to` as the `.bdf` made of it. The node is also an IPP printer at `/ipp/print/<name>` of the manager's own server: on a serverless deployment, which has no ports, the only one. [CUPS server nodes](entries.md#cups) has the setup.

| Setting | |
|---|---|
| Forwards to (`to`) | every node chosen gets what is printed here, made `.bdf` |
| Port (`port`) | the port the server listens on; none given: the first port no node has, from 6310 (none on a serverless deployment) |
| TLS (`tls_cert`, `tls_key`) | the files of a certificate and its key: `ipps`; without them, plain `ipp` |
| Bonjour (`dnssd`) | advertises the printer on the local network with DNS-SD (Bonjour) as an AirPrint printer: Macs list it under Add Printer → Default, iPhones and iPads in Print. With a port only (none on a serverless deployment); the sidebar checks it for a new CUPS server |
| Metadata (`meta`) | as an uploader's. The job attributes of a Print-Job request are metadata, `ipp.media`, `ipp.copies` … |
| Apps allowed (`apps`) | the apps that may print here with their id and secret as IPP Basic ([Apps](entries.md#apps)) |

What computers print is a document: a CUPS server has no Accepts, Sends on or Formats of an uploader, and sends on the `.bdf` it makes. A file to send on as received goes to an uploader or an S3 trigger.

## Queue {#queue}

Forwards the job to every node of `to`: a virtual printer in front of others, and a fan-out when there are several. What reaches it goes on as it is.

| Setting | |
|---|---|
| Forwards to (`to`) | the nodes |

## Class {#class}

Sends the job to one of its `members`: in turn (`round-robin`, the default), or the first (`failover`). A member that does not take the job — a printer out of reach — is skipped for the next member at the next try.

| Setting | |
|---|---|
| Members (`members`) | the nodes, in order; none yet: a job gets an error here |
| Policy (`policy`) | `round-robin` or `failover` |

## Filter {#filter}

Sends the job on by rules: the first rule whose match fits the job forwards it to the rule's nodes; a job no rule fits goes to `otherwise`, or ends there (`dropped`) when there is none. What reaches it goes on as it is. [Rules of a filter](flow.md#rules) has what a rule looks at.

| Setting | |
|---|---|
| Rules (`rules`) | in order: a match (name, user, format, pages, metadata, kind of payload) and the nodes |
| Otherwise (`otherwise`) | the nodes a job no rule fits goes to |

## Assert {#assert}

Checks the job's metadata: a job that passes every check goes on to every node of `to`, as it reached the assert; one that fails a check goes no further, and goes to `on_error` with what failed. [Checks](flow.md#checks) has what a check looks at.

An assert right after the node a job is sent to — an entry, or a queue sent to from the pages or the API — or that node itself is the sender's gate: the job is checked before it goes anywhere, and the sender's request is answered only once its checks pass. A job that fails is refused to the sender with what failed — 422 of REST, `client-error-bad-request` of IPP ([Uploaders](entries.md#uploaders)) — and kept, aborted, with a leg `refused` at the assert. An assert further on checks the job as it reaches it, after the sender was answered: a job it refuses is aborted there.

| Setting | |
|---|---|
| Passes to (`to`) | the nodes a job that passes goes to; none yet: a job that passes gets an error here |
| Checks (`checks`) | every one must hold: a key, and what its value must be |
| On fail (`on_error`) | the nodes a job that fails goes to, with what reached the assert: a [notify](#notify) node tells `{{meta.$error}}`, what failed. None: the manager's default (`docflow.on_error`), or the job stops here |

On the canvas an assert has two ports: `pass`, its `to`, and `fail`, its `on_error`.

## Branch {#branch}

A switch: the first of its cases whose checks the job passes sends it to the case's nodes, and the others are not tried; a job no case fits goes to `default`, or ends there (`dropped`) when there is none. What reaches it goes on as it is. A case checks as an assert does ([Checks](flow.md#checks)).

| Setting | |
|---|---|
| Cases (`cases`) | in order: a name (`name`, the label of its port; `case <n>` when it has none), the checks (`when`), every one of which must hold, and the nodes (`to`). A case checks something; a case with no node yet gives a job that fits it an error |
| Default (`default`) | the nodes a job no case fits goes to |

On the canvas a branch has a port for each case, labelled with its name, and one for `default`: drag from a port to add a node to that arm. The leg names the arm the job took.

## Template {#template}

Fills one of its templates with the data that reaches it, and sends the form on as a `.bdf` document to every node of `to`. It takes data only: JSON, YAML, CSV, an Excel workbook. [Templates and data](templates.md) has the templates and the data.

| Setting | |
|---|---|
| Forwards to (`to`) | the nodes the filled form goes to |
| The data as it came to (`original`) | the nodes the data goes to as it reached the node, traced with the form ([Originals and traces](flow.md#traces)) |
| Template (`template`) | the template to fill, by its id; none yet: upload one once the node is there (a job meanwhile fails here) |
| Template from (`template_from`) | instead: the metadata key whose value names the template (a sender's `doc-type`) |
| List name (`list_name`) | the name data that is a list of records (a CSV file, a sheet, a JSON array) gets in the template; the file's name by default |
| Strict (`strict`) | refuse data that lacks a name the template looks up |

The step is done like a delivery: by the process right after the job is routed, or by the tick; it fails for good, with no retry, when the template is missing or the data cannot be read ([Deliveries](inboxes.md#deliveries)).

## Convert {#convert}

Makes a `.bdf` document of a file that reaches it as received — a PDF, an Office document, a picture — and sends it on to every node of `to`; a `.bdf` goes on as it is. JSON and YAML have no document of their own (fill a template with them), and a CSV or a Parquet file has no pages to print yet. The step is done like a delivery.

| Setting | |
|---|---|
| Forwards to (`to`) | the nodes the document goes to |
| The file as it came to (`original`) | the nodes the file goes to as it reached the node, traced with the document ([Originals and traces](flow.md#traces)) |
| Formats (`formats`) | the formats it converts; none converts any the converter reads |

## Sign {#sign}

Signs the `.bdf` document that reaches it with Docflow's key, as a witness, and sends the signed document on to every node of `to`. The person who sent the job holds no key: Docflow vouches for them. The signature says what Docflow saw — the role, the node, the job, who sent it as Docflow verified it and when it was received — and `bdf verify` shows it. It finds a change made to the document after it; it does not prevent one. The key is Docflow's (`DOCFLOW_SIGN_KEY_FILE` and `DOCFLOW_SIGN_CERT_FILE`, with `DOCFLOW_SIGN_TSA` for a time-stamp); without one the node fails its jobs. The step is done like a delivery.

| Setting | |
|---|---|
| Forwards to (`to`) | the nodes the signed document goes to |
| Role (`sign_role`) | the role the signature says Docflow signs in (`witness` when empty) |

## Inbox {#inbox}

Holds the job until someone releases it to nodes, or discards it: a person on the Inboxes page ([Inboxes](inboxes.md#inboxes)), or an application through the API at `/api/inbox/<name>` ([The inboxes' API](inboxes.md#api)). Released and discarded items stay shown for `keep`. A [downloader](#download) in its `to` lists each job while the inbox holds it (a stocked file), for applications to fetch, and no longer once it is released or discarded. Who works it is its scope's people and its apps ([Who works an inbox](inboxes.md#workers)).

| Setting | |
|---|---|
| Releases to (`to`) | the nodes an item may be released to; none: any node, which only the role topology releases to. Its downloaders are not among them |
| Keep (`keep`) | how long ended items stay (14 days) |
| Apps allowed (`apps`) | the apps that work it with their token, with the scope `inbox` |
| Scope and owner | of a `group` inbox, its members work it; of a `personal` one, its owner; a `shared` one is the role topology's to work |

## Display {#display}

Shows the job as a slip on the devices enrolled at the display — the tablets of a location — and prints on paper, on its fallback printer, what nobody opened within `fallback.after`, or every unopened slip while no device is connected. It takes `.bdf` documents only. [Display boards and devices](displays.md) has the devices and the board.

| Setting | |
|---|---|
| Fallback printer, after (`fallback`) | the printer, and how long an unopened slip waits before it is printed; 0: only while no device is connected |
| Keep (`keep`) | how long slips done or canceled stay on the board (14 days) |
| Offline after (`offline_after`) | how long a device may be silent and still count as connected (2 minutes) |

## Printer {#printer}

Prints the job on paper, on a printer the manager reaches over the network. It takes `.bdf` documents only. A print that fails is tried again ([Deliveries](inboxes.md#deliveries)).

| Setting | |
|---|---|
| URL (`url`) | `ipp://` or `ipps://` (an IPP printer), `socket://host:9100`, or `lpd://host/queue`. A queue of a CUPS server is a [CUPS server (sink)](#cups-out) |
| Route (`route`) | how its jobs reach it: `direct`, the manager sends them (the default). `branch`, a `bdf printd` of the site fetching them, is not there yet: a delivery to such a printer fails |
| Preview (`preview`) | the URL of this printer on a fake printer farm's API (`tools/fakeprinter`): what the fake printer received is shown beside the node, as the farm read it back |

## CUPS server (sink) {#cups-out}

Prints the job to a queue of a CUPS server over IPP, as a computer prints to it: a printer or a class there. The CUPS server's filters and drivers make what its printer takes, of the PDF (or PWG Raster) the manager sends, as the queue says what it supports. It takes `.bdf` documents only; a print that fails is tried again ([Deliveries](inboxes.md#deliveries)). Its counterpart, a [CUPS server (source)](#cups), takes what computers print.

| Setting | |
|---|---|
| Server (`url`) | the queue's URL: `ipp://host:631/printers/<queue>`, `/classes/<class>` for a class, `ipps://` over TLS. In the sidebar, Find queues asks the server for its queues (CUPS-Get-Printers) to pick one |
| Credential (`credential`) | `user:password`, when the server asks for them; none: no authentication ([Credentials](#credentials)) |

The keys `ipp.job-account-id` and the others go on to the queue as they go on to a printer ([What goes on to the printer](flow.md#ipp-attributes)).

## Archive {#archive}

Writes what reaches it — a `.bdf` document, or any file or data as received — and a JSON record of the job (its metadata, payloads and legs) to a store, each object once; the manager never changes or deletes them. Give the bucket a default retention (S3 Object Lock) for an audit trail nobody can change. Once written, the job goes on to every node of `to`: a downloader, which lists it ([Downloaders](downloads.md)).

| Setting | |
|---|---|
| Store (`store`) | a documents store: `bbolt:path`, `s3://bucket?region=…`, `r2://bucket?endpoint=…`, `gcs://bucket` |
| Collection (`collection`) | the prefix of its objects: `archive-<name>` by default |
| Forwards to (`to`) | the nodes the job goes on to once written: downloaders |

## Sink {#sink}

Keeps the latest documents that reach it, to look at on the topology: the first page shows on the node, every page in the sidebar. An older one goes when a newer comes. A sink is where a document is seen without a printer or a device: behind an uploader or a CUPS server while they are set up.

| Setting | |
|---|---|
| Latest (`latest`) | how many it keeps (10) |

## Downloader {#download}

The kind `download`. Lists what an inbox or an archive right before it keeps — a stocked file, the only kind it takes — for applications and people to fetch at `GET /api/download/<name>`: each entry under the owner its metadata names. An archive's stays for good; an inbox's while the inbox holds the job. An edge to it from any other node is refused, and nothing is printed to it. [Downloaders](downloads.md) has who reads what, backends and frontends. To send documents to a web service instead, the manager making the request, use [HTTP (push)](#http).

| Setting | |
|---|---|
| Owner's key (`by`) | the metadata key whose value says whose a document is: a customer's id or email the sender gave. Empty: a person reads what they could see of the job |
| Subject (`subject`) | what of a person is compared with the owner: `email` (the default), `sub`, `id` or `group` ([Whose a document is](downloads.md#owners)) |
| Fields (`fields`) | metadata keys copied into the entries listed, for the reader to show or pick by (up to 16) |
| Entries (`entries`) | how many of the latest entries are listed for everyone and for each owner (1000); an older entry is still read by its job's number |
| Apps allowed (`apps`) | the backends that read every entry with their token (the scope `download`); a frontend reads as the person and needs no listing |

## S3 Trigger {#s3-trigger}

An entry for the objects put in a bucket. It does not listen to the bucket: an S3 event goes to an EventBridge rule (an API destination) or a Lambda, which posts the object's reference to `POST /api/s3-trigger/<name>` with the node's key, made and rotated as an uploader's. The body is `{"bucket": …, "key": …}`, an S3 notification (`Records`, one object), or the event EventBridge sends (`detail`). The node reads the object itself, and the job goes on with `s3.bucket` and `s3.key` in its metadata.

| Setting | |
|---|---|
| Bucket (`url`) | `s3://bucket?region=…`; `&endpoint=…` for a service other than AWS's (MinIO, R2) |
| Credential (`credential`) | `access-key:secret-key`, or with `:session-token`; none takes AWS's own (the environment, the role) ([Credentials](#credentials)) |
| Prefix (`path`) | the keys it takes are under this prefix; a key elsewhere, or another bucket, is refused (403) |
| Forwards to, Accepts, Sends on, Formats, Metadata | as an uploader's |

## Render {#render}

Draws the `.bdf` document that reaches it as a PDF, or one page of it as a PNG image (150 dpi), for an exit whose far side has no viewer of BDF: a mail's reader, a web service. The file goes on to `to` as any other file. Done like a delivery, in the background.

| Setting | |
|---|---|
| Forwards to (`to`) | the nodes the file goes to |
| Format (`format`) | `pdf` (the default) or `png`; SVG is not there yet |
| Page (`page`) | the page of a PNG (1) |

## Notify {#notify}

Tells that a job passed it, and sends what reached it on to `to`, if any. Failures reach it along [on_error](#on-error) edges, with the error. It sends no document.

| Setting | |
|---|---|
| Via (`via`) | `slack`: an incoming webhook, whose URL is the credential; `mail`: a mail, with a Mail node's settings and no attachment; `http`: a JSON object (`text`, `job`, `name`, `node`, `error`, `error_node`, `meta`) posted to `url`, the credential its Authorization |
| Text (`text`) | what it says, with `{{meta.<key>}}`; by default the job and the node, or of a failure the node and the error |
| Forwards to (`to`) | the nodes the job goes on to |

## Mail (SMTP) {#smtp}

Mails what reaches it, as an attachment. The Message-ID is the delivery's id, the same on every try.

| Setting | |
|---|---|
| Server (`url`) | `smtp://host:587` (STARTTLS when the server offers it) or `smtps://host:465` |
| Credential (`credential`) | `user:password`, given only over TLS or to the machine itself |
| From, Recipients, Subject, Text (`from`, `recipients`, `mail_subject`, `text`) | the mail, with `{{meta.<key>}}` |

## S3 {#s3}

Writes what reaches it to a bucket, once, and does not look back: an object there already is taken as written. What must be kept, with a record, goes to an [archive](#archive).

| Setting | |
|---|---|
| Bucket (`url`), Credential (`credential`) | as an S3 Trigger's |
| Key (`path`) | with `{{meta.<key>}}`: `{{meta.$job}}/{{meta.$name}}` by default; a value cannot leave the prefix (no `..`, no leading `/`) |

## HTTP (push) {#http}

The kind `http`. Sends what reaches it to a web service in a request, the document and the job's metadata as `body` says, with the delivery's id as the `Idempotency-Key`. To another Docflow, the URL is its `/api/upload/<node>` and the credential its key or an app's token: the `.bdf` goes there as it is, with its name and its metadata as the sender's, in either body ([Uploaders](entries.md#uploaders)). 408, 429 and 5xx are tried again; any other answer fails for good. The other way round, applications fetching documents, is [downloader](#download).

| Setting | |
|---|---|
| URL (`url`) | `http(s)://…`, with `{{meta.<key>}}` (each value URL-escaped) |
| Method (`method`) | `POST` (the default) or `PUT` |
| Body (`body`) | `direct` (the default): the document is the body as it is, with its `Content-Type` and its name in `Content-Disposition`, and the metadata goes in `X-Meta-<key>` headers. `multipart`: the body is `multipart/form-data`, the metadata a JSON object in the part `meta` (its values as given), then the document in the part `file`, with its name and type. Either way, the manager's keys (`$…`) are not sent |
| Headers (`headers`) | more headers, with `{{meta.<key>}}`; with a `multipart` body, not `Content-Type` or `Content-Disposition` |
| Credential (`credential`) | the Authorization: a token, sent as Bearer, or a whole value (`Basic …`) |

Render, Notify, Mail (SMTP), S3 and HTTP have no form of their own in the sidebar: their settings are a JSON box there. An S3 trigger has a form, as an uploader has.

## On error {#on-error}

A send or a step that fails for good — a refusal, a missing template, the last of its tries — goes along the node's `on_error` edges, drawn dashed, with what reached the node and `$error` and `$error_node` in what the next nodes see. A node without them sends its failures to the manager's default (`docflow.on_error`, `DOCFLOW_ON_ERROR`), unless the job is on an error's way already, so that a notify node that fails does not tell of itself. With neither, the job stops there, aborted. The job is aborted either way: the edges tell, they do not mend.

## Credentials and `{{meta.<key>}}` {#credentials}

A node names its credential (`credential`); the secret is never stored with it. The value is the environment's `DOCFLOW_CRED_<name>` (`-` as `_`), else the file `<name>` in `docflow.credential_dir`. The settings of the exits and notify nodes may use `{{meta.<key>}}` — nothing more, no expressions — with the job's metadata ([Metadata](flow.md#metadata)) and `$node`, `$delivery`, `$name` (the document's name), `$error` and `$error_node`; a key not there is empty. A header loses its line breaks, a key its `..`. Every exit sends at least once: a send whose outcome is not known is sent again, with the same delivery id, for the far side to tell.
