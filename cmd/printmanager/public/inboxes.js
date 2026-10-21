// The inboxes: the items each holds, released to a node of the inbox's
// choice (or any, when it names none) or discarded, with the worker's
// name for the record.

import { $, api, badge, button, el, onChange, sender, since, status } from "./admin.js";

async function load() {
  try {
    const topo = await api("GET", "/api/topology");
    const inboxes = topo.filter((n) => n.kind === "inbox");
    const cards = [];
    for (const n of inboxes) {
      const v = await api("GET", "/api/inboxes/" + n.name);
      const card = el("div", undefined, { className: "card bg-base-100 shadow-sm" });
      const body = el("div", undefined, { className: "card-body gap-3" });
      const head = el("div", undefined, { className: "flex flex-wrap items-center gap-2" });
      head.append(el("h3", v.title, { className: "card-title" }), badge(n.name, "badge-ghost font-mono"), badge("held " + v.held.length, v.held.length ? "badge-warning" : "badge-ghost"), badge("released or discarded " + v.ended.length));
      head.append(el("a", "Edit the node", { className: "btn btn-xs btn-ghost", href: "/admin/topology#" + encodeURIComponent(n.name) }));
      body.append(head);
      const table = el("table", undefined, { className: "table table-sm" });
      const tbody = el("tbody");
      for (const it of v.held) {
        const tr = el("tr");
        tr.append(el("td", "#" + it.job, { className: "font-mono" }), el("td", it.name || "", { className: "font-semibold" }), el("td", sender(it), { className: "text-sm" }), el("td", since(it.received), { className: "text-sm text-base-content/60" }), el("td"), el("td", undefined, { className: "whitespace-nowrap" }));
        const choice = tr.children[4];
        const boxes = [];
        if (v.to.length) {
          for (const to of v.to) {
            const l = el("label", undefined, { className: "mr-3 inline-flex items-center gap-1 text-sm" });
            const cb = el("input", undefined, { type: "checkbox", className: "checkbox checkbox-sm", checked: true, value: to });
            boxes.push(cb);
            l.append(cb, to);
            choice.append(l);
          }
        } else {
          const free = el("input", undefined, { className: "input input-sm w-48", placeholder: "nodes (comma-separated)" });
          boxes.push(free);
          choice.append(free);
        }
        const targets = () => (v.to.length ? boxes.filter((b) => b.checked).map((b) => b.value) : boxes[0].value.split(","));
        tr.lastChild.append(
          button("Release", () => api("POST", `/api/items/${it.id}/release`, { to: targets(), worker: $("worker").value }), "btn btn-sm btn-primary mr-1"),
          button("Discard", () => api("POST", `/api/items/${it.id}/discard`, { worker: $("worker").value }), "btn btn-sm btn-ghost text-error"),
        );
        tbody.append(tr);
      }
      if (!v.held.length) {
        const tr = el("tr");
        tr.append(el("td", "Nothing held", { colSpan: 6, className: "text-base-content/50" }));
        tbody.append(tr);
      }
      table.append(tbody);
      body.append(el("div", undefined, { className: "overflow-x-auto" }));
      body.lastChild.append(table);
      if (v.ended.length) {
        const details = el("details", undefined, { className: "collapse collapse-arrow bg-base-200" });
        details.append(el("summary", `Released or discarded: ${v.ended.length}`, { className: "collapse-title text-sm" }));
        const list = el("ul", undefined, { className: "collapse-content text-sm" });
        for (const it of v.ended) {
          const what = it.state === "released" ? `→ ${it.released.to.join(", ")} (${it.released.worker || ""} ${since(it.released.at)})` : `discarded (${it.discarded.worker || ""} ${since(it.discarded.at)})`;
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
      p.append("No inboxes yet. Add one on the ", el("a", "topology", { className: "link", href: "/admin/topology" }), ".");
      cards.push(p);
    }
    $("inboxes").replaceChildren(...cards);
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
