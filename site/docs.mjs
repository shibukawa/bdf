// Builds the documentation pages of the site from the Markdown (GFM) under
// docs/ (nav.mjs lists them), into docs/ under the site root: each page with
// the site's header, the documentation's sidebar and the outline of its own
// h2 and h3 headings.
//
// Headings get GitHub's anchor ids, so links into a document work as they do
// on GitHub. Links between the documents point at their pages, links to the
// site at its pages, pictures under docs/images/ at their copies; other
// relative links (source files, directories) point at the repository on
// GitHub. Mermaid blocks are drawn in the browser by mermaid from a CDN.
import { statSync } from "node:fs";
import { cp, mkdir, readFile, writeFile } from "node:fs/promises";
import { dirname, join, posix } from "node:path";
import { Marked, Renderer } from "marked";
import { root } from "../examples/common/build.mjs";
import { LANGS, PAGES, REPO, SECTIONS, SITE, UI, output, source } from "./nav.mjs";

const MERMAID = "https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.esm.min.mjs";
const IMAGES = "docs/images/";

export const esc = (s) => s.replace(/[&<>"]/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" })[c]);

/** The way from a published file to another, both from the site's root. */
export const relative = (from, to) => posix.relative(posix.dirname(from), to) || posix.basename(to);

/** The way from a published file to the site's root: "", "../", "../../". */
export const toRoot = (from) => "../".repeat(from.split("/").length - 1);

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

/** Every Markdown file that is a page, with the file it is published as. */
const PUBLISHED = new Map(PAGES.flatMap((p) => LANGS.filter((l) => source(p, l)).map((l) => [source(p, l), output(p, l)])));

