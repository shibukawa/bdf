// Builds the documentation pages of the demo site: the READMEs and the
// documents under docs/, rendered from Markdown (GFM) into docs/ under the
// site root, with a shared navigation bar, a table of contents from the h2 and
// h3 headings and one stylesheet (docs.css).
//
// Headings get GitHub's anchor ids, so links into a document work as they do
// on GitHub. Links between these documents point at their pages; other
// relative links (source files, directories) point at the repository on
// GitHub. Mermaid blocks are drawn in the browser by mermaid from a CDN.
import { statSync } from "node:fs";
import { mkdir, readFile, writeFile } from "node:fs/promises";
import { join, posix } from "node:path";
import { Marked, Renderer } from "marked";
import { root } from "./build.mjs";

const REPO = "https://github.com/shibukawa/bdf";
const SITE = "https://shibukawa.github.io/bdf/";
const MERMAID = "https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.esm.min.mjs";

/** The pages: Markdown source, output file, language and navigation label. */
const PAGES = [
  { src: "README.md", out: "index.html", lang: "en", nav: "Overview", label: "Overview" },
  { src: "README.ja.md", out: "ja.html", lang: "ja", nav: "日本語", label: "概要" },
  { src: "docs/api.md", out: "api.html", lang: "ja", nav: "API" },
  { src: "docs/spec.md", out: "spec.html", lang: "ja", nav: "Spec (仕様)" },
  { src: "docs/design.md", out: "design.html", lang: "ja", nav: "Design (設計)" },
];

/** The words of the page around the document, in the document's language. */
const UI = {
  en: { skip: "Skip to content", site: "Site", toc: "Contents", source: "Source" },
  ja: { skip: "本文へ移動", site: "サイト", toc: "目次", source: "原文" },
};

const esc = (s) => s.replace(/[&<>"]/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" })[c]);

/**
 * GitHub's heading anchors (github-slugger): lower case, anything but
 * letters, marks, numbers, connectors, spaces and hyphens removed, spaces
 * turned into hyphens; a repeated anchor gets -1, -2, … appended.
 */
function slugger() {
  const seen = new Map();
  return (text) => {
    const base = text.toLowerCase().replace(/[^\p{L}\p{M}\p{N}\p{Pc} -]/gu, "").replace(/ /g, "-");
    let id = base;
    while (seen.has(id)) {
      seen.set(base, seen.get(base) + 1);
      id = `${base}-${seen.get(base)}`;
    }
    seen.set(id, 0);
    return id;
  };
}

