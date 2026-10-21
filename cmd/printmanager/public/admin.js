// The administrator's page: the topology (printers, displays, endpoints,
// queues, classes, inboxes, archives) drawn on a canvas and listed in a
// table, the templates, the inboxes, the jobs of a destination, the
// deliveries waiting, and the devices of the displays with their
// enrollment codes. Everything goes through the API, authorized
// with the session's cookie (the sign-in at /).

const $ = (id) => document.getElementById(id);

// status shows a message for a while, as a toast.
function status(text, bad) {
  if (!text) return;
  const t = document.createElement("div");
  t.role = "status";
  t.className = "alert shadow-lg " + (bad ? "alert-error" : "alert-success");
  t.textContent = text;
  $("toasts").append(t);
  setTimeout(() => t.remove(), bad ? 6000 : 3000);
}

// api calls the API with the session's cookie; a session that ended
// sends the page back to the sign-in.
async function api(method, path, body, raw) {
  const res = await fetch(path, {
    method,
    headers: body && !raw ? { "Content-Type": "application/json" } : {},
    body: body ? (raw ? body : JSON.stringify(body)) : undefined,
  });
  if (res.status === 401) {
    location.href = "/";
    throw new Error("signed out");
  }
  if (res.status === 204) return null;
  const out = await res.json();
  if (!res.ok) throw new Error(out.error || res.statusText);
  return out;
}

function el(tag, text, attrs = {}) {
  const e = document.createElement(tag);
  if (text !== undefined) e.textContent = text;
  Object.assign(e, attrs);
  return e;
}

function button(text, fn, className = "btn btn-sm") {
  const b = el("button", text, { type: "button", className });
  b.addEventListener("click", async () => {
    try {
      await fn();
      await load();
    } catch (e) {
      status(e.message, true);
    }
  });
  return b;
}

function since(t) {
  if (!t) return "—";
  return new Date(t).toLocaleString();
}

function stat(title, value) {
  const s = el("div", undefined, { className: "stat px-4 py-2" });
  s.append(el("div", title, { className: "stat-title" }), el("div", String(value), { className: "stat-value text-2xl" }));
  return s;
}

function badge(text, cls = "badge-ghost") {
  return el("span", text, { className: "badge " + cls });
}

const KINDS = { endpoint: "Endpoint", queue: "Queue", class: "Class", inbox: "Inbox", archive: "Archive", display: "Display", printer: "Printer" };
const KIND_BADGE = { display: "badge-info badge-soft", printer: "badge-neutral badge-soft", endpoint: "badge-accent badge-soft" };
const STATES = { processing: "processing", held: "held", completed: "completed", canceled: "canceled", aborted: "aborted" };
const ROUTING = ["endpoint", "queue", "class", "inbox", "archive"];
let topo = []; // the nodes, as /api/topology lists them
let layout = {}; // name -> {x, y}
let jobsDest = ""; // the destination whose jobs are shown
let worker = ""; // who releases items
let selected = ""; // the node in the form

