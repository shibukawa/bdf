# print-server: printing invoices through a spool

[`examples/print-server`](https://github.com/shibukawa/bdf/tree/main/examples/print-server) — a small business application that prints invoices. It fills an invoice template with what is typed on its page, puts the document in the queue of a [print server](../print-server.md) it embeds, and shows each job as the spool sends it and the printer prints it. The page answers at once; the printing happens behind it, and survives a restart.

```sh
cd examples/print-server
go run . -fake-printer          # prints on bdf's fake IPP printer, in the same process
open http://127.0.0.1:8083/

go run . -printer ipps://printer.example.local/ipp/print   # on a printer of the network
```

## What the server does

- **Fills the template.** `POST /invoices` reads the customer, the items (a line each: name, quantity, price), a note and the user, and fills [`converter/xlsx/testdata/form-rows.xlsx`](https://github.com/shibukawa/bdf/tree/main/converter/xlsx/testdata) with them ([`form.Fill`](../forms.md)): an Excel workbook whose tags (`{{no}}`, `{{customer}}`, `{{#items}}` … `{{/items}}`) take the data, laid out as the pages Excel prints. The result is a `.bdf` in memory.
- **Puts it in the queue.** The application opens a print server from a configuration it writes itself (one printer, one queue `invoices` whose defaults print on both sides), with [`server.New`](https://github.com/shibukawa/bdf/tree/main/print/server), and submits the document with `Spool().Submit`, for the user the page names. It does not wait for the printer.
- **Shares the API.** The print server's `http.Handler` is mounted under `/api/`: the page reads `/api/v1/jobs` and follows `/api/v1/events`, and `bdf jobs -server http://127.0.0.1:8083` works on the same jobs.
- **Keeps the jobs.** The spool is `data/spool.db`. Stop the program while a job is being sent and run it again: the job is checked against the printer and printed once.

## Why this shape

A business system that prints should not hold a request open until the printer is done, nor lose a job because the printer was off. Embedding the print server gives the application the spool's guarantees (one file, retries, no job printed twice, accounting) without another process; a larger site runs [`bdf printd`](../print-server.md) on its own and talks to it with [`print/server/client`](https://github.com/shibukawa/bdf/tree/main/print/server/client) or plain HTTP, as `curl` does.
