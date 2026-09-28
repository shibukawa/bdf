# secure-reader: books in protected mode

[`examples/secure-reader`](https://github.com/shibukawa/bdf/tree/main/examples/secure-reader) — a site where readers log in and read the PDF books they own. A book is never sent as a file but in [protected mode](../architecture/protection.md#protected-mode): the browser asks for ten pages at a time, and each answer is sealed for a key pair the browser made for that one request. Once a segment is open, neither side keeps the key that opens it.

```sh
cd examples/secure-reader
node web/build.mjs     # builds web/main.js and web/lib/worker.js (no converter)
go run .
open http://127.0.0.1:8084/    # alice / alice-pass owns both books, bob / bob-pass one
```

## What the server does

At startup `main.go` converts each PDF of `-books` (two sample books, 32 and 23 pages, made by `books/gen.mjs`) to bdf once and keeps it under `-cache`, with a cover picture. Nothing under `-cache` is served as a file. It answers:

- `GET /login`, `POST /login`, `POST /logout` — the demo accounts' passwords are kept as PBKDF2 hashes; the session cookie is `HttpOnly` and `SameSite=Strict`, and the server keeps only the hash of it, so a copy of the session table logs nobody in
- `GET /` — the reader's shelf: every book, marked as the reader's own or as a sample of its first ten pages
- `GET /covers/NAME` — a book's cover, to a reader who is logged in
- `POST /segments/NAME` — the book, ten pages at a time: [`segment.Handler`](https://github.com/shibukawa/bdf/tree/main/segment) with the site's own `Open` (is the reader logged in, does the book exist) and `Allow` (the first segment of any book is a sample, the others go to the book's owners; more than 20 segments a minute is answered 429)
- `GET /read/` — the viewer's page and scripts, which hold nothing of any book

Every answer carries a Content-Security-Policy that allows the site's own scripts only, and `http.CrossOriginProtection` refuses forms and segment requests from other sites.

## What the browser does

`web/main.ts` opens the book with [`MiniViewer`](https://github.com/shibukawa/bdf/tree/main/examples/miniviewer) as `{ kind: "segments", url: "/segments/NAME" }`. The rendering worker then does the rest (`SegmentLoader` of `@bdf/core`): for each request it makes a P-256 key pair whose private key cannot be exported, sends the public key with the page it needs and the segments it holds, opens the answer and lets the key go. Pages that come into view fetch their segment first; the segment three pages ahead is fetched in advance. When the server refuses pages (a sample's end, reading too fast), the page says so.

## Why this shape

Sealing the whole book once with one key (as bdf's password encryption does, [spec §3.5](../spec.md#35-暗号化)) would leave a reader holding a ciphertext and its key: whoever takes the key later opens everything that was ever sent. Here every segment has its own key, agreed with ECDH between two keys that are both thrown away, so a recording of the traffic stays sealed even if the server's TLS key, a reader's password or a session cookie is taken later. What the server still decides is who reads which pages, and how fast.

The cost is small: sealing a segment takes well under a millisecond (the outline of the book and the key agreement included, [design.md §3.30](../design.md#330-保護モード-区間ごとに封印して配るsegmentログインした読者向け)), and fonts and pictures that several pages share are sent once, with the first segment that needs them. What it gives up is a shared CDN cache — every answer is sealed for one reader — and offline reading, which would need a key kept on the device.

It does not stop a reader who may read a book from copying what the screen shows, and the server keeps the books in the clear. The format side is [spec §3.6](../spec.md#36-保護モード区間の文書segment).
