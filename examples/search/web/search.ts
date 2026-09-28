// The search box on search's index page (main.go's indexTemplate): queries
// /search (a thin proxy in front of Meilisearch, so its API key stays on
// the server) and lists the hits, each linking to the page it was found on.
interface Hit {
  file: string;
  bdf: string;
  view: string;
  page?: number;
  title: string;
  snippet: string; // HTML: Meilisearch wraps the match in <mark>
}

const form = document.getElementById("searchForm") as HTMLFormElement;
const input = document.getElementById("q") as HTMLInputElement;
const list = document.getElementById("hits") as HTMLUListElement;

async function search(q: string) {
  if (!q) {
    list.replaceChildren();
    return;
  }
  const res = await fetch(`/search?q=${encodeURIComponent(q)}`);
  const hits = (await res.json()) as Hit[] | { error: string };
  if (!Array.isArray(hits)) {
    list.replaceChildren(errorItem(hits.error));
    return;
  }
  if (!hits.length) {
    const li = document.createElement("li");
    li.textContent = "No matches.";
    list.replaceChildren(li);
    return;
  }
  list.replaceChildren(...hits.map(hitItem));
}

function hitItem(h: Hit): HTMLLIElement {
  const li = document.createElement("li");
  const a = document.createElement("a");
  // the target page (main.ts) reads the view and page off the hash, after it opens the document
  const hash = new URLSearchParams();
  if (h.view) hash.set("view", h.view);
  if (h.page) hash.set("page", String(h.page));
  const frag = hash.toString() ? `#${hash}` : "";
  a.href = `/view/?src=/cache/${encodeURIComponent(h.bdf)}${frag}`;
  a.textContent = `${h.title}${h.page ? ` — page ${h.page}` : ""}`;
  const p = document.createElement("p");
  p.className = "snippet";
  p.innerHTML = h.snippet; // Meilisearch's own <mark> highlighting; the query text, not arbitrary input
  li.append(a, p);
  return li;
}

function errorItem(message: string): HTMLLIElement {
  const li = document.createElement("li");
  li.textContent = `search unavailable: ${message}`;
  return li;
}

form.addEventListener("submit", (e) => {
  e.preventDefault();
  search(input.value).catch((e) => list.replaceChildren(errorItem((e as Error).message ?? String(e))));
});