async function load() {
  try {
    const printers = await api("GET", "/api/printers");
    const rows = printers.map((p) => {
      const tr = el("tr");
      tr.append(el("td", p.name, { className: "font-semibold" }), el("td", p.url, { className: "font-mono text-sm break-all" }), el("td"), el("td"));
      tr.children[2].append(badge(p.route));
      tr.lastChild.append(button("Delete", () => api("DELETE", "/api/printers/" + p.name), "btn btn-sm btn-ghost text-error"));
      return tr;
    });
    if (!rows.length) {
      const tr = el("tr");
      tr.append(el("td", "No printers yet", { colSpan: 4, className: "text-base-content/50" }));
      rows.push(tr);
    }
    $("printers").replaceChildren(...rows);

    const locations = await api("GET", "/api/locations");
    const cards = [];
    for (const l of locations) {
      const card = el("div", undefined, { className: "card bg-base-100 shadow-sm" });
      const body = el("div", undefined, { className: "card-body gap-3" });
      const head = el("div", undefined, { className: "flex flex-wrap items-center gap-2" });
      head.append(el("h3", l.title || l.name, { className: "card-title" }), badge(l.name, "badge-ghost font-mono"));
      head.append(
        l.fallback
          ? badge(`fallback: ${l.fallback.printer}${l.fallback.after ? " after " + l.fallback.after : " while no device is connected"}`, "badge-info badge-soft")
          : badge("no fallback printer"),
      );
      body.append(head, el("p", `IPP: /ipp/print/${l.name}`, { className: "font-mono text-sm text-base-content/60" }));
      const board = await api("GET", `/api/locations/${l.name}/board`);
      const stats = el("div", undefined, { className: "stats stats-horizontal bg-base-200" });
      stats.append(stat("Open", board.open.length), stat("Done and canceled", board.ended.length));
      body.append(stats);
      const devices = await api("GET", `/api/locations/${l.name}/devices`);
      const table = el("table", undefined, { className: "table table-sm" });
      const tbody = el("tbody");
      for (const d of devices) {
        const tr = el("tr");
        tr.append(el("td", d.name, { className: "font-semibold" }), el("td", "last seen " + since(d.last_seen), { className: "text-sm text-base-content/60" }), el("td"));
        tr.lastChild.append(button("Remove", () => api("DELETE", "/api/devices/" + d.id), "btn btn-xs btn-ghost text-error"));
        tbody.append(tr);
      }
      if (!devices.length) {
        const tr = el("tr");
        tr.append(el("td", "No devices yet", { className: "text-base-content/50" }));
        tbody.append(tr);
      }
      table.append(tbody);
      body.append(table);
      const actions = el("div", undefined, { className: "card-actions justify-end" });
      actions.append(
        button("Delete", () => api("DELETE", "/api/locations/" + l.name), "btn btn-sm btn-ghost text-error"),
        button(
          "Add a device (enrollment code)",
          async () => {
            const e = await api("POST", `/api/locations/${l.name}/enrollments`);
            $("qr").innerHTML = e.qr; // SVG the manager drew
            $("enroll-code").textContent = e.code.replace(/(.{5})/, "$1-");
            $("enroll-url").textContent = e.url;
            $("enroll-url").href = e.url;
            $("enroll-expires").textContent = `Valid once, until ${since(e.expires)}`;
            $("enroll-modal").showModal();
          },
          "btn btn-sm btn-primary",
        ),
      );
      body.append(actions);
      card.append(body);
      cards.push(card);
    }
    if (!cards.length) cards.push(el("p", "No displays yet", { className: "text-base-content/50" }));
    $("locations").replaceChildren(...cards);
    await loadTopology();
    await loadTemplates();
  } catch (e) {
    status(e.message, true);
  }
}

// The topology: the canvas, the table, the inboxes, the jobs of a
// destination and the deliveries waiting.

// targets are the nodes an edge goes to from a node, and what kind of
// edge: to (a queue's, an endpoint's, an inbox's), member (a class's) or
// fallback (a display's printer).
function edgesOf(n) {
  const out = [];
  for (const to of n.to || []) out.push({ to, kind: "to" });
  for (const to of n.members || []) out.push({ to, kind: "member" });
  if (n.fallback) out.push({ to: n.fallback.printer, kind: "fallback" });
  return out;
}

function edgesText(n) {
  switch (n.kind) {
    case "queue":
    case "endpoint":
      return "→ " + (n.to || []).join(", ");
    case "inbox":
      return n.to && n.to.length ? "→ " + n.to.join(", ") + " (a person chooses)" : "(a person names the node)";
    case "class":
      return (n.members || []).join(" | ") + (n.policy === "failover" ? " (first)" : " (in turn)");
    case "display":
      return n.fallback ? "fallback → " + n.fallback.printer : "";
  }
  return "";
}

