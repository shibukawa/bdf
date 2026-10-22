# Node kinds

Every node has a kind, and the settings of its kind. The entries — endpoints and CUPS servers — take documents and data from outside; queues, classes, filters and inboxes route; template and convert nodes make a document; displays, printers, archives and sinks are where a job ends. Every node, whatever its kind, is also an IPP printer at `/ipp/print/<name>` ([IPP](entries.md#ipp)).

## Every node {#common}

| Setting | |
|---|---|
| Name | the node's name ([Names and titles](topology.md#names)) |
| Title | shown instead of the name |
| Scope, Owner | who may send to the node from the pages, and whose its jobs are: `shared` (everyone signed in), `group` with the group as owner, or `personal` with the person's email as owner ([Scopes of nodes](roles.md#scopes)) |

A line in the sidebar tells what the node takes and what it sends on ([Kinds of payload](flow.md#kinds)).

## Endpoint {#endpoint}

A REST entry for applications: `POST /in/<name>` with the endpoint's own key takes a document as the body, or data for a template node, and forwards the job to every node of `to`. [Endpoints](entries.md#endpoints) has the requests.

| Setting | |
|---|---|
| Forwards to (`to`) | every node chosen gets the job |
| Accepts (`accept`) | `file`: documents (the default); `data`: JSON, YAML, CSV, an Excel workbook or a Parquet file, for a template node after it |
| Sends on (`output`) | of a document: `bdf`, the `.bdf` made of it (the default), or `file`, the file as received, for a convert node or an archive after it |
| Formats (`formats`) | the formats it takes, as the converter names them (`pdf`, `docx`, `image`, `csv`, `bdf`, `json`, `yaml` …); none takes any. A file of another format is refused (415), and the refusal is in the audit log |
| Metadata (`meta`) | keys the endpoint adds to every job it takes; they win over the sender's ([Metadata](flow.md#metadata)) |
| Port (`port`) | a port of the process's own: every POST to it is to this endpoint (none on a serverless deployment) |
| Apps allowed (`apps`) | the apps that may send here with their own credentials, besides the key ([Apps](entries.md#apps)) |

The key is shown once, when the endpoint is made or the key rotated; the manager keeps only its hash.

## CUPS server {#cups}

An IPP entry for computers and business systems: a CUPS server on a port of the process's own, with one printer, the node itself. `lpstat -h <host>:<port> -p` lists it, computers print to `ipp://<host>:<port>/printers/<name>`, and what it receives goes on to every node of `to`. [CUPS server nodes](entries.md#cups) has the setup.

| Setting | |
|---|---|
| Forwards to (`to`) | every node chosen gets what is printed here |
| Port (`port`) | the port the server listens on (none on a serverless deployment) |
| TLS (`tls_cert`, `tls_key`) | the files of a certificate and its key: `ipps`; without them, plain `ipp` |
| Accepts, Sends on, Formats, Metadata | as an endpoint's. The job attributes of a Print-Job request are metadata, `ipp.media`, `ipp.copies` … |

## Queue {#queue}

Forwards the job to every node of `to`: a virtual printer in front of others, and a fan-out when there are several. What reaches it goes on as it is.

| Setting | |
|---|---|
| Forwards to (`to`) | the nodes |

## Class {#class}

Sends the job to one of its `members`: in turn (`round-robin`, the default), or the first (`failover`). A member that does not take the job — a printer out of reach — is skipped for the next member at the next try.

| Setting | |
|---|---|
| Members (`members`) | the nodes, in order |
| Policy (`policy`) | `round-robin` or `failover` |

## Filter {#filter}

Sends the job on by rules: the first rule whose match fits the job forwards it to the rule's nodes; a job no rule fits goes to `otherwise`, or ends there (`dropped`) when there is none. What reaches it goes on as it is. [Rules of a filter](flow.md#rules) has what a rule looks at.

| Setting | |
|---|---|
| Rules (`rules`) | in order: a match (name, user, format, pages, metadata, kind of payload) and the nodes |
| Otherwise (`otherwise`) | the nodes a job no rule fits goes to |

## Template {#template}

Fills one of its templates with the data that reaches it, and sends the form on as a `.bdf` document to every node of `to`. It takes data only: JSON, YAML, CSV, an Excel workbook. [Templates and data](templates.md) has the templates and the data.

| Setting | |
|---|---|
| Forwards to (`to`) | the nodes the filled form goes to |
| Template (`template`) | the template to fill, by its id |
| Template from (`template_from`) | instead: the metadata key whose value names the template (a sender's `doc-type`) |
| List name (`list_name`) | the name data that is a list of records (a CSV file, a sheet, a JSON array) gets in the template; the file's name by default |
| Strict (`strict`) | refuse data that lacks a name the template looks up |

The step is done like a delivery: by the process right after the job is routed, or by the tick; it fails for good, with no retry, when the template is missing or the data cannot be read ([Deliveries](inboxes.md#deliveries)).

## Convert {#convert}

Makes a `.bdf` document of a file that reaches it as received — a PDF, an Office document, a picture — and sends it on to every node of `to`; a `.bdf` goes on as it is. JSON and YAML have no document of their own (fill a template with them), and a CSV or a Parquet file has no pages to print yet. The step is done like a delivery.

| Setting | |
|---|---|
| Forwards to (`to`) | the nodes the document goes to |
| Formats (`formats`) | the formats it converts; none converts any the converter reads |

## Inbox {#inbox}

Holds the job until a person releases it to nodes, or discards it, on Admin → Inboxes ([Inboxes](inboxes.md#inboxes)). Released and discarded items stay shown for `keep`.

| Setting | |
|---|---|
| Releases to (`to`) | the nodes an item may be released to; none: any node |
| Keep (`keep`) | how long ended items stay (14 days) |

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
| URL (`url`) | `ipp://` or `ipps://` (an IPP printer, or a queue of a CUPS server: `ipp://host:631/printers/<queue>`), `socket://host:9100`, or `lpd://host/queue` |
| Route (`route`) | how its jobs reach it: `direct`, the manager sends them (the default). `branch`, a `bdf printd` of the site fetching them, is not there yet: a delivery to such a printer fails |
| Preview (`preview`) | the URL of this printer on a fake printer farm's API (`tools/fakeprinter`): what the fake printer received is shown beside the node, as the farm read it back |

## Archive {#archive}

Writes what reaches it — a `.bdf` document, or any file or data as received — and a JSON record of the job (its metadata, payloads and legs) to a store, each object once; the manager never changes or deletes them. Give the bucket a default retention (S3 Object Lock) for an audit trail nobody can change.

| Setting | |
|---|---|
| Store (`store`) | a documents store: `bbolt:path`, `s3://bucket?region=…`, `r2://bucket?endpoint=…`, `gcs://bucket` |
| Collection (`collection`) | the prefix of its objects: `archive-<name>` by default |

## Sink {#sink}

Keeps the latest documents that reach it, to look at on the topology: the first page shows on the node, every page in the sidebar. An older one goes when a newer comes. A sink is where a document is seen without a printer or a device: behind an endpoint or a CUPS server while they are set up.

| Setting | |
|---|---|
| Latest (`latest`) | how many it keeps (10) |
