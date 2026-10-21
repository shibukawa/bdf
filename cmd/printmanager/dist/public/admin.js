// Shared by the administrator's pages (topology.js, inboxes.js,
// admin-jobs.js, deliveries.js, audit.js): the API with the session's
// cookie (the sign-in at /), toasts, small builders of elements, and the
// nav's own tab.

export const $ = (id) => document.getElementById(id);

// status shows a message for a while, as a toast.
export function status(text, bad) {
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
export async function api(method, path, body, raw) {
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

export function el(tag, text, attrs = {}) {
  const e = document.createElement(tag);
  if (text !== undefined) e.textContent = text;
  Object.assign(e, attrs);
  return e;
}

// changed tells the page that something changed: it loads again.
export function changed() {
  document.dispatchEvent(new CustomEvent("pm:changed"));
}

export function onChange(fn) {
  document.addEventListener("pm:changed", fn);
}

// button runs fn and tells the page something changed; an error is a
// toast.
export function button(text, fn, className = "btn btn-sm") {
  const b = el("button", text, { type: "button", className });
  b.addEventListener("click", async () => {
    try {
      await fn();
      changed();
    } catch (e) {
      status(e.message, true);
    }
  });
  return b;
}

export function since(t) {
  if (!t) return "—";
  return new Date(t).toLocaleString();
}

export function badge(text, cls = "badge-ghost") {
  return el("span", text, { className: "badge " + cls });
}

export const KINDS = { endpoint: "Endpoint", queue: "Queue", class: "Class", inbox: "Inbox", archive: "Archive", display: "Display", printer: "Printer" };
export const KIND_BADGE = {
  endpoint: "badge-accent badge-soft",
  queue: "badge-primary badge-soft",
  class: "badge-primary badge-soft",
  inbox: "badge-warning badge-soft",
  archive: "badge-neutral badge-soft",
  display: "badge-info badge-soft",
  printer: "badge-neutral badge-soft",
};

export function kindBadge(kind) {
  return badge(KINDS[kind] || kind, KIND_BADGE[kind] || "badge-ghost");
}

// sender is who sent a job, as verified, and the name they gave.
export function sender(j) {
  const by = j.submitted_by || {};
  const who = by.name || by.id || by.kind || "";
  return j.claimed_user && j.claimed_user !== who ? `${who} (as ${j.claimed_user})` : who;
}

export function legsOf(j) {
  return (j.legs || [])
    .map((l) => l.node + ":" + l.state + (l.member ? "(" + l.member + ")" : "") + (l.error ? " ! " + l.error : ""))
    .join(" → ");
}

export const ENDED = ["completed", "canceled", "aborted"];

export function stateBadge(state) {
  const cls = state === "completed" ? "badge-success badge-soft" : state === "held" ? "badge-warning" : state === "processing" ? "badge-info" : "badge-ghost";
  return el("span", state, { className: "badge whitespace-nowrap " + cls });
}

// cancellable tells a job the administrator may still cancel: one not
// ended, or one whose slip is still open on a board.
export function cancellable(j) {
  return !ENDED.includes(j.state) || (j.legs || []).some((l) => l.kind === "display" && l.state === "shown");
}

// the nav's own tab
for (const a of document.querySelectorAll("#admin-nav a")) {
  if (a.getAttribute("href") === location.pathname) a.classList.add("tab-active");
}
