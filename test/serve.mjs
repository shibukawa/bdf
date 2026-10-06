// Minimal static server with Range support, used by tests and the demo.
import { createServer } from "node:http";
import { realpath, stat, open } from "node:fs/promises";
import { join, extname, normalize, resolve, sep } from "node:path";
import { pipeline } from "node:stream/promises";

const TYPES = {
  ".html": "text/html; charset=utf-8", ".js": "text/javascript", ".mjs": "text/javascript", ".map": "application/json",
  ".json": "application/json", ".png": "image/png", ".svg": "image/svg+xml", ".bdf": "application/octet-stream", ".css": "text/css", ".ts": "text/plain",
  ".wasm": "application/wasm", ".ttf": "font/ttf", ".otf": "font/otf", ".ttc": "font/collection", ".otc": "font/collection",
};

export async function serve(root, port = 0) {
  root = await realpath(resolve(root));
  const inside = (file) => file === root || file.startsWith(root + sep);
  const server = createServer(async (req, res) => {
    let fh;
    try {
      const url = new URL(req.url, "http://x");
      let path = decodeURIComponent(url.pathname);
      if (path.endsWith("/")) path += "index.html";
      let file = normalize(join(root, path));
      // under the root: a directory beside it whose name starts with the root's is not
      if (!inside(file)) { res.writeHead(403).end(); return; }
      const type = TYPES[extname(file)] ?? "application/octet-stream";
      // A symlink inside the root must not expose a file outside it.
      file = await realpath(file).catch(() => null);
      if (!file) { res.writeHead(404).end("not found"); return; }
      if (!inside(file)) { res.writeHead(403).end(); return; }
      // Opening a FIFO can wait indefinitely; reject non-files first.
      if (!(await stat(file)).isFile()) { res.writeHead(404).end("not found"); return; }
      fh = await open(file);
      const st = await fh.stat();
      if (!st.isFile()) { res.writeHead(404).end("not found"); return; }
      const range = req.headers.range?.match(/^bytes=(\d*)-(\d*)$/);
      let start = 0, end = st.size - 1, status = 200;
      if (range && (range[1] || range[2])) {
        const first = Number(range[1]), last = Number(range[2]);
        start = range[1] ? first : Math.max(0, st.size - last);
        end = range[1] && range[2] ? Math.min(last, st.size - 1) : st.size - 1;
        if (!Number.isSafeInteger(first) || !Number.isSafeInteger(last) || start > end || start >= st.size) {
          res.writeHead(416, { "Content-Range": `bytes */${st.size}`, "Content-Length": 0 }).end();
          return;
        }
        status = 206;
      }
      res.writeHead(status, {
        "Content-Type": type, "Content-Length": end - start + 1, "Accept-Ranges": "bytes",
        ...(status === 206 ? { "Content-Range": `bytes ${start}-${end}/${st.size}` } : {}),
        "Cache-Control": "no-store", "Access-Control-Allow-Origin": "*",
      });
      if (req.method === "HEAD" || st.size === 0) { res.end(); return; }
      // Backpressure bounds memory per request, including when the client is
      // slow. The handle also closes on an invalid range or disconnected client.
      await pipeline(fh.createReadStream({ start, end, autoClose: false }), res);
    } catch (e) {
      if (res.headersSent) res.destroy(e);
      else res.writeHead(500).end("server error");
    } finally {
      await fh?.close().catch(() => {});
    }
  });
  return new Promise((resolve) => server.listen(port, "127.0.0.1", () => resolve({ server, port: server.address().port })));
}
