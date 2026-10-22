// The administrator's jobs: every job, or those sent to one node, with
// the route each took; a job not ended (or whose slip is still open) is
// canceled here.

import { $, api, button, cancellable, el, KINDS, legsOf, onChange, payloadBadge, sender, since, stateBadge, status } from "./admin.js";

let dest = new URLSearchParams(location.search).get("dest") || "";

async function loadNodes() {
  const topo = await api("GET", "/api/topology");
  if (dest && !topo.some((n) => n.name === dest)) dest = "";
  const options = [el("option", "All jobs", { value: "", selected: dest === "" })];
  for (const n of topo) options.push(el("option", `${n.name} (${KINDS[n.kind] || n.kind})`, { value: n.name, selected: n.name === dest }));
  $("jobs-dest").replaceChildren(...options);
}

async function loadJobs() {
  const jobs = dest ? await api("GET", `/api/destinations/${encodeURIComponent(dest)}/jobs?limit=100`) : await api("GET", "/api/my/jobs?limit=100");
  const rows = jobs.map((j) => {
    const tr = el("tr");
    tr.append(
      el("td", "#" + j.number, { className: "font-mono" }),
      el("td", undefined, { className: "font-semibold", title: j.meta ? Object.entries(j.meta).map(([k, v]) => k + ": " + (typeof v === "string" ? v : JSON.stringify(v))).join("\n") : "" }),
      el("td", j.dest, { className: "whitespace-nowrap font-mono text-sm" }),
      el("td", sender(j), { className: "text-sm" }),
      el("td", since(j.received), { className: "whitespace-nowrap text-sm text-base-content/60" }),
      el("td"),
      el("td", legsOf(j), { className: "font-mono text-xs" }),
      el("td", undefined, { className: "whitespace-nowrap" }),
    );
    tr.children[1].append(j.name || "", " ", payloadBadge(j));
    tr.children[5].append(stateBadge(j.state));
    tr.lastChild.append(el("a", "Document", { className: "btn btn-xs btn-ghost mr-1", href: `/api/jobs/${j.id}/document` }));
    if (cancellable(j)) tr.lastChild.append(button("Cancel", () => api("POST", `/api/jobs/${j.id}/cancel`), "btn btn-xs btn-ghost text-error"));
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
    await loadNodes();
    await loadJobs();
  } catch (e) {
    status(e.message, true);
  }
}

$("jobs-dest").addEventListener("change", (ev) => {
  dest = ev.target.value;
  history.replaceState(null, "", dest ? "?dest=" + encodeURIComponent(dest) : location.pathname);
  loadJobs().catch((e) => status(e.message, true));
});
$("refresh").addEventListener("click", load);
onChange(load);
load();
setInterval(() => {
  if (!document.hidden) loadJobs().catch(() => {});
}, 10000);
