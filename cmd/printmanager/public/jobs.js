// A person's page: print a document to a node, and see their jobs (an
// administrator: every job), through the API with the session's cookie.

const $ = (id) => document.getElementById(id);

function status(text, bad) {
  if (!text) return;
  const t = document.createElement("div");
  t.role = "status";
  t.className = "alert shadow-lg " + (bad ? "alert-error" : "alert-success");
  t.textContent = text;
  $("toasts").append(t);
  setTimeout(() => t.remove(), bad ? 6000 : 3000);
}

async function api(method, path, body) {
  const res = await fetch(path, { method, body });
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

const KINDS = { endpoint: "Endpoint", queue: "Queue", class: "Class", inbox: "Inbox", archive: "Archive", display: "Display", printer: "Printer" };
const ENDED = ["completed", "canceled", "aborted"];

function legsOf(j) {
  return (j.legs || [])
    .map((l) => l.node + ":" + l.state + (l.member ? "(" + l.member + ")" : "") + (l.error ? " ! " + l.error : ""))
    .join(" → ");
}

function stateBadge(state) {
  const cls = state === "completed" ? "badge-success badge-soft" : state === "held" ? "badge-warning" : state === "processing" ? "badge-info" : "badge-ghost";
  return el("span", state, { className: "badge " + cls });
}

async function loadDestinations() {
  const ds = await api("GET", "/api/destinations");
  const chosen = $("print-dest").value;
  $("print-dest").replaceChildren(...ds.map((d) => el("option", `${d.title && d.title !== d.name ? d.title + " · " : ""}${d.name} (${KINDS[d.kind] || d.kind})`, { value: d.name, selected: d.name === chosen })));
}

async function loadJobs() {
  const jobs = await api("GET", "/api/my/jobs?limit=100");
  const rows = jobs.map((j) => {
    const tr = el("tr");
    tr.append(
      el("td", "#" + j.number, { className: "font-mono" }),
      el("td", j.name || "", { className: "font-semibold" }),
      el("td", j.dest, { className: "font-mono text-sm" }),
      el("td", j.user || "", { className: "text-sm" }),
      el("td", new Date(j.received).toLocaleString(), { className: "text-sm text-base-content/60" }),
      el("td"),
      el("td", legsOf(j), { className: "font-mono text-xs" }),
      el("td"),
    );
    tr.children[5].append(stateBadge(j.state));
    const actions = tr.lastChild;
    const a = el("a", "Document", { className: "btn btn-xs btn-ghost mr-1", href: `/api/jobs/${j.id}/document` });
    actions.append(a);
    if (!ENDED.includes(j.state) || (j.legs || []).some((l) => l.kind === "display" && l.state === "shown")) {
      const b = el("button", "Cancel", { type: "button", className: "btn btn-xs btn-ghost text-error" });
      b.addEventListener("click", async () => {
        try {
          await api("POST", `/api/jobs/${j.id}/cancel`);
          await loadJobs();
        } catch (e) {
          status(e.message, true);
        }
      });
      actions.append(b);
    }
    return tr;
  });
  if (!rows.length) {
    const tr = el("tr");
    tr.append(el("td", "No jobs yet", { colSpan: 8, className: "text-base-content/50" }));
    rows.push(tr);
  }
  $("jobs").replaceChildren(...rows);
}

async function load() {
  try {
    await loadDestinations();
    await loadJobs();
  } catch (e) {
    status(e.message, true);
  }
}

$("print-form").addEventListener("submit", async (ev) => {
  ev.preventDefault();
  const file = $("print-file").files[0];
  if (!file) return;
  const name = $("print-name").value || file.name;
  try {
    const j = await api("POST", `/api/destinations/${$("print-dest").value}/jobs?name=${encodeURIComponent(name)}`, file);
    status(`Job #${j.number} is ${j.state}`);
    $("print-file").value = "";
    $("print-name").value = "";
    await loadJobs();
  } catch (e) {
    status(e.message, true);
  }
});
$("refresh").addEventListener("click", load);
load();
setInterval(() => {
  if (!document.hidden) loadJobs().catch(() => {});
}, 10000);
