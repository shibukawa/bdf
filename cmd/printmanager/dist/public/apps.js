// The apps: applications that send with credentials of their own
// (client credentials). An administrator registers one with its scopes
// and gets its secret once; an endpoint lists the apps it takes.

import { $, api, badge, button, el, onChange, since, status } from "./admin.js";

const SCOPES = { print: "print: send documents", "print.template": "print.template: fill templates", "jobs.cancel": "jobs.cancel: cancel its own jobs", tick: "tick: POST /api/tick", admin: "admin: the API as an administrator" };

let editing = ""; // the app in the form

function showSecret(id, secret) {
  $("secret-id").textContent = id;
  $("secret-value").textContent = secret;
  $("secret-usage").textContent = `curl -u '${id}:${secret}' -d grant_type=client_credentials ${location.origin}/oauth/token`;
  $("secret-modal").showModal();
}

async function load() {
  try {
    const [apps, topo] = await Promise.all([api("GET", "/api/apps"), api("GET", "/api/topology")]);
    const rows = apps.map((a) => {
      const tr = el("tr", undefined, { className: a.disabled ? "opacity-60" : "" });
      const lists = topo.filter((n) => (n.apps || []).includes(a.id)).map((n) => n.name);
      tr.append(el("td", a.id, { className: "font-mono font-semibold" }), el("td", a.title || ""), el("td"), el("td", a.group || "", { className: "text-sm" }), el("td", lists.join(", "), { className: "font-mono text-xs" }), el("td", a.last_used ? since(a.last_used) : "never", { className: "whitespace-nowrap text-sm text-base-content/60" }), el("td", undefined, { className: "whitespace-nowrap" }));
      for (const s of a.scopes || []) tr.children[2].append(badge(s, s === "admin" ? "badge-error badge-soft mr-1" : "badge-ghost mr-1"));
      if (a.disabled) tr.children[2].append(badge("disabled", "badge-warning"));
      const actions = tr.lastChild;
      actions.append(
        button("Edit", async () => edit(a), "btn btn-xs btn-ghost mr-1"),
        button("Rotate secret", async () => showSecret(a.id, (await api("POST", `/api/apps/${a.id}/secret`)).secret), "btn btn-xs btn-ghost mr-1"),
        button(a.disabled ? "Enable" : "Disable", () => api("PUT", "/api/apps/" + a.id, { title: a.title, scopes: a.scopes, group: a.group, disabled: !a.disabled }), "btn btn-xs btn-ghost mr-1"),
        button("Delete", () => api("DELETE", "/api/apps/" + a.id), "btn btn-xs btn-ghost text-error"),
      );
      return tr;
    });
    if (!rows.length) {
      const tr = el("tr");
      tr.append(el("td", "No apps yet", { colSpan: 7, className: "text-base-content/50" }));
      rows.push(tr);
    }
    $("apps").replaceChildren(...rows);
  } catch (e) {
    status(e.message, true);
  }
}

function scopeBoxes(checked) {
  const box = $("app-scopes");
  box.replaceChildren(el("span", "Scopes", { className: "text-sm font-semibold" }));
  for (const [s, label] of Object.entries(SCOPES)) {
    const l = el("label", undefined, { className: "inline-flex items-center gap-1 text-sm" });
    l.append(el("input", undefined, { type: "checkbox", className: "checkbox checkbox-sm", name: "scope", value: s, checked: checked.includes(s) }), label);
    box.append(l);
  }
}

function edit(a) {
  editing = a.id;
  $("app-id").value = a.id;
  $("app-id").readOnly = true;
  $("app-title").value = a.title || "";
  $("app-group").value = a.group || "";
  scopeBoxes(a.scopes || []);
  $("app-submit").textContent = "Save " + a.id;
  $("app-cancel").hidden = false;
  $("app-form").scrollIntoView({ behavior: "smooth", block: "nearest" });
}

function reset() {
  editing = "";
  $("app-form").reset();
  $("app-id").readOnly = false;
  scopeBoxes(["print"]);
  $("app-submit").textContent = "Add";
  $("app-cancel").hidden = true;
}

$("app-form").addEventListener("submit", async (ev) => {
  ev.preventDefault();
  const id = $("app-id").value.trim();
  const body = { title: $("app-title").value.trim(), scopes: [...document.querySelectorAll('input[name="scope"]:checked')].map((c) => c.value) };
  if ($("app-group").value.trim()) body.group = $("app-group").value.trim();
  try {
    const a = await api("PUT", "/api/apps/" + encodeURIComponent(id), body);
    reset();
    if (a.secret) showSecret(id, a.secret);
    else status(`${id} saved`);
    await load();
  } catch (e) {
    status(e.message, true);
  }
});
$("app-cancel").addEventListener("click", reset);
$("refresh").addEventListener("click", load);
onChange(load);
reset();
load();
