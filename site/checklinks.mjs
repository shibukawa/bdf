// Checks the documentation's internal links: after `node site/build.mjs`,
// every <a href="...#fragment"> from one page of the site to another must
// resolve to an id in the target page, and every plain internal link must
// resolve to a file that exists. External links (http:, mailto:) are not
// checked.
//
//   node site/build.mjs && node site/checklinks.mjs [out dir]
import { readFile, readdir, stat } from "node:fs/promises";
import { dirname, join, normalize, posix, resolve } from "node:path";

const out = resolve(process.argv[2] ?? "site/dist");

async function files(dir, ext, acc = []) {
  for (const e of await readdir(dir, { withFileTypes: true })) {
    const p = join(dir, e.name);
    if (e.isDirectory()) await files(p, ext, acc);
    else if (e.name.endsWith(ext)) acc.push(p);
  }
  return acc;
}

const idCache = new Map();
async function ids(file) {
  if (idCache.has(file)) return idCache.get(file);
  let html;
  try {
    html = await readFile(file, "utf8");
  } catch {
    idCache.set(file, null);
    return null;
  }
  const set = new Set([...html.matchAll(/\sid="([^"]+)"/g)].map((m) => m[1]));
  idCache.set(file, set);
  return set;
}

let problems = 0;
for (const file of await files(out, ".html")) {
  // samples/ holds input fixtures (test documents opened by the converters),
  // not authored pages: their own markup (e.g. a fake article's nav links)
  // is content to convert, not a link to check.
  if (posix.relative(out, file).startsWith("samples" + posix.sep)) continue;
  const html = await readFile(file, "utf8");
  const dir = dirname(file);
  for (const m of html.matchAll(/<a\s[^>]*href="([^"]+)"/g)) {
    const href = m[1];
    if (/^([a-z][a-z\d+.-]*:|#)/i.test(href)) continue; // external, or same-page fragment (browsers handle those)
    const i = href.search(/[?#]/);
    const pathPart = i < 0 ? href : href.slice(0, i);
    const hash = href.includes("#") ? href.slice(href.indexOf("#") + 1) : "";
    if (!pathPart) continue; // "?x" or "#foo" on the same page
    let target = normalize(join(dir, decodeURI(pathPart)));
    if (!target.startsWith(out)) continue; // left the site (should not happen for a relative link)
    let st = await stat(target).catch(() => null);
    if (st?.isDirectory()) {
      target = join(target, "index.html"); // a directory link (e.g. "../viewer/") serves its index.html
      st = await stat(target).catch(() => null);
    }
    if (!st || !st.isFile()) {
      console.log(`${posix.relative(out, file)}: broken link ${href} (no file ${posix.relative(out, target)})`);
      problems++;
      continue;
    }
    if (hash) {
      const set = await ids(target);
      if (set && !set.has(hash)) {
        console.log(`${posix.relative(out, file)}: broken anchor ${href} (no id="${hash}" in ${posix.relative(out, target)})`);
        problems++;
      }
    }
  }
}
console.log(problems ? `${problems} problem(s)` : "all internal links resolve");
process.exit(problems ? 1 : 0);
