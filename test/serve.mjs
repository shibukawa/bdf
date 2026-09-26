// Minimal static server with Range support, used by tests and the demo.
import { createServer } from "node:http";
import { stat, open } from "node:fs/promises";
import { join, extname, normalize } from "node:path";

const TYPES = {
  ".html": "text/html; charset=utf-8", ".js": "text/javascript", ".mjs": "text/javascript", ".map": "application/json",
  ".json": "application/json", ".png": "image/png", ".bdf": "application/octet-stream", ".css": "text/css", ".ts": "text/plain",
  ".wasm": "application/wasm", ".ttf": "font/ttf", ".otf": "font/otf", ".ttc": "font/collection", ".otc": "font/collection",
};

export function serve(root, port = 0) {
  const server = createServer(async (req, res) => {
    try {
      const url = new URL(req.url, "http://x");
      let path = decodeURIComponent(url.pathname);
      if (path.endsWith("/")) path += "index.html";
      const file = normalize(join(root, path));
      if (!file.startsWith(root)) { res.writeHead(403).end(); return; }
      const st = await stat(file).catch(() => null);
      if (!st || !st.isFile()) { res.writeHead(404).end("not found"); return; }
      const type = TYPES[extname(file)] ?? "application/octet-stream";
      const range = req.headers.range?.match(/^bytes=(\d+)-(\d*)$/);
      let start = 0, end = st.size - 1, status = 200;
      if (range) {
        start = Number(range[1]);
        end = range[2] ? Math.min(Number(range[2]), st.size - 1) : st.size - 1;
        status = 206;
      }
      const fh = await open(file);
      const buf = Buffer.alloc(end - start + 1);
      await fh.read(buf, 0, buf.length, start);
      await fh.close();
      res.writeHead(status, {
        "Content-Type": type, "Content-Length": buf.length, "Accept-Ranges": "bytes",
        ...(status === 206 ? { "Content-Range": `bytes ${start}-${end}/${st.size}` } : {}),
        "Cache-Control": "no-store", "Access-Control-Allow-Origin": "*",
      });
      res.end(buf);
    } catch (e) {
      res.writeHead(500).end(String(e));
    }
  });
  return new Promise((resolve) => server.listen(port, "127.0.0.1", () => resolve({ server, port: server.address().port })));
}
