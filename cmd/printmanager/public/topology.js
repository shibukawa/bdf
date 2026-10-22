// The topology: everyone's page, and the editor of a person with the
// role topology. The nodes and the edges are drawn on the canvas
// (canvas.js: panned and zoomed, with the activity of the jobs the
// person may see running along the edges, polled with its revision).
// Editing: nodes of every kind are placed, connected by dragging from a
// node's output port onto another node, and set up in the sidebar (the
// fields of the node's kind, and what is at it: an endpoint's templates
// and key, a display's devices and enrollment codes); everything goes
// through the API, and the server checks a node when it is put (its
// edges, its name, no cycle, and that every node takes the kinds of
// payload that reach it along a new edge; a change of settings may
// leave an edge broken, which the canvas draws dashed and the sidebar
// explains). Looking: a node shows the jobs sent to it that the person
// sees, and a form to print there when they may.

import { $, api, button, changed, el, keptSection, kindBadge, KINDS, legsOf, listedSection, onChange, PAYLOADS, sender, since, stateBadge, status } from "./admin.js";
import { createCanvas, iconOf, legend, payloadColor, roleColor, typeOf } from "./canvas.js";

const SCHEMES = ["ipp", "ipps", "socket", "lpd"];
const ACCEPTS = { file: "documents (PDF, Office, images, .bdf …)", data: "data (JSON, YAML, CSV, Excel, Parquet)" };
const OUTPUTS = { bdf: "the .bdf document made of it", file: "the file as received (a convert node makes the document later)" };
const KINDS_OF_PAYLOAD = { "": "any", bdf: ".bdf document", data: "data", any: "any other file" };
const SUBJECTS = { email: "email: the person's email, as the IdP verified it", sub: "sub: the person's subject at the IdP", id: "id: the person's ID here", group: "group: any of the person's groups" };
const POLL_MS = 4000;

let topo = []; // the nodes, as /api/topology lists them
let formats = []; // the formats an endpoint may take
let me = null; // /api/me: whether this process has ports of its own
let apps = []; // the apps, for an endpoint's allow-list
let selected = null; // the node in the sidebar
let draft = null; // a node being added: {kind}
let fitted = false;
let etag = ""; // the activity's revision
let live = null; // the activity, as /api/topology/live last told it
// editable: the server rendered the editor's controls (the palette) for
// a person with the role topology; anyone else looks, and prints where
// they may. The API checks the role itself.
const editable = !!$("palette");
let dests = []; // looking: the nodes the person may send to
let jobs = []; // looking: the jobs the person sees, latest first

const canvas = createCanvas($("canvas"), {
  editable,
  onSelect(name) {
    selected = name;
    draft = null;
    renderSidebar();
    $("sidebar").scrollTop = 0;
  },
  onMove(name, p) {
    api("PUT", "/api/layout", { positions: { [name]: p } }).catch((e) => status(e.message, true));
  },
  onConnect(from, to) {
    addEdge(from, to).catch((e) => status(e.message, true));
  },
  onReject(from, to, why) {
    status(`${from.name} → ${to}: ${why}`, true);
  },
  onRemoveEdge: removeEdge,
  draftLabel: () => $("f-name")?.value || "",
  draftNote: () => $("f-description")?.value || "",
});

async function load() {
  try {
    const [t, l] = await Promise.all([api("GET", "/api/topology"), api("GET", "/api/layout")]);
    topo = t;
    if (editable) {
      if (!formats.length) formats = await api("GET", "/api/formats");
      if (!me) me = await api("GET", "/api/me");
      apps = await api("GET", "/api/apps");
      $("node-names").replaceChildren(...topo.map((n) => el("option", undefined, { value: n.name })));
    } else {
      [dests, jobs] = await Promise.all([api("GET", "/api/destinations"), api("GET", "/api/my/jobs?limit=200")]);
    }
    if (selected && !topo.some((n) => n.name === selected)) selected = null;
    canvas.setData(topo, l.positions || {});
    if (!fitted) {
      canvas.fit();
      fitted = true;
    }
    renderSidebar();
  } catch (e) {
    status(e.message, true);
  }
}

// poll reads the activity (the jobs the person sees; every job, for an
// administrator) when its revision changed, and, looking, the jobs with
// it.
async function poll() {
  try {
    const res = await fetch("/api/topology/live", { headers: etag ? { "If-None-Match": etag } : {} });
    if (res.status === 401) {
      location.href = "/";
      return;
    }
    if (res.status === 200) {
      etag = res.headers.get("ETag") || "";
      live = await res.json();
      canvas.setLive(live);
      if (!editable) {
        jobs = await api("GET", "/api/my/jobs?limit=200");
        renderSidebar();
      }
    } else if (res.status !== 304) {
      throw new Error("live " + res.status);
    }
    $("live").textContent = "Live · " + new Date().toLocaleTimeString(undefined, { hour: "2-digit", minute: "2-digit", second: "2-digit" });
    $("live").className = "badge badge-success badge-soft badge-sm";
  } catch (e) {
    $("live").textContent = "Not connected: " + e.message;
    $("live").className = "badge badge-error badge-sm";
  }
}

$("zoom-in").addEventListener("click", () => canvas.zoomBy(1.25));
$("zoom-out").addEventListener("click", () => canvas.zoomBy(0.8));
$("zoom-fit").addEventListener("click", () => canvas.fit());
if (editable) $("refresh").addEventListener("click", load);
legend($("legend"));
// the palette's kinds get their icons, in their roles' colors
for (const a of document.querySelectorAll("#palette a[data-kind]")) a.prepend(kindIcon(a.dataset.kind, 18));

// kindIcon is the icon of a kind, in its role's color.
function kindIcon(kind, size) {
  const i = iconOf(kind, size);
  i.style.color = roleColor(kind);
  i.classList.add("shrink-0");
  return i;
}

// brokenEdges names the edges that no longer carry what the node after
// them takes.
const brokenEdges = (t) => t.flatMap((n) => Object.entries(n.broken || {}).map(([to, kinds]) => `${n.name} → ${to} (${kinds.map((k) => PAYLOADS[k] || k).join(", ")})`));

// addEdge connects a node to a target through the API, which checks the
// edge (the target exists, no cycle): a class's member, a display's
// fallback printer, or the node's to.
async function addEdge(from, to) {
  if (from.name === to) return;
  const n = await api("GET", "/api/nodes/" + from.name);
  if (n.kind === "display") {
    n.fallback = { ...(n.fallback || {}), printer: to };
  } else if (n.kind === "filter") {
    if ((n.otherwise || []).includes(to)) return;
    n.otherwise = [...(n.otherwise || []), to];
  } else {
    const field = n.kind === "class" ? "members" : "to";
    const list = n[field] || [];
    if (list.includes(to)) return;
    n[field] = [...list, to];
  }
  delete n.name;
  await api("PUT", "/api/nodes/" + from.name, n);
  status(`${from.name} → ${to}`);
  await load();
}

async function removeEdge(from, e) {
  const what = e.kind === "fallback" ? `the fallback printer ${e.to} of ${from.name}` : `the edge ${from.name} → ${e.to}`;
  if (!confirm(`Remove ${what}?`)) return;
  try {
    const n = await api("GET", "/api/nodes/" + from.name);
    if (e.kind === "fallback") n.fallback = null;
    else if (n.kind === "filter") {
      n.otherwise = (n.otherwise || []).filter((x) => x !== e.to);
      n.rules = (n.rules || []).map((r) => ({ ...r, to: (r.to || []).filter((x) => x !== e.to) })).filter((r) => r.to.length);
    } else {
      const field = e.kind === "member" ? "members" : "to";
      n[field] = (n[field] || []).filter((x) => x !== e.to);
    }
    delete n.name;
    await api("PUT", "/api/nodes/" + from.name, n);
    await load();
  } catch (err) {
    status(err.message, true);
  }
}