function settingsText(n) {
  switch (n.kind) {
    case "printer":
      return n.url + (n.route ? " (" + n.route + ")" : "");
    case "archive":
      return n.store + " / " + n.collection;
    case "endpoint":
      return "POST " + n.in + (n.templates && n.templates.length ? " · templates: " + n.templates.join(", ") : " · any template");
    case "inbox":
    case "display":
      return n.keep ? "keep " + n.keep : "";
  }
  return "";
}

function legsOf(j) {
  return (j.legs || [])
    .map((l) => l.node + ":" + l.state + (l.member ? "(" + l.member + ")" : "") + (l.error ? " ! " + l.error : ""))
    .join(" → ");
}

async function loadTopology() {
  topo = await api("GET", "/api/topology");
  layout = (await api("GET", "/api/layout")).positions || {};
  drawCanvas();
  const rows = topo.map((n) => {
    const tr = el("tr", undefined, { className: n.name === selected ? "bg-primary/10" : "" });
    const name = el("td", undefined, { className: "font-semibold" });
    name.append(el("div", n.title && n.title !== n.name ? `${n.title} (${n.name})` : n.name), el("div", n.ipp, { className: "font-mono text-xs text-base-content/60" }));
    tr.append(name, el("td"), el("td", edgesText(n), { className: "font-mono text-sm" }), el("td", settingsText(n), { className: "font-mono text-xs break-all" }), el("td"));
    tr.children[1].append(badge(KINDS[n.kind] || n.kind, KIND_BADGE[n.kind] || "badge-primary badge-soft"), " ", scopeBadge(n));
    if (ROUTING.includes(n.kind)) {
      tr.lastChild.append(button("Edit", async () => select(n.name), "btn btn-xs btn-ghost mr-1"));
      if (n.kind === "endpoint") tr.lastChild.append(button("Rotate key", () => rotateKey(n.name), "btn btn-xs btn-ghost mr-1"));
      tr.lastChild.append(button("Delete", () => api("DELETE", "/api/nodes/" + n.name), "btn btn-xs btn-ghost text-error"));
    }
    return tr;
  });
  if (!rows.length) {
    const tr = el("tr");
    tr.append(el("td", "No nodes yet", { colSpan: 5, className: "text-base-content/50" }));
    rows.push(tr);
  }
  $("topology").replaceChildren(...rows);

  // the destination whose jobs are shown
  const sel = $("jobs-dest");
  if (!topo.some((n) => n.name === jobsDest)) jobsDest = topo.length ? topo[0].name : "";
  sel.replaceChildren(...topo.map((n) => el("option", n.name + " (" + (KINDS[n.kind] || n.kind) + ")", { value: n.name, selected: n.name === jobsDest })));

  await loadInboxes(topo.filter((n) => n.kind === "inbox"));
  await loadJobs();
  await loadDeliveries();
}

async function rotateKey(name) {
  const r = await api("POST", `/api/nodes/${name}/key`);
  showKey(name, r.key);
}

function showKey(name, key) {
  $("key-name").textContent = name;
  $("key-value").textContent = key;
  $("key-usage").textContent = `curl -X POST -H "Authorization: Bearer ${key}" --data-binary @document.pdf "${location.origin}/in/${name}?name=document"`;
  $("key-modal").showModal();
}

// The canvas: nodes as boxes, edges as arrows. Nodes without a position
// are laid out in columns by their depth from the entries.

const NODE_W = 168;
const NODE_H = 46;
const COL = 230;
const ROW = 72;
const PAD = 24;
const SVG = "http://www.w3.org/2000/svg";
const FILL = { endpoint: "#fde68a", queue: "#ddd6fe", class: "#c7d2fe", inbox: "#fed7aa", archive: "#e5e7eb", display: "#bfdbfe", printer: "#d1d5db" };

function svgEl(tag, attrs = {}) {
  const e = document.createElementNS(SVG, tag);
  for (const [k, v] of Object.entries(attrs)) e.setAttribute(k, v);
  return e;
}