/** Where a link in the Markdown file src points on the site. */
function rewrite(href, src, image) {
  if (href.startsWith(SITE)) return `../${href.slice(SITE.length)}`;
  // absolute URLs, fragments and root-relative paths stay as they are
  if (/^([a-z][a-z\d+.-]*:|[/#])/i.test(href)) return href;
  const i = href.indexOf("#");
  const hash = i < 0 ? "" : href.slice(i);
  const path = posix.normalize(posix.join(posix.dirname(src), i < 0 ? href : href.slice(0, i))).replace(/\/$/, "");
  const page = PAGES.find((p) => p.src === path);
  if (page) return page.out + hash;
  if (path.startsWith("..")) return href;
  let kind = "blob";
  if (image) kind = "raw";
  else if (statSync(join(root, decodeURI(path)), { throwIfNoEntry: false })?.isDirectory()) kind = "tree";
  return `${REPO}/${kind}/main/${path}${hash}`;
}

/** Render one Markdown file: the body HTML, its title (the first h1), the h2/h3 outline and whether it has diagrams. */
function render(p, md) {
  const slug = slugger();
  const doc = { h1: "", toc: [], mermaid: false };
  const marked = new Marked({ gfm: true });
  marked.use({
    walkTokens(t) {
      if (t.type === "link" || t.type === "image") t.href = rewrite(t.href, p.src, t.type === "image");
    },
    renderer: {
      heading({ tokens, depth }) {
        const text = this.parser.parseInline(tokens, this.parser.textRenderer);
        const id = slug(text);
        if (depth === 1) doc.h1 ||= text;
        if (depth === 2 || depth === 3) doc.toc.push({ depth, id, text });
        return `<h${depth} id="${esc(id)}">${this.parser.parseInline(tokens)}</h${depth}>\n`;
      },
      code({ text, lang }) {
        if (lang?.match(/^\S*/)[0] !== "mermaid") return false;
        doc.mermaid = true;
        // escaped, so it still reads as the diagram's source if mermaid does not load
        return `<pre class="mermaid">${esc(text)}</pre>\n`;
      },
      // tables scroll sideways on their own when they are wider than the text
      table(token) {
        return `<div class="table-wrap">${Renderer.prototype.table.call(this, token)}</div>\n`;
      },
    },
  });
  doc.body = marked.parse(md);
  return doc;
}

/** The outline: h2 headings with their h3 headings nested. */
function outline(toc) {
  const items = [];
  for (const h of toc) {
    if (h.depth === 3 && items.length) items.at(-1).sub.push(h);
    else items.push({ ...h, sub: [] });
  }
  const link = (h) => `<a href="#${esc(h.id)}">${esc(h.text)}</a>`;
  const list = (hs) => `<ul>\n${hs.map((h) => `<li>${link(h)}${h.sub?.length ? list(h.sub) : ""}</li>\n`).join("")}</ul>`;
  return list(items);
}

function page(p, doc) {
  const ui = UI[p.lang];
  // a document titled with the project's name is named after the page
  const title = /^bdf$/i.test(doc.h1) ? `BDF – ${p.label}` : /\bbdf\b/i.test(doc.h1) ? doc.h1 : `${doc.h1} – BDF`;
  const nav = PAGES.map((q) => {
    const current = q === p ? ' aria-current="page"' : "";
    // the language switch is labelled in the language it switches to
    const lang = q.label && q.lang !== p.lang ? ` lang="${q.lang}" hreflang="${q.lang}"` : "";
    return `<li><a href="${q.out}"${current}${lang}>${esc(q.nav)}</a></li>`;
  });
  return `<!doctype html>
<html lang="${p.lang}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="color-scheme" content="light dark">
<title>${esc(title)}</title>
<link rel="stylesheet" href="docs.css">
</head>
<body>
<a class="skip" href="#content">${ui.skip}</a>
<header class="bar">
<nav aria-label="${ui.site}">
<a class="home" href="index.html">BDF</a>
<ul>
<li><a href="../">Viewer</a></li>
${nav.join("\n")}
<li><a href="${REPO}">GitHub</a></li>
</ul>
</nav>
</header>
<div class="layout">
<nav class="toc" aria-labelledby="toc-title">
<details id="toc" open><summary id="toc-title">${ui.toc}</summary>
${outline(doc.toc)}
</details>
</nav>
<script>
// the outline is a sidebar on wide screens and folded above the text on narrow ones
{
  const toc = document.getElementById("toc"), wide = matchMedia("(min-width: 64rem)");
  const fit = () => { toc.open = wide.matches; };
  fit();
  wide.addEventListener("change", fit);
}
</script>
<main id="content" tabindex="-1">
${doc.body}</main>
</div>
<footer class="foot">${ui.source}: <a href="${REPO}/blob/main/${p.src}">${p.src}</a></footer>
${doc.mermaid ? mermaidScript() : ""}</body>
</html>
`;
}

/** Draws the diagrams in the current color scheme, and again when it changes. */
function mermaidScript() {
  return `<script type="module">
import mermaid from "${MERMAID}";
const dark = matchMedia("(prefers-color-scheme: dark)");
const nodes = [...document.querySelectorAll("pre.mermaid")];
const sources = nodes.map((n) => n.textContent);
async function draw() {
  nodes.forEach((n, i) => { n.textContent = sources[i]; n.removeAttribute("data-processed"); });
  mermaid.initialize({ startOnLoad: false, theme: dark.matches ? "dark" : "default" });
  await mermaid.run({ nodes });
}
dark.addEventListener("change", draw);
await draw();
</script>
`;
}

const CSS = `/* The documentation pages of the demo site (examples/viewer/docs.mjs). */
:root {
  color-scheme: light dark;
  --bg: #ffffff;
  --fg: #1f2328;
  --muted: #59636e;
  --border: #d1d9e0;
  --subtle: #f6f8fa;
  --link: #0969da;
  --focus: #1a73e8;
  --sans: system-ui, -apple-system, "Segoe UI", "Hiragino Sans", "Hiragino Kaku Gothic ProN", "Noto Sans JP", "Noto Sans CJK JP", Meiryo, sans-serif;
  --mono: ui-monospace, SFMono-Regular, Menlo, Consolas, "Liberation Mono", monospace;
  /* the width of the text; with the outline beside it on wide screens */
  --page: 50rem;
}
@media (prefers-color-scheme: dark) {
  :root {
    --bg: #0d1117;
    --fg: #e6edf3;
    --muted: #9198a1;
    --border: #3d444d;
    --subtle: #161b22;
    --link: #4493f8;
    --focus: #58a6ff;
  }
}
@media (min-width: 64rem) {
  :root { --page: 68rem; }
}
*, *::before, *::after { box-sizing: border-box; }
html { text-size-adjust: 100%; -webkit-text-size-adjust: 100%; }
body { margin: 0; background: var(--bg); color: var(--fg); font: 1rem/1.7 var(--sans); overflow-wrap: break-word; }
html:lang(ja) body { line-height: 1.8; }
a { color: var(--link); }
:focus-visible { outline: 3px solid var(--focus); outline-offset: 2px; }
main:focus { outline: none; }

.skip { position: absolute; top: 8px; left: 16px; z-index: 1; padding: .5rem 1rem; background: var(--bg); border: 1px solid var(--border); border-radius: 6px; transform: translateY(-200%); }
.skip:focus { transform: none; }

.bar { background: var(--subtle); border-bottom: 1px solid var(--border); }
.bar nav, .layout, .foot { max-width: calc(var(--page) + 32px); margin: 0 auto; padding-inline: 16px; }
.bar nav { display: flex; flex-wrap: wrap; align-items: baseline; gap: .25rem 1.5rem; padding-block: .6rem; }
/* the links wrap beside the home link on narrow screens */
.bar ul { display: flex; flex: 1 1 0; flex-wrap: wrap; gap: .25rem 1rem; min-width: 0; margin: 0; padding: 0; list-style: none; }
.bar a { text-decoration: none; }
.bar a:hover { text-decoration: underline; }
.bar .home { color: var(--fg); font-size: 1.125rem; font-weight: 700; }
.bar a[aria-current="page"] { color: var(--fg); font-weight: 600; text-decoration: underline 2px; text-underline-offset: .35em; }

.toc { margin-top: 1rem; padding: .5rem .75rem; border: 1px solid var(--border); border-radius: 6px; font-size: .9375rem; line-height: 1.5; }
.toc summary { cursor: pointer; font-weight: 600; }
.toc ul { margin: .25rem 0; padding: 0; list-style: none; }
.toc ul ul { margin: 0; padding-left: 1rem; }
.toc a { display: block; padding: .2rem 0; color: var(--fg); text-decoration: none; }
.toc a:hover { color: var(--link); text-decoration: underline; }
@media (min-width: 64rem) {
  .layout { display: grid; grid-template-columns: 15rem minmax(0, 50rem); gap: 3rem; }
  .toc { position: sticky; top: 0; align-self: start; max-height: 100vh; overflow-y: auto; margin: 0; padding: 1.5rem 0; border: 0; }
}

main { padding: 1rem 0 3rem; }
main h1, main h2, main h3, main h4 { line-height: 1.3; scroll-margin-top: 1rem; }
main h1 { font-size: clamp(1.625rem, 1rem + 3vw, 2rem); margin: .5rem 0 1rem; }
main h2 { font-size: 1.5rem; margin: 2.5rem 0 1rem; padding-bottom: .3rem; border-bottom: 1px solid var(--border); }
main h3 { font-size: 1.25rem; margin: 2rem 0 .75rem; }
main h4 { font-size: 1rem; margin: 1.5rem 0 .5rem; }
main p, main ul, main ol, main pre, main blockquote, .table-wrap { margin: 0 0 1rem; }
main li > ul, main li > ol { margin: 0; }
main ul, main ol { padding-left: 1.75rem; }
main li + li { margin-top: .25rem; }
main hr { border: 0; border-top: 1px solid var(--border); margin: 2rem 0; }
main img { max-width: 100%; }
main blockquote { padding: 0 1rem; color: var(--muted); border-left: .25rem solid var(--border); }
code { font-family: var(--mono); font-size: .875em; padding: .1em .35em; background: var(--subtle); border-radius: 4px; }
pre { overflow-x: auto; padding: .75rem 1rem; font: .875rem/1.5 var(--mono); background: var(--subtle); border: 1px solid var(--border); border-radius: 6px; }
pre code { padding: 0; font-size: inherit; background: none; }
pre.mermaid[data-processed] { padding: 0; font: inherit; text-align: center; background: none; border: 0; }
.table-wrap { overflow-x: auto; }
table { border-collapse: collapse; font-size: .9375rem; line-height: 1.6; }
th, td { padding: .4rem .75rem; vertical-align: top; border: 1px solid var(--border); }
th { background: var(--subtle); }
th:not([align]) { text-align: start; }

.foot { padding-block: 1rem 2rem; color: var(--muted); font-size: .875rem; border-top: 1px solid var(--border); }
`;

/** Write the documentation pages into out/docs. */
export async function buildDocs(out) {
  const dir = join(out, "docs");
  await mkdir(dir, { recursive: true });
  await writeFile(join(dir, "docs.css"), CSS);
  await Promise.all(PAGES.map(async (p) => {
    const doc = render(p, await readFile(join(root, p.src), "utf8"));
    await writeFile(join(dir, p.out), page(p, doc));
  }));
  console.log(`docs: ${PAGES.map((p) => p.out).join(", ")}`);
}