// Adding a node: the palette makes a draft, placed at the centre of the
// view clear of the others; the sidebar takes its name and settings, and
// Add puts it.
if (editable) $("palette").addEventListener("click", (ev) => {
  const a = ev.target.closest("a[data-kind]");
  if (!a) return;
  ev.preventDefault();
  document.activeElement?.blur();
  draft = { kind: a.dataset.kind, name: "", title: "" };
  selected = null;
  canvas.setDraft(draft, canvas.freeSpot(canvas.center()));
  renderSidebar();
  $("sidebar").scrollTop = 0;
  $("f-name")?.focus();
});

// The sidebar: the node selected (or the draft), its kind's fields, and
// what is at it.

function field(label, input, hint) {
  const f = el("fieldset", undefined, { className: "fieldset" });
  f.append(el("legend", label, { className: "fieldset-legend" }), input);
  if (hint) f.append(el("p", hint, { className: "label whitespace-normal" }));
  return f;
}

function input(id, value, attrs = {}) {
  // list is read-only as a property: an attribute
  const { list, ...rest } = attrs;
  const e = el("input", undefined, { id, className: "input input-sm w-full", value: value ?? "", ...rest });
  if (list) e.setAttribute("list", list);
  return e;
}

// picker is the nodes a node may send to, as checkboxes (name: the
// checkboxes' name; except the node itself and the entries), each with a
// dot of what it takes; one that would not take what the node carries
// says so.
function picker(name, chosen, except) {
  const box = el("div", undefined, { className: "flex flex-wrap gap-x-3 gap-y-1" });
  const from = draft || topo.find((x) => x.name === selected) || {};
  const carries = typeOf(from).out;
  const names = topo.filter((x) => x.name !== except && x.kind !== "cups").map((x) => x.name);
  for (const x of names) {
    const target = topo.find((t) => t.name === x);
    const t = typeOf(target);
    const bad = carries.filter((k) => !t.accepts.includes(k));
    const text = (kinds) => kinds.map((k) => PAYLOADS[k] || k).join(" or ");
    const l = el("label", undefined, { className: "inline-flex items-center gap-1 text-sm" + (bad.length ? " text-error" : ""), title: bad.length ? `${x} takes ${text(t.accepts) || "nothing now"}, but this carries ${text(bad)}` : `may be sent ${t.accepts.length === 3 ? "any kind" : text(t.accepts)}` });
    const dot = el("span", undefined, { className: "inline-block size-2.5 rounded-full border border-base-100" });
    dot.style.background = t.dead ? "#dc2626" : t.in.length ? payloadColor(t.in) : t.anyKind ? "#475569" : "#f8fafc";
    if (!t.in.length && !t.anyKind && !t.dead) dot.style.borderColor = "#94a3b8";
    l.append(el("input", undefined, { type: "checkbox", className: "checkbox checkbox-xs", name, value: x, checked: chosen.includes(x) }), dot, x + (bad.length ? " (would not take it)" : ""));
    box.append(l);
  }
  // one chosen that is not there any more stays chosen, so that it is not lost unseen
  for (const x of chosen) {
    if (names.includes(x)) continue;
    const l = el("label", undefined, { className: "inline-flex items-center gap-1 text-sm text-error" });
    l.append(el("input", undefined, { type: "checkbox", className: "checkbox checkbox-xs", name, value: x, checked: true }), x + " (gone)");
    box.append(l);
  }
  if (!names.length) box.append(el("span", "No other node yet", { className: "text-xs text-base-content/60" }));
  return box;
}

const picked = (name) => [...document.querySelectorAll(`input[name="${name}"]:checked`)].map((c) => c.value);

// ruleRow is a rule of a filter in the sidebar.
function ruleRow(r = { match: {}, to: [] }) {
  const row = el("div", undefined, { className: "rule grid grid-cols-2 gap-1 rounded-box bg-base-200 p-2" });
  const m = r.match || {};
  row.append(
    el("input", undefined, { className: "input input-xs", placeholder: "name matches (RE2)", value: m.name || "", "data-k": "name" }),
    el("input", undefined, { className: "input input-xs", placeholder: "user matches", value: m.user || "", "data-k": "user" }),
    el("input", undefined, { className: "input input-xs", placeholder: "format matches (pdf, bdf …)", value: m.format || "", "data-k": "format" }),
    el("div", undefined, { className: "flex gap-1" }),
  );
  row.lastChild.append(el("input", undefined, { className: "input input-xs w-full", placeholder: "pages ≥", value: m.min_pages || "", inputMode: "numeric", "data-k": "min_pages" }), el("input", undefined, { className: "input input-xs w-full", placeholder: "pages ≤", value: m.max_pages || "", inputMode: "numeric", "data-k": "max_pages" }));
  const meta = el("input", undefined, { className: "input input-xs", placeholder: "metadata matches: key=RE2, key2=RE2", value: Object.entries(m.meta || {}).map(([k, v]) => k + "=" + v).join(", "), "data-k": "meta", title: "The sender's keys, the entry's, and $job, $entry, $sender, $file, $kind, $format, $pages …" });
  const kind = el("select", undefined, { className: "select select-xs", "data-k": "kind" });
  for (const [v, label] of Object.entries(KINDS_OF_PAYLOAD)) kind.append(el("option", v ? "payload: " + label : "payload: any kind", { value: v, selected: v === (m.kind || "") }));
  row.append(meta, kind);
  const to = el("input", undefined, { className: "input input-xs col-span-2", placeholder: "→ nodes, comma-separated", value: (r.to || []).join(", "), "data-k": "to" });
  to.setAttribute("list", "node-names");
  const drop = el("button", "remove the rule", { type: "button", className: "btn btn-xs btn-ghost col-span-2 justify-self-end text-error" });
  drop.addEventListener("click", () => row.remove());
  row.append(to, drop);
  return row;
}

// rulesOf reads the filter's rules from the sidebar.
function rulesOf() {
  const out = [];
  for (const row of document.querySelectorAll("#f-rules .rule")) {
    const v = (k) => row.querySelector(`[data-k="${k}"]`).value.trim();
    const match = {};
    if (v("name")) match.name = v("name");
    if (v("user")) match.user = v("user");
    if (v("format")) match.format = v("format");
    if (v("min_pages")) match.min_pages = Number(v("min_pages"));
    if (v("max_pages")) match.max_pages = Number(v("max_pages"));
    if (v("kind")) match.kind = v("kind");
    const meta = {};
    for (const pair of v("meta").split(",")) {
      const i = pair.indexOf("=");
      if (i > 0) meta[pair.slice(0, i).trim()] = pair.slice(i + 1).trim();
    }
    if (Object.keys(meta).length) match.meta = meta;
    const to = names(v("to"));
    if (!Object.keys(match).length && !to.length) continue;
    out.push({ match, to });
  }
  return out;
}

function selectInput(id, options, value) {
  const s = el("select", undefined, { id, className: "select select-sm w-full" });
  for (const [v, label] of options) s.append(el("option", label, { value: v, selected: v === value }));
  return s;
}

