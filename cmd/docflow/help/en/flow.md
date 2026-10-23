# What flows

What a job carries along the edges and what goes with it: its payloads and their kinds, its metadata, and the rules of a filter, which look at them.

## Kinds of payload {#kinds}

What flows along an edge is of one of four kinds:

| Kind | | Color |
|---|---|---|
| `bdf` | a `.bdf` document: what printers, displays and the pages show | blue |
| `data` | data for a template: JSON, YAML, CSV, an Excel workbook or a Parquet file | green |
| `any` | any other file as it was received — a PDF, an Office document, a picture — for a convert node or an archive | amber |
| `stocked` | a stocked file, of any format: what an inbox holds or an archive wrote, read from there by an [HTTP (pull)](exports.md) node, which takes nothing else. Where it is kept is where it comes from and how long it is there: an inbox's while it holds the job, an archive's for good | cyan |

What a node sends on:

- an entry (an endpoint, a CUPS server) sends what its settings say: with `accept: data`, data; with `output: file`, the file as received; else the `.bdf` made of the document (the default);
- a template node and a convert node send the `.bdf` they made, and along `original` what reached them, as it was;
- a display sends what its fallback printer prints, `.bdf`;
- an inbox and an archive send a stocked file to the HTTP (pull) nodes of their `to`, and what reached them to the others;
- every other node sends on what reached it.

What a node takes: a printer and a display take `bdf` only; a template node takes `data` only; every other node takes anything. What may be sent to a node now is narrower: what it takes, within what every node after it takes. A queue in front of a printer takes `.bdf` alone, as it would carry nothing else that the printer could print; a queue, a class, a filter or an inbox with no edge yet takes the kind of its first edge, on either side. The canvas colors the ports and the edges by these kinds — blue `bdf`, green `data`, amber `any`, purple several — and draws a port not typed yet hollow and dashed.

When a node is saved, the kinds that reach every node are inferred from the entries along the edges, and a new edge that would bring a kind a node does not take is refused; the message names the path, "laser (printer) takes bdf, but any reaches it: raw → fan → laser". `GET /api/topology` shows each node's `takes`, `emits`, `accepts` (what may be sent to it now) and `broken`.

## Broken edges {#broken}

A change of a node's settings alone — an entry that takes data now, a queue made a template node — is saved even when it leaves an edge carrying what the node after it does not take. The edge stays, broken: the canvas draws it dashed and red with a mark, the sidebar of the node it leaves says why, and the save says which edges it broke. Mend it by putting a convert or a template node between, by changing the setting back, or by removing the edge.

A job's payload is checked again at every node it reaches, so that a broken edge, or a topology changed while the job was on its way, does not break the job: a node that cannot take the payload gets a leg in `error`.

## Payloads {#payloads}

A job carries its payloads: what was received — its `kind`, `format`, `pages`, `size`, `sha256`, and the file's `name` when the sender gave one — and what the template and the convert nodes made of it, each saying which node made it. Each leg says which payload the node got and which it made. The job's document, on the Jobs pages, in a display's slip and at a printer, is the last `.bdf` made, or the one received; `GET /api/jobs/<n>/document` downloads it, and `?payload=<id>` any other.

The formats are the converter's names — `bdf`, `pdf`, `docx`, `xlsx`, `pptx`, `image` … — and `json`, `yaml`, `csv`, `xlsx` and `parquet` for data; `GET /api/formats` lists what an entry may take. What is sent is told by its name or its content type (JSON, YAML, CSV), or by its content. An Excel workbook or a CSV file is data when it is passed on as received (an entry with `accept: data` or `output: file`), and a document when an entry makes a `.bdf` of it.

## Originals and traces {#traces}

A template or a convert node sends the `.bdf` it makes along `to` and, along `original`, what it took as it was — the data, the file received: an archive of originals beside an archive of what was made, say. An edge along `original` carries what reaches the node, so a file received does not go to a printer that way; the canvas draws it with a short dash.

Both sides go on with a `trace` of the step:

| | |
|---|---|
| `id` | the step's id, the same on both sides |
| `node`, `kind` | the template or convert node |
| `role` | the side: `source` along `original`, `output` along `to` |
| `source`, `output` | both payloads — `payload`, `kind`, `format`, `name`, `sha256` — with the `key` and the `object` an archive names each by (`<received day>.<job id>.<sha256 head>`; the record is `<key>.json`) |
| `fingerprints` | every step of the payload's lineage, oldest first: `id`, `node`, `role` |