function autoLayout() {
  const byName = Object.fromEntries(topo.map((n) => [n.name, n]));
  const incoming = {};
  for (const n of topo) for (const e of edgesOf(n)) incoming[e.to] = (incoming[e.to] || 0) + 1;
  const depth = {};
  const walk = (name, d) => {
    if ((depth[name] ?? -1) >= d) return;
    depth[name] = d;
    if (d > 20) return;
    for (const e of edgesOf(byName[name] || {})) if (byName[e.to]) walk(e.to, d + 1);
  };
  for (const n of topo) if (!incoming[n.name]) walk(n.name, 0);
  for (const n of topo) if (depth[n.name] === undefined) walk(n.name, 0);
  const columns = {};
  for (const n of topo) {
    const d = depth[n.name];
    columns[d] = columns[d] || [];
    columns[d].push(n.name);
  }
  const taken = new Set(Object.values(layout).map((p) => `${Math.round(p.x)},${Math.round(p.y)}`));
  for (const [d, names] of Object.entries(columns)) {
    let row = 0;
    for (const name of names) {
      if (layout[name]) continue;
      let p;
      do {
        p = { x: PAD + d * COL, y: PAD + row * ROW };
        row++;
      } while (taken.has(`${p.x},${p.y}`));
      layout[name] = p;
    }
  }
}

function pos(name) {
  return layout[name] || { x: PAD, y: PAD };
}

function drawCanvas() {
  autoLayout();
  const svg = $("canvas");
  svg.replaceChildren();
  const defs = svgEl("defs");
  const marker = svgEl("marker", { id: "arrow", viewBox: "0 0 10 10", refX: "10", refY: "5", markerWidth: "8", markerHeight: "8", orient: "auto-start-reverse" });
  marker.append(svgEl("path", { d: "M 0 0 L 10 5 L 0 10 z", fill: "#64748b" }));
  defs.append(marker);
  svg.append(defs);
  let maxX = 600;
  let maxY = 200;
  for (const n of topo) {
    const p = pos(n.name);
    maxX = Math.max(maxX, p.x + NODE_W + PAD);
    maxY = Math.max(maxY, p.y + NODE_H + PAD);
  }
  svg.setAttribute("width", maxX);
  svg.setAttribute("height", maxY);
  svg.setAttribute("viewBox", `0 0 ${maxX} ${maxY}`);
  const edges = svgEl("g");
  for (const n of topo) {
    for (const e of edgesOf(n)) {
      if (!topo.some((x) => x.name === e.to)) continue;
      const a = pos(n.name);
      const b = pos(e.to);
      const x1 = a.x + NODE_W;
      const y1 = a.y + NODE_H / 2;
      const x2 = b.x;
      const y2 = b.y + NODE_H / 2;
      const dx = Math.max(40, Math.abs(x2 - x1) / 2);
      const d = `M ${x1} ${y1} C ${x1 + dx} ${y1}, ${x2 - dx} ${y2}, ${x2} ${y2}`;
      const hit = svgEl("path", { d, fill: "none", stroke: "transparent", "stroke-width": "14", class: e.kind === "fallback" ? "" : "cursor-pointer" });
      const line = svgEl("path", { d, fill: "none", stroke: "#64748b", "stroke-width": "2", "marker-end": "url(#arrow)", "pointer-events": "none" });
      if (e.kind === "member") line.setAttribute("stroke-dasharray", "6 4");
      if (e.kind === "fallback") line.setAttribute("stroke-dasharray", "2 4");
      const title = svgEl("title");
      title.textContent = e.kind === "fallback" ? `${n.name} prints on ${e.to} as a fallback` : `${n.name} → ${e.to} (click to remove)`;
      hit.append(title);
      if (e.kind !== "fallback") hit.addEventListener("click", () => removeEdge(n, e));
      edges.append(hit, line);
    }
  }
  svg.append(edges);
  for (const n of topo) {
    const p = pos(n.name);
    const g = svgEl("g", { transform: `translate(${p.x} ${p.y})`, class: "cursor-grab", "data-name": n.name });
    g.append(svgEl("rect", { width: NODE_W, height: NODE_H, rx: "8", fill: FILL[n.kind] || "#e5e7eb", stroke: n.name === selected ? "#4f46e5" : "#94a3b8", "stroke-width": n.name === selected ? "3" : "1.5" }));
    const name = svgEl("text", { x: "10", y: "19", "font-size": "13", "font-weight": "600", fill: "#111827" });
    name.textContent = n.title && n.title !== n.name ? n.title : n.name;
    const kind = svgEl("text", { x: "10", y: "36", "font-size": "11", fill: "#374151" });
    kind.textContent = (KINDS[n.kind] || n.kind) + (n.title && n.title !== n.name ? " · " + n.name : "") + (n.held ? ` · ${n.held} held` : "");
    g.append(name, kind);
    g.append(svgEl("circle", { cx: "0", cy: NODE_H / 2, r: "5", fill: "#f8fafc", stroke: "#64748b", "stroke-width": "1.5" }));
    if (["queue", "endpoint", "inbox", "class"].includes(n.kind)) {
      const port = svgEl("circle", { cx: NODE_W, cy: NODE_H / 2, r: "7", fill: "#4f46e5", stroke: "#fff", "stroke-width": "2", class: "cursor-crosshair", "data-port": n.name });
      const t = svgEl("title");
      t.textContent = "Drag onto a node to connect";
      port.append(t);
      g.append(port);
    }
    g.addEventListener("pointerdown", (ev) => startDrag(ev, n));
    svg.append(g);
  }
}

