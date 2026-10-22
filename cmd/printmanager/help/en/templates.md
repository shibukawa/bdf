# Templates and data

A template node fills a form with data and sends the form on as a `.bdf` document: an application sends the data of an invoice, and the printers get the invoice. The templates are the node's own.

## The templates of a node {#templates}

A template is uploaded in the node's sidebar on Admin → Topology (Templates: an id and a file), or with `PUT /api/nodes/<node>/templates/<id>` with the file as the body ([Templates](api.md#templates)). It goes with the node: deleting the node deletes its templates.

A template is one of what `bdf fill` takes:

| Template | |
|---|---|
| `.bdft` | a compiled template, from `bdf template compile` |
| Excel, Word, PowerPoint | an Office document with tags in its cells or its text |
| a PDF form | its fields filled |
| pdfme | a pdfme template (JSON) |
| text | a text template, for printers in their text mode |

The sidebar lists each template with its format, how many fields its tags look up (their names when hovering), its file, and Delete. `GET /api/nodes/<node>/templates/<id>` tells the fields, and a sample of the data when the template has one: a start for what a sender sends.

## Which template is filled {#which}

The node's `template` names the template to fill. Instead, `template_from` names a metadata key whose value names it: with `template_from` set to `doc-type`, a sender's `X-Meta-Doc-Type: invoice` header picks the template `invoice` ([Metadata](flow.md#metadata)). The filling is a step of the node, done like a delivery; a template that is missing is a failure no retry mends, and the job is `aborted` with the reason in its leg ([Deliveries](inboxes.md#deliveries)).

## The data {#data}

The data reaches the node as a payload of the kind `data`: from an endpoint with `accept: data` ([Endpoints](entries.md#endpoints)), or sent to the node itself through the API. A JSON or YAML object is the data: its names are what the template's tags look up. A list of records — a CSV file, a sheet of an Excel workbook, a JSON array — is given a name, the node's `list_name` or the file's name without its extension, and the template looks the records up under it (a table that repeats a row per record). With `strict`, data that lacks a name the template looks up is refused instead of filled with nothing. A Parquet file is not filled into a template yet.

The filled form is a `.bdf` document, a new payload of the job: the printers, the displays and the archives after the node get it, and the Jobs pages show it as the job's document ([Payloads](flow.md#payloads)).