An archive writes the trace into its record, and an export into its entries, so that from either object the other is found by its name in the archive it went to. The manager does not look for it, nor know whether it was written. The trace goes with the payload through queues, filters, classes, inboxes and an archive, until the next template or convert node, whose trace keeps the fingerprints of the old: two objects many steps apart — a PDF received, and the stamped `.bdf` made of the `.bdf` made of it — are matched by the step ids their fingerprints have in common. Both sides are of the same job, so an export after each archive lists them under the same job number.

## Metadata {#metadata}

Metadata goes with a job through the topology. It comes from three sources, each key once, the later winning:

1. **The sender's.** `X-Meta-<Key>` headers (the key in lower case) and `?meta.<key>=` parameters of a request to an endpoint or to the API, and the JSON object in the part `meta` of its `multipart/form-data` body ([Endpoints](entries.md#endpoints)), and the job attributes of an IPP request as `ipp.<name>` (`ipp.media`, `ipp.copies`, `ipp.job-name` …). These are claims: a filter routes by them and a template node picks its template by them, as the administrator set up; nothing is authorized by them.
2. **The entry's.** The `meta` of the endpoint or the CUPS server the job came in at.
3. **The manager's**, under keys that start with `$`, which a sender cannot write: `$job` (the number), `$received` (when), `$entry` (the node the job was sent to), `$source` (the way it came in), `$sender` and `$sender_kind` (who, as verified: a person, an app, a token, an endpoint, anonymous), `$file` (the job's name).

A key is letters, digits, `_`, `.`, `-` and `/`, up to 64 characters; a job has at most 64 keys and 16 KiB of metadata. The Jobs pages show a job's metadata when hovering its name, and an archive's record keeps it.

At a node, a filter's rules and a template node's `template_from` see the job's metadata, its `name` and `user`, and the payload that reached the node: `$kind`, `$format`, `$pages`, `$size`, `$sha256` and `$payload` (its id), and the trace of the step it came out of or went into: `$trace`, `$trace_role`, `$trace_node`, `$trace_peer` and `$fingerprints` ([Originals and traces](#traces)).

### What goes on to the printer {#ipp-attributes}

A few keys go on to the printer as the IPP attributes of the same names, when a printer node prints the job: `ipp.job-account-id` and `ipp.job-accounting-user-id` (the account the print is charged to, PWG 5100.7), `ipp.document-metadata` (`key=value` strings with Dublin Core's keys or keys that begin with `x-`, PWG 5100.13) and `ipp.x-…` (attributes of one's own). An IPP request's own attributes give them, or the metadata of a request to an endpoint or the API: `?meta.ipp.job-account-id=sales`. A value the attribute does not take is left out; whether the printer uses the rest is the printer's.

## Rules of a filter {#rules}

A rule is a match and the nodes a job that matches goes to. Every condition of the match must hold:

| Condition | |
|---|---|
| name | a regular expression (RE2) on the job's name |
| user | on the name the sender gave, and on the verified sender's label |
| format | on the format received: `pdf`, `docx` … |
| min pages, max pages | bounds on the pages of the payload |
| metadata | by key, a regular expression on the value; a key the job lacks does not fit |
| kind of payload | `bdf`, `data` or `any` |

The rules are tried in order, and the first that fits sends the job on; the leg names the rule. A job no rule fits goes to the filter's otherwise, or ends there, `dropped`, when there is none. A rule with no condition fits every job.

## Checks {#checks}

An [assert](nodes.md#assert) and the cases of a [branch](nodes.md#branch) look at the job by checks. A check is a key of the job's metadata, as the nodes see it ([Metadata](#metadata)): a sender's key (`doc-type` of `X-Meta-Doc-Type`, `ipp.job-account-id` of IPP), an entry's, or the manager's (`$pages`, `$format`, `$sender` …), and what its value must be. Every condition given must hold:

| Condition | |
|---|---|
| key (`key`) | the key must have a value that is not empty |
| optional (`optional`) | instead: a job without the key passes, and the conditions below are of a value it has |
| matches (`pattern`) | a regular expression (RE2) the value matches somewhere in it; `^…$` for all of it |
| min, max (`min`, `max`) | the value is a number within the bounds, both included |
| one of (`enum`) | the value is one of these, as written |
| message (`message`) | of an assert: what a failure says, instead of the check's own words |

A failure says the key and why: `doc-type: required`, `amount: 12000 is more than 10000`, `doc-type: "memo" is not one of invoice, order`. An assert checks every check, so that the sender learns all that failed at once; the job's leg at the assert and `{{meta.$error}}` of a notify node after it are those failures, one after another. The OpenAPI document of an entry tells the keys the asserts right after it check, with their constraints, required unless optional ([OpenAPI](api.md#openapi)).
