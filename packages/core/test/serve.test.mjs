// The static server of the tests and the demo (test/serve.mjs).
import { test } from "node:test";
import assert from "node:assert/strict";
import { get } from "node:http";
import { mkdtemp, mkdir, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { serve } from "../../../test/serve.mjs";

// http.get sends the path as it is; fetch would take the dots out of it
const request = (port, path, headers = {}) => new Promise((resolve, reject) => {
  get({ host: "127.0.0.1", port, path, headers }, (res) => {
    let body = "";
    res.on("data", (d) => (body += d));
    res.on("end", () => resolve({ status: res.statusCode, body }));
  }).on("error", reject);
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
    let port;
    ({ server, port } = await serve(join(dir, "site")));
    assert.deepEqual(await request(port, "/"), { status: 200, body: "index" });
    assert.deepEqual(await request(port, "/sub/file.json"), { status: 200, body: "{}" });
    assert.deepEqual(await request(port, "/sub/../index.html"), { status: 200, body: "index" });
    assert.equal((await request(port, "/index.html", { Range: "bytes=1-2" })).body, "nd");
    // a directory beside the root whose name starts with the root's
    for (const path of ["/..%2fsite-private/secret.txt", "/..%2Fsite-private%2Fsecret.txt", "/sub/..%2f..%2fsite-private/secret.txt", "/..%2fother/secret.txt", "/..%2f..%2f..%2f..%2fetc/passwd"]) {
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