// Dragging a node moves it; dragging from its port draws an edge.
let drag = null;

function canvasPoint(ev) {
  const svg = $("canvas");
  const r = svg.getBoundingClientRect();
  return { x: ev.clientX - r.left, y: ev.clientY - r.top };
}

function startDrag(ev, n) {
  if (ev.button !== 0) return;
  ev.preventDefault();
  const at = canvasPoint(ev);
  const p = pos(n.name);
  const svg = $("canvas");
  if (ev.target.dataset.port) {
    const line = svgEl("path", { fill: "none", stroke: "#4f46e5", "stroke-width": "2", "stroke-dasharray": "4 3", "pointer-events": "none" });
    svg.append(line);
    drag = { kind: "edge", from: n, line, x1: p.x + NODE_W, y1: p.y + NODE_H / 2 };
  } else {
    drag = { kind: "move", node: n, dx: at.x - p.x, dy: at.y - p.y, moved: false };
  }
  svg.setPointerCapture(ev.pointerId);
}

$("canvas").addEventListener("pointermove", (ev) => {
  if (!drag) return;
  const at = canvasPoint(ev);
  if (drag.kind === "move") {
    layout[drag.node.name] = { x: Math.max(0, Math.round(at.x - drag.dx)), y: Math.max(0, Math.round(at.y - drag.dy)) };
    drag.moved = true;
    drawCanvas();
  } else {
    drag.line.setAttribute("d", `M ${drag.x1} ${drag.y1} L ${at.x} ${at.y}`);
  }
});

$("canvas").addEventListener("pointerup", async (ev) => {
  if (!drag) return;
  const d = drag;
  drag = null;
  try {
    if (d.kind === "move") {
      if (d.moved) {
        await api("PUT", "/api/layout", { positions: { [d.node.name]: layout[d.node.name] } });
      } else {
        select(d.node.name);
      }
      return;
    }
    d.line.remove();
    const target = document.elementFromPoint(ev.clientX, ev.clientY)?.closest("g[data-name]");
    if (!target) return;
    await addEdge(d.from, target.dataset.name);
  } catch (e) {
    status(e.message, true);
  }
});

