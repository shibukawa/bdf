// The administrator's audit page: who did what, and the recent jobs with
// their routes, through the API with the session's cookie.

const $ = (id) => document.getElementById(id);

async function api(path) {
  const res = await fetch(path);
  if (res.status === 401) {
    location.href = "/";
    throw new Error("signed out");
  }
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

function legsOf(j) {
  return (j.legs || [])
    .map((l) => l.node + ":" + l.state + (l.member ? "(" + l.member + ")" : "") + (l.error ? " ! " + l.error : ""))
    .join(" → ");
}

async function load() {
  const log = await api("/api/audit?limit=500");
  const rows = log.map((e) => {
    const tr = el("tr");
    tr.append(
      el("td", new Date(e.at).toLocaleString(), { className: "text-sm text-base-content/60 whitespace-nowrap" }),
      el("td", e.who, { className: "font-mono text-sm" }),
      el("td", e.action, { className: "font-semibold" }),
      el("td", e.subject || "", { className: "font-mono text-sm" }),
      el("td", e.detail || "", { className: "text-sm break-all" }),
    );
    return tr;
  });
  if (!rows.length) {
    const tr = el("tr");
    tr.append(el("td", "Nothing yet", { colSpan: 5, className: "text-base-content/50" }));
    rows.push(tr);
  }
  $("audit").replaceChildren(...rows);

  const jobs = await api("/api/my/jobs?limit=200");
  const jrows = jobs.map((j) => {
    const tr = el("tr");
    tr.append(
      el("td", "#" + j.number, { className: "font-mono" }),
      el("td", j.name || "", { className: "font-semibold" }),
      el("td", j.dest, { className: "font-mono text-sm" }),
      el("td", j.user || "", { className: "text-sm" }),
      el("td", j.source, { className: "text-sm" }),
      el("td", new Date(j.received).toLocaleString(), { className: "text-sm text-base-content/60 whitespace-nowrap" }),
      el("td", j.state),
      el("td", legsOf(j), { className: "font-mono text-xs" }),
    );
    return tr;
  });
  if (!jrows.length) {
    const tr = el("tr");
    tr.append(el("td", "No jobs yet", { colSpan: 8, className: "text-base-content/50" }));
    jrows.push(tr);
  }
  $("jobs").replaceChildren(...jrows);
}

load().catch((e) => {
  const tr = el("tr");
  tr.append(el("td", e.message, { colSpan: 5, className: "text-error" }));
  $("audit").replaceChildren(tr);
});