// metaRows are the metadata an entry adds to every job it takes, as
// rows of key and value in the sidebar.
function metaRow(k = "", v = "") {
  const row = el("div", undefined, { className: "meta grid grid-cols-[1fr_1fr_auto] gap-1" });
  row.append(el("input", undefined, { className: "input input-xs", placeholder: "key (site)", value: k, "data-k": "key" }), el("input", undefined, { className: "input input-xs", placeholder: "value (tokyo)", value: v, "data-k": "value" }));
  const drop = el("button", "×", { type: "button", className: "btn btn-xs btn-ghost text-error" });
  drop.addEventListener("click", () => row.remove());
  row.append(drop);
  return row;
}

function metaField(n) {
  const rows = el("div", undefined, { id: "f-meta", className: "flex flex-col gap-1" });
  for (const [k, v] of Object.entries(n.meta || {})) rows.append(metaRow(k, typeof v === "string" ? v : JSON.stringify(v)));
  const add = el("button", "+ a key", { type: "button", className: "btn btn-xs w-fit" });
  add.addEventListener("click", () => rows.append(metaRow()));
  const box = el("div", undefined, { className: "flex flex-col gap-1" });
  box.append(rows, add);
  return field("Metadata added", box, "Keys this entry adds to every job it takes (a site, a department): they win over the sender's (X-Meta-* headers, ?meta.<key>=), and filters route by them. Keys starting with $ are the manager's.");
}