// addEdge adds a target to a node (its to, or its members) through the
// API, which checks the edge (the node exists, no cycle).
async function addEdge(from, to) {
  if (from.name === to) return;
  const n = await api("GET", "/api/nodes/" + from.name);
  const field = n.kind === "class" ? "members" : "to";
  const list = n[field] || [];
  if (list.includes(to)) return;
  n[field] = [...list, to];
  delete n.name;
  await api("PUT", "/api/nodes/" + from.name, n);
  status(`${from.name} → ${to}`);
  await loadTopology();
}

async function removeEdge(from, e) {
  if (!confirm(`Remove the edge ${from.name} → ${e.to}?`)) return;
  try {
    const n = await api("GET", "/api/nodes/" + from.name);
    const field = e.kind === "member" ? "members" : "to";
    n[field] = (n[field] || []).filter((x) => x !== e.to);
    delete n.name;
    await api("PUT", "/api/nodes/" + from.name, n);
    await loadTopology();
  } catch (err) {
    status(err.message, true);
  }
}

// sender is who sent a job, as verified, and the name they gave.
function sender(j) {
  const by = j.submitted_by || {};
  const who = by.name || by.id || by.kind || "";
  return j.claimed_user && j.claimed_user !== who ? `${who} (as ${j.claimed_user})` : who;
}

// scopeBadge tells who may send to a node, and whose its jobs are.
function scopeBadge(n) {
  const scope = n.scope || "shared";
  const cls = { group: "badge badge-sm badge-soft badge-secondary whitespace-nowrap", personal: "badge badge-sm badge-soft badge-accent whitespace-nowrap" }[scope] || "badge badge-sm badge-ghost whitespace-nowrap";
  return el("span", scope + (n.owner ? ": " + n.owner : ""), { className: cls });
}

// scoped puts the scope and owner of a form in a body.
function scoped(form, body) {
  body.scope = $(form + "-scope").value;
  if ($(form + "-owner").value) body.owner = $(form + "-owner").value;
  return body;
}

// select puts a routing node in the form, to edit it.
async function select(name) {
  const n = topo.find((x) => x.name === name);
  if (!n || !ROUTING.includes(n.kind)) {
    selected = "";
    drawCanvas();
    return;
  }
  selected = name;
  $("node-kind").value = n.kind;
  $("node-name").value = n.name;
  $("node-title").value = n.title || "";
  $("node-to").value = (n.kind === "class" ? n.members : n.to || []).join(", ");
  $("node-policy").value = n.policy || "round-robin";
  $("node-templates").value = (n.templates || []).join(", ");
  $("node-keep").value = n.keep || "";
  $("node-store").value = n.store || "";
  $("node-collection").value = n.collection || "";
  $("node-scope").value = n.scope || "shared";
  $("node-owner").value = n.owner || "";
  nodeFields();
  drawCanvas();
  $("node-form").scrollIntoView({ behavior: "smooth", block: "nearest" });
}

async function loadInboxes(inboxes) {
  const cards = [];
  if (inboxes.length) {
    const head = el("div", undefined, { className: "flex flex-wrap items-center gap-3" });
    const input = el("input", undefined, { id: "worker", className: "input input-sm w-48", placeholder: "name (for the record)", value: worker });
    input.addEventListener("input", () => (worker = input.value));
    const label = el("label", undefined, { className: "flex items-center gap-2 text-sm" });
    label.append("Worker", input);
    head.append(el("h2", "Inboxes", { className: "text-xl font-bold" }), label);
    cards.push(head);
  }
  for (const n of inboxes) {
    const v = await api("GET", "/api/inboxes/" + n.name);
    const card = el("div", undefined, { className: "card bg-base-100 shadow-sm" });
    const body = el("div", undefined, { className: "card-body gap-3" });
    const head = el("div", undefined, { className: "flex flex-wrap items-center gap-2" });
    head.append(el("h3", v.title, { className: "card-title" }), badge(n.name, "badge-ghost font-mono"), badge("held " + v.held.length, v.held.length ? "badge-warning" : "badge-ghost"), badge("released or discarded " + v.ended.length));
    body.append(head);
    const table = el("table", undefined, { className: "table table-sm" });
    const tbody = el("tbody");
    for (const it of v.held) {
      const tr = el("tr");
      tr.append(el("td", "#" + it.job, { className: "font-mono" }), el("td", it.name || "", { className: "font-semibold" }), el("td", sender(it), { className: "text-sm" }), el("td", since(it.received), { className: "text-sm text-base-content/60" }), el("td"), el("td"));
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
        button("Release", () => api("POST", `/api/items/${it.id}/release`, { to: targets(), worker }), "btn btn-sm btn-primary mr-1"),
        button("Discard", () => api("POST", `/api/items/${it.id}/discard`, { worker }), "btn btn-sm btn-ghost text-error"),
      );
      tbody.append(tr);
    }
    if (!v.held.length) {
      const tr = el("tr");
      tr.append(el("td", "Nothing held", { colSpan: 6, className: "text-base-content/50" }));
      tbody.append(tr);
    }
    table.append(tbody);
    body.append(table);
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
  $("inboxes").replaceChildren(...cards);
}

