// A person's page: print a document to a source they may send to (the
// role topology: to any node), and see their jobs and their groups' (an
// administrator: every job), through the API with the session's cookie.

import { KINDS } from "./admin.js";
import { _, locale } from "./i18n.js";

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
    throw new Error(_("signed out"));
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

const ENDED = ["completed", "canceled", "aborted"];
const STATES = { processing: _("processing"), held: _("held"), completed: _("completed"), canceled: _("canceled"), aborted: _("aborted") };

function legsOf(j) {
  return (j.legs || [])
    .map((l) => l.node + ":" + l.state + (l.member ? "(" + l.member + ")" : "") + (l.error ? " ! " + l.error : ""))
    .join(" → ");
}

// sender is who sent a job, as verified, and the name they gave.
function sender(j) {
  const by = j.submitted_by || {};
  const who = by.name || by.id || by.kind || "";
  return j.claimed_user && j.claimed_user !== who ? _("{who} (as {claimed})", { who, claimed: j.claimed_user }) : who;
}

let me = null;

// mine tells a job the person sent.
function mine(j) {
  const by = j.submitted_by || {};
  return me && by.kind === me.kind && by.id === me.id;
}

function stateBadge(state) {
  const cls = state === "completed" ? "badge-success badge-soft" : state === "held" ? "badge-warning" : state === "processing" ? "badge-info" : "badge-ghost";
  return el("span", STATES[state] || state, { className: "badge whitespace-nowrap " + cls, title: state });
}

// loadDestinations lists the sources, where documents come into the
// topology; with the role topology, every other node after them, apart.
async function loadDestinations() {
  const ds = await api("GET", "/api/destinations");
  const chosen = $("print-dest").value;
  const option = (d) => el("option", (d.title && d.title !== d.name ? d.title + " · " : "") + _("{name} ({kind})", { name: d.name, kind: KINDS[d.kind] || d.kind }), { value: d.name, selected: d.name === chosen });
  const group = (label, list) => {
    const g = el("optgroup", undefined, { label });
    g.append(...list.map(option));
    return g;
  };
  const sources = ds.filter((d) => d.source);
  const others = ds.filter((d) => !d.source);
  if (!ds.length) {
    $("print-dest").replaceChildren(el("option", _("No source you may send to"), { value: "", disabled: true, selected: true }));
  } else if (!others.length) {
    $("print-dest").replaceChildren(...sources.map(option));
  } else {
    $("print-dest").replaceChildren(...(sources.length ? [group(_("Sources"), sources)] : []), group(_("Other nodes, to try them out (the role topology)"), others));
  }
}

async function loadJobs() {
  const jobs = await api("GET", "/api/my/jobs?limit=100");
  const rows = jobs.map((j) => {
    const tr = el("tr");
    tr.append(
      el("td", "#" + j.number, { className: "font-mono" }),
      el("td", j.name || "", { className: "font-semibold", title: (j.format ? _("{kind} · received as {format}", { kind: j.kind || "bdf", format: j.format }) : j.kind || "bdf") + (j.meta ? "\n" + Object.entries(j.meta).map(([k, v]) => k + ": " + (typeof v === "string" ? v : JSON.stringify(v))).join("\n") : "") }),
      el("td", j.dest, { className: "whitespace-nowrap font-mono text-sm" }),
      el("td", sender(j), { className: "text-sm" }),
      el("td", j.group || "—", { className: "hidden whitespace-nowrap text-sm md:table-cell" }),
      el("td", new Date(j.received).toLocaleString(locale), { className: "hidden text-sm text-base-content/60 sm:table-cell" }),
      el("td"),
      el("td", legsOf(j), { className: "hidden font-mono text-xs lg:table-cell" }),
      el("td", undefined, { className: "whitespace-nowrap" }),
    );
    tr.children[6].append(stateBadge(j.state));
    const actions = tr.lastChild;
    const a = el("a", _("Document"), { className: "btn btn-xs btn-ghost mr-1", href: `/api/jobs/${j.id}/document` });
    actions.append(a);
    // a reprint, to the job's node or to another chosen in the form above
    const re = el("button", _("Reprint"), { type: "button", className: "btn btn-xs btn-ghost mr-1", title: _("Send the document again, to {node} (hold Shift: to the node chosen above)", { node: j.dest }) });
    re.addEventListener("click", async (ev) => {
      try {
        const to = ev.shiftKey ? $("print-dest").value : j.dest;
        const again = await api("POST", `/api/jobs/${j.id}/reprint`, JSON.stringify({ to }));
        status(_("Job #{n} is {state} (a reprint of #{of} at {node})", { n: again.number, state: STATES[again.state] || again.state, of: j.number, node: to }));
        await loadJobs();
      } catch (e) {
        status(e.message, true);
      }
    });
    actions.append(re);
    // the sender and an administrator cancel; the group only sees
    if ((mine(j) || me?.admin) && (!ENDED.includes(j.state) || (j.legs || []).some((l) => l.kind === "display" && l.state === "shown"))) {
      const b = el("button", _("Cancel"), { type: "button", className: "btn btn-xs btn-ghost text-error" });
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
    tr.append(el("td", _("No jobs yet"), { colSpan: 9, className: "text-base-content/50" }));
    rows.push(tr);
  }
  $("jobs").replaceChildren(...rows);
}

async function load() {
  try {
    me = me || (await api("GET", "/api/me"));
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
    status(_("Job #{n} is {state}", { n: j.number, state: STATES[j.state] || j.state }));
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
