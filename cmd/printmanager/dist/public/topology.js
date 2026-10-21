// The editor of the topology. Nodes of every kind are placed on the
// canvas (canvas.js: panned and zoomed, with the activity running along
// the edges), connected by dragging from a node's output port onto
// another node, and set up in the sidebar: the fields of the node's kind,
// and what is at it (an endpoint's templates and key, a display's
// devices and enrollment codes). Everything goes through the API; the
// server checks a node when it is put (its edges, its name, no cycle).

import { $, api, button, changed, el, keptSection, kindBadge, KINDS, onChange, since, status } from "./admin.js";
import { createCanvas } from "./canvas.js";

const SCHEMES = ["ipp", "ipps", "socket", "lpd"];
const ACCEPTS = { both: "documents, and templates to fill", file: "documents only", template: "templates to fill only" };
const POLL_MS = 4000;

let topo = []; // the nodes, as /api/topology lists them
let formats = []; // the formats an endpoint may take
let me = null; // /api/me: whether this process has ports of its own
let apps = []; // the apps, for an endpoint's allow-list
let selected = null; // the node in the sidebar
let draft = null; // a node being added: {kind}
let fitted = false;
let etag = ""; // the activity's revision

const canvas = createCanvas($("canvas"), {
  editable: true,
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
  onRemoveEdge: removeEdge,
  draftLabel: () => $("f-name")?.value || "",
});