async function loadJobs() {
  const rows = [];
  if (jobsDest) {
    const jobs = await api("GET", `/api/destinations/${jobsDest}/jobs?limit=50`);
    for (const j of jobs) {
      const tr = el("tr");
      tr.append(el("td", "#" + j.number, { className: "font-mono" }), el("td", j.name || "", { className: "font-semibold" }), el("td", sender(j), { className: "text-sm" }), el("td", since(j.received), { className: "text-sm text-base-content/60" }), el("td"), el("td", legsOf(j), { className: "font-mono text-xs" }), el("td"));
      tr.children[4].append(badge(STATES[j.state] || j.state, j.state === "completed" ? "badge-success badge-soft" : j.state === "held" ? "badge-warning" : j.state === "processing" ? "badge-info" : "badge-ghost"));
      if (!["completed", "canceled", "aborted"].includes(j.state) || (j.legs || []).some((l) => l.kind === "display" && l.state === "shown")) {
        tr.lastChild.append(button("Cancel", () => api("POST", `/api/jobs/${j.id}/cancel`), "btn btn-xs btn-ghost text-error"));
      }
      rows.push(tr);
    }
  }
  if (!rows.length) {
    const tr = el("tr");
    tr.append(el("td", "No jobs yet", { colSpan: 7, className: "text-base-content/50" }));
    rows.push(tr);
  }
  $("jobs").replaceChildren(...rows);
}

async function loadDeliveries() {
  const ds = await api("GET", "/api/deliveries");
  const rows = ds.map((d) => {
    const tr = el("tr");
    tr.append(el("td", "#" + d.job, { className: "font-mono" }), el("td", d.node), el("td", d.printer || "", { className: "font-mono text-sm" }), el("td", String(d.attempts || 0)), el("td", d.retry_at ? since(d.retry_at) : "now", { className: "text-sm" }), el("td", d.last_error || "", { className: "text-xs text-error break-all" }), el("td"));
    tr.lastChild.append(button("Retry now", () => api("POST", `/api/deliveries/${d.id}/retry`), "btn btn-xs"));
    return tr;
  });
  if (!rows.length) {
    const tr = el("tr");
    tr.append(el("td", "Nothing waiting", { colSpan: 7, className: "text-base-content/50" }));
    rows.push(tr);
  }
  $("deliveries").replaceChildren(...rows);
}

async function loadTemplates() {
  const ts = await api("GET", "/api/templates");
  const rows = ts.map((t) => {
    const tr = el("tr");
    const fields = (t.fields || []).map((f) => f.name).join(", ");
    tr.append(el("td", t.id, { className: "font-mono font-semibold" }), el("td", t.format), el("td", fields, { className: "font-mono text-xs break-all", title: fields }), el("td", `${Math.round(t.size / 1024)} KiB`, { className: "text-sm" }), el("td", since(t.uploaded), { className: "text-sm text-base-content/60" }), el("td"));
    tr.lastChild.append(button("Delete", () => api("DELETE", "/api/templates/" + t.id), "btn btn-xs btn-ghost text-error"));
    return tr;
  });
  if (!rows.length) {
    const tr = el("tr");
    tr.append(el("td", "No templates yet", { colSpan: 6, className: "text-base-content/50" }));
    rows.push(tr);
  }
  $("templates").replaceChildren(...rows);
}

