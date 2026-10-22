# Where documents come in

A job enters the network at an entry: the form on the Jobs page, any node's IPP printer, an endpoint an application POSTs to, a CUPS server node computers print to, or the API. Who may send where is [Scopes of nodes](roles.md#scopes).

## From the pages {#pages}

The form on [Jobs](jobs.md#print), and the form of a node on the [Topology](topology.md#view) page, send a document to a node you may send to, as yourself: the job is yours, and its group the node's.

## IPP {#ipp}

Every node is an IPP printer at `ipp://<host>:<port>/ipp/print/<name>` (`ipps://` when the manager is behind TLS): a job printed to a queue's name is forwarded, to a printer's name printed, to a display's name shown. A computer adds the node as an IPP printer; `bdf printers -url ipp://<host>:<port>/ipp/print` lists every node, as a CUPS client's `lpstat` does (CUPS-Get-Printers). The job attributes of the request (`media`, `copies`, `job-name` …) and its `document-metadata` go with the job as metadata, `ipp.<name>`; the account (`job-account-id`, `job-accounting-user-id`), `document-metadata` and attributes of one's own (`x-…`) go on to the printers the job reaches ([What goes on to the printer](flow.md#ipp-attributes)).

Who may print over IPP is the operator's `submit_tokens`: with any set, the client gives a token's name and secret as Basic authentication (an app gives its id and secret the same way); with none, anyone who reaches the manager may send. The requesting-user-name the client gives is shown as "(as …)", never checked.

## Endpoints {#endpoints}

An endpoint is a REST entry for an application: a business system, a WMS. It is a node on the topology ([Endpoint](nodes.md#endpoint)); when it is saved for the first time its key is shown once, and Rotate key in its sidebar makes a new one and shows it once (the old one stops working). The manager keeps only the key's hash.

The application sends with the key as a bearer token:

```sh
curl -X POST -H "Authorization: Bearer bdfk_…" \
     --data-binary @statement.pdf \
     "https://<host>/in/erp?name=statement&meta.dept=sales"
```

- The body is the document, or the data for an endpoint with `accept: data`; the name (`?name=`) and the `Content-Type` tell its format when the content does not.
- `?name=` names the job; `?user=` names who sent it, shown as "(as …)".
- `X-Meta-<Key>` headers and `?meta.<key>=` parameters are the job's metadata ([Metadata](flow.md#metadata)).
- The answer is the job (201), with its number and its legs so far. A document of a format the endpoint does not take is refused (415), and the refusal is in the audit log; a wrong key is refused (401).

With a `port`, a process listens on that port too, and every POST to it is to this endpoint: for an application that cannot set a path. A serverless deployment has no ports.

## CUPS server nodes {#cups}

A CUPS server node ([CUPS server](nodes.md#cups)) is for computers and business systems that print: a CUPS server on a port of the process's own, with one printer, the node itself. A computer adds the printer with `ipp://<host>:<port>/printers/<name>` (`ipps://` when the node has a certificate), as it adds any CUPS printer, and what it prints there goes on to the node's `to`: a printer makes the node a relay, a sink shows what was printed, a filter keeps only what fits. `lpstat -h <host>:<port> -p` lists it. The node's sidebar shows the URL and the command.

## Apps {#apps}

An app is an application that sends with credentials of its own rather than an endpoint's key: the credentials are the app's, and the audit log and the jobs say which app sent. An administrator registers one on Admin → Apps: an id, a title, the group whose members see its jobs (optional), and its scopes:

| Scope | |
|---|---|
| `print` | send documents |
| `print.template` | fill templates: send data |
| `jobs.cancel` | cancel its own jobs |
| `tick` | call `POST /api/tick`: a scheduler |
| `admin` | the API as an administrator: a script |

The secret is shown once, when the app is added, and Rotate secret makes a new one (the old one stops working). Disable refuses the app's credentials and the tokens out; Delete removes it.

The app takes a token at `POST /oauth/token`, the OAuth 2.0 client credentials grant, and sends with it as a bearer to the endpoints that list the app among their apps (Apps allowed, in the endpoint's sidebar), or prints over IPP with its id and secret as Basic authentication:

```sh
curl -u 'erp:<secret>' -d grant_type=client_credentials https://<host>/oauth/token
# {"access_token": "…", "token_type": "Bearer", "expires_in": 600, "scope": "print"}
curl -X POST -H "Authorization: Bearer <access_token>" --data-binary @statement.pdf "https://<host>/in/erp?name=statement"
```

A token lasts ten minutes. An app sees its own jobs (`GET /api/my/jobs`) and cancels them with the scope `jobs.cancel`; its jobs are its group's when the node gives them none.

## Scripts {#scripts}

A script sends with the operator's admin token, or a submit token's `name:secret`, as a bearer: `POST /api/destinations/<node>/jobs` with the document as the body takes `?name=`, `?user=` and the metadata as an endpoint does ([Jobs](api.md#jobs)).
