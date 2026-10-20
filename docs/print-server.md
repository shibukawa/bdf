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

- **printers**: the URL (`ipp://`, `ipps://`, `http://`, `https://`; a CUPS queue is `ipp://host:631/printers/<name>`), the user and password of a printer that asks for Basic authentication (`password-env` or `password-file`), and the certificate of an `ipps` printer: `fingerprint` pins it; without one the first certificate seen is recorded in the spool and required after, as CUPS does; `insecure: true` takes any. A printer without IPP is `socket://host:9100` or `lpd://host/queue`, with the `format` it takes (`ps`, the default, or `pdf`) and `pjl` (a PJL header before the document; by default for PDF only). It is sent PostScript or PDF with the ticket in the file, tells nothing of its jobs, and is connected to only while it is not known to be up.
- **queues**: the printers in order, the policy, the ticket's `defaults` (IPP's names: `copies`, `sides`, `media`, `media-source`, `media-type`, `print-color-mode`, `print-quality`, `printer-resolution`, `print-scaling`, `orientation-requested`, `finishings`, `output-bin` …), which a job's ticket asks over, and how it is sent.
- **retention**: how long the documents of the jobs done are kept (7 days), and their records, which are the accounting (400 days).
- **limits**: the size of a document (256 MiB; beyond it 413), the jobs not done and their documents (10,000 and 10 GiB; beyond them 507), the pages of a job, how many jobs are converted or drawn at once (0: the number of CPUs), how long a printer is retried, and how many submissions and event streams are taken at once.
- Durations take `s`, `m`, `h` and `d`; sizes `KiB`, `MiB`, `GiB` (or `KB` …). An unknown key is an error, so that a misspelled one is not ignored. Relative paths are the configuration file's directory's. The configuration is read when the server starts.

`bdf printd -config printd.yaml -check` reads the configuration, asks each printer what it supports, and prints a table, without opening the spool (exit status 1 when a printer does not answer).

## Who may do what

- **On the loopback, without tokens** (the default), anyone of the machine may do anything; `X-BDF-User` (which `bdf print` sets to the user of the process) names the user of a job.
- **Beyond the loopback** the server does not start without TLS and tokens (unless `insecure: true` says so).
- **A token** has a name (the user of its jobs), a secret of 16 characters at least (from `file`, `env`, or `value` in the file itself, which is not recommended), the queues it may submit to (all when none) and roles: `submit` (submit jobs and act on its own), `manage` (act on others' jobs and read the accounting), `admin` (pause, resume, reject and accept queues, discover printers, take a backup). It is sent as `Authorization: Bearer …` and compared in constant time.
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

Out of `/api/v1`: `/healthz` (200 while the server serves), `/metrics` (Prometheus' text: the jobs by queue and state, the jobs done, sheets, the printers' state, bytes sent, retries, the size of the spool), and `/`, the page of the state (it only reads; a token can be given on it).

From Go, [`print/server/client`](https://github.com/shibukawa/bdf/tree/main/print/server/client) is the client of the API (`Submit`, `Job`, `Jobs`, `Cancel` …, `Wait`, `Events`, `Accounting`), and depends on nothing else of bdf. A program can also embed the server: [`print/server`](https://github.com/shibukawa/bdf/tree/main/print/server) opens the spool of a configuration and gives its `http.Handler` (see [the sample](examples/print-server.md)).

## Encrypted documents

An encrypted `.bdf` stays encrypted in the spool. Its password, given with the job (`bdf print` sends `$BDF_PASSWORD`, checked first), is kept in memory only: after a restart the job waits (`pending-held`, `job-password-wait`) until it is given again (`bdf jobs password <id>`). A PDF or an Office document that needs a password is converted with it, and the converted `.bdf` is encrypted with the same password. `bdf print` does not send a password over a connection without TLS to another machine unless `-allow-plaintext`.

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