$("jobs-dest").addEventListener("change", async (ev) => {
  jobsDest = ev.target.value;
  try {
    await loadJobs();
  } catch (e) {
    status(e.message, true);
  }
});

// the node form shows the fields of the kind chosen
function nodeFields() {
  const kind = $("node-kind").value;
  $("node-to").hidden = kind === "archive";
  $("node-to").placeholder = kind === "class" ? "members (comma-separated): mfp-3f, mfp-4f" : kind === "inbox" ? "a person may release to (comma-separated; empty: any node)" : "forwards to (comma-separated): wh-a, wh-b";
  $("node-policy").hidden = kind !== "class";
  $("node-templates").hidden = kind !== "endpoint";
  $("node-keep").hidden = kind !== "inbox";
  $("node-store").hidden = kind !== "archive";
  $("node-collection").hidden = kind !== "archive";
}
$("node-kind").addEventListener("change", nodeFields);
nodeFields();
const names = (s) => s.split(",").map((x) => x.trim()).filter(Boolean);
$("node-form").addEventListener("submit", async (ev) => {
  ev.preventDefault();
  const kind = $("node-kind").value;
  const body = scoped("node", { kind, title: $("node-title").value });
  if (kind === "class") {
    body.members = names($("node-to").value);
    body.policy = $("node-policy").value;
  } else if (kind === "archive") {
    body.store = $("node-store").value;
    if ($("node-collection").value) body.collection = $("node-collection").value;
  } else {
    body.to = names($("node-to").value);
    if (kind === "inbox" && $("node-keep").value) body.keep = $("node-keep").value;
    if (kind === "endpoint") body.templates = names($("node-templates").value);
  }
  try {
    const name = $("node-name").value;
    const n = await api("PUT", "/api/nodes/" + name, body);
    ev.target.reset();
    selected = "";
    nodeFields();
    if (n.key) showKey(name, n.key);
    await load();
  } catch (e) {
    status(e.message, true);
  }
});
$("template-form").addEventListener("submit", async (ev) => {
  ev.preventDefault();
  const file = $("template-file").files[0];
  if (!file) return;
  try {
    await api("PUT", `/api/templates/${$("template-id").value}?name=${encodeURIComponent(file.name)}`, file, true);
    ev.target.reset();
    await loadTemplates();
  } catch (e) {
    status(e.message, true);
  }
});
$("tick").addEventListener("click", async () => {
  try {
    const r = await api("POST", "/api/tick");
    status(`tick: delivered ${r.delivered}, failed ${r.undelivered}, printed ${r.printed} (${r.failed} failed), removed ${r.removed + r.items + r.jobs}`);
    await load();
  } catch (e) {
    status(e.message, true);
  }
});
$("printer-form").addEventListener("submit", async (ev) => {
  ev.preventDefault();
  try {
    await api("PUT", "/api/printers/" + $("printer-name").value, scoped("printer", { url: $("printer-url").value }));
    ev.target.reset();
    await load();
  } catch (e) {
    status(e.message, true);
  }
});
$("location-form").addEventListener("submit", async (ev) => {
  ev.preventDefault();
  const body = scoped("location", { title: $("location-title").value });
  if ($("location-printer").value) {
    body.fallback = { printer: $("location-printer").value };
    if ($("location-after").value) body.fallback.after = $("location-after").value;
  }
  try {
    await api("PUT", "/api/locations/" + $("location-name").value, body);
    ev.target.reset();
    await load();
  } catch (e) {
    status(e.message, true);
  }
});
load();