async function load() {
  try {
    const [t, l] = await Promise.all([api("GET", "/api/topology"), api("GET", "/api/layout")]);
    topo = t;
    if (!formats.length) formats = await api("GET", "/api/formats");
    if (!me) me = await api("GET", "/api/me");
    apps = await api("GET", "/api/apps");
    $("node-names").replaceChildren(...topo.map((n) => el("option", undefined, { value: n.name })));
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

// poll reads the activity (every job, for an administrator) when its
// revision changed.
async function poll() {
  try {
    const res = await fetch("/api/topology/live", { headers: etag ? { "If-None-Match": etag } : {} });
    if (res.status === 200) {
      etag = res.headers.get("ETag") || "";
      canvas.setLive(await res.json());
    }
  } catch {}
}

$("zoom-in").addEventListener("click", () => canvas.zoomBy(1.25));
$("zoom-out").addEventListener("click", () => canvas.zoomBy(0.8));
$("zoom-fit").addEventListener("click", () => canvas.fit());
$("refresh").addEventListener("click", load);

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
$("palette").addEventListener("click", (ev) => {
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
// checkboxes' name; except the node itself and the entries).
function picker(name, chosen, except) {
  const box = el("div", undefined, { className: "flex flex-wrap gap-x-3 gap-y-1" });
  const names = topo.filter((x) => x.name !== except && x.kind !== "cups").map((x) => x.name);
  for (const x of names) {
    const l = el("label", undefined, { className: "inline-flex items-center gap-1 text-sm" });
    l.append(el("input", undefined, { type: "checkbox", className: "checkbox checkbox-xs", name, value: x, checked: chosen.includes(x) }), x);
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

function renderSidebar() {
  const aside = $("sidebar");
  const n = draft || topo.find((x) => x.name === selected);
  if (!n) {
    aside.replaceChildren(el("p", "Click a node to edit it here, or add one with “+ Add a node”. Drag from a node's right port onto another node to connect them; a display's port sets its fallback printer.", { className: "text-sm text-base-content/60" }));
    return;
  }
  const kind = n.kind;
  const form = el("form", undefined, { className: "flex flex-col gap-1" });
  const head = el("div", undefined, { className: "flex flex-wrap items-center gap-2" });
  head.append(kindBadge(kind));
  if (draft) head.append(el("span", "new", { className: "badge badge-outline" }));
  else head.append(el("h3", n.name, { className: "font-mono text-lg font-semibold" }), scopeBadge(n));
  form.append(head);
  if (draft) {
    const name = input("f-name", "", { placeholder: "intake, both, tokyo, review, wh-a, laser, audit", required: true, pattern: "[A-Za-z0-9][A-Za-z0-9_.\\-]{0,63}", autocomplete: "off" });
    name.addEventListener("input", () => canvas.draw());
    form.append(field("Name", name, "Letters, digits, . _ and -. It is the node's IPP destination (/ipp/print/<name>) and does not change later."));
  }
  form.append(field("Title", input("f-title", n.title || "", { placeholder: "shown instead of the name" })));
  switch (kind) {
    case "endpoint": {
      form.append(field("Forwards to", picker("to", n.to || [], n.name), "Every node chosen gets the job. The port on the canvas adds one too."));
      form.append(field("Accepts", selectInput("f-accept", Object.entries(ACCEPTS), n.accept || "both")));
      const fmts = el("div", undefined, { className: "flex flex-wrap gap-x-3 gap-y-1" });
      for (const f of formats) {
        const l = el("label", undefined, { className: "inline-flex items-center gap-1 text-sm" });
        l.append(el("input", undefined, { type: "checkbox", className: "checkbox checkbox-xs", name: "fmt", value: f, checked: (n.formats || []).includes(f) }), f);
        fmts.append(l);
      }
      form.append(field("Formats taken", fmts, "None checked: every format the converter reads. A document of another format is refused, and the refusal is in the audit log."));
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
      form.append(portField(n, `Computers print to ipp://<this host>:<port>/printers/${n.name || "<name>"}, the one printer of this CUPS server; lpstat -h <host>:<port> -p lists it.`, true));
      form.append(field("TLS certificate (file)", input("f-tls-cert", n.tls_cert || "", { placeholder: "/etc/bdf/cups.crt; empty: plain ipp" })));
      form.append(field("TLS key (file)", input("f-tls-key", n.tls_key || "", { placeholder: "/etc/bdf/cups.key" })));
      break;
    case "sink":
      form.append(field("Latest", input("f-latest", n.latest || "", { placeholder: "10", inputMode: "numeric" }), "How many of the latest documents it keeps; an older one goes when a newer comes. Their first pages show here and on the canvas."));
      break;
    case "queue":
      form.append(field("Forwards to", picker("to", n.to || [], n.name), "Every node chosen gets the job: a fan-out when there are several."));
      break;
    case "filter": {
      const rules = el("div", undefined, { id: "f-rules", className: "flex flex-col gap-2" });
      for (const r of n.rules || []) rules.append(ruleRow(r));
      const add = el("button", "+ a rule", { type: "button", className: "btn btn-xs w-fit" });
      add.addEventListener("click", () => rules.append(ruleRow()));
      const box = el("div", undefined, { className: "flex flex-col gap-2" });
      box.append(rules, add);
      form.append(field("Rules, in order", box, "The first rule that fits the job sends it to its nodes. name, user (the sender) and format are regular expressions; pages are bounds."));
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
  form.append(field("Scope", selectInput("f-scope", [["shared", "shared: anyone signed in may send"], ["group", "group: its members send and see the jobs"], ["personal", "personal: one person"]], n.scope || "shared")));
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
  const body = { kind: n.kind, title: $("f-title").value.trim(), scope: $("f-scope").value };
  if ($("f-owner").value.trim()) body.owner = $("f-owner").value.trim();
  const dur = (id) => ($(id).value.trim() ? $(id).value.trim() : undefined);
  switch (n.kind) {
    case "endpoint":
      body.to = picked("to");
      body.accept = $("f-accept").value;
      body.formats = [...document.querySelectorAll('input[name="fmt"]:checked')].map((c) => c.value);
      body.apps = [...document.querySelectorAll('input[name="app"]:checked')].map((c) => c.value);
      if (dur("f-port")) body.port = Number(dur("f-port"));
      break;
    case "cups":
      body.to = picked("to");
      body.port = Number($("f-port").value);
      if (dur("f-tls-cert")) body.tls_cert = dur("f-tls-cert");
      if (dur("f-tls-key")) body.tls_key = dur("f-tls-key");
      break;
    case "sink":
      if (dur("f-latest")) body.latest = Number(dur("f-latest"));
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
    $("sidebar").scrollTop = 0;
  } catch (e) {
    status(e.message, true);
  }
}

async function remove(name) {
  if (!confirm(`Delete ${name}? A node another one sends to stays. An endpoint's templates go with it; a display's slips and devices stay until it is added again.`)) return;
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
  links.append(el("span", "IPP: " + n.ipp, { className: "font-mono text-xs text-base-content/60 break-all" }));
  const row = el("div", undefined, { className: "flex flex-wrap gap-3" });
  row.append(el("a", "Jobs sent here", { className: "link", href: "/admin/jobs?dest=" + encodeURIComponent(n.name) }));
  if (n.kind === "inbox") row.append(el("a", `${n.held || 0} held: open the inbox`, { className: "link", href: "/admin/inboxes" }));
  if (n.pending) row.append(el("a", `${n.pending} deliveries waiting`, { className: "link", href: "/admin/deliveries" }));
  links.append(row);
  aside.append(links);
  if (n.kind === "endpoint") {
    const entry = section("Entry");
    entry.append(hint("Applications POST a document, or JSON naming a template and its data, with the key as a bearer token."));
    entry.append(el("pre", `POST ${location.origin}${n.in}`, { className: "overflow-x-auto rounded-box bg-base-200 p-2 font-mono text-xs break-all whitespace-pre-wrap" }));
    const key = section("Key");
    key.append(hint("Shown once, when the endpoint was made or the key rotated; the manager keeps only its hash. Rotating it makes a new one and the old one stops working."));
    key.append(button("Rotate key", async () => showKey(n.name, (await api("POST", `/api/nodes/${n.name}/key`)).key), "btn btn-sm btn-outline w-fit"));
    const tpls = section("Templates");
    aside.append(entry, key, tpls);
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
  box.append(hint(`This endpoint fills them with the data it is sent: POST ${n.in} with JSON naming the template ("template") and its data ("data"). A compiled template (.bdft), or Excel, Word, PowerPoint, a PDF form, pdfme or text.`));
  if (n.accept === "file") box.append(el("p", "This endpoint takes documents only: its templates are not filled until Accepts changes.", { className: "text-xs text-warning" }));
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
$("key-copy").addEventListener("click", async () => {
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


document.addEventListener("keydown", (ev) => {
  const typing = ["INPUT", "TEXTAREA", "SELECT"].includes(document.activeElement?.tagName);
  if (ev.key === "Escape" && draft) {
    draft = null;
    canvas.setDraft(null, null);
    renderSidebar();
  } else if (ev.key === "Delete" && selected && !typing) {
    remove(selected);
  }
});

// a link from another page names a node: /admin/topology#<name>
if (location.hash.length > 1) {
  selected = decodeURIComponent(location.hash.slice(1));
  canvas.setSelected(selected);
}
onChange(load);
load().then(poll);
setInterval(() => {
  if (!document.hidden) poll();
}, POLL_MS);
