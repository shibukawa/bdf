// The apps: applications that send with credentials of their own
// (client credentials). An administrator registers one with its scopes
// and gets its secret once; an endpoint lists the apps it takes.

import { $, _, api, badge, button, el, onChange, since, status } from "./admin.js";

const SCOPES = { print: _("print: send documents"), "print.template": _("print.template: fill templates"), "jobs.cancel": _("jobs.cancel: cancel its own jobs"), export: _("export: read the exports"), tick: "tick: POST /api/tick", admin: _("admin: the API as an administrator") };

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
      const where = a.public ? (a.redirect_uris || []).join(", ") : lists.join(", ");
      tr.append(el("td", a.id, { className: "font-mono font-semibold" }), el("td", a.title || ""), el("td"), el("td", a.group || "", { className: "text-sm" }), el("td", where, { className: "font-mono text-xs break-all" }), el("td", a.last_used ? since(a.last_used) : _("never"), { className: "whitespace-nowrap text-sm text-base-content/60" }), el("td", undefined, { className: "whitespace-nowrap" }));
      if (a.public) tr.children[2].append(badge(_("public client"), "badge-info badge-soft mr-1"));
      for (const s of a.scopes || []) tr.children[2].append(badge(s, s === "admin" ? "badge-error badge-soft mr-1" : "badge-ghost mr-1"));
      if (a.disabled) tr.children[2].append(badge(_("disabled"), "badge-warning"));
      const actions = tr.lastChild;
      const saved = { title: a.title, scopes: a.scopes, group: a.group, public: a.public, redirect_uris: a.redirect_uris };
      actions.append(button(_("Edit"), async () => edit(a), "btn btn-xs btn-ghost mr-1"));
      if (!a.public) actions.append(button(_("Rotate secret"), async () => showSecret(a.id, (await api("POST", `/api/apps/${a.id}/secret`)).secret), "btn btn-xs btn-ghost mr-1"));
      actions.append(
        button(a.disabled ? _("Enable") : _("Disable"), () => api("PUT", "/api/apps/" + a.id, { ...saved, disabled: !a.disabled }), "btn btn-xs btn-ghost mr-1"),
        button(_("Delete"), () => api("DELETE", "/api/apps/" + a.id), "btn btn-xs btn-ghost text-error"),
      );
      return tr;
    });
    if (!rows.length) {
      const tr = el("tr");
      tr.append(el("td", _("No apps yet"), { colSpan: 7, className: "text-base-content/50" }));
      rows.push(tr);
    }
    $("apps").replaceChildren(...rows);
  } catch (e) {
    status(e.message, true);
  }
}

function scopeBoxes(checked) {
  const box = $("app-scopes");
  box.replaceChildren(el("span", _("Scopes"), { className: "text-sm font-semibold" }));
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
  $("app-public").checked = !!a.public;
  $("app-redirects").value = (a.redirect_uris || []).join("\n");
  publicFields();
  scopeBoxes(a.scopes || []);
  $("app-submit").textContent = _("Save {id}", { id: a.id });
  $("app-cancel").hidden = false;
  $("app-form").scrollIntoView({ behavior: "smooth", block: "nearest" });
}

// publicFields shows the redirect URIs of a public client.
function publicFields() {
  $("app-redirects").hidden = !$("app-public").checked;
}

function reset() {
  editing = "";
  $("app-form").reset();
  $("app-id").readOnly = false;
  publicFields();
  scopeBoxes(["print"]);
  $("app-submit").textContent = _("Add");
  $("app-cancel").hidden = true;
}

$("app-form").addEventListener("submit", async (ev) => {
  ev.preventDefault();
  const id = $("app-id").value.trim();
  const body = { title: $("app-title").value.trim(), scopes: [...document.querySelectorAll('input[name="scope"]:checked')].map((c) => c.value) };
  if ($("app-group").value.trim()) body.group = $("app-group").value.trim();
  if ($("app-public").checked) {
    body.public = true;
    body.redirect_uris = $("app-redirects").value.split(/\s+/).map((s) => s.trim()).filter(Boolean);
  }
  try {
    const a = await api("PUT", "/api/apps/" + encodeURIComponent(id), body);
    reset();
    if (a.secret) showSecret(id, a.secret);
    else if (a.public) status(_("{id} saved: a public client, with no secret. Its frontend sends people to {url} (PKCE) and exchanges the code at /oauth/token", { id, url: location.origin + "/oauth/authorize" }));
    else status(_("{name} saved", { name: id }));
    await load();
  } catch (e) {
    status(e.message, true);
  }
});
$("app-cancel").addEventListener("click", reset);
$("app-public").addEventListener("change", publicFields);
$("refresh").addEventListener("click", load);
onChange(load);
reset();
load();
