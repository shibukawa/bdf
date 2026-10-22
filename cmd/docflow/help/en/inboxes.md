# Inboxes and deliveries

Two pages of the administrator's are about jobs on their way: Inboxes, where a person sends jobs on by hand, and Deliveries, where what the manager sends on its own is watched and retried.

## Inboxes {#inboxes}

An inbox is a node that holds every job that reaches it until a person releases it or discards it: a filter worked by hand — a check before printing, a batch sent on at the end of the day. Admin → Inboxes shows every inbox with the items it holds, oldest first:

- **Release** sends the item on to nodes: those of the inbox's `to`, each checked by default, or, when the inbox names none, the names typed in the box, comma-separated. The job goes on from the nodes chosen as if the inbox had forwarded it.
- **Discard** ends the item: the job is not sent on from here.
- The name in **Worker** is written on the record and in the audit log, as who released or discarded.

Released and discarded items are shown under the list for the inbox's keep (14 days by default), with where each went and who sent it. Edit the node opens the inbox on the topology. A job canceled on the Jobs pages leaves its inboxes.

## Deliveries {#deliveries}

A delivery is what the manager sends on its own: a job to a printer or to an archive, and the step of a template or a convert node. It is made when the job is routed, and done right after by the process, or by the tick. Admin → Deliveries lists the deliveries waiting or failed, with the job, the node, the printer (of a class, the member tried), the state, how many tries, when the next is due, and the last error.

- A delivery that fails — a printer out of reach, a store that does not answer — is tried again after a wait that grows: 30 seconds, then a minute, two, four, eight, then every 15 minutes, up to 20 tries, about four hours; a class tries its next member. After the last try the delivery fails for good, and the job is `aborted`.
- A failure no retry mends — a template missing, data that cannot be read, a format with no converter — fails for good at once.
- **Retry now** sends a delivery again now: one waiting for its time, or one that failed for good.

## The tick {#tick}

The tick does what is due: it sends the deliveries that wait or are to be tried again, prints on paper the slips nobody opened in time (and every unopened slip of a display none of whose devices is connected), removes the slips and the items ended longer ago than their display or inbox keeps them, the jobs that ended and nothing shows any more, and the enrollment codes expired. A process ticks every 30 seconds (the operator's `tick_interval`); on a serverless platform a scheduler calls `POST /api/tick` every minute ([Serverless](operate.md#serverless)). **Tick now** on Admin → Deliveries runs it and reports what it did.
