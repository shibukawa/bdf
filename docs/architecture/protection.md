# Passwords and protected mode

There are two ways to limit who can read a bdf document. A document encrypted with a password can be fetched by anyone, but only opened by those who know the password. In protected mode, the server keeps the document and hands it to a logged-in reader a few pages at a time, each answer sealed with a key made for that request.

| | Encrypted with a password | Protected mode |
|---|---|---|
| How it is served | As a single file or a split layout; range requests and CDNs work | The server answers segment by segment (ten pages each), every answer made for one reader's one request |
| Key | One key, derived from the password | Agreed for each request between the reader's browser and the server, with keys thrown away afterwards |
| If a key is taken later | Everything sent so far can be opened | A recording of the traffic stays sealed |
| Offline | Readable with the password | Not readable (it would need a key kept on the device) |
| What the server decides | Nothing, once the file is out | Who reads which pages, and how fast |

## Password-protected input and encryption

A password-protected input (a read-password Office document, a user-password PDF) is opened with its password and converted. The bdf output is then encrypted with the same password, sealed part by part (AES-256-GCM), so range requests and split layouts keep working. The viewer decrypts with WebCrypto, and a server never stores the password ([spec §3.5](../spec.md#35-暗号化)).

## Protected mode

The server keeps the document and hands it to a logged-in reader segment by segment (ten pages each by default): each segment is sealed for a key pair the reader's browser made for that one request, and neither side keeps the key afterwards. A recording of the traffic stays sealed even if the server's keys or the reader's session are taken later. The Go package `segment` answers the requests, `SegmentLoader` of `@bdf/core` makes them, and the rendering worker fetches the pages it is asked to draw ([spec §3.6](../spec.md#36-保護モード区間の文書segment), [secure-reader](../examples/secure-reader.md)).

It does not stop a reader who may read a document from copying what the screen shows, and the server keeps the documents in the clear. [secure-reader](../examples/secure-reader.md) is a working site built this way; it also covers what sealing a segment costs and what is given up (a shared CDN cache, offline reading).
