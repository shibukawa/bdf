// The deliveries waiting: jobs on their way to printers and archives,
// their tries, and the tick that sends what is due.

import { $, _, api, badge, button, el, onChange, since, status } from "./admin.js";

// STATES are the states of a delivery, as the page says them.
const STATES = { pending: _("pending"), failed: _("failed"), done: _("done"), canceled: _("canceled") };

async function load() {
  try {
    const ds = await api("GET", "/api/deliveries");
    const rows = ds.map((d) => {
      const tr = el("tr");
      tr.append(
        el("td", "#" + d.job, { className: "font-mono" }),
        el("td", d.node, { className: "font-mono text-sm" }),
        el("td", d.printer || "", { className: "font-mono text-sm" }),
        el("td"),
        el("td", String(d.attempts || 0)),
        el("td", d.retry_at ? since(d.retry_at) : _("now"), { className: "whitespace-nowrap text-sm" }),
        el("td", d.last_error || "", { className: "text-xs text-error break-all" }),
        el("td"),
      );
      tr.children[3].append(badge(STATES[d.state] || d.state, d.state === "failed" ? "badge-error" : d.lease ? "badge-info" : "badge-ghost"));
      tr.lastChild.append(button(_("Retry now"), () => api("POST", `/api/deliveries/${d.id}/retry`), "btn btn-xs"));
      return tr;
    });
    if (!rows.length) {
      const tr = el("tr");
      tr.append(el("td", _("Nothing waiting"), { colSpan: 8, className: "text-base-content/50" }));
      rows.push(tr);
    }
    $("deliveries").replaceChildren(...rows);
  } catch (e) {
    status(e.message, true);
  }
}

$("tick").addEventListener("click", async () => {
  try {
    const r = await api("POST", "/api/tick");
    status(_("tick: delivered {delivered}, failed {undelivered}, printed {printed} ({failed} failed), removed {removed}", { delivered: r.delivered, undelivered: r.undelivered, printed: r.printed, failed: r.failed, removed: r.removed + r.items + r.jobs }));
    await load();
  } catch (e) {
    status(e.message, true);
  }
});
$("refresh").addEventListener("click", load);
onChange(load);
load();
setInterval(() => {
  if (!document.hidden) load();
}, 10000);
