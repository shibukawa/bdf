// The inboxes the person works (every one for the role topology, their
// groups' and their own for the others), through the API applications
// use too (/api/inbox): the items each holds, with the first page of
// each (a click shows every page in the lightbox, with the download),
// released to a node of the inbox's choice (or to any, when it names
// none: the role topology's alone) or discarded, with the worker's name
// for the record.

import { $, _, api, badge, button, el, onChange, pagesText, preview, sender, since, status } from "./admin.js";

let me = null; // /api/me: whether the person has the role topology

async function load() {
  try {
    if (!me) me = await api("GET", "/api/me");
    const inboxes = await api("GET", "/api/inbox");
    const cards = [];
    for (const n of inboxes) {
      const base = "/api/inbox/" + encodeURIComponent(n.name);
      const v = await api("GET", base);
      const card = el("div", undefined, { className: "card bg-base-100 shadow-sm", id: n.name });
      const body = el("div", undefined, { className: "card-body gap-3" });
      const head = el("div", undefined, { className: "flex flex-wrap items-center gap-2" });
      head.append(el("h3", v.title, { className: "card-title" }), badge(n.name, "badge-ghost font-mono"), badge(_("held {n}", { n: v.held.length }), v.held.length ? "badge-warning" : "badge-ghost"), badge(_("released or discarded {n}", { n: v.ended.length })));
      if (me.topology) head.append(el("a", _("Edit the node"), { className: "btn btn-xs btn-ghost", href: "/topology#" + encodeURIComponent(n.name) }));
      body.append(head);
      if (n.description) body.append(el("p", n.description, { className: "text-sm text-base-content/60" }));
      const table = el("table", undefined, { className: "table table-sm" });
      const tbody = el("tbody");
      for (const it of v.held) {
        const at = `${base}/${it.job}`;
        const tr = el("tr");
        const thumb = el("td", undefined, { className: "w-20" });
        if (it.pages) {
          const doc = { title: `#${it.job} ${it.name || ""}`, sub: `${pagesText(it.pages)} · ${since(it.received)}`, pages: it.pages, page: (n, w) => `${at}/pages/${n}.png?w=${w}`, href: `${at}/document` };
          const b = el("button", undefined, { type: "button", className: "cursor-pointer", title: _("Show the pages, with the download") });
          b.append(el("img", undefined, { src: `${at}/thumbnail?size=128`, alt: _("page 1 of {doc}", { doc: doc.title }), className: "h-16 w-auto rounded bg-white shadow-sm", loading: "lazy" }));
          b.addEventListener("click", () => preview(doc));
          thumb.append(b);
        } else {
          thumb.append(el("a", _("Document"), { className: "link text-xs", href: `${at}/document` }));
        }
        tr.append(thumb, el("td", "#" + it.job, { className: "font-mono" }), el("td", it.name || "", { className: "font-semibold" }), el("td", sender(it), { className: "text-sm" }), el("td", since(it.received), { className: "text-sm text-base-content/60" }), el("td"), el("td", undefined, { className: "whitespace-nowrap" }));
        const choice = tr.children[5];
        const boxes = [];
        if (v.to.length) {
          for (const to of v.to) {
            const l = el("label", undefined, { className: "mr-3 inline-flex items-center gap-1 text-sm" });
            const cb = el("input", undefined, { type: "checkbox", className: "checkbox checkbox-sm", checked: true, value: to });
            boxes.push(cb);
            l.append(cb, to);
            choice.append(l);
          }
        } else if (me.topology) {
          const free = el("input", undefined, { className: "input input-sm w-48", placeholder: _("nodes (comma-separated)") });
          boxes.push(free);
          choice.append(free);
        } else {
          choice.append(el("span", _("This inbox names no nodes: the role topology releases what it holds."), { className: "text-xs text-base-content/60" }));
        }
        const targets = () => (v.to.length ? boxes.filter((b) => b.checked).map((b) => b.value) : boxes[0].value.split(","));
        if (boxes.length) tr.lastChild.append(button(_("Release"), () => api("POST", `${at}/release`, { to: targets(), worker: $("worker").value }), "btn btn-sm btn-primary mr-1"));
        tr.lastChild.append(button(_("Discard"), () => api("POST", `${at}/discard`, { worker: $("worker").value }), "btn btn-sm btn-ghost text-error"));
        tbody.append(tr);
      }
      if (!v.held.length) {
        const tr = el("tr");
        tr.append(el("td", _("Nothing held"), { colSpan: 7, className: "text-base-content/50" }));
        tbody.append(tr);
      }
      table.append(tbody);
      body.append(el("div", undefined, { className: "overflow-x-auto" }));
      body.lastChild.append(table);
      if (v.ended.length) {
        const details = el("details", undefined, { className: "collapse collapse-arrow bg-base-200" });
        details.append(el("summary", _("Released or discarded: {n}", { n: v.ended.length }), { className: "collapse-title text-sm" }));
        const list = el("ul", undefined, { className: "collapse-content text-sm" });
        for (const it of v.ended) {
          const what = it.state === "released" ? `→ ${it.released.to.join(", ")} (${it.released.worker || ""} ${since(it.released.at)})` : _("discarded ({who} {when})", { who: it.discarded.worker || "", when: since(it.discarded.at) });
          list.append(el("li", `#${it.job} ${it.name || ""} ${what}`));
        }
        details.append(list);
        body.append(details);
      }
      card.append(body);
      cards.push(card);
    }
    if (!cards.length) {
      const p = el("p", undefined, { className: "text-base-content/50" });
      if (me.topology) p.append(_("No inboxes yet."), el("a", _("Add one on the topology."), { className: "link ms-1", href: "/topology" }));
      else p.append(_("No inbox is yours to work: an inbox of one of your groups shows here."));
      cards.push(p);
    }
    $("inboxes").replaceChildren(...cards);
    if (location.hash) document.getElementById(decodeURIComponent(location.hash.slice(1)))?.scrollIntoView();
  } catch (e) {
    status(e.message, true);
  }
}

try {
  $("worker").value = localStorage.getItem("bdf-pm-worker") || "";
} catch {}
$("worker").addEventListener("change", () => {
  try {
    localStorage.setItem("bdf-pm-worker", $("worker").value);
  } catch {}
});
$("refresh").addEventListener("click", load);
onChange(load);
load();
