# Sign-in, roles and audit

## Sign-in {#sign-in}

Every page but the devices' `/tablet` is for a person signed in at the organization's identity provider (OpenID Connect). Anyone the provider signs in gets in, with no registration here, when the operator's rules let them (an issuer, an email domain, a group); the groups the provider says a person is in are copied at each sign-in, and who is an administrator is decided on each request. Sign out ends the session here, and at the next sign-in the provider asks which account to use.

## Roles {#roles}

| Role | Pages | |
|---|---|---|
| everyone signed in | Jobs, Topology, Help | print to the nodes they may send to; see their jobs and their groups' |
| `topology` | Topology as the editor; Admin → Apps, Inboxes, Jobs, Deliveries | lay out the network and what is at it — nodes, templates, keys, apps, devices, inboxes, deliveries — and see every job |
| `audit` | Audit | the audit log, and every job (reading) |
| administrator | all | both roles, and the API as the admin token |

Who has a role is the operator's configuration: a person by their identity or verified email, or the members of a group ([Sign-in](operate.md#oidc)). A page of a role you do not have says so.

## Scopes of nodes {#scopes}

Every node has a scope and an owner, set in its sidebar on the topology:

| Scope | Who may send to it from the pages | Whose its jobs are |
|---|---|---|
| `shared` (the default) | everyone signed in | the sender's alone |
| `group`, with a group as owner | the group's members | the group's: its members see them |
| `personal`, with a person's email as owner | that person | theirs |

A display without a group of its own gives its jobs its fallback printer's group. An app's jobs are the app's group's when the node gives them none. The scope says nothing to applications and IPP clients: an endpoint takes what its key or its apps send, and IPP what the submit tokens allow ([Where documents come in](entries.md)).

## Who sees a job {#who-sees}

A job records who sent it as the manager verified it — the person signed in, the app, the token, the endpoint, or anonymous — apart from the name the sender gave, shown as "(as …)". The sender sees the job, so do the members of its group, and an administrator or an auditor sees every job. The sender and an administrator cancel it; an app cancels its own with the scope `jobs.cancel`. A reprint is a new job of the person who reprinted.

## Audit {#audit}

The Audit page (the role `audit`) lists who did what, newest first: when, who, the action, its subject and a detail; `GET /api/audit` is the same log. What is recorded:

| Action | |
|---|---|
| `login`, `login-refused` | a sign-in, and one the rules refused |
| `put-node`, `delete-node` | a node saved or deleted |
| `rotate-key` | an endpoint's key rotated |
| `put-template`, `delete-template` | a template uploaded or deleted |
| `put-app`, `delete-app`, `rotate-app-secret` | an app registered, changed or deleted, its secret rotated |
| `enrollment-code`, `enroll`, `remove-device` | a code made for a display, a device enrolled, a device removed |
| `refuse-document` | an entry refused a document for its format |
| `token-refused` | an app's credentials refused at `/oauth/token` |
| `oauth-code`, `user-token`, `code-refused` | a code given to a public client for a person, a token of the person made of it, an exchange of a code refused |
| `export-document` | a document downloaded from an export, and by whom |
| `release-item`, `discard-item` | an inbox's item released or discarded, and by whom |
| `cancel-job`, `reprint` | a job canceled, a job reprinted |

Under the log, the recent jobs and the route each took.
