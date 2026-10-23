# BDF Docflow

Docflow is a software-defined document network. A document enters at an entry — an application's endpoint, a CUPS server a computer prints to, an S3 trigger a bucket's events reach, or the form on the Jobs page — and flows through the nodes an administrator composed: queues that fan out, classes that send to one member in turn, filters that route by what a job is, inboxes that hold a job until a person sends it on, template nodes that fill a form with data, convert nodes that make a document of a file. It ends on display boards (the tablets enrolled at a location), on printers, in archives that keep a copy nobody changes, in sinks that keep the latest documents to look at, and in exports that applications and people fetch them from. Paper is one destination among others.

## The pages {#pages}

| Page | Who | What |
|---|---|---|
| Jobs | everyone signed in | print a document to a node you may send to; your jobs and your groups' ([Jobs](jobs.md)) |
| Topology | everyone signed in | the network as it is laid out, with the activity of the jobs you may see; print at a node ([Looking](topology.md#view)) |
| Topology, as the editor | the role topology | the same page: nodes placed, connected and set up ([The editor](topology.md#editor), [Node kinds](nodes.md)); Admin in the bar leads there |
| Admin → Apps | the role topology | the applications with credentials of their own: backends that send and read, and frontends in people's browsers ([Apps](entries.md#apps), [Frontends](exports.md#frontends)) |
| Admin → Inboxes | the role topology | release or discard what the inboxes hold ([Inboxes](inboxes.md#inboxes)) |
| Admin → Jobs | the role topology | every job, or those sent to one node, and the route each took |
| Admin → Deliveries | the role topology | what waits for a printer, an archive or a step, and the tick ([Deliveries](inboxes.md#deliveries)) |
| Audit | the role audit | who did what ([Audit](roles.md#audit)) |
| Help | everyone signed in | this manual, in English and Japanese |
| `/tablet` | a device | the display board of a location: no sign-in, the device's own key ([Display boards](displays.md)) |

An administrator has every role. Who gets in, and who has which role, is in [Sign-in, roles and audit](roles.md).

## How a job flows {#flow}

A job is what a sender sent to a node — a document, or data for a template — and its way through the network. The manager keeps the document once, and a *leg* for every node the job reached: what became of it there (routed on, shown, held, printed, archived, failed …). The job's state follows its legs: `processing` while a delivery is pending, `held` while an inbox holds it, `completed` when every node took it, `canceled`, or `aborted` when a delivery failed for good. [States](jobs.md#states) lists them; [What flows](flow.md) tells what a job carries and the metadata that goes with it.

The process holds nothing: the records and the documents are in a store, so the same manager runs as one process on a site's server, as several, or as a function. What is due — deliveries, slips to print on paper, old records to remove — is done right after a job is routed, and by the *tick*, every half minute in a process or when a scheduler calls it ([The tick](inboxes.md#tick)).

## Where to start {#start}

- **You print.** [Jobs](jobs.md). The nodes you may send to are those shared with everyone, your groups' and your own ([Scopes](roles.md#scopes)).
- **You lay out the network.** [Topology](topology.md), then [Node kinds](nodes.md) for each node's settings, [What flows](flow.md) for what may be connected to what, [Templates and data](templates.md) when forms are filled from data, and [Display boards](displays.md) for the tablets.
- **You connect an application.** [Where documents come in](entries.md): endpoints and their keys, apps and their tokens, CUPS server nodes and IPP; [Exports](exports.md) for applications and people that fetch documents; then [The API](api.md).
- **You run the service.** [Running the manager](operate.md).