// metaOf reads the metadata rows: a value that is JSON (a number, true,
// an object) is kept as that, else as text.
function metaOf() {
  const out = {};
  for (const row of document.querySelectorAll("#f-meta .meta")) {
    const k = row.querySelector('[data-k="key"]').value.trim();
    const v = row.querySelector('[data-k="value"]').value;
    if (!k) continue;
    try {
      out[k] = /^[-\d\[{"tfn]/.test(v.trim()) ? JSON.parse(v) : v;
    } catch {
      out[k] = v;
    }
  }
  return out;
}

// formatsField is the formats an entry takes, or a convert node converts.
function formatsField(n, label, what) {
  const fmts = el("div", undefined, { className: "flex flex-wrap gap-x-3 gap-y-1" });
  for (const f of formats) {
    const l = el("label", undefined, { className: "inline-flex items-center gap-1 text-sm" });
    l.append(el("input", undefined, { type: "checkbox", className: "checkbox checkbox-xs", name: "fmt", value: f, checked: (n.formats || []).includes(f) }), f);
    fmts.append(l);
  }
  return field(label, fmts, what);
}

// entryFields are what an entry takes and sends on: documents or data,
// the formats, and whether a document goes on as the .bdf made of it or
// as received.
function entryFields(form, n) {
  const accept = selectInput("f-accept", Object.entries(ACCEPTS), n.accept || "file");
  const output = selectInput("f-output", Object.entries(OUTPUTS), n.output || "bdf");
  const sync = () => (output.disabled = accept.value === "data");
  accept.addEventListener("change", sync);
  setTimeout(sync);
  form.append(field("Accepts", accept, "Data is for a template node after this entry: the edges leaving it carry data."));
  form.append(field("Sends on", output, "Of a document: converted here (printers and displays take .bdf), or as received for a convert node or an archive after this. Data always goes on as it is."));
  form.append(formatsField(n, "Formats taken", "None checked: every format the converter reads (and json, yaml for data). A file of another format is refused, and the refusal is in the audit log."));
}

// kindsLine tells what a node takes and sends on, as the API inferred
// it from the entries, what may be sent to it now (what the nodes after
// it take), and the edges from it that no longer connect.
function kindsLine(n) {
  const box = el("div", undefined, { className: "flex flex-col gap-1 text-xs text-base-content/60" });
  const text = (kinds) => kinds.map((k) => PAYLOADS[k] || k).join(" or ");
  const parts = [];
  if (n.takes) parts.push("takes " + (PAYLOADS[n.takes] || n.takes));
  const acc = n.accepts || [];
  if (n.accepts && acc.length < 3 && !(acc.length === 1 && acc[0] === n.takes)) parts.push(acc.length ? `may be sent ${text(acc)} only (what is after it takes)` : "nothing fits it now: the nodes after it take different kinds");
  if (n.emits && n.emits.length) parts.push("sends on " + text(n.emits));
  if (parts.length) box.append(el("p", parts.join(" · ")));
  for (const [to, kinds] of Object.entries(n.broken || {})) {
    const target = topo.find((x) => x.name === to);
    box.append(el("p", `⚠ → ${to} no longer connects: ${to} takes ${target ? PAYLOADS[target.takes] || target.takes : "something else"}, but ${text(kinds)} reaches it from here. Put a convert or a template node between, change what this node takes, or remove the edge; a job gets an error there meanwhile.`, { className: "text-error" }));
  }
  return box.childElementCount ? box : null;
}

const names = (s) => s.split(",").map((x) => x.trim()).filter(Boolean);

function hint(text) {
  return el("p", text, { className: "text-xs text-base-content/60" });
}

function section(title) {
  const d = el("div", undefined, { className: "mt-3 flex flex-col gap-2 border-t border-base-300 pt-3" });
  d.append(el("h4", title, { className: "font-semibold" }));
  return d;
}

// scopeBadge tells who may send to a node, and whose its jobs are.
function scopeBadge(n) {
  const scope = n.scope || "shared";
  const cls = { group: "badge badge-sm badge-soft badge-secondary whitespace-nowrap", personal: "badge badge-sm badge-soft badge-accent whitespace-nowrap" }[scope] || "badge badge-sm badge-ghost whitespace-nowrap";
  return el("span", scope + (n.owner ? ": " + n.owner : ""), { className: cls });
}

// parsePrinterURL splits a printer's URL into the sidebar's fields.
function parsePrinterURL(u) {
  if (!u) return { scheme: "ipp", host: "", port: "", path: "" };
  try {
    const x = new URL(u);
    return { scheme: x.protocol.replace(":", ""), host: x.hostname, port: x.port, path: x.pathname === "/" && !x.search ? "" : x.pathname + x.search };
  } catch {
    return { scheme: "ipp", host: u, port: "", path: "" };
  }
}

// composeURL is the printer's URL of the sidebar's fields.
function composeURL() {
  const scheme = $("f-scheme").value;
  const host = $("f-host").value.trim();
  const port = $("f-port").value.trim();
  let path = $("f-path").value.trim();
  if (!host) return "";
  if (!path && scheme.startsWith("ipp")) path = "/ipp/print";
  if (path && !path.startsWith("/")) path = "/" + path;
  return `${scheme}://${host}${port ? ":" + port : ""}${path}`;
}

// portField is the port of an entry node, when this process has ports.
function portField(n, what, required = false) {
  const none = me && me.ports === false;
  const f = field("Port (this process only)", input("f-port", n.port || "", { placeholder: required ? "6310" : "empty: the manager's own port only", inputMode: "numeric", required, disabled: none }), none ? "Not available here (serverless): applications call the manager's own API, and there is no IPP port." : what);
  return f;
}

// Looking: the sidebar of a node is its counts, a form to print there
// when the person may, and the jobs there they see.

function countsOf(name) {
  const c = (live && live.nodes[name]) || {};
  const out = [];
  if (c.open) out.push([String(c.open), "open"]);
  if (c.held) out.push([String(c.held), "held"]);
  if (c.pending) out.push([String(c.pending), "waiting"]);
  if (c.failed) out.push([String(c.failed), "failed"]);
  out.push([String(c.recent || 0), "in the last hour"]);
  return out;
}

function viewSidebar() {
  const aside = $("sidebar");
  const n = topo.find((x) => x.name === selected);
  if (!n) {
    const p = el("p", undefined, { className: "text-sm text-base-content/60" });
    p.append("Click a node to see the jobs sent there that you may see, and to print there when you may. ");
    if (live && Object.keys(live.nodes).length === 0) p.append("Nothing of yours is in the network right now.");
    aside.replaceChildren(p, el("a", "The manual: this page, and what each node does", { className: "link mt-2 block text-sm", href: "/help/topology#view", target: "_blank", rel: "noopener" }));
    return;
  }
  const head = el("div", undefined, { className: "flex flex-wrap items-center gap-2" });
  head.append(kindIcon(n.kind, 20), kindBadge(n.kind), el("h3", n.title && n.title !== n.name ? n.title : n.name, { className: "text-lg font-semibold" }));
  if (n.title && n.title !== n.name) head.append(el("span", n.name, { className: "font-mono text-sm text-base-content/60" }));
  const stats = el("div", undefined, { className: "stats stats-horizontal flex-wrap bg-base-200" });
  for (const [v, label] of countsOf(n.name)) {
    const s = el("div", undefined, { className: "stat px-3 py-2" });
    s.append(el("div", v, { className: "stat-value text-xl" }), el("div", label, { className: "stat-title text-xs" }));
    stats.append(s);
  }
  const parts = [head];
  if (n.description) parts.push(el("p", n.description, { className: "whitespace-pre-wrap text-sm text-base-content/80" }));
  parts.push(stats);
  const out = [];
  for (const to of n.to || []) out.push(to);
  for (const [i, r] of (n.rules || []).entries()) out.push((r.to || []).join(", ") + ` (rule ${i + 1})`);
  for (const to of n.otherwise || []) out.push(to + " (otherwise)");
  for (const m of n.members || []) out.push(m + (n.policy === "failover" ? " (first)" : " (in turn)"));
  if (n.fallback && n.fallback.printer) out.push(n.fallback.printer + " (fallback)");
  if (out.length) parts.push(el("p", "→ " + out.join(", "), { className: "font-mono text-xs text-base-content/60" }));
  const k = kindsLine(n);
  if (k) parts.push(k);
  const dest = dests.find((d) => d.name === n.name);
  if (dest) {
    const form = el("form", undefined, { className: "flex flex-col gap-2 border-t border-base-300 pt-3" });
    form.append(el("h4", "Print here", { className: "font-semibold" }));
    const file = el("input", undefined, { type: "file", className: "file-input file-input-sm w-full", required: true });
    const name = el("input", undefined, { className: "input input-sm w-full", placeholder: "job name (the file's)" });
    form.append(file, name, el("button", "Print", { type: "submit", className: "btn btn-sm btn-primary w-fit" }));
    form.addEventListener("submit", async (ev) => {
      ev.preventDefault();
      const f = file.files[0];
      if (!f) return;
      try {
        const j = await api("POST", `/api/destinations/${encodeURIComponent(n.name)}/jobs?name=${encodeURIComponent(name.value || f.name)}`, f, true);
        status(`Job #${j.number} is ${j.state}`);
        file.value = "";
        name.value = "";
        etag = "";
        await poll();
      } catch (e) {
        status(e.message, true);
      }
    });
    parts.push(form);
  } else {
    parts.push(el("p", "You may not print here.", { className: "text-xs text-base-content/60" }));
  }
  const mine = jobs.filter((j) => j.dest === n.name || (j.legs || []).some((l) => l.node === n.name));
  const sec = el("div", undefined, { className: "flex flex-col gap-2 border-t border-base-300 pt-3" });
  sec.append(el("h4", `Jobs here (${mine.length})`, { className: "font-semibold" }));
  if (!mine.length) sec.append(el("p", "None of yours.", { className: "text-xs text-base-content/60" }));
  const list = el("ul", undefined, { className: "flex flex-col gap-2" });
  for (const j of mine.slice(0, 20)) {
    const li = el("li", undefined, { className: "rounded-box bg-base-200 p-2 text-sm" });
    const row = el("div", undefined, { className: "flex flex-wrap items-center gap-2" });
    row.append(el("span", "#" + j.number, { className: "font-mono" }), el("span", j.name || "", { className: "font-semibold" }), stateBadge(j.state));
    li.append(row);
    li.append(el("div", `${sender(j)} · ${since(j.received)}`, { className: "text-xs text-base-content/60" }));
    li.append(el("div", legsOf(j), { className: "font-mono text-xs text-base-content/60 break-all", title: legsOf(j) }));
    const links = el("div", undefined, { className: "flex gap-3" });
    links.append(el("a", "Document", { className: "link text-xs", href: `/api/jobs/${j.id}/document` }));
    if (dest) {
      const re = el("button", "Reprint here", { type: "button", className: "link text-xs" });
      re.addEventListener("click", async () => {
        try {
          const again = await api("POST", `/api/jobs/${j.id}/reprint`, { to: n.name });
          status(`Job #${again.number} is ${again.state}`);
          etag = "";
          await poll();
        } catch (e) {
          status(e.message, true);
        }
      });
      links.append(re);
    }
    li.append(links);
    list.append(li);
  }
  sec.append(list);
  parts.push(sec);
  aside.replaceChildren(...parts);
  if (n.kind === "sink") keptSection(n.name).then((box) => aside.append(box)).catch((e) => status(e.message, true));
  if (n.kind === "export") listedSection(n.name).then((box) => aside.append(box)).catch((e) => status(e.message, true));
}

function renderSidebar() {
  if (!editable) return viewSidebar();
  const aside = $("sidebar");
  const n = draft || topo.find((x) => x.name === selected);
  if (!n) {
    aside.replaceChildren(
      el("p", "Click a node to edit it here, or add one with “+ Add a node”. Drag from a node's right port onto another node to connect them; a display's port sets its fallback printer. A node's title bar is colored by its role (source, control flow, filter, transform, sink), with an icon of its kind. The ports and the edges are colored by what flows: blue a .bdf document, green data (for a template node), amber any file as received (for a convert node), purple several kinds. A queue, a class, a filter or an inbox is typed by its first edge, on either side: in front of a printer it carries .bdf, and takes .bdf alone; one with no edge yet is hollow. While you drag an edge, the nodes that would not take it fade and the wire turns red over them. An edge drawn dashed and red no longer connects: a change of settings left it carrying what the node after it does not take.", { className: "text-sm text-base-content/60" }),
      el("a", "The manual: the editor, and every kind of node", { className: "link mt-2 block text-sm", href: "/help/topology#editor", target: "_blank", rel: "noopener" }),
    );
    return;
  }
  const kind = n.kind;
  const form = el("form", undefined, { className: "flex flex-col gap-1" });
  const head = el("div", undefined, { className: "flex flex-wrap items-center gap-2" });
  head.append(kindIcon(kind, 20), kindBadge(kind));
  if (draft) head.append(el("span", "new", { className: "badge badge-outline" }));
  else head.append(el("h3", n.name, { className: "font-mono text-lg font-semibold" }), scopeBadge(n));
  // the manual's section of the kind (the language is the browser's)
  head.append(el("a", "Help", { className: "link link-hover ml-auto text-xs text-base-content/60", href: "/help/nodes#" + kind, target: "_blank", rel: "noopener", title: "The manual: " + (KINDS[kind] || kind) + " nodes" }));
  form.append(head);
  if (draft) {
    const name = input("f-name", "", { placeholder: "intake, both, tokyo, review, wh-a, laser, audit", required: true, pattern: "[A-Za-z0-9][A-Za-z0-9_.\\-]{0,63}", autocomplete: "off" });
    name.addEventListener("input", () => canvas.draw());
    form.append(field("Name", name, "Letters, digits, . _ and -. It is the node's IPP destination (/ipp/print/<name>) and does not change later."));
  }
  form.append(field("Title", input("f-title", n.title || "", { placeholder: "shown in the node beside its kind" })));
  const note = el("textarea", undefined, { id: "f-description", className: "textarea textarea-sm w-full", rows: 3, value: n.description || "", placeholder: "what it is for: the invoices of the ERP, PDF only; printed on the 3rd floor …" });
  note.addEventListener("input", () => canvas.draw());
  form.append(field("Description", note, "A note shown in the node, under its [kind], as in a C4 diagram: the first lines on the canvas, all of it here."));
  if (!draft) {
    const k = kindsLine(n);
    if (k) form.append(k);
  }
  switch (kind) {
    case "endpoint": {
      form.append(field("Forwards to", picker("to", n.to || [], n.name), "Every node chosen gets the job. The port on the canvas adds one too."));
      entryFields(form, n);
      form.append(metaField(n));
      form.append(portField(n, "Applications may POST to http://<this host>:<port>/ as well as to " + n.in + "."));
      const appBoxes = el("div", undefined, { className: "flex flex-wrap gap-x-3 gap-y-1" });
      for (const a of apps) {
        const l = el("label", undefined, { className: "inline-flex items-center gap-1 text-sm" });
        l.append(el("input", undefined, { type: "checkbox", className: "checkbox checkbox-xs", name: "app", value: a.id, checked: (n.apps || []).includes(a.id) }), a.id + (a.title ? " · " + a.title : ""));
        appBoxes.append(l);
      }
      if (!apps.length) appBoxes.append(el("a", "No apps yet: register one", { className: "link text-sm", href: "/admin/apps" }));
      form.append(field("Apps allowed", appBoxes, "Apps send here with their token (POST /oauth/token) or IPP Basic, within their scopes; the endpoint's own key works as well."));
      break;
    }
    case "cups":
      form.append(field("Forwards to", picker("to", n.to || [], n.name), "Every node chosen gets what computers print here: a printer makes it a relay, a sink shows the pages, a filter keeps only what fits. The port on the canvas adds one too."));
      entryFields(form, n);
      form.append(metaField(n));
      form.append(portField(n, `Computers print to ipp://<this host>:<port>/printers/${n.name || "<name>"}, the one printer of this CUPS server; lpstat -h <host>:<port> -p lists it.`, true));
      form.append(field("TLS certificate (file)", input("f-tls-cert", n.tls_cert || "", { placeholder: "/etc/bdf/cups.crt; empty: plain ipp" })));
      form.append(field("TLS key (file)", input("f-tls-key", n.tls_key || "", { placeholder: "/etc/bdf/cups.key" })));
      break;
    case "sink":
      form.append(field("Latest", input("f-latest", n.latest || "", { placeholder: "10", inputMode: "numeric" }), "How many of the latest documents it keeps; an older one goes when a newer comes. Their first pages show here and on the canvas."));
      break;
    case "export": {
      form.append(field("Owner's key", input("f-by", n.by || "", { placeholder: "customer_id (a metadata key); empty: a person reads what they could see of the job", autocomplete: "off" }), "The metadata key whose value says whose a document is: a person reads the documents whose value is them, under the subject below. A backend reads every one."));
      form.append(field("Subject", selectInput("f-subject", Object.entries(SUBJECTS), n.subject || "email"), "What of a person is compared with the owner's value."));
      form.append(field("Fields", input("f-fields", (n.fields || []).join(", "), { placeholder: "invoice_no, doc_type (metadata keys, comma-separated)" }), "Metadata copied into the entries listed, for the reader to show or pick by."));
      form.append(field("Entries", input("f-entries", n.entries || "", { placeholder: "1000", inputMode: "numeric" }), "How many of the latest entries are listed for one owner, and for everyone; an older entry is still read by its job's number."));
      const appBoxes = el("div", undefined, { className: "flex flex-wrap gap-x-3 gap-y-1" });
      for (const a of apps) {
        if (a.public) continue;
        const l = el("label", undefined, { className: "inline-flex items-center gap-1 text-sm" });
        l.append(el("input", undefined, { type: "checkbox", className: "checkbox checkbox-xs", name: "app", value: a.id, checked: (n.apps || []).includes(a.id) }), a.id + (a.title ? " · " + a.title : ""));
        appBoxes.append(l);
      }
      if (!appBoxes.childElementCount) appBoxes.append(el("a", "No backend apps yet: register one", { className: "link text-sm", href: "/admin/apps" }));
      form.append(field("Apps allowed", appBoxes, "Backends that read every entry with their token (the scope export). A frontend (a public client) reads as the person using it, with a token of theirs: it needs no listing here."));
      break;
    }
    case "queue":
      form.append(field("Forwards to", picker("to", n.to || [], n.name), "Every node chosen gets the job: a fan-out when there are several."));
      break;
    case "template": {
      form.append(field("Forwards to", picker("to", n.to || [], n.name), "The form filled (.bdf) goes to every node chosen."));
      const ids = n.templates || [];
      const options = [["", "— named by the metadata below —"], ...ids.map((id) => [id, id])];
      if (n.template && !ids.includes(n.template)) options.push([n.template, n.template + " (not uploaded yet)"]);
      form.append(field("Template", selectInput("f-template", options, n.template || ""), draft ? "Add the node first, then upload its templates below and choose one here." : "One of the templates uploaded below, filled with the data that reaches the node."));
      form.append(field("Or the metadata key naming it", input("f-template-from", n.template_from || "", { placeholder: "doc_type: the sender's X-Meta-Doc-Type names the template" }), "A job whose metadata lacks the key, or names a template that is not here, fails here (no retry)."));
      form.append(field("Name of a list", input("f-list-name", n.list_name || "", { placeholder: "the file's name without its extension (orders for orders.csv)" }), "Data that is a list of records (a CSV file, a sheet, a JSON array) is the value of this name for the template; an object is the data itself."));
      const strict = el("label", undefined, { className: "inline-flex items-center gap-2 text-sm" });
      strict.append(el("input", undefined, { id: "f-strict", type: "checkbox", className: "checkbox checkbox-sm", checked: !!n.strict }), "Strict: refuse data that lacks a name the template looks up (else the tag is left empty, with a warning)");
      form.append(strict);
      break;
    }
    case "convert":
      form.append(field("Forwards to", picker("to", n.to || [], n.name), "The .bdf document made goes to every node chosen; a .bdf that reaches the node goes on as it is."));
      form.append(formatsField(n, "Formats converted", "None checked: every format the converter reads. A file of another format fails here (no retry); JSON and YAML have no document of their own: fill a template with them."));
      break;
    case "filter": {
      const rules = el("div", undefined, { id: "f-rules", className: "flex flex-col gap-2" });
      for (const r of n.rules || []) rules.append(ruleRow(r));
      const add = el("button", "+ a rule", { type: "button", className: "btn btn-xs w-fit" });
      add.addEventListener("click", () => rules.append(ruleRow()));
      const box = el("div", undefined, { className: "flex flex-col gap-2" });
      box.append(rules, add);
      form.append(field("Rules, in order", box, "The first rule that fits the job sends it to its nodes. name, user (the sender), format and the values of the metadata are regular expressions; pages are bounds; the kind of payload is .bdf, data or any other file. What reaches the filter goes on as it is."));
      form.append(field("Otherwise", picker("otherwise", n.otherwise || [], n.name), "Where a job no rule fits goes; none: it ends here (dropped)."));
      break;
    }
    case "class":
      form.append(field("Members", input("f-members", (n.members || []).join(", "), { placeholder: "mfp-3f, mfp-4f", list: "node-names", required: true })));
      form.append(field("Policy", selectInput("f-policy", [["round-robin", "in turn (round-robin)"], ["failover", "the first member (failover)"]], n.policy || "round-robin"), "A printer that does not take the job is skipped for the next member."));
      break;
    case "inbox":
      form.append(field("A person releases to", picker("to", n.to || [], n.name), "The nodes a person may release to; none chosen: any node."));
      form.append(field("Keep", input("f-keep", n.keep || "", { placeholder: "14d" }), "How long released and discarded items stay shown."));
      break;
    case "archive":
      form.append(field("Forwards to", picker("to", n.to || [], n.name), "Once written, the job goes on to every node chosen: an export after the archive lists the document for applications and people to fetch. The port on the canvas adds one too."));
      form.append(field("Store", input("f-store", n.store || "", { placeholder: "s3://bucket?region=… or bbolt:path", required: true }), "A documents URL: bbolt, s3, r2 or gcs. Each job's .bdf and record are written once; the manager never changes or deletes them."));
      form.append(field("Collection", input("f-collection", n.collection || "", { placeholder: "archive-" + (n.name || "<name>") })));
      break;
    case "display": {
      const printers = topo.filter((x) => x.kind === "printer").map((x) => [x.name, x.name + (x.title ? " · " + x.title : "")]);
      form.append(field("Fallback printer", selectInput("f-fallback", [["", "none"], ...printers], n.fallback ? n.fallback.printer : ""), "Prints the slips nobody opened in time, and every unopened slip while no device is connected. The port on the canvas sets it too."));
      form.append(field("Print after", input("f-after", n.fallback && n.fallback.after ? n.fallback.after : "", { placeholder: "10m; empty: only while no device is connected" })));
      form.append(field("Keep", input("f-keep", n.keep || "", { placeholder: "14d" }), "How long slips done or canceled stay on the devices."));
      form.append(field("Offline after", input("f-offline", n.offline_after || "", { placeholder: "2m" }), "How long a device may be silent and still count as connected."));
      break;
    }
    case "printer": {
      const u = parsePrinterURL(n.url);
      const row = el("div", undefined, { className: "grid grid-cols-[6rem_1fr_5rem] gap-1" });
      row.append(selectInput("f-scheme", SCHEMES.map((s) => [s, s + "://"]), SCHEMES.includes(u.scheme) ? u.scheme : "ipp"), input("f-host", u.host, { placeholder: "host or address", required: true }), input("f-port", u.port, { placeholder: "port", inputMode: "numeric" }));
      form.append(field("Protocol and address", row, "ipp and ipps are IPP printers; socket is port 9100 (raw), lpd a line printer daemon."));
      form.append(field("Path", input("f-path", u.path, { placeholder: "/ipp/print (ipp); a queue name (lpd)" })));
      const preview = el("p", "", { className: "font-mono text-xs text-base-content/60 break-all" });
      const update = () => (preview.textContent = composeURL() || "—");
      form.append(preview);
      setTimeout(update);
      for (const id of ["f-scheme", "f-host", "f-port", "f-path"]) form.querySelector("#" + id).addEventListener("input", update);
      form.append(field("Route", selectInput("f-route", [["direct", "direct: the manager sends to it"], ["branch", "branch: a bdf printd of the site fetches (not yet)"]], n.route || "direct")));
      form.append(field("Preview (a fake printer)", input("f-preview", n.preview || "", { placeholder: "http://127.0.0.1:8650/api/printers/eve" }), "The printer on a tools/fakeprinter farm's API: what it received is shown below. For development and demonstrations."));
      break;
    }
  }
  const scopes = kind === "export"
    ? [["shared", "shared: anyone signed in may read their own"], ["group", "group: its members may read their own"], ["personal", "personal: one person"]]
    : [["shared", "shared: anyone signed in may send"], ["group", "group: its members send and see the jobs"], ["personal", "personal: one person"]];
  form.append(field("Scope", selectInput("f-scope", scopes, n.scope || "shared")));
  form.append(field("Owner", input("f-owner", n.owner || "", { placeholder: "a group (a value of the groups claim), or a person's email" })));
  const actions = el("div", undefined, { className: "mt-2 flex flex-wrap gap-2" });
  actions.append(el("button", draft ? "Add" : "Save", { type: "submit", className: "btn btn-sm btn-primary" }));
  if (draft) {
    const cancel = el("button", "Cancel", { type: "button", className: "btn btn-sm btn-ghost" });
    cancel.addEventListener("click", () => {
      draft = null;
      canvas.setDraft(null, null);
      renderSidebar();
    });
    actions.append(cancel);
  } else {
    const del = el("button", "Delete", { type: "button", className: "btn btn-sm btn-ghost text-error" });
    del.addEventListener("click", () => remove(n.name));
    actions.append(del);
  }
  form.append(actions);
  form.addEventListener("submit", save);
  aside.replaceChildren(form);
  if (!draft) extras(n, aside);
}

async function save(ev) {
  ev.preventDefault();
  const n = draft || topo.find((x) => x.name === selected);
  if (!n) return;
  const name = draft ? $("f-name").value.trim() : n.name;
  const body = { kind: n.kind, title: $("f-title").value.trim(), description: $("f-description").value.trim(), scope: $("f-scope").value };
  if ($("f-owner").value.trim()) body.owner = $("f-owner").value.trim();
  const dur = (id) => ($(id).value.trim() ? $(id).value.trim() : undefined);
  switch (n.kind) {
    case "endpoint":
      body.to = picked("to");
      body.accept = $("f-accept").value;
      body.output = $("f-output").value;
      body.formats = [...document.querySelectorAll('input[name="fmt"]:checked')].map((c) => c.value);
      body.meta = metaOf();
      body.apps = [...document.querySelectorAll('input[name="app"]:checked')].map((c) => c.value);
      if (dur("f-port")) body.port = Number(dur("f-port"));
      break;
    case "cups":
      body.to = picked("to");
      body.accept = $("f-accept").value;
      body.output = $("f-output").value;
      body.formats = [...document.querySelectorAll('input[name="fmt"]:checked')].map((c) => c.value);
      body.meta = metaOf();
      body.port = Number($("f-port").value);
      if (dur("f-tls-cert")) body.tls_cert = dur("f-tls-cert");
      if (dur("f-tls-key")) body.tls_key = dur("f-tls-key");
      break;
    case "template":
      body.to = picked("to");
      if ($("f-template").value) body.template = $("f-template").value;
      if (dur("f-template-from")) body.template_from = dur("f-template-from");
      if (dur("f-list-name")) body.list_name = dur("f-list-name");
      if ($("f-strict").checked) body.strict = true;
      break;
    case "convert":
      body.to = picked("to");
      body.formats = [...document.querySelectorAll('input[name="fmt"]:checked')].map((c) => c.value);
      break;
    case "sink":
      if (dur("f-latest")) body.latest = Number(dur("f-latest"));
      break;
    case "export":
      if (dur("f-by")) body.by = dur("f-by");
      body.subject = $("f-subject").value;
      body.fields = names($("f-fields").value);
      if (dur("f-entries")) body.entries = Number(dur("f-entries"));
      body.apps = [...document.querySelectorAll('input[name="app"]:checked')].map((c) => c.value);
      break;
    case "queue":
      body.to = picked("to");
      break;
    case "filter":
      body.rules = rulesOf();
      body.otherwise = picked("otherwise");
      break;
    case "class":
      body.members = names($("f-members").value);
      body.policy = $("f-policy").value;
      break;
    case "inbox":
      body.to = picked("to");
      body.keep = dur("f-keep");
      break;
    case "archive":
      body.to = picked("to");
      body.store = $("f-store").value.trim();
      if ($("f-collection").value.trim()) body.collection = $("f-collection").value.trim();
      break;
    case "display":
      if ($("f-fallback").value) {
        body.fallback = { printer: $("f-fallback").value };
        if (dur("f-after")) body.fallback.after = dur("f-after");
      }
      body.keep = dur("f-keep");
      body.offline_after = dur("f-offline");
      break;
    case "printer":
      body.url = composeURL();
      body.route = $("f-route").value;
      if (dur("f-preview")) body.preview = dur("f-preview");
      break;
  }
  try {
    const was = new Set(brokenEdges(topo));
    const res = await api("PUT", "/api/nodes/" + encodeURIComponent(name), body);
    if (draft) {
      await api("PUT", "/api/layout", { positions: { [name]: canvas.draftPos() } });
      draft = null;
      selected = name;
      canvas.setSelected(name);
    }
    status(`${name} saved`);
    if (res.key) showKey(name, res.key);
    await load();
    // a change of settings may leave an edge carrying what the node
    // after it does not take: said here, and drawn dashed
    const broke = brokenEdges(topo).filter((e) => !was.has(e));
    if (broke.length) status(`Saved, but ${broke.join("; ")} no longer connect${broke.length === 1 ? "s" : ""}: the node after does not take what reaches it (the dashed edge). Put a convert or a template node between, or remove the edge.`, true);
    $("sidebar").scrollTop = 0;
  } catch (e) {
    status(e.message, true);
  }
}

async function remove(name) {
  if (!confirm(`Delete ${name}? A node another one sends to stays. A template node's templates go with it; a display's slips and devices stay until it is added again.`)) return;
  try {
    await api("DELETE", "/api/nodes/" + encodeURIComponent(name));
    selected = null;
    status(`${name} deleted`);
    await load();
  } catch (e) {
    status(e.message, true);
  }
}

// extras are what is at a saved node: its jobs, an endpoint's entry, key
// and templates, a display's devices.
function extras(n, aside) {
  const links = el("div", undefined, { className: "mt-3 flex flex-col gap-1 border-t border-base-300 pt-3 text-sm" });
  // an export is read from, never sent to: no IPP destination
  if (n.kind !== "export") links.append(el("span", "IPP: " + n.ipp, { className: "font-mono text-xs text-base-content/60 break-all" }));
  const row = el("div", undefined, { className: "flex flex-wrap gap-3" });
  if (n.kind !== "export") row.append(el("a", "Jobs sent here", { className: "link", href: "/admin/jobs?dest=" + encodeURIComponent(n.name) }));
  if (n.kind === "inbox") row.append(el("a", `${n.held || 0} held: open the inbox`, { className: "link", href: "/admin/inboxes" }));
  if (n.pending) row.append(el("a", `${n.pending} deliveries waiting`, { className: "link", href: "/admin/deliveries" }));
  links.append(row);
  aside.append(links);
  if (n.kind === "endpoint") {
    const entry = section("Entry");
    entry.append(hint(n.accept === "data" ? "Applications POST the data (JSON, YAML, CSV, an Excel workbook or a Parquet file) as the body, with the key as a bearer token; ?name= names the job, X-Meta-* headers and ?meta.<key>= give metadata." : "Applications POST a document as the body, with the key as a bearer token; ?name= names the job, X-Meta-* headers and ?meta.<key>= give metadata."));
    entry.append(el("pre", `POST ${location.origin}${n.in}`, { className: "overflow-x-auto rounded-box bg-base-200 p-2 font-mono text-xs break-all whitespace-pre-wrap" }));
    const key = section("Key");
    key.append(hint("Shown once, when the endpoint was made or the key rotated; the manager keeps only its hash. Rotating it makes a new one and the old one stops working."));
    key.append(button("Rotate key", async () => showKey(n.name, (await api("POST", `/api/nodes/${n.name}/key`)).key), "btn btn-sm btn-outline w-fit"));
    aside.append(entry, key);
  }
  if (n.kind === "template") {
    const tpls = section("Templates");
    aside.append(tpls);
    templatesSection(n, tpls).catch((e) => status(e.message, true));
  }
  if (n.kind === "display") {
    const dev = section("Devices");
    aside.append(dev);
    devicesSection(n, dev).catch((e) => status(e.message, true));
  }
  if (n.kind === "cups" && n.port) {
    const sec = section("Clients");
    sec.append(hint(`As on a CUPS server with one printer, ${n.name}: what is printed there goes on to ${(n.to || []).join(", ") || "nothing yet"}. Add the printer on a computer with the URL, or list it:`));
    sec.append(el("pre", `ipp://${location.hostname}:${n.port}/printers/${n.name}
lpstat -h ${location.hostname}:${n.port} -p`, { className: "overflow-x-auto rounded-box bg-base-200 p-2 font-mono text-xs whitespace-pre-wrap" }));
    aside.append(sec);
  }
  if (n.kind === "sink") {
    keptSection(n.name).then((sec) => aside.append(sec)).catch((e) => status(e.message, true));
  }
  if (n.kind === "export") {
    const sec = section("Readers");
    sec.append(hint("Applications and people fetch what is listed here: the entries newest first (?owner= one owner's, ?limit=), one entry, the archive's record of it, and the document itself. A backend sends its token (POST /oauth/token, client credentials, the scope export); a frontend sends the token it got for the person (a public client, /oauth/authorize with PKCE); the pages use their session."));
    sec.append(el("pre", `GET ${location.origin}/out/${n.name}
GET ${location.origin}/out/${n.name}/<job>
GET ${location.origin}/out/${n.name}/<job>/record
GET ${location.origin}/out/${n.name}/<job>/document`, { className: "overflow-x-auto rounded-box bg-base-200 p-2 font-mono text-xs whitespace-pre-wrap" }));
    aside.append(sec);
    listedSection(n.name).then((box) => aside.append(box)).catch((e) => status(e.message, true));
  }
  if (n.kind === "printer" && n.preview) {
    const sec = section("Preview");
    aside.append(sec);
    previewSection(n, sec).catch((e) => status(e.message, true));
  }
}

// previewSection shows what the printer's fake printer received: the
// latest jobs, and the pages read back of the newest.
async function previewSection(n, box) {
  box.append(hint("What the fake printer on the farm received, as it read it back."));
  const jobs = await api("GET", `/api/nodes/${n.name}/preview`);
  if (!jobs.length) {
    box.append(el("p", "Nothing received yet.", { className: "text-xs text-base-content/60" }));
    return;
  }
  const latest = jobs[0];
  box.append(el("p", `#${latest.id} ${latest.name || ""} · ${latest.state || ""} · ${since(latest.created)}`, { className: "text-sm font-semibold" }));
  const pages = el("div", undefined, { className: "flex flex-col gap-2" });
  for (let k = 1; k <= (latest.pages || 0); k++) {
    pages.append(el("img", undefined, { src: `/api/nodes/${n.name}/preview/${latest.id}/pages/${k}.png`, alt: `page ${k}`, className: "w-full rounded bg-white shadow-sm", loading: "lazy" }));
  }
  if (!latest.pages) pages.append(el("p", "The farm could not read it back as pages.", { className: "text-xs text-base-content/60" }));
  box.append(pages);
  if (jobs.length > 1) {
    const list = el("ul", undefined, { className: "text-xs text-base-content/60" });
    for (const j of jobs.slice(1)) list.append(el("li", `#${j.id} ${j.name || ""} · ${since(j.created)}`));
    box.append(el("p", "Before it:", { className: "text-xs font-semibold" }), list);
  }
}

async function templatesSection(n, box) {
  box.append(hint("This node fills one of them (the one chosen above, or the one the metadata names) with the data that reaches it from a data endpoint. A compiled template (.bdft), or Excel, Word, PowerPoint, a PDF form, pdfme or text."));
  const ts = await api("GET", `/api/nodes/${n.name}/templates`);
  const table = el("table", undefined, { className: "table table-xs" });
  const tbody = el("tbody");
  for (const t of ts) {
    const tr = el("tr");
    const fields = (t.fields || []).map((f) => f.name).join(", ");
    tr.append(el("td", t.id, { className: "font-mono font-semibold" }), el("td", t.format), el("td", `${(t.fields || []).length} fields`, { title: fields }), el("td", `${Math.max(1, Math.round(t.size / 1024))} KiB`), el("td"));
    tr.lastChild.append(el("a", "File", { className: "btn btn-xs btn-ghost", href: `/api/nodes/${n.name}/templates/${t.id}/file` }), button("Delete", () => api("DELETE", `/api/nodes/${n.name}/templates/${t.id}`), "btn btn-xs btn-ghost text-error"));
    tbody.append(tr);
  }
  if (!ts.length) {
    const tr = el("tr");
    tr.append(el("td", "No templates yet", { colSpan: 5, className: "text-base-content/50" }));
    tbody.append(tr);
  }
  table.append(tbody);
  const wrap = el("div", undefined, { className: "overflow-x-auto" });
  wrap.append(table);
  box.append(wrap);
  const form = el("form", undefined, { className: "flex flex-wrap gap-1" });
  const id = el("input", undefined, { className: "input input-sm w-28", placeholder: "id (invoice)", required: true, pattern: "[A-Za-z0-9][A-Za-z0-9_\\-]{0,63}" });
  const file = el("input", undefined, { type: "file", className: "file-input file-input-sm min-w-0 flex-1", required: true });
  form.append(id, file, el("button", "Upload", { type: "submit", className: "btn btn-sm" }));
  form.addEventListener("submit", async (ev) => {
    ev.preventDefault();
    const f = file.files[0];
    if (!f) return;
    try {
      await api("PUT", `/api/nodes/${n.name}/templates/${encodeURIComponent(id.value)}?name=${encodeURIComponent(f.name)}`, f, true);
      status(`${id.value} uploaded`);
      changed();
    } catch (e) {
      status(e.message, true);
    }
  });
  box.append(form);
}

async function devicesSection(n, box) {
  box.append(hint(`${n.open || 0} open slip${n.open === 1 ? "" : "s"} on the board. Only the devices enrolled here see its documents; a device keeps a key of its own.`));
  const ds = await api("GET", `/api/nodes/${n.name}/devices`);
  const table = el("table", undefined, { className: "table table-xs" });
  const tbody = el("tbody");
  for (const d of ds) {
    const tr = el("tr");
    tr.append(el("td", d.name, { className: "font-semibold" }), el("td", "seen " + since(d.last_seen), { className: "text-base-content/60" }), el("td"));
    tr.lastChild.append(button("Remove", () => api("DELETE", "/api/devices/" + d.id), "btn btn-xs btn-ghost text-error"));
    tbody.append(tr);
  }
  if (!ds.length) {
    const tr = el("tr");
    tr.append(el("td", "No devices yet", { colSpan: 3, className: "text-base-content/50" }));
    tbody.append(tr);
  }
  table.append(tbody);
  box.append(table);
  const add = el("button", "Add a device (QR code)", { type: "button", className: "btn btn-sm btn-primary w-fit" });
  add.addEventListener("click", async () => {
    try {
      const e = await api("POST", `/api/nodes/${n.name}/enrollments`);
      $("qr").innerHTML = e.qr; // SVG the manager drew
      $("enroll-code").textContent = e.code.replace(/(.{5})/, "$1-");
      $("enroll-url").textContent = e.url;
      $("enroll-url").href = e.url;
      $("enroll-expires").textContent = `Valid once, until ${since(e.expires)}`;
      $("enroll-modal").showModal();
    } catch (err) {
      status(err.message, true);
    }
  });
  box.append(add);
}

function showKey(name, key) {
  $("key-name").textContent = name;
  $("key-value").textContent = key;
  $("key-usage").textContent = `curl -X POST -H "Authorization: Bearer ${key}" --data-binary @document.pdf "${location.origin}/in/${name}?name=document"`;
  $("key-copy").textContent = "Copy";
  $("key-modal").showModal();
}

// The Copy button puts the key on the clipboard. Without the clipboard
// API (a page that is not a secure context) the key is selected, so
// that the keyboard's copy takes it.
if (editable) $("key-copy").addEventListener("click", async () => {
  const key = $("key-value").textContent;
  let copied = false;
  try {
    await navigator.clipboard.writeText(key);
    copied = true;
  } catch {
    const whole = document.createRange();
    whole.selectNodeContents($("key-value"));
    const sel = getSelection();
    sel.removeAllRanges();
    sel.addRange(whole);
    try {
      copied = document.execCommand("copy");
    } catch {}
  }
  $("key-copy").textContent = copied ? "Copied" : "Selected: press Ctrl+C or ⌘C";
  setTimeout(() => ($("key-copy").textContent = "Copy"), 2500);
});


if (editable) document.addEventListener("keydown", (ev) => {
  const typing = ["INPUT", "TEXTAREA", "SELECT"].includes(document.activeElement?.tagName);
  if (ev.key === "Escape" && draft) {
    draft = null;
    canvas.setDraft(null, null);
    renderSidebar();
  } else if (ev.key === "Delete" && selected && !typing) {
    remove(selected);
  }
});

// a link from another page names a node: /topology#<name>
if (location.hash.length > 1) {
  selected = decodeURIComponent(location.hash.slice(1));
  canvas.setSelected(selected);
}
onChange(load);
load().then(poll);
setInterval(() => {
  if (!document.hidden) poll();
}, POLL_MS);
if (!editable) {
  // looking, the nodes are read again now and then (editing, Refresh)
  setInterval(() => {
    if (!document.hidden) load();
  }, 60000);
}
document.addEventListener("visibilitychange", () => {
  if (!document.hidden) poll();
});
