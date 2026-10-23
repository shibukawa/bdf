# Jobs

The Jobs page is everyone's: a form that prints a document, and the jobs you may see.

## Print a document {#print}

Choose a node, a file and, if you like, a name for the job (the file's name otherwise), and press Print. The nodes listed are those you may send to: the shared ones, your groups' and your own ([Scopes](roles.md#scopes)). The job is received at that node and flows from it along the edges the administrator drew; where it went and what became of it is in the list below the form, and on the [Topology](topology.md#view) page.

What a node takes depends on its kind ([Kinds of payload](flow.md#kinds)). A printer, a display, a queue or an inbox takes a document — PDF, Office documents (Word, Excel, PowerPoint), images and `.bdf`, among the formats the converter reads — and it is made into a `.bdf` document as it is received. A template node takes data (JSON, YAML, CSV, an Excel workbook) and fills a form with it; a convert node takes a file as it is; an endpoint takes what its settings say. A file of a format the node does not take is refused, and the page says so.

## Your jobs {#list}

The list shows the jobs you sent and the jobs sent to your groups' nodes ([Who sees a job](roles.md#who-sees)); an administrator or an auditor sees every job. Refresh reloads it.

| Column | |
|---|---|
| No. | the job's number, `#12`: the same on every page, in the audit log and in the archives |
| Name | the job's name; hovering shows what was received (its kind and format) and the job's metadata |
| Sent to | the node the sender sent it to |
| Sent by | who sent it, as the manager verified: the person signed in, an app, a token, an endpoint, or anonymous. When the sender named someone else (`?user=`, IPP's requesting-user-name), the name is shown after it as "(as …)": claimed, never checked |
| Group | the group whose members see the job: the node's owner group when the job was sent |
| Received | when |
| State | see [States](#states) |
| Route | the legs: each node the job reached and what became of it there, as `node:state`, with the member a class chose or the rule a filter matched in parentheses, and an error after `!` |

**Document** downloads the `.bdf` document of the job: the one received, or the last one a template or a convert node made ([Payloads](flow.md#payloads)).

**Reprint** sends the document again, as a new job of yours, to the node it went to; hold Shift while clicking to send it to the node chosen in the form above instead. The new job carries the sender's metadata of the old one, and you must be allowed to send to the node.

**Cancel** is the sender's and an administrator's; the group's members only see the job. A job is canceled while it has not ended, or while its slip is still open on a display board. Canceling takes its slips off the boards, discards its items in the inboxes, and withdraws its deliveries not yet sent; what a printer already printed stays printed.

## States {#states}

A job's state follows its legs:

| State | |
|---|---|
| `processing` | a delivery is pending: a printer, an archive, a template or a convert step has it and has not finished |
| `held` | an inbox holds it, and nothing else is pending |
| `completed` | every node took it: shown, printed, archived, kept, released or dropped |
| `canceled` | the sender or an administrator canceled it |
| `aborted` | a delivery failed for good: after every retry, or at once when no retry could mend it ([Deliveries](inboxes.md#deliveries)) |

A leg's state is what the node did with the job:

| Leg | Nodes | |
|---|---|---|
| `routed` | queue, class, filter, endpoint, CUPS server | sent on along the edges |
| `shown` | display | a slip on the devices; `canceled` when the job was |
| `held`, `released`, `discarded` | inbox | held; sent on by a person; discarded by a person, or by a cancellation |
| `pending`, `done`, `failed`, `canceled` | printer, archive, template, convert | a delivery: waiting or being sent; done; given up; withdrawn by a cancellation |
| `kept` | sink | among the latest the sink keeps |
| `listed` | HTTP (pull) | listed for its owner, who the leg names ([HTTP (pull)](exports.md)) |
| `dropped` | filter | no rule fit, and there is no otherwise |
| `error` | any | the node could not take the job, and the error says why: a kind of payload it does not take, a node that is gone |

A job that ended is kept 14 days by default (the operator's `job_keep`), with its document while something still shows it; then it is removed.
