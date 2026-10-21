# Print server

`bdf printd` is a print server: it takes jobs over HTTP, keeps them in one file, sends them to IPP printers (and to printers of 9100 and LPD) in order, follows each until the printer is done, and retries a printer that is not reached. A business system hands it a document and goes on; the server does the waiting. It speaks IPP to the printers the way [`bdf print`](architecture/cli.md#print-printers) does, without a driver of the system.

```sh
bdf printd -config printd.yaml                           # in the foreground; Ctrl-C or SIGTERM stops it
bdf print -server http://127.0.0.1:8640 -queue office invoice.bdf     # submit, print the job's id, return
bdf print -server http://127.0.0.1:8640 -queue office -wait 10m report.xlsx   # and follow it until printed
bdf jobs -server http://127.0.0.1:8640                    # the jobs not done
bdf printers -server http://127.0.0.1:8640                # the queues and the printers
```

`$BDF_SERVER` and `$BDF_TOKEN` stand in for `-server` and the token. A browser at `http://127.0.0.1:8640/` shows the queues, the printers and the jobs as they change.

## What it does

- **Jobs survive a restart.** A job is recorded, with its document, before the submission is answered. The jobs, the documents, the accounting and the fingerprints of the printers' certificates are all in one file, `spool.db` (the [yuruidb](https://github.com/shibukawa/yuruidb) store over bbolt, written with `fsync`). One process opens it; a second `bdf printd` on the same file stops with an error.
- **No job is printed twice.** Before a document is sent, the attempt is recorded; the printer's job id is recorded as soon as the printer answers `Create-Job`, and that the printer has the whole document as soon as it has. A server that stops at any point — a crash, a power cut — checks each job that was at a printer against the printer before it sends anything again: the job is followed if the printer has it, sent again if it does not, and held for a person (`bdf-duplicate-check`) in the one case where the printer has more than one job that may be it. Each attempt's job name at the printer begins with `#<job>.<attempt>`, which is what the server looks for. A printer of 9100 or LPD keeps no record of its jobs: the server records when it begins to send what the printer may print (the first byte of a 9100 job, the control file of an LPD job); a job stopped before that is sent again, one stopped after it is held for a person (`bdf-duplicate-check`), and one the printer had is taken as printed.
- **The printer waits, the job does not fail.** A printer that is off, busy or out of reach is tried again after 10 s, 30 s, 1 min, 2 min, 5 min, then every 10 min, for 24 hours (`retry-for`); the job's reason is `bdf-retry-wait` meanwhile, and `bdf-retry-exhausted` when it is aborted. A stopped printer (out of paper, a jam) keeps the job it has; the next waits. A job the printer aborts is not sent again by itself: `bdf jobs retry` does.
- **Two jobs at a printer.** A printer is given up to two jobs of a queue at once (`max-inflight`), so that the next is there when the first is printed; the rest wait in the queue, where they can still be reordered, held, moved to another queue, or canceled without the printer.
- **Queues of several printers.** A queue's printers are tried in order: `failover` (the default) sends to the first printer that is up, `any` to the first that has room.
- **Other formats are converted.** A PDF, an Office document, an image or any other format of `bdf generate` is converted to `.bdf` in the spool (once: a job sent again is not converted again), as the pages it prints: an Excel workbook as Excel prints it, HTML and Markdown as A4 pages. HTML, Markdown and EPUB do not fetch images from the network.
- **States are IPP's.** A job is `pending`, `pending-held`, `processing`, `processing-stopped`, `canceled`, `aborted` or `completed`, with IPP's `job-state-reasons` (`job-queued`, `job-printing`, `job-password-wait` …) and a few of bdf's own (`bdf-retry-wait`, `bdf-retry-exhausted`, `bdf-unsupported`, `bdf-duplicate-check`).
- **Order**: the highest priority first (`job-priority`, 1 to 100, 50 by default), then the oldest. `job-hold-until` holds a job (`indefinite`, `day-time`, `evening`, `night`, `weekend`, `second-shift`, `third-shift`) or until a time.

## Configuration

```yaml
listen: 127.0.0.1:8640              # the default
spool: /var/lib/bdf/spool.db        # default: spool.db next to this file
user: bdf                           # the requesting-user-name the printers see
tls: { cert: /etc/bdf/tls.crt, key: /etc/bdf/tls.key }   # required beyond the loopback
tokens:
  - { name: erp, file: /etc/bdf/erp.token, queues: [office, warehouse], roles: [submit], act-as: true }
  - { name: ops, env: BDF_OPS_TOKEN, roles: [submit, manage, admin] }
printers:
  mfp1: { url: "ipps://mfp1.example.local/ipp/print", fingerprint: "sha256:…" }
  mfp2: { url: "ipp://192.0.2.20/ipp/print", username: print, password-env: MFP2_PASSWORD }
  old-laser: { url: "socket://192.0.2.30:9100" }                     # no IPP: PostScript over the raw port
  plotter: { url: "lpd://192.0.2.31/lp1", format: pdf, pjl: true }   # PDF over LPD, with a PJL header
  labels: { url: "socket://192.0.2.40:9100", profile: zd421-203 }    # ZPL labels
  register: { url: "serial:///dev/ttyUSB0?baud=19200", format: escpos }   # ESC/POS receipts
  zebra-usb: { url: "ipp://localhost:631/printers/zebra", profile: zd421-300 }   # a raw queue of CUPS
queues:
  office:
    printers: [mfp1, mfp2]
    policy: failover                # or any
    defaults: { sides: two-sided-long-edge, print-color-mode: monochrome }
    format: auto                    # auto (PDF when the printer takes it, else PostScript, PWG Raster, Apple Raster, PCLm), pdf, ps, pwg-raster, urf or pclm
    raster-if-lossy: true
    max-inflight: 2
  warehouse: { printers: [mfp2] }
retention: { documents: 7d, records: 400d }
limits: { max-document: 256MiB, max-jobs: 10000, max-bytes: 10GiB, max-impressions: 10000,
          renders: 0, retry-for: 24h, uploads: 8, events: 64 }
fonts: { dirs: [/usr/share/fonts/bdf], no-system: false }
temp-dir: /var/tmp/bdf
```

- **printers**: the URL (`ipp://`, `ipps://`, `http://`, `https://`; a CUPS queue is `ipp://host:631/printers/<name>`), the user and password of a printer that asks for Basic authentication (`password-env` or `password-file`), and the certificate of an `ipps` printer: `fingerprint` pins it; without one the first certificate seen is recorded in the spool and required after, as CUPS does; `insecure: true` takes any. A printer without IPP is `socket://host:9100` or `lpd://host/queue`, with the `format` it takes (`ps`, the default, or `pdf`) and `pjl` (a PJL header before the document; by default for PDF only). It is sent PostScript or PDF with the ticket in the file, tells nothing of its jobs, and is connected to only while it is not known to be up. A printer of labels, receipts or forms ([below](#labels-receipts-and-forms)) has its language as `format` (`zpl`, `sbpl`, `tpcl`, `escpos`, `escp`) or its model as `profile`, over `socket://`, `lpd://`, `serial://`, `usb://` or `winspool://` — or over IPP to a raw queue of CUPS, which passes the data on.
- **queues**: the printers in order, the policy, the ticket's `defaults` (IPP's names: `copies`, `sides`, `media`, `media-source`, `media-type`, `print-color-mode`, `print-quality`, `printer-resolution`, `print-scaling`, `orientation-requested`, `finishings`, `output-bin` …, and what a job says about itself: `job-account-id`, `job-accounting-user-id`, `document-metadata`, and `bdf-attributes`, attributes of one's own named `x-…`), which a job's ticket asks over, and how it is sent.
- **retention**: how long the documents of the jobs done are kept (7 days), and their records, which are the accounting (400 days).
- **limits**: the size of a document (256 MiB; beyond it 413), the jobs not done and their documents (10,000 and 10 GiB; beyond them 507), the pages of a job, how many jobs are converted or drawn at once (0: the number of CPUs), how long a printer is retried, and how many submissions and event streams are taken at once.
- Durations take `s`, `m`, `h` and `d`; sizes `KiB`, `MiB`, `GiB` (or `KB` …). An unknown key is an error, so that a misspelled one is not ignored. Relative paths are the configuration file's directory's. The configuration is read when the server starts.

`bdf printd -config printd.yaml -check` reads the configuration, asks each printer what it supports, and prints a table, without opening the spool (exit status 1 when a printer does not answer).

`ipp:` shows the queues as IPP printers ([below](#as-an-ipp-printer)); `upstreams:` and printers of `proxy: true` fetch jobs across a NAT ([below](#branches-behind-a-nat)).

## Who may do what

- **On the loopback, without tokens** (the default), anyone of the machine may do anything; `X-BDF-User` (which `bdf print` sets to the user of the process) names the user of a job.
- **Beyond the loopback** the server does not start without TLS and tokens (unless `insecure: true` says so).
- **A token** has a name (the user of its jobs), a secret of 16 characters at least (from `file`, `env`, or `value` in the file itself, which is not recommended), the queues it may submit to (all when none) and roles: `submit` (submit jobs and act on its own), `manage` (act on others' jobs and read the accounting), `admin` (pause, resume, reject and accept queues, discover printers, take a backup), `proxy` (fetch the jobs of its queues as the server of a branch). It is sent as `Authorization: Bearer …` and compared in constant time.
- **`act-as`** lets a token name the user of a job with `X-BDF-User`: a business system that prints for its users. Without `manage`, a token sees and acts on the jobs of the user it acts as.
- The log has no document, password or token, and no user or password of a printer's URL.

## The API

JSON, under `/api/v1`. An error is `{"error": {"code": "…", "message": "…"}}` with an HTTP status: 400 `invalid` or `bad-document`, 401 `unauthorized`, 403 `forbidden` or `wrong-password`, 404 `not-found` or `no-queue`, 409 `not-possible` (not in the job's state), 413 `too-large`, 415 `unknown-format`, 503 `not-accepting` or `busy`, 507 `full`.

| Method and path | |
|---|---|
| `POST /queues/{q}/jobs` | Submit: `multipart/form-data` with `document` (the file) and, in any order, `ticket` and `options` (JSON), `name`, `password`, `priority`, `hold-until`, `hold-until-time`. 201 and the job |
| `GET /jobs?queue=&user=&state=&after=&limit=` | The jobs not done; `state=all` with those done, or a state. `after` and `next` page by id |
| `GET /jobs/{id}` | A job |
| `POST /jobs/{id}/cancel`, `hold`, `release`, `retry`, `password`, `move` | Cancel (at the printer too), hold, release (`{"password"}` too), send again, give the password, move (`{"queue"}`) |
| `PATCH /jobs/{id}` | `{"priority": n}` |
| `GET /queues`, `GET /queues/{q}` | Queues and how many jobs are in each state |
| `POST /queues/{q}/pause`, `resume`, `reject`, `accept` | Stop sending, send again, stop taking jobs, take them again (admin) |
| `GET /printers`, `GET /printers/{p}` | Printers, their state, and what one supports |
| `GET /discover?timeout=` | IPP printers on the network, by DNS-SD (admin) |
| `GET /accounting?from=&to=&format=jsonl\|csv` | The jobs done: who printed what, where, and how many sheets (manage) |
| `GET /backup` | A consistent copy of `spool.db` (admin) |
| `GET /events?job=` | Server-Sent Events: `job`, `queue` and `printer`, with the JSON of each |

```sh
curl -F document=@invoice.bdf -F 'ticket={"sides":"two-sided-long-edge","copies":2}' \
     -H "Authorization: Bearer $TOKEN" https://print.example.local:8640/api/v1/queues/office/jobs
```

A ticket also carries what a job says about itself, which the server keeps with the job, tells back over IPP, and sends to the printer: `job-account-id` and `job-accounting-user-id` (the account the job is charged to and its user, asked of a printer that says it takes them), `document-metadata` (`key=value` strings about the document, with the keys of Dublin Core, `title`, `creator`, `identifier`, `subject`, `description`, `date` …, or keys of one's own that begin with `x-`; a printer keeps them with the job), and `bdf-attributes`, job attributes of one's own named `x-…` (`{"x-cost-center": "cc-42", "x-tags": ["a", "b"]}`), sent as they are; a printer that does not know one says so and prints without it (with `strict`, the job fails). `bdf print` has `-account`, `-accounting-user`, `-metadata key=value` and `-attribute x-name=value` for them.

Out of `/api/v1`: `/healthz` (200 while the server serves), `/metrics` (Prometheus' text: the jobs by queue and state, the jobs done, sheets, the printers' state, bytes sent, retries, the size of the spool), and `/`, the page of the state (it only reads; a token can be given on it).

From Go, [`print/server/client`](https://github.com/shibukawa/bdf/tree/main/print/server/client) is the client of the API (`Submit`, `Job`, `Jobs`, `Cancel` …, `Wait`, `Events`, `Accounting`), and depends on nothing else of bdf. A program can also embed the server: [`print/server`](https://github.com/shibukawa/bdf/tree/main/print/server) opens the spool of a configuration and gives its `http.Handler` (see [the sample](examples/print-server.md)).

## Encrypted documents

An encrypted `.bdf` stays encrypted in the spool. Its password, given with the job (`bdf print` sends `$BDF_PASSWORD`, checked first), is kept in memory only: after a restart the job waits (`pending-held`, `job-password-wait`) until it is given again (`bdf jobs password <id>`). A PDF or an Office document that needs a password is converted with it, and the converted `.bdf` is encrypted with the same password. `bdf print` does not send a password over a connection without TLS to another machine unless `-allow-plaintext`.

## As an IPP printer

With `ipp:` in the configuration, each queue is also an IPP Everywhere printer that every system adds without a driver: Windows (its IPP class driver, as Windows Ready Print adds printers), macOS and iOS (AirPrint), Android (Mopria, the default print service), ChromeOS and Linux (CUPS). They send PDF (listed first, so vectors are sent), Apple Raster or PWG Raster, JPEG or PNG; the server converts what it takes into `.bdf`, as it does a submission, and prints it on the queue's printers. Raster pages become pictures of the paper they were drawn for.

```yaml
ipp:
  listen: ":8631"              # the default; ipp and ipps on the one port
  host: print.example.com      # the name in the printers' URIs and the certificate: the machine's, in .local, by default
  auth: none                   # none (as the printers of a LAN) or basic (a token's name and secret): basic when there are tokens
  require-tls: false           # refuse ipp (HTTP 426): only ipps; so with basic
  tls: { cert: /etc/bdf/ipp.crt, key: /etc/bdf/ipp.key }   # without it, one the server signs itself (ipp-tls.pem next to the spool)
  queues: [office, warehouse]  # the queues shown: all by default
  default: office              # also at /ipp/print
  location: "2F"               # printer-location, under a queue's own
  organization: Example Inc.
  dns-sd: true                 # advertise on the local links (mDNS): the default
queues:
  office: { printers: [mfp1], info: "Office printer", location: "2F north" }
```

- **The printers**: `ipp://host:8631/ipp/print/<queue>` and `ipps://…` (the same port: the first byte of a connection tells TLS). The operations are those IPP Everywhere requires (Print-Job, Validate-Job, Create-Job, Send-Document, Close-Job, Cancel-Job, Cancel-My-Jobs, Get-Job-Attributes, Get-Jobs, Get-Printer-Attributes, Identify-Printer) and Hold-Job and Release-Job. A job's id is the spool's (`bdf jobs` shows the same numbers), its state the spool's. CUPS's `ipptool` checks of IPP/1.1, IPP/2.0 and IPP Everywhere's attributes pass (`test/print/ipp-everywhere-server.test`).
- **What a job asks** goes into its ticket: `copies`, `sides`, `media` or `media-col` (a size is matched to the paper of the queue), `print-color-mode`, `print-quality`, `printer-resolution`, `page-ranges`, `print-scaling`, `orientation-requested`, `finishings`, `output-bin`, `job-hold-until`, `job-priority`, and what the job says about itself: `job-account-id`, `job-accounting-user-id`, attributes of one's own named `x-…`, and `document-metadata` with the document (Print-Job, Send-Document), which Get-Job-Attributes tells back and the printer behind is sent. What the queue does not support is told back; with `ipp-attribute-fidelity` the job is refused. `document-password` opens an encrypted PDF, Office document or `.bdf`, and is kept in memory only, as a submission's.
- **What the printers say they support** is the queue's first printer's (its paper, trays, finishings, colour, both sides), with what the server does itself (copies, page ranges, scaling) whatever the printer. A printer that has not answered yet, or of 9100 or LPD, is shown with an office's paper.
- **Who**: with `auth: none` anyone of the network prints, under the `requesting-user-name` the client gives, and cancels the jobs of that name; others' jobs are listed without their names. With `auth: basic` the client gives a token's name and its secret (the system asks once and keeps them); ipp from another machine is refused, so that the secret goes over TLS, and only `_ipps._tcp` is advertised. A token's queues and roles hold as in the API.
- **DNS-SD**: each queue is advertised as `_ipp._tcp` and `_ipps._tcp`, with the subtypes `_print` (IPP Everywhere) and `_universal` (AirPrint), on a host name of the server's own (`<machine>-bdf.local`; the machine's name is the system responder's), with the TXT record clients choose by (`pdl`, `URF`, `Color`, `Duplex`, `PaperMax`, `TLS`, `air` …), updated once the printers behind have answered. The server shares the port of mDNS with Bonjour, Avahi and Windows.
- **Other subnets**: mDNS stays on its link. `bdf printd -config printd.yaml -dns-sd example.com` writes the same records for the company's DNS zone (wide-area DNS-SD: `ipp.host` is then the server's name in the zone); and `bdf printd -config printd.yaml -airprint-profile print.example.com > printers.mobileconfig` writes a configuration profile for iPhones and iPads (an AirPrint entry of each queue, and the certificate the server signed itself as a root to trust), for an MDM to send. The page of the printers (`https://host:8631/`) links the same profile, made for the address it was opened at (over https only: over plain HTTP another machine could put its own certificate in it). The certificate the server makes is a leaf, signed by itself, that signs nothing else: a device that trusts it trusts this server only.

Adding the printers:

| System | |
|---|---|
| Windows 10 and 11 | Settings → Bluetooth & devices → Printers & scanners → Add device: the queue is listed (mDNS). Elsewhere: "Add a printer using an IP address or hostname" → IPP device → `https://host:8631/ipp/print/office` (or `http://`). The IPP class driver is used; no driver is installed |
| macOS | System Settings → Printers & Scanners → Add Printer: the queue is listed (AirPrint). Elsewhere: the IP tab, protocol IPP, address `host:8631`, queue `ipp/print/office`. From the shell: `lpadmin -p office -E -v ipps://host:8631/ipp/print/office -m everywhere` |
| iPhone, iPad | The queue is listed when printing (AirPrint) on the same subnet; elsewhere the profile above |
| Android | The default print service (Mopria) lists it; elsewhere add it by its address in the service's settings |
| ChromeOS | Settings → Printers → Add printer: listed, or by address with IPP or IPPS and the path `ipp/print/office` |
| Linux | `lpadmin -p office -E -v ipps://host:8631/ipp/print/office -m everywhere`, or the desktop's printer settings, which list it |

## Branches behind a NAT

A branch whose printers the head office's server does not reach (a NAT, no VPN) runs `bdf printd` of its own, which fetches the jobs of a queue of the head office and prints them on its printers: it connects out, the head office never connects in. The two speak IPP's Shared Infrastructure Extensions (PWG 5100.18), so another proxy of that standard can fetch the jobs too.

```yaml
# the head office: a queue whose printer is the branch's proxy
printers:
  osaka: { proxy: true }
queues:
  osaka: { printers: [osaka] }
ipp: { listen: ":8631", host: print.example.com }
tokens:
  - { name: osaka-proxy, file: /etc/bdf/osaka.token, roles: [proxy], queues: [osaka] }
```

```yaml
# the branch
printers:
  mfp: { url: "ipp://192.0.2.10/ipp/print" }
queues:
  local: { printers: [mfp] }
upstreams:
  hq:
    url: ipps://print.example.com:8631/ipp/print/osaka
    queue: local
    username: osaka-proxy
    password-file: /etc/bdf/osaka.token
    poll: 10s              # the default
    # fingerprint: "sha256:…"   pins the head office's certificate; without it the first one seen is kept
```

- Jobs submitted to the head office's queue (by the API, `bdf print -server`, or IPP) wait there as `job-fetchable`. The branch asks every `poll`, takes each job (`Fetch-Job`, `Acknowledge-Job`), fetches its document (the `.bdf` as the head office keeps it, converted already; with its password, which the branch keeps in memory only), puts it in its queue, and tells the head office it has it (`Acknowledge-Document`). From then on it reports the job's state (`Update-Job-Status`): `bdf jobs` at the head office shows the branch's printing, completion and impressions.
- **Not printed twice.** A job is bound to the branch that took it and offered to no other. The branch keeps the head office's job's uuid with its own job, and a job it has already (it stopped between taking a job and saying so) is acknowledged, not taken again. A job the head office gave a branch that no longer has it (the branch lost its spool) is aborted at the head office for a person to retry, not sent again by itself.
- **Canceled at the head office** is canceled at the branch (at its next report, or `Update-Active-Jobs`, which the two compare their jobs with each `poll`).
- The head office shows the branch's printer as the queue's (its state, paper and colour, as the branch reports them; down when the branch has not reported for two minutes), and so does the head office's IPP printer of that queue.
- The branch's token has the role `proxy` (it reads the documents of the queue): only for the queues of `queues`, and only with Basic authentication over ipps. A printer of proxies is the printer of one queue only. A job the branch cannot take (a document it cannot read, too large) is reported aborted; one retried at the head office is taken again.

## Labels, receipts and forms

A queue of a printer of labels, receipts or continuous forms prints what any other queue takes — a filled form, a PDF, an image — in the printer's own language ([`bdf print`](architecture/cli.md#print-printers) has the details):

| Language | Printers | What is sent |
|---|---|---|
| ZPL | Zebra, and the printers that take ZPL II | each label one format: `^GFA` of the dots in ZPL's ASCII compression, the QR codes and barcodes as `^BQ`, `^BC`, `^BE` …, `^PQ` copies |
| SBPL | SATO | each label one job `<A>` … `<Z>` between STX and ETX: `<G>H` graphics, `<2D30>`, `<B>`, `<BG>` symbols, `<Q>` copies |
| TPCL | Toshiba TEC | `{D` size, `{SG;` graphic in nibble mode, `{XB`/`{RB` symbols (not Code 128), `{XS;I` issue |
| ESC/POS | Epson TM, receipt printers | bands of `GS v 0` as they are drawn, `GS V` cut, `ESC p` drawer |
| ESC/P | Epson VP, dot-matrix printers | 24-pin bit images (`ESC *`, 180 dpi), `ESC C` page length, `ESC N` perforation skip, form feed |

The queue's `defaults` take the settings of such printers with the ticket's names: `print-darkness`, `print-speed` and `media-tracking` (IPP's), and `bdf-label-gap`, `bdf-cut`, `bdf-cash-drawer`, `bdf-perforation-skip`, `bdf-no-printer-barcodes`, and `bdf-print-mode` (`raster` or `text`: the characters of a form of a template of text in the printer's fonts; ESC/P prints such a form in the text mode by default). Copies are the printer's (`^PQ`, `<Q>`, `{XS`), or sent again for receipts and forms.

The commands are those of the makers' programming references (Zebra's ZPL II guide, SATO's CL4NX Plus and CT4-LX references, Toshiba TEC's B-EX4, B-FV4 and BX400 interface specifications, Epson's ESC/POS and ESC/P references). They have been checked against readers of our own that print them back to dots, not on printers: try a label before a batch, and if a printer prints a symbol wrongly, `bdf-no-printer-barcodes: true` draws it as dots instead (readable at 203 dpi). SBPL is wrapped in STX and ETX, as SATO's status protocols ask. Windows reaches USB printers through the queue of their driver (`winspool:///Name`, passed on as `RAW`), macOS through CUPS (`ipp://localhost:631/printers/<queue>`), Linux through either or `usb:///dev/usb/lp0`.

## Running it as a service

`bdf printd` runs in the foreground, logs to the standard error (JSON; `-log text` for people), and stops on SIGINT or SIGTERM (Ctrl-C or Ctrl-Break on Windows): it stops taking requests, lets the documents being sent finish for 30 seconds, and closes the spool. What is cut is put right when it starts again. It does not install itself as a service: the system's service manager keeps it running.

**Linux (systemd)**: `/etc/systemd/system/bdf-printd.service`, then `systemctl enable --now bdf-printd`.

```ini
[Unit]
Description=bdf print server
After=network-online.target
Wants=network-online.target

[Service]
# a user of its own (useradd --system bdf), who may read the token files
User=bdf
ExecStart=/usr/local/bin/bdf printd -config /etc/bdf/printd.yaml
Restart=on-failure
# spool: /var/lib/bdf/spool.db in printd.yaml
StateDirectory=bdf
TimeoutStopSec=45

[Install]
WantedBy=multi-user.target
```

**macOS (launchd)**: `/Library/LaunchDaemons/io.github.shibukawa.bdf-printd.plist`, then `sudo launchctl bootstrap system /Library/LaunchDaemons/io.github.shibukawa.bdf-printd.plist`.

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key><string>io.github.shibukawa.bdf-printd</string>
  <key>ProgramArguments</key>
  <array>
    <string>/usr/local/bin/bdf</string><string>printd</string>
    <string>-config</string><string>/usr/local/etc/bdf/printd.yaml</string>
  </array>
  <key>RunAtLoad</key><true/>
  <key>KeepAlive</key><dict><key>SuccessfulExit</key><false/></dict>
  <key>StandardErrorPath</key><string>/usr/local/var/log/bdf-printd.log</string>
</dict>
</plist>
```

**Windows**: a service wrapper runs it, such as [NSSM](https://nssm.cc/) or [WinSW](https://github.com/winsw/winsw).

```bat
nssm install bdf-printd "C:\Program Files\bdf\bdf.exe" printd -config "C:\ProgramData\bdf\printd.yaml"
nssm set bdf-printd AppStopMethodConsole 30000
nssm set bdf-printd AppStderr "C:\ProgramData\bdf\printd.log"
nssm start bdf-printd
```

```xml
<!-- bdf-printd.xml, next to WinSW renamed bdf-printd.exe: bdf-printd.exe install, then start -->
<service>
  <id>bdf-printd</id>
  <name>bdf print server</name>
  <executable>C:\Program Files\bdf\bdf.exe</executable>
  <arguments>printd -config "C:\ProgramData\bdf\printd.yaml"</arguments>
  <stoptimeout>30 sec</stoptimeout>
  <onfailure action="restart" delay="10 sec"/>
  <log mode="roll"/>
</service>
```

Without a wrapper, the Task Scheduler can start it at boot: `schtasks /create /tn bdf-printd /sc onstart /ru SYSTEM /tr "\"C:\Program Files\bdf\bdf.exe\" printd -config C:\ProgramData\bdf\printd.yaml"` (it is then not restarted when it fails).

## Keeping the spool

- **Backup**: `GET /api/v1/backup` (admin) gives a consistent copy of `spool.db` while the server runs; the copy opens as a spool. It is made to a temporary file first (`temp-dir`), so a slow download does not hold the spool.
- **Size**: the file does not shrink when documents are removed after their retention; the space is used again. `bdf printd -config printd.yaml -compact` rewrites it without that space, while the server is stopped.
- **Not shared**: one `spool.db` is one server's. It must not be on a network file system that several servers open.