/** Where a link in the Markdown file src points, from the page out. */
function rewrite(href, src, out, image) {
  if (href.startsWith(SITE)) return toRoot(out) + href.slice(SITE.length);
  // absolute URLs, fragments and root-relative paths stay as they are
  if (/^([a-z][a-z\d+.-]*:|[/#])/i.test(href)) return href;
  const i = href.search(/[?#]/);
  const tail = i < 0 ? "" : href.slice(i);
  const path = posix.normalize(posix.join(posix.dirname(src), i < 0 ? href : href.slice(0, i))).replace(/\/$/, "");
  const page = PUBLISHED.get(path);
  if (page) return relative(out, page) + tail;
  if (path.startsWith(IMAGES)) return relative(out, path) + tail;
  if (path.startsWith("..")) return href;
  // a link to a file that is not there (say a page's .html, which is not in
  // the repository) would be a dead link on GitHub
  const entry = statSync(join(root, decodeURI(path)), { throwIfNoEntry: false });
  if (!entry) throw new Error(`${src}: the link ${href} points at ${path}, which is neither a page nor a file of the repository`);
  let kind = "blob";
  if (image) kind = "raw";
  else if (entry.isDirectory()) kind = "tree";
  return `${REPO}/${kind}/main/${path}${tail}`;
}

/** Render one Markdown file: the body HTML, its title (the first h1), the h2/h3 outline and whether it has diagrams. */
function render(src, out, md) {
  const slug = slugger();
  const doc = { h1: "", toc: [], mermaid: false };
  const marked = new Marked({ gfm: true });
  marked.use({
    walkTokens(t) {
      if (t.type === "link" || t.type === "image") t.href = rewrite(t.href, src, out, t.type === "image");
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
      // pictures are fetched as they come into view, and links to pictures of the site open them
      image({ href, title, text }) {
        return `<img src="${esc(href)}" alt="${esc(text)}"${title ? ` title="${esc(title)}"` : ""} loading="lazy">`;
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

/**
 * The header of every page of the site. out is the page, from the site's
 * root; current names the page of the header it is ("viewer", "thumbnail",
 * "text", "docs"); other is the same page in the other language.
 */
export function header({ lang, out, current, other }) {
  const ui = UI[lang], up = toRoot(out);
  const docs = lang === "ja" ? "docs/index.ja.html" : "docs/index.html";
  const home = lang === "ja" ? "index.ja.html" : "./";
  const item = (id, href, label) => `<li><a href="${up}${href}"${id === current ? ' aria-current="page"' : ""}>${label}</a></li>`;
  const switchTo = lang === "ja" ? "en" : "ja";
  return `<a class="skip" href="#content">${ui.skip}</a>
<header class="bar">
<nav aria-label="${ui.site}">
<a class="home" href="${up}${home}">BDF</a>
<ul>
${item("viewer", "viewer/", ui.viewer)}
${item("thumbnail", "thumbnail/", ui.thumbnail)}
${item("text", "text/", ui.text)}
${item("docs", docs, ui.docsNav)}
<li><a href="${REPO}">GitHub</a></li>
</ul>
${other ? `<a class="lang" href="${esc(relative(out, other))}" lang="${switchTo}" hreflang="${switchTo}">${ui.other}</a>` : ""}
</nav>
</header>`;
}

/** The head of every page of the site but its title. */
export function head({ out, other, lang }) {
  const switchTo = lang === "ja" ? "en" : "ja";
  return `<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="color-scheme" content="light dark">
<link rel="stylesheet" href="${toRoot(out)}site.css">${other ? `\n<link rel="alternate" hreflang="${switchTo}" href="${esc(relative(out, other))}">` : ""}`;
}

/** The sidebar: the pages of the documentation in the language of the page, by section. */
function sidebar(p, lang, out) {
  const ui = UI[lang];
  const sections = SECTIONS.map((s) => {
    const items = s.pages.map((q) => {
      const foreign = q.only && q.only !== lang;
      const current = q === p ? ' aria-current="page"' : "";
      return `<li><a href="${esc(relative(out, output(q, lang)))}"${current}${foreign ? ` hreflang="${q.only}"` : ""}>${esc(q.label[lang])}${foreign ? ui.onlyJa : ""}</a></li>`;
    });
    return `<li><span>${esc(s.title[lang])}</span>\n<ul>\n${items.join("\n")}\n</ul></li>`;
  });
  return `<nav class="sidebar" aria-label="${ui.docs}">
<details id="menu" open><summary>${ui.menu}</summary>
<ul>
${sections.join("\n")}
</ul>
</details>
</nav>`;
}

function page(p, lang, doc) {
  const ui = UI[lang];
  const src = source(p, lang), out = output(p, lang);
  // the reference documents are in Japanese alone: the way back to the other language is its first page
  const otherLang = lang === "ja" ? "en" : "ja";
  const other = output(p.only ? PAGES[0] : p, otherLang);
  const title = /\bbdf\b/i.test(doc.h1) ? doc.h1 : `${doc.h1} – BDF`;
  // a page's own outline is worth its room from three headings on
  const toc = doc.toc.length >= 3 ? `<nav class="toc" aria-labelledby="toc-title">
<h2 id="toc-title">${ui.toc}</h2>
${outline(doc.toc)}
</nav>` : "";
  return `<!doctype html>
<html lang="${lang}">
<head>
${head({ out, other, lang })}
<title>${esc(title)}</title>
</head>
<body class="doc">
${header({ lang, out, current: "docs", other })}
<div class="layout${toc ? "" : " no-toc"}">
${sidebar(p, lang, out)}
<script>
// the menu is a sidebar on wide screens and folded above the text on narrow ones
{
  const menu = document.getElementById("menu"), wide = matchMedia("(min-width: 60rem)");
  const fit = () => { menu.open = wide.matches; };
  fit();
  wide.addEventListener("change", fit);
}
</script>
<main id="content" tabindex="-1">
${doc.body}<footer class="foot">${ui.source}: <a href="${REPO}/blob/main/${src}">${src}</a></footer>
</main>
${toc}
</div>
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

/** A page that sends its readers on to another: the address a page had before. */
const moved = (to) => `<!doctype html>
<meta charset="utf-8">
<title>BDF</title>
<link rel="canonical" href="${to}">
<meta http-equiv="refresh" content="0; url=${to}">
<a href="${to}">${to}</a>
`;

/** Write the documentation pages, and the pictures they show, into out/docs. */
export async function buildDocs(out) {
  const written = [];
  await Promise.all(PAGES.flatMap((p) => LANGS.map(async (lang) => {
    const src = source(p, lang);
    if (!src) return;
    const file = output(p, lang);
    const doc = render(src, file, await readFile(join(root, src), "utf8"));
    await mkdir(dirname(join(out, file)), { recursive: true });
    await writeFile(join(out, file), page(p, lang, doc));
    written.push(file);
  })));
  await cp(join(root, IMAGES), join(out, IMAGES), { recursive: true });
  // the Japanese overview was docs/ja.html; the features were a page of their
  // own before they went into why and the architecture
  await writeFile(join(out, "docs/ja.html"), moved("index.ja.html"));
  await writeFile(join(out, "docs/features.html"), moved("why.html"));
  await writeFile(join(out, "docs/features.ja.html"), moved("why.ja.html"));
  console.log(`docs: ${written.length} pages`);
  return written.sort();
}
