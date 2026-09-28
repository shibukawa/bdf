// The static server of the tests and the demo (test/serve.mjs).
import { test } from "node:test";
import assert from "node:assert/strict";
import { request as httpRequest } from "node:http";
import { mkdtemp, mkdir, rm, symlink, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { serve } from "../../../test/serve.mjs";

// http.request sends the path as it is; fetch would take the dots out of it
const request = (port, path, headers = {}, method = "GET") => new Promise((resolve, reject) => {
  httpRequest({ host: "127.0.0.1", port, path, headers, method }, (res) => {
    let body = "";
    res.on("data", (d) => (body += d));
    res.on("error", reject);
    res.on("end", () => resolve({ status: res.statusCode, body, headers: res.headers }));
  }).on("error", reject).end();
});

test("only files under the root are served", async () => {
  const dir = await mkdtemp(join(tmpdir(), "bdf-serve-"));
  let server;
  try {
    for (const d of ["site/sub", "site-private", "other"]) await mkdir(join(dir, d), { recursive: true });
    await writeFile(join(dir, "site/index.html"), "index");
    await writeFile(join(dir, "site/sub/file.json"), "{}");
    await writeFile(join(dir, "site-private/secret.txt"), "secret");
    await writeFile(join(dir, "other/secret.txt"), "secret");
    await symlink("../other/secret.txt", join(dir, "site/leak.txt"));
    await symlink("../other", join(dir, "site/leak-dir"), "dir");
    await symlink("sub/file.json", join(dir, "site/alias.json"));
    await writeFile(join(dir, "site/content"), "<p>text</p>");
    await symlink("content", join(dir, "site/alias.html"));
    await symlink("site", join(dir, "site-link"), "dir");
    let port;
    ({ server, port } = await serve(join(dir, "site-link")));
    for (const [path, body] of [["/", "index"], ["/sub/file.json", "{}"], ["/sub/../index.html", "index"], ["/alias.json", "{}"]]) {
      const res = await request(port, path);
      assert.equal(res.status, 200, path);
      assert.equal(res.body, body, path);
    }
    assert.equal((await request(port, "/index.html", { Range: "bytes=1-2" })).body, "nd");
    assert.equal((await request(port, "/alias.html")).headers["content-type"], "text/html; charset=utf-8");
    assert.equal((await request(port, "/sub")).status, 404);
    // a directory beside the root whose name starts with the root's
    for (const path of ["/leak.txt", "/leak-dir/secret.txt", "/..%2fsite-private/secret.txt", "/..%2Fsite-private%2Fsecret.txt", "/sub/..%2f..%2fsite-private/secret.txt", "/..%2fother/secret.txt", "/..%2f..%2f..%2f..%2fetc/passwd"]) {
      const res = await request(port, path);
      assert.equal(res.status, 403, path);
      assert.equal(res.body, "", path);
    }
    assert.equal((await request(port, "/../site-private/secret.txt")).status, 404);
  } finally {
    server?.close();
    await rm(dir, { recursive: true, force: true });
  }
});

test("streamed files, HEAD and byte ranges return the expected bytes", async () => {
  const dir = await mkdtemp(join(tmpdir(), "bdf-ranges-"));
  let server;
  try {
    const large = "0123456789".repeat(50000);
    await writeFile(join(dir, "large.bdf"), large);
    await writeFile(join(dir, "small.bdf"), "index");
    await writeFile(join(dir, "empty.bdf"), "");
    let port;
    ({ server, port } = await serve(dir));
    const full = await request(port, "/large.bdf");
    assert.equal(full.body, large);
    assert.equal(full.headers["content-length"], String(large.length));
    const head = await request(port, "/large.bdf", {}, "HEAD");
    assert.equal(head.status, 200);
    assert.equal(head.body, "");
    assert.equal(head.headers["content-length"], String(large.length));
    for (const [range, start, end] of [["1-2", 1, 2], ["1-", 1, 4], ["-2", 3, 4], ["0-99", 0, 4]]) {
      const res = await request(port, "/small.bdf", { Range: `bytes=${range}` });
      assert.equal(res.status, 206, range);
      assert.equal(res.body, "index".slice(start, end + 1), range);
      assert.equal(res.headers["content-range"], `bytes ${start}-${end}/5`, range);
    }
    for (const range of ["3-1", "5-", "-0", "9007199254740992-"]) {
      const res = await request(port, "/small.bdf", { Range: `bytes=${range}` });
      assert.equal(res.status, 416, range);
      assert.equal(res.body, "", range);
      assert.equal(res.headers["content-range"], "bytes */5", range);
    }
    const empty = await request(port, "/empty.bdf");
    assert.equal(empty.status, 200);
    assert.equal(empty.headers["content-length"], "0");
    assert.equal((await request(port, "/empty.bdf", { Range: "bytes=0-" })).status, 416);
    assert.equal((await request(port, "/small.bdf")).body, "index");
  } finally {
    server?.close();
    await rm(dir, { recursive: true, force: true });
  }
});
